/**
 * Canonical naming helpers for JetStream streams, NATS subjects, and
 * consumer durables. All four read process.env.GO_ENV; if undefined or empty,
 * "dev" is used as the default.
 *
 * The SERVICE slot is always a real deployed-service identifier (ASSETS,
 * EVENTS, WORKFLOW, ROUTER, TRIGGERS, MAPEXVAULT, MAPEXIAM, JSEXECUTOR,
 * JSWORKFLOWEXECUTOR, HTTPGATEWAY, plus MAPEXOSGOKIT for the cross-service
 * DLQ). Stream names use UPPERCASE env prefix; subjects and durables use
 * lowercase env prefix to match NATS conventions.
 */

/**
 * Returns the runtime environment prefix, reading process.env.GO_ENV.
 * If undefined or empty, returns "dev".
 *
 * Example: getEnv() returns "dev" by default, "prod" when GO_ENV=prod.
 */
export function getEnv(): string {
	const v = process.env.GO_ENV;
	if (!v) return 'dev';
	return v;
}

/**
 * Builds a canonical JetStream stream name following the pattern
 * ${ENV}-MAPEXOS-{SERVICE}-{CONTEXT}. Env, service, and context are uppercased
 * independently of the input casing. When context is empty, the trailing dash
 * is omitted (returns ${ENV}-MAPEXOS-{SERVICE}).
 *
 * Example: streamName("JSWORKFLOWEXECUTOR", "CODE") returns
 * "DEV-MAPEXOS-JSWORKFLOWEXECUTOR-CODE" when GO_ENV=dev.
 */
export function streamName(service: string, context: string): string {
	const env = getEnv().toUpperCase();
	const svc = service.toUpperCase();
	if (!context) return `${env}-MAPEXOS-${svc}`;
	return `${env}-MAPEXOS-${svc}-${context.toUpperCase()}`;
}

/**
 * Builds a canonical NATS subject following the pattern
 * ${env}.mapexos.{service}.{action}. Env, service, and action are lowercased
 * independently of the input casing.
 *
 * Example: subject("workflow", "code") returns "dev.mapexos.workflow.code"
 * when GO_ENV=dev.
 */
export function subject(service: string, action: string): string {
	const env = getEnv().toLowerCase();
	return `${env}.mapexos.${service.toLowerCase()}.${action.toLowerCase()}`;
}

/**
 * Builds a canonical JetStream consumer durable name following the pattern
 * ${env}-{service}-{context}-consumer. Env, service, and context are
 * lowercased independently of the input casing.
 *
 * Example: durable("jsworkflowexecutor", "workflow-code") returns
 * "dev-jsworkflowexecutor-workflow-code-consumer" when GO_ENV=dev.
 */
export function durable(service: string, context: string): string {
	const env = getEnv().toLowerCase();
	return `${env}-${service.toLowerCase()}-${context.toLowerCase()}-consumer`;
}

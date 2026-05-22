export const NATS_CONNECTION_TOKEN = Symbol('nats_connection');
export const NATS_BUS_TOKEN = Symbol('nats_bus');

/**
 * Resolve the runtime environment prefix for canonical stream/subject names.
 * Reads process.env.GO_ENV; if undefined or empty, returns "dev".
 */
function getEnv(): string {
  const v = process.env.GO_ENV;
  if (!v) return 'dev';
  return v;
}

/** DLQ stream name — resolves to e.g. "DEV-MAPEXOS-MAPEXOSGOKIT-DLQ". */
export const DLQ_STREAM = `${getEnv().toUpperCase()}-MAPEXOS-MAPEXOSGOKIT-DLQ`;

/** DLQ subject for publishing dead letter messages — resolves to e.g. "dev.mapexos.mapexosgokit.dlq". */
export const DLQ_SUBJECT = `${getEnv().toLowerCase()}.mapexos.mapexosgokit.dlq`;

/** Default retry policy with exponential backoff (1s, 5s, 30s, 2m, 10m) */
export const DEFAULT_RETRY_BACKOFF = [1000, 5000, 30000, 120000, 600000];

/** Default max retries before sending to DLQ */
export const DEFAULT_MAX_RETRIES = 5;

/** DLQ retention in nanoseconds (30 days) */
export const DLQ_MAX_AGE_NANOS = 30 * 24 * 60 * 60 * 1000000000;

export type ConfigDefinition = {
	key: string;
	env: string;
	type: 'string' | 'int' | 'bool' | 'array' | 'json';
	default: any;
	/**
	 * Marks keys that carry credentials, secrets, or any value that must
	 * never be the hardcoded `default` in a non-dev environment. When
	 * `sensitive` is true and the resolved value still equals `default`
	 * while `NODE_ENV` is not dev/development/empty, ConfigModule refuses
	 * to start the process. See findSensitiveDefaultsInUse in ./validation.
	 */
	sensitive?: boolean;
};
export interface RedisConfig {
	host: string;
	port: number;
	username?: string;
	password?: string;
	db?: number;
	lock?: RedlockConfig;
}

export interface RedlockConfig {
	driftFactor?: number;
	retryCount?: number;
	retryDelay?: number;
	retryJitter?: number;
	automaticExtensionThreshold?: number;
}

export interface GetData {
	cacheKey: string;
	cacheTTL: number;
	callback: () => Promise<any>;
}

export interface GetDataWithLock extends GetData {
	lockTTL?: number;
	retryPolicy?: {
		maxRetry?: number;
		retryDelay?: number;
	};
}
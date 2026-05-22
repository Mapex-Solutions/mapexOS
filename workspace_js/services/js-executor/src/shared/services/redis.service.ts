import type { Logger } from '@mapexos/microservices';
import { RedisService as BaseRedisService } from '@mapexos/infrastructure';

/** Redis wrapper service with logging for all cache operations. */
export class RedisService {
	private readonly service: BaseRedisService;

	constructor(
		private readonly logger: Logger,
		service: BaseRedisService,
	) {
		this.service = service;
	}

	get redlock() {
		return this.service.redlock;
	}

	get instance() {
		return this.service.instance;
	}

	async getData(params: { cacheKey: string; cacheTTL: number; callback: () => Promise<any> }): Promise<any> {
		this.logger.debug(`[SERVICE:Redis] Getting data from cache: ${params.cacheKey}`);
		return this.service.getData(params);
	}

	async getDataWithLock(params: {
		cacheKey: string;
		cacheTTL: number;
		callback: () => Promise<any>;
		lockTTL?: number;
		retryPolicy?: { maxRetry?: number; retryDelay?: number };
	}): Promise<any> {
		this.logger.debug(`[SERVICE:Redis] Getting data with lock from cache: ${params.cacheKey}`);
		return this.service.getDataWithLock(params);
	}

	async setEx(cacheKey: string, expireIn: number, data: any): Promise<void> {
		this.logger.debug(`[SERVICE:Redis] Setting data in cache: ${cacheKey}`);
		return this.service.setEx(cacheKey, expireIn, data);
	}

	async addToList(listName: string, element: string, ttl = 86400): Promise<boolean> {
		this.logger.debug(`[SERVICE:Redis] Adding element to list: ${listName}`);
		return this.service.addToList(listName, element, ttl);
	}

	async elementExistsInList(listName: string, element: string): Promise<boolean> {
		this.logger.debug(`[SERVICE:Redis] Checking if element exists in list: ${listName}`);
		return this.service.elementExistsInList(listName, element);
	}

	async del(key: string): Promise<number> {
		this.logger.debug(`[SERVICE:Redis] Deleting key: ${key}`);
		return this.service.del(key);
	}

	async delMany(keys: string[]): Promise<number> {
		this.logger.debug(`[SERVICE:Redis] Deleting ${keys.length} keys`);
		return this.service.delMany(keys);
	}

	async getBuffer(key: string): Promise<Buffer | null> {
		this.logger.debug(`[SERVICE:Redis] Getting buffer: ${key}`);
		return this.service.getBuffer(key);
	}

	async setBuffer(key: string, ttl: number, data: Buffer): Promise<void> {
		this.logger.debug(`[SERVICE:Redis] Setting buffer: ${key}`);
		return this.service.setBuffer(key, ttl, data);
	}

	async exists(key: string): Promise<boolean> {
		this.logger.debug(`[SERVICE:Redis] Checking if key exists: ${key}`);
		return this.service.exists(key);
	}

	async get(key: string): Promise<string | null> {
		this.logger.debug(`[SERVICE:Redis] Getting string: ${key}`);
		return this.service.get(key);
	}

	async set(key: string, ttl: number, value: string): Promise<void> {
		this.logger.debug(`[SERVICE:Redis] Setting string: ${key}`);
		return this.service.set(key, ttl, value);
	}
}
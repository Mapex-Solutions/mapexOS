import Redis from 'ioredis';
import Redlock, { Lock } from 'redlock';
import { RedisConfig, GetData, GetDataWithLock } from './types';
import { getInfraLogger } from '../logger';

function delay(ms: number): Promise<void> {
	return new Promise(resolve => setTimeout(resolve, ms));
}

export class RedisService {
	private redlockInstance?: Redlock;

	constructor(
		private readonly redis: Redis,
		config?: RedisConfig
	) {
		if (config?.lock) {
			this.redlockInstance = new Redlock([redis as any], config.lock);
		}
	}

	get redlock(): Redlock | undefined {
		return this.redlockInstance;
	}

	get instance(): Redis {
		return this.redis;
	}

	private async getAndParseData(cacheKey: string): Promise<any> {
		const data = await this.redis.get(cacheKey);

		if (data) {
			try {
				return JSON.parse(data);
			} catch {
				return undefined;
			}
		}
		return undefined;
	}

	async getData(params: GetData): Promise<any> {
		const { cacheKey, cacheTTL, callback } = params;

		const data = await this.getAndParseData(cacheKey);
		if (data !== undefined) return data;

		try {
			const dbData = await callback();
			await this.redis.set(cacheKey, JSON.stringify(dbData), 'EX', cacheTTL);
			return dbData;
		} catch (error) {
			throw error;
		}
	}

	async getDataWithLock(params: GetDataWithLock): Promise<any> {
		const { cacheKey, cacheTTL, callback, lockTTL = 3000, retryPolicy = {} } = params;
		const { maxRetry = 30, retryDelay = 150 } = retryPolicy;

		if (!this.redlockInstance) {
			throw new Error('Redlock is not configured. Please provide lock configuration when creating RedisService.');
		}

		let lock: Lock | null = null;

		const data = await this.getAndParseData(cacheKey);
		if (data !== undefined) return data;

		try {
			lock = await this.redlockInstance.acquire([`__LOCK:${cacheKey}`], lockTTL);
		} catch (err) {
			for (let i = 0; i < maxRetry; i++) {
				await delay(retryDelay);

				const data = await this.getAndParseData(cacheKey);
				if (data !== undefined) return data;
			}

			getInfraLogger().error('[INFRA:REDIS] Could not acquire lock, cache empty after retries');
			throw new Error('Could not acquire lock and cache is still empty after retries.');
		}

		try {
			const data = await this.getAndParseData(cacheKey);
			if (data !== undefined) return data;

			const dbData = await callback();
			await this.redis.set(cacheKey, JSON.stringify(dbData), 'EX', cacheTTL);
			return dbData;
		} catch (error) {
			throw error;
		} finally {
			if (lock) {
				try {
					await lock.release();
				} catch (unlockError) {
					getInfraLogger().error({ err: unlockError }, '[INFRA:REDIS] Failed to release lock');
				}
			}
		}
	}

	async setEx(cacheKey: string, expireIn: number, data: any): Promise<void> {
		await this.redis.setex(
			cacheKey,
			expireIn,
			JSON.stringify(data),
		);
	}

	async addToList(listName: string, element: string, ttl = 86400): Promise<boolean> {
		const exists = await this.redis.exists(listName);
		if (!exists) {
			await this.redis.rpush(listName, element);
			await this.redis.expire(listName, ttl);
			return true;
		}

		await this.redis.rpush(listName, element);
		return true;
	}

	async elementExistsInList(listName: string, element: string): Promise<boolean> {
		const exists = await this.redis.exists(listName);
		if (!exists) return false;

		const list = await this.redis.lrange(listName, 0, -1);
		return list.includes(element);
	}

	/**
	 * Checks the Redis connection health by issuing a PING command.
	 *
	 * @returns A promise that resolves when the connection is healthy
	 * @throws If the connection is down or unreachable
	 */
	async ping(): Promise<void> {
		await this.redis.ping();
	}

	async disconnect(): Promise<void> {
		await this.redis.disconnect();
	}

	/**
	 * Deletes a key from Redis.
	 *
	 * @param key - The key to delete
	 * @returns Number of keys deleted (0 or 1)
	 */
	async del(key: string): Promise<number> {
		return this.redis.del(key);
	}

	/**
	 * Deletes multiple keys from Redis.
	 *
	 * @param keys - Array of keys to delete
	 * @returns Number of keys deleted
	 */
	async delMany(keys: string[]): Promise<number> {
		if (keys.length === 0) return 0;
		return this.redis.del(...keys);
	}

	/**
	 * Gets a value as Buffer from Redis.
	 * Useful for binary data like compiled bytecode.
	 *
	 * @param key - The key to get
	 * @returns Buffer or null if key doesn't exist
	 */
	async getBuffer(key: string): Promise<Buffer | null> {
		return this.redis.getBuffer(key);
	}

	/**
	 * Sets a Buffer value with TTL.
	 * Useful for binary data like compiled bytecode.
	 *
	 * @param key - The key to set
	 * @param ttl - Time to live in seconds
	 * @param data - Buffer data to store
	 */
	async setBuffer(key: string, ttl: number, data: Buffer): Promise<void> {
		await this.redis.setex(key, ttl, data);
	}

	/**
	 * Checks if a key exists in Redis.
	 *
	 * @param key - The key to check
	 * @returns true if key exists, false otherwise
	 */
	async exists(key: string): Promise<boolean> {
		const result = await this.redis.exists(key);
		return result === 1;
	}

	/**
	 * Gets a string value from Redis.
	 *
	 * @param key - The key to get
	 * @returns String value or null if key doesn't exist
	 */
	async get(key: string): Promise<string | null> {
		return this.redis.get(key);
	}

	/**
	 * Sets a string value with TTL.
	 *
	 * @param key - The key to set
	 * @param ttl - Time to live in seconds
	 * @param value - String value to store
	 */
	async set(key: string, ttl: number, value: string): Promise<void> {
		await this.redis.setex(key, ttl, value);
	}
}
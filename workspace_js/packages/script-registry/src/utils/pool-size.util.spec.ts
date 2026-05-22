import { getPoolSize, getCpuCount } from './pool-size.util';
import { cpus } from 'os';

jest.mock('os', () => ({
	cpus: jest.fn(),
}));

const mockCpus = cpus as jest.MockedFunction<typeof cpus>;

describe('pool-size.util', () => {
	describe('getPoolSize', () => {
		it('should return explicit value when configValue > 0', () => {
			expect(getPoolSize(4)).toBe(4);
			expect(getPoolSize(10)).toBe(10);
			expect(getPoolSize(1)).toBe(1);
		});

		it('should auto-detect based on CPU cores when configValue is 0', () => {
			mockCpus.mockReturnValue(Array(4).fill({}) as any);
			expect(getPoolSize(0)).toBe(4);
		});

		it('should enforce minimum of 2 when CPU count is 1', () => {
			mockCpus.mockReturnValue(Array(1).fill({}) as any);
			expect(getPoolSize(0)).toBe(2);
		});

		it('should enforce maximum of 8 when CPU count exceeds 8', () => {
			mockCpus.mockReturnValue(Array(16).fill({}) as any);
			expect(getPoolSize(0)).toBe(8);

			mockCpus.mockReturnValue(Array(32).fill({}) as any);
			expect(getPoolSize(0)).toBe(8);
		});

		it('should return CPU count when within bounds (2-8)', () => {
			mockCpus.mockReturnValue(Array(2).fill({}) as any);
			expect(getPoolSize(0)).toBe(2);

			mockCpus.mockReturnValue(Array(6).fill({}) as any);
			expect(getPoolSize(0)).toBe(6);

			mockCpus.mockReturnValue(Array(8).fill({}) as any);
			expect(getPoolSize(0)).toBe(8);
		});
	});

	describe('getCpuCount', () => {
		it('should return the number of CPU cores', () => {
			mockCpus.mockReturnValue(Array(8).fill({}) as any);
			expect(getCpuCount()).toBe(8);

			mockCpus.mockReturnValue(Array(4).fill({}) as any);
			expect(getCpuCount()).toBe(4);
		});
	});
});

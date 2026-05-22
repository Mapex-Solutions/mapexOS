/**
 * Script Compression Utility Tests
 *
 * Tests bytecode compression/decompression: threshold behavior,
 * round-trip integrity, edge cases, invalid input handling.
 */

import { compressBytecode, decompressBytecode } from './script-compression.util';

describe('compressBytecode', () => {
	it('should compress bytecode above threshold (512 bytes)', () => {
		const bytecode = new ArrayBuffer(1024);
		new Uint8Array(bytecode).fill(0x42);

		const result = compressBytecode(bytecode);

		expect(result.compressed).toBe(true);
		expect(result.originalSize).toBe(1024);
		expect(result.finalSize).toBeLessThan(1024); // gzip should compress repeated bytes well
		expect(result.payload[0]).toBe(0x01); // COMPRESSED_FLAG
	});

	it('should not compress bytecode below threshold (< 512 bytes)', () => {
		const bytecode = new ArrayBuffer(100);
		new Uint8Array(bytecode).fill(0x42);

		const result = compressBytecode(bytecode);

		expect(result.compressed).toBe(false);
		expect(result.originalSize).toBe(100);
		expect(result.finalSize).toBe(101); // 1 byte flag + 100 bytes data
		expect(result.payload[0]).toBe(0x00); // UNCOMPRESSED_FLAG
	});

	it('should handle empty bytecode', () => {
		const bytecode = new ArrayBuffer(0);

		const result = compressBytecode(bytecode);

		expect(result.compressed).toBe(false);
		expect(result.originalSize).toBe(0);
		expect(result.payload.length).toBe(1); // Just the flag byte
	});

	it('should handle bytecode exactly at threshold', () => {
		const bytecode = new ArrayBuffer(512);
		new Uint8Array(bytecode).fill(0xAA);

		const result = compressBytecode(bytecode);

		expect(result.compressed).toBe(true);
		expect(result.originalSize).toBe(512);
	});
});

describe('decompressBytecode', () => {
	it('should decompress compressed bytecode', () => {
		const original = new ArrayBuffer(1024);
		new Uint8Array(original).fill(0x42);

		const { payload } = compressBytecode(original);
		const result = decompressBytecode(payload);

		expect(result.wasCompressed).toBe(true);
		expect(result.bytecode.byteLength).toBe(1024);
		expect(new Uint8Array(result.bytecode).every(b => b === 0x42)).toBe(true);
	});

	it('should return uncompressed bytecode as-is', () => {
		const original = new ArrayBuffer(100);
		new Uint8Array(original).fill(0x42);

		const { payload } = compressBytecode(original);
		const result = decompressBytecode(payload);

		expect(result.wasCompressed).toBe(false);
		expect(result.bytecode.byteLength).toBe(100);
		expect(new Uint8Array(result.bytecode).every(b => b === 0x42)).toBe(true);
	});

	it('should throw on empty buffer', () => {
		const emptyBuffer = Buffer.alloc(0);

		expect(() => decompressBytecode(emptyBuffer)).toThrow('Empty buffer');
	});

	it('should throw on invalid compression flag', () => {
		const invalidBuffer = Buffer.from([0xFF, 0x01, 0x02, 0x03]);

		expect(() => decompressBytecode(invalidBuffer)).toThrow('Invalid compression flag');
	});
});

describe('round-trip integrity', () => {
	it('should preserve large bytecode through compress/decompress', () => {
		const bytecode = new ArrayBuffer(4096);
		const view = new Uint8Array(bytecode);
		for (let i = 0; i < view.length; i++) {
			view[i] = i % 256;
		}

		const { payload } = compressBytecode(bytecode);
		const { bytecode: result } = decompressBytecode(payload);

		expect(result.byteLength).toBe(4096);
		const resultView = new Uint8Array(result);
		for (let i = 0; i < resultView.length; i++) {
			expect(resultView[i]).toBe(i % 256);
		}
	});

	it('should preserve small bytecode through compress/decompress', () => {
		const bytecode = new ArrayBuffer(10);
		new Uint8Array(bytecode).set([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);

		const { payload } = compressBytecode(bytecode);
		const { bytecode: result } = decompressBytecode(payload);

		expect(new Uint8Array(result)).toEqual(new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]));
	});
});

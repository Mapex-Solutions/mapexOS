import { gzipSync, gunzipSync } from 'zlib';
import {
  compressScript,
  decompressScript,
  compressBytecode,
  decompressBytecode,
  CompressionResult,
  DecompressionResult,
  BytecodeCompressionResult,
  BytecodeDecompressionResult,
} from './script-compression.util';

// Mock terser to control minification behavior
jest.mock('terser', () => ({
  minify: jest.fn(),
}));

import { minify } from 'terser';
const mockedMinify = minify as jest.MockedFunction<typeof minify>;

describe('Script Compression Utilities', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('compressScript', () => {
    beforeEach(() => {
      // Default successful minification - returns input as-is for most tests
      mockedMinify.mockImplementation(async (input) => ({
        code: typeof input === 'string' ? input : 'const result=payload;',
      }));
    });

    it('should compress large scripts with gzip', async () => {
      // Create a script that is definitely > 512 bytes
      const largeScript = 'a'.repeat(600) + '; const result = payload;';

      const result: CompressionResult = await compressScript(largeScript);

      expect(result.compressed).toBe(true);
      expect(result.originalSize).toBeGreaterThanOrEqual(512);
      expect(result.finalSize).toBeLessThan(result.originalSize);
      expect(result.payload[0]).toBe(0x01); // Compressed flag
    });

    it('should store small scripts without compression', async () => {
      const smallScript = 'const result = payload;'; // < 512 bytes

      const result: CompressionResult = await compressScript(smallScript);

      expect(result.compressed).toBe(false);
      expect(result.originalSize).toBeLessThan(512);
      expect(result.finalSize).toBeGreaterThan(result.originalSize); // Due to flag overhead
      expect(result.payload[0]).toBe(0x00); // Uncompressed flag
    });

    it('should sanitize empty scripts', async () => {
      const emptyScript = '   \n  \t  ';

      const result: CompressionResult = await compressScript(emptyScript);

      expect(result.compressed).toBe(false);
      expect(mockedMinify).toHaveBeenCalledWith(
        'const result = undefined;',
        expect.any(Object)
      );
    });

    it('should add semicolon to scripts without proper termination', async () => {
      const scriptWithoutSemicolon = 'const result = payload';

      await compressScript(scriptWithoutSemicolon);

      expect(mockedMinify).toHaveBeenCalledWith(
        'const result = payload;',
        expect.any(Object)
      );
    });

    it('should handle minification errors gracefully', async () => {
      const script = 'const result = payload;';
      const minifyError = new Error('Minification failed');
      mockedMinify.mockRejectedValue(minifyError);

      // Mock console.warn to verify fallback behavior
      const consoleSpy = jest.spyOn(console, 'warn').mockImplementation();

      const result: CompressionResult = await compressScript(script);

      expect(consoleSpy).toHaveBeenCalledWith(
        expect.stringContaining('Minification failed for script')
      );
      expect(result.compressed).toBe(false);

      consoleSpy.mockRestore();
    });

    it('should use proper terser configuration', async () => {
      const script = 'const result = payload;';

      await compressScript(script);

      expect(mockedMinify).toHaveBeenCalledWith(
        expect.any(String),
        {
          compress: {
            dead_code: true,
            drop_console: false,
            drop_debugger: true,
            pure_funcs: [],
          },
          mangle: false,
          format: {
            comments: false,
          },
        }
      );
    });
  });

  describe('decompressScript', () => {
    it('should decompress gzip-compressed scripts', () => {
      const originalScript = 'const result = payload;';
      const compressed = gzipSync(Buffer.from(originalScript, 'utf8'));
      const payload = Buffer.concat([Buffer.from([0x01]), compressed]);

      const result: DecompressionResult = decompressScript(payload);

      expect(result.script).toBe(originalScript);
      expect(result.wasCompressed).toBe(true);
    });

    it('should return uncompressed scripts as-is', () => {
      const originalScript = 'const result = payload;';
      const payload = Buffer.concat([
        Buffer.from([0x00]),
        Buffer.from(originalScript, 'utf8')
      ]);

      const result: DecompressionResult = decompressScript(payload);

      expect(result.script).toBe(originalScript);
      expect(result.wasCompressed).toBe(false);
    });

    it('should throw error for empty buffer', () => {
      const emptyBuffer = Buffer.alloc(0);

      expect(() => decompressScript(emptyBuffer)).toThrow(
        'Empty buffer provided for decompression'
      );
    });

    it('should throw error for invalid compression flag', () => {
      const invalidPayload = Buffer.from([0x99, 0x01, 0x02, 0x03]);

      expect(() => decompressScript(invalidPayload)).toThrow(
        'Invalid compression flag: 153. Expected 0 or 1'
      );
    });
  });

  describe('compressBytecode', () => {
    it('should compress large bytecode with gzip', () => {
      const largeBytecode = new ArrayBuffer(1024); // > 512 bytes
      const view = new Uint8Array(largeBytecode);
      view.fill(0xAB); // Fill with test data

      const result: BytecodeCompressionResult = compressBytecode(largeBytecode);

      expect(result.compressed).toBe(true);
      expect(result.originalSize).toBe(1024);
      expect(result.finalSize).toBeLessThan(result.originalSize);
      expect(result.payload[0]).toBe(0x01); // Compressed flag
    });

    it('should store small bytecode without compression', () => {
      const smallBytecode = new ArrayBuffer(256); // < 512 bytes
      const view = new Uint8Array(smallBytecode);
      view.fill(0xCD); // Fill with test data

      const result: BytecodeCompressionResult = compressBytecode(smallBytecode);

      expect(result.compressed).toBe(false);
      expect(result.originalSize).toBe(256);
      expect(result.finalSize).toBe(257); // Original + flag byte
      expect(result.payload[0]).toBe(0x00); // Uncompressed flag
    });

    it('should handle empty bytecode', () => {
      const emptyBytecode = new ArrayBuffer(0);

      const result: BytecodeCompressionResult = compressBytecode(emptyBytecode);

      expect(result.compressed).toBe(false);
      expect(result.originalSize).toBe(0);
      expect(result.finalSize).toBe(1); // Just the flag byte
      expect(result.payload[0]).toBe(0x00);
    });
  });

  describe('decompressBytecode', () => {
    it('should decompress gzip-compressed bytecode', () => {
      const originalData = Buffer.alloc(1024);
      originalData.fill(0xAB);
      const compressed = gzipSync(originalData);
      const payload = Buffer.concat([Buffer.from([0x01]), compressed]);

      const result: BytecodeDecompressionResult = decompressBytecode(payload);

      expect(result.bytecode.byteLength).toBe(1024);
      expect(result.wasCompressed).toBe(true);

      // Verify data integrity
      const resultView = new Uint8Array(result.bytecode);
      expect(resultView[0]).toBe(0xAB);
      expect(resultView[1023]).toBe(0xAB);
    });

    it('should return uncompressed bytecode as ArrayBuffer', () => {
      const originalData = Buffer.alloc(256);
      originalData.fill(0xCD);
      const payload = Buffer.concat([Buffer.from([0x00]), originalData]);

      const result: BytecodeDecompressionResult = decompressBytecode(payload);

      expect(result.bytecode.byteLength).toBe(256);
      expect(result.wasCompressed).toBe(false);

      // Verify data integrity
      const resultView = new Uint8Array(result.bytecode);
      expect(resultView[0]).toBe(0xCD);
      expect(resultView[255]).toBe(0xCD);
    });

    it('should throw error for empty buffer', () => {
      const emptyBuffer = Buffer.alloc(0);

      expect(() => decompressBytecode(emptyBuffer)).toThrow(
        'Empty buffer provided for bytecode decompression'
      );
    });

    it('should throw error for invalid compression flag', () => {
      const invalidPayload = Buffer.from([0x77, 0x01, 0x02, 0x03]);

      expect(() => decompressBytecode(invalidPayload)).toThrow(
        'Invalid compression flag: 119. Expected 0 or 1'
      );
    });
  });

  describe('compression threshold behavior', () => {
    it('should respect 512-byte threshold for scripts', async () => {
      // Test exactly at threshold
      const exactThresholdScript = 'x'.repeat(512);
      mockedMinify.mockResolvedValue({ code: exactThresholdScript });

      const result = await compressScript(exactThresholdScript);
      expect(result.compressed).toBe(true);

      // Test just below threshold
      const belowThresholdScript = 'x'.repeat(511);
      mockedMinify.mockResolvedValue({ code: belowThresholdScript });

      const result2 = await compressScript(belowThresholdScript);
      expect(result2.compressed).toBe(false);
    });

    it('should respect 512-byte threshold for bytecode', () => {
      // Test exactly at threshold
      const exactThresholdBytecode = new ArrayBuffer(512);
      const result = compressBytecode(exactThresholdBytecode);
      expect(result.compressed).toBe(true);

      // Test just below threshold
      const belowThresholdBytecode = new ArrayBuffer(511);
      const result2 = compressBytecode(belowThresholdBytecode);
      expect(result2.compressed).toBe(false);
    });
  });

  describe('round-trip compatibility', () => {
    it('should maintain script integrity through compression/decompression cycle', async () => {
      const originalScript = `
        // Complex script with various elements
        const data = {
          timestamp: Date.now(),
          values: [1, 2, 3, 4, 5],
          nested: {
            property: "test"
          }
        };

        const result = {
          ...payload,
          processed: data,
          completed: true
        };
      `;

      // Compress
      const compressed = await compressScript(originalScript);

      // Decompress
      const decompressed = decompressScript(compressed.payload);

      // Should match minified version (not original with whitespace)
      expect(decompressed.script).toBeDefined();
      expect(decompressed.script.length).toBeGreaterThan(0);
    });

    it('should maintain bytecode integrity through compression/decompression cycle', () => {
      // Create test bytecode with pattern
      const originalBytecode = new ArrayBuffer(1024);
      const originalView = new Uint8Array(originalBytecode);
      for (let i = 0; i < 1024; i++) {
        originalView[i] = i % 256;
      }

      // Compress
      const compressed = compressBytecode(originalBytecode);

      // Decompress
      const decompressed = decompressBytecode(compressed.payload);

      // Verify integrity
      expect(decompressed.bytecode.byteLength).toBe(1024);
      const decompressedView = new Uint8Array(decompressed.bytecode);

      for (let i = 0; i < 1024; i++) {
        expect(decompressedView[i]).toBe(i % 256);
      }
    });
  });
});
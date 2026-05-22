import { gzipSync, gunzipSync } from 'zlib';
import { minify } from 'terser';

const THRESHOLD_BYTES = 512;
const COMPRESSED_FLAG = 0x01;
const UNCOMPRESSED_FLAG = 0x00;

export interface CompressionResult {
  payload: Buffer;
  compressed: boolean;
  originalSize: number;
  finalSize: number;
}

export interface DecompressionResult {
  script: string;
  wasCompressed: boolean;
}

export interface BytecodeCompressionResult {
  payload: Buffer;
  compressed: boolean;
  originalSize: number;
  finalSize: number;
}

export interface BytecodeDecompressionResult {
  bytecode: ArrayBuffer;
  wasCompressed: boolean;
}

/**
 * Sanitizes script code to ensure it's valid JavaScript before minification
 *
 * @param script - The script code to sanitize
 * @returns Sanitized script code
 */
function sanitizeScript(script: string): string {
  // Remove leading/trailing whitespace
  let sanitized = script.trim();

  // If the script is empty or just whitespace, return a minimal valid script
  if (!sanitized) {
    return 'const result = undefined;';
  }

  // Ensure the script ends with a semicolon for better minification
  if (!sanitized.endsWith(';') && !sanitized.endsWith('}')) {
    sanitized += ';';
  }

  return sanitized;
}

/**
 * Compresses script using minification and conditional gzip compression
 *
 * @param script - The script code to compress
 * @returns CompressionResult with payload and metadata
 */
export async function compressScript(script: string): Promise<CompressionResult> {
  // Step 1: Sanitize the script
  const sanitized = sanitizeScript(script);

  // Step 2: Minify the script with better error handling
  let minified: string;
  try {
    const minifyResult = await minify(sanitized, {
      compress: {
        dead_code: true,
        drop_console: false,
        drop_debugger: true,
        pure_funcs: [],
      },
      mangle: false, // Keep variable names for debugging
      format: {
        comments: false,
      },
    });

    minified = minifyResult.code || sanitized;
  } catch (error) {
    // If minification fails, use the sanitized version
    console.warn(`Minification failed for script, using original: ${error.message}`);
    minified = sanitized;
  }

  const raw = Buffer.from(minified, 'utf8');
  const originalSize = raw.length;

  // Step 3: Check threshold for compression
  if (raw.length >= THRESHOLD_BYTES) {
    // Compress with gzip
    const gz = gzipSync(raw);
    const payload = Buffer.concat([Buffer.from([COMPRESSED_FLAG]), gz]);

    return {
      payload,
      compressed: true,
      originalSize,
      finalSize: payload.length,
    };
  } else {
    // Store without compression
    const payload = Buffer.concat([Buffer.from([UNCOMPRESSED_FLAG]), raw]);

    return {
      payload,
      compressed: false,
      originalSize,
      finalSize: payload.length,
    };
  }
}

/**
 * Decompresses script based on flag
 *
 * @param buf - The buffer containing the compressed/uncompressed script
 * @returns DecompressionResult with script and metadata
 */
export function decompressScript(buf: Buffer): DecompressionResult {
  if (buf.length === 0) {
    throw new Error('Empty buffer provided for decompression');
  }

  const flag = buf[0];
  const data = buf.slice(1);

  if (flag === UNCOMPRESSED_FLAG) {
    // Uncompressed data
    return {
      script: data.toString('utf8'),
      wasCompressed: false,
    };
  } else if (flag === COMPRESSED_FLAG) {
    // Compressed data - decompress with gunzip
    const decompressed = gunzipSync(data);
    return {
      script: decompressed.toString('utf8'),
      wasCompressed: true,
    };
  } else {
    throw new Error(`Invalid compression flag: ${flag}. Expected ${UNCOMPRESSED_FLAG} or ${COMPRESSED_FLAG}`);
  }
}

/**
 * Compresses bytecode using conditional gzip compression (no minification needed)
 *
 * @param bytecode - The ArrayBuffer bytecode to compress
 * @returns BytecodeCompressionResult with payload and metadata
 */
export function compressBytecode(bytecode: ArrayBuffer): BytecodeCompressionResult {
  const raw = Buffer.from(bytecode);
  const originalSize = raw.length;

  // Check threshold for compression
  if (raw.length >= THRESHOLD_BYTES) {
    // Compress with gzip
    const gz = gzipSync(raw);
    const payload = Buffer.concat([Buffer.from([COMPRESSED_FLAG]), gz]);

    return {
      payload,
      compressed: true,
      originalSize,
      finalSize: payload.length,
    };
  } else {
    // Store without compression
    const payload = Buffer.concat([Buffer.from([UNCOMPRESSED_FLAG]), raw]);

    return {
      payload,
      compressed: false,
      originalSize,
      finalSize: payload.length,
    };
  }
}

/**
 * Decompresses bytecode based on flag
 *
 * @param buf - The buffer containing the compressed/uncompressed bytecode
 * @returns BytecodeDecompressionResult with bytecode and metadata
 */
export function decompressBytecode(buf: Buffer): BytecodeDecompressionResult {
  if (buf.length === 0) {
    throw new Error('Empty buffer provided for bytecode decompression');
  }

  const flag = buf[0];
  const data = buf.slice(1);

  if (flag === UNCOMPRESSED_FLAG) {
    // Uncompressed data
    return {
      bytecode: data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength),
      wasCompressed: false,
    };
  } else if (flag === COMPRESSED_FLAG) {
    // Compressed data - decompress with gunzip
    const decompressed = gunzipSync(data);
    return {
      bytecode: decompressed.buffer.slice(decompressed.byteOffset, decompressed.byteOffset + decompressed.byteLength),
      wasCompressed: true,
    };
  } else {
    throw new Error(`Invalid compression flag: ${flag}. Expected ${UNCOMPRESSED_FLAG} or ${COMPRESSED_FLAG}`);
  }
}

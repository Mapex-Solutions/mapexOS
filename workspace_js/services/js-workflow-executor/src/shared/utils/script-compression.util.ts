import { gzipSync, gunzipSync } from 'zlib';

const THRESHOLD_BYTES = 512;
const COMPRESSED_FLAG = 0x01;
const UNCOMPRESSED_FLAG = 0x00;

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
 * Compresses bytecode using conditional gzip compression
 *
 * @param bytecode - The ArrayBuffer bytecode to compress
 * @returns BytecodeCompressionResult with payload and metadata
 */
export function compressBytecode(bytecode: ArrayBuffer): BytecodeCompressionResult {
  const raw = Buffer.from(bytecode);
  const originalSize = raw.length;

  if (raw.length >= THRESHOLD_BYTES) {
    const gz = gzipSync(raw);
    const payload = Buffer.concat([Buffer.from([COMPRESSED_FLAG]), gz]);

    return {
      payload,
      compressed: true,
      originalSize,
      finalSize: payload.length,
    };
  } else {
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
    return {
      bytecode: data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength),
      wasCompressed: false,
    };
  } else if (flag === COMPRESSED_FLAG) {
    const decompressed = gunzipSync(data);
    return {
      bytecode: decompressed.buffer.slice(decompressed.byteOffset, decompressed.byteOffset + decompressed.byteLength),
      wasCompressed: true,
    };
  } else {
    throw new Error(`Invalid compression flag: ${flag}. Expected ${UNCOMPRESSED_FLAG} or ${COMPRESSED_FLAG}`);
  }
}

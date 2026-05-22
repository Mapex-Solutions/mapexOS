# Script Compression and Bytecode Utilities

## Overview
This document describes how JS‑Executor compresses **script source code** and **V8 bytecode** for caching and transport efficiency. Compression is conditional and optimized for Redis/MinIO storage.

## Compression Strategy
- **Text scripts**: sanitize + minify (Terser) + conditional gzip
- **Bytecode**: conditional gzip (no minification)
- **Threshold**: 512 bytes → gzip if >= 512 bytes, store raw if < 512 bytes
- **Flags**:
  - `0x00` = uncompressed
  - `0x01` = gzip compressed

## File Location
```
src/modules/scripts/application/utils/
├── script-compression.util.ts    # Main compression utilities
├── mapexValidatorCode.util.ts     # Validator code injection
└── index.ts                       # Utils exports
```

## Data Shapes
### CompressionResult
```ts
interface CompressionResult {
  payload: Buffer;        // Data with compression flag
  compressed: boolean;    // Whether compression was applied
  originalSize: number;   // Size before compression
  finalSize: number;      // Size after compression
}
```

### DecompressionResult
```ts
interface DecompressionResult {
  script: string;
  wasCompressed: boolean;
}
```

### BytecodeCompressionResult
```ts
interface BytecodeCompressionResult {
  payload: Buffer;
  compressed: boolean;
  originalSize: number;
  finalSize: number;
}
```

### BytecodeDecompressionResult
```ts
interface BytecodeDecompressionResult {
  bytecode: ArrayBuffer;
  wasCompressed: boolean;
}
```

## Core Functions
### sanitizeScript(script: string): string
- Ensures valid JS syntax before minification.
- Returns a safe fallback if input is empty.

### compressScript(script: string): Promise<CompressionResult>
- Sanitizes and minifies with Terser.
- Applies gzip if size >= 512 bytes.
- Prepends compression flag.

### decompressScript(buf: Buffer): DecompressionResult
- Reads flag + decompresses if needed.
- Throws on invalid flags.

### compressBytecode(bytecode: ArrayBuffer): BytecodeCompressionResult
- Converts to Buffer + gzip if >= 512 bytes.
- Prepends compression flag.

### decompressBytecode(buf: Buffer): BytecodeDecompressionResult
- Reads flag + decompresses if needed.
- Returns ArrayBuffer.

## Why this matters
- Reduces cache size and speeds up retrieval.
- Ensures scripts/bytecode can be stored and transported efficiently.
- Preserves fast paths for small scripts (no gzip overhead).

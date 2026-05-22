/**
 * Public cross-module type aliases for the engine bounded context.
 * Other modules MUST import these through `@modules/engine/application/ports`
 * — NEVER directly from `domain/types` (enforces the DDD entity boundary
 * per /js-arch-back §8).
 */
export type { ScriptSet } from '@modules/engine/domain/types';

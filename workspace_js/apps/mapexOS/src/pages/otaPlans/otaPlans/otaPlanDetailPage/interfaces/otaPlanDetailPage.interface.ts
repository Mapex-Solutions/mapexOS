import type { OTAExecutionState } from '@mapexos/schemas';

/** Option for the executions state filter select */
export interface OtaExecutionStateOption {
	label: string;
	value: OTAExecutionState | null;
}

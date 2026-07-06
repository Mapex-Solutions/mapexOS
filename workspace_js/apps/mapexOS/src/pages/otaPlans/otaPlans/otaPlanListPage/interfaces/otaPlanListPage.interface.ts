import type { OTAPlanStatus } from '@mapexos/schemas';

/**
 * Active filters applied to the plans list request
 */
export interface OtaPlanListFilters {
	name?: string | undefined;
	status?: OTAPlanStatus | undefined;
}

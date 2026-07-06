import type { OtaPlanFormData } from '../interfaces/createEditOtaPlan.interface';

/**
 * 1-based wizard leaves, grouped as:
 *   Setup   { Details, Schedule }
 *   Rollout { FromTemplate, Devices, ToTemplate, Firmware }
 *   Finalization { Review }
 *
 * The FROM template + devices are chosen before the TO template + firmware:
 * you pick who moves, then where they move to and what they run.
 */
export const STEP = {
	DETAILS: 1,
	SCHEDULE: 2,
	FROM_TEMPLATE: 3,
	DEVICES: 4,
	TO_TEMPLATE: 5,
	FIRMWARE: 6,
	REVIEW: 7,
} as const;

export const TOTAL_STEPS = 7;

/** Platform rollout defaults (mirror the backend config defaults) */
export const ROLLOUT_DEFAULTS = {
	ratePerMinute: 60,
	abortThresholdPct: 20,
	abortMinExecuted: 50,
} as const;

/** Page size when listing the source template's candidate devices */
export const ASSETS_FETCH_PER_PAGE = 100;

/**
 * Fresh form state. A factory (not a shared const) because the nested
 * rolloutConfig would otherwise be shared across visits by a shallow spread.
 */
export function createInitialOtaPlanFormData(): OtaPlanFormData {
	return {
		name: '',
		description: '',
		startAt: null,
		maxTime: null,
		rolloutConfig: { ...ROLLOUT_DEFAULTS },
		sourceTemplateId: null,
		selectedAssetIds: [],
		targetTemplateId: null,
		version: '',
		file: null,
		firmwareId: null,
		firmwareReady: false,
		checksumHex: '',
		size: 0,
	};
}

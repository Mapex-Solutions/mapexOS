import type { OTARolloutConfig } from '@mapexos/schemas';

/**
 * Wizard form state, filled across the grouped steps:
 *   Setup (details + schedule)
 *   Rollout (source template -> devices -> target template -> firmware)
 *   Finalization (review)
 */
export interface OtaPlanFormData {
	/** Setup — details */
	name: string;
	description: string;

	/** Setup — schedule. startAt null means "run immediately"; ISO otherwise. */
	startAt: string | null;
	maxTime: string | null;
	rolloutConfig: OTARolloutConfig;

	/** Rollout — source template + the devices on it to update */
	sourceTemplateId: string | null;
	selectedAssetIds: string[];

	/** Rollout — target template + the firmware artifact (must be READY) */
	targetTemplateId: string | null;
	version: string;
	file: File | null;
	firmwareId: string | null;
	firmwareReady: boolean;
	checksumHex: string;
	size: number;
}

/** Option shape for the template selects */
export interface TemplateOption {
	label: string;
	value: string;
}

/** Row shape for the device selection list */
export interface SelectableAsset {
	id: string;
	name: string;
	assetUUID: string;
}

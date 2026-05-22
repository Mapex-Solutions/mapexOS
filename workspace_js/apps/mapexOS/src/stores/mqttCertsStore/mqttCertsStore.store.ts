import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { RevokedCertRow } from './types';

/**
 * Pinia store for the mqttcerts revoked list per asset.
 * Cache + refresh on invalidate.
 *
 * NOTE: the actual API call wiring lives in Step4Connectivity (or a
 * composable that holds the api instance) — this store just owns the
 * cache and exposes setters/getters.
 */
export const useMqttCertsStore = defineStore('mqttCerts', () => {
	const revokedByAsset = ref<Map<string, RevokedCertRow[]>>(new Map());
	const loadingByAsset = ref<Map<string, boolean>>(new Map());

	function revokedFor(assetUUID: string): RevokedCertRow[] {
		return revokedByAsset.value.get(assetUUID) ?? [];
	}
	function isLoading(assetUUID: string): boolean {
		return loadingByAsset.value.get(assetUUID) ?? false;
	}
	function setRows(assetUUID: string, rows: RevokedCertRow[]): void {
		revokedByAsset.value.set(assetUUID, rows);
	}
	function setLoading(assetUUID: string, v: boolean): void {
		loadingByAsset.value.set(assetUUID, v);
	}
	function clear(assetUUID: string): void {
		revokedByAsset.value.delete(assetUUID);
		loadingByAsset.value.delete(assetUUID);
	}

	return { revokedByAsset, loadingByAsset, revokedFor, isLoading, setRows, setLoading, clear };
});

import { ref } from 'vue';
import apis from '@services/mapex';
import type { RetentionPolicyResponse } from '@mapexos/schemas';

/**
 * Data composable for the admin retention policies page.
 *
 * Wraps apis.events.retention.listRetentionPolicies + upsertRetentionPolicy.
 * Keeps the page dumb: it binds to `items` + `loading` + `error` and calls
 * `fetch()` / `update()` as needed.
 */
export function useRetentionPolicies() {
	const items = ref<RetentionPolicyResponse[]>([]);
	const loading = ref(false);
	const error = ref<string | null>(null);

	async function fetch() {
		loading.value = true;
		error.value = null;
		try {
			const result = await apis.events.retention.listRetentionPolicies({ perPage: 100 });
			items.value = result.items ?? [];
		} catch (e) {
			error.value = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			loading.value = false;
		}
	}

	async function update(type: string, name: string, retentionDays: number) {
		await apis.events.retention.upsertRetentionPolicy({
			type,
			name,
			retentionDays,
			enabled: true,
		});
		await fetch();
	}

	return { items, loading, error, fetch, update };
}

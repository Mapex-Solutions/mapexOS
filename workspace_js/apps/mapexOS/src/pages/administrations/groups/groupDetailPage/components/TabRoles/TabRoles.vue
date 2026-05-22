<script setup lang="ts">
defineOptions({
	name: 'TabRoles'
});

/** TYPE IMPORTS */
import type { DataRowColumn } from '@components/cards';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { DataRow } from '@components/cards';

/** COMPOSABLES */
import { useGroupDetailTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import type { RoleInfo } from './interfaces/TabRoles.interface';

interface TabRolesProps {
	/** Role IDs from group */
	roleIds: string[];
	/** Loading state from parent */
	loading?: boolean;
}

/** PROPS */
const props = withDefaults(defineProps<TabRolesProps>(), {
	roleIds: () => [],
	loading: false,
});

/** COMPOSABLES */
const t = useGroupDetailTranslations();
const router = useRouter();
const logger = useLogger('TabRoles');

/** STATE */
const roles = ref<RoleInfo[]>([]);
const loadingRoles = ref(false);

/** COMPUTED */

/**
 * Columns configuration for DataRow
 */
const columns = computed<DataRowColumn[]>(() => [
	{
		key: 'avatar',
		label: '',
		type: 'avatar',
		visible: 'always',
		width: 56,
		icon: () => 'admin_panel_settings',
		color: () => 'primary',
	},
	{
		key: 'name',
		label: t.roles.columns.name.value,
		type: 'text',
		visible: 'always',
		width: 200,
		ellipsis: true,
		secondaryKey: 'description',
	},
	{
		key: 'isSystem',
		label: t.roles.columns.type.value,
		type: 'chip',
		visible: 'laptop',
		width: 120,
		color: (value: boolean) => value ? 'purple' : 'blue',
		icon: (value: boolean) => value ? 'lock' : 'lock_open',
		format: (value: boolean) => value ? 'SYSTEM' : 'CUSTOM',
	},
	{
		key: 'scope',
		label: t.roles.columns.scope.value,
		type: 'chip',
		visible: 'laptop',
		width: 120,
		color: (value: string) => value === 'global' ? 'teal' : 'orange',
		icon: (value: string) => value === 'global' ? 'public' : 'place',
		format: (value: string) => value?.toUpperCase() || '—',
	},
]);

/** FUNCTIONS */

/**
 * Fetch role details for each roleId
 *
 * @returns {Promise<void>}
 */
async function fetchRoles(): Promise<void> {
	if (!props.roleIds.length || !apis.mapexOS?.roles) {
		roles.value = [];
		return;
	}

	loadingRoles.value = true;

	try {
		// Fetch all roles and filter by roleIds
		// This is more efficient than N individual requests
		const response = await apis.mapexOS.roles.list({ perPage: 100 });
		const roleIdsSet = new Set(props.roleIds);

		roles.value = (response.items || [])
			.filter((role: any) => roleIdsSet.has(role.id))
			.map((role: any) => ({
				id: role.id || '',
				name: role.name || '—',
				description: role.description || '',
				isSystem: role.isSystem ?? false,
				isTemplate: role.isTemplate ?? false,
				scope: role.scope || '',
			}));

		logger.debug('Fetched roles for group:', { count: roles.value.length });
	} catch (err: any) {
		logger.error('Error fetching roles:', err);
		notifyFail({ message: t.roles.loadError.value });
	} finally {
		loadingRoles.value = false;
	}
}

/**
 * Navigate to role detail page
 *
 * @param {RoleInfo} role - Role clicked
 */
function viewRole(role: RoleInfo): void {
	if (!role.id) return;
	void router.push(`/roles/detail/${role.id}`);
}

/** WATCHERS */
watch(() => props.roleIds, () => {
	void fetchRoles();
}, { deep: true });

/** LIFECYCLE */
onMounted(() => {
	void fetchRoles();
});
</script>

<template>
	<div class="tab-roles">
		<!-- Loading -->
		<div v-if="loadingRoles || props.loading" class="flex flex-center q-pa-lg">
			<q-spinner color="primary" size="40px" />
		</div>

		<!-- Roles List -->
		<div v-else-if="roles.length > 0" class="section">
			<!-- Section Header -->
			<div class="section-header q-mb-md">
				<q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
				<span class="text-subtitle1 text-weight-medium">{{ t.roles.title.value }}</span>
				<q-badge color="primary" class="q-ml-sm">{{ roles.length }}</q-badge>
			</div>

			<!-- Roles Rows -->
			<div
				v-for="(role, index) in roles"
				:key="role.id || `role-${index}`"
				class="q-mb-xs"
			>
				<DataRow
					:data="role"
					:columns="columns"
					:show-actions="false"
					@click="viewRole"
				/>
			</div>
		</div>

		<!-- Empty State -->
		<div v-else class="section text-center q-pa-xl">
			<q-icon name="admin_panel_settings" size="64px" color="grey-4" class="q-mb-md" />
			<div class="text-h6 text-grey-8 q-mb-sm">{{ t.roles.empty.title.value }}</div>
			<div class="text-body2 text-grey-6">{{ t.roles.empty.description.value }}</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.tab-roles {
	min-height: 200px;

	.section {
		background: var(--mapex-surface-bg);
		border-radius: var(--mapex-radius-md);
		padding: 16px;
	}

	.section-header {
		display: flex;
		align-items: center;
		border-bottom: 1px solid var(--mapex-card-border);
		padding-bottom: 8px;
	}
}
</style>

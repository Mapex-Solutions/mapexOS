<script setup lang="ts">
defineOptions({
	name: 'TabMembers'
});

/** TYPE IMPORTS */
import type { TabMembersProps, GroupMemberInfo } from '../../interfaces';
import type { DataRowColumn } from '@components/cards';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';

/** COMPOSABLES */
import { useGroupDetailTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** PROPS */
const props = defineProps<TabMembersProps>();

/** COMPOSABLES */
const t = useGroupDetailTranslations();
const router = useRouter();
const logger = useLogger('TabMembers');

/** STATE */
const members = ref<GroupMemberInfo[]>([]);
const loadingMembers = ref(false);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const itemsPerPage = ref(20);

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
		icon: () => 'person',
		color: () => 'primary',
	},
	{
		key: 'userFullName',
		label: t.members.columns.name.value,
		type: 'text',
		visible: 'always',
		width: 200,
		ellipsis: true,
		secondaryKey: 'userEmail',
	},
	{
		key: 'addedAt',
		label: t.members.columns.addedAt.value,
		type: 'chip',
		visible: 'laptop',
		width: 180,
		color: 'blue-6',
		icon: 'schedule',
		format: (value: any) => {
			if (!value) return '—';
			try {
				return new Date(value).toLocaleDateString('pt-BR', {
					day: '2-digit',
					month: '2-digit',
					year: 'numeric',
				});
			} catch {
				return '—';
			}
		},
	},
]);

/** FUNCTIONS */

/**
 * Fetch group members from API
 *
 * @returns {Promise<void>}
 */
async function fetchMembers(): Promise<void> {
	if (!props.groupId || !apis.mapexOS?.groups) {
		return;
	}

	loadingMembers.value = true;

	try {
		// API expects pathParams and queryParams as separate objects
		const response = await apis.mapexOS.groups.getMembers(
			{ groupId: props.groupId },
			{ page: currentPage.value, perPage: itemsPerPage.value }
		);

		// Map response to GroupMemberInfo with full name
		members.value = (response.items || []).map((member: any) => ({
			id: member.id || '',
			userId: member.userId || '',
			userEmail: member.userEmail || '',
			userFirstName: member.userFirstName || '',
			userLastName: member.userLastName || '',
			userFullName: `${member.userFirstName || ''} ${member.userLastName || ''}`.trim() || member.userEmail || '—',
			addedAt: member.addedAt || '',
			addedBy: member.addedBy || '',
		}));

		// Update pagination
		if (response.pagination) {
			totalPages.value = response.pagination.totalPages || 1;
			totalItems.value = response.pagination.totalItems || 0;
		}
	} catch (err: any) {
		logger.error('Error fetching members:', err);
		notifyFail({ message: t.members.loadError.value });
	} finally {
		loadingMembers.value = false;
	}
}

/**
 * Handle page change
 *
 * @param {number} page - New page number
 */
function handlePageChange(page: number): void {
	currentPage.value = page;
	void fetchMembers();
}

/**
 * Navigate to user detail page
 *
 * @param {GroupMemberInfo} member - Member clicked
 */
function viewMember(member: GroupMemberInfo): void {
	if (!member.userId) return;
	void router.push(`/users/detail/${member.userId}`);
}

/** WATCHERS */
watch(() => props.groupId, () => {
	currentPage.value = 1;
	void fetchMembers();
});

/** LIFECYCLE */
onMounted(() => {
	void fetchMembers();
});
</script>

<template>
	<div class="tab-members">
		<!-- Loading -->
		<div v-if="loadingMembers" class="flex flex-center q-pa-lg">
			<q-spinner color="primary" size="40px" />
		</div>

		<!-- Members List -->
		<div v-else-if="members.length > 0" class="section">
			<!-- Section Header -->
			<div class="section-header q-mb-md">
				<q-icon name="people" color="primary" size="sm" class="q-mr-sm" />
				<span class="text-subtitle1 text-weight-medium">{{ t.members.title.value }}</span>
				<q-badge color="primary" class="q-ml-sm">{{ totalItems }}</q-badge>
			</div>

			<!-- Members Rows -->
			<div
				v-for="(member, index) in members"
				:key="member.id || `member-${index}`"
				class="q-mb-xs"
			>
				<DataRow
					:data="member"
					:columns="columns"
					:show-actions="false"
					@click="viewMember"
				/>
			</div>

			<!-- Pagination -->
			<ListPagination
				v-if="totalPages > 1"
				v-model="currentPage"
				:total-pages="totalPages"
				@change="handlePageChange"
			/>
		</div>

		<!-- Empty State -->
		<div v-else class="section text-center q-pa-xl">
			<q-icon name="people" size="64px" color="grey-4" class="q-mb-md" />
			<div class="text-h6 text-grey-8 q-mb-sm">{{ t.members.empty.title.value }}</div>
			<div class="text-body2 text-grey-6">{{ t.members.empty.description.value }}</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.tab-members {
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

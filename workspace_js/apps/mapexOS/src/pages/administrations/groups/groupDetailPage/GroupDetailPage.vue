<script setup lang="ts">
defineOptions({
	name: 'GroupDetailPage'
});

/** TYPE IMPORTS */
import type { GroupDetailData } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { AppTabs } from '@components/tabs';
import { TabInfo, TabRoles, TabMembers } from './components';

/** COMPOSABLES */
import { useGroupDetailTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import { TAB, DEFAULT_TAB, getTabsConfig } from './constants';

/** COMPOSABLES & STORES */
const t = useGroupDetailTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('GroupDetailPage');
const orgStore = useOrganizationStore();

/** GROUP ID FROM ROUTE */
const groupId = computed(() => route.params.id as string);

/** STATE */
const activeTab = ref(DEFAULT_TAB);
const group = ref<GroupDetailData | null>(null);
const loading = ref(false);

/** COMPUTED */

/**
 * Tabs configuration with reactive translations and badges
 */
const tabs = computed(() => {
	const config = getTabsConfig(t);

	// Add badge to roles tab dynamically
	const rolesTab = config.find(tab => tab.name === TAB.ROLES);
	if (rolesTab && group.value?.roleIds?.length) {
		rolesTab.badge = group.value.roleIds.length;
	}

	// Add badge to members tab dynamically
	const membersTab = config.find(tab => tab.name === TAB.MEMBERS);
	if (membersTab && group.value?.membersCount) {
		membersTab.badge = group.value.membersCount;
	}

	return config;
});

/**
 * Page title based on group name
 */
const pageTitle = computed(() => {
	if (!group.value) return t.page.title.value;
	return group.value.name || t.page.title.value;
});

/** FUNCTIONS */

/**
 * Fetch group details from API
 *
 * @returns {Promise<void>}
 */
async function fetchGroupDetails(): Promise<void> {
	if (!groupId.value) {
		notifyFail({ message: t.errors.idMissing.value });
		void router.back();
		return;
	}

	if (!apis.mapexOS?.groups) {
		notifyFail({ message: t.errors.apiNotInitialized.value });
		return;
	}

	loading.value = true;

	try {
		const response = await apis.mapexOS.groups.getById({ groupId: groupId.value });

		// Get organization name from store
		const currentOrg = orgStore.flatList.find((org: any) => org.id === response.orgId);

		// Build group detail data with explicit optional property handling
		const groupData: GroupDetailData = {
			id: response.id || '',
			name: response.name || '',
			enabled: response.enabled ?? false,
		};

		// Add optional fields only if they have values
		if (response.description) groupData.description = response.description;
		if (response.orgId) groupData.orgId = response.orgId;
		if (currentOrg?.name) groupData.organizationName = currentOrg.name;
		if (response.pathKey) groupData.pathKey = response.pathKey;
		if (response.membersCount !== undefined) groupData.membersCount = response.membersCount;
		if ((response as any).roleIds?.length) groupData.roleIds = (response as any).roleIds;
		if (response.created) groupData.created = response.created;
		if (response.updated) groupData.updated = response.updated;

		group.value = groupData;
		logger.debug('Group details loaded:', { groupId: groupData.id, name: groupData.name });
	} catch (err: any) {
		logger.error('Error fetching group:', err);
		notifyFail({ message: t.page.loadError.value });
		void router.back();
	} finally {
		loading.value = false;
	}
}

/**
 * Navigate to edit group page
 *
 * @returns {void}
 */
function editGroup(): void {
	if (!groupId.value) return;
	void router.push(`/groups/edit/${groupId.value}`);
}

/**
 * Navigate back (uses browser history)
 *
 * @returns {void}
 */
function goBack(): void {
	void router.back();
}

/** LIFECYCLE HOOKS */
onMounted(() => {
	void fetchGroupDetails();
});
</script>

<template>
	<q-page class="q-pa-lg">
		<!-- Loading State -->
		<div v-if="loading" class="flex flex-center q-pa-xl">
			<q-spinner color="primary" size="50px" />
			<span class="q-ml-md text-grey-7">{{ t.page.loading.value }}</span>
		</div>

		<!-- Content -->
		<div v-else>
			<!-- Header Section -->
			<PageHeader
				icon="group"
				icon-color="primary"
				:title="pageTitle"
				:description="t.page.description.value"
				:button="{
					label: t.page.backButton.value,
					icon: 'arrow_back',
					flat: true,
					onClick: goBack,
				}"
			>
				<template #actions>
					<q-btn
						unelevated
						rounded
						color="primary"
						icon="edit"
						:label="t.page.editButton.value"
						@click="editGroup"
					/>
				</template>
			</PageHeader>

			<!-- Tabs Navigation -->
			<q-card flat bordered class="rounded-borders">
				<AppTabs v-model="activeTab" :tabs="tabs" :separator="false">
					<!-- Tab Panels -->
					<q-tab-panels v-model="activeTab" animated class="bg-transparent">
						<!-- Info Tab -->
						<q-tab-panel :name="TAB.INFO" class="q-pa-lg">
							<TabInfo :group="group" :loading="loading" />
						</q-tab-panel>

						<!-- Roles Tab -->
						<q-tab-panel :name="TAB.ROLES" class="q-pa-lg">
							<TabRoles :role-ids="group?.roleIds || []" :loading="loading" />
						</q-tab-panel>

						<!-- Members Tab -->
						<q-tab-panel :name="TAB.MEMBERS" class="q-pa-lg">
							<TabMembers :group-id="groupId" :loading="loading" />
						</q-tab-panel>
					</q-tab-panels>
				</AppTabs>
			</q-card>
		</div>
	</q-page>
</template>

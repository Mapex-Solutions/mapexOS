<script setup lang="ts">
defineOptions({
  name: 'UserDetailPage'
});

/** TYPE IMPORTS */
import type { UserDetailData } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { AppTabs } from '@components/tabs';
import {
  TabProfile,
  TabAccess,
  TabGroups,
} from './components';

/** COMPOSABLES */
import { useUserDetailTranslations } from '@composables/i18n';
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
const t = useUserDetailTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('UserDetailPage');
const orgStore = useOrganizationStore();

/** USER ID FROM ROUTE */
const userId = computed(() => route.params.id as string);

/** STATE */
const activeTab = ref(DEFAULT_TAB);
const user = ref<UserDetailData | null>(null);
const loading = ref(false);

/** COMPUTED */

/**
 * Tabs configuration with reactive translations and badges
 * Groups tab only shows if user belongs to groups
 */
const tabs = computed(() => {
  const config = getTabsConfig(t);

  // Filter out groups tab if user has no groups
  const hasGroups = user.value?.groupsCount && user.value.groupsCount > 0;
  const filteredConfig = hasGroups
    ? config
    : config.filter(tab => tab.name !== TAB.GROUPS);

  // Add badge to groups tab dynamically
  if (hasGroups && user.value?.groupsCount) {
    const groupsTab = filteredConfig.find(tab => tab.name === TAB.GROUPS);
    if (groupsTab) {
      groupsTab.badge = user.value.groupsCount;
    }
  }

  return filteredConfig;
});

/**
 * Page title based on user name
 */
const pageTitle = computed(() => {
  if (!user.value) return t.page.title.value;
  const firstName = user.value.firstName || '';
  const lastName = user.value.lastName || '';
  const fullName = `${firstName} ${lastName}`.trim();
  return fullName || user.value.email || t.page.title.value;
});

/** FUNCTIONS */

/**
 * Fetch complete user details from API (single request for all tabs)
 *
 * @returns {Promise<void>}
 */
async function fetchUserDetails(): Promise<void> {
  if (!userId.value) {
    notifyFail({ message: t.errors.idMissing.value });
    void router.push('/users');
    return;
  }

  if (!apis.mapexOS?.users) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;

  try {
    const response = await apis.mapexOS.users.getById({ userId: userId.value });

    // Get organization name from store
    const currentOrg = orgStore.flatList.find((org: any) => org.id === orgStore.selectedOrganizationId);

    // Build user detail data with proper type handling
    const userData: UserDetailData = {
      id: response.id || '',
      email: response.email || '',
      enabled: response.enabled ?? false,
      organizationName: currentOrg?.name || 'Unknown',
    };

    // Add optional fields
    if (response.firstName) userData.firstName = response.firstName;
    if (response.lastName) userData.lastName = response.lastName;
    if (response.phone) userData.phone = response.phone;
    if (response.jobTitle) userData.jobTitle = response.jobTitle;
    if (response.avatar) userData.avatar = response.avatar;
    if (response.changePasswordNextLogin !== undefined) {
      userData.changePasswordNextLogin = response.changePasswordNextLogin;
    }
    if (response.created) userData.created = response.created;
    if (response.updated) userData.updated = response.updated;

    // Handle authProvider with exactOptionalPropertyTypes compatibility
    if (response.authProvider) {
      const authProvider: UserDetailData['authProvider'] = {
        type: response.authProvider.type,
      };
      if (response.authProvider.externalId) {
        authProvider.externalId = response.authProvider.externalId;
      }
      if (response.authProvider.metadata) {
        authProvider.metadata = response.authProvider.metadata;
      }
      userData.authProvider = authProvider;
    }

    // Handle new fields from API with explicit mapping for type compatibility
    if (typeof response.groupsCount === 'number') {
      userData.groupsCount = response.groupsCount;
    }
    if (Array.isArray(response.groups)) {
      userData.groups = response.groups.map(g => ({
        id: g.id,
        name: g.name,
        ...(g.description && { description: g.description }),
      }));
    }
    if (Array.isArray(response.memberships)) {
      userData.memberships = response.memberships.map(m => ({
        orgId: m.orgId,
        orgName: m.orgName,
        orgType: m.orgType,
        scope: m.scope,
        roleNames: m.roleNames,
        via: m.via,
      }));
    }

    user.value = userData;
    logger.debug('User details loaded:', { userId: userData.id, groupsCount: userData.groupsCount });
  } catch (err: any) {
    logger.error('Error fetching user:', err);
    notifyFail({ message: t.page.loadError.value });
    void router.push('/users');
  } finally {
    loading.value = false;
  }
}

/**
 * Navigate to edit user page
 *
 * @returns {void}
 */
function editUser(): void {
  if (!userId.value) return;
  void router.push(`/users/edit/${userId.value}`);
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
  void fetchUserDetails();
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
        icon="person"
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
            @click="editUser"
          />
        </template>
      </PageHeader>

      <!-- Tabs Navigation -->
      <q-card flat bordered class="rounded-borders">
        <AppTabs v-model="activeTab" :tabs="tabs" :separator="false">
          <!-- Tab Panels -->
          <q-tab-panels v-model="activeTab" animated class="bg-transparent">
            <!-- Profile Tab -->
            <q-tab-panel :name="TAB.PROFILE" class="q-pa-lg">
              <TabProfile :user="user" :loading="loading" />
            </q-tab-panel>

            <!-- Access Tab -->
            <q-tab-panel :name="TAB.ACCESS" class="q-pa-lg">
              <TabAccess :user="user" :loading="loading" />
            </q-tab-panel>

            <!-- Groups Tab -->
            <q-tab-panel :name="TAB.GROUPS" class="q-pa-lg">
              <TabGroups :user="user" :loading="loading" />
            </q-tab-panel>
          </q-tab-panels>
        </AppTabs>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
defineOptions({
  name: 'SystemSettingsPage'
});

// 1. Types
import type { OrganizationResponse } from '@mapexos/schemas';
import type { GeneralSettingsSavePayload } from './interfaces';

// 2. Vue
import { ref, computed, onMounted } from 'vue';

// 3. Components
import { PageHeader } from '@components/headers';
import { AppTabs } from '@components/tabs';
import { GeneralSettings } from '../components/GeneralSettings';
import { PersistenceSettings } from '../components/PersistenceSettings';

// 4. Composables
import { useSettingsTranslations } from '@composables/i18n';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';

// 5. Utils
import { notifySuccess, notifyFail } from '@utils/alert';

// 6. Services
import { apis } from '@services/mapex';

const { page, tabs, errors, notifications } = useSettingsTranslations();
const orgStore = useOrganizationStore();
const logger = useLogger('SystemSettingsPage');

// State
const currentTab = ref('general');
const organizationData = ref<OrganizationResponse | null>(null);
const loading = ref(false);

/**
 * Tabs configuration
 */
const settingsTabs = computed(() => [
  { name: 'general', label: tabs.general.value, icon: 'settings' },
  { name: 'persistence', label: tabs.persistence.value, icon: 'storage' },
]);

/**
 * Fetch organization data from API
 * Single source of truth - called once on mount
 * Prevents multiple fetches when switching tabs
 */
async function fetchOrganizationData() {
  if (!apis.mapexOS?.organizations) {
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  if (!orgStore.selectedOrganizationId) {
    notifyFail({ message: errors.noOrganization.value });
    return;
  }

  loading.value = true;

  try {
    const response = await apis.mapexOS.organizations.getById({
      organizationId: orgStore.selectedOrganizationId,
    });

    organizationData.value = response;
  } catch (err: any) {
    logger.error('Error fetching organization data:', err);
    notifyFail({ message: err.message || 'Failed to load organization data' });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle save from GeneralSettings component
 * Updates organization data via API and refreshes
 */
async function handleSaveGeneral(payload: GeneralSettingsSavePayload) {
  if (!apis.mapexOS?.organizations) {
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  if (!orgStore.selectedOrganizationId) return;

  loading.value = true;

  try {
    await apis.mapexOS.organizations.update(
      { organizationId: orgStore.selectedOrganizationId },
      payload as any
    );

    notifySuccess({ message: notifications.savedSuccessfully.value });

    // Refresh data after successful save
    await fetchOrganizationData();
  } catch (err: any) {
    logger.error('Error saving organization data:', err);
    notifyFail({ message: err.message || 'Failed to save organization data' });
  } finally {
    loading.value = false;
  }
}

// Fetch data on component mount
onMounted(() => void fetchOrganizationData());
</script>

<template>
  <q-page class="q-pt-lg">
    <!-- Header Section -->
    <PageHeader
      icon="settings"
      iconColor="primary"
      :title="page.title.value"
      :description="page.description.value"
      :info="page.info.value"
    />

    <!-- Tabs Container -->
    <q-card class="rounded-borders">
      <AppTabs v-model="currentTab" :tabs="settingsTabs">
        <q-tab-panels v-model="currentTab" animated>
          <!-- General Tab -->
          <q-tab-panel name="general" class="q-px-lg q-py-md">
            <GeneralSettings
              :organization-data="organizationData"
              :loading="loading"
              @save="handleSaveGeneral"
            />
          </q-tab-panel>

          <!-- Persistence Tab -->
          <q-tab-panel name="persistence" class="q-px-lg q-py-md">
            <PersistenceSettings />
          </q-tab-panel>
        </q-tab-panels>
      </AppTabs>
    </q-card>
  </q-page>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

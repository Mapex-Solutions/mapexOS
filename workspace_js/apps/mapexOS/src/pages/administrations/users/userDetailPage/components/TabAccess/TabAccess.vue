<script setup lang="ts">
defineOptions({
  name: 'TabAccess'
});

/** TYPE IMPORTS */
import type { DataRowColumn } from '@components/cards';
import type { TabAccessProps, UserMembershipInfo } from '../../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { DataRow } from '@components/cards';

/** COMPOSABLES */
import { useUserDetailTranslations } from '@composables/i18n';

/** PROPS */
const props = defineProps<TabAccessProps>();

/** COMPOSABLES & STORES */
const t = useUserDetailTranslations();

/** COMPUTED */

/**
 * Memberships from user data (passed from parent via API response)
 */
const memberships = computed<UserMembershipInfo[]>(() => props.user?.memberships || []);

/**
 * Loading state from props
 */
const loading = computed(() => props.loading);

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
    icon: (value: any, row: UserMembershipInfo) => row.orgType === 'company' ? 'business' : 'domain',
    color: (value: any, row: UserMembershipInfo) => row.orgType === 'company' ? 'blue' : 'teal',
  },
  {
    key: 'orgName',
    label: t.access.columns.organization.value,
    type: 'text',
    visible: 'always',
    width: 200,
    ellipsis: true,
  },
  {
    key: 'roleNames',
    label: t.access.columns.roles.value,
    type: 'chips',
    visible: 'laptop',
    width: 250,
    color: () => 'primary',
  },
  {
    key: 'scope',
    label: t.access.columns.scope.value,
    type: 'chip',
    visible: 'laptop',
    width: 120,
    color: (value: string) => value === 'recursive' ? 'blue' : 'orange',
    icon: (value: string) => value === 'recursive' ? 'account_tree' : 'place',
    format: (value: string) => value?.toUpperCase() || '—',
  },
  {
    key: 'via',
    label: t.access.columns.via.value,
    type: 'chip',
    visible: 'laptop',
    width: 140,
    color: (value: string) => value === 'direct' ? 'green' : 'purple',
    icon: (value: string) => value === 'direct' ? 'person' : 'groups',
    format: (value: string) => value || '—',
  },
]);

/** FUNCTIONS */

/**
 * Navigate to organization detail page (if needed)
 *
 * @param {UserMembershipInfo} membership - Membership clicked
 */
function viewMembership(membership: UserMembershipInfo): void {
  // For now, just log - could navigate to org detail or membership detail
  console.log('View membership:', membership);
}
</script>

<template>
  <div class="tab-access">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-lg">
      <q-spinner color="primary" size="40px" />
    </div>

    <!-- Memberships List -->
    <div v-else-if="memberships.length > 0" class="section">
      <!-- Section Header -->
      <div class="section-header q-mb-md">
        <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
        <span class="text-subtitle1 text-weight-medium">{{ t.access.title.value }}</span>
        <q-badge color="primary" class="q-ml-sm">{{ memberships.length }}</q-badge>
      </div>

      <!-- Membership Rows -->
      <div
        v-for="(membership, index) in memberships"
        :key="membership.orgId || `membership-${index}`"
        class="q-mb-xs"
      >
        <DataRow
          :data="membership"
          :columns="columns"
          :show-actions="false"
          @click="viewMembership"
        />
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="section text-center q-pa-xl">
      <q-icon name="admin_panel_settings" size="64px" color="grey-4" class="q-mb-md" />
      <div class="text-h6 text-grey-8 q-mb-sm">{{ t.access.empty.title.value }}</div>
      <div class="text-body2 text-grey-6">{{ t.access.empty.description.value }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.tab-access {
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

<template>
  <!-- Invisible backdrop for click outside detection -->
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="drawer-backdrop"
      @click="close"
    />
  </Teleport>

  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
    @keydown.esc="close"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon size="sm" class="q-mr-sm" name="person" color="primary" />
      <q-toolbar-title class="text-weight-medium">{{ t.drawer.title.value }}</q-toolbar-title>

      <q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
        <AppTooltip :content="t.drawer.close.value" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner size="3em" class="q-mb-md" color="primary" />
          <div class="text-grey-7">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            {{ t.drawer.error.value }}
          </q-banner>
        </div>

        <!-- User Data -->
        <div v-else-if="user" class="q-px-md q-py-lg">

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="info" color="primary" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">
                {{ getUserFullName() || '-' }}
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.email.value }}</div>
              <div class="field-value">{{ user.email || '-' }}</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.status.value }}</div>
              <div class="field-value">
                <DetailChip
                  :color="user.enabled ? 'positive' : 'negative'"
                  size="sm"
                  :label="user.enabled ? t.status.active.value.toUpperCase() : t.status.inactive.value.toUpperCase()"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.avatar.value }}</div>
              <div class="field-value">
                <div v-if="user.avatar" class="row items-center">
                  <q-avatar size="48px" class="q-mr-sm">
                    <img :src="user.avatar" />
                  </q-avatar>
                  <code class="text-caption bg-grey-2 q-pa-xs rounded-borders">{{ user.avatar }}</code>
                </div>
                <div v-else class="text-grey-7">{{ t.drawer.empty.avatar.value }}</div>
              </div>
            </div>
          </div>

          <!-- Contact & Job Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="work" color="primary" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.contact.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Phone & Job Title (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.phone.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      v-if="user.phone"
                      icon="phone"
                      color="green"
                      size="sm"
                      :label="user.phone"
                    />
                    <span v-else class="text-grey-7">{{ t.drawer.empty.phone.value }}</span>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.jobTitle.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      v-if="user.jobTitle"
                      icon="badge"
                      color="blue"
                      size="sm"
                      :label="user.jobTitle"
                    />
                    <span v-else class="text-grey-7">{{ t.drawer.empty.jobTitle.value }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Authentication Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="security" color="primary" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.authentication.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Auth Provider & Change Password (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.authProvider.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="getAuthProviderIcon()"
                      :color="getAuthProviderColorName()"
                      size="sm"
                      :label="getAuthProviderLabel()"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.changePasswordNextLogin.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="user.changePasswordNextLogin ? 'lock_reset' : 'lock'"
                      :color="user.changePasswordNextLogin ? 'warning' : 'grey'"
                      size="sm"
                      :label="user.changePasswordNextLogin ? t.filters.options.yes.value.toUpperCase() : t.filters.options.no.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- External ID (full width) -->
            <div v-if="user.authProvider?.externalId" class="field-row">
              <div class="field-label">{{ t.drawer.fields.externalId.value }}</div>
              <div class="field-value text-grey-8">
                <code class="bg-grey-2 q-pa-xs rounded-borders">{{ user.authProvider.externalId }}</code>
              </div>
            </div>
          </div>

          <!-- Groups Section -->
          <div v-if="user.groups && user.groups.length > 0" class="section q-mb-md">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="groups" color="primary" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.groups.value }} ({{ user.groupsCount || user.groups.length }})</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="groups-chips-container">
              <div v-for="group in user.groups" :key="group.id" class="group-chip-wrapper">
                <DetailChip
                  :label="group.name"
                  icon="group"
                  color="teal"
                  size="sm"
                />
                <AppTooltip
                  v-if="group.description"
                  :content="group.description"
                  :delay="500"
                  anchor="top middle"
                  self="bottom middle"
                />
              </div>
            </div>
          </div>

          <!-- Organization Access Section -->
          <div v-if="user.memberships && user.memberships.length > 0" class="section">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="domain" color="primary" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.orgAccess.value }} ({{ getOrganizedAccess().length }})</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Organization Cards -->
            <div class="org-cards-container">
              <div
                v-for="orgAccess in getOrganizedAccess()"
                :key="orgAccess.orgId"
                class="org-access-card"
              >
                <!-- Card Header: Org Name + Metadata -->
                <div class="org-card-header">
                  <div class="org-card-title">
                    <span class="org-name">{{ orgAccess.orgName }}</span>
                    <AppTooltip :content="orgAccess.orgName" />
                  </div>
                  <div class="org-card-meta">
                    <span class="org-type">{{ orgAccess.orgType }}</span>
                    <span class="meta-separator">·</span>
                    <span :class="['org-scope', orgAccess.scope === 'recursive' ? 'scope-recursive' : 'scope-local']">
                      {{ orgAccess.scope }}
                    </span>
                  </div>
                </div>

                <!-- Card Body: Access Sources -->
                <div class="org-card-body">
                  <!-- Via Groups -->
                  <div v-if="orgAccess.groups.length > 0" class="access-source">
                    <q-icon name="groups" size="xs" color="teal-6" class="source-icon" />
                    <div class="source-content">
                      <span
                        v-for="(group, idx) in orgAccess.groups.slice(0, 3)"
                        :key="group"
                        class="source-item"
                      >
                        {{ group }}<span v-if="idx < Math.min(orgAccess.groups.length, 3) - 1" class="item-separator"> · </span>
                      </span>
                      <q-btn
                        v-if="orgAccess.groups.length > 3"
                        flat
                        dense
                        no-caps
                        size="xs"
                        color="teal"
                        class="see-more-btn"
                      >
                        +{{ orgAccess.groups.length - 3 }}
                        <q-menu>
                          <q-list dense class="see-more-menu">
                            <q-item-label header class="text-teal-8">
                              <q-icon name="groups" size="xs" class="q-mr-xs" />
                              {{ t.drawer.viaGroups.allGroups.value }}
                            </q-item-label>
                            <q-item v-for="group in orgAccess.groups" :key="group" dense>
                              <q-item-section>{{ group }}</q-item-section>
                            </q-item>
                          </q-list>
                        </q-menu>
                      </q-btn>
                    </div>
                  </div>

                  <!-- Via Direct -->
                  <div v-if="orgAccess.directRoles.length > 0" class="access-source access-source--direct">
                    <q-icon name="person" size="xs" color="amber-8" class="source-icon" />
                    <div class="source-content">
                      <span
                        v-for="(role, idx) in orgAccess.directRoles.slice(0, 3)"
                        :key="role"
                        class="source-item source-item--role"
                      >
                        {{ role }}<span v-if="idx < Math.min(orgAccess.directRoles.length, 3) - 1" class="item-separator"> · </span>
                      </span>
                      <q-btn
                        v-if="orgAccess.directRoles.length > 3"
                        flat
                        dense
                        no-caps
                        size="xs"
                        color="amber-8"
                        class="see-more-btn"
                      >
                        +{{ orgAccess.directRoles.length - 3 }}
                        <q-menu>
                          <q-list dense class="see-more-menu">
                            <q-item-label header class="text-amber-9">
                              <q-icon name="person" size="xs" class="q-mr-xs" />
                              {{ t.drawer.directAccess.allRoles.value }}
                            </q-item-label>
                            <q-item v-for="role in orgAccess.directRoles" :key="role" dense>
                              <q-item-section>{{ role }}</q-item-section>
                            </q-item>
                          </q-list>
                        </q-menu>
                      </q-btn>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Timestamps Section (last - less important) -->
          <div class="section q-mt-md">
            <div class="section-header">
              <q-icon size="sm" class="q-mr-sm" name="schedule" color="grey-6" />
              <span class="text-subtitle2 text-grey-7">{{ t.drawer.sections.timestamps.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row field-row--compact">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value text-grey-7">{{ formatDate(user.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row field-row--compact">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value text-grey-7">{{ formatDate(user.updated) || '-' }}</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </q-scroll-area>
    </div>

    <!-- Footer Actions -->
    <q-separator />
    <div class="drawer-footer">
      <q-space />
      <q-btn
        unelevated
        icon="edit"
        color="primary"
        :label="t.drawer.edit.value"
        :disable="!user"
        @click="handleEdit"
      />
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'UserDetailsDrawer'
});

/** TYPE IMPORTS */
import type {
  UserDetailsDrawerProps,
  UserDetailsDrawerEmits,
  UserGroupInfo,
  UserMembershipInfo,
  OrganizedAccess,
} from './interfaces/userDetailsDrawer.interface';
import type { UserResponse } from '@mapexos/schemas';

/** Extended UserResponse with enriched fields */
interface UserResponseEnriched extends UserResponse {
  groupsCount?: number;
  groups?: UserGroupInfo[];
  memberships?: UserMembershipInfo[];
}

/** VUE IMPORTS */
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useUsersTranslations, useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

const props = defineProps<UserDetailsDrawerProps>();
const emit = defineEmits<UserDetailsDrawerEmits>();

const t = useUsersTranslations();
const errors = useCommonErrors();
const logger = useLogger('UserDetailsDrawer');

// Component state
const user = ref<UserResponseEnriched | null>(null);
const loading = ref(false);
const error = ref(false);

// Handle ESC key to close drawer
function handleEscKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});

/**
 * Fetch user details by ID from API
 * IMPORTANT: This ensures we get complete user data, not just the projection from the list
 */
async function fetchUserDetails(userId: string) {
  if (!apis.mapexOS?.users) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  user.value = null;

  try {
    const response = await apis.mapexOS.users.getById({ userId });
    user.value = response as UserResponseEnriched;
  } catch (err: any) {
    logger.error('Error fetching user details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

// Single watcher for both props to prevent duplicate fetches
// Uses a guard to prevent multiple concurrent requests
let fetchInProgress = false;

watch(
  () => [props.modelValue, props.userId] as const,
  async ([isOpen, userId]) => {
    // Reset state when drawer closes
    if (!isOpen) {
      user.value = null;
      error.value = false;
      return;
    }

    // Only fetch when drawer opens (or userId changes while open)
    if (isOpen && userId && !fetchInProgress) {
      fetchInProgress = true;
      try {
        await fetchUserDetails(userId);
      } finally {
        fetchInProgress = false;
      }
    }
  },
  { immediate: true }
);

/**
 * Get user's full name
 */
function getUserFullName(): string {
  if (!user.value) return '';
  const firstName = user.value.firstName || '';
  const lastName = user.value.lastName || '';
  return `${firstName} ${lastName}`.trim() || user.value.email || 'Unknown';
}

/**
 * Get auth provider display label
 */
function getAuthProviderLabel(): string {
  const type = user.value?.authProvider?.type || 'internal';
  return t.drawer.authProviders[type]?.value.toUpperCase() || type.toUpperCase();
}

/**
 * Get auth provider color name (DetailChip format)
 * @returns {string} Color name compatible with DetailChip
 */
function getAuthProviderColorName(): 'grey' | 'red' | 'blue' | 'orange' {
  const type = user.value?.authProvider?.type;
  const colorMap: Record<string, 'grey' | 'red' | 'blue' | 'orange'> = {
    internal: 'grey',
    google: 'red',
    github: 'grey',
    microsoft: 'blue',
    keycloak: 'orange',
  };
  return colorMap[type || 'internal'] || 'grey';
}

/**
 * Get auth provider icon
 */
function getAuthProviderIcon(): string {
  const type = user.value?.authProvider?.type;
  const iconMap: Record<string, string> = {
    internal: 'lock',
    google: 'g_translate',
    github: 'code',
    microsoft: 'business',
    keycloak: 'vpn_key',
  };
  return iconMap[type || 'internal'] || 'lock';
}

/**
 * Format date using Quasar date utils
 */
function formatDate(dateValue: any): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Organize memberships by organization for BigTech card UI
 * Groups memberships by org and separates group access from direct access
 * @returns {OrganizedAccess[]} Array of organized access per organization
 */
function getOrganizedAccess(): OrganizedAccess[] {
  if (!user.value?.memberships) return [];

  const orgMap = new Map<string, OrganizedAccess>();

  user.value.memberships.forEach(membership => {
    const orgId = membership.orgId;

    if (!orgMap.has(orgId)) {
      orgMap.set(orgId, {
        orgId,
        orgName: membership.orgName,
        orgType: membership.orgType,
        scope: membership.scope,
        groups: [],
        directRoles: [],
      });
    }

    const orgAccess = orgMap.get(orgId)!;

    // Update scope to recursive if any membership is recursive
    if (membership.scope === 'recursive') {
      orgAccess.scope = 'recursive';
    }

    // Check if access is via group or direct
    const viaLower = membership.via.toLowerCase();
    if (viaLower.includes('group:')) {
      // Extract group name from "Group: Group Name"
      const groupName = membership.via.replace(/^group:\s*/i, '').trim();
      if (groupName && !orgAccess.groups.includes(groupName)) {
        orgAccess.groups.push(groupName);
      }
    } else {
      // Direct access - add roles
      membership.roleNames.forEach(roleName => {
        if (!orgAccess.directRoles.includes(roleName)) {
          orgAccess.directRoles.push(roleName);
        }
      });
    }
  });

  return Array.from(orgMap.values());
}

/**
 * Close drawer
 */
function close() {
  emit('update:modelValue', false);
}

/**
 * Handle edit action
 */
function handleEdit() {
  if (!user.value?.id) return;
  emit('edit', user.value.id);
  close();
}
</script>

<style lang="scss" scoped>
// Flex layout for drawer content
:deep(.q-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// Drawer Header
.drawer-header {
  flex-shrink: 0;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Backdrop (teleported to body, needs :global) - transparent, just for click detection
:global(.drawer-backdrop) {
  position: fixed;
  top: 0;
  left: 0;
  right: 450px; // Leave space for drawer (450px width)
  bottom: 0;
  background: transparent;
  z-index: 5999; // Below q-drawer (6000)
  cursor: default;
}

// Drawer Content
.drawer-content {
  flex: 1;
  min-height: 0; // Important for flex children with overflow
  overflow: hidden;

  :deep(.q-scrollarea__content) {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

// Section Styling
.section {
  .section-header {
    display: flex;
    align-items: center;
    color: var(--q-primary);
    margin-bottom: 8px;
  }
}

// Field Row Styling
.field-row {
  display: flex;
  flex-direction: column;
  padding: 10px 0;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-child {
    border-bottom: none;
  }

  .field-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    margin-bottom: 4px;
    letter-spacing: 0.8px;
  }

  .field-value {
    font-size: 0.9rem;
    color: var(--mapex-text-primary);
    word-break: break-word;
    line-height: 1.4;

    code {
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
    }
  }

  // Compact variant for less important info (timestamps)
  &--compact {
    padding: 6px 0;
    border-bottom: none;

    .field-label {
      font-size: 0.65rem;
      margin-bottom: 2px;
    }

    .field-value {
      font-size: 0.8rem;
    }
  }
}

// Custom Scrollbar
:deep(.q-scrollarea__content) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: var(--mapex-radius-lg);
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(var(--q-primary-rgb), 0.3);
    border-radius: var(--mapex-radius-lg);
    transition: background var(--mapex-transition-base) ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}

// Groups Chips Container - horizontal wrap with scroll
.groups-chips-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-height: 120px;
  overflow-y: auto;
  padding: 4px 0;

  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--mapex-scrollbar-thumb);
    border-radius: var(--mapex-radius-xs);
  }

  .group-chip-wrapper {
    display: inline-flex;
    position: relative;
  }
}

// Access Summary Stats
.access-summary {
  .access-stat {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: rgba(var(--q-primary-rgb), 0.08);
    border-radius: var(--mapex-radius-md);

    .stat-value {
      font-size: 1.1rem;
      font-weight: 600;
      color: var(--q-primary);
    }

    .stat-label {
      font-size: 0.75rem;
      color: var(--mapex-text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
  }
}

// Access Table Layout
.access-table {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  overflow: hidden;

  .access-table-header {
    display: grid;
    grid-template-columns: 1.2fr 0.6fr 1.2fr;
    gap: 8px;
    padding: 10px 12px;
    background: var(--mapex-surface-bg);
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.8px;
    color: var(--mapex-text-secondary);
  }

  .access-table-body {
    max-height: 200px;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-thumb {
      background: var(--mapex-scrollbar-thumb);
      border-radius: var(--mapex-radius-xs);
    }
  }

  .access-table-row {
    display: grid;
    grid-template-columns: 1.2fr 0.6fr 1.2fr;
    gap: 8px;
    padding: 10px 12px;
    border-bottom: 1px solid var(--mapex-divider);
    align-items: start;

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: var(--mapex-surface-bg);
    }

    .col-org,
    .col-scope,
    .col-roles {
      display: flex;
      flex-direction: column;
      gap: 4px;
      min-width: 0; // Allow flex children to shrink
    }

    .col-org {
      .org-chip-wrapper {
        display: inline-flex;
        max-width: 100%;
      }

      .org-via {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 100%;
      }
    }

    .col-scope {
      justify-content: center;
      align-items: flex-start;
    }

    .roles-chips {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
    }
  }
}

// Via Groups Summary (outside table)
.via-groups-summary {
  padding: 8px 12px;
  background: rgba(var(--mapex-primary-rgb), 0.05);
  border-radius: var(--mapex-radius-sm);
  border-left: 3px solid var(--q-teal-6);

  .via-groups-list {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 2px;

    .via-group-item {
      font-size: 0.8rem;
      color: var(--mapex-text-primary);
      font-weight: 500;
    }
  }
}

// Organization Cards Container - BigTech Style
.org-cards-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.org-access-card {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-lg);
  overflow: hidden;
  transition: box-shadow var(--mapex-transition-base) ease, border-color var(--mapex-transition-base) ease;

  &:hover {
    border-color: rgba(var(--q-primary-rgb), 0.3);
    box-shadow: 0 2px 8px var(--mapex-elevation-shadow);
  }
}

.org-card-header {
  padding: 12px 14px;
  background: linear-gradient(135deg, rgba(var(--q-primary-rgb), 0.04) 0%, transparent 100%);
  border-bottom: 1px solid var(--mapex-divider);
}

.org-card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;

  .org-name {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
  }
}

.org-card-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;

  .org-type {
    color: var(--mapex-text-secondary);
    font-weight: 500;
  }

  .meta-separator {
    color: var(--mapex-text-muted);
  }

  .org-scope {
    font-weight: 600;
    padding: 1px 6px;
    border-radius: var(--mapex-radius-xs);

    &.scope-recursive {
      background: rgba(var(--mapex-primary-rgb), 0.12);
      color: var(--q-positive);
    }

    &.scope-local {
      background: rgba(var(--mapex-primary-rgb), 0.12);
      color: var(--q-info);
    }
  }
}

.org-card-body {
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.access-source {
  display: flex;
  align-items: flex-start;
  gap: 8px;

  .source-icon {
    flex-shrink: 0;
    margin-top: 2px;
  }

  .source-content {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 2px;
    font-size: 0.82rem;
    color: var(--mapex-text-primary);
    line-height: 1.5;
  }

  .source-item {
    font-weight: 500;

    &--role {
      color: var(--q-amber-9);
    }
  }

  .item-separator {
    color: var(--mapex-text-muted);
    margin: 0 1px;
  }

  .see-more-btn {
    font-weight: 600;
    margin-left: 4px;
    padding: 0 6px;
    min-height: 20px;
    border-radius: var(--mapex-radius-xs);

    &:hover {
      background: var(--mapex-surface-bg);
    }
  }
}

.access-source--direct {
  .source-content {
    color: var(--q-amber-9);
  }
}

.see-more-menu {
  min-width: 180px;
  max-height: 250px;
  overflow-y: auto;

  .q-item {
    min-height: 32px;
    font-size: 0.85rem;
  }

  .q-item-label--header {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    padding: 8px 16px 4px;
    display: flex;
    align-items: center;
  }
}
</style>

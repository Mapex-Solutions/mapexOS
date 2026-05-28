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
      <q-icon
        name="admin_panel_settings"
        size="sm"
        color="primary"
        class="q-mr-sm"
      />
      <q-toolbar-title class="text-weight-medium">{{ t.drawer.title.value }}</q-toolbar-title>

      <q-btn
        flat
        round
        dense
        icon="close"
        class="drawer-close-btn"
        @click="close"
      >
        <AppTooltip :content="t.drawer.close.value" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner
            size="3em"
            color="primary"
            class="q-mb-md"
          />
          <div class="text-grey-7">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner
            rounded
            class="bg-negative text-white"
          >
            <template #avatar>
              <q-icon
                name="error"
                color="white"
              />
            </template>
            {{ t.drawer.error.value }}
          </q-banner>
        </div>

        <!-- Role Data -->
        <div v-else-if="role" class="q-px-md q-py-lg">

          <!-- System Role Warning -->
          <div v-if="role.isSystem" class="q-mb-md">
            <q-banner
              rounded
              class="bg-warning text-white"
            >
              <template #avatar>
                <q-icon
                  name="lock"
                  color="white"
                />
              </template>
              {{ t.drawer.systemRoleWarning.value }}
            </q-banner>
          </div>

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="info"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ role?.name || '-' }}</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.type.value }}</div>
              <div class="field-value">
                <DetailChip
                  :icon="role?.isSystem ? 'lock' : 'lock_open'"
                  :color="role?.isSystem ? 'purple' : 'blue'"
                  size="sm"
                  :label="role?.isSystem ? t.drawer.system.yes.value.toUpperCase() : t.drawer.system.no.value.toUpperCase()"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ role?.description || t.drawer.empty.description.value }}
              </div>
            </div>
          </div>

          <!-- Scope & Organization Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="public"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.scope.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Scope & Organization (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.scope.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="role?.scope === 'global' ? 'public' : 'place'"
                      :color="role?.scope === 'global' ? 'purple' : 'orange'"
                      size="sm"
                      :label="role?.scope?.toUpperCase() || 'N/A'"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.organization.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      icon="domain"
                      color="green"
                      size="sm"
                      :label="(role as any)?.organizationName || role?.orgId || '-'"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Path Key (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.pathKey.value }}</div>
              <div class="field-value text-grey-8">
                <code class="bg-grey-2 q-pa-xs rounded-borders">{{ role?.pathKey || '-' }}</code>
              </div>
            </div>
          </div>

          <!-- Permissions Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="vpn_key"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.permissions.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.permissionsCount.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="vpn_key"
                  color="blue"
                  size="sm"
                  :label="String(role?.permissions?.length || 0)"
                />
              </div>
            </div>

            <!-- Permissions Accordion -->
            <div v-if="groupedPermissions.length > 0" class="q-mt-md permissions-accordion">
              <q-list bordered separator>
                <q-expansion-item
                  v-for="group in groupedPermissions"
                  expand-separator
                  header-class="text-primary"
                  :key="group.resource"
                  :icon="getResourceIcon(group.resource)"
                  :label="formatResourceName(group.resource)"
                  :caption="`${group.count} ${group.count === 1 ? 'permission' : 'permissions'}`"
                >
                  <template #header>
                    <q-item-section avatar>
                      <q-icon
                        :name="getResourceIcon(group.resource)"
                        color="primary"
                      />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ formatResourceName(group.resource) }}</q-item-label>
                      <q-item-label caption>{{ group.count }} {{ group.count === 1 ? 'permission' : 'permissions' }}</q-item-label>
                    </q-item-section>
                  </template>
                  <q-card flat>
                    <q-card-section class="q-pa-sm">
                      <div class="row q-col-gutter-xs">
                        <div
                          v-for="(action, idx) in group.actions"
                          class="col-auto"
                          :key="idx"
                        >
                          <DetailChip
                            icon="check_circle"
                            color="secondary"
                            size="sm"
                            :label="formatAction(action)"
                          />
                        </div>
                      </div>
                    </q-card-section>
                  </q-card>
                </q-expansion-item>
              </q-list>
            </div>
            <div v-else class="text-grey-7 text-caption q-mt-sm">
              {{ t.drawer.empty.permissions.value }}
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header">
              <q-icon
                name="schedule"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value">{{ formatDate(role?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(role?.updated) }}</div>
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
        :disable="!role || isSystemRole"
        @click="handleEdit"
      >
        <AppTooltip
          v-if="isSystemRole"
          :content="t.drawer.systemRoleTooltip.value"
        />
      </q-btn>
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'RoleDetailsDrawer'
});

/** TYPE IMPORTS */
import type { RoleDetailsDrawerProps, RoleDetailsDrawerEmits, PermissionGroup } from './interfaces';
import type { RoleResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useRolesTranslations, useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

const props = defineProps<RoleDetailsDrawerProps>();
const emit = defineEmits<RoleDetailsDrawerEmits>();

const t = useRolesTranslations();
const errors = useCommonErrors();
const logger = useLogger('RoleDetailsDrawer');

// Component state
const role = ref<RoleResponse | null>(null);
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
 * Fetch role details by ID from API
 * IMPORTANT: This ensures we get complete role data, not just the projection from the list
 */
async function fetchRoleDetails(roleId: string) {
  if (!apis.mapexOS?.roles) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  role.value = null;

  try {
    const response = await apis.mapexOS.roles.getById({ roleId });
    role.value = response;
  } catch (err: any) {
    logger.error('Error fetching role details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

// Watch for roleId changes to fetch new data
watch(() => props.roleId, (newRoleId) => {
  if (newRoleId && props.modelValue) {
    void fetchRoleDetails(newRoleId);
  }
}, { immediate: true });

// Watch for drawer open/close to fetch data when opened
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.roleId) {
    void fetchRoleDetails(props.roleId);
  } else if (!isOpen) {
    // Reset state when drawer closes
    role.value = null;
    error.value = false;
  }
});

// Check if role is system role
const isSystemRole = computed(() => {
  return role.value?.isSystem === true;
});

/**
 * Group permissions by resource
 * Example: ['users.list', 'users.read', 'assets.create'] becomes:
 * {
 *   'users': ['list', 'read'],
 *   'assets': ['create']
 * }
 */
const groupedPermissions = computed<PermissionGroup[]>(() => {
  if (!role.value?.permissions) return [];

  // Handle wildcard permission (mapex.* or similar)
  if (role.value.permissions.includes('mapex.*')) {
    return [{
      resource: 'mapex',
      actions: ['*'],
      count: 1,
    }];
  }

  const groups: Record<string, string[]> = {};

  role.value.permissions.forEach((permission: string) => {
    const parts = permission.split('.');

    // Handle resource.action pattern
    if (parts.length === 2 && parts[0] && parts[1]) {
      const resource = parts[0];
      const action = parts[1];
      if (!groups[resource]) groups[resource] = [];
      groups[resource].push(action);
    }
    // Handle resource.subresource.action pattern (e.g., events.raw.list)
    else if (parts.length === 3 && parts[0] && parts[1] && parts[2]) {
      const resource = `${parts[0]}.${parts[1]}`;
      const action = parts[2];
      if (!groups[resource]) groups[resource] = [];
      groups[resource].push(action);
    }
    // Handle other patterns
    else if (parts.length > 0) {
      const resource = parts.slice(0, -1).join('.');
      const action = parts[parts.length - 1];
      if (resource && action) {
        if (!groups[resource]) groups[resource] = [];
        groups[resource].push(action);
      }
    }
  });

  // Convert to array and sort
  return Object.entries(groups)
    .map(([resource, actions]) => ({
      resource,
      actions: actions.sort(),
      count: actions.length,
    }))
    .sort((a, b) => a.resource.localeCompare(b.resource));
});

/**
 * Get icon for resource type
 */
function getResourceIcon(resource: string): string {
  const iconMap: Record<string, string> = {
    'auth': 'login',
    'users': 'person',
    'roles': 'admin_panel_settings',
    'organizations': 'domain',
    'groups': 'group',
    'memberships': 'card_membership',
    'assets': 'router',
    'assettemplates': 'memory',
    'lists': 'list_alt',
    'routegroups': 'account_tree',
    'datasources': 'storage',
    'events': 'event',
    'events.raw': 'data_object',
    'events.processed': 'done_all',
    'events.js_executor': 'code',
    'events.router': 'alt_route',
    'events.business_rule': 'rule',
    'events.audit': 'history',
    'events.notifications': 'notifications',
    'mapex': 'admin_panel_settings',
  };
  return iconMap[resource] || 'vpn_key';
}

/**
 * Format action for display
 */
function formatAction(action: string): string {
  // Handle wildcard
  if (action === '*') return 'All operations';

  // Capitalize first letter
  return action.charAt(0).toUpperCase() + action.slice(1);
}

/**
 * Format resource name for display
 */
function formatResourceName(resource: string): string {
  // Handle special cases
  const nameMap: Record<string, string> = {
    'auth': 'Authentication',
    'assettemplates': 'Asset Templates',
    'routegroups': 'Route Groups',
    'datasources': 'Data Sources',
    'events.raw': 'Events - Raw',
    'events.processed': 'Events - Processed',
    'events.js_executor': 'Events - JS Executor',
    'events.router': 'Events - Router',
    'events.business_rule': 'Events - Business Rules',
    'events.audit': 'Events - Audit',
    'events.notifications': 'Events - Notifications',
    'mapex': 'Full Platform Access',
  };

  if (nameMap[resource]) return nameMap[resource];

  // Default: capitalize first letter
  return resource.charAt(0).toUpperCase() + resource.slice(1);
}

// Format date using Quasar date utils
function formatDate(dateValue: any): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

// Close drawer
function close() {
  emit('update:modelValue', false);
}

// Handle edit action
function handleEdit() {
  if (isSystemRole.value) return;
  if (!role.value?.id) return;
  emit('edit', role.value.id);
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
}

// Permissions Accordion
.permissions-accordion {
  border-radius: var(--mapex-radius-md);
  overflow: hidden;

  .q-item {
    transition: all var(--mapex-transition-base) ease;
  }

  :deep(.q-expansion-item__container) {
    border-radius: var(--mapex-radius-md);
    margin-bottom: 4px;
    overflow: hidden;
  }

  :deep(.q-item__section--avatar) {
    min-width: 40px;
  }

  :deep(.q-expansion-item__content) {
    background: var(--mapex-surface-bg);
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
</style>

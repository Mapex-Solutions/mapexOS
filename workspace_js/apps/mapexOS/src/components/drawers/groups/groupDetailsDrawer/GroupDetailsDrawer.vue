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
        name="groups"
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

        <!-- Group Data -->
        <div v-else-if="group" class="q-px-md q-py-lg">

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
              <div class="field-value text-weight-medium">{{ group?.name || '-' }}</div>
            </div>

            <!-- Status -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.enabled.value }}</div>
              <div class="field-value">
                <DetailChip
                  :icon="group?.enabled ? 'check_circle' : 'cancel'"
                  :color="group?.enabled ? 'green' : 'red'"
                  size="sm"
                  :label="group?.enabled ? t.drawer.status.enabled.value.toUpperCase() : t.drawer.status.disabled.value.toUpperCase()"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ group?.description || t.drawer.empty.description.value }}
              </div>
            </div>
          </div>

          <!-- Organization Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="domain"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.organization.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.organization.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="domain"
                  color="green"
                  size="sm"
                  :label="(group as any)?.organizationName || group?.orgId || '-'"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.pathKey.value }}</div>
              <div class="field-value text-grey-8">
                <code class="bg-grey-2 q-pa-xs rounded-borders">{{ group?.pathKey || '-' }}</code>
              </div>
            </div>
          </div>

          <!-- Roles Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="admin_panel_settings"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.roles.value }}</span>
              <q-badge v-if="rolesList.length > 0" color="primary" class="q-ml-sm">{{ rolesList.length }}</q-badge>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Roles List -->
            <div v-if="rolesList.length > 0" class="roles-list">
              <q-list bordered separator>
                <q-item
                  v-for="role in displayedRoles"
                  :key="role.id"
                  dense
                >
                  <q-item-section avatar>
                    <q-avatar
                      size="32px"
                      :color="role.isSystem ? 'purple' : 'primary'"
                      text-color="white"
                    >
                      <q-icon name="admin_panel_settings" size="18px" />
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-medium">{{ role.name }}</q-item-label>
                    <q-item-label
                      v-if="role.description"
                      caption
                      lines="1"
                    >{{ role.description }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-badge
                      :color="role.isSystem ? 'purple' : 'blue'"
                      :label="role.isSystem ? 'SYSTEM' : 'CUSTOM'"
                    />
                  </q-item-section>
                </q-item>
              </q-list>
              <!-- View More Roles -->
              <q-btn
                v-if="rolesList.length > MAX_DISPLAY_ITEMS"
                flat
                dense
                color="primary"
                class="full-width q-mt-sm"
                :label="`+ ${rolesList.length - MAX_DISPLAY_ITEMS} more`"
                @click="navigateToDetail('roles')"
              />
            </div>
            <div v-else class="text-grey-7 text-caption q-mt-sm">
              {{ t.drawer.empty.roles.value }}
            </div>
          </div>

          <!-- Members Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="person"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.members.value }}</span>
              <q-badge v-if="group?.membersCount" color="primary" class="q-ml-sm">{{ group.membersCount }}</q-badge>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Members List -->
            <div v-if="membersList.length > 0" class="members-list">
              <q-list bordered separator>
                <q-item
                  v-for="member in displayedMembers"
                  :key="member.id"
                  dense
                >
                  <q-item-section avatar>
                    <q-avatar
                      size="32px"
                      color="primary"
                      text-color="white"
                    >
                      {{ getInitials(member.name) }}
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-medium">{{ member.name }}</q-item-label>
                    <q-item-label
                      v-if="member.email"
                      caption
                    >{{ member.email }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
              <!-- View More Members -->
              <q-btn
                v-if="membersList.length > MAX_DISPLAY_ITEMS"
                flat
                dense
                color="primary"
                class="full-width q-mt-sm"
                :label="`+ ${membersList.length - MAX_DISPLAY_ITEMS} more`"
                @click="navigateToDetail('members')"
              />
            </div>
            <div v-else class="text-grey-7 text-caption q-mt-sm">
              {{ t.drawer.empty.members.value }}
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
                  <div class="field-value">{{ formatDate(group?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(group?.updated) }}</div>
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
        :disable="!group"
        @click="handleEdit"
      >
      </q-btn>
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'GroupDetailsDrawer'
});

/** TYPE IMPORTS */
import type { GroupDetailsDrawerProps, GroupDetailsDrawerEmits, MemberDisplayItem, RoleDisplayItem } from './interfaces';
import type { GroupResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { date } from 'quasar';

/** CONSTANTS */
const MAX_DISPLAY_ITEMS = 3;

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useGroupsTranslations, useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** PROPS & EMITS */
const props = defineProps<GroupDetailsDrawerProps>();
const emit = defineEmits<GroupDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();
const errors = useCommonErrors();
const router = useRouter();
const logger = useLogger('GroupDetailsDrawer');
const orgStore = useOrganizationStore();

/** STATE */
const group = ref<GroupResponse | null>(null);
const membersList = ref<MemberDisplayItem[]>([]);
const rolesList = ref<RoleDisplayItem[]>([]);
const loading = ref(false);
const error = ref(false);

/** FUNCTIONS */

/**
 * Handle ESC key to close drawer
 *
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

/**
 * Fetch group details by ID from API
 * IMPORTANT: This ensures we get complete group data, not just the projection from the list
 *
 * @param {string} groupId - The group ID to fetch
 */
async function fetchGroupDetails(groupId: string): Promise<void> {
  if (!apis.mapexOS?.groups) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  group.value = null;
  membersList.value = [];
  rolesList.value = [];

  try {
    const response = await apis.mapexOS.groups.getById({ groupId });

    // Enrich with organization name
    const organization = orgStore.flatList.find((org: any) => org.id === response.orgId);
    group.value = {
      ...response,
      organizationName: organization?.name || 'Unknown',
    } as any;

    // Fetch role details if roleIds available
    const roleIds = (response as any).roleIds as string[] | undefined;
    if (roleIds?.length) {
      await fetchGroupRoles(roleIds);
    }

    // Fetch member details if membersCount > 0
    if (response.membersCount && response.membersCount > 0) {
      await fetchGroupMembers(groupId);
    }
  } catch (err: any) {
    logger.error('Error fetching group details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Fetch role details for display
 *
 * @param {string[]} roleIds - Array of role IDs to fetch
 */
async function fetchGroupRoles(roleIds: string[]): Promise<void> {
  if (!apis.mapexOS?.roles || !roleIds.length) {
    rolesList.value = [];
    return;
  }

  try {
    // Fetch all roles and filter by roleIds
    const response = await apis.mapexOS.roles.list({ perPage: 100 });
    const roleIdsSet = new Set(roleIds);

    rolesList.value = (response.items || [])
      .filter((role: any) => roleIdsSet.has(role.id))
      .map((role: any) => ({
        id: role.id || '',
        name: role.name || '',
        description: role.description || '',
        isSystem: role.isSystem ?? false,
      }));

    logger.debug('Fetched roles for group:', { count: rolesList.value.length });
  } catch (err: any) {
    logger.warn('Failed to fetch roles:', err);
    rolesList.value = [];
  }
}

/**
 * Fetch group members via paginated API
 * Loads first page of members (max 100)
 *
 * @param {string} groupId - The group ID to fetch members for
 */
async function fetchGroupMembers(groupId: string): Promise<void> {
  if (!apis.mapexOS?.groups) return;

  try {
    const membersResponse = await apis.mapexOS.groups.getMembers(
      { groupId },
      { page: 1, perPage: 100 },
    );

    // Extract user IDs from member response
    const memberIds = (membersResponse.items || [])
      .map((m: any) => m.userId)
      .filter(Boolean);

    if (memberIds.length > 0) {
      await fetchMemberDetails(memberIds);
    }
  } catch (err: any) {
    logger.warn('Failed to fetch group members:', err);
  }
}

/**
 * Fetch member details for display
 *
 * @param {string[]} memberIds - Array of member IDs
 */
async function fetchMemberDetails(memberIds: string[]): Promise<void> {
  if (!apis.mapexOS?.users) {
    // If users API not available, show IDs as fallback
    membersList.value = memberIds.map(id => ({
      id,
      name: id,
    }));
    return;
  }

  try {
    // Fetch user details for each member
    const members: MemberDisplayItem[] = [];

    for (const memberId of memberIds) {
      try {
        const user = await apis.mapexOS.users.getById({ userId: memberId });
        // Build full name from firstName and lastName
        const fullName = [user.firstName, user.lastName].filter(Boolean).join(' ') || user.email || memberId;
        const member: MemberDisplayItem = {
          id: memberId,
          name: fullName,
        };
        if (user.email) {
          member.email = user.email;
        }
        if (user.avatar) {
          member.avatar = user.avatar;
        }
        members.push(member);
      } catch {
        // If user fetch fails, add with ID as name
        members.push({
          id: memberId,
          name: memberId,
        });
      }
    }

    membersList.value = members;
  } catch (err: any) {
    logger.warn('Error fetching member details:', err);
    // Fallback to showing IDs
    membersList.value = memberIds.map(id => ({
      id,
      name: id,
    }));
  }
}

/**
 * Get initials from name
 *
 * @param {string} name - Full name
 * @returns {string} Initials (up to 2 characters)
 */
function getInitials(name: string): string {
  if (!name) return '?';
  const parts = name.split(' ').filter(Boolean);
  if (parts.length >= 2) {
    const firstPart = parts[0] || '';
    const lastPart = parts[parts.length - 1] || '';
    return (firstPart.charAt(0) + lastPart.charAt(0)).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
}

/**
 * Format date using Quasar date utils
 *
 * @param {any} dateValue - Date value to format
 * @returns {string} Formatted date string
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
 * Close drawer
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Handle edit action
 */
function handleEdit(): void {
  if (!group.value?.id) return;
  emit('edit', group.value.id);
  close();
}

/**
 * Navigate to group detail page with specific tab
 *
 * @param {string} tab - Tab to navigate to ('roles' or 'members')
 */
function navigateToDetail(tab: 'roles' | 'members'): void {
  if (!group.value?.id) return;
  close();
  void router.push(`/groups/${group.value.id}?tab=${tab}`);
}

/** COMPUTED */

/**
 * Roles list limited to MAX_DISPLAY_ITEMS
 */
const displayedRoles = computed(() => {
  return rolesList.value.slice(0, MAX_DISPLAY_ITEMS);
});

/**
 * Members list limited to MAX_DISPLAY_ITEMS
 */
const displayedMembers = computed(() => {
  return membersList.value.slice(0, MAX_DISPLAY_ITEMS);
});

/** WATCHERS */

// Watch for groupId changes to fetch new data
watch(() => props.groupId, (newGroupId) => {
  if (newGroupId && props.modelValue) {
    void fetchGroupDetails(newGroupId);
  }
}, { immediate: true });

// Watch for drawer open/close to fetch data when opened
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.groupId) {
    void fetchGroupDetails(props.groupId);
  } else if (!isOpen) {
    // Reset state when drawer closes
    group.value = null;
    membersList.value = [];
    error.value = false;
  }
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
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

// Roles List
.roles-list {
  border-radius: var(--mapex-radius-md);
  overflow: hidden;

  .q-item {
    transition: all var(--mapex-transition-base) ease;

    &:hover {
      background: var(--mapex-surface-bg);
    }
  }
}

// Members List
.members-list {
  border-radius: var(--mapex-radius-md);
  overflow: hidden;

  .q-item {
    transition: all var(--mapex-transition-base) ease;

    &:hover {
      background: var(--mapex-surface-bg);
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
</style>

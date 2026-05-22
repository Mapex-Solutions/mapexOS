<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="vpn_key" color="primary" class="q-mr-xs" />
        {{ t.sections.permissions.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.permissions.value }}
      </div>
    </div>

    <!-- Selected Count -->
    <div class="q-mb-md">
      <DetailChip
        :value="`${selectedCount} ${selectedCount === 1 ? t.labels.permissionSelected.value : t.labels.permissionsSelected.value}`"
        icon="check_circle"
        color="primary"
        size="md"
      />
    </div>

    <!-- Permissions by Group -->
    <div class="permissions-list">
      <div
        v-for="group in groupedPermissions"
        :key="group.label"
        class="q-mb-lg"
      >
        <!-- Group Header -->
        <div class="group-header q-mb-sm">
          <q-icon :name="group.icon" color="grey-8" size="sm" class="q-mr-sm" />
          <span class="text-subtitle2 text-weight-bold text-grey-8">{{ group.label }}</span>
          <q-badge
            :label="`${getGroupGrantedCount(group)} / ${getGroupTotalCount(group)}`"
            color="grey-4"
            text-color="grey-8"
            class="q-ml-sm"
          />
          <q-btn
            flat
            round
            dense
            size="sm"
            :icon="collapsedGroups.has(group.label) ? 'expand_more' : 'expand_less'"
            :color="collapsedGroups.has(group.label) ? 'grey-5' : 'grey-7'"
            class="q-ml-xs"
            @click="toggleGroupVisibility(group.label)"
          >
            <AppTooltip
              :content="collapsedGroups.has(group.label) ? 'Show permissions' : 'Hide permissions'"
            />
          </q-btn>
        </div>

        <!-- Resources in group -->
        <q-slide-transition>
          <q-list v-show="!collapsedGroups.has(group.label)" bordered separator class="rounded-borders">
          <q-expansion-item
            v-for="item in group.items"
            :key="item.resource.resource"
            expand-separator
            :header-class="item.resource.enabled ? 'bg-primary-subtle' : ''"
          >
            <template #header>
              <q-item-section avatar>
                <q-checkbox
                  :model-value="item.resource.enabled"
                  color="primary"
                  @click.stop
                  @update:model-value="onResourceToggle(item.globalIndex)"
                />
              </q-item-section>
              <q-item-section avatar>
                <q-icon
                  :name="item.resource.icon"
                  :color="item.resource.enabled ? 'primary' : 'grey-6'"
                />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">
                  {{ item.resource.label }}
                </q-item-label>
                <q-item-label caption>
                  {{ getResourceCaption(item.resource) }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="row q-gutter-xs">
                  <q-btn
                    flat
                    dense
                    size="sm"
                    :label="t.labels.selectAll.value"
                    color="primary"
                    @click.stop="onToggleAllActions(item.globalIndex, true)"
                  />
                  <q-btn
                    flat
                    dense
                    size="sm"
                    :label="t.labels.deselectAll.value"
                    color="grey-7"
                    @click.stop="onToggleAllActions(item.globalIndex, false)"
                  />
                </div>
              </q-item-section>
            </template>

            <q-card flat>
              <q-card-section class="q-pa-md">
                <div class="row q-col-gutter-md">
                  <div
                    v-for="(action, actionIndex) in item.resource.actions"
                    :key="action.name"
                    class="col-12 col-sm-6 col-md-4"
                  >
                    <q-item
                      tag="label"
                      class="action-item rounded-borders"
                      :class="{ 'action-item--selected': action.granted }"
                    >
                      <q-item-section avatar>
                        <q-checkbox
                          :model-value="action.granted"
                          color="primary"
                          @update:model-value="onActionToggle(item.globalIndex, actionIndex)"
                        />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">
                          {{ action.label }}
                        </q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-icon
                          :name="getActionIcon(action.name)"
                          :color="action.granted ? 'primary' : 'grey-5'"
                          size="xs"
                        />
                      </q-item-section>
                    </q-item>
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </q-expansion-item>
        </q-list>
        </q-slide-transition>
      </div>
    </div>

    <!-- Validation message -->
    <div v-if="showValidationError" class="text-negative text-caption q-mt-md">
      <q-icon name="warning" size="xs" class="q-mr-xs" />
      {{ t.validation.permissionsRequired.value }}
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2Permissions',
});

/** TYPE IMPORTS */
import type { QForm } from 'quasar';
import type { ResourcePermission, PermissionGroup } from '../../interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useRolesTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { PERMISSION_GROUPS } from '../../constants';

/** PROPS & EMITS */
interface Props {
  resourcePermissions: ResourcePermission[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'resource-toggle', index: number): void;
  (e: 'action-toggle', resourceIndex: number, actionIndex: number): void;
  (e: 'toggle-all-actions', resourceIndex: number, granted: boolean): void;
}>();

/** COMPOSABLES & STORES */
const t = useRolesTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showValidationError = ref(false);
const collapsedGroups = ref<Set<string>>(new Set());

/** LOCAL IMPORTS */
import type { GroupedResourceItem, GroupedPermissionSection } from './interfaces/Step2Permissions.interface';

/** COMPUTED */

/**
 * Group resources by PERMISSION_GROUPS, preserving global indices for event emission
 */
const groupedPermissions = computed((): GroupedPermissionSection[] => {
  return PERMISSION_GROUPS
    .map((group: PermissionGroup) => ({
      label: group.label,
      icon: group.icon,
      items: group.resources
        .map((resourceKey: string) => {
          const globalIndex = props.resourcePermissions.findIndex(
            (r: ResourcePermission) => r.resource === resourceKey,
          );
          if (globalIndex === -1) return null;
          return {
            resource: props.resourcePermissions[globalIndex],
            globalIndex,
          };
        })
        .filter((item): item is GroupedResourceItem => item !== null),
    }))
    .filter((group: GroupedPermissionSection) => group.items.length > 0);
});

/**
 * Count total selected permissions
 */
const selectedCount = computed(() => {
  let count = 0;
  props.resourcePermissions.forEach(resource => {
    resource.actions.forEach(action => {
      if (action.granted) count++;
    });
  });
  return count;
});

/** FUNCTIONS */

/**
 * Toggle visibility of a permission group's resource list
 *
 * @param {string} groupLabel - Group label used as key
 */
function toggleGroupVisibility(groupLabel: string): void {
  const next = new Set(collapsedGroups.value);
  if (next.has(groupLabel)) {
    next.delete(groupLabel);
  } else {
    next.add(groupLabel);
  }
  collapsedGroups.value = next;
}

/**
 * Get granted permission count for a group
 *
 * @param {GroupedPermissionSection} group - Permission group
 * @returns {number} Number of granted actions
 */
function getGroupGrantedCount(group: GroupedPermissionSection): number {
  let count = 0;
  group.items.forEach(item => {
    item.resource.actions.forEach(action => {
      if (action.granted) count++;
    });
  });
  return count;
}

/**
 * Get total permission count for a group
 *
 * @param {GroupedPermissionSection} group - Permission group
 * @returns {number} Total number of actions
 */
function getGroupTotalCount(group: GroupedPermissionSection): number {
  let count = 0;
  group.items.forEach(item => {
    count += item.resource.actions.length;
  });
  return count;
}

/**
 * Get caption text for resource
 *
 * @param {ResourcePermission} resource - Resource permission object
 * @returns {string} Caption text
 */
function getResourceCaption(resource: ResourcePermission): string {
  const grantedCount = resource.actions.filter(a => a.granted).length;
  const totalCount = resource.actions.length;
  return `${grantedCount} / ${totalCount} ${t.labels.actionsSelected.value}`;
}

/**
 * Get icon for action type
 *
 * @param {string} actionName - Action name
 * @returns {string} Icon name
 */
function getActionIcon(actionName: string): string {
  const iconMap: Record<string, string> = {
    list: 'list',
    create: 'add_circle',
    read: 'visibility',
    update: 'edit',
    delete: 'delete',
  };
  return iconMap[actionName] || 'check_circle';
}

/**
 * Handle resource toggle
 *
 * @param {number} index - Resource index
 */
function onResourceToggle(index: number): void {
  showValidationError.value = false;
  emit('resource-toggle', index);
}

/**
 * Handle action toggle
 *
 * @param {number} resourceIndex - Resource index
 * @param {number} actionIndex - Action index
 */
function onActionToggle(resourceIndex: number, actionIndex: number): void {
  showValidationError.value = false;
  emit('action-toggle', resourceIndex, actionIndex);
}

/**
 * Handle toggle all actions
 *
 * @param {number} resourceIndex - Resource index
 * @param {boolean} granted - Whether to grant or revoke
 */
function onToggleAllActions(resourceIndex: number, granted: boolean): void {
  showValidationError.value = false;
  emit('toggle-all-actions', resourceIndex, granted);
}

/**
 * Validate that at least one permission is selected
 *
 * @returns {boolean} Whether validation passed
 */
function validate(): boolean {
  const hasPermissions = selectedCount.value > 0;
  if (!hasPermissions) {
    showValidationError.value = true;
  }
  return hasPermissions;
}

/** EXPOSE */
defineExpose({
  formRef,
  validate,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.group-header {
  display: flex;
  align-items: center;
  padding: 4px 0;
}

.permissions-list {
  :deep(.q-expansion-item__container) {
    margin-bottom: 4px;
    border-radius: var(--mapex-radius-md);
    overflow: hidden;
  }

  :deep(.q-item__section--avatar) {
    min-width: 40px;
  }
}

.bg-primary-subtle {
  background-color: rgba(var(--q-primary-rgb), 0.08);
}

.action-item {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);
  min-height: 56px;

  &:hover {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.04);
  }

  &--selected {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.08);
  }
}
</style>

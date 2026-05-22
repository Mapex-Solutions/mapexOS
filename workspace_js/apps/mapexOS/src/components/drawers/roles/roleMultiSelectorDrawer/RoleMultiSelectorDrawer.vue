<script setup lang="ts">
defineOptions({
  name: 'RoleMultiSelectorDrawer'
});

/** TYPE IMPORTS */
import type { RoleMultiSelectorDrawerProps, RoleMultiSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { SelectableChip } from '@components/chips';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** PROPS & EMITS */
const props = withDefaults(defineProps<RoleMultiSelectorDrawerProps>(), {
  selectedRoleIds: () => [],
});

const emit = defineEmits<RoleMultiSelectorDrawerEmits>();

/** STATE */
const loading = ref(false);
const loadingMore = ref(false);
const roles = ref<any[]>([]);
const selectedRoles = ref<any[]>([]);
const scrollAreaRef = ref<any>(null);

/** PAGINATION STATE */
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(15);

/** FILTER STATE */
const filters = ref({
  name: undefined as string | undefined,
  isSystem: undefined as boolean | undefined,
  isTemplate: undefined as boolean | undefined,
});

/** COMPUTED */

/**
 * Drawer visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Type filter options (System vs Custom)
 */
const typeOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'System', value: true },
  { label: 'Custom', value: false },
]);

/**
 * Scope filter options (Template = shared, Non-template = local)
 */
const templateOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'Shared (Templates)', value: true },
  { label: 'Local Only', value: false },
]);

/**
 * Selected role IDs for quick lookup
 */
const selectedRoleIdsSet = computed(() => new Set(selectedRoles.value.map(r => r.id)));

/**
 * Check if confirm button should be enabled
 */
const canConfirm = computed(() => selectedRoles.value.length > 0);

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    // Initialize selection from props
    selectedRoles.value = [];
    void fetchRoles();
  }
});

/**
 * Watch selectedRoleIds prop and sync with loaded roles
 */
watch(() => [props.selectedRoleIds, roles.value], () => {
  if (props.selectedRoleIds && props.selectedRoleIds.length > 0 && roles.value.length > 0) {
    const matchingRoles = roles.value.filter(r => props.selectedRoleIds?.includes(r.id));
    // Merge with existing selection, avoiding duplicates
    const existingIds = new Set(selectedRoles.value.map(r => r.id));
    matchingRoles.forEach(role => {
      if (!existingIds.has(role.id)) {
        selectedRoles.value.push(role);
      }
    });
  }
}, { deep: true });

/** FUNCTIONS */

/**
 * Fetch roles from API with current filters and pagination
 * @param {boolean} append - If true, append to existing list (for infinite scroll)
 * @returns {Promise<void>}
 */
async function fetchRoles(append: boolean = false): Promise<void> {
  if (append) {
    loadingMore.value = true;
  } else {
    loading.value = true;
    currentPage.value = 1;
  }

  try {
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: perPage.value,
    };

    // Add filters conditionally
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.isSystem === 'boolean') {
      queryParams.isSystem = filters.value.isSystem;
    }
    if (typeof filters.value.isTemplate === 'boolean') {
      queryParams.isTemplate = filters.value.isTemplate;
    }

    const response = await apis.mapexOS.roles.list(queryParams);

    if (append) {
      roles.value = [...roles.value, ...(response.items || [])];
    } else {
      roles.value = response.items || [];
    }

    totalPages.value = response.pagination?.totalPages || 1;
    totalItems.value = response.pagination?.totalItems || 0;
  } catch (error: any) {
    handleApiError({
      error,
      customMessage: 'Failed to fetch roles'
    });
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
}

/**
 * Infinite scroll handler
 * @param {any} info - Scroll information
 */
function onScroll(info: any): void {
  const scrollPosition = info.verticalPosition;
  const scrollSize = info.verticalSize;
  const containerSize = info.verticalContainerSize;

  if (scrollPosition + containerSize >= scrollSize * 0.8) {
    if (!loadingMore.value && currentPage.value < totalPages.value) {
      currentPage.value++;
      void fetchRoles(true);
    }
  }
}

/**
 * Toggle role selection
 * @param {any} role - Role to toggle
 */
function toggleRoleSelection(role: any): void {
  const index = selectedRoles.value.findIndex(r => r.id === role.id);
  if (index >= 0) {
    selectedRoles.value.splice(index, 1);
  } else {
    selectedRoles.value.push(role);
  }
}

/**
 * Check if role is selected
 * @param {any} role - Role to check
 * @returns {boolean} True if selected
 */
function isSelected(role: any): boolean {
  return selectedRoleIdsSet.value.has(role.id);
}

/**
 * Remove a role from selection
 * @param {any} role - Role to remove
 */
function removeFromSelection(role: any): void {
  const index = selectedRoles.value.findIndex(r => r.id === role.id);
  if (index >= 0) {
    selectedRoles.value.splice(index, 1);
  }
}

/**
 * Clear all selections
 */
function clearSelection(): void {
  selectedRoles.value = [];
}

/**
 * Filter change handler
 * Resets pagination and refetches roles
 */
function onFilterChange(): void {
  currentPage.value = 1;
  void fetchRoles();
}

/**
 * Confirm selection and close drawer
 */
function handleConfirm(): void {
  emit('confirm', [...selectedRoles.value]);
  showDialog.value = false;
}

/**
 * Cancel handler
 */
function handleCancel(): void {
  emit('cancel');
  showDialog.value = false;
}

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    handleCancel();
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<template>
  <q-dialog v-model="showDialog" position="right" maximized>
    <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
      <!-- Header with Padding -->
      <q-card-section class="q-pb-sm">
        <div class="row items-center">
          <q-icon name="admin_panel_settings" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h6">Select Roles</div>
          <q-space />
          <q-btn icon="close" flat round dense class="rounded-borders" @click="handleCancel" />
        </div>
      </q-card-section>

      <!-- Info Banner -->
      <q-card-section class="q-pt-none q-pb-md">
        <q-banner dense class="bg-teal-1 text-teal-9 rounded-borders">
          <template #avatar>
            <q-icon name="info" color="teal-6" size="sm" />
          </template>
          <div class="text-caption">
            Select one or more roles to assign. Use filters to find roles and click to toggle selection.
          </div>
        </q-banner>
      </q-card-section>

      <!-- Selected Roles Preview -->
      <q-card-section v-if="selectedRoles.length > 0" class="q-pt-none q-pb-md">
        <div class="text-overline text-grey-7 q-mb-sm">
          <q-icon name="check_circle" size="xs" class="q-mr-xs" />
          Selected ({{ selectedRoles.length }})
          <q-btn
            flat
            dense
            size="xs"
            label="Clear all"
            color="red-6"
            class="q-ml-sm"
            @click="clearSelection"
          />
        </div>
        <div class="row q-gutter-xs">
          <SelectableChip
            v-for="role in selectedRoles"
            :key="role.id"
            :label="role.name"
            color="primary"
            size="sm"
            @remove="removeFromSelection(role)"
          />
        </div>
      </q-card-section>

      <!-- Filters -->
      <q-card-section class="q-py-md">
        <div class="text-overline text-grey-7 q-mb-md">
          <q-icon name="filter_list" size="xs" class="q-mr-xs" />
          Filters
        </div>
        <div class="row q-col-gutter-md">
          <!-- Search - Full width -->
          <div class="col-12">
            <q-input
              v-model="filters.name"
              outlined
              dense
              label="Search by name"
              placeholder="Type to search..."
              clearable
              class="rounded-borders"
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Type (System/Custom) and Scope (Template/Local) - Same line on desktop -->
          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.isSystem"
              outlined
              dense
              label="Type"
              class="rounded-borders"
              :options="typeOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="lock" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.isTemplate"
              outlined
              dense
              label="Scope"
              class="rounded-borders"
              :options="templateOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="share" />
              </template>
            </q-select>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="admin_panel_settings" size="xs" class="q-mr-xs" />
          Results
        </div>
      </q-card-section>

      <!-- Roles List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="loading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="text-grey-7 q-mt-md">Loading roles...</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="roles.length === 0" class="q-pa-md text-center">
          <q-icon name="inbox" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">No roles found</div>
        </div>

        <!-- Roles List with Infinite Scroll -->
        <q-scroll-area
          v-else
          ref="scrollAreaRef"
          style="height: 100%;"
          @scroll="onScroll"
        >
          <q-list separator>
            <q-item
              v-for="role in roles"
              :key="role.id || `role-${Math.random()}`"
              clickable
              :active="isSelected(role)"
              @click="toggleRoleSelection(role)"
            >
              <q-item-section avatar>
                <q-checkbox
                  :model-value="isSelected(role)"
                  color="primary"
                  @update:model-value="toggleRoleSelection(role)"
                  @click.stop
                />
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ role.name || 'Unnamed Role' }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  <span v-if="role.isTemplate">Template</span>
                  <span v-if="role.isSystem"> • System</span>
                  <span v-if="role.description"> • {{ role.description }}</span>
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-badge
                  v-if="role.isSystem"
                  color="blue-grey-6"
                  label="SYSTEM"
                />
                <q-badge
                  v-else-if="role.isTemplate"
                  color="teal-6"
                  label="TEMPLATE"
                />
                <q-badge
                  v-else
                  color="primary"
                  label="LOCAL"
                />
              </q-item-section>
            </q-item>
          </q-list>

          <!-- Load More Indicator -->
          <div v-if="loadingMore" class="q-pa-md text-center">
            <q-spinner color="primary" size="2em" />
          </div>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator />
      <q-card-actions class="row items-center q-px-md q-py-md">
        <div class="text-caption text-grey-7">
          <q-icon name="admin_panel_settings" size="xs" class="q-mr-xs" />
          {{ totalItems }} {{ totalItems === 1 ? 'role' : 'roles' }}
        </div>
        <q-space />
        <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders q-mr-sm" @click="handleCancel" />
        <q-btn
          unelevated
          dense
          label="Confirm Selection"
          color="primary"
          size="sm"
          class="rounded-borders"
          :disable="!canConfirm"
          @click="handleConfirm"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

/* Hover effects for list items */
:deep(.q-item) {
  transition: all var(--mapex-transition-base) ease;
}

:deep(.q-item:hover) {
  background-color: var(--mapex-surface-bg);
}

:deep(.q-item.q-item--active) {
  background-color: var(--mapex-active-bg) !important;
  border-left: 3px solid var(--q-primary);
}

/* Better spacing for filter inputs */
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

/* Smooth transitions */
:deep(.q-badge),
:deep(.q-chip) {
  transition: all var(--mapex-transition-base) ease;
}

/* Footer padding (ensure proper spacing) */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

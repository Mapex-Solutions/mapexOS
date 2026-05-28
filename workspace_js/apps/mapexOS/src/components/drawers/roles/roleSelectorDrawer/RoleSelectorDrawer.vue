<script setup lang="ts">
defineOptions({
  name: 'RoleSelectorDrawer'
});

/** TYPE IMPORTS */
import type { RoleSelectorDrawerProps, RoleSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPOSABLES */
import { useCommonPlaceholders } from '@composables/i18n';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** COMPOSABLES & STORES */
const { placeholders } = useCommonPlaceholders();

/** PROPS & EMITS */
const props = withDefaults(defineProps<RoleSelectorDrawerProps>(), {
  selectedRoleId: null,
});

const emit = defineEmits<RoleSelectorDrawerEmits>();

/** STATE */
const loading = ref(false);
const loadingMore = ref(false);
const roles = ref<any[]>([]);
const scrollAreaRef = ref<any>(null);

/** PAGINATION STATE */
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(15);

/** FILTER STATE */
const filters = ref({
  name: undefined as string | undefined,
  enabled: undefined as boolean | undefined,
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
 * Status filter options
 */
const statusOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false },
]);

/**
 * Template filter options
 */
const templateOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'Templates Only', value: true },
  { label: 'Non-Templates Only', value: false },
]);

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchRoles();
  }
});

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
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
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
 * Select role and close drawer
 * @param {any} role - Role to select
 */
function selectRole(role: any): void {
  emit('select', role);
  showDialog.value = false;
}

/**
 * Check if role is selected
 * @param {any} role - Role to check
 * @returns {boolean} True if selected
 */
function isSelected(role: any): boolean {
  return role.id === props.selectedRoleId;
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
          <div class="text-h6">Select Role</div>
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
            Use filters below to find the role you want to assign. Click on a role to select it.
          </div>
        </q-banner>
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
              :placeholder="placeholders.typeToSearch.value"
              clearable
              class="rounded-borders"
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Status and Template - Same line on desktop -->
          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.enabled"
              outlined
              dense
              label="Status"
              class="rounded-borders"
              :options="statusOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="toggle_on" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.isTemplate"
              outlined
              dense
              label="Type"
              class="rounded-borders"
              :options="templateOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="library_books" />
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
              @click="selectRole(role)"
            >
              <q-item-section avatar>
                <q-avatar
                  :color="role.enabled ? 'primary' : 'grey-5'"
                  icon="admin_panel_settings"
                  text-color="white"
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
                  :color="role.enabled ? 'green-6' : 'red-6'"
                  :label="role.enabled ? 'ACTIVE' : 'INACTIVE'"
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
        <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders" @click="handleCancel" />
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

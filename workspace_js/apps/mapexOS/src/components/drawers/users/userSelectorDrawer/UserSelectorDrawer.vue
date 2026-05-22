<script setup lang="ts">
defineOptions({
  name: 'UserSelectorDrawer'
});

/** TYPE IMPORTS */
import type { UserSelectorDrawerProps, UserSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** PROPS & EMITS */
const props = withDefaults(defineProps<UserSelectorDrawerProps>(), {
  selectedUserId: null,
});

const emit = defineEmits<UserSelectorDrawerEmits>();

/** STATE */
const loading = ref(false);
const loadingMore = ref(false);
const users = ref<any[]>([]);
const scrollAreaRef = ref<any>(null);

/** PAGINATION STATE */
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(15);

/** FILTER STATE */
const filters = ref({
  email: undefined as string | undefined,
  firstName: undefined as string | undefined,
  lastName: undefined as string | undefined,
  enabled: undefined as boolean | undefined,
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

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchUsers();
  }
});

/** FUNCTIONS */

/**
 * Fetch users from API with current filters and pagination
 * @param {boolean} append - If true, append to existing list (for infinite scroll)
 * @returns {Promise<void>}
 */
async function fetchUsers(append: boolean = false): Promise<void> {
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
    if (filters.value.email) {
      queryParams.email = filters.value.email;
    }
    if (filters.value.firstName) {
      queryParams.firstName = filters.value.firstName;
    }
    if (filters.value.lastName) {
      queryParams.lastName = filters.value.lastName;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }

    const response = await apis.mapexOS.users.list(queryParams);

    if (append) {
      users.value = [...users.value, ...(response.items || [])];
    } else {
      users.value = response.items || [];
    }

    totalPages.value = response.pagination?.totalPages || 1;
    totalItems.value = response.pagination?.totalItems || 0;
  } catch (error: any) {
    handleApiError({
      error,
      customMessage: 'Failed to fetch users'
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
      void fetchUsers(true);
    }
  }
}

/**
 * Select user and close drawer
 * @param {any} user - User to select
 */
function selectUser(user: any): void {
  emit('select', user);
  showDialog.value = false;
}

/**
 * Check if user is selected
 * @param {any} user - User to check
 * @returns {boolean} True if selected
 */
function isSelected(user: any): boolean {
  return user.id === props.selectedUserId;
}

/**
 * Filter change handler
 * Resets pagination and refetches users
 */
function onFilterChange(): void {
  currentPage.value = 1;
  void fetchUsers();
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

/**
 * Get user display name
 * @param {any} user - User object
 * @returns {string} Display name
 */
function getUserDisplayName(user: any): string {
  if (user.firstName || user.lastName) {
    return `${user.firstName || ''} ${user.lastName || ''}`.trim();
  }
  return user.email || 'Unnamed User';
}

/**
 * Get user initials for avatar
 * @param {any} user - User object
 * @returns {string} Initials
 */
function getUserInitials(user: any): string {
  if (user.firstName && user.lastName) {
    return `${user.firstName[0]}${user.lastName[0]}`.toUpperCase();
  }
  if (user.firstName) {
    return user.firstName.substring(0, 2).toUpperCase();
  }
  if (user.email) {
    return user.email.substring(0, 2).toUpperCase();
  }
  return '??';
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
          <q-icon name="person" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h6">Select User</div>
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
            Use filters below to find the user you want to select. Click on a user to select them.
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
          <!-- Email Search - Full width -->
          <div class="col-12">
            <q-input
              v-model="filters.email"
              outlined
              dense
              label="Search by email"
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

          <!-- First Name and Last Name - Same line -->
          <div class="col-12 col-sm-6">
            <q-input
              v-model="filters.firstName"
              outlined
              dense
              label="First name"
              clearable
              class="rounded-borders"
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="person" />
              </template>
            </q-input>
          </div>

          <div class="col-12 col-sm-6">
            <q-input
              v-model="filters.lastName"
              outlined
              dense
              label="Last name"
              clearable
              class="rounded-borders"
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="person" />
              </template>
            </q-input>
          </div>

          <!-- Status - Full width -->
          <div class="col-12">
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
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="person" size="xs" class="q-mr-xs" />
          Results
        </div>
      </q-card-section>

      <!-- Users List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="loading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="text-grey-7 q-mt-md">Loading users...</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="users.length === 0" class="q-pa-md text-center">
          <q-icon name="inbox" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">No users found</div>
        </div>

        <!-- Users List with Infinite Scroll -->
        <q-scroll-area
          v-else
          ref="scrollAreaRef"
          style="height: 100%;"
          @scroll="onScroll"
        >
          <q-list separator>
            <q-item
              v-for="user in users"
              :key="user.id || `user-${Math.random()}`"
              clickable
              :active="isSelected(user)"
              @click="selectUser(user)"
            >
              <q-item-section avatar>
                <q-avatar
                  :color="user.enabled ? 'primary' : 'grey-5'"
                  text-color="white"
                >
                  {{ getUserInitials(user) }}
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ getUserDisplayName(user) }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ user.email || 'No email' }}
                  <span v-if="user.jobTitle"> • {{ user.jobTitle }}</span>
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-badge
                  :color="user.enabled ? 'green-6' : 'red-6'"
                  :label="user.enabled ? 'ACTIVE' : 'INACTIVE'"
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
          <q-icon name="person" size="xs" class="q-mr-xs" />
          {{ totalItems }} {{ totalItems === 1 ? 'user' : 'users' }}
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

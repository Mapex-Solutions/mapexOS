<script setup lang="ts">
defineOptions({
  name: 'UserMultiSelectorDrawer',
});

/** TYPE IMPORTS */
import type {
  UserMultiSelectorDrawerProps,
  UserMultiSelectorDrawerEmits,
  UserSelectorItem,
  UserFilterMode,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPONENTS */
import { SelectableChip } from '@components/chips';

/** UTILS */
import { handleApiError } from '@utils/error';

/** LOCAL IMPORTS */
import { USER_SELECTOR_DEFAULTS, FILTER_MODE_OPTIONS } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<UserMultiSelectorDrawerProps>(), {
  excludeUserIds: () => [],
  selectedUserIds: () => [],
});

const emit = defineEmits<UserMultiSelectorDrawerEmits>();

/** STATE */
const loading = ref(false);
const loadingMore = ref(false);
const users = ref<UserSelectorItem[]>([]);
const selectedUsers = ref<UserSelectorItem[]>([]);
const scrollAreaRef = ref<any>(null);
const searchQuery = ref('');
const filterMode = ref<UserFilterMode>('name');
let debounceTimer: ReturnType<typeof setTimeout> | null = null;

/** PAGINATION STATE */
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(USER_SELECTOR_DEFAULTS.PER_PAGE);

/** COMPUTED */

/**
 * Drawer visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Excluded user IDs set for quick lookup
 */
const excludedIdsSet = computed(() => new Set(props.excludeUserIds));

/**
 * Selected user IDs set for quick lookup
 */
const selectedUserIdsSet = computed(() => new Set(selectedUsers.value.map(u => u.id)));

/**
 * Filter mode options for toggle
 */
const filterModeOptions = computed(() =>
  FILTER_MODE_OPTIONS.map(opt => ({
    value: opt.value,
    slot: opt.value,
  }))
);

/**
 * Check if confirm button should be enabled
 */
const canConfirm = computed(() => selectedUsers.value.length > 0);

/**
 * Filtered users excluding already selected users elsewhere
 */
const displayUsers = computed(() =>
  users.value.filter(u => !excludedIdsSet.value.has(u.id))
);

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    // Reset state
    selectedUsers.value = [];
    searchQuery.value = '';
    filterMode.value = 'name';
    currentPage.value = 1;
    void fetchUsers();
  }
});

/**
 * Watch selectedUserIds prop and sync with loaded users
 */
watch(() => [props.selectedUserIds, users.value], () => {
  if (props.selectedUserIds && props.selectedUserIds.length > 0 && users.value.length > 0) {
    const matchingUsers = users.value.filter(u => props.selectedUserIds?.includes(u.id));
    const existingIds = new Set(selectedUsers.value.map(u => u.id));
    matchingUsers.forEach(user => {
      if (!existingIds.has(user.id)) {
        selectedUsers.value.push(user);
      }
    });
  }
}, { deep: true });

/** FUNCTIONS */

/**
 * Fetch users from API with current filters and pagination
 *
 * @param {boolean} append - If true, append to existing list (for infinite scroll)
 * @returns {Promise<void>}
 */
async function fetchUsers(append: boolean = false): Promise<void> {
  if (!apis.mapexOS?.users) {
    return;
  }

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
      projection: 'id,firstName,lastName,email',
    };

    // Add filter based on mode
    if (searchQuery.value.trim()) {
      if (filterMode.value === 'name') {
        queryParams.firstName = searchQuery.value.trim();
      } else {
        queryParams.email = searchQuery.value.trim();
      }
    }

    const response = await apis.mapexOS.users.list(queryParams);

    const mappedUsers: UserSelectorItem[] = (response.items || []).map((user: any) => ({
      id: user.id,
      firstName: user.firstName || '',
      lastName: user.lastName || '',
      email: user.email || '',
    }));

    if (append) {
      users.value = [...users.value, ...mappedUsers];
    } else {
      users.value = mappedUsers;
    }

    totalPages.value = response.pagination?.totalPages || 1;
    totalItems.value = response.pagination?.totalItems || 0;
  } catch (error: any) {
    handleApiError({
      error,
      customMessage: 'Failed to fetch users',
    });
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
}

/**
 * Handle search input with debounce
 */
function onSearchInput(): void {
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }
  debounceTimer = setTimeout(() => {
    currentPage.value = 1;
    void fetchUsers();
  }, USER_SELECTOR_DEFAULTS.DEBOUNCE_MS);
}

/**
 * Handle filter mode change
 */
function onFilterModeChange(): void {
  currentPage.value = 1;
  void fetchUsers();
}

/**
 * Infinite scroll handler
 *
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
 * Toggle user selection
 *
 * @param {UserSelectorItem} user - User to toggle
 */
function toggleUserSelection(user: UserSelectorItem): void {
  const index = selectedUsers.value.findIndex(u => u.id === user.id);
  if (index >= 0) {
    selectedUsers.value.splice(index, 1);
  } else {
    selectedUsers.value.push(user);
  }
}

/**
 * Check if user is selected
 *
 * @param {UserSelectorItem} user - User to check
 * @returns {boolean} True if selected
 */
function isSelected(user: UserSelectorItem): boolean {
  return selectedUserIdsSet.value.has(user.id);
}

/**
 * Remove a user from selection
 *
 * @param {UserSelectorItem} user - User to remove
 */
function removeFromSelection(user: UserSelectorItem): void {
  const index = selectedUsers.value.findIndex(u => u.id === user.id);
  if (index >= 0) {
    selectedUsers.value.splice(index, 1);
  }
}

/**
 * Clear all selections
 */
function clearSelection(): void {
  selectedUsers.value = [];
}

/**
 * Get initials from user
 *
 * @param {UserSelectorItem} user - User object
 * @returns {string} User initials
 */
function getInitials(user: UserSelectorItem): string {
  const first = user.firstName?.charAt(0) || '';
  const last = user.lastName?.charAt(0) || '';
  if (first || last) {
    return (first + last).toUpperCase();
  }
  return user.email?.charAt(0).toUpperCase() || '?';
}

/**
 * Get display name for user
 *
 * @param {UserSelectorItem} user - User object
 * @returns {string} Display name
 */
function getDisplayName(user: UserSelectorItem): string {
  const name = `${user.firstName || ''} ${user.lastName || ''}`.trim();
  return name || user.email;
}

/**
 * Confirm selection and close drawer
 */
function handleConfirm(): void {
  emit('confirm', [...selectedUsers.value]);
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
 *
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
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }
});
</script>

<template>
  <q-dialog v-model="showDialog" position="right" maximized>
    <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
      <!-- Header -->
      <q-card-section class="q-pb-sm">
        <div class="row items-center">
          <q-icon name="group_add" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h6">Select Users</div>
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
            Select users to add as members. Use the toggle to switch between searching by name or email.
          </div>
        </q-banner>
      </q-card-section>

      <!-- Selected Users Preview -->
      <q-card-section v-if="selectedUsers.length > 0" class="q-pt-none q-pb-md">
        <div class="text-overline text-grey-7 q-mb-sm">
          <q-icon name="check_circle" size="xs" class="q-mr-xs" />
          Selected ({{ selectedUsers.length }})
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
            v-for="user in selectedUsers"
            :key="user.id"
            :label="getDisplayName(user)"
            color="primary"
            size="sm"
            @remove="removeFromSelection(user)"
          />
        </div>
      </q-card-section>

      <!-- Search Filter with Toggle -->
      <q-card-section class="q-py-md">
        <div class="text-overline text-grey-7 q-mb-md">
          <q-icon name="search" size="xs" class="q-mr-xs" />
          Search
        </div>
        <div class="row q-col-gutter-md items-center">
          <div class="col">
            <q-input
              v-model="searchQuery"
              outlined
              dense
              :placeholder="filterMode === 'name' ? 'Search by name...' : 'Search by email...'"
              clearable
              class="rounded-borders"
              @update:model-value="onSearchInput"
            >
              <template #prepend>
                <q-icon :name="filterMode === 'name' ? 'person' : 'email'" />
              </template>
            </q-input>
          </div>
          <div class="col-auto">
            <q-btn-toggle
              v-model="filterMode"
              toggle-color="primary"
              :options="filterModeOptions"
              rounded
              dense
              unelevated
              class="rounded-borders"
              @update:model-value="onFilterModeChange"
            >
              <template #name>
                <q-icon name="person" size="xs" class="q-mr-xs" />
                Name
              </template>
              <template #email>
                <q-icon name="email" size="xs" class="q-mr-xs" />
                Email
              </template>
            </q-btn-toggle>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="people" size="xs" class="q-mr-xs" />
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
        <div v-else-if="displayUsers.length === 0" class="q-pa-md text-center">
          <q-icon name="person_off" size="4em" color="grey-5" />
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
              v-for="user in displayUsers"
              :key="user.id"
              clickable
              :active="isSelected(user)"
              @click="toggleUserSelection(user)"
            >
              <q-item-section avatar>
                <q-checkbox
                  :model-value="isSelected(user)"
                  color="primary"
                  @update:model-value="toggleUserSelection(user)"
                  @click.stop
                />
              </q-item-section>

              <q-item-section avatar>
                <q-avatar
                  :color="isSelected(user) ? 'primary' : 'grey-4'"
                  :text-color="isSelected(user) ? 'white' : 'grey-8'"
                  size="md"
                >
                  {{ getInitials(user) }}
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ getDisplayName(user) }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ user.email }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-icon
                  v-if="isSelected(user)"
                  name="check_circle"
                  color="primary"
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
          <q-icon name="people" size="xs" class="q-mr-xs" />
          {{ totalItems }} {{ totalItems === 1 ? 'user' : 'users' }}
        </div>
        <q-space />
        <q-btn
          flat
          dense
          label="Cancel"
          color="grey-7"
          size="sm"
          class="rounded-borders q-mr-sm"
          @click="handleCancel"
        />
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

/* Footer padding */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

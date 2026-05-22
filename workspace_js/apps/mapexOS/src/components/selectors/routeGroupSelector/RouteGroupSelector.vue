<script setup lang="ts">
defineOptions({
  name: 'RouteGroupSelector'
});

/** TYPE IMPORTS */
import type { RouteGroupResponse } from '@mapexos/schemas';
import type { RouteGroupSelectorProps, RouteGroupSelectorEmits } from './interfaces/routeGroupSelector.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { SelectableChip } from '@components/chips';

/** UTILS */
import { handleApiError } from '@utils/error';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useRouteGroupSelectorTranslations } from '@src/composables/i18n/components/selectors/useRouteGroupSelectorTranslations';

/** PROPS & EMITS */
const props = defineProps<RouteGroupSelectorProps>();
const emit = defineEmits<RouteGroupSelectorEmits>();

const t = useRouteGroupSelectorTranslations();

// Filters
const filters = ref({
  search: '',
  isSystem: null as boolean | null,
  isTemplate: null as boolean | null,
});

// Data
const routeGroups = ref<RouteGroupResponse[]>([]);
const selectedRouteGroups = ref<RouteGroupResponse[]>([]);

// Pagination
const currentPage = ref(1);
const totalPages = ref(1);
const totalRouteGroups = ref(0);
const perPage = 30;

// Loading states
const isLoading = ref(false);
const isLoadingMore = ref(false);

// Scroll ref
const scrollAreaRef = ref<any>(null);

// Template Type filter options
const templateTypeOptions = computed(() => [
  { label: t.filters.templateTypeOptions.all.value, value: null },
  { label: t.filters.templateTypeOptions.system.value, value: true },
  { label: t.filters.templateTypeOptions.custom.value, value: false },
]);

// Template Source filter options
const templateSourceOptions = computed(() => [
  { label: t.filters.templateSourceOptions.all.value, value: null },
  { label: t.filters.templateSourceOptions.shared.value, value: true },
  { label: t.filters.templateSourceOptions.local.value, value: false },
]);

// Load route groups with current filters
async function loadRouteGroups(append = false) {
  if (append) {
    isLoadingMore.value = true;
  } else {
    isLoading.value = true;
    routeGroups.value = [];
  }

  try {
    if (!apis.router) {
      throw new Error('Router API not configured');
    }

    const response = await apis.router.routegroup.list({
      page: currentPage.value,
      perPage,
      ...(filters.value.search && { name: filters.value.search }),
      ...(typeof filters.value.isSystem === 'boolean' && { isSystem: filters.value.isSystem }),
      ...(typeof filters.value.isTemplate === 'boolean' && { isTemplate: filters.value.isTemplate }),
      sort: 'name:asc',
    });

    if (append) {
      routeGroups.value.push(...response.items);
    } else {
      routeGroups.value = response.items;
    }

    totalRouteGroups.value = response.pagination.totalItems;
    totalPages.value = response.pagination.totalPages;
  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.errors.loadFailed.value,
      timeout: 5000
    });
  } finally {
    isLoading.value = false;
    isLoadingMore.value = false;
  }
}

// Infinite scroll handler
function onScroll(info: any) {
  const scrollPosition = info.verticalPosition;
  const scrollSize = info.verticalSize;
  const containerSize = info.verticalContainerSize;

  // Load more when scrolled to 80%
  if (scrollPosition + containerSize >= scrollSize * 0.8) {
    if (!isLoadingMore.value && currentPage.value < totalPages.value) {
      currentPage.value++;
      void loadRouteGroups(true);
    }
  }
}

// Toggle route group selection
function toggleRouteGroup(routeGroup: RouteGroupResponse) {
  const index = selectedRouteGroups.value.findIndex(rg => rg.id === routeGroup.id);

  if (index === -1) {
    // Add to selection
    selectedRouteGroups.value.push(routeGroup);
  } else {
    // Remove from selection
    selectedRouteGroups.value.splice(index, 1);
  }

  // Emit updated IDs
  const ids = selectedRouteGroups.value.map(rg => rg.id || '').filter(id => id);
  emit('update:modelValue', ids);
  emit('update:selectedRouteGroups', selectedRouteGroups.value);
}

// Check if route group is selected
function isSelected(routeGroup: RouteGroupResponse): boolean {
  return selectedRouteGroups.value.some(rg => rg.id === routeGroup.id);
}

// Filter change handler
function onFilterChange() {
  currentPage.value = 1;
  void loadRouteGroups();
}

// Search watcher (debounce handled by q-input debounce prop)
watch(() => filters.value.search, () => {
  currentPage.value = 1;
  void loadRouteGroups();
});

// Watch for external model value changes
watch(() => props.modelValue, (newValue) => {
  if (!newValue || newValue.length === 0) {
    selectedRouteGroups.value = [];
  } else {
    // Find route groups in current list
    const selected = routeGroups.value.filter(rg => newValue.includes(rg.id || ''));
    selectedRouteGroups.value = selected;
  }
}, { immediate: true });

// Initial load
onMounted(() => {
  void loadRouteGroups();
});
</script>

<template>
  <div class="route-group-selector">
    <!-- Filters Row: Search + Type + Source -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-sm-4">
        <q-input
          v-model="filters.search"
          outlined
          dense
          debounce="1000"
          :placeholder="t.search.placeholder.value"
        >
          <template v-slot:prepend>
            <q-icon name="search" color="primary" />
          </template>
          <template v-slot:append>
            <q-icon
              v-if="filters.search"
              name="close"
              @click="filters.search = ''"
              class="cursor-pointer"
            />
          </template>
        </q-input>
      </div>
      <div class="col-12 col-sm-4">
        <q-select
          v-model="filters.isSystem"
          outlined
          dense
          clearable
          emit-value
          map-options
          :label="t.filters.templateType.value"
          :options="templateTypeOptions"
          @update:model-value="onFilterChange"
        >
          <template v-slot:prepend>
            <q-icon name="build" color="primary" />
          </template>
        </q-select>
      </div>
      <div class="col-12 col-sm-4">
        <q-select
          v-model="filters.isTemplate"
          outlined
          dense
          clearable
          emit-value
          map-options
          :label="t.filters.templateSource.value"
          :options="templateSourceOptions"
          @update:model-value="onFilterChange"
        >
          <template v-slot:prepend>
            <q-icon name="content_copy" color="primary" />
          </template>
        </q-select>
      </div>
    </div>

    <!-- Results Count -->
    <div class="text-caption text-grey-7 q-mb-sm">
      <q-icon name="info" size="xs" class="q-mr-xs" />
      {{ t.results.found.value.replace('{count}', String(totalRouteGroups)) }}
    </div>

    <!-- Route Group List with Virtual Scroll -->
    <q-scroll-area
      :style="{ height: '450px', border: '1px solid var(--mapex-card-border)', borderRadius: 'var(--mapex-radius-md)' }"
      @scroll="onScroll"
      ref="scrollAreaRef"
      class="bg-grey-1"
    >
      <div class="q-pa-sm">
        <!-- Loading Skeleton -->
        <div v-if="isLoading">
          <q-skeleton
            v-for="i in 5"
            :key="i"
            height="100px"
            class="q-mb-sm rounded-borders"
          />
        </div>

        <!-- Route Group Cards -->
        <q-card
          v-for="routeGroup in routeGroups"
          :key="routeGroup.id || ''"
          flat
          bordered
          clickable
          class="route-group-card q-mb-sm"
          :class="{ 'selected': isSelected(routeGroup) }"
          @click="toggleRouteGroup(routeGroup)"
        >
          <q-card-section class="row items-start q-pa-md">
            <!-- Selection Checkbox -->
            <div class="col-auto q-mr-md q-pt-xs">
              <q-checkbox
                :model-value="isSelected(routeGroup)"
                color="primary"
                size="md"
              />
            </div>

            <!-- Route Group Info -->
            <div class="col" style="min-width: 0;">
              <!-- Line 1: Name and Status Badge -->
              <div class="row items-center q-mb-xs no-wrap">
                <div class="text-subtitle1 text-weight-medium ellipsis" style="flex: 1; min-width: 0;">
                  {{ routeGroup.name }}
                </div>
                <q-badge
                  :color="routeGroup.enabled ? 'positive' : 'grey'"
                  :label="routeGroup.enabled ? t.badge.active.value : t.badge.inactive.value"
                  class="q-ml-sm"
                  style="flex-shrink: 0;"
                />
              </div>

              <!-- Line 2: Description -->
              <div class="text-body2 text-grey-7 q-mb-xs ellipsis-2-lines">
                {{ routeGroup.description || t.noDescription.value }}
              </div>

              <!-- Line 3: Routers count -->
              <div class="text-caption text-grey-8 row items-center no-wrap">
                <q-icon name="route" size="xs" class="q-mr-xs" />
                <span>{{ (routeGroup.routers?.length || 0) }} {{ t.routers.value }}</span>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Loading More Indicator -->
        <div v-if="isLoadingMore" class="text-center q-pa-md">
          <q-spinner color="primary" size="md" />
          <div class="text-caption text-grey-7 q-mt-sm">
            {{ t.loading.more.value }}
          </div>
        </div>

        <!-- Empty State -->
        <div
          v-if="!isLoading && routeGroups.length === 0"
          class="text-center q-pa-xl"
        >
          <q-icon name="inbox" size="64px" color="grey-5" />
          <div class="text-h6 text-grey-7 q-mt-md">
            {{ t.emptyState.title.value }}
          </div>
          <div class="text-caption text-grey-6">
            {{ t.emptyState.subtitle.value }}
          </div>
        </div>
      </div>
    </q-scroll-area>

    <!-- Selected Route Groups Preview -->
    <div v-if="selectedRouteGroups.length > 0" class="q-mt-md">
      <q-card flat bordered class="bg-blue-1">
        <q-card-section class="q-pa-md">
          <div class="text-caption text-weight-medium text-primary q-mb-xs">
            <q-icon name="check_circle" size="xs" class="q-mr-xs" />
            {{ t.selected.label.value }} ({{ selectedRouteGroups.length }})
          </div>
          <div class="row q-gutter-xs">
            <SelectableChip
              v-for="(routeGroup, index) in selectedRouteGroups"
              :key="routeGroup.id ?? `route-group-${index}`"
              :label="routeGroup.name"
              color="primary"
              size="sm"
              @remove="toggleRouteGroup(routeGroup)"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.route-group-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);
  background-color: white;
}

.route-group-card:hover {
  border-color: var(--q-primary);
  box-shadow: var(--mapex-shadow-sm);
}

.route-group-card.selected {
  border-color: var(--q-primary);
  border-width: 2px;
  background-color: rgba(59, 109, 94, 0.05);
}

.cursor-pointer {
  cursor: pointer;
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

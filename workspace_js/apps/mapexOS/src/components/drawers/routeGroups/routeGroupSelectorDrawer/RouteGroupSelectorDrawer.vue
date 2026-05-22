<script setup lang="ts">
defineOptions({
  name: 'RouteGroupSelectorDrawer'
});

/** TYPE IMPORTS */
import type { RouteGroupResponse } from '@mapexos/schemas';
import type { RouteGroupSelectorDrawerProps, RouteGroupSelectorDrawerEmits } from './interfaces/routeGroupSelectorDrawer.interface';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useRouters } from '@composables/routing/routers';
import { useRouteGroupSelectorTranslations } from '@src/composables/i18n/components/selectors/useRouteGroupSelectorTranslations';

/** PROPS & EMITS */
const props = withDefaults(defineProps<RouteGroupSelectorDrawerProps>(), {
  selectedRouteGroupIds: () => [],
  multiSelect: true,
});

const emit = defineEmits<RouteGroupSelectorDrawerEmits>();

const t = useRouteGroupSelectorTranslations();

/** COMPOSABLES & STORES */
const {
  routeGroups,
  isLoading,
  isLoadingMore,
  filters,
  pagination,
  fetchRouteGroups,
} = useRouters();

/** STATE */
const selectedRouteGroups = ref<RouteGroupResponse[]>([]);
const scrollAreaRef = ref<any>(null);

/** COMPUTED */

/**
 * Template Type filter options
 */
const templateTypeOptions = computed(() => [
  { label: t.filters.templateTypeOptions.all.value, value: undefined },
  { label: t.filters.templateTypeOptions.system.value, value: true },
  { label: t.filters.templateTypeOptions.custom.value, value: false },
]);

/**
 * Template Source filter options
 */
const templateSourceOptions = computed(() => [
  { label: t.filters.templateSourceOptions.all.value, value: undefined },
  { label: t.filters.templateSourceOptions.shared.value, value: true },
  { label: t.filters.templateSourceOptions.local.value, value: false },
]);

/**
 * Dialog visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Total route groups count from pagination
 */
const totalRouteGroups = computed(() => pagination.value.totalItems);

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    filters.value.kinds = props.allowedRouterKinds;
    void fetchRouteGroups();
  }
});

/** FUNCTIONS */

/**
 * Infinite scroll handler
 * @param {any} info - Scroll information
 */
function onScroll(info: any): void {
  const scrollPosition = info.verticalPosition;
  const scrollSize = info.verticalSize;
  const containerSize = info.verticalContainerSize;

  if (scrollPosition + containerSize >= scrollSize * 0.8) {
    if (!isLoadingMore.value && pagination.value.currentPage < pagination.value.totalPages) {
      pagination.value.currentPage++;
      void fetchRouteGroups(true);
    }
  }
}

/**
 * Toggle route group selection
 * @param {RouteGroupResponse} routeGroup - Route group to toggle
 */
function toggleRouteGroup(routeGroup: RouteGroupResponse): void {
  if (!props.multiSelect) {
    selectedRouteGroups.value = [routeGroup];
    emit('select', selectedRouteGroups.value);
    showDialog.value = false;
    return;
  }

  const index = selectedRouteGroups.value.findIndex(rg => rg.id === routeGroup.id);

  if (index >= 0) {
    selectedRouteGroups.value.splice(index, 1);
  } else {
    selectedRouteGroups.value.push(routeGroup);
  }
}

/**
 * Check if route group is selected
 * @param {RouteGroupResponse} routeGroup - Route group to check
 * @returns {boolean} True if selected
 */
function isSelected(routeGroup: RouteGroupResponse): boolean {
  return selectedRouteGroups.value.some(rg => rg.id === routeGroup.id);
}

/**
 * Confirm selection (multi-select mode)
 */
function confirmSelection(): void {
  emit('select', selectedRouteGroups.value);
  showDialog.value = false;
}

/**
 * Filter change handler
 * Resets pagination and refetches route groups
 */
function onFilterChange(): void {
  pagination.value.currentPage = 1;
  void fetchRouteGroups();
}

/**
 * Cancel handler
 */
function handleCancel(): void {
  emit('cancel');
  showDialog.value = false;
}
</script>

<template>
  <q-dialog v-model="showDialog" position="right" maximized>
    <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
      <!-- Header with Padding -->
      <q-card-section class="q-pb-sm">
        <div class="row items-center">
          <div class="text-h6">{{ multiSelect ? t.title.plural.value : t.title.singular.value }}</div>
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
            {{ multiSelect ? 'Select one or more route groups. Click Confirm when done.' : 'Click on a route group to select it.' }}
          </div>
        </q-banner>
      </q-card-section>

      <!-- Filters -->
      <q-card-section class="q-py-md">
        <div class="text-overline text-grey-7 q-mb-md">
          <q-icon name="filter_list" size="xs" class="q-mr-xs" />
          {{ t.filters.title?.value || 'Filters' }}
        </div>
        <div class="row q-col-gutter-md">
          <!-- Search - Full width -->
          <div class="col-12">
            <q-input
              v-model="filters.name"
              outlined
              dense
              :label="t.filters.search.label.value"
              :placeholder="t.filters.search.placeholder.value"
              clearable
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Template Type and Template Source - Same line on desktop -->
          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.isSystem"
              outlined
              dense
              :label="t.filters.templateType.value"
              :options="templateTypeOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="widgets" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.isTemplate"
              outlined
              dense
              :label="t.filters.templateSource.value"
              :options="templateSourceOptions"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="source" />
              </template>
            </q-select>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="alt_route" size="xs" class="q-mr-xs" />
          {{ t.resultsTitle?.value || 'Results' }}
        </div>
      </q-card-section>

      <!-- Route Groups List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="isLoading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="text-grey-7 q-mt-md">{{ t.loading.value }}</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="routeGroups.length === 0" class="q-pa-md text-center">
          <q-icon name="inbox" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">{{ t.empty.value }}</div>
        </div>

        <!-- Route Groups List with Infinite Scroll -->
        <q-scroll-area
          v-else
          ref="scrollAreaRef"
          style="height: 100%;"
          @scroll="onScroll"
        >
          <q-list separator>
            <q-item
              v-for="routeGroup in routeGroups"
              :key="routeGroup.id || `rg-${Math.random()}`"
              clickable
              :active="isSelected(routeGroup)"
              @click="toggleRouteGroup(routeGroup)"
            >
              <q-item-section v-if="multiSelect" avatar>
                <q-checkbox
                  :model-value="isSelected(routeGroup)"
                  color="primary"
                  @click.stop="toggleRouteGroup(routeGroup)"
                />
              </q-item-section>

              <q-item-section avatar v-else>
                <q-avatar color="teal-6" icon="alt_route" text-color="white" />
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ routeGroup.name }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ routeGroup.description || 'No description' }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-icon
                  v-if="routeGroup.description"
                  name="info"
                  color="primary"
                  size="sm"
                  class="cursor-pointer"
                >
                  <AppTooltip :content="routeGroup.description" />
                </q-icon>
              </q-item-section>
            </q-item>
          </q-list>

          <!-- Load More Indicator -->
          <div v-if="isLoadingMore" class="q-pa-md text-center">
            <q-spinner color="primary" size="2em" />
          </div>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator />
      <q-card-actions class="row items-center q-px-md q-py-md">
        <div class="text-caption text-grey-7">
          <q-icon name="alt_route" size="xs" class="q-mr-xs" />
          {{ totalRouteGroups }} {{ totalRouteGroups === 1 ? 'route group' : 'route groups' }}
        </div>
        <q-space />
        <template v-if="multiSelect">
          <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders" @click="handleCancel" />
          <q-btn
            flat
            dense
            label="Confirm Selection"
            color="primary"
            size="sm"
            class="rounded-borders"
            :disable="selectedRouteGroups.length === 0"
            @click="confirmSelection"
          />
        </template>
        <q-btn v-else flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders" @click="handleCancel" />
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

/* Checkbox styling */
:deep(.q-checkbox__inner) {
  border-radius: var(--mapex-radius-xs);
}

/* Better chip styling */
:deep(.q-chip) {
  font-weight: 500;
}

/* Footer padding (ensure proper spacing) */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'AssetSelectorDrawer'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { AssetResponse } from '@mapexos/schemas';
import type { AssetSelectorDrawerProps, AssetSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAssets } from '@composables/assets/assets';
import { useCommonPlaceholders } from '@composables/i18n';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AssetSelectorDrawerProps>(), {
  selectedAssetId: null,
  multiSelect: false,
});

const emit = defineEmits<AssetSelectorDrawerEmits>();

/** COMPOSABLES & STORES */
const {
  assets,
  isLoading,
  isLoadingMore,
  filters,
  pagination,
  categoryOptions: categoryOptionsData,
  manufacturerOptions: manufacturerOptionsData,
  modelOptions: modelOptionsData,
  loadingCategories,
  loadingManufacturers,
  loadingModels,
  fetchAssets,
  loadCategories,
  handleCategoryChange,
  handleManufacturerChange,
} = useAssets();
const { placeholders } = useCommonPlaceholders();

/** STATE */
const selectedAsset = ref<AssetResponse | null>(null);
const scrollAreaRef = ref<any>(null);

/** COMPUTED */

/**
 * Status filter options
 */
const statusOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false },
]);

/**
 * Drawer visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Category options for dropdown (alias from composable)
 */
const categoryOptions = categoryOptionsData;

/**
 * Manufacturer options for dropdown (alias from composable)
 */
const manufacturerOptions = manufacturerOptionsData;

/**
 * Model options for dropdown (alias from composable)
 */
const modelOptions = modelOptionsData;

/**
 * Total assets count from pagination
 */
const totalAssets = computed(() => pagination.value.totalItems);

/** WATCHERS */

/**
 * Watch category changes and cascade reset
 */
watch(() => filters.value.categoryId, async (newVal) => {
  await handleCategoryChange(newVal);
});

/**
 * Watch manufacturer changes and cascade reset
 */
watch(() => filters.value.manufacturerId, async (newVal) => {
  await handleManufacturerChange(newVal);
});

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    // Apply assetTemplateId filter if provided (pre-filter by template type)
    filters.value.assetTemplateId = props.assetTemplateId || undefined;
    void fetchAssets();
    void loadCategories();
  } else {
    // Reset template filter on close
    filters.value.assetTemplateId = undefined;
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
      void fetchAssets(true);
    }
  }
}

/**
 * Select asset
 * @param {AssetResponse} asset - Asset to select
 */
function selectAsset(asset: AssetResponse): void {
  selectedAsset.value = asset;
  emit('select', asset);
  showDialog.value = false;
}

/**
 * Filter change handler
 * Resets pagination and refetches assets
 */
function onFilterChange(): void {
  pagination.value.currentPage = 1;
  void fetchAssets();
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
          <div class="text-h6">Select Asset</div>
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
            Use filters below to narrow down your search. Click on an asset to select it.
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
              class="rounded-borders"
              label="Search by name"
              :placeholder="placeholders.typeToSearch.value"
              clearable
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Category and Status - 2 columns on tablet+ -->
          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.categoryId"
              outlined
              dense
              clearable
              class="rounded-borders"
              label="Category"
              :options="categoryOptions"
              :loading="loadingCategories"
              option-label="label"
              option-value="value"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="category" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.status"
              outlined
              dense
              class="rounded-borders"
              label="Status"
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

          <!-- Manufacturer and Model - 2 columns on tablet+ -->
          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.manufacturerId"
              outlined
              dense
              clearable
              class="rounded-borders"
              label="Manufacturer"
              :options="manufacturerOptions"
              :loading="loadingManufacturers"
              :disable="!filters.categoryId"
              option-label="label"
              option-value="value"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="factory" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.modelId"
              outlined
              dense
              clearable
              class="rounded-borders"
              label="Model"
              :options="modelOptions"
              :loading="loadingModels"
              :disable="!filters.manufacturerId"
              option-label="label"
              option-value="value"
              emit-value
              map-options
              @update:model-value="onFilterChange"
            >
              <template #prepend>
                <q-icon name="precision_manufacturing" />
              </template>
            </q-select>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="inventory_2" size="xs" class="q-mr-xs" />
          Results
        </div>
      </q-card-section>

      <!-- Assets List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="isLoading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="text-grey-7 q-mt-md">Loading assets...</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="assets.length === 0" class="q-pa-md text-center">
          <q-icon name="inbox" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">No assets found</div>
        </div>

        <!-- Assets List with Infinite Scroll -->
        <q-scroll-area
          v-else
          ref="scrollAreaRef"
          style="height: 100%;"
          @scroll="onScroll"
        >
          <q-list separator>
            <q-item
              v-for="asset in assets"
              :key="asset.id || `asset-${Math.random()}`"
              clickable
              :active="selectedAsset?.id === asset.id"
              @click="selectAsset(asset)"
            >
              <q-item-section avatar>
                <q-avatar :color="asset.enabled ? 'positive' : 'grey-5'" icon="sensors" text-color="white" />
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ asset.name }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ asset.assetUUID || 'No UUID' }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-icon
                  v-if="asset.assetUUID"
                  name="info"
                  color="primary"
                  size="sm"
                  class="cursor-pointer"
                >
                  <AppTooltip :content="`UUID: ${asset.assetUUID}`" />
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
          <q-icon name="inventory_2" size="xs" class="q-mr-xs" />
          {{ totalAssets }} {{ totalAssets === 1 ? 'asset' : 'assets' }}
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

/* Better chip styling */
:deep(.q-chip) {
  font-weight: 500;
}

/* Footer padding (ensure proper spacing) */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'AssetSelector'
});

/** TYPE IMPORTS */
import type { AssetResponse } from '@mapexos/schemas';
import type { AssetSelectorProps, AssetSelectorEmits } from './interfaces/assetSelector.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** UTILS */
import { handleApiError } from '@utils/error';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';

const logger = useLogger('AssetSelector');

/** PROPS & EMITS */
const props = withDefaults(defineProps<AssetSelectorProps>(), {
  label: 'Select Asset',
  required: false,
});

const emit = defineEmits<AssetSelectorEmits>();

// Filters
const filters = ref({
  search: '',
  status: null as boolean | null,
  categoryId: undefined as string | undefined,
  manufacturerId: undefined as string | undefined,
  model: undefined as string | undefined,
  assetUUID: undefined as string | undefined,
});

// Dynamic filter options
const categoryOptions = ref<Array<{ label: string; value: string }>>([]);
const manufacturerOptions = ref<Array<{ label: string; value: string }>>([]);
const modelOptions = ref<Array<{ label: string; value: string }>>([]);

const loadingCategories = ref(false);
const loadingManufacturers = ref(false);
const loadingModels = ref(false);

// Data
const assets = ref<AssetResponse[]>([]);
const selectedAsset = ref<AssetResponse | null>(null);

// Pagination
const currentPage = ref(1);
const totalPages = ref(1);
const totalAssets = ref(0);
const perPage = 30;

// Loading states
const isLoading = ref(false);
const isLoadingMore = ref(false);

// Dialog state
const showDialog = ref(false);

// Scroll ref
const scrollAreaRef = ref<any>(null);

// Status filter options
const statusOptions = computed(() => [
  { label: 'All', value: null },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false },
]);

// Load assets with current filters
async function loadAssets(append = false) {
  if (append) {
    isLoadingMore.value = true;
  } else {
    isLoading.value = true;
    assets.value = [];
  }

  try {
    if (!apis.assets) {
      throw new Error('Assets API not configured');
    }

    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage,
      sort: 'name:asc',
    };

    // Add active filters
    if (filters.value.search) queryParams.name = filters.value.search;
    if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;
    if (filters.value.categoryId) queryParams.categoryId = filters.value.categoryId;
    if (filters.value.manufacturerId) queryParams.manufacturerId = filters.value.manufacturerId;
    if (filters.value.model) queryParams.modelId = filters.value.model;

    const response = await apis.assets.asset.list(queryParams);

    if (append) {
      assets.value.push(...response.items);
    } else {
      assets.value = response.items;
    }

    totalAssets.value = response.pagination.totalItems;
    totalPages.value = response.pagination.totalPages;
  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: 'Failed to load assets',
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
      void loadAssets(true);
    }
  }
}

// Select asset
async function selectAsset(asset: AssetResponse) {
  // Fetch complete asset details to get assetIdPath (list API may not include it)
  try {
    if (apis.assets && asset.id) {
      const fullAsset = await apis.assets.asset.getById({ assetId: asset.id });
      selectedAsset.value = fullAsset;
      emit('update:modelValue', fullAsset.id || null);
      emit('update:selectedAsset', fullAsset);
    } else {
      selectedAsset.value = asset;
      emit('update:modelValue', asset.id || null);
      emit('update:selectedAsset', asset);
    }
  } catch (error) {
    // Fallback to list asset if getById fails
    logger.warn('Failed to fetch complete asset, using list data:', error);
    selectedAsset.value = asset;
    emit('update:modelValue', asset.id || null);
    emit('update:selectedAsset', asset);
  }
  showDialog.value = false;
}

// Clear selection
function clearSelection() {
  selectedAsset.value = null;
  emit('update:modelValue', null);
  emit('update:selectedAsset', null);
}

// Filter change handler
function onFilterChange() {
  currentPage.value = 1;
  void loadAssets();
}

// Load categories from API
async function loadCategories() {
  if (!apis.mapexOS?.lists) return;

  try {
    loadingCategories.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_category',
      page: 1,
      perPage: 1000,
    });

    categoryOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (error: any) {
    logger.error('Error loading categories:', error);
  } finally {
    loadingCategories.value = false;
  }
}

// Load manufacturers from API (requires categoryId)
async function loadManufacturers() {
  if (!apis.mapexOS?.lists || !filters.value.categoryId) {
    manufacturerOptions.value = [];
    return;
  }

  try {
    loadingManufacturers.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_manufacturer',
      parentId: filters.value.categoryId,
      page: 1,
      perPage: 1000,
    });

    manufacturerOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (error: any) {
    logger.error('Error loading manufacturers:', error);
  } finally {
    loadingManufacturers.value = false;
  }
}

// Load models from API (requires categoryId and manufacturerId)
async function loadModels() {
  if (!apis.mapexOS?.lists || !filters.value.categoryId || !filters.value.manufacturerId) {
    modelOptions.value = [];
    return;
  }

  try {
    loadingModels.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_model',
      parentId: filters.value.manufacturerId,
      page: 1,
      perPage: 1000,
    });

    modelOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.name, // Model uses name as value, not id
    }));
  } catch (error: any) {
    logger.error('Error loading models:', error);
  } finally {
    loadingModels.value = false;
  }
}

// Watch category changes to load manufacturers
watch(() => filters.value.categoryId, async (newVal) => {
  if (newVal) {
    filters.value.manufacturerId = undefined;
    filters.value.model = undefined;
    await loadManufacturers();
  } else {
    filters.value.manufacturerId = undefined;
    filters.value.model = undefined;
    manufacturerOptions.value = [];
    modelOptions.value = [];
  }
});

// Watch manufacturer changes to load models
watch(() => filters.value.manufacturerId, async (newVal) => {
  if (newVal && filters.value.categoryId) {
    filters.value.model = undefined;
    await loadModels();
  } else {
    filters.value.model = undefined;
    modelOptions.value = [];
  }
});

// Open dialog and load assets
function openDialog() {
  showDialog.value = true;
  if (assets.value.length === 0) {
    void loadAssets();
  }
  if (categoryOptions.value.length === 0) {
    void loadCategories();
  }
}

// Watch modelValue changes from parent
watch(() => props.modelValue, async (newValue) => {
  if (newValue && (!selectedAsset.value || selectedAsset.value.id !== newValue)) {
    // Load the selected asset details
    try {
      if (!apis.assets) return;

      const asset = await apis.assets.asset.getById({ assetId: newValue });
      selectedAsset.value = asset;
      // Emit selectedAsset so parent can access full asset data (e.g., assetIdPath)
      emit('update:selectedAsset', asset);
    } catch (error: any) {
      logger.error('Error loading asset:', error);
    }
  } else if (!newValue) {
    selectedAsset.value = null;
    emit('update:selectedAsset', null);
  }
});

// Load initial asset if modelValue is provided
onMounted(async () => {
  if (props.modelValue) {
    try {
      if (!apis.assets) return;

      const asset = await apis.assets.asset.getById({ assetId: props.modelValue });
      selectedAsset.value = asset;
      // Emit selectedAsset so parent can access full asset data (e.g., assetIdPath)
      emit('update:selectedAsset', asset);
    } catch (error: any) {
      logger.error('Error loading initial asset:', error);
    }
  }
});
</script>

<template>
  <div>
    <!-- Selected Asset Display -->
    <q-input
      :model-value="selectedAsset?.name || ''"
      outlined
      readonly
      :label="label"
      :rules="required ? [(val: any) => !!val || 'Asset is required'] : []"
      class="rounded-borders cursor-pointer"
      @click="openDialog"
    >
      <template #prepend>
        <q-icon name="sensors" color="primary" />
      </template>
      <template #append>
        <q-icon
          v-if="selectedAsset"
          name="close"
          class="cursor-pointer"
          @click.stop="clearSelection"
        />
        <q-icon
          name="search"
          class="cursor-pointer"
          @click.stop="openDialog"
        />
      </template>
    </q-input>

    <!-- Asset Selection Dialog -->
    <q-dialog v-model="showDialog" position="right" maximized>
      <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
        <!-- Header with Padding -->
        <q-card-section class="q-pb-sm">
          <div class="row items-center">
            <div class="text-h6">Select Asset</div>
            <q-space />
            <q-btn icon="close" flat round dense class="rounded-borders" @click="showDialog = false" />
          </div>
        </q-card-section>

        <!-- Info Banner -->
        <q-card-section class="q-pt-none q-pb-md">
          <q-banner dense class="bg-blue-1 text-blue-9 rounded-borders">
            <template #avatar>
              <q-icon name="info" color="blue-6" size="sm" />
            </template>
            <div class="text-caption">
              Use filters below to narrow down your search. Select an asset to continue.
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
                v-model="filters.search"
                outlined
                dense
                class="rounded-borders"
                label="Search by name"
                placeholder="Type to search..."
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
                v-model="filters.model"
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
          <q-btn flat dense label="Close" color="grey-7" size="sm" class="rounded-borders" @click="showDialog = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
}

/* Hover effects for list items */
:deep(.q-item) {
  transition: var(--mapex-transition-base);
}

:deep(.q-item:hover) {
  background-color: var(--mapex-surface-highlight);
}

:deep(.q-item.q-item--active) {
  background-color: rgba(25, 118, 210, 0.08) !important;
  border-left: 3px solid var(--q-primary);
}

/* Better spacing for filter inputs */
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

/* Smooth transitions */
:deep(.q-badge),
:deep(.q-chip) {
  transition: var(--mapex-transition-base);
}

/* Footer padding (ensure proper spacing) */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

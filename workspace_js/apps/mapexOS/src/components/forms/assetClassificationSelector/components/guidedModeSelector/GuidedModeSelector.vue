<script setup lang="ts">
defineOptions({
  name: 'GuidedModeSelector'
});

import type { GuidedModeSelectorProps, GuidedModeSelectorEmits, ListOption } from './interfaces';
import type { AssetClassification } from '../../interfaces';

import { ref, watch, computed } from 'vue';

import { useAssetClassificationSelectorTranslations } from '@composables/i18n/components/forms/assetClassificationSelector/useAssetClassificationSelectorTranslations';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { notifyFail } from '@utils/alert';

const logger = useLogger('GuidedModeSelector');

const props = defineProps<GuidedModeSelectorProps>();
const emit = defineEmits<GuidedModeSelectorEmits>();

const t = useAssetClassificationSelectorTranslations();

// State
const selectedCategoryId = ref<string | undefined>(undefined);
const selectedManufacturerId = ref<string | undefined>(undefined);
const selectedModelId = ref<string | undefined>(undefined);
// Defaults to "1.0.0" — the field is optional in the UI but the backend
// always receives a value so users can advance without ever filling it in.
const version = ref<string>('1.0.0');

const categories = ref<ListOption[]>([]);
const manufacturers = ref<ListOption[]>([]);
const models = ref<ListOption[]>([]);

const loadingCategories = ref(false);
const loadingManufacturers = ref(false);
const loadingModels = ref(false);

// Pagination state
const categoryPage = ref(1);
const manufacturerPage = ref(1);
const modelPage = ref(1);

const categoryHasMore = ref(true);
const manufacturerHasMore = ref(true);
const modelHasMore = ref(true);

const categoryTotalPages = ref(1);
const manufacturerTotalPages = ref(1);
const modelTotalPages = ref(1);

// Computed
const manufacturerDisabled = computed(() => !selectedCategoryId.value || props.disabled);
const modelDisabled = computed(() => !selectedManufacturerId.value || props.disabled);

// Initialize from modelValue
watch(() => props.modelValue, (newValue) => {
  if (newValue) {
    selectedCategoryId.value = newValue.categoryId;
    selectedManufacturerId.value = newValue.manufacturerId;
    selectedModelId.value = newValue.modelId;
    version.value = newValue.version || '1.0.0';

    // Pre-populate dropdown options with the selected items (for edit mode)
    // This ensures the names are displayed correctly even if these items aren't in the first page
    if (newValue.categoryId && newValue.categoryName) {
      const categoryExists = categories.value.some(c => c.id === newValue.categoryId);
      if (!categoryExists) {
        categories.value = [
          { id: newValue.categoryId, name: newValue.categoryName, value: newValue.categoryId },
          ...categories.value
        ];
      }
    }

    if (newValue.manufacturerId && newValue.manufacturerName) {
      const manufacturerExists = manufacturers.value.some(m => m.id === newValue.manufacturerId);
      if (!manufacturerExists) {
        manufacturers.value = [
          { id: newValue.manufacturerId, name: newValue.manufacturerName, value: newValue.manufacturerId },
          ...manufacturers.value
        ];
      }
      // Also fetch the full manufacturer list for the selected category
      if (newValue.categoryId) {
        void fetchManufacturers(newValue.categoryId);
      }
    }

    if (newValue.modelId && newValue.modelName) {
      const modelExists = models.value.some(m => m.id === newValue.modelId);
      if (!modelExists) {
        models.value = [
          { id: newValue.modelId, name: newValue.modelName, value: newValue.modelId },
          ...models.value
        ];
      }
      // Also fetch the full model list for the selected manufacturer
      if (newValue.manufacturerId) {
        void fetchModels(newValue.manufacturerId);
      }
    }
  } else {
    selectedCategoryId.value = undefined;
    selectedManufacturerId.value = undefined;
    selectedModelId.value = undefined;
    version.value = '1.0.0';
  }
}, { immediate: true });

// Fetch categories on mount
async function fetchCategories(page = 1, append = false) {
  if (!apis.mapexOS?.lists) {
    return;
  }

  try {
    loadingCategories.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_category',
      page,
      perPage: 20,
    });

    const newItems = response.items.map((item: any) => ({
      id: item.id,
      name: item.name,
      value: item.id,
    }));

    if (append) {
      categories.value = [...categories.value, ...newItems];
    } else {
      categories.value = newItems;
    }

    // Update pagination state
    categoryPage.value = page;
    categoryTotalPages.value = response.pagination?.totalPages || 1;
    categoryHasMore.value = page < (response.pagination?.totalPages || 1);
  } catch (err: any) {
    logger.error('Error fetching categories:', err);
    notifyFail({ message: t.errors.loadFailed.value });
  } finally {
    loadingCategories.value = false;
  }
}

// Fetch manufacturers when category changes
async function fetchManufacturers(categoryId: string, page = 1, append = false) {
  if (!apis.mapexOS?.lists || !categoryId) {
    return;
  }

  try {
    loadingManufacturers.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_manufacturer',
      parentId: categoryId,
      page,
      perPage: 20,
    });

    const newItems = response.items.map((item: any) => ({
      id: item.id,
      name: item.name,
      value: item.value,
    }));

    if (append) {
      manufacturers.value = [...manufacturers.value, ...newItems];
    } else {
      manufacturers.value = newItems;
    }

    // Update pagination state
    manufacturerPage.value = page;
    manufacturerTotalPages.value = response.pagination?.totalPages || 1;
    manufacturerHasMore.value = page < (response.pagination?.totalPages || 1);
  } catch (err: any) {
    logger.error('Error fetching manufacturers:', err);
    notifyFail({ message: t.errors.loadFailed.value });
  } finally {
    loadingManufacturers.value = false;
  }
}

// Fetch models when manufacturer changes
async function fetchModels(manufacturerId: string, page = 1, append = false) {
  if (!apis.mapexOS?.lists || !manufacturerId) {
    return;
  }

  try {
    loadingModels.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_model',
      parentId: manufacturerId,
      page,
      perPage: 20,
    });

    const newItems = response.items.map((item: any) => ({
      id: item.id,
      name: item.name,
      value: item.value,
    }));

    if (append) {
      models.value = [...models.value, ...newItems];
    } else {
      models.value = newItems;
    }

    // Update pagination state
    modelPage.value = page;
    modelTotalPages.value = response.pagination?.totalPages || 1;
    modelHasMore.value = page < (response.pagination?.totalPages || 1);
  } catch (err: any) {
    logger.error('Error fetching models:', err);
    notifyFail({ message: t.errors.loadFailed.value });
  } finally {
    loadingModels.value = false;
  }
}

// Infinite scroll handlers
async function onCategoryScroll(index: number, done: (stop?: boolean) => void) {
  if (!categoryHasMore.value || loadingCategories.value) {
    done(true);
    return;
  }

  await fetchCategories(categoryPage.value + 1, true);
  done(!categoryHasMore.value);
}

async function onManufacturerScroll(index: number, done: (stop?: boolean) => void) {
  if (!manufacturerHasMore.value || loadingManufacturers.value || !selectedCategoryId.value) {
    done(true);
    return;
  }

  await fetchManufacturers(selectedCategoryId.value, manufacturerPage.value + 1, true);
  done(!manufacturerHasMore.value);
}

async function onModelScroll(index: number, done: (stop?: boolean) => void) {
  if (!modelHasMore.value || loadingModels.value || !selectedManufacturerId.value) {
    done(true);
    return;
  }

  await fetchModels(selectedManufacturerId.value, modelPage.value + 1, true);
  done(!modelHasMore.value);
}

// Handle category selection
function handleCategoryChange(categoryId: string | null) {
  selectedCategoryId.value = categoryId || undefined;
  selectedManufacturerId.value = undefined;
  selectedModelId.value = undefined;
  version.value = '1.0.0';
  manufacturers.value = [];
  models.value = [];

  // Reset manufacturer pagination
  manufacturerPage.value = 1;
  manufacturerHasMore.value = true;

  if (categoryId) {
    void fetchManufacturers(categoryId);
  }

  emitValue();
}

// Handle manufacturer selection
function handleManufacturerChange(manufacturerId: string | null) {
  selectedManufacturerId.value = manufacturerId || undefined;
  selectedModelId.value = undefined;
  version.value = '1.0.0';
  models.value = [];

  // Reset model pagination
  modelPage.value = 1;
  modelHasMore.value = true;

  if (manufacturerId) {
    void fetchModels(manufacturerId);
  }

  emitValue();
}

// Handle model selection
function handleModelChange(modelId: string | null) {
  selectedModelId.value = modelId || undefined;
  emitValue();
}

// Handle version input — keep "1.0.0" as fallback when the user empties the field
function handleVersionInput() {
  if (!version.value || !version.value.trim()) {
    version.value = '1.0.0';
  }
  emitValue();
}

// Emit updated value
function emitValue() {
  if (selectedCategoryId.value && selectedManufacturerId.value && selectedModelId.value && version.value) {
    const category = categories.value.find(c => c.id === selectedCategoryId.value);
    const manufacturer = manufacturers.value.find(m => m.id === selectedManufacturerId.value);
    const model = models.value.find(m => m.id === selectedModelId.value);

    const classification: AssetClassification = {
      categoryId: selectedCategoryId.value,
      manufacturerId: selectedManufacturerId.value,
      modelId: selectedModelId.value,
      version: version.value,
      categoryName: category?.name ?? undefined,
      manufacturerName: manufacturer?.name ?? undefined,
      modelName: model?.name ?? undefined,
    };

    emit('update:modelValue', classification);
  } else {
    emit('update:modelValue', undefined);
  }
}

// Load categories on mount
void fetchCategories();
</script>

<template>
  <div class="row q-col-gutter-md">
    <!-- Category Select -->
    <div class="col-12 col-lg-6">
      <q-select
        v-model="selectedCategoryId"
        outlined
        dense
        emit-value
        map-options
        clearable
        :label="t.guided.category.label.value"
        :placeholder="t.guided.category.placeholder.value"
        :hint="t.guided.category.hint.value"
        :options="categories"
        :loading="loadingCategories"
        :disable="disabled"
        :rules="required ? [(val: any) => !!val || t.guided.category.required.value] : []"
        option-value="id"
        option-label="name"
        @update:model-value="handleCategoryChange"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey">
              {{ t.guided.noOptions.value }}
            </q-item-section>
          </q-item>
        </template>
        <template #after-options>
          <q-infinite-scroll
            v-if="categoryHasMore"
            @load="onCategoryScroll"
            :offset="50"
          >
            <template #loading>
              <div class="row justify-center q-my-sm">
                <q-spinner color="primary" size="sm" />
              </div>
            </template>
          </q-infinite-scroll>
        </template>
      </q-select>
    </div>

    <!-- Manufacturer Select -->
    <div class="col-12 col-lg-6">
      <q-select
        v-model="selectedManufacturerId"
        outlined
        dense
        emit-value
        map-options
        clearable
        :label="t.guided.manufacturer.label.value"
        :placeholder="t.guided.manufacturer.placeholder.value"
        :hint="manufacturerDisabled ? t.guided.manufacturer.disabled.value : t.guided.manufacturer.hint.value"
        :options="manufacturers"
        :loading="loadingManufacturers"
        :disable="manufacturerDisabled"
        :rules="required ? [(val: any) => !!val || t.guided.manufacturer.required.value] : []"
        option-value="id"
        option-label="name"
        @update:model-value="handleManufacturerChange"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey">
              {{ t.guided.noOptions.value }}
            </q-item-section>
          </q-item>
        </template>
        <template #after-options>
          <q-infinite-scroll
            v-if="manufacturerHasMore && !manufacturerDisabled"
            @load="onManufacturerScroll"
            :offset="50"
          >
            <template #loading>
              <div class="row justify-center q-my-sm">
                <q-spinner color="primary" size="sm" />
              </div>
            </template>
          </q-infinite-scroll>
        </template>
      </q-select>
    </div>

    <!-- Model Select -->
    <div class="col-12 col-lg-6">
      <q-select
        v-model="selectedModelId"
        outlined
        dense
        emit-value
        map-options
        clearable
        :label="t.guided.model.label.value"
        :placeholder="t.guided.model.placeholder.value"
        :hint="modelDisabled ? t.guided.model.disabled.value : t.guided.model.hint.value"
        :options="models"
        :loading="loadingModels"
        :disable="modelDisabled"
        :rules="required ? [(val: any) => !!val || t.guided.model.required.value] : []"
        option-value="id"
        option-label="name"
        @update:model-value="handleModelChange"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey">
              {{ t.guided.noOptions.value }}
            </q-item-section>
          </q-item>
        </template>
        <template #after-options>
          <q-infinite-scroll
            v-if="modelHasMore && !modelDisabled"
            @load="onModelScroll"
            :offset="50"
          >
            <template #loading>
              <div class="row justify-center q-my-sm">
                <q-spinner color="primary" size="sm" />
              </div>
            </template>
          </q-infinite-scroll>
        </template>
      </q-select>
    </div>

    <!-- Version Input -->
    <div class="col-12 col-lg-6">
      <q-input
        v-model="version"
        outlined
        dense
        :label="t.guided.version.label.value"
        :placeholder="t.guided.version.placeholder.value"
        :hint="t.guided.version.hint.value"
        @update:model-value="handleVersionInput"
      />
    </div>
  </div>
</template>

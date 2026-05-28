<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateSelectorDialog'
});

/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { AssetTemplateSelectorDialogProps, AssetTemplateSelectorDialogEmits } from './interfaces';

/** VUE IMPORTS */
import { computed, watch } from 'vue';

/** COMPONENTS */
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';

/** COMPOSABLES */
import { useAssetTemplates } from '@composables/assets/assetTemplates';
import { useCommonPlaceholders } from '@composables/i18n';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AssetTemplateSelectorDialogProps>(), {
  selectedTemplateIds: () => [],
  multiSelect: true,
});

const emit = defineEmits<AssetTemplateSelectorDialogEmits>();

/** COMPOSABLES & STORES */
const {
  templates,
  isLoading,
  isLoadingMore,
  filters,
  pagination,
  categoryOptions,
  manufacturerOptions,
  modelOptions,
  loadingCategories,
  loadingManufacturers,
  loadingModels,
  fetchTemplates,
  loadCategories,
  handleCategoryChange,
  handleManufacturerChange,
} = useAssetTemplates();
const { placeholders } = useCommonPlaceholders();

/** COMPUTED */

/**
 * Total templates count from pagination
 */
const totalTemplates = computed(() => pagination.value.totalItems);

/**
 * Whether more pages are available
 */
const hasMorePages = computed(() => pagination.value.currentPage < pagination.value.totalPages);

/**
 * Dialog title based on multiSelect mode
 */
const dialogTitle = computed(() =>
  props.multiSelect ? 'Select Asset Templates' : 'Select Asset Template',
);

/**
 * Info banner text
 */
const bannerText = computed(() =>
  props.multiSelect
    ? 'Select one or more asset templates. Click Confirm when done.'
    : 'Click on a template to select it.',
);

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
 * Watch dialog open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchTemplates();
    void loadCategories();
  }
});

/** FUNCTIONS */

/**
 * Handle template selection from GenericSelectorDialog
 *
 * @param {any[]} items - Selected items
 */
function handleSelect(items: any[]): void {
  emit('select', items as AssetTemplateResponse[]);
}

/**
 * Handle cancel action
 */
function handleCancel(): void {
  emit('cancel');
}

/**
 * Handle search query
 *
 * @param {string} query - Search query
 */
function handleSearch(query: string): void {
  filters.value.name = query || undefined;
  pagination.value.currentPage = 1;
  void fetchTemplates();
}

/**
 * Handle filter change (category, status, manufacturer, model)
 * Resets pagination and refetches
 */
function onFilterChange(): void {
  pagination.value.currentPage = 1;
  void fetchTemplates();
}

/**
 * Handle load more (infinite scroll)
 */
function handleLoadMore(): void {
  if (!isLoadingMore.value && hasMorePages.value) {
    pagination.value.currentPage++;
    void fetchTemplates(true);
  }
}
</script>

<template>
  <GenericSelectorDialog
    :model-value="modelValue"
    :title="dialogTitle"
    icon="description"
    icon-color="primary"
    :items="templates"
    item-key="id"
    :multi-select="multiSelect"
    :selected-ids="selectedTemplateIds"
    :loading="isLoading"
    :loading-more="isLoadingMore"
    :has-more-pages="hasMorePages"
    :total-items="totalTemplates"
    :search-placeholder="placeholders.searchByName.value"
    :info-banner="{ text: bannerText }"
    empty-text="No templates found"
    empty-icon="inbox"
    results-icon="description"
    footer-icon="description"
    item-noun-singular="template"
    item-noun-plural="templates"
    @update:model-value="emit('update:modelValue', $event)"
    @select="handleSelect"
    @cancel="handleCancel"
    @search="handleSearch"
    @load-more="handleLoadMore"
  >
    <!-- Domain-specific filters -->
    <template #filters>
      <!-- Category -->
      <div class="col-12 col-sm-6">
        <q-select
          v-model="filters.categoryId"
          outlined
          dense
          clearable
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

      <!-- Status -->
      <div class="col-12 col-sm-6">
        <q-select
          v-model="filters.status"
          outlined
          dense
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

      <!-- Manufacturer -->
      <div class="col-12 col-sm-6">
        <q-select
          v-model="filters.manufacturerId"
          outlined
          dense
          clearable
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

      <!-- Model -->
      <div class="col-12 col-sm-6">
        <q-select
          v-model="filters.modelId"
          outlined
          dense
          clearable
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
    </template>

    <!-- Item rendering -->
    <template #item="{ item }">
      <q-item-section avatar>
        <q-avatar
          :color="item.enabled ? 'primary' : 'grey-5'"
          icon="description"
          text-color="white"
        />
      </q-item-section>
      <q-item-section>
        <q-item-label>{{ item.name }}</q-item-label>
        <q-item-label caption>
          {{ item.categoryName }} • {{ item.manufacturerName }} • {{ item.modelName }}
          <span v-if="item.version"> • v{{ item.version }}</span>
        </q-item-label>
      </q-item-section>
    </template>
  </GenericSelectorDialog>
</template>

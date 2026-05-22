<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateSelector'
});

/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { AssetTemplateSelectorProps, AssetTemplateSelectorEmits } from './interfaces/assetTemplateSelector.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** UTILS */
import { handleApiError } from '@utils/error';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useAssetTemplateSelectorTranslations } from '@src/composables/i18n/components/selectors/useAssetTemplateSelectorTranslations';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('AssetTemplateSelector');

/** PROPS & EMITS */
const props = defineProps<AssetTemplateSelectorProps>();
const emit = defineEmits<AssetTemplateSelectorEmits>();

const t = useAssetTemplateSelectorTranslations();

// Filters
const filters = ref({
  search: '',
  isSystem: null as boolean | null,
  isTemplate: null as boolean | null,
  categoryId: null as string | null,
  manufacturerId: null as string | null,
  modelId: null as string | null,
});

// Dynamic filter options state
const categoryOptions = ref<Array<{ label: string; value: string }>>([]);
const manufacturerOptions = ref<Array<{ label: string; value: string }>>([]);
const modelOptions = ref<Array<{ label: string; value: string }>>([]);

const loadingCategories = ref(false);
const loadingManufacturers = ref(false);
const loadingModels = ref(false);

// Data
const templates = ref<AssetTemplateResponse[]>([]);
const selectedTemplate = ref<AssetTemplateResponse | null>(null);

// Pagination
const currentPage = ref(1);
const totalPages = ref(1);
const totalTemplates = ref(0);
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

// Load templates with current filters
async function loadTemplates(append = false) {
  if (append) {
    isLoadingMore.value = true;
  } else {
    isLoading.value = true;
    templates.value = [];
  }

  try {
    if (!apis.assets) {
      throw new Error('Assets API not configured');
    }

    const response = await apis.assets.assetTemplate.list({
      page: currentPage.value,
      perPage,
      ...(filters.value.search && { name: filters.value.search }),
      ...(typeof filters.value.isSystem === 'boolean' && { isSystem: filters.value.isSystem }),
      ...(typeof filters.value.isTemplate === 'boolean' && { isTemplate: filters.value.isTemplate }),
      ...(filters.value.categoryId && { categoryId: filters.value.categoryId }),
      ...(filters.value.manufacturerId && { manufacturerId: filters.value.manufacturerId }),
      ...(filters.value.modelId && { modelId: filters.value.modelId }),
      sort: 'name:asc',
    });

    if (append) {
      templates.value.push(...response.items);
    } else {
      templates.value = response.items;
    }

    totalTemplates.value = response.pagination.totalItems;
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
      void loadTemplates(true);
    }
  }
}

// Select template
function selectTemplate(template: AssetTemplateResponse) {
  selectedTemplate.value = template;
  logger.debug('Selected template:', template);
  logger.debug('Template ID:', template.id);
  logger.debug('Emitting modelValue:', template.id || null);
  emit('update:modelValue', template.id || null);
  emit('update:selectedTemplate', template);
}

// Filter change handler
function onFilterChange() {
  currentPage.value = 1;
  void loadTemplates();
}

// Fetch categories from API
async function fetchCategories() {
  if (!apis.mapexOS?.lists) {
    return;
  }

  try {
    loadingCategories.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_category',
      page: 1,
      perPage: 100,
    });

    categoryOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching categories:', err);
  } finally {
    loadingCategories.value = false;
  }
}

// Fetch manufacturers based on selected category
async function fetchManufacturers(categoryId: string) {
  if (!apis.mapexOS?.lists || !categoryId) {
    manufacturerOptions.value = [];
    return;
  }

  try {
    loadingManufacturers.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_manufacturer',
      parentId: categoryId,
      page: 1,
      perPage: 100,
    });

    manufacturerOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching manufacturers:', err);
    manufacturerOptions.value = [];
  } finally {
    loadingManufacturers.value = false;
  }
}

// Fetch models based on selected manufacturer
async function fetchModels(manufacturerId: string) {
  if (!apis.mapexOS?.lists || !manufacturerId) {
    modelOptions.value = [];
    return;
  }

  try {
    loadingModels.value = true;
    const response = await apis.mapexOS.lists.list({
      type: 'asset_model',
      parentId: manufacturerId,
      page: 1,
      perPage: 100,
    });

    modelOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching models:', err);
    modelOptions.value = [];
  } finally {
    loadingModels.value = false;
  }
}

// Handle cascading filter changes
function onCategoryChange(value: string | null) {
  // Reset dependent fields
  filters.value.manufacturerId = null;
  filters.value.modelId = null;
  manufacturerOptions.value = [];
  modelOptions.value = [];

  // Fetch manufacturers if category selected
  if (value) {
    void fetchManufacturers(value);
  }

  // Reload templates
  onFilterChange();
}

function onManufacturerChange(value: string | null) {
  // Reset dependent field
  filters.value.modelId = null;
  modelOptions.value = [];

  // Fetch models if manufacturer selected
  if (value) {
    void fetchModels(value);
  }

  // Reload templates
  onFilterChange();
}

// Search watcher (debounce handled by q-input debounce prop)
watch(() => filters.value.search, () => {
  currentPage.value = 1;
  void loadTemplates();
});

// Watch for external model value changes
watch(() => props.modelValue, (newValue) => {
  if (!newValue) {
    selectedTemplate.value = null;
  } else if (selectedTemplate.value?.id !== newValue) {
    // Find template in current list
    const found = templates.value.find(t => t.id === newValue);
    if (found) {
      selectedTemplate.value = found;
    }
  }
});

// Initial load
onMounted(() => {
  void loadTemplates();
  void fetchCategories();
});
</script>

<template>
  <div class="asset-template-selector">
    <!-- Filters Row 1: Search + Template Type + Template Source -->
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

    <!-- Cascading Filters Row 2: Category + Manufacturer + Model -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-select
          v-model="filters.categoryId"
          outlined
          dense
          clearable
          emit-value
          map-options
          label="Category"
          :options="categoryOptions"
          :loading="loadingCategories"
          @update:model-value="onCategoryChange"
        >
          <template v-slot:prepend>
            <q-icon name="category" color="primary" />
          </template>
        </q-select>
      </div>
      <div class="col-12 col-md-4">
        <q-select
          v-model="filters.manufacturerId"
          outlined
          dense
          clearable
          emit-value
          map-options
          label="Manufacturer"
          :options="manufacturerOptions"
          :loading="loadingManufacturers"
          :disable="!filters.categoryId"
          @update:model-value="onManufacturerChange"
        >
          <template v-slot:prepend>
            <q-icon name="factory" color="primary" />
          </template>
        </q-select>
      </div>
      <div class="col-12 col-md-4">
        <q-select
          v-model="filters.modelId"
          outlined
          dense
          clearable
          emit-value
          map-options
          label="Model"
          :options="modelOptions"
          :loading="loadingModels"
          :disable="!filters.manufacturerId"
          @update:model-value="onFilterChange"
        >
          <template v-slot:prepend>
            <q-icon name="memory" color="primary" />
          </template>
        </q-select>
      </div>
    </div>

    <!-- Results Count -->
    <div class="text-caption text-grey-7 q-mb-sm">
      <q-icon name="info" size="xs" class="q-mr-xs" />
      {{ t.results.found.value.replace('{count}', String(totalTemplates)) }}
    </div>

    <!-- Template List with Virtual Scroll -->
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

        <!-- Template Cards -->
        <q-card
          v-for="template in templates"
          :key="template.id || ''"
          flat
          bordered
          clickable
          class="template-card q-mb-sm"
          :class="{ 'selected': props.modelValue === template.id }"
          @click="selectTemplate(template)"
        >
          <q-card-section class="row items-start q-pa-md">
            <!-- Selection Radio -->
            <div class="col-auto q-mr-md q-pt-xs">
              <q-radio
                :model-value="props.modelValue"
                :val="template.id"
                color="primary"
                size="md"
              />
            </div>

            <!-- Template Info -->
            <div class="col" style="min-width: 0;">
              <!-- Line 1: Name and Status Badge -->
              <div class="row items-center q-mb-xs no-wrap">
                <div class="text-subtitle1 text-weight-medium ellipsis" style="flex: 1; min-width: 0;">
                  {{ template.name }}
                </div>
                <q-badge
                  :color="template.enabled ? 'positive' : 'grey'"
                  :label="template.enabled ? t.badge.active.value : t.badge.inactive.value"
                  class="q-ml-sm"
                  style="flex-shrink: 0;"
                />
              </div>

              <!-- Line 2: Description -->
              <div class="text-body2 text-grey-7 q-mb-xs ellipsis description-text">
                {{ template.description || t.noDescription.value }}
                <AppTooltip
                  v-if="template.description"
                  :content="template.description"
                  :offset="[0, 8]"
                  max-width="400px"
                />
              </div>

              <!-- Line 3: Manufacturer, Model, Version -->
              <div class="text-caption text-grey-8 row items-center no-wrap">
                <q-icon name="factory" size="xs" class="q-mr-xs" />
                <span class="ellipsis" style="max-width: 200px;">{{ template.manufacturerName }}</span>
                <q-separator vertical inset class="q-mx-sm" />
                <q-icon name="devices" size="xs" class="q-mr-xs" />
                <span class="ellipsis" style="max-width: 180px;">{{ template.modelName }}</span>
                <template v-if="template.version">
                  <q-separator vertical inset class="q-mx-sm" />
                  <q-icon name="label" size="xs" class="q-mr-xs" />
                  <span style="max-width: 60px;">v{{ template.version }}</span>
                </template>
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
          v-if="!isLoading && templates.length === 0"
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

    <!-- Selected Template Preview -->
    <div v-if="selectedTemplate" class="q-mt-md">
      <q-card flat bordered class="bg-blue-1">
        <q-card-section class="q-pa-md">
          <div class="text-caption text-weight-medium text-primary q-mb-xs">
            <q-icon name="check_circle" size="xs" class="q-mr-xs" />
            {{ t.selected.label.value }}
          </div>
          <div class="row items-center">
            <div class="text-body1 text-weight-medium text-primary">
              {{ selectedTemplate.name }}
            </div>
          </div>
          <div class="text-caption text-grey-8 q-mt-xs">
            {{ selectedTemplate.manufacturerName }} • {{ selectedTemplate.modelName }}
            <template v-if="selectedTemplate.version">• v{{ selectedTemplate.version }}</template>
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

.template-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);
  background-color: white;
}

.template-card:hover {
  border-color: var(--q-primary);
  box-shadow: var(--mapex-shadow-sm);
}

.template-card.selected {
  border-color: var(--q-primary);
  border-width: 2px;
  background-color: rgba(59, 109, 94, 0.05);
}

.cursor-pointer {
  cursor: pointer;
}

.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.description-text {
  max-width: 400px;
  display: block;
}
</style>

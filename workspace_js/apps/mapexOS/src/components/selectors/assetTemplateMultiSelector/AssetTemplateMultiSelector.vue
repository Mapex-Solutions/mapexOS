<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateMultiSelector'
});

/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { AssetTemplateMultiSelectorProps, AssetTemplateMultiSelectorEmits } from './interfaces/assetTemplateMultiSelector.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useAssetTemplates } from '@composables/assets/assetTemplates';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('AssetTemplateMultiSelector');

/** PROPS & EMITS */
const props = withDefaults(defineProps<AssetTemplateMultiSelectorProps>(), {
  label: 'Select Asset Templates',
});

const emit = defineEmits<AssetTemplateMultiSelectorEmits>();

/** COMPOSABLES & STORES */
const {
  templates,
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
  fetchTemplates,
  loadCategories,
  handleCategoryChange,
  handleManufacturerChange,
} = useAssetTemplates();

/** STATE */
const selectedTemplates = ref<AssetTemplateResponse[]>([]);
const showDialog = ref(false);
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
 * Total templates count from pagination
 */
const totalTemplates = computed(() => pagination.value.totalItems);

/**
 * Selected template IDs
 */
const selectedIds = computed(() => selectedTemplates.value.map(t => t.id).filter(Boolean) as string[]);

/**
 * Extracted paths from selected templates
 */
const extractedPaths = computed(() => {
  return selectedTemplates.value
    .filter(t => t.assetIdPath)
    .map(t => ({
      templateId: t.id || '',
      templateName: t.name || '',
      assetIdPath: t.assetIdPath || '',
    }));
});

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
 * Watch extracted paths and emit
 */
watch(extractedPaths, (newPaths) => {
  emit('update:extractedPaths', newPaths);
}, { deep: true });

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
      void fetchTemplates(true);
    }
  }
}

/**
 * Toggle template selection
 * @param {AssetTemplateResponse} template - Template to toggle
 */
function toggleTemplate(template: AssetTemplateResponse): void {
  const index = selectedTemplates.value.findIndex(t => t.id === template.id);

  if (index >= 0) {
    selectedTemplates.value.splice(index, 1);
  } else {
    selectedTemplates.value.push(template);
  }

  emit('update:modelValue', selectedIds.value);
  emit('update:selectedTemplates', selectedTemplates.value);
}

/**
 * Check if template is selected
 * @param {AssetTemplateResponse} template - Template to check
 * @returns {boolean} True if selected
 */
function isSelected(template: AssetTemplateResponse): boolean {
  return selectedTemplates.value.some(t => t.id === template.id);
}

/**
 * Remove template from selection
 * @param {string} templateId - Template ID to remove
 */
function removeTemplate(templateId: string): void {
  const index = selectedTemplates.value.findIndex(t => t.id === templateId);
  if (index >= 0) {
    selectedTemplates.value.splice(index, 1);
    emit('update:modelValue', selectedIds.value);
    emit('update:selectedTemplates', selectedTemplates.value);
  }
}

/**
 * Clear all selections
 */
function clearAll(): void {
  selectedTemplates.value = [];
  emit('update:modelValue', []);
  emit('update:selectedTemplates', []);
}

/**
 * Filter change handler
 * Resets pagination and refetches templates
 */
function onFilterChange(): void {
  pagination.value.currentPage = 1;
  void fetchTemplates();
}

/**
 * Open dialog and load templates
 */
function openDialog(): void {
  showDialog.value = true;
  if (templates.value.length === 0) {
    void fetchTemplates();
  }
  if (categoryOptions.value.length === 0) {
    void loadCategories();
  }
}

// Watch modelValue changes from parent
watch(() => props.modelValue, async (newValue) => {
  if (newValue && newValue.length > 0) {
    // Load the selected templates
    try {
      if (!apis.assets) return;

      const loadedTemplates: AssetTemplateResponse[] = [];
      for (const id of newValue) {
        const template = await apis.assets.assetTemplate.getById({ assetTemplateId: id });
        loadedTemplates.push(template);
      }
      selectedTemplates.value = loadedTemplates;
    } catch (error: any) {
      logger.error('Error loading templates:', error);
    }
  } else if (newValue.length === 0) {
    selectedTemplates.value = [];
  }
});

// Load initial templates if modelValue is provided
onMounted(async () => {
  if (props.modelValue && props.modelValue.length > 0) {
    try {
      if (!apis.assets) return;

      const loadedTemplates: AssetTemplateResponse[] = [];
      for (const id of props.modelValue) {
        const template = await apis.assets.assetTemplate.getById({ assetTemplateId: id });
        loadedTemplates.push(template);
      }
      selectedTemplates.value = loadedTemplates;
    } catch (error: any) {
      logger.error('Error loading initial templates:', error);
    }
  }
});
</script>

<template>
  <div>
    <!-- Selected Templates Display -->
    <div class="q-mb-sm">
      <div class="text-caption text-grey-7 q-mb-xs">{{ label }}</div>

      <!-- Selected templates list -->
      <div v-if="selectedTemplates.length > 0" class="q-mt-md q-mb-md">
        <div class="text-caption text-grey-7 q-mb-sm">
          <q-icon name="check_circle" color="positive" size="xs" />
          Selected templates:
        </div>
        <q-list bordered class="rounded-borders template-list">
          <q-item v-for="(template, idx) in selectedTemplates" :key="idx" class="template-list-item">
            <q-item-section avatar>
              <q-icon name="description" color="purple-6" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ template.name }}</q-item-label>
              <q-item-label caption>
                <code class="path-code">{{ template.assetIdPath || 'No UUID path' }}</code>
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-btn
                flat
                round
                dense
                size="sm"
                icon="delete"
                color="negative"
                @click="removeTemplate(template.id || '')"
              >
                <AppTooltip content="Remove template" />
              </q-btn>
            </q-item-section>
          </q-item>
        </q-list>
      </div>

      <!-- Empty state -->
      <div v-else class="text-grey-6 text-caption q-mb-sm">
        No templates selected
      </div>

      <!-- Action buttons -->
      <div class="row q-gutter-sm q-mt-md">
        <q-btn
          outline
          dense
          color="primary"
          icon="add"
          label="Select Templates"
          size="sm"
          class="rounded-borders"
          @click="openDialog"
        />
        <q-btn
          v-if="selectedTemplates.length > 0"
          flat
          dense
          color="negative"
          icon="clear"
          label="Clear All"
          size="sm"
          class="rounded-borders"
          @click="clearAll"
        />
      </div>
    </div>

    <!-- Template Selection Dialog -->
    <q-dialog v-model="showDialog" position="right" maximized>
      <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
        <!-- Header with Padding -->
        <q-card-section class="q-pb-sm">
          <div class="row items-center">
            <div class="text-h6">Select Asset Templates</div>
            <q-space />
            <q-btn icon="close" flat round dense class="rounded-borders" @click="showDialog = false" />
          </div>
        </q-card-section>

        <!-- Info Banner -->
        <q-card-section class="q-pt-none q-pb-md">
          <q-banner dense class="bg-purple-1 text-purple-9 rounded-borders">
            <template #avatar>
              <q-icon name="info" color="purple-6" size="sm" />
            </template>
            <div class="text-caption">
              Select one or more asset templates. Click Done when finished.
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
            <q-icon name="description" size="xs" class="q-mr-xs" />
            Results
          </div>
        </q-card-section>

        <!-- Templates List -->
        <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
          <!-- Loading state -->
          <div v-if="isLoading" class="q-pa-md text-center">
            <q-spinner color="primary" size="3em" />
            <div class="text-grey-7 q-mt-md">Loading templates...</div>
          </div>

          <!-- Empty state -->
          <div v-else-if="templates.length === 0" class="q-pa-md text-center">
            <q-icon name="inbox" size="4em" color="grey-5" />
            <div class="text-grey-7 q-mt-md">No templates found</div>
          </div>

          <!-- Templates List with Infinite Scroll -->
          <q-scroll-area
            v-else
            ref="scrollAreaRef"
            style="height: 100%;"
            @scroll="onScroll"
          >
            <q-list separator>
              <q-item
                v-for="template in templates"
                :key="template.id || `template-${Math.random()}`"
                clickable
                :active="isSelected(template)"
                @click="toggleTemplate(template)"
              >
                <q-item-section avatar>
                  <q-checkbox
                    :model-value="isSelected(template)"
                    color="primary"
                    @click.stop="toggleTemplate(template)"
                  />
                </q-item-section>

                <q-item-section avatar>
                  <q-avatar color="purple-6" icon="description" text-color="white" />
                </q-item-section>

                <q-item-section>
                  <q-item-label>{{ template.name }}</q-item-label>
                  <q-item-label caption class="text-grey-7">
                    {{ template.assetIdPath || 'No UUID path' }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-icon
                    v-if="template.assetIdPath"
                    name="info"
                    color="primary"
                    size="sm"
                    class="cursor-pointer"
                  >
                    <AppTooltip :content="`UUID Path: ${template.assetIdPath}`" />
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
            <q-icon name="description" size="xs" class="q-mr-xs" />
            {{ totalTemplates }} {{ totalTemplates === 1 ? 'template' : 'templates' }}
          </div>
          <q-space />
          <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders q-mr-sm" @click="showDialog = false" />
          <q-btn flat dense label="Confirm Selection" color="primary" size="sm" class="rounded-borders" @click="showDialog = false" />
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

/* Footer padding (ensure proper spacing) */
:deep(.q-card__actions) {
  padding: 16px !important;
}

/* Checkbox styling */
:deep(.q-checkbox__inner) {
  border-radius: var(--mapex-radius-xs);
}

/* Selected templates list spacing */
.template-list {
  background-color: var(--mapex-surface-bg);
}

.template-list-item {
  padding: 12px 16px !important;
  min-height: 60px;
}

.template-list-item:not(:last-child) {
  border-bottom: 1px solid var(--mapex-divider);
}

.path-code {
  background-color: var(--mapex-surface-highlight);
  padding: 4px 8px;
  border-radius: var(--mapex-radius-xs);
  font-family: 'Courier New', Consolas, monospace;
  font-size: 0.9em;
  color: var(--q-primary);
  font-weight: 500;
}
</style>

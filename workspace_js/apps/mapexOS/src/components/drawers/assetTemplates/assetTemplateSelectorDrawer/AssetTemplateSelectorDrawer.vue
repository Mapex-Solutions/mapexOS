<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateSelectorDrawer'
});

/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { DrawerScrollInfo } from '@components/drawers/common/genericDrawer';
import type { AssetTemplateSelectorDrawerProps, AssetTemplateSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { GenericDrawer } from '@components/drawers/common/genericDrawer';
import { AssetTemplateDetailsModal } from '@components/assetTemplates';
import { AppTooltip } from '@components/tooltips';
import { InfoBanner } from '@components/banners';

/** COMPOSABLES */
import { useAssetTemplates } from '@composables/assets/assetTemplates';
import { useCommonPlaceholders, useAssetTemplatesTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AssetTemplateSelectorDrawerProps>(), {
  selectedTemplateIds: () => [],
  multiSelect: true,
});

const emit = defineEmits<AssetTemplateSelectorDrawerEmits>();

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
  loadMore,
  loadCategories,
  handleCategoryChange,
  handleManufacturerChange,
} = useAssetTemplates();

const logger = useLogger('AssetTemplateSelectorDrawer');
const { placeholders } = useCommonPlaceholders();
const t = useAssetTemplatesTranslations();

/** STATE */
const selectedTemplates = ref<AssetTemplateResponse[]>([]);
const detailTemplate = ref<AssetTemplateResponse | null>(null);
const showDetailsModal = ref(false);

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
 * Drawer title based on multiSelect mode
 */
const drawerTitle = computed(() =>
  props.multiSelect ? 'Select Asset Templates' : 'Select Asset Template',
);

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
    void fetchTemplates();
    void loadCategories();
  }
});

/** FUNCTIONS */

/**
 * Toggle template selection
 *
 * @param {AssetTemplateResponse} template - Template to toggle
 */
function toggleTemplate(template: AssetTemplateResponse): void {
  if (!props.multiSelect) {
    selectedTemplates.value = [template];
    emit('select', selectedTemplates.value);
    emit('update:modelValue', false);
    return;
  }

  const index = selectedTemplates.value.findIndex(t => t.id === template.id);

  if (index >= 0) {
    selectedTemplates.value.splice(index, 1);
  } else {
    selectedTemplates.value.push(template);
  }
}

/**
 * Open the shared details modal for a template. Rows carry only a summary, so
 * fetch the full template (scripts, dynamic/available fields) first; fall back
 * to the summary on failure.
 *
 * @param {AssetTemplateResponse} template - Template to inspect
 * @returns {Promise<void>}
 */
async function openDetails(template: AssetTemplateResponse): Promise<void> {
  try {
    detailTemplate.value = template.id
      ? await apis.assets.assetTemplate.getById({ assetTemplateId: template.id })
      : template;
  } catch (err) {
    logger.error('Failed to load template details:', err);
    detailTemplate.value = template;
  }
  showDetailsModal.value = true;
}

/**
 * Select the template shown in the details modal, then close the modal.
 *
 * @param {AssetTemplateResponse} template - Template to select
 * @returns {void}
 */
function selectFromModal(template: AssetTemplateResponse): void {
  showDetailsModal.value = false;
  if (!isSelected(template)) toggleTemplate(template);
}

/**
 * Check if template is selected
 *
 * @param {AssetTemplateResponse} template - Template to check
 * @returns {boolean} True if selected
 */
function isSelected(template: AssetTemplateResponse): boolean {
  return selectedTemplates.value.some(t => t.id === template.id);
}

/**
 * Confirm selection
 */
function confirmSelection(): void {
  logger.debug('Confirming selection:', selectedTemplates.value);
  emit('select', selectedTemplates.value);
  emit('update:modelValue', false);
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
 * Infinite-scroll handler: load the next page when the drawer content is
 * scrolled near the bottom. loadMore is a no-op when there is no next page
 * or a fetch is already in flight.
 *
 * @param {DrawerScrollInfo} info - Vertical scroll metrics from the drawer.
 */
function onScroll(info: DrawerScrollInfo): void {
  const reachedThreshold =
    info.verticalPosition + info.verticalContainerSize >= info.verticalSize * 0.8;

  if (reachedThreshold) {
    void loadMore();
  }
}

/**
 * Cancel handler
 */
function handleCancel(): void {
  emit('cancel');
  emit('update:modelValue', false);
}
</script>

<template>
  <GenericDrawer
    :model-value="modelValue"
    :title="drawerTitle"
    icon="description"
    :width="600"
    @update:model-value="emit('update:modelValue', $event)"
    @scroll="onScroll"
    @close="handleCancel"
  >
    <!-- Info Banner -->
    <InfoBanner variant="info" dense class="q-mb-md">
      <div class="text-caption">
        {{ multiSelect ? 'Select one or more asset templates. Click Confirm when done.' : 'Click on a template to select it.' }}
      </div>
    </InfoBanner>

    <!-- Filters -->
    <div class="q-mb-md">
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

        <!-- Category and Status -->
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

        <!-- Manufacturer and Model -->
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
      </div>
    </div>

    <q-separator class="q-mb-md" />

    <!-- Results Header -->
    <div class="text-overline text-grey-7 q-mb-sm">
      <q-icon name="description" size="xs" class="q-mr-xs" />
      Results
    </div>

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

    <!-- Templates List -->
    <q-list v-else separator class="template-list">
      <q-item
        v-for="template in templates"
        :key="template.id || `template-${Math.random()}`"
        clickable
        :active="isSelected(template)"
        @click="toggleTemplate(template)"
      >
        <q-item-section v-if="multiSelect" avatar>
          <q-checkbox
            :model-value="isSelected(template)"
            color="primary"
            @click.stop="toggleTemplate(template)"
          />
        </q-item-section>

        <q-item-section v-else avatar>
          <q-avatar
            :color="template.enabled ? 'primary' : 'grey-5'"
            icon="description"
            text-color="white"
          />
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ template.name }}</q-item-label>
          <q-item-label caption class="text-grey-7">
            {{ template.categoryName }} • {{ template.manufacturerName }} • {{ template.modelName }}
            <span v-if="template.version"> • v{{ template.version }}</span>
          </q-item-label>
        </q-item-section>

        <q-item-section side>
          <q-btn
            flat
            round
            dense
            size="sm"
            icon="info"
            color="primary"
            @click.stop="openDetails(template)"
          >
            <AppTooltip :content="t.drawer.title.value" />
          </q-btn>
        </q-item-section>
      </q-item>

      <!-- Load More Indicator -->
      <div v-if="isLoadingMore" class="q-pa-md text-center">
        <q-spinner color="primary" size="2em" />
      </div>
    </q-list>

    <!-- Footer -->
    <template #footer>
      <div class="text-caption text-grey-7">
        <q-icon name="description" size="xs" class="q-mr-xs" />
        {{ totalTemplates }} {{ totalTemplates === 1 ? 'template' : 'templates' }}
      </div>
      <q-space />
      <q-btn flat dense no-caps label="Cancel" color="grey-7" size="sm" @click="handleCancel" />
      <q-btn
        v-if="multiSelect"
        flat
        dense
        no-caps
        label="Confirm Selection"
        color="primary"
        size="sm"
        :disable="selectedTemplates.length === 0"
        @click="confirmSelection"
      />
    </template>
  </GenericDrawer>

  <!-- Shared details viewer, opened from a row's info button -->
  <AssetTemplateDetailsModal v-model="showDetailsModal" :template="detailTemplate">
    <template #actions="{ template: detail }">
      <q-btn
        unelevated
        no-caps
        icon="check"
        color="primary"
        :label="t.drawer.selectTemplate.value"
        :disable="!detail"
        @click="detail && selectFromModal(detail)"
      />
    </template>
  </AssetTemplateDetailsModal>
</template>

<style lang="scss" scoped>
// Template list hover effects
.template-list {
  :deep(.q-item) {
    transition: all var(--mapex-transition-base) ease;
  }

  :deep(.q-item:hover) {
    background-color: var(--mapex-surface-bg);
  }

  :deep(.q-item.q-item--active) {
    background-color: rgba(var(--q-primary-rgb), 0.08) !important;
    border-left: 3px solid var(--q-primary);
  }
}

// Better spacing for filter inputs
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

// Checkbox styling
:deep(.q-checkbox__inner) {
  border-radius: var(--mapex-radius-xs);
}
</style>

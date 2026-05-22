<template>
  <GenericDrawer
    :model-value="modelValue"
    :title="t.title.value"
    icon="schema"
    :close-tooltip="t.closeTooltip.value"
    @update:model-value="emit('update:modelValue', $event)"
    @close="emit('update:modelValue', false)"
  >
    <!-- Source Toggle -->
    <div class="filter-field">
      <div class="filter-field__label">
        <q-icon name="source" size="xs" color="primary" class="q-mr-sm" />
        <span>{{ t.sourceLabel.value }}</span>
      </div>
      <q-btn-toggle
        v-model="sourceType"
        spread
        no-caps
        rounded
        unelevated
        toggle-color="primary"
        color="grey-2"
        text-color="grey-8"
        class="filter-toggle"
        :options="sourceOptions"
        @update:model-value="handleSourceChange"
      />
    </div>

    <!-- Autocomplete Search -->
    <div class="filter-field">
      <div class="filter-field__label">
        <q-icon name="search" size="xs" color="primary" class="q-mr-sm" />
        <span>{{ sourceType === 'asset' ? t.searchAsset.value : t.searchAssetTemplate.value }}</span>
      </div>
      <q-select
        v-model="selectedSourceId"
        outlined
        dense
        use-input
        hide-selected
        fill-input
        :input-debounce="AUTOCOMPLETE_DEBOUNCE"
        emit-value
        map-options
        clearable
        class="filter-input"
        :options="autocompleteOptions"
        :loading="loadingAutocomplete"
        :placeholder="sourceType === 'asset' ? t.searchAsset.value : t.searchAssetTemplate.value"
        option-value="id"
        option-label="label"
        :virtual-scroll-item-size="48"
        @filter="handleAutocomplete"
        @update:model-value="handleSourceSelect"
        @virtual-scroll="handleVirtualScroll"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey">
              {{ t.noResults.value }}
            </q-item-section>
          </q-item>
        </template>
        <template #option="{ itemProps, opt }">
          <q-item v-bind="itemProps">
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white" size="sm">
                {{ opt.label?.charAt(0)?.toUpperCase() || '?' }}
              </q-avatar>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ opt.label }}</q-item-label>
              <q-item-label v-if="opt.caption" caption>{{ opt.caption }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <!-- Resolving Template Spinner -->
    <div v-if="loadingTemplate" class="row items-center q-pa-sm q-gutter-sm">
      <q-spinner color="primary" size="xs" />
      <span class="text-caption text-grey-7">{{ t.resolvingTemplate.value }}</span>
    </div>

    <!-- Filters Section (Progressive Disclosure) -->
    <template v-if="resolvedTemplateId">
      <!-- Filters Header with count -->
      <div class="filters-header q-mb-sm">
        <q-icon name="tune" size="xs" color="secondary" class="q-mr-sm" />
        <span>{{ t.filtersSection.value }}</span>
        <q-badge
          :label="`${activeFilters.length} / ${availableFields.length}`"
          color="grey-3"
          text-color="grey-7"
          class="q-ml-sm"
        />
      </div>

      <!-- Loading Fields -->
      <div v-if="loadingFields" class="row items-center q-pa-sm q-gutter-sm">
        <q-spinner color="secondary" size="xs" />
        <span class="text-caption text-grey-7">{{ t.loadingFields.value }}</span>
      </div>

      <!-- No Filterable Fields in Template -->
      <div v-else-if="availableFields.length === 0" class="text-caption text-grey-6 q-pa-sm">
        {{ t.noFields.value }}
      </div>

      <!-- Filter Builder -->
      <template v-else>
        <!-- Active Filter Cards -->
        <div v-if="activeFilters.length > 0" class="q-gutter-md q-mt-md q-mb-sm">
          <div
            v-for="filter in activeFilters"
            :key="filter.key"
            class="filter-card"
          >
            <!-- Card Header: type icon + label + remove button -->
            <div class="filter-card__header">
              <q-icon
                :name="getFieldTypeIcon(filter.type)"
                size="18px"
                :color="getFieldTypeColor(filter.type)"
              />
              <span class="filter-card__label">{{ filter.label }}</span>
              <q-space />
              <q-btn
                flat
                round
                dense
                size="sm"
                icon="close"
                color="grey-7"
                class="filter-card__close"
                @click="removeFilter(filter.key)"
              />
            </div>

            <!-- Card Controls -->
            <div class="filter-card__controls">
              <!-- Boolean: Toggle True/False (no operator) -->
              <q-btn-toggle
                v-if="filter.type === 'boolean'"
                v-model="filter.value"
                spread
                no-caps
                rounded
                unelevated
                toggle-color="secondary"
                color="grey-2"
                text-color="grey-8"
                clearable
                size="sm"
                :options="[
                  { label: t.booleanTrue.value, value: true },
                  { label: t.booleanFalse.value, value: false },
                ]"
              />

              <!-- Non-boolean: Operator + Value(s) -->
              <template v-else>
                <!-- Operator Select -->
                <q-select
                  v-model="filter.operator"
                  :options="getOperatorOptions(filter.type)"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="filter-input"
                  option-value="value"
                  option-label="translatedLabel"
                >
                  <template #selected-item="{ opt }">
                    <div class="row items-center no-wrap">
                      <q-icon :name="opt.icon" size="xs" class="q-mr-xs" />
                      <span>{{ opt.translatedLabel }}</span>
                    </div>
                  </template>
                  <template #option="{ itemProps, opt }">
                    <q-item v-bind="itemProps" dense>
                      <q-item-section avatar>
                        <q-icon :name="opt.icon" size="xs" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label>{{ opt.translatedLabel }}</q-item-label>
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <!-- Value Inputs -->
                <div class="q-mt-sm">
                  <!-- Between: Two inputs side by side -->
                  <div v-if="filter.operator === 'between'" class="row q-gutter-sm">
                    <div class="col">
                      <q-input
                        v-model="filter.value"
                        :type="getInputType(filter.type)"
                        :label="t.operators.rangeFrom.value"
                        outlined
                        dense
                        clearable
                        class="filter-input"
                      />
                    </div>
                    <div class="col">
                      <q-input
                        v-model="filter.endValue"
                        :type="getInputType(filter.type)"
                        :label="t.operators.rangeTo.value"
                        outlined
                        dense
                        clearable
                        class="filter-input"
                      />
                    </div>
                  </div>

                  <!-- Single value input -->
                  <q-input
                    v-else
                    v-model="filter.value"
                    :type="getInputType(filter.type)"
                    outlined
                    dense
                    clearable
                    class="filter-input"
                  />
                </div>
              </template>
            </div>
          </div>
        </div>

        <!-- Empty State with centered Add button -->
        <div v-else class="empty-state">
          <q-icon name="filter_list_off" size="48px" color="grey-4" />
          <div class="text-subtitle2 text-grey-6 q-mt-md">{{ t.emptyTitle.value }}</div>
          <div class="text-caption text-grey-5 q-mt-xs">{{ t.emptyDescription.value }}</div>
          <q-btn
            round
            unelevated
            color="primary"
            icon="add"
            size="md"
            class="q-mt-lg"
          >
            <!-- Add Filter Menu (empty state) -->
            <q-menu class="add-filter-menu" :offset="[0, 4]">
              <div class="q-pa-sm">
                <q-input
                  v-model="fieldSearchTerm"
                  :placeholder="t.searchFields.value"
                  outlined
                  dense
                  clearable
                  autofocus
                  class="filter-input"
                >
                  <template #prepend>
                    <q-icon name="search" size="xs" />
                  </template>
                </q-input>
              </div>
              <q-separator />
              <div
                v-if="Object.keys(availableFieldsByType).length === 0 && fieldSearchTerm"
                class="text-caption text-grey-6 q-pa-md text-center"
              >
                {{ t.noResults.value }}
              </div>
              <q-list class="add-filter-menu__list">
                <template v-for="(fields, typeKey) in availableFieldsByType" :key="typeKey">
                  <q-item-label header class="text-overline text-grey-6 q-pb-none">
                    {{ t.fieldTypeHeaders[typeKey as keyof typeof t.fieldTypeHeaders]?.value || typeKey }}
                  </q-item-label>
                  <q-item
                    v-for="field in fields"
                    :key="field.key"
                    v-close-popup
                    clickable
                    @click="addFilter(field)"
                  >
                    <q-item-section avatar>
                      <q-icon :name="getFieldTypeIcon(field.type)" size="18px" :color="getFieldTypeColor(field.type)" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-body2">{{ field.label }}</q-item-label>
                    </q-item-section>
                  </q-item>
                </template>
              </q-list>
            </q-menu>
          </q-btn>
        </div>

        <!-- Add Filter Button (when filters exist) -->
        <q-btn
          v-if="activeFilters.length > 0 && !allFieldsAdded"
          flat
          no-caps
          dense
          color="primary"
          icon="add"
          :label="t.addFilter.value"
          class="q-mt-sm full-width"
        >
          <!-- Add Filter Menu (with filters) -->
          <q-menu class="add-filter-menu" anchor="bottom left" self="top left" :offset="[0, 4]">
            <div class="q-pa-sm">
              <q-input
                v-model="fieldSearchTerm"
                :placeholder="t.searchFields.value"
                outlined
                dense
                clearable
                autofocus
                class="filter-input"
              >
                <template #prepend>
                  <q-icon name="search" size="xs" />
                </template>
              </q-input>
            </div>
            <q-separator />
            <div
              v-if="Object.keys(availableFieldsByType).length === 0 && fieldSearchTerm"
              class="text-caption text-grey-6 q-pa-md text-center"
            >
              {{ t.noResults.value }}
            </div>
            <q-list class="add-filter-menu__list">
              <template v-for="(fields, typeKey) in availableFieldsByType" :key="typeKey">
                <q-item-label header class="text-overline text-grey-6 q-pb-none">
                  {{ t.fieldTypeHeaders[typeKey as keyof typeof t.fieldTypeHeaders]?.value || typeKey }}
                </q-item-label>
                <q-item
                  v-for="field in fields"
                  :key="field.key"
                  v-close-popup
                  clickable
                  @click="addFilter(field)"
                >
                  <q-item-section avatar>
                    <q-icon :name="getFieldTypeIcon(field.type)" size="18px" :color="getFieldTypeColor(field.type)" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-body2">{{ field.label }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-list>
          </q-menu>
        </q-btn>
      </template>
    </template>

    <!-- Footer Actions -->
    <template #footer>
      <q-btn
        flat
        no-caps
        color="grey-7"
        icon="refresh"
        class="col"
        :label="t.buttons.reset.value"
        @click="handleReset"
      >
        <AppTooltip :content="t.buttons.resetTooltip.value" />
      </q-btn>
      <q-btn
        unelevated
        no-caps
        :color="activeFilters.length > 0 ? 'warning' : 'primary'"
        icon="check"
        class="col"
        :disable="!resolvedTemplateId"
        :label="t.buttons.apply.value"
        @click="handleApply"
      >
        <AppTooltip :content="t.buttons.applyTooltip.value" />
      </q-btn>
    </template>
  </GenericDrawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'DynamicFiltersDrawer'
});

/** TYPE IMPORTS */
import type {
  DynamicFiltersDrawerProps,
  DynamicFiltersDrawerEmits,
  DynamicSourceType,
  DynamicFilterField,
  DynamicFieldDefinition,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { GenericDrawer } from '../genericDrawer';

/** COMPOSABLES */
import { useDynamicFiltersDrawerTranslations } from '@composables/i18n/components/drawers/dynamicFiltersDrawer';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert/notify';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_SOURCE_TYPE,
  EVA_TYPE_MAP,
  AUTOCOMPLETE_MIN_CHARS,
  AUTOCOMPLETE_DEBOUNCE,
  AUTOCOMPLETE_PER_PAGE,
  DEFAULT_OPERATOR_BY_TYPE,
  EVA_OPERATORS_BY_TYPE,
} from './constants';

/** PROPS & EMITS */
const props = defineProps<DynamicFiltersDrawerProps>();
const emit = defineEmits<DynamicFiltersDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useDynamicFiltersDrawerTranslations();
const logger = useLogger('DynamicFiltersDrawer');

/** STATE */
const sourceType = ref<DynamicSourceType>(DEFAULT_SOURCE_TYPE);
const selectedSourceId = ref<string | null>(null);
const selectedSourceName = ref('');
const resolvedTemplateId = ref<string | undefined>();
const resolvedTemplateName = ref('');
const availableFields = ref<DynamicFieldDefinition[]>([]);
const activeFilters = ref<DynamicFilterField[]>([]);
const fieldSearchTerm = ref('');
const loadingAutocomplete = ref(false);
const loadingTemplate = ref(false);
const loadingFields = ref(false);
const autocompleteOptions = ref<Array<{ id: string; label: string; caption?: string }>>([]);

/** Asset template infinite scroll state */
const assetTemplatePage = ref(1);
const assetTemplateHasMore = ref(false);
const assetTemplateCurrentSearch = ref('');

/** COMPUTED */

/**
 * Source toggle options built from translations
 */
const sourceOptions = computed(() => [
  { label: t.sourceAsset.value, value: 'asset' },
  { label: t.sourceAssetTemplate.value, value: 'assetTemplate' },
]);


/**
 * Available fields not yet added as active filters, filtered by search term
 */
const availableFieldsFiltered = computed(() => {
  const addedKeys = new Set(activeFilters.value.map(f => f.key));
  let fields = availableFields.value.filter(f => !addedKeys.has(f.key));
  if (fieldSearchTerm.value) {
    const term = fieldSearchTerm.value.toLowerCase();
    fields = fields.filter(f => f.label.toLowerCase().includes(term));
  }
  return fields;
});

/**
 * Available fields grouped by type for the add-filter menu
 */
const availableFieldsByType = computed(() => {
  const groups: Record<string, DynamicFieldDefinition[]> = {};
  for (const field of availableFieldsFiltered.value) {
    (groups[field.type] ??= []).push(field);
  }
  return groups;
});

/**
 * Whether all available fields have been added as active filters
 */
const allFieldsAdded = computed(() =>
  availableFields.value.length > 0 &&
  activeFilters.value.length >= availableFields.value.length,
);

/** WATCHERS */

/**
 * Emit pending-change when active filters change
 * Pending = user has added filters but not yet applied
 */
watch(
  () => activeFilters.value.length,
  (count) => {
    emit('pending-change', count > 0);
  },
);

/**
 * Populate demo data when demo prop becomes true (tour mode)
 * Shows a realistic example with asset selected and filters pre-filled
 */
watch(
  () => props.demo,
  (isDemo) => {
    if (isDemo) {
      sourceType.value = 'asset';
      selectedSourceId.value = 'demo-asset-001';
      selectedSourceName.value = 'Sensor Station - DEMO';
      resolvedTemplateId.value = 'demo-template-001';
      resolvedTemplateName.value = 'IoT Sensor Template';
      autocompleteOptions.value = [
        { id: 'demo-asset-001', label: 'Sensor Station - DEMO', caption: 'IoT Sensor Template' },
      ];

      availableFields.value = [
        { key: 'temperature', label: 'temperature', type: 'number', originalType: 'number', fieldId: 0 },
        { key: 'humidity', label: 'humidity', type: 'number', originalType: 'number', fieldId: 1 },
        { key: 'batteryLevel', label: 'batteryLevel', type: 'number', originalType: 'number', fieldId: 2 },
        { key: 'deviceType', label: 'deviceType', type: 'string', originalType: 'string', fieldId: 3 },
        { key: 'deviceID', label: 'deviceID', type: 'string', originalType: 'string', fieldId: 4 },
        { key: 'isActive', label: 'isActive', type: 'boolean', originalType: 'bool', fieldId: 5 },
        { key: 'lastSeen', label: 'lastSeen', type: 'date', originalType: 'date', fieldId: 6 },
      ];

      activeFilters.value = [
        { key: 'temperature', label: 'temperature', type: 'number', operator: 'gte', value: '25.5', fieldId: 0, originalType: 'number' },
        { key: 'isActive', label: 'isActive', type: 'boolean', operator: 'eq', value: true, fieldId: 5, originalType: 'bool' },
      ];
    } else {
      handleReset();
    }
  },
);

/** FUNCTIONS */

/**
 * Map operator labelKey to translated label for display
 *
 * @param {string} labelKey - i18n key like 'operators.equals'
 * @returns {string} Translated label
 */
function translateOperatorLabel(labelKey: string): string {
  const key = labelKey.replace('operators.', '') as keyof typeof t.operators;
  return t.operators[key]?.value ?? labelKey;
}

/**
 * Get operator options for a field type, with translated labels
 *
 * @param {string} fieldType - Field type (string, number, boolean, date)
 * @returns {Array} Operator options with translatedLabel
 */
function getOperatorOptions(fieldType: string) {
  const ops = EVA_OPERATORS_BY_TYPE[fieldType] ?? EVA_OPERATORS_BY_TYPE.string ?? [];
  return ops.map(op => ({
    ...op,
    translatedLabel: translateOperatorLabel(op.labelKey),
  }));
}

/**
 * Get icon for field type
 *
 * @param {string} fieldType - Field type
 * @returns {string} Material icon name
 */
function getFieldTypeIcon(fieldType: string): string {
  switch (fieldType) {
    case 'number': return 'tag';
    case 'string': return 'text_fields';
    case 'boolean': return 'toggle_on';
    case 'date': return 'event';
    default: return 'help_outline';
  }
}

/**
 * Get color for field type
 *
 * @param {string} fieldType - Field type
 * @returns {string} Quasar color name
 */
function getFieldTypeColor(fieldType: string): string {
  switch (fieldType) {
    case 'number': return 'blue-7';
    case 'string': return 'teal-7';
    case 'boolean': return 'orange-7';
    case 'date': return 'purple-7';
    default: return 'grey-7';
  }
}

/** Quasar q-input valid type values */
type QInputType = 'text' | 'number' | 'textarea' | 'email' | 'search' | 'tel' | 'file' | 'url' | 'password' | 'date' | 'time' | 'datetime-local';

/**
 * Get HTML input type for a field type
 *
 * @param {string} fieldType - Field type
 * @returns {QInputType} Quasar q-input type attribute
 */
function getInputType(fieldType: string): QInputType {
  switch (fieldType) {
    case 'number': return 'number';
    case 'date': return 'datetime-local';
    default: return 'text';
  }
}

/**
 * Add a field from the available pool as an active filter
 *
 * @param {DynamicFieldDefinition} field - Field definition to add
 */
function addFilter(field: DynamicFieldDefinition): void {
  activeFilters.value.push({
    key: field.key,
    label: field.label,
    type: field.type,
    operator: DEFAULT_OPERATOR_BY_TYPE[field.type] || 'eq',
    value: field.type === 'boolean' ? null : '',
    fieldId: field.fieldId,
    ...(field.originalType != null && { originalType: field.originalType }),
  });
}

/**
 * Remove an active filter by key (field returns to the available pool)
 *
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  const idx = activeFilters.value.findIndex(f => f.key === key);
  if (idx !== -1) activeFilters.value.splice(idx, 1);
}

/**
 * Handle autocomplete search based on current source type
 *
 * @param {string} searchTerm - Search term from input
 * @param {Function} update - Callback to update options
 * @param {Function} abort - Callback to abort
 */
function handleAutocomplete(
  searchTerm: string,
  update: (callback: () => void) => void,
  abort: () => void,
): void {
  if (sourceType.value === 'asset' && searchTerm.length < AUTOCOMPLETE_MIN_CHARS) {
    abort();
    return;
  }

  loadingAutocomplete.value = true;

  if (sourceType.value === 'asset') {
    void apis.assets.asset.list({
      name: searchTerm,
      perPage: AUTOCOMPLETE_PER_PAGE,
    }).then((response) => {
      update(() => {
        autocompleteOptions.value = (response.items || []).map((item: any) => ({
          id: item.id,
          label: item.name || item.id,
          caption: item.assetTemplateName || undefined,
        }));
        loadingAutocomplete.value = false;
      });
    }).catch((error: any) => {
      logger.error('Error searching asset templates', error);
      update(() => {
        autocompleteOptions.value = [];
        loadingAutocomplete.value = false;
      });
    });
  } else {
    // assetTemplate: reset pagination, load first page
    assetTemplatePage.value = 1;
    assetTemplateCurrentSearch.value = searchTerm;
    void apis.assets.assetTemplate.list({
      name: searchTerm || undefined,
      perPage: AUTOCOMPLETE_PER_PAGE,
      page: 1,
    }).then((response) => {
      update(() => {
        autocompleteOptions.value = (response.items || []).map((item: any) => ({
          id: item.id,
          label: item.name || item.id,
        }));
        assetTemplateHasMore.value = (response.items?.length ?? 0) >= AUTOCOMPLETE_PER_PAGE;
        loadingAutocomplete.value = false;
      });
    }).catch((error: any) => {
      logger.error('Error searching asset templates', error);
      update(() => {
        autocompleteOptions.value = [];
        assetTemplateHasMore.value = false;
        loadingAutocomplete.value = false;
      });
    });
  }
}

/**
 * Load next page of asset templates and append to options (infinite scroll)
 */
function loadMoreAssetTemplates(): void {
  if (!assetTemplateHasMore.value || loadingAutocomplete.value) return;

  loadingAutocomplete.value = true;
  const nextPage = assetTemplatePage.value + 1;

  void apis.assets.assetTemplate.list({
    name: assetTemplateCurrentSearch.value || undefined,
    perPage: AUTOCOMPLETE_PER_PAGE,
    page: nextPage,
  }).then((response) => {
    const newItems = (response.items || []).map((item: any) => ({
      id: item.id,
      label: item.name || item.id,
    }));
    autocompleteOptions.value = [...autocompleteOptions.value, ...newItems];
    assetTemplatePage.value = nextPage;
    assetTemplateHasMore.value = newItems.length >= AUTOCOMPLETE_PER_PAGE;
    loadingAutocomplete.value = false;
  }).catch((error: any) => {
    logger.error('Error loading more asset templates', error);
    loadingAutocomplete.value = false;
  });
}

/**
 * Handle virtual scroll to trigger infinite scroll when near the end
 *
 * @param {{ to: number }} param - Virtual scroll event with last visible index
 */
function handleVirtualScroll({ to }: { to: number }): void {
  if (sourceType.value !== 'assetTemplate') return;
  const lastIndex = autocompleteOptions.value.length - 1;
  if (to >= lastIndex - 3) {
    loadMoreAssetTemplates();
  }
}

/**
 * Handle source type change (Asset / Asset Template)
 * Resets selection, fields, active filters, and infinite scroll state
 */
function handleSourceChange(): void {
  selectedSourceId.value = null;
  selectedSourceName.value = '';
  resolvedTemplateId.value = undefined;
  resolvedTemplateName.value = '';
  availableFields.value = [];
  activeFilters.value = [];
  fieldSearchTerm.value = '';
  autocompleteOptions.value = [];
  assetTemplatePage.value = 1;
  assetTemplateHasMore.value = false;
  assetTemplateCurrentSearch.value = '';
}

/**
 * Handle source entity selection from autocomplete
 * Resolves the asset template and loads available fields
 *
 * @param {string | null} value - Selected entity ID (or null when cleared)
 */
async function handleSourceSelect(value: string | null): Promise<void> {
  if (!value) {
    selectedSourceName.value = '';
    resolvedTemplateId.value = undefined;
    resolvedTemplateName.value = '';
    availableFields.value = [];
    activeFilters.value = [];
    fieldSearchTerm.value = '';
    return;
  }

  const selectedOption = autocompleteOptions.value.find(opt => opt.id === value);
  selectedSourceName.value = selectedOption?.label || value;

  loadingTemplate.value = true;

  try {
    if (sourceType.value === 'asset') {
      await resolveTemplateFromAsset(value);
    } else {
      await resolveTemplateFromAssetTemplate(value);
    }
  } catch (error: any) {
    logger.error('Error resolving template', error);
    resolvedTemplateId.value = undefined;
    resolvedTemplateName.value = '';
    availableFields.value = [];
    activeFilters.value = [];
  } finally {
    loadingTemplate.value = false;
  }
}

/**
 * Resolve asset template from an asset ID
 * Uses asset.assetTemplateId directly, then fetches the template
 *
 * @param {string} assetId - Asset ID
 */
async function resolveTemplateFromAsset(assetId: string): Promise<void> {
  const asset = await apis.assets.asset.getById({ assetId });

  if (!asset.assetTemplateId) {
    resolvedTemplateId.value = undefined;
    resolvedTemplateName.value = '';
    availableFields.value = [];
    activeFilters.value = [];
    notifyFail({ message: t.noTemplateForAsset.value });
    return;
  }

  resolvedTemplateId.value = asset.assetTemplateId;
  resolvedTemplateName.value = asset.assetTemplateName || asset.assetTemplateId;
  await loadAvailableFields(asset.assetTemplateId);
}

/**
 * Resolve template directly from a selected asset template ID
 * No extra fetch needed — the template ID is already known
 *
 * @param {string} templateId - Asset template ID selected by the user
 */
async function resolveTemplateFromAssetTemplate(templateId: string): Promise<void> {
  const selectedOption = autocompleteOptions.value.find(opt => opt.id === templateId);
  resolvedTemplateId.value = templateId;
  resolvedTemplateName.value = selectedOption?.label || templateId;
  await loadAvailableFields(templateId);
}

/**
 * Load available fields from an asset template into the pool
 * Populates availableFields but NOT activeFilters (user chooses which to add)
 *
 * @param {string} templateId - Asset template ID
 */
async function loadAvailableFields(templateId: string): Promise<void> {
  loadingFields.value = true;
  activeFilters.value = [];
  fieldSearchTerm.value = '';

  try {
    const template = await apis.assets.assetTemplate.getById({ assetTemplateId: templateId });

    if (template.dynamicFields?.length) {
      availableFields.value = template.dynamicFields
        .filter((f: any) => f.type !== 'geo' && f.status !== 0)
        .map((f: any) => ({
          key: f.field,
          label: f.field,
          type: EVA_TYPE_MAP[f.type] || 'string',
          originalType: f.type,
          fieldId: f.fieldId,
        }));
    } else {
      availableFields.value = [];
    }
  } catch (error: any) {
    logger.error('Error loading dynamic fields', error);
    availableFields.value = [];
  } finally {
    loadingFields.value = false;
  }
}

/**
 * Handle apply button click
 * Emits DynamicFiltersResult with active filters that have values
 */
function handleApply(): void {
  if (!resolvedTemplateId.value || !selectedSourceId.value) return;

  emit('apply', {
    sourceType: sourceType.value,
    sourceId: selectedSourceId.value,
    sourceName: selectedSourceName.value,
    assetTemplateId: resolvedTemplateId.value,
    templateName: resolvedTemplateName.value,
    fields: activeFilters.value.map(f => ({ ...f })),
  });
  emit('update:modelValue', false);
}

/**
 * Handle reset button click
 * Clears all active filters and emits reset event
 */
function handleReset(): void {
  sourceType.value = DEFAULT_SOURCE_TYPE;
  selectedSourceId.value = null;
  selectedSourceName.value = '';
  resolvedTemplateId.value = undefined;
  resolvedTemplateName.value = '';
  availableFields.value = [];
  activeFilters.value = [];
  fieldSearchTerm.value = '';
  autocompleteOptions.value = [];
  assetTemplatePage.value = 1;
  assetTemplateHasMore.value = false;
  assetTemplateCurrentSearch.value = '';
  emit('reset');
}

/**
 * Handle Enter key to apply filters when drawer is open
 *
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEnterKey(event: KeyboardEvent): void {
  if (event.key === 'Enter' && resolvedTemplateId.value) {
    handleApply();
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEnterKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEnterKey);
});
</script>

<style lang="scss" scoped>
.filter-field {
  padding-bottom: 16px;
  margin-bottom: 16px;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-of-type {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 8px;
  }

  &__label {
    display: flex;
    align-items: center;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--mapex-text-secondary);
    margin-bottom: 12px;
    letter-spacing: 0.01em;
  }
}

.filters-header {
  display: flex;
  align-items: center;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--mapex-text-secondary);
  padding-top: 4px;
  letter-spacing: 0.01em;
}

.filter-toggle {
  :deep(.q-btn) {
    font-weight: 500;
    padding: 6px 16px;
  }

  :deep(.q-btn--active) {
    font-weight: 600;
  }
}

.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }

  :deep(.q-field__label) {
    font-weight: 400;
  }
}

.filter-card {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-lg);
  overflow: hidden;

  &__header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 8px 8px 14px;
    background: var(--mapex-surface-elevated);
    border-bottom: 1px solid var(--mapex-divider);
  }

  &__label {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    letter-spacing: 0.01em;
  }

  &__close {
    opacity: 0.6;
    transition: opacity var(--mapex-transition-fast);

    &:hover {
      opacity: 1;
    }
  }

  &__controls {
    padding: 12px 14px;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 480px);
  padding: 40px 24px;
}

.add-filter-menu {
  min-width: 260px;
  max-width: 320px;

  &__list {
    max-height: 300px;
    overflow-y: auto;

    :deep(.q-item__section--avatar) {
      min-width: 32px;
    }
  }
}
</style>

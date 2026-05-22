<script setup lang="ts">
defineOptions({
  name: 'ListFilter'
});

import type { FilterListProps, FilterListEmitEvents } from './interfaces';

import { ref, reactive, computed, watch } from 'vue';
import { useListFilterTranslations } from '@composables/i18n';
import { SelectableChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { InputFilter } from './components/inputFilter';
import { SelectFilter } from './components/selectFilter';
import { MultiSelectFilter } from './components/multiSelectFilter';
import { DateRangeFilter } from './components/dateRangeFilter';
import { UserSelectFilter } from './components/userSelectFilter';

const props = withDefaults(defineProps<FilterListProps>(), {
  autoApply: false,
  showIncludeChildren: false,
  includeChildrenInitial: false,
  watchFields: () => [],
});

const emit = defineEmits<FilterListEmitEvents>();
const translations = useListFilterTranslations();

// Expand/collapse state
const isExpanded = ref(false);

function toggle() {
  isExpanded.value = !isExpanded.value;
}

// Include children toggle state
const includeChildren = ref(props.includeChildrenInitial ?? false);

// Local state for items - this ensures reactivity
const localItems = ref([...props.items]);

// Reactive object to hold current filter values
const filterValues = reactive<Record<string, any>>({});

// Initialize filter values
function initializeFilterValues() {
  localItems.value.forEach(cfg => {
    // Use defaultValue if provided
    if (cfg.defaultValue !== undefined) {
      filterValues[cfg.key] = cfg.defaultValue;
    } else if (cfg.type === 'select' || cfg.type === 'user-select') {
      filterValues[cfg.key] = null;
    } else if (cfg.type === 'multiselect') {
      filterValues[cfg.key] = [];
    } else if (cfg.type === 'daterange') {
      filterValues[cfg.key] = { from: null, to: null };
    } else {
      filterValues[cfg.key] = '';
    }
  });
}

// Initialize on mount
initializeFilterValues();

// Watch for changes in items prop to sync with local state
watch(() => props.items, (newItems) => {
  localItems.value = [...newItems];
}, { deep: true });

// Handle field change and emit to parent if field is being watched
function handleFieldChange(fieldKey: string, newValue: any) {
  if (props.watchFields && props.watchFields.includes(fieldKey)) {
    emit('fieldChange', fieldKey, newValue);
  }
}

// Determine if any filter is active
const hasActiveFilters = computed(() =>
    Object.entries(filterValues).some(([key, val]) => {
      const cfg = localItems.value.find(item => item.key === key);
      if (cfg?.type === 'multiselect') {
        return Array.isArray(val) && val.length > 0;
      }
      if (cfg?.type === 'daterange') {
        return val?.from != null || val?.to != null;
      }
      return val !== '' && val != null;
    }),
);

// Count active filters
const activeFiltersCount = computed(() => {
  return Object.entries(filterValues).filter(([key, val]) => {
    const cfg = localItems.value.find(item => item.key === key);
    if (cfg?.type === 'multiselect') {
      return Array.isArray(val) && val.length > 0;
    }
    if (cfg?.type === 'daterange') {
      return val?.from != null || val?.to != null;
    }
    return val !== '' && val != null;
  }).length;
});

// Get active filter chips for display when collapsed
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  Object.entries(filterValues).forEach(([key, val]) => {
    const cfg = localItems.value.find(item => item.key === key);
    if (!cfg) return;

    let displayValue = '';

    if (cfg.type === 'multiselect' && Array.isArray(val) && val.length > 0) {
      displayValue = `${val.length} selected`;
    } else if (cfg.type === 'daterange') {
      // Only show if at least one date is filled
      if (val?.from || val?.to) {
        const parts = [];
        if (val.from) parts.push(val.from);
        if (val.to) parts.push(val.to);
        displayValue = parts.join(' → ');
      }
    } else if (cfg.type === 'select' && val != null) {
      const option = cfg.options?.find((opt: any) => opt.value === val);
      displayValue = option?.label || String(val);
    } else if (typeof val === 'string' && val !== '') {
      displayValue = val;
    } else if (typeof val === 'number') {
      displayValue = String(val);
    } else if (typeof val === 'boolean') {
      displayValue = val ? 'Yes' : 'No';
    }

    if (displayValue) {
      chips.push({ key, label: cfg.label, value: displayValue });
    }
  });

  return chips;
});

// Clear all filters
function clearFilters() {
  localItems.value.forEach(cfg => {
    if (cfg.type === 'select' || cfg.type === 'user-select') {
      filterValues[cfg.key] = null;
    } else if (cfg.type === 'multiselect') {
      filterValues[cfg.key] = [];
    } else if (cfg.type === 'daterange') {
      filterValues[cfg.key] = { from: null, to: null };
    } else {
      filterValues[cfg.key] = '';
    }
  });
  emit('reset');
}

// Remove individual filter
function removeFilter(key: string) {
  const cfg = localItems.value.find(item => item.key === key);
  if (!cfg) return;

  if (cfg.type === 'select' || cfg.type === 'user-select') {
    filterValues[key] = null;
  } else if (cfg.type === 'multiselect') {
    filterValues[key] = [];
  } else if (cfg.type === 'daterange') {
    filterValues[key] = { from: null, to: null };
  } else {
    filterValues[key] = '';
  }

  if (props.autoApply) {
    applyFilters();
  }
}

// Apply filters and emit values
function applyFilters() {
  const values = { ...filterValues };

  // Add includeChildren to the emitted values if the toggle is shown
  if (props.showIncludeChildren) {
    values.includeChildren = includeChildren.value;
  }

  emit('apply', values);
}

// Auto-apply filters when values change (optional, can be enabled per config)
// Note: Individual fields use Quasar's native debounce prop (500ms)
if (props.autoApply) {
  watch(filterValues, () => {
    applyFilters();
  }, { deep: true });
}
</script>

<template>
  <q-card class="rounded-borders" :class="!isExpanded ? 'q-mb-sm' : 'q-mb-lg'">
    <q-card-section class="q-pa-none">
      <!-- Header -->
      <div class="row items-center justify-between q-px-md q-py-md cursor-pointer" @click="toggle">
        <div class="row items-center">
          <q-icon name="tune" size="sm" color="primary" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-medium text-primary">{{ translations.listFilter.title.value }}</div>
          <q-badge
              v-if="hasActiveFilters"
              color="primary"
              text-color="white"
              class="q-ml-sm"
              :label="`${activeFiltersCount} ${translations.listFilter.active.value}`"
          />
          <q-icon
              size="sm"
              color="grey-6"
              class="transition-transform q-ml-sm"
              :name="isExpanded ? 'expand_less' : 'expand_more'"
              :class="{ 'rotate-180': isExpanded }"
          />
        </div>

        <!-- Include Children Toggle -->
        <div v-if="props.showIncludeChildren" class="row items-center" @click.stop>
          <q-toggle
              v-model="includeChildren"
              icon="account_tree"
              color="primary"
              :label="translations.listFilter.includeChildren.value"
              @update:model-value="applyFilters"
          >
            <AppTooltip :content="translations.listFilter.includeChildrenTooltip.value" />
          </q-toggle>
        </div>
      </div>

      <!-- Active Filter Chips (when collapsed) -->
      <div v-if="!isExpanded && hasActiveFilters" class="q-px-md q-pb-md">
        <div class="row q-gutter-xs">
          <SelectableChip
            v-for="chip in activeFilterChips"
            :key="chip.key"
            :label="`${chip.label}: ${chip.value}`"
            icon="filter_alt"
            color="primary"
            size="sm"
            dense
            @remove="removeFilter(chip.key)"
          />
        </div>
      </div>

      <!-- Expandable Content -->
      <q-slide-transition>
        <div v-show="isExpanded" class="q-px-md q-pb-md border-top">
          <q-separator/>
          <div class="row q-col-gutter-md q-mt-sm">
            <div
                v-for="(cfg, index) in localItems"
                :class="cfg.grid || 'col-12 col-md-4'"
                :key="`${cfg.key}-${index}`"
            >
              <!-- Input filter -->
              <InputFilter
                  v-if="cfg.type === 'input'"
                  v-model="filterValues[cfg.key]"
                  :label="cfg.label"
                  :icon="cfg.icon"
                  :clearable="cfg.clearable ?? true"
                  :disabled="cfg.disabled ?? false"
                  :debounce="props.autoApply ? 500 : 0"
                  :mask="cfg.mask"
                  :type="cfg.inputType"
                  :placeholder="cfg.placeholder"
                  @update:model-value="(val) => handleFieldChange(cfg.key, val)"
                  @enter="applyFilters"
              />

              <!-- Select filter -->
              <SelectFilter
                  v-else-if="cfg.type === 'select'"
                  v-model="filterValues[cfg.key]"
                  :label="cfg.label"
                  :icon="cfg.icon"
                  :options="cfg.options || []"
                  :clearable="cfg.clearable ?? true"
                  :disabled="cfg.disabled ?? false"
                  @update:model-value="(val) => handleFieldChange(cfg.key, val)"
                  @enter="applyFilters"
              />

              <!-- Multi-select filter -->
              <MultiSelectFilter
                  v-else-if="cfg.type === 'multiselect'"
                  v-model="filterValues[cfg.key]"
                  :label="cfg.label"
                  :icon="cfg.icon"
                  :options="cfg.options || []"
                  :clearable="cfg.clearable ?? true"
                  :disabled="cfg.disabled ?? false"
                  @update:model-value="(val) => handleFieldChange(cfg.key, val)"
                  @enter="applyFilters"
              />

              <!-- Date range filter -->
              <DateRangeFilter
                  v-else-if="cfg.type === 'daterange'"
                  v-model="filterValues[cfg.key]"
                  :label="cfg.label"
                  :icon="cfg.icon"
                  :clearable="cfg.clearable ?? true"
                  :disabled="cfg.disabled ?? false"
                  @enter="applyFilters"
              />

              <!-- User select filter -->
              <UserSelectFilter
                  v-else-if="cfg.type === 'user-select'"
                  v-model="filterValues[cfg.key]"
                  :label="cfg.label"
                  :icon="cfg.icon"
                  :clearable="cfg.clearable ?? true"
                  :disabled="cfg.disabled ?? false"
                  :placeholder="cfg.placeholder"
                  @update:model-value="(val) => handleFieldChange(cfg.key, val)"
              />
            </div>
          </div>

          <!-- Actions -->
          <div class="row justify-end items-center q-mt-md q-pt-md border-top">
            <q-btn
                flat
                icon="refresh"
                color="grey-7"
                class="q-px-md rounded-borders"
                :label="translations.listFilter.clear.value"
                :disable="!hasActiveFilters"
                @click="clearFilters"
            />
            <q-btn
                unelevated
                icon="search"
                color="primary"
                class="q-px-md q-ml-sm rounded-borders"
                :label="translations.listFilter.apply.value"
                @click="applyFilters"
            />
          </div>
        </div>
      </q-slide-transition>
    </q-card-section>
  </q-card>
  <div v-if="!isExpanded" class="q-mb-lg text-grey-7">
    <q-icon name="info" size="sm" color="grey-5" />
    {{ translations.listFilter.hint.value }}
  </div>
</template>

<style lang="scss" scoped>

/* Separator */
.q-separator {
  background: var(--mapex-divider);
}

/* Rounded borders */
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
}

.hover-bg {
  transition: background-color 0.2s ease;

  &:hover {
    background-color: var(--mapex-surface-highlight);
  }
}

.hover-bg-section {
  transition: background-color 0.2s ease;
  padding: 4px 8px;
  border-radius: var(--mapex-radius-xs);

  &:hover {
    background-color: var(--mapex-surface-highlight);
  }
}
</style>
<template>
  <GenericDrawer
    :model-value="modelValue"
    :title="drawerTitle"
    icon="tune"
    :width="width"
    :close-tooltip="t.closeTooltip.value"
    @update:model-value="emit('update:modelValue', $event)"
    @close="emit('update:modelValue', false)"
  >
    <!-- Filter Fields -->
    <div
      v-for="(field, index) in fields"
      :key="field.key"
      class="filter-field"
      :class="{ 'filter-field--last': index === fields.length - 1 }"
    >
      <!-- Field Label -->
      <div class="filter-field__label">
        <q-icon :name="field.icon" size="xs" color="primary" class="q-mr-sm" />
        <span>{{ field.label }}</span>
      </div>

      <!-- Toggle Field (with options) -->
      <q-btn-toggle
        v-if="field.type === 'toggle'"
        v-model="localValues[field.key]"
        spread
        no-caps
        rounded
        unelevated
        toggle-color="primary"
        color="grey-2"
        text-color="grey-8"
        class="filter-toggle"
        :options="field.options"
        :disable="field.disabled"
        @update:model-value="(val) => handleFieldChange(field.key, val)"
      />

      <!-- Switch Field (simple boolean) -->
      <q-toggle
        v-else-if="field.type === 'switch'"
        v-model="localValues[field.key]"
        color="primary"
        :label="field.placeholder"
        :disable="field.disabled"
        @update:model-value="(val) => handleFieldChange(field.key, val)"
      />

      <!-- Autocomplete Field -->
      <q-select
        v-else-if="field.type === 'autocomplete'"
        v-model="localValues[field.key]"
        outlined
        dense
        use-input
        hide-selected
        fill-input
        input-debounce="300"
        emit-value
        map-options
        clearable
        class="filter-input"
        :options="autocompleteOptions[field.key] || []"
        :loading="autocompleteLoading[field.key]"
        :label="field.placeholder || 'Search...'"
        :disable="field.disabled"
        option-value="id"
        option-label="label"
        @filter="(val, update, abort) => handleAutocomplete(field, val, update, abort)"
        @update:model-value="(val) => handleAutocompleteSelect(field.key, val)"
      >
        <template #prepend>
          <q-icon :name="field.icon" color="primary" />
        </template>
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey">
              {{ t.autocomplete.noOption.value }}
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

      <!-- Input Field -->
      <q-input
        v-else-if="field.type === 'input'"
        v-model="localValues[field.key]"
        outlined
        dense
        clearable
        class="filter-input"
        :type="field.inputType || 'text'"
        :placeholder="field.placeholder"
        :disable="field.disabled"
        @update:model-value="(val) => handleFieldChange(field.key, val)"
      >
        <template #prepend>
          <q-icon :name="field.icon" color="primary" />
        </template>
      </q-input>

      <!-- Select Field -->
      <q-select
        v-else-if="field.type === 'select'"
        v-model="localValues[field.key]"
        outlined
        dense
        emit-value
        map-options
        clearable
        class="filter-input"
        :options="field.options"
        :label="field.placeholder"
        :disable="field.disabled"
        :loading="field.loading"
        @update:model-value="(val) => handleFieldChange(field.key, val)"
      >
        <template #prepend>
          <q-icon :name="field.icon" color="primary" />
        </template>
      </q-select>

    </div>

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
        :color="hasPendingChanges ? 'warning' : 'primary'"
        icon="check"
        class="col"
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
  name: 'AdvancedFiltersDrawer'
});

/** TYPE IMPORTS */
import type {
  AdvancedFiltersDrawerProps,
  AdvancedFiltersDrawerEmits,
  FilterField,
  FilterAutocompleteOption,
  FilterValues,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { GenericDrawer } from '../genericDrawer';

/** COMPOSABLES */
import { useAdvancedFiltersDrawerTranslations } from '@composables/i18n/components/drawers/advancedFiltersDrawer';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AdvancedFiltersDrawerProps>(), {
  width: 380,
});
const emit = defineEmits<AdvancedFiltersDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useAdvancedFiltersDrawerTranslations();

/** COMPUTED */

/**
 * Drawer title - uses prop if provided, otherwise uses translation
 */
const drawerTitle = computed(() => props.title || t.title.value);

/** STATE */
const localValues = ref<FilterValues>({});
const appliedValuesSnapshot = ref<FilterValues>({});
const autocompleteOptions = ref<Record<string, FilterAutocompleteOption[]>>({});
const autocompleteLoading = ref<Record<string, boolean>>({});
const autocompleteLabels = ref<Record<string, string>>({});

/**
 * Check if there are pending (unapplied) changes
 * Compares current local values with the snapshot taken when drawer opened
 */
const hasPendingChanges = computed(() => {
  return JSON.stringify(localValues.value) !== JSON.stringify(appliedValuesSnapshot.value);
});

/** FUNCTIONS */

/**
 * Handle autocomplete search
 * @param {FilterField} field - The filter field config
 * @param {string} searchTerm - Search term from input
 * @param {Function} update - Callback to update options
 * @param {Function} abort - Callback to abort
 */
async function handleAutocomplete(
  field: FilterField,
  searchTerm: string,
  update: (callback: () => void) => void,
  abort: () => void
): Promise<void> {
  if (field.type !== 'autocomplete') {
    abort();
    return;
  }

  if (searchTerm.length < 2) {
    abort();
    return;
  }

  autocompleteLoading.value[field.key] = true;

  try {
    const options = await field.fetchOptions(searchTerm);
    update(() => {
      autocompleteOptions.value[field.key] = options;
      autocompleteLoading.value[field.key] = false;
    });
  } catch {
    autocompleteLoading.value[field.key] = false;
    abort();
  }
}

/**
 * Handle field value change - emits event for cascading filters
 * @param {string} key - Field key
 * @param {any} value - New value
 */
function handleFieldChange(key: string, value: any): void {
  emit('field-change', key, value);
}

/**
 * Handle autocomplete selection
 * @param {string} key - Field key
 * @param {string | null} value - Selected value
 */
function handleAutocompleteSelect(key: string, value: string | null): void {
  if (value) {
    const options = autocompleteOptions.value[key] || [];
    const selected = options.find(opt => opt.id === value);
    if (selected) {
      autocompleteLabels.value[key] = selected.label;
    }
  } else {
    delete autocompleteLabels.value[key];
  }
  // Emit field change for cascading filters
  emit('field-change', key, value);
}

/**
 * Handle apply button click
 */
function handleApply(): void {
  // Build values with labels for autocomplete fields
  const values: FilterValues = { ...localValues.value };

  // Add labels for autocomplete fields
  for (const key of Object.keys(autocompleteLabels.value)) {
    if (values[key]) {
      values[`${key}Label`] = autocompleteLabels.value[key];
    }
  }

  emit('apply', values);
  emit('update:modelValue', false);
}

/**
 * Handle reset button click
 */
function handleReset(): void {
  // Reset all local values to null
  for (const field of props.fields) {
    localValues.value[field.key] = null;
  }
  autocompleteLabels.value = {};
  emit('reset');
}

/**
 * Initialize local values from props
 */
function initializeValues(): void {
  const values: FilterValues = {};
  for (const field of props.fields) {
    values[field.key] = props.values[field.key] ?? null;
    // Also restore label if available
    const labelKey = `${field.key}Label`;
    if (props.values[labelKey]) {
      autocompleteLabels.value[field.key] = props.values[labelKey];
    }
  }
  localValues.value = values;
}

/** WATCHERS */

// Watch for values changes from parent
watch(() => props.values, () => {
  initializeValues();
}, { deep: true, immediate: true });

// Watch for drawer open - take snapshot of applied values
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    initializeValues();
    // Take snapshot of current applied values for comparison
    appliedValuesSnapshot.value = JSON.parse(JSON.stringify(localValues.value));
  }
});

// Emit pending state changes
watch(hasPendingChanges, (pending) => {
  emit('pending-change', pending);
});

/** FUNCTIONS - KEYBOARD */

/**
 * Handle Enter key to apply filters when drawer is open
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEnterKey(event: KeyboardEvent): void {
  if (event.key === 'Enter' && props.modelValue) {
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
// Filter field group
.filter-field {
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--mapex-divider);

  &--last {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 0;
  }

  // Field label: neutral grey text, icon stays primary (green)
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

// Toggle button styling
.filter-toggle {
  :deep(.q-btn) {
    font-weight: 500;
    padding: 6px 16px;
  }

  :deep(.q-btn--active) {
    font-weight: 600;
  }
}

// Filter input styling
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }

  :deep(.q-field__label) {
    font-weight: 400;
  }
}
</style>

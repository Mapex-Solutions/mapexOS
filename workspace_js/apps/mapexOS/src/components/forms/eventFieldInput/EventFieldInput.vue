<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { EventFieldInputProps, EventFieldInputEmits, FieldValue } from './interfaces';

defineOptions({
  name: 'EventFieldInput'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useEventFieldInputTranslations } from '@composables/i18n/components/forms/useEventFieldInputTranslations';

/** LOCAL IMPORTS (constants and handlers ONLY - NO types here!) */
import { DEFAULT_FIELD_VALUE } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<EventFieldInputProps>(), {
  modelValue: () => DEFAULT_FIELD_VALUE,
  label: 'Value',
  placeholder: 'Enter value',
  disabled: false,
  showTypeSelector: true,
  hasTemplates: false,
  templateCount: 0,
  hasStateFields: false,
  stateFieldCount: 0,
  stateFields: () => [],
});

const emit = defineEmits<EventFieldInputEmits>();

/** COMPOSABLES & STORES */
const t = useEventFieldInputTranslations();

/** COMPUTED */

/**
 * Field type options with translated labels
 */
const fieldTypeOptions = computed(() => [
  { label: t.typeOptions.event.value, value: 'event', icon: 'event', color: 'blue-6' },
  { label: t.typeOptions.state.value, value: 'state', icon: 'storage', color: 'purple-6' },
  { label: t.typeOptions.variable.value, value: 'variable', icon: 'code', color: 'orange-6' },
  { label: t.typeOptions.literal.value, value: 'literal', icon: 'format_quote', color: 'green-6' }
]);

/**
 * Current field value with fallback to default
 */
const fieldValue = computed({
  get: () => props.modelValue || DEFAULT_FIELD_VALUE,
  set: (value: FieldValue) => emit('update:modelValue', value),
});

/**
 * Current field input mode (type)
 */
const fieldInputMode = computed({
  get: () => fieldValue.value.type,
  set: (type: 'event' | 'state' | 'variable' | 'literal') => {
    const newValue: FieldValue = {
      ...fieldValue.value,
      type,
      value: type === 'literal' ? fieldValue.value.value : ''
    };

    // Set mode only for 'event' type, otherwise remove it
    if (type === 'event') {
      newValue.mode = fieldValue.value.mode || 'dynamic';
    } else {
      delete newValue.mode;
    }

    fieldValue.value = newValue;
  }
});

/**
 * Current mode option details (icon, color, etc.)
 */
const currentModeOption = computed(() => {
  return fieldTypeOptions.value.find(opt => opt.value === fieldInputMode.value) || fieldTypeOptions.value[3];
});

/**
 * Computed placeholder based on field type and mode
 */
const computedPlaceholder = computed(() => {
  if (fieldValue.value.type === 'state') {
    return t.placeholders.selectOrTypeState.value;
  }
  return props.placeholder;
});

/**
 * Computed label based on field type
 */
const computedLabel = computed(() => {
  if (fieldValue.value.type === 'event') {
    return t.labels.eventField.value;
  }
  if (fieldValue.value.type === 'state') {
    return t.labels.stateName.value;
  }
  return props.label;
});

/**
 * Available state field names for autocomplete
 */
const availableStateNames = computed(() => {
  return props.stateFields.map((field) => field.name);
});

/**
 * Computed for the value string (for literal and variable types)
 * This ensures that changes to the value trigger the parent update
 */
const inputValue = computed({
  get: () => fieldValue.value.value,
  set: (newValue: string) => {
    fieldValue.value = {
      ...fieldValue.value,
      value: newValue
    };
  }
});

/** FUNCTIONS */

/**
 * Get current input mode (dynamic or manual)
 * Defaults to 'dynamic' if not specified
 */
const currentMode = computed(() => {
  if (fieldValue.value.type !== 'event') return null;
  return fieldValue.value.mode || 'dynamic';
});

/**
 * Check if input should be readonly
 * Only readonly when type='event' and mode='dynamic'
 */
const isReadonly = computed(() => {
  return fieldValue.value.type === 'event' && currentMode.value === 'dynamic';
});

/**
 * Handle input click - opens appropriate selector
 * Only works in dynamic mode
 * If no templates selected → opens template selector
 * If templates selected → opens event field selector
 */
function handleInputClick(): void {
  // Only open drawers in dynamic mode
  if (currentMode.value !== 'dynamic') return;

  if (!props.hasTemplates) {
    emit('openTemplateSelector');
  } else {
    emit('openEventSelector');
  }
}

/**
 * Switch to manual mode
 */
function switchToManualMode(): void {
  fieldValue.value = {
    ...fieldValue.value,
    mode: 'manual'
  };
}

/**
 * Switch to dynamic mode
 */
function switchToDynamicMode(): void {
  fieldValue.value = {
    ...fieldValue.value,
    mode: 'dynamic'
  };
}



</script>

<template>
  <div class="row q-col-gutter-sm">
    <!-- Type Selector (if enabled) -->
    <div v-if="showTypeSelector" class="col-auto">
      <q-select
        v-model="fieldInputMode"
        :options="fieldTypeOptions"
        outlined
        dense
        style="min-width: 140px"
        option-label="label"
        option-value="value"
        emit-value
        map-options
        options-dense
        :disable="disabled"
      >
        <template #prepend>
          <q-icon v-if="currentModeOption" :name="currentModeOption.icon" :color="currentModeOption.color" />
        </template>
        <template #option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section avatar>
              <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ scope.opt.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <!-- Value Input/Select -->
    <div class="col">
      <!-- STATE Type - Use q-select with autocomplete -->
      <q-select
        v-if="fieldInputMode === 'state'"
        v-model="inputValue"
        outlined
        dense
        use-input
        hide-selected
        fill-input
        input-debounce="0"
        :label="computedLabel"
        :placeholder="computedPlaceholder"
        :options="availableStateNames"
        :disable="disabled"
      >
        <template #prepend>
          <q-icon name="storage" color="purple-6" />
        </template>
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey-7 text-center">
              <q-item-label caption>{{ t.empty.noStates.value }}</q-item-label>
              <q-item-label caption class="q-mt-xs">{{ t.empty.goToLocalState.value }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
        <template #append>
          <q-icon
            v-if="inputValue"
            name="close"
            class="cursor-pointer"
            @click.stop="inputValue = ''"
          >
            <AppTooltip :content="t.tooltips.clear.value" />
          </q-icon>
        </template>
      </q-select>

      <!-- EVENT Type - Input that can be readonly (dynamic) or editable (manual) -->
      <q-input
        v-else-if="fieldValue.type === 'event'"
        v-model="inputValue"
        outlined
        dense
        :readonly="isReadonly"
        :label="computedLabel"
        :placeholder="computedPlaceholder"
        :class="{ 'cursor-pointer': isReadonly }"
        :disable="disabled"
        @click="handleInputClick"
      >
        <!-- Type Icon Prepend -->
        <template #prepend>
          <q-icon name="event" color="blue-6" />
        </template>

        <!-- Append Slot - Template Badge & Menu -->
        <template #append>
          <q-chip
            v-if="hasTemplates"
            dense
            size="sm"
            color="blue-1"
            text-color="blue-9"
          >
            <q-icon name="folder" size="xs" class="q-mr-xs" />
            {{ templateCount }}
            <AppTooltip :content="t.tooltips.templatesSelected.value.replace('{count}', (templateCount || 0).toString())" />
          </q-chip>

          <!-- Options Menu (Kebab) -->
          <q-btn
            flat
            round
            dense
            size="sm"
            icon="more_vert"
            color="grey-7"
            class="q-ml-xs"
            @click.stop
          >
            <AppTooltip :content="t.tooltips.options.value" />

            <q-menu
              anchor="bottom end"
              self="top end"
              :offset="[0, 4]"
            >
              <q-list dense style="min-width: 240px; padding: 8px 0;">
                <!-- DYNAMIC MODE OPTIONS -->
                <template v-if="currentMode === 'dynamic'">
                  <!-- Select Field (when templates are selected) -->
                  <q-item
                    v-if="hasTemplates"
                    clickable
                    v-close-popup
                    class="q-py-sm q-px-md context-menu-item"
                    @click="emit('openEventSelector')"
                  >
                    <q-item-section avatar>
                      <q-icon name="search" color="primary" size="sm" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ t.menu.selectField.title.value }}</q-item-label>
                      <q-item-label caption class="text-grey-7">
                        {{ t.menu.selectField.description.value.replace('{count}', (templateCount || 0).toString()) }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>

                  <q-separator v-if="hasTemplates" class="q-my-xs" />

                  <!-- Search Templates -->
                  <q-item
                    clickable
                    v-close-popup
                    class="q-py-sm q-px-md context-menu-item"
                    @click="emit('openTemplateSelector')"
                  >
                    <q-item-section avatar>
                      <q-icon name="add" color="primary" size="sm" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ t.menu.searchTemplates.title.value }}</q-item-label>
                      <q-item-label caption class="text-grey-7">
                        {{ t.menu.searchTemplates.description.value }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>

                  <q-separator class="q-my-xs" />

                  <!-- Switch to Manual Mode -->
                  <q-item
                    clickable
                    v-close-popup
                    class="q-py-sm q-px-md context-menu-item"
                    @click="switchToManualMode"
                  >
                    <q-item-section avatar>
                      <q-icon name="edit" color="orange-6" size="sm" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ t.menu.manualMode.title.value }}</q-item-label>
                      <q-item-label caption class="text-grey-7">
                        {{ t.menu.manualMode.description.value }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </template>

                <!-- MANUAL MODE OPTIONS -->
                <template v-else-if="currentMode === 'manual'">
                  <!-- Switch to Dynamic Mode -->
                  <q-item
                    clickable
                    v-close-popup
                    class="q-py-sm q-px-md context-menu-item"
                    @click="switchToDynamicMode"
                  >
                    <q-item-section avatar>
                      <q-icon name="folder" color="primary" size="sm" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ t.menu.dynamicMode.title.value }}</q-item-label>
                      <q-item-label caption class="text-grey-7">
                        {{ t.menu.dynamicMode.description.value }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </template>
              </q-list>
            </q-menu>
          </q-btn>
        </template>
      </q-input>

      <!-- OTHER types (literal, variable) - Use q-input -->
      <q-input
        v-else
        v-model="inputValue"
        outlined
        dense
        :label="computedLabel"
        :placeholder="computedPlaceholder"
        :disable="disabled"
      >
        <!-- Type Icon Prepend -->
        <template #prepend>
          <q-icon v-if="currentModeOption" :name="currentModeOption.icon" :color="currentModeOption.color" />
        </template>
      </q-input>
    </div>
  </div>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>

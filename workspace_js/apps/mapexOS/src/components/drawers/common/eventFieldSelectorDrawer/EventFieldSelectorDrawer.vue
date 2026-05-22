<script setup lang="ts">
defineOptions({
  name: 'EventFieldSelectorDrawer'
});

/** TYPE IMPORTS */
import type {
  EventFieldSelectorDrawerProps,
  EventFieldSelectorDrawerEmits,
  FieldInfo,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { GenericDrawer } from '@components/drawers/common/genericDrawer';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';

/** Stub for removed rule page state — template field operations are no longer available */
function useRuleAddPageState() {
  return {
    getTemplateFields: (id: string): string[] => { void id; return []; },
    getTemplateName: (id: string): string => { void id; return ''; },
    removeTemplate: (id: string): void => { void id; },
  };
}

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { SEARCH_DEBOUNCE_DELAY, EMPTY_STATE_MESSAGES, LOADING_MESSAGES } from './constants';

/** PROPS & EMITS */
const props = defineProps<EventFieldSelectorDrawerProps>();
const emit = defineEmits<EventFieldSelectorDrawerEmits>();

/** COMPOSABLES & STORES */
const { getTemplateFields, getTemplateName, removeTemplate } = useRuleAddPageState();
const logger = useLogger('EventFieldSelectorDrawer');

/** STATE */

/**
 * Loading state indicator
 */
const loading = ref(false);

/**
 * Search query string
 */
const searchQuery = ref('');

/**
 * Array of ALL fields from ALL templates (NO deduplication)
 * Each field maintains its template association
 */
const allFields = ref<FieldInfo[]>([]);

/**
 * Currently selected template ID for viewing
 */
const selectedTemplateForView = ref<string | undefined>(undefined);

/**
 * Template names cache
 * Key: templateId, Value: template name
 */
const templateNames = ref<Map<string, string>>(new Map());

/** COMPUTED */

/**
 * Options for template selector dropdown
 * @returns {Array} Array of template options with id, name, and value
 */
const selectedTemplateOptions = computed(() => {
  const uniqueTemplates = new Map<string, string>();

  allFields.value.forEach((fieldInfo) => {
    if (!uniqueTemplates.has(fieldInfo.templateId)) {
      uniqueTemplates.set(fieldInfo.templateId, fieldInfo.templateName);
    }
  });

  return Array.from(uniqueTemplates.entries())
    .map(([id, name]) => ({
      id,
      name,
      value: id,
    }))
    .sort((a, b) => a.name.localeCompare(b.name));
});

/**
 * Get name of currently selected template
 * @returns {string} Template name or empty string
 */
const selectedTemplateName = computed(() => {
  if (!selectedTemplateForView.value) return '';
  return templateNames.value.get(selectedTemplateForView.value) || '';
});

/**
 * Filter fields for the currently selected template only
 * @returns {string[]} Array of field paths for selected template
 */
const currentTemplateFields = computed((): string[] => {
  if (!selectedTemplateForView.value) return [];

  const query = searchQuery.value.toLowerCase().trim();

  return allFields.value
    .filter((fieldInfo) => fieldInfo.templateId === selectedTemplateForView.value)
    .map((fieldInfo) => fieldInfo.path)
    .filter((path) => query === '' ? true : path.toLowerCase().includes(query));
});

/**
 * Check if any fields are available
 * @returns {boolean} True if there are fields
 */
const hasFields = computed(() => allFields.value.length > 0);

/**
 * Check if result set is empty
 * @returns {boolean} True if no results
 */
const isEmpty = computed(() => {
  if (!selectedTemplateForView.value) return true;
  return currentTemplateFields.value.length === 0;
});

/**
 * Get contextual empty message
 * @returns {string} Empty state message
 */
const emptyMessage = computed(() => {
  if (props.selectedTemplates.length === 0) {
    return EMPTY_STATE_MESSAGES.NO_TEMPLATES;
  }
  if (!hasFields.value) {
    return EMPTY_STATE_MESSAGES.NO_FIELDS;
  }
  if (!selectedTemplateForView.value) {
    return 'Please select a template from the dropdown above.';
  }
  if (searchQuery.value.trim() !== '') {
    return EMPTY_STATE_MESSAGES.NO_RESULTS;
  }
  return EMPTY_STATE_MESSAGES.NO_FIELDS;
});

/**
 * Count total visible fields for selected template
 * @returns {number} Total field count
 */
const totalFieldsCount = computed(() => currentTemplateFields.value.length);

/** WATCHERS */

/**
 * Watch drawer open state and fetch fields when opened
 */
watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen && props.selectedTemplates.length > 0) {
      fetchAvailableFields();
    }
  }
);

/**
 * Watch selected templates and refetch when they change
 */
watch(
  () => props.selectedTemplates,
  (newTemplates) => {
    if (newTemplates.length > 0 && props.modelValue) {
      fetchAvailableFields();
    }
  },
  { deep: true }
);

/**
 * Auto-select first template when templates change
 */
watch(
  () => selectedTemplateOptions.value,
  (options) => {
    if (options.length > 0 && !selectedTemplateForView.value && options[0]?.id) {
      selectedTemplateForView.value = options[0].id;
    }
  },
  { immediate: true }
);

/** FUNCTIONS */

/**
 * Fetch available fields from all selected templates (from cache)
 * Stores ALL fields without deduplication
 *
 * @returns {void}
 */
function fetchAvailableFields(): void {
  loading.value = true;
  allFields.value = [];

  try {
    for (const templateId of props.selectedTemplates) {
      const fields = getTemplateFields(templateId);
      const templateName = getTemplateName(templateId);

      fields.forEach((fieldPath: string) => {
        allFields.value.push({
          path: fieldPath,
          templateId,
          templateName,
        });
      });

      templateNames.value.set(templateId, templateName);
    }

    if (selectedTemplateOptions.value.length > 0 && selectedTemplateOptions.value[0]?.id) {
      selectedTemplateForView.value = selectedTemplateOptions.value[0].id;
    }
  } catch (error: any) {
    logger.error('Failed to fetch available fields:', error);
  } finally {
    loading.value = false;
  }
}

/**
 * Handle search input with debounce
 *
 * @param {string | number | null} query - Search query
 * @returns {void}
 */
let searchTimeout: NodeJS.Timeout | undefined;
function handleSearch(query: string | number | null): void {
  if (searchTimeout) {
    clearTimeout(searchTimeout);
  }

  searchTimeout = setTimeout(() => {
    searchQuery.value = query ? String(query) : '';
  }, SEARCH_DEBOUNCE_DELAY);
}

/**
 * Handle field selection
 *
 * @param {string} field - Selected field path
 * @returns {void}
 */
function handleFieldSelect(field: string): void {
  emit('select', field);
  emit('update:modelValue', false);
}

/**
 * Open template manager
 *
 * @returns {void}
 */
function openTemplateManager(): void {
  emit('manage-templates');
}

/**
 * Handle template removal
 *
 * @param {string} templateId - Template ID to remove
 * @returns {void}
 */
function handleRemoveTemplate(templateId: string): void {
  removeTemplate(templateId);

  allFields.value = allFields.value.filter((fieldInfo) => fieldInfo.templateId !== templateId);
  templateNames.value.delete(templateId);

  if (selectedTemplateForView.value === templateId) {
    if (selectedTemplateOptions.value.length > 0 && selectedTemplateOptions.value[0]?.id) {
      selectedTemplateForView.value = selectedTemplateOptions.value[0].id;
    } else {
      selectedTemplateForView.value = undefined;
    }
  }
}

/** LIFECYCLE HOOKS */

onBeforeUnmount(() => {
  if (searchTimeout) {
    clearTimeout(searchTimeout);
  }
});
</script>

<template>
  <GenericDrawer
    :model-value="modelValue"
    title="Select Event Field"
    icon="search"
    :width="600"
    @update:model-value="emit('update:modelValue', $event)"
    @close="emit('update:modelValue', false)"
  >
    <!-- Search Box -->
    <q-input
      outlined
      dense
      clearable
      placeholder="Search fields..."
      :model-value="searchQuery"
      class="q-mb-md"
      @update:model-value="handleSearch"
    >
      <template #prepend>
        <q-icon name="search" color="grey-7" />
      </template>
    </q-input>

    <!-- Template Selector Dropdown -->
    <div v-if="selectedTemplateOptions.length > 0" class="q-mb-md">
      <div class="text-caption text-grey-7 q-mb-sm">Select Asset Template:</div>
      <q-select
        v-model="selectedTemplateForView"
        outlined
        dense
        label="Asset Template"
        :options="selectedTemplateOptions"
        option-label="name"
        option-value="id"
        emit-value
        map-options
      >
        <template #prepend>
          <q-icon name="memory" color="primary" />
        </template>
        <template #append>
          <q-btn
            v-if="selectedTemplateForView"
            flat
            round
            dense
            size="sm"
            icon="close"
            color="grey-7"
            @click.stop="handleRemoveTemplate(selectedTemplateForView)"
          />
        </template>
      </q-select>
    </div>

    <q-separator v-if="selectedTemplateOptions.length > 0" class="q-mb-md" />

    <!-- Loading State -->
    <div v-if="loading" class="q-pa-lg text-center">
      <q-spinner color="primary" size="lg" />
      <div class="text-caption text-grey-7 q-mt-md">
        {{ LOADING_MESSAGES.FETCHING }}
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="isEmpty" class="q-pa-lg text-center">
      <q-icon name="inbox" size="xl" color="grey-5" />
      <div class="text-body2 text-grey-7 q-mt-md">
        {{ emptyMessage }}
      </div>
      <q-btn
        v-if="props.selectedTemplates.length === 0"
        flat
        color="primary"
        label="Select Templates"
        icon="add"
        class="q-mt-md"
        @click="openTemplateManager"
      />
    </div>

    <!-- Fields List for Selected Template -->
    <div v-else>
      <!-- Fields Header -->
      <div class="row items-center q-mb-md">
        <q-icon name="list_alt" color="primary" size="sm" class="q-mr-xs" />
        <div class="text-subtitle2 text-grey-8">
          Fields from {{ selectedTemplateName }}
        </div>
      </div>

      <!-- Fields List -->
      <q-list separator class="field-list">
        <q-item
          v-for="(field, index) in currentTemplateFields"
          :key="`${selectedTemplateForView}-${field}-${index}`"
          clickable
          :active="field === props.currentValue"
          @click="handleFieldSelect(field)"
        >
          <q-item-section>
            <q-item-label class="text-code">{{ field }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-btn
              flat
              dense
              size="sm"
              color="primary"
              label="SELECT"
              @click.stop="handleFieldSelect(field)"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- Footer -->
    <template #footer>
      <div class="text-caption text-grey-7">
        {{ totalFieldsCount }} field{{ totalFieldsCount !== 1 ? 's' : '' }} available
        <template v-if="selectedTemplateName"> from {{ selectedTemplateName }}</template>
      </div>
      <q-space />
      <q-btn
        flat
        dense
        no-caps
        label="Manage Templates"
        icon="settings"
        size="sm"
        color="primary"
        @click="openTemplateManager"
      />
    </template>
  </GenericDrawer>
</template>

<style lang="scss" scoped>
.text-code {
  font-family: 'Roboto Mono', monospace;
  font-size: 0.875rem;
  color: var(--mapex-text-primary);
}

// Field list hover effects
.field-list {
  :deep(.q-item) {
    transition: all var(--mapex-transition-base) ease;
    border-radius: var(--mapex-radius-sm);
  }

  :deep(.q-item:hover) {
    background-color: var(--mapex-surface-bg);
  }

  :deep(.q-item.q-item--active) {
    background-color: rgba(var(--q-primary-rgb), 0.08) !important;
    border-left: 3px solid var(--q-primary);
  }
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'TriggerSelectorDialog'
});

/** TYPE IMPORTS */
import type { TriggerResponse } from '@mapexos/schemas';
import type { TriggerSelectorDialogProps, TriggerSelectorDialogEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** PROPS & EMITS */
const props = withDefaults(defineProps<TriggerSelectorDialogProps>(), {
  selectedTriggerId: null,
});

const emit = defineEmits<TriggerSelectorDialogEmits>();

/** COMPOSABLES & STORES */
const tsTitle = useTS({ titleCase: true });
const ts = useTS({ capitalize: true });
const tsRaw = useTS({ capitalize: false });
const bp = 'components.dialogs.triggerSelector';

/** STATE */

/**
 * All triggers fetched from API
 */
const triggers = ref<TriggerResponse[]>([]);

/**
 * Loading state
 */
const loading = ref(false);

/**
 * Client-side search query (filtering done locally)
 */
const searchQuery = ref('');

/**
 * Selected category filter
 */
const selectedCategory = ref<string | undefined>(undefined);

/** COMPUTED */

/**
 * Category options extracted from fetched triggers
 */
const categoryOptions = computed(() => {
  const cats = new Set<string>();
  triggers.value.forEach(t => {
    if (t.category) cats.add(t.category);
  });
  return [
    { label: tsTitle(`${bp}.allFilter`), value: undefined },
    ...Array.from(cats).sort().map(c => ({ label: c.charAt(0).toUpperCase() + c.slice(1), value: c })),
  ];
});

/**
 * Filtered triggers (client-side search + category filter)
 */
const filteredTriggers = computed(() => {
  let result = triggers.value;

  if (selectedCategory.value) {
    result = result.filter(t => t.category === selectedCategory.value);
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(t =>
      t.name?.toLowerCase().includes(q) ||
      t.description?.toLowerCase().includes(q),
    );
  }

  return result;
});

/**
 * Pre-selected trigger IDs for highlighting
 */
const selectedIds = computed(() =>
  props.selectedTriggerId ? [props.selectedTriggerId] : [],
);

/** WATCHERS */

/**
 * Fetch triggers when dialog opens
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && triggers.value.length === 0) {
    void fetchTriggers();
  }
});

/** FUNCTIONS */

/**
 * Fetch triggers from API
 *
 * @returns {Promise<void>}
 */
async function fetchTriggers(): Promise<void> {
  loading.value = true;
  try {
    const response = await apis.triggers?.trigger.list({
      page: 1,
      pageSize: 100,
    });
    if (response?.items) {
      triggers.value = response.items;
    }
  } catch (error) {
    handleApiError(error, {
      defaultMessage: ts(`${bp}.fetchError`),
    });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle trigger selection from GenericSelectorDialog
 *
 * @param {any[]} items - Selected items (single-select, array of 1)
 */
function handleSelect(items: any[]): void {
  const trigger = items[0] as TriggerResponse;
  if (trigger) {
    emit('select', trigger);
  }
}

/**
 * Handle search query (client-side filtering)
 *
 * @param {string} query - Search query
 */
function handleSearch(query: string): void {
  searchQuery.value = query;
}

/**
 * Handle category filter change
 */
function onCategoryChange(): void {
  // Client-side filtering — no API call needed
}

/**
 * Get icon for trigger category
 *
 * @param {string} category - Trigger category
 * @returns {string} Material icon name
 */
function getCategoryIcon(category: string): string {
  const icons: Record<string, string> = {
    email: 'email',
    slack: 'chat',
    teams: 'groups',
    http: 'http',
    mqtt: 'router',
    rabbitmq: 'cloud_queue',
    nats: 'cloud',
    websocket: 'cable',
    custom: 'code',
  };
  return icons[category] || 'notifications';
}

/**
 * Get color for trigger category
 *
 * @param {string} category - Trigger category
 * @returns {string} Quasar color name
 */
function getCategoryColor(category: string): string {
  const colors: Record<string, string> = {
    email: 'blue',
    slack: 'purple',
    teams: 'indigo',
    http: 'green',
    mqtt: 'orange',
    rabbitmq: 'deep-orange',
    nats: 'cyan',
    websocket: 'teal',
    custom: 'grey',
  };
  return colors[category] || 'primary';
}
</script>

<template>
  <GenericSelectorDialog
    :model-value="modelValue"
    :title="tsTitle(`${bp}.title`)"
    icon="notifications_active"
    icon-color="warning"
    :items="filteredTriggers"
    item-key="id"
    :multi-select="false"
    :selected-ids="selectedIds"
    :loading="loading"
    :search-placeholder="tsRaw(`${bp}.searchPlaceholder`)"
    :info-banner="{ text: ts(`${bp}.infoBanner`) }"
    :empty-text="ts(`${bp}.emptyText`)"
    empty-icon="inbox"
    results-icon="notifications_active"
    footer-icon="notifications_active"
    :item-noun-singular="tsRaw(`${bp}.itemSingular`)"
    :item-noun-plural="tsRaw(`${bp}.itemPlural`)"
    :total-items="filteredTriggers.length"
    :active-item-style="{ backgroundColor: 'rgba(255, 152, 0, 0.08)', borderColor: 'var(--q-warning)' }"
    @update:model-value="emit('update:modelValue', $event)"
    @select="handleSelect"
    @cancel="emit('cancel')"
    @search="handleSearch"
  >
    <!-- Category filter -->
    <template #filters>
      <div class="col-12">
        <q-select
          v-model="selectedCategory"
          outlined
          dense
          clearable
          :label="tsTitle(`${bp}.category`)"
          :options="categoryOptions"
          option-label="label"
          option-value="value"
          emit-value
          map-options
          @update:model-value="onCategoryChange"
        >
          <template #prepend>
            <q-icon name="category" />
          </template>
        </q-select>
      </div>
    </template>

    <!-- Item rendering -->
    <template #item="{ item }">
      <q-item-section avatar>
        <q-avatar
          :icon="getCategoryIcon(item.category || 'custom')"
          :color="getCategoryColor(item.category || 'custom')"
          text-color="white"
          size="md"
        />
      </q-item-section>
      <q-item-section>
        <q-item-label class="text-weight-medium">{{ item.name }}</q-item-label>
        <q-item-label caption lines="2">{{ item.description }}</q-item-label>
        <q-item-label caption class="q-mt-xs">
          <q-badge
            :color="getCategoryColor(item.category || 'custom')"
            :label="item.category"
            dense
          />
        </q-item-label>
      </q-item-section>
    </template>
  </GenericSelectorDialog>
</template>

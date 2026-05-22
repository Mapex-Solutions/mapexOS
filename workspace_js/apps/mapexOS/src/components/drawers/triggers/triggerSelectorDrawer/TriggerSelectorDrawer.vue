<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { TriggerResponse } from '@mapexos/schemas';

defineOptions({
  name: 'TriggerSelectorDrawer'
});

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** PROPS & EMITS */
const props = withDefaults(defineProps<{
  /** Whether the drawer is open */
  modelValue: boolean;
  /** Pre-selected trigger ID for highlighting */
  selectedTriggerId?: string | null;
}>(), {
  selectedTriggerId: null,
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  'select': [trigger: TriggerResponse];
  'cancel': [];
}>();

/** STATE */
const loading = ref(false);
const searchQuery = ref('');
const triggers = ref<TriggerResponse[]>([]);
const selectedCategory = ref<string>('all');

/** COMPUTED */
const showDrawer = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

const categories = computed(() => {
  const cats = new Set<string>();
  cats.add('all');
  triggers.value.forEach(trigger => {
    if (trigger.category) {
      cats.add(trigger.category);
    }
  });
  return Array.from(cats);
});

const filteredTriggers = computed(() => {
  let result = triggers.value;

  // Filter by category
  if (selectedCategory.value !== 'all') {
    result = result.filter(t => t.category?.toString() === selectedCategory.value);
  }

  // Filter by search query
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(t =>
      t.name?.toLowerCase().includes(query) ||
      t.description?.toLowerCase().includes(query)
    );
  }

  return result;
});

/** FUNCTIONS */
async function fetchTriggers(): Promise<void> {
  loading.value = true;
  try {
    const response = await apis.triggers?.trigger.list({
      page: 1,
      pageSize: 100
    });

    if (response?.items) {
      triggers.value = response.items;
    }
  } catch (error: any) {
    handleApiError({ error, customMessage: 'Failed to fetch triggers' });
  } finally {
    loading.value = false;
  }
}

function handleTriggerSelect(trigger: TriggerResponse): void {
  emit('select', trigger);
  close();
}

function close(): void {
  emit('update:modelValue', false);
}

function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

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

/**
 * Check if trigger is currently selected
 * @param {TriggerResponse} trigger - Trigger to check
 * @returns {boolean} True if selected
 */
function isSelected(trigger: TriggerResponse): boolean {
  return trigger.id === props.selectedTriggerId;
}

/** WATCHERS */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && triggers.value.length === 0) {
    void fetchTriggers();
  }
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<template>
  <q-drawer
    v-model="showDrawer"
    side="right"
    overlay
    elevated
    :width="700"
    class="trigger-selector-drawer"
  >
    <!-- Header -->
    <q-toolbar class="bg-grey-2">
      <q-icon
        name="notifications_active"
        size="md"
        color="primary"
        class="q-mr-sm"
      />
      <q-toolbar-title class="text-weight-medium text-grey-9">
        Select Trigger
      </q-toolbar-title>
      <q-btn
        flat
        round
        dense
        icon="close"
        color="grey-7"
        @click="close"
      >
        <AppTooltip content="Close" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Search Box -->
    <div class="q-pa-md">
      <q-input
        v-model="searchQuery"
        outlined
        dense
        clearable
        placeholder="Search triggers..."
      >
        <template #prepend>
          <q-icon name="search" color="grey-7" />
        </template>
      </q-input>
    </div>

    <q-separator />

    <!-- Category Filter -->
    <div class="q-pa-md">
      <div class="text-caption text-grey-7 q-mb-sm">Filter by category:</div>
      <q-btn-toggle
        v-model="selectedCategory"
        toggle-color="primary"
        :options="categories.map(cat => ({ label: cat.toUpperCase(), value: cat }))"
        size="sm"
        unelevated
        class="category-toggle"
      />
    </div>

    <q-separator />

    <!-- Loading State -->
    <div v-if="loading" class="q-pa-lg text-center">
      <q-spinner
        color="primary"
        size="lg"
      />
      <div class="text-caption text-grey-7 q-mt-md">
        Loading triggers...
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="filteredTriggers.length === 0"
      class="q-pa-lg text-center"
    >
      <q-icon
        name="inbox"
        size="xl"
        color="grey-5"
      />
      <div class="text-body2 text-grey-7 q-mt-md">
        {{ searchQuery.trim() ? 'No triggers found' : 'No triggers available' }}
      </div>
    </div>

    <!-- Triggers List -->
    <div v-else class="q-pa-md">
      <q-list separator>
        <q-item
          v-for="trigger in filteredTriggers"
          :key="trigger.id || trigger.name || 'unknown'"
          clickable
          :active="isSelected(trigger)"
          @click="handleTriggerSelect(trigger)"
          class="trigger-item"
        >
          <q-item-section avatar>
            <q-avatar
              :icon="getCategoryIcon(trigger.category || 'custom')"
              :color="getCategoryColor(trigger.category || 'custom')"
              text-color="white"
              size="md"
            />
          </q-item-section>

          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ trigger.name }}
            </q-item-label>
            <q-item-label caption lines="2">
              {{ trigger.description }}
            </q-item-label>
            <q-item-label caption class="q-mt-xs">
              <DetailChip
                :value="trigger.category"
                :color="getCategoryColor(trigger.category || 'custom') as any"
                size="sm"
                dense
              />
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-btn
              flat
              dense
              color="primary"
              label="SELECT"
              @click.stop="handleTriggerSelect(trigger)"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- Footer with Count -->
    <div
      v-if="!loading && filteredTriggers.length > 0"
      class="absolute-bottom bg-grey-2 q-pa-sm text-center text-caption text-grey-7"
    >
      {{ filteredTriggers.length }} trigger{{ filteredTriggers.length !== 1 ? 's' : '' }} available
    </div>
  </q-drawer>
</template>

<style scoped lang="scss">
.trigger-selector-drawer {
  max-width: 100vw;
}

.trigger-item {
  transition: all var(--mapex-transition-base) ease;

  &:hover {
    background-color: var(--mapex-surface-bg);
  }
}

:deep(.q-item.q-item--active) {
  background-color: var(--mapex-active-bg) !important;
  border-left: 3px solid var(--q-warning);
}

.category-toggle {
  :deep(.q-btn) {
    font-size: 0.75rem;
  }
}
</style>

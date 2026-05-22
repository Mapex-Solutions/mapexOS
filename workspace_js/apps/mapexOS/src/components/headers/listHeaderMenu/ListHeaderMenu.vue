<template>
  <div class="row items-center no-wrap q-gutter-md">

    <!-- Relative timestamp + inline Atualizar link (Linear/GitHub pattern) -->
    <div
      v-if="props.showRefresh !== false"
      class="last-updated text-caption text-grey-7 row items-center no-wrap"
    >
      <span v-if="relativeTimeLabel" class="last-updated__time">{{ relativeTimeLabel }}</span>
      <span v-if="relativeTimeLabel" class="last-updated__separator">·</span>
      <q-btn
        flat
        dense
        no-caps
        size="sm"
        padding="2px 6px"
        color="grey-7"
        icon="refresh"
        :label="t.refresh.value"
        :disable="props.refreshing === true"
        :class="{ 'refresh-spinning': props.refreshing === true }"
        class="last-updated__refresh"
        @click="emit('refresh')"
      />
    </div>

    <!-- Stats + menu pill (unchanged) -->
    <q-btn
      outline
      rounded
      class="q-py-xs q-px-md"
      color="primary"
      :label="buttonLabel"
      :icon="props.icon || 'list'"
    >
      <q-menu>
        <q-list style="min-width: 250px">
          <!-- Items per page section -->
          <template v-if="props.showItemsPerPage !== false && itemsPerPageOptions.length > 0">
            <q-item-label header>{{ t.itemsPerPage.value }}</q-item-label>
            <q-item
              v-for="option in itemsPerPageOptions"
              v-close-popup
              clickable
              :key="option"
              @click="emit('update:itemsPerPage', option)"
            >
              <q-item-section>
                <q-item-label>{{ option }} items</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-icon v-if="props.itemsPerPage === option" name="check" color="primary" />
              </q-item-section>
            </q-item>

            <q-separator v-if="props.showColumnVisibility !== false && localColumns.length > 0" class="q-my-sm" />
          </template>

          <!-- Column visibility section -->
          <template v-if="props.showColumnVisibility !== false && localColumns.length > 0">
            <q-item-label header>{{ t.visibleColumns.value }}</q-item-label>
            <q-item
              v-for="column in localColumns"
              :key="column.key"
              tag="label"
            >
              <q-item-section>
                <q-item-label>{{ column.label }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-checkbox
                  :model-value="column.visible"
                  @update:model-value="toggleColumn(column.key, $event)"
                />
              </q-item-section>
            </q-item>
          </template>
        </q-list>
      </q-menu>
    </q-btn>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ListHeaderMenu'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuProps, ListHeaderMenuEmits, ListHeaderMenuColumn } from './interfaces';

/** VUE IMPORTS */
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPOSABLES */
import { useListHeaderMenuTranslations } from '@composables/i18n/components/headers';

/** PROPS & EMITS */
const props = withDefaults(defineProps<ListHeaderMenuProps>(), {
  itemLabelPlural: '',
  icon: 'list',
  itemsPerPageOptions: () => [10, 15, 25, 50, 100],
  columns: () => [],
  showItemsPerPage: true,
  showColumnVisibility: true,
  filtered: false,
  showRefresh: true,
  refreshing: false,
});
const emit = defineEmits<ListHeaderMenuEmits>();

/** COMPOSABLES & STORES */
const t = useListHeaderMenuTranslations();

/** CONSTANTS */
const TICK_INTERVAL_MS = 5000;

/** STATE */
const localColumns = ref<ListHeaderMenuColumn[]>([...props.columns]);
const now = ref<number>(Date.now());
let tickHandle: ReturnType<typeof setInterval> | null = null;

/** COMPUTED */

/**
 * Computed button label
 * Shows: "X of Y ITEMS" where X is items per page and Y is total
 */
const buttonLabel = computed(() => {
  const totalCount = props.itemsCount;
  const perPage = props.itemsPerPage;
  const label = totalCount === 1
    ? props.itemLabel
    : (props.itemLabelPlural || `${props.itemLabel}s`);

  const filteredSuffix = props.filtered ? ` ${t.filtered.value}` : '';

  if (perPage && totalCount > perPage) {
    return `${perPage} OF ${totalCount} ${label.toUpperCase()}${filteredSuffix}`;
  }

  return `${totalCount} ${label.toUpperCase()}${filteredSuffix}`;
});

/**
 * Relative time label for the last update.
 * Returns "" when lastUpdatedAt is not provided so the caption stays empty.
 */
const relativeTimeLabel = computed(() => {
  if (!props.lastUpdatedAt) return '';

  const updatedMs = typeof props.lastUpdatedAt === 'number'
    ? props.lastUpdatedAt
    : props.lastUpdatedAt.getTime();

  const diffSec = Math.max(0, Math.floor((now.value - updatedMs) / 1000));

  if (diffSec < 5) return t.lastUpdatedNow.value;
  if (diffSec < 60) return t.lastUpdatedSeconds(diffSec);
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) return t.lastUpdatedMinutes(diffMin);
  const diffHours = Math.floor(diffMin / 60);
  return t.lastUpdatedHours(diffHours);
});

/** WATCHERS */

watch(() => props.columns, (newColumns) => {
  localColumns.value = [...newColumns];
}, { deep: true });

/** FUNCTIONS */

/**
 * Toggle column visibility
 * @param {string} key - Column key to toggle
 * @param {boolean} visible - New visibility state
 */
function toggleColumn(key: string, visible: boolean): void {
  const column = localColumns.value.find(col => col.key === key);
  if (column) {
    column.visible = visible;
    emit('update:columns', localColumns.value);
  }
}

/** LIFECYCLE HOOKS */

onMounted(() => {
  tickHandle = setInterval(() => { now.value = Date.now(); }, TICK_INTERVAL_MS);
});

onBeforeUnmount(() => {
  if (tickHandle !== null) clearInterval(tickHandle);
});
</script>

<style scoped lang="scss">
.last-updated {
  gap: 6px;
  line-height: 1;

  &__separator {
    opacity: 0.6;
  }

  &__refresh {
    font-weight: 500;
    text-transform: none;
    letter-spacing: normal;
  }
}

.refresh-spinning :deep(.q-icon) {
  animation: refresh-spin 0.8s linear infinite;
}

@keyframes refresh-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>

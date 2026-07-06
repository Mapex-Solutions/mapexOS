<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepAssetSelectionProps, StepAssetSelectionEmits, SelectableAsset } from './interfaces';

defineOptions({
  name: 'StepAssetSelection'
});

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { ListPagination } from '@components/navigation';

/** COMPOSABLES */
import { useCreateEditMigrationPlanTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { handleApiError } from '@utils/error';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<StepAssetSelectionProps>();
const emit = defineEmits<StepAssetSelectionEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditMigrationPlanTranslations();
const logger = useLogger('StepAssetSelection');

/** CONSTANTS */
// Rows shown per page in the selection list.
const PER_PAGE = 10;
// Page size used to paginate through the whole result set when selecting all;
// the asset query caps perPage at 100, so select-all fetches in 100-id batches.
const SELECT_ALL_PER_PAGE = 100;

/** STATE */
const assets = ref<SelectableAsset[]>([]);
// Selected ids across all pages, seeded from the incoming form model.
const selectedIds = ref<Set<string>>(new Set(props.assetIds));
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const loading = ref(false);
const selectingAll = ref(false);

/** COMPUTED */
const selectedCount = computed(() => selectedIds.value.size);

// Every asset on the current page is selected.
const pageAllSelected = computed(
  () => assets.value.length > 0 && assets.value.every(a => selectedIds.value.has(a.id))
);

/** FUNCTIONS */
/**
 * Emit the current selection and its validity to the parent.
 * @returns {void}
 */
function emitSelection(): void {
  emit('update:assetIds', Array.from(selectedIds.value));
  emit('update:valid', selectedIds.value.size > 0);
}

/**
 * Toggle a single asset in the selection.
 * @param assetId - Asset id to toggle
 * @param checked - Whether the asset should be selected
 * @returns {void}
 */
function toggleAsset(assetId: string, checked: boolean): void {
  const next = new Set(selectedIds.value);
  if (checked) {
    next.add(assetId);
  } else {
    next.delete(assetId);
  }
  selectedIds.value = next;
  emitSelection();
}

/**
 * Select or clear every asset on the current page.
 * @param checked - Whether the page rows should be selected
 * @returns {void}
 */
function toggleSelectPage(checked: boolean): void {
  const next = new Set(selectedIds.value);
  assets.value.forEach(a => (checked ? next.add(a.id) : next.delete(a.id)));
  selectedIds.value = next;
  emitSelection();
}

/**
 * Clear the whole selection.
 * @returns {void}
 */
function clearSelection(): void {
  selectedIds.value = new Set();
  emitSelection();
}

/**
 * Fetch a single page of assets scoped to the source template.
 * @param page - 1-based page number
 * @returns {Promise<void>}
 */
async function fetchAssets(page: number): Promise<void> {
  if (!apis.assets || !props.fromTemplateId) {
    assets.value = [];
    return;
  }

  try {
    loading.value = true;
    const response = await apis.assets.asset.list({
      page,
      perPage: PER_PAGE,
      assetTemplateId: props.fromTemplateId,
      sort: 'name:asc',
    });

    assets.value = response.items.map(a => ({
      id: a.id || '',
      name: a.name || '',
      assetUUID: a.assetUUID || '',
    }));
    totalItems.value = response.pagination?.totalItems ?? 0;
    totalPages.value = response.pagination?.totalPages ?? 1;
  } catch (error: unknown) {
    logger.error('Error fetching assets:', error);
    handleApiError(error, { defaultMessage: t.steps.assets.loadError.value, timeout: 5000 });
    assets.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * Select every asset matching the current template filter by paginating through
 * the whole result set and collecting their ids.
 * @returns {Promise<void>}
 */
async function selectAllMatching(): Promise<void> {
  if (!apis.assets || !props.fromTemplateId) return;

  try {
    selectingAll.value = true;
    const ids: string[] = [];
    let page = 1;
    let pages = 1;

    do {
      const response = await apis.assets.asset.list({
        page,
        perPage: SELECT_ALL_PER_PAGE,
        assetTemplateId: props.fromTemplateId,
      });
      response.items.forEach(a => a.id && ids.push(a.id));
      pages = response.pagination?.totalPages ?? 1;
      page++;
    } while (page <= pages);

    selectedIds.value = new Set(ids);
    emitSelection();
  } catch (error: unknown) {
    logger.error('Error selecting all assets:', error);
    handleApiError(error, { defaultMessage: t.steps.assets.loadError.value, timeout: 5000 });
  } finally {
    selectingAll.value = false;
  }
}

/**
 * Handle pagination navigation.
 * @param page - New 1-based page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchAssets(page);
}

/** WATCHERS */
// A new source template invalidates any prior selection: reset and reload.
watch(
  () => props.fromTemplateId,
  () => {
    selectedIds.value = new Set();
    currentPage.value = 1;
    emitSelection();
    void fetchAssets(1);
  }
);

/** LIFECYCLE HOOKS */
onMounted(() => {
  emit('update:valid', selectedIds.value.size > 0);
  void fetchAssets(1);
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="inventory_2" color="primary" class="q-mr-xs" />
        {{ t.steps.assets.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.assets.subtitle.value }}
      </div>
    </div>

    <!-- Selection toolbar -->
    <div class="row items-center q-col-gutter-sm q-mb-sm">
      <div class="col text-caption text-grey-7">
        {{ t.steps.assets.counter(selectedCount, totalItems) }}
      </div>
      <div class="col-auto row items-center q-gutter-sm">
        <q-btn
          flat
          dense
          no-caps
          size="sm"
          color="primary"
          icon="done_all"
          :label="t.steps.assets.selectAll(totalItems)"
          :loading="selectingAll"
          :disable="loading || totalItems === 0"
          @click="selectAllMatching"
        />
        <q-btn
          v-if="selectedCount > 0"
          flat
          dense
          no-caps
          size="sm"
          color="grey-7"
          icon="clear"
          :label="t.steps.assets.clearSelection.value"
          @click="clearSelection"
        />
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="2.5em" />
    </div>

    <!-- Empty -->
    <div v-else-if="assets.length === 0" class="text-caption text-grey-7 q-pa-md text-center">
      {{ t.steps.assets.empty.value }}
    </div>

    <!-- Asset list -->
    <q-list v-else bordered separator class="rounded-borders asset-list">
      <q-item class="bg-grey-1">
        <q-item-section side>
          <q-checkbox
            :model-value="pageAllSelected"
            dense
            @update:model-value="toggleSelectPage"
          />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-caption text-weight-medium text-grey-7">
            {{ t.steps.assets.selectPage.value }}
          </q-item-label>
        </q-item-section>
      </q-item>

      <q-item v-for="asset in assets" :key="asset.id" dense>
        <q-item-section side>
          <q-checkbox
            :model-value="selectedIds.has(asset.id)"
            dense
            @update:model-value="(checked: boolean) => toggleAsset(asset.id, checked)"
          />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ asset.name }}</q-item-label>
          <q-item-label caption class="text-mono">{{ asset.assetUUID }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Pagination -->
    <ListPagination
      v-if="!loading && totalPages > 1"
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.asset-list {
  max-height: 420px;
  overflow-y: auto;
}

.text-mono {
  font-family: monospace;
}
</style>

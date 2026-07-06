<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateMarketplacePage'
});

/** TYPE IMPORTS */
import type {
  AssetTemplateCatalogItem,
  AssetTemplateCatalogQuery,
  AssetTemplateFacets
} from '@mapexos/schemas';
import type { MarketplacePageFilters, MarketplaceSelection } from './interfaces';

/** VUE IMPORTS */
import { ref, reactive, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { ListCardEmpty } from '@components/cards';
import { MarketplaceCard } from './components/MarketplaceCard';
import { AssetTemplateMarketplaceDetailModal } from './components/AssetTemplateMarketplaceDetailModal';

/** COMPOSABLES */
import { useAssetTemplateMarketplaceTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES & STORES */
const t = useAssetTemplateMarketplaceTranslations();
const { locale } = useI18n();
const logger = useLogger('AssetTemplateMarketplacePage');

/** STATE */
const PER_PAGE = 24;
const items = ref<AssetTemplateCatalogItem[]>([]);
const total = ref(0);
const page = ref(0);
const loaded = ref(false);
const error = ref(false);
const filterKey = ref(0);

const filters = reactive<MarketplacePageFilters>({
  search: '',
  category: null,
  vendor: null,
  model: null,
  version: null,
});

const facets = ref<AssetTemplateFacets>({
  categories: [],
  vendors: [],
  models: [],
  versions: [],
});

const selected = ref<MarketplaceSelection | null>(null);
const modalOpen = ref(false);

/** FUNCTIONS */

/**
 * Build the catalog list query for a given page from the active filters and
 * locale. Empty facet/search values are omitted so the server returns all.
 * @param {number} targetPage - The 1-based page to request.
 * @returns {AssetTemplateCatalogQuery} The query payload.
 */
function buildQuery(targetPage: number): AssetTemplateCatalogQuery {
  return {
    page: targetPage,
    perPage: PER_PAGE,
    lang: locale.value,
    search: filters.search || undefined,
    category: filters.category ?? undefined,
    vendor: filters.vendor ?? undefined,
    model: filters.model ?? undefined,
    version: filters.version ?? undefined,
  };
}

/**
 * Infinite-scroll loader. Appends the next page of results and tells Quasar to
 * stop once every catalog item has been loaded (or the server returns none).
 * @param {number} _index - Quasar's 1-based load index (unused; own counter tracks pages).
 * @param {(stop?: boolean) => void} done - Quasar callback; done(true) halts further loads.
 */
async function onLoad(_index: number, done: (stop?: boolean) => void): Promise<void> {
  if (!apis.assetTemplatesMarketplace) {
    done(true);
    return;
  }

  const next = page.value + 1;

  try {
    const res = await apis.assetTemplatesMarketplace.marketplace.list(buildQuery(next));
    items.value = next === 1 ? res.data.items : [...items.value, ...res.data.items];
    total.value = res.data.total;
    page.value = next;
    error.value = false;
    done(res.data.items.length === 0 || items.value.length >= total.value);
  } catch (err) {
    logger.error('Failed to load the asset template marketplace catalog', err);
    error.value = true;
    done(true);
  } finally {
    loaded.value = true;
  }
}

/**
 * Reset pagination and remount the infinite-scroll list so it reloads from the
 * first page under the current filters.
 */
function applyFilters(): void {
  items.value = [];
  total.value = 0;
  page.value = 0;
  loaded.value = false;
  error.value = false;
  filterKey.value += 1;
}

/**
 * Fetch the facet options for the current drill-down (vendor narrows models,
 * model narrows versions) in the active locale.
 */
async function loadFacets(): Promise<void> {
  if (!apis.assetTemplatesMarketplace) return;

  try {
    const res = await apis.assetTemplatesMarketplace.marketplace.facets({
      vendor: filters.vendor ?? undefined,
      model: filters.model ?? undefined,
      lang: locale.value,
    });
    facets.value = res.data;
  } catch (err) {
    logger.error('Failed to load the asset template marketplace facets', err);
  }
}

/**
 * Vendor facet changed: clear the dependent model/version, refresh the
 * drill-down options, and reload the grid.
 */
function onVendorChange(): void {
  filters.model = null;
  filters.version = null;
  void loadFacets();
  applyFilters();
}

/**
 * Model facet changed: clear the dependent version, refresh the version
 * options, and reload the grid.
 */
function onModelChange(): void {
  filters.version = null;
  void loadFacets();
  applyFilters();
}

/**
 * Clear every filter and the search query, refresh the facet options, and
 * reload the grid from the full catalog. Backs the empty-state action.
 */
function clearAllFilters(): void {
  filters.search = '';
  filters.category = null;
  filters.vendor = null;
  filters.model = null;
  filters.version = null;
  void loadFacets();
  applyFilters();
}

/**
 * Retry loading the catalog after a load error. Backs the error-state action.
 */
function reload(): void {
  applyFilters();
}

/**
 * Open the detail modal for the selected catalog item.
 * @param {AssetTemplateCatalogItem} item - The catalog item to preview.
 */
function handleView(item: AssetTemplateCatalogItem): void {
  selected.value = { vendor: item.vendor, slug: item.slug };
  modalOpen.value = true;
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadFacets();
});
</script>

<template>
  <q-page class="q-pt-lg marketplace-page">
    <PageHeader
      icon="storefront"
      icon-color="primary"
      :title="t.page.title.value"
      :description="t.page.subtitle.value"
    />

    <!-- Filters -->
    <div class="marketplace-page__filters">
      <q-input
        v-model="filters.search"
        dense
        outlined
        clearable
        debounce="300"
        class="marketplace-page__search"
        :placeholder="t.search.placeholder.value"
        @update:model-value="applyFilters"
      >
        <template #prepend>
          <q-icon name="search" />
        </template>
      </q-input>

      <q-select
        v-model="filters.category"
        dense
        outlined
        clearable
        emit-value
        map-options
        option-value="value"
        option-label="label"
        class="marketplace-page__facet"
        :label="t.facets.category.value"
        :options="facets.categories"
        @update:model-value="applyFilters"
      />

      <q-select
        v-model="filters.vendor"
        dense
        outlined
        clearable
        emit-value
        map-options
        option-value="value"
        option-label="label"
        class="marketplace-page__facet"
        :label="t.facets.vendor.value"
        :options="facets.vendors"
        @update:model-value="onVendorChange"
      />

      <q-select
        v-model="filters.model"
        dense
        outlined
        clearable
        emit-value
        map-options
        option-value="value"
        option-label="label"
        class="marketplace-page__facet"
        :label="t.facets.model.value"
        :options="facets.models"
        :disable="!filters.vendor"
        @update:model-value="onModelChange"
      />

      <q-select
        v-model="filters.version"
        dense
        outlined
        clearable
        emit-value
        map-options
        option-value="value"
        option-label="label"
        class="marketplace-page__facet"
        :label="t.facets.version.value"
        :options="facets.versions"
        :disable="!filters.model"
        @update:model-value="applyFilters"
      />
    </div>

    <!-- Error -->
    <div v-if="error" class="row">
      <ListCardEmpty
        icon="cloud_off"
        :title="t.error.title.value"
        :description="t.error.description.value"
        :button-label="t.error.button.value"
        button-icon="refresh"
        @button-click="reload"
      />
    </div>

    <!-- Empty -->
    <div v-else-if="loaded && items.length === 0" class="row">
      <ListCardEmpty
        icon="storefront"
        :title="t.empty.title.value"
        :description="t.empty.description.value"
        :button-label="t.empty.button.value"
        button-icon="filter_alt_off"
        @button-click="clearAllFilters"
      />
    </div>

    <!-- Grid + infinite scroll -->
    <q-infinite-scroll v-else :key="filterKey" :offset="250" @load="onLoad">
      <div class="row q-col-gutter-md">
        <div
          v-for="item in items"
          :key="item.id"
          class="col-12 col-sm-6 col-md-4 col-lg-3"
        >
          <MarketplaceCard :item="item" @view="handleView" />
        </div>
      </div>

      <template #loading>
        <div class="row justify-center q-my-md">
          <q-spinner color="primary" size="32px" />
        </div>
      </template>
    </q-infinite-scroll>

    <!-- Detail modal -->
    <AssetTemplateMarketplaceDetailModal
      v-model="modalOpen"
      :vendor="selected?.vendor ?? null"
      :slug="selected?.slug ?? null"
    />
  </q-page>
</template>

<style scoped>
.marketplace-page__filters {
  display: flex;
  flex-wrap: wrap;
  gap: var(--mapex-spacing-md);
  margin-bottom: var(--mapex-spacing-lg);
}

.marketplace-page__search {
  flex: 2 1 260px;
}

.marketplace-page__facet {
  flex: 1 1 160px;
}
</style>

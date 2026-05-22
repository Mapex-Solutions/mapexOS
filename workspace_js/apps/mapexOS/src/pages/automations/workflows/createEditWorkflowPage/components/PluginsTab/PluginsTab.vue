<script setup lang="ts">
/** TYPE IMPORTS */
import type {
  RegistryEntry,
  InstalledPlugin,
  InstalledPluginGroup,
  PluginDetailState,
} from './interfaces';
import type { PluginNodeTypeSummary } from '@components/dialogs/pluginDetailDialog/interfaces';
import type { PluginCategory, PluginCredentialDefinition } from '@src/components/workflow/interfaces';
/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { PluginCard } from '../PluginCard';
import { PluginDetailDialog } from '@components/dialogs/pluginDetailDialog';
import { CredentialManagerDialog } from '../CredentialManagerDialog';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { PLUGIN_CDN_BASE_URL, MARKETPLACE_PAGE_SIZE } from './constants';
import { loadManifest, convertManifestToPlugin } from '../../utils/manifestLoader';

/** COMPOSABLES & STORES */
const pluginRegistry = usePluginRegistryStore();
const t = useCreateEditWorkflowTranslations();
const { addInstalledPlugin, removeInstalledPlugin } = useWorkflowEditorState();

/** STATE */

/**
 * Active sub-tab
 */
const subTab = ref<'installed' | 'marketplace'>('installed');

/**
 * Registry entries fetched from API
 */
const registryEntries = ref<RegistryEntry[]>([]);

/**
 * Loading state for registry fetch
 */
const loading = ref(false);

/**
 * Plugin ID currently being installed
 */
const installingId = ref<string | null>(null);

/**
 * Search query for marketplace
 */
const searchQuery = ref('');

/**
 * Active category filter for marketplace (null = show all)
 */
const marketplaceCategoryFilter = ref<PluginCategory | null>(null);

/**
 * Current page for marketplace pagination (1-indexed)
 */
const marketplacePage = ref(1);

/**
 * Total marketplace plugins from API (for server-side pagination)
 */
const marketplaceTotal = ref(0);

/**
 * Plugin ID currently being enabled/disabled via PATCH
 */
const togglingId = ref<string | null>(null);

/**
 * Loading state for installed plugins API fetch
 */
const installedLoading = ref(false);

/**
 * Installed plugins fetched from the workflow service API
 */
const installedApiPlugins = ref<InstalledPlugin[]>([]);

/**
 * Category filter for non-core installed plugins (null = show all)
 */
const installedCategoryFilter = ref<PluginCategory | null>(null);

/**
 * Plugin detail dialog visibility
 */
const showDetailDialog = ref(false);

/**
 * Detail dialog state — populated from core plugin or registry entry
 */
const detailState = ref<PluginDetailState | null>(null);

/**
 * Whether the detail dialog manifest is loading
 */
const detailLoading = ref(false);

/**
 * Node type summaries for the detail dialog
 */
const detailNodeTypes = ref<PluginNodeTypeSummary[]>([]);

/**
 * Credential manager dialog visibility
 */
const showCredentialDialog = ref(false);

/**
 * Plugin ID for the active credential dialog
 */
const credentialPluginId = ref('');

/**
 * Plugin name for the active credential dialog
 */
const credentialPluginName = ref('');

/**
 * Credential definition for the active credential dialog
 */
const credentialDefs = ref<PluginCredentialDefinition[]>([]);

/** COMPUTED */

/**
 * Installed plugin IDs for quick lookup (excluding core).
 * Merges store (session) + API (persisted) sources.
 */
const installedIds = computed(() => {
  const ids = new Set<string>();
  for (const [id] of pluginRegistry.plugins) {
    if (!id.startsWith('core-')) ids.add(id);
  }
  for (const p of installedApiPlugins.value) {
    ids.add(p.id);
  }
  return ids;
});

/**
 * All installed plugins — core from store + non-core from API
 */
const installedPlugins = computed((): InstalledPlugin[] => {
  return [...corePlugins.value, ...installedApiPlugins.value];
});

/**
 * Core plugins from the local plugin registry store
 */
const corePlugins = computed((): InstalledPlugin[] => {
  const result: InstalledPlugin[] = [];
  for (const [, plugin] of pluginRegistry.plugins) {
    if (!plugin.id.startsWith('core-')) continue;
    result.push({
      id: plugin.id,
      name: plugin.name,
      version: plugin.version,
      icon: plugin.icon,
      brandIcon: '',
      color: '#4FA3E4',
      description: '',
      author: 'MapexOS',
      tags: [],
      category: plugin.category,
      nodeCount: plugin.nodeTypes.length,
      isCore: true,
    });
  }
  return result;
});

/**
 * Non-core (marketplace installed) plugins from the API
 */
const nonCorePlugins = computed(() => installedApiPlugins.value);

/**
 * Category options for non-core installed plugins
 */
const nonCoreCategoryOptions = computed(() => {
  const cats = new Set<PluginCategory>();
  for (const p of nonCorePlugins.value) {
    cats.add(p.category);
  }
  return Array.from(cats).map((cat) => ({
    label: getCategoryLabel(cat),
    value: cat,
  }));
});

/**
 * Non-core plugins filtered by category, then grouped
 */
const nonCoreGroups = computed((): InstalledPluginGroup[] => {
  let plugins = nonCorePlugins.value;

  if (installedCategoryFilter.value) {
    plugins = plugins.filter((p) => p.category === installedCategoryFilter.value);
  }

  const groupMap = new Map<PluginCategory, InstalledPlugin[]>();
  for (const plugin of plugins) {
    const list = groupMap.get(plugin.category);
    if (list) {
      list.push(plugin);
    } else {
      groupMap.set(plugin.category, [plugin]);
    }
  }

  const groups: InstalledPluginGroup[] = [];
  for (const [category, items] of groupMap) {
    groups.push({
      category,
      label: getCategoryLabel(category),
      plugins: items,
    });
  }

  return groups.sort((a, b) => a.label.localeCompare(b.label));
});

/**
 * Marketplace category options from registry entries
 */
const marketplaceCategoryOptions = computed(() => {
  const cats = new Set<PluginCategory>();
  for (const entry of registryEntries.value) {
    cats.add(entry.category);
  }
  return Array.from(cats).map((cat) => ({
    label: getCategoryLabel(cat),
    value: cat,
  }));
});

/**
 * Total pages for marketplace pagination (server-side)
 */
const marketplaceTotalPages = computed(() =>
  Math.ceil(marketplaceTotal.value / MARKETPLACE_PAGE_SIZE),
);

/** WATCHERS */

watch([searchQuery, marketplaceCategoryFilter], () => {
  marketplacePage.value = 1;
  void fetchPlugins();
});

watch(marketplacePage, () => {
  void fetchPlugins();
});

watch(subTab, () => {
  marketplaceCategoryFilter.value = null;
  installedCategoryFilter.value = null;
  searchQuery.value = '';
  marketplacePage.value = 1;
  if (subTab.value === 'marketplace') {
    // Fetch CDN registry + API enabled list in parallel
    void fetchPlugins();
    void fetchInstalledPlugins();
  } else {
    void fetchInstalledPlugins();
  }
});

/** FUNCTIONS */

/**
 * Check if a plugin is already installed
 *
 * @param {string} pluginId - Plugin ID to check
 * @returns {boolean} True if installed
 */
function isInstalled(pluginId: string): boolean {
  return installedIds.value.has(pluginId);
}

/**
 * Get the brand icon URL for a registry entry
 *
 * @param {string} brandIcon - Relative brand icon path
 * @returns {string} Full URL to the brand icon
 */
function getBrandIconUrl(brandIcon: string): string {
  return `${PLUGIN_CDN_BASE_URL}/${brandIcon}`;
}

/**
 * Get the translated label for a plugin category
 *
 * @param {PluginCategory} category - Category key
 * @returns {string} Translated label
 */
function getCategoryLabel(category: PluginCategory): string {
  const categoryMap: Record<string, string> = {
    triggers: t.pluginCatalog.categories.triggers.value,
    logic: t.pluginCatalog.categories.logic.value,
    state: t.pluginCatalog.categories.state.value,
    flow_control: t.pluginCatalog.categories.flowControl.value,
    timers: t.pluginCatalog.categories.timers.value,
    integrations: t.pluginCatalog.categories.integrations.value,
    observability: t.pluginCatalog.categories.observability.value,
    annotations: t.pluginCatalog.categories.annotations.value,
    custom: t.pluginCatalog.categories.custom.value,
  };
  return categoryMap[category] ?? category;
}

/**
 * Convert a CDN registry entry to a UI RegistryEntry
 *
 * @param {Record<string, unknown>} raw - Raw entry from registry.json
 * @returns {RegistryEntry} UI-friendly registry entry
 */
function convertCdnEntryToRegistryEntry(raw: Record<string, unknown>): RegistryEntry {
  const metadata = raw.metadata as Record<string, string> | undefined;
  return {
    id: (raw.pluginId as string) ?? (raw.id as string) ?? '',
    name: (raw.name as string) ?? '',
    version: (raw.version as string) ?? '0.0.0',
    category: ((raw.category as string) ?? 'integrations'),
    icon: (raw.icon as string) ?? 'extension',
    brandIcon: metadata?.brandIcon ?? '',
    color: metadata?.color ?? (raw.color as string) ?? '#666',
    description: (raw.description as string) ?? '',
    author: (raw.author as string) ?? '',
    tags: (raw.tags as string[]) ?? [],
    manifestUrl: (raw.manifestUrl as string) ?? '',
    docsUrl: metadata?.docs ?? '',
    requiresCredentials: raw.requiresCredentials === true || raw.credentials !== undefined,
    nodeCount: (raw.nodeCount as number) ?? 0,
    triggerCount: (raw.triggerCount as number) ?? 0,
    enabled: (raw.enabled as boolean) ?? false,
    isSystem: (raw.isSystem as boolean) ?? false,
  };
}

/**
 * Open the detail dialog for a marketplace plugin (fetches manifest)
 *
 * @param {RegistryEntry} entry - Registry entry to show details for
 * @returns {void}
 */
function openDetailDialog(entry: RegistryEntry): void {
  detailState.value = {
    id: entry.id,
    name: entry.name,
    author: entry.author,
    version: entry.version,
    description: entry.description,
    brandIconUrl: getBrandIconUrl(entry.brandIcon),
    icon: entry.icon,
    color: entry.color,
    category: entry.category,
    tags: entry.tags,
    isCore: false,
    manifestUrl: entry.manifestUrl,
  };
  detailNodeTypes.value = [];
  detailLoading.value = true;
  showDetailDialog.value = true;
  void fetchDetailNodeTypes(entry);
}

/**
 * Open the detail dialog for a core (installed) plugin — reads node types from store
 *
 * @param {string} pluginId - Core plugin ID
 * @returns {void}
 */
function openCorePluginDetail(pluginId: string): void {
  const plugin = pluginRegistry.plugins.get(pluginId);
  if (!plugin) return;

  detailState.value = {
    id: plugin.id,
    name: plugin.name,
    author: 'MapexOS',
    version: plugin.version,
    description: '',
    brandIconUrl: '',
    icon: plugin.icon,
    color: '#4FA3E4',
    category: plugin.category,
    tags: [],
    isCore: true,
    manifestUrl: null,
  };
  detailNodeTypes.value = plugin.nodeTypes
    .filter((nt) => !nt.catalogHidden)
    .map((nt) => ({
      type: nt.type,
      label: nt.label,
      icon: nt.icon,
      color: nt.color,
      description: nt.description,
      inputCount: nt.inputs?.length ?? 0,
      outputCount: nt.outputs?.length ?? 0,
    }));
  detailLoading.value = false;
  showDetailDialog.value = true;
}

/**
 * Open the detail dialog for a non-core installed plugin.
 * Uses data from InstalledPlugin + plugin registry store (no CDN dependency).
 * Falls back to CDN registry entry if available.
 *
 * @param {InstalledPlugin} plugin - Installed plugin data
 */
function openInstalledPluginDetail(plugin: InstalledPlugin): void {
  // All data comes from the API (DB manifest) — never from CDN
  const registered = pluginRegistry.plugins.get(plugin.id);

  detailState.value = {
    id: plugin.id,
    name: plugin.name,
    author: plugin.author,
    version: plugin.version,
    description: plugin.description,
    brandIconUrl: plugin.brandIcon ? getBrandIconUrl(plugin.brandIcon) : '',
    icon: plugin.icon,
    color: plugin.color,
    category: plugin.category,
    tags: plugin.tags,
    isCore: plugin.isCore,
    manifestUrl: null,
  };

  // Get node types from the local registry if available
  if (registered) {
    detailNodeTypes.value = registered.nodeTypes
      .filter((nt) => !nt.catalogHidden)
      .map((nt) => ({
        type: nt.type,
        label: nt.label,
        icon: nt.icon,
        color: nt.color ?? plugin.color,
        description: nt.description,
        inputCount: nt.inputs?.length ?? 0,
        outputCount: nt.outputs?.length ?? 0,
      }));
    detailLoading.value = false;
  } else {
    detailNodeTypes.value = [];
    detailLoading.value = false;
  }

  showDetailDialog.value = true;
}

/**
 * Fetch full plugin manifest from API and extract node type summaries for the detail dialog
 *
 * @param {RegistryEntry} entry - Registry entry to fetch details for
 * @returns {Promise<void>}
 */
async function fetchDetailNodeTypes(entry: RegistryEntry): Promise<void> {
  try {
    // Fetch full manifest from CDN
    const plugin = await loadManifest(entry.manifestUrl);

    detailNodeTypes.value = plugin.nodeTypes.map((nt) => ({
      type: nt.type,
      label: nt.label,
      icon: nt.icon,
      color: nt.color ?? entry.color,
      description: nt.description,
      inputCount: nt.inputs?.length ?? 0,
      outputCount: nt.outputs?.length ?? 0,
    }));
  } catch (error) {
    console.error('[PluginsTab] Failed to load plugin details:', error);
  } finally {
    detailLoading.value = false;
  }
}

/**
 * Fetch installed (enabled) plugins from the workflow service API
 *
 * @returns {Promise<void>}
 */
async function fetchInstalledPlugins(): Promise<void> {
  installedLoading.value = true;
  try {
    const plugins = await apis.workflows.plugin.getEnabled();

    installedApiPlugins.value = (plugins ?? [])
      .filter((p) => !p.pluginId.startsWith('core-'))
      .map((p) => {
        const meta = p.metadata as Record<string, string> | undefined;
        const item: InstalledPlugin = {
          id: p.pluginId,
          name: p.name,
          version: p.version,
          icon: p.icon ?? 'extension',
          brandIcon: meta?.brandIcon ?? '',
          color: p.color ?? '#666',
          description: p.description ?? '',
          author: p.author ?? '',
          tags: p.tags ?? [],
          category: ((p.category ?? 'integrations')),
          nodeCount: p.nodeTypes?.length ?? 0,
          isCore: false,
        };
        if (p.id !== undefined) item.mongoId = p.id;
        return item;
      });
  } catch (error) {
    console.warn('[PluginsTab] Failed to fetch installed plugins from API:', error);
    installedApiPlugins.value = [];
  } finally {
    installedLoading.value = false;
  }
}

/**
 * Fetch marketplace plugins from CDN registry
 *
 * @returns {Promise<void>}
 */
async function fetchPlugins(): Promise<void> {
  loading.value = true;
  try {
    const response = await fetch(`${PLUGIN_CDN_BASE_URL}/registry.json`);

    if (!response.ok) {
      throw new Error(`CDN returned ${response.status}`);
    }

    const data = await response.json() as { plugins?: Record<string, unknown>[] };
    const allEntries = (data.plugins ?? []).map(convertCdnEntryToRegistryEntry);

    // Client-side filtering (CDN is static, no server-side filtering)
    let filtered = allEntries;

    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase();
      filtered = filtered.filter(
        (e) => e.name.toLowerCase().includes(q) || e.description.toLowerCase().includes(q),
      );
    }

    if (marketplaceCategoryFilter.value) {
      filtered = filtered.filter((e) => e.category === marketplaceCategoryFilter.value);
    }

    // Client-side pagination
    marketplaceTotal.value = filtered.length;
    const start = (marketplacePage.value - 1) * MARKETPLACE_PAGE_SIZE;
    registryEntries.value = filtered.slice(start, start + MARKETPLACE_PAGE_SIZE);
  } catch (error) {
    console.error('[PluginsTab] Failed to fetch plugins from CDN:', error);
    registryEntries.value = [];
    marketplaceTotal.value = 0;
  } finally {
    loading.value = false;
  }
}

/**
 * Install a plugin: fetch manifest from CDN → persist via API → register locally.
 * If the API call fails, no local registration happens (rollback).
 *
 * @param {RegistryEntry} entry - Plugin to install
 * @returns {Promise<void>}
 */
async function enablePlugin(entry: RegistryEntry): Promise<void> {
  installingId.value = entry.id;
  try {
    // 1. Fetch raw manifest JSON from CDN
    const cdnUrl = `${PLUGIN_CDN_BASE_URL}/${entry.manifestUrl}`;
    const cdnResponse = await fetch(cdnUrl);
    if (!cdnResponse.ok) {
      throw new Error(`CDN returned ${cdnResponse.status}`);
    }
    const rawManifest = await cdnResponse.json() as Record<string, unknown>;

    // 2. Persist to backend — POST /api/v1/plugins
    const saved = await apis.workflows.plugin.create(
      rawManifest as unknown as Parameters<typeof apis.workflows.plugin.create>[0],
    );

    // 3. Convert API response to WorkflowPlugin and register locally
    const plugin = convertManifestToPlugin(saved);
    pluginRegistry.registerPlugin(plugin);
    addInstalledPlugin(entry.id);

    // 4. Add to installedApiPlugins for immediate UI feedback
    const savedMeta = saved.metadata as Record<string, string> | undefined;
    const installedEntry: InstalledPlugin = {
      id: saved.pluginId,
      name: saved.name,
      version: saved.version,
      icon: saved.icon ?? 'extension',
      brandIcon: savedMeta?.brandIcon ?? '',
      color: saved.color ?? '#666',
      description: saved.description ?? '',
      author: saved.author ?? '',
      tags: saved.tags ?? [],
      category: ((saved.category ?? 'integrations')),
      nodeCount: saved.nodeTypes?.length ?? 0,
      isCore: false,
    };
    if (saved.id !== undefined) installedEntry.mongoId = saved.id;
    installedApiPlugins.value.push(installedEntry);

    notifySuccess({ message: t.pluginsTab.installSuccessMsg(entry.name) });
  } catch (error) {
    console.error('[PluginsTab] Failed to install plugin:', error);
    notifyFail({ message: t.pluginsTab.installFailed.value });
  } finally {
    installingId.value = null;
  }
}

/**
 * Uninstall a plugin: delete from backend API → remove locally.
 * If the API call fails, no local removal happens (rollback).
 *
 * @param {string} pluginId - Plugin string ID (e.g., 'telegram')
 * @param {string} pluginName - Plugin display name for notification
 * @returns {Promise<void>}
 */
async function disablePlugin(pluginId: string, pluginName: string): Promise<void> {
  togglingId.value = pluginId;
  try {
    // 1. Find the MongoDB ObjectID for this plugin
    const installed = installedApiPlugins.value.find((p) => p.id === pluginId);
    if (!installed?.mongoId) {
      throw new Error(`No MongoDB ID found for plugin "${pluginId}"`);
    }

    // 2. Delete from backend — DELETE /api/v1/plugins/:id
    await apis.workflows.plugin.delete({ id: installed.mongoId });

    // 3. On success, remove from local state
    pluginRegistry.unregisterPlugin(pluginId);
    removeInstalledPlugin(pluginId);
    installedApiPlugins.value = installedApiPlugins.value.filter((p) => p.id !== pluginId);

    notifySuccess({ message: t.pluginsTab.uninstallSuccessMsg(pluginName) });
  } catch (error) {
    console.error('[PluginsTab] Failed to uninstall plugin:', error);
    notifyFail({ message: t.pluginsTab.installFailed.value });
  } finally {
    togglingId.value = null;
  }
}

/**
 * Open the credential manager dialog for a plugin.
 * Looks up the credential definition from the plugin registry.
 *
 * @param {InstalledPlugin} plugin - Installed plugin
 */
function openCredentialDialog(plugin: InstalledPlugin): void {
  const registered = pluginRegistry.plugins.get(plugin.id);
  if (!registered?.credentials) return;

  credentialPluginId.value = plugin.id;
  credentialPluginName.value = plugin.name;
  credentialDefs.value = registered.credentials ?? [];
  showCredentialDialog.value = true;
}

/**
 * Check if a plugin has credential definitions
 *
 * @param {string} pluginId - Plugin ID
 * @returns {boolean} True if the plugin requires credentials
 */
function hasCredentials(pluginId: string): boolean {
  const registered = pluginRegistry.plugins.get(pluginId);
  return !!registered?.credentials;
}

/** LIFECYCLE HOOKS */

onMounted(() => {
  // Always fetch enabled plugins from API (populates installedIds)
  void fetchInstalledPlugins();

  if (subTab.value === 'marketplace') {
    void fetchPlugins();
  }
});
</script>

<template>
  <div class="plugins-tab">
    <!-- Sub-tabs Header -->
    <div class="row items-center q-mb-lg">
      <q-tabs
        v-model="subTab"
        dense
        no-caps
        class="text-grey-8 plugins-subtabs"
        active-color="primary"
        indicator-color="primary"
      >
        <q-tab name="installed" icon="extension" :label="t.pluginsTab.installedTab.value">
          <q-badge
            v-if="installedPlugins.length > 0"
            color="primary"
            floating
            rounded
          >
            {{ installedPlugins.length }}
          </q-badge>
        </q-tab>
        <q-tab name="marketplace" icon="storefront" :label="t.pluginsTab.marketplaceTab.value" />
      </q-tabs>
      <q-space />
      <span class="text-caption text-grey-6">
        <q-icon name="info" size="xs" class="q-mr-xs" />
        {{ subTab === 'marketplace'
          ? t.pluginsTab.marketplaceHelp.value
          : t.pluginsTab.installedHelp.value
        }}
      </span>
    </div>

    <!-- ══════════════════════════════════════════ -->
    <!-- ── Installed Sub-tab ── -->
    <!-- ══════════════════════════════════════════ -->
    <div v-if="subTab === 'installed'">
      <!-- Loading -->
      <div v-if="installedLoading" class="flex flex-center q-pa-xl">
        <q-spinner color="primary" size="40px" />
      </div>

      <!-- Empty State -->
      <div v-else-if="installedPlugins.length === 0" class="text-center q-pa-xl">
        <q-icon name="extension_off" size="48px" color="grey-5" />
        <div class="text-subtitle1 text-grey-6 q-mt-sm">{{ t.pluginsTab.noPluginsInstalled.value }}</div>
        <div class="text-caption text-grey-5">{{ t.pluginsTab.noPluginsInstalledDesc.value }}</div>
        <q-btn
          flat
          no-caps
          color="primary"
          icon="storefront"
          class="q-mt-md"
          :label="t.pluginsTab.marketplaceTab.value"
          @click="subTab = 'marketplace'"
        />
      </div>

      <template v-else>
        <!-- ── Core Plugins (compact list) ── -->
        <div v-if="corePlugins.length > 0" class="q-mb-lg">
          <div class="row items-center q-mb-xs plugins-tab__section-header">
            <span class="text-subtitle2 text-weight-bold text-uppercase">
              {{ t.pluginsTab.corePlugins.value }}
            </span>
            <q-badge outline color="grey-6" class="q-ml-sm text-caption">
              {{ corePlugins.length }}
            </q-badge>
          </div>

          <q-list bordered separator dense class="plugins-tab__list">
            <q-item
              v-for="plugin in corePlugins"
              :key="plugin.id"
              clickable
              @click="openCorePluginDetail(plugin.id)"
            >
              <q-item-section avatar>
                <q-icon :name="plugin.icon" size="22px" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium text-body2">{{ plugin.name }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="row items-center q-gutter-sm">
                  <span class="text-caption text-grey-6">
                    {{ plugin.nodeCount }} {{ plugin.nodeCount === 1 ? 'node' : 'nodes' }}
                  </span>
                  <q-icon name="chevron_right" size="xs" color="grey-6" />
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </div>

        <!-- ── Non-Core Plugins (cards by category) ── -->
        <div v-if="nonCorePlugins.length > 0">
          <div class="row items-center q-mb-md">
            <span class="text-subtitle2 text-weight-bold text-uppercase plugins-tab__section-header">
              {{ t.pluginsTab.installedPlugins.value }}
            </span>
            <q-badge outline color="grey-6" class="q-ml-sm text-caption">
              {{ nonCorePlugins.length }}
            </q-badge>
            <q-space />
            <q-select
              v-if="nonCoreCategoryOptions.length > 1"
              v-model="installedCategoryFilter"
              outlined
              dense
              clearable
              emit-value
              map-options
              class="plugins-tab__category-filter"
              :options="nonCoreCategoryOptions"
              :label="t.pluginsTab.filterByCategory.value"
            >
              <template #prepend>
                <q-icon name="filter_list" size="xs" />
              </template>
            </q-select>
          </div>

          <!-- Empty after filter -->
          <div v-if="nonCoreGroups.length === 0" class="text-center q-pa-lg">
            <q-icon name="filter_alt_off" size="48px" color="grey-5" />
            <div class="text-subtitle1 text-grey-6 q-mt-sm">{{ t.pluginsTab.noCategoryResults.value }}</div>
          </div>

          <!-- Grouped cards -->
          <div v-for="group in nonCoreGroups" :key="group.category" class="q-mb-lg">
            <div v-if="nonCoreGroups.length > 1" class="text-caption text-uppercase text-grey-6 text-weight-bold q-mb-sm">
              {{ group.label }}
            </div>
            <div class="row q-col-gutter-md">
              <div
                v-for="plugin in group.plugins"
                :key="plugin.id"
                class="col-12 col-sm-6 col-md-4"
              >
                <PluginCard
                  :name="plugin.name"
                  :description="plugin.description"
                  :author="plugin.author"
                  :version="plugin.version"
                  :brand-icon-url="plugin.brandIcon ? getBrandIconUrl(plugin.brandIcon) : ''"
                  :icon="plugin.icon"
                  :color="plugin.color"
                  :category-label="getCategoryLabel(plugin.category)"
                  :tags="plugin.tags"
                  :node-count="plugin.nodeCount"
                  installed
                  :installing="false"
                  :install-disabled="false"
                  :install-label="t.pluginsTab.install.value"
                  :installing-label="t.pluginsTab.installing.value"
                  :installed-label="t.pluginsTab.installed.value"
                  :details-label="t.pluginsTab.details.value"
                  @details="openInstalledPluginDetail(plugin)"
                />
                <!-- Plugin actions -->
                <div class="row justify-end q-mt-xs q-gutter-xs">
                  <q-btn
                    v-if="hasCredentials(plugin.id)"
                    flat
                    dense
                    no-caps
                    color="primary"
                    icon="vpn_key"
                    size="sm"
                    :label="t.pluginsTab.credentialsTitle.value"
                    @click="openCredentialDialog(plugin)"
                  />
                  <q-btn
                    flat
                    dense
                    no-caps
                    color="negative"
                    icon="delete_outline"
                    size="sm"
                    :label="t.pluginsTab.uninstall.value"
                    :loading="togglingId === plugin.id"
                    @click="disablePlugin(plugin.id, plugin.name)"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- ══════════════════════════════════════════ -->
    <!-- ── Marketplace Sub-tab ── -->
    <!-- ══════════════════════════════════════════ -->
    <div v-else>
      <!-- Search + Category Filter -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-xs-12 col-sm-6 col-md-6">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            clearable
            :placeholder="t.pluginsTab.searchPlaceholder.value"
          >
            <template #prepend>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
        <div class="col-xs-12 col-sm-6 col-md-6">
          <q-select
            v-model="marketplaceCategoryFilter"
            outlined
            dense
            clearable
            emit-value
            map-options
            :options="marketplaceCategoryOptions"
            :label="t.pluginsTab.filterByCategory.value"
          >
            <template #prepend>
              <q-icon name="filter_list" size="xs" />
            </template>
          </q-select>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex flex-center q-pa-xl">
        <q-spinner color="primary" size="40px" />
      </div>

      <!-- Empty -->
      <div v-else-if="registryEntries.length === 0" class="text-center q-pa-xl">
        <q-icon name="extension_off" size="48px" color="grey-5" />
        <div class="text-subtitle1 text-grey-6 q-mt-sm">{{ t.pluginsTab.noPluginsAvailable.value }}</div>
      </div>

      <!-- Plugin Cards Grid -->
      <template v-else>
        <div class="row q-col-gutter-md">
          <div
            v-for="entry in registryEntries"
            :key="entry.id"
            class="col-12 col-sm-6 col-md-4"
          >
            <PluginCard
              :name="entry.name"
              :description="entry.description"
              :author="entry.author"
              :version="entry.version"
              :brand-icon-url="getBrandIconUrl(entry.brandIcon)"
              :icon="entry.icon"
              :color="entry.color"
              :category-label="getCategoryLabel(entry.category)"
              :tags="entry.tags"
              :node-count="entry.nodeCount"
              :installed="isInstalled(entry.id)"
              :installing="installingId === entry.id"
              :install-disabled="installingId !== null"
              :install-label="t.pluginsTab.install.value"
              :installing-label="t.pluginsTab.installing.value"
              :installed-label="t.pluginsTab.installed.value"
              :details-label="t.pluginsTab.details.value"
              @install="enablePlugin(entry)"
              @details="openDetailDialog(entry)"
            />
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="marketplaceTotalPages > 1" class="flex flex-center q-mt-lg">
          <q-pagination
            v-model="marketplacePage"
            color="primary"
            :max="marketplaceTotalPages"
            :max-pages="7"
            boundary-numbers
          />
        </div>
      </template>
    </div>

    <!-- ── Credential Manager Dialog ── -->
    <CredentialManagerDialog
      v-if="credentialDefs.length > 0"
      v-model="showCredentialDialog"
      :plugin-id="credentialPluginId"
      :plugin-name="credentialPluginName"
      :credential-defs="credentialDefs"
    />

    <!-- ── Plugin Detail Dialog ── -->
    <PluginDetailDialog
      v-if="detailState"
      v-model="showDetailDialog"
      :name="detailState.name"
      :author="detailState.author"
      :version="detailState.version"
      :description="detailState.description"
      :brand-icon-url="detailState.brandIconUrl"
      :icon="detailState.icon"
      :color="detailState.color"
      :category-label="getCategoryLabel(detailState.category)"
      :tags="detailState.tags"
      :loading="detailLoading"
      :node-types="detailNodeTypes"
      :installed="detailState.isCore || isInstalled(detailState.id)"
      :installing="installingId === detailState.id"
      :install-disabled="installingId !== null"
      :install-label="t.pluginsTab.install.value"
      :installing-label="t.pluginsTab.installing.value"
      :installed-label="t.pluginsTab.installed.value"
      :node-types-label="t.pluginsTab.nodeTypes.value"
      :loading-label="t.pluginsTab.loadingManifest.value"
      :inputs-label="t.pluginsTab.inputs.value"
      :outputs-label="t.pluginsTab.outputs.value"
      @install="enablePlugin(registryEntries.find((e) => e.id === detailState!.id)!)"
    />
  </div>
</template>

<style lang="scss" scoped>
.plugins-subtabs {
  background: transparent;
  border-radius: var(--mapex-radius-md);

  :deep(.q-tab) {
    min-height: 36px;
    padding: 0 16px;
    border-radius: var(--mapex-radius-sm);
    margin-right: 4px;

    &:hover:not(.q-tab--active) {
      background: var(--mapex-surface-bg);
    }
  }

  :deep(.q-tabs__content) {
    padding: 4px;
    background: var(--mapex-page-bg);
    border-radius: var(--mapex-radius-md);
  }

  :deep(.q-tab--active) {
    background: var(--mapex-surface-elevated);
    box-shadow: var(--mapex-shadow-xs);
  }

  :deep(.q-tab__indicator) {
    display: none;
  }
}

.plugins-tab {
  &__section-header {
    padding: 4px 0;
    color: var(--mapex-text-secondary);
  }

  &__category-filter {
    max-width: 220px;
  }

  &__list {
    background: var(--mapex-card-bg);
    border-color: var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
  }
}
</style>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginCatalogProps, PluginCatalogEmits } from './interfaces/PluginCatalog.interface';
import type { CatalogGroup } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import PluginCatalogGroup from '../PluginCatalogGroup/PluginCatalogGroup.vue';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** PROPS & EMITS */
defineProps<PluginCatalogProps>();
const emit = defineEmits<PluginCatalogEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();
const pluginRegistry = usePluginRegistryStore();

/** STATE */

/**
 * Search query for filtering node types
 */
const search = ref('');

/** COMPUTED */

/**
 * Catalog groups filtered by search
 */
const filteredGroups = computed<CatalogGroup[]>(() => {
  const groups = pluginRegistry.catalog;
  if (!search.value.trim()) return groups;

  const query = search.value.toLowerCase();
  return groups
    .map(group => ({
      ...group,
      nodeTypes: group.nodeTypes.filter(
        nt =>
          nt.label.toLowerCase().includes(query) ||
          nt.description.toLowerCase().includes(query),
      ),
    }))
    .filter(group => group.nodeTypes.length > 0);
});
</script>

<template>
  <div class="plugin-catalog">
    <!-- Header -->
    <div class="plugin-catalog__header">
      <template v-if="!collapsed">
        <span class="text-subtitle2 text-weight-medium">{{ t.pluginCatalog.title.value }}</span>
        <q-space />
      </template>
      <q-btn
        flat
        dense
        round
        :icon="collapsed ? 'chevron_right' : 'chevron_left'"
        size="sm"
        @click="emit('toggle-collapse')"
      >
        <AppTooltip :content="collapsed ? t.pluginCatalog.expand.value : t.pluginCatalog.collapse.value" />
      </q-btn>
    </div>

    <!-- Search (only when expanded) -->
    <div v-if="!collapsed" class="plugin-catalog__search">
      <q-input
        v-model="search"
        dense
        outlined
        :placeholder="t.pluginCatalog.searchPlaceholder.value"
        class="q-px-sm"
      >
        <template #prepend>
          <q-icon name="search" size="xs" />
        </template>
        <template v-if="search" #append>
          <q-icon name="close" size="xs" class="cursor-pointer" @click="search = ''" />
        </template>
      </q-input>
    </div>

    <!-- Groups -->
    <div class="plugin-catalog__groups">
      <template v-if="filteredGroups.length > 0">
        <PluginCatalogGroup
          v-for="group in filteredGroups"
          :key="group.category"
          :group="group"
          :collapsed="collapsed"
        />
      </template>

      <!-- Empty search result -->
      <div v-else-if="search.trim()" class="text-center q-pa-md text-grey-6">
        <q-icon name="search_off" size="24px" />
        <p class="text-caption q-mt-sm">{{ t.pluginCatalog.noNodesFound.value }}</p>
      </div>

      <!-- No plugins registered -->
      <div v-else class="text-center q-pa-md text-grey-6">
        <q-icon name="extension" size="24px" />
        <p class="text-caption q-mt-sm">{{ t.pluginCatalog.noPlugins.value }}</p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.plugin-catalog {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--mapex-surface-bg);

  &__header {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid var(--mapex-card-border);
    min-height: 40px;
  }

  &__search {
    padding: 8px 4px;
    border-bottom: 1px solid var(--mapex-card-border);
  }

  &__groups {
    flex: 1;
    overflow-y: auto;
  }
}
</style>

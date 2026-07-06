<template>
  <div class="available-fields-list">
    <!-- Empty state -->
    <div v-if="fields.length === 0" class="empty-state text-grey-7">
      {{ t.availableFields.empty.value }}
    </div>

    <!-- Field paths -->
    <div v-else class="paths">
      <code
        v-for="(path, index) in fields"
        :key="index"
        class="path-chip"
      >{{ path }}</code>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'AvailableFieldsList'
});

/** TYPE IMPORTS */
import type { AvailableFieldsListProps } from './interfaces';

/** COMPOSABLES */
import { useAssetTemplateFieldsTranslations } from '@composables/i18n/components/assetTemplates/useAssetTemplateFieldsTranslations';

/** PROPS & EMITS */
defineProps<AvailableFieldsListProps>();

/** COMPOSABLES & STORES */
const t = useAssetTemplateFieldsTranslations();
</script>

<style lang="scss" scoped>
.available-fields-list {
  width: 100%;

  .empty-state {
    padding: 16px;
    text-align: center;
    font-size: 0.9rem;
  }

  .paths {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .path-chip {
    font-family: 'Courier New', monospace;
    font-size: 0.8rem;
    color: var(--mapex-text-primary);
    background: var(--mapex-surface-sunken);
    padding: 4px 8px;
    border-radius: var(--mapex-radius-sm);
    word-break: break-all;
  }
}
</style>

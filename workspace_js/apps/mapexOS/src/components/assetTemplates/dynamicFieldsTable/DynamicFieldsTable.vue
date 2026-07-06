<template>
  <div class="dynamic-fields-table">
    <!-- Empty state -->
    <div v-if="fields.length === 0" class="empty-state text-grey-7">
      {{ t.dynamicFields.empty.value }}
    </div>

    <!-- Fields table -->
    <q-markup-table v-else flat bordered dense class="fields-table">
      <thead>
        <tr>
          <th class="text-left">{{ t.dynamicFields.headers.field.value }}</th>
          <th class="text-left">{{ t.dynamicFields.headers.type.value }}</th>
          <th class="text-left">{{ t.dynamicFields.headers.value.value }}</th>
          <th class="text-left">{{ t.dynamicFields.headers.status.value }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, index) in fields" :key="row.fieldId ?? index">
          <td class="text-left text-weight-medium">{{ row.field }}</td>
          <td class="text-left">
            <DetailChip color="blue" size="sm" :label="row.type" />
          </td>
          <td class="text-left">
            <code v-if="row.value" class="value-path">{{ row.value }}</code>
            <span v-else class="text-grey-6">-</span>
          </td>
          <td class="text-left">
            <DetailChip
              :color="isActive(row) ? 'green' : 'grey'"
              size="sm"
              :label="isActive(row) ? t.dynamicFields.status.active.value : t.dynamicFields.status.deprecated.value"
            />
          </td>
        </tr>
      </tbody>
    </q-markup-table>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'DynamicFieldsTable'
});

/** TYPE IMPORTS */
import type { DynamicField, DynamicFieldsTableProps } from './interfaces';

/** COMPONENTS */
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useAssetTemplateFieldsTranslations } from '@composables/i18n/components/assetTemplates/useAssetTemplateFieldsTranslations';

/** PROPS & EMITS */
defineProps<DynamicFieldsTableProps>();

/** COMPOSABLES & STORES */
const t = useAssetTemplateFieldsTranslations();

/** FUNCTIONS */

/**
 * Whether a dynamic field is active (status equal to 1)
 * @param {DynamicField} field - Dynamic field row
 * @returns {boolean} True when the field is active
 */
function isActive(field: DynamicField): boolean {
  return field.status === 1;
}
</script>

<style lang="scss" scoped>
.dynamic-fields-table {
  width: 100%;

  .empty-state {
    padding: 16px;
    text-align: center;
    font-size: 0.9rem;
  }

  .fields-table {
    border-radius: var(--mapex-radius-md);

    th {
      font-size: 0.7rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.8px;
      color: var(--mapex-text-secondary);
    }

    td {
      font-size: 0.85rem;
      color: var(--mapex-text-primary);
    }

    .value-path {
      font-family: 'Courier New', monospace;
      font-size: 0.8rem;
      background: var(--mapex-surface-sunken);
      padding: 2px 6px;
      border-radius: var(--mapex-radius-sm);
    }
  }
}
</style>

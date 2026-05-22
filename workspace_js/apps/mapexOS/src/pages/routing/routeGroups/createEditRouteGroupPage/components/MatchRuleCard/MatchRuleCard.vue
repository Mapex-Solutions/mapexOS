<script setup lang="ts">
defineOptions({
  name: 'MatchRuleCard'
});

/** TYPE IMPORTS */
import type { MatchRuleCardProps, MatchRuleCardEmits } from './interfaces/MatchRuleCard.interface';
import type { MatchRule } from '@interfaces/routing/routeGroups.interface';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** PROPS & EMITS */
const props = defineProps<MatchRuleCardProps>();
const emit = defineEmits<MatchRuleCardEmits>();

// Update rule fields - handle nullable Quasar input values
function updateField<K extends keyof MatchRule>(field: K, value: MatchRule[K] | string | number | null) {
  if (value === null) return; // Ignore null values from Quasar
  emit('update:rule', { ...props.rule, [field]: value as MatchRule[K] });
}
</script>

<template>
  <q-card
    flat
    bordered
    class="bg-white"
  >
    <q-card-section class="q-pa-md">
      <div class="row q-col-gutter-md">
        <!-- Field -->
        <div class="col-12 col-sm-4">
          <q-input
            outlined
            dense
            class="rounded-borders"
            :label="t.createEdit.routersStep.routerCard.conditionalRouting.rule.field.label.value"
            :placeholder="t.createEdit.routersStep.routerCard.conditionalRouting.rule.field.placeholder.value"
            :hint="t.createEdit.routersStep.routerCard.conditionalRouting.rule.field.hint.value"
            :model-value="rule.field"
            @update:model-value="(val) => updateField('field', val)"
          >
            <template #prepend>
              <q-icon name="label" color="secondary" size="xs" />
            </template>
          </q-input>
        </div>

        <!-- Operator -->
        <div class="col-12 col-sm-3">
          <q-select
            outlined
            dense
            emit-value
            map-options
            class="rounded-borders"
            :label="t.createEdit.routersStep.routerCard.conditionalRouting.rule.operator.label.value"
            :options="matchOperatorOptions"
            :hint="t.createEdit.routersStep.routerCard.conditionalRouting.rule.operator.placeholder.value"
            option-label="label"
            option-value="value"
            :model-value="rule.operator"
            @update:model-value="(val) => updateField('operator', val)"
          >
            <template #prepend>
              <q-icon name="compare_arrows" color="secondary" size="xs" />
            </template>
          </q-select>
        </div>

        <!-- Value -->
        <div class="col-12 col-sm-4">
          <q-input
            outlined
            dense
            class="rounded-borders"
            :label="t.createEdit.routersStep.routerCard.conditionalRouting.rule.value.label.value"
            :placeholder="t.createEdit.routersStep.routerCard.conditionalRouting.rule.value.placeholder.value"
            :hint="rule.operator === 'in' || rule.operator === 'nin'
              ? t.createEdit.routersStep.routerCard.conditionalRouting.rule.value.arrayHint.value
              : t.createEdit.routersStep.routerCard.conditionalRouting.rule.value.hint.value"
            :model-value="rule.value"
            @update:model-value="(val) => updateField('value', val)"
          >
            <template #prepend>
              <q-icon name="edit" color="secondary" size="xs" />
            </template>
          </q-input>
        </div>

        <!-- Remove Rule Button -->
        <div v-if="ruleIndex > 0" class="col-12 col-sm-1 flex items-start q-pt-lg">
          <q-btn
            flat
            round
            icon="delete"
            color="negative"
            @click="emit('delete')"
          >
            <AppTooltip :content="t.createEdit.routersStep.routerCard.conditionalRouting.removeRule.value" />
          </q-btn>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'MatchConfiguration'
});

/** TYPE IMPORTS */
import type { MatchConfigurationProps, MatchConfigurationEmits } from './interfaces/MatchConfiguration.interface';
import type { MatchRule } from '@interfaces/routing/routeGroups.interface';

/** COMPONENTS */
import { MatchRuleCard } from '../MatchRuleCard';

/** PROPS & EMITS */
const props = defineProps<MatchConfigurationProps>();
const emit = defineEmits<MatchConfigurationEmits>();

// Update policy
function updatePolicy(policy: string) {
  if (props.match) {
    emit('update:match', {
      ...props.match,
      policy: policy as 'all' | 'any',
    });
  }
}

// Update rule
function updateRule(ruleIndex: number, updatedRule: MatchRule) {
  if (props.match) {
    const newRules = [...props.match.rules];
    newRules[ruleIndex] = updatedRule;
    emit('update:match', {
      ...props.match,
      rules: newRules,
    });
  }
}
</script>

<template>
  <q-card
    flat
    class="bg-grey-2"
  >
    <q-card-section class="q-pb-md">
      <!-- Match Policy -->
      <div class="q-mb-md">
        <div class="text-caption text-grey-7 q-mb-sm text-weight-medium">
          {{ t.createEdit.routersStep.routerCard.conditionalRouting.policy.label.value }}
        </div>
        <q-btn-toggle
          unelevated
          spread
          no-caps
          toggle-color="secondary"
          color="white"
          text-color="grey-8"
          :options="matchPolicyOptions"
          :model-value="match?.policy"
          @update:model-value="updatePolicy"
        />
      </div>

      <!-- Match Rules -->
      <div class="text-caption text-grey-7 q-mb-sm text-weight-medium">
        {{ t.createEdit.routersStep.routerCard.conditionalRouting.rule.field.label.value }}s
      </div>
      <div class="q-gutter-sm">
        <MatchRuleCard
          v-for="(rule, ruleIndex) in match?.rules"
          :key="ruleIndex"
          :rule="rule"
          :rule-index="ruleIndex"
          :match-operator-options="matchOperatorOptions"
          :t="t"
          @update:rule="(updatedRule) => updateRule(ruleIndex, updatedRule)"
          @delete="emit('remove-rule', ruleIndex)"
        />
      </div>

      <!-- Add Rule Button -->
      <div class="q-mt-md">
        <q-btn
          outline
          no-caps
          unelevated
          color="secondary"
          icon="add"
          :label="t.createEdit.routersStep.routerCard.conditionalRouting.addRule.value"
          @click="emit('add-rule')"
        />
      </div>
    </q-card-section>
  </q-card>
</template>

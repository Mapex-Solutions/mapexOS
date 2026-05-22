<script setup lang="ts">
defineOptions({
  name: 'Step2RoutersConfig',
});

/** TYPE IMPORTS */
import type { Step2RoutersConfigProps, Step2RoutersConfigEmits } from './interfaces/Step2RoutersConfig.interface';
import type { RouterFormState } from '../../interfaces';

/** COMPONENTS */
import { RouterCard } from '../RouterCard';

/** PROPS & EMITS */
const props = defineProps<Step2RoutersConfigProps>();
const emit = defineEmits<Step2RoutersConfigEmits>();

/** FUNCTIONS */

/**
 * Update a router in the forms array
 *
 * @param {number} index - Index of the router to update
 * @param {RouterFormState} updatedRouter - Updated router data
 */
function updateRouter(index: number, updatedRouter: RouterFormState): void {
  const newRouters = [...props.routerForms];
  newRouters[index] = updatedRouter;
  emit('update:routerForms', newRouters);
}
</script>

<template>
  <div>
    <div class="q-mb-lg">
      <div class="text-subtitle1 text-weight-medium q-mb-sm">
        <q-icon name="route" color="primary" class="q-mr-xs" />
        {{ t.createEdit.routersStep.title.value }}
      </div>
      <div class="text-body2 text-grey-7 q-mb-md">
        {{ t.createEdit.routersStep.subtitle.value }}
      </div>
    </div>

    <!-- Router Cards -->
    <div class="q-gutter-md">
      <RouterCard
        v-for="(routerForm, index) in routerForms"
        :key="routerForm.id"
        :router="routerForm"
        :index="index"
        :router-kind-options="routerKindOptions"
        :match-policy-options="matchPolicyOptions"
        :match-operator-options="matchOperatorOptions"
        :t="t"
        @update:router="(updatedRouter) => updateRouter(index, updatedRouter)"
        @remove="emit('remove-router', routerForm.id)"
        @kind-change="(kind) => emit('router-kind-change', routerForm.id, kind)"
        @toggle-conditional="(enabled) => emit('toggle-conditional-routing', routerForm.id, enabled)"
        @add-match-rule="emit('add-match-rule', routerForm.id)"
        @remove-match-rule="(ruleIndex) => emit('remove-match-rule', routerForm.id, ruleIndex)"
      />

      <!-- Add Router Button -->
      <div id="route-group-add-router" class="q-mt-md">
        <q-btn
          color="primary"
          icon="add"
          :label="t.createEdit.routersStep.addRouter.value"
          @click="emit('add-router')"
        />
      </div>

      <!-- Warning if no routers -->
      <q-banner
        v-if="routerForms.length === 0"
        rounded
        class="bg-warning text-white q-mt-md"
      >
        <template #avatar>
          <q-icon name="warning" color="white" size="sm" />
        </template>
        <div class="text-body2">
          <strong>{{ t.createEdit.routersStep.noRoutersWarning.value }}</strong> {{ t.createEdit.routersStep.noRoutersHint.value }}
        </div>
      </q-banner>
    </div>
  </div>
</template>

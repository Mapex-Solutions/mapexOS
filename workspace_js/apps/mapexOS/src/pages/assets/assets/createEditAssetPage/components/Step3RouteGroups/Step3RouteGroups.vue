<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="route" color="primary" class="q-mr-xs" />
        {{ t.steps.step3.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step3.subtitle.value }}
      </div>
    </div>

    <q-banner rounded class="bg-blue-1 text-primary q-mb-md">
      <template #avatar>
        <q-icon name="info" color="primary" />
      </template>
      <div class="text-subtitle2 text-weight-medium q-mb-xs">{{ t.steps.step3.banner.title.value }}</div>
      <div class="text-body2">
        {{ t.steps.step3.banner.description.value }}
      </div>
    </q-banner>

    <!-- Selected Route Groups Display -->
    <div class="q-mb-sm">
      <div class="text-caption text-grey-7 q-mb-xs">{{ t.steps.step3.labels.selectedHeader.value }}</div>

      <!-- Selected chips -->
      <div v-if="selectedRouteGroups.length > 0" class="row q-gutter-xs q-mb-md">
        <span v-for="routeGroup in selectedRouteGroups" :key="routeGroup.id || routeGroup.name || ''">
          <SelectableChip
            :label="routeGroup.name"
            icon="alt_route"
            color="primary"
            size="sm"
            @remove="removeRouteGroup(routeGroup.id || '')"
          >
            <AppTooltip>
              <div class="text-caption">
                <div v-if="routeGroup.isSystem"><strong>{{ t.steps.step3.labels.systemTemplate.value }}</strong></div>
                <div v-if="routeGroup.isTemplate"><strong>{{ t.steps.step3.labels.sharedTemplate.value }}</strong></div>
                <div>{{ t.steps.step3.labels.routersConfigured(routeGroup.routers?.length || 0) }}</div>
              </div>
            </AppTooltip>
          </SelectableChip>
        </span>
      </div>

      <!-- Empty state -->
      <div v-else class="text-grey-6 text-caption q-mb-md">
        {{ t.steps.step3.labels.noneSelected.value }}
      </div>

      <!-- Action buttons -->
      <div class="row q-gutter-sm">
        <q-btn
          outline
          dense
          class="rounded-borders"
          color="primary"
          icon="add"
          :label="t.steps.step3.buttons.select.value"
          size="sm"
          data-testid="asset-routegroup-select-btn"
          @click="showRouteGroupDrawer = true"
        />
        <q-btn
          v-if="selectedRouteGroups.length > 0"
          flat
          dense
          class="rounded-borders"
          color="negative"
          icon="clear"
          :label="t.steps.step3.buttons.clearAll.value"
          size="sm"
          data-testid="asset-routegroup-clear-btn"
          @click="clearAll"
        />
      </div>
    </div>

    <!-- Route Groups Details (if selected) -->
    <q-card v-if="selectedRouteGroups.length > 0" flat bordered class="q-mt-md rounded-borders" data-testid="asset-routegroup-count">
      <q-card-section class="q-pa-md">
        <div class="row items-center q-mb-sm">
          <q-icon name="alt_route" color="primary" size="sm" class="q-mr-xs" />
          <div class="text-subtitle2 text-weight-medium">
            {{ t.steps.step3.labels.selectedCount(selectedRouteGroups.length) }}
          </div>
        </div>
        <div class="text-caption text-grey-7">
          {{ t.steps.step3.labels.routingExplanation.value }}
        </div>
      </q-card-section>
    </q-card>

    <!-- Route Group Selector Drawer -->
    <RouteGroupSelectorDrawer
      v-model="showRouteGroupDrawer"
      :multi-select="true"
      @select="handleRouteGroupsSelect"
    />
  </q-form>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step3RouteGroupsProps, Step3RouteGroupsEmits } from './interfaces/Step3RouteGroups.interface';

defineOptions({
  name: 'Step3RouteGroups'
});

import type { RouteGroupResponse } from '@mapexos/schemas';
import type { QForm } from 'quasar';

import { ref, reactive, watch } from 'vue';

import { SelectableChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { RouteGroupSelectorDrawer } from '@components/drawers';

const props = defineProps<Step3RouteGroupsProps>();
const emit = defineEmits<Step3RouteGroupsEmits>();

const formRef = ref<QForm | null>(null);
const showRouteGroupDrawer = ref(false);
const selectedRouteGroups = ref<RouteGroupResponse[]>([]);

const t = useAddAssetTranslations();

const localData = reactive({
  routeGroupIds: props.modelValue.routeGroupIds || [],
});

// Initialize selectedRouteGroups from modelValue if available
watch(() => props.modelValue, (newVal) => {
  localData.routeGroupIds = newVal.routeGroupIds || [];

  // Preserve selectedRouteGroups array if it exists in state
  if (newVal.selectedRouteGroups && newVal.selectedRouteGroups.length > 0) {
    selectedRouteGroups.value = newVal.selectedRouteGroups;
  }
}, { deep: true, immediate: true });

function handleRouteGroupsSelect(routeGroups: RouteGroupResponse[]) {
  selectedRouteGroups.value = routeGroups;
  localData.routeGroupIds = routeGroups.map(rg => rg.id).filter(Boolean) as string[];

  emit('update:modelValue', {
    routeGroupIds: localData.routeGroupIds,
    selectedRouteGroups: routeGroups // Store full route group objects
  });
  emit('routeGroupsSelected', routeGroups);
}

function removeRouteGroup(routeGroupId: string) {
  selectedRouteGroups.value = selectedRouteGroups.value.filter(rg => rg.id !== routeGroupId);
  localData.routeGroupIds = selectedRouteGroups.value.map(rg => rg.id).filter(Boolean) as string[];

  emit('update:modelValue', {
    routeGroupIds: localData.routeGroupIds,
    selectedRouteGroups: selectedRouteGroups.value
  });
  emit('routeGroupsSelected', selectedRouteGroups.value);
}

function clearAll() {
  selectedRouteGroups.value = [];
  localData.routeGroupIds = [];

  emit('update:modelValue', {
    routeGroupIds: [],
    selectedRouteGroups: []
  });
  emit('routeGroupsSelected', []);
}

defineExpose({
  formRef,
});
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

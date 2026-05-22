<script setup lang="ts">
defineOptions({ name: 'Step2Definition' });

/** TYPE IMPORTS */
import type { WorkflowInstanceFormData } from '../../interfaces';
import type { DefinitionResponse } from '@mapexos/schemas';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { WorkflowDefinitionSelectorDrawer } from '@components/drawers/automations/workflowDefinitionSelectorDrawer';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: WorkflowInstanceFormData;
}>();
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<WorkflowInstanceFormData>): void;
  (e: 'definition-selected', definition: DefinitionResponse): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const drawerOpen = ref(false);
const selectedDefinition = ref<DefinitionResponse | null>(props.modelValue.selectedDefinition || null);

/** FUNCTIONS */

/**
 * Handle definition selection from drawer
 * @param {DefinitionResponse} definition - The selected definition
 * @returns {void}
 */
function onDefinitionSelected(definition: DefinitionResponse): void {
  selectedDefinition.value = definition;

  emit('update:modelValue', {
    definitionId: definition._id || null,
    definitionVersion: definition.definitionVersion || 1,
    selectedDefinition: definition,
    externalInputs: {},
  });
  emit('definition-selected', definition);
}

/**
 * Clear the current definition selection
 * @returns {void}
 */
function clearSelection(): void {
  selectedDefinition.value = null;
  emit('update:modelValue', {
    definitionId: null,
    definitionVersion: 1,
    selectedDefinition: null,
    externalInputs: {},
  });
}

/** COMPUTED */

/**
 * Translated health label based on definition status
 */
const healthLabel = computed(() => {
  const status = String(selectedDefinition.value?.status || '');
  if (status === 'plugin_missing') return t.review.healthPluginMissing.value;
  if (status === 'invalid') return t.review.healthInvalid.value;
  return t.review.healthValid.value;
});

/** WATCHERS */
watch(() => props.modelValue.selectedDefinition, (newVal) => {
  if (newVal) {
    selectedDefinition.value = newVal;
  }
}, { deep: true });

defineExpose({ formRef });
</script>

<template>
  <q-form ref="formRef" greedy>
    <div class="row q-col-gutter-md">
      <!-- Definition selector input -->
      <div class="col-12">
        <q-input
          :model-value="selectedDefinition?.name || ''"
          outlined
          dense
          readonly
          class="rounded-borders cursor-pointer"
          data-testid="instance-definition-input"
          :label="t.fields.definition.value + ' *'"
          :placeholder="t.fields.definitionPlaceholder.value"
          :rules="[() => !!selectedDefinition || t.fields.requiredField.value]"
          @click="drawerOpen = true"
        >
          <template #prepend>
            <q-icon name="account_tree" color="primary" />
          </template>
          <template #append>
            <q-icon name="chevron_right" color="grey-6" />
          </template>
        </q-input>
      </div>

      <!-- Selected definition card -->
      <div v-if="selectedDefinition" class="col-12">
        <q-card flat bordered class="definition-card">
          <!-- Header -->
          <q-card-section class="q-pb-sm">
            <div class="row items-center no-wrap">
              <q-avatar size="36px" color="primary" text-color="white" icon="account_tree" class="q-mr-sm" />
              <div class="col">
                <div class="text-subtitle1 text-weight-bold">{{ selectedDefinition.name }}</div>
                <div v-if="selectedDefinition.description" class="text-caption text-secondary ellipsis-2-lines">
                  {{ selectedDefinition.description }}
                </div>
              </div>
              <DetailChip color="blue" size="sm" dense :label="`v${selectedDefinition.definitionVersion}`" class="q-mr-xs" />
              <q-btn flat round dense icon="close" size="sm" color="grey-6" @click="clearSelection">
                <AppTooltip :content="t.review.deselect.value" />
              </q-btn>
            </div>
          </q-card-section>

          <q-separator />

          <!-- Stats grid -->
          <q-card-section class="q-py-sm">
            <div class="row q-col-gutter-sm">
              <!-- Health -->
              <div class="col-6 col-sm-3">
                <div class="stat-item">
                  <div class="stat-label">{{ t.review.health.value }}</div>
                  <DetailChip
                    dense
                    size="xs"
                    :color="selectedDefinition.status === 'valid' ? 'positive' : selectedDefinition.status === 'invalid' ? 'negative' : 'warning'"
                    :label="healthLabel"
                  />
                </div>
              </div>

              <!-- Nodes -->
              <div class="col-6 col-sm-3">
                <div class="stat-item">
                  <div class="stat-label">{{ t.review.nodes.value }}</div>
                  <div class="stat-value">
                    <q-icon name="hub" size="xs" color="grey-7" class="q-mr-xs" />
                    {{ selectedDefinition.nodes?.length || 0 }}
                  </div>
                </div>
              </div>

              <!-- States -->
              <div class="col-6 col-sm-3">
                <div class="stat-item">
                  <div class="stat-label">{{ t.review.states.value }}</div>
                  <div class="stat-value">
                    <q-icon name="data_object" size="xs" color="grey-7" class="q-mr-xs" />
                    {{ (selectedDefinition as any).states?.length || 0 }}
                  </div>
                </div>
              </div>

              <!-- External Inputs -->
              <div class="col-6 col-sm-3">
                <div class="stat-item">
                  <div class="stat-label">{{ t.review.inputs.value }}</div>
                  <div class="stat-value">
                    <q-icon name="input" size="xs" color="grey-7" class="q-mr-xs" />
                    {{ (selectedDefinition as any).externalInputs?.length || 0 }}
                  </div>
                </div>
              </div>
            </div>
          </q-card-section>

          <!-- Plugins section -->
          <q-separator />
          <q-card-section class="q-py-sm">
            <div class="stat-label q-mb-xs">{{ t.review.plugins.value }}</div>
            <div v-if="selectedDefinition.installedPlugins?.length" class="row q-gutter-xs">
              <DetailChip
                v-for="plugin in selectedDefinition.installedPlugins"
                :key="plugin"
                dense
                size="sm"
                color="purple"
                icon="extension"
                :label="plugin"
              />
            </div>
            <div v-else class="text-caption text-muted">
              <q-icon name="extension_off" size="xs" class="q-mr-xs" />
              {{ t.review.coreOnly.value }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Definition Selector Drawer (generic pattern) -->
    <WorkflowDefinitionSelectorDrawer
      v-model="drawerOpen"
      :selected-definition-id="selectedDefinition?._id ?? ''"
      @select="onDefinitionSelected"
      @cancel="drawerOpen = false"
    />
  </q-form>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.definition-card {
  border-radius: var(--mapex-radius-md);
  border-color: var(--mapex-border-color);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.stat-label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--mapex-text-secondary);
}

.stat-value {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--mapex-text-primary);
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

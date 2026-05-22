<script setup lang="ts">
defineOptions({ name: 'Step3ExternalInputs' });

/** TYPE IMPORTS */
import type { WorkflowInstanceFormData, ExternalInputDefinition } from '../../interfaces';
import type { AssetResponse } from '@mapexos/schemas';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed, reactive, watch } from 'vue';

/** COMPONENTS */
import { AssetSelectorDrawer } from '@components/drawers/assets/assetSelectorDrawer';

/** COMPOSABLES */
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: WorkflowInstanceFormData;
}>();
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<WorkflowInstanceFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const localInputs = reactive<Record<string, any>>({ ...props.modelValue.externalInputs });

// Asset selector drawer state
const assetDrawerOpen = ref(false);
const activeAssetField = ref<ExternalInputDefinition | null>(null);

// Store selected asset display names (resolved from API)
const assetDisplayNames = reactive<Record<string, string>>({});

/** COMPUTED */

/**
 * External input definitions from the selected definition
 */
const inputDefinitions = computed((): ExternalInputDefinition[] => {
  const def = props.modelValue.selectedDefinition as any;
  if (!def?.externalInputs || !Array.isArray(def.externalInputs)) return [];
  return def.externalInputs as ExternalInputDefinition[];
});

/**
 * Whether external inputs exist
 */
const hasInputs = computed(() => inputDefinitions.value.length > 0);

/** WATCHERS */
watch(() => props.modelValue.externalInputs, (newVal) => {
  Object.keys(localInputs).forEach(k => delete localInputs[k]);
  Object.assign(localInputs, newVal || {});
}, { deep: true });

// Initialize defaults + resolve asset names for edit mode
watch(inputDefinitions, (defs) => {
  for (const def of defs) {
    if (localInputs[def.field] === undefined && def.defaultValue !== undefined && def.defaultValue !== null) {
      localInputs[def.field] = def.defaultValue;
    }
  }
  updateValue();
  void resolveAssetNames();
}, { immediate: true });

/** FUNCTIONS */

/**
 * Emit updated external inputs to parent
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', { externalInputs: { ...localInputs } });
}

/**
 * Get the appropriate input type for standard fields
 * @param {string} type - The variable type
 * @returns {'text' | 'number' | 'textarea'} Input type
 */
function getInputType(type: string): 'text' | 'number' | 'textarea' {
  switch (type) {
    case 'number': return 'number';
    case 'json': return 'textarea';
    default: return 'text';
  }
}

/**
 * Open asset selector drawer for a specific field
 * @param {ExternalInputDefinition} inputDef - The field definition
 * @returns {void}
 */
function openAssetSelector(inputDef: ExternalInputDefinition): void {
  activeAssetField.value = inputDef;
  assetDrawerOpen.value = true;
}

/**
 * Handle asset selection from drawer.
 * Extracts the value from the asset using the configured fieldPath.
 * @param {AssetResponse} asset - Selected asset
 * @returns {void}
 */
function onAssetSelected(asset: AssetResponse): void {
  if (!activeAssetField.value) return;

  const field = activeAssetField.value.field;
  const fieldPath = activeAssetField.value.fieldPath || 'assetUUID';

  // Extract the value from the asset based on fieldPath
  localInputs[field] = (asset as any)[fieldPath] || asset.assetUUID || '';
  assetDisplayNames[field] = asset.name || asset.assetUUID || '';

  activeAssetField.value = null;
  updateValue();
}

/**
 * Resolve asset display names for fields that already have values (edit mode).
 * For each assetFromTemplate field with a value, fetches the asset by UUID to get its name.
 * @returns {Promise<void>}
 */
async function resolveAssetNames(): Promise<void> {
  const defs = inputDefinitions.value;

  for (const def of defs) {
    if (def.type !== 'assetFromTemplate') continue;

    const currentValue = localInputs[def.field];
    if (!currentValue || assetDisplayNames[def.field]) continue;

    try {
      // The value stored is the asset UUID (or whatever fieldPath points to)
      // We need to find the asset to show its name
      const response = await apis.assets.asset.list({
        assetUUID: currentValue,
        perPage: 1,
      });

      if (response.items?.length > 0) {
        const asset = response.items[0];
        assetDisplayNames[def.field] = asset?.name || currentValue;
      } else {
        assetDisplayNames[def.field] = currentValue;
      }
    } catch {
      // Fallback to raw value if lookup fails
      assetDisplayNames[def.field] = currentValue;
    }
  }
}

defineExpose({ formRef });
</script>

<template>
  <q-form ref="formRef" greedy>
    <!-- No external inputs -->
    <div v-if="!hasInputs" class="text-center q-pa-xl">
      <q-icon name="check_circle" color="positive" size="48px" class="q-mb-md" />
      <div class="text-h6 text-secondary">{{ t.fields.noExternalInputs.value }}</div>
    </div>

    <!-- External inputs form -->
    <div v-else>
      <div class="q-mb-md">
        <div class="text-subtitle1 text-weight-medium q-mb-xs">
          <q-icon name="input" color="primary" class="q-mr-xs" />
          {{ t.fields.externalInputsTitle.value }}
        </div>
        <div class="text-body2 text-secondary">
          {{ t.fields.externalInputsDescription.value }}
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <div
          v-for="inputDef in inputDefinitions"
          :key="inputDef.field"
          class="col-12"
        >
          <!-- Boolean type → toggle -->
          <q-item
            v-if="inputDef.type === 'boolean'"
            tag="label"
            class="rounded-borders input-card"
          >
            <q-item-section avatar>
              <q-icon :name="inputDef.icon || 'toggle_on'" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium">
                {{ inputDef.label }}
                <span v-if="inputDef.required" class="text-negative q-ml-xs">*</span>
              </q-item-label>
              <q-item-label v-if="inputDef.description" caption>{{ inputDef.description }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-toggle
                v-model="localInputs[inputDef.field]"
                color="primary"
                @update:model-value="updateValue"
              />
            </q-item-section>
          </q-item>

          <!-- assetFromTemplate → readonly input + asset selector drawer -->
          <q-input
            v-else-if="inputDef.type === 'assetFromTemplate'"
            :model-value="assetDisplayNames[inputDef.field] || localInputs[inputDef.field] || ''"
            outlined
            dense
            readonly
            class="rounded-borders cursor-pointer"
            :label="inputDef.label + (inputDef.required ? ' *' : '')"
            :hint="inputDef.description || t.fields.optionalField.value"
            :rules="inputDef.required ? [(val: any) => !!val || t.fields.requiredField.value] : []"
            @click="openAssetSelector(inputDef)"
          >
            <template #prepend>
              <q-icon :name="inputDef.icon || 'devices'" color="primary" />
            </template>
            <template #append>
              <q-icon name="chevron_right" color="grey-6" />
            </template>
          </q-input>

          <!-- String / Number / JSON → standard input -->
          <q-input
            v-else
            v-model="localInputs[inputDef.field]"
            outlined
            dense
            class="rounded-borders"
            :type="getInputType(inputDef.type)"
            :label="inputDef.label + (inputDef.required ? ' *' : '')"
            :hint="inputDef.description || (inputDef.required ? t.fields.requiredField.value : t.fields.optionalField.value)"
            :rules="inputDef.required ? [(val: any) => (val !== '' && val !== null && val !== undefined) || t.fields.requiredField.value] : []"
            @update:model-value="updateValue"
          >
            <template #prepend>
              <q-icon :name="inputDef.icon || 'input'" color="primary" />
            </template>
          </q-input>
        </div>
      </div>
    </div>

    <!-- Asset Selector Drawer -->
    <AssetSelectorDrawer
      v-model="assetDrawerOpen"
      :asset-template-id="activeAssetField?.assetTemplateId ?? ''"
      @select="onAssetSelected"
      @cancel="assetDrawerOpen = false"
    />
  </q-form>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.input-card {
  background: rgba(var(--q-primary-rgb), 0.04);
  border: 1px solid rgba(var(--q-primary-rgb), 0.12);
  border-radius: var(--mapex-radius-md);
}

.text-secondary {
  color: var(--mapex-text-secondary);
}
</style>

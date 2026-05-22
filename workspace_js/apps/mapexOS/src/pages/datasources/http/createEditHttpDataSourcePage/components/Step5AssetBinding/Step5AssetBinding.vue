<template>
  <div class="row q-col-gutter-md">
    <!-- Binding Mode Selection -->
    <div class="col-12">
      <div class="row items-center q-mb-sm">
        <q-icon name="link" color="primary" class="q-mr-xs" />
        <div class="text-subtitle2 text-weight-medium">{{ t.assetBinding.title.value }}</div>
      </div>
      <q-select
        v-model="localData.bindingMode"
        outlined
        dense
        emit-value
        map-options
        class="rounded-borders"
        :label="`${t.assetBinding.bindingMode.value} *`"
        :options="ASSET_BINDING_OPTIONS"
        option-label="label"
        option-value="value"
        :rules="[(val: any) => !!val || t.assetBinding.bindingModeRequired.value]"
        @update:model-value="handleBindingModeChange"
      >
        <template #prepend>
          <q-icon name="settings_input_component" />
        </template>
      </q-select>
    </div>

    <!-- ========================================== -->
    <!-- DIRECT MODE: Fixed Asset                  -->
    <!-- ========================================== -->
    <template v-if="localData.bindingMode === 'fixedAssetId'">
      <div class="col-12">
        <q-banner dense class="bg-blue-1 text-blue-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="info" color="blue-6" />
          </template>
          <div class="text-caption">
            {{ t.assetBinding.fixedAsset.banner.value }}
          </div>
        </q-banner>

        <AssetSelector
          v-model="localData.directAssetId"
          :label="`${t.assetBinding.fixedAsset.selectAsset.value} *`"
          :required="true"
          @update:model-value="updateValue"
          @update:selected-asset="handleSelectedAssetChange"
        />
      </div>
    </template>

    <!-- ========================================== -->
    <!-- DYNAMIC MODE: UUID Field Mapping          -->
    <!-- ========================================== -->
    <template v-if="localData.bindingMode === 'uuidField'">
      <div class="col-12">
        <q-banner dense class="bg-purple-1 text-purple-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="info" color="purple-6" />
          </template>
          <div class="text-caption">
            {{ t.assetBinding.uuidField.banner.value }}
          </div>
        </q-banner>
      </div>

      <!-- Path entries -->
      <div
        v-for="(mapping, idx) in localData.customUuidPaths"
        :key="idx"
        class="col-12"
      >
        <q-card flat bordered class="rounded-borders">
          <q-card-section class="q-pa-md">
            <div class="row items-start q-col-gutter-sm">
              <div class="col">
                <q-input
                  v-model="mapping.path"
                  outlined
                  dense
                  :label="`${t.assetBinding.uuidField.uuidJsonPath.value} *`"
                  :placeholder="t.assetBinding.uuidField.uuidJsonPathPlaceholder.value"
                  :hint="t.assetBinding.uuidField.uuidJsonPathHint.value"
                  class="rounded-borders"
                  :rules="[
                    (val: any) => !!val || t.assetBinding.uuidField.pathRequired.value,
                    (val: any) => !val || /^[a-zA-Z0-9_.]+$/.test(val) || t.assetBinding.uuidField.invalidPathFormat.value
                  ]"
                  @update:model-value="updateValue"
                >
                  <template #prepend>
                    <q-icon name="mdi-code-json" color="primary" />
                  </template>
                </q-input>

                <!-- Select path from template button -->
                <q-btn
                  flat
                  dense
                  no-caps
                  size="sm"
                  icon="description"
                  color="purple-7"
                  class="q-mt-xs"
                  :label="t.assetBinding.uuidField.selectFromTemplate.value"
                  @click="openTemplatePicker(idx)"
                />
              </div>

              <div v-if="localData.customUuidPaths.length > 1" class="col-auto q-pt-xs">
                <q-btn
                  flat
                  dense
                  round
                  icon="delete"
                  color="negative"
                  size="sm"
                  @click="removeCustomPath(idx)"
                >
                  <AppTooltip :content="t.assetBinding.uuidField.removePath.value" />
                </q-btn>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12">
        <q-btn
          flat
          dense
          icon="add"
          :label="t.assetBinding.uuidField.addPath.value"
          color="primary"
          size="sm"
          @click="addCustomPath"
        />
      </div>

      <!-- Template drawer (lateral, right side) — fills the active path on selection -->
      <AssetTemplateSelectorDrawer
        v-model="showTemplatePicker"
        :multi-select="false"
        @select="handleTemplateSelect"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step5AssetBinding'
});

/** TYPE IMPORTS */
import type { AssetResponse, AssetTemplateResponse } from '@mapexos/schemas';
import type { StepEmits, StepProps } from '../../interfaces/httpDataSource.interface';

/** VUE IMPORTS */
import { reactive, watch, computed, ref } from 'vue';

/** COMPONENTS */
import AssetSelector from '@components/selectors/assetSelector/AssetSelector.vue';
import { AssetTemplateSelectorDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

/** UTILS */
import { notifyFail } from '@utils/alert/notify';

/** LOCAL IMPORTS */
import { ASSET_BINDING_OPTIONS } from '../../constants/httpDataSourceConstants';

/** PROPS & EMITS */
const props = defineProps<StepProps>();
const emit = defineEmits<StepEmits>();

/** COMPOSABLES & STORES */
const t = useHttpDataSourceCreateEditTranslations();

/** STATE */
const showTemplatePicker = ref(false);
const activePathIndex = ref(0);

const localData = reactive({
  bindingMode: props.modelValue.bindingMode || null,
  directAssetId: props.modelValue.directAssetId || null,
  directAssetIdPath: props.modelValue.directAssetIdPath || null,
  customUuidPaths: props.modelValue.customUuidPaths || [{ path: '' }],
  payloadExample: props.modelValue.payloadExample || '',
});

/** COMPUTED */
const finalUuidPaths = computed(() =>
  localData.customUuidPaths
    .map(c => c.path.trim())
    .filter((p, i, arr) => p !== '' && arr.indexOf(p) === i)
);

/** WATCHERS */
watch(() => props.modelValue, (newVal) => {
  localData.bindingMode = newVal.bindingMode || null;
  localData.directAssetId = newVal.directAssetId || null;
  localData.directAssetIdPath = newVal.directAssetIdPath || null;
  localData.customUuidPaths = newVal.customUuidPaths || [{ path: '' }];
  localData.payloadExample = newVal.payloadExample || '';
}, { deep: true, immediate: true });

/** FUNCTIONS */

/**
 * Open the template drawer for a specific path row
 * @param {number} index - Index of the path row to fill on selection
 * @returns {void}
 */
function openTemplatePicker(index: number): void {
  activePathIndex.value = index;
  showTemplatePicker.value = true;
}

/**
 * Handle template selection from the drawer.
 * Extracts assetIdPath and fills the active path row, then closes the drawer.
 * @param {AssetTemplateResponse[]} templates - Selected templates (single select)
 * @returns {void}
 */
function handleTemplateSelect(templates: AssetTemplateResponse[]): void {
  const template = templates[0];
  if (!template) return;

  if (!template.assetIdPath) {
    notifyFail({ message: t.assetBinding.uuidField.noTemplatePath.value });
    return;
  }

  const entry = localData.customUuidPaths[activePathIndex.value];
  if (!entry) return;

  entry.path = template.assetIdPath;
  showTemplatePicker.value = false;
  updateValue();
}

/**
 * Handle binding mode change — resets mode-specific fields
 * @returns {void}
 */
function handleBindingModeChange(): void {
  if (localData.bindingMode === 'fixedAssetId') {
    localData.customUuidPaths = [{ path: '' }];
    localData.payloadExample = '';
  } else if (localData.bindingMode === 'uuidField') {
    localData.directAssetId = null;
    localData.directAssetIdPath = null;
  }
  updateValue();
}

/**
 * Handle selected asset change — extracts assetIdPath for backend use
 * @param {AssetResponse | null} asset - Selected asset or null
 * @returns {void}
 */
function handleSelectedAssetChange(asset: AssetResponse | null): void {
  localData.directAssetIdPath = asset?.assetIdPath ?? null;
  updateValue();
}

/**
 * Add a new empty path entry
 * @returns {void}
 */
function addCustomPath(): void {
  localData.customUuidPaths.push({ path: '' });
  updateValue();
}

/**
 * Remove a path entry by index
 * @param {number} index - Index of the path to remove
 * @returns {void}
 */
function removeCustomPath(index: number): void {
  if (localData.customUuidPaths.length > 1) {
    localData.customUuidPaths.splice(index, 1);
    updateValue();
  }
}

/**
 * Emit updated values to parent component
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', {
    ...props.modelValue,
    ...localData,
    finalUuidPaths: finalUuidPaths.value,
  });
}
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

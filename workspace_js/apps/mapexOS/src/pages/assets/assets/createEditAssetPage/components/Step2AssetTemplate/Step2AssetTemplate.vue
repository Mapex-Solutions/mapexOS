<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="mdi-file-document" color="primary" class="q-mr-xs" />
        {{ t.steps.step2.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step2.subtitle.value }}
      </div>
    </div>

    <!-- Selected Template Display -->
    <q-input
      :model-value="selectedTemplate?.name || ''"
      outlined
      dense
      readonly
      class="rounded-borders cursor-pointer"
      data-testid="asset-template-input"
      :label="`${t.steps.step2.fields.assetTemplateId.label.value} *`"
      :placeholder="t.steps.step2.fields.assetTemplateId.placeholder.value"
      :rules="[(val: any) => !!val || t.steps.step2.fields.assetTemplateId.required.value]"
      @click="showTemplateDrawer = true"
    >
      <template #prepend>
        <q-icon name="description" color="primary" />
      </template>
      <template #append>
        <q-icon
          v-if="selectedTemplate"
          name="close"
          class="cursor-pointer"
          data-testid="asset-template-clear-btn"
          @click.stop="clearTemplate"
        >
          <AppTooltip :content="t.steps.step2.fields.assetTemplateId.clearTooltip.value" />
        </q-icon>
        <q-icon
          name="search"
          class="cursor-pointer"
          @click.stop="showTemplateDrawer = true"
        >
          <AppTooltip :content="t.steps.step2.fields.assetTemplateId.searchTooltip.value" />
        </q-icon>
      </template>
    </q-input>

    <!-- Template Details (if selected) -->
    <q-card v-if="selectedTemplate" flat bordered class="q-mt-md rounded-borders" data-testid="asset-template-card">
      <q-card-section class="q-pa-md">
        <div class="row items-center q-mb-sm">
          <q-icon name="description" color="primary" size="sm" class="q-mr-xs" />
          <div class="text-subtitle2 text-weight-medium">{{ t.steps.step2.preview.detailsTitle.value }}</div>
        </div>
        <div class="text-body2 q-mb-xs">
          <strong>{{ t.steps.step2.preview.nameLabel.value }}:</strong> {{ selectedTemplate.name }}
        </div>
        <div v-if="selectedTemplate.description" class="text-body2 q-mb-xs">
          <strong>{{ t.steps.step2.preview.descriptionLabel.value }}:</strong> {{ selectedTemplate.description }}
        </div>
        <div v-if="selectedTemplate.assetIdPath" class="text-body2 q-mb-sm">
          <strong>{{ t.steps.step2.preview.uuidPathLabel.value }}:</strong> <code class="path-code">{{ selectedTemplate.assetIdPath }}</code>
        </div>

        <!-- Chips Section -->
        <div class="row q-gutter-xs">
          <!-- Status Chip -->
          <DetailChip
            :color="selectedTemplate.enabled ? 'positive' : 'grey'"
            size="sm"
            :label="selectedTemplate.enabled ? t.steps.step2.preview.statusActive.value : t.steps.step2.preview.statusInactive.value"
          />

          <!-- Assets IoT Chip (if template has this category indicator) -->
          <DetailChip
            v-if="selectedTemplate.categoryName"
            icon="category"
            color="blue"
            size="sm"
            :label="selectedTemplate.categoryName"
          />

          <!-- Manufacturer Chip -->
          <DetailChip
            v-if="selectedTemplate.manufacturerName"
            icon="factory"
            color="purple"
            size="sm"
            :label="selectedTemplate.manufacturerName"
          />

          <!-- Model/Template Type Chip -->
          <DetailChip
            v-if="selectedTemplate.modelName"
            icon="precision_manufacturing"
            color="teal"
            size="sm"
            :label="selectedTemplate.modelName"
          />
        </div>
      </q-card-section>
    </q-card>

    <!-- Asset Template Selector Drawer -->
    <AssetTemplateSelectorDrawer
      v-model="showTemplateDrawer"
      :multi-select="false"
      @select="handleTemplateSelect"
    />
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2AssetTemplate'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { Step2AssetTemplateProps, Step2AssetTemplateEmits } from './interfaces/Step2AssetTemplate.interface';
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, reactive, watch } from 'vue';

/** COMPONENTS */
import { AssetTemplateSelectorDrawer } from '@components/drawers';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';

/** PROPS & EMITS */
const props = defineProps<Step2AssetTemplateProps>();
const emit = defineEmits<Step2AssetTemplateEmits>();

/** COMPOSABLES & STORES */
const t = useAddAssetTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showTemplateDrawer = ref(false);
const selectedTemplate = ref<AssetTemplateResponse | null>(null);
const localData = reactive({
  assetTemplateId: props.modelValue.assetTemplateId || '',
});

/** WATCHERS */

/**
 * Watch for changes in modelValue to sync local state
 * Preserves selectedTemplate object if it exists in state
 */
watch(() => props.modelValue, (newVal) => {
  localData.assetTemplateId = newVal.assetTemplateId || '';

  if (newVal.selectedTemplate) {
    selectedTemplate.value = newVal.selectedTemplate;
  }
}, { deep: true, immediate: true });

/** FUNCTIONS */

/**
 * Handle template selection from drawer
 * Updates local state and emits changes to parent
 *
 * @param {AssetTemplateResponse[]} templates - Selected templates array (single select)
 */
function handleTemplateSelect(templates: AssetTemplateResponse[]): void {
  const template = templates[0];
  if (!template) return;

  selectedTemplate.value = template;
  localData.assetTemplateId = template.id || '';

  emit('update:modelValue', {
    assetTemplateId: localData.assetTemplateId,
    selectedTemplate: template
  });
  emit('templateSelected', template);
}

/**
 * Clear selected template
 * Resets local state and emits changes to parent
 */
function clearTemplate(): void {
  selectedTemplate.value = null;
  localData.assetTemplateId = '';

  emit('update:modelValue', {
    assetTemplateId: '',
    selectedTemplate: null
  });
  emit('templateSelected', null);
}

/** EXPOSE */
defineExpose({
  formRef,
});
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
}

.path-code {
  background-color: var(--mapex-submenu-bg);
  padding: 2px 6px;
  border-radius: var(--mapex-radius-xs);
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
  color: var(--q-info);
}
</style>

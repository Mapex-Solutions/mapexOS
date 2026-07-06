<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="devices" color="primary" class="q-mr-xs" />
        {{ t.steps.devices.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.devices.subtitle.value }}
      </div>
    </div>

    <div class="row items-center q-mb-sm">
      <div class="col text-caption text-grey-7">
        {{ modelValue.selectedAssetIds.length }} {{ t.devices.selected.value }}
      </div>
      <div class="col-auto">
        <q-checkbox
          :model-value="allSelected"
          dense
          :label="t.devices.selectAll.value"
          @update:model-value="toggleSelectAll"
        />
      </div>
    </div>

    <div v-if="loadingAssets" class="row justify-center q-my-md">
      <q-spinner color="primary" size="2em" />
    </div>

    <div v-else-if="assets.length === 0" class="text-caption text-grey-7 q-pa-md text-center">
      {{ t.devices.noAssets.value }}
    </div>

    <q-list v-else bordered separator class="rounded-borders device-list">
      <q-item v-for="asset in assets" :key="asset.id" dense>
        <q-item-section side>
          <q-checkbox
            :model-value="modelValue.selectedAssetIds.includes(asset.id)"
            dense
            @update:model-value="(checked) => toggleAsset(asset.id, checked)"
          />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ asset.name }}</q-item-label>
          <q-item-label caption class="text-mono">{{ asset.assetUUID }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'StepDevices'
});

/** TYPE IMPORTS */
import type { OtaPlanFormData, SelectableAsset } from '../../interfaces/createEditOtaPlan.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPOSABLES */
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import { ASSETS_FETCH_PER_PAGE } from '../../constants/createEditOtaPlan.constant';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: OtaPlanFormData;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: OtaPlanFormData];
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();
const logger = useLogger('StepDevices');

/** STATE */
const assets = ref<SelectableAsset[]>([]);
const loadingAssets = ref(false);

/** COMPUTED */

/** All listed devices are selected */
const allSelected = computed(
  () => assets.value.length > 0 && props.modelValue.selectedAssetIds.length === assets.value.length,
);

/** FUNCTIONS */

/**
 * Merge a partial change into the form state
 */
function update(partial: Partial<OtaPlanFormData>): void {
  emit('update:modelValue', { ...props.modelValue, ...partial });
}

/**
 * Toggle one device in the selection
 */
function toggleAsset(assetId: string, checked: boolean): void {
  const selected = new Set(props.modelValue.selectedAssetIds);
  if (checked) {
    selected.add(assetId);
  } else {
    selected.delete(assetId);
  }
  update({ selectedAssetIds: Array.from(selected) });
}

/**
 * Select or clear every listed device
 */
function toggleSelectAll(checked: boolean): void {
  update({ selectedAssetIds: checked ? assets.value.map((a) => a.id) : [] });
}

/**
 * Load the devices currently bound to the source template
 */
async function fetchAssets(templateId: string): Promise<void> {
  try {
    loadingAssets.value = true;
    const response = await apis.assets.asset.list({
      page: 1,
      perPage: ASSETS_FETCH_PER_PAGE,
      assetTemplateId: templateId,
    });
    assets.value = response.items.map((a) => ({
      id: a.id || '',
      name: a.name || '',
      assetUUID: a.assetUUID || '',
    }));
  } catch (err: unknown) {
    logger.error('Error fetching assets:', err);
    assets.value = [];
  } finally {
    loadingAssets.value = false;
  }
}

/** WATCHERS */

// A different source template resets the selection and reloads the candidates.
watch(
  () => props.modelValue.sourceTemplateId,
  (templateId, previous) => {
    if (templateId === previous) return;
    if (previous !== undefined) {
      update({ selectedAssetIds: [] });
    }
    assets.value = [];
    if (templateId) {
      void fetchAssets(templateId);
    }
  },
);

/** LIFECYCLE HOOKS */
onMounted(async () => {
  if (props.modelValue.sourceTemplateId) {
    await fetchAssets(props.modelValue.sourceTemplateId);
  }
});
</script>

<style lang="scss" scoped>
.device-list {
  max-height: 360px;
  overflow-y: auto;
}

.text-mono {
  font-family: monospace;
}
</style>

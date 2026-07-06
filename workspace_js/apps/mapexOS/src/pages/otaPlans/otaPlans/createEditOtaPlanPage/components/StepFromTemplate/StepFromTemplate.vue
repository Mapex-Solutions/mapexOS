<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepFromTemplateProps, StepFromTemplateEmits } from './interfaces';

defineOptions({
  name: 'StepFromTemplate'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AssetTemplateSelector } from '@components/selectors/assetTemplateSelector';

/** COMPOSABLES */
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';

/** PROPS & EMITS */
const props = defineProps<StepFromTemplateProps>();
const emit = defineEmits<StepFromTemplateEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();

/** COMPUTED */
const selectedTemplateId = computed({
  get: () => props.modelValue,
  set: (value: string | null) => emit('update:modelValue', value),
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="swap_horiz" color="primary" class="q-mr-xs" />
        {{ t.steps.fromTemplate.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.fromTemplate.subtitle.value }}
      </div>
    </div>

    <AssetTemplateSelector v-model="selectedTemplateId" />
  </div>
</template>

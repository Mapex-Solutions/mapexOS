<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepToTemplateProps, StepToTemplateEmits } from './interfaces';

defineOptions({
  name: 'StepToTemplate'
});

/** VUE IMPORTS */
import { computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { AssetTemplateSelector } from '@components/selectors/assetTemplateSelector';

/** COMPOSABLES */
import { useCreateEditMigrationPlanTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<StepToTemplateProps>();
const emit = defineEmits<StepToTemplateEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditMigrationPlanTranslations();

/** COMPUTED */
// Any template is a valid target, including the source template.
const selectedTemplateId = computed({
  get: () => props.modelValue,
  set: (value: string | null) => emit('update:modelValue', value),
});

/** WATCHERS */
watch(
  () => props.modelValue,
  value => emit('update:valid', !!value)
);

/** LIFECYCLE HOOKS */
onMounted(() => emit('update:valid', !!props.modelValue));
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="swap_horiz" color="primary" class="q-mr-xs" />
        {{ t.steps.toTemplate.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.toTemplate.subtitle.value }}
      </div>
    </div>

    <AssetTemplateSelector v-model="selectedTemplateId" />
  </div>
</template>

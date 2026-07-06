<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepDetailsProps, StepDetailsEmits } from './interfaces';

defineOptions({
  name: 'StepDetails'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useCreateEditMigrationPlanTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<StepDetailsProps>();
const emit = defineEmits<StepDetailsEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditMigrationPlanTranslations();

/** COMPUTED */
const name = computed({
  get: () => props.name,
  set: (value: string) => emit('update:name', value),
});

const description = computed({
  get: () => props.description,
  set: (value: string) => emit('update:description', value),
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="info" color="primary" class="q-mr-xs" />
        {{ t.steps.details.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.details.subtitle.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="name"
          outlined
          dense
          class="rounded-borders"
          :label="t.steps.details.fields.name.label.value + ' *'"
          :placeholder="t.steps.details.fields.name.placeholder.value"
          :hint="t.steps.details.fields.name.hint.value"
          :rules="[(val) => !!val || t.steps.details.fields.name.required.value]"
        >
          <template v-slot:prepend>
            <q-icon name="label" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="description"
          outlined
          dense
          type="textarea"
          rows="3"
          class="rounded-borders"
          :label="t.steps.details.fields.description.label.value"
          :placeholder="t.steps.details.fields.description.placeholder.value"
          :hint="t.steps.details.fields.description.hint.value"
        >
          <template v-slot:prepend>
            <q-icon name="notes" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

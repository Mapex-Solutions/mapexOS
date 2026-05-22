<script setup lang="ts">
defineOptions({
  name: 'TeamsConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** PROPS & EMITS */
const props = defineProps<{ modelValue: Trigger }>();
const emit = defineEmits<{ 'update:modelValue': [value: Trigger] }>();

/** COMPOSABLES & STORES */
const translations = useCreateEditTriggerTranslations();

/** STATE */
const config = ref({
  webhookUrl: props.modelValue.config?.webhookUrl || '',
  title: props.modelValue.config?.title || '',
  text: props.modelValue.config?.text || '',
  themeColor: props.modelValue.config?.themeColor || '0078D7',
});

/** WATCHERS */
watch(config, (newConfig) => {
  emit('update:modelValue', { ...props.modelValue, config: newConfig });
}, { deep: true });
</script>

<template>
  <div class="teams-config">
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="config.webhookUrl"
          outlined
          dense
          :label="translations.step4Configs.teams.webhookUrlLabel.value"
          :placeholder="translations.step4Configs.teams.webhookUrlPlaceholder.value"
          :hint="translations.step4Configs.teams.webhookUrlHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.teams.webhookUrlRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="link" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.title"
          outlined
          dense
          :label="translations.step4Configs.teams.titleLabel.value"
          :placeholder="translations.step4Configs.teams.titlePlaceholder.value"
          :rules="[(val: string) => !!val || translations.step4Configs.teams.titleRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="title" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.text"
          outlined
          dense
          type="textarea"
          :label="translations.step4Configs.teams.textLabel.value"
          :placeholder="translations.step4Configs.teams.textPlaceholder.value"
          :hint="translations.step4Configs.teams.textHint.value"
          rows="8"
          :rules="[(val: string) => !!val || translations.step4Configs.teams.textRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="message" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.themeColor"
          outlined
          dense
          :label="translations.step4Configs.teams.themeColorLabel.value"
          :placeholder="translations.step4Configs.teams.themeColorPlaceholder.value"
          :hint="translations.step4Configs.teams.themeColorHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="palette" />
          </template>
          <template v-slot:append>
            <div
              class="color-preview"
              :style="{ backgroundColor: `#${config.themeColor}` }"
            ></div>
          </template>
        </q-input>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.teams-config {
  .color-preview {
    width: 24px;
    height: 24px;
    border-radius: var(--mapex-radius-xs);
    border: 1px solid var(--mapex-card-border);
  }
}
</style>

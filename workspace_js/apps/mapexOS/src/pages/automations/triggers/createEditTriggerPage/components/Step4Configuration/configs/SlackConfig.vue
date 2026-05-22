<script setup lang="ts">
defineOptions({
  name: 'SlackConfig'
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
  channel: props.modelValue.config?.channel || '',
  username: props.modelValue.config?.username || 'MapexOS',
  iconEmoji: props.modelValue.config?.iconEmoji || ':robot_face:',
  message: props.modelValue.config?.message || '',
});

/** WATCHERS */
watch(config, (newConfig) => {
  emit('update:modelValue', { ...props.modelValue, config: newConfig });
}, { deep: true });
</script>

<template>
  <div class="slack-config">
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="config.webhookUrl"
          outlined
          dense
          :label="translations.step4Configs.slack.webhookUrlLabel.value"
          :placeholder="translations.step4Configs.slack.webhookUrlPlaceholder.value"
          :hint="translations.step4Configs.slack.webhookUrlHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.slack.webhookUrlRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="mdi-slack" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.channel"
          outlined
          dense
          :label="translations.step4Configs.slack.channelLabel.value"
          :placeholder="translations.step4Configs.slack.channelPlaceholder.value"
          :hint="translations.step4Configs.slack.channelHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="tag" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.username"
          outlined
          dense
          :label="translations.step4Configs.slack.usernameLabel.value"
          :placeholder="translations.step4Configs.slack.usernamePlaceholder.value"
          :hint="translations.step4Configs.slack.usernameHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="person" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.iconEmoji"
          outlined
          dense
          :label="translations.step4Configs.slack.iconEmojiLabel.value"
          :placeholder="translations.step4Configs.slack.iconEmojiPlaceholder.value"
          :hint="translations.step4Configs.slack.iconEmojiHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="sentiment_satisfied_alt" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.message"
          outlined
          dense
          type="textarea"
          :label="translations.step4Configs.slack.messageLabel.value"
          :placeholder="translations.step4Configs.slack.messagePlaceholder.value"
          :hint="translations.step4Configs.slack.messageHint.value"
          rows="8"
          :rules="[(val: string) => !!val || translations.step4Configs.slack.messageRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="message" />
          </template>
        </q-input>
      </div>
    </div>
  </div>
</template>

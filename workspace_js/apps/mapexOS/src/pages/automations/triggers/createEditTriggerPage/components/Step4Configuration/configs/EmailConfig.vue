<script setup lang="ts">
defineOptions({
  name: 'EmailConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Trigger];
}>();

/** COMPOSABLES & STORES */
const translations = useCreateEditTriggerTranslations();

/** STATE */
const config = ref({
  smtpHost: props.modelValue.config?.smtpHost || '',
  smtpPort: props.modelValue.config?.smtpPort ?? 587,
  username: props.modelValue.config?.username || '',
  password: props.modelValue.config?.password || '',
  fromAddr: props.modelValue.config?.fromAddr || '',
  to: props.modelValue.config?.to || '',
  cc: props.modelValue.config?.cc || '',
  bcc: props.modelValue.config?.bcc || '',
  subject: props.modelValue.config?.subject || '',
  body: props.modelValue.config?.body || '',
  htmlBody: props.modelValue.config?.htmlBody || '',
});

/** CONSTANTS */
const OPTIONAL_STRING_FIELDS = ['cc', 'bcc', 'username', 'password', 'body', 'htmlBody'] as const;

/** WATCHERS */
watch(
  config,
  (newConfig) => {
    // Strip empty optional fields — backend Zod schema rejects "" (StringAndNotBeEmptyOrOptional accepts undefined, not empty)
    const cleaned: Record<string, unknown> = { ...newConfig };
    for (const key of OPTIONAL_STRING_FIELDS) {
      if (!cleaned[key]) delete cleaned[key];
    }
    emit('update:modelValue', {
      ...props.modelValue,
      config: cleaned,
    });
  },
  { deep: true }
);
</script>

<template>
  <div class="email-config">
    <!-- SMTP Server Section -->
    <div class="text-subtitle2 text-weight-medium text-primary q-mb-sm">
      <q-icon name="dns" class="q-mr-xs" />
      {{ translations.step4Configs.email.smtpSectionTitle.value }}
    </div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-8">
        <q-input
          v-model="config.smtpHost"
          outlined
          dense
          :label="translations.step4Configs.email.smtpHostLabel.value"
          :placeholder="translations.step4Configs.email.smtpHostPlaceholder.value"
          :hint="translations.step4Configs.email.smtpHostHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.email.smtpHostRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="dns" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-4">
        <q-input
          v-model.number="config.smtpPort"
          outlined
          dense
          type="number"
          :label="translations.step4Configs.email.smtpPortLabel.value"
          :placeholder="translations.step4Configs.email.smtpPortPlaceholder.value"
          :hint="translations.step4Configs.email.smtpPortHint.value"
          :rules="[
            (val: number) => (val !== null && val !== undefined && val !== 0) || translations.step4Configs.email.smtpPortRequired.value,
            (val: number) => (val >= 1 && val <= 65535) || translations.step4Configs.email.smtpPortInvalid.value,
          ]"
        >
          <template v-slot:prepend>
            <q-icon name="lan" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.fromAddr"
          outlined
          dense
          :label="translations.step4Configs.email.fromAddrLabel.value"
          :placeholder="translations.step4Configs.email.fromAddrPlaceholder.value"
          :hint="translations.step4Configs.email.fromAddrHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.email.fromAddrRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="alternate_email" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.username"
          outlined
          dense
          autocomplete="off"
          :label="translations.step4Configs.email.usernameLabel.value"
          :placeholder="translations.step4Configs.email.usernamePlaceholder.value"
          :hint="translations.step4Configs.email.usernameHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="person" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.password"
          outlined
          dense
          type="password"
          autocomplete="new-password"
          :label="translations.step4Configs.email.passwordLabel.value"
          :placeholder="translations.step4Configs.email.passwordPlaceholder.value"
          :hint="translations.step4Configs.email.passwordHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
        </q-input>
      </div>
    </div>

    <!-- Email Content Section -->
    <div class="text-subtitle2 text-weight-medium text-primary q-mb-sm">
      <q-icon name="mail" class="q-mr-xs" />
      {{ translations.step4Configs.email.contentSectionTitle.value }}
    </div>
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="config.to"
          outlined
          dense
          :label="translations.step4Configs.email.toLabel.value"
          :placeholder="translations.step4Configs.email.toPlaceholder.value"
          :hint="translations.step4Configs.email.toHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.email.toRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="email" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.cc"
          outlined
          dense
          :label="translations.step4Configs.email.ccLabel.value"
          :placeholder="translations.step4Configs.email.ccPlaceholder.value"
        >
          <template v-slot:prepend>
            <q-icon name="content_copy" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.bcc"
          outlined
          dense
          :label="translations.step4Configs.email.bccLabel.value"
          :placeholder="translations.step4Configs.email.bccPlaceholder.value"
        >
          <template v-slot:prepend>
            <q-icon name="visibility_off" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.subject"
          outlined
          dense
          :label="translations.step4Configs.email.subjectLabel.value"
          :placeholder="translations.step4Configs.email.subjectPlaceholder.value"
          :rules="[(val: string) => !!val || translations.step4Configs.email.subjectRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="title" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.body"
          outlined
          dense
          type="textarea"
          :label="translations.step4Configs.email.bodyLabel.value"
          :placeholder="translations.step4Configs.email.bodyPlaceholder.value"
          :rules="[() => !!config.body || !!config.htmlBody || translations.step4Configs.email.bodyRequired.value]"
          rows="6"
        >
          <template v-slot:prepend>
            <q-icon name="description" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.htmlBody"
          outlined
          dense
          type="textarea"
          :label="translations.step4Configs.email.htmlBodyLabel.value"
          :placeholder="translations.step4Configs.email.htmlBodyPlaceholder.value"
          :hint="translations.step4Configs.email.htmlBodyHint.value"
          :rules="[() => !!config.body || !!config.htmlBody || translations.step4Configs.email.bodyRequired.value]"
          rows="6"
        >
          <template v-slot:prepend>
            <q-icon name="code" />
          </template>
        </q-input>
      </div>
    </div>
  </div>
</template>

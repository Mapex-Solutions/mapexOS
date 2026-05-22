<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type {
  StandardizedPayloadHelpProps,
  StandardizedPayloadHelpEmits,
} from './interfaces';

defineOptions({
  name: 'StandardizedPayloadHelpModal'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** EXTERNAL IMPORTS */
import { copyToClipboard } from 'quasar';

/** COMPONENTS */
import { BaseButton } from '@components/buttons';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert/notify';

/** COMPOSABLES */
import { useStandardizedPayloadHelpTranslations } from '@src/composables/i18n/components/dialogs/useStandardizedPayloadHelpTranslations';

/** PROPS & EMITS */
const props = defineProps<StandardizedPayloadHelpProps>();
const emit = defineEmits<StandardizedPayloadHelpEmits>();

/** COMPOSABLES & STORES */
const t = useStandardizedPayloadHelpTranslations();

/** COMPUTED */
const isOpen = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/** CONSTANTS */
const structureCode = `{
  eventType: string,    // Required
  eventId: string,      // Required
  data: object,         // Required
  metadata?: object,    // Optional
  created: string       // Required (ISO 8601)
}`;

const completeExampleCode = `const standardPayload = {
  eventType: "sensor.reading",
  eventId: \`\${payload.deviceId}-\${Date.now()}\`,
  data: {
    temperature: payload.temperature,
    humidity: payload.humidity || null,
    battery: payload.battery || null,
    deviceId: payload.deviceId
  },
  metadata: {
    source: "asset-template",
    version: "1.0.0",
    processingTime: new Date().toISOString()
  },
  created: new Date(payload.timestamp).toISOString()
};

return standardPayload;`;

/** FUNCTIONS */

/**
 * Copy code to clipboard with user notification
 * @param {string} code - Code string to copy
 */
const copyCode = async (code: string): Promise<void> => {
  try {
    await copyToClipboard(code);
    notifySuccess({
      message: t.notifications.copySuccess.value,
      timeout: 2000,
    });
  } catch {
    notifyFail({
      message: t.notifications.copyFail.value,
      timeout: 2000,
    });
  }
};

/**
 * Close modal dialog
 */
const closeModal = (): void => {
  isOpen.value = false;
};
</script>

<template>
  <q-dialog v-model="isOpen">
    <q-card class="standardized-payload-help-modal" style="width: 900px; max-width: 90vw;">
      <!-- Header -->
      <q-card-section class="modal-header">
        <div class="row items-center">
          <q-icon name="check_circle" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h5 text-weight-medium text-grey-9">{{ t.header.title.value }}</div>
          <q-space />
          <BaseButton flat round dense icon="close" color="grey-7" v-close-popup>
            <AppTooltip :content="t.footer.close.value" />
          </BaseButton>
        </div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pt-md">
        <q-scroll-area style="height: 70vh; max-height: 600px;">
          <div class="q-gutter-lg">
            <!-- Overview Section -->
            <q-card bordered flat class="section-card">
              <q-card-section>
                <div class="text-h6 text-weight-medium q-mb-md">
                  <q-icon name="info" color="primary" class="q-mr-xs" />
                  {{ t.overview.title.value }}
                </div>
                <div class="text-body1 q-mb-md">
                  {{ t.overview.description.value }}
                </div>
                <q-separator class="q-my-md" />
                <div class="text-subtitle2 text-weight-bold q-mb-sm">{{ t.overview.whyTitle.value }}</div>
                <ul class="text-body2">
                  <li v-for="(benefit, index) in t.overview.benefits.value" :key="index">{{ benefit }}</li>
                </ul>
              </q-card-section>
            </q-card>

            <!-- Structure Section -->
            <q-card bordered flat class="section-card">
              <q-card-section>
                <div class="row items-center q-mb-md">
                  <div class="text-h6 text-weight-medium">
                    <q-icon name="data_object" color="primary" class="q-mr-xs" />
                    {{ t.structure.title.value }}
                  </div>
                  <q-space />
                  <BaseButton
                    icon="content_copy"
                    flat
                    round
                    dense
                    color="primary"
                    @click="copyCode(structureCode)"
                  >
                    <AppTooltip :content="t.structure.copyTooltip.value" />
                  </BaseButton>
                </div>
                <q-card bordered flat class="bg-grey-10 text-white code-block q-mb-md">
                  <q-card-section class="q-pa-md">
                    <pre class="q-ma-none"><code>{
  eventType: string,    // {{ t.structure.comment1.value }}
  eventId: string,      // {{ t.structure.comment2.value }}
  data: object,         // {{ t.structure.comment3.value }}
  metadata?: object,    // {{ t.structure.comment4.value }}
  created: string       // {{ t.structure.comment5.value }}
}</code></pre>
                  </q-card-section>
                </q-card>
              </q-card-section>
            </q-card>

            <!-- Fields Description -->
            <q-card bordered flat class="section-card">
              <q-card-section>
                <div class="text-h6 text-weight-medium q-mb-md">
                  <q-icon name="list" color="primary" class="q-mr-xs" />
                  {{ t.fields.title.value }}
                </div>

                <!-- eventType -->
                <div class="field-description q-mb-lg">
                  <div class="row items-center q-mb-sm">
                    <DetailChip :label="t.fields.required.value" color="negative" size="sm" dense />
                    <span class="text-h6 text-weight-medium q-ml-sm">{{ t.fields.eventType.name.value }}</span>
                    <span class="text-caption text-grey-7 q-ml-sm">{{ t.fields.eventType.type.value }}</span>
                  </div>
                  <p class="text-body2">
                    {{ t.fields.eventType.description.value }}
                  </p>
                  <div class="text-caption text-weight-bold q-mb-xs">{{ t.fields.eventType.examplesTitle.value }}</div>
                  <ul class="text-body2">
                    <li v-for="(example, index) in t.fields.eventType.examples.value" :key="index">{{ example }}</li>
                  </ul>
                </div>

                <!-- eventId -->
                <div class="field-description q-mb-lg">
                  <div class="row items-center q-mb-sm">
                    <DetailChip :label="t.fields.required.value" color="negative" size="sm" dense />
                    <span class="text-h6 text-weight-medium q-ml-sm">{{ t.fields.eventId.name.value }}</span>
                    <span class="text-caption text-grey-7 q-ml-sm">{{ t.fields.eventId.type.value }}</span>
                  </div>
                  <p class="text-body2">
                    {{ t.fields.eventId.description.value }}
                  </p>
                  <div class="text-caption text-weight-bold q-mb-xs">{{ t.fields.eventId.examplesTitle.value }}</div>
                  <ul class="text-body2">
                    <li v-for="(example, index) in t.fields.eventId.examples.value" :key="index">{{ example }}</li>
                  </ul>
                </div>

                <!-- data -->
                <div class="field-description q-mb-lg">
                  <div class="row items-center q-mb-sm">
                    <DetailChip :label="t.fields.required.value" color="negative" size="sm" dense />
                    <span class="text-h6 text-weight-medium q-ml-sm">{{ t.fields.data.name.value }}</span>
                    <span class="text-caption text-grey-7 q-ml-sm">{{ t.fields.data.type.value }}</span>
                  </div>
                  <p class="text-body2">
                    {{ t.fields.data.description.value }}
                  </p>
                  <div class="text-caption text-weight-bold q-mb-xs">{{ t.fields.data.exampleTitle.value }}</div>
                  <q-card bordered flat class="bg-grey-10 text-white code-block">
                    <q-card-section class="q-pa-sm">
                      <pre class="q-ma-none"><code>{
  temperature: 23.5,
  humidity: 65.2,
  pressure: 1013.25,
  battery: 85,
  deviceId: "SENSOR001"
}</code></pre>
                    </q-card-section>
                  </q-card>
                </div>

                <!-- metadata -->
                <div class="field-description q-mb-lg">
                  <div class="row items-center q-mb-sm">
                    <DetailChip :label="t.fields.optional.value" color="grey" size="sm" dense />
                    <span class="text-h6 text-weight-medium q-ml-sm">{{ t.fields.metadata.name.value }}</span>
                    <span class="text-caption text-grey-7 q-ml-sm">{{ t.fields.metadata.type.value }}</span>
                  </div>
                  <p class="text-body2">
                    {{ t.fields.metadata.description.value }}
                  </p>
                  <div class="text-caption text-weight-bold q-mb-xs">{{ t.fields.metadata.exampleTitle.value }}</div>
                  <q-card bordered flat class="bg-grey-10 text-white code-block">
                    <q-card-section class="q-pa-sm">
                      <pre class="q-ma-none"><code>{
  source: "asset-template",
  version: "1.0.0",
  processingTime: "2024-01-15T10:30:01.234Z",
  gateway: "GW-001",
  rssi: -65
}</code></pre>
                    </q-card-section>
                  </q-card>
                </div>

                <!-- created -->
                <div class="field-description">
                  <div class="row items-center q-mb-sm">
                    <DetailChip :label="t.fields.required.value" color="negative" size="sm" dense />
                    <span class="text-h6 text-weight-medium q-ml-sm">{{ t.fields.created.name.value }}</span>
                    <span class="text-caption text-grey-7 q-ml-sm">{{ t.fields.created.type.value }}</span>
                  </div>
                  <p class="text-body2">
                    {{ t.fields.created.description.value }}
                  </p>
                  <div class="text-caption text-weight-bold q-mb-xs">{{ t.fields.created.validFormatsTitle.value }}</div>
                  <ul class="text-body2">
                    <li v-for="(format, index) in t.fields.created.formats.value" :key="index">{{ format }}</li>
                  </ul>
                </div>
              </q-card-section>
            </q-card>

            <!-- Complete Example -->
            <q-card bordered flat class="section-card">
              <q-card-section>
                <div class="row items-center q-mb-md">
                  <div class="text-h6 text-weight-medium">
                    <q-icon name="code" color="primary" class="q-mr-xs" />
                    {{ t.example.title.value }}
                  </div>
                  <q-space />
                  <BaseButton
                    icon="content_copy"
                    flat
                    round
                    dense
                    color="primary"
                    @click="copyCode(completeExampleCode)"
                  >
                    <AppTooltip :content="t.example.copyTooltip.value" />
                  </BaseButton>
                </div>
                <div class="text-body2 text-grey-7 q-mb-md">
                  {{ t.example.description.value }}
                </div>
                <q-card bordered flat class="bg-grey-10 text-white code-block">
                  <q-card-section class="q-pa-md">
                    <pre class="q-ma-none"><code>// Input payload
const inputPayload = {
  deviceId: "SENSOR001",
  timestamp: "2024-01-15T10:30:00Z",
  temp_c: 23.5,
  humidity: 65.2,
  battery: 85
};

// Conversion Script
const standardPayload = {
  eventType: "sensor.reading",
  eventId: `${payload.deviceId}-${Date.now()}`,
  data: {
    temperature: payload.temp_c,
    humidity: payload.humidity || null,
    battery: payload.battery || null,
    deviceId: payload.deviceId
  },
  metadata: {
    source: "asset-template",
    version: "1.0.0",
    processingTime: new Date().toISOString()
  },
  created: new Date(payload.timestamp).toISOString()
};

return standardPayload;

// Output (StandardizedPayload)
{
  "eventType": "sensor.reading",
  "eventId": "SENSOR001-1705315800000",
  "data": {
    "temperature": 23.5,
    "humidity": 65.2,
    "battery": 85,
    "deviceId": "SENSOR001"
  },
  "metadata": {
    "source": "asset-template",
    "version": "1.0.0",
    "processingTime": "2024-01-15T10:30:01.234Z"
  },
  "created": "2024-01-15T10:30:00.000Z"
}</code></pre>
                  </q-card-section>
                </q-card>
              </q-card-section>
            </q-card>

            <!-- Common Mistakes -->
            <q-card bordered flat class="section-card">
              <q-card-section>
                <div class="text-h6 text-weight-medium q-mb-md">
                  <q-icon name="warning" color="warning" class="q-mr-xs" />
                  {{ t.mistakes.title.value }}
                </div>

                <q-list separator>
                  <q-item v-for="(mistake, index) in t.mistakes.items.value" :key="index">
                    <q-item-section avatar>
                      <q-icon name="close" color="negative" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ mistake.label }}</q-item-label>
                      <q-item-label caption>
                        {{ mistake.caption }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card-section>
            </q-card>
          </div>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator />
      <q-card-section class="row items-center q-py-sm bg-grey-2">
        <q-icon name="info" size="20px" color="primary" class="q-mr-sm" />
        <div class="text-caption text-grey-7">
          {{ t.footer.info.value }}
        </div>
        <q-space />
        <BaseButton color="primary" :label="t.footer.close.value" @click="closeModal" />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style scoped lang="scss">
.modal-header {
  padding: 20px 24px !important;
  background: var(--mapex-surface-bg);
  border-bottom: 1px solid var(--mapex-card-border);
}

.standardized-payload-help-modal {
  .section-card {
    border-radius: var(--mapex-radius-md);
    transition: box-shadow 0.3s ease;

    &:hover {
      box-shadow: var(--mapex-shadow-md);
    }
  }

  .code-block {
    border-radius: var(--mapex-radius-sm);
    overflow: hidden;

    pre {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 13px;
      line-height: 1.6;
      overflow-x: auto;
      white-space: pre-wrap;
      word-wrap: break-word;
    }

    code {
      font-family: inherit;
    }
  }

  .field-description {
    padding: 12px;
    background-color: var(--mapex-surface-elevated);
    border-radius: var(--mapex-radius-md);
    border-left: 4px solid var(--q-primary);

    code {
      background-color: var(--mapex-surface-highlight);
      padding: 2px 6px;
      border-radius: var(--mapex-radius-xs);
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 12px;
    }

    ul {
      margin: 4px 0;
      padding-left: 20px;
    }

    li {
      margin-bottom: 4px;
    }
  }
}
</style>

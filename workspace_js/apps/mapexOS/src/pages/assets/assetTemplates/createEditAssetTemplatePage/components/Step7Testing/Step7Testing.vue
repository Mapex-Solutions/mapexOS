<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { Step7TestingProps, Step7TestingEmits } from './interfaces/Step7Testing.interface';
import type { AssetTemplateData, TestResults } from '../../interfaces';

defineOptions({
  name: 'Step7Testing'
});

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { AvailableFieldsList } from '@components/lists/availableFieldsList';

/** COMPOSABLES */
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

/** LOCAL IMPORTS (constants and handlers ONLY - NO types here!) */
import { executeScriptTests, formatJSON } from '../../handlers';
import { normalizePayloadPaths } from '../../utils';

const props = defineProps<Step7TestingProps>();
const emit = defineEmits<Step7TestingEmits>();

/** COMPOSABLES & STORES */
const t = useAddAssetTemplateTranslations();

/** STATE */
const testing = ref(false);
const showPayloadJson = ref(false);

/** COMPUTED */

/**
 * Computes available fields from the test results StandardizedPayload
 * Only returns fields when test was successful and payload exists
 * @returns {string[]} Array of normalized field paths
 */
const availableFields = computed(() => {
  if (!props.testResults.success || !props.testResults.standardizedPayload) {
    return [];
  }

  return normalizePayloadPaths(props.testResults.standardizedPayload);
});

/** WATCHERS */

/**
 * Watch availableFields and update formData when changed
 * Updates the fields array in the form data
 */
watch(availableFields, (fields) => {
  const updatedData: AssetTemplateData = {
    ...props.modelValue,
    availableFields: fields,
  };
  emit('update:modelValue', updatedData);
});

/** FUNCTIONS */

// Run test
async function runTest() {
  testing.value = true;

  const results = await executeScriptTests(
    props.modelValue.scriptProcessor || '',
    props.modelValue.scriptValidator || '',
    props.modelValue.scriptConversion || '',
    props.modelValue.scriptTest || '{}'
  );

  emit('update:testResults', results);
  testing.value = false;
}

// Clear test results
function clearTest() {
  const emptyResults: TestResults = {
    executed: false,
    success: false,
    steps: [],
    output: null,
    logs: [],
  };
  emit('update:testResults', emptyResults);
  showPayloadJson.value = false;
}

// Open help modal
function openHelpModal() {
  emit('showStandardizedPayloadHelp');
}
</script>

<template>
  <div>
    <!-- Info Banner -->
    <q-banner dense rounded class="bg-blue-1 text-primary q-mb-md">
      <template v-slot:avatar>
        <q-icon name="info" color="primary" />
      </template>
      <div class="text-body2">
        {{ t.steps.step7.banner.standardizedPayloadInfo.value }}
        <q-btn
          flat
          dense
          no-caps
          size="sm"
          color="primary"
          icon-right="open_in_new"
          :label="t.steps.step7.banner.viewFormat.value"
          @click="openHelpModal"
        />
      </div>
    </q-banner>

    <!-- Test Execution Card -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section>
        <div class="row items-center">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-grey-8">
              {{ t.steps.step7.testSection.usingPayloadFromInline.value }} • <a href="#" @click.prevent="$router.push({ hash: '#step-6' })" class="text-primary">{{ t.steps.step7.testSection.editLink.value }}</a>
            </div>
          </div>
          <div class="col-auto">
            <q-btn
              unelevated
              color="primary"
              icon="play_arrow"
              label="Run Test"
              :loading="testing"
              @click="runTest"
            />
            <q-btn
              v-if="testResults.executed"
              flat
              color="grey-7"
              icon="refresh"
              size="sm"
              class="q-ml-sm"
              @click="clearTest"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Test Results -->
    <div v-if="testResults.executed">
      <!-- Status Card -->
      <q-card flat bordered class="q-mb-md" :class="testResults.success ? 'bg-green-1' : 'bg-red-1'">
        <q-card-section>
          <div class="row items-center q-gutter-sm">
            <div class="col-auto">
              <q-icon
                :name="testResults.success ? 'check_circle' : 'error'"
                :color="testResults.success ? 'positive' : 'negative'"
                size="32px"
              />
            </div>
            <div class="col">
              <div class="text-subtitle1 text-weight-medium" :class="testResults.success ? 'text-positive' : 'text-negative'">
                {{ testResults.success ? 'All Tests Passed' : 'Test Failed' }}
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Success Results -->
      <template v-if="testResults.success">
        <!-- Available Fields - DESTACADO -->
        <q-card v-if="availableFields.length > 0" flat bordered class="q-mb-md highlight-card">
          <q-card-section class="bg-blue-1">
            <div class="text-subtitle1 text-weight-medium text-primary">
              <q-icon name="list_alt" color="primary" class="q-mr-xs" />
              Available Fields ({{ availableFields.length }})
            </div>
          </q-card-section>
          <q-separator />
          <q-card-section>
            <AvailableFieldsList
              :fields="availableFields"
              :max-height="400"
            />
          </q-card-section>
        </q-card>

        <!-- Standardized Payload - Collapsible -->
        <q-card v-if="testResults.standardizedPayload" flat bordered class="q-mb-md">
          <q-card-section class="cursor-pointer" @click="showPayloadJson = !showPayloadJson">
            <div class="row items-center">
              <div class="col">
                <div class="text-subtitle2 text-weight-medium">
                  <q-icon name="code" size="sm" class="q-mr-xs" />
                  Generated StandardizedPayload
                </div>
              </div>
              <div class="col-auto">
                <q-icon :name="showPayloadJson ? 'expand_less' : 'expand_more'" size="sm" />
              </div>
            </div>
          </q-card-section>

          <q-slide-transition>
            <div v-show="showPayloadJson">
              <q-separator />
              <q-card-section class="bg-grey-1">
                <pre class="payload-json">{{ formatJSON(testResults.standardizedPayload) }}</pre>
              </q-card-section>
            </div>
          </q-slide-transition>
        </q-card>
      </template>

      <!-- Error Results -->
      <template v-if="!testResults.success">
        <q-card flat bordered class="q-mb-md error-card">
          <q-card-section class="bg-red-1">
            <div class="text-subtitle1 text-weight-medium text-negative q-mb-sm">
              <q-icon name="error" size="sm" class="q-mr-xs" />
              Error Details
            </div>
            <div v-for="(step, index) in testResults.steps" :key="index">
              <div v-if="!step.success" class="q-mb-sm">
                <div class="text-subtitle2 text-weight-medium text-negative q-mb-xs">
                  {{ step.name }}
                </div>
                <q-card flat bordered class="bg-white">
                  <q-card-section>
                    <div class="text-body2 text-negative">
                      {{ step.error }}
                    </div>
                    <div v-if="step.details" class="q-mt-sm">
                      <div class="text-caption text-weight-medium text-grey-8 q-mb-xs">{{ t.steps.step7.testSection.detailsLabel.value }}</div>
                      <pre class="error-details">{{ formatJSON(step.details) }}</pre>
                    </div>
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </template>
    </div>

    <!-- Empty State -->
    <q-card v-else flat bordered class="empty-state-card">
      <q-card-section class="text-center q-py-lg">
        <q-icon name="science" size="48px" color="grey-5" class="q-mb-sm" />
        <div class="text-subtitle1 text-grey-7 q-mb-xs">No test results yet</div>
        <div class="text-body2 text-grey-6">
          Click "Run Test" to validate your scripts
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<style scoped>
.highlight-card {
  border: 2px solid var(--q-primary);
  border-radius: var(--mapex-radius-md);
}

.error-card {
  border: 2px solid var(--q-negative);
  border-radius: var(--mapex-radius-md);
}

.empty-state-card {
  border-radius: var(--mapex-radius-md);
  min-height: 300px;
}

.payload-json {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  padding: 16px;
  background: var(--mapex-surface-bg);
  border-radius: var(--mapex-radius-xs);
  max-height: 500px;
  overflow-y: auto;
}

.error-details {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  padding: 12px;
  background: var(--mapex-surface-bg);
  border-radius: var(--mapex-radius-xs);
  max-height: 300px;
  overflow-y: auto;
}
</style>

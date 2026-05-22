<script setup lang="ts">
defineOptions({
  name: 'ExecuteResultDialog'
});

/** VUE IMPORTS */
import { computed } from 'vue';
import { useRouter } from 'vue-router';

/** COMPOSABLES */
import { useWorkflowInstanceListPageTranslations } from '@composables/i18n/pages/automations/workflowInstances/workflowInstanceListPage';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: boolean;
  loading: boolean;
  result: { workflowUUID: string; status: string; errorInfo?: any } | null;
  error: string | null;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
}>();

/** COMPOSABLES & STORES */
const router = useRouter();
const t = useWorkflowInstanceListPageTranslations();

/** COMPUTED */
const dialogOpen = computed({
  get: () => props.modelValue,
  set: (val: boolean) => emit('update:modelValue', val),
});

const statusColor = computed(() => {
  if (!props.result) return 'grey-6';
  switch (props.result.status) {
    case 'completed': return 'green-6';
    case 'waiting': return 'orange-6';
    default: return 'red-6';
  }
});

/** FUNCTIONS */
function navigateToExecutions(): void {
  void router.push('/logs/workflow_executions');
}
</script>

<template>
  <q-dialog v-model="dialogOpen">
    <q-card style="min-width: 400px;">
      <q-card-section class="row items-center">
        <q-icon name="play_arrow" color="green-7" size="sm" class="q-mr-sm" />
        <div class="text-h6">{{ t.executeResult.title.value }}</div>
        <q-space />
        <q-btn flat round dense icon="close" v-close-popup />
      </q-card-section>

      <q-separator />

      <!-- Loading -->
      <q-card-section v-if="loading" class="column items-center q-pa-xl">
        <q-spinner color="primary" size="40px" />
        <div class="text-body2 text-grey-7 q-mt-md">{{ t.executeResult.executing.value }}</div>
      </q-card-section>

      <!-- Error -->
      <q-card-section v-else-if="error">
        <q-banner class="bg-red-1 text-red-8" rounded>
          <template #avatar>
            <q-icon name="error" color="red" />
          </template>
          {{ error }}
        </q-banner>
      </q-card-section>

      <!-- Result -->
      <q-card-section v-else-if="result">
        <div class="row items-center q-mb-md">
          <div class="text-body2 text-grey-7 q-mr-sm">{{ t.executeResult.statusLabel.value }}:</div>
          <q-badge
            :color="statusColor"
            :label="result.status.toUpperCase()"
          />
        </div>

        <div class="row items-center q-mb-md">
          <div class="text-body2 text-grey-7 q-mr-sm">{{ t.executeResult.uuidLabel.value }}:</div>
          <code class="text-body2">{{ result.workflowUUID }}</code>
        </div>

        <q-banner v-if="result.errorInfo" class="bg-red-1 text-red-8 q-mb-md" rounded>
          <template #avatar>
            <q-icon name="error" color="red" />
          </template>
          <div class="text-weight-medium">{{ result.errorInfo.code }}</div>
          <div>{{ result.errorInfo.message }}</div>
        </q-banner>

        <q-btn
          flat
          color="primary"
          icon="open_in_new"
          :label="t.executeResult.viewExecutionsButton.value"
          no-caps
          @click="navigateToExecutions"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

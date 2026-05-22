<script setup lang="ts">
/** TYPE IMPORTS */
import type { CredentialSelectorProps, CredentialSelectorEmits } from './interfaces/CredentialSelector.interface';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';

/** COMPONENTS */
import { CredentialManagerDialog } from '../CredentialManagerDialog';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<CredentialSelectorProps>();
const emit = defineEmits<CredentialSelectorEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Available credentials loaded from API
 */
const credentials = ref<Array<{ id: string; name: string }>>([]);

/**
 * Loading state for credential fetch
 */
const loading = ref(false);

/**
 * Whether the credential manager dialog is open
 */
const showManager = ref(false);

/** COMPUTED */

/**
 * Options for the q-select dropdown
 */
const credentialOptions = computed(() =>
  credentials.value.map((c) => ({
    label: c.name,
    value: c.id,
  })),
);

/** WATCHERS */

watch(showManager, (visible) => {
  if (!visible) {
    // Re-fetch credentials after manager closes (may have created/deleted)
    void fetchCredentials();
  }
});

/** FUNCTIONS */

/**
 * Fetch credentials for this plugin from the API
 *
 * @returns {Promise<void>}
 */
async function fetchCredentials(): Promise<void> {
  loading.value = true;
  try {
    const response = await apis.vault.credential.list({ pluginId: props.pluginId });
    const items = response?.items ?? [];
    credentials.value = (items as Array<{ id: string; name: string }>).map((c) => ({
      id: c.id,
      name: c.name,
    }));
  } catch (error) {
    console.error('[CredentialSelector] Failed to fetch credentials:', error);
    credentials.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * Handle credential selection change
 *
 * @param {string | null} value - Selected credential ID
 */
function handleSelect(value: string | null): void {
  emit('update:modelValue', value);
}

/** LIFECYCLE HOOKS */

onMounted(() => {
  void fetchCredentials();
});
</script>

<template>
  <div class="credential-selector q-mb-md">
    <div class="text-caption text-weight-medium q-mb-xs" style="color: var(--mapex-text-secondary)">
      <q-icon name="vpn_key" size="14px" class="q-mr-xs" />
      {{ t.credentials.credentialSelector.value }}
    </div>
    <div class="row no-wrap q-gutter-sm items-center">
      <q-select
        :model-value="modelValue"
        :options="credentialOptions"
        :loading="loading"
        :placeholder="credentials.length > 0 ? t.credentials.selectCredential.value : t.credentials.noCredentialsAvailable.value"
        outlined
        dense
        clearable
        emit-value
        map-options
        class="col"
        @update:model-value="handleSelect"
      />
      <q-btn
        flat
        dense
        round
        icon="add"
        size="sm"
        color="primary"
        @click="showManager = true"
      >
        <q-tooltip>{{ t.credentials.newCredential.value }}</q-tooltip>
      </q-btn>
    </div>

    <!-- Inline Credential Manager Dialog -->
    <CredentialManagerDialog
      v-model="showManager"
      :plugin-id="pluginId"
      :plugin-name="pluginName"
      :credential-defs="credentialDefs"
    />
  </div>
</template>

<style lang="scss" scoped>
.credential-selector {
  padding: 8px 0;
  border-bottom: 1px solid var(--mapex-card-border);
}
</style>

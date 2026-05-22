<script setup lang="ts">
/** TYPE IMPORTS */
import type { CredentialManagerDialogProps, CredentialManagerDialogEmits, CredentialListItem } from './interfaces/CredentialManagerDialog.interface';

/** VUE IMPORTS */
import type { PluginCredentialDefinition } from '@components/workflow/interfaces';
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { CredentialForm } from '../CredentialForm';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { CREDENTIAL_DIALOG_MAX_WIDTH } from './constants';

/** PROPS & EMITS */
const props = defineProps<CredentialManagerDialogProps>();
const emit = defineEmits<CredentialManagerDialogEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Dialog view mode: list or form
 */
const mode = ref<'list' | 'form'>('list');

/**
 * Credential list loaded from API
 */
const credentials = ref<CredentialListItem[]>([]);

/**
 * Loading state for list fetch
 */
const loading = ref(false);

/**
 * MongoDB ID of the credential being edited (null = creating new)
 */
const editingId = ref<string | null>(null);

/**
 * Name of the credential being edited (for form pre-fill)
 */
const editingName = ref('');

/**
 * ID of the credential currently being tested
 */
const testingId = ref<string | null>(null);

/**
 * ID of the credential currently being deleted
 */
const deletingId = ref<string | null>(null);

/**
 * Selected credential type ID (for create mode when multiple types exist)
 */
const selectedTypeId = ref<string>('');

/** COMPUTED */

/**
 * Whether the plugin has multiple credential types
 */
const hasMultipleTypes = computed(() => props.credentialDefs.length > 1);

/**
 * Options for the credential type selector
 */
const typeOptions = computed(() =>
  props.credentialDefs.map((def) => ({
    label: def.name,
    value: def.id,
  })),
);

/**
 * Currently selected credential definition
 */
const selectedCredentialDef = computed((): PluginCredentialDefinition | null => {
  if (!selectedTypeId.value) return props.credentialDefs[0] ?? null;
  return props.credentialDefs.find((d) => d.id === selectedTypeId.value) ?? null;
});

/** WATCHERS */

watch(() => props.modelValue, (visible) => {
  if (visible) {
    mode.value = 'list';
    editingId.value = null;
    editingName.value = '';
    selectedTypeId.value = props.credentialDefs[0]?.id ?? '';
    void fetchCredentials();
  }
}, { immediate: true });

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
    credentials.value = (items as Array<{ id: string; name: string; created?: string }>).map((c) => ({
      id: c.id,
      name: c.name,
      created: c.created ?? '',
      testStatus: null,
    }));
  } catch (error) {
    console.error('[CredentialManager] Failed to fetch credentials:', error);
    credentials.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * Open the form for creating a new credential
 */
function openNewForm(): void {
  editingId.value = null;
  editingName.value = '';
  mode.value = 'form';
}

/**
 * Open the form for editing an existing credential
 *
 * @param {CredentialListItem} cred - Credential to edit
 */
function openEditForm(cred: CredentialListItem): void {
  editingId.value = cred.id;
  editingName.value = cred.name;
  mode.value = 'form';
}

/**
 * Handle save from the credential form (create or update)
 *
 * @param {{ name: string; data: Record<string, unknown> }} payload - Form data
 * @returns {Promise<void>}
 */
async function handleSave(payload: { name: string; data: Record<string, unknown> }): Promise<void> {
  try {
    if (editingId.value) {
      // Update existing (2 args: pathParams, bodyParams)
      await apis.vault.credential.update(
        { credentialId: editingId.value },
        { name: payload.name, data: payload.data },
      );
    } else {
      // Create new
      await apis.vault.credential.create({
        pluginId: props.pluginId,
        type: 'manual',
        credentialDefId: selectedTypeId.value,
        name: payload.name,
        data: payload.data,
      });
    }
    notifySuccess({ message: t.credentials.saveSuccess.value });
    mode.value = 'list';
    void fetchCredentials();
  } catch (error) {
    console.error('[CredentialManager] Save failed:', error);
    notifyFail({ message: t.credentials.saveFailed.value });
  }
}

/**
 * Test a credential by calling the test endpoint
 *
 * @param {CredentialListItem} cred - Credential to test
 * @returns {Promise<void>}
 */
async function handleTest(cred: CredentialListItem): Promise<void> {
  testingId.value = cred.id;
  try {
    const result = await apis.vault.credential.test({ credentialId: cred.id });
    const item = credentials.value.find((c) => c.id === cred.id);
    if (item) item.testStatus = result.success;

    if (result.success) {
      notifySuccess({ message: t.credentials.testSuccess.value });
    } else {
      notifyFail({ message: t.credentials.testFailedMsg('') });
    }
  } catch (error) {
    const item = credentials.value.find((c) => c.id === cred.id);
    if (item) item.testStatus = false;
    notifyFail({ message: t.credentials.testFailedMsg(String(error)) });
  } finally {
    testingId.value = null;
  }
}

/**
 * Delete a credential after confirmation
 *
 * @param {CredentialListItem} cred - Credential to delete
 * @returns {Promise<void>}
 */
async function handleDelete(cred: CredentialListItem): Promise<void> {
  deletingId.value = cred.id;
  try {
    await apis.vault.credential.delete({ credentialId: cred.id });
    notifySuccess({ message: t.credentials.deleteSuccess.value });
    credentials.value = credentials.value.filter((c) => c.id !== cred.id);
  } catch (error) {
    console.error('[CredentialManager] Delete failed:', error);
    notifyFail({ message: t.credentials.saveFailed.value });
  } finally {
    deletingId.value = null;
  }
}

/**
 * Close the dialog
 */
function closeDialog(): void {
  emit('update:modelValue', false);
}
</script>

<template>
  <q-dialog
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <q-card class="credential-manager-dialog" :style="{ maxWidth: CREDENTIAL_DIALOG_MAX_WIDTH, width: '100%' }">
      <!-- Header -->
      <q-card-section class="row items-center no-wrap q-pb-sm">
        <q-icon name="vpn_key" size="24px" color="primary" class="q-mr-sm" />
        <div class="text-h6">
          {{ mode === 'form'
            ? (editingId ? t.credentials.editCredential.value : t.credentials.newCredential.value)
            : t.credentials.titleMsg(pluginName)
          }}
        </div>
        <q-space />
        <q-btn flat dense round icon="close" @click="closeDialog" />
      </q-card-section>

      <q-separator />

      <!-- List Mode -->
      <q-card-section v-if="mode === 'list'">
        <!-- Loading -->
        <div v-if="loading" class="flex flex-center q-pa-lg">
          <q-spinner color="primary" size="32px" />
        </div>

        <!-- Empty state -->
        <div v-else-if="credentials.length === 0" class="text-center q-pa-lg">
          <q-icon name="vpn_key_off" size="48px" color="grey-5" />
          <div class="text-subtitle1 text-grey-6 q-mt-sm">
            {{ t.credentials.noCredentials.value }}
          </div>
          <div class="text-caption text-grey-5">
            {{ t.credentials.noCredentialsDescMsg(pluginName) }}
          </div>
        </div>

        <!-- Credential list -->
        <q-list v-else separator>
          <q-item v-for="cred in credentials" :key="cred.id">
            <q-item-section avatar>
              <q-icon
                name="vpn_key"
                :color="cred.testStatus === true ? 'positive' : cred.testStatus === false ? 'negative' : 'grey-6'"
              />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium">{{ cred.name }}</q-item-label>
              <q-item-label v-if="cred.created" caption>
                {{ new Date(cred.created).toLocaleDateString() }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row q-gutter-xs">
                <!-- Test -->
                <q-btn
                  flat
                  dense
                  no-caps
                  size="sm"
                  color="primary"
                  icon="play_arrow"
                  :label="t.credentials.test.value"
                  :loading="testingId === cred.id"
                  @click="handleTest(cred)"
                />
                <!-- Edit -->
                <q-btn
                  flat
                  dense
                  round
                  size="sm"
                  color="grey-7"
                  icon="edit"
                  @click="openEditForm(cred)"
                />
                <!-- Delete -->
                <q-btn
                  flat
                  dense
                  round
                  size="sm"
                  color="negative"
                  icon="delete_outline"
                  :loading="deletingId === cred.id"
                  @click="handleDelete(cred)"
                />
              </div>
            </q-item-section>
          </q-item>
        </q-list>

        <!-- New credential button -->
        <div class="q-mt-md">
          <q-btn
            unelevated
            no-caps
            color="primary"
            icon="add"
            :label="t.credentials.newCredential.value"
            @click="openNewForm"
          />
        </div>
      </q-card-section>

      <!-- Form Mode -->
      <q-card-section v-else>
        <!-- Credential Type Selector (only for create with multiple types) -->
        <q-select
          v-if="hasMultipleTypes && !editingId"
          v-model="selectedTypeId"
          :options="typeOptions"
          :label="t.credentials.credentialType.value"
          outlined
          dense
          emit-value
          map-options
          class="q-mb-md"
        />

        <CredentialForm
          v-if="selectedCredentialDef"
          :credential-def="selectedCredentialDef"
          :initial-name="editingName"
          :is-edit="editingId !== null"
          @save="handleSave"
          @cancel="mode = 'list'"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style lang="scss" scoped>
.credential-manager-dialog {
  background: var(--mapex-card-bg);
  border-radius: var(--mapex-radius-lg);
}
</style>

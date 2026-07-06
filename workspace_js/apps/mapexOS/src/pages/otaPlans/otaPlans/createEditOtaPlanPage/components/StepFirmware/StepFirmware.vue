<template>
  <div class="q-gutter-md">
    <!-- Version -->
    <q-input
      :model-value="modelValue.version"
      outlined
      dense
      class="rounded-borders"
      :label="t.firmware.version.value + ' *'"
      :hint="t.firmware.versionHint.value"
      :disable="modelValue.firmwareReady"
      :rules="[(v: string) => !!v || t.validation.required.value]"
      @update:model-value="(v) => update({ version: String(v ?? '') })"
    >
      <template #prepend>
        <q-icon name="new_releases" color="primary" />
      </template>
    </q-input>

    <!-- Drag & drop zone -->
    <div
      class="dropzone"
      :class="{ 'dropzone--active': dragging, 'dropzone--disabled': modelValue.firmwareReady }"
      @click="openPicker"
      @dragover.prevent="dragging = true"
      @dragleave.prevent="dragging = false"
      @drop.prevent="onDrop"
    >
      <q-icon name="cloud_upload" size="42px" color="primary" class="q-mb-sm" />
      <div class="text-body2 text-weight-medium">{{ t.firmware.dropzone.title.value }}</div>
      <div class="text-caption text-grey-7">{{ t.firmware.dropzone.hint.value }}</div>
      <input
        ref="fileInput"
        type="file"
        accept=".bin"
        class="hidden-input"
        @change="onPicked"
      />
    </div>

    <!-- File facts -->
    <div v-if="modelValue.file" class="row q-gutter-sm">
      <DetailChip icon="attach_file" color="grey" :label="`${t.firmware.fileFacts.name.value}: ${modelValue.file.name}`" />
      <DetailChip icon="straighten" color="blue" :label="`${t.firmware.fileFacts.size.value}: ${formatSize(modelValue.size)}`" />
      <DetailChip
        v-if="modelValue.checksumHex"
        icon="tag"
        color="teal"
        :label="`${t.firmware.fileFacts.checksum.value}: ${modelValue.checksumHex.slice(0, 16)}…`"
      />
    </div>

    <!-- Upload action / state -->
    <div class="row items-center q-gutter-sm">
      <BaseButton
        v-if="!modelValue.firmwareReady"
        color="primary"
        icon="cloud_upload"
        :label="uploadLabel"
        :loading="uploadPhase !== 'idle'"
        :disable="!canUpload"
        @click="uploadFirmware"
      />
      <DetailChip
        v-else
        icon="check_circle"
        color="positive"
        :label="t.firmware.ready.value"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'StepFirmware'
});

/** TYPE IMPORTS */
import type { OtaPlanFormData } from '../../interfaces/createEditOtaPlan.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { BaseButton } from '@components/buttons';
import { DetailChip } from '@components/chips/DetailChip';

/** COMPOSABLES */
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: OtaPlanFormData;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: OtaPlanFormData];
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();
const logger = useLogger('StepFirmware');

/** STATE */
const uploadPhase = ref<'idle' | 'computing' | 'uploading' | 'confirming'>('idle');
const dragging = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

/** COMPUTED */

/** Everything the artifact declaration needs must be present before uploading */
const canUpload = computed(() =>
  !!(props.modelValue.targetTemplateId && props.modelValue.version && props.modelValue.file),
);

/** Upload button label follows the phase */
const uploadLabel = computed(() => {
  if (uploadPhase.value === 'computing') return t.firmware.computing.value;
  if (uploadPhase.value === 'uploading') return t.firmware.uploading.value;
  if (uploadPhase.value === 'confirming') return t.firmware.confirming.value;
  return t.firmware.upload.value;
});

/** FUNCTIONS */

/**
 * Merge a partial change into the form state
 */
function update(partial: Partial<OtaPlanFormData>): void {
  emit('update:modelValue', { ...props.modelValue, ...partial });
}

/** Open the native picker (the dropzone is also clickable) */
function openPicker(): void {
  if (props.modelValue.firmwareReady) return;
  fileInput.value?.click();
}

/** Handle a picked file from the hidden input */
function onPicked(event: Event): void {
  const input = event.target as HTMLInputElement;
  setFile(input.files?.[0] ?? null);
  input.value = '';
}

/** Handle a dropped file */
function onDrop(event: DragEvent): void {
  dragging.value = false;
  if (props.modelValue.firmwareReady) return;
  setFile(event.dataTransfer?.files?.[0] ?? null);
}

/**
 * A new file resets the upload lifecycle (a different binary needs a fresh
 * artifact — objects are immutable platform-side)
 */
function setFile(file: File | null): void {
  update({
    file,
    size: file?.size ?? 0,
    checksumHex: '',
    firmwareId: null,
    firmwareReady: false,
  });
}

/** Human-readable byte size */
function formatSize(bytes: number): string {
  if (!bytes) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

/**
 * Full upload lifecycle: checksum -> declare (init) -> presigned PUT straight
 * to object storage -> finalize. The step is complete only after finalize
 * confirms the stored object. The target template was chosen in the previous
 * step; this one owns only the binary.
 */
async function uploadFirmware(): Promise<void> {
  const f = props.modelValue;
  if (!f.file || !f.targetTemplateId || !f.version) return;

  try {
    uploadPhase.value = 'computing';
    // hex is shown to the operator (human-readable); base64 is the store's
    // native x-amz-checksum-sha256 encoding used for the presign + upload.
    const { hex, base64 } = await apis.assets.otaPlans.sha256OfFirmware(f.file);
    update({ checksumHex: hex });

    uploadPhase.value = 'uploading';
    // The store checksum is the S3 base64 of the digest — it must be the SAME
    // value the backend signs into the presigned PUT and the browser sends in
    // the x-amz-checksum-sha256 header, or MinIO rejects the upload (403).
    const { firmwareId, uploadUrl } = await apis.assets.otaPlans.firmwareInit({
      targetTemplateId: f.targetTemplateId,
      version: f.version,
      filename: f.file.name,
      size: f.file.size,
      sha256: base64,
    });
    await apis.assets.otaPlans.uploadFirmwareBinary(uploadUrl, f.file, base64);

    uploadPhase.value = 'confirming';
    await apis.assets.otaPlans.firmwareComplete({ firmwareId });

    update({ firmwareId, checksumHex: hex, firmwareReady: true });
    notifySuccess({ message: t.firmware.uploadSuccess.value });
  } catch (err: unknown) {
    logger.error('Firmware upload failed:', err);
    notifyFail({ message: t.firmware.uploadError.value });
    update({ firmwareId: null, firmwareReady: false });
  } finally {
    uploadPhase.value = 'idle';
  }
}
</script>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.dropzone {
  border: 2px dashed var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  padding: var(--mapex-spacing-xl) var(--mapex-spacing-md);
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s ease, background-color 0.2s ease;

  &:hover,
  &--active {
    border-color: var(--mapex-active-border);
    background-color: var(--mapex-active-bg);
  }

  &--disabled {
    opacity: 0.6;
    cursor: default;
  }
}

.hidden-input {
  display: none;
}
</style>

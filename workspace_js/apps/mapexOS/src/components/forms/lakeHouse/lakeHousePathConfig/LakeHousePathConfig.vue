<script setup lang="ts">
defineOptions({
  name: 'LakeHousePathConfig'
});

import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';

import { computed } from 'vue';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { VueDraggable } from 'vue-draggable-plus';
import { notifySuccess, notifyFail } from '@utils/alert/notify';

import { DEFAULT_PATH_CONFIG } from './constants';

// this ref is automatically tied to `modelValue` + `update:modelValue`
const modelRef = defineModel<LakeHouseConfigProps>({
  default: () => ({ pathConfig: DEFAULT_PATH_CONFIG }),
});


// Available partition options
const availablePartitions = [
  { value: 'year', label: 'Year (YYYY)', example: 'year=2025' },
  { value: 'month', label: 'Month (MM)', example: 'month=12' },
  { value: 'day', label: 'Day (DD)', example: 'day=15' },
  { value: 'hour', label: 'Hour (HH)', example: 'hour=14' },
  { value: 'asset_id', label: 'Asset ID', example: 'asset_id=sensor-001' },
  { value: 'asset_type', label: 'Asset Type', example: 'asset_type=temperature' },
];

const compressionOptions = [
  { value: 'gzip', label: 'GZIP (Recommended)' },
  { value: 'snappy', label: 'Snappy' },
  { value: 'lz4', label: 'LZ4' },
  { value: 'none', label: 'No Compression' },
];

// Computed properties
const fullPathPreview = computed(() => {
  const base = modelRef.value.pathConfig.basePath || 'datalake';
  const partitions = modelRef.value.pathConfig.partitions.map(p => getPartitionExample(p)).join('/');
  const prefix = modelRef.value.pathConfig.filePrefix || 'data';
  const compression = modelRef.value.pathConfig.compression || 'gzip';

  const pathStructure = partitions ? `${base}/${partitions}/` : `${base}/`;
  const fileName = `${prefix}_20251215_140000.json${compression !== 'none' ? `.${compression}` : ''}`;

  return `${pathStructure}${fileName}`;
});

// Methods
const isPartitionSelected = (partitionValue: string): boolean => {
  return modelRef.value.pathConfig.partitions.includes(partitionValue);
};

const addPartition = (partitionValue: string) => {
  if (!isPartitionSelected(partitionValue)) {
    modelRef.value.pathConfig.partitions.push(partitionValue);
  }
};

const removePartition = (index: number) => {
  modelRef.value.pathConfig.partitions.splice(index, 1);
};

const getPartitionLabel = (partitionValue: string): string => {
  const partition = availablePartitions.find(p => p.value === partitionValue);
  return partition?.label || partitionValue;
};

const getPartitionExample = (partitionValue: string): string => {
  const partition = availablePartitions.find(p => p.value === partitionValue);
  return partition?.example || `${partitionValue}=example`;
};

const copyPathToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(fullPathPreview.value);
    notifySuccess({ message: 'Path copied to clipboard!' });
  } catch {
    notifyFail({ message: 'Error copying path' });
  }
};

// Initialize default values if not present
if (!modelRef.value.pathConfig) {
  modelRef.value.pathConfig = {
    basePath: 'datalake',
    partitions: [],
    compression: 'gzip',
    filePrefix: 'data_export_',
    maxFileSize: 100,
  };
}
</script>

<template>
  <div class="path-config">
    <div class="row q-col-gutter-md">
      <!-- Base Path Configuration -->
      <div class="col-12">
        <q-input
            v-model="modelRef.pathConfig.basePath"
            outlined
            placeholder="datalake"
            label="Base Path *"
            hint="Root directory where data will be stored (e.g., datalake, iot-data)"
            :rules="[
            val => !!val || 'Base path is required',
            val => !val.startsWith('/') || 'Path should not start with /',
            val => !val.endsWith('/') || 'Path should not end with /',
            val => /^[a-zA-Z0-9_-]+$/.test(val) || 'Use only letters, numbers, _ and -'
          ]"
        />
      </div>

      <!-- Compression -->
      <div class="col-12 col-md-6">
        <q-select
            v-model="modelRef.pathConfig.compression"
            outlined
            emit-value
            map-options
            option-value="value"
            option-label="label"
            label="Compression"
            hint="Compression type applied to JSON files"
            :options="compressionOptions"
        />
      </div>

      <!-- Custom Prefix -->
      <div class="col-12 col-md-6">
        <q-input
            v-model="modelRef.pathConfig.filePrefix"
            outlined
            placeholder="data_export_"
            label="File Prefix (optional)"
            hint="Prefix added to file names"
        />
      </div>

      <!-- Max File Size -->
      <div class="col-12">
        <q-input
            v-model.number="modelRef.pathConfig.maxFileSize"
            outlined
            type="number"
            min="1"
            max="1000"
            label="Maximum File Size (MB)"
            hint="Maximum size of each exported JSON file"
        />
      </div>

      <!-- Partitions Configuration -->
      <div class="col-12">
        <div class="text-subtitle1 text-weight-medium q-mb-md">Partition Configuration</div>
        <div class="text-body2 text-grey-6 q-mb-md">
          Partitions organize data into subdirectories for better performance and organization
        </div>

        <!-- Available Partitions -->
        <div class="q-mb-md">
          <div class="text-subtitle2 q-mb-sm">Available Partitions</div>
          <div class="row q-col-gutter-xs">
            <div
                v-for="partition in availablePartitions"
                :key="partition.value"
                class="col-auto"
            >
              <q-btn
                  outline
                  size="sm"
                  class="partition-btn"
                  :label="partition.label"
                  :color="isPartitionSelected(partition.value) ? 'grey-4' : 'primary'"
                  :text-color="isPartitionSelected(partition.value) ? 'grey-6' : 'primary'"
                  :disable="isPartitionSelected(partition.value)"
                  @click="addPartition(partition.value)"
              >
                <q-icon size="14px" class="q-mr-xs" name="add" />
              </q-btn>
            </div>
          </div>
        </div>

        <!-- Selected Partitions (Draggable with vue-draggable-plus) -->
        <div v-if="modelRef.pathConfig.partitions.length > 0" class="q-mb-md">
          <div class="text-subtitle2 q-mb-sm">
            <q-icon name="drag_indicator" class="q-mr-xs"/>
            Selected Partitions (drag to reorder)
          </div>

          <div class="selected-partitions-container">
            <VueDraggable
                v-model="modelRef.pathConfig.partitions"
                class="draggable-list"
                ghost-class="partition-ghost"
                chosen-class="partition-chosen"
                drag-class="partition-drag"
                handle=".drag-handle"
                :animation="200"
            >
              <div
                  v-for="(partition, index) in modelRef.pathConfig.partitions"
                  :key="`${partition}-${index}`"
                  class="partition-item q-mb-sm"
              >
                <q-card
                    flat
                    bordered
                    class="partition-card"
                >
                  <q-card-section class="q-pa-md">
                    <div class="row items-center no-wrap">
                      <!-- Drag Handle -->
                      <div class="col-auto">
                        <q-icon
                            size="20px"
                            class="drag-handle cursor-move"
                            color="grey-6"
                            name="drag_indicator"
                        />
                      </div>

                      <!-- Partition Info -->
                      <div class="col q-ml-sm">
                        <div class="text-body2 text-weight-medium">
                          {{ getPartitionLabel(partition) }}
                        </div>
                        <div class="text-caption text-grey-6">
                          {{ getPartitionExample(partition) }}
                        </div>
                      </div>

                      <!-- Order Number -->
                      <div class="col-auto">
                        <DetailChip
                          :label="`${index + 1}`"
                          color="primary"
                          size="sm"
                          class="order-chip"
                        />
                      </div>

                      <!-- Remove Button -->
                      <div class="col-auto q-ml-sm">
                        <q-btn
                            flat
                            round
                            size="sm"
                            class="rounded-borders"
                            icon="close"
                            color="negative"
                            @click="removePartition(index)"
                        >
                          <AppTooltip content="Remove partition" />
                        </q-btn>
                      </div>
                    </div>
                  </q-card-section>
                </q-card>
              </div>
            </VueDraggable>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="empty-partitions rounded-borders">
          <q-card flat bordered class="text-center q-pa-xl">
            <q-icon size="64px" color="grey-4" name="folder_open"/>
            <div class="text-h6 text-grey-6 q-mt-md">
              No partitions selected
            </div>
            <div class="text-body2 text-grey-5 q-mt-sm">
              Add partitions to better organize your data
            </div>
            <q-btn
                outline
                class="q-mt-md"
                label="Add First Partition"
                color="primary"
                @click="addPartition('year')"
            />
          </q-card>
        </div>
      </div>

      <!-- Path Preview -->
      <div class="col-12">
        <q-card flat bordered class="bg-blue-1">
          <q-card-section>
            <div class="text-subtitle2 text-primary q-mb-sm">
              <q-icon class="q-mr-xs" name="preview"/>
              Complete Path Preview
            </div>

            <!-- Full Path Example -->
            <div class="path-preview-container q-mb-md">
              <code class="path-preview">{{ fullPathPreview }}</code>
              <q-btn
                  flat
                  round
                  size="sm"
                  class="q-ml-sm"
                  icon="content_copy"
                  color="primary"
                  @click="copyPathToClipboard"
              >
                <AppTooltip content="Copy path" />
              </q-btn>
            </div>

            <!-- Additional Info -->
            <div class="text-caption text-grey-6 q-mt-md">
              <q-icon size="14px" class="q-mr-xs" name="info" />
              This is an example of how data will be organized in storage
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.path-config {
  .partition-btn {
    border-radius: var(--mapex-radius-xl);
    transition: var(--mapex-transition-base);

    &:hover:not(:disabled) {
      transform: translateY(-1px);
      box-shadow: var(--mapex-shadow-sm);
    }
  }

  .selected-partitions-container {
    .draggable-list {
      min-height: 60px;
    }

    .partition-item {
      transition: var(--mapex-transition-base);

      &:hover {
        transform: translateY(-1px);
      }
    }

    .partition-card {
      border-radius: var(--mapex-radius-md);
      transition: var(--mapex-transition-base);

      &:hover {
        box-shadow: var(--mapex-shadow-md);
      }
    }

    .drag-handle {
      transition: color 0.2s ease;

      &:hover {
        color: var(--q-primary) !important;
      }
    }

    .order-chip {
      font-weight: 600;
    }

    // Dragging states
    .partition-ghost {
      opacity: 0.5;
      transform: rotate(2deg);
    }

    .partition-chosen {
      transform: scale(1.02);
      box-shadow: var(--mapex-shadow-lg);
    }

    .partition-drag {
      transform: rotate(5deg);
      opacity: 0.8;
    }
  }

  .empty-partitions {
    .q-card {
      border: 2px dashed var(--mapex-card-border);
      transition: var(--mapex-transition-base);

      &:hover {
        border-color: var(--mapex-card-hover-border);
        background-color: var(--mapex-surface-highlight);
      }
    }
  }

  .path-preview-container {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .path-preview {
    display: inline-block;
    background: var(--mapex-surface-bg);
    padding: 12px 16px;
    border-radius: var(--mapex-radius-md);
    font-family: 'Roboto Mono', 'Courier New', monospace;
    font-size: 13px;
    color: var(--q-info);
    border: 1px solid var(--mapex-card-border);
    word-break: break-all;
    flex: 1;
    font-weight: 500;
  }

  .path-breakdown {
    background: var(--mapex-page-bg);
    padding: 16px;
    border-radius: var(--mapex-radius-md);
    border: 1px solid var(--mapex-card-border);

    .q-chip {
      font-weight: 500;
    }
  }
}

// Global dragging cursor
.vue-draggable-plus-ghost {
  cursor: grabbing !important;
}
</style>
<script setup lang="ts">
defineOptions({
  name: 'LakeHouseSelectorDrawer'
});

/** TYPE IMPORTS */
import type { LakeHouseItem, LakeHouseSelectorDrawerProps, LakeHouseSelectorDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useCommonPlaceholders } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('LakeHouseSelectorDrawer');
const { placeholders } = useCommonPlaceholders();

/** PROPS & EMITS */
const props = withDefaults(defineProps<LakeHouseSelectorDrawerProps>(), {
  selectedLakeHouseId: null,
});

const emit = defineEmits<LakeHouseSelectorDrawerEmits>();

/** STATE */
const loading = ref(false);
const lakeHouses = ref<LakeHouseItem[]>([]);

/** FILTER STATE */
const filters = ref({
  name: undefined as string | undefined,
  status: undefined as boolean | undefined,
  type: undefined as string | undefined,
});

/** COMPUTED */

/**
 * Drawer visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Status filter options
 */
const statusOptions = computed(() => [
  { label: 'All', value: undefined },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false },
]);

/**
 * Type filter options
 */
const typeOptions = computed(() => [
  { label: 'All Types', value: undefined },
  { label: 'AWS S3', value: 'aws-s3' },
  { label: 'Azure Blob', value: 'azure-blob' },
  { label: 'Google Cloud Storage', value: 'gcp-storage' },
  { label: 'MinIO', value: 'minio' },
]);

/**
 * Filtered data lakes based on current filters
 */
const filteredLakeHouses = computed(() => {
  return lakeHouses.value.filter(dl => {
    // Filter by name
    if (filters.value.name) {
      const nameMatch = dl.name.toLowerCase().includes(filters.value.name.toLowerCase());
      const descMatch = dl.description?.toLowerCase().includes(filters.value.name.toLowerCase());
      if (!nameMatch && !descMatch) return false;
    }

    // Filter by status
    if (typeof filters.value.status === 'boolean') {
      if (dl.status !== filters.value.status) return false;
    }

    // Filter by type
    if (filters.value.type) {
      if (dl.type !== filters.value.type) return false;
    }

    return true;
  });
});

/** WATCHERS */

/**
 * Watch drawer open state and load data
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchLakeHouses();
  }
});

/** FUNCTIONS */

/**
 * Get icon name based on data lake type
 * @param {string} type - Data lake type
 * @returns {string} Icon name
 */
function getLakeHouseIcon(type: string): string {
  switch (type) {
    case 'aws-s3':
      return 'mdi-aws';
    case 'azure-blob':
      return 'mdi-microsoft-azure';
    case 'gcp-storage':
      return 'mdi-google-cloud';
    case 'minio':
      return 'mdi-database';
    default:
      return 'storage';
  }
}

/**
 * Get icon color based on data lake type
 * @param {string} type - Data lake type
 * @returns {string} Color class
 */
function getLakeHouseIconColor(type: string): string {
  switch (type) {
    case 'aws-s3':
      return 'orange-6';
    case 'azure-blob':
      return 'blue-6';
    case 'gcp-storage':
      return 'red-6';
    case 'minio':
      return 'purple-6';
    default:
      return 'purple-6';
  }
}

/**
 * Get type label for display
 * @param {string} type - Data lake type
 * @returns {string} Type label
 */
function getTypeLabel(type: string): string {
  switch (type) {
    case 'aws-s3':
      return 'AWS S3';
    case 'azure-blob':
      return 'Azure Blob';
    case 'gcp-storage':
      return 'GCP Storage';
    case 'minio':
      return 'MinIO';
    default:
      return type;
  }
}

/**
 * Fetch data lakes - using stub data for now
 * TODO: Replace with actual API call when available
 * @returns {Promise<void>}
 */
async function fetchLakeHouses(): Promise<void> {
  loading.value = true;

  try {
    // Simulated delay for loading state
    await new Promise(resolve => setTimeout(resolve, 300));

    // Using stub data - TODO: Replace with API call
    lakeHouses.value = [
      {
        id: 'dl-1',
        type: 'aws-s3',
        name: 'Amazon S3',
        description: 'Amazon S3 data lake configured for daily partitioned uploads.',
        status: true,
        pathConfig: {
          maxFileSize: 100,
          partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
        },
        credentials: {
          bucket: 'my-amazon-s3-bucket',
          region: 'us-east-1'
        },
        frequency: {
          interval: 1,
          type: 'day',
          time: '00:00'
        },
      },
      {
        id: 'dl-2',
        type: 'azure-blob',
        name: 'Azure Blob Storage',
        description: 'Azure Blob Storage data lake with scalable object storage.',
        status: true,
        pathConfig: {
          maxFileSize: 100,
          partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
        },
        credentials: {
          bucket: 'my-azure-blob-container',
          region: 'eastus'
        },
        frequency: {
          interval: 1,
          type: 'day',
          time: '00:00'
        },
      },
      {
        id: 'dl-3',
        type: 'gcp-storage',
        name: 'Google Cloud Storage',
        description: 'Google Cloud Storage data lake for high-throughput processing.',
        status: true,
        pathConfig: {
          maxFileSize: 100,
          partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
        },
        credentials: {
          bucket: 'my-gcp-storage-bucket',
          region: 'us-central1'
        },
        frequency: {
          interval: 1,
          type: 'day',
          time: '00:00'
        },
      },
      {
        id: 'dl-4',
        type: 'minio',
        name: 'MinIO',
        description: 'MinIO instance for local S3-compatible object storage.',
        status: false,
        pathConfig: {
          maxFileSize: 100,
          partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
        },
        credentials: {
          bucket: 'my-minio-bucket',
          region: 'minio-local'
        },
        frequency: {
          interval: 1,
          type: 'day',
          time: '00:00'
        },
      },
    ];
  } catch (error: any) {
    logger.error('Failed to fetch data lakes:', error);
  } finally {
    loading.value = false;
  }
}

/**
 * Select data lake and close drawer
 * @param {LakeHouseItem} lakeHouse - Data lake to select
 */
function selectLakeHouse(lakeHouse: LakeHouseItem): void {
  emit('select', lakeHouse);
  showDialog.value = false;
}

/**
 * Check if data lake is selected
 * @param {LakeHouseItem} lakeHouse - Data lake to check
 * @returns {boolean} True if selected
 */
function isSelected(lakeHouse: LakeHouseItem): boolean {
  return lakeHouse.id === props.selectedLakeHouseId;
}

/**
 * Cancel handler
 */
function handleCancel(): void {
  emit('cancel');
  showDialog.value = false;
}

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    handleCancel();
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<template>
  <q-dialog v-model="showDialog" position="right" maximized>
    <q-card style="width: 600px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
      <!-- Header with Padding -->
      <q-card-section class="q-pb-sm">
        <div class="row items-center">
          <q-icon name="storage" size="md" color="purple-6" class="q-mr-sm" />
          <div class="text-h6">Select Data Lake</div>
          <q-space />
          <q-btn icon="close" flat round dense class="rounded-borders" @click="handleCancel" />
        </div>
      </q-card-section>

      <!-- Info Banner -->
      <q-card-section class="q-pt-none q-pb-md">
        <q-banner dense class="bg-purple-1 text-purple-9 rounded-borders">
          <template #avatar>
            <q-icon name="info" color="purple-6" size="sm" />
          </template>
          <div class="text-caption">
            Data Lakes are cloud storage destinations for processed event data, supporting AWS S3, Azure Blob, GCP Storage, and MinIO.
          </div>
        </q-banner>
      </q-card-section>

      <!-- Filters -->
      <q-card-section class="q-py-md">
        <div class="text-overline text-grey-7 q-mb-md">
          <q-icon name="filter_list" size="xs" class="q-mr-xs" />
          Filters
        </div>
        <div class="row q-col-gutter-md">
          <!-- Search - Full width -->
          <div class="col-12">
            <q-input
              v-model="filters.name"
              outlined
              dense
              label="Search by name"
              :placeholder="placeholders.typeToSearch.value"
              clearable
              class="rounded-borders"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Type -->
          <div class="col-6">
            <q-select
              v-model="filters.type"
              outlined
              dense
              label="Type"
              class="rounded-borders"
              :options="typeOptions"
              emit-value
              map-options
            >
              <template #prepend>
                <q-icon name="category" />
              </template>
            </q-select>
          </div>

          <!-- Status -->
          <div class="col-6">
            <q-select
              v-model="filters.status"
              outlined
              dense
              label="Status"
              class="rounded-borders"
              :options="statusOptions"
              emit-value
              map-options
            >
              <template #prepend>
                <q-icon name="toggle_on" />
              </template>
            </q-select>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="storage" size="xs" class="q-mr-xs" />
          Results
        </div>
      </q-card-section>

      <!-- Data Lakes List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="loading" class="q-pa-md text-center">
          <q-spinner color="purple-6" size="3em" />
          <div class="text-grey-7 q-mt-md">Loading data lakes...</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="filteredLakeHouses.length === 0" class="q-pa-md text-center">
          <q-icon name="inbox" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">No data lakes found</div>
        </div>

        <!-- Data Lakes List -->
        <q-scroll-area
          v-else
          style="height: 100%;"
        >
          <q-list separator>
            <q-item
              v-for="lakeHouse in filteredLakeHouses"
              :key="lakeHouse.id || `dl-${Math.random()}`"
              clickable
              :active="isSelected(lakeHouse)"
              @click="selectLakeHouse(lakeHouse)"
            >
              <q-item-section avatar>
                <q-avatar
                  :color="getLakeHouseIconColor(lakeHouse.type)"
                  :icon="getLakeHouseIcon(lakeHouse.type)"
                  text-color="white"
                />
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ lakeHouse.name || 'Unnamed Data Lake' }}</q-item-label>
                <q-item-label caption class="text-grey-7">
                  <span v-if="lakeHouse.description">{{ lakeHouse.description }}</span>
                  <span v-else class="text-grey-5">No description</span>
                </q-item-label>
                <q-item-label caption class="q-mt-xs">
                  <DetailChip
                    :value="getTypeLabel(lakeHouse.type)"
                    :color="getLakeHouseIconColor(lakeHouse.type) as any"
                    size="sm"
                    dense
                  />
                  <span v-if="lakeHouse.credentials?.bucket" class="q-ml-sm text-grey-6">
                    <q-icon name="folder" size="xs" class="q-mr-xs" />
                    {{ lakeHouse.credentials.bucket }}
                  </span>
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-badge
                  :color="lakeHouse.status ? 'green-6' : 'red-6'"
                  :label="lakeHouse.status ? 'ACTIVE' : 'INACTIVE'"
                />
              </q-item-section>
            </q-item>
          </q-list>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator />
      <q-card-actions class="row items-center q-px-md q-py-md">
        <div class="text-caption text-grey-7">
          <q-icon name="storage" size="xs" class="q-mr-xs" />
          {{ filteredLakeHouses.length }} {{ filteredLakeHouses.length === 1 ? 'data lake' : 'data lakes' }}
        </div>
        <q-space />
        <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders" @click="handleCancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

/* Hover effects for list items */
:deep(.q-item) {
  transition: all var(--mapex-transition-base) ease;
}

:deep(.q-item:hover) {
  background-color: var(--mapex-surface-bg);
}

:deep(.q-item.q-item--active) {
  background-color: var(--mapex-active-bg) !important;
  border-left: 3px solid var(--q-accent);
}

/* Better spacing for filter inputs */
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

/* Smooth transitions */
:deep(.q-badge),
:deep(.q-chip) {
  transition: all var(--mapex-transition-base) ease;
}

/* Footer padding */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>

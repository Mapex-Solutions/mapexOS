<template>
  <!-- Invisible backdrop for click outside detection -->
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="drawer-backdrop"
      @click="close"
    />
  </Teleport>

  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
    @keydown.esc="close"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon
        name="domain"
        size="sm"
        color="primary"
        class="q-mr-sm"
      />
      <q-toolbar-title class="text-weight-medium">{{ t.drawer.title.value }}</q-toolbar-title>

      <q-btn
        flat
        round
        dense
        icon="close"
        class="drawer-close-btn"
        @click="close"
      >
        <AppTooltip :content="t.drawer.close.value" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner
            size="3em"
            color="primary"
            class="q-mb-md"
          />
          <div class="text-grey-7">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner
            rounded
            class="bg-negative text-white"
          >
            <template #avatar>
              <q-icon
                name="error"
                color="white"
              />
            </template>
            {{ t.drawer.error.value }}
          </q-banner>
        </div>

        <!-- Customer Data -->
        <div v-else-if="customer" class="q-px-md q-py-lg">

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="info"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ customer?.name || '-' }}</div>
            </div>

            <!-- Type & Status (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.type.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="getTypeIcon(customer?.type || '')"
                      :color="getTypeColorName(customer?.type || '')"
                      size="sm"
                      :label="getTypeLabel(customer?.type || '')"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.status.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="customer?.enabled ? 'positive' : 'negative'"
                      size="sm"
                      :label="customer?.enabled ? t.status.active.value.toUpperCase() : t.status.inactive.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Hierarchy Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="account_tree"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.hierarchy.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.customerOf.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="domain"
                  color="green"
                  size="sm"
                  :label="(customer as any)?.organizationName || customer?.customerId || '-'"
                />
              </div>
            </div>

            <!-- Code & Path Key (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.code.value }}</div>
                  <div class="field-value text-grey-8">
                    <code class="bg-grey-2 q-pa-xs rounded-borders">{{ (customer as any)?.code || '-' }}</code>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.pathKey.value }}</div>
                  <div class="field-value text-grey-8">
                    <code class="bg-grey-2 q-pa-xs rounded-borders">{{ (customer as any)?.pathKey || '-' }}</code>
                  </div>
                </div>
              </div>
            </div>

            <!-- Depth & Child Count (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.depth.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      color="teal"
                      size="sm"
                      :label="String((customer as any)?.depth || 0)"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.childCount.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      color="orange"
                      size="sm"
                      :label="String((customer as any)?.childCount || 0)"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Address Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                name="location_on"
                size="sm"
                color="primary"
                class="q-mr-sm"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.address.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Street (full width) -->
            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.street.value }}</div>
              <div class="field-value">{{ (customer as any)?.address?.street || '-' }}</div>
            </div>

            <!-- City & State (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.city.value }}</div>
                  <div class="field-value">{{ customer?.address?.city || '-' }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.state.value }}</div>
                  <div class="field-value">{{ customer?.address?.state || '-' }}</div>
                </div>
              </div>
            </div>

            <!-- Country & Zip Code (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.country.value }}</div>
                  <div class="field-value">{{ customer?.address?.country || '-' }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.zipCode.value }}</div>
                  <div class="field-value">{{ customer?.address?.zipCode || '-' }}</div>
                </div>
              </div>
            </div>

            <div v-if="!hasAddress" class="text-grey-7 text-caption q-mt-sm">
              {{ t.drawer.empty.address.value }}
            </div>
          </div>

          <!-- Contact Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon
                size="sm"
                name="phone"
                class="q-mr-sm"
                color="primary"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.contact.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.phone.value }}</div>
              <div class="field-value">
                <DetailChip
                  v-if="customer?.phone"
                  icon="phone"
                  color="blue"
                  size="sm"
                  :label="customer.phone"
                />
                <span v-else class="text-grey-7">{{ t.drawer.empty.phone.value }}</span>
              </div>
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header">
              <q-icon
                size="sm"
                name="schedule"
                class="q-mr-sm"
                color="primary"
              />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value">{{ formatDate((customer as any)?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate((customer as any)?.updated) }}</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </q-scroll-area>
    </div>

    <!-- Footer Actions -->
    <q-separator />
    <div class="drawer-footer">
      <q-space />
      <q-btn
        unelevated
        icon="edit"
        color="primary"
        :label="t.drawer.edit.value"
        :disable="!customer"
        @click="handleEdit"
      />
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'CustomerDetailsDrawer'
});

/** TYPE IMPORTS */
import type { CustomerDetailsDrawerProps, CustomerDetailsDrawerEmits } from './interfaces/customerDetailsDrawer.interface';
import type { OrganizationResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useCustomersTranslations, useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

const props = defineProps<CustomerDetailsDrawerProps>();
const emit = defineEmits<CustomerDetailsDrawerEmits>();

const t = useCustomersTranslations();
const errors = useCommonErrors();
const logger = useLogger('CustomerDetailsDrawer');

// Component state
const customer = ref<OrganizationResponse | null>(null);
const loading = ref(false);
const error = ref(false);

// Handle ESC key to close drawer
function handleEscKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});

/**
 * Fetch customer details by ID from API
 */
async function fetchCustomerDetails(customerId: string) {
  if (!apis.mapexOS?.organizations) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  customer.value = null;

  try {
    const response = await apis.mapexOS.organizations.getById({ organizationId: customerId });
    customer.value = response;
  } catch (err: any) {
    logger.error('Error fetching customer details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

// Watch for customerId changes
watch(() => props.customerId, (newCustomerId) => {
  if (newCustomerId && props.modelValue) {
    void fetchCustomerDetails(newCustomerId);
  }
}, { immediate: true });

// Watch for drawer open/close
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.customerId) {
    void fetchCustomerDetails(props.customerId);
  } else if (!isOpen) {
    customer.value = null;
    error.value = false;
  }
});

// Check if customer has address
const hasAddress = computed(() => {
  return !!((customer.value as any)?.address?.street || customer.value?.address?.city ||
            customer.value?.address?.state || customer.value?.address?.country ||
            customer.value?.address?.zipCode);
});

/**
 * Get icon for organization type
 */
function getTypeIcon(type: string): string {
  const iconMap: Record<string, string> = {
    customer: 'domain',
    site: 'location_on',
    building: 'apartment',
    floor: 'layers',
    zone: 'place',
  };
  return iconMap[type] || 'domain';
}

/**
 * Get color name for organization type (DetailChip format)
 * @param {string} type - Organization type
 * @returns {string} Color name compatible with DetailChip
 */
function getTypeColorName(type: string): 'green' | 'orange' | 'blue' | 'teal' | 'primary' {
  const colorMap: Record<string, 'green' | 'orange' | 'blue' | 'teal' | 'primary'> = {
    customer: 'green',
    site: 'orange',
    building: 'blue',
    floor: 'teal',
    zone: 'green',
  };
  return colorMap[type] || 'primary';
}

/**
 * Get label for organization type
 */
function getTypeLabel(type: string): string {
  if (!type) return '-';
  return t.drawer.type[type as keyof typeof t.drawer.type]?.value.toUpperCase() || type.toUpperCase();
}

/**
 * Format date using Quasar date utils
 */
function formatDate(dateValue: any): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

// Close drawer
function close() {
  emit('update:modelValue', false);
}

// Handle edit action
function handleEdit() {
  if (!customer.value?.id) return;
  emit('edit', customer.value.id);
  close();
}
</script>

<style lang="scss" scoped>
// Flex layout for drawer content
:deep(.q-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// Drawer Header
.drawer-header {
  flex-shrink: 0;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Backdrop (teleported to body, needs :global) - transparent, just for click detection
:global(.drawer-backdrop) {
  position: fixed;
  top: 0;
  left: 0;
  right: 450px; // Leave space for drawer (450px width)
  bottom: 0;
  background: transparent;
  z-index: 5999; // Below q-drawer (6000)
  cursor: default;
}

// Drawer Content
.drawer-content {
  flex: 1;
  min-height: 0; // Important for flex children with overflow
  overflow: hidden;

  :deep(.q-scrollarea__content) {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

// Section Styling
.section {
  .section-header {
    display: flex;
    align-items: center;
    color: var(--q-primary);
    margin-bottom: 8px;
  }
}

// Field Row Styling
.field-row {
  display: flex;
  flex-direction: column;
  padding: 10px 0;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-child {
    border-bottom: none;
  }

  .field-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    margin-bottom: 4px;
    letter-spacing: 0.8px;
  }

  .field-value {
    font-size: 0.9rem;
    color: var(--mapex-text-primary);
    word-break: break-word;
    line-height: 1.4;

    code {
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
    }
  }
}

// Custom Scrollbar
:deep(.q-scrollarea__content) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: var(--mapex-radius-lg);
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(var(--q-primary-rgb), 0.3);
    border-radius: var(--mapex-radius-lg);
    transition: background var(--mapex-transition-base) ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}
</style>

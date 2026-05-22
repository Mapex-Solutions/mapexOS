<script setup lang="ts">
defineOptions({
  name: 'PersistenceSettings'
});

/** TYPE IMPORTS */
import type { RetentionPolicyResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useSettingsTranslations } from '@composables/i18n';
import { useCommonActions } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { dialogConfirm } from '@utils/alert';
import { showInfo } from '@utils/modal';
import { notifySuccess, notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

const { persistence } = useSettingsTranslations();
const { actions } = useCommonActions();
const logger = useLogger('PersistenceSettings');

/** LOCAL IMPORTS */
import type { RetentionPolicyLimits, RetentionPolicy } from './interfaces/persistenceSettings.interface';

const retentionPoliciesLimits: Record<string, RetentionPolicyLimits> = {
  eventsRaw: {
    name: 'Raw Events',
    defaultDays: 1,
    minDays: 1,
    maxDays: 3,
  },
  eventsJsExecutor: {
    name: 'JS Executor Events',
    defaultDays: 3,
    minDays: 1,
    maxDays: 3,
  },
  eventsRouter: {
    name: 'Router Events',
    defaultDays: 1,
    minDays: 1,
    maxDays: 3,
  },
  events: {
    name: 'Events',
    defaultDays: 90,
    minDays: 7,
    maxDays: 365,
  },
  eventsWorkflow: {
    name: 'Workflow Execution Events',
    defaultDays: 7,
    minDays: 1,
    maxDays: 365,
  },
  eventsAudit: {
    name: 'Audit Events',
    defaultDays: 365,
    minDays: 365,
    maxDays: 2555,
  },
  eventsNotifications: {
    name: 'Notification Events',
    defaultDays: 30,
    minDays: 1,
    maxDays: 365,
  },
};

/** Policy icon mapping */
const POLICY_ICONS: Record<string, { icon: string; color: string }> = {
  eventsRaw: { icon: 'input', color: 'orange-8' },
  eventsJsExecutor: { icon: 'code', color: 'purple-8' },
  eventsRouter: { icon: 'alt_route', color: 'cyan-8' },
  events: { icon: 'event_note', color: 'blue-8' },
  eventsWorkflow: { icon: 'account_tree', color: 'purple-8' },
  eventsAudit: { icon: 'verified_user', color: 'indigo-8' },
  eventsNotifications: { icon: 'notifications', color: 'amber-8' },
};

/** Group keys */
const PIPELINE_KEYS = ['eventsRaw', 'eventsJsExecutor', 'eventsRouter', 'eventsWorkflow'];
const STORAGE_KEYS = ['events', 'eventsNotifications', 'eventsAudit'];

/**
 * Initialize policies from API data or defaults
 * Creates local editable state from retention policies data
 *
 * @param {RetentionPolicyResponse[]} apiPolicies - Policies from API response
 * @returns {RetentionPolicy[]} Initialized local policy state
 */
function initializePolicies(apiPolicies?: RetentionPolicyResponse[]): RetentionPolicy[] {
  const policyMap = new Map<string, RetentionPolicyResponse>();
  if (apiPolicies) {
    for (const p of apiPolicies) {
      if (p.type) {
        policyMap.set(p.type, p);
      }
    }
  }

  return Object.entries(retentionPoliciesLimits).map(([key, limits]) => {
    const apiPolicy = policyMap.get(key);
    const currentValue = apiPolicy?.retentionDays ?? limits.defaultDays;
    return {
      key,
      name: limits.name,
      currentDays: currentValue,
      defaultDays: currentValue,
      minDays: limits.minDays,
      maxDays: limits.maxDays,
    };
  });
}

/** STATE */
const loading = ref(false);
const saving = ref(false);
const policies = ref<RetentionPolicy[]>(initializePolicies());

/** COMPUTED */
const pipelineEvents = computed(() =>
  policies.value.filter(p => PIPELINE_KEYS.includes(p.key))
);

const storageEvents = computed(() =>
  policies.value.filter(p => STORAGE_KEYS.includes(p.key))
);

const hasErrors = computed(() => {
  return policies.value.some(policy => validatePolicy(policy) !== true);
});

const hasChanges = computed(() => {
  return policies.value.some(policy => policy.currentDays !== policy.defaultDays);
});

/**
 * Check if pipeline group has any low retention policies
 * @returns {boolean} True if group has LakeHouse warning
 */
const hasPipelineLakeHouseWarning = computed(() => {
  return pipelineEvents.value.some(p => isLowRetentionPolicy(p.key));
});

/**
 * Check if storage group has audit compliance warning
 * @returns {boolean} True if group has compliance warning
 */
const hasStorageComplianceWarning = computed(() => {
  return storageEvents.value.some(p => isAuditPolicy(p.key));
});

/** FUNCTIONS */

/**
 * Fetch retention policies from events API
 */
async function fetchPolicies() {
  if (!apis.events?.retention) {
    logger.error('Events retention API not initialized');
    return;
  }

  loading.value = true;

  try {
    const response = await apis.events.retention.listRetentionPolicies({
      perPage: 50,
    });

    policies.value = initializePolicies(response.items);
  } catch (err: any) {
    logger.error('Error fetching retention policies:', err);
    notifyFail({ message: err.message || 'Failed to load retention policies' });
  } finally {
    loading.value = false;
  }
}

/**
 * Get translated event type name
 * @param {string} key - Policy key
 * @returns {string} Translated name
 */
function getEventTypeName(key: string): string {
  const eventTypeKey = key as keyof typeof persistence.eventTypes;
  return persistence.eventTypes[eventTypeKey]?.value || key;
}

/**
 * Show event info modal with description
 * @param {string} eventKey - Policy key
 */
async function showEventInfo(eventKey: string) {
  const descriptionKey = eventKey as keyof typeof persistence.eventDescriptions;
  const info = persistence.eventDescriptions[descriptionKey];

  await showInfo(
    info.title.value,
    info.description.value
  );
}

/**
 * Validate a single retention policy value
 * @param {RetentionPolicy} policy - Policy to validate
 * @returns {string | boolean} True if valid, error message string if invalid
 */
function validatePolicy(policy: RetentionPolicy): string | boolean {
  if (!policy.currentDays && policy.currentDays !== 0) {
    return persistence.validation.required.value;
  }

  const value = Number(policy.currentDays);

  if (!Number.isInteger(value)) {
    return persistence.validation.integer.value;
  }

  if (value < policy.minDays) {
    return persistence.validation.minValue.value.replace('{min}', policy.minDays.toString());
  }

  if (value > policy.maxDays) {
    return persistence.validation.maxValue.value.replace('{max}', policy.maxDays.toString());
  }

  return true;
}

/**
 * Get hint text showing min/max range for a policy
 * @param {RetentionPolicy} policy - Policy to get hint for
 * @returns {string} Formatted hint string
 */
function getDaysHint(policy: RetentionPolicy): string {
  return persistence.form.daysHint.value
    .replace('{min}', policy.minDays.toString())
    .replace('{max}', policy.maxDays.toString());
}

/**
 * Get default label text for a policy
 * @param {RetentionPolicy} policy - Policy to get label for
 * @returns {string} Formatted default label
 */
function getDefaultLabel(policy: RetentionPolicy): string {
  return persistence.form.defaultLabel.value
    .replace('{days}', policy.defaultDays.toString());
}

/**
 * Get icon config for a policy key
 * @param {string} key - Policy key
 * @returns {{ icon: string; color: string }} Icon and color
 */
function getPolicyIcon(key: string): { icon: string; color: string } {
  return POLICY_ICONS[key] || { icon: 'event', color: 'grey-8' };
}

/**
 * Check if policy is an audit policy
 * @param {string} key - Policy key
 * @returns {boolean} True if audit policy
 */
function isAuditPolicy(key: string): boolean {
  return key === 'eventsAudit';
}

/**
 * Check if policy is a low retention LakeHouse-eligible policy
 * @param {string} key - Policy key
 * @returns {boolean} True if low retention policy
 */
function isLowRetentionPolicy(key: string): boolean {
  return ['eventsRaw', 'eventsJsExecutor', 'eventsRouter', 'eventsBusinessRule'].includes(key);
}

/**
 * Save retention policies via events API
 * Upserts each changed policy individually
 */
async function handleSave() {
  if (hasErrors.value || !apis.events?.retention) {
    return;
  }

  saving.value = true;

  try {
    const changedPolicies = policies.value.filter(p => p.currentDays !== p.defaultDays);

    for (const policy of changedPolicies) {
      await apis.events.retention.upsertRetentionPolicy({
        type: policy.key,
        name: policy.name,
        retentionDays: policy.currentDays,
        enabled: true,
      });
    }

    notifySuccess({ message: persistence.messages.savedSuccessfully.value });

    // Refresh data after successful save
    await fetchPolicies();
  } catch (err: any) {
    logger.error('Error saving retention policies:', err);
    notifyFail({ message: err.message || 'Failed to save retention policies' });
  } finally {
    saving.value = false;
  }
}

/**
 * Confirm reset action with user
 */
async function confirmReset() {
  const confirmed = await dialogConfirm({
    title: actions.confirm.value,
    message: persistence.messages.confirmReset.value,
  });

  if (confirmed) {
    handleReset();
  }
}

/**
 * Reset policies to original values from last fetch
 * Does not save - just resets local state
 */
function handleReset() {
  policies.value = policies.value.map(p => ({
    ...p,
    currentDays: p.defaultDays,
  }));
}

/** LIFECYCLE HOOKS */
onMounted(() => void fetchPolicies());
</script>

<template>
  <div class="persistence-settings">
    <!-- Header -->
    <div class="row items-center q-mt-md q-mb-lg q-pl-md">
      <q-icon name="storage" size="sm" color="primary" class="q-mr-sm"/>
      <div class="text-subtitle1 text-weight-medium text-primary">{{ persistence.title.value }}</div>
    </div>

    <div class="q-px-md q-mb-md">
      <p class="text-body2 text-grey-7 q-ma-none">
        {{ persistence.description.value }}
      </p>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Content -->
    <q-form v-else @submit="handleSave">
      <!-- Group 1: Pipeline Events -->
      <div class="policy-group q-mb-lg q-px-md">
        <div class="group-header group-header--pipeline q-mb-sm">
          <q-icon name="device_hub" size="xs" color="orange-8" class="q-mr-sm" />
          <div>
            <div class="text-subtitle2 text-weight-bold">{{ persistence.groups.pipelineEvents.value }}</div>
            <div class="text-caption text-grey-6">{{ persistence.groups.pipelineDescription.value }}</div>
          </div>
        </div>

        <!-- LakeHouse Warning Banner -->
        <div v-if="hasPipelineLakeHouseWarning" class="warning-banner warning-banner--pipeline q-mb-sm">
          <q-icon name="warning" size="xs" color="orange-8" class="q-mr-sm" />
          <span class="text-caption">{{ persistence.messages.lakeHouseWarning.value }}</span>
        </div>

        <!-- Pipeline Rows -->
        <div
          v-for="policy in pipelineEvents"
          :key="policy.key"
          class="policy-row policy-row--pipeline"
        >
          <div class="policy-row__icon">
            <q-avatar
              :color="getPolicyIcon(policy.key).color"
              text-color="white"
              size="32px"
              font-size="16px"
              :icon="getPolicyIcon(policy.key).icon"
            />
          </div>

          <div class="policy-row__name">
            <div class="text-body2 text-weight-medium">{{ getEventTypeName(policy.key) }}</div>
            <div class="text-caption text-grey-6">{{ getDefaultLabel(policy) }}</div>
          </div>

          <div class="policy-row__input">
            <q-input
              v-model.number="policy.currentDays"
              outlined
              dense
              hide-bottom-space
              type="number"
              :suffix="persistence.form.daysLabel.value"
              :rules="[() => validatePolicy(policy)]"
              class="compact-input"
            />
          </div>

          <div class="policy-row__limits text-caption text-grey-6">
            {{ getDaysHint(policy) }}
          </div>

          <div class="policy-row__action">
            <q-icon
              name="info_outline"
              color="primary"
              size="20px"
              class="cursor-pointer"
              @click="showEventInfo(policy.key)"
            >
              <AppTooltip :content="persistence.tooltips.moreInfo.value" />
            </q-icon>
          </div>
        </div>
      </div>

      <!-- Group 2: Storage & Compliance -->
      <div class="policy-group q-mb-lg q-px-md">
        <div class="group-header group-header--storage q-mb-sm">
          <q-icon name="cloud_done" size="xs" color="blue-8" class="q-mr-sm" />
          <div>
            <div class="text-subtitle2 text-weight-bold">{{ persistence.groups.storageCompliance.value }}</div>
            <div class="text-caption text-grey-6">{{ persistence.groups.storageDescription.value }}</div>
          </div>
        </div>

        <!-- Compliance Warning Banner -->
        <div v-if="hasStorageComplianceWarning" class="warning-banner warning-banner--storage q-mb-sm">
          <q-icon name="verified_user" size="xs" color="blue-8" class="q-mr-sm" />
          <span class="text-caption">{{ persistence.messages.complianceWarning.value }}</span>
        </div>

        <!-- Storage Rows -->
        <div
          v-for="policy in storageEvents"
          :key="policy.key"
          class="policy-row policy-row--storage"
        >
          <div class="policy-row__icon">
            <q-avatar
              :color="getPolicyIcon(policy.key).color"
              text-color="white"
              size="32px"
              font-size="16px"
              :icon="getPolicyIcon(policy.key).icon"
            />
          </div>

          <div class="policy-row__name">
            <div class="text-body2 text-weight-medium">{{ getEventTypeName(policy.key) }}</div>
            <div class="text-caption text-grey-6">{{ getDefaultLabel(policy) }}</div>
          </div>

          <div class="policy-row__input">
            <q-input
              v-model.number="policy.currentDays"
              outlined
              dense
              hide-bottom-space
              type="number"
              :suffix="persistence.form.daysLabel.value"
              :rules="[() => validatePolicy(policy)]"
              class="compact-input"
            />
          </div>

          <div class="policy-row__limits text-caption text-grey-6">
            {{ getDaysHint(policy) }}
          </div>

          <div class="policy-row__action">
            <q-icon
              name="info_outline"
              color="primary"
              size="20px"
              class="cursor-pointer"
              @click="showEventInfo(policy.key)"
            >
              <AppTooltip :content="persistence.tooltips.moreInfo.value" />
            </q-icon>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="row justify-end q-gutter-sm q-mt-lg q-px-md">
        <q-btn
          unelevated
          color="grey-7"
          icon="refresh"
          class="rounded-borders"
          :label="persistence.resetButton.value"
          :disable="!hasChanges"
          @click="confirmReset"
        />
        <q-btn
          type="submit"
          unelevated
          color="primary"
          icon="save"
          class="rounded-borders"
          :label="actions.saveChanges.value"
          :loading="saving"
          :disable="hasErrors || !hasChanges"
        />
      </div>
    </q-form>
  </div>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
    transform: scale(1.1);
  }
}

.persistence-settings {
  .text-caption {
    line-height: 1.4;
  }
}

// Group headers
.group-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: var(--mapex-radius-md) var(--mapex-radius-md) 0 0;

  &--pipeline {
    background-color: rgba(var(--q-warning-rgb), 0.06);
    border-left: 3px solid var(--q-warning);
  }

  &--storage {
    background-color: rgba(var(--q-info-rgb), 0.06);
    border-left: 3px solid var(--q-info);
  }
}

// Warning banners
.warning-banner {
  display: flex;
  align-items: flex-start;
  padding: 10px 16px;
  border-radius: var(--mapex-radius-xs);

  &--pipeline {
    background-color: rgba(var(--q-warning-rgb), 0.08);
    border-left: 3px solid var(--q-warning);
  }

  &--storage {
    background-color: rgba(var(--q-info-rgb), 0.08);
    border-left: 3px solid var(--q-info);
  }
}

// Policy rows — DataRow style
.policy-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-xs);
  margin-bottom: 8px;
  transition: background-color 0.15s ease;

  &:hover {
    background-color: var(--mapex-surface-highlight);
  }

  &--pipeline {
    border-left: 3px solid var(--q-warning);
  }

  &--storage {
    border-left: 3px solid var(--q-info);
  }

  &__icon {
    flex-shrink: 0;
    width: 32px;
  }

  &__name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }

  &__input {
    flex-shrink: 0;
    width: 120px;
  }

  &__limits {
    flex-shrink: 0;
    width: 140px;
    text-align: center;
  }

  &__action {
    flex-shrink: 0;
    width: 24px;
    display: flex;
    justify-content: center;
  }
}

// Compact input override
.compact-input {
  :deep(.q-field__control) {
    height: 36px;
  }

  :deep(.q-field__label) {
    font-size: 12px;
  }

  :deep(.q-field__native) {
    font-size: 13px;
    padding: 0 8px;
  }
}

// Responsive: stack on small screens
@media (max-width: 768px) {
  .policy-row {
    flex-wrap: wrap;
    gap: 8px;

    &__name {
      flex: 1 1 calc(100% - 56px);
    }

    &__input {
      width: 100%;
      order: 10;
    }

    &__limits {
      width: auto;
      text-align: left;
      order: 11;
    }

    &__action {
      position: absolute;
      right: 16px;
      top: 10px;
    }

    position: relative;
  }
}
</style>

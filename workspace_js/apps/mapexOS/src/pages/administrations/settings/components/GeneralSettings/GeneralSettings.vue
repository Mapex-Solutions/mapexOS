<script setup lang="ts">
defineOptions({
  name: 'GeneralSettings'
});

// 1. Types
import type { GeneralSettingsProps, GeneralSettingsEmits, LocalGeneralData } from './interfaces';
import type { GeneralSettingsSavePayload } from '../../settingsPage/interfaces';

// 2. Vue
import { ref, computed, watch } from 'vue';
import { date } from 'quasar';

// 3. Components
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

// 4. Composables
import { useSettingsTranslations } from '@composables/i18n';
import { useOrganizationStore } from '@stores/organization';

// 5. Utils
import { showWarning } from '@utils/modal';

const { general, overview } = useSettingsTranslations();
const orgStore = useOrganizationStore();

// Props & Emits
const props = defineProps<GeneralSettingsProps>();
const emit = defineEmits<GeneralSettingsEmits>();

// Modal state
const showAccessPolicyInfo = ref(false);

// Local editable state (deep clone from props)
// This allows us to edit without mutating props
const localData = ref<LocalGeneralData>({
  // Basic Info (editable)
  name: '',
  phone: '',
  enabled: true,

  // Address (editable)
  city: '',
  state: '',
  country: '',
  zipCode: '',

  // Access Policy (editable)
  rolePolicy: 'merge',
  defaultScope: 'recursive',

  // Readonly fields (for display only)
  type: '',
  code: '',
  pathKey: '',
  depth: 0,
  childCount: 0,
  parentOrgId: null,
  authProviderType: '',
  created: '',
  updated: '',
});

/**
 * Watch props.organizationData and update local state
 * This creates a deep clone for editing without mutating props
 */
watch(
  () => props.organizationData,
  (newData) => {
    if (!newData) return;

    // Map API response to local editable fields
    localData.value = {
      // Basic Info (editable)
      name: newData.name || '',
      phone: newData.phone || '',
      enabled: newData.enabled ?? true,

      // Address (editable)
      city: newData.address?.city || '',
      state: newData.address?.state || '',
      country: newData.address?.country || '',
      zipCode: newData.address?.zipCode || '',

      // Access Policy (editable)
      rolePolicy: (newData as any).accessPolicy?.rolePolicy || 'merge',
      defaultScope: (newData as any).accessPolicy?.defaultScope || 'recursive',

      // Readonly fields (for display)
      type: newData.type || '',
      code: (newData as any).code || '',
      pathKey: (newData as any).pathKey || '',
      depth: (newData as any).depth || 0,
      childCount: (newData as any).childCount || 0,
      parentOrgId: (newData as any).parentOrgId || null,
      authProviderType: (newData as any).authConfig?.providerType || '',
      created: (newData as any).created || '',
      updated: (newData as any).updated || '',
    };
  },
  { immediate: true, deep: true }
);

/**
 * Save organization data
 * Emits save event to parent with sanitized payload
 * Parent handles API call and refresh
 */
function saveOrganizationData() {
  // Build payload with only editable fields
  const payload: GeneralSettingsSavePayload = {
    name: localData.value.name,
    phone: localData.value.phone,
    enabled: localData.value.enabled,
    address: {
      city: localData.value.city,
      state: localData.value.state,
      country: localData.value.country,
      zipCode: localData.value.zipCode,
    },
    accessPolicy: {
      rolePolicy: localData.value.rolePolicy,
      defaultScope: localData.value.defaultScope,
    },
  };

  // Emit save event to parent (no API call here)
  emit('save', payload);
}

// Helper: Get type icon based on organization type
function getTypeIcon(type: string): string {
  const iconMap: Record<string, string> = {
    vendor: 'business',
    customer: 'domain',
    site: 'location_on',
    building: 'apartment',
    floor: 'layers',
    zone: 'place',
  };
  return iconMap[type] || 'domain';
}

// Helper: Get type color based on organization type
function getTypeColor(type: string): string {
  const colorMap: Record<string, string> = {
    vendor: 'purple-6',
    customer: 'green-6',
    site: 'orange-6',
    building: 'blue-6',
    floor: 'teal-6',
    zone: 'green-7',
  };
  return colorMap[type] || 'primary';
}

// Helper: Get auth provider icon
function getAuthProviderIcon(provider: string): string {
  const iconMap: Record<string, string> = {
    internal: 'vpn_key',
    google: 'login',
    github: 'code',
    microsoft: 'business',
    keycloak: 'shield',
  };
  return iconMap[provider] || 'key';
}

// Helper: Format date timestamp
function formatDate(dateString: string): string {
  if (!dateString) return 'N/A';
  return date.formatDate(dateString, 'DD/MM/YYYY HH:mm');
}

// Computed: Get parent organization name from store
const parentOrgName = computed(() => {
  if (!localData.value.parentOrgId) return null;
  const parent = orgStore.flatList.find(org => org.id === localData.value.parentOrgId);
  return parent?.name || 'Unknown';
});

// Computed: Role Policy options
const rolePolicyOptions = computed(() => [
  { label: general.accessPolicyModal.rolePolicyMerge.value, value: 'merge', icon: 'merge' },
  { label: general.accessPolicyModal.rolePolicyStrict.value, value: 'strict', icon: 'lock' },
]);

// Computed: Default Scope options
const defaultScopeOptions = computed(() => [
  { label: general.accessPolicyModal.defaultScopeRecursive.value, value: 'recursive', icon: 'sync_alt' },
  { label: general.accessPolicyModal.defaultScopeLocal.value, value: 'local', icon: 'account_tree' },
]);

// Handle status toggle change with confirmation
async function handleStatusChange(newValue: boolean) {
  // If disabling the organization, show warning modal
  if (!newValue) {
    const confirmed = await showWarning(
      general.disableModal.title.value,
      general.disableModal.message.value,
      general.disableModal.confirmButton.value,
      general.disableModal.cancelButton.value
    );

    if (!confirmed) {
      // User cancelled, revert the value
      return;
    }
  }

  // User confirmed or is enabling, update the value
  localData.value.enabled = newValue;
}
</script>

<template>
  <div class="general-settings">
    <!-- Page Header -->
    <div class="row items-center q-mt-md q-mb-lg">
      <q-icon name="business_center" size="sm" color="primary" class="q-mr-sm"/>
      <div class="text-subtitle1 text-weight-medium text-primary">{{ general.title.value }}</div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="props.loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Form with 2-Column Layout -->
    <q-form v-else @submit="saveOrganizationData">
      <div class="row q-col-gutter-lg">

        <!-- LEFT COLUMN: Editable Cards -->
        <div class="col-12 col-md-7">
          <div class="column q-gutter-md">

            <!-- CARD 1: Basic Information -->
            <q-card flat bordered class="rounded-borders section-card basic-card">
              <q-card-section>
                <div class="row items-center q-mb-md">
                  <q-icon name="info" size="sm" color="primary" class="q-mr-sm"/>
                  <div class="text-subtitle2 text-weight-medium text-primary">{{ general.sections.basicInfo.value }}</div>
                </div>

                <div class="row q-col-gutter-md">
                  <!-- Enabled Toggle -->
                  <div class="col-12">
                    <div class="row items-center q-gutter-sm">
                      <q-icon name="toggle_on" color="primary" size="sm" />
                      <span class="text-body2">{{ general.fields.status.value }}</span>
                      <q-toggle
                        :model-value="localData.enabled"
                        :label="localData.enabled ? general.fields.enabled.value : general.fields.disabled.value"
                        :color="localData.enabled ? 'green-6' : 'grey-5'"
                        checked-icon="check"
                        unchecked-icon="clear"
                        @update:model-value="handleStatusChange"
                      />
                    </div>
                  </div>

                  <!-- Name Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.name"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.organizationName.value"
                      :rules="[val => !!val || general.validation.nameRequired.value]"
                    >
                      <template v-slot:prepend>
                        <q-icon name="business" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <!-- Phone Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.phone"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.phone.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="phone" color="primary" />
                      </template>
                    </q-input>
                  </div>
                </div>
              </q-card-section>
            </q-card>

            <!-- CARD 2: Address Information -->
            <q-card flat bordered class="rounded-borders section-card address-card">
              <q-card-section>
                <div class="row items-center q-mb-md">
                  <q-icon name="location_on" size="sm" color="primary" class="q-mr-sm"/>
                  <div class="text-subtitle2 text-weight-medium text-primary">{{ general.sections.address.value }}</div>
                </div>

                <div class="row q-col-gutter-md">
                  <!-- City Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.city"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.city.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="location_city" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <!-- State Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.state"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.state.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="map" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <!-- Country Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.country"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.country.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="public" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <!-- Zip Code Input -->
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="localData.zipCode"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="general.fields.zipCode.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="markunread_mailbox" color="primary" />
                      </template>
                    </q-input>
                  </div>
                </div>
              </q-card-section>
            </q-card>

            <!-- CARD 3: Access Policy -->
            <q-card flat bordered class="rounded-borders section-card access-card">
              <q-card-section>
                <div class="row items-center q-mb-md">
                  <q-icon name="security" size="sm" color="primary" class="q-mr-sm"/>
                  <div class="text-subtitle2 text-weight-medium text-primary">{{ general.sections.accessPolicy.value }}</div>
                  <q-btn
                    flat
                    round
                    dense
                    icon="info"
                    size="sm"
                    color="primary"
                    class="q-ml-xs"
                    @click="showAccessPolicyInfo = true"
                  >
                    <AppTooltip :content="general.messages.infoTooltip.value" />
                  </q-btn>
                </div>

                <div class="row q-col-gutter-md">
                  <!-- Auth Provider (Readonly) -->
                  <div class="col-12">
                    <div class="readonly-field">
                      <q-icon name="vpn_key" size="sm" color="grey-7" class="q-mr-xs" />
                      <span class="field-label">{{ general.fields.authProvider.value }}:</span>
                      <DetailChip
                        :value="localData.authProviderType || 'N/A'"
                        :icon="getAuthProviderIcon(localData.authProviderType)"
                        color="deep-purple"
                        size="sm"
                        dense
                        class="q-ml-sm"
                      />
                    </div>
                  </div>

                  <!-- Role Policy Select -->
                  <div class="col-12 col-md-6">
                    <q-select
                      v-model="localData.rolePolicy"
                      outlined
                      dense
                      emit-value
                      map-options
                      class="rounded-borders"
                      :options="rolePolicyOptions"
                      :label="general.fields.rolePolicy.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="admin_panel_settings" color="primary" />
                      </template>
                      <template v-slot:option="scope">
                        <q-item v-bind="scope.itemProps">
                          <q-item-section avatar>
                            <q-icon :name="scope.opt.icon" />
                          </q-item-section>
                          <q-item-section>
                            <q-item-label>{{ scope.opt.label }}</q-item-label>
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>
                  </div>

                  <!-- Default Scope Select -->
                  <div class="col-12 col-md-6">
                    <q-select
                      v-model="localData.defaultScope"
                      outlined
                      dense
                      emit-value
                      map-options
                      class="rounded-borders"
                      :options="defaultScopeOptions"
                      :label="general.fields.defaultScope.value"
                    >
                      <template v-slot:prepend>
                        <q-icon name="account_tree" color="primary" />
                      </template>
                      <template v-slot:option="scope">
                        <q-item v-bind="scope.itemProps">
                          <q-item-section avatar>
                            <q-icon :name="scope.opt.icon" />
                          </q-item-section>
                          <q-item-section>
                            <q-item-label>{{ scope.opt.label }}</q-item-label>
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>
                  </div>
                </div>
              </q-card-section>
            </q-card>

          </div>
        </div>

        <!-- RIGHT COLUMN: Organization Overview (Readonly) -->
        <div class="col-12 col-md-5">
          <q-card flat bordered class="rounded-borders section-card overview-card">
            <q-card-section>
              <!-- Card Header -->
              <div class="row items-center q-mb-lg">
                <q-icon name="dashboard" size="sm" color="blue-6" class="q-mr-sm"/>
                <div class="text-subtitle2 text-weight-medium" style="color: var(--q-info);">{{ overview.organizationOverview.value }}</div>
              </div>

              <!-- Hierarchy Section -->
              <div class="overview-section q-mb-md">
                <div class="row items-center q-mb-sm">
                  <q-icon name="account_tree" size="xs" color="grey-7" class="q-mr-xs"/>
                  <span class="text-caption text-weight-bold text-grey-7 text-uppercase">{{ general.sections.hierarchy.value }}</span>
                </div>
                <q-separator class="q-mb-sm" />

                <div class="overview-grid">
                  <!-- Parent Organization -->
                  <div class="info-item">
                    <div class="info-label">{{ general.fields.parentOrganization.value }}</div>
                    <div class="info-value">
                      <DetailChip
                        v-if="parentOrgName"
                        :value="parentOrgName"
                        icon="domain"
                        color="indigo"
                        size="sm"
                        dense
                      />
                      <DetailChip
                        v-else
                        :value="general.fields.rootOrganization.value"
                        icon="check_circle"
                        color="green"
                        size="sm"
                        dense
                      />
                    </div>
                  </div>

                  <!-- Type -->
                  <div class="info-item">
                    <div class="info-label">{{ general.fields.type.value }}</div>
                    <div class="info-value">
                      <DetailChip
                        :value="localData.type || 'N/A'"
                        :icon="getTypeIcon(localData.type)"
                        :color="getTypeColor(localData.type) as any"
                        size="sm"
                        dense
                      />
                    </div>
                  </div>

                  <!-- Metrics Row -->
                  <div class="row q-col-gutter-sm q-mt-xs">
                    <div class="col-4">
                      <div class="metric-item">
                        <div class="metric-value">
                          <q-badge color="orange-6" class="metric-badge">{{ localData.childCount }}</q-badge>
                        </div>
                        <div class="metric-label">{{ general.fields.childOrganizations.value }}</div>
                      </div>
                    </div>
                    <div class="col-4">
                      <div class="metric-item">
                        <div class="metric-value">
                          <q-badge color="teal-6" class="metric-badge">{{ localData.depth }}</q-badge>
                        </div>
                        <div class="metric-label">{{ general.fields.depth.value }}</div>
                      </div>
                    </div>
                    <div class="col-4">
                      <div class="metric-item">
                        <div class="metric-value">
                          <q-badge color="blue-6" class="metric-badge">{{ localData.code || 'N/A' }}</q-badge>
                        </div>
                        <div class="metric-label">{{ general.fields.code.value }}</div>
                      </div>
                    </div>
                  </div>

                  <!-- Path Key -->
                  <div class="info-item q-mt-sm">
                    <div class="info-label">{{ general.fields.path.value }}</div>
                    <div class="info-value">
                      <q-badge color="indigo-6">{{ localData.pathKey || 'N/A' }}</q-badge>
                    </div>
                  </div>
                </div>
              </div>

              <!-- System Section -->
              <div class="overview-section">
                <div class="row items-center q-mb-sm">
                  <q-icon name="info_outline" size="xs" color="grey-7" class="q-mr-xs"/>
                  <span class="text-caption text-weight-bold text-grey-7 text-uppercase">{{ general.sections.system.value }}</span>
                </div>
                <q-separator class="q-mb-sm" />

                <div class="overview-grid">
                  <!-- Created -->
                  <div class="info-item">
                    <div class="info-label">
                      <q-icon name="event" size="xs" color="grey-6" class="q-mr-xs" />
                      {{ general.fields.created.value }}
                    </div>
                    <div class="info-value text-body2">{{ formatDate(localData.created) }}</div>
                  </div>

                  <!-- Updated -->
                  <div class="info-item">
                    <div class="info-label">
                      <q-icon name="update" size="xs" color="grey-6" class="q-mr-xs" />
                      {{ general.fields.lastUpdated.value }}
                    </div>
                    <div class="info-value text-body2">{{ formatDate(localData.updated) }}</div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

      </div>

      <!-- Save Button -->
      <div class="row justify-center q-mt-lg">
        <q-btn
          type="submit"
          unelevated
          icon="save"
          color="primary"
          size="md"
          class="rounded-borders"
          :label="general.buttons.saveChanges.value"
          :loading="props.loading"
        />
      </div>
    </q-form>

    <!-- Access Policy Information Modal -->
    <q-dialog v-model="showAccessPolicyInfo" class="access-policy-info-dialog">
      <q-card style="min-width: 500px; max-width: 700px;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-primary">
            <q-icon name="security" class="q-mr-sm" />
            {{ general.accessPolicyModal.title.value }}
          </div>
          <q-space />
          <q-btn v-close-popup flat round dense icon="close" />
        </q-card-section>

        <q-separator class="q-mt-sm" />

        <q-card-section class="q-pt-md">
          <!-- Role Policy Section -->
          <div class="info-section q-mb-lg">
            <div class="row items-center q-mb-sm">
              <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
              <div class="text-subtitle1 text-weight-medium">{{ general.accessPolicyModal.rolePolicyTitle.value }}</div>
            </div>
            <div class="text-body2 text-grey-8 q-mb-md">
              {{ general.accessPolicyModal.rolePolicyDescription.value }}
            </div>

            <!-- Merge Option -->
            <div class="policy-option q-mb-md">
              <div class="row items-center q-mb-xs">
                <q-icon name="merge" color="green-6" size="sm" class="q-mr-sm" />
                <span class="text-weight-medium text-body2">{{ general.accessPolicyModal.rolePolicyMerge.value }}</span>
              </div>
              <div class="text-body2 text-grey-7 q-pl-lg">
                {{ general.accessPolicyModal.rolePolicyMergeDescription.value }}
              </div>
            </div>

            <!-- Strict Option -->
            <div class="policy-option">
              <div class="row items-center q-mb-xs">
                <q-icon name="lock" color="orange-6" size="sm" class="q-mr-sm" />
                <span class="text-weight-medium text-body2">{{ general.accessPolicyModal.rolePolicyStrict.value }}</span>
              </div>
              <div class="text-body2 text-grey-7 q-pl-lg">
                {{ general.accessPolicyModal.rolePolicyStrictDescription.value }}
              </div>
            </div>
          </div>

          <q-separator class="q-my-md" />

          <!-- Default Scope Section -->
          <div class="info-section">
            <div class="row items-center q-mb-sm">
              <q-icon name="account_tree" color="primary" size="sm" class="q-mr-sm" />
              <div class="text-subtitle1 text-weight-medium">{{ general.accessPolicyModal.defaultScopeTitle.value }}</div>
            </div>
            <div class="text-body2 text-grey-8 q-mb-md">
              {{ general.accessPolicyModal.defaultScopeDescription.value }}
            </div>

            <!-- Recursive Option -->
            <div class="policy-option q-mb-md">
              <div class="row items-center q-mb-xs">
                <q-icon name="sync_alt" color="blue-6" size="sm" class="q-mr-sm" />
                <span class="text-weight-medium text-body2">{{ general.accessPolicyModal.defaultScopeRecursive.value }}</span>
              </div>
              <div class="text-body2 text-grey-7 q-pl-lg">
                {{ general.accessPolicyModal.defaultScopeRecursiveDescription.value }}
              </div>
            </div>

            <!-- Local Option -->
            <div class="policy-option">
              <div class="row items-center q-mb-xs">
                <q-icon name="place" color="purple-6" size="sm" class="q-mr-sm" />
                <span class="text-weight-medium text-body2">{{ general.accessPolicyModal.defaultScopeLocal.value }}</span>
              </div>
              <div class="text-body2 text-grey-7 q-pl-lg">
                {{ general.accessPolicyModal.defaultScopeLocalDescription.value }}
              </div>
            </div>
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right">
          <q-btn v-close-popup flat color="primary" :label="general.accessPolicyModal.closeButton.value" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.general-settings {
  .text-caption {
    line-height: 1.4;
  }

  // Section cards with subtle shadow and hover effect
  .section-card {
    transition: var(--mapex-transition-slow);
    box-shadow: var(--mapex-shadow-xs);

    &:hover {
      box-shadow: var(--mapex-shadow-sm);
    }
  }

  // Colored left borders per card type
  .basic-card {
    border-left: 3px solid var(--q-positive);
  }

  .address-card {
    border-left: 3px solid var(--q-warning);
  }

  .access-card {
    border-left: 3px solid var(--q-accent);
  }

  .overview-card {
    border-left: 3px solid var(--q-info);
    position: sticky;
    top: 16px;

    .overview-section {
      .overview-grid {
        display: flex;
        flex-direction: column;
        gap: 10px;
      }
    }

    .info-item {
      .info-label {
        font-size: 0.75rem;
        color: var(--mapex-text-muted);
        font-weight: 500;
        text-transform: uppercase;
        letter-spacing: 0.025em;
        margin-bottom: 4px;
        display: flex;
        align-items: center;
      }

      .info-value {
        font-size: 0.875rem;
        color: var(--mapex-text-primary);
      }
    }

    .metric-item {
      text-align: center;
      padding: 8px 4px;
      background: var(--mapex-submenu-bg);
      border-radius: var(--mapex-radius-md);

      .metric-badge {
        font-size: 0.8rem;
        padding: 4px 10px;
      }

      .metric-label {
        font-size: 0.7rem;
        color: var(--mapex-text-muted);
        font-weight: 500;
        text-transform: uppercase;
        margin-top: 6px;
        letter-spacing: 0.025em;
      }
    }
  }

  // Readonly field styling
  .readonly-field {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    background: var(--mapex-submenu-bg);
    border-radius: var(--mapex-radius-md);
    min-height: 40px;

    .field-label {
      font-size: 0.875rem;
      color: var(--mapex-text-secondary);
      font-weight: 500;
    }

    .field-value {
      font-size: 0.875rem;
      color: var(--mapex-text-primary);
      font-weight: 400;
    }
  }
}

// Access Policy Info Dialog Styling
.access-policy-info-dialog {
  .info-section {
    .policy-option {
      padding: 12px;
      background: var(--mapex-submenu-bg);
      border-radius: var(--mapex-radius-md);
      border-left: 3px solid transparent;
      transition: var(--mapex-transition-base);

      &:hover {
        background: var(--mapex-submenu-bg);
        border-left-color: var(--q-primary);
      }
    }
  }
}
</style>

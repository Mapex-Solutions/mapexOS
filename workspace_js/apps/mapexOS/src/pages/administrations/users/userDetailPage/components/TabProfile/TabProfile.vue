<script setup lang="ts">
defineOptions({
  name: 'TabProfile'
});

/** TYPE IMPORTS */
import type { TabProfileProps, UserDetailData } from '../../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useUserDetailTranslations } from '@composables/i18n';

/** PROPS */
const props = defineProps<TabProfileProps>();

/** COMPOSABLES & STORES */
const t = useUserDetailTranslations();

/** COMPUTED */

/**
 * User data from props (passed from parent)
 */
const user = computed<UserDetailData | null>(() => props.user);

/**
 * Loading state from props
 */
const loading = computed(() => props.loading);

/** FUNCTIONS */

/**
 * Get user's full name
 *
 * @returns {string} Full name or email
 */
function getUserFullName(): string {
  if (!user.value) return '';
  const firstName = user.value.firstName || '';
  const lastName = user.value.lastName || '';
  return `${firstName} ${lastName}`.trim() || user.value.email || 'Unknown';
}

/**
 * Get user initials for avatar
 *
 * @returns {string} Initials (up to 2 characters)
 */
function getUserInitials(): string {
  if (!user.value) return '?';
  const firstName = user.value.firstName || '';
  const lastName = user.value.lastName || '';
  if (firstName.length > 0 && lastName.length > 0) {
    return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase();
  }
  const email = user.value.email || '';
  return email.length > 0 ? email.charAt(0).toUpperCase() : '?';
}

/**
 * Get auth provider display label
 *
 * @returns {string} Auth provider label
 */
function getAuthProviderLabel(): string {
  const type = user.value?.authProvider?.type || 'internal';
  const authProviders = t.profile.authProviders as Record<string, { value: string }>;
  return authProviders[type]?.value.toUpperCase() || type.toUpperCase();
}

/**
 * Get auth provider color
 *
 * @returns {string} Color name
 */
function getAuthProviderColor(): 'grey' | 'red' | 'blue' | 'orange' {
  const type = user.value?.authProvider?.type;
  const colorMap: Record<string, 'grey' | 'red' | 'blue' | 'orange'> = {
    internal: 'grey',
    google: 'red',
    github: 'grey',
    microsoft: 'blue',
    keycloak: 'orange',
  };
  return colorMap[type || 'internal'] || 'grey';
}

/**
 * Get auth provider icon
 *
 * @returns {string} Icon name
 */
function getAuthProviderIcon(): string {
  const type = user.value?.authProvider?.type;
  const iconMap: Record<string, string> = {
    internal: 'lock',
    google: 'g_translate',
    github: 'code',
    microsoft: 'business',
    keycloak: 'vpn_key',
  };
  return iconMap[type || 'internal'] || 'lock';
}

/**
 * Format date using Quasar date utils
 *
 * @param {any} dateValue - Date to format
 * @returns {string} Formatted date
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

</script>

<template>
  <div class="tab-profile">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-lg">
      <q-spinner color="primary" size="40px" />
    </div>

    <!-- User Profile Content -->
    <div v-else-if="user">
      <!-- Basic Info Section -->
      <div class="section q-mb-lg">
        <div class="section-header q-mb-md">
          <q-icon name="person" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.profile.sections.basicInfo.value }}</span>
        </div>

        <div class="row q-col-gutter-md">
          <!-- Avatar & Name Column -->
          <div class="col-12 col-md-4">
            <div class="avatar-container text-center q-pa-md">
              <q-avatar size="100px" color="primary" text-color="white" class="q-mb-sm shadow-2">
                <img v-if="user.avatar" :src="user.avatar" />
                <span v-else class="text-h4">{{ getUserInitials() }}</span>
              </q-avatar>
              <div class="text-subtitle1 text-weight-bold">{{ getUserFullName() }}</div>
              <div v-if="user.jobTitle" class="text-caption text-grey-7">{{ user.jobTitle }}</div>
            </div>
          </div>

          <!-- Info Fields Column -->
          <div class="col-12 col-md-8">
            <div class="row q-col-gutter-md">
              <!-- Email -->
              <div class="col-12 col-md-6">
                <div class="info-item">
                  <div class="info-label text-grey-7">{{ t.profile.fields.email.value }}</div>
                  <div class="info-value text-weight-medium">{{ user.email || '—' }}</div>
                </div>
              </div>

              <!-- Phone -->
              <div class="col-12 col-md-6">
                <div class="info-item">
                  <div class="info-label text-grey-7">{{ t.profile.fields.phone.value }}</div>
                  <div class="info-value">{{ user.phone || '—' }}</div>
                </div>
              </div>

              <!-- Status -->
              <div class="col-12 col-md-6">
                <div class="info-item">
                  <div class="info-label text-grey-7">{{ t.profile.fields.status.value }}</div>
                  <div class="info-value">
                    <DetailChip
                      :label="user.enabled ? t.profile.status.active.value : t.profile.status.inactive.value"
                      :color="user.enabled ? 'green' : 'red'"
                      size="sm"
                    />
                  </div>
                </div>
              </div>

              <!-- Organization -->
              <div class="col-12 col-md-6">
                <div class="info-item">
                  <div class="info-label text-grey-7">{{ t.profile.fields.organization.value }}</div>
                  <div class="info-value">{{ user.organizationName || '—' }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Authentication Section -->
      <div class="section q-mb-lg">
        <div class="section-header q-mb-md">
          <q-icon name="security" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.profile.sections.authentication.value }}</span>
        </div>

        <div class="row q-col-gutter-md">
          <!-- Auth Provider -->
          <div class="col-12 col-md-6">
            <div class="info-item">
              <div class="info-label text-grey-7">{{ t.profile.fields.authProvider.value }}</div>
              <div class="info-value">
                <DetailChip
                  :icon="getAuthProviderIcon()"
                  :color="getAuthProviderColor()"
                  size="sm"
                  :label="getAuthProviderLabel()"
                />
              </div>
            </div>
          </div>

          <!-- Change Password -->
          <div class="col-12 col-md-6">
            <div class="info-item">
              <div class="info-label text-grey-7">{{ t.profile.fields.changePasswordNextLogin.value }}</div>
              <div class="info-value">
                <DetailChip
                  :icon="user.changePasswordNextLogin ? 'lock_reset' : 'lock'"
                  :color="user.changePasswordNextLogin ? 'orange' : 'grey'"
                  size="sm"
                  :label="user.changePasswordNextLogin ? t.profile.values.yes.value : t.profile.values.no.value"
                />
              </div>
            </div>
          </div>

          <!-- External ID -->
          <div v-if="user.authProvider?.externalId" class="col-12">
            <div class="info-item">
              <div class="info-label text-grey-7">{{ t.profile.fields.externalId.value }}</div>
              <code class="external-id">{{ user.authProvider.externalId }}</code>
            </div>
          </div>
        </div>
      </div>

      <!-- Timestamps Section -->
      <div class="section">
        <div class="section-header q-mb-md">
          <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.profile.sections.timestamps.value }}</span>
        </div>

        <div class="row q-col-gutter-md">
          <!-- Created -->
          <div class="col-12 col-md-6">
            <div class="info-item">
              <div class="info-label text-grey-7">{{ t.profile.fields.created.value }}</div>
              <div class="info-value">{{ formatDate(user.created) }}</div>
            </div>
          </div>

          <!-- Updated -->
          <div class="col-12 col-md-6">
            <div class="info-item">
              <div class="info-label text-grey-7">{{ t.profile.fields.updated.value }}</div>
              <div class="info-value">{{ formatDate(user.updated) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- No Data -->
    <div v-else class="section text-center q-pa-xl text-grey-6">
      <q-icon name="warning" size="48px" class="q-mb-md" />
      <div>{{ t.profile.error.value }}</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.tab-profile {
  .section {
    background: var(--mapex-surface-bg);
    border-radius: var(--mapex-radius-md);
    padding: 16px;
  }

  .section-header {
    display: flex;
    align-items: center;
    border-bottom: 1px solid var(--mapex-card-border);
    padding-bottom: 8px;
  }

  .avatar-container {
    background: white;
    border-radius: var(--mapex-radius-md);
    border: 1px solid var(--mapex-card-border);
  }

  .info-item {
    .info-label {
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      margin-bottom: 4px;
    }

    .info-value {
      font-size: 0.95rem;
    }
  }

  .external-id {
    display: inline-block;
    background: var(--mapex-card-border);
    padding: 4px 8px;
    border-radius: var(--mapex-radius-xs);
    font-size: 0.85rem;
  }
}
</style>

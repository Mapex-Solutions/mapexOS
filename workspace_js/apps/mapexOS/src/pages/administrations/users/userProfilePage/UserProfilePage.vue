<script setup lang="ts">
defineOptions({
  name: 'UserProfilePage'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { ReviewSectionDef } from '@components/forms/review/interfaces';
import type {
  UserProfileData,
  PasswordData,
  ProfileSection,
  UserGroupInfo,
  UserMembershipInfo,
} from './interfaces';

/** VUE IMPORTS */
import { ref, reactive, computed, onMounted } from 'vue';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useUserProfileTranslations } from '@composables/i18n';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert/notify';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import {
  DEFAULT_ACTIVE_SECTION,
  SECTION,
  PASSWORD_MIN_LENGTH,
  NOTIFICATION_TIMEOUT,
  PROFILE_SECTIONS_CONFIG,
} from './constants';

/** COMPOSABLES & STORES */
const t = useUserProfileTranslations();

/** STATE */
const currentSection = ref(DEFAULT_ACTIVE_SECTION);
const personalForm = ref(null);
const passwordForm = ref(null);
const loading = ref(true);
const saving = ref(false);
const errorMessage = ref('');

const userData = reactive<UserProfileData>({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  jobTitle: '',
});

const passwordData = reactive<PasswordData>({
  current: '',
  new: '',
  confirm: '',
});

const groups = ref<UserGroupInfo[]>([]);
const memberships = ref<UserMembershipInfo[]>([]);
const groupsCount = ref(0);

/** COMPUTED */

/**
 * Profile sections with reactive translations for StepperVertical
 */
const profileSections = computed((): ProfileSection[] => PROFILE_SECTIONS_CONFIG(t));

/**
 * Current section header info
 */
const currentSectionHeader = computed(() => {
  const section = profileSections.value[currentSection.value - 1];
  return section || { title: '', description: '', icon: 'settings' };
});

/**
 * Preview data for FormReview component
 */
const previewData = computed((): ReviewSectionDef[] => [
  {
    stepNumber: SECTION.PERSONAL,
    label: t.review.sections.personal.value,
    icon: { name: 'person', color: 'primary' },
    fields: [
      {
        label: t.review.fields.firstName.value,
        value: userData.firstName || '-',
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.lastName.value,
        value: userData.lastName || '-',
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.email.value,
        value: userData.email || '-',
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.phone.value,
        value: userData.phone || '-',
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.jobTitle.value,
        value: userData.jobTitle || '-',
        type: 'text',
        colSize: 6,
      },
    ],
  },
  {
    stepNumber: SECTION.PASSWORD,
    label: t.review.sections.password.value,
    icon: { name: 'lock', color: 'primary' },
    fields: [
      {
        label: t.review.fields.passwordChanged.value,
        value: passwordData.new ? t.review.fields.yes.value : t.review.fields.no.value,
        type: 'badge',
        badgeColors: {
          [t.review.fields.yes.value]: 'positive',
          [t.review.fields.no.value]: 'grey',
        },
        colSize: 12,
      },
    ],
  },
  {
    stepNumber: SECTION.GROUPS_ACCESS,
    label: t.review.sections.groupsAccess.value,
    icon: { name: 'shield', color: 'primary' },
    fields: [
      {
        label: t.review.fields.groupsCount.value,
        value: groups.value.length.toString(),
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.membershipsCount.value,
        value: memberships.value.length.toString(),
        type: 'text',
        colSize: 6,
      },
    ],
  },
]);

/** FUNCTIONS */

/**
 * Fetch current user profile from API
 *
 * @returns {Promise<void>}
 */
async function fetchUserProfile(): Promise<void> {
  loading.value = true;
  errorMessage.value = '';

  try {
    const user = await apis.mapexOS.users.me();

    if (user) {
      userData.firstName = user.firstName || '';
      userData.lastName = user.lastName || '';
      userData.email = user.email || '';
      userData.phone = user.phone || '';
      userData.jobTitle = user.jobTitle || '';
      groups.value = user.groups || [];
      memberships.value = user.memberships || [];
      groupsCount.value = user.groupsCount ?? groups.value.length;
    }
  } catch {
    errorMessage.value = t.messages.errorLoading.value;
    notifyFail({ message: t.messages.errorLoading.value, timeout: NOTIFICATION_TIMEOUT });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle section change from StepperVertical
 *
 * @param {number} sectionNumber - Section number (1-based)
 */
function handleSectionChange(sectionNumber: number): void {
  currentSection.value = sectionNumber;
}

/**
 * Update personal information via API
 *
 * @returns {Promise<void>}
 */
async function updatePersonalInfo(): Promise<void> {
  if (!personalForm.value) return;

  saving.value = true;
  try {
    await apis.mapexOS.users.updateMe({
      firstName: userData.firstName,
      lastName: userData.lastName,
      phone: userData.phone || undefined,
      jobTitle: userData.jobTitle || undefined,
    });
    notifySuccess({ message: t.messages.personalInfoUpdated.value, timeout: NOTIFICATION_TIMEOUT });
  } catch {
    notifyFail({ message: t.messages.errorLoading.value, timeout: NOTIFICATION_TIMEOUT });
  } finally {
    saving.value = false;
  }
}

/**
 * Update user password via API
 *
 * @returns {Promise<void>}
 */
async function updatePassword(): Promise<void> {
  if (!passwordForm.value) return;

  saving.value = true;
  try {
    await apis.mapexOS.users.updateMe({
      password: passwordData.new,
    });
    notifySuccess({ message: t.messages.passwordUpdated.value, timeout: NOTIFICATION_TIMEOUT });
    passwordData.current = '';
    passwordData.new = '';
    passwordData.confirm = '';
  } catch {
    notifyFail({ message: t.messages.errorLoading.value, timeout: NOTIFICATION_TIMEOUT });
  } finally {
    saving.value = false;
  }
}

/**
 * Save profile from review step
 *
 * @returns {Promise<void>}
 */
async function saveProfile(): Promise<void> {
  saving.value = true;
  try {
    const payload: Record<string, any> = {
      firstName: userData.firstName,
      lastName: userData.lastName,
      phone: userData.phone || undefined,
      jobTitle: userData.jobTitle || undefined,
    };

    if (passwordData.new) {
      payload.password = passwordData.new;
    }

    await apis.mapexOS.users.updateMe(payload);
    notifySuccess({ message: t.messages.profileSaved.value, timeout: NOTIFICATION_TIMEOUT });

    if (passwordData.new) {
      passwordData.current = '';
      passwordData.new = '';
      passwordData.confirm = '';
    }
  } catch {
    notifyFail({ message: t.messages.errorLoading.value, timeout: NOTIFICATION_TIMEOUT });
  } finally {
    saving.value = false;
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void fetchUserProfile();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Header Section -->
    <PageHeader
      icon="account_circle"
      iconColor="primary"
      :title="t.page.title.value"
      :description="t.page.description.value"
      :button="{ label: t.page.backButton.value, icon: 'arrow_back', flat: true, to: '/' }"
    />

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-xl" style="min-height: 400px">
      <q-spinner color="primary" size="60px" />
    </div>

    <!-- Error State -->
    <q-banner v-else-if="errorMessage" class="bg-negative text-white q-mb-lg rounded-borders">
      <template #avatar>
        <q-icon name="error" color="white" />
      </template>
      {{ errorMessage }}
      <template #action>
        <q-btn flat color="white" :label="t.buttons.save.value" @click="fetchUserProfile" />
      </template>
    </q-banner>

    <!-- Content -->
    <div v-else class="row q-col-gutter-lg">
      <!-- Navigation Stepper -->
      <div class="col-12 col-md-4">
        <StepperVertical
          :title="t.sections.profileSettings.value"
          :subtitle="t.sectionDescriptions.profileSettings.value"
          header-icon="settings"
          :info-text="t.stepper.infoText.value"
          :current-step-label="t.stepper.currentSection.value"
          :current-step="currentSection"
          :steps="profileSections"
          :allow-step-navigation="true"
          @step-click="handleSectionChange"
        />
      </div>

      <!-- Content Card -->
      <div class="col-12 col-md-8">
        <q-card class="rounded-borders shadow-2">
          <!-- Card Header -->
          <q-card-section class="bg-grey-1 q-pb-md">
            <div class="text-h6 text-weight-bold text-primary">
              <q-icon
                size="sm"
                class="q-mr-xs"
                :name="currentSectionHeader.icon"
                color="primary"
              />
              {{ currentSectionHeader.title }}
            </div>
            <div class="text-caption text-grey-7">{{ currentSectionHeader.description }}</div>
          </q-card-section>

          <!-- Personal Information -->
          <div v-if="currentSection === SECTION.PERSONAL">
            <q-form ref="personalForm" @submit="updatePersonalInfo">
              <q-card-section class="q-pa-lg">
                <div class="row q-col-gutter-md">
                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="userData.firstName"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="`${t.fields.firstName.value} *`"
                      :rules="[val => !!val || t.validation.firstNameRequired.value]"
                    >
                      <template #prepend>
                        <q-icon name="person" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="userData.lastName"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="`${t.fields.lastName.value} *`"
                      :rules="[val => !!val || t.validation.lastNameRequired.value]"
                    >
                      <template #prepend>
                        <q-icon name="person" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12">
                    <q-input
                      v-model="userData.email"
                      outlined
                      dense
                      disable
                      type="email"
                      class="rounded-borders"
                      :label="t.fields.email.value"
                    >
                      <template #prepend>
                        <q-icon name="email" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="userData.phone"
                      outlined
                      dense
                      type="tel"
                      class="rounded-borders"
                      :label="t.fields.phone.value"
                    >
                      <template #prepend>
                        <q-icon name="phone" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-md-6">
                    <q-input
                      v-model="userData.jobTitle"
                      outlined
                      dense
                      class="rounded-borders"
                      :label="t.fields.jobTitle.value"
                    >
                      <template #prepend>
                        <q-icon name="work" color="primary" />
                      </template>
                    </q-input>
                  </div>
                </div>
              </q-card-section>

              <q-separator />

              <q-card-actions align="right" class="q-px-md q-pb-md q-pt-sm">
                <q-btn
                  unelevated
                  type="submit"
                  icon-right="save"
                  class="rounded-borders"
                  color="primary"
                  :label="t.buttons.saveChanges.value"
                  :loading="saving"
                  :ripple="false"
                />
              </q-card-actions>
            </q-form>
          </div>

          <!-- Password -->
          <div v-else-if="currentSection === SECTION.PASSWORD">
            <q-form ref="passwordForm" @submit="updatePassword">
              <q-card-section class="q-pa-lg">
                <div class="row q-col-gutter-md">
                  <div class="col-12">
                    <q-input
                      v-model="passwordData.current"
                      outlined
                      dense
                      type="password"
                      class="rounded-borders"
                      :label="`${t.fields.currentPassword.value} *`"
                      :rules="[val => !!val || t.validation.currentPasswordRequired.value]"
                    >
                      <template #prepend>
                        <q-icon name="lock" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12">
                    <q-input
                      v-model="passwordData.new"
                      outlined
                      dense
                      type="password"
                      class="rounded-borders"
                      :label="`${t.fields.newPassword.value} *`"
                      :rules="[
                        val => !!val || t.validation.newPasswordRequired.value,
                        val => val.length >= PASSWORD_MIN_LENGTH || t.validation.passwordMinLength.value
                      ]"
                    >
                      <template #prepend>
                        <q-icon name="lock" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12">
                    <q-input
                      v-model="passwordData.confirm"
                      outlined
                      dense
                      type="password"
                      class="rounded-borders"
                      :label="`${t.fields.confirmPassword.value} *`"
                      :rules="[
                        val => !!val || t.validation.confirmPasswordRequired.value,
                        val => val === passwordData.new || t.validation.passwordsDoNotMatch.value
                      ]"
                    >
                      <template #prepend>
                        <q-icon name="lock" color="primary" />
                      </template>
                    </q-input>
                  </div>
                </div>
              </q-card-section>

              <q-separator />

              <q-card-actions align="right" class="q-px-md q-pb-md q-pt-sm">
                <q-btn
                  unelevated
                  type="submit"
                  icon-right="save"
                  class="rounded-borders"
                  color="primary"
                  :label="t.buttons.updatePassword.value"
                  :loading="saving"
                  :ripple="false"
                />
              </q-card-actions>
            </q-form>
          </div>

          <!-- Groups & Access -->
          <div v-else-if="currentSection === SECTION.GROUPS_ACCESS">
            <q-card-section class="q-pa-lg">
              <!-- Groups Section -->
              <div class="section q-mb-lg">
                <div class="section-header q-mb-md">
                  <q-icon name="groups" color="primary" size="sm" class="q-mr-sm" />
                  <span class="text-subtitle1 text-weight-medium">{{ t.groupsAccess.groupsTitle.value }}</span>
                  <q-badge color="primary" class="q-ml-sm">{{ groupsCount }}</q-badge>
                </div>

                <!-- Groups Empty State -->
                <div v-if="groups.length === 0" class="text-center q-pa-lg">
                  <q-icon name="groups" size="48px" color="grey-4" class="q-mb-sm" />
                  <div class="text-body2 text-grey-6">{{ t.groupsAccess.noGroups.value }}</div>
                </div>

                <!-- Groups Card Grid -->
                <div v-else class="row q-col-gutter-md">
                  <div v-for="group in groups" :key="group.id" class="col-12 col-sm-4">
                    <div class="access-card">
                      <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">
                        <q-icon name="groups" color="primary" size="xs" class="q-mr-xs" />
                        {{ t.groupsAccess.groupsTitle.value }}
                      </div>
                      <div class="text-subtitle2 text-weight-medium text-grey-9">{{ group.name }}</div>
                      <div v-if="group.description" class="text-caption text-grey-6 q-mt-xs ellipsis-2-lines">
                        {{ group.description }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Memberships Section -->
              <div class="section">
                <div class="section-header q-mb-md">
                  <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
                  <span class="text-subtitle1 text-weight-medium">{{ t.groupsAccess.membershipsTitle.value }}</span>
                  <q-badge color="primary" class="q-ml-sm">{{ memberships.length }}</q-badge>
                </div>

                <!-- Memberships Empty State -->
                <div v-if="memberships.length === 0" class="text-center q-pa-lg">
                  <q-icon name="admin_panel_settings" size="48px" color="grey-4" class="q-mb-sm" />
                  <div class="text-body2 text-grey-6">{{ t.groupsAccess.noMemberships.value }}</div>
                </div>

                <!-- Membership Cards -->
                <div v-else class="row q-col-gutter-md">
                  <div
                    v-for="(membership, index) in memberships"
                    :key="membership.orgId || `membership-${index}`"
                    class="col-12 col-sm-4"
                  >
                    <div class="access-card">
                      <!-- Organization -->
                      <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">
                        <q-icon
                          :name="membership.orgType === 'company' ? 'business' : 'domain'"
                          :color="membership.orgType === 'company' ? 'blue' : 'teal'"
                          size="xs"
                          class="q-mr-xs"
                        />
                        {{ t.groupsAccess.columns.organization.value }}
                      </div>
                      <div class="text-subtitle2 text-weight-medium text-grey-9 q-mb-sm">
                        {{ membership.orgName }}
                      </div>

                      <!-- Roles -->
                      <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">
                        {{ t.groupsAccess.columns.roles.value }}
                      </div>
                      <div class="q-mb-sm">
                        <q-chip
                          v-for="role in membership.roleNames"
                          :key="role"
                          dense
                          color="primary"
                          text-color="white"
                          size="sm"
                          :label="role"
                        />
                      </div>

                      <!-- Scope + Via -->
                      <div class="row q-col-gutter-sm">
                        <div class="col-6">
                          <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">
                            {{ t.groupsAccess.columns.scope.value }}
                          </div>
                          <q-chip
                            dense
                            outline
                            size="sm"
                            :color="membership.scope === 'recursive' ? 'blue' : 'orange'"
                            :icon="membership.scope === 'recursive' ? 'account_tree' : 'place'"
                            :label="membership.scope?.toUpperCase()"
                          />
                        </div>
                        <div class="col-6">
                          <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">
                            {{ t.groupsAccess.columns.via.value }}
                          </div>
                          <q-chip
                            dense
                            outline
                            size="sm"
                            :color="membership.via === 'direct' ? 'green' : 'purple'"
                            :icon="membership.via === 'direct' ? 'person' : 'groups'"
                            :label="membership.via"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </div>

          <!-- Review & Save -->
          <div v-else-if="currentSection === SECTION.REVIEW">
            <q-card-section class="q-pa-lg">
              <div class="q-mb-md">
                <div class="text-subtitle1 text-weight-medium q-mb-xs">
                  <q-icon name="check_circle" color="primary" class="q-mr-xs" />
                  {{ t.review.title.value }}
                </div>
                <div class="text-body2 text-grey-7">
                  {{ t.review.subtitle.value }}
                </div>
              </div>
              <FormReview :sections="previewData" @edit-section="handleSectionChange" />
            </q-card-section>

            <q-separator />

            <q-card-actions align="right" class="q-px-md q-pb-md q-pt-sm">
              <q-btn
                unelevated
                icon-right="save"
                class="rounded-borders"
                color="primary"
                :label="t.buttons.save.value"
                :loading="saving"
                :ripple="false"
                @click="saveProfile"
              />
            </q-card-actions>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
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

.access-card {
  background: white;
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  padding: 14px 16px;
  height: 100%;
  transition: var(--mapex-transition-base);

  &:hover {
    box-shadow: var(--mapex-shadow-md);
  }

  &__title {
    display: flex;
    align-items: center;
    font-size: 13px;
    font-weight: 600;
    color: var(--mapex-text-primary);
    margin-bottom: 8px;
  }

  &__value {
    font-size: 12px;
    color: var(--mapex-text-secondary);
  }
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

:deep(.q-field__control) {
  height: 44px;
}
</style>

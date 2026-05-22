<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="checklist" color="primary" class="q-mr-xs" />
        {{ t.sections.review.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.review.value }}
      </div>
    </div>

    <!-- Personal Information Section -->
    <q-card flat bordered class="q-mb-md rounded-borders" data-testid="review-personal-section">
      <q-card-section>
        <div class="row items-center justify-between q-mb-md">
          <div class="text-subtitle2 text-weight-medium">
            <q-icon name="person" color="primary" class="q-mr-sm" />
            {{ t.sections.personalInfo.value }}
          </div>
          <q-btn
            flat
            dense
            color="primary"
            icon="edit"
            data-testid="review-edit-personal-btn"
            :label="t.buttons.edit.value"
            @click="emit('edit-section', STEP.PERSONAL)"
          />
        </div>

        <div class="row q-col-gutter-md">
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.firstName.value }}</div>
              <div class="review-field__value">{{ modelValue.firstName || '-' }}</div>
            </div>
          </div>
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.lastName.value }}</div>
              <div class="review-field__value">{{ modelValue.lastName || '-' }}</div>
            </div>
          </div>
          <div class="col-12">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.email.value }}</div>
              <div class="review-field__value">{{ modelValue.email || '-' }}</div>
            </div>
          </div>
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.phone.value }}</div>
              <div class="review-field__value">{{ modelValue.phone || '-' }}</div>
            </div>
          </div>
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.jobTitle.value }}</div>
              <div class="review-field__value">{{ modelValue.jobTitle || '-' }}</div>
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Security Information Section -->
    <q-card flat bordered class="q-mb-md rounded-borders" data-testid="review-security-section">
      <q-card-section>
        <div class="row items-center justify-between q-mb-md">
          <div class="text-subtitle2 text-weight-medium">
            <q-icon name="security" color="primary" class="q-mr-sm" />
            {{ t.sections.security.value }}
          </div>
          <q-btn
            flat
            dense
            color="primary"
            icon="edit"
            data-testid="review-edit-security-btn"
            :label="t.buttons.edit.value"
            @click="emit('edit-section', STEP.SECURITY)"
          />
        </div>

        <div class="row q-col-gutter-md">
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.password.value }}</div>
              <div class="review-field__value">
                {{ modelValue.password ? '••••••••' : '-' }}
              </div>
            </div>
          </div>
          <div class="col-12 col-sm-6">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.changePasswordNextLogin.value }}</div>
              <div class="review-field__value">
                <DetailChip
                  :label="modelValue.changePasswordNextLogin ? t.status.yes.value : t.status.no.value"
                  :color="modelValue.changePasswordNextLogin ? 'warning' : 'grey'"
                  size="sm"
                />
              </div>
            </div>
          </div>
          <div class="col-12">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.enabled.value }}</div>
              <div class="review-field__value">
                <DetailChip
                  :label="modelValue.enabled ? t.status.enabled.value : t.status.disabled.value"
                  :color="modelValue.enabled ? 'positive' : 'negative'"
                  :icon="modelValue.enabled ? 'check_circle' : 'cancel'"
                  size="sm"
                />
              </div>
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Access Configuration Section (only for create mode) -->
    <q-card v-if="!isEditMode" flat bordered class="q-mb-md rounded-borders" data-testid="review-access-section">
      <q-card-section>
        <div class="row items-center justify-between q-mb-md">
          <div class="text-subtitle2 text-weight-medium">
            <q-icon name="admin_panel_settings" color="primary" class="q-mr-sm" />
            {{ t.sections.access.value }}
          </div>
          <q-btn
            flat
            dense
            color="primary"
            icon="edit"
            data-testid="review-edit-access-btn"
            :label="t.buttons.edit.value"
            @click="emit('edit-section', STEP.ACCESS)"
          />
        </div>

        <div class="row q-col-gutter-md">
          <!-- Access Type -->
          <div class="col-12">
            <div class="review-field">
              <div class="review-field__label">{{ t.fields.accessType?.value || 'Access Type' }}</div>
              <div class="review-field__value">
                <DetailChip
                  :label="modelValue.accessType === 'group' ? 'Group Membership' : 'Direct Assignment'"
                  :color="modelValue.accessType === 'group' ? 'positive' : 'warning'"
                  :icon="modelValue.accessType === 'group' ? 'group' : 'person'"
                  size="sm"
                />
              </div>
            </div>
          </div>

          <!-- Group Membership Details -->
          <template v-if="modelValue.accessType === 'group' && modelValue.selectedGroup">
            <!-- Group Access Mode -->
            <div class="col-12">
              <div class="review-field">
                <div class="review-field__label">{{ t.reviewLabels.groupMode.value }}</div>
                <div class="review-field__value">
                  <DetailChip
                    :label="modelValue.selectedGroup.mode === 'existing' ? 'Use Existing Group' : 'Create New Group'"
                    :color="modelValue.selectedGroup.mode === 'existing' ? 'primary' : 'amber'"
                    :icon="modelValue.selectedGroup.mode === 'existing' ? 'group' : 'group_add'"
                    size="sm"
                  />
                </div>
              </div>
            </div>

            <!-- Existing Group Details -->
            <template v-if="modelValue.selectedGroup.mode === 'existing' && modelValue.selectedGroup.existingGroup">
              <div class="col-12">
                <div class="review-field">
                  <div class="review-field__label">{{ t.fields.group?.value || 'Group' }}</div>
                  <div class="review-field__value">
                    <DetailChip
                      :label="modelValue.selectedGroup.existingGroup.groupName"
                      color="primary"
                      icon="group"
                      size="sm"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12">
                <q-banner rounded class="bg-blue-1 text-blue-9">
                  <template #avatar>
                    <q-icon name="info" color="blue-6" size="sm" />
                  </template>
                  <div class="text-caption">
                    User will inherit all roles from this group
                  </div>
                </q-banner>
              </div>
            </template>

            <!-- New Group Details -->
            <template v-if="modelValue.selectedGroup.mode === 'new' && modelValue.selectedGroup.newGroup">
              <div class="col-12">
                <div class="review-field">
                  <div class="review-field__label">{{ t.reviewLabels.newGroupName.value }}</div>
                  <div class="review-field__value">
                    <DetailChip
                      :label="modelValue.selectedGroup.newGroup.name"
                      color="amber"
                      icon="group_add"
                      size="sm"
                    />
                  </div>
                </div>
              </div>
              <div v-if="modelValue.selectedGroup.newGroup.description" class="col-12">
                <div class="review-field">
                  <div class="review-field__label">{{ t.reviewLabels.description.value }}</div>
                  <div class="review-field__value">{{ modelValue.selectedGroup.newGroup.description }}</div>
                </div>
              </div>
              <div v-if="modelValue.selectedGroup.newGroup.roleNames?.length" class="col-12">
                <div class="review-field">
                  <div class="review-field__label">{{ t.fields.roles.value }}</div>
                  <div class="review-field__value row items-center q-gutter-xs">
                    <DetailChip
                      v-for="(role, index) in modelValue.selectedGroup.newGroup.roleNames"
                      :key="index"
                      :label="role"
                      color="secondary"
                      icon="badge"
                      size="sm"
                    />
                  </div>
                </div>
              </div>
            </template>
          </template>

          <!-- Direct Assignment Details -->
          <template v-else-if="modelValue.accessType === 'direct' && modelValue.directMembership">
            <div class="col-12">
              <div class="review-field">
                <div class="review-field__label">{{ t.fields.organization.value }}</div>
                <div class="review-field__value">
                  <DetailChip
                    :label="modelValue.directMembership.orgName || modelValue.directMembership.orgId"
                    color="primary"
                    icon="business"
                    size="sm"
                  />
                </div>
              </div>
            </div>
            <div class="col-12">
              <div class="review-field">
                <div class="review-field__label">{{ t.fields.roles.value }}</div>
                <div class="review-field__value row items-center q-gutter-xs">
                  <template v-if="modelValue.directMembership.roleNames?.length">
                    <DetailChip
                      v-for="(role, index) in modelValue.directMembership.roleNames"
                      :key="index"
                      :label="role"
                      color="secondary"
                      icon="badge"
                      size="sm"
                    />
                  </template>
                  <span v-else>-</span>
                </div>
              </div>
            </div>
            <div class="col-12">
              <div class="review-field">
                <div class="review-field__label">{{ t.fields.scope.value }}</div>
                <div class="review-field__value">
                  <DetailChip
                    :label="modelValue.directMembership.scope === 'recursive' ? 'Recursive' : 'Local'"
                    :color="modelValue.directMembership.scope === 'recursive' ? 'purple' : 'teal'"
                    :icon="modelValue.directMembership.scope === 'recursive' ? 'account_tree' : 'location_on'"
                    size="sm"
                  />
                </div>
              </div>
            </div>
          </template>

          <!-- No Selection -->
          <template v-else>
            <div class="col-12">
              <div class="text-grey-6 text-italic">
                {{ t.messages.noAccessConfigured?.value || 'No access configured' }}
              </div>
            </div>
          </template>
        </div>
      </q-card-section>
    </q-card>

    <!-- Summary Banner -->
    <q-banner rounded class="bg-primary-subtle">
      <template #avatar>
        <q-icon name="info" color="primary" />
      </template>
      <div class="text-body2">
        {{ isEditMode ? t.messages.reviewEditSummary.value : t.messages.reviewCreateSummary.value }}
      </div>
    </q-banner>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step4Review',
});

/** TYPE IMPORTS */
import type { Step4ReviewProps } from './interfaces/Step4Review.interface';

/** COMPONENTS */
import { DetailChip } from '@components/chips/DetailChip';

/** COMPOSABLES */
import { useAddUserTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

withDefaults(defineProps<Step4ReviewProps>(), {
  isEditMode: false,
});

const emit = defineEmits<{
  (e: 'edit-section', step: number): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddUserTranslations();
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.review-field {
  &__label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--q-grey-6);
    margin-bottom: 4px;
    letter-spacing: 0.5px;
  }

  &__value {
    font-size: 0.95rem;
    color: var(--q-dark);
  }
}

.bg-primary-subtle {
  background-color: rgba(var(--q-primary-rgb), 0.08);
}
</style>

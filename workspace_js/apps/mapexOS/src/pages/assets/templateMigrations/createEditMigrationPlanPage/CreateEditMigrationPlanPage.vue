<script setup lang="ts">
defineOptions({
  name: 'CreateEditMigrationPlanPage'
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { GroupedStep } from '@components/steppers';
import type { MigrationPlanCreateRequest } from '@mapexos/schemas';
import type { MigrationPlanForm } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import { StepperVertical, buildStepperTree } from '@components/steppers';
import { StepDetails, StepSchedule, StepFromTemplate, StepAssetSelection, StepToTemplate } from './components';

/** COMPOSABLES */
import { useCreateEditMigrationPlanTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_MIGRATION_PLAN_FORM,
  MIGRATION_WIZARD_TOTAL_STEPS,
  MIGRATION_EDIT_TOTAL_STEPS,
  MIGRATION_EDITABLE_STATUSES,
  MIGRATION_ROUTE_BASE,
} from './constants';

/** COMPOSABLES & STORES */
const t = useCreateEditMigrationPlanTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditMigrationPlanPage');

/** STATE */
// A route id switches the wizard into edit mode, where only name/description/
// schedule are mutable — the template and asset selection are fixed at creation.
const planId = computed(() => (route.params.id as string | undefined) ?? null);
const isEditMode = computed(() => !!planId.value);
// Shared reactive model filled across the wizard steps.
const form = ref<MigrationPlanForm>({ ...DEFAULT_MIGRATION_PLAN_FORM });
const currentStep = ref(1);
// The schedule step owns its own validity because "run now" and an empty
// scheduled datetime both leave scheduleAt null but differ in validity.
const scheduleValid = ref(true);
// Per-step validity for the two steps whose model alone cannot express it.
const assetsValid = ref(false);
const toTemplateValid = ref(false);
// Guards the submit while the request is in flight.
const submitting = ref(false);

/** COMPUTED */
// Flat, ordered leaf steps drive navigation and the form-card header; `group`
// collapses them into the stepper tree. Edit mode exposes only the two mutable
// steps (Details, Schedule); create mode adds the immutable selection steps.
const translatedSteps = computed<GroupedStep[]>(() => {
  const setup = { id: 'setup', label: t.groups.setup.value, icon: 'mdi-cog-outline' };
  const migration = { id: 'migration', label: t.groups.migration.value, icon: 'mdi-swap-horizontal' };
  const steps: GroupedStep[] = [
    { title: t.steps.details.label.value, icon: 'mdi-information', description: t.steps.details.description.value, group: setup },
    { title: t.steps.schedule.label.value, icon: 'mdi-calendar-clock', description: t.steps.schedule.description.value, group: setup },
  ];
  if (!isEditMode.value) {
    steps.push(
      { title: t.steps.fromTemplate.label.value, icon: 'mdi-database-arrow-right-outline', description: t.steps.fromTemplate.description.value, group: migration },
      { title: t.steps.assets.label.value, icon: 'mdi-cube-outline', description: t.steps.assets.description.value, group: migration },
      { title: t.steps.toTemplate.label.value, icon: 'mdi-database-arrow-left-outline', description: t.steps.toTemplate.description.value, group: migration },
    );
  }
  return steps;
});

// Grouped tree for the vertical stepper; navigation still runs on translatedSteps.
const stepperTree = computed(() => buildStepperTree(translatedSteps.value));

// Whether the current step is complete enough to advance.
const currentStepValid = computed(() => {
  switch (currentStep.value) {
    case 1:
      return !!form.value.name.trim();
    case 2:
      return scheduleValid.value;
    case 3:
      return !!form.value.fromTemplateId;
    case 4:
      return assetsValid.value;
    case 5:
      return toTemplateValid.value;
    default:
      return true;
  }
});

// All steps that must be valid before the plan can be persisted. Edit mode only
// covers name + schedule; create mode covers the full selection flow.
const allStepsValid = computed(() => {
  const base = !!form.value.name.trim() && scheduleValid.value;
  if (isEditMode.value) return base;
  return base && !!form.value.fromTemplateId && assetsValid.value && !!form.value.toTemplateId;
});

const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: isEditMode.value ? MIGRATION_EDIT_TOTAL_STEPS : MIGRATION_WIZARD_TOTAL_STEPS,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  disableNextButton: !currentStepValid.value,
  disableSaveButton: !allStepsValid.value || submitting.value,
  loadingSaveButton: submitting.value,
}));

const buttonLabels = computed(() => ({
  previous: t.navigation.previous.value,
  next: t.navigation.next.value,
  save: t.navigation.save.value,
}));

const headerTitle = computed(() => (isEditMode.value ? t.page.editTitle.value : t.page.title.value));
const headerDescription = computed(() => (isEditMode.value ? t.page.editDescription.value : t.page.description.value));

/** WATCHERS */
// A different source template invalidates the earlier asset selection.
watch(
  () => form.value.fromTemplateId,
  () => {
    form.value.assetIds = [];
    assetsValid.value = false;
  }
);

/** FUNCTIONS */
/**
 * Move to the target step, blocking any forward move while the current step is
 * invalid. Backward navigation is always allowed.
 * @param step - Target 1-based step index
 * @returns {void}
 */
function changeStep(step: number): void {
  if (step > currentStep.value && !currentStepValid.value) return;
  currentStep.value = step;
}

/**
 * Load the existing plan into the form for edit mode. A plan that is no longer
 * editable (already running or terminal) is rejected back to its detail page,
 * mirroring the backend's 409 guard.
 * @returns {Promise<void>}
 */
async function loadPlan(): Promise<void> {
  if (!planId.value || !apis.assets) return;
  try {
    const plan = await apis.assets.migration.get({ id: planId.value });
    if (!MIGRATION_EDITABLE_STATUSES.includes(plan.status)) {
      notifyFail({ message: t.notifications.notEditable.value });
      await router.push(`${MIGRATION_ROUTE_BASE}/${planId.value}`);
      return;
    }
    form.value = {
      name: plan.name,
      description: plan.description,
      scheduleAt: plan.scheduleAt ?? null,
      fromTemplateId: plan.fromTemplateId,
      toTemplateId: plan.toTemplateId,
      assetIds: [],
    };
  } catch (err: unknown) {
    logger.error('Error loading migration plan:', err);
    notifyFail({ message: t.notifications.loadError.value });
    await router.push(MIGRATION_ROUTE_BASE);
  }
}

/**
 * Persist the wizard form: PATCH the mutable fields in edit mode, otherwise
 * create a new plan. Routes back to the list on success.
 * @returns {Promise<void>}
 */
async function submitPlan(): Promise<void> {
  if (!apis.assets) {
    notifyFail({ message: t.notifications.apiNotInitialized.value });
    return;
  }
  if (!allStepsValid.value) return;

  try {
    submitting.value = true;

    if (isEditMode.value && planId.value) {
      await apis.assets.migration.update({ id: planId.value }, {
        name: form.value.name.trim(),
        description: form.value.description,
        scheduleAt: form.value.scheduleAt,
      });
      notifySuccess({ message: t.notifications.updateSuccess.value });
      await router.push(MIGRATION_ROUTE_BASE);
      return;
    }

    const payload: MigrationPlanCreateRequest = {
      name: form.value.name.trim(),
      description: form.value.description,
      fromTemplateId: form.value.fromTemplateId as string,
      toTemplateId: form.value.toTemplateId as string,
      assetIds: form.value.assetIds,
      scheduleAt: form.value.scheduleAt,
    };

    await apis.assets.migration.create(payload);
    notifySuccess({ message: t.notifications.createSuccess.value });
    await router.push(MIGRATION_ROUTE_BASE);
  } catch (err: unknown) {
    logger.error('Error saving migration plan:', err);
    notifyFail({ message: isEditMode.value ? t.notifications.updateError.value : t.notifications.createError.value });
  } finally {
    submitting.value = false;
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  if (isEditMode.value) void loadPlan();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Header Section -->
    <PageHeader
      icon="swap_horiz"
      iconColor="primary"
      :title="headerTitle"
      :description="headerDescription"
      :button="{ label: t.page.back.value, icon: 'arrow_back', flat: true, to: MIGRATION_ROUTE_BASE }"
    />

    <!-- Content -->
    <div class="row q-col-gutter-lg">
      <!-- Progress Stepper Vertical -->
      <div class="col-12 col-md-4">
        <StepperVertical
          :title="t.stepper.title.value"
          :subtitle="t.stepper.subtitle.value"
          :info-text="t.stepper.requiredInfo.value"
          :current-step-label="t.stepper.currentStep.value"
          :current-step="currentStep"
          :steps="stepperTree"
          step-id-prefix="step"
          @step-click="changeStep"
        />
      </div>

      <!-- Form Card -->
      <div class="col-12 col-md-8">
        <FormCard
          :header="translatedSteps[currentStep - 1] as unknown as FormCardHeader"
          :navigation="formNavigation"
          :button-labels="buttonLabels"
          @previous="changeStep"
          @next="changeStep"
          @save="submitPlan"
        >
          <template #form>
            <!-- STEP 1: DETAILS -->
            <StepDetails
              v-if="currentStep === 1"
              v-model:name="form.name"
              v-model:description="form.description"
            />

            <!-- STEP 2: SCHEDULE -->
            <StepSchedule
              v-else-if="currentStep === 2"
              v-model="form.scheduleAt"
              @update:valid="scheduleValid = $event"
            />

            <!-- STEP 3: FROM TEMPLATE -->
            <StepFromTemplate
              v-else-if="currentStep === 3"
              v-model="form.fromTemplateId"
            />

            <!-- STEP 4: ASSETS -->
            <StepAssetSelection
              v-else-if="currentStep === 4"
              v-model:asset-ids="form.assetIds"
              :from-template-id="form.fromTemplateId"
              @update:valid="assetsValid = $event"
            />

            <!-- STEP 5: TO TEMPLATE -->
            <StepToTemplate
              v-else-if="currentStep === 5"
              v-model="form.toTemplateId"
              @update:valid="toTemplateValid = $event"
            />
          </template>
        </FormCard>
      </div>
    </div>
  </q-page>
</template>

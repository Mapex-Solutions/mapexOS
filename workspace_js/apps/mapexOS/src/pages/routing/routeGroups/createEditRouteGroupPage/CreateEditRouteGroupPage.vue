<script setup lang="ts">
defineOptions({
  name: 'CreateEditRouteGroupPage'
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { RouteGroupCreate } from '@interfaces/routing/routeGroups.interface';
import type { RouterFormState } from './interfaces';
import type { PageTourStep } from '@composables/tour';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import { Step1BasicInfo } from './components/Step1BasicInfo';
import { Step2RoutersConfig } from './components/Step2RoutersConfig';
import { Step3Review } from './components/Step3Review';

/** COMPOSABLES */
import { useRouteGroupsTranslations } from '@composables/i18n/pages/routing/routeGroups/useRouteGroupsTranslations';
import { usePageTour } from '@composables/tour';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ENABLED,
  DEFAULT_IS_TEMPLATE,
  DEFAULT_ROUTER_KIND,
  INITIAL_STEP,
  TOTAL_STEPS,
  STEP,
  ROUTE_GROUP_TOUR_STEPS,
  TOUR_STEP_NAVIGATION,
  TOUR_ROUTER_REQUIRED_STEPS,
} from './constants';
import {
  handleAddRouter,
  handleRemoveRouter,
  handleRouterKindChange,
  handleAddMatchRule,
  handleRemoveMatchRule,
  handleToggleConditionalRouting,
  handleChangeStep,
  handleSave,
  handleLoadRouteGroup,
} from './handlers';

/** COMPOSABLES & STORES */
const router = useRouter();
const route = useRoute();
const t = useRouteGroupsTranslations();
const orgStore = useOrganizationStore();

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const routeGroupId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** STATE */
const currentStep = ref(INITIAL_STEP);
const step1FormRef = ref<InstanceType<typeof Step1BasicInfo> | null>(null);

const formData = ref<RouteGroupCreate>({
  name: '',
  description: '',
  version: '1.0.0',
  enabled: DEFAULT_ENABLED,
  isTemplate: DEFAULT_IS_TEMPLATE,
  routers: [],
});

const routerForms = ref<RouterFormState[]>([
  {
    id: `router-initial-${Date.now()}`,
    kind: DEFAULT_ROUTER_KIND,
    hasConditionalRouting: false,
    saveEvent: {},
  },
]);

/** COMPUTED */

/**
 * Page title changes based on mode
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.createEdit.page.titleEdit.value : t.createEdit.page.title.value
);

/**
 * Page description changes based on mode
 */
const pageDescription = computed(() =>
  isEditMode.value
    ? t.createEdit.page.descriptionEdit.value
    : t.createEdit.page.description.value
);

/**
 * Check if user can create templates
 */
const canCreateTemplate = computed(() => orgStore.isVendor || orgStore.isCustomer);

/**
 * Form navigation configuration
 */
const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: TOTAL_STEPS,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  showCancelButton: true,
  disableNextButton: isSaving.value,
  disableSaveButton: isSaving.value,
  loadingSaveButton: isSaving.value,
}));

/**
 * Button labels - changes based on mode
 */
const buttonLabels = computed(() => ({
  previous: t.createEdit.navigation.previous.value,
  next: t.createEdit.navigation.next.value,
  save: isEditMode.value
    ? t.createEdit.navigation.update.value
    : t.createEdit.navigation.save.value,
}));

/**
 * Translated steps array
 */
const translatedSteps = computed(() => t.createEdit.steps.value);

/** FUNCTIONS */

/**
 * Add a new router to the forms
 *
 * @returns {void}
 */
function addRouter(): void {
  handleAddRouter(routerForms);
}

/**
 * Remove a router from the forms
 *
 * @param {string} routerId - Router ID to remove
 * @returns {void}
 */
function removeRouter(routerId: string): void {
  handleRemoveRouter(routerForms, routerId);
}

/**
 * Handle router kind change
 *
 * @param {string} routerId - Router ID
 * @param {string} newKind - New kind value
 * @returns {void}
 */
function onRouterKindChange(routerId: string, newKind: string): void {
  handleRouterKindChange(routerForms, routerId, newKind);
}

/**
 * Add a match rule to a router
 *
 * @param {string} routerId - Router ID
 * @returns {void}
 */
function addMatchRule(routerId: string): void {
  handleAddMatchRule(routerForms, routerId);
}

/**
 * Remove a match rule from a router
 *
 * @param {string} routerId - Router ID
 * @param {number} ruleIndex - Rule index to remove
 * @returns {void}
 */
function removeMatchRule(routerId: string, ruleIndex: number): void {
  handleRemoveMatchRule(routerForms, routerId, ruleIndex);
}

/**
 * Toggle conditional routing for a router
 *
 * @param {string} routerId - Router ID
 * @param {boolean} enabled - Whether to enable conditional routing
 * @returns {void}
 */
function toggleConditionalRouting(routerId: string, enabled: boolean): void {
  handleToggleConditionalRouting(routerForms, routerId, enabled);
}

/**
 * Change to a different step
 * Validates current step before allowing navigation in create mode
 *
 * @param {number} step - Target step number
 * @returns {Promise<void>}
 */
async function changeStep(step: number): Promise<void> {
  // In edit mode, allow free navigation between steps
  if (isEditMode.value) {
    currentStep.value = step;
    return;
  }

  // In create mode, validate before advancing
  await handleChangeStep(currentStep, step, step1FormRef, routerForms);
}

/**
 * Save or update the route group
 *
 * @returns {Promise<void>}
 */
async function submitForm(): Promise<void> {
  await handleSave(
    isSaving,
    isEditMode,
    routeGroupId,
    formData,
    routerForms,
    router,
    t,
  );
}

/**
 * Load route group data for edit mode
 *
 * @returns {Promise<void>}
 */
async function loadRouteGroupData(): Promise<void> {
  await handleLoadRouteGroup(
    isLoading,
    isEditMode,
    routeGroupId,
    formData,
    routerForms,
    currentStep,
    router,
    t,
  );
}

/** TOUR FUNCTIONS */

/**
 * Build tour steps with resolved translations and navigation callbacks
 * Each step navigates to the appropriate wizard step when highlighted
 *
 * @returns {PageTourStep[]} Tour steps with translations and navigation
 */
function buildTourSteps(): PageTourStep[] {
  return ROUTE_GROUP_TOUR_STEPS.map((step, index) => {
    const key = step.translationKey as keyof typeof t.tour;
    const translation = t.tour[key];

    const result: PageTourStep = {
      element: step.element,
      title: translation.title.value,
      description: translation.description.value,
    };

    if (step.side) result.side = step.side;
    if (step.align) result.align = step.align;

    // Add navigation callback for step changes
    result.onHighlightStarted = () => {
      // Navigate to wizard step if needed
      const targetStep = TOUR_STEP_NAVIGATION[index];
      if (targetStep && currentStep.value !== targetStep) {
        currentStep.value = targetStep;
      }

      // Ensure router exists for router-related steps
      if (TOUR_ROUTER_REQUIRED_STEPS.includes(index)) {
        if (routerForms.value.length === 0) {
          addRouter();
        }
      }
    };

    return result;
  });
}

/** PAGE TOUR */
const { startTour } = usePageTour({
  tourId: 'route-group-builder',
  steps: buildTourSteps,
  onTourEnd: () => {
    // Reset to step 1 when tour ends
    currentStep.value = STEP.BASIC_INFO;
  },
});

/**
 * Handle start tour event from PageHeader
 */
function handleStartTour(): void {
  // Ensure we start on step 1
  currentStep.value = STEP.BASIC_INFO;
  startTour();
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadRouteGroupData();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl" style="min-height: 400px;">
      <div class="column items-center">
        <q-spinner color="primary" size="50px" />
        <span class="q-mt-md text-grey-7">{{ t.createEdit.page.loading.value }}</span>
      </div>
    </div>

    <!-- Form Content -->
    <div v-else>
      <!-- Header Section -->
      <div id="route-group-header">
        <PageHeader
          icon="route"
          iconColor="primary"
          :title="pageTitle"
          :description="pageDescription"
          :button="{
            label: t.createEdit.page.backButton.value,
            icon: 'arrow_back',
            flat: true,
            to: '/routing/route_groups',
          }"
          :tour="{ enabled: true }"
          @start-tour="handleStartTour"
        />
      </div>

      <!-- Content -->
      <div class="row q-col-gutter-lg">
        <!-- Progress Stepper Vertical -->
        <div id="route-group-stepper" class="col-12 col-md-4">
          <StepperVertical
            :title="t.createEdit.stepper.title.value"
            :subtitle="t.createEdit.stepper.subtitle.value"
            :info-text="t.createEdit.stepper.requiredInfo.value"
            :current-step-label="t.createEdit.stepper.currentStep.value"
            :current-step="currentStep"
            :steps="translatedSteps"
            :allow-step-navigation="isEditMode"
            step-id-prefix="route-group-step"
            @step-click="changeStep"
          />
        </div>

        <!-- Form Card -->
        <div class="col-12 col-md-8">
          <FormCard
            :header="translatedSteps[currentStep - 1] as unknown as FormCardHeader"
            :navigation="formNavigation"
            :button-labels="buttonLabels"
            save-button-id="route-group-save-button"
            @previous="changeStep"
            @next="changeStep"
            @save="submitForm"
          >
            <template #form>
              <!-- STEP 1: BASIC INFORMATION -->
              <Step1BasicInfo
                v-if="currentStep === STEP.BASIC_INFO"
                ref="step1FormRef"
                :form-data="formData"
                :status-options="t.statusOptions.value"
                :can-create-template="canCreateTemplate"
                :t="t"
                @update:form-data="formData = $event"
              />

              <!-- STEP 2: ROUTERS CONFIGURATION -->
              <Step2RoutersConfig
                v-else-if="currentStep === STEP.ROUTERS_CONFIG"
                :router-forms="routerForms"
                :router-kind-options="t.routerKindOptions.value"
                :match-policy-options="t.matchPolicyOptions.value"
                :match-operator-options="t.matchOperatorOptions.value"
                :t="t"
                @update:router-forms="routerForms = $event"
                @add-router="addRouter"
                @remove-router="removeRouter"
                @router-kind-change="onRouterKindChange"
                @toggle-conditional-routing="toggleConditionalRouting"
                @add-match-rule="addMatchRule"
                @remove-match-rule="removeMatchRule"
              />

              <!-- STEP 3: REVIEW -->
              <Step3Review
                v-else-if="currentStep === STEP.REVIEW"
                :form-data="formData"
                :router-forms="routerForms"
                @edit-section="changeStep"
              />
            </template>
          </FormCard>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.router-card {
  border-radius: var(--mapex-radius-md);
}
</style>

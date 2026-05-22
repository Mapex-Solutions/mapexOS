import type { TourStepDefinition } from '@composables/tour';

/**
 * Tour step definitions for Route Group Builder
 * Guides users through creating route groups step by step
 *
 * Flow: Overview → Stepper → Step1 (Basic Info) → Step2 (Routers) → Step3 (Review) → Save
 */
export const ROUTE_GROUP_TOUR_STEPS: TourStepDefinition[] = [
  // 1. Page Overview
  {
    element: '#route-group-header',
    translationKey: 'welcome',
    side: 'bottom',
    align: 'start',
  },
  // 2. Stepper Overview
  {
    element: '#route-group-stepper',
    translationKey: 'stepperOverview',
    side: 'right',
    align: 'start',
  },
  // 3. Step 1: Basic Info
  {
    element: '#route-group-step-1',
    translationKey: 'step1Overview',
    side: 'bottom',
    align: 'center',
  },
  // 4. Name Field
  {
    element: '#route-group-field-name',
    translationKey: 'fieldName',
    side: 'bottom',
    align: 'start',
  },
  // 5. Status Field
  {
    element: '#route-group-field-status',
    translationKey: 'fieldStatus',
    side: 'bottom',
    align: 'start',
  },
  // 7. Step 2: Routers Config
  {
    element: '#route-group-step-2',
    translationKey: 'step2Overview',
    side: 'bottom',
    align: 'center',
  },
  // 8. Add Router Button
  {
    element: '#route-group-add-router',
    translationKey: 'addRouter',
    side: 'top',
    align: 'center',
  },
  // 9. Router Card (first one)
  {
    element: '#route-group-router-card-0',
    translationKey: 'routerCard',
    side: 'bottom',
    align: 'start',
  },
  // 10. Router Kind Select
  {
    element: '#route-group-router-kind-0',
    translationKey: 'routerKind',
    side: 'bottom',
    align: 'start',
  },
  // 11. Conditional Routing Toggle
  {
    element: '#route-group-conditional-routing-0',
    translationKey: 'conditionalRouting',
    side: 'bottom',
    align: 'start',
  },
  // 12. Step 3: Review
  {
    element: '#route-group-step-3',
    translationKey: 'step3Overview',
    side: 'bottom',
    align: 'center',
  },
  // 13. Save Button
  {
    element: '#route-group-save-button',
    translationKey: 'saveButton',
    side: 'left',
    align: 'center',
  },
];

/**
 * Step names mapped to tour step indices
 * Used for navigation during tour
 */
export const TOUR_STEP_NAVIGATION: Record<number, number | null> = {
  0: null,    // welcome - no step change
  1: null,    // stepperOverview - no step change
  2: 1,       // step1Overview - go to step 1
  3: 1,       // fieldName - stay on step 1
  4: 1,       // fieldStatus - stay on step 1
  5: 2,       // step2Overview - go to step 2
  6: 2,       // addRouter - stay on step 2
  7: 2,       // routerCard - stay on step 2
  8: 2,       // routerKind - stay on step 2
  9: 2,       // conditionalRouting - stay on step 2
  10: 3,      // step3Overview - go to step 3
  11: null,   // saveButton - no step change
};

/**
 * Step indices that require a router to exist
 * Used to automatically add a router during tour
 */
export const TOUR_ROUTER_REQUIRED_STEPS: number[] = [7, 8, 9];

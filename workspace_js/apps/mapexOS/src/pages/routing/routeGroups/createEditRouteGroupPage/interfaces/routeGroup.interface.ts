import type { Router } from '@interfaces/routing/routeGroups.interface';

/**
 * Router form state with additional UI properties
 */
export interface RouterFormState extends Router {
  /** Temporary ID for tracking the router in the form */
  id: string;

  /** Whether conditional routing is enabled for this router */
  hasConditionalRouting: boolean;

  /** Display name for the selected data lake (UI only) */
  lakeHouseName?: string;

  /** Display name for the selected notification/trigger (UI only) */
  notificationName?: string;

  /** Display name for the selected workflow definition (UI only) */
  workflowName?: string;
}

/**
 * Form navigation state for the CreateEditRouteGroupPage stepper
 */
export interface FormNavigation {
  /** Current active step number */
  currentStep: number;

  /** Total number of steps in the form */
  totalSteps: number;

  /** Whether to show the previous button */
  showPreviousButton: boolean;

  /** Whether to show the next button */
  showNextButton: boolean;

  /** Whether to show the save button */
  showSaveButton: boolean;

  /** Whether to show the cancel button */
  showCancelButton: boolean;

  /** Whether the next button should be disabled */
  disableNextButton: boolean;

  /** Whether the save button should be disabled */
  disableSaveButton: boolean;

  /** Whether the save button should show loading state */
  loadingSaveButton: boolean;
}

/**
 * Button labels for the form navigation
 */
export interface ButtonLabels {
  /** Label for the previous button */
  previous: string;

  /** Label for the next button */
  next: string;

  /** Label for the save button */
  save: string;
}

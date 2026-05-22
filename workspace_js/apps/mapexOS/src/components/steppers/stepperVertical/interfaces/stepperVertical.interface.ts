/**
 * Interface for each step in the stepper
 */
export interface StepperVerticalItem {
  /** Display label for the step */
  title: string;
  /** Description text for the step */
  description: string;
  /** Material icon name for the step */
  icon: string;
  /** Optional custom color for the step icon */
  color?: string;
  /** Optional custom background color for the step icon */
  backgroundColor?: string;
}

/**
 * Type for stepper mode
 */
export type StepperMode = 'creating' | 'editing';

/**
 * Interface for custom step style
 */
export interface StepperVerticalIconStyle {
  backgroundColor?: string;
  borderColor?: string;
}

// Define props with TypeScript
export interface StepperVerticalProps {
  /** Array of step items to display */
  steps: StepperVerticalItem[];
  /** Current active step (1-based index) */
  currentStep?: number;
  /** Header title for the stepper card */
  title?: string;
  /** Subtitle/description below the header title */
  subtitle?: string;
  /** Icon displayed in the header */
  headerIcon?: string;
  /** Informational text displayed at the bottom */
  infoText?: string;
  /** Label for "Current Step" text displayed at the bottom */
  currentStepLabel?: string;
  /** Mode of operation - creating or editing */
  mode?: StepperMode;
  /** Custom style for the height */
  fullHeight?: boolean;
  /**
   * Allow navigation to any step by clicking
   * - true: User can click on any step to navigate (edit mode behavior)
   * - false: User must go step-by-step, can only go back (create mode behavior)
   * Default: false (step-by-step navigation)
   */
  allowStepNavigation?: boolean;
  /**
   * Prefix for generating step IDs (e.g., "step" generates "step-1", "step-2", etc.)
   * Useful for targeting steps with Driver.js tour or CSS selectors
   */
  stepIdPrefix?: string;
}
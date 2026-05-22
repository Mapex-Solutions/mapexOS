/**
 * Icon definition for review sections
 */
export interface ReviewIconDef {
  /** Icon name (Material Icons or MDI) */
  name: string;
  /** Icon color */
  color?: string;
}

/**
 * Field types supported in review
 */
export type ReviewFieldType = 'text' | 'badge' | 'chip' | 'datetime' | 'json' | 'boolean';

/**
 * Field definition for review sections
 */
export interface ReviewFieldDef {
  /** Field label */
  label: string;
  /** Field value */
  value: unknown;
  /** Field type for rendering */
  type: ReviewFieldType;
  /** Badge/chip colors - string for single color, object for value-based colors */
  badgeColors?: string | Record<string, string>;
  /** Date format: 'date', 'time', or 'datetime' */
  format?: string;
  /** Column size (1-12), default is 6 */
  colSize?: number;
  /** Icon to display with chip */
  icon?: string;
}

/**
 * Section definition for review
 */
export interface ReviewSectionDef {
  /** Step number to navigate to when editing */
  stepNumber: number;
  /** Section label/title */
  label: string;
  /** Section icon */
  icon: ReviewIconDef;
  /** Fields in this section */
  fields: ReviewFieldDef[];
  /** Optional data-testid for E2E testing */
  testId?: string;
}

/**
 * Props for FormReview component
 */
export interface FormReviewProps {
  /** Sections to display */
  sections: ReviewSectionDef[];
  /** Enable edit mode with edit buttons */
  editMode?: boolean;
  /** Show success banner at the bottom */
  showSuccessBanner?: boolean;
  /** Custom success banner message */
  successMessage?: string;
  /** Description text shown at the top */
  description?: string;
}

/**
 * Emits for FormReview component
 */
export interface FormReviewEmits {
  (e: 'editSection', stepNumber: number): void;
}

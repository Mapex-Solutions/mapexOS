export interface FormCardHeader {
  icon: string;
  iconColor?: string;
  title: string;
  description: string;
}

export interface FormCardNavigation {
  currentStep: number;
  totalSteps?: number;
  showPreviousButton?: boolean;
  showNextButton?: boolean;
  showSaveButton?: boolean;
  disableNextButton?: boolean;
  disableSaveButton?: boolean;
  loadingSaveButton?: boolean;
}

export interface FormCardButtonLabels {
  previous?: string;
  next?: string;
  save?: string;
}

export interface FormCardProps {
  header: FormCardHeader;
  navigation: FormCardNavigation;
  buttonLabels?: FormCardButtonLabels;
  /** Optional ID for the save button (used for tour highlighting) */
  saveButtonId?: string;
}

export interface FormCardEmits {
  (e: 'save'): void;
  (e: 'previous', n: number): number;
  (e: 'next', n: number): number;
}
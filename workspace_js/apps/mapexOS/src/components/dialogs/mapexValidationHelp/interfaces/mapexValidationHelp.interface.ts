export interface MapexValidationHelpProps {
  modelValue: boolean;
}

export interface MapexValidationHelpEmits {
  (e: 'update:modelValue', value: boolean): void;
}

/**
 * Code example entry for validation help tabs
 */
export interface CodeExample {
  title: string;
  description: string;
  code: string;
  keywords: string[];
}

/**
 * Tab content definition for validation help modal
 */
export interface TabContent {
  icon: string;
  color: string;
  sections: CodeExample[];
}

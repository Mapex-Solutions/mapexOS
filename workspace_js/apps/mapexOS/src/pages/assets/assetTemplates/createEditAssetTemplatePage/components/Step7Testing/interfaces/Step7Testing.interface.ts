import type { AssetTemplateData } from '../../../interfaces';
import type { TestResults } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step7TestingProps {
  modelValue: AssetTemplateData;
  testResults: TestResults;
}

export interface Step7TestingEmits {
  (e: 'update:testResults', value: TestResults): void;
  (e: 'update:modelValue', value: AssetTemplateData): void;
  (e: 'showStandardizedPayloadHelp'): void;
}

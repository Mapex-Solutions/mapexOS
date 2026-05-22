import type { StepDefinition } from '../../../interfaces';

export interface StepProgressProps {
  steps: StepDefinition[];
  currentStep: number;
}

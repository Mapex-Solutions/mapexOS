import type { AssetAttributeForm } from '../../../interfaces';

export interface AttributesSectionProps {
  /** Current list of custom attributes (v-model). */
  modelValue: AssetAttributeForm[];
}

export interface AttributesSectionEmits {
  (e: 'update:modelValue', value: AssetAttributeForm[]): void;
}

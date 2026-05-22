export interface CustomerDetailsDrawerProps {
  modelValue: boolean;
  customerId: string | null;
}

export interface CustomerDetailsDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'edit', customerId: string): void;
}

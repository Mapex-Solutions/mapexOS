import type { RouterFormState } from '../../../interfaces';

/**
 * Props for Step2RoutersConfig component
 */
export interface Step2RoutersConfigProps {
  /** Array of router form states */
  routerForms: RouterFormState[];

  /** Router kind options for dropdown */
  routerKindOptions: Array<{
    label: string;
    value: string;
    icon: string;
    color: string;
    description: string;
  }>;

  /** Match policy options for conditional routing */
  matchPolicyOptions: Array<{
    label: string;
    value: string;
    description: string;
  }>;

  /** Match operator options for rules */
  matchOperatorOptions: Array<{
    label: string;
    value: string;
    description: string;
  }>;

  /** Translation composable */
  t: any;
}

/**
 * Emits for Step2RoutersConfig component
 */
export interface Step2RoutersConfigEmits {
  (e: 'update:routerForms', value: RouterFormState[]): void;
  (e: 'add-router'): void;
  (e: 'remove-router', routerId: string): void;
  (e: 'router-kind-change', routerId: string, kind: string): void;
  (e: 'toggle-conditional-routing', routerId: string, enabled: boolean): void;
  (e: 'add-match-rule', routerId: string): void;
  (e: 'remove-match-rule', routerId: string, ruleIndex: number): void;
}

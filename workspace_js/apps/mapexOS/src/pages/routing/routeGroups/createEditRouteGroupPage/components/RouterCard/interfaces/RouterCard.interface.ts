import type { RouterFormState } from '../../../interfaces';

/**
 * Props for RouterCard component
 */
export interface RouterCardProps {
  /** Router form state data */
  router: RouterFormState;

  /** Index of this router in the array */
  index: number;

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
 * Emits for RouterCard component
 */
export interface RouterCardEmits {
  (e: 'update:router', value: RouterFormState): void;
  (e: 'remove'): void;
  (e: 'kind-change', kind: string): void;
  (e: 'toggle-conditional', enabled: boolean): void;
  (e: 'add-match-rule'): void;
  (e: 'remove-match-rule', ruleIndex: number): void;
}

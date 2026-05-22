import type { Router } from '@interfaces/routing/routeGroups.interface';

/**
 * Props for MatchConfiguration component
 */
export interface MatchConfigurationProps {
  /** Match configuration from the router */
  match: Router['match'];

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
 * Emits for MatchConfiguration component
 */
export interface MatchConfigurationEmits {
  (e: 'update:match', value: Router['match']): void;
  (e: 'add-rule'): void;
  (e: 'remove-rule', ruleIndex: number): void;
}

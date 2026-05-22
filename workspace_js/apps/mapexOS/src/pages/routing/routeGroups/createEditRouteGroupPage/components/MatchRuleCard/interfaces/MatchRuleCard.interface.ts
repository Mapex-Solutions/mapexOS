import type { MatchRule } from '@interfaces/routing/routeGroups.interface';

/**
 * Props for MatchRuleCard component
 */
export interface MatchRuleCardProps {
  /** Match rule data */
  rule: MatchRule;

  /** Index of this rule in the array */
  ruleIndex: number;

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
 * Emits for MatchRuleCard component
 */
export interface MatchRuleCardEmits {
  (e: 'update:rule', value: MatchRule): void;
  (e: 'delete'): void;
}

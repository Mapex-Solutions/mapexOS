/**
 * Durability options for workflow state variables
 */
export const DURABILITY_OPTIONS = [
  { label: 'Ephemeral', value: false },
  { label: 'Durable', value: true },
] as const;

/**
 * Visual variants for the End node based on termination mode
 */
export const END_NODE_VARIANTS = {
  /** Success termination — purple check circle */
  success: {
    icon: 'check_circle',
    hex: '#7B1FA2',
  },
  /** Error termination — red cancel icon */
  error: {
    icon: 'cancel',
    hex: '#D32F2F',
  },
} as const;

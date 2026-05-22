import { useTS } from '@utils/translation';

/**
 * Common message translations
 * Used for notifications, dialogs, and user feedback
 *
 * @example
 * ```ts
 * const { getMessage } = useCommonMessages();
 * notifySuccess({
 *   message: getMessage('deletedSuccessfully', { item: 'User' })
 * });
 * ```
 */
export function useCommonMessages() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Get a specific message with interpolation
     */
    getMessage: (key: string, params: Record<string, unknown> = {}) => {
      return ts(`common.messages.${key}`, params);
    },

    /**
     * Common messages
     */
    messages: {
      success: () => ts('common.messages.success'),
      error: () => ts('common.messages.error'),
      loading: () => ts('common.messages.loading'),
      noData: () => ts('common.messages.noData'),

      // Parameterized messages
      confirmDelete: (item: string) => ts('common.messages.confirmDelete', { item }),
      deletedSuccessfully: (item: string) => ts('common.messages.deletedSuccessfully', { item }),
      savedSuccessfully: (item: string) => ts('common.messages.savedSuccessfully', { item }),
      updatedSuccessfully: (item: string) => ts('common.messages.updatedSuccessfully', { item }),
      createdSuccessfully: (item: string) => ts('common.messages.createdSuccessfully', { item }),
    },
  };
}

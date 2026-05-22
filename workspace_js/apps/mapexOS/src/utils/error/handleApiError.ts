import { notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import type { ApiError, HandleApiErrorOptions } from './types';

const logger = useLogger('handleApiError');

/**
 * Centralized API error handler with i18n support
 *
 * IMPORTANT: This function does NOT use i18n composables internally.
 * You must provide custom messages or a default message.
 *
 * @example Basic usage with custom messages
 * ```ts
 * import { handleApiError } from '@utils/error';
 * import { useMyTranslations } from '@composables/i18n/useMyTranslations';
 *
 * const t = useMyTranslations();
 *
 * try {
 *   await apis.assets.assetTemplate.create(data);
 * } catch (error) {
 *   handleApiError(error, {
 *     customMessages: {
 *       409: t.notifications.alreadyExists,
 *       422: t.notifications.validationFailed,
 *     },
 *     defaultMessage: t.notifications.creationFailed
 *   });
 * }
 * ```
 *
 * @example Usage with global error translations
 * ```ts
 * import { useErrorTranslations } from '@composables/i18n/common/useErrorTranslations';
 *
 * const globalErrors = useErrorTranslations();
 *
 * try {
 *   await apis.users.user.delete({ userId });
 * } catch (error) {
 *   handleApiError(error, {
 *     customMessages: {
 *       404: globalErrors.http[404],
 *       500: globalErrors.http[500],
 *     },
 *     defaultMessage: globalErrors.http.unknown
 *   });
 * }
 * ```
 *
 * @param error - The error object from the API call
 * @param options - Configuration options for error handling
 */
export function handleApiError(error: any, options: HandleApiErrorOptions = {}) {
  const {
    customMessages = {},
    defaultMessage = 'An error occurred. Please try again.',
    timeout = 5000,
    logError = true,
    onError,
  } = options;

  // Log error if enabled
  if (logError) {
    logger.error('API Error:', error);
  }

  // Call custom error handler if provided
  if (onError) {
    onError(error as ApiError);
  }

  // Determine the error message
  let errorMessage: string;

  if (error.response?.status) {
    const statusCode = error.response.status;

    // Check for custom message for this status code
    if (customMessages[statusCode]) {
      const customMsg = customMessages[statusCode];
      errorMessage = typeof customMsg === 'string'
        ? customMsg
        : customMsg.value;
    }
    // Fallback to default message
    else {
      errorMessage = typeof defaultMessage === 'string'
        ? defaultMessage
        : defaultMessage.value;
    }
  }
  // Network error (no response)
  else if (!error.response) {
    if (customMessages.network) {
      const networkMsg = customMessages.network;
      errorMessage = typeof networkMsg === 'string'
        ? networkMsg
        : networkMsg.value;
    } else {
      errorMessage = typeof defaultMessage === 'string'
        ? defaultMessage
        : defaultMessage.value;
    }
  }
  // Unknown error
  else {
    errorMessage = typeof defaultMessage === 'string'
      ? defaultMessage
      : defaultMessage.value;
  }

  // Display notification
  notifyFail({
    message: errorMessage,
    timeout,
  });
}

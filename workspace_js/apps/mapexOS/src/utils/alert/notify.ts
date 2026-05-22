/* eslint no-unused-expressions: ["off"] */
import { Notify, type QNotifyCreateOptions } from 'quasar';

export type NotifyParams = {
	message: QNotifyCreateOptions['message'];
	html?: boolean | undefined;
	timeout?: number | undefined;
	actions?: QNotifyCreateOptions['actions'] | undefined;
};

// Notification manager to prevent stacking
let activeNotification: (() => void) | null = null;
const DISMISS_DELAY = 300; // ms to wait before showing new notification

/**
 * Creates a notification using the Quasar Notify plugin with specified options.
 * Automatically dismisses previous notification to prevent stacking.
 *
 * @param base - An object containing the base properties for the notification, specifically 'color' and 'icon'.
 * @param param1 - An object containing additional notification parameters.
 * @param param1.message - The message to be displayed in the notification.
 * @param param1.html - Optional. If true, the message will be treated as HTML.
 * @param param1.timeout - Optional. The duration in milliseconds before the notification automatically closes.
 * @param param1.actions - Optional. An array of actions to be displayed in the notification.
 *
 * @returns void
 */
function createNotify(
	base: Pick<QNotifyCreateOptions, 'color' | 'icon'>,
	{ message, html, timeout, actions }: NotifyParams,
): void {
	// Dismiss previous notification if exists
	if (activeNotification) {
		activeNotification();
		activeNotification = null;
	}

	// Wait a bit before showing new notification for smooth transition
	setTimeout(() => {
		const opts = {
			...base,
			message,
			progress: true,
			...(html !== undefined && { html }),
			...(timeout !== undefined && { timeout }),
			...(actions !== undefined && { actions }),
		} as QNotifyCreateOptions;

		// Store the dismiss function
		activeNotification = Notify.create(opts);
	}, DISMISS_DELAY);
}


/**
 * Displays a success notification with a positive color and check icon.
 *
 * @param params - An object containing the notification parameters.
 * @param params.message - The message to be displayed in the notification.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.timeout - Optional. The duration in milliseconds before the notification automatically closes.
 * @param params.actions - Optional. An array of actions to be displayed in the notification.
 *
 * @returns void
 */
export function notifySuccess(params: NotifyParams): void {
	createNotify({ color: 'positive', icon: 'check' }, params);
}

/**
 * Displays an informational notification with a primary color and info icon.
 *
 * @param params - An object containing the notification parameters.
 * @param params.message - The message to be displayed in the notification.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.timeout - Optional. The duration in milliseconds before the notification automatically closes.
 * @param params.actions - Optional. An array of actions to be displayed in the notification.
 *
 * @returns void
 */
export function notifyInfo(params: NotifyParams): void {
	createNotify({ color: 'primary', icon: 'info' }, params);
}

/**
 * Displays a failure notification with a red color and report problem icon.
 *
 * @param params - An object containing the notification parameters.
 * @param params.message - The message to be displayed in the notification.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.timeout - Optional. The duration in milliseconds before the notification automatically closes.
 * @param params.actions - Optional. An array of actions to be displayed in the notification.
 *
 * @returns void
 */
export function notifyFail(params: NotifyParams): void {
	createNotify({ color: 'red-4', icon: 'report_problem' }, params);
}

/**
 * Displays a warning notification with an orange color and warning icon.
 *
 * @param params - An object containing the notification parameters.
 * @param params.message - The message to be displayed in the notification.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.timeout - Optional. The duration in milliseconds before the notification automatically closes.
 * @param params.actions - Optional. An array of actions to be displayed in the notification.
 *
 * @returns void
 */
export function notifyWarning(params: NotifyParams): void {
	createNotify({ color: 'warning', icon: 'warning' }, params);
}

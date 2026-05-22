import { Dialog, type QDialogOptions } from 'quasar';

export type DialogParams = {
	title: string;
	message: string;
	html?: boolean | undefined;
	persistent?: boolean | undefined;
	ok?: QDialogOptions['ok'] | undefined;
	cancel?: QDialogOptions['cancel'] | undefined;
};

/**
 * Creates a confirmation dialog using the Quasar Dialog plugin with specified options.
 * Returns a promise that resolves to true if the user confirms, false if canceled.
 *
 * @param options - An object containing the dialog configuration options.
 * @param options.title - The title to be displayed in the dialog.
 * @param options.message - The message to be displayed in the dialog.
 * @param options.html - Optional. If true, the message will be treated as HTML.
 * @param options.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC.
 * @param options.ok - Optional. Configuration for the OK button (label, color, etc).
 * @param options.cancel - Optional. Configuration for the Cancel button (label, color, etc).
 *
 * @returns Promise<boolean> - Resolves to true when OK is clicked, false when Cancel is clicked.
 */
function createDialog(
	baseOptions: Partial<QDialogOptions>,
	{ title, message, html, persistent, ok, cancel }: DialogParams,
): Promise<boolean> {
	return new Promise((resolve) => {
		const opts = {
			...baseOptions,
			title,
			message,
			...(html !== undefined && { html }),
			...(persistent !== undefined && { persistent }),
			...(ok !== undefined && { ok }),
			...(cancel !== undefined && { cancel }),
		} as QDialogOptions;

		Dialog.create(opts)
			.onOk(() => resolve(true))
			.onCancel(() => resolve(false))
			.onDismiss(() => resolve(false));
	});
}

/**
 * Shows a standard confirmation dialog with primary color styling.
 * Useful for general yes/no or ok/cancel confirmations.
 *
 * @param params - An object containing the dialog parameters.
 * @param params.title - The title to be displayed in the dialog.
 * @param params.message - The message to be displayed in the dialog.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC. Defaults to true.
 * @param params.ok - Optional. Configuration for the OK button. Defaults to { label: 'OK', color: 'primary' }.
 * @param params.cancel - Optional. Configuration for the Cancel button. Defaults to { label: 'Cancel', flat: true }.
 *
 * @returns Promise<boolean> - Resolves to true when OK is clicked, false when Cancel is clicked.
 *
 * @example
 * const confirmed = await dialogConfirm({
 *   title: 'Confirm Action',
 *   message: 'Are you sure you want to proceed?'
 * });
 *
 * if (confirmed) {
 *   // User clicked OK
 * }
 */
export function dialogConfirm(params: DialogParams): Promise<boolean> {
	const defaultParams = {
		...params,
		persistent: params.persistent ?? true,
		ok: params.ok ?? { label: 'OK', color: 'primary' },
		cancel: params.cancel ?? { label: 'Cancel', flat: true },
	};

	return createDialog({ cancel: true }, defaultParams);
}

/**
 * Shows a delete confirmation dialog with danger styling (negative color).
 * Specialized for delete operations with destructive actions.
 *
 * @param params - An object containing the dialog parameters.
 * @param params.title - The title to be displayed in the dialog.
 * @param params.message - The message to be displayed in the dialog.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC. Defaults to true.
 * @param params.ok - Optional. Configuration for the Delete button. Defaults to { label: 'Delete', color: 'negative' }.
 * @param params.cancel - Optional. Configuration for the Cancel button. Defaults to { label: 'Cancel', flat: true }.
 *
 * @returns Promise<boolean> - Resolves to true when Delete is clicked, false when Cancel is clicked.
 *
 * @example
 * const confirmed = await dialogDelete({
 *   title: 'Confirm Deletion',
 *   message: `Are you sure you want to delete "${itemName}"?`
 * });
 *
 * if (confirmed) {
 *   // Perform delete operation
 * }
 */
export function dialogDelete(params: DialogParams): Promise<boolean> {
	const defaultParams = {
		...params,
		persistent: params.persistent ?? true,
		ok: params.ok ?? { label: 'Delete', color: 'negative' },
		cancel: params.cancel ?? { label: 'Cancel', flat: true },
	};

	return createDialog({ cancel: true }, defaultParams);
}

/**
 * Shows a warning dialog with warning color styling.
 * Useful for warning messages that require user acknowledgment or confirmation.
 *
 * @param params - An object containing the dialog parameters.
 * @param params.title - The title to be displayed in the dialog.
 * @param params.message - The message to be displayed in the dialog.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC. Defaults to true.
 * @param params.ok - Optional. Configuration for the OK button. Defaults to { label: 'OK', color: 'warning' }.
 * @param params.cancel - Optional. Configuration for the Cancel button. Defaults to { label: 'Cancel', flat: true }.
 *
 * @returns Promise<boolean> - Resolves to true when OK is clicked, false when Cancel is clicked.
 *
 * @example
 * const confirmed = await dialogWarning({
 *   title: 'Warning',
 *   message: 'This action may have unexpected consequences. Continue?'
 * });
 *
 * if (confirmed) {
 *   // User acknowledged warning
 * }
 */
export function dialogWarning(params: DialogParams): Promise<boolean> {
	const defaultParams = {
		...params,
		persistent: params.persistent ?? true,
		ok: params.ok ?? { label: 'OK', color: 'warning' },
		cancel: params.cancel ?? { label: 'Cancel', flat: true },
	};

	return createDialog({ cancel: true }, defaultParams);
}

/**
 * Shows an informational dialog with primary color styling.
 * Useful for displaying information that requires user acknowledgment.
 *
 * @param params - An object containing the dialog parameters.
 * @param params.title - The title to be displayed in the dialog.
 * @param params.message - The message to be displayed in the dialog.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC. Defaults to false.
 * @param params.ok - Optional. Configuration for the OK button. Defaults to { label: 'OK', color: 'primary' }.
 * @param params.cancel - Optional. Not typically used for info dialogs, but can be provided.
 *
 * @returns Promise<boolean> - Resolves to true when OK is clicked.
 *
 * @example
 * await dialogInfo({
 *   title: 'Information',
 *   message: 'Your changes have been saved successfully.'
 * });
 */
export function dialogInfo(params: DialogParams): Promise<boolean> {
	const defaultParams = {
		...params,
		persistent: params.persistent ?? false,
		ok: params.ok ?? { label: 'OK', color: 'primary' },
	};

	return createDialog({}, defaultParams);
}

/**
 * Shows an error dialog with negative color styling.
 * Useful for displaying error messages that require user acknowledgment.
 *
 * @param params - An object containing the dialog parameters.
 * @param params.title - The title to be displayed in the dialog.
 * @param params.message - The message to be displayed in the dialog.
 * @param params.html - Optional. If true, the message will be treated as HTML.
 * @param params.persistent - Optional. If true, the dialog cannot be dismissed by clicking outside or pressing ESC. Defaults to true.
 * @param params.ok - Optional. Configuration for the OK button. Defaults to { label: 'OK', color: 'negative' }.
 * @param params.cancel - Optional. Not typically used for error dialogs, but can be provided.
 *
 * @returns Promise<boolean> - Resolves to true when OK is clicked.
 *
 * @example
 * await dialogError({
 *   title: 'Error',
 *   message: 'Failed to save changes. Please try again.'
 * });
 */
export function dialogError(params: DialogParams): Promise<boolean> {
	const defaultParams = {
		...params,
		persistent: params.persistent ?? true,
		ok: params.ok ?? { label: 'OK', color: 'negative' },
	};

	return createDialog({}, defaultParams);
}

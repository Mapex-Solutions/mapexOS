import { Dialog } from 'quasar';

export type ModalType = 'info' | 'alert' | 'warning' | 'confirm';

export interface ModalOptions {
  type?: ModalType;
  title: string;
  message: string;
  okLabel?: string;
  cancelLabel?: string;
  persistent?: boolean;
  html?: boolean;
}

/**
 * Show a modal dialog with customizable type and styling
 *
 * @param options - Modal configuration options
 * @returns Promise that resolves to true if user clicks OK, false if cancelled
 */
export function showModal(options: ModalOptions): Promise<boolean> {
  const {
    type = 'info',
    title,
    message,
    okLabel = 'OK',
    cancelLabel = 'Cancel',
    persistent = false,
  } = options;

  // Determine icon and color based on type
  const config = getModalConfig(type);

  return new Promise((resolve) => {
    const dialogConfig: any = {
      title: `<div class="row items-center">
        <q-icon name="${config.icon}" color="${config.color}" size="sm" class="q-mr-sm" />
        <span>${title}</span>
      </div>`,
      html: true,
      message,
      ok: {
        label: okLabel,
        color: config.okColor,
        flat: false,
        unelevated: true,
      },
      persistent,
      class: `modal-${type}`,
    };

    // Add cancel button only for confirm and warning types
    if (type === 'confirm' || type === 'warning') {
      dialogConfig.cancel = {
        label: cancelLabel,
        color: 'grey-7',
        flat: true,
      };
    }

    Dialog.create(dialogConfig)
      .onOk(() => resolve(true))
      .onCancel(() => resolve(false))
      .onDismiss(() => resolve(false));
  });
}

/**
 * Get modal configuration based on type
 */
function getModalConfig(type: ModalType) {
  const configs = {
    info: {
      icon: 'info',
      color: 'blue-6',
      okColor: 'primary',
    },
    alert: {
      icon: 'error',
      color: 'red-6',
      okColor: 'negative',
    },
    warning: {
      icon: 'warning',
      color: 'orange-6',
      okColor: 'warning',
    },
    confirm: {
      icon: 'help',
      color: 'primary',
      okColor: 'primary',
    },
  };

  return configs[type];
}

/**
 * Show an info modal
 */
export function showInfo(title: string, message: string, okLabel = 'OK'): Promise<boolean> {
  return showModal({ type: 'info', title, message, okLabel });
}

/**
 * Show an alert modal
 */
export function showAlert(title: string, message: string, okLabel = 'OK'): Promise<boolean> {
  return showModal({ type: 'alert', title, message, okLabel });
}

/**
 * Show a warning modal with cancel option
 */
export function showWarning(
  title: string,
  message: string,
  okLabel = 'Continue',
  cancelLabel = 'Cancel'
): Promise<boolean> {
  return showModal({ type: 'warning', title, message, okLabel, cancelLabel });
}

/**
 * Show a confirmation modal
 */
export function showConfirm(
  title: string,
  message: string,
  okLabel = 'Confirm',
  cancelLabel = 'Cancel'
): Promise<boolean> {
  return showModal({ type: 'confirm', title, message, okLabel, cancelLabel });
}

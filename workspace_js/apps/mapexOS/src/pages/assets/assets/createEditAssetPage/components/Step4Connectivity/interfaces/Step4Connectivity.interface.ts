import type { AssetFormData } from '../../../interfaces';

export interface Step4ConnectivityProps {
  modelValue: AssetFormData;

  /**
   * True when the form is editing an existing asset. Used only for
   * the password field's placeholder/hint copy ("leave blank to keep
   * existing hash" on edit vs "leave blank for cert-only" on create);
   * password is always optional regardless.
   */
  isEditMode?: boolean;
}

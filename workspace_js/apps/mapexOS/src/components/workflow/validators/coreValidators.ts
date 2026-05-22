/** TYPE IMPORTS */
import type { ValidationResult, FieldSourceValue } from '../interfaces';

/** UTILS */
import { isFieldSourceEmpty } from '@src/utils/workflow';

// ────────────────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────────────────

/**
 * Create a valid result
 *
 * @returns {ValidationResult} Valid result
 */
function valid(): ValidationResult {
  return { valid: true, errors: [] };
}

/**
 * Create an invalid result with errors
 *
 * @param {string[]} errors - Error i18n keys
 * @returns {ValidationResult} Invalid result
 */
function invalid(errors: string[]): ValidationResult {
  return { valid: false, errors };
}

/**
 * Recursively validate all condition items in a group/items structure.
 * Returns i18n keys with parameter placeholders ({name}, {groupName}).
 *
 * @param {unknown[]} items - Array of ConditionGroupItem-like objects
 * @param {string[]} errors - Error array to push into
 * @param {string} prefix - Label prefix for error messages
 */
function validateConditionItems(items: unknown[], errors: string[], prefix: string): void {
  for (const raw of items) {
    const item = raw as { type?: string; data?: Record<string, unknown> };
    if (!item?.data) continue;

    if (item.type === 'condition') {
      const field = item.data.field as FieldSourceValue | undefined;
      const value = item.data.value as FieldSourceValue | undefined;
      const name = (item.data.name as string) || 'Condition';

      if (isFieldSourceEmpty(field)) {
        errors.push(`fieldSourceIncomplete::${prefix}${name}`);
      }
      if (isFieldSourceEmpty(value)) {
        errors.push(`comparisonValueIncomplete::${prefix}${name}`);
      }
    } else if (item.type === 'group') {
      const subItems = item.data.items as unknown[] | undefined;
      const groupName = (item.data.name as string) || 'Group';
      if (subItems?.length) {
        validateConditionItems(subItems, errors, `${prefix}${groupName} > `);
      }
    }
  }
}

// ────────────────────────────────────────────────────────────────────────────
// Validators
// ────────────────────────────────────────────────────────────────────────────

/**
 * Validate core/trigger_event node config.
 * Requires a trigger to be selected and all variables to have values.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateTriggerEvent(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const triggerId = config.triggerId as string | undefined;

  if (!triggerId?.trim()) {
    errors.push('triggerRequired');
  }

  const variables = config.variables as Record<string, { path?: string; value?: FieldSourceValue }> | undefined;
  if (variables) {
    for (const [path, v] of Object.entries(variables)) {
      if (isFieldSourceEmpty(v?.value)) {
        errors.push(`variableIncomplete::${path}`);
      }
    }
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/condition node config.
 * Requires at least 1 condition item; each item's field and value must be complete.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateCondition(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const items = config.items as unknown[] | undefined;

  if (!items || items.length === 0) {
    errors.push('atLeastOneCondition');
    return invalid(errors);
  }

  validateConditionItems(items, errors, '');

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/set_state node config.
 * Requires targetField; valueSource must be complete unless operation is 'remove'.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateSetState(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const operation = config.operation as string | undefined;
  const targetField = config.targetField as string | undefined;
  const valueSource = config.valueSource as FieldSourceValue | undefined;

  if (!targetField?.trim()) {
    errors.push('targetFieldRequired');
  }

  if (operation !== 'remove') {
    if (isFieldSourceEmpty(valueSource)) {
      errors.push('valueSourceIncomplete');
    }
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/code node config.
 * Requires script that is not empty or just the default comment template.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateCode(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const script = config.script as string | undefined;
  const DEFAULT_SCRIPT = '// Access: state, event, inputs, nodes\n\nreturn {};';

  if (!script?.trim() || script.trim() === DEFAULT_SCRIPT.trim()) {
    errors.push('scriptRequired');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/switch node config.
 * Requires at least 1 case; each case's conditions must be complete.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateSwitch(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const cases = config.cases as Array<{ id?: string; name?: string; condition?: { items?: unknown[] } }> | undefined;

  if (!cases || cases.length === 0) {
    errors.push('atLeastOneCase');
    return invalid(errors);
  }

  for (const c of cases) {
    const caseName = c.name || c.id || 'Case';
    const items = c.condition?.items;
    if (!items || items.length === 0) {
      errors.push(`caseNeedsCondition::${caseName}`);
    } else {
      validateConditionItems(items, errors, `${caseName} > `);
    }
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/subworkflow node config.
 * Requires workflowId to be selected.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateSubworkflow(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const workflowId = config.workflowId as string | undefined;

  if (!workflowId?.trim()) {
    errors.push('workflowRequired');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/loop node config.
 * Requires source to be a complete field source.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateLoop(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const source = config.source as FieldSourceValue | undefined;

  if (isFieldSourceEmpty(source)) {
    errors.push('loopSourceIncomplete');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/end node config.
 * If terminateWithError is true, errorCode is required.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateEnd(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const terminateWithError = config.terminateWithError as boolean | undefined;
  const errorCode = config.errorCode as string | undefined;

  if (terminateWithError && !errorCode?.trim()) {
    errors.push('errorCodeRequired');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/goto node config.
 * Requires pairLabel to be defined.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateGoto(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const pairLabel = config.pairLabel as string | undefined;

  if (!pairLabel?.trim()) {
    errors.push('portalLabelRequired');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/wait_signal node config.
 * Requires signalName to be defined.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateWaitSignal(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const signalName = config.signalName as string | undefined;

  if (!signalName?.trim()) {
    errors.push('signalNameRequired');
  }

  return errors.length ? invalid(errors) : valid();
}

/**
 * Validate core/wait_for node config.
 * Requires field and compareTo value to be defined.
 * Unary operators (isEmpty, isNotEmpty, isTrue, isFalse) skip compareTo check.
 * Note: WaitFor uses { source, value } shape — NOT FieldSourceValue.
 *
 * @param {Record<string, unknown>} config - Node config
 * @returns {ValidationResult} Validation result
 */
export function validateWaitFor(config: Record<string, unknown>): ValidationResult {
  const errors: string[] = [];
  const field = config.field as string | undefined;
  const operator = config.operator as string | undefined;
  const compareTo = config.compareTo as { source?: string; value?: string } | undefined;

  if (!field?.trim()) {
    errors.push('fieldRequired');
  }

  const unaryOperators = ['isEmpty', 'isNotEmpty', 'isTrue', 'isFalse'];
  if (!unaryOperators.includes(operator ?? '')) {
    if (!compareTo?.value?.toString().trim()) {
      errors.push('compareToIncomplete');
    }
  }

  return errors.length ? invalid(errors) : valid();
}

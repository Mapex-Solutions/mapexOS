import { inject } from 'vue';
import type { IWorkflowEditorContext } from '@src/components/workflow/interfaces';

/** Chave unica para o provide/inject */
export const WORKFLOW_CONTEXT_KEY = Symbol('workflow-editor-context');

/**
 * Acessa o contexto do workflow editor.
 * DEVE ser chamado dentro de um componente descendente de onde o provide() foi feito.
 *
 * @returns {IWorkflowEditorContext} Contexto do editor
 * @throws {Error} Se chamado fora da arvore do workflow editor
 */
export function useWorkflowContext(): IWorkflowEditorContext {
  const ctx = inject<IWorkflowEditorContext>(WORKFLOW_CONTEXT_KEY);
  if (!ctx) {
    throw new Error(
      '[workflow-sdk] useWorkflowContext() chamado fora da arvore do workflow editor. ' +
      'O componente deve ser descendente de onde provide(WORKFLOW_CONTEXT_KEY) foi chamado.',
    );
  }
  return ctx;
}

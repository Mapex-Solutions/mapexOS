import type { WorkflowNode, WorkflowEdge, WorkflowVariable, ExternalSignal } from './workflowNode.interface';
import type { PluginNodeType } from './workflowPlugin.interface';

/**
 * Referencia reativa read-only (structural type).
 *
 * Tanto Ref<T> quanto ComputedRef<T> satisfazem este contrato,
 * evitando incompatibilidades com symbols internos do Vue 3.5+
 * (RefSymbol, ComputedRefSymbol) que vue-tsc nao resolve em .vue files.
 */
export interface ReadonlyReactiveRef<T> {
  readonly value: T;
}

/**
 * Contexto do editor disponivel para plugins via inject().
 *
 * O host app provê esta interface via provide().
 * Plugins consomem via useWorkflowContext().
 */
export interface IWorkflowEditorContext {
  /** Lista reativa de todos os nodes do workflow */
  nodes: ReadonlyReactiveRef<WorkflowNode[]>;

  /** Lista reativa de todas as edges do workflow */
  edges: ReadonlyReactiveRef<WorkflowEdge[]>;

  /** Lista reativa de todas as state variables do workflow */
  states: ReadonlyReactiveRef<WorkflowVariable[]>;

  /** Lista reativa de todos os external signals do workflow */
  externalSignals: ReadonlyReactiveRef<ExternalSignal[]>;

  /**
   * Atualiza a configuracao de um node
   *
   * @param {string} nodeId - ID do node a atualizar
   * @param {Record<string, unknown>} config - Propriedades de config a mergear
   */
  updateNodeConfig: (nodeId: string, config: Record<string, unknown>) => void;

  /**
   * Adiciona uma nota (TextNote) vinculada a um node
   *
   * @param {string} nodeId - ID do node pai
   */
  addNoteToNode: (nodeId: string) => void;

  /**
   * Registra um snapshot para undo/redo
   *
   * @param {string} label - Descricao da acao
   */
  pushSnapshot: (label: string) => void;

  /**
   * Busca a definicao de um node type pelo identificador
   *
   * @param {string} type - Tipo do node (ex: 'core/condition')
   * @returns {PluginNodeType | undefined} Definicao do tipo ou undefined
   */
  getNodeType: (type: string) => PluginNodeType | undefined;
}

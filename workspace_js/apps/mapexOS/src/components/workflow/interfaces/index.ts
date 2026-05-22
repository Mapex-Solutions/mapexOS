/**
 * SDK interfaces — public contracts for workflow plugins
 */
export type {
  // Core types
  NodePropertyType,
  NodePropertyDefinition,
  PropertyRendering,
  PropertyFetchOptions,
  FetchOptionsRule,
  DynamicNodeFormProps,
  DynamicNodeFormEmits,
  PluginCategory,
  ValidationResult,
  PluginNodeType,
  WorkflowPlugin,
  CatalogGroup,
  HandleDefinition,
  HandleResolver,
  HandleOverrides,
  ResolvedHandles,
  Disposable,
  PluginActivationContext,

  // Action contract
  ActionDef,
  HttpActionDef,
  ActionOutputDef,
  OperationDefinition,

  // Marketplace plugin types
  PluginMetadata,
  PluginCredentialDefinition,
  CredentialFieldDefinition,
} from './workflowPlugin.interface';

export type {
  NodeTimeoutConfig,
  NodeErrorHandlerConfig,
  WorkflowNode,
  WorkflowEdge,
  WorkflowVariable,
  CaptureField,
  ExternalSignal,
  ExternalVariable,
} from './workflowNode.interface';

export type {
  BaseWorkflowNodeProps,
  WorkflowNodeComponentProps,
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
} from './nodeConfig.interface';

export type { ReadonlyReactiveRef, IWorkflowEditorContext } from './workflowContext.interface';

export type {
  SourceType,
  FieldSourceValue,
  NodeOutputOption,
  SourceTypeOption,
  AssetStatusFieldOption,
} from './fieldSource.interface';

export type {
  GroupLogicOperator,
  WorkflowConditionItem,
  ConditionGroupItem,
  WorkflowConditionGroup,
} from './conditionNode.interface';

export type { SwitchCase } from './switchNode.interface';

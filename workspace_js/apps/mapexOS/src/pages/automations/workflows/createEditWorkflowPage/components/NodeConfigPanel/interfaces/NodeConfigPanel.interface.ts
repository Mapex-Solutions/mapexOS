/**
 * Props for NodeConfigPanel component
 */
export interface NodeConfigPanelProps {
  /** Selected node ID */
  nodeId: string;
}

/**
 * Emits for NodeConfigPanel component
 */
export interface NodeConfigPanelEmits {
  (e: 'close'): void;
}

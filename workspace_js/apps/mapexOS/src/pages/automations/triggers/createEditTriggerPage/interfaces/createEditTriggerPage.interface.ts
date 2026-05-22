/**
 * Import types from schemas package
 */
import type { TriggerTypeEnum, TriggerCategoryEnum } from '@mapexos/schemas';

/**
 * Re-export schema types for UI usage
 */
export type TriggerCategory = TriggerCategoryEnum;
export type TriggerType = TriggerTypeEnum;

/**
 * Technical trigger types (📡)
 */
export type TechnicalTriggerType = Extract<TriggerTypeEnum,
  TriggerTypeEnum.HTTP |
  TriggerTypeEnum.MQTT |
  TriggerTypeEnum.RABBITMQ |
  TriggerTypeEnum.NATS |
  TriggerTypeEnum.WEBSOCKET
>;

/**
 * Communication trigger types (💬)
 */
export type CommunicationTriggerType = Extract<TriggerTypeEnum,
  TriggerTypeEnum.EMAIL |
  TriggerTypeEnum.TEAMS |
  TriggerTypeEnum.SLACK
>;

/**
 * Trigger entity structure
 * This is the unified structure for both technical and communication triggers
 */
export interface Trigger {
  id?: string;
  name: string;
  description?: string;
  triggerType: TriggerType;
  category: TriggerCategory;
  enabled: boolean;
  isTemplate?: boolean;
  config: Record<string, any>;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Category option for Step 1
 */
export interface CategoryOption {
  value: TriggerCategory;
  label: string;
  description: string;
  icon: string;
  emoji: string;
}

/**
 * Trigger type option for Step 2
 */
export interface TriggerTypeOption {
  value: TriggerType;
  label: string;
  description: string;
  icon: string;
  category: TriggerCategory;
}

/**
 * Form state for managing stepper navigation
 */
export interface TriggerFormState {
  selectedCategory: TriggerCategory | null;
  selectedType: TriggerType | null;
  isCreating: boolean;
  currentStep: number;
}

// Shared JSON type
type JsonData = object | string;

// Props for the extended JSON+Event drawer
export interface JsonAndEventDrawerProps {
  show: boolean;
  jsonData: JsonData;
  editable?: boolean;
  title: string;
  subtitle?: string;
  eventData?: any;
}

// Emits for the extended JSON+Event drawer
export interface JsonAndEventDrawerEmit {
  (e: 'update:show', value: boolean): void;
  (e: 'save', updated: JsonData): void;
  /** When the user switches to the event tab and data is missing */
  (e: 'fetch-event', eventId: string): void;
}
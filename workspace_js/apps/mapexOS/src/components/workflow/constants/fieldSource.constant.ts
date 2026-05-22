/** TYPE IMPORTS */
import type { FieldSourceValue, SourceTypeOption, AssetStatusFieldOption } from '../interfaces';

/**
 * Available source type options for field source selectors.
 * Used by FieldSourceSelector and any plugin UI that needs source type selection.
 */
export const SOURCE_TYPE_OPTIONS: readonly SourceTypeOption[] = [
  { label: 'Asset Status', value: 'assetStatus', icon: 'monitor_heart', color: 'orange-7' },
  { label: 'Event Field', value: 'event', icon: 'event', color: 'blue-6' },
  { label: 'Fetch Options', value: 'fetchOptions', icon: 'cloud_download', color: 'purple-6' },
  { label: 'Input', value: 'input', icon: 'input', color: 'cyan-6' },
  { label: 'Literal', value: 'literal', icon: 'format_quote', color: 'green-6' },
  { label: 'Node Output', value: 'nodeOutput', icon: 'output', color: 'teal-6' },
  { label: 'State', value: 'state', icon: 'storage', color: 'purple-6' },
] as const;

/**
 * Default field source value — literal type with empty value.
 */
export const DEFAULT_FIELD_SOURCE_VALUE: FieldSourceValue = {
  type: 'literal',
  value: '',
};

/**
 * Predefined fields available in health monitoring events (sensor.offline / sensor.online).
 * Used by FieldSourceSelector when mode is 'assetStatus'.
 */
export const ASSET_STATUS_FIELD_OPTIONS: readonly AssetStatusFieldOption[] = [
  { label: 'Event Type', value: 'eventType', icon: 'label', type: 'string', availability: 'all' },
  { label: 'Source', value: 'source', icon: 'source', type: 'string', availability: 'all' },
  { label: 'Asset UUID', value: 'assetUUID', icon: 'fingerprint', type: 'string', availability: 'all' },
  { label: 'Asset Name', value: 'assetName', icon: 'badge', type: 'string', availability: 'all' },
  { label: 'Last Seen At', value: 'lastSeenAt', icon: 'schedule', type: 'timestamp', availability: 'offline' },
  { label: 'Offline Since', value: 'offlineSince', icon: 'timer_off', type: 'timestamp', availability: 'offline' },
  { label: 'Threshold (min)', value: 'thresholdMinutes', icon: 'timelapse', type: 'number', availability: 'offline' },
  { label: 'Miss Count', value: 'missCount', icon: 'error_outline', type: 'number', availability: 'offline' },
] as const;

import type { Ref } from 'vue';

import { notifyFail, notifySuccess } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';

import type { HttpDataSource } from '../interfaces/httpDataSource.interface';

const logger = useLogger('useHttpDataSourceHandlers');

export function useHttpDataSourceHandlers(dataSource: Ref<HttpDataSource>, t?: any) {
  /**
   * Add a new time interval to the working hours configuration
   * Adds a default interval from 09:00 to 17:00
   * @returns {void}
   */
  function addInterval(): void {
    dataSource.value.timeIntervals.push({ startTime: '09:00', endTime: '17:00' });
  }

  /**
   * Remove a time interval from the working hours configuration
   * Prevents removal if only one interval remains
   * @param {number} index - Index of the interval to remove
   * @returns {void}
   */
  function removeInterval(index: number): void {
    if (dataSource.value.timeIntervals.length > 1) {
      dataSource.value.timeIntervals.splice(index, 1);
    }
  }

  /**
   * Add a new custom UUID path for field mapping in asset binding
   * Adds an empty path that can be configured by the user
   * @returns {void}
   */
  function addMapping(): void {
    dataSource.value.customUuidPaths.push({ path: '' });
  }

  /**
   * Remove a custom UUID path from field mapping configuration
   * Prevents removal if only one path remains
   * @param {number} index - Index of the path to remove
   * @returns {void}
   */
  function removeMapping(index: number): void {
    if (dataSource.value.customUuidPaths.length > 1) {
      dataSource.value.customUuidPaths.splice(index, 1);
    }
  }

  /**
   * Test UUID path extraction against example payload
   * Validates that configured paths can extract values from the JSON payload
   * Shows success notification with extracted values or error if invalid
   * @returns {void}
   */
  function testMapping(): void {
    try {
      const payload = JSON.parse(dataSource.value.payloadExample);
      const mappedValues: string[] = [];

      // Test finalUuidPaths if available
      const paths = dataSource.value.finalUuidPaths || [];
      paths.forEach((path: string) => {
        const value = getValueFromPath(payload, path);
        mappedValues.push(`${path}: ${value || 'not found'}`);
      });

      if (mappedValues.length > 0) {
        notifySuccess({ message: `UUID extraction test: ${mappedValues.join(', ')}`, timeout: 5000 });
      } else {
        notifyFail({ message: t?.errors?.noUuidPaths?.value ?? 'No UUID paths configured to test' });
      }
    } catch {
      notifyFail({ message: t?.errors?.invalidPayloadOrPath?.value ?? 'Invalid JSON payload or path format' });
    }
  }

  /**
   * Extract a value from an object using a dot-notation path
   * @param {any} obj - The object to extract value from
   * @param {string} path - Dot-notation path (e.g., 'device.uuid')
   * @returns {any} The extracted value or undefined if not found
   */
  function getValueFromPath(obj: any, path: string): any {
    return path.split('.').reduce((current, prop) => current?.[prop], obj);
  }

  /**
   * Save the HTTP data source (placeholder for API integration)
   * TODO: Integrate with actual API endpoint
   * @returns {boolean} True if successful, false otherwise
   */
  function saveDataSource(): boolean {
    try {
      // TODO: Call API to save dataSource.value
      logger.debug('Creating data source:', dataSource.value);
      notifySuccess({ message: t?.notifications?.createSuccess?.value ?? 'Data Source created successfully!' });
      return true;
    } catch {
      notifyFail({ message: t?.notifications?.createFailed?.value ?? 'Failed to create data source' });
      return false;
    }
  }

  return {
    addInterval,
    removeInterval,
    addMapping,
    removeMapping,
    testMapping,
    saveDataSource,
  };
}

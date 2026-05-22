/**
 * Script Tester Handler
 * Handles testing of asset template scripts via js-executor API
 */

import { apis } from '@services/mapex';
import type { TestResults } from '../interfaces';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('ScriptTesterHandler');

/**
 * Execute script tests via js-executor API
 *
 * @param scriptProcessor - Preprocessor script (optional)
 * @param scriptValidator - Validation script
 * @param scriptConversion - Conversion script
 * @param testInput - Test input JSON string
 * @returns TestResults with execution status and output
 */
export async function executeScriptTests(
  scriptProcessor: string,
  scriptValidator: string,
  scriptConversion: string,
  testInput: string
): Promise<TestResults> {
  const testResults: TestResults = {
    executed: false,
    success: false,
    steps: [],
    output: null,
    logs: [],
  };

  try {
    // Parse test input
    let eventPayload: any;
    try {
      eventPayload = JSON.parse(testInput || '{}');
    } catch {
      testResults.success = false;
      testResults.steps = [{
        name: 'Parse Test Input',
        success: false,
        error: 'Invalid JSON in test input'
      }];
      testResults.executed = true;
      return testResults;
    }

    // Check if API is configured
    if (!apis.jsExecutor) {
      testResults.steps = [{
        name: 'API Error',
        success: false,
        error: 'JS Executor API not configured'
      }];
      testResults.executed = true;
      return testResults;
    }

    // Call API to test scripts
    const response = await apis.jsExecutor.scripts.test({
      debugEnabled: true,
      decode: scriptProcessor || '',
      validation: scriptValidator || '',
      transform: scriptConversion || '',
      event: eventPayload
    });

    logger.debug('API Response:', response);

    // API always returns 200, but actual status is in response.success
    if (response.success) {
      // ✅ Success - Display returned payload
      testResults.success = true;
      testResults.steps = response.steps || [];
      testResults.output = response.output || response.result || null;
      testResults.standardizedPayload = response.standardizedPayload || null;
      testResults.responseData = response || null;
      testResults.newPayload = response.data || null;
      testResults.logs = response.logs || [];

      logger.debug('Test Success! StandardizedPayload:', testResults.standardizedPayload);
    } else {
      // ❌ Error - Display API error message
      testResults.success = false;

      // Extract error message from response
      let errorMessage = 'Script execution failed';
      let errorDetails: any = null;

      if (response.error) {
        errorMessage = typeof response.error === 'string'
          ? response.error
          : response.error.message || errorMessage;
        errorDetails = response.error.details || null;
      } else if (response.message) {
        errorMessage = response.message;
      }

      // If API returns steps with errors, use them
      if (response.steps && Array.isArray(response.steps)) {
        testResults.steps = response.steps;
      } else {
        // Create a single error step
        testResults.steps = [{
          name: 'Script Execution Error',
          success: false,
          error: errorMessage,
          details: errorDetails
        }];
      }

      // Log error for debugging
      logger.error('Test Failed:', {
        error: errorMessage,
        details: errorDetails,
        steps: testResults.steps
      });
    }

  } catch (error: any) {
    // Network error or unexpected error (rare)
    logger.error('Unexpected error during test:', error);
    testResults.success = false;
    testResults.steps = [{
      name: 'Network Error',
      success: false,
      error: error.message || 'Failed to communicate with test API'
    }];
  }

  testResults.executed = true;
  return testResults;
}

/**
 * Format JSON for display
 *
 * @param obj - Object to format
 * @returns Formatted JSON string
 */
export function formatJSON(obj: any): string {
  return JSON.stringify(obj, null, 2);
}

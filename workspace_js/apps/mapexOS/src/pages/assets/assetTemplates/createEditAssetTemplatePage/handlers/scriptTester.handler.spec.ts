import { describe, it, expect, vi, beforeEach } from 'vitest';
import { executeScriptTests, formatJSON } from './scriptTester.handler';
import { apis } from '@services/mapex';

/**
 * Mock useLogger to prevent side effects
 */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
    warn: vi.fn(),
    info: vi.fn(),
  }),
}));

/**
 * Mock jsExecutor API (not in global setup)
 */
vi.mock('@services/mapex', () => ({
  apis: {
    jsExecutor: {
      scripts: {
        test: vi.fn(),
      },
    },
  },
}));

const mockTest = vi.mocked(apis.jsExecutor.scripts.test);

describe('scriptTester.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('executeScriptTests', () => {
    it('returns parse error when test input is invalid JSON', async () => {
      const result = await executeScriptTests('', '', '', '{invalid');

      expect(result.executed).toBe(true);
      expect(result.success).toBe(false);
      expect(result.steps).toHaveLength(1);
      expect(result.steps[0]!.name).toBe('Parse Test Input');
      expect(result.steps[0]!.error).toBe('Invalid JSON in test input');
    });

    it('parses empty string as empty object', async () => {
      mockTest.mockResolvedValue({ success: true, steps: [], logs: [] });

      await executeScriptTests('', '', '', '');

      expect(mockTest).toHaveBeenCalledWith(
        expect.objectContaining({ event: {} }),
      );
    });

    it('returns API error when jsExecutor is not configured', async () => {
      // Temporarily remove jsExecutor
      const original = apis.jsExecutor;
      (apis as any).jsExecutor = undefined;

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.executed).toBe(true);
      expect(result.success).toBe(false);
      expect(result.steps[0]!.error).toBe('JS Executor API not configured');

      // Restore
      (apis as any).jsExecutor = original;
    });

    it('calls API with correct parameters', async () => {
      mockTest.mockResolvedValue({ success: true, steps: [], logs: [] });

      await executeScriptTests('decode_code', 'validate_code', 'transform_code', '{"key":"value"}');

      expect(mockTest).toHaveBeenCalledWith({
        debugEnabled: true,
        decode: 'decode_code',
        validation: 'validate_code',
        transform: 'transform_code',
        event: { key: 'value' },
      });
    });

    it('maps successful API response correctly', async () => {
      mockTest.mockResolvedValue({
        success: true,
        steps: [{ name: 'Decode', success: true }],
        output: { decoded: true },
        standardizedPayload: { id: '123' },
        data: { converted: true },
        logs: ['log1', 'log2'],
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.executed).toBe(true);
      expect(result.success).toBe(true);
      expect(result.steps).toEqual([{ name: 'Decode', success: true }]);
      expect(result.output).toEqual({ decoded: true });
      expect(result.standardizedPayload).toEqual({ id: '123' });
      expect(result.newPayload).toEqual({ converted: true });
      expect(result.logs).toEqual(['log1', 'log2']);
    });

    it('uses response.result as fallback for output', async () => {
      mockTest.mockResolvedValue({
        success: true,
        steps: [],
        result: { fallback: true },
        logs: [],
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.output).toEqual({ fallback: true });
    });

    it('handles failed API response with string error', async () => {
      mockTest.mockResolvedValue({
        success: false,
        error: 'Syntax error on line 5',
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.executed).toBe(true);
      expect(result.success).toBe(false);
      expect(result.steps[0]!.error).toBe('Syntax error on line 5');
    });

    it('handles failed API response with object error', async () => {
      mockTest.mockResolvedValue({
        success: false,
        error: { message: 'TypeError', details: { line: 10 } },
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.success).toBe(false);
      expect(result.steps[0]!.error).toBe('TypeError');
      expect(result.steps[0]!.details).toEqual({ line: 10 });
    });

    it('handles failed API response with message field', async () => {
      mockTest.mockResolvedValue({
        success: false,
        message: 'Validation failed',
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.success).toBe(false);
      expect(result.steps[0]!.error).toBe('Validation failed');
    });

    it('uses API steps when response includes them on failure', async () => {
      const apiSteps = [
        { name: 'Decode', success: true },
        { name: 'Validate', success: false, error: 'Failed' },
      ];
      mockTest.mockResolvedValue({
        success: false,
        error: 'Script failed',
        steps: apiSteps,
      });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.steps).toEqual(apiSteps);
    });

    it('falls back to default error message when no error info provided', async () => {
      mockTest.mockResolvedValue({ success: false });

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.steps[0]!.error).toBe('Script execution failed');
    });

    it('handles network/unexpected errors', async () => {
      mockTest.mockRejectedValue(new Error('Network timeout'));

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.executed).toBe(true);
      expect(result.success).toBe(false);
      expect(result.steps[0]!.name).toBe('Network Error');
      expect(result.steps[0]!.error).toBe('Network timeout');
    });

    it('handles unexpected error without message', async () => {
      mockTest.mockRejectedValue({});

      const result = await executeScriptTests('', '', '', '{}');

      expect(result.steps[0]!.error).toBe('Failed to communicate with test API');
    });

    it('initializes testResults with correct defaults', async () => {
      mockTest.mockResolvedValue({ success: true, steps: [], logs: [] });

      const result = await executeScriptTests('', '', '', '{}');

      // Verify the base structure is always present
      expect(result).toHaveProperty('executed');
      expect(result).toHaveProperty('success');
      expect(result).toHaveProperty('steps');
      expect(result).toHaveProperty('output');
      expect(result).toHaveProperty('logs');
    });
  });

  describe('formatJSON', () => {
    it('formats object with 2-space indentation', () => {
      const result = formatJSON({ key: 'value' });
      expect(result).toBe('{\n  "key": "value"\n}');
    });

    it('formats null', () => {
      expect(formatJSON(null)).toBe('null');
    });

    it('formats array', () => {
      const result = formatJSON([1, 2, 3]);
      expect(result).toBe('[\n  1,\n  2,\n  3\n]');
    });

    it('formats nested object', () => {
      const result = formatJSON({ a: { b: 'c' } });
      expect(result).toContain('"a"');
      expect(result).toContain('"b"');
    });
  });
});

import { describe, it, expect, vi } from 'vitest';
import { ref, inject } from 'vue';
import { useWorkflowContext, WORKFLOW_CONTEXT_KEY } from './useWorkflowContext';

/**
 * Mock vue's inject to simulate provide/inject behavior.
 */
vi.mock('vue', async () => {
  // eslint-disable-next-line @typescript-eslint/consistent-type-imports
  const actual = await vi.importActual<typeof import('vue')>('vue');
  return {
    ...actual,
    inject: vi.fn(),
  };
});

const mockedInject = vi.mocked(inject);

describe('useWorkflowContext', () => {
  describe('WORKFLOW_CONTEXT_KEY', () => {
    it('is a Symbol', () => {
      expect(typeof WORKFLOW_CONTEXT_KEY).toBe('symbol');
    });

    it('has a descriptive name', () => {
      expect(WORKFLOW_CONTEXT_KEY.toString()).toContain('workflow-editor-context');
    });
  });

  describe('when context is provided', () => {
    it('returns the injected context', () => {
      const mockContext = {
        nodes: ref([]),
        edges: ref([]),
        states: ref([]),
        updateNodeConfig: vi.fn(),
        addNoteToNode: vi.fn(),
        pushSnapshot: vi.fn(),
        getNodeType: vi.fn(),
      };

      mockedInject.mockReturnValue(mockContext);

      const result = useWorkflowContext();

      expect(mockedInject).toHaveBeenCalledWith(WORKFLOW_CONTEXT_KEY);
      expect(result).toBe(mockContext);
    });

    it('returns a context with all expected properties', () => {
      const mockContext = {
        nodes: ref([{ id: '1' }]),
        edges: ref([{ id: 'e1' }]),
        states: ref([{ name: 'x' }]),
        updateNodeConfig: vi.fn(),
        addNoteToNode: vi.fn(),
        pushSnapshot: vi.fn(),
        getNodeType: vi.fn(),
      };

      mockedInject.mockReturnValue(mockContext);

      const ctx = useWorkflowContext();

      expect(ctx.nodes.value).toHaveLength(1);
      expect(ctx.edges.value).toHaveLength(1);
      expect(ctx.states.value).toHaveLength(1);
      expect(typeof ctx.updateNodeConfig).toBe('function');
      expect(typeof ctx.addNoteToNode).toBe('function');
      expect(typeof ctx.pushSnapshot).toBe('function');
      expect(typeof ctx.getNodeType).toBe('function');
    });
  });

  describe('when context is NOT provided', () => {
    it('throws an error with a descriptive message', () => {
      mockedInject.mockReturnValue(undefined);

      expect(() => useWorkflowContext()).toThrowError(
        '[workflow-sdk] useWorkflowContext() chamado fora da arvore do workflow editor.',
      );
    });

    it('the error message mentions provide(WORKFLOW_CONTEXT_KEY)', () => {
      mockedInject.mockReturnValue(undefined);

      expect(() => useWorkflowContext()).toThrowError(/provide\(WORKFLOW_CONTEXT_KEY\)/);
    });
  });
});

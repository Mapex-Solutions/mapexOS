import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import type { TriggerListPageFilters } from './triggerListPage.handler';
import {
  fetchTriggersHandler,
  handleFilterApplyHandler,
  handlePageChangeHandler,
} from './triggerListPage.handler';
import { apis } from '@services/mapex';

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

function makeFilters(overrides: Partial<TriggerListPageFilters> = {}): TriggerListPageFilters {
  return {
    name: undefined,
    status: undefined,
    includeChildren: undefined,
    category: undefined,
    triggerType: undefined,
    ...overrides,
  };
}

describe('triggerListPage.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure triggers API mock exists
    (apis as any).triggers = {
      trigger: {
        list: vi.fn(),
      },
    };
  });

  // ── fetchTriggersHandler ─────────────────────────────────────────────
  describe('fetchTriggersHandler', () => {
    it('should return triggers and pagination data on success', async () => {
      (apis.triggers as any).trigger.list.mockResolvedValue({
        items: [
          { id: 't-1', name: 'Trigger1' },
          { id: 't-2', name: 'Trigger2' },
        ],
        pagination: { totalPages: 2, totalItems: 15 },
      });

      const result = await fetchTriggersHandler(makeFilters(), 1, 10);

      expect(result.triggers).toHaveLength(2);
      expect(result.totalPages).toBe(2);
      expect(result.totalItems).toBe(15);
    });

    it('should include all filter params when set', async () => {
      (apis.triggers as any).trigger.list.mockResolvedValue({
        items: [],
        pagination: { totalPages: 0, totalItems: 0 },
      });

      const filters = makeFilters({
        name: 'mqtt',
        status: true,
        includeChildren: true,
        category: 'technical',
        triggerType: 'mqtt',
      });

      await fetchTriggersHandler(filters, 1, 10);

      const callArgs = (apis.triggers as any).trigger.list.mock.calls[0][0];
      expect(callArgs.name).toBe('mqtt');
      expect(callArgs.enabled).toBe(true);
      expect(callArgs.includeChildren).toBe(true);
      expect(callArgs.category).toBe('technical');
      expect(callArgs.triggerType).toBe('mqtt');
    });

    it('should not include undefined filters in query params', async () => {
      (apis.triggers as any).trigger.list.mockResolvedValue({
        items: [],
        pagination: { totalPages: 0, totalItems: 0 },
      });

      await fetchTriggersHandler(makeFilters(), 1, 10);

      const callArgs = (apis.triggers as any).trigger.list.mock.calls[0][0];
      expect(callArgs).toEqual({ page: 1, perPage: 10 });
    });

    it('should return empty result on error', async () => {
      (apis.triggers as any).trigger.list.mockRejectedValue(new Error('Fail'));

      const result = await fetchTriggersHandler(makeFilters(), 1, 10);

      expect(result.triggers).toEqual([]);
      expect(result.totalPages).toBe(0);
      expect(result.totalItems).toBe(0);
    });

    it('should pass correct pagination params', async () => {
      (apis.triggers as any).trigger.list.mockResolvedValue({
        items: [],
        pagination: { totalPages: 5, totalItems: 50 },
      });

      await fetchTriggersHandler(makeFilters(), 3, 25);

      const callArgs = (apis.triggers as any).trigger.list.mock.calls[0][0];
      expect(callArgs.page).toBe(3);
      expect(callArgs.perPage).toBe(25);
    });
  });

  // ── handleFilterApplyHandler ─────────────────────────────────────────
  describe('handleFilterApplyHandler', () => {
    it('should update filters, reset page, and call fetch', () => {
      const filters = makeFilters();
      const currentPage = ref(3);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { name: 'test', status: true, category: 'communication', triggerType: 'email' },
        filters,
        currentPage,
        fetchCallback,
      );

      expect(filters.name).toBe('test');
      expect(filters.status).toBe(true);
      expect(filters.category).toBe('communication');
      expect(filters.triggerType).toBe('email');
      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });

    it('should set includeChildren filter', () => {
      const filters = makeFilters();
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { includeChildren: true },
        filters,
        currentPage,
        fetchCallback,
      );

      expect(filters.includeChildren).toBe(true);
    });

    it('should clear filter values when empty', () => {
      const filters = makeFilters({ name: 'old', category: 'technical' });
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { name: '', category: '' },
        filters,
        currentPage,
        fetchCallback,
      );

      expect(filters.name).toBeUndefined();
      expect(filters.category).toBeUndefined();
    });
  });

  // ── handlePageChangeHandler ──────────────────────────────────────────
  describe('handlePageChangeHandler', () => {
    it('should set current page and call fetch', () => {
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handlePageChangeHandler(5, currentPage, fetchCallback);

      expect(currentPage.value).toBe(5);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });
});

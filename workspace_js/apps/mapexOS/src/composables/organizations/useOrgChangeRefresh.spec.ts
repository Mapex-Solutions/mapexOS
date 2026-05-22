import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { nextTick } from 'vue';
import { useOrgChangeRefresh } from './useOrgChangeRefresh';
import { useOrganizationStore } from '@stores/organization';

/**
 * Mock useLogger
 */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
  }),
}));

describe('useOrgChangeRefresh', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('returns an object (empty for now)', () => {
    const result = useOrgChangeRefresh(vi.fn());
    expect(result).toBeDefined();
    expect(typeof result).toBe('object');
  });

  it('does not call callback on initial mount (no org change)', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);
    await nextTick();

    expect(callback).not.toHaveBeenCalled();
  });

  it('calls callback when organization changes from one ID to another', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);

    // Change org
    store.selectedOrganizationId = 'org-2';
    await nextTick();

    expect(callback).toHaveBeenCalledTimes(1);
  });

  it('does not call callback when org changes from null', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = null;

    useOrgChangeRefresh(callback);

    // Set initial org (from null to something)
    store.selectedOrganizationId = 'org-1';
    await nextTick();

    expect(callback).not.toHaveBeenCalled();
  });

  it('does not call callback when org changes to null', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);

    store.selectedOrganizationId = null;
    await nextTick();

    expect(callback).not.toHaveBeenCalled();
  });

  it('calls callback multiple times for multiple org changes', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);

    store.selectedOrganizationId = 'org-2';
    await nextTick();

    store.selectedOrganizationId = 'org-3';
    await nextTick();

    expect(callback).toHaveBeenCalledTimes(2);
  });

  it('handles async callbacks', async () => {
    const callback = vi.fn().mockResolvedValue(undefined);
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);

    store.selectedOrganizationId = 'org-2';
    await nextTick();

    expect(callback).toHaveBeenCalledTimes(1);
  });

  it('does not call callback when org stays the same', async () => {
    const callback = vi.fn();
    const store = useOrganizationStore();
    store.selectedOrganizationId = 'org-1';

    useOrgChangeRefresh(callback);

    // Set same value
    store.selectedOrganizationId = 'org-1';
    await nextTick();

    expect(callback).not.toHaveBeenCalled();
  });

  describe('immediate option', () => {
    it('triggers watch immediately when immediate is true', async () => {
      const callback = vi.fn();
      const store = useOrganizationStore();
      store.selectedOrganizationId = 'org-1';

      useOrgChangeRefresh(callback, { immediate: true });
      await nextTick();

      // immediate: true fires the watcher with (newVal, undefined)
      // But the guard checks newOrgId && oldOrgId, so it should NOT call callback
      expect(callback).not.toHaveBeenCalled();
    });
  });
});

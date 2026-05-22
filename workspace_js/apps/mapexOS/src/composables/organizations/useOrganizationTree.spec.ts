import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useOrganizationTree } from './useOrganizationTree';
import { useOrganizationStore } from '@stores/organization';
import type { OrganizationTreeNode } from '@stores/organization/types';

/**
 * Helper to create a mock tree node
 */
function makeNode(overrides: Partial<OrganizationTreeNode> = {}): OrganizationTreeNode {
  return {
    id: 'org-1',
    name: 'Org A',
    type: 'vendor',
    pathKey: '000001',
    scope: 'inherited',
    membershipId: 'm-1',
    roleIds: [],
    depth: 0,
    enabled: true,
    children: [],
    ...overrides,
  };
}

describe('useOrganizationTree', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  describe('initial state', () => {
    it('returns empty treeNodes when store is empty', () => {
      const { treeNodes, filters } = useOrganizationTree();

      expect(treeNodes.value).toEqual([]);
      expect(filters.value.enabled).toBe('active');
      expect(filters.value.types).toEqual([]);
      expect(filters.value.name).toBe('');
    });
  });

  describe('filteredTreeNodes', () => {
    it('returns all active nodes when default filter (active) is set', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'Active Org', enabled: true }),
        makeNode({ id: '2', name: 'Inactive Org', enabled: false }),
      ];

      const { treeNodes } = useOrganizationTree();

      // Default filter is 'active', so only enabled nodes should show
      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
    });

    it('returns all nodes when enabled filter is "all"', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'Active', enabled: true }),
        makeNode({ id: '2', name: 'Inactive', enabled: false }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.enabled = 'all';

      expect(treeNodes.value).toHaveLength(2);
    });

    it('returns only inactive nodes when enabled filter is "inactive"', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'Active', enabled: true }),
        makeNode({ id: '2', name: 'Inactive', enabled: false }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.enabled = 'inactive';

      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('2');
    });

    it('filters by name (case-insensitive)', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'MAPEX Corp', enabled: true }),
        makeNode({ id: '2', name: 'Other Corp', enabled: true }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.name = 'mapex';

      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
    });

    it('filters by type', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'Vendor A', type: 'vendor', enabled: true }),
        makeNode({ id: '2', name: 'Customer A', type: 'customer', enabled: true }),
        makeNode({ id: '3', name: 'Site A', type: 'site', enabled: true }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.types = ['vendor', 'customer'];

      expect(treeNodes.value).toHaveLength(2);
    });

    it('preserves parent hierarchy when child matches name filter', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({
          id: '1',
          name: 'Parent Corp',
          enabled: true,
          children: [
            makeNode({ id: '2', name: 'MAPEX Site', type: 'site', enabled: true }),
          ],
        }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.name = 'mapex';

      // Parent should be included because child matches
      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
      expect(treeNodes.value[0]!.children).toHaveLength(1);
      expect(treeNodes.value[0]!.children![0]!.id).toBe('2');
    });

    it('preserves parent hierarchy when child matches type filter', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({
          id: '1',
          name: 'Vendor Corp',
          type: 'vendor',
          enabled: true,
          children: [
            makeNode({ id: '2', name: 'Site A', type: 'site', enabled: true }),
          ],
        }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.types = ['site'];

      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
      expect(treeNodes.value[0]!.children).toHaveLength(1);
    });

    it('preserves parent when child is active and filter is active', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({
          id: '1',
          name: 'Inactive Parent',
          enabled: false,
          children: [
            makeNode({ id: '2', name: 'Active Child', enabled: true }),
          ],
        }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.enabled = 'active';

      // Parent should be included because child is active
      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
    });

    it('excludes subtree when no node matches any filter', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({
          id: '1',
          name: 'Parent',
          type: 'vendor',
          enabled: true,
          children: [
            makeNode({ id: '2', name: 'Child', type: 'customer', enabled: true }),
          ],
        }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.name = 'nonexistent';

      expect(treeNodes.value).toHaveLength(0);
    });

    it('combines multiple filters (name + type + enabled)', () => {
      const store = useOrganizationStore();
      store.treeNodes = [
        makeNode({ id: '1', name: 'MAPEX Vendor', type: 'vendor', enabled: true }),
        makeNode({ id: '2', name: 'MAPEX Site', type: 'site', enabled: true }),
        makeNode({ id: '3', name: 'MAPEX Disabled', type: 'vendor', enabled: false }),
        makeNode({ id: '4', name: 'Other Vendor', type: 'vendor', enabled: true }),
      ];

      const { treeNodes, filters } = useOrganizationTree();
      filters.value.name = 'mapex';
      filters.value.types = ['vendor'];
      filters.value.enabled = 'active';

      expect(treeNodes.value).toHaveLength(1);
      expect(treeNodes.value[0]!.id).toBe('1');
    });

    it('returns empty array when store treeNodes is empty', () => {
      const store = useOrganizationStore();
      store.treeNodes = [];

      const { treeNodes } = useOrganizationTree();
      expect(treeNodes.value).toEqual([]);
    });
  });

  describe('selectOrganization', () => {
    it('delegates to store.selectOrganization', () => {
      const store = useOrganizationStore();
      const selectSpy = vi.spyOn(store, 'selectOrganization');

      const { selectOrganization } = useOrganizationTree();
      selectOrganization('org-123');

      expect(selectSpy).toHaveBeenCalledWith('org-123');
    });
  });

  describe('refreshCoverage', () => {
    it('delegates to store.refreshCoverage', async () => {
      const store = useOrganizationStore();
      const refreshSpy = vi.spyOn(store, 'refreshCoverage').mockResolvedValue();

      const { refreshCoverage } = useOrganizationTree();
      await refreshCoverage();

      expect(refreshSpy).toHaveBeenCalled();
    });
  });

  describe('computed store properties', () => {
    it('loading reflects store loading state', () => {
      const store = useOrganizationStore();
      store.loading = true;

      const { loading } = useOrganizationTree();
      expect(loading.value).toBe(true);
    });

    it('error reflects store error state', () => {
      const store = useOrganizationStore();
      store.error = 'Something failed';

      const { error } = useOrganizationTree();
      expect(error.value).toBe('Something failed');
    });

    it('selectedOrganizationId reflects store state', () => {
      const store = useOrganizationStore();
      store.selectedOrganizationId = 'org-42';

      const { selectedOrganizationId } = useOrganizationTree();
      expect(selectedOrganizationId.value).toBe('org-42');
    });
  });
});

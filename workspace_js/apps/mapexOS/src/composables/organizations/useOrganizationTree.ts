import { computed, ref } from 'vue';
import { useOrganizationStore } from '@stores/organization';
import type { OrganizationTreeNode } from '@stores/organization/types';
import type { OrganizationFilters } from 'src/types/organization';

/**
 * Organization Tree Composable
 *
 * Wrapper around organization store providing filtering and tree management.
 * Integrates with backend API through Pinia store.
 *
 * Features:
 * - Access organization tree from store
 * - Apply name, type, and enabled filters
 * - Reactive filtering on store data
 * - Select/deselect organizations
 *
 * @example
 * ```ts
 * const { treeNodes, loading, filters, applyFilters, selectOrganization } = useOrganizationTree();
 *
 * // Apply filter
 * filters.value.name = 'MAPEX';
 * applyFilters();
 *
 * // Select organization
 * selectOrganization('org-id');
 * ```
 */
export function useOrganizationTree() {
  const store = useOrganizationStore();

  const filters = ref<OrganizationFilters>({
    enabled: 'active',
    types: [],
    name: '',
  });

  /**
   * Check if a node or any of its descendants match the search term
   * Recursive function to enable hierarchical search
   */
  function nodeOrChildrenMatchName(node: OrganizationTreeNode, searchTerm: string): boolean {
    // Check current node
    if (node.name.toLowerCase().includes(searchTerm)) {
      return true;
    }

    // Check children recursively
    if (node.children && node.children.length > 0) {
      return node.children.some(child => nodeOrChildrenMatchName(child, searchTerm));
    }

    return false;
  }

  /**
   * Check if a node or any of its descendants match the type filter
   * Enables showing parent nodes when children match type
   */
  function nodeOrChildrenMatchType(node: OrganizationTreeNode, types: string[]): boolean {
    // Check current node
    if (types.includes(node.type)) {
      return true;
    }

    // Check children recursively
    if (node.children && node.children.length > 0) {
      return node.children.some(child => nodeOrChildrenMatchType(child, types));
    }

    return false;
  }

  /**
   * Check if a node or any of its descendants match the enabled filter
   * Enables showing parent nodes when children match enabled status
   */
  function nodeOrChildrenMatchEnabled(
    node: OrganizationTreeNode,
    enabledFilter: 'all' | 'active' | 'inactive'
  ): boolean {
    // Check current node
    if (enabledFilter === 'all') {
      return true;
    }

    const nodeMatches =
      (enabledFilter === 'active' && node.enabled) ||
      (enabledFilter === 'inactive' && !node.enabled);

    if (nodeMatches) {
      return true;
    }

    // Check children recursively
    if (node.children && node.children.length > 0) {
      return node.children.some(child => nodeOrChildrenMatchEnabled(child, enabledFilter));
    }

    return false;
  }

  /**
   * Filter tree nodes recursively with hierarchical logic
   * Returns tree with only nodes matching filters OR parent nodes needed for hierarchy
   *
   * Key Logic:
   * - If a child matches, parent is included to maintain tree structure
   * - Filters are applied recursively to entire subtree
   * - All filter types (name, type, enabled) work together
   */
  function filterTreeNodes(nodes: OrganizationTreeNode[]): OrganizationTreeNode[] {
    return nodes
      .map(node => {
        // First, recursively filter children
        const filteredChildren = node.children
          ? filterTreeNodes(node.children)
          : [];

        // Check if current node OR its children match all active filters
        let nodeOrChildrenMatch = true;

        // Filter by enabled status
        if (filters.value.enabled && filters.value.enabled !== 'all') {
          nodeOrChildrenMatch = nodeOrChildrenMatch &&
            nodeOrChildrenMatchEnabled(node, filters.value.enabled);
        }

        // Filter by types (only if types filter is active)
        if (filters.value.types && filters.value.types.length > 0) {
          nodeOrChildrenMatch = nodeOrChildrenMatch &&
            nodeOrChildrenMatchType(node, filters.value.types);
        }

        // Filter by name (only if name filter is active)
        if (filters.value.name && filters.value.name.trim() !== '') {
          const searchTerm = filters.value.name.toLowerCase();
          nodeOrChildrenMatch = nodeOrChildrenMatch &&
            nodeOrChildrenMatchName(node, searchTerm);
        }

        // Include node if it matches filters OR has filtered children (to maintain hierarchy)
        if (nodeOrChildrenMatch || filteredChildren.length > 0) {
          return {
            ...node,
            children: filteredChildren,
          };
        }

        return null;
      })
      .filter((node): node is NonNullable<typeof node> => node !== null);
  }

  /**
   * Computed filtered tree nodes from store
   * Applies filters reactively to store's treeNodes
   */
  const filteredTreeNodes = computed(() => {
    if (!store.treeNodes || store.treeNodes.length === 0) {
      return [];
    }
    return filterTreeNodes(store.treeNodes);
  });

  /**
   * Select an organization (set as active context)
   */
  function selectOrganization(orgId: string) {
    store.selectOrganization(orgId);
  }

  /**
   * Refresh coverage from API
   */
  async function refreshCoverage() {
    await store.refreshCoverage();
  }

  return {
    // Computed from store with filters applied
    treeNodes: filteredTreeNodes,

    // Direct store access
    loading: computed(() => store.loading),
    error: computed(() => store.error),
    selectedOrganizationId: computed(() => store.selectedOrganizationId),
    selectedOrganization: computed(() => store.selectedOrganization),

    // Filters
    filters,

    // Actions
    selectOrganization,
    refreshCoverage,
  };
}

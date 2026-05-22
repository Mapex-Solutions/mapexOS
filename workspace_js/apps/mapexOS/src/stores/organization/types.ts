export type OrganizationType = 'vendor' | 'customer' | 'site' | 'building' | 'floor' | 'zone';
export type OrganizationScope = 'inherited' | 'recursive';

export interface OrganizationCoverageItem {
  id: string;
  name: string;
  type: OrganizationType;
  pathKey: string;
  scope: OrganizationScope;
  membershipId: string;
  roleIds: string[];
}

export interface OrganizationCoverageResponse {
  lastUpdated: string;
  organizations: OrganizationCoverageItem[];
}

export interface OrganizationTreeNode {
  id: string;
  name: string;
  type: OrganizationType;
  pathKey: string;
  scope: OrganizationScope;
  membershipId: string;
  roleIds: string[];
  depth: number;
  enabled: boolean;
  children?: OrganizationTreeNode[];
}

export interface OrganizationState {
  // Data
  coverage: OrganizationCoverageResponse | null;
  treeNodes: OrganizationTreeNode[];
  flatList: OrganizationCoverageItem[];

  // UI State
  loading: boolean;
  error: string | null;
  lastUpdated: string | null;

  // Selected organization (current context)
  selectedOrganizationId: string | null;
  selectedOrganizationName: string | null;
}

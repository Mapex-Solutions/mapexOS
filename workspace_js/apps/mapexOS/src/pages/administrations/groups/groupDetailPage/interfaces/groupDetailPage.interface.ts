/**
 * Tab configuration for GroupDetailPage
 */
export interface GroupDetailTab {
	/** Unique key for the tab */
	name: string;
	/** Display label */
	label: string;
	/** Icon name */
	icon: string;
	/** Optional badge count */
	badge?: number;
	/** Badge color */
	badgeColor?: string;
}

/**
 * Member info from API response
 */
export interface GroupMemberInfo {
	/** Member ID (junction table) */
	id: string;
	/** User ID */
	userId: string;
	/** User email */
	userEmail?: string;
	/** User first name */
	userFirstName?: string;
	/** User last name */
	userLastName?: string;
	/** When the user was added */
	addedAt?: string;
	/** Who added the user */
	addedBy?: string;
}

/**
 * Group detail data for display (matches API response)
 */
export interface GroupDetailData {
	id: string;
	name: string;
	description?: string;
	enabled: boolean;
	orgId?: string;
	organizationName?: string;
	pathKey?: string;
	membersCount?: number;
	roleIds?: string[];
	created?: string;
	updated?: string;
}

/**
 * Props for TabInfo component
 */
export interface TabInfoProps {
	/** Group data to display */
	group: GroupDetailData | null;
	/** Loading state */
	loading: boolean;
}

/**
 * Props for TabMembers component
 */
export interface TabMembersProps {
	/** Group ID to fetch members */
	groupId: string;
	/** Loading state */
	loading: boolean;
}

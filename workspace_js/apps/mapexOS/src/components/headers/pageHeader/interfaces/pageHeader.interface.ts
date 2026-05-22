/**
 * Interface for the header action button.
 * You can add more fields as needed.
 */
export interface PageHeaderButton {
	label: string;          // Button text
	icon?: string;          // Icon name for the button
	color?: string;         // Button color
	flat?: boolean;         // Flat style toggle
	rounded?: boolean;      // Rounded corners toggle
	unelevated?: boolean;    // Unelevated style toggle
	ripple?: boolean;       // Ripple effect toggle
	to?: string;            // Route to navigate (Vue Router)
	onClick?: () => void;   // Custom click handler
	id?: string;            // HTML id attribute for targeting (e.g., tours)
}

/**
 * Interface for info item in the modal
 */
export interface PageHeaderInfoItem {
	icon?: string;
	color?: string;
	title?: string;
	text: string;
}

/**
 * Interface for the info button configuration
 */
export interface PageHeaderInfo {
	title: string;              // Modal title
	description: string;        // Modal description
	items?: PageHeaderInfoItem[]; // List of features/items
	docsUrl?: string;           // Documentation URL
	docsLabel?: string;         // Documentation button label
}

/**
 * Tour button configuration for PageHeader
 */
export interface PageHeaderTour {
	/** Whether tour button is enabled */
	enabled: boolean;

	/** Optional tooltip text for the tour button */
	tooltip?: string;
}

/**
 * Interface for the header props.
 */
export interface PageHeaderProps {
	icon?: string;              // Main icon for the header
	iconColor?: string;         // Icon color
	title: string;              // Main title (required)
	description?: string;       // Subtitle/description
	button?: PageHeaderButton | undefined;      // Button properties
	info?: PageHeaderInfo;      // Info modal configuration
	tour?: PageHeaderTour;      // Tour button configuration
}
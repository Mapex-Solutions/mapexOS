/**
 * Interface for the header action button.
 * You can add more fields as needed.
 */
export interface ListHeadersButton {
	label: string;          // Button text
	icon?: string;          // Icon name for the button
	color?: string;         // Button color
	flat?: boolean;         // Flat style toggle
	rounded?: boolean;      // Rounded corners toggle
	unelevated?: boolean;    // Unelevated style toggle
	ripple?: boolean;       // Ripple effect toggle
	to?: string;            // Route to navigate (Vue Router)
}

/**
 * Interface for the header props.
 */
export interface ListHeadersProps {
	title: string;             // Main title (required)
	icon: string;              // icon for the header
	button: ListHeadersButton; // Button properties
}
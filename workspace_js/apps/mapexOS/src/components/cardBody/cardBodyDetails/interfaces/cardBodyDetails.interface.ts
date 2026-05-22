interface IconsGroup {
	icon?: string;
	color?: string;
	tooltip?: string; // Tooltip for the icon
}

interface CardBodyDetailItem {
	name?: string;

	icon?: string; // Icon for the item
	iconColor?: string; // Color for the icon

	color?: string; // Value color for the item
	value?: string | number; // Value to be displayed
	cols?: number;

	size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'; // Size of the item (xs, sm, md, lg, xl)
	type?: 'text' | 'icon' | 'iconsGroup' | 'card' | 'chip';
	icons?: IconsGroup[];
	tooltip?: string; // Tooltip for the item
}

interface CardContainerProps {
	color?: string;
}

export interface CarbBodyDetailsProps {
	tenantName?: string;
	title: string;
	items: CardBodyDetailItem[];
	container?: CardContainerProps;
}
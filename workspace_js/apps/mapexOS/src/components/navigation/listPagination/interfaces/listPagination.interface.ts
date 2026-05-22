/**
 * ListPagination Component Interfaces
 *
 * Responsive pagination component for list pages.
 * Automatically adjusts the number of visible pages based on screen size.
 */

/**
 * Props interface for ListPagination component
 */
export interface ListPaginationProps {
	/**
	 * Current active page (v-model)
	 */
	modelValue: number;

	/**
	 * Total number of pages
	 */
	totalPages: number;

	/**
	 * Optional: Custom color for pagination
	 * @default 'primary'
	 */
	color?: string;

	/**
	 * Optional: Custom active color
	 * @default 'primary'
	 */
	activeColor?: string;
}

/**
 * Emits interface for ListPagination component
 */
export interface ListPaginationEmits {
	/**
	 * Emitted when page changes
	 * @param e - Event name
	 * @param value - New page number
	 */
	(e: 'update:modelValue', value: number): void;

	/**
	 * Emitted when page changes (alternative event)
	 * @param e - Event name
	 * @param value - New page number
	 */
	(e: 'change', value: number): void;
}

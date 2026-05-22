/**
 * EventStorePage Constants
 */

import type { TourStepDefinition } from '@composables/tour';

/**
 * Default limit for cursor pagination
 */
export const DEFAULT_LIMIT = 15;

/**
 * Default column visibility state for event store list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
	assetName: true,
	templateName: true,
	threadId: true,
	source: true,
	created: true,
} as const;

/**
 * Default fallback color
 */
export const DEFAULT_COLOR = 'grey-6';

/**
 * Tour step definitions for the Event Store page
 * Pattern: header → searchInput → advancedFiltersBtn → advancedFiltersOpen → dynamicFiltersBtn → dynamicFiltersOpen → results
 * Text comes from translations, these define targeting and positioning
 */
export const EVENT_STORE_TOUR_STEPS: TourStepDefinition[] = [
	{
		element: '#page-header-section',
		translationKey: 'header',
		side: 'bottom',
		align: 'start',
	},
	{
		element: '#filter-search-input',
		translationKey: 'searchInput',
		side: 'bottom',
		align: 'start',
	},
	{
		element: '#advanced-filters-btn',
		translationKey: 'advancedFiltersBtn',
		side: 'bottom',
		align: 'end',
	},
	{
		element: '.drawer-content',
		translationKey: 'advancedFiltersOpen',
		side: 'left',
		align: 'start',
	},
	{
		element: '#dynamic-filters-btn',
		translationKey: 'dynamicFiltersBtn',
		side: 'bottom',
		align: 'end',
	},
	{
		element: '.drawer-content',
		translationKey: 'dynamicFiltersOpen',
		side: 'left',
		align: 'start',
	},
	{
		element: '#results-section',
		translationKey: 'results',
		side: 'top',
		align: 'start',
	},
];

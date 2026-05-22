/**
 * UserListPage Constants
 */

import type { TourStepDefinition, TourTransition } from '@composables/tour';
import type { UserResponse } from '@mapexos/schemas';

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for users list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  email: true,
  jobTitle: true,
  groups: true,
  created: true,
} as const;

/**
 * Default filter values for users list
 */
export const FILTER_DEFAULTS = {
  email: undefined,
  firstName: undefined,
  lastName: undefined,
  enabled: undefined,
  includeChildren: undefined,
} as const;

/**
 * API projection fields for users query
 */
export const USERS_PROJECTION = 'firstName,lastName,email,enabled,jobTitle,avatar,groupsCount' as const;

/**
 * Tour button step definition (conditional - shown when PageHeader has tour.enabled)
 * This step is prepended to the tour to explain the tour button functionality
 */
export const TOUR_BUTTON_STEP: TourStepDefinition = {
  element: '#tour-start-btn',
  translationKey: 'tourButton',
  side: 'bottom',
  align: 'start',
};

/**
 * Tour step definitions for the Users list page
 * Standard pattern: header → filters → advancedFiltersBtn → advancedFiltersOpen → results → addNew
 * Text comes from translations, these define targeting and positioning
 */
export const USER_LIST_TOUR_STEPS: TourStepDefinition[] = [
  {
    element: '#page-header-section',
    translationKey: 'header',
    side: 'bottom',
    align: 'start',
  },
  {
    element: '#filter-section',
    translationKey: 'filters',
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
    element: '#results-section',
    translationKey: 'results',
    side: 'top',
    align: 'start',
  },
  {
    element: '#add-user-btn',
    translationKey: 'addNew',
    side: 'bottom',
    align: 'end',
  },
];

/**
 * Tour step for row actions demonstration
 * Shows the action menu on a demo row to explain available actions
 */
export const ROW_ACTIONS_STEP: TourStepDefinition = {
  element: '.data-row-actions-menu',
  translationKey: 'rowActions',
  side: 'left',
  align: 'start',
};

/**
 * Demo user for tour demonstration
 * This user is injected at the top of the list during the tour
 */
export const DEMO_USER: UserResponse = {
  id: 'demo-tour-user',
  firstName: 'Jane',
  lastName: 'Doe',
  email: 'jane.doe@example.com',
  enabled: true,
  jobTitle: 'Software Engineer',
  avatar: 'JD',
  groupsCount: 3,
};

/**
 * Tour transition config - navigates to create user page after last step
 */
export const USER_LIST_TOUR_TRANSITION: TourTransition = {
  targetRoute: '/users/add',
  triggerAtStep: 6, // Last step (0-indexed) - adjusted for rowActions step
};
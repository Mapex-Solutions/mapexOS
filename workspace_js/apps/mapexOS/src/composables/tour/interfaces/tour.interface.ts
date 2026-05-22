/**
 * Tour System Interfaces
 *
 * Shared interfaces for the page tour and multi-page tour system.
 * Used by usePageTour composable and page components.
 */

import type { Side, Alignment } from 'driver.js';

/**
 * Structural definition of a tour step (without text - text comes from translations)
 * Used in constants files to define element targeting and positioning
 */
export interface TourStepDefinition {
  /** CSS selector targeting the element (#id or .class) */
  element: string;

  /** Popover position relative to element */
  side?: Side;

  /** Popover alignment relative to element */
  align?: Alignment;

  /** i18n translation key suffix (e.g., 'filters', 'results') */
  translationKey: string;
}

/**
 * Complete tour step with resolved text (passed to usePageTour)
 * Translation composable resolves the text before passing to the composable
 */
export interface PageTourStep {
  /** CSS selector targeting the element */
  element: string;

  /** Resolved title text (from translations) */
  title: string;

  /** Resolved description text (from translations) */
  description: string;

  /** Popover position relative to element */
  side?: Side;

  /** Popover alignment relative to element */
  align?: Alignment;

  /** Hook called when this step's element is highlighted */
  onHighlightStarted?: () => void;

  /** Custom handler for Next button click. Receives a moveNext function to manually advance. */
  onNextClick?: (moveNext: () => void) => void;
}

/**
 * Multi-page transition configuration
 * Defines when and where to navigate during a tour
 */
export interface TourTransition {
  /** Target route path (e.g., '/users/add') */
  targetRoute: string;

  /** Step index (0-based) that triggers navigation on "Next" click */
  triggerAtStep: number;
}

/**
 * Configuration passed to usePageTour()
 */
export interface PageTourConfig {
  /** Unique identifier for this tour */
  tourId: string;

  /** Getter function that returns tour steps with resolved translations */
  steps: () => PageTourStep[];

  /** Optional multi-page transition config */
  transition?: TourTransition;

  /** Auto-start tour (default: detects via route query param) */
  autoStart?: boolean;

  /** Callback when tour ends (finish, cancel, skip, or close) */
  onTourEnd?: () => void;
}

/**
 * State saved to sessionStorage for multi-page tour transitions
 */
export interface TourSessionState {
  /** Tour flow identifier */
  flowId: string;

  /** Whether the destination page should auto-start its tour */
  shouldAutoStart: boolean;

  /** Timestamp for TTL expiration check */
  timestamp: number;
}

/**
 * Tour configuration for PageHeader component
 */
export interface PageHeaderTour {
  /** Whether tour button is enabled */
  enabled: boolean;

  /** Optional tooltip text for the tour button */
  tooltip?: string;
}

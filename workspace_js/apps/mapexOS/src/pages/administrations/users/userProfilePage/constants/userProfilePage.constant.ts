/**
 * UserProfilePage Constants
 */

import type { ProfileSection } from '../interfaces';

/**
 * Default active section on page load (1-based index)
 */
export const DEFAULT_ACTIVE_SECTION = 1;

/**
 * Total number of sections
 */
export const TOTAL_SECTIONS = 4;

/**
 * Section enum for easy reference
 */
export const SECTION = {
  PERSONAL: 1,
  PASSWORD: 2,
  GROUPS_ACCESS: 3,
  REVIEW: 4,
} as const;

/**
 * Password validation minimum length
 */
export const PASSWORD_MIN_LENGTH = 8;

/**
 * Email validation regex pattern
 */
export const EMAIL_VALIDATION_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Notification timeout duration in milliseconds
 */
export const NOTIFICATION_TIMEOUT = 2000;

/**
 * Profile sections configuration for StepperVertical
 * Uses StepperVerticalItem format (title, description, icon)
 */
export const PROFILE_SECTIONS_CONFIG = (t: any): ProfileSection[] => [
  {
    title: t.navigation.personal.label.value,
    icon: 'person',
    description: t.navigation.personal.description.value,
  },
  {
    title: t.navigation.password.label.value,
    icon: 'lock',
    description: t.navigation.password.description.value,
  },
  {
    title: t.navigation.groupsAccess.label.value,
    icon: 'shield',
    description: t.navigation.groupsAccess.description.value,
  },
  {
    title: t.navigation.review.label.value,
    icon: 'checklist',
    description: t.navigation.review.description.value,
  },
];

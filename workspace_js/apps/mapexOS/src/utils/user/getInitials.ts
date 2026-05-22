/**
 * Extract user initials from first and last name
 *
 * @param {string | undefined} firstName - User first name
 * @param {string | undefined} lastName - User last name
 * @param {string | undefined} email - Fallback email
 * @returns {string} Uppercase initials (e.g. "TA") or "?" as fallback
 */
export function getInitials(
  firstName?: string,
  lastName?: string,
  email?: string,
): string {
  const first = firstName?.charAt(0) || '';
  const last = lastName?.charAt(0) || '';

  if (first || last) {
    return (first + last).toUpperCase();
  }

  return email?.charAt(0).toUpperCase() || '?';
}

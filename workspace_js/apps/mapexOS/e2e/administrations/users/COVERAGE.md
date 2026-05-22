# E2E Test Coverage: User Registration

**Suite:** `administration.users.add`
**Page:** `/users/add` (Create User Wizard - 4 steps)
**Run:** `npx playwright test -g "administration.users.add"`

---

## Tests

| # | Test | Category | What it covers |
|---|------|----------|----------------|
| 1 | `happy-path: complete full registration flow` | Integration | Full 4-step wizard: fill personal info, security, select group via drawer, verify review, save, assert success toast + redirect to `/users` |
| 2 | `navigation: Next and Previous buttons work correctly` | Navigation | Previous hidden on Step 1, visible from Step 2+. Next visible on Steps 1-3, hidden on Step 4. Save visible only on Step 4. Data persists after Previous/Next. |
| 3 | `validation.step1: block navigation when required fields are empty` | Validation | Next disabled with no fields. Still disabled with only firstName. Still disabled with firstName+lastName. Enabled after firstName+lastName+email. |
| 4 | `validation.step2: block navigation when password is empty or too short` | Validation | Next disabled with empty password. Still disabled with short password (<8 chars). Enabled after valid password + confirmation. |
| 5 | `validation.step3: block navigation when no group is selected` | Validation | Next disabled when access type is "group" but no group has been selected from the drawer. |
| 6 | `review: edit button navigates back to the correct step` | Review | Completes all steps, verifies all 3 review sections are visible (personal, security, access). Clicks edit on personal section, verifies Step 1 renders with persisted data. |
| 7 | `tour: pre-fill form with demo data` | Tour Mode | Navigates to `/users/add?tour=true`. Closes Driver.js tour overlay. Verifies Step 1 fields are pre-filled with `DEMO_USER_FORM_DATA` (John / Doe / john.doe@example.com). |

---

## Run Commands

```bash
# All user registration tests
npx playwright test -g "administration.users.add"

# By category
npx playwright test -g "happy-path"
npx playwright test -g "navigation"
npx playwright test -g "validation.step1"
npx playwright test -g "validation.step2"
npx playwright test -g "validation.step3"
npx playwright test -g "review"
npx playwright test -g "tour"

# Headed (visible browser)
npx playwright test -g "administration.users.add" --headed --workers=1
```

---

## File Structure

```
e2e/administrations/users/
  user-registration.spec.ts   # Test cases
  user-registration.po.ts     # Page Object (UserRegistrationPage)
  user-registration.data.ts   # Test data factories
  COVERAGE.md                 # This file
```

---

## Page Object: `UserRegistrationPage`

| Method | Description |
|--------|-------------|
| `goto(options?)` | Navigate to `/users/add`. Pass `{ tour: true }` for tour mode. |
| `fillStep1(data)` | Fill first name, last name, email, phone (optional), job title (optional). |
| `fillStep2(data)` | Fill password, confirm password. Toggle change-password-next-login if set. |
| `fillStep3WithGroup()` | Select "group" access type, open group drawer, pick first group. |
| `verifyReview(data)` | Assert personal, security and access review sections are visible and contain expected data. |
| `next()` / `previous()` / `save()` | Wizard navigation buttons. |
| `isNextDisabled()` / `isSaveDisabled()` | Check button disabled state. |
| `clickEditOnReview(section)` | Click edit button on a review section (`personal`, `security`, `access`). |

---

## Components Covered

| Component | data-testid |
|-----------|-------------|
| FormCard (buttons) | `wizard-previous-btn`, `wizard-next-btn`, `wizard-save-btn` |
| Step1Personal | `user-firstname-input`, `user-lastname-input`, `user-email-input`, `user-phone-input`, `user-jobtitle-input` |
| Step2Security | `user-password-input`, `user-confirm-password-input`, `user-change-pwd-checkbox` |
| Step3Access | `user-access-type-group`, `user-group-select-btn` |
| Step4Review | `review-personal-section`, `review-security-section`, `review-access-section`, `review-edit-personal-btn`, `review-edit-security-btn`, `review-edit-access-btn` |

---

## Known Limitations

- **Happy path** creates a real user via API. Requires a running backend with valid credentials (`vendor@mapex.global`).
- **Step 3** always picks the first group from the drawer. Requires at least one group to exist in the system.
- **Review edit** only tests navigating back to Step 1. Navigating back through Steps 2-3 after editing is limited by `confirmPassword` being a local field that resets on re-mount.
- **Tour mode** requires Driver.js overlay to be dismissed before asserting field values.

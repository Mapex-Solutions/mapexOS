# E2E Test Coverage: Customer Registration

**Suite:** `administration.customers.add`
**Page:** `/customers/add` (Create Customer Wizard - 4 steps for customer/site types)
**Run:** `npx playwright test -g "administration.customers.add"`

---

## Tests

| # | Test | Category | What it covers |
|---|------|----------|----------------|
| 1 | `happy-path: complete full customer registration flow` | Integration | Full 4-step wizard: fill name+phone, fill address, select access policy (strict/local), verify review, save, assert success toast + redirect to `/customers` |
| 2 | `navigation: Next and Previous buttons work correctly` | Navigation | Previous hidden on Step 1, visible from Step 2+. Next visible on Steps 1-3, hidden on Step 4. Save visible only on Step 4. Data persists after Previous/Next. |
| 3 | `validation.step1: block navigation when name is empty or too short` | Validation | Next disabled with empty name. Still disabled with name < 3 chars. Enabled after valid name (>= 3 chars). |
| 4 | `validation.step2: allow navigation with empty address fields` | Validation | All address fields are optional — Next should be enabled even with empty fields. |
| 5 | `access-policy: card selection updates correctly` | Interaction | Click "Merge" card → rolePolicy changes. Click "Recursive" card → defaultScope changes. Verify visual selection state via CSS class. |
| 6 | `review: all sections are visible and contain expected data` | Review | Complete all steps, verify basic + address + access policy sections visible with correct data (name, country, city, policy labels). |

---

## Run Commands

```bash
# All customer registration tests
npx playwright test -g "administration.customers.add"

# By category
npx playwright test -g "happy-path"
npx playwright test -g "navigation"
npx playwright test -g "validation.step1"
npx playwright test -g "validation.step2"
npx playwright test -g "access-policy"
npx playwright test -g "review"

# Headed (visible browser)
npx playwright test -g "administration.customers.add" --headed --workers=1

# All administration tests (users + customers)
npx playwright test -g "administration"
```

---

## File Structure

```
e2e/administrations/customers/
  customer-registration.spec.ts   # Test cases (6 tests)
  customer-registration.po.ts     # Page Object (CustomerRegistrationPage)
  customer-registration.data.ts   # Test data factories
  COVERAGE.md                     # This file
```

---

## Page Object: `CustomerRegistrationPage`

| Method | Description |
|--------|-------------|
| `goto()` | Navigate to `/customers/add`. |
| `fillStep1(data)` | Fill name, phone (optional), toggle enabled if false. |
| `fillStep2(data)` | Fill address fields: country, state, city, zipCode (all optional). |
| `fillStep3(data)` | Click card for rolePolicy + defaultScope selection. |
| `verifyReview(data)` | Assert basic, address and access review sections are visible and contain expected data. |
| `next()` / `previous()` / `save()` | Wizard navigation buttons. |
| `isNextDisabled()` / `isSaveDisabled()` | Check button disabled state. |
| `clickEditOnReview(section)` | Click edit button on a review section (`basic`, `address`, `access`). |

---

## Components Covered

| Component | data-testid |
|-----------|-------------|
| FormCard (buttons) | `wizard-previous-btn`, `wizard-next-btn`, `wizard-save-btn` |
| Step1Basic | `customer-name-input`, `customer-phone-input`, `customer-enabled-checkbox` |
| Step2Address | `customer-country-input`, `customer-state-input`, `customer-city-input`, `customer-zipcode-input` |
| Step3AccessPolicy | `customer-role-policy-strict`, `customer-role-policy-merge`, `customer-scope-local`, `customer-scope-recursive` |
| Step4Review (FormReview) | `review-basic-section`, `review-address-section`, `review-access-section`, `review-basic-section-edit-btn`, `review-address-section-edit-btn`, `review-access-section-edit-btn` |

---

## Known Limitations

- **Happy path** creates a real organization via API. Requires a running backend with valid credentials (`vendor@mapex.global`).
- **Customer type only** — tests cover the 4-step flow (with address). Types without address (building, floor, zone) use a 3-step flow not covered here.
- **Access policy defaults** — "Strict" and "Local" are pre-selected by default; tests verify changing selection but don't test all combinations.
- **Review edit** is available via `clickEditOnReview()` but not explicitly tested in the current suite (covered functionally through navigation test).

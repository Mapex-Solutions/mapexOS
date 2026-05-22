# E2E Test Coverage: Asset Registration

**Suite:** `assets.add`
**Page:** `/assets/add` (Create Asset Wizard - 5 steps)
**Run:** `npx playwright test -g "assets.add"`

---

## Tests

| # | Test | Category | What it covers |
|---|------|----------|----------------|
| 1 | `happy-path: complete full asset registration flow` | Integration | Full 5-step wizard with HTTP protocol: fill identification (name, assetId, status, description), select template via drawer, select route group via drawer, set connectivity (HTTP), verify review, save, assert success toast + redirect to `/assets` |
| 2 | `navigation: Next and Previous buttons work correctly` | Navigation | Previous hidden on Step 1, visible from Step 2+. Next visible on Steps 1-4, hidden on Step 5. Save visible only on Step 5. Data persists after Previous/Next. |
| 3 | `validation.step1: block navigation when required fields are empty` | Validation | Click Next with empty form → stays on Step 1. Fill only name, click Next → stays on Step 1 (assetId required). Fill both name + assetId, click Next → advances to Step 2. |
| 4 | `validation.step2: block navigation when no template is selected` | Validation | Next disabled when no asset template is selected via drawer. |
| 5 | `validation.step3: block navigation when no route group is selected` | Validation | Next disabled when no route group is selected via drawer. |
| 6 | `review: all sections are visible and contain expected data` | Review | Complete all steps, verify identification + template + route groups + connectivity sections visible with correct data (name, assetId, protocol). |

---

## Run Commands

```bash
# All asset registration tests
npx playwright test -g "assets.add"

# By category
npx playwright test -g "happy-path"
npx playwright test -g "navigation"
npx playwright test -g "validation.step1"
npx playwright test -g "validation.step2"
npx playwright test -g "validation.step3"
npx playwright test -g "review"

# Headed (visible browser)
npx playwright test -g "assets.add" --headed --workers=1

# All suites (administration + assets)
npx playwright test -g "administration|assets"
```

---

## File Structure

```
e2e/assets/
  asset-registration.spec.ts   # Test cases (6 tests)
  asset-registration.po.ts     # Page Object (AssetRegistrationPage)
  asset-registration.data.ts   # Test data factories
  COVERAGE.md                  # This file
```

---

## Page Object: `AssetRegistrationPage`

| Method | Description |
|--------|-------------|
| `goto()` | Navigate to `/assets/add`, wait for load. |
| `fillStep1(data)` | Fill name, assetId, select status (optional), description (optional), toggle debug (optional). |
| `fillStep2WithTemplate()` | Click template input -> open drawer -> pick first template -> confirm -> wait for drawer close. |
| `fillStep3WithRouteGroups()` | Click "Select Route Groups" -> open drawer -> pick first group -> confirm -> wait for drawer close. |
| `fillStep4(data)` | Select protocol. If MQTT: fill username, password, clientId. Optionally fill lat/lng. |
| `verifyReview(data)` | Assert 4 review sections visible with expected data (name, assetId, protocol). |
| `next()` / `previous()` / `save()` | Wizard navigation buttons. |
| `isNextDisabled()` / `isSaveDisabled()` | Check button disabled state. |
| `clickEditOnReview(section)` | Click edit button on a review section (`identification`, `template`, `routegroups`, `connectivity`). |

---

## Components Covered

| Component | data-testid |
|-----------|-------------|
| FormCard (buttons) | `wizard-previous-btn`, `wizard-next-btn`, `wizard-save-btn` |
| Step1Identification | `asset-name-input`, `asset-id-input`, `asset-status-select`, `asset-description-input`, `asset-debug-toggle` |
| Step2AssetTemplate | `asset-template-input`, `asset-template-clear-btn`, `asset-template-card` |
| Step3RouteGroups | `asset-routegroup-select-btn`, `asset-routegroup-clear-btn`, `asset-routegroup-count` |
| Step4Connectivity | `asset-protocol-select`, `asset-mqtt-username-input`, `asset-mqtt-password-input`, `asset-mqtt-clientid-input`, `asset-latitude-input`, `asset-longitude-input` |
| Step5Review (FormReview) | `review-identification-section`, `review-template-section`, `review-routegroups-section`, `review-connectivity-section`, `review-*-section-edit-btn` |

---

## Known Limitations

- **Happy path** creates a real asset via API. Requires a running backend with valid credentials (`vendor@mapex.global`).
- **Template/Route Group selection** picks the first available item from the drawer — tests assume at least one template and one route group exist in the system.
- **MQTT protocol** is not tested in the current suite (only HTTP). A future test could use `createMqttAssetData()` to cover the MQTT-specific fields.
- **Review edit** is available via `clickEditOnReview()` but not explicitly tested in the current suite (covered functionally through navigation test).
- **LoRaWAN protocol** option is disabled in the UI and not tested.

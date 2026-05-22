# Phase 0 — IAM bootstrap (seed admin login)

## What this test proves

The seed admin user can sign in via mapexIam and the resulting JWT
carries org-context coverage for the seed root organization. The
ClientSet propagates that bearer + `X-Org-Context` header to every
service so every subsequent phase can drive the platform as the seeded
admin without re-issuing logins.

The phase:

1. Posts admin credentials (`admin@mapex.local`) to `/auth/login`; publishes JWT + organizationID to the bag.
2. Asserts the JWT is structurally valid (parse, signature, future expiration).
3. Asserts the JWT carries access to `MapexosOrgID` via the IAM coverage endpoint.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/
```

## Requirements

- Live `mapexos` service on `MAPEXOS_URL` (default `http://localhost:5000`).
- Seed admin user (`admin@mapex.local` / `mapex@123`) provisioned by the canonical mongodb-init seed.
- Coverage build job already produced the wildcard `mapex.*` entry for the seed admin anchored at the seed root org.

// Constants used across the organizations module e2e suite.
//
// Lives in its own file (with _test.go suffix so it only compiles for
// `go test`) so the actual test functions stay focused on assertions.
// Other modules adapting this template change these constants in one
// place rather than hunting them across helpers and tests.
package e2e

// nonExistentOrgID is a syntactically valid but never-allocated ObjectID
// used for "not found" tests. The 24-hex shape passes the validator so
// 404 reflects the persistence path, not the input shape.
const nonExistentOrgID = "0000000000000000deadbeef"

// listFixtureCount is the number of orgs created for each list test.
// 15 was picked because perPage=1 across 15 pages exposes pagination
// math edge cases (page 1, mid pages, page 15) better than 5 or 10.
const listFixtureCount = 15

// orgNamePrefix is the deterministic prefix every list-test fixture uses
// in its Name field. Sharing one prefix lets list queries scope the
// universe to "this test only" via ?name=<prefix>-<runID>, eliminating
// interference from pre-seeded data and from parallel runs.
const orgNamePrefix = "saga-customer"

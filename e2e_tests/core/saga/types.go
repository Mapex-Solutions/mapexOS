// Package saga holds the saga runner used by every E2E journey.
//
// A journey is a sequence of Items run against the live stack. Each Item is
// either a Step (mutates the system, may publish output keys to the bag,
// optionally registers a Compensate) or an Assert (read-only verification of
// state). The runner walks the list in order and, on success or failure, runs
// every registered Compensate in reverse so the environment is left clean for
// the next run.
package saga

// Item is the unit of work executed by Run. Step and Assert satisfy it.
// The runner does not distinguish between them at execution time; the
// distinction matters for rollback (only Step.Compensate is invoked).
type Item interface {
	// GetName returns a stable identifier surfaced in test logs so failures
	// point at the offending step without opening source.
	GetName() string

	// Execute performs the action (Step.Do) or the verification (Assert.Check).
	Execute(c *Context) error

	// Rollback invokes the compensation registered by the Step; Asserts are
	// no-ops because they never mutated the system.
	Rollback(c *Context) error
}

// Step represents a mutating action against the live stack. The runner
// records executed Steps and replays Compensate in reverse order on
// completion (success or failure) so resources created during the journey
// are removed before the next run.
type Step struct {
	// Name is the qualified identifier of the step (e.g.,
	// "mapexIam/organizations.CreateOrganization") logged by the runner.
	Name string

	// Do performs the mutation. Returning an error halts the journey and
	// triggers rollback of every previously executed step.
	Do func(c *Context) error

	// Compensate undoes the mutation Do performed. It is optional — steps
	// that publish to NATS or read from an external system may have nothing
	// to undo. Compensate must be idempotent: when Do never ran (because a
	// previous step failed), Compensate is still invoked and must no-op.
	Compensate func(c *Context) error
}

// GetName satisfies Item.
func (s Step) GetName() string { return s.Name }

// Execute satisfies Item by delegating to Do.
func (s Step) Execute(c *Context) error { return s.Do(c) }

// Rollback satisfies Item by delegating to Compensate when set.
func (s Step) Rollback(c *Context) error {
	if s.Compensate == nil {
		return nil
	}
	return s.Compensate(c)
}

// Assert represents a read-only verification of system state. Asserts have
// no Compensate because they never mutate the world.
type Assert struct {
	// Name is the qualified identifier of the assertion (e.g.,
	// "mapexIam/auth.AssertJwtValid") logged by the runner.
	Name string

	// Check returns nil when the expectation holds. Returning an error
	// halts the journey and triggers rollback of executed Steps.
	Check func(c *Context) error
}

// GetName satisfies Item.
func (a Assert) GetName() string { return a.Name }

// Execute satisfies Item by delegating to Check.
func (a Assert) Execute(c *Context) error { return a.Check(c) }

// Rollback is a no-op because Asserts do not mutate state.
func (a Assert) Rollback(c *Context) error { return nil }


package services

/* MatchEvaluator */

// MatchEvaluator is the stateful domain service that evaluates MatchConfig
// against arbitrary event payloads. It owns the operator vocabulary used
// for history rendering.
type MatchEvaluator struct {
	operatorTextMap map[string]string
}

/* MatchRule */

// MatchRule defines a single conditional rule for event routing evaluation.
// Independent from any specific storage schema — the application layer adapts
// external match configurations (e.g., routegroup) into this domain type.
type MatchRule struct {
	Field    string
	Operator string
	Value    interface{}
}

/* MatchConfig */

// MatchConfig defines the matching strategy evaluated by the MatchEvaluator.
// Policy "all" = AND logic; Policy "any" = OR logic.
type MatchConfig struct {
	Policy string
	Rules  []MatchRule
}

/* ConditionResult */

// ConditionResult holds the result of a single condition evaluation.
type ConditionResult struct {
	Field    string
	Operator string
	Expected interface{}
	Actual   interface{}
	Passed   bool
}

/* EvaluationResult */

// EvaluationResult holds the result of match evaluation.
type EvaluationResult struct {
	ShouldProcess bool
	Conditions    []ConditionResult
	History       []string
}

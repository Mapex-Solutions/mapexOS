package ports

import (
	domainServices "router/src/modules/events/domain/services"
)

// MatchEvaluatorPort defines the contract for evaluating match conditions.
//
// This port interface enables Hexagonal Architecture by decoupling the application
// layer from the concrete MatchEvaluator domain service, improving testability.
type MatchEvaluatorPort interface {
	// Evaluate evaluates all match rules based on policy (all/any).
	// Returns structured results with individual condition evaluations.
	Evaluate(event interface{}, matchConfig *domainServices.MatchConfig) (*domainServices.EvaluationResult, error)
}

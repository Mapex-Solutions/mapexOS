package services

import (
	"workflow/src/modules/engine/application/ports"
	"workflow/src/modules/engine/domain/evaluators"
	"workflow/src/modules/engine/domain/operators"
	"workflow/src/modules/engine/domain/operators/comparison"
	"workflow/src/modules/engine/domain/operators/datetime"
	"workflow/src/modules/engine/domain/operators/group"
	"workflow/src/modules/engine/domain/operators/stringops"
)

// Compile-time checks
var _ ports.ConditionEvaluatorPort = (*evaluators.ConditionEvaluator)(nil)
var _ ports.ValueResolverPort = (*evaluators.ValueResolver)(nil)

/*
 * New creates and provides the engine ports (ConditionEvaluator, ValueResolver).
 * Initializes the OperatorRegistry with all 22 built-in operators at startup.
 */
func New() (ports.ConditionEvaluatorPort, ports.ValueResolverPort) {
	reg := operators.NewOperatorRegistry()
	registerAllOperators(reg)

	return evaluators.NewConditionEvaluator(reg), evaluators.NewValueResolver()
}

/*
 * registerAllOperators registers all built-in operators to the registry.
 * Called once at startup — registry is immutable after initialization.
 *
 * Categories:
 * - Comparison: 6 condition + 1 between (7 total)
 * - String: 5 condition
 * - DateTime: 4 condition + 2 between (6 total)
 * - Group: 4
 */
func registerAllOperators(reg *operators.OperatorRegistry) {
	/* Comparison operators */
	reg.RegisterCondition(&comparison.EqualsOperator{})
	reg.RegisterCondition(&comparison.NotEqualsOperator{})
	reg.RegisterCondition(&comparison.GreaterThanOperator{})
	reg.RegisterCondition(&comparison.GreaterThanEqualsOperator{})
	reg.RegisterCondition(&comparison.LessThanOperator{})
	reg.RegisterCondition(&comparison.LessThanEqualsOperator{})
	reg.RegisterCondition(&comparison.BetweenOperator{})
	reg.RegisterBetween(&comparison.BetweenOperator{})

	/* String operators */
	reg.RegisterCondition(&stringops.ContainsOperator{})
	reg.RegisterCondition(&stringops.NotContainsOperator{})
	reg.RegisterCondition(&stringops.StartsWithOperator{})
	reg.RegisterCondition(&stringops.EndsWithOperator{})
	reg.RegisterCondition(&stringops.RegexOperator{})

	/* DateTime operators */
	reg.RegisterCondition(&datetime.BeforeDateOperator{})
	reg.RegisterCondition(&datetime.AfterDateOperator{})
	reg.RegisterCondition(&datetime.BeforeTimeOperator{})
	reg.RegisterCondition(&datetime.AfterTimeOperator{})
	reg.RegisterCondition(&datetime.BetweenDateOperator{})
	reg.RegisterBetween(&datetime.BetweenDateOperator{})
	reg.RegisterCondition(&datetime.BetweenTimeOperator{})
	reg.RegisterBetween(&datetime.BetweenTimeOperator{})

	/* Group operators */
	reg.RegisterGroup(&group.AndOperator{})
	reg.RegisterGroup(&group.OrOperator{})
	reg.RegisterGroup(&group.NandOperator{})
	reg.RegisterGroup(&group.NorOperator{})
}

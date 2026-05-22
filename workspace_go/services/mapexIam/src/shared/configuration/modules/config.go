package configMod

import (
	"mapexIam/src/modules/auth"
	"mapexIam/src/modules/authorization_cache"
	"mapexIam/src/modules/cache_invalidation"
	"mapexIam/src/modules/groups"
	"mapexIam/src/modules/lists"
	"mapexIam/src/modules/memberships"
	"mapexIam/src/modules/onboarding_orchestrator"
	"mapexIam/src/modules/organizations"
	"mapexIam/src/modules/roles"
	"mapexIam/src/modules/users"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Core modules (no dependencies): lists, organizations
//  2. Authorization core: roles, groups, memberships
//  3. Users module (depends on memberships for multi-tenant filtering)
//  4. Auth: authentication module
//  5. Orchestrators: onboarding_orchestrator (ALWAYS LAST!)
var Modules = []common.ModuleConfig{
	// Core modules (no dependencies on other MapexOS modules)
	{
		Name:             "lists",
		Lazy:             false,
		InitRepositories: lists.InitRepositories,
		InitServices:     lists.InitServices,
		InitInterfaces:   lists.InitInterfaces,
	},
	{
		Name:             "organizations",
		Lazy:             false,
		InitRepositories: organizations.InitRepositories,
		InitServices:     organizations.InitServices,
		InitInterfaces:   organizations.InitInterfaces,
	},

	// Shared authorization cache (used by multiple modules)
	{
		Name:             "authorization_cache",
		Lazy:             false,
		InitRepositories: authorization_cache.InitRepositories,
		InitServices:     nil, // No services - only repository for cache invalidation
		InitInterfaces:   nil, // No HTTP routes or consumers
	},

	// Authorization core modules
	{
		Name:             "roles",
		Lazy:             false,
		InitRepositories: roles.InitRepositories,
		InitServices:     roles.InitServices,
		InitInterfaces:   roles.InitInterfaces,
	},
	{
		Name:             "groups",
		Lazy:             false,
		InitRepositories: groups.InitRepositories,
		InitServices:     groups.InitServices,
		InitInterfaces:   groups.InitInterfaces,
	},
	{
		Name:             "memberships",
		Lazy:             false,
		InitRepositories: memberships.InitRepositories,
		InitServices:     memberships.InitServices,
		InitInterfaces:   memberships.InitInterfaces,
	},

	// Users module (depends on memberships for multi-tenant filtering)
	{
		Name:             "users",
		Lazy:             false,
		InitRepositories: users.InitRepositories,
		InitServices:     users.InitServices,
		InitInterfaces:   users.InitInterfaces,
	},

	// Authentication module
	{
		Name:             "auth",
		Lazy:             false,
		InitRepositories: auth.InitRepositories,
		InitServices:     auth.InitServices,
		InitInterfaces:   auth.InitInterfaces,
	},

	// Cache invalidation module (listens to events from other modules)
	{
		Name:          "cache_invalidation",
		Lazy:          false,
		InitServices:  cache_invalidation.InitServices,
		InitListeners: cache_invalidation.InitListeners,
	},

	// ORCHESTRATORS (ALWAYS LAST!)
	// These coordinate other modules and must be initialized after all dependencies
	{
		Name:             "onboarding_orchestrator",
		Lazy:             false,
		InitRepositories: nil, // Orchestrators don't have repositories
		InitServices:     onboarding_orchestrator.InitServices,
		InitInterfaces:   onboarding_orchestrator.InitInterfaces,
	},
}

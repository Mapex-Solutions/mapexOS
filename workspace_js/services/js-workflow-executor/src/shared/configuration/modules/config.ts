/**
 * Module configuration following workspace_go pattern
 * Defines the order and initialization functions for all modules
 */

import * as engineModule from '@modules/engine/module';
import * as scriptsModule from '@modules/scripts/module';
import * as eventsModule from '@modules/events/module';

export interface ModuleConfig {
	/** Module identifier */
	name: string;
	/** Reserved for future lazy loading */
	lazy: boolean;
	/** Register repositories in DI (optional) */
	initRepositories?: () => void;
	/** Register services in DI (optional) */
	initServices?: () => void;
	/** Register HTTP routes (optional) */
	initInterfaces?: () => void;
	/** Start NATS listeners (optional, may be async for stream initialization) */
	initListeners?: () => void | Promise<void>;
}

/**
 * Modules initialization order
 * Order matters — modules depending on others come after their dependencies
 *
 * engine  → Piscina worker pool (no dependencies)
 * scripts → WorkflowScriptService (depends on engine)
 * events  → NATS consumers (depends on scripts)
 */
export const Modules: ModuleConfig[] = [
	// Engine module — Piscina Worker Pool (must be first)
	{
		name: 'engine',
		lazy: false,
		initRepositories: undefined,
		initServices: engineModule.initServices,
		initInterfaces: undefined,
		initListeners: engineModule.initListeners,
	},

	// Scripts module — WorkflowScriptService (depends on engine)
	{
		name: 'scripts',
		lazy: false,
		initRepositories: undefined,
		initServices: scriptsModule.initServices,
		initInterfaces: undefined,
		initListeners: undefined,
	},

	// Events module — NATS consumers (depends on scripts)
	{
		name: 'events',
		lazy: false,
		initRepositories: undefined,
		initServices: undefined,
		initInterfaces: undefined,
		initListeners: eventsModule.initListeners,
	},
];

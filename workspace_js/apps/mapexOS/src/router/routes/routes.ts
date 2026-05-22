import type { RouteRecordRaw } from 'vue-router';

import { AllAssets, AssetsManager } from './assets';
import { Triggers, Workflows, WorkflowInstances } from './automations';
import { Customers, UserProfile, Users, Settings, Lists, Roles, Groups, AccessAudit, AllRetentionPolicies } from './administrations';
import { AnalyticDashboard } from './dashboards';
import { LakeHouse } from './lakeHouse';
import { HttpDataSources } from './datasources';
import { PageNotFound, NoOrganization, Forbidden } from './erros';
import { EventsRoutes } from './events';
import { Login } from './login';
import { RawLogs } from './logs';
import { RoutingRoutes } from './routing';

export const routes: RouteRecordRaw[] = [
	// assets
	AllAssets, AssetsManager,

	// automations
	Triggers, Workflows, WorkflowInstances,

	// administrations
	Customers, UserProfile, Users, Settings, Lists, Roles, Groups, AccessAudit, AllRetentionPolicies,

	// dashboards
	AnalyticDashboard,

	// data lake
	LakeHouse,

	// data sources
	HttpDataSources,

	// erros
	PageNotFound,
	NoOrganization,
	Forbidden,

	// login
	Login,

	// routing
	...RoutingRoutes,

	// logs
	RawLogs,

	// events
	EventsRoutes,
];
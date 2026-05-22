export const RawLogs = {
	path: '/logs',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: 'event_tracer',
			component: () => import('pages/logs/eventTracerPage/EventTracerPage.vue'),
			meta: { isProtected: true, permissions: ['events.list'] },
		},
		{
			path: 'raw_logs',
			component: () => import('pages/logs/assetRawLogsPage/AssetRawLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.raw.list'] },
		},
		{
			path: 'connectivity',
			component: () => import('pages/logs/assetConnectivityLogsPage/AssetConnectivityLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.asset_status.list'] },
		},
		{
			path: 'js_exec_logs',
			component: () => import('pages/logs/jsExecLogsPage/JsExecLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.js_executor.list'] },
		},
		{
			path: 'router_logs',
			component: () => import('pages/logs/routerLogsPage/RouterLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.router.list'] },
		},
		{
			path: 'audit_logs',
			component: () => import('pages/logs/auditLogsPage/AuditLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.audit.list'] },
		},
		{
			path: 'triggers_log',
			component: () => import('pages/logs/triggerLogsPage/TriggerLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.trigger.list'] },
		},
		{
			path: 'notifications',
			component: () => import('pages/logs/notificationsLogsPage/NotificationsLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.notifications.list'] },
		},
		{
			path: 'workflow_executions',
			component: () => import('pages/logs/workflowExecutionsPage/WorkflowExecutionsPage.vue'),
			meta: { isProtected: true, permissions: ['workflows.instances.list'] },
		},
		{
			path: 'dlq',
			component: () => import('pages/logs/dlqLogsPage/DlqLogsPage.vue'),
			meta: { isProtected: true, permissions: ['events.dlq.list'] },
		},
	],
};
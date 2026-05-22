export const AccessAudit = {
	path: '/access-audit',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'access-audit',
			component: () => import('pages/administrations/accessAudit/accessAuditPage/AccessAuditPage.vue'),
			meta: { isProtected: true, permissions: ['events.audit.list'] },
		},
	],
};

export const AllRetentionPolicies = {
	path: '/administrations/retention',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () =>
				import('pages/administrations/retentionPoliciesPage/RetentionPoliciesPage.vue'),
			meta: { isProtected: true, permissions: ['retention.list'] },
		},
	],
};

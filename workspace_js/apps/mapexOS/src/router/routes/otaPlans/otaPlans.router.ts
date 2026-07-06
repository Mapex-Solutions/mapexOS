export const OtaPlans = {
	path: '/ota_plans',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/otaPlans/otaPlans/otaPlanListPage/OtaPlanListPage.vue'),
			meta: { isProtected: true, permissions: ['ota_plans.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/otaPlans/otaPlans/createEditOtaPlanPage/CreateEditOtaPlanPage.vue'),
			meta: { isProtected: true, permissions: ['ota_plans.create'] },
		},
		{
			path: ':id',
			component: () => import('pages/otaPlans/otaPlans/otaPlanDetailPage/OtaPlanDetailPage.vue'),
			meta: { isProtected: true, permissions: ['ota_plans.read'] },
		},
	],
};

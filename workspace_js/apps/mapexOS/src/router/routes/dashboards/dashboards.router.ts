export const AnalyticDashboard = {
	path: '/home',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/dashboards/dashboardAdm/DashboardAdm.vue'),
			meta: { isProtected: true },
		},
	],
};
export const Settings = {
	path: '/admin/settings',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/administrations/settings/settingsPage/SystemSettingsPage.vue'),
			meta: { isProtected: true },
		},
	],
};

export const HttpDataSources = {
	path: '/data_sources/http',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/datasources/http/httpDataSourcesListPage/HttpDataSourcesListPage.vue'),
			meta: { isProtected: true, permissions: ['datasources.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/datasources/http/createEditHttpDataSourcePage/CreateEditHttpDataSourcePage.vue'),
			meta: { isProtected: true, permissions: ['datasources.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/datasources/http/createEditHttpDataSourcePage/CreateEditHttpDataSourcePage.vue'),
			meta: { isProtected: true, permissions: ['datasources.update'] },
		},
	],
};

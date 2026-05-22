export const Lists = {
	path: '/admin/lists',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/administrations/lists/listsListPage/ListsListPage.vue'),
			meta: { isProtected: true, permissions: ['lists.lists'] },
		},
	],
};

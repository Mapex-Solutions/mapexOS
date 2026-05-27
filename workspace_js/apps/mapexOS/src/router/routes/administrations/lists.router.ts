export const Lists = {
	path: '/admin/lists',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'lists-list',
			component: () => import('pages/administrations/lists/listsListPage/ListsListPage.vue'),
			meta: { isProtected: true, permissions: ['lists.lists'] },
		},
		{
			path: 'add',
			name: 'lists-add',
			component: () => import('pages/administrations/lists/createEditListPage/CreateEditListPage.vue'),
			meta: { isProtected: true, permissions: ['lists.create'] },
		},
		{
			path: 'edit/:id',
			name: 'lists-edit',
			component: () => import('pages/administrations/lists/createEditListPage/CreateEditListPage.vue'),
			meta: { isProtected: true, permissions: ['lists.update'] },
		},
	],
};

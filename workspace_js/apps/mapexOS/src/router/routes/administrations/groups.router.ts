export const Groups = {
	path: '/groups',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'groups-list',
			component: () => import('pages/administrations/groups/groupsListPage/GroupsListPage.vue'),
			meta: { isProtected: true, permissions: ['groups.list'] },
		},
		{
			path: 'add',
			name: 'groups-add',
			component: () => import('pages/administrations/groups/createEditGroupPage/CreateEditGroupPage.vue'),
			meta: { isProtected: true, permissions: ['groups.create'] },
		},
		{
			path: 'edit/:id',
			name: 'groups-edit',
			component: () => import('pages/administrations/groups/createEditGroupPage/CreateEditGroupPage.vue'),
			meta: { isProtected: true, permissions: ['groups.update'] },
		},
		{
			path: 'detail/:id',
			name: 'groups-detail',
			component: () => import('pages/administrations/groups/groupDetailPage/GroupDetailPage.vue'),
			meta: { isProtected: true, permissions: ['groups.read'] },
		},
	],
};

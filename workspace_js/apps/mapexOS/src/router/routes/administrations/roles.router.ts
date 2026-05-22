export const Roles = {
	path: '/roles',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'roles-list',
			component: () => import('pages/administrations/roles/rolesListPage/RolesListPage.vue'),
			meta: { isProtected: true, permissions: ['roles.list'] },
		},
		{
			path: 'add',
			name: 'roles-add',
			component: () => import('pages/administrations/roles/createEditRolePage/CreateEditRolePage.vue'),
			meta: { isProtected: true, permissions: ['roles.create'] },
		},
		{
			path: 'edit/:id',
			name: 'roles-edit',
			component: () => import('pages/administrations/roles/createEditRolePage/CreateEditRolePage.vue'),
			meta: { isProtected: true, permissions: ['roles.update'] },
		},
	],
};

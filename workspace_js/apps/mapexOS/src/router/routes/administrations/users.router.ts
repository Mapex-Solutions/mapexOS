export const Users = {
	path: '/users',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'users-list',
			component: () => import('pages/administrations/users/userListPage/UserListPage.vue'),
			meta: { isProtected: true, permissions: ['users.list'] },
		},
		{
			path: 'add',
			name: 'users-add',
			component: () => import('pages/administrations/users/createEditUserPage/CreateEditUserPage.vue'),
			meta: { isProtected: true, permissions: ['users.create'] },
		},
		{
			path: 'edit/:id',
			name: 'users-edit',
			component: () => import('pages/administrations/users/createEditUserPage/CreateEditUserPage.vue'),
			meta: { isProtected: true, permissions: ['users.update'] },
		},
		{
			path: 'detail/:id',
			name: 'users-detail',
			component: () => import('pages/administrations/users/userDetailPage/UserDetailPage.vue'),
			meta: { isProtected: true, permissions: ['users.read'] },
		},
	],
};

export const UserProfile = {
	path: '/my_profile',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/administrations/users/userProfilePage/UserProfilePage.vue'),
			meta: { isProtected: true },
		},
	],
};
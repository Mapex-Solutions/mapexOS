export const Login = {
	path: '/',
	component: () => import('layouts/login/LoginLayout.vue'),
	meta: { isPublic: true },
	children: [
		{
			path: '',
			component: () => import('pages/login/LoginPage.vue'),
			meta: { isPublic: true },
		},
		{
			path: 'change-password',
			name: 'change-password',
			component: () => import('pages/changePassword/ChangePasswordPage.vue'),
			meta: { isPublic: true },
		},
	],
};
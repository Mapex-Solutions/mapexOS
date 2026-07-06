export const TemplateMigrations = {
	path: '/assets_template/migrations',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/assets/templateMigrations/migrationPlansListPage/MigrationPlansListPage.vue'),
			meta: { isProtected: true, permissions: ['templatemigrations.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/assets/templateMigrations/createEditMigrationPlanPage/CreateEditMigrationPlanPage.vue'),
			meta: { isProtected: true, permissions: ['templatemigrations.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/assets/templateMigrations/createEditMigrationPlanPage/CreateEditMigrationPlanPage.vue'),
			meta: { isProtected: true, permissions: ['templatemigrations.update'] },
		},
		{
			path: ':id',
			component: () => import('pages/assets/templateMigrations/migrationPlanDetailPage/MigrationPlanDetailPage.vue'),
			meta: { isProtected: true, permissions: ['templatemigrations.read'] },
		},
	],
};

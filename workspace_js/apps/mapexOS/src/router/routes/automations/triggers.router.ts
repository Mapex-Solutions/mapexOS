export const Triggers = {
	path: '/triggers',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/automations/triggers/triggerListPage/TriggerListPage.vue'),
			meta: { isProtected: true, permissions: ['triggers.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/automations/triggers/createEditTriggerPage/CreateEditTriggerPage.vue'),
			meta: { isProtected: true, permissions: ['triggers.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/automations/triggers/createEditTriggerPage/CreateEditTriggerPage.vue'),
			meta: { isProtected: true, permissions: ['triggers.update'] },
		},
	],
};

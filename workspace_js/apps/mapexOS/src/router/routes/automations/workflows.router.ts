export const Workflows = {
	path: '/workflows',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/automations/workflows/workflowListPage/WorkflowListPage.vue'),
			meta: { isProtected: true, permissions: ['workflows.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/automations/workflows/createEditWorkflowPage/CreateEditWorkflowPage.vue'),
			meta: { isProtected: true, permissions: ['workflows.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/automations/workflows/createEditWorkflowPage/CreateEditWorkflowPage.vue'),
			meta: { isProtected: true, permissions: ['workflows.update'] },
		},
	],
};

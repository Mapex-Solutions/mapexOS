export const WorkflowInstances = {
	path: '/workflow_instances',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/automations/workflowInstances/workflowInstanceListPage/WorkflowInstanceListPage.vue'),
			meta: { isProtected: true, permissions: ['workflows.instances.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/automations/workflowInstances/createEditWorkflowInstancePage/CreateEditWorkflowInstancePage.vue'),
			meta: { isProtected: true, permissions: ['workflows.instances.create'] },
		},
		{
			path: ':id',
			component: () => import('pages/automations/workflowInstances/createEditWorkflowInstancePage/CreateEditWorkflowInstancePage.vue'),
			meta: { isProtected: true, permissions: ['workflows.instances.read'] },
		},
	],
};

export const RouteGroupsRouter = {
	path: '/routing/route_groups',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/routing/routeGroups/routeGroupsListPage/RouteGroupsListPage.vue'),
			meta: { isProtected: true, permissions: ['routegroups.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/routing/routeGroups/createEditRouteGroupPage/CreateEditRouteGroupPage.vue'),
			meta: { isProtected: true, permissions: ['routegroups.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/routing/routeGroups/createEditRouteGroupPage/CreateEditRouteGroupPage.vue'),
			meta: { isProtected: true, permissions: ['routegroups.update'] },
		},
	],
};

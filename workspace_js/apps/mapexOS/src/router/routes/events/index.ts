export const EventsRoutes = {
	path: '/events',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: 'store',
			component: () => import('pages/events/eventStorePage/EventStorePage.vue'),
			meta: { isProtected: true, permissions: ['events.list'] },
		},
	],
};

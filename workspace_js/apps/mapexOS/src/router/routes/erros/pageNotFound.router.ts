export const PageNotFound = {
	path: '/:catchAll(.*)*',
	component: () => import('pages/errors/ErrorNotFound.vue'),
	meta: { isPublic: true },
};
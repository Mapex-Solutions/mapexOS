export const LakeHouse = {
	path: '/lakehouse',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/lakeHouse/lakeHouseListPage/LakeHouseListPage.vue'),
			meta: { isProtected: true, permissions: ['retention.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/lakeHouse/addLakeHouse/AddLakeHouse.vue'),
			meta: { isProtected: true, permissions: ['retention.update'] },
		},
	],
};
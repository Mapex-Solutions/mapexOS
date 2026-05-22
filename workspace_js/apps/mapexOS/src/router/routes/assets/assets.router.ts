export const AllAssets = {
	path: '/assets',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/assets/assets/assetListPage/AssetsListPage.vue'),
			meta: { isProtected: true, permissions: ['assets.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/assets/assets/createEditAssetPage/CreateEditAssetPage.vue'),
			meta: { isProtected: true, permissions: ['assets.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/assets/assets/createEditAssetPage/CreateEditAssetPage.vue'),
			meta: { isProtected: true, permissions: ['assets.update'] },
		},
	],
};
export const AssetTemplateMarketplace = {
	path: '/assets_template/marketplace',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/assets/assetTemplates/assetTemplateMarketplacePage/AssetTemplateMarketplacePage.vue'),
			meta: { isProtected: true, permissions: ['assettemplates.list'] },
		},
	],
};

export const AssetsManager = {
	path: '/assets_template',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			component: () => import('pages/assets/assetTemplates/assetTemplateListPage/AssetTemplateListPage.vue'),
			meta: { isProtected: true, permissions: ['assettemplates.list'] },
		},
		{
			path: 'add',
			component: () => import('pages/assets/assetTemplates/createEditAssetTemplatePage/CreateEditAssetTemplatePage.vue'),
			meta: { isProtected: true, permissions: ['assettemplates.create'] },
		},
		{
			path: 'edit/:id',
			component: () => import('pages/assets/assetTemplates/createEditAssetTemplatePage/CreateEditAssetTemplatePage.vue'),
			meta: { isProtected: true, permissions: ['assettemplates.update'] },
		},
	],
};
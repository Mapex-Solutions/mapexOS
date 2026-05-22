export const Customers = {
	path: '/customers',
	component: () => import('layouts/main/MainLayout.vue'),
	meta: { isProtected: true },
	children: [
		{
			path: '',
			name: 'customers-list',
			component: () => import('pages/administrations/customers/customerListPage/CustomersListPage.vue'),
			meta: { isProtected: true, permissions: ['organizations.list'] },
		},
		{
			path: 'add',
			name: 'customers-add',
			component: () => import('pages/administrations/customers/createEditCustomerPage/CreateEditCustomerPage.vue'),
			meta: { isProtected: true, permissions: ['organizations.create'] },
			// Supports query params: ?parentId=<orgId>
			// Child type is auto-derived from parent type via CHILD_TYPE_MAP
		},
		{
			path: 'edit/:id',
			name: 'customers-edit',
			component: () => import('pages/administrations/customers/createEditCustomerPage/CreateEditCustomerPage.vue'),
			meta: { isProtected: true, permissions: ['organizations.update'] },
		},
	],
};
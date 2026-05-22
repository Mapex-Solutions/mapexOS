export const NoOrganization = {
  path: '/errors/no-organization',
  component: () => import('layouts/login/LoginLayout.vue'),
  meta: { isPublic: true },
  children: [
    {
      path: '',
      component: () => import('pages/errors/NoOrganization.vue'),
      meta: { isPublic: true },
    },
  ],
};

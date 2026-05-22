export const Forbidden = {
  path: '/errors/forbidden',
  component: () => import('layouts/login/LoginLayout.vue'),
  meta: { isPublic: true },
  children: [
    {
      path: '',
      component: () => import('pages/errors/ForbiddenPage.vue'),
      meta: { isPublic: true },
    },
  ],
};

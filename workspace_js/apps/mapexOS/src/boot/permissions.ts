import { defineBoot } from '#q-app/wrappers';
import { vPermission } from '@src/directives/permission.directive';

export default defineBoot(({ app }) => {
  app.directive('permission', vPermission);
});

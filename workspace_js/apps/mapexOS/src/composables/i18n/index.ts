/**
 * i18n Composables
 *
 * Centralized translations organized by scope:
 * - common: Reusable translations (actions, labels, messages)
 * - pages: Page-specific translations (mirrors pages/ structure)
 *
 * Structure pattern:
 * - Page file: src/pages/{section}/{feature}/{page}/ComponentPage.vue
 * - JSON file: src/i18n/{locale}/pages/{section}/{feature}.json
 * - Composable: src/composables/i18n/pages/{section}/{feature}/useFeatureTranslations.ts
 *
 * @example
 * ```ts
 * // Import common translations
 * import { useCommonActions, useCommonLabels } from '@composables/i18n';
 *
 * // Import page-specific translations
 * import { useSettingsTranslations } from '@composables/i18n';
 * ```
 */

// Common composables
export * from './common';

// Layout composables
export { useMainLayoutTranslations } from './layout/useMainLayoutTranslations';

// Composable translations
export { useOnboardingTranslations } from './composables/useOnboardingTranslations';

// Component composables
export { useListFilterTranslations } from './components/useListFilterTranslations';
export { useOrganizationTreeDrawerTranslations } from './components/useOrganizationTreeDrawerTranslations';
export { useListDrawerTranslations } from './components/drawers/listDrawer';

// Page composables (following pages/ directory structure)
// Assets
export { useAssetsTranslations } from './pages/assets/assets/useAssetsTranslations';
export { useAssetTemplatesTranslations } from './pages/assets/assetTemplates/useAssetTemplatesTranslations';
export { useAddAssetTranslations } from './pages/assets/addAsset/useAddAssetTranslations';
export { useAddAssetTemplateTranslations } from './pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

// DataSources
export { useHttpDataSourcesTranslations } from './pages/datasources/http/useHttpDataSourcesTranslations';

// Administrations
export { useSettingsTranslations } from './pages/administrations/settings/useSettingsTranslations';
export { useUsersTranslations } from './pages/administrations/users/useUsersTranslations';
export { useAddUserTranslations } from './pages/administrations/users/useAddUserTranslations';
export { useUserProfileTranslations } from './pages/administrations/users/useUserProfileTranslations';
export { useUserDetailTranslations } from './pages/administrations/users/useUserDetailTranslations';
export { useCustomersTranslations } from './pages/administrations/customers/useCustomersTranslations';
export { useAddCustomerTranslations } from './pages/administrations/customers/useAddCustomerTranslations';
export { useListsTranslations } from './pages/administrations/lists/useListsTranslations';
export { useGroupsTranslations } from './pages/administrations/groups/useGroupsTranslations';
export { useGroupDetailTranslations } from './pages/administrations/groups/useGroupDetailTranslations';
export { useRolesTranslations } from './pages/administrations/roles/useRolesTranslations';
export { useAccessAuditTranslations } from './pages/administrations/accessAudit/useAccessAuditTranslations';

// Login, Change Password & Error pages
export { useLoginTranslations } from './pages/login/useLoginTranslations';
export { useChangePasswordTranslations } from './pages/changePassword';
export { useErrorTranslations } from './pages/error/useErrorTranslations';
export { useNoOrganizationTranslations } from './pages/errors/useNoOrganizationTranslations';

// Routing
export { useRouteGroupsTranslations } from './pages/routing/routeGroups/useRouteGroupsTranslations';

// Logs
export { useNotificationsLogsPageTranslations } from './pages/logs/notificationsLogsPage';
export { useAuditLogsPageTranslations } from './pages/logs/auditLogsPage';

// Automations
export { useWorkflowListPageTranslations } from './pages/automations/workflows/workflowListPage';
export { useWorkflowInstanceListPageTranslations } from './pages/automations/workflowInstances/workflowInstanceListPage';

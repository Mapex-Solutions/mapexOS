import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Migration plan detail page translations.
 *
 * Structure mirrors:
 * - Page: src/pages/assets/templateMigrations/migrationPlanDetailPage/MigrationPlanDetailPage.vue
 * - JSON: src/i18n/{locale}/pages/templateMigrations.json (detail section)
 */
export function useMigrationPlanDetailTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  const base = 'pages.templateMigrations.detail';
  const statusesBase = 'pages.templateMigrations.statuses';

  return {
    page: {
      back: computed(() => ts(`${base}.page.back`)),
      loadError: computed(() => ts(`${base}.page.loadError`)),
      notFound: computed(() => ts(`${base}.page.notFound`)),
    },

    header: {
      route: computed(() => ts(`${base}.header.route`)),
      schedule: computed(() => ts(`${base}.header.schedule`)),
      immediate: computed(() => ts(`${base}.header.immediate`)),
      created: computed(() => ts(`${base}.header.created`)),
      updated: computed(() => ts(`${base}.header.updated`)),
    },

    sections: {
      overview: computed(() => tsTitle(`${base}.sections.overview`)),
      overviewSubtitle: computed(() => ts(`${base}.sections.overviewSubtitle`)),
      migration: computed(() => tsTitle(`${base}.sections.migration`)),
      status: computed(() => tsTitle(`${base}.sections.status`)),
      templates: computed(() => tsTitle(`${base}.sections.templates`)),
      schedule: computed(() => tsTitle(`${base}.sections.schedule`)),
    },

    fields: {
      status: computed(() => ts(`${base}.fields.status`)),
      fromTemplate: computed(() => ts(`${base}.fields.fromTemplate`)),
      toTemplate: computed(() => ts(`${base}.fields.toTemplate`)),
    },

    progress: {
      title: computed(() => tsTitle(`${base}.progress.title`)),
      total: computed(() => ts(`${base}.progress.total`)),
      migrated: computed(() => ts(`${base}.progress.migrated`)),
      failed: computed(() => ts(`${base}.progress.failed`)),
    },

    executions: {
      title: computed(() => tsTitle(`${base}.executions.title`)),
      itemLabel: computed(() => tsRaw(`${base}.executions.itemLabel`)),
      itemLabelPlural: computed(() => tsRaw(`${base}.executions.itemLabelPlural`)),
      filterStatus: computed(() => ts(`${base}.executions.filterStatus`)),
      allStatus: computed(() => ts(`${base}.executions.allStatus`)),
      empty: computed(() => ts(`${base}.executions.empty`)),
      emptyTitle: computed(() => ts(`${base}.executions.emptyTitle`)),
      resetFilter: computed(() => ts(`${base}.executions.resetFilter`)),
      columns: {
        assetId: computed(() => ts(`${base}.executions.columns.assetId`)),
        status: computed(() => ts(`${base}.executions.columns.status`)),
        error: computed(() => ts(`${base}.executions.columns.error`)),
        attempts: computed(() => ts(`${base}.executions.columns.attempts`)),
        updated: computed(() => ts(`${base}.executions.columns.updated`)),
      },
    },

    actions: {
      edit: computed(() => ts(`${base}.actions.edit`)),
      cancel: computed(() => ts(`${base}.actions.cancel`)),
      refresh: computed(() => ts(`${base}.actions.refresh`)),
    },

    cancelDialog: {
      title: computed(() => ts(`${base}.cancelDialog.title`)),
      message: computed(() => ts(`${base}.cancelDialog.message`)),
      confirm: computed(() => ts(`${base}.cancelDialog.confirm`)),
      success: computed(() => ts(`${base}.cancelDialog.success`)),
      error: computed(() => ts(`${base}.cancelDialog.error`)),
    },

    errors: {
      apiNotInitialized: computed(() => ts(`${base}.errors.apiNotInitialized`)),
    },

    /** Plan lifecycle status label, keyed by the wire status value (reuses the list statuses) */
    planStatusLabel: (status: string) => ts(`${statusesBase}.${status}`),

    /** Execution status label, keyed by the wire status value */
    executionStatusLabel: (status: string) => ts(`${base}.executions.statuses.${status}`),
  };
}

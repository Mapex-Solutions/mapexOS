import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Event Tracer Page
 * Provides all translated strings for the event tracer interface
 * @returns {Object} Translation object with page header, search, phases, stages, etc.
 */
export function useEventTracerPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.eventTracerPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.eventTracerPage.pageHeader.description')),
    },
    search: {
      label: computed(() => ts('pages.logs.eventTracerPage.search.label')),
      placeholder: computed(() => tsRaw('pages.logs.eventTracerPage.search.placeholder')),
      button: computed(() => ts('pages.logs.eventTracerPage.search.button')),
      hint: computed(() => tsRaw('pages.logs.eventTracerPage.search.hint')),
    },
    filters: {
      label: computed(() => ts('pages.logs.eventTracerPage.filters.label')),
      searchPlaceholder: computed(() => tsRaw('pages.logs.eventTracerPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.eventTracerPage.filters.advancedFilters')),
      pendingFilters: computed(() => tsRaw('pages.logs.eventTracerPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.eventTracerPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.eventTracerPage.filters.allStatus')),
      eventTrackerId: computed(() => ts('pages.logs.eventTracerPage.filters.eventTrackerId')),
      dateRange: computed(() => ts('pages.logs.eventTracerPage.filters.dateRange')),
      includeChildren: computed(() => ts('pages.logs.eventTracerPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.eventTracerPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.eventTracerPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.eventTracerPage.filters.options.no')),
      },
    },
    phases: {
      ingestion: computed(() => ts('pages.logs.eventTracerPage.phases.ingestion')),
      routing: computed(() => ts('pages.logs.eventTracerPage.phases.routing')),
    },
    stages: {
      raw: computed(() => ts('pages.logs.eventTracerPage.stages.raw')),
      jsExec: computed(() => ts('pages.logs.eventTracerPage.stages.jsExec')),
      router: computed(() => ts('pages.logs.eventTracerPage.stages.router')),
      directTrigger: computed(() => ts('pages.logs.eventTracerPage.stages.directTrigger')),
      trigger: computed(() => ts('pages.logs.eventTracerPage.stages.trigger')),
      saveEvent: computed(() => ts('pages.logs.eventTracerPage.stages.saveEvent')),
      lakeHouse: computed(() => ts('pages.logs.eventTracerPage.stages.lakeHouse')),
    },
    stageStatus: {
      success: computed(() => ts('pages.logs.eventTracerPage.stageStatus.success')),
      failed: computed(() => ts('pages.logs.eventTracerPage.stageStatus.failed')),
      noData: computed(() => ts('pages.logs.eventTracerPage.stageStatus.noData')),
      notImplemented: computed(() => ts('pages.logs.eventTracerPage.stageStatus.notImplemented')),
      matched: computed(() => ts('pages.logs.eventTracerPage.stageStatus.matched')),
      notMatched: computed(() => ts('pages.logs.eventTracerPage.stageStatus.notMatched')),
    },
    routing: {
      destinations: computed(() => ts('pages.logs.eventTracerPage.routing.destinations')),
      routers: computed(() => ts('pages.logs.eventTracerPage.routing.routers')),
      matched: computed(() => ts('pages.logs.eventTracerPage.routing.matched')),
      published: computed(() => ts('pages.logs.eventTracerPage.routing.published')),
      conditions: computed(() => ts('pages.logs.eventTracerPage.routing.conditions')),
    },
    timeline: {
      title: computed(() => tsTitle('pages.logs.eventTracerPage.timeline.title')),
      executionTime: computed(() => ts('pages.logs.eventTracerPage.timeline.executionTime')),
      timestamp: computed(() => ts('pages.logs.eventTracerPage.timeline.timestamp')),
      viewDetails: computed(() => ts('pages.logs.eventTracerPage.timeline.viewDetails')),
      totalTime: computed(() => ts('pages.logs.eventTracerPage.timeline.totalTime')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.eventTracerPage.drawer.title')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.eventTracerPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.eventTracerPage.empty.description')),
      searchFirst: computed(() => ts('pages.logs.eventTracerPage.empty.searchFirst')),
      searchFirstDescription: computed(() => tsRaw('pages.logs.eventTracerPage.empty.searchFirstDescription')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.eventTracerPage.messages.loadFailed')),
      invalidUuid: computed(() => tsRaw('pages.logs.eventTracerPage.messages.invalidUuid')),
      minLength: computed(() => tsRaw('pages.logs.eventTracerPage.messages.minLength')),
      loading: computed(() => tsRaw('pages.logs.eventTracerPage.messages.loading')),
      fetchingRaw: computed(() => tsRaw('pages.logs.eventTracerPage.messages.fetchingRaw')),
      fetchingJsExec: computed(() => tsRaw('pages.logs.eventTracerPage.messages.fetchingJsExec')),
      fetchingRouter: computed(() => tsRaw('pages.logs.eventTracerPage.messages.fetchingRouter')),
      fetchingTriggers: computed(() => tsRaw('pages.logs.eventTracerPage.messages.fetchingTriggers')),
    },
    summary: {
      title: computed(() => tsTitle('pages.logs.eventTracerPage.summary.title')),
      stagesProcessed: computed(() => ts('pages.logs.eventTracerPage.summary.stagesProcessed')),
      totalExecutionTime: computed(() => ts('pages.logs.eventTracerPage.summary.totalExecutionTime')),
      status: computed(() => ts('pages.logs.eventTracerPage.summary.status')),
      allSuccess: computed(() => ts('pages.logs.eventTracerPage.summary.allSuccess')),
      hasFailed: computed(() => tsRaw('pages.logs.eventTracerPage.summary.hasFailed')),
      totalTriggers: computed(() => ts('pages.logs.eventTracerPage.summary.totalTriggers')),
    },
    actions: {
      trackEvent: computed(() => ts('pages.logs.eventTracerPage.actions.trackEvent')),
      viewJson: computed(() => ts('pages.logs.eventTracerPage.actions.viewJson')),
      copyId: computed(() => ts('pages.logs.eventTracerPage.actions.copyId')),
      refresh: computed(() => ts('pages.logs.eventTracerPage.actions.refresh')),
    },
  };
}

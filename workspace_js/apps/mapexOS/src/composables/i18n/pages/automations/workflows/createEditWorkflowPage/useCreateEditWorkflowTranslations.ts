import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Maps validation i18n keys to their interpolation parameter name.
 * Only keys that contain a dynamic parameter need an entry here.
 */
const VALIDATION_PARAM_MAP: Record<string, string> = {
  variableIncomplete: 'path',
  fieldSourceIncomplete: 'name',
  comparisonValueIncomplete: 'name',
  caseNeedsCondition: 'caseName',
  propRequired: 'propName',
  nodesHaveErrors: 'count',
  gotoSenderNeedsReceiver: 'label',
  gotoReceiverNeedsSender: 'label',
};

/**
 * Translations composable for the Create/Edit Workflow page and all its sub-components
 */
export function useCreateEditWorkflowTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  const basePath = 'pages.automations.createEditWorkflow';

  /**
   * Resolve a raw validation error key into a translated string.
   * Supports parameterized keys using `::` separator (e.g., `'variableIncomplete::myVar'`).
   *
   * @param {string} errorKey - Raw validation error key from validators
   * @returns {string} Translated error message
   */
  function resolveValidationError(errorKey: string): string {
    const separatorIndex = errorKey.indexOf('::');

    if (separatorIndex !== -1) {
      const key = errorKey.substring(0, separatorIndex);
      const param = errorKey.substring(separatorIndex + 2);
      const paramName = VALIDATION_PARAM_MAP[key];

      if (paramName) {
        return ts(`${basePath}.validation.${key}`, { [paramName]: param });
      }
    }

    return ts(`${basePath}.validation.${errorKey}`);
  }

  return {
    resolveValidationError,
    page: {
      title: computed(() => tsTitle(`${basePath}.page.title`)),
      titleEdit: computed(() => tsTitle(`${basePath}.page.titleEdit`)),
      description: computed(() => ts(`${basePath}.page.description`)),
      descriptionEdit: computed(() => ts(`${basePath}.page.descriptionEdit`)),
      loading: computed(() => ts(`${basePath}.page.loading`)),
      back: computed(() => tsTitle(`${basePath}.page.back`)),
    },
    tabs: {
      general: computed(() => tsTitle(`${basePath}.tabs.general`)),
      data: computed(() => tsTitle(`${basePath}.tabs.data`)),
      state: computed(() => tsTitle(`${basePath}.tabs.state`)),
      workflow: computed(() => tsTitle(`${basePath}.tabs.workflow`)),
      jsonDebug: computed(() => tsTitle(`${basePath}.tabs.jsonDebug`)),
      plugins: computed(() => tsTitle(`${basePath}.tabs.plugins`)),
    },
    buttons: {
      save: computed(() => tsTitle(`${basePath}.buttons.save`)),
      update: computed(() => tsTitle(`${basePath}.buttons.update`)),
      validate: computed(() => tsTitle(`${basePath}.buttons.validate`)),
    },
    statusBanner: {
      pluginMissing: computed(() => ts(`${basePath}.statusBanner.pluginMissing`)),
    },
    tooltips: {
      validate: computed(() => ts(`${basePath}.tooltips.validate`)),
      saveReady: computed(() => ts(`${basePath}.tooltips.saveReady`)),
      saveDisabled: computed(() => ts(`${basePath}.tooltips.saveDisabled`)),
      unsavedChanges: computed(() => ts(`${basePath}.tooltips.unsavedChanges`)),
    },
    validation: {
      nameRequired: computed(() => ts(`${basePath}.validation.nameRequired`)),
      nodesRequired: computed(() => ts(`${basePath}.validation.nodesRequired`)),
      nodesHaveErrors: computed(() => ts(`${basePath}.validation.nodesHaveErrors`)),
      fixErrors: computed(() => ts(`${basePath}.validation.fixErrors`)),
      nameIsRequired: computed(() => ts(`${basePath}.validation.nameIsRequired`)),
      mustBe1to10: computed(() => tsRaw(`${basePath}.validation.mustBe1to10`)),
      mustBe1to10float: computed(() => tsRaw(`${basePath}.validation.mustBe1to10float`)),
      fieldNameRequired: computed(() => ts(`${basePath}.validation.fieldNameRequired`)),
    },
    notifications: {
      apiNotAvailable: computed(() => ts(`${basePath}.notifications.apiNotAvailable`)),
      loadFailed: computed(() => ts(`${basePath}.notifications.loadFailed`)),
      savedSuccess: computed(() => ts(`${basePath}.notifications.savedSuccess`)),
      updatedSuccess: computed(() => ts(`${basePath}.notifications.updatedSuccess`)),
      saveFailed: computed(() => ts(`${basePath}.notifications.saveFailed`)),
      jsonCopied: computed(() => ts(`${basePath}.notifications.jsonCopied`)),
      copyFailed: computed(() => ts(`${basePath}.notifications.copyFailed`)),
      workflowUpdatedFromJson: computed(() => ts(`${basePath}.notifications.workflowUpdatedFromJson`)),
      invalidWorkflowName: computed(() => ts(`${basePath}.notifications.invalidWorkflowName`)),
      invalidWorkflowNodes: computed(() => ts(`${basePath}.notifications.invalidWorkflowNodes`)),
      invalidWorkflowEdges: computed(() => ts(`${basePath}.notifications.invalidWorkflowEdges`)),
    },
    dialogs: {
      unsavedChangesTitle: computed(() => tsTitle(`${basePath}.dialogs.unsavedChangesTitle`)),
      unsavedChangesCancel: computed(() => ts(`${basePath}.dialogs.unsavedChangesCancel`)),
      unsavedChangesLeave: computed(() => ts(`${basePath}.dialogs.unsavedChangesLeave`)),
    },
    generalTab: {
      basicInfo: computed(() => tsTitle(`${basePath}.generalTab.basicInfo`)),
      name: computed(() => tsTitle(`${basePath}.generalTab.name`)),
      description: computed(() => tsTitle(`${basePath}.generalTab.description`)),
      status: computed(() => tsTitle(`${basePath}.generalTab.status`)),
      enabled: computed(() => tsTitle(`${basePath}.generalTab.enabled`)),
      disabled: computed(() => tsTitle(`${basePath}.generalTab.disabled`)),
      isTemplate: computed(() => ts(`${basePath}.generalTab.isTemplate`)),
      sharedWithChildren: computed(() => ts(`${basePath}.generalTab.sharedWithChildren`)),
      timezone: computed(() => tsTitle(`${basePath}.generalTab.timezone`)),
      timezoneType: computed(() => tsTitle(`${basePath}.generalTab.timezoneType`)),
      timezoneValue: computed(() => tsTitle(`${basePath}.generalTab.timezoneValue`)),
      timezoneHintIana: computed(() => tsRaw(`${basePath}.generalTab.timezoneHintIana`)),
      timezoneHintVariable: computed(() => tsRaw(`${basePath}.generalTab.timezoneHintVariable`)),
      timezoneInfo: computed(() => ts(`${basePath}.generalTab.timezoneInfo`)),
      retryPolicy: computed(() => tsTitle(`${basePath}.generalTab.retryPolicy`)),
      retryDescription: computed(() => ts(`${basePath}.generalTab.retryDescription`)),
      maxAttempts: computed(() => tsTitle(`${basePath}.generalTab.maxAttempts`)),
      initialInterval: computed(() => tsTitle(`${basePath}.generalTab.initialInterval`)),
      initialIntervalHint: computed(() => tsRaw(`${basePath}.generalTab.initialIntervalHint`)),
      backoffMultiplier: computed(() => tsTitle(`${basePath}.generalTab.backoffMultiplier`)),
      maxInterval: computed(() => tsTitle(`${basePath}.generalTab.maxInterval`)),
      maxIntervalHint: computed(() => tsRaw(`${basePath}.generalTab.maxIntervalHint`)),
      nonRetryableErrors: computed(() => tsTitle(`${basePath}.generalTab.nonRetryableErrors`)),
      nonRetryableErrorsHint: computed(() => ts(`${basePath}.generalTab.nonRetryableErrorsHint`)),
      retryInfo: computed(() => tsRaw(`${basePath}.generalTab.retryInfo`)),
    },
    variablesTab: {
      inputsTab: computed(() => tsTitle(`${basePath}.variablesTab.inputsTab`)),
      stateTab: computed(() => tsTitle(`${basePath}.variablesTab.stateTab`)),
      captureFieldsTab: computed(() => tsTitle(`${basePath}.variablesTab.captureFieldsTab`)),
      signalsTab: computed(() => tsTitle(`${basePath}.variablesTab.signalsTab`)),
      inputsHelp: computed(() => ts(`${basePath}.variablesTab.inputsHelp`)),
      stateHelp: computed(() => ts(`${basePath}.variablesTab.stateHelp`)),
      captureFieldsHelp: computed(() => ts(`${basePath}.variablesTab.captureFieldsHelp`)),
      signalsHelp: computed(() => ts(`${basePath}.variablesTab.signalsHelp`)),
    },
    variables: {
      editTitle: computed(() => tsTitle(`${basePath}.variables.editTitle`)),
      addTitle: computed(() => tsTitle(`${basePath}.variables.addTitle`)),
      name: computed(() => tsTitle(`${basePath}.variables.name`)),
      nameHint: computed(() => tsRaw(`${basePath}.variables.nameHint`)),
      type: computed(() => tsTitle(`${basePath}.variables.type`)),
      description: computed(() => tsTitle(`${basePath}.variables.description`)),
      descriptionHint: computed(() => ts(`${basePath}.variables.descriptionHint`)),
      defaultValue: computed(() => tsTitle(`${basePath}.variables.defaultValue`)),
      defaultValueJson: computed(() => tsTitle(`${basePath}.variables.defaultValueJson`)),
      persistence: computed(() => tsTitle(`${basePath}.variables.persistence`)),
      ephemeralLabel: computed(() => ts(`${basePath}.variables.ephemeralLabel`)),
      durableLabel: computed(() => ts(`${basePath}.variables.durableLabel`)),
      ephemeralTooltip: computed(() => ts(`${basePath}.variables.ephemeralTooltip`)),
      durableTooltip: computed(() => ts(`${basePath}.variables.durableTooltip`)),
      update: computed(() => tsTitle(`${basePath}.variables.update`)),
      add: computed(() => tsTitle(`${basePath}.variables.add`)),
      cancel: computed(() => tsTitle(`${basePath}.variables.cancel`)),
      emptyTitle: computed(() => tsTitle(`${basePath}.variables.emptyTitle`)),
      emptyDescription: computed(() => ts(`${basePath}.variables.emptyDescription`)),
      defaultLabel: computed(() => tsRaw(`${basePath}.variables.defaultLabel`)),
    },
    externalInputs: {
      editTitle: computed(() => tsTitle(`${basePath}.externalInputs.editTitle`)),
      addTitle: computed(() => tsTitle(`${basePath}.externalInputs.addTitle`)),
      field: computed(() => tsTitle(`${basePath}.externalInputs.field`)),
      fieldHint: computed(() => tsRaw(`${basePath}.externalInputs.fieldHint`)),
      label: computed(() => tsTitle(`${basePath}.externalInputs.label`)),
      labelHint: computed(() => ts(`${basePath}.externalInputs.labelHint`)),
      icon: computed(() => tsTitle(`${basePath}.externalInputs.icon`)),
      pickIcon: computed(() => tsTitle(`${basePath}.externalInputs.pickIcon`)),
      manualIcon: computed(() => ts(`${basePath}.externalInputs.manualIcon`)),
      type: computed(() => tsTitle(`${basePath}.externalInputs.type`)),
      description: computed(() => tsTitle(`${basePath}.externalInputs.description`)),
      descriptionHint: computed(() => ts(`${basePath}.externalInputs.descriptionHint`)),
      defaultValue: computed(() => tsTitle(`${basePath}.externalInputs.defaultValue`)),
      defaultValueJson: computed(() => tsTitle(`${basePath}.externalInputs.defaultValueJson`)),
      required: computed(() => tsTitle(`${basePath}.externalInputs.required`)),
      requiredHint: computed(() => ts(`${basePath}.externalInputs.requiredHint`)),
      update: computed(() => tsTitle(`${basePath}.externalInputs.update`)),
      add: computed(() => tsTitle(`${basePath}.externalInputs.add`)),
      cancel: computed(() => tsTitle(`${basePath}.externalInputs.cancel`)),
      emptyTitle: computed(() => tsTitle(`${basePath}.externalInputs.emptyTitle`)),
      emptyDescription: computed(() => ts(`${basePath}.externalInputs.emptyDescription`)),
      requiredBadge: computed(() => tsRaw(`${basePath}.externalInputs.requiredBadge`)),
      optionalBadge: computed(() => tsRaw(`${basePath}.externalInputs.optionalBadge`)),
      sectionIdentity: computed(() => tsTitle(`${basePath}.externalInputs.sectionIdentity`)),
      sectionTypeConfig: computed(() => tsTitle(`${basePath}.externalInputs.sectionTypeConfig`)),
      sectionMetadata: computed(() => tsTitle(`${basePath}.externalInputs.sectionMetadata`)),
      literalValue: computed(() => tsTitle(`${basePath}.externalInputs.literalValue`)),
      literalValueHint: computed(() => ts(`${basePath}.externalInputs.literalValueHint`)),
      literalInfo: computed(() => ts(`${basePath}.externalInputs.literalInfo`)),
      assetTemplate: computed(() => tsTitle(`${basePath}.externalInputs.assetTemplate`)),
      assetTemplateHint: computed(() => ts(`${basePath}.externalInputs.assetTemplateHint`)),
      selectAssetTemplate: computed(() => ts(`${basePath}.externalInputs.selectAssetTemplate`)),
      fieldPath: computed(() => tsTitle(`${basePath}.externalInputs.fieldPath`)),
      fieldPathHint: computed(() => ts(`${basePath}.externalInputs.fieldPathHint`)),
      assetFromTemplateInfo: computed(() => ts(`${basePath}.externalInputs.assetFromTemplateInfo`)),
      eventFieldsBannerTitle: computed(() => tsTitle(`${basePath}.externalInputs.eventFieldsBannerTitle`)),
      eventFieldsBannerDescription: computed(() => tsRaw(`${basePath}.externalInputs.eventFieldsBannerDescription`)),
    },
    captureFields: {
      editTitle: computed(() => tsTitle(`${basePath}.captureFields.editTitle`)),
      addTitle: computed(() => tsTitle(`${basePath}.captureFields.addTitle`)),
      fieldName: computed(() => tsTitle(`${basePath}.captureFields.fieldName`)),
      type: computed(() => tsTitle(`${basePath}.captureFields.type`)),
      description: computed(() => tsTitle(`${basePath}.captureFields.description`)),
      descriptionHint: computed(() => ts(`${basePath}.captureFields.descriptionHint`)),
      update: computed(() => tsTitle(`${basePath}.captureFields.update`)),
      add: computed(() => tsTitle(`${basePath}.captureFields.add`)),
      cancel: computed(() => tsTitle(`${basePath}.captureFields.cancel`)),
      emptyTitle: computed(() => tsTitle(`${basePath}.captureFields.emptyTitle`)),
      emptyDescription: computed(() => ts(`${basePath}.captureFields.emptyDescription`)),
    },
    externalSignals: {
      editTitle: computed(() => tsTitle(`${basePath}.externalSignals.editTitle`)),
      addTitle: computed(() => tsTitle(`${basePath}.externalSignals.addTitle`)),
      name: computed(() => tsTitle(`${basePath}.externalSignals.name`)),
      description: computed(() => tsTitle(`${basePath}.externalSignals.description`)),
      descriptionHint: computed(() => ts(`${basePath}.externalSignals.descriptionHint`)),
      update: computed(() => tsTitle(`${basePath}.externalSignals.update`)),
      add: computed(() => tsTitle(`${basePath}.externalSignals.add`)),
      cancel: computed(() => tsTitle(`${basePath}.externalSignals.cancel`)),
      emptyTitle: computed(() => tsTitle(`${basePath}.externalSignals.emptyTitle`)),
      emptyDescription: computed(() => ts(`${basePath}.externalSignals.emptyDescription`)),
    },
    jsonDebug: {
      title: computed(() => tsTitle(`${basePath}.jsonDebug.title`)),
      descriptionView: computed(() => ts(`${basePath}.jsonDebug.descriptionView`)),
      descriptionEdit: computed(() => ts(`${basePath}.jsonDebug.descriptionEdit`)),
      viewMode: computed(() => tsTitle(`${basePath}.jsonDebug.viewMode`)),
      editMode: computed(() => tsTitle(`${basePath}.jsonDebug.editMode`)),
      copyTooltip: computed(() => ts(`${basePath}.jsonDebug.copyTooltip`)),
      cancel: computed(() => tsTitle(`${basePath}.jsonDebug.cancel`)),
      applyChanges: computed(() => tsTitle(`${basePath}.jsonDebug.applyChanges`)),
      guidelines: computed(() => tsTitle(`${basePath}.jsonDebug.guidelines`)),
      expectedStructure: computed(() => ts(`${basePath}.jsonDebug.expectedStructure`)),
      structureName: computed(() => tsRaw(`${basePath}.jsonDebug.structureName`)),
      structureDescription: computed(() => tsRaw(`${basePath}.jsonDebug.structureDescription`)),
      structureVariables: computed(() => tsRaw(`${basePath}.jsonDebug.structureVariables`)),
      structureCaptureFields: computed(() => tsRaw(`${basePath}.jsonDebug.structureCaptureFields`)),
      structureExternalVariables: computed(() => tsRaw(`${basePath}.jsonDebug.structureExternalVariables`)),
      structureNodes: computed(() => tsRaw(`${basePath}.jsonDebug.structureNodes`)),
      structureEdges: computed(() => tsRaw(`${basePath}.jsonDebug.structureEdges`)),
      editModeTitle: computed(() => ts(`${basePath}.jsonDebug.editModeTitle`)),
      editModeHint1: computed(() => ts(`${basePath}.jsonDebug.editModeHint1`)),
      editModeHint2: computed(() => ts(`${basePath}.jsonDebug.editModeHint2`)),
      editModeHint3: computed(() => ts(`${basePath}.jsonDebug.editModeHint3`)),
      editModeHint4: computed(() => ts(`${basePath}.jsonDebug.editModeHint4`)),
      tipsTitle: computed(() => ts(`${basePath}.jsonDebug.tipsTitle`)),
      tipSearch: computed(() => tsRaw(`${basePath}.jsonDebug.tipSearch`)),
      tipSelectAll: computed(() => tsRaw(`${basePath}.jsonDebug.tipSelectAll`)),
      tipFormat: computed(() => ts(`${basePath}.jsonDebug.tipFormat`)),
    },
    nodeConfig: {
      configTab: computed(() => tsTitle(`${basePath}.nodeConfig.configTab`)),
      notesTab: computed(() => tsTitle(`${basePath}.nodeConfig.notesTab`)),
      nodeInfo: computed(() => tsRaw(`${basePath}.nodeConfig.nodeInfo`)),
      nodeId: computed(() => tsRaw(`${basePath}.nodeConfig.nodeId`)),
      connections: computed(() => ts(`${basePath}.nodeConfig.connections`)),
      connectionsFormat: computed(() => tsRaw(`${basePath}.nodeConfig.connectionsFormat`)),
      position: computed(() => ts(`${basePath}.nodeConfig.position`)),
      configJson: computed(() => tsTitle(`${basePath}.nodeConfig.configJson`)),
      noConfigForm: computed(() => ts(`${basePath}.nodeConfig.noConfigForm`)),
      /**
       * Format validation error count message
       *
       * @param {number} count - Number of validation errors
       * @returns {string} Formatted error count string
       */
      validationErrorCount: (count: number) => ts(`${basePath}.nodeConfig.validationErrorCount`, { count }),
      notesInfo: computed(() => ts(`${basePath}.nodeConfig.notesInfo`)),
      notesPlaceholder: computed(() => ts(`${basePath}.nodeConfig.notesPlaceholder`)),
      delete: computed(() => tsTitle(`${basePath}.nodeConfig.delete`)),
      apply: computed(() => tsTitle(`${basePath}.nodeConfig.apply`)),
      nodeNotFound: computed(() => tsTitle(`${basePath}.nodeConfig.nodeNotFound`)),
    },
    canvasToolbar: {
      autoOrganize: computed(() => tsTitle(`${basePath}.canvasToolbar.autoOrganize`)),
      unlockCanvas: computed(() => tsTitle(`${basePath}.canvasToolbar.unlockCanvas`)),
      lockCanvas: computed(() => tsTitle(`${basePath}.canvasToolbar.lockCanvas`)),
      undo: computed(() => tsRaw(`${basePath}.canvasToolbar.undo`)),
      redo: computed(() => tsRaw(`${basePath}.canvasToolbar.redo`)),
      minimap: computed(() => tsTitle(`${basePath}.canvasToolbar.minimap`)),
      grid: computed(() => tsTitle(`${basePath}.canvasToolbar.grid`)),
      keyboardShortcuts: computed(() => tsTitle(`${basePath}.canvasToolbar.keyboardShortcuts`)),
      exitFullscreen: computed(() => tsTitle(`${basePath}.canvasToolbar.exitFullscreen`)),
      fullscreen: computed(() => tsTitle(`${basePath}.canvasToolbar.fullscreen`)),
      shortcutsTitle: computed(() => tsTitle(`${basePath}.canvasToolbar.shortcutsTitle`)),
      shortcutsDescription: computed(() => ts(`${basePath}.canvasToolbar.shortcutsDescription`)),
      close: computed(() => tsTitle(`${basePath}.canvasToolbar.close`)),
      hotkeyEscape: computed(() => ts(`${basePath}.canvasToolbar.hotkeyEscape`)),
      hotkeyCopy: computed(() => ts(`${basePath}.canvasToolbar.hotkeyCopy`)),
      hotkeyPaste: computed(() => ts(`${basePath}.canvasToolbar.hotkeyPaste`)),
      hotkeyUndo: computed(() => ts(`${basePath}.canvasToolbar.hotkeyUndo`)),
      hotkeyRedo: computed(() => ts(`${basePath}.canvasToolbar.hotkeyRedo`)),
      hotkeyDelete: computed(() => ts(`${basePath}.canvasToolbar.hotkeyDelete`)),
      hotkeyDuplicate: computed(() => ts(`${basePath}.canvasToolbar.hotkeyDuplicate`)),
      hotkeyBoxSelect: computed(() => ts(`${basePath}.canvasToolbar.hotkeyBoxSelect`)),
      hotkeyMultiSelect: computed(() => ts(`${basePath}.canvasToolbar.hotkeyMultiSelect`)),
    },
    pluginCatalog: {
      title: computed(() => tsTitle(`${basePath}.pluginCatalog.title`)),
      expand: computed(() => ts(`${basePath}.pluginCatalog.expand`)),
      collapse: computed(() => ts(`${basePath}.pluginCatalog.collapse`)),
      searchPlaceholder: computed(() => tsRaw(`${basePath}.pluginCatalog.searchPlaceholder`)),
      noNodesFound: computed(() => tsTitle(`${basePath}.pluginCatalog.noNodesFound`)),
      noPlugins: computed(() => tsTitle(`${basePath}.pluginCatalog.noPlugins`)),
      categories: {
        triggers: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.triggers`)),
        logic: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.logic`)),
        state: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.state`)),
        flowControl: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.flow_control`)),
        timers: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.timers`)),
        integrations: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.integrations`)),
        observability: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.observability`)),
        annotations: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.annotations`)),
        custom: computed(() => tsTitle(`${basePath}.pluginCatalog.categories.custom`)),
      },
    },
    pluginsTab: {
      installedTab: computed(() => tsTitle(`${basePath}.pluginsTab.installedTab`)),
      marketplaceTab: computed(() => tsTitle(`${basePath}.pluginsTab.marketplaceTab`)),
      installedHelp: computed(() => ts(`${basePath}.pluginsTab.installedHelp`)),
      marketplaceHelp: computed(() => ts(`${basePath}.pluginsTab.marketplaceHelp`)),
      install: computed(() => tsTitle(`${basePath}.pluginsTab.install`)),
      installed: computed(() => tsTitle(`${basePath}.pluginsTab.installed`)),
      uninstall: computed(() => tsTitle(`${basePath}.pluginsTab.uninstall`)),
      configure: computed(() => tsTitle(`${basePath}.pluginsTab.configure`)),
      installing: computed(() => ts(`${basePath}.pluginsTab.installing`)),
      searchPlaceholder: computed(() => tsRaw(`${basePath}.pluginsTab.searchPlaceholder`)),
      noPluginsInstalled: computed(() => tsTitle(`${basePath}.pluginsTab.noPluginsInstalled`)),
      noPluginsInstalledDesc: computed(() => ts(`${basePath}.pluginsTab.noPluginsInstalledDesc`)),
      noPluginsAvailable: computed(() => tsTitle(`${basePath}.pluginsTab.noPluginsAvailable`)),
      credentialsTitle: computed(() => tsTitle(`${basePath}.pluginsTab.credentialsTitle`)),
      credentialsSave: computed(() => tsTitle(`${basePath}.pluginsTab.credentialsSave`)),
      credentialsCancel: computed(() => tsTitle(`${basePath}.pluginsTab.credentialsCancel`)),
      credentialsMockInfo: computed(() => ts(`${basePath}.pluginsTab.credentialsMockInfo`)),
      filterByCategory: computed(() => tsTitle(`${basePath}.pluginsTab.filterByCategory`)),
      noCategoryResults: computed(() => ts(`${basePath}.pluginsTab.noCategoryResults`)),
      installFailed: computed(() => ts(`${basePath}.pluginsTab.installFailed`)),
      details: computed(() => tsTitle(`${basePath}.pluginsTab.details`)),
      corePlugins: computed(() => tsTitle(`${basePath}.pluginsTab.corePlugins`)),
      installedPlugins: computed(() => tsTitle(`${basePath}.pluginsTab.installedPlugins`)),
      nodeTypes: computed(() => tsTitle(`${basePath}.pluginsTab.nodeTypes`)),
      inputs: computed(() => tsRaw(`${basePath}.pluginsTab.inputs`)),
      outputs: computed(() => tsRaw(`${basePath}.pluginsTab.outputs`)),
      properties: computed(() => tsTitle(`${basePath}.pluginsTab.properties`)),
      required: computed(() => tsRaw(`${basePath}.pluginsTab.required`)),
      noProperties: computed(() => ts(`${basePath}.pluginsTab.noProperties`)),
      loadingManifest: computed(() => ts(`${basePath}.pluginsTab.loadingManifest`)),
      /**
       * Format install success message with plugin name
       *
       * @param {string} name - Plugin name
       * @returns {string} Formatted success message
       */
      installSuccessMsg: (name: string) => ts(`${basePath}.pluginsTab.installSuccess`, { name }),
      /**
       * Format uninstall success message with plugin name
       *
       * @param {string} name - Plugin name
       * @returns {string} Formatted success message
       */
      uninstallSuccessMsg: (name: string) => ts(`${basePath}.pluginsTab.uninstallSuccess`, { name }),
    },
    timeout: {
      sectionTitle: computed(() => tsTitle(`${basePath}.timeout.sectionTitle`)),
      duration: computed(() => tsTitle(`${basePath}.timeout.duration`)),
      unit: computed(() => tsTitle(`${basePath}.timeout.unit`)),
      enableOutput: computed(() => ts(`${basePath}.timeout.enableOutput`)),
      enableOutputHint: computed(() => tsRaw(`${basePath}.timeout.enableOutputHint`)),
      units: {
        seconds: computed(() => ts(`${basePath}.timeout.units.seconds`)),
        minutes: computed(() => ts(`${basePath}.timeout.units.minutes`)),
        hours: computed(() => ts(`${basePath}.timeout.units.hours`)),
        days: computed(() => ts(`${basePath}.timeout.units.days`)),
        months: computed(() => ts(`${basePath}.timeout.units.months`)),
        years: computed(() => ts(`${basePath}.timeout.units.years`)),
      },
    },
    errorHandler: {
      sectionTitle: computed(() => tsTitle(`${basePath}.errorHandler.sectionTitle`)),
      banner: computed(() => ts(`${basePath}.errorHandler.banner`)),
      enabled: computed(() => ts(`${basePath}.errorHandler.enabled`)),
      maxAttempts: computed(() => tsTitle(`${basePath}.errorHandler.maxAttempts`)),
      maxAttemptsHint: computed(() => ts(`${basePath}.errorHandler.maxAttemptsHint`)),
      initialInterval: computed(() => tsTitle(`${basePath}.errorHandler.initialInterval`)),
      intervalUnit: computed(() => tsTitle(`${basePath}.errorHandler.intervalUnit`)),
      backoffMultiplier: computed(() => tsTitle(`${basePath}.errorHandler.backoffMultiplier`)),
      backoffMultiplierHint: computed(() => ts(`${basePath}.errorHandler.backoffMultiplierHint`)),
      units: {
        seconds: computed(() => ts(`${basePath}.errorHandler.units.seconds`)),
        minutes: computed(() => ts(`${basePath}.errorHandler.units.minutes`)),
        hours: computed(() => ts(`${basePath}.errorHandler.units.hours`)),
      },
    },
    credentials: {
      newCredential: computed(() => tsTitle(`${basePath}.credentials.newCredential`)),
      editCredential: computed(() => tsTitle(`${basePath}.credentials.editCredential`)),
      credentialType: computed(() => tsTitle(`${basePath}.credentials.credentialType`)),
      credentialName: computed(() => tsTitle(`${basePath}.credentials.credentialName`)),
      credentialNameHint: computed(() => ts(`${basePath}.credentials.credentialNameHint`)),
      save: computed(() => tsTitle(`${basePath}.credentials.save`)),
      cancel: computed(() => tsTitle(`${basePath}.credentials.cancel`)),
      test: computed(() => tsTitle(`${basePath}.credentials.test`)),
      testing: computed(() => ts(`${basePath}.credentials.testing`)),
      testSuccess: computed(() => ts(`${basePath}.credentials.testSuccess`)),
      deleteSuccess: computed(() => ts(`${basePath}.credentials.deleteSuccess`)),
      saveSuccess: computed(() => ts(`${basePath}.credentials.saveSuccess`)),
      saveFailed: computed(() => ts(`${basePath}.credentials.saveFailed`)),
      noCredentials: computed(() => tsTitle(`${basePath}.credentials.noCredentials`)),
      secretPlaceholder: computed(() => tsRaw(`${basePath}.credentials.secretPlaceholder`)),
      credentialSelector: computed(() => tsTitle(`${basePath}.credentials.credentialSelector`)),
      selectCredential: computed(() => ts(`${basePath}.credentials.selectCredential`)),
      noCredentialsAvailable: computed(() => ts(`${basePath}.credentials.noCredentialsAvailable`)),
      /**
       * Format credential title with plugin name
       *
       * @param {string} name - Plugin name
       * @returns {string} Formatted title
       */
      titleMsg: (name: string) => tsTitle(`${basePath}.credentials.title`, { name }),
      /**
       * Format test failed message with error
       *
       * @param {string} error - Error message
       * @returns {string} Formatted error message
       */
      testFailedMsg: (error: string) => ts(`${basePath}.credentials.testFailed`, { error }),
      /**
       * Format no credentials description with plugin name
       *
       * @param {string} name - Plugin name
       * @returns {string} Formatted description
       */
      noCredentialsDescMsg: (name: string) => ts(`${basePath}.credentials.noCredentialsDesc`, { name }),
    },
  };
}

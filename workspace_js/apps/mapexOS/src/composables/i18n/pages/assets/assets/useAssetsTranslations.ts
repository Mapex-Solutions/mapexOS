import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Assets list page translations
 *
 * Structure mirrors:
 * - File: src/pages/assets/assets/assetListPage/AssetsListPage.vue
 * - JSON: src/i18n/{locale}/pages/assets/assets.json
 * - Composable: src/composables/i18n/pages/assets/assets/useAssetsTranslations.ts
 *
 * Provides all translations for the Assets List page including:
 * - Page header (title, description, button)
 * - Filter items (labels, options)
 * - DataRow column definitions (reactive)
 * - Menu column labels
 * - Empty state
 * - Success/error messages
 * - Dialog content
 *
 * @example
 * ```ts
 * // In AssetsListPage.vue
 * const {
 *   page,
 *   filters,
 *   columns,
 *   menuColumns,
 *   empty,
 *   messages
 * } = useAssetsTranslations();
 *
 * <PageHeader :title="page.title.value" :description="page.description.value" />
 * <DataRow :columns="columns.value" />
 * ```
 */
export function useAssetsTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     * Mirrors: pages.assets.assets
     */
    page: {
      title: computed(() => tsTitle('pages.assets.assets.title')),
      description: computed(() => ts('pages.assets.assets.description')),
      addButton: computed(() => ts('pages.assets.assets.addButton')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.assets.assets.info.title'),
        description: ts('pages.assets.assets.info.description'),
        items: [
          {
            icon: 'description',
            color: 'blue-6',
            title: ts('pages.assets.assets.info.items.templates.title'),
            text: ts('pages.assets.assets.info.items.templates.text'),
          },
          {
            icon: 'sync_alt',
            color: 'purple-6',
            title: ts('pages.assets.assets.info.items.protocols.title'),
            text: ts('pages.assets.assets.info.items.protocols.text'),
          },
          {
            icon: 'domain',
            color: 'indigo-6',
            title: ts('pages.assets.assets.info.items.organization.title'),
            text: ts('pages.assets.assets.info.items.organization.text'),
          },
          {
            icon: 'speed',
            color: 'green-6',
            title: ts('pages.assets.assets.info.items.realtime.title'),
            text: ts('pages.assets.assets.info.items.realtime.text'),
          },
          {
            icon: 'place',
            color: 'orange-6',
            title: ts('pages.assets.assets.info.items.location.title'),
            text: ts('pages.assets.assets.info.items.location.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/assets',
        docsLabel: ts('pages.assets.assets.info.docsLabel'),
      })),
      listTitle: computed(() => tsTitle('pages.assets.assets.listTitle')),
      itemLabel: computed(() => ts('pages.assets.assets.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.assets.assets.itemLabelPlural')),
    },

    /**
     * Filter translations with reactive options
     * Mirrors: pages.assets.assets.filters
     */
    filters: {
      label: computed(() => ts('pages.assets.assets.filters.label')),
      searchPlaceholder: computed(() => ts('pages.assets.assets.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.assets.assets.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.assets.assets.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.assets.assets.filters.clearAll')),
      allStatus: computed(() => ts('pages.assets.assets.filters.allStatus')),
      assetName: computed(() => ts('pages.assets.assets.filters.assetName')),
      assetUUID: computed(() => ts('pages.assets.assets.filters.assetUUID')),
      deviceUUID: computed(() => ts('pages.assets.assets.filters.deviceUUID')),
      status: computed(() => ts('pages.assets.assets.filters.status')),
      template: computed(() => ts('pages.assets.assets.filters.template')),
      category: computed(() => ts('pages.assets.assets.filters.category')),
      manufacturer: computed(() => ts('pages.assets.assets.filters.manufacturer')),
      model: computed(() => ts('pages.assets.assets.filters.model')),
      assetType: computed(() => ts('pages.assets.assets.filters.assetType')),
      protocol: computed(() => ts('pages.assets.assets.filters.protocol')),
      createdDate: computed(() => ts('pages.assets.assets.filters.createdDate')),
      includeChildren: computed(() => ts('pages.assets.assets.filters.includeChildren')),
      filterByUUID: computed(() => ts('pages.assets.assets.filters.filterByUUID')),
      filterByCategory: computed(() => ts('pages.assets.assets.filters.filterByCategory')),
      filterByManufacturer: computed(() => ts('pages.assets.assets.filters.filterByManufacturer')),
      filterByModel: computed(() => ts('pages.assets.assets.filters.filterByModel')),

      options: {
        all: computed(() => ts('pages.assets.assets.filters.options.all')),
        active: computed(() => ts('pages.assets.assets.filters.options.active')),
        inactive: computed(() => ts('pages.assets.assets.filters.options.inactive')),
        enabled: computed(() => ts('pages.assets.assets.filters.options.enabled')),
        disabled: computed(() => ts('pages.assets.assets.filters.options.disabled')),
        yes: computed(() => ts('pages.assets.assets.filters.options.yes')),
        no: computed(() => ts('pages.assets.assets.filters.options.no')),
      },
    },

    /**
     * Menu column labels (for ListHeaderMenu)
     * Mirrors: pages.assets.assets.menuColumns
     */
    menuColumns: {
      uuid: computed(() => ts('pages.assets.assets.columns.assetUUID')),
      type: computed(() => ts('pages.assets.assets.menuColumns.type')),
      protocol: computed(() => ts('pages.assets.assets.columns.protocol')),
      category: computed(() => ts('pages.assets.assets.menuColumns.category')),
      manufacturerModel: computed(() => ts('pages.assets.assets.columns.manufacturer')),
      debug: computed(() => ts('pages.assets.assets.columns.debug')),
      status: computed(() => ts('pages.assets.assets.healthMonitoring.status')),
      organization: computed(() => ts('pages.assets.assets.columns.organization')),
    },

    /**
     * Status badge translations
     * Mirrors: pages.assets.assets.statusBadge
     */
    statusBadge: {
      online: computed(() => ts('pages.assets.assets.statusBadge.online')),
      offline: computed(() => ts('pages.assets.assets.statusBadge.offline')),
      unknown: computed(() => ts('pages.assets.assets.statusBadge.unknown')),
      lastSeenPrefix: computed(() => ts('pages.assets.assets.statusBadge.lastSeenPrefix')),
      lastSeenNever: computed(() => ts('pages.assets.assets.statusBadge.lastSeenNever')),
    },

    /**
     * Empty state translations
     * Mirrors: pages.assets.assets.empty
     */
    empty: {
      title: computed(() => ts('pages.assets.assets.empty.title')),
      description: computed(() => ts('pages.assets.assets.empty.description')),
    },

    /**
     * Status labels
     * Mirrors: pages.assets.assets.status
     */
    status: {
      active: computed(() => ts('pages.assets.assets.status.active')),
      inactive: computed(() => ts('pages.assets.assets.status.inactive')),
    },

    /**
     * Message translations
     * Mirrors: pages.assets.assets.messages
     */
    messages: {
      deletedSuccessfully: computed(() => ts('pages.assets.assets.messages.deletedSuccessfully')),
      confirmDelete: (name: string) => ts('pages.assets.assets.messages.confirmDelete', { name }),
    },

    /**
     * Dialog translations
     * Mirrors: pages.assets.assets.dialog
     */
    dialog: {
      deleteTitle: computed(() => ts('pages.assets.assets.dialog.deleteTitle')),
    },

    /**
     * Error message translations
     * Mirrors: pages.assets.assets.errors
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.assets.assets.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.assets.assets.errors.idMissing')),
      noOrganization: computed(() => ts('pages.assets.assets.errors.noOrganization')),
    },

    /**
     * Drawer translations for asset details view
     * Mirrors: pages.assets.assets.drawer
     */
    drawer: {
      title: computed(() => ts('pages.assets.assets.drawer.title')),
      close: computed(() => ts('pages.assets.assets.drawer.close')),
      edit: computed(() => ts('pages.assets.assets.drawer.edit')),
      duplicate: computed(() => ts('pages.assets.assets.drawer.duplicate')),
      loading: computed(() => ts('pages.assets.assets.drawer.loading')),
      error: computed(() => ts('pages.assets.assets.drawer.error')),

      sections: {
        basicInfo: computed(() => ts('pages.assets.assets.drawer.sections.basicInfo')),
        configuration: computed(() => ts('pages.assets.assets.drawer.sections.configuration')),
        protocol: computed(() => ts('pages.assets.assets.drawer.sections.protocol')),
        auth: computed(() => ts('pages.assets.assets.drawer.sections.auth')),
        location: computed(() => ts('pages.assets.assets.drawer.sections.location')),
        routing: computed(() => ts('pages.assets.assets.drawer.sections.routing')),
        healthMonitoring: computed(() => ts('pages.assets.assets.drawer.sections.healthMonitoring')),
        organization: computed(() => ts('pages.assets.assets.drawer.sections.organization')),
        timestamps: computed(() => ts('pages.assets.assets.drawer.sections.timestamps')),
      },

      fields: {
        name: computed(() => ts('pages.assets.assets.drawer.fields.name')),
        uuid: computed(() => ts('pages.assets.assets.drawer.fields.uuid')),
        status: computed(() => ts('pages.assets.assets.drawer.fields.status')),
        description: computed(() => ts('pages.assets.assets.drawer.fields.description')),
        template: computed(() => ts('pages.assets.assets.drawer.fields.template')),
        category: computed(() => ts('pages.assets.assets.drawer.fields.category')),
        manufacturer: computed(() => ts('pages.assets.assets.drawer.fields.manufacturer')),
        model: computed(() => ts('pages.assets.assets.drawer.fields.model')),
        version: computed(() => ts('pages.assets.assets.drawer.fields.version')),
        protocolType: computed(() => ts('pages.assets.assets.drawer.fields.protocolType')),
        clientId: computed(() => ts('pages.assets.assets.drawer.fields.clientId')),
        username: computed(() => ts('pages.assets.assets.drawer.fields.username')),
        tokenExpiresAt: computed(() => ts('pages.assets.assets.drawer.fields.tokenExpiresAt')),
        latitude: computed(() => ts('pages.assets.assets.drawer.fields.latitude')),
        longitude: computed(() => ts('pages.assets.assets.drawer.fields.longitude')),
        routeGroups: computed(() => ts('pages.assets.assets.drawer.fields.routeGroups')),
        healthStatus: computed(() => ts('pages.assets.assets.drawer.fields.healthStatus')),
        lastSeen: computed(() => ts('pages.assets.assets.drawer.fields.lastSeen')),
        threshold: computed(() => ts('pages.assets.assets.drawer.fields.threshold')),
        requiredMisses: computed(() => ts('pages.assets.assets.drawer.fields.requiredMisses')),
        organization: computed(() => ts('pages.assets.assets.drawer.fields.organization')),
        pathKey: computed(() => ts('pages.assets.assets.drawer.fields.pathKey')),
        customer: computed(() => ts('pages.assets.assets.drawer.fields.customer')),
        created: computed(() => ts('pages.assets.assets.drawer.fields.created')),
        updated: computed(() => ts('pages.assets.assets.drawer.fields.updated')),
      },

      empty: {
        description: computed(() => ts('pages.assets.assets.drawer.empty.description')),
        location: computed(() => ts('pages.assets.assets.drawer.empty.location')),
        routeGroups: computed(() => ts('pages.assets.assets.drawer.empty.routeGroups')),
      },

      regenerateMqttToken: {
        button: computed(() => ts('pages.assets.assets.drawer.regenerateMqttToken.button')),
        confirmTitle: computed(() => ts('pages.assets.assets.drawer.regenerateMqttToken.confirmTitle')),
        confirmBody: computed(() => tsRaw('pages.assets.assets.drawer.regenerateMqttToken.confirmBody')),
        success: computed(() => tsRaw('pages.assets.assets.drawer.regenerateMqttToken.success')),
        failed: computed(() => tsRaw('pages.assets.assets.drawer.regenerateMqttToken.failed')),
      },

      auth: {
        password: {
          label: computed(() => ts('pages.assets.assets.drawer.auth.password.label')),
          set: computed(() => ts('pages.assets.assets.drawer.auth.password.set')),
          notSet: computed(() => ts('pages.assets.assets.drawer.auth.password.notSet')),
          hint: computed(() => tsRaw('pages.assets.assets.drawer.auth.password.hint')),
        },
        certificate: {
          label: computed(() => ts('pages.assets.assets.drawer.auth.certificate.label')),
          active: computed(() => ts('pages.assets.assets.drawer.auth.certificate.active')),
          noActive: computed(() => ts('pages.assets.assets.drawer.auth.certificate.noActive')),
          expired: computed(() => ts('pages.assets.assets.drawer.auth.certificate.expired')),
          fields: {
            serial: computed(() => ts('pages.assets.assets.drawer.auth.certificate.fields.serial')),
            fingerprint: computed(() => ts('pages.assets.assets.drawer.auth.certificate.fields.fingerprint')),
            subjectCN: computed(() => ts('pages.assets.assets.drawer.auth.certificate.fields.subjectCN')),
            issued: computed(() => ts('pages.assets.assets.drawer.auth.certificate.fields.issued')),
            expires: computed(() => ts('pages.assets.assets.drawer.auth.certificate.fields.expires')),
          },
          actions: {
            generate: computed(() => ts('pages.assets.assets.drawer.auth.certificate.actions.generate')),
            revoke: computed(() => ts('pages.assets.assets.drawer.auth.certificate.actions.revoke')),
            regenerate: computed(() => ts('pages.assets.assets.drawer.auth.certificate.actions.regenerate')),
          },
          revokeConfirmTitle: computed(() => ts('pages.assets.assets.drawer.auth.certificate.revokeConfirmTitle')),
          revokeConfirmBody: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.revokeConfirmBody')),
          revokeSuccess: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.revokeSuccess')),
          revokeFailed: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.revokeFailed')),
          issueSuccess: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.issueSuccess')),
          issueFailed: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.issueFailed')),
          dialog: {
            title: computed(() => ts('pages.assets.assets.drawer.auth.certificate.dialog.title')),
            warning: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.dialog.warning')),
            replaceWarning: computed(() => tsRaw('pages.assets.assets.drawer.auth.certificate.dialog.replaceWarning')),
            generateButton: computed(() => ts('pages.assets.assets.drawer.auth.certificate.dialog.generateButton')),
            skipButton: computed(() => ts('pages.assets.assets.drawer.auth.certificate.dialog.skipButton')),
          },
        },
        revoked: {
          title: computed(() => ts('pages.assets.assets.drawer.auth.revoked.title')),
          retentionNotice: computed(() => tsRaw('pages.assets.assets.drawer.auth.revoked.retentionNotice')),
          empty: computed(() => tsRaw('pages.assets.assets.drawer.auth.revoked.empty')),
          columns: {
            serial: computed(() => ts('pages.assets.assets.drawer.auth.revoked.columns.serial')),
            reason: computed(() => ts('pages.assets.assets.drawer.auth.revoked.columns.reason')),
            revokedAt: computed(() => ts('pages.assets.assets.drawer.auth.revoked.columns.revokedAt')),
          },
        },
      },
    },

    /**
     * DataRow column definitions with reactive translations
     * IMPORTANT: Returns a computed ref so columns update when language changes
     *
     * Usage: <DataRow :columns="columns.value" />
     */
    columns: computed(() => {
      return [
        {
          key: 'icon',
          label: '',
          type: 'avatar',
          visible: 'always',
          width: 56,
          icon: (value, row) => row.icon || 'sensors',
          color: (value, row) => row.enabled ? 'primary' : 'grey-5',
          tooltip: (value, row) =>
            row.enabled
              ? ts('pages.assets.assets.status.active')
              : ts('pages.assets.assets.status.inactive'),
        },
        {
          key: 'name',
          label: ts('pages.assets.assets.columns.name'),
          type: 'text',
          visible: 'always',
          width: 200,
          ellipsis: true,
          secondaryKey: 'description',
        },
        {
          key: 'organizationName',
          label: ts('pages.assets.assets.columns.organization'),
          type: 'chip',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'manufacturerName',
          label: ts('pages.assets.assets.columns.manufacturer'),
          type: 'text',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 220,
          ellipsis: true,
          secondaryKey: 'modelName',
        },
        {
          key: 'assetUUID',
          label: ts('pages.assets.assets.columns.assetUUID'),
          type: 'code',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 180,
          ellipsis: true,
        },
        {
          key: 'protocol.type',
          label: ts('pages.assets.assets.columns.protocol'),
          type: 'chip',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 120,
          format: (value: any) => value?.toUpperCase() || 'N/A',
          color: (value) => {
            const colors: Record<string, string> = {
              'mqtt': 'purple-6',
              'http': 'blue-6',
              'lorawan': 'orange-6',
            };
            return colors[value?.toLowerCase()] || 'grey-6';
          },
          // Surface a "no certificate" warning under the chip when a
          // cert-mode MQTT asset has no active mTLS certificate — the
          // device cannot connect in that state. Password-mode MQTT
          // assets never trigger the warning because they do not need
          // a cert by design. The drawer's Auth section is the place
          // to fix the gap.
          secondary: (value: any, row: any) => {
            const isMqtt = value?.toLowerCase() === 'mqtt';
            const isCert = row?.protocol?.mqtt?.authType === 'cert';
            if (isMqtt && isCert && !row?.currentCert?.serial) {
              return ts('pages.assets.assets.warnings.noCertSecondary');
            }
            return '';
          },
          tooltip: (value: any, row: any) => {
            const isMqtt = value?.toLowerCase() === 'mqtt';
            const isCert = row?.protocol?.mqtt?.authType === 'cert';
            if (isMqtt && isCert && !row?.currentCert?.serial) {
              return ts('pages.assets.assets.warnings.noCertTooltip');
            }
            return undefined as any;
          },
        },
        {
          key: 'debugEnabled',
          label: ts('pages.assets.assets.columns.debug'),
          type: 'chip',
          visible: 'laptop',
          width: 80,
          format: (value: any) => value ? 'ON' : 'OFF',
          color: (value: any) => value ? 'orange-6' : 'grey-5',
          icon: (value: any) => value ? 'bug_report' : 'bug_report',
          align: 'center',
        },
        {
          key: 'healthStatus',
          label: ts('pages.assets.assets.healthMonitoring.status'),
          type: 'chip',
          visible: 'laptop',
          width: 140,
          format: (value: any) => {
            const labels: Record<string, string> = {
              online: ts('pages.assets.assets.statusBadge.online'),
              offline: ts('pages.assets.assets.statusBadge.offline'),
              unknown: ts('pages.assets.assets.statusBadge.unknown'),
            };
            return labels[value] || labels.unknown;
          },
          color: (value: any) => {
            const colors: Record<string, string> = {
              online: 'green',
              offline: 'red',
              unknown: 'grey-5',
            };
            return colors[value] || 'grey-5';
          },
          icon: (value: any) => {
            if (value === 'online') return 'wifi';
            if (value === 'offline') return 'wifi_off';
            return 'help_outline';
          },
          secondary: (value: any, row: any) => {
            const iso = row?.lastSeenAt ?? row?.healthStatusChangedAt;
            // No timestamp AND unknown status → don't noise the row with a
            // placeholder; the em-dash chip already says "no data yet".
            if (!iso) {
              if (value === 'unknown' || !value) return '';
              return ts('pages.assets.assets.statusBadge.lastSeenNever');
            }
            const elapsedMs = Date.now() - new Date(iso).getTime();
            if (Number.isNaN(elapsedMs) || elapsedMs < 0) return '';
            const prefix = ts('pages.assets.assets.statusBadge.lastSeenPrefix');
            const sec = Math.floor(elapsedMs / 1000);
            if (sec < 60) return `${prefix} ${sec}s ago`;
            const min = Math.floor(sec / 60);
            if (min < 60) return `${prefix} ${min}m ago`;
            const hr = Math.floor(min / 60);
            if (hr < 24) return `${prefix} ${hr}h ago`;
            const days = Math.floor(hr / 24);
            return `${prefix} ${days}d ago`;
          },
          align: 'center',
        },
      ] as DataRowColumn[];
    }),

    /**
     * Filter items with reactive translations
     * Returns computed FilterListItem array for ListFilter component
     *
     * Contract-Based Pattern (ZodAssetQuerySchema):
     * Row 1 (Primary Filters): Search, Status, Include Children, UUID
     * Row 2 (Domain Filters): Template, Category, Asset Type
     *
     * Usage: <ListFilter :items="filterItems.value" />
     */
    filterItems: computed(() => {
      return [
        // Row 1: Primary Filters
        {
          key: 'name',
          type: 'input',
          label: ts('pages.assets.assets.filters.assetName'),
          icon: 'search',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'status',
          type: 'select',
          label: ts('pages.assets.assets.filters.status'),
          icon: 'toggle_on',
          options: [
            { label: ts('pages.assets.assets.filters.options.all'), value: undefined },
            { label: ts('pages.assets.assets.filters.options.active'), value: true },
            { label: ts('pages.assets.assets.filters.options.inactive'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.assets.assets.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.assets.assets.filters.options.all'), value: null },
            { label: ts('pages.assets.assets.filters.options.yes'), value: true },
            { label: ts('pages.assets.assets.filters.options.no'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },

        // Row 2: Domain-Specific Filters (Cascading: Category → Manufacturer → Model)
        {
          key: 'assetUUID',
          type: 'input',
          label: ts('pages.assets.assets.filters.assetUUID'),
          icon: 'fingerprint',
          grid: 'col-12 col-md-3'
        },
        {
          key: 'category',
          type: 'input',
          label: ts('pages.assets.assets.filters.category'),
          icon: 'category',
          grid: 'col-12 col-md-3'
        },
        {
          key: 'manufacture',
          type: 'input',
          label: ts('pages.assets.assets.filters.manufacturer'),
          icon: 'factory',
          grid: 'col-12 col-md-3',
          disabled: true,
        },
        {
          key: 'model',
          type: 'input',
          label: ts('pages.assets.assets.filters.model'),
          icon: 'memory',
          grid: 'col-12 col-md-3',
          disabled: true,
        },
      ] as FilterListItem[];
    }),
  };
}

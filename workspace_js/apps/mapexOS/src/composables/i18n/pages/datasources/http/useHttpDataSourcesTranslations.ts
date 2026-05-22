import type { FilterField } from '@components/drawers';
import type { DataRowColumn } from '@components/cards';
import type { PageHeaderInfo } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useHttpDataSourcesTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	return {
		pageHeader: {
			title: computed(() => tsTitle('pages.datasources.http.pageHeader.title')),
			description: computed(() => ts('pages.datasources.http.pageHeader.description')),
			button: computed(() => ts('pages.datasources.http.pageHeader.button')),
			info: computed((): PageHeaderInfo => ({
				title: ts('pages.datasources.http.pageHeader.info.title'),
				description: ts('pages.datasources.http.pageHeader.info.description'),
				items: [
					{
						icon: 'download',
						color: 'green-6',
						title: ts('pages.datasources.http.pageHeader.info.items.pull.title'),
						text: ts('pages.datasources.http.pageHeader.info.items.pull.text'),
					},
					{
						icon: 'upload',
						color: 'orange-6',
						title: ts('pages.datasources.http.pageHeader.info.items.push.title'),
						text: ts('pages.datasources.http.pageHeader.info.items.push.text'),
					},
					{
						icon: 'lock',
						color: 'blue-6',
						title: ts('pages.datasources.http.pageHeader.info.items.authentication.title'),
						text: ts('pages.datasources.http.pageHeader.info.items.authentication.text'),
					},
					{
						icon: 'link',
						color: 'purple-6',
						title: ts('pages.datasources.http.pageHeader.info.items.assetBinding.title'),
						text: ts('pages.datasources.http.pageHeader.info.items.assetBinding.text'),
					},
				],
				docsUrl: 'https://docs.mapexos.com/datasources/http',
				docsLabel: ts('pages.datasources.http.pageHeader.info.docsLabel'),
			})),
		},

		menuColumns: {
			organization: computed(() => ts('pages.datasources.http.columns.organization')),
			assetBind: computed(() => ts('pages.datasources.http.columns.assetBind')),
			auth: computed(() => ts('pages.datasources.http.columns.auth')),
			mode: computed(() => ts('pages.datasources.http.columns.mode')),
		},

		filters: {
			label: computed(() => ts('pages.datasources.http.filters.label')),
			searchPlaceholder: computed(() => ts('pages.datasources.http.filters.searchPlaceholder')),
			allStatus: computed(() => ts('pages.datasources.http.filters.allStatus')),
			advancedFilters: computed(() => ts('pages.datasources.http.filters.advancedFilters')),
			pendingFilters: computed(() => ts('pages.datasources.http.filters.pendingFilters')),
			clearAll: computed(() => ts('pages.datasources.http.filters.clearAll')),
			name: computed(() => ts('pages.datasources.http.filters.name')),
			mode: computed(() => ts('pages.datasources.http.filters.mode')),
			status: computed(() => ts('pages.datasources.http.filters.status')),
			includeChildren: computed(() => ts('pages.datasources.http.filters.includeChildren')),
			authType: computed(() => ts('pages.datasources.http.filters.authType')),
			assetBindType: computed(() => ts('pages.datasources.http.filters.assetBindType')),
			options: {
				active: computed(() => ts('pages.datasources.http.filters.statusOptions.active')),
				inactive: computed(() => ts('pages.datasources.http.filters.statusOptions.inactive')),
				pull: computed(() => ts('pages.datasources.http.filters.modeOptions.pull')),
				push: computed(() => ts('pages.datasources.http.filters.modeOptions.push')),
				yes: computed(() => ts('pages.datasources.http.filters.includeChildrenOptions.yes')),
				no: computed(() => ts('pages.datasources.http.filters.includeChildrenOptions.no')),
				none: computed(() => ts('pages.datasources.http.filters.authTypeOptions.none')),
				apiKey: computed(() => ts('pages.datasources.http.filters.authTypeOptions.apiKey')),
				ipWhitelist: computed(() => ts('pages.datasources.http.filters.authTypeOptions.ipWhitelist')),
				jwt: computed(() => ts('pages.datasources.http.filters.authTypeOptions.jwt')),
				oauth2: computed(() => ts('pages.datasources.http.filters.authTypeOptions.oauth2')),
				fixed: computed(() => ts('pages.datasources.http.filters.assetBindTypeOptions.fixed')),
				dynamic: computed(() => ts('pages.datasources.http.filters.assetBindTypeOptions.dynamic')),
			},
		},

		advancedFilters: computed((): FilterField[] => [
			{
				key: 'includeChildren',
				type: 'toggle',
				label: ts('pages.datasources.http.filters.includeChildren'),
				icon: 'account_tree',
				options: [
					{ label: ts('pages.datasources.http.filters.allStatus'), value: null },
					{ label: ts('pages.datasources.http.filters.includeChildrenOptions.yes'), value: true },
					{ label: ts('pages.datasources.http.filters.includeChildrenOptions.no'), value: false },
				],
			},
			{
				key: 'mode',
				type: 'select',
				label: ts('pages.datasources.http.filters.mode'),
				icon: 'sync_alt',
				options: [
					{ label: ts('pages.datasources.http.filters.allStatus'), value: null },
					{ label: ts('pages.datasources.http.filters.modeOptions.pull'), value: 'pull' },
					{ label: ts('pages.datasources.http.filters.modeOptions.push'), value: 'push' },
				],
			},
			{
				key: 'auth',
				type: 'select',
				label: ts('pages.datasources.http.filters.authType'),
				icon: 'lock',
				options: [
					{ label: ts('pages.datasources.http.filters.allStatus'), value: null },
					{ label: ts('pages.datasources.http.filters.authTypeOptions.none'), value: 'none' },
					{ label: ts('pages.datasources.http.filters.authTypeOptions.apiKey'), value: 'apiKey' },
					{ label: ts('pages.datasources.http.filters.authTypeOptions.ipWhitelist'), value: 'ip_whitelist' },
					{ label: ts('pages.datasources.http.filters.authTypeOptions.jwt'), value: 'jwt' },
					{ label: ts('pages.datasources.http.filters.authTypeOptions.oauth2'), value: 'oauth2' },
				],
			},
			{
				key: 'assetBind',
				type: 'select',
				label: ts('pages.datasources.http.filters.assetBindType'),
				icon: 'link',
				options: [
					{ label: ts('pages.datasources.http.filters.allStatus'), value: null },
					{ label: ts('pages.datasources.http.filters.assetBindTypeOptions.fixed'), value: 'fixedAssetId' },
					{ label: ts('pages.datasources.http.filters.assetBindTypeOptions.dynamic'), value: 'uuidField' },
				],
			},
		]),

		columns: computed((): DataRowColumn[] => [
			{
				key: 'icon',
				label: '',
				type: 'avatar',
				visible: 'always',
				width: 56,
				icon: () => 'settings_input_antenna',
				color: (_val, row) => row.enabled ? 'primary' : 'grey-5',
				tooltip: (_val, row) =>
					row.enabled
						? ts('pages.datasources.http.drawer.status.active')
						: ts('pages.datasources.http.drawer.status.inactive'),
			},
			{
				key: 'name',
				label: ts('pages.datasources.http.columns.name'),
				type: 'text',
				visible: 'always',
				width: 250,
				ellipsis: true,
				secondaryKey: 'description',
			},
			{
				key: 'organizationName',
				label: ts('pages.datasources.http.columns.organization'),
				type: 'chip',
				visible: 'laptop',
				width: 180,
				ellipsis: true,
				color: 'indigo-6',
				icon: 'domain',
			},
			{
				key: 'assetBind.type',
				label: ts('pages.datasources.http.columns.assetBind'),
				type: 'chip',
				visible: 'laptop',
				width: 130,
				format: (val) => val === 'fixedAssetId' ? 'FIXED' : 'DYNAMIC',
				color: (val) => val === 'fixedAssetId' ? 'blue-6' : 'purple-6',
			},
			{
				key: 'auth.type',
				label: ts('pages.datasources.http.columns.auth'),
				type: 'chip',
				visible: 'laptop',
				width: 130,
				format: (val) => val ? val.toUpperCase() : 'NONE',
			},
			{
				key: 'mode',
				label: ts('pages.datasources.http.columns.mode'),
				type: 'chip',
				visible: 'laptop',
				width: 120,
				format: (val) => val ? val.toUpperCase() : '-',
				color: (val) => val === 'pull' ? 'green-6' : 'orange-6',
			},
		]),

		headerMenu: {
			itemLabel: computed(() => ts('pages.datasources.http.headerMenu.itemLabel')),
			itemLabelPlural: computed(() => ts('pages.datasources.http.headerMenu.itemLabelPlural')),
		},

		list: {
			title: computed(() => ts('pages.datasources.http.list.title')),
		},

		empty: {
			title: computed(() => ts('pages.datasources.http.empty.title')),
			description: computed(() => ts('pages.datasources.http.empty.description')),
		},

		deleteDialog: {
			title: computed(() => ts('pages.datasources.http.deleteDialog.title')),
			message: (name: string) => ts('pages.datasources.http.deleteDialog.message', { name }),
		},

		notifications: {
			deleteSuccess: computed(() => ts('pages.datasources.http.notifications.deleteSuccess')),
			deleteFailed: computed(() => ts('pages.datasources.http.notifications.deleteFailed')),
			loadFailed: computed(() => ts('pages.datasources.http.notifications.loadFailed')),
			endpointCopied: computed(() => ts('pages.datasources.http.notifications.endpointCopied')),
		},

		errors: {
			apiNotInitialized: computed(() => ts('pages.datasources.http.errors.apiNotInitialized')),
			noUuidPaths: computed(() => ts('pages.datasources.http.errors.noUuidPaths')),
			invalidPayloadOrPath: computed(() => ts('pages.datasources.http.errors.invalidPayloadOrPath')),
		},

		actions: {
			copyEndpoint: computed(() => ts('pages.datasources.http.actions.copyEndpoint')),
		},

		drawer: {
			title: computed(() => ts('pages.datasources.http.drawer.title')),
			close: computed(() => ts('pages.datasources.http.drawer.close')),
			edit: computed(() => ts('pages.datasources.http.drawer.edit')),
			loading: computed(() => ts('pages.datasources.http.drawer.loading')),
			error: computed(() => ts('pages.datasources.http.drawer.error')),

			sections: {
				basicInfo: computed(() => ts('pages.datasources.http.drawer.sections.basicInfo')),
				configuration: computed(() => ts('pages.datasources.http.drawer.sections.configuration')),
				authentication: computed(() => ts('pages.datasources.http.drawer.sections.authentication')),
				assetBinding: computed(() => ts('pages.datasources.http.drawer.sections.assetBinding')),
				timestamps: computed(() => ts('pages.datasources.http.drawer.sections.timestamps')),
			},

			fields: {
				name: computed(() => ts('pages.datasources.http.drawer.fields.name')),
				description: computed(() => ts('pages.datasources.http.drawer.fields.description')),
				status: computed(() => ts('pages.datasources.http.drawer.fields.status')),
				mode: computed(() => ts('pages.datasources.http.drawer.fields.mode')),
				protocol: computed(() => ts('pages.datasources.http.drawer.fields.protocol')),
				authType: computed(() => ts('pages.datasources.http.drawer.fields.authType')),
				assetBindType: computed(() => ts('pages.datasources.http.drawer.fields.assetBindType')),
				fixedAssetId: computed(() => ts('pages.datasources.http.drawer.fields.fixedAssetId')),
				uuidField: computed(() => ts('pages.datasources.http.drawer.fields.uuidField')),
				created: computed(() => ts('pages.datasources.http.drawer.fields.created')),
				updated: computed(() => ts('pages.datasources.http.drawer.fields.updated')),
			},

			empty: {
				description: computed(() => ts('pages.datasources.http.drawer.empty.description')),
			},

			status: {
				active: computed(() => ts('pages.datasources.http.drawer.status.active')),
				inactive: computed(() => ts('pages.datasources.http.drawer.status.inactive')),
			},

			modes: {
				pull: computed(() => ts('pages.datasources.http.drawer.modes.pull')),
				push: computed(() => ts('pages.datasources.http.drawer.modes.push')),
			},

			authTypes: {
				none: computed(() => ts('pages.datasources.http.drawer.authTypes.none')),
				apiKey: computed(() => ts('pages.datasources.http.drawer.authTypes.apiKey')),
				ipWhitelist: computed(() => ts('pages.datasources.http.drawer.authTypes.ipWhitelist')),
				jwt: computed(() => ts('pages.datasources.http.drawer.authTypes.jwt')),
				oauth2: computed(() => ts('pages.datasources.http.drawer.authTypes.oauth2')),
			},
		},
	};
}

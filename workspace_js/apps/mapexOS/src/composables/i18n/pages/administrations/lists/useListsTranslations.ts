import type { DataRowColumn } from '@components/cards';
import type { PageHeaderInfo } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useListsTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	return {
		/**
		 * CreateEditListPage header translations (title/description/buttons)
		 */
		page: {
			title: computed(() => tsTitle('pages.administrations.lists.page.titleCreate')),
			titleCreate: computed(() => tsTitle('pages.administrations.lists.page.titleCreate')),
			titleEdit: computed(() => tsTitle('pages.administrations.lists.page.titleEdit')),
			descriptionCreate: computed(() => ts('pages.administrations.lists.page.descriptionCreate')),
			backButton: computed(() => ts('pages.administrations.lists.page.backButton')),
			addButton: computed(() => ts('pages.administrations.lists.page.addButton')),
		},

		pageHeader: {
			title: computed(() => tsTitle('pages.administrations.lists.pageHeader.title')),
			description: computed(() => ts('pages.administrations.lists.pageHeader.description')),
			button: computed(() => ts('pages.administrations.lists.pageHeader.button')),
			info: computed((): PageHeaderInfo => ({
				title: ts('pages.administrations.lists.pageHeader.info.title'),
				description: ts('pages.administrations.lists.pageHeader.info.description'),
				items: [
					{
						icon: 'category',
						color: 'blue-6',
						title: ts('pages.administrations.lists.pageHeader.info.items.types.title'),
						text: ts('pages.administrations.lists.pageHeader.info.items.types.text'),
					},
					{
						icon: 'account_tree',
						color: 'green-6',
						title: ts('pages.administrations.lists.pageHeader.info.items.hierarchy.title'),
						text: ts('pages.administrations.lists.pageHeader.info.items.hierarchy.text'),
					},
					{
						icon: 'lock',
						color: 'purple-6',
						title: ts('pages.administrations.lists.pageHeader.info.items.systemLists.title'),
						text: ts('pages.administrations.lists.pageHeader.info.items.systemLists.text'),
					},
					{
						icon: 'tune',
						color: 'indigo-6',
						title: ts('pages.administrations.lists.pageHeader.info.items.customLists.title'),
						text: ts('pages.administrations.lists.pageHeader.info.items.customLists.text'),
					},
					{
						icon: 'content_copy',
						color: 'orange-6',
						title: ts('pages.administrations.lists.pageHeader.info.items.templates.title'),
						text: ts('pages.administrations.lists.pageHeader.info.items.templates.text'),
					},
				],
				docsUrl: 'https://docs.mapexos.com/lists',
				docsLabel: ts('pages.administrations.lists.pageHeader.info.docsLabel'),
			})),
		},

		menuColumns: {
			organization: computed(() => ts('pages.administrations.lists.menuColumns.organization')),
			parent: computed(() => ts('pages.administrations.lists.menuColumns.parent')),
			type: computed(() => ts('pages.administrations.lists.menuColumns.type')),
			source: computed(() => ts('pages.administrations.lists.menuColumns.source')),
			scope: computed(() => ts('pages.administrations.lists.menuColumns.scope')),
		},

		filters: {
			label: computed(() => ts('pages.administrations.lists.filters.label')),
			searchPlaceholder: computed(() => ts('pages.administrations.lists.filters.searchPlaceholder')),
			allSource: computed(() => ts('pages.administrations.lists.filters.allSource')),
			advancedFilters: computed(() => ts('pages.administrations.lists.filters.advancedFilters')),
			pendingFilters: computed(() => ts('pages.administrations.lists.filters.pendingFilters')),
			clearAll: computed(() => ts('pages.administrations.lists.filters.clearAll')),
			name: computed(() => ts('pages.administrations.lists.filters.name')),
			category: computed(() => ts('pages.administrations.lists.filters.category')),
			manufacturer: computed(() => ts('pages.administrations.lists.filters.manufacturer')),
			type: computed(() => ts('pages.administrations.lists.filters.type')),
			isSystem: computed(() => ts('pages.administrations.lists.filters.isSystem')),
			isTemplate: computed(() => ts('pages.administrations.lists.filters.isTemplate')),
			includeChildren: computed(() => ts('pages.administrations.lists.filters.includeChildren')),
			options: {
				all: computed(() => ts('pages.administrations.lists.filters.categoryOptions.all')),
				iot: computed(() => ts('pages.administrations.lists.filters.categoryOptions.iot')),
				categories: computed(() => ts('pages.administrations.lists.filters.typeOptions.categories')),
				manufacturers: computed(() => ts('pages.administrations.lists.filters.typeOptions.manufacturers')),
				models: computed(() => ts('pages.administrations.lists.filters.typeOptions.models')),
				system: computed(() => ts('pages.administrations.lists.filters.isSystemOptions.system')),
				custom: computed(() => ts('pages.administrations.lists.filters.isSystemOptions.custom')),
				shared: computed(() => ts('pages.administrations.lists.filters.isTemplateOptions.shared')),
				local: computed(() => ts('pages.administrations.lists.filters.isTemplateOptions.local')),
				yes: computed(() => ts('pages.administrations.lists.filters.includeChildrenOptions.yes')),
				no: computed(() => ts('pages.administrations.lists.filters.includeChildrenOptions.no')),
			},
		},

		listHeader: {
			title: computed(() => tsTitle('pages.administrations.lists.listHeader.title')),
			itemLabel: computed(() => ts('pages.administrations.lists.listHeader.itemLabel')),
			itemLabelPlural: computed(() => ts('pages.administrations.lists.listHeader.itemLabelPlural')),
		},

		columns: computed((): DataRowColumn[] => [
			{
				key: 'icon',
				label: '',
				type: 'avatar',
				visible: 'always',
				width: 56,
				icon: (value: any, row: any) => {
					const iconMap: Record<string, string> = {
						assetType: 'devices',
						assetGroup: 'folder',
						priority: 'priority_high',
						status: 'toggle_on',
						severity: 'warning',
					};
					return iconMap[row.type] || 'list';
				},
				// List items have no deactivation flow in the UI yet, so the
				// chip/icon always renders the active styling and "Ativo" tooltip.
				color: () => 'primary',
				tooltip: () => ts('pages.administrations.lists.status.active'),
			},
			{
				key: 'name',
				label: ts('pages.administrations.lists.columns.name'),
				type: 'text',
				visible: 'always',
				width: 250,
				ellipsis: true,
				secondaryKey: 'value',
			},
			{
				key: 'parentName',
				label: ts('pages.administrations.lists.columns.parent'),
				type: 'chip',
				visible: 'laptop',
				width: 180,
				ellipsis: true,
				format: (value: any) => value || '—',
				color: () => 'indigo-6',
				icon: 'account_tree',
			},
			{
				key: 'type',
				label: ts('pages.administrations.lists.columns.type'),
				type: 'chip',
				visible: 'laptop',
				width: 120,
				format: (value: any) => {
					if (!value) return '—';
					return ts(`pages.administrations.lists.listTypes.${value}`).toUpperCase();
				},
				color: (value: any) => {
					const colorMap: Record<string, string> = {
						manufacturers: 'blue-6',
						assets: 'green-6',
					};
					return colorMap[value] || 'grey-6';
				},
				icon: 'category',
			},
			{
				key: 'organizationName',
				label: ts('pages.administrations.lists.columns.organization'),
				type: 'chip',
				visible: 'laptop',
				width: 180,
				ellipsis: true,
				color: 'indigo-6',
				icon: 'domain',
			},
			{
				key: 'isSystem',
				label: ts('pages.administrations.lists.columns.source'),
				type: 'chip',
				visible: 'laptop',
				width: 120,
				format: (value: any) => value
					? ts('pages.administrations.lists.filters.isSystemOptions.system').toUpperCase()
					: ts('pages.administrations.lists.filters.isSystemOptions.custom').toUpperCase(),
				color: (value: any) => value ? 'orange-6' : 'green-6',
				icon: (value: any) => value ? 'lock' : 'edit',
			},
			{
				key: 'isTemplate',
				label: ts('pages.administrations.lists.columns.scope'),
				type: 'chip',
				visible: 'laptop',
				width: 150,
				format: (value: any) => value
					? ts('pages.administrations.lists.filters.isTemplateOptions.shared').toUpperCase()
					: ts('pages.administrations.lists.filters.isTemplateOptions.local').toUpperCase(),
				color: (value: any) => value ? 'purple-6' : 'grey-6',
				icon: (value: any) => value ? 'content_copy' : 'folder',
			},
		]),

		empty: {
			title: computed(() => ts('pages.administrations.lists.empty.title')),
			description: computed(() => ts('pages.administrations.lists.empty.description')),
		},

		dialog: {
			confirmDelete: {
				title: computed(() => ts('pages.administrations.lists.dialog.confirmDelete.title')),
				message: (name: string) => ts('pages.administrations.lists.dialog.confirmDelete.message', { name }),
			},
		},

		notifications: {
			deleted: computed(() => ts('pages.administrations.lists.notifications.deleted')),
			deleteError: computed(() => ts('pages.administrations.lists.notifications.deleteError')),
			loadFailed: computed(() => ts('pages.administrations.lists.notifications.loadFailed')),
			systemListEdit: computed(() => ts('pages.administrations.lists.notifications.systemListEdit')),
			systemListDelete: computed(() => ts('pages.administrations.lists.notifications.systemListDelete')),
			sharedListEdit: computed(() => ts('pages.administrations.lists.notifications.sharedListEdit')),
			sharedListDelete: computed(() => ts('pages.administrations.lists.notifications.sharedListDelete')),
		},

		errors: {
			apiNotInitialized: computed(() => ts('pages.administrations.lists.errors.apiNotInitialized')),
			idMissing: computed(() => ts('pages.administrations.lists.errors.idMissing')),
		},

		actions: {
			edit: computed(() => ts('pages.administrations.lists.actions.edit')),
			view: computed(() => ts('pages.administrations.lists.actions.view')),
			delete: computed(() => ts('pages.administrations.lists.actions.delete')),
		},

		/**
		 * CreateEditListPage - stepper labels
		 */
		steps: {
			basicInfo: computed(() => tsTitle('pages.administrations.lists.createEdit.steps.basicInfo')),
			review: computed(() => tsTitle('pages.administrations.lists.createEdit.steps.review')),
		},

		stepDescriptions: {
			basicInfo: computed(() => ts('pages.administrations.lists.createEdit.stepDescriptions.basicInfo')),
			review: computed(() => ts('pages.administrations.lists.createEdit.stepDescriptions.review')),
		},

		sections: {
			basicInfo: computed(() => tsTitle('pages.administrations.lists.createEdit.sections.basicInfo')),
			progressSteps: computed(() => tsTitle('pages.administrations.lists.createEdit.sections.progressSteps')),
		},

		formDescriptions: {
			basicInfo: computed(() => ts('pages.administrations.lists.createEdit.formDescriptions.basicInfo')),
			isTemplate: computed(() => ts('pages.administrations.lists.createEdit.formDescriptions.isTemplate')),
		},

		/**
		 * CreateEditListPage - field labels
		 */
		fields: {
			type: computed(() => ts('pages.administrations.lists.createEdit.fields.type')),
			parent: computed(() => ts('pages.administrations.lists.createEdit.fields.parent')),
			name: computed(() => ts('pages.administrations.lists.createEdit.fields.name')),
			value: computed(() => ts('pages.administrations.lists.createEdit.fields.value')),
			enabled: computed(() => ts('pages.administrations.lists.createEdit.fields.enabled')),
			isTemplate: computed(() => ts('pages.administrations.lists.createEdit.fields.isTemplate')),
			status: {
				label: computed(() => ts('pages.administrations.lists.createEdit.fields.status.label')),
				placeholder: computed(() => ts('pages.administrations.lists.createEdit.fields.status.placeholder')),
				hint: computed(() => ts('pages.administrations.lists.createEdit.fields.status.hint')),
				required: computed(() => ts('pages.administrations.lists.createEdit.fields.status.required')),
			},
		},

		typeOptions: {
			asset_category: computed(() => ts('pages.administrations.lists.createEdit.typeOptions.asset_category')),
			asset_manufacturer: computed(() => ts('pages.administrations.lists.createEdit.typeOptions.asset_manufacturer')),
			asset_model: computed(() => ts('pages.administrations.lists.createEdit.typeOptions.asset_model')),
		},

		statusOptions: {
			active: {
				label: computed(() => ts('pages.administrations.lists.createEdit.statusOptions.active.label')),
				value: true,
			},
			inactive: {
				label: computed(() => ts('pages.administrations.lists.createEdit.statusOptions.inactive.label')),
				value: false,
			},
		},

		hints: {
			type: computed(() => ts('pages.administrations.lists.createEdit.hints.type')),
			parentForManufacturer: computed(() => ts('pages.administrations.lists.createEdit.hints.parentForManufacturer')),
			parentForModel: computed(() => ts('pages.administrations.lists.createEdit.hints.parentForModel')),
			name: computed(() => ts('pages.administrations.lists.createEdit.hints.name')),
			value: computed(() => ts('pages.administrations.lists.createEdit.hints.value')),
			isTemplate: computed(() => ts('pages.administrations.lists.createEdit.hints.isTemplate')),
		},

		validation: {
			typeRequired: computed(() => ts('pages.administrations.lists.createEdit.validation.typeRequired')),
			parentRequired: computed(() => ts('pages.administrations.lists.createEdit.validation.parentRequired')),
			nameRequired: computed(() => ts('pages.administrations.lists.createEdit.validation.nameRequired')),
			nameMaxLength: computed(() => ts('pages.administrations.lists.createEdit.validation.nameMaxLength')),
			valueRequired: computed(() => ts('pages.administrations.lists.createEdit.validation.valueRequired')),
			valueMaxLength: computed(() => ts('pages.administrations.lists.createEdit.validation.valueMaxLength')),
			valuePattern: computed(() => ts('pages.administrations.lists.createEdit.validation.valuePattern')),
		},

		buttons: {
			back: computed(() => ts('pages.administrations.lists.createEdit.buttons.back')),
			next: computed(() => ts('pages.administrations.lists.createEdit.buttons.next')),
			createList: computed(() => ts('pages.administrations.lists.createEdit.buttons.createList')),
			updateList: computed(() => ts('pages.administrations.lists.createEdit.buttons.updateList')),
		},

		messages: {
			loadingList: computed(() => ts('pages.administrations.lists.createEdit.messages.loadingList')),
			completeAllSteps: computed(() => ts('pages.administrations.lists.createEdit.messages.completeAllSteps')),
			allFieldsRequired: computed(() => ts('pages.administrations.lists.createEdit.messages.allFieldsRequired')),
			currentStep: computed(() => ts('pages.administrations.lists.createEdit.messages.currentStep')),
			systemListWarning: computed(() => ts('pages.administrations.lists.createEdit.messages.systemListWarning')),
			typeLocked: computed(() => ts('pages.administrations.lists.createEdit.messages.typeLocked')),
			noParentForCategory: computed(() => ts('pages.administrations.lists.createEdit.messages.noParentForCategory')),
			parentLoadFailed: computed(() => ts('pages.administrations.lists.createEdit.messages.parentLoadFailed')),
		},

		reviewStep: {
			subtitle: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.subtitle')),
			successMessage: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.successMessage')),
			successMessageEdit: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.successMessageEdit')),
			sections: {
				basicInfo: computed(() => tsTitle('pages.administrations.lists.createEdit.reviewStep.sections.basicInfo')),
			},
			fields: {
				type: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.type')),
				parent: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.parent')),
				name: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.name')),
				value: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.value')),
				enabled: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.enabled')),
				isTemplate: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.fields.isTemplate')),
			},
			values: {
				notProvided: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.values.notProvided')),
				notSelected: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.values.notSelected')),
				noParent: computed(() => ts('pages.administrations.lists.createEdit.reviewStep.values.noParent')),
			},
		},

		createEditNotifications: {
			created: computed(() => ts('pages.administrations.lists.createEdit.notifications.created')),
			updated: computed(() => ts('pages.administrations.lists.createEdit.notifications.updated')),
			createFailed: computed(() => ts('pages.administrations.lists.createEdit.notifications.createFailed')),
			updateFailed: computed(() => ts('pages.administrations.lists.createEdit.notifications.updateFailed')),
			loadFailed: computed(() => ts('pages.administrations.lists.createEdit.notifications.loadFailed')),
			alreadyExists: computed(() => ts('pages.administrations.lists.createEdit.notifications.alreadyExists')),
			forbidden: computed(() => ts('pages.administrations.lists.createEdit.notifications.forbidden')),
		},
	};
}

import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Composable for HTTP Data Source CreateEdit page translations
 * Provides all translated strings for the create/edit form
 */
export function useHttpDataSourceCreateEditTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	const basePath = 'pages.datasources.http.createEdit';

	return {
		/** Page header translations */
		page: {
			title: computed(() => tsTitle(`${basePath}.page.title`)),
			titleEdit: computed(() => tsTitle(`${basePath}.page.titleEdit`)),
			description: computed(() => ts(`${basePath}.page.description`)),
			descriptionEdit: computed(() => ts(`${basePath}.page.descriptionEdit`)),
			backButton: computed(() => ts(`${basePath}.page.backButton`)),
			loading: computed(() => ts(`${basePath}.page.loading`)),
		},

		/** Stepper translations */
		stepper: {
			title: computed(() => ts(`${basePath}.stepper.title`)),
			subtitle: computed(() => ts(`${basePath}.stepper.subtitle`)),
			infoText: computed(() => ts(`${basePath}.stepper.infoText`)),
			currentStepLabel: computed(() => ts(`${basePath}.stepper.currentStepLabel`)),
		},

		/** Steps configuration */
		steps: computed(() => [
			{
				title: ts(`${basePath}.steps.basicInfo.title`),
				icon: 'info',
				description: ts(`${basePath}.steps.basicInfo.description`),
			},
			{
				title: ts(`${basePath}.steps.workingHours.title`),
				icon: 'schedule',
				description: ts(`${basePath}.steps.workingHours.description`),
			},
			{
				title: ts(`${basePath}.steps.authentication.title`),
				icon: 'lock',
				description: ts(`${basePath}.steps.authentication.description`),
			},
			{
				title: ts(`${basePath}.steps.assetBinding.title`),
				icon: 'device_unknown',
				description: ts(`${basePath}.steps.assetBinding.description`),
			},
			{
				title: ts(`${basePath}.steps.review.title`),
				icon: 'check_circle',
				description: ts(`${basePath}.steps.review.description`),
			},
		]),

		/** Step progress side panel labels */
		progress: {
			completeAllSteps: computed(() => ts(`${basePath}.steps.progress.completeAllSteps`)),
		},

		/** Protocol config (Step 2) banner labels */
		protocolConfig: {
			gatewayBannerTitle: computed(() => ts(`${basePath}.steps.protocolConfig.gatewayBannerTitle`)),
			gatewayBannerBody: computed(() => ts(`${basePath}.steps.protocolConfig.gatewayBannerBody`)),
		},

		/** Navigation button labels */
		navigation: {
			back: computed(() => ts(`${basePath}.navigation.back`)),
			next: computed(() => ts(`${basePath}.navigation.next`)),
			create: computed(() => ts(`${basePath}.navigation.create`)),
			update: computed(() => ts(`${basePath}.navigation.update`)),
		},

		/** Notification messages */
		notifications: {
			createSuccess: computed(() => ts(`${basePath}.notifications.createSuccess`)),
			updateSuccess: computed(() => ts(`${basePath}.notifications.updateSuccess`)),
			createFailed: computed(() => ts(`${basePath}.notifications.createFailed`)),
			updateFailed: computed(() => ts(`${basePath}.notifications.updateFailed`)),
			loadFailed: computed(() => ts(`${basePath}.notifications.loadFailed`)),
			copiedToClipboard: computed(() => ts(`${basePath}.notifications.copiedToClipboard`)),
			copyFailed: computed(() => ts(`${basePath}.notifications.copyFailed`)),
		},

		/** Module-level error messages */
		errors: {
			apiNotInitialized: computed(() => ts('pages.datasources.http.errors.apiNotInitialized')),
			noUuidPaths: computed(() => ts('pages.datasources.http.errors.noUuidPaths')),
			invalidPayloadOrPath: computed(() => ts('pages.datasources.http.errors.invalidPayloadOrPath')),
		},

		/** Step 1: Basic Information */
		basicInfo: {
			name: computed(() => ts(`${basePath}.basicInfo.name`)),
			nameRequired: computed(() => ts(`${basePath}.basicInfo.nameRequired`)),
			description: computed(() => ts(`${basePath}.basicInfo.description`)),
			status: computed(() => ts(`${basePath}.basicInfo.status`)),
			statusOptions: computed(() => [
				{ label: ts(`${basePath}.basicInfo.statusOptions.enabled`), value: true },
				{ label: ts(`${basePath}.basicInfo.statusOptions.disabled`), value: false },
			]),
			debugMode: computed(() => ts(`${basePath}.basicInfo.debugMode`)),
			debugModeOptions: computed(() => [
				{ label: ts(`${basePath}.basicInfo.statusOptions.enabled`), value: true },
				{ label: ts(`${basePath}.basicInfo.statusOptions.disabled`), value: false },
			]),
			debugModeTooltip: computed(() => ts(`${basePath}.basicInfo.debugModeTooltip`)),
		},

		/** Step 2: Working Hours */
		workingHours: {
			title: computed(() => ts(`${basePath}.workingHours.title`)),
			enableWorkingHours: computed(() => ts(`${basePath}.workingHours.enableWorkingHours`)),
			daysOfWeek: {
				monday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.monday`)),
				tuesday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.tuesday`)),
				wednesday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.wednesday`)),
				thursday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.thursday`)),
				friday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.friday`)),
				saturday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.saturday`)),
				sunday: computed(() => ts(`${basePath}.workingHours.daysOfWeek.sunday`)),
			},
			timeIntervals: computed(() => ts(`${basePath}.workingHours.timeIntervals`)),
			startTime: computed(() => ts(`${basePath}.workingHours.startTime`)),
			endTime: computed(() => ts(`${basePath}.workingHours.endTime`)),
			startTimeRequired: computed(() => ts(`${basePath}.workingHours.startTimeRequired`)),
			endTimeRequired: computed(() => ts(`${basePath}.workingHours.endTimeRequired`)),
			invalidTimeFormat: computed(() => ts(`${basePath}.workingHours.invalidTimeFormat`)),
			removeInterval: computed(() => ts(`${basePath}.workingHours.removeInterval`)),
			addInterval: computed(() => ts(`${basePath}.workingHours.addInterval`)),
			timezone: computed(() => ts(`${basePath}.workingHours.timezone`)),
			timezonePlaceholder: computed(() => ts(`${basePath}.workingHours.timezonePlaceholder`)),
			timezoneRequired: computed(() => ts(`${basePath}.workingHours.timezoneRequired`)),
		},

		/** Step 2: Rate Limit */
		rateLimit: {
			title: computed(() => ts(`${basePath}.rateLimit.title`)),
			enableRateLimit: computed(() => ts(`${basePath}.rateLimit.enableRateLimit`)),
			limitType: computed(() => ts(`${basePath}.rateLimit.limitType`)),
			limitTypeRequired: computed(() => ts(`${basePath}.rateLimit.limitTypeRequired`)),
			value: computed(() => ts(`${basePath}.rateLimit.value`)),
			valueRequired: computed(() => ts(`${basePath}.rateLimit.valueRequired`)),
			burstCapacity: computed(() => ts(`${basePath}.rateLimit.burstCapacity`)),
			actionOnExceed: computed(() => ts(`${basePath}.rateLimit.actionOnExceed`)),
			actionRequired: computed(() => ts(`${basePath}.rateLimit.actionRequired`)),
		},

		/** Step 3: Authentication */
		authentication: {
			title: computed(() => ts(`${basePath}.authentication.title`)),
			authTypeLabel: computed(() => ts(`${basePath}.authentication.authTypeLabel`)),
			authTypeRequired: computed(() => ts(`${basePath}.authentication.authTypeRequired`)),

			apiKey: {
				banner: computed(() => ts(`${basePath}.authentication.apiKey.banner`)),
				headerName: computed(() => ts(`${basePath}.authentication.apiKey.headerName`)),
				headerNamePlaceholder: computed(() => ts(`${basePath}.authentication.apiKey.headerNamePlaceholder`)),
				headerNameHint: computed(() => ts(`${basePath}.authentication.apiKey.headerNameHint`)),
				headerNameRequired: computed(() => ts(`${basePath}.authentication.apiKey.headerNameRequired`)),
				value: computed(() => ts(`${basePath}.authentication.apiKey.value`)),
				valuePlaceholder: computed(() => ts(`${basePath}.authentication.apiKey.valuePlaceholder`)),
				valueHint: computed(() => ts(`${basePath}.authentication.apiKey.valueHint`)),
				valueRequired: computed(() => ts(`${basePath}.authentication.apiKey.valueRequired`)),
				copyTooltip: computed(() => ts(`${basePath}.authentication.apiKey.copyTooltip`)),
				generateTooltip: computed(() => ts(`${basePath}.authentication.apiKey.generateTooltip`)),
			},

			jwt: {
				banner: computed(() => ts(`${basePath}.authentication.jwt.banner`)),
				headerName: computed(() => ts(`${basePath}.authentication.jwt.headerName`)),
				headerNamePlaceholder: computed(() => ts(`${basePath}.authentication.jwt.headerNamePlaceholder`)),
				headerNameHint: computed(() => ts(`${basePath}.authentication.jwt.headerNameHint`)),
				headerNameRequired: computed(() => ts(`${basePath}.authentication.jwt.headerNameRequired`)),
				secretKey: computed(() => ts(`${basePath}.authentication.jwt.secretKey`)),
				secretKeyPlaceholder: computed(() => ts(`${basePath}.authentication.jwt.secretKeyPlaceholder`)),
				secretKeyHint: computed(() => ts(`${basePath}.authentication.jwt.secretKeyHint`)),
				secretKeyRequired: computed(() => ts(`${basePath}.authentication.jwt.secretKeyRequired`)),
				copyTooltip: computed(() => ts(`${basePath}.authentication.jwt.copyTooltip`)),
				generateTooltip: computed(() => ts(`${basePath}.authentication.jwt.generateTooltip`)),
			},

			ipWhitelist: {
				banner: computed(() => ts(`${basePath}.authentication.ipWhitelist.banner`)),
				ipAddress: computed(() => ts(`${basePath}.authentication.ipWhitelist.ipAddress`)),
				ipAddressPlaceholder: computed(() => ts(`${basePath}.authentication.ipWhitelist.ipAddressPlaceholder`)),
				ipAddressHint: computed(() => ts(`${basePath}.authentication.ipWhitelist.ipAddressHint`)),
				ipAddressRequired: computed(() => ts(`${basePath}.authentication.ipWhitelist.ipAddressRequired`)),
				cidrMask: computed(() => ts(`${basePath}.authentication.ipWhitelist.cidrMask`)),
				cidrMaskPlaceholder: computed(() => ts(`${basePath}.authentication.ipWhitelist.cidrMaskPlaceholder`)),
				cidrMaskHint: computed(() => ts(`${basePath}.authentication.ipWhitelist.cidrMaskHint`)),
				cidrMaskRequired: computed(() => ts(`${basePath}.authentication.ipWhitelist.cidrMaskRequired`)),
				addButton: computed(() => ts(`${basePath}.authentication.ipWhitelist.addButton`)),
				removeTooltip: computed(() => ts(`${basePath}.authentication.ipWhitelist.removeTooltip`)),
				emptyMessage: computed(() => ts(`${basePath}.authentication.ipWhitelist.emptyMessage`)),
				addedMessage: (ip: string) => ts(`${basePath}.authentication.ipWhitelist.addedMessage`, { ip }),
				removedMessage: (ip: string) => ts(`${basePath}.authentication.ipWhitelist.removedMessage`, { ip }),
				duplicateError: computed(() => ts(`${basePath}.authentication.ipWhitelist.duplicateError`)),
				invalidIpv4Cidr: computed(() => ts(`${basePath}.authentication.ipWhitelist.invalidIpv4Cidr`)),
				invalidIpv6Cidr: computed(() => ts(`${basePath}.authentication.ipWhitelist.invalidIpv6Cidr`)),
				invalidIpFormat: computed(() => ts(`${basePath}.authentication.ipWhitelist.invalidIpFormat`)),
			},

			oauth2: {
				banner: computed(() => ts(`${basePath}.authentication.oauth2.banner`)),
				jwksUrl: computed(() => ts(`${basePath}.authentication.oauth2.jwksUrl`)),
				jwksUrlPlaceholder: computed(() => ts(`${basePath}.authentication.oauth2.jwksUrlPlaceholder`)),
				jwksUrlHint: computed(() => ts(`${basePath}.authentication.oauth2.jwksUrlHint`)),
				jwksUrlRequired: computed(() => ts(`${basePath}.authentication.oauth2.jwksUrlRequired`)),
				invalidUrl: computed(() => ts(`${basePath}.authentication.oauth2.invalidUrl`)),
			},

			none: {
				banner: computed(() => ts(`${basePath}.authentication.none.banner`)),
			},
		},

		/** Step 4: Asset Binding */
		assetBinding: {
			title: computed(() => ts(`${basePath}.assetBinding.title`)),
			bindingMode: computed(() => ts(`${basePath}.assetBinding.bindingMode`)),
			bindingModeRequired: computed(() => ts(`${basePath}.assetBinding.bindingModeRequired`)),

			fixedAsset: {
				banner: computed(() => ts(`${basePath}.assetBinding.fixedAsset.banner`)),
				selectAsset: computed(() => ts(`${basePath}.assetBinding.fixedAsset.selectAsset`)),
			},

			uuidField: {
				banner: computed(() => ts(`${basePath}.assetBinding.uuidField.banner`)),
				uuidJsonPath: computed(() => ts(`${basePath}.assetBinding.uuidField.uuidJsonPath`)),
				uuidJsonPathPlaceholder: computed(() => ts(`${basePath}.assetBinding.uuidField.uuidJsonPathPlaceholder`)),
				uuidJsonPathHint: computed(() => ts(`${basePath}.assetBinding.uuidField.uuidJsonPathHint`)),
				pathRequired: computed(() => ts(`${basePath}.assetBinding.uuidField.pathRequired`)),
				invalidPathFormat: computed(() => ts(`${basePath}.assetBinding.uuidField.invalidPathFormat`)),
				selectFromTemplate: computed(() => ts(`${basePath}.assetBinding.uuidField.selectFromTemplate`)),
				noTemplatePath: computed(() => ts(`${basePath}.assetBinding.uuidField.noTemplatePath`)),
				templateSearchPlaceholder: computed(() => ts(`${basePath}.assetBinding.uuidField.templateSearchPlaceholder`)),
				removePath: computed(() => ts(`${basePath}.assetBinding.uuidField.removePath`)),
				addPath: computed(() => ts(`${basePath}.assetBinding.uuidField.addPath`)),
			},
		},

		/** Step 5: Review */
		review: {
			subtitle: computed(() => ts(`${basePath}.review.subtitle`)),
			successMessage: computed(() => ts(`${basePath}.review.successMessage`)),
			sections: {
				basicInfo: computed(() => ts(`${basePath}.review.sections.basicInfo`)),
				workingHours: computed(() => ts(`${basePath}.review.sections.workingHours`)),
				authentication: computed(() => ts(`${basePath}.review.sections.authentication`)),
				assetBinding: computed(() => ts(`${basePath}.review.sections.assetBinding`)),
			},
			fields: {
				name: computed(() => ts(`${basePath}.review.fields.name`)),
				description: computed(() => ts(`${basePath}.review.fields.description`)),
				status: computed(() => ts(`${basePath}.review.fields.status`)),
				workingHoursEnabled: computed(() => ts(`${basePath}.review.fields.workingHoursEnabled`)),
				daysOfWeek: computed(() => ts(`${basePath}.review.fields.daysOfWeek`)),
				timeInterval: computed(() => ts(`${basePath}.review.fields.timeInterval`)),
				timezone: computed(() => ts(`${basePath}.review.fields.timezone`)),
				rateLimitEnabled: computed(() => ts(`${basePath}.review.fields.rateLimitEnabled`)),
				rateLimitType: computed(() => ts(`${basePath}.review.fields.rateLimitType`)),
				rateLimitValue: computed(() => ts(`${basePath}.review.fields.rateLimitValue`)),
				burstCapacity: computed(() => ts(`${basePath}.review.fields.burstCapacity`)),
				actionOnExceed: computed(() => ts(`${basePath}.review.fields.actionOnExceed`)),
				authType: computed(() => ts(`${basePath}.review.fields.authType`)),
				apiKeyHeader: computed(() => ts(`${basePath}.review.fields.apiKeyHeader`)),
				apiKeyValue: computed(() => ts(`${basePath}.review.fields.apiKeyValue`)),
				jwtSecret: computed(() => ts(`${basePath}.review.fields.jwtSecret`)),
				ipAddresses: computed(() => ts(`${basePath}.review.fields.ipAddresses`)),
				jwksUrl: computed(() => ts(`${basePath}.review.fields.jwksUrl`)),
				bindingMode: computed(() => ts(`${basePath}.review.fields.bindingMode`)),
				selectedAsset: computed(() => ts(`${basePath}.review.fields.selectedAsset`)),
				uuidPaths: computed(() => ts(`${basePath}.review.fields.uuidPaths`)),
			},
			values: {
				enabled: computed(() => ts(`${basePath}.review.values.enabled`)),
				disabled: computed(() => ts(`${basePath}.review.values.disabled`)),
				notConfigured: computed(() => ts(`${basePath}.review.values.notConfigured`)),
				none: computed(() => ts(`${basePath}.review.values.none`)),
				noDescription: computed(() => ts(`${basePath}.review.values.noDescription`)),
				fixedAsset: computed(() => ts(`${basePath}.review.values.fixedAsset`)),
				uuidField: computed(() => ts(`${basePath}.review.values.uuidField`)),
				perSecond: computed(() => ts(`${basePath}.review.values.perSecond`)),
				perMinute: computed(() => ts(`${basePath}.review.values.perMinute`)),
				perHour: computed(() => ts(`${basePath}.review.values.perHour`)),
				drop: computed(() => ts(`${basePath}.review.values.drop`)),
				queue: computed(() => ts(`${basePath}.review.values.queue`)),
			},
			basicInfo: computed(() => ts(`${basePath}.review.basicInfo`)),
			name: computed(() => ts(`${basePath}.review.name`)),
			enabled: computed(() => ts(`${basePath}.review.enabled`)),
			yes: computed(() => ts(`${basePath}.review.yes`)),
			no: computed(() => ts(`${basePath}.review.no`)),
			description: computed(() => ts(`${basePath}.review.description`)),
			notAvailable: computed(() => ts(`${basePath}.review.notAvailable`)),
			protocolMode: computed(() => ts(`${basePath}.review.protocolMode`)),
			mode: computed(() => ts(`${basePath}.review.mode`)),
			protocol: computed(() => ts(`${basePath}.review.protocol`)),
			pushFixed: computed(() => ts(`${basePath}.review.pushFixed`)),
			httpFixed: computed(() => ts(`${basePath}.review.httpFixed`)),
			protocolBanner: computed(() => ts(`${basePath}.review.protocolBanner`)),
			authentication: computed(() => ts(`${basePath}.review.authentication`)),
			authType: computed(() => ts(`${basePath}.review.authType`)),
			headerApiKey: computed(() => ts(`${basePath}.review.headerApiKey`)),
			secretKey: computed(() => ts(`${basePath}.review.secretKey`)),
			notConfigured: computed(() => ts(`${basePath}.review.notConfigured`)),
			ipWhitelist: computed(() => ts(`${basePath}.review.ipWhitelist`)),
			noIpConfigured: computed(() => ts(`${basePath}.review.noIpConfigured`)),
			jwksUrl: computed(() => ts(`${basePath}.review.jwksUrl`)),
			noAuthWarning: computed(() => ts(`${basePath}.review.noAuthWarning`)),
			assetBinding: computed(() => ts(`${basePath}.review.assetBinding`)),
			bindingMode: computed(() => ts(`${basePath}.review.bindingMode`)),
			selectedAsset: computed(() => ts(`${basePath}.review.selectedAsset`)),
			uuidPaths: computed(() => ts(`${basePath}.review.uuidPaths`)),
			noUuidPaths: computed(() => ts(`${basePath}.review.noUuidPaths`)),
			assetTemplates: computed(() => ts(`${basePath}.review.assetTemplates`)),
			templatesSelected: (count: number) => ts(`${basePath}.review.templatesSelected`, { count }),
			noTemplatesSelected: computed(() => ts(`${basePath}.review.noTemplatesSelected`)),
			payloadExample: computed(() => ts(`${basePath}.review.payloadExample`)),
		},
	};
}

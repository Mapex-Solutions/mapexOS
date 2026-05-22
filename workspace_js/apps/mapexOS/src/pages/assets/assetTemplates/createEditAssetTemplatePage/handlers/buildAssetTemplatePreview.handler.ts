import type { ReviewSectionDef } from '@components/forms/review/interfaces';
import type { AssetTemplateData } from '../interfaces';

/**
 * Builds the "Basic Information" section for review
 */
function buildBasicInformationSection(data: AssetTemplateData): ReviewSectionDef {
	return {
		stepNumber: 1,
		label: 'Basic Information',
		icon: { name: 'info', color: 'primary' },
		fields: [
			{ label: 'Name', value: data.name, type: 'text', colSize: 6 },
			{
				label: 'Status',
				value: data.enabled ? 'Active' : 'Inactive',
				type: 'badge',
				badgeColors: { Active: 'positive', Inactive: 'negative' },
				colSize: 6,
			},
			{ label: 'Description', value: data.description || 'No description provided', type: 'text', colSize: 12 },
			{ label: 'Category', value: data.categoryName || 'Not specified', type: 'text', colSize: 3 },
			{ label: 'Manufacturer', value: data.manufacturerName || 'Not specified', type: 'text', colSize: 3 },
			{ label: 'Model', value: data.modelName || 'Not specified', type: 'text', colSize: 3 },
			{ label: 'Version', value: data.version || 'Not specified', type: 'text', colSize: 3 },
		],
	};
}

/**
 * Builds the "Asset ID Path" section for review
 */
function buildAssetIdPathSection(data: AssetTemplateData): ReviewSectionDef {
	return {
		stepNumber: 2,
		label: 'Asset ID Path',
		icon: { name: 'route', color: 'secondary' },
		fields: [
			{
				label: 'Asset ID Path',
				value: data.assetIdPath,
				type: 'text',
				colSize: 12,
			},
		],
	};
}

/**
 * Builds the "Scripts Summary" section for review
 */
function buildScriptsSummarySection(data: AssetTemplateData): ReviewSectionDef {
	return {
		stepNumber: 3,
		label: 'Scripts Summary',
		icon: { name: 'code', color: 'primary' },
		fields: [
			{
				label: 'Preprocessor Script',
				value: data.scriptProcessor ? 'Configured' : 'Not configured',
				type: 'badge',
				badgeColors: { Configured: 'positive', 'Not configured': 'grey' },
				colSize: 6,
			},
			{
				label: 'Validation Script',
				value: data.scriptValidator ? 'Configured' : 'Not configured',
				type: 'badge',
				badgeColors: { Configured: 'positive', 'Not configured': 'negative' },
				colSize: 6,
			},
			{
				label: 'Conversion Script',
				value: data.scriptConversion ? 'Configured' : 'Not configured',
				type: 'badge',
				badgeColors: { Configured: 'positive', 'Not configured': 'negative' },
				colSize: 6,
			},
			{
				label: 'Test Script',
				value: data.scriptTest ? 'Configured' : 'Not configured',
				type: 'badge',
				badgeColors: { Configured: 'positive', 'Not configured': 'grey' },
				colSize: 6,
			},
		],
	};
}

/**
 * Main builder that composes all sections for FormReview
 */
export function buildAssetTemplatePreview(data: AssetTemplateData): ReviewSectionDef[] {
	return [
		buildBasicInformationSection(data),
		buildAssetIdPathSection(data),
		buildScriptsSummarySection(data),
	];
}

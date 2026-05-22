import type { StepperVerticalItem } from '@components/steppers';

export const STEPS: StepperVerticalItem[] = [
	{
		title: 'Data Lake Information',
		icon: 'mdi-database',
		description: 'Add your data lake name, status, and description',
	},
	{
		title: 'Select Data Lake Type',
		icon: 'mdi-cloud-outline',
		description: 'Choose between AWS, Azure, GCP, or MinIO',
	},
	{
		title: 'Access Credentials',
		icon: 'mdi-key',
		description: 'Enter keys and secrets for authentication',
	},
	{
		title: 'Path Configuration',
		icon: 'mdi-folder-outline',
		description: 'Define the path where data will be stored',
	},
	{
		title: 'Export Frequency',
		icon: 'mdi-calendar-clock',
		description: 'Set the interval for exporting data',
	},
	{
		title: 'Review & Confirm',
		icon: 'mdi-clipboard-check',
		description: 'Verify all details before saving',
	},
];
import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useListDrawerTranslations() {
	const ts = useTS({ capitalize: true });

	return {
		drawer: {
			title: computed(() => ts('pages.administrations.lists.drawer.title')),
			close: computed(() => ts('pages.administrations.lists.drawer.close')),
			loading: computed(() => ts('pages.administrations.lists.drawer.loading')),
			error: computed(() => ts('pages.administrations.lists.drawer.error')),

			sections: {
				basicInfo: computed(() => ts('pages.administrations.lists.drawer.sections.basicInfo')),
				ids: computed(() => ts('pages.administrations.lists.drawer.sections.ids')),
				metadata: computed(() => ts('pages.administrations.lists.drawer.sections.metadata')),
				timestamps: computed(() => ts('pages.administrations.lists.drawer.sections.timestamps')),
			},

			fields: {
				name: computed(() => ts('pages.administrations.lists.drawer.fields.name')),
				value: computed(() => ts('pages.administrations.lists.drawer.fields.value')),
				category: computed(() => ts('pages.administrations.lists.drawer.fields.category')),
				type: computed(() => ts('pages.administrations.lists.drawer.fields.type')),
				isSystem: computed(() => ts('pages.administrations.lists.drawer.fields.isSystem')),
				isTemplate: computed(() => ts('pages.administrations.lists.drawer.fields.isTemplate')),
				id: computed(() => ts('pages.administrations.lists.drawer.fields.id')),
				parentType: computed(() => ts('pages.administrations.lists.drawer.fields.parentType')),
				parentName: computed(() => ts('pages.administrations.lists.drawer.fields.parentName')),
				orgId: computed(() => ts('pages.administrations.lists.drawer.fields.orgId')),
				created: computed(() => ts('pages.administrations.lists.drawer.fields.created')),
				updated: computed(() => ts('pages.administrations.lists.drawer.fields.updated')),
			},

			values: {
				yes: computed(() => ts('pages.administrations.lists.drawer.values.yes')),
				no: computed(() => ts('pages.administrations.lists.drawer.values.no')),
			},

			empty: {
				parentType: computed(() => ts('pages.administrations.lists.drawer.empty.parentType')),
				parentName: computed(() => ts('pages.administrations.lists.drawer.empty.parentName')),
				orgId: computed(() => ts('pages.administrations.lists.drawer.empty.orgId')),
			},
		},
	};
}

<script setup lang="ts">
defineOptions({
	name: 'TabInfo'
});

/** TYPE IMPORTS */
import type { TabInfoProps } from '../../interfaces';

/** COMPONENTS */
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useGroupDetailTranslations } from '@composables/i18n';

/** PROPS */
defineProps<TabInfoProps>();

/** COMPOSABLES */
const t = useGroupDetailTranslations();

/**
 * Format date for display
 *
 * @param {string | undefined} dateString - ISO date string
 * @returns {string} Formatted date
 */
function formatDate(dateString?: string): string {
	if (!dateString) return '—';
	return new Date(dateString).toLocaleDateString('pt-BR', {
		day: '2-digit',
		month: '2-digit',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
	});
}
</script>

<template>
	<div class="tab-info">
		<!-- Loading State -->
		<div v-if="loading" class="flex flex-center q-pa-lg">
			<q-spinner color="primary" size="40px" />
		</div>

		<!-- Content -->
		<div v-else-if="group">
			<!-- Basic Info Section -->
			<div class="section q-mb-lg">
				<div class="section-header q-mb-md">
					<q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
					<span class="text-subtitle1 text-weight-medium">{{ t.info.basicInfo.title.value }}</span>
				</div>

				<div class="row q-col-gutter-md">
					<!-- Name -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.name.value }}</div>
							<div class="info-value text-weight-medium">{{ group.name || '—' }}</div>
						</div>
					</div>

					<!-- Description -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.description.value }}</div>
							<div class="info-value">{{ group.description || '—' }}</div>
						</div>
					</div>

					<!-- Status -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.status.value }}</div>
							<div class="info-value">
								<DetailChip
									:label="group.enabled ? t.info.status.enabled.value : t.info.status.disabled.value"
									:color="group.enabled ? 'green' : 'red'"
									size="sm"
								/>
							</div>
						</div>
					</div>

					<!-- Members Count -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.membersCount.value }}</div>
							<div class="info-value">
								<DetailChip
									:label="String(group.membersCount ?? 0)"
									color="blue"
									icon="people"
									size="sm"
								/>
							</div>
						</div>
					</div>

					<!-- Organization -->
					<div v-if="group.organizationName" class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.organization.value }}</div>
							<div class="info-value">{{ group.organizationName }}</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Dates Section -->
			<div class="section">
				<div class="section-header q-mb-md">
					<q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
					<span class="text-subtitle1 text-weight-medium">{{ t.info.dates.title.value }}</span>
				</div>

				<div class="row q-col-gutter-md">
					<!-- Created -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.created.value }}</div>
							<div class="info-value">{{ formatDate(group.created) }}</div>
						</div>
					</div>

					<!-- Updated -->
					<div class="col-12 col-md-6">
						<div class="info-item">
							<div class="info-label text-grey-7">{{ t.info.fields.updated.value }}</div>
							<div class="info-value">{{ formatDate(group.updated) }}</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- No Data -->
		<div v-else class="text-center q-pa-lg text-grey-6">
			<q-icon name="warning" size="48px" class="q-mb-md" />
			<div>{{ t.info.noData.value }}</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.tab-info {
	.section {
		background: var(--mapex-surface-bg);
		border-radius: var(--mapex-radius-md);
		padding: 16px;
	}

	.section-header {
		display: flex;
		align-items: center;
		border-bottom: 1px solid var(--mapex-card-border);
		padding-bottom: 8px;
	}

	.info-item {
		.info-label {
			font-size: 0.75rem;
			text-transform: uppercase;
			letter-spacing: 0.5px;
			margin-bottom: 4px;
		}

		.info-value {
			font-size: 0.95rem;
		}
	}
}
</style>

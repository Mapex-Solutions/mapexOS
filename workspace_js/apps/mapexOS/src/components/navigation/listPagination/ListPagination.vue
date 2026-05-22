<script setup lang="ts">
defineOptions({
  name: 'ListPagination'
});

import type { ListPaginationProps, ListPaginationEmits } from './interfaces';

import { computed } from 'vue';
import { useQuasar } from 'quasar';

const props = withDefaults(defineProps<ListPaginationProps>(), {
	color: 'primary',
	activeColor: 'primary',
});

const emit = defineEmits<ListPaginationEmits>();

const $q = useQuasar();

/**
 * Compute max visible pages based on screen size
 * Mobile: 3 pages
 * Tablet: 5 pages
 * Laptop: 7 pages
 * Desktop: 10 pages
 */
const maxPages = computed(() => {
	if ($q.screen.lt.sm) {
		// Mobile (< 600px)
		return 3;
	} else if ($q.screen.lt.md) {
		// Tablet (600px - 1024px)
		return 5;
	} else if ($q.screen.lt.lg) {
		// Laptop (1024px - 1440px)
		return 7;
	} else {
		// Desktop (>= 1440px)
		return 10;
	}
});

/**
 * Show pagination only if there's more than 1 page
 */
const shouldShowPagination = computed(() => props.totalPages > 1);

/**
 * Handle page change
 */
function handlePageChange(page: number) {
	emit('update:modelValue', page);
	emit('change', page);
}
</script>

<template>
	<div v-if="shouldShowPagination" class="row justify-center q-mt-lg q-mb-lg">
		<q-pagination
			:model-value="modelValue"
			direction-links
			boundary-links
			:color="color"
			:active-color="activeColor"
			class="rounded-borders"
			:max="totalPages"
			:max-pages="maxPages"
			@update:model-value="handlePageChange"
		/>
	</div>
</template>

<style lang="scss" scoped>
// Optional: Add custom styles if needed
</style>

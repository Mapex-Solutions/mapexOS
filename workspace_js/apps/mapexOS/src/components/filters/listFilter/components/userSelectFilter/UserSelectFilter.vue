<script setup lang="ts">
defineOptions({
	name: 'UserSelectFilter'
});

/** TYPE IMPORTS */
import type { UserSelectFilterProps, UserSelectFilterEmits, SelectedUserInfo } from './interfaces';

/** VUE IMPORTS */
import { ref, watch, onMounted } from 'vue';

/** COMPONENTS */
import { UserSelectorDrawer } from '@components/drawers';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = withDefaults(defineProps<UserSelectFilterProps>(), {
	clearable: true,
	disabled: false,
	placeholder: 'Click to select...',
});

const emit = defineEmits<UserSelectFilterEmits>();

/** STATE */
const showDrawer = ref(false);
const selectedUser = ref<SelectedUserInfo | null>(null);
const loadingUser = ref(false);

/** FUNCTIONS */

/**
 * Open the user selector drawer
 */
function openDrawer(): void {
	if (props.disabled) return;
	showDrawer.value = true;
}

/**
 * Handle user selection from drawer
 *
 * @param {any} user - Selected user object
 */
function handleUserSelect(user: any): void {
	const userInfo: SelectedUserInfo = {
		id: user.id,
		name: getUserDisplayName(user),
	};
	if (user.email) {
		userInfo.email = user.email;
	}
	selectedUser.value = userInfo;
	emit('update:modelValue', user.id);
}

/**
 * Clear the selected user
 */
function clearSelection(): void {
	selectedUser.value = null;
	emit('update:modelValue', null);
}

/**
 * Get user display name from user object
 *
 * @param {any} user - User object
 * @returns {string} Display name
 */
function getUserDisplayName(user: any): string {
	if (user.firstName || user.lastName) {
		return `${user.firstName || ''} ${user.lastName || ''}`.trim();
	}
	return user.email || 'Unknown User';
}

/**
 * Fetch user details by ID to display name
 *
 * @param {string} userId - User ID to fetch
 */
async function fetchUserDetails(userId: string): Promise<void> {
	if (!apis.mapexOS?.users) return;

	loadingUser.value = true;
	try {
		const user = await apis.mapexOS.users.getById({ userId });
		const userInfo: SelectedUserInfo = {
			id: userId,
			name: getUserDisplayName(user),
		};
		if (user.email) {
			userInfo.email = user.email;
		}
		selectedUser.value = userInfo;
	} catch {
		// If fetch fails, just show the ID
		selectedUser.value = {
			id: userId,
			name: userId,
		};
	} finally {
		loadingUser.value = false;
	}
}

/** WATCHERS */

// Watch for external modelValue changes (e.g., reset)
watch(() => props.modelValue, (newValue) => {
	if (!newValue) {
		selectedUser.value = null;
	} else if (!selectedUser.value || selectedUser.value.id !== newValue) {
		// Fetch user details if ID changed externally
		void fetchUserDetails(newValue);
	}
}, { immediate: true });

/** LIFECYCLE HOOKS */
onMounted(() => {
	// If there's an initial value, fetch user details
	if (props.modelValue) {
		void fetchUserDetails(props.modelValue);
	}
});
</script>

<template>
	<div class="user-select-filter">
		<!-- Clickable Wrapper -->
		<div
			class="field-wrapper"
			:class="{ 'field-wrapper--disabled': disabled }"
			@click="openDrawer"
		>
			<q-field
				outlined
				dense
				readonly
				class="rounded-borders"
				:label="label"
				:disable="disabled"
				stack-label
			>
				<template #prepend>
					<q-icon color="primary" :name="icon" />
				</template>

				<template #control>
					<div class="self-center full-width no-outline row items-center">
						<!-- Loading State -->
						<template v-if="loadingUser">
							<q-spinner size="xs" color="grey-6" class="q-mr-sm" />
							<span class="text-grey-6">Loading...</span>
						</template>

						<!-- Selected User -->
						<template v-else-if="selectedUser">
							<span class="text-body2">{{ selectedUser.name }}</span>
						</template>

						<!-- Placeholder -->
						<template v-else>
							<span class="text-grey-6">{{ placeholder }}</span>
						</template>
					</div>
				</template>

				<template v-if="clearable && selectedUser && !disabled" #append>
					<q-icon
						name="cancel"
						class="cursor-pointer clear-icon"
						color="grey-6"
						@click.stop="clearSelection"
					/>
				</template>
			</q-field>
		</div>

		<!-- User Selector Drawer -->
		<UserSelectorDrawer
			v-model="showDrawer"
			:selected-user-id="modelValue"
			@select="handleUserSelect"
		/>
	</div>
</template>

<style lang="scss" scoped>
.user-select-filter {
	.field-wrapper {
		cursor: pointer;

		&--disabled {
			cursor: not-allowed;
			opacity: 0.6;
			pointer-events: none;
		}

		// Make q-field not capture clicks, let wrapper handle it
		:deep(.q-field) {
			pointer-events: none;
		}

		// But allow the clear icon to be clickable
		:deep(.q-field__append) {
			pointer-events: auto;
		}
	}

	.rounded-borders {
		border-radius: var(--mapex-radius-md);
	}

	.clear-icon {
		cursor: pointer;
		z-index: 1;

		&:hover {
			color: var(--q-negative) !important;
		}
	}

	:deep(.q-field__control) {
		cursor: pointer;
	}

	:deep(.q-field--outlined .q-field__control) {
		border-radius: var(--mapex-radius-md);
	}
}
</style>

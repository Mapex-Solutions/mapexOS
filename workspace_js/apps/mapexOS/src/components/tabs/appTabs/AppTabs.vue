<script setup lang="ts">
defineOptions({
  name: 'AppTabs'
});

/** TYPE IMPORTS */
import type { AppTabsProps, AppTabsEmits, AppTabItem } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AppTabsProps>(), {
  bordered: false,
  align: 'left',
  separator: true,
  variant: 'default',
});

const emit = defineEmits<AppTabsEmits>();

/** COMPUTED */

/**
 * Active tab value for v-model binding
 */
const activeTab = computed({
  get: () => props.modelValue,
  set: (value: string) => {
    emit('update:modelValue', value);
    emit('change', value);
  },
});

/**
 * Whether to use pill variant styling
 */
const isPill = computed(() => props.variant === 'pill');

/** FUNCTIONS */

/**
 * Get badge color for a tab item
 *
 * @param {AppTabItem} tab - Tab configuration
 * @returns {string} Quasar color name
 */
function getBadgeColor(tab: AppTabItem): string {
  return tab.badgeColor || 'primary';
}

/**
 * Check if tab should show badge
 *
 * @param {AppTabItem} tab - Tab configuration
 * @returns {boolean} Whether to show badge
 */
function shouldShowBadge(tab: AppTabItem): boolean {
  return tab.badge !== undefined && tab.badge > 0;
}
</script>

<template>
  <div class="app-tabs" :class="{ 'app-tabs--bordered': bordered, 'app-tabs--pill': isPill }">
    <q-tabs
      v-model="activeTab"
      :narrow-indicator="!isPill"
      :class="isPill ? 'app-tabs__pill-tabs' : 'app-tabs__default-tabs'"
      :indicator-color="isPill ? 'transparent' : 'primary'"
      :align="align"
      no-caps
    >
      <q-tab
        v-for="tab in tabs"
        :key="tab.name"
        :id="tab.id"
        :name="tab.name"
        :icon="isPill ? tab.icon : undefined"
        :label="isPill ? tab.label : undefined"
        :disable="tab.disabled"
        :class="isPill ? 'app-tabs__pill-tab' : 'app-tabs__default-tab'"
      >
        <!-- Custom template for default variant (icon left of label) -->
        <template v-if="!isPill" #default>
          <div class="row items-center no-wrap">
            <q-icon
              v-if="tab.icon"
              size="22px"
              class="q-mr-sm"
              :name="tab.icon"
            />
            <div class="tab-label">{{ tab.label }}</div>
            <q-badge
              v-if="shouldShowBadge(tab)"
              :color="getBadgeColor(tab)"
              text-color="white"
              rounded
              class="q-ml-sm"
              style="font-size: 10px; padding: 2px 6px;"
            >
              {{ tab.badge }}
            </q-badge>
          </div>
        </template>
      </q-tab>
    </q-tabs>

    <q-separator v-if="separator && !isPill" />

    <!-- Slot for tab panels -->
    <slot />
  </div>
</template>

<style scoped lang="scss">
.app-tabs {
  &--bordered {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    overflow: hidden;
  }

  // ==========================================
  // DEFAULT VARIANT - Enterprise Style
  // ==========================================
  .app-tabs__default-tabs {
    background: var(--mapex-page-bg);
    border-bottom: 2px solid var(--mapex-divider);

    :deep(.q-tabs__content) {
      padding: 0 8px;
    }

    :deep(.q-tab__indicator) {
      height: 3px;
      border-radius: var(--mapex-radius-xs) var(--mapex-radius-xs) 0 0;
      bottom: -2px;
    }
  }

  .app-tabs__default-tab {
    min-height: 52px;
    padding: 0 20px;
    font-weight: 500;
    color: var(--mapex-text-secondary);
    transition: var(--mapex-transition-base);
    border-radius: var(--mapex-radius-md) var(--mapex-radius-md) 0 0;
    margin: 4px 2px 0;

    .tab-label {
      font-size: 14px;
      font-weight: 500;
      letter-spacing: 0.01em;
    }

    &:hover:not(.q-tab--active) {
      background: var(--mapex-surface-bg);
      color: var(--mapex-text-primary);
    }

    :deep(.q-tab__content) {
      min-width: unset;
    }

    :deep(.q-focus-helper) {
      border-radius: var(--mapex-radius-md) var(--mapex-radius-md) 0 0;
    }
  }

  :deep(.q-tab--active.app-tabs__default-tab) {
    color: var(--q-primary);
    background: var(--mapex-surface-elevated);
    box-shadow: 0 -1px 3px var(--mapex-elevation-shadow);

    .tab-label {
      font-weight: 600;
    }
  }

  // ==========================================
  // PILL VARIANT - Segmented Control Style
  // ==========================================
  &--pill {
    .app-tabs__pill-tabs {
      background: transparent;

      :deep(.q-tabs__content) {
        padding: 4px;
        background: var(--mapex-page-bg);
        border-radius: var(--mapex-radius-md);
      }

      :deep(.q-tab__indicator) {
        display: none;
      }
    }

    .app-tabs__pill-tab {
      min-height: 56px;
      padding: 8px 16px;
      border-radius: var(--mapex-radius-sm);
      margin-right: 4px;
      color: var(--mapex-text-primary);

      &:last-child {
        margin-right: 0;
      }

      &:hover:not(.q-tab--active) {
        background: var(--mapex-surface-bg);
      }

      :deep(.q-tab__content) {
        min-width: unset;
      }
    }

    :deep(.q-tab--active) {
      background: var(--mapex-surface-elevated);
      box-shadow: 0 1px 3px var(--mapex-elevation-shadow);
      color: var(--q-primary);
    }
  }
}
</style>

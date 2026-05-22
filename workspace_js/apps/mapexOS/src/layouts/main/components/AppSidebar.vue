<script setup lang="ts">
/** TYPE IMPORTS */
import type { MenuItem } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';
import { useRoute } from 'vue-router';

/** COMPOSABLES */
import { useMainLayoutTranslations } from '@composables/i18n';

/** PROPS & EMITS */
defineProps<{
  isOpen: boolean;
  miniState: boolean;
  menuList: MenuItem[];
}>();

const emit = defineEmits<{
  'update:isOpen': [value: boolean];
}>();

/** COMPOSABLES & STORES */
const route = useRoute();
const t = useMainLayoutTranslations();

/** COMPUTED */

/**
 * Current route path
 */
const currentRoute = computed(() => route.path);

/** FUNCTIONS */

/**
 * Get child menu item icon with fallback
 *
 * @param {MenuItem} child - Child menu item
 * @returns {string} Icon name
 */
function getChildIcon(child: MenuItem): string {
  return child.icon || 'chevron_right';
}
</script>

<template>
  <q-drawer
    bordered
    show-if-above
    class="sidebar-drawer"
    :style="{ background: 'var(--mapex-sidebar-bg)' }"
    :model-value="isOpen"
    :mini="miniState"
    :mini-width="70"
    :width="240"
    :breakpoint="500"
    @update:model-value="emit('update:isOpen', $event)"
  >
    <!-- Logo Section -->
    <div class="column items-center q-py-md logo-container">
      <q-img
        v-if="!miniState"
        class="full-logo"
        src="/mapex-logo.png"
        width="120"
        height="50"
        fit="contain"
      />
      <q-img
        v-else
        class="mini-logo"
        src="/only-logo.png"
        width="40"
        height="40"
        fit="contain"
      />
    </div>

    <q-separator/>

    <!-- Navigation Menu -->
    <q-list
      id="sidebar-menu"
      padding
      class="menu-list"
    >
      <template v-for="(item, index) in menuList" :key="index">
        <!-- Items with children -->
        <template v-if="item.children">
          <q-expansion-item
            v-if="!miniState"
            dense
            expand-icon="keyboard_arrow_down"
            expand-icon-class="text-grey-7"
            :icon="item.icon"
            :label="item.label"
            :default-opened="currentRoute.startsWith(item.children[0]?.to ?? '')"
            :header-class="'menu-parent-item'"
            :icon-color="currentRoute.startsWith(item.children[0]?.to ?? '') ? 'primary' : 'grey-7'"
          >
            <q-list class="submenu-list">
              <template v-for="(child, cIndex) in item.children" :key="cIndex">
                <!-- Separator -->
                <q-separator v-if="child.separator" class="q-my-xs" />
                <!-- Menu Item -->
                <q-item
                  v-else
                  v-ripple
                  clickable
                  active-class="active-menu-item"
                  class="submenu-item"
                  :to="child.to"
                  :active="currentRoute === child.to"
                >
                  <q-item-section avatar>
                    <q-icon
                      :name="getChildIcon(child)"
                      :color="currentRoute === child.to ? 'primary' : 'grey-7'"
                    />
                  </q-item-section>
                  <q-item-section>{{ child.label }}</q-item-section>
                </q-item>
              </template>
            </q-list>
          </q-expansion-item>

          <q-item
            v-else
            clickable
            active-class="active-menu-item"
            class="mini-menu-item"
            :active="currentRoute.startsWith(item.children[0]?.to ?? '')"
          >
            <q-item-section avatar>
              <q-icon
                :name="item.icon"
                :color="currentRoute.startsWith(item.children[0]?.to ?? '') ? 'primary' : 'grey-7'"
              />
            </q-item-section>

            <q-menu
              v-if="miniState"
              class="submenu-popup"
              anchor="bottom right"
              self="top start"
              :offset="[0, 0]"
            >
              <q-list class="submenu-list">
                <q-item-label header>{{ item.label }}</q-item-label>
                <q-separator class="q-my-sm"/>
                <template v-for="(child, cIndex) in item.children" :key="cIndex">
                  <!-- Separator -->
                  <q-separator v-if="child.separator" class="q-my-xs" />
                  <!-- Menu Item -->
                  <q-item
                    v-else
                    v-ripple
                    v-close-popup
                    clickable
                    active-class="active-menu-item"
                    :to="child.to"
                    :active="currentRoute === child.to"
                  >
                    <q-item-section avatar>
                      <q-icon
                        :name="getChildIcon(child)"
                        :color="currentRoute === child.to ? 'primary' : 'grey-7'"
                      />
                    </q-item-section>
                    <q-item-section>{{ child.label }}</q-item-section>
                  </q-item>
                </template>
              </q-list>
            </q-menu>
          </q-item>
        </template>

        <!-- Regular items without children -->
        <q-item
          v-else
          v-ripple
          clickable
          active-class="active-menu-item"
          class="regular-menu-item"
          :to="item.to"
          :active="currentRoute === item.to"
          :class="{ 'mini-regular-item': miniState }"
        >
          <q-item-section avatar>
            <q-icon
              :name="item.icon"
              :color="currentRoute === item.to ? 'primary' : 'grey-7'"
            />
          </q-item-section>

          <q-item-section v-if="!miniState">
            <q-item-label>{{ item.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
    </q-list>

    <!-- Version Info -->
    <div class="absolute-bottom q-pa-sm text-center">
      <div class="version-info">
        <q-icon
          class="q-mr-xs"
          name="info"
          size="xs"
          color="primary"
        />
        <span>{{ t.version.value }} 1.0.0</span>
      </div>
    </div>
  </q-drawer>
</template>

<style lang="scss" scoped>
// Sidebar styles
.sidebar-drawer {
  transition: var(--mapex-transition-slow);

  .logo-container {
    height: 80px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: var(--mapex-transition-slow);

    .full-logo, .mini-logo {
      transition: opacity 0.3s ease;
    }
  }
}

// Menu list styling
.menu-list {
  height: calc(100vh - 140px);
  overflow-y: auto;
  padding-bottom: 60px;

  &::-webkit-scrollbar {
    width: 5px;
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--mapex-scrollbar-thumb);
    border-radius: var(--mapex-radius-sm);

    &:hover {
      background: var(--mapex-scrollbar-thumb-hover);
    }
  }

  .regular-menu-item, :deep(.menu-parent-item) {
    border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
    margin: 4px 0;
    transition: var(--mapex-transition-base);

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.05);
    }
  }

  // Expansion item icons - respect the icon-color prop or default to grey-7
  :deep(.q-expansion-item) {
    // Default grey for all expansion item icons
    .q-expansion-item__container > .q-item .q-item__section--side > .q-icon {
      color: var(--mapex-text-secondary) !important;
    }

    // Override with primary when icon-color is set to primary
    &[class*="text-primary"] .q-expansion-item__container > .q-item .q-item__section--side > .q-icon,
    .q-expansion-item__container > .q-item .q-item__section--side > .q-icon.text-primary {
      color: var(--q-primary) !important;
    }
  }

  .mini-menu-item {
    padding: 8px 0;
    justify-content: center;

    .q-item__section--avatar {
      min-width: 0;
      padding-left: 0;
      justify-content: center;
    }
  }

  .mini-regular-item {
    padding: 12px 0;
    justify-content: center;

    .q-item__section--avatar {
      min-width: 0;
      padding-left: 0;
      justify-content: center;
    }
  }
}

// Active menu styling
.active-menu-item {
  background: rgba(var(--q-primary-rgb), 0.1);
  border-right: 3px solid var(--q-primary);
  color: var(--q-primary);
}

// Submenu styling
.submenu-list {
  padding-left: 16px !important;
  background: var(--mapex-submenu-bg);

  .q-item {
    padding: 8px 16px;
    min-height: 40px;
    border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
    margin: 2px 0;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.05);
    }
  }
}

.submenu-popup {
  min-width: 200px;
  background: var(--mapex-popup-bg);
  border-radius: var(--mapex-radius-md);
  box-shadow: 0 4px 20px var(--mapex-elevation-shadow);
}

// Version info styling
.version-info {
  opacity: 0.7;
  font-size: 0.8rem;
  transition: var(--mapex-transition-base);

  .q-icon {
    font-size: 1rem;
    vertical-align: middle;
  }
}
</style>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { MenuItem } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import AppHeader from './components/AppHeader.vue';
import AppSidebar from './components/AppSidebar.vue';
import { OrganizationTreeDrawer } from '@components/drawers';

/** COMPOSABLES */
import { useOnboarding } from '@composables/onboarding/useOnboarding';
import { useMainLayoutTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';
import { useTheme } from '@composables/theme';

/** UTILS */
import { isBelowMinResolution, isResolutionWarningDismissed, dismissResolutionWarning } from '@utils/screen';

/** STORES */
import { usePermissionStore } from '@stores/permission';

/** LOCAL IMPORTS */
import { buildMenuList } from './constants';

/** COMPOSABLES & STORES */
const router = useRouter();
const logger = useLogger('MainLayout');
const permStore = usePermissionStore();
const { startTour, shouldAutoStart } = useOnboarding();
const layoutT = useMainLayoutTranslations();
useTheme();

/**
 * Filter menu items based on user permissions.
 * Items without permissions are always visible.
 * Parent groups are hidden if no visible children remain.
 *
 * @param {MenuItem[]} items - Menu items to filter
 * @returns {MenuItem[]} Filtered menu items
 */
function filterByPermissions(items: MenuItem[]): MenuItem[] {
  return items.reduce<MenuItem[]>((acc, item) => {
    // Separators pass through
    if (item.separator) {
      acc.push(item);
      return acc;
    }

    // Check parent-level permission
    if (item.permissions?.length && !permStore.hasAnyPermission(item.permissions)) {
      return acc;
    }

    // Filter children if present
    if (item.children) {
      const filteredChildren = filterByPermissions(item.children);
      // Only show parent if it has at least one non-separator child
      const hasVisibleChild = filteredChildren.some(c => !c.separator);
      if (hasVisibleChild) {
        acc.push({ ...item, children: filteredChildren });
      }
      return acc;
    }

    acc.push(item);
    return acc;
  }, []);
}

/** COMPUTED */

/**
 * Menu filtered by current user permissions
 */
const filteredMenu = computed(() => filterByPermissions(buildMenuList(layoutT.menu)));

/** STATE */
const drawerOpen = ref(false);
const miniState = ref(true);
const showOrgTreeDrawer = ref(false);
const showResolutionWarning = ref(false);
const screenWidth = ref(window.innerWidth);
const screenHeight = ref(window.innerHeight);

/** FUNCTIONS */

/**
 * Dismiss resolution warning and persist to localStorage
 */
function handleDismissResolutionWarning(): void {
  showResolutionWarning.value = false;
  dismissResolutionWarning();
}

/**
 * Toggle drawer based on screen size
 */
function toggleDrawer() {
  if (window.innerWidth <= 1024) {
    drawerOpen.value = !drawerOpen.value;
  } else {
    miniState.value = !miniState.value;
  }
}

/**
 * Open organization tree drawer
 */
function openOrgTreeDrawer() {
  showOrgTreeDrawer.value = !showOrgTreeDrawer.value;
}

/**
 * Handle organization selection
 */
function onOrganizationSelected(orgId: string) {
  logger.debug('Organization selected:', orgId);
}

/**
 * Handle start tour from header menu
 * Navigates to dashboard first if not already there, then starts the tour
 */
async function handleStartTour(): Promise<void> {
  const currentPath = router.currentRoute.value.path;

  // Navigate to home first if not already there
  if (currentPath !== '/home' && currentPath !== '/') {
    logger.debug('Navigating to home before starting tour');
    await router.push('/home');
    // Wait for navigation and render to complete
    await nextTick();
    // Small delay to ensure page is fully rendered
    await new Promise(resolve => setTimeout(resolve, 300));
  }

  startTour();
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  // Check screen resolution and show warning if below HD and not dismissed
  if (isBelowMinResolution() && !isResolutionWarningDismissed()) {
    showResolutionWarning.value = true;
  }

  if (shouldAutoStart()) {
    void nextTick(() => startTour());
  }
});
</script>

<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Header -->
    <AppHeader
      @toggle-drawer="toggleDrawer"
      @open-org-tree="openOrgTreeDrawer"
      @start-tour="handleStartTour"
    />

    <!-- Sidebar -->
    <AppSidebar
      v-model:is-open="drawerOpen"
      :mini-state="miniState"
      :menu-list="filteredMenu"
    />

    <!-- Page Container -->
    <q-page-container id="page-container" class="q-my-md">
      <!-- Resolution Warning Banner -->
      <q-banner
        v-if="showResolutionWarning"
        dense
        inline-actions
        class="resolution-banner text-white"
      >
        <template #avatar>
          <q-icon name="monitor" color="white" />
        </template>
        For the best experience, we recommend a minimum screen resolution of <strong>1366 x 768 (HD)</strong>.
        Your current resolution is <strong>{{ screenWidth }} x {{ screenHeight }}</strong>.
        <template #action>
          <q-btn
            flat
            dense
            label="OK"
            color="white"
            @click="handleDismissResolutionWarning"
          />
        </template>
      </q-banner>

      <div class="container">
        <router-view/>
      </div>
    </q-page-container>

    <!-- Organization Tree Drawer -->
    <OrganizationTreeDrawer
      v-model="showOrgTreeDrawer"
      @select="onOrganizationSelected"
    />
  </q-layout>
</template>

<style lang="scss">
// Resolution warning banner (inside q-page-container, respects sidebar)
.resolution-banner {
  background: linear-gradient(135deg, var(--q-warning) 0%, var(--q-warning) 100%);
  font-size: 13px;
}

// Container styles
.container {
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
  padding: 0 16px;
}

@media (min-width: 1024px) {
  .container {
    padding: 0 24px;
  }
}

@media (min-width: 1440px) {
  .container {
    padding: 0 32px;
  }
}

// Global styles
body {
  background-color: transparent;
}

// Menu parent item styling (global for expansion items)
:deep(.menu-parent-item) {
  border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
  margin: 4px 0;
  transition: var(--mapex-transition-base);
  padding: 8px 16px;

  &:hover {
    background: rgba(var(--q-primary-rgb), 0.05);
  }
}

// Submenu list styling (global for expansion items)
:deep(.q-expansion-item__content) {
  background: var(--mapex-submenu-bg);
}

// Icon sizing and default colors
:deep(.q-item__section--side > .q-icon) {
  font-size: 24px;
}

// Ensure expansion item icons respect the icon-color prop
:deep(.q-expansion-item__toggle-icon) {
  font-size: 24px;
}
</style>

<script setup lang="ts">
/** VUE IMPORTS */
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useMainLayoutTranslations } from '@composables/i18n';
import { useBreadcrumbs } from '../composables';
import { useTheme } from '@composables/theme';

/** LOCAL IMPORTS */
import { buildMenuList } from '../constants';

/** UTILS */
import { notifyInfo } from '@utils/alert/notify';
import { getInitials } from '@utils/user';

/** SERVICES */
import { useAuth } from '@services/auth';

/** STORES */
import { useOrganizationStore } from '@stores/organization';
import { useAuthStore } from '@stores/auth';

/** EMITS */
const emit = defineEmits<{
  'toggle-drawer': [];
  'open-org-tree': [];
  'start-tour': [];
}>();

/** COMPOSABLES & STORES */
const auth = useAuth();
const authStore = useAuthStore();
const { locale } = useI18n();
const t = useMainLayoutTranslations();
const organizationStore = useOrganizationStore();
const translatedMenu = computed(() => buildMenuList(t.menu));
const { breadcrumbs } = useBreadcrumbs(translatedMenu);
const { isDark, toggleTheme } = useTheme();

/** COMPUTED */

/**
 * User initials computed from auth store (e.g. "TA" for Thiago Anselmo)
 */
const userInitials = computed(() =>
  getInitials(authStore.user?.firstName, authStore.user?.lastName, authStore.user?.email)
);

/**
 * Selected organization from store
 */
const selectedOrganization = computed(() => organizationStore.selectedOrganization);

/**
 * Language options with reactive labels
 */
const languageList = computed(() => [
  { value: 'en-US', label: t.languageSelector.languages.english.value, icon: '/flag-usa.svg' },
  { value: 'pt-BR', label: t.languageSelector.languages.portuguese.value, icon: '/flag-br.svg' },
]);

/**
 * Current locale
 */
const currentLocale = computed(() => locale.value);

/** FUNCTIONS */

/**
 * Change application language and persist to localStorage
 *
 * @param {string} lang - Locale code (e.g. 'en-US', 'pt-BR')
 */
function changeLanguage(lang: string): void {
  locale.value = lang;
  localStorage.setItem('user-locale', lang);
  const langLabel = languageList.value.find(l => l.value === lang)?.label;
  notifyInfo({
    message: `${t.languageSelector.changed.value} ${langLabel}`,
  });
}

/**
 * Handle user logout
 */
function handleLogout(): void {
  auth.logout();
}
</script>

<template>
  <q-header
    elevated
    class="header-container"
  >
    <q-toolbar class="q-px-lg q-py-sm">
      <!-- Menu Toggle and Separator -->
      <div class="row items-center">
        <q-btn
          flat
          dense
          round
          no-wrap
          class="q-mr-sm"
          icon="menu"
          color="primary"
          aria-label="Menu"
          @click="emit('toggle-drawer')"
        />
        <q-separator
          vertical
          class="q-mx-sm"
        />
      </div>

      <!-- Breadcrumbs -->
      <q-breadcrumbs
        id="header-breadcrumbs"
        class="text-primary"
        active-color="primary"
        separator-color="grey-6"
      >
        <q-breadcrumbs-el
          icon="dashboard"
          to="/home"
          :label="breadcrumbs.length === 0 ? t.breadcrumbs.home.value : undefined"
        />
        <q-breadcrumbs-el
          v-for="(crumb, index) in breadcrumbs"
          :key="index"
          :label="crumb.label"
          :icon="crumb.icon"
          :to="crumb.to"
        />
      </q-breadcrumbs>

      <q-space/>

      <!-- Organization Indicator -->
      <div v-if="selectedOrganization" id="header-org-selector" class="q-mr-md org-indicator">
        <q-chip
          clickable
          outline
          class="org-chip"
          color="primary"
          icon="business"
          @click="emit('open-org-tree')"
        >
          <span class="text-weight-medium org-name">{{ selectedOrganization.name }}</span>
          <AppTooltip>
            <div class="column">
              <span class="text-weight-bold">{{ selectedOrganization.name }}</span>
              <span class="text-caption q-mt-xs">{{ t.orgIndicator.tooltip.value }}</span>
            </div>
          </AppTooltip>
        </q-chip>
      </div>

      <!-- Action Buttons -->
      <div class="row items-center q-gutter-sm">

        <!-- Theme Toggle -->
        <q-btn
          id="header-theme-toggle"
          flat
          round
          :icon="isDark ? 'light_mode' : 'dark_mode'"
          :color="isDark ? 'amber' : 'primary'"
          @click="toggleTheme"
        >
          <AppTooltip>
            {{ isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode' }}
          </AppTooltip>
        </q-btn>

        <!-- Language Selector -->
        <q-btn
          id="header-lang-selector"
          flat
          round
          color="primary"
        >
          <q-icon name="language"/>
          <q-menu
            :class="isDark ? 'bg-dark' : 'bg-white'"
            anchor="bottom right"
            self="top right"
          >
            <q-list style="min-width: 150px">
              <q-item
                v-for="lang in languageList"
                v-close-popup
                clickable
                active-class="bg-primary text-white"
                :key="lang.value"
                :active="currentLocale === lang.value"
                @click="changeLanguage(lang.value)"
              >
                <q-item-section avatar>
                  <q-img
                    width="24px"
                    :src="lang.icon"
                  />
                </q-item-section>
                <q-item-section>{{ lang.label }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>

        <!-- User Menu -->
        <q-btn
          id="header-user-menu"
          flat
          round
        >
          <q-avatar
            size="32px"
            color="primary"
            text-color="white"
            class="text-caption text-weight-bold"
          >
            {{ userInitials }}
          </q-avatar>
          <q-menu
            anchor="bottom right"
            self="top right"
          >
            <q-list style="min-width: 200px">
              <q-item
                v-close-popup
                clickable
                to="/my_profile"
              >
                <q-item-section avatar>
                  <q-icon
                    name="person"
                    color="primary"
                  />
                </q-item-section>
                <q-item-section>{{ t.userMenu.profile.value }}</q-item-section>
              </q-item>
              <q-item
                v-close-popup
                clickable
                to="/admin/settings"
              >
                <q-item-section avatar>
                  <q-icon
                    name="settings"
                    color="primary"
                  />
                </q-item-section>
                <q-item-section>{{ t.userMenu.settings.value }}</q-item-section>
              </q-item>
              <q-separator/>
              <q-item
                v-close-popup
                clickable
                href="https://mapexos.io/"
                target="_blank"
              >
                <q-item-section avatar>
                  <q-icon
                    name="menu_book"
                    color="primary"
                  />
                </q-item-section>
                <q-item-section>{{ t.userMenu.docs.value }}</q-item-section>
              </q-item>
              <q-item
                v-close-popup
                clickable
                @click="emit('start-tour')"
              >
                <q-item-section avatar>
                  <q-icon
                    name="mdi-map-marker-path"
                    color="primary"
                  />
                </q-item-section>
                <q-item-section>{{ t.userMenu.startTour.value }}</q-item-section>
              </q-item>
              <q-separator/>
              <q-item
                v-close-popup
                clickable
                @click="handleLogout"
              >
                <q-item-section avatar>
                  <q-icon
                    name="logout"
                    color="negative"
                  />
                </q-item-section>
                <q-item-section>{{ t.userMenu.logout.value }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </q-toolbar>
  </q-header>
</template>

<style lang="scss" scoped>
.header-container {
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);
}

.org-indicator {
  .org-chip {
    max-width: 280px;
    transition: var(--mapex-transition-base);

    &:hover {
      transform: translateY(-1px);
      box-shadow: var(--mapex-shadow-sm);
    }

    .org-name {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      max-width: 220px;
      display: inline-block;
    }
  }
}
</style>

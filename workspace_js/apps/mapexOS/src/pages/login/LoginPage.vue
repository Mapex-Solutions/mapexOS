<script setup lang="ts">
defineOptions({
  name: 'LoginPage'
});

/** TYPE IMPORTS */
import type { QForm } from 'quasar';
import type { SupportedLocale, LanguageOption } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

/** COMPOSABLES */
import { useLoginTranslations } from '@composables/i18n';
import { useTheme } from '@composables/theme';

/** STORES */
import { useAuthStore } from '@stores/auth';

/** LOCAL IMPORTS */
import {
  DEFAULT_LANGUAGE,
  LANGUAGE_FLAGS,
  DEFAULT_REDIRECT_PATH,
  LOGIN_DELAY_MS,
  EMAIL_REGEX,
  LOCALE_STORAGE_KEY,
} from './constants';

/** COMPOSABLES & STORES */
const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const { locale } = useI18n();
const { isDark, toggleTheme } = useTheme();
const { welcome, features, form, validation, errors, languages, aria } = useLoginTranslations();

/** STATE */
const email = ref<string>('');
const password = ref<string>('');
const rememberMe = ref<boolean>(false);
const showPassword = ref<boolean>(false);
const selectedLanguage = ref<SupportedLocale>(DEFAULT_LANGUAGE);
const loading = ref<boolean>(false);
const errorMessage = ref<string>('');
const formRef = ref<QForm | null>(null);

/** COMPUTED */

/**
 * Language options with reactive labels
 * Labels update automatically when locale changes
 */
const languageOptions = computed<LanguageOption[]>(() => [
  { label: languages.english.value, value: 'en-US' as const, flag: LANGUAGE_FLAGS['en-US'] },
  { label: languages.portuguese.value, value: 'pt-BR' as const, flag: LANGUAGE_FLAGS['pt-BR'] }
]);

/** FUNCTIONS */

/**
 * Parse authentication error from API response
 * @param {unknown} err - Error object from auth API
 * @returns {string} Formatted error message
 */
function parseAuthError(err: unknown): string {
  const anyErr = err as any;
  const apiMessage: string | undefined = anyErr?.response?.data?.errors.join(',')
  return apiMessage || errors.default.value;
}

/**
 * Change language and persist preference
 * @param {SupportedLocale} newLocale - New locale to apply
 * @returns {void}
 */
function changeLanguage(newLocale: SupportedLocale): void {
  selectedLanguage.value = newLocale;
  locale.value = newLocale;
  localStorage.setItem(LOCALE_STORAGE_KEY, newLocale);
}

/**
 * Handle login form submission
 * @returns {Promise<void>}
 */
async function handleLogin(): Promise<void> {
  errorMessage.value = '';
  try {
    loading.value = true;

    // Validate inputs before calling the API
    const ok = await formRef.value?.validate();
    if (!ok) return;

    // Optional micro-delay for UX polish (spinner visibility)
    await new Promise((r) => setTimeout(r, LOGIN_DELAY_MS));
    await authStore.login(email.value, password.value, rememberMe.value);

    // Check if user must change password before entering the app
    if (authStore.user?.changePasswordNextLogin) {
      const redirectPath = (route.query.redirect as string) || DEFAULT_REDIRECT_PATH;
      await router.push({ path: '/change-password', query: { redirect: redirectPath } });
      return;
    }

    // Redirect to the original requested page or default to /home
    const redirectPath = (route.query.redirect as string) || DEFAULT_REDIRECT_PATH;
    await router.push({ path: redirectPath });
  } catch (err) {
    errorMessage.value = parseAuthError(err);
  } finally {
    loading.value = false;
  }
}

/** LIFECYCLE HOOKS */

/**
 * Load saved language preference on mount
 */
onMounted(() => {
  const savedLocale = localStorage.getItem(LOCALE_STORAGE_KEY) as SupportedLocale | null;
  if (savedLocale && (savedLocale === 'en-US' || savedLocale === 'pt-BR')) {
    selectedLanguage.value = savedLocale;
  }
});
</script>

<template>
  <q-page class="login-page">
    <div class="row full-width full-height">
      <!-- Left side with illustration -->
      <div class="col-12 col-md-8 illustration-side">
        <div class="illustration-container">
          <img alt="IoT Background" src="/iot-illustration.png" class="illustration-image"/>
          <div class="overlay">
            <div class="overlay-content q-pa-xl">
              <h2 class="text-h3 text-weight-bold text-white q-mb-sm">{{ welcome.title.value }}</h2>
              <p class="text-h5 text-weight-medium text-white q-mb-lg" style="opacity: 0.9;">{{ welcome.tagline.value }}</p>
              <p class="text-h6 text-white q-mb-xl">{{ welcome.subtitle.value }}</p>
              <div class="features q-gutter-y-md">
                <div class="feature-item">
                  <q-icon name="sensors" size="sm" color="white" class="q-mr-sm" />
                  <span>{{ features.monitoring.value }}</span>
                </div>
                <div class="feature-item">
                  <q-icon name="analytics" size="sm" color="white" class="q-mr-sm" />
                  <span>{{ features.analytics.value }}</span>
                </div>
                <div class="feature-item">
                  <q-icon name="security" size="sm" color="white" class="q-mr-sm" />
                  <span>{{ features.security.value }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side with login form -->
      <div class="col-12 col-md-4 login-side">
        <div class="login-container">
          <!-- Theme & Language selectors -->
          <div class="absolute-top-right q-pa-md row items-center q-gutter-x-xs">
            <q-btn
              flat
              round
              :icon="isDark ? 'light_mode' : 'dark_mode'"
              :color="isDark ? 'amber' : 'grey-7'"
              @click="toggleTheme"
            />
            <q-btn flat round>
              <q-icon name="language" class="login-icon-muted" />
              <q-menu anchor="bottom right" self="top right">
                <q-list style="min-width: 150px">
                  <q-item
                    v-for="lang in languageOptions"
                    v-close-popup
                    clickable
                    :key="lang.value"
                    @click="changeLanguage(lang.value)"
                  >
                    <q-item-section avatar>
                      <q-img width="24px" :src="lang.flag"/>
                    </q-item-section>
                    <q-item-section>{{ lang.label }}</q-item-section>
                    <q-item-section v-if="selectedLanguage === lang.value" side>
                      <q-icon name="check" color="primary"/>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>

          <!-- Login form -->
          <div class="login-form q-pa-lg">
            <div class="text-center q-mb-xl">
              <img src="/mapex-logo.png" alt="Mapex Logo" class="login-logo q-mb-md"/>
              <p class="text-subtitle1 login-subtitle">{{ form.title.value }}</p>
            </div>

            <!-- Enhanced Error Banner -->
            <transition name="slide-down">
              <div v-if="errorMessage" class="error-banner-wrapper q-mb-lg">
                <q-banner
                  rounded
                  class="error-banner"
                  role="alert"
                  aria-live="polite"
                >
                  <div class="error-content">
                    <div class="error-icon-wrapper">
                      <q-icon name="warning" size="24px" />
                    </div>
                    <div class="error-text">
                      <div class="error-title">{{ errors.title.value }}</div>
                      <div class="error-message">{{ errorMessage }}</div>
                    </div>
                    <q-btn
                      flat
                      round
                      size="sm"
                      class="error-close-btn"
                      icon="close"
                      :aria-label="aria.closeError.value"
                      @click="errorMessage = ''"
                    />
                  </div>
                </q-banner>
              </div>
            </transition>

            <q-form ref="formRef" class="q-gutter-y-md" @submit.prevent="handleLogin">
              <q-input
                v-model="email"
                outlined
                dense
                type="email"
                class="login-input"
                data-testid="login-email-input"
                :label="form.email.label.value"
                :disable="loading"
                :rules="[
                  (val) => !!val || validation.emailRequired.value,
                  (val) => EMAIL_REGEX.test(val) || validation.emailInvalid.value
                ]"
              >
                <template #prepend>
                  <q-icon name="email" color="primary"/>
                </template>
              </q-input>

              <q-input
                v-model="password"
                outlined
                dense
                class="login-input"
                data-testid="login-password-input"
                :type="showPassword ? 'text' : 'password'"
                :label="form.password.label.value"
                :disable="loading"
                :rules="[ (val) => !!val || validation.passwordRequired.value ]"
                @keyup.enter="handleLogin"
              >
                <template #prepend>
                  <q-icon name="lock" color="primary"/>
                </template>
                <template #append>
                  <q-icon
                    :name="showPassword ? 'visibility_off' : 'visibility'"
                    class="cursor-pointer"
                    @click="showPassword = !showPassword"
                  />
                </template>
              </q-input>

              <div class="row items-center justify-between q-mt-sm">
                <q-checkbox v-model="rememberMe" dense color="primary" :label="form.rememberMe.value" :disable="loading"/>
                <q-btn flat class="text-caption" color="primary" :label="form.forgotPassword.value" :disable="loading"/>
              </div>

              <q-btn
                unelevated
                type="submit"
                size="lg"
                class="full-width q-py-sm q-mt-lg"
                color="primary"
                data-testid="login-submit-btn"
                :loading="loading"
                :disable="loading"
              >
                {{ form.signIn.value }}
              </q-btn>
            </q-form>
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: var(--mapex-page-bg);
}

.illustration-side {
  height: 100vh;
  position: relative;
}

.illustration-container {
  height: 100%;
  position: relative;
  overflow: hidden;
  border-radius: 0 2rem 2rem 0;
}

.illustration-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.6) 100%);
  display: flex;
  align-items: center;
}

.feature-item {
  display: flex;
  align-items: center;
  color: white;
  font-size: 1.1rem;
}

.login-side {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: var(--mapex-page-bg);
}

.login-container {
  width: 100%;
  max-width: 440px;
  margin: 0 auto;
}

.login-logo {
  max-width: 180px;
  height: auto;
}

.login-subtitle {
  color: var(--mapex-text-secondary);
}

.login-icon-muted {
  color: var(--mapex-text-secondary);
}

.login-form {
  background: var(--mapex-surface-elevated);
  border-radius: var(--mapex-radius-xl);
  border: 1px solid var(--mapex-card-border);
  box-shadow: 0 4px 6px -1px var(--mapex-elevation-shadow), 0 2px 4px -1px var(--mapex-elevation-shadow);
}

.login-input {
  :deep(.q-field__control) {
    height: 56px;
    border-radius: var(--mapex-radius-md);
  }
  :deep(.q-field__marginal) {
    height: 56px;
  }
}

/* Enhanced Error Banner Styles */
.error-banner-wrapper {
  margin-bottom: 1.5rem;
}

.error-banner {
  background: var(--mapex-surface-elevated);
  border: 1px solid var(--q-negative);
  border-radius: var(--mapex-radius-lg);
  padding: 0;
  box-shadow: var(--mapex-shadow-md);

  :deep(.q-banner__content) {
    padding: 0;
  }
}

.error-content {
  display: flex;
  align-items: flex-start;
  padding: 16px 20px;
  gap: 12px;
  width: 100%;
}

.error-icon-wrapper {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  background: rgba(var(--q-negative-rgb), 0.1);
  border-radius: var(--mapex-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--q-negative);
  margin-top: 2px;
}

.error-text {
  flex: 1;
  min-width: 0;
}

.error-title {
  font-weight: 600;
  color: var(--q-negative);
  font-size: 14px;
  margin-bottom: 4px;
  line-height: 1.2;
}

.error-message {
  color: var(--q-negative);
  font-size: 13px;
  line-height: 1.4;
  word-wrap: break-word;
}

.error-close-btn {
  flex-shrink: 0;
  color: var(--q-negative);
  margin: -4px -4px 0 0;
  transition: var(--mapex-transition-base);

  &:hover {
    background: rgba(var(--q-negative-rgb), 0.1);
    transform: scale(1.1);
  }
}

/* Enhanced slide-down animation */
.slide-down-enter-active {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.slide-down-leave-active {
  transition: all 0.25s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}

@media (max-width: 1023px) {
  .illustration-side { display: none; }
  .login-container { padding: 2rem; }

  .error-content {
    padding: 14px 16px;
    gap: 10px;
  }

  .error-icon-wrapper {
    width: 36px;
    height: 36px;
  }

  .error-title {
    font-size: 13px;
  }

  .error-message {
    font-size: 12px;
  }
}
</style>
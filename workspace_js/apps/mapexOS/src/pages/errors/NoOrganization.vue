<template>
  <q-page class="no-org-page flex flex-center">
    <div class="q-pa-md" style="max-width: 600px; width: 100%;">
      <q-card flat bordered class="no-org-card shadow-2">
        <q-card-section class="text-center q-pt-xl q-pb-md">
          <q-icon
            size="100px"
            name="business_center"
            class="no-org-icon"
          />
        </q-card-section>

        <q-card-section class="text-center">
          <div class="text-h4 text-weight-medium q-mb-md no-org-title">
            {{ t.page.title.value }}
          </div>
          <div class="text-h6 q-mb-lg no-org-description">
            {{ t.page.description.value }}
          </div>
          <div class="text-body1 no-org-message">
            {{ t.page.message.value }}
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions class="q-pa-md justify-center">
          <q-btn
            unelevated
            size="md"
            class="q-px-lg"
            icon="logout"
            :label="t.actions.logout.value"
            color="primary"
            @click="handleLogout"
          />
        </q-card-actions>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
defineOptions({
  name: 'NoOrganization'
});

/** VUE IMPORTS */
import { useRouter } from 'vue-router';

/** COMPOSABLES */
import { useNoOrganizationTranslations } from '@composables/i18n';

/** STORES */
import { useAuthStore } from '@stores/auth';

/** COMPOSABLES & STORES */
const router = useRouter();
const authStore = useAuthStore();
const t = useNoOrganizationTranslations();

/** FUNCTIONS */

/**
 * Handle logout action
 * Clears auth state and redirects to login page
 */
function handleLogout() {
  authStore.logout();
  void router.push('/');
}
</script>

<style lang="scss" scoped>
.no-org-page {
  background-color: var(--mapex-page-bg);
}

.no-org-card {
  background-color: var(--mapex-surface-elevated);
  border-color: var(--mapex-card-border);
}

.no-org-icon {
  color: var(--mapex-text-muted);
}

.no-org-title {
  color: var(--mapex-text-primary);
}

.no-org-description {
  color: var(--mapex-text-secondary);
}

.no-org-message {
  color: var(--mapex-text-secondary);
}
</style>

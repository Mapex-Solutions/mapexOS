<template>
  <div class="error-page">
    <div class="error-container">
      
      <!-- Mapex Logo -->
      <div class="logo-container q-mb-xl">
        <img
          :src="logoSrc"
          :alt="aria.logoAlt.value"
          width="180px"
        />
      </div>

      <div class="text-h2 text-primary q-mb-md">{{ notFound.code.value }}</div>
      <div class="text-h5 text-primary q-mb-xl">{{ notFound.title.value }}</div>
      <div class="text-body1 q-mb-xl">
        {{ notFound.description.value }}
      </div>

      <q-btn
        size="md"
        class="q-px-xl"
        :label="notFound.goBack.value"
        color="primary"
        @click="goBack"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ErrorNotFound'
});

/** VUE IMPORTS */
import { ref } from 'vue';
import { useRouter } from 'vue-router';

/** COMPOSABLES */
import { useErrorTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import logoImage from '/mapex-logo.png';

/** COMPOSABLES & STORES */
const router = useRouter();
const { notFound, aria } = useErrorTranslations();

/** STATE */
const logoSrc = ref(logoImage);

/** FUNCTIONS */

/**
 * Navigate back to previous page
 */
function goBack(): void {
  router.go(-1);
}
</script>

<style lang="scss" scoped>
.error-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--mapex-page-bg);
}

.error-container {
  max-width: 400px;
  text-align: center;
  padding: 2rem;
  background-color: var(--mapex-surface-elevated);
  border-radius: var(--mapex-radius-md);
  border: 1px solid var(--mapex-card-border);
  box-shadow: 0 1px 5px var(--mapex-elevation-shadow);
}

.logo-container {
  margin: 0 auto;
}

.q-btn {
  font-weight: 500;
  border-radius: var(--mapex-radius-xs);
}
</style>

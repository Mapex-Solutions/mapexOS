<script setup lang="ts">
/** TYPE IMPORTS */
import type { InfoBannerProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { INFO_BANNER_VARIANT_ICON } from './constants';

defineOptions({
  name: 'InfoBanner'
});

/** PROPS & EMITS */
const props = withDefaults(defineProps<InfoBannerProps>(), {
  variant: 'info',
  dense: false
});

/** COMPUTED */
const resolvedIcon = computed(() => props.icon ?? INFO_BANNER_VARIANT_ICON[props.variant]);
</script>

<template>
  <q-banner
    rounded
    :dense="dense"
    :class="['info-banner', `info-banner--${variant}`]"
  >
    <template #avatar>
      <q-icon :name="resolvedIcon" />
    </template>
    <div v-if="title" class="info-banner__title">{{ title }}</div>
    <div class="info-banner__body">
      <slot />
    </div>
  </q-banner>
</template>

<style scoped lang="scss">
/* Compound `.info-banner.q-banner` lifts specificity above Quasar's own
   `.q-banner` so the soft fill and intent color win. The whole banner is
   tinted with the variant foreground — icon and text inherit currentColor,
   matching the legacy colored banners while staying theme-aware. */
.info-banner.q-banner {
  background: var(--mapex-neutral-soft);
  color: var(--mapex-neutral-fg);
  border: 1px solid var(--mapex-neutral-border);
  border-radius: var(--mapex-radius-md);
}

.info-banner--info.q-banner {
  background: var(--mapex-info-soft);
  border-color: var(--mapex-info-border);
  color: var(--mapex-info-fg);
}

.info-banner--success.q-banner {
  background: var(--mapex-success-soft);
  border-color: var(--mapex-success-border);
  color: var(--mapex-success-fg);
}

.info-banner--warning.q-banner {
  background: var(--mapex-warning-soft);
  border-color: var(--mapex-warning-border);
  color: var(--mapex-warning-fg);
}

.info-banner--danger.q-banner {
  background: var(--mapex-danger-soft);
  border-color: var(--mapex-danger-border);
  color: var(--mapex-danger-fg);
}

.info-banner--neutral.q-banner {
  background: var(--mapex-neutral-soft);
  border-color: var(--mapex-neutral-border);
  color: var(--mapex-neutral-fg);
}

.info-banner__title {
  font-weight: var(--mapex-font-weight-semibold);
  margin-bottom: var(--mapex-spacing-2xs);
}

.info-banner__body {
  font-size: var(--mapex-font-sm);
}
</style>

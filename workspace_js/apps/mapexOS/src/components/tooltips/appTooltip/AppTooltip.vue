<template>
  <q-tooltip
    v-if="shouldShow"
    class="mapex-tooltip"
    :delay="delay"
    :hide-delay="hideDelay"
    :anchor="anchor"
    :self="self"
    :offset="offset"
    :max-width="maxWidth"
  >
    <slot>{{ content }}</slot>
  </q-tooltip>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { AppTooltipProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useQuasar } from 'quasar';

/** LOCAL IMPORTS */
import { TOOLTIP_DEFAULTS } from './constants';

/** PROPS */
const props = withDefaults(defineProps<AppTooltipProps>(), {
  showOnMobile: false,
  delay: TOOLTIP_DEFAULTS.DELAY,
  hideDelay: TOOLTIP_DEFAULTS.HIDE_DELAY,
  maxWidth: TOOLTIP_DEFAULTS.MAX_WIDTH,
  anchor: TOOLTIP_DEFAULTS.ANCHOR,
  self: TOOLTIP_DEFAULTS.SELF,
  disabled: false,
});

/** COMPOSABLES & STORES */
const $q = useQuasar();

/** COMPUTED */

/**
 * Check if device is touch-enabled
 * @returns True if device has touch capabilities
 */
const isTouchDevice = computed(
  () => $q.platform.is.mobile || $q.platform.has.touch,
);

/**
 * Determine if tooltip should render
 * @returns True if tooltip should be shown
 */
const shouldShow = computed(() => {
  if (props.disabled) return false;
  if (isTouchDevice.value && !props.showOnMobile) return false;
  return true;
});
</script>

<style lang="scss">
.mapex-tooltip {
  background: var(--mapex-tooltip-bg) !important;
  color: var(--mapex-tooltip-text) !important;
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-sm);
  font-size: 0.8125rem;
  line-height: 1.4;
  padding: 6px 10px;

  .text-caption,
  .text-grey-5,
  .text-grey-6 {
    color: var(--mapex-tooltip-text-secondary) !important;
  }
}
</style>

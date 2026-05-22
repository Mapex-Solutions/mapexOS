<template>
  <q-chip
    :class="chipClasses"
    :style="chipStyles"
    :clickable="clickable"
    :outline="outline"
    :square="square"
    :dense="dense"
    :color="chipColor"
    :text-color="chipTextColor"
  >
    <q-icon
      v-if="icon"
      :name="icon"
      :size="sizeConfig.iconSize"
      class="chip-icon"
    />
    <span class="chip-label">{{ displayLabel }}</span>
  </q-chip>
</template>

<script setup lang="ts">
defineOptions({
  name: 'DetailChip'
});

/** TYPE IMPORTS */
import type { DetailChipProps } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** LOCAL IMPORTS */
import {
  DETAIL_CHIP_COLORS,
  DETAIL_CHIP_TEXT_COLORS,
  DETAIL_CHIP_SIZES,
  DETAIL_CHIP_DEFAULTS,
} from './constants';

/** PROPS */
const props = withDefaults(defineProps<DetailChipProps>(), {
  color: DETAIL_CHIP_DEFAULTS.color,
  size: DETAIL_CHIP_DEFAULTS.size,
  clickable: DETAIL_CHIP_DEFAULTS.clickable,
  outline: DETAIL_CHIP_DEFAULTS.outline,
  square: DETAIL_CHIP_DEFAULTS.square,
  dense: DETAIL_CHIP_DEFAULTS.dense,
  rounded: DETAIL_CHIP_DEFAULTS.rounded,
});

/** COMPUTED */

/**
 * Display label (uses label or value prop)
 * @returns {string} Label text to display
 */
const displayLabel = computed(() => {
  if (props.label !== undefined) return props.label;
  if (props.value !== undefined) return String(props.value);
  return '';
});

/**
 * Get size configuration based on size prop
 * @returns {object} Size configuration with fontSize, padding, and iconSize
 */
const sizeConfig = computed(() => {
  return DETAIL_CHIP_SIZES[props.size];
});

/**
 * Get chip color from color mapping
 * @returns {string} Quasar color value
 */
const chipColor = computed(() => {
  return DETAIL_CHIP_COLORS[props.color];
});

/**
 * Get chip text color (custom or default based on background)
 * @returns {string} Text color value
 */
const chipTextColor = computed(() => {
  if (props.textColor) return props.textColor;
  return DETAIL_CHIP_TEXT_COLORS[props.color];
});

/**
 * Computed CSS classes for the chip
 * @returns {string} CSS class string
 */
const chipClasses = computed(() => {
  const classes = ['detail-chip', `detail-chip--${props.size}`];
  if (props.dense) classes.push('detail-chip--dense');
  if (props.rounded) classes.push('detail-chip--rounded');
  return classes.join(' ');
});

/**
 * Computed inline styles for the chip
 * @returns {object} CSS style object
 */
const chipStyles = computed(() => {
  // When dense, use Quasar's default sizing
  if (props.dense) {
    return {};
  }
  return {
    fontSize: sizeConfig.value.fontSize,
    padding: sizeConfig.value.padding,
  };
});
</script>

<style lang="scss" scoped>
.detail-chip {
  font-weight: 500;
  letter-spacing: 0.3px;
  transition: var(--mapex-transition-base);

  // Icon spacing
  .chip-icon {
    margin-right: 6px;
  }

  // Label styling
  .chip-label {
    line-height: 1.2;
    white-space: nowrap;
  }

  // Size variants with better spacing
  &--xs {
    .chip-icon {
      margin-right: 4px;
    }
  }

  &--sm {
    .chip-icon {
      margin-right: 5px;
    }
  }

  &--md {
    .chip-icon {
      margin-right: 6px;
    }
  }

  &--lg {
    .chip-icon {
      margin-right: 8px;
    }
  }

  // Dense mode - use Quasar's default sizing
  &--dense {
    .chip-icon {
      margin-right: 4px;
    }
  }

  // Rounded corners modifier
  &--rounded {
    border-radius: var(--mapex-radius-xl);
  }

  // Hover effect for clickable chips
  &:deep(.q-chip--clickable) {
    &:hover {
      transform: translateY(-1px);
      box-shadow: var(--mapex-shadow-sm);
    }

    &:active {
      transform: translateY(0);
    }
  }
}
</style>

<template>
  <q-chip
    :class="chipClasses"
    :style="chipStyles"
    :clickable="clickable"
    :removable="removable"
    :disable="disable"
    :dense="dense"
    :outline="outline"
    :square="square"
    :color="chipColor"
    :text-color="chipTextColor"
    @remove="emit('remove')"
    @click="handleClick"
  >
    <q-icon
      v-if="icon"
      :name="icon"
      :size="sizeConfig.iconSize"
      class="chip-icon"
    />
    <span class="chip-label">{{ label }}</span>
  </q-chip>
</template>

<script setup lang="ts">
defineOptions({
  name: 'SelectableChip',
});

/** TYPE IMPORTS */
import type { SelectableChipProps, SelectableChipEmits } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** LOCAL IMPORTS */
import {
  DETAIL_CHIP_COLORS,
  DETAIL_CHIP_TEXT_COLORS,
  DETAIL_CHIP_SIZES,
} from '../DetailChip';
import { SELECTABLE_CHIP_DEFAULTS } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<SelectableChipProps>(), {
  color: SELECTABLE_CHIP_DEFAULTS.color,
  size: SELECTABLE_CHIP_DEFAULTS.size,
  removable: SELECTABLE_CHIP_DEFAULTS.removable,
  dense: SELECTABLE_CHIP_DEFAULTS.dense,
  disable: SELECTABLE_CHIP_DEFAULTS.disable,
  clickable: SELECTABLE_CHIP_DEFAULTS.clickable,
  outline: SELECTABLE_CHIP_DEFAULTS.outline,
  square: SELECTABLE_CHIP_DEFAULTS.square,
});

const emit = defineEmits<SelectableChipEmits>();

/** COMPUTED */

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
  const classes = ['selectable-chip', `selectable-chip--${props.size}`];
  if (props.dense) classes.push('selectable-chip--dense');
  return classes.join(' ');
});

/**
 * Computed inline styles for the chip
 * @returns {object} CSS style object
 */
const chipStyles = computed(() => {
  if (props.dense) {
    return {};
  }
  return {
    fontSize: sizeConfig.value.fontSize,
    padding: sizeConfig.value.padding,
  };
});

/** FUNCTIONS */

/**
 * Handle chip click event
 */
function handleClick(): void {
  if (props.clickable && !props.disable) {
    emit('click');
  }
}
</script>

<style lang="scss" scoped>
.selectable-chip {
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

  // Hover effect for clickable chips
  &:deep(.q-chip--clickable:not(.q-chip--disabled)) {
    &:hover {
      transform: translateY(-1px);
      box-shadow: var(--mapex-shadow-sm);
    }

    &:active {
      transform: translateY(0);
    }
  }

  // Remove button styling
  &:deep(.q-chip__icon--remove) {
    opacity: 0.8;
    transition: opacity 0.2s ease;

    &:hover {
      opacity: 1;
    }
  }
}
</style>

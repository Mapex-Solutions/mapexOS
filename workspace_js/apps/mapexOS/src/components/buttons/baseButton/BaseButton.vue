<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { QBtn } from 'quasar';

defineOptions({
  name: 'BaseButton'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** PROPS & EMITS */
interface Props {
  /** All QBtn props - will be bound directly to q-btn */
  [key: string]: any;
}

const props = defineProps<Props>();

/** COMPUTED */
/**
 * Merges user-provided class with MapexOS default styles
 * Always applies rounded-borders for consistent button appearance
 */
const buttonClass = computed(() => {
  const userClass = props.class || '';
  return `rounded-borders ${userClass}`.trim();
});

/**
 * Creates a new props object excluding 'class' since we handle it separately
 * This allows proper v-bind spreading while customizing the class
 */
const qBtnProps = computed(() => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { class: _unusedClass, ...rest } = props;
  return rest;
});
</script>

<template>
  <q-btn
    v-bind="qBtnProps"
    :class="buttonClass"
  >
    <!-- Forward all slots to q-btn -->
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData || {}" />
    </template>
  </q-btn>
</template>

<style scoped lang="scss">
/**
 * MapexOS Button Styles
 * Ensures consistent rounded borders across all buttons
 */
:deep(.rounded-borders) {
  border-radius: var(--mapex-radius-md) !important;
}
</style>

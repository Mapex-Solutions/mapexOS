<script setup lang="ts">
defineOptions({
  name: 'GroupFrameNode',
});

/** TYPE IMPORTS */
import type { WorkflowNodeComponentProps } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';
import { NodeResizer } from '@vue-flow/node-resizer';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS */
import { FRAME_COLOR_OPTIONS } from './constants';

/** PROPS & EMITS */
const props = defineProps<WorkflowNodeComponentProps>();

/** COMPOSABLES & STORES */
const { updateNodeConfig } = useWorkflowContext();
const { t } = usePluginI18n('core-annotations');

/** COMPUTED */

/**
 * Frame title from config
 */
const title = computed<string>(
  () => (props.data.config?.title as string) || '',
);

/**
 * Frame description from config
 */
const description = computed<string>(
  () => (props.data.config?.description as string) || '',
);

/**
 * Frame color name from config
 */
const colorName = computed<string>(
  () => (props.data.config?.color as string) || 'blue-grey',
);

/**
 * Resolved hex color from FRAME_COLOR_OPTIONS or custom hex
 */
const colorHex = computed<string>(() => {
  const opt = FRAME_COLOR_OPTIONS.find(o => o.value === colorName.value);
  if (opt) return opt.hex;
  return colorName.value.startsWith('#') ? colorName.value : '#78909c';
});

/**
 * Frame width from config (px)
 */
const width = computed<number>(
  () => (props.data.config?.width as number) || 300,
);

/**
 * Frame height from config (px)
 */
const height = computed<number>(
  () => (props.data.config?.height as number) || 200,
);

/**
 * Frame inline style — dimensions and dynamic border/bg color
 */
const frameStyle = computed(() => ({
  width: `${width.value}px`,
  height: `${height.value}px`,
  '--frame-color': colorHex.value,
}));

/** FUNCTIONS */

/**
 * Handle resize event from NodeResizer — persist new dimensions to config
 *
 * @param {object} event - Resize event with dimensions
 */
function handleResize(event: { params: { width: number; height: number } }): void {
  updateNodeConfig(props.id, {
    width: Math.round(event.params.width),
    height: Math.round(event.params.height),
  });
}
</script>

<template>
  <div
    class="group-frame"
    :class="{ 'group-frame--selected': selected }"
    :style="frameStyle"
  >
    <NodeResizer
      :min-width="150"
      :min-height="100"
      :color="colorHex"
      :handle-style="{ width: '8px', height: '8px' }"
      @resize="handleResize"
    />

    <!-- Header -->
    <div class="group-frame__header">
      <q-icon name="dashboard" size="14px" :style="{ color: colorHex }" />
      <span class="group-frame__title">{{ title || t('nodes.group_frame.defaultTitle') }}</span>
    </div>

    <!-- Description -->
    <div v-if="description" class="group-frame__description">
      {{ description }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.group-frame {
  --frame-color: #78909c;

  position: relative;
  border: 1.5px dashed var(--frame-color);
  border-radius: var(--mapex-radius-lg);
  background: color-mix(in srgb, var(--frame-color) 6%, transparent);
  padding: 8px 10px;
  cursor: grab;
  transition: border-color var(--mapex-transition-base),
              background var(--mapex-transition-base);

  &--selected {
    border-color: var(--frame-color);
    background: color-mix(in srgb, var(--frame-color) 10%, transparent);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--frame-color) 25%, transparent);
  }

  &__header {
    display: flex;
    align-items: center;
    gap: 4px;
    pointer-events: none;
  }

  &__title {
    font-size: 11px;
    font-weight: 700;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: calc(100% - 24px);
    letter-spacing: 0.3px;
  }

  &__description {
    margin-top: 2px;
    font-size: 9px;
    color: var(--mapex-text-muted);
    line-height: 1.3;
    white-space: pre-wrap;
    word-break: break-word;
    pointer-events: none;
    max-height: 40px;
    overflow: hidden;
  }
}
</style>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { BaseWorkflowNodeProps } from '@src/components/workflow/interfaces/nodeConfig.interface';
import type { HandleOverrides } from '@src/components/workflow/interfaces/workflowPlugin.interface';

/** VUE IMPORTS */
import { computed } from 'vue';
import { Handle, Position } from '@vue-flow/core';
import { useI18n } from 'vue-i18n';

/** COMPOSABLES */
import { useWorkflowContext } from '@src/composables/workflow';

/** LOCAL IMPORTS */
import { POSITION_OPTIONS } from '@src/components/workflow/constants/workflowSdk.constant';

/**
 * Local extension — forces @vue/compiler-sfc to resolve hasErrors
 * even when external interface cache is stale.
 */
interface Props extends BaseWorkflowNodeProps {
  /** Whether the node has validation errors */
  hasErrors?: boolean;
}

/** PROPS & EMITS */
const props = withDefaults(defineProps<Props>(), {
  label: '',
  hasErrors: false,
});

/** COMPOSABLES & STORES */
const { updateNodeConfig, nodes, addNoteToNode, pushSnapshot } = useWorkflowContext();
const { t } = useI18n();

/** COMPUTED */

/**
 * Root element inline style — applies dynamic --node-color when colorHex is provided
 */
const rootStyle = computed(() => {
  const style: Record<string, string> = { 'min-width': nodeMinWidth.value };
  if (props.colorHex) {
    style['--node-color'] = props.colorHex;
  }
  return style;
});

/**
 * Max handle count on any single edge (top or bottom).
 * Used to calculate min-width so handles have room to spread.
 */
const maxHandleCount = computed(() => {
  const inputCount = (props.inputs || []).length;
  const outputCount = (props.outputs || []).length;
  return Math.max(inputCount, outputCount);
});

/**
 * Dynamic min-width based on handle count.
 * Single handle: no min-width needed. Multiple: 50px per handle.
 */
const nodeMinWidth = computed(() => {
  if (maxHandleCount.value <= 1) return '44px';
  return `${maxHandleCount.value * 50}px`;
});

/** FUNCTIONS */

/**
 * Map handle position string to Vue Flow Position enum
 *
 * @param {string} pos - Position string
 * @returns {Position} Vue Flow Position
 */
function mapPosition(pos: string): Position {
  const map: Record<string, Position> = {
    top: Position.Top,
    bottom: Position.Bottom,
    left: Position.Left,
    right: Position.Right,
  };
  return map[pos] || Position.Bottom;
}

/**
 * Compute inline style for distributing multiple handles along an edge.
 * When there are N handles on the same edge, they are evenly distributed.
 *
 * @param {number} index - Handle index in the array
 * @param {number} total - Total handles on this edge
 * @param {string} position - Edge position (top/bottom/left/right)
 * @returns {Record<string, string>} CSS style object
 */
function getHandleStyle(index: number, total: number, position: string): Record<string, string> {
  if (total <= 1) return {};

  const pct = ((index + 1) / (total + 1)) * 100;

  if (position === 'top' || position === 'bottom') {
    return { left: `${pct}%` };
  }
  return { top: `${pct}%` };
}

/**
 * Compute inline style for a handle label positioned near its handle
 *
 * @param {number} index - Handle index
 * @param {number} total - Total handles on this edge
 * @param {string} position - Edge position
 * @returns {Record<string, string>} CSS style object
 */
function getLabelStyle(index: number, total: number, position: string): Record<string, string> {
  const pct = ((index + 1) / (total + 1)) * 100;

  if (position === 'top' || position === 'bottom') {
    return { left: `${pct}%`, transform: 'translateX(-50%)' };
  }
  return { top: `${pct}%`, transform: 'translateY(-50%)' };
}

/**
 * Update a handle's position via config.__handleOverrides
 *
 * @param {string} handleId - Handle ID to update
 * @param {string} position - New position (top/bottom/left/right)
 * @returns {void}
 */
function setHandlePosition(handleId: string, position: string): void {
  const node = nodes.value.find(n => n.id === props.id);
  if (!node) return;

  const current = (node.config.__handleOverrides as HandleOverrides) || {};
  const updated: HandleOverrides = {
    ...current,
    [handleId]: {
      ...current[handleId],
      position: position as 'top' | 'bottom' | 'left' | 'right',
    },
  };
  updateNodeConfig(props.id, { __handleOverrides: updated });
}

</script>

<template>
  <div
    class="wf-node"
    :class="{
      'wf-node--selected': selected,
      'wf-node--circle': shape === 'circle',
      'wf-node--error': hasErrors,
    }"
    :style="rootStyle"
  >
    <!-- Input handle labels (above handles) -->
    <template v-if="(inputs || []).length >= 1">
      <span
        v-for="(input, idx) in (inputs || [])"
        :key="`label-in-${input.id}`"
        class="wf-node__handle-label wf-node__handle-label--top"
        :style="{
          ...getLabelStyle(idx, (inputs || []).length, input.position),
          ...(input.color ? { color: input.color } : {}),
        }"
      >{{ input.label }}</span>
    </template>

    <!-- Input handles (with right-click position menu) -->
    <Handle
      v-for="(input, idx) in (inputs || [])"
      :key="input.id"
      type="target"
      :position="mapPosition(input.position)"
      :id="input.id"
      :style="{
        ...getHandleStyle(idx, (inputs || []).length, input.position),
        ...(input.color ? { '--handle-color': input.color } : {}),
      }"
      class="wf-node__handle wf-node__handle--target"
    >
      <q-menu context-menu class="wf-node__handle-menu">
        <q-list dense>
          <q-item
            v-for="opt in POSITION_OPTIONS"
            :key="opt.value"
            v-close-popup
            clickable
            :active="input.position === opt.value"
            @click="setHandlePosition(input.id, opt.value)"
          >
            <q-item-section side>
              <q-icon :name="opt.icon" size="14px" />
            </q-item-section>
            <q-item-section>{{ opt.label }}</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </Handle>


    <!-- Icon square + right-click context menu -->
    <div class="wf-node__icon">
      <q-icon :name="icon" size="20px" />

      <!-- Error badge (top-right of icon) -->
      <div v-if="hasErrors" class="wf-node__error-badge">
        <q-icon name="priority_high" size="10px" />
      </div>

      <!-- Node context menu -->
      <q-menu context-menu class="wf-node__context-menu">
        <q-list dense>
          <q-item
            v-close-popup
            clickable
            @click="pushSnapshot('Add note'); addNoteToNode(id)"
          >
            <q-item-section side>
              <q-icon name="sticky_note_2" size="xs" />
            </q-item-section>
            <q-item-section>{{ t('wf.common.addNote') }}</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </div>

    <!-- Label (always visible, tight below icon) -->
    <div class="wf-node__label">{{ label }}</div>

    <!-- Output handles (with right-click position menu) -->
    <Handle
      v-for="(output, idx) in (outputs || [])"
      :key="output.id"
      type="source"
      :position="mapPosition(output.position)"
      :id="output.id"
      :style="{
        ...getHandleStyle(idx, (outputs || []).length, output.position),
        ...(output.color ? { '--handle-color': output.color } : {}),
      }"
      class="wf-node__handle wf-node__handle--source"
    >
      <q-menu context-menu class="wf-node__handle-menu">
        <q-list dense>
          <q-item
            v-for="opt in POSITION_OPTIONS"
            :key="opt.value"
            v-close-popup
            clickable
            :active="output.position === opt.value"
            @click="setHandlePosition(output.id, opt.value)"
          >
            <q-item-section side>
              <q-icon :name="opt.icon" size="14px" />
            </q-item-section>
            <q-item-section>{{ opt.label }}</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </Handle>

    <!-- Output handle labels (below handles) -->
    <template v-if="(outputs || []).length >= 1">
      <span
        v-for="(output, idx) in (outputs || [])"
        :key="`label-out-${output.id}`"
        class="wf-node__handle-label wf-node__handle-label--bottom"
        :style="{
          ...getLabelStyle(idx, (outputs || []).length, output.position),
          ...(output.color ? { color: output.color } : {}),
        }"
      >{{ output.label }}</span>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.wf-node {
  --node-color: var(--q-primary);
  --node-border: var(--mapex-wf-node-border, rgba(255, 255, 255, 0.15));
  --node-border-hover: var(--mapex-wf-node-border-hover, rgba(255, 255, 255, 0.3));
  --node-border-selected: var(--mapex-wf-node-border-selected, rgba(255, 255, 255, 0.4));

  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  cursor: pointer;

  /* Icon square — the hero element */
  &__icon {
    position: relative;
    width: 44px;
    height: 44px;
    border-radius: var(--mapex-radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(
      135deg,
      var(--node-color),
      color-mix(in srgb, var(--node-color) 70%, black)
    );
    color: var(--mapex-wf-text-on-accent, #fff);
    border: 1.5px solid var(--node-border);
    box-shadow: var(--mapex-shadow-sm);
    transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1),
                box-shadow 0.25s ease,
                border-color 0.25s ease;
  }

  &:hover &__icon {
    transform: scale(1.08);
    border-color: var(--node-border-hover);
    box-shadow:
      0 4px 20px color-mix(in srgb, var(--node-color) 40%, transparent),
      0 0 0 1px color-mix(in srgb, var(--node-color) 40%, transparent);
  }

  &--selected &__icon {
    border-color: var(--node-border-selected);
    transform: scale(1.08);
    box-shadow:
      0 0 0 3px color-mix(in srgb, var(--node-color) 70%, transparent),
      0 0 12px color-mix(in srgb, var(--node-color) 50%, transparent),
      0 6px 24px color-mix(in srgb, var(--node-color) 35%, transparent);
  }

  &--selected &__label {
    color: var(--node-color);
    font-weight: var(--mapex-font-weight-bold);
  }

  /* Error state — subtle thin border + badge indicator */
  &--error &__icon {
    border-color: var(--q-negative);
    box-shadow: 0 0 8px color-mix(in srgb, var(--q-negative) 25%, transparent);
  }

  /* Error badge — warning icon on top-right corner of icon */
  &__error-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    width: 14px;
    height: 14px;
    background: var(--q-negative);
    color: #fff;
    border-radius: var(--mapex-radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3;
    box-shadow: var(--mapex-shadow-sm);
  }

  /* Circle shape variant (e.g., Start node) */
  &--circle &__icon {
    border-radius: 50%;
  }

  /* Label — tight below icon */
  &__label {
    margin-top: var(--mapex-spacing-xs);
    font-size: var(--mapex-font-2xs);
    font-weight: var(--mapex-font-weight-semibold);
    color: var(--mapex-text-primary);
    text-align: center;
    white-space: nowrap;
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: var(--mapex-line-height-tight);
  }

  /* Handle labels — hidden by default, visible on node hover */
  &__handle-label {
    position: absolute;
    font-size: 6px;
    font-weight: var(--mapex-font-weight-semibold);
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    white-space: nowrap;
    pointer-events: none;
    letter-spacing: 0.3px;
    z-index: 1;
    opacity: 0;
    transition: var(--mapex-transition-fast);

    &--top {
      top: -14px;
    }

    &--bottom {
      bottom: -14px;
    }
  }

  &:hover &__handle-label {
    opacity: 1;
  }

  /* Handles — small dots with theme-aware border, per-handle color via --handle-color */
  &__handle {
    width: 8px !important;
    height: 8px !important;
    background: var(--handle-color, var(--mapex-wf-edge-color, var(--q-primary))) !important;
    border: 1.5px solid var(--mapex-wf-handle-border, #fff) !important;
    border-radius: var(--mapex-radius-full) !important;
    box-shadow: 0 0 0 1px var(--mapex-wf-edge-glow, rgba(59, 109, 94, 0.4));
    transition: transform 0.15s, box-shadow 0.15s, width 0.15s, height 0.15s;
    z-index: 2;

    &:hover {
      width: 12px !important;
      height: 12px !important;
      background: var(--handle-color, var(--mapex-wf-edge-hover, var(--q-primary))) !important;
      box-shadow: 0 0 8px var(--mapex-wf-edge-glow, rgba(59, 109, 94, 0.6)),
                  0 0 0 2px var(--handle-color, var(--mapex-wf-edge-color, var(--q-primary)));
    }

    &--target {
      top: -2px !important;
    }

    &--source {
      bottom: -2px !important;
    }
  }

  /* Context menu for handle position */
  &__handle-menu {
    background: var(--mapex-surface-elevated) !important;
    border: 1px solid var(--mapex-card-border) !important;
    border-radius: var(--mapex-radius-md) !important;
    box-shadow: var(--mapex-shadow-md) !important;
    min-width: 110px !important;
  }


  /* Node context menu */
  &__context-menu {
    background: var(--mapex-surface-elevated) !important;
    border: 1px solid var(--mapex-card-border) !important;
    border-radius: var(--mapex-radius-md) !important;
    box-shadow: var(--mapex-shadow-md) !important;
    min-width: 130px !important;
  }
}

</style>

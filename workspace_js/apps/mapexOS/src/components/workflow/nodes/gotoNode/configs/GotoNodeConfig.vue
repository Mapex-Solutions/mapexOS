<script setup lang="ts">
defineOptions({
  name: 'GotoNodeConfig',
});

/** TYPE IMPORTS */
import type {
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
} from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { GOTO_COLOR_OPTIONS } from '../constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { nodes } = useWorkflowContext();
const { t } = usePluginI18n('core-flow-control');

/** STATE */

/**
 * Whether the custom hex input is shown
 */
const showHexInput = ref(false);

/**
 * Custom hex color buffer
 */
const hexBuffer = ref('');

/** COMPUTED */

/**
 * Current role from config
 */
const role = computed<'sender' | 'receiver'>(
  () => (props.config.role as 'sender' | 'receiver') || 'sender',
);

/**
 * Current pair label from config
 */
const pairLabel = computed<string>(
  () => (props.config.pairLabel as string) || '',
);

/**
 * Effective pair color — receivers inherit from matched sender dynamically
 */
const pairColor = computed<string>(() => {
  const own = (props.config.pairColor as string) || 'deep-purple-6';
  if (role.value === 'receiver' && pairLabel.value) {
    const sender = nodes.value.find(n =>
      n.type === 'core/goto' &&
      n.id !== props.config._nodeId &&
      (n.config?.role as string) === 'sender' &&
      (n.config?.pairLabel as string) === pairLabel.value,
    );
    if (sender) return (sender.config?.pairColor as string) || 'deep-purple-6';
  }
  return own;
});

/**
 * Whether the current color is a custom hex (not in presets)
 */
const isCustomColor = computed(() =>
  !GOTO_COLOR_OPTIONS.some(o => o.value === pairColor.value),
);

/**
 * Hex color resolved from GOTO_COLOR_OPTIONS or used directly if custom
 */
const colorHex = computed<string>(() => {
  const opt = GOTO_COLOR_OPTIONS.find(o => o.value === pairColor.value);
  if (opt) return opt.hex;
  // pairColor IS the hex value for custom colors
  return pairColor.value.startsWith('#') ? pairColor.value : '#5e35b1';
});

/**
 * Existing sender labels — unique pairLabels from goto sender nodes (excluding self)
 */
const senderLabels = computed(() => {
  const result: Array<{ label: string; color: string; hex: string }> = [];
  const seen = new Set<string>();

  for (const node of nodes.value) {
    if (node.type !== 'core/goto') continue;
    if (node.id === props.config._nodeId) continue;
    if ((node.config?.role as string) !== 'sender') continue;
    const lbl = (node.config?.pairLabel as string) || '';
    if (!lbl || seen.has(lbl)) continue;
    seen.add(lbl);

    const clr = (node.config?.pairColor as string) || 'deep-purple-6';
    const hex = GOTO_COLOR_OPTIONS.find(o => o.value === clr)?.hex ?? '#5e35b1';
    result.push({ label: lbl, color: clr, hex });
  }
  return result.sort((a, b) => a.label.localeCompare(b.label));
});

/**
 * Matched pair nodes — other goto nodes with the same pairLabel
 */
const matchedPairs = computed(() => {
  if (!pairLabel.value) return [];
  return nodes.value
    .filter(n =>
      n.type === 'core/goto' &&
      n.id !== props.config._nodeId &&
      (n.config?.pairLabel as string) === pairLabel.value,
    )
    .map(n => ({
      id: n.id,
      label: n.label || n.id,
      role: (n.config?.role as string) || 'sender',
      color: (n.config?.pairColor as string) || 'deep-purple-6',
    }));
});

/** FUNCTIONS */

/**
 * Emit config update with merged values
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}

/**
 * Update the node role (sender/receiver)
 *
 * @param {'sender' | 'receiver'} newRole - New role value
 */
function updateRole(newRole: 'sender' | 'receiver'): void {
  emitUpdate({ role: newRole, pairLabel: '' });
}

/**
 * Update the pair label (sender types freely, receiver picks from senders)
 *
 * @param {string} label - New pair label
 */
function updatePairLabel(label: string): void {
  emitUpdate({ pairLabel: label });
}

/**
 * Handle receiver selecting a sender label — also auto-match color
 *
 * @param {string} label - Selected sender label
 */
function selectSenderLabel(label: string): void {
  const sender = senderLabels.value.find(s => s.label === label);
  if (sender) {
    emitUpdate({ pairLabel: label, pairColor: sender.color });
  } else {
    emitUpdate({ pairLabel: label });
  }
}

/**
 * Update the pair color from preset
 *
 * @param {string} color - Quasar color name from presets
 */
function updatePairColor(color: string): void {
  showHexInput.value = false;
  emitUpdate({ pairColor: color });
}

/**
 * Open the custom hex input
 */
function openHexInput(): void {
  hexBuffer.value = colorHex.value;
  showHexInput.value = true;
}

/**
 * Apply the custom hex color
 */
function applyHexColor(): void {
  const hex = hexBuffer.value.trim();
  if (/^#[0-9a-fA-F]{6}$/.test(hex)) {
    emitUpdate({ pairColor: hex });
  }
}
</script>

<template>
  <div class="goto-config">
    <!-- ROLE section -->
    <div class="goto-config__section">
      <div class="goto-config__section-label">{{ t('nodes.goto.config.roleSection') }}</div>
      <q-btn-toggle
        :model-value="role"
        spread
        no-caps
        dense
        unelevated
        toggle-color="primary"
        :options="[
          { value: 'sender', label: t('nodes.goto.config.sender'), icon: 'near_me' },
          { value: 'receiver', label: t('nodes.goto.config.receiver'), icon: 'place' },
        ]"
        class="goto-config__role-toggle"
        @update:model-value="updateRole"
      />
      <div class="goto-config__hint">
        {{ role === 'sender'
          ? t('nodes.goto.config.senderHint')
          : t('nodes.goto.config.receiverHint')
        }}
      </div>
    </div>

    <!-- SENDER: Label input (free text) -->
    <div v-if="role === 'sender'" class="goto-config__section">
      <div class="goto-config__section-label">{{ t('nodes.goto.config.labelSection') }}</div>
      <q-input
        :model-value="pairLabel"
        outlined
        dense
        :placeholder="t('nodes.goto.config.labelPlaceholder')"
        @update:model-value="(val: string | number | null) => updatePairLabel(String(val ?? ''))"
      >
        <template #prepend>
          <q-icon name="label" size="xs" :style="{ color: colorHex }" />
        </template>
      </q-input>
      <div class="goto-config__hint">
        {{ t('nodes.goto.config.labelHint') }}
      </div>
    </div>

    <!-- RECEIVER: Select from existing senders -->
    <div v-if="role === 'receiver'" class="goto-config__section">
      <div class="goto-config__section-label">{{ t('nodes.goto.config.targetSenderSection') }}</div>

      <div v-if="senderLabels.length === 0" class="goto-config__empty-state">
        <q-icon name="link_off" size="xs" color="grey-6" class="q-mr-sm" />
        <span>{{ t('nodes.goto.config.noSendersAvailable') }}</span>
      </div>

      <q-select
        v-else
        :model-value="pairLabel"
        outlined
        dense
        emit-value
        map-options
        options-dense
        :options="senderLabels.map(s => ({ value: s.label, label: s.label, hex: s.hex }))"
        option-value="value"
        option-label="label"
        :placeholder="t('nodes.goto.config.selectSenderLabel')"
        @update:model-value="selectSenderLabel"
      >
        <template #prepend>
          <q-icon name="place" size="xs" :style="{ color: colorHex }" />
        </template>
        <template #option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section side>
              <div
                class="goto-config__color-dot"
                :style="{ background: scope.opt.hex }"
              />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ scope.opt.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
        <template #selected-item="scope">
          <div class="row items-center no-wrap q-gutter-xs">
            <div
              class="goto-config__color-dot"
              :style="{ background: scope.opt?.hex || colorHex }"
            />
            <span>{{ scope.opt?.label || pairLabel }}</span>
          </div>
        </template>
      </q-select>
    </div>

    <!-- COLOR section (sender only — receiver inherits sender color) -->
    <div v-if="role === 'sender'" class="goto-config__section">
      <div class="goto-config__section-label">{{ t('nodes.goto.config.colorSection') }}</div>
      <div class="goto-config__color-grid">
        <div
          v-for="opt in GOTO_COLOR_OPTIONS"
          :key="opt.value"
          class="goto-config__color-swatch"
          :class="{ 'goto-config__color-swatch--active': pairColor === opt.value }"
          :style="{ background: opt.hex }"
          @click="updatePairColor(opt.value)"
        >
          <q-icon
            v-if="pairColor === opt.value"
            name="check"
            size="12px"
            color="white"
          />
        </div>
        <!-- Custom hex swatch -->
        <div
          class="goto-config__color-swatch"
          :class="{ 'goto-config__color-swatch--active': isCustomColor }"
          :style="{ background: isCustomColor ? colorHex : 'var(--mapex-surface-elevated)', border: isCustomColor ? undefined : '1px dashed var(--mapex-text-muted)' }"
          @click="openHexInput"
        >
          <q-icon
            :name="isCustomColor ? 'check' : 'colorize'"
            :size="isCustomColor ? '12px' : '14px'"
            :color="isCustomColor ? 'white' : 'grey-6'"
          />
        </div>
      </div>

      <!-- Custom hex input -->
      <div v-if="showHexInput" class="goto-config__hex-row">
        <q-input
          v-model="hexBuffer"
          outlined
          dense
          placeholder="#ff5722"
          mask="\#XXXXXX"
          class="goto-config__hex-input"
          @keyup.enter="applyHexColor"
        >
          <template #prepend>
            <div
              class="goto-config__color-dot"
              :style="{ background: hexBuffer }"
            />
          </template>
        </q-input>
        <q-btn
          flat
          dense
          icon="check"
          color="positive"
          size="sm"
          @click="applyHexColor"
        />
      </div>
    </div>

    <!-- MATCHED PAIRS section -->
    <div v-if="pairLabel" class="goto-config__section">
      <div class="goto-config__section-label">{{ t('nodes.goto.config.matchedPairsSection') }}</div>

      <div v-if="matchedPairs.length === 0" class="goto-config__empty-state">
        <q-icon name="link_off" size="xs" color="grey-6" class="q-mr-sm" />
        <span>{{ t('nodes.goto.config.noMatchingGoto') }} "{{ pairLabel }}"</span>
      </div>

      <div v-else class="goto-config__pairs-list">
        <div
          v-for="pair in matchedPairs"
          :key="pair.id"
          class="goto-config__pair-item"
        >
          <q-icon
            :name="pair.role === 'sender' ? 'near_me' : 'place'"
            size="14px"
            :style="{ color: GOTO_COLOR_OPTIONS.find(o => o.value === pair.color)?.hex || '#5e35b1' }"
          />
          <span class="goto-config__pair-label">{{ pair.label }}</span>
          <q-badge
            :label="pair.role === 'sender' ? t('nodes.goto.config.sendBadge') : t('nodes.goto.config.recvBadge')"
            :style="{ background: GOTO_COLOR_OPTIONS.find(o => o.value === pair.color)?.hex || '#5e35b1' }"
            class="goto-config__pair-badge"
          />
        </div>
      </div>
    </div>

    <!-- INFO banner -->
    <div class="goto-config__info">
      <q-icon name="info" color="grey-6" size="xs" class="q-mr-sm" />
      <span>
        {{ t('nodes.goto.config.matchedPairsHint') }}
      </span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.goto-config {
  &__section {
    margin-bottom: 16px;
  }

  &__section-label {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--mapex-text-secondary);
    margin-bottom: 6px;
    text-transform: uppercase;
  }

  &__role-toggle {
    width: 100%;
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
  }

  &__hint {
    font-size: 0.7rem;
    color: var(--mapex-text-muted);
    margin-top: 4px;
    font-style: italic;
  }

  &__color-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 8px;
  }

  &__color-swatch {
    width: 100%;
    aspect-ratio: 1;
    border-radius: var(--mapex-radius-md);
    cursor: pointer;
    border: 2px solid transparent;
    transition:
      transform var(--mapex-transition-base),
      border-color var(--mapex-transition-base);
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      transform: scale(1.12);
    }

    &--active {
      border-color: var(--mapex-wf-text-on-accent);
      box-shadow: var(--mapex-wf-selection-ring);
    }
  }

  &__color-dot {
    width: 10px;
    height: 10px;
    border-radius: var(--mapex-radius-full);
    flex-shrink: 0;
  }

  &__hex-row {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-top: 8px;
  }

  &__hex-input {
    flex: 1;
  }

  &__empty-state {
    display: flex;
    align-items: center;
    padding: 8px 10px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-1);
    font-size: 0.75rem;
    color: var(--mapex-text-muted);
  }

  &__pairs-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  &__pair-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 8px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-1);
    border: 1px solid var(--mapex-card-border);
  }

  &__pair-label {
    flex: 1;
    font-size: 0.8rem;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &__pair-badge {
    font-size: 0.6rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  &__info {
    display: flex;
    align-items: flex-start;
    padding: 10px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    line-height: 1.4;
  }
}
</style>

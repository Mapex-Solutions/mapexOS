<script setup lang="ts">
defineOptions({
  name: 'GroupFrameNodeConfig',
});

/** TYPE IMPORTS */
import type {
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
} from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { FRAME_COLOR_OPTIONS } from '../constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { t } = usePluginI18n('core-annotations');

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
 * Current title from config
 */
const title = computed<string>(
  () => (props.config.title as string) || '',
);

/**
 * Current description from config
 */
const description = computed<string>(
  () => (props.config.description as string) || '',
);

/**
 * Current color name from config
 */
const colorName = computed<string>(
  () => (props.config.color as string) || 'blue-grey',
);

/**
 * Whether the current color is a custom hex (not in presets)
 */
const isCustomColor = computed(() =>
  !FRAME_COLOR_OPTIONS.some(o => o.value === colorName.value),
);

/**
 * Hex color resolved from FRAME_COLOR_OPTIONS or custom
 */
const colorHex = computed<string>(() => {
  const opt = FRAME_COLOR_OPTIONS.find(o => o.value === colorName.value);
  if (opt) return opt.hex;
  return colorName.value.startsWith('#') ? colorName.value : '#78909c';
});

/**
 * Current width from config
 */
const width = computed<number>(
  () => (props.config.width as number) || 300,
);

/**
 * Current height from config
 */
const height = computed<number>(
  () => (props.config.height as number) || 200,
);

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
 * Update the frame title
 *
 * @param {string | number | null} val - Input value
 */
function updateTitle(val: string | number | null): void {
  emitUpdate({ title: String(val ?? '') });
}

/**
 * Update the frame description
 *
 * @param {string | number | null} val - Input value
 */
function updateDescription(val: string | number | null): void {
  emitUpdate({ description: String(val ?? '') });
}

/**
 * Update the frame color from preset
 *
 * @param {string} color - Color value from presets
 */
function updateColor(color: string): void {
  showHexInput.value = false;
  emitUpdate({ color });
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
    emitUpdate({ color: hex });
  }
}

/**
 * Update the frame width
 *
 * @param {string | number | null} val - Input value
 */
function updateWidth(val: string | number | null): void {
  const num = Number(val);
  if (!isNaN(num) && num >= 150) {
    emitUpdate({ width: num });
  }
}

/**
 * Update the frame height
 *
 * @param {string | number | null} val - Input value
 */
function updateHeight(val: string | number | null): void {
  const num = Number(val);
  if (!isNaN(num) && num >= 100) {
    emitUpdate({ height: num });
  }
}
</script>

<template>
  <div class="frame-config">
    <!-- TITLE section -->
    <div class="frame-config__section">
      <div class="frame-config__section-label">{{ t('nodes.group_frame.config.titleSection') }}</div>
      <q-input
        :model-value="title"
        outlined
        dense
        :placeholder="t('nodes.group_frame.config.titlePlaceholder')"
        @update:model-value="updateTitle"
      >
        <template #prepend>
          <q-icon name="title" size="xs" :style="{ color: colorHex }" />
        </template>
      </q-input>
    </div>

    <!-- DESCRIPTION section -->
    <div class="frame-config__section">
      <div class="frame-config__section-label">{{ t('nodes.group_frame.config.descriptionSection') }}</div>
      <q-input
        :model-value="description"
        outlined
        dense
        type="textarea"
        autogrow
        :placeholder="t('nodes.group_frame.config.descriptionPlaceholder')"
        @update:model-value="updateDescription"
      />
    </div>

    <!-- COLOR section -->
    <div class="frame-config__section">
      <div class="frame-config__section-label">{{ t('nodes.group_frame.config.colorSection') }}</div>
      <div class="frame-config__color-grid">
        <div
          v-for="opt in FRAME_COLOR_OPTIONS"
          :key="opt.value"
          class="frame-config__color-swatch"
          :class="{ 'frame-config__color-swatch--active': colorName === opt.value }"
          :style="{ background: opt.hex }"
          @click="updateColor(opt.value)"
        >
          <q-icon
            v-if="colorName === opt.value"
            name="check"
            size="12px"
            color="white"
          />
        </div>
        <!-- Custom hex swatch -->
        <div
          class="frame-config__color-swatch"
          :class="{ 'frame-config__color-swatch--active': isCustomColor }"
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
      <div v-if="showHexInput" class="frame-config__hex-row">
        <q-input
          v-model="hexBuffer"
          outlined
          dense
          placeholder="#78909c"
          mask="\#XXXXXX"
          class="frame-config__hex-input"
          @keyup.enter="applyHexColor"
        >
          <template #prepend>
            <div
              class="frame-config__color-dot"
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

    <!-- SIZE section -->
    <div class="frame-config__section">
      <div class="frame-config__section-label">{{ t('nodes.group_frame.config.sizeSection') }}</div>
      <div class="frame-config__size-row">
        <q-input
          :model-value="width"
          outlined
          dense
          type="number"
          :min="150"
          :label="t('nodes.group_frame.config.widthLabel')"
          class="frame-config__size-input"
          @update:model-value="updateWidth"
        />
        <q-input
          :model-value="height"
          outlined
          dense
          type="number"
          :min="100"
          :label="t('nodes.group_frame.config.heightLabel')"
          class="frame-config__size-input"
          @update:model-value="updateHeight"
        />
      </div>
      <div class="frame-config__hint">
        {{ t('nodes.group_frame.config.resizeHint') }}
      </div>
    </div>

    <!-- INFO banner -->
    <div class="frame-config__info">
      <q-icon name="info" color="grey-6" size="xs" class="q-mr-sm" />
      <span>
        {{ t('nodes.group_frame.config.infoHint') }}
      </span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.frame-config {
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

  &__color-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
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

  &__size-row {
    display: flex;
    gap: 8px;
  }

  &__size-input {
    flex: 1;
  }

  &__hint {
    font-size: 0.7rem;
    color: var(--mapex-text-muted);
    margin-top: 4px;
    font-style: italic;
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

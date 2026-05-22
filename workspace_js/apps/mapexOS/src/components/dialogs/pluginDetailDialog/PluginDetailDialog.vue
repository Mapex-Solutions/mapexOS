<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginDetailDialogProps, PluginDetailDialogEmits } from './interfaces/pluginDetailDialog.interface';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** PROPS & EMITS */
const props = defineProps<PluginDetailDialogProps>();
const emit = defineEmits<PluginDetailDialogEmits>();

/** STATE */

/**
 * Internal dialog state synced with v-model
 */
const isOpen = ref(props.modelValue);

/** WATCHERS */

watch(() => props.modelValue, (val) => {
  isOpen.value = val;
});

watch(isOpen, (val) => {
  emit('update:modelValue', val);
});

/** FUNCTIONS */

/**
 * Close the dialog
 *
 * @returns {void}
 */
function handleClose(): void {
  emit('update:modelValue', false);
}
</script>

<template>
  <q-dialog v-model="isOpen" @hide="handleClose">
    <q-card class="plugin-detail-dialog">
      <!-- Header -->
      <q-card-section class="plugin-detail-dialog__header">
        <div class="row items-center no-wrap">
          <q-avatar v-if="brandIconUrl" size="48px" class="q-mr-md" square>
            <img :src="brandIconUrl" :alt="name" />
          </q-avatar>
          <q-icon v-else :name="icon" size="48px" :style="{ color: color }" class="q-mr-md" />
          <div class="col">
            <div class="plugin-detail-dialog__title">
              {{ name }}
            </div>
            <div class="plugin-detail-dialog__subtitle">
              {{ author }} &middot; v{{ version }}
            </div>
          </div>
          <span
            class="plugin-detail-dialog__category-badge q-mr-md"
            :style="{ color: color, borderColor: color }"
          >
            {{ categoryLabel }}
          </span>
          <q-btn
            v-close-popup
            flat
            round
            dense
            icon="close"
            class="plugin-detail-dialog__close-btn"
          />
        </div>
        <div v-if="description" class="plugin-detail-dialog__description q-mt-sm">
          {{ description }}
        </div>

        <!-- Tags -->
        <div v-if="tags.length > 0" class="row q-gutter-xs q-mt-sm">
          <span
            v-for="tag in tags"
            :key="tag"
            class="plugin-detail-dialog__tag"
          >
            {{ tag }}
          </span>
        </div>
      </q-card-section>

      <div class="plugin-detail-dialog__separator" />

      <!-- Node Types Section -->
      <q-card-section>
        <div class="row items-center q-mb-md">
          <q-icon name="account_tree" size="sm" class="q-mr-sm" style="color: var(--mapex-primary)" />
          <span class="plugin-detail-dialog__section-title">
            {{ nodeTypesLabel }}
          </span>
          <span v-if="!loading" class="plugin-detail-dialog__count-badge q-ml-sm">
            {{ nodeTypes.length }}
          </span>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex flex-center q-pa-lg">
          <q-spinner size="32px" class="q-mr-sm" style="color: var(--mapex-primary)" />
          <span class="plugin-detail-dialog__loading-text">{{ loadingLabel }}</span>
        </div>

        <!-- Node list (flat — name + description only) -->
        <q-list v-else dense separator class="plugin-detail-dialog__node-list">
          <q-item v-for="node in nodeTypes" :key="node.type">
            <q-item-section avatar>
              <q-icon :name="node.icon" :style="{ color: node.color }" size="24px" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="plugin-detail-dialog__node-label">
                {{ node.label }}
              </q-item-label>
              <q-item-label class="plugin-detail-dialog__node-description">
                {{ node.description }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row q-gutter-xs">
                <span class="plugin-detail-dialog__io-badge">
                  {{ node.inputCount }} {{ inputsLabel }}
                </span>
                <span class="plugin-detail-dialog__io-badge">
                  {{ node.outputCount }} {{ outputsLabel }}
                </span>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <div class="plugin-detail-dialog__separator" />

      <!-- Footer: Install action -->
      <q-card-actions class="q-pa-md" align="right">
        <q-btn
          v-if="installed"
          flat
          dense
          no-caps
          icon="check_circle"
          disable
          :label="installedLabel"
          class="plugin-detail-dialog__installed-btn"
        />
        <q-btn
          v-else
          unelevated
          no-caps
          icon="download"
          :label="installing ? installingLabel : installLabel"
          :loading="installing"
          :disable="installDisabled"
          class="plugin-detail-dialog__install-btn"
          @click="emit('install')"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style lang="scss" scoped>
.plugin-detail-dialog {
  min-width: 550px;
  max-width: 700px;
  background: var(--mapex-surface-elevated);
  border-radius: var(--mapex-radius-md);
  box-shadow: var(--mapex-shadow-lg);

  &__header {
    padding: 20px 24px !important;
    background: var(--mapex-surface-bg);
    border-bottom: 1px solid var(--mapex-divider);
  }

  &__title {
    font-size: var(--mapex-font-lg);
    font-weight: var(--mapex-font-weight-bold);
    color: var(--mapex-text-primary);
    line-height: var(--mapex-line-height-tight);
  }

  &__subtitle {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
    line-height: var(--mapex-line-height-base);
  }

  &__description {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-secondary);
    line-height: var(--mapex-line-height-relaxed);
  }

  &__category-badge {
    display: inline-flex;
    align-items: center;
    padding: 2px 8px;
    font-size: var(--mapex-font-xs);
    font-weight: var(--mapex-font-weight-medium);
    border: 1px solid;
    border-radius: var(--mapex-radius-xs);
    white-space: nowrap;
  }

  &__close-btn {
    color: var(--mapex-text-muted);

    &:hover {
      color: var(--mapex-text-primary);
    }
  }

  &__tag {
    display: inline-flex;
    align-items: center;
    padding: 1px 6px;
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-xs);
  }

  &__separator {
    height: 1px;
    background: var(--mapex-divider);
  }

  &__section-title {
    font-size: var(--mapex-font-sm);
    font-weight: var(--mapex-font-weight-bold);
    color: var(--mapex-text-primary);
  }

  &__count-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    padding: 0 6px;
    font-size: var(--mapex-font-xs);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--mapex-primary);
    border: 1px solid var(--mapex-primary);
    border-radius: var(--mapex-radius-xs);
  }

  &__loading-text {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }

  &__node-list {
    background: var(--mapex-card-bg, var(--mapex-surface-bg));
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-sm);

    :deep(.q-separator) {
      background: var(--mapex-divider);
    }
  }

  &__node-label {
    font-size: var(--mapex-font-sm);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--mapex-text-primary);
  }

  &__node-description {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }

  &__io-badge {
    display: inline-flex;
    align-items: center;
    padding: 1px 6px;
    font-size: var(--mapex-font-2xs);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--mapex-text-secondary);
    background: var(--mapex-surface-sunken);
    border-radius: var(--mapex-radius-xs);
    white-space: nowrap;
  }

  &__installed-btn {
    color: var(--q-positive);
  }

  &__install-btn {
    background: var(--mapex-primary);
    color: var(--mapex-surface-bg);
  }

  // Override Quasar defaults for q-card-section padding
  :deep(.q-card__section) {
    padding: 16px 24px;
  }

  :deep(.q-card__actions) {
    padding: 12px 24px;
  }

  // Ensure q-item separators use our token
  :deep(.q-item + .q-separator) {
    background: var(--mapex-divider);
  }
}
</style>

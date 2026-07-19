<script setup lang="ts">
defineOptions({
  name: 'MarketplaceCard'
});

/** TYPE IMPORTS */
import type { MarketplaceCardProps, MarketplaceCardEmits } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { BaseButton } from '@components/buttons';

/** COMPOSABLES */
import { useAssetTemplateMarketplaceTranslations } from '@composables/i18n';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<MarketplaceCardProps>();
const emit = defineEmits<MarketplaceCardEmits>();

/** COMPOSABLES & STORES */
const t = useAssetTemplateMarketplaceTranslations();

/** COMPUTED */

/**
 * Resolved image URL for the template, when the catalog item carries one.
 */
const imageUrl = computed<string | undefined>(() => {
  if (!props.item.hasImage || !props.item.image) return undefined;
  return apis.assetTemplatesMarketplace?.marketplace.assetUrl(
    props.item.vendor,
    props.item.slug,
    props.item.image,
  );
});

/**
 * Material icon fallback shown when the item has no image.
 */
const iconName = computed(() => props.item.icon || 'memory');

/**
 * Vendor display name — the card title. Falls back to the vendor slug.
 */
const vendorName = computed(() => props.item.vendorName || props.item.vendor);

/**
 * Category chip text: the friendly label when resolved, else the raw slug.
 */
const categoryText = computed(() => props.categoryLabel || props.item.category);

// Cap the model at a fixed character count so version always lands at the end of
// the row (no accordion between short and long model names); the full name stays
// reachable via the tooltip.
const MODEL_MAX_CHARS = 18;

/**
 * Whether the model name exceeds the display cap and is shown truncated.
 */
const modelTruncated = computed(() => (props.item.model?.length ?? 0) > MODEL_MAX_CHARS);

/**
 * The model text as rendered: truncated with an ellipsis when over the cap.
 */
const modelDisplay = computed(() =>
  modelTruncated.value ? `${props.item.model.slice(0, MODEL_MAX_CHARS).trimEnd()}…` : props.item.model,
);

/** FUNCTIONS */

/**
 * Emit the view event with the current catalog item.
 */
function handleView(): void {
  emit('view', props.item);
}

/**
 * Emit the install event; the parent runs the install so the card stays
 * presentational. Stops propagation so it does not also open the preview.
 */
function handleInstall(): void {
  emit('install', props.item);
}

/**
 * Emit the uninstall event; the parent runs the removal so the card stays
 * presentational. Stops propagation so it does not also open the preview.
 */
function handleUninstall(): void {
  emit('uninstall', props.item);
}
</script>

<template>
  <q-card
    flat
    bordered
    class="marketplace-card cursor-pointer"
    @click="handleView"
  >
    <AppTooltip :content="t.card.view.value" />

    <!-- Media: image or icon, with the category pill overlaid top-right -->
    <div class="marketplace-card__media">
      <q-img
        v-if="imageUrl"
        :src="imageUrl"
        :ratio="16 / 9"
        fit="contain"
        class="marketplace-card__image"
      />
      <div v-else class="marketplace-card__icon">
        <q-icon :name="iconName" size="44px" color="primary" />
      </div>

      <div v-if="categoryText" class="marketplace-card__category">
        <q-icon name="category" size="14px" />
        <span>{{ categoryText }}</span>
      </div>
    </div>

    <q-card-section class="marketplace-card__body">
      <!-- Vendor is the title -->
      <div class="marketplace-card__vendor">
        <q-icon name="factory" size="18px" class="marketplace-card__vendor-icon" />
        <span class="marketplace-card__vendor-name">{{ vendorName }}</span>
      </div>

      <!-- Template name is the subtitle -->
      <div class="marketplace-card__name">{{ item.name }}</div>

      <!-- Model (truncated + tooltip) on the left, version pinned to the end -->
      <div class="marketplace-card__meta">
        <div v-if="item.model" class="marketplace-card__meta-item marketplace-card__meta-item--model">
          <span class="marketplace-card__meta-label">{{ t.card.model.value }}</span>
          <span class="marketplace-card__meta-value marketplace-card__meta-value--model">
            {{ modelDisplay }}
            <AppTooltip v-if="modelTruncated" :content="item.model" />
          </span>
        </div>
        <div v-if="item.version" class="marketplace-card__meta-item marketplace-card__meta-item--version">
          <span class="marketplace-card__meta-label">{{ t.card.version.value }}</span>
          <span class="marketplace-card__meta-value">{{ item.version }}</span>
        </div>
      </div>
    </q-card-section>

    <q-separator />

    <!-- Footer: install when the template is not yet installed, uninstall when it is -->
    <div class="marketplace-card__footer">
      <BaseButton
        v-if="installed"
        outline
        no-caps
        color="negative"
        icon="delete"
        class="marketplace-card__action"
        :label="installing ? t.uninstall.uninstalling.value : t.uninstall.button.value"
        :loading="installing"
        @click.stop="handleUninstall"
      />
      <BaseButton
        v-else
        unelevated
        no-caps
        color="primary"
        icon="download"
        class="marketplace-card__action"
        :label="installing ? t.install.installing.value : t.install.button.value"
        :loading="installing"
        @click.stop="handleInstall"
      />
    </div>
  </q-card>
</template>

<style scoped>
.marketplace-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-radius: var(--mapex-radius-md);
  transition: box-shadow var(--mapex-transition-base), transform var(--mapex-transition-base);
}

.marketplace-card:hover {
  box-shadow: var(--mapex-shadow-lg);
  transform: translateY(-2px);
}

/* Media */
.marketplace-card__media {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--mapex-spacing-md);
  background: var(--mapex-surface-sunken);
  border-bottom: 1px solid var(--mapex-divider);
}

.marketplace-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 88px;
}

.marketplace-card__image {
  max-width: 100%;
}

.marketplace-card__category {
  position: absolute;
  top: var(--mapex-spacing-sm);
  right: var(--mapex-spacing-sm);
  display: inline-flex;
  align-items: center;
  gap: var(--mapex-spacing-xs);
  padding: 2px 8px;
  border-radius: var(--mapex-radius-pill, 999px);
  background: var(--mapex-surface-elevated);
  border: 1px solid var(--mapex-divider);
  color: var(--mapex-text-secondary);
  font-size: var(--mapex-font-2xs);
  font-weight: var(--mapex-font-weight-semibold);
  text-transform: capitalize;
}

/* Body */
.marketplace-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--mapex-spacing-xs);
  flex: 1;
}

.marketplace-card__vendor {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--mapex-text-primary);
}

.marketplace-card__vendor-icon {
  color: var(--mapex-primary);
  flex-shrink: 0;
}

.marketplace-card__vendor-name {
  font-size: var(--mapex-font-md, 1rem);
  font-weight: var(--mapex-font-weight-bold, 700);
  line-height: 1.25;
}

.marketplace-card__name {
  color: var(--mapex-text-secondary);
  font-size: var(--mapex-font-sm);
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Model / version metadata — version stays pinned to the row end */
.marketplace-card__meta {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--mapex-spacing-md);
  margin-top: var(--mapex-spacing-xs);
}

.marketplace-card__meta-item {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.marketplace-card__meta-item--version {
  flex-shrink: 0;
  text-align: right;
}

.marketplace-card__meta-value--model {
  white-space: nowrap;
}

.marketplace-card__meta-label {
  font-size: var(--mapex-font-2xs);
  font-weight: var(--mapex-font-weight-semibold);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--mapex-text-secondary);
}

.marketplace-card__meta-value {
  font-size: var(--mapex-font-sm);
  font-weight: var(--mapex-font-weight-medium, 500);
  color: var(--mapex-text-primary);
}

/* Footer */
.marketplace-card__footer {
  padding: var(--mapex-spacing-sm) var(--mapex-spacing-md) var(--mapex-spacing-md);
}

.marketplace-card__action {
  width: 100%;
}
</style>

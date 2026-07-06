<script setup lang="ts">
defineOptions({
  name: 'MarketplaceCard'
});

/** TYPE IMPORTS */
import type { MarketplaceCardProps, MarketplaceCardEmits } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

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
 * Field-count label, e.g. "12 fields".
 */
const fieldsLabel = computed(() => t.card.fields(props.item.fieldCount));

/** FUNCTIONS */

/**
 * Emit the view event with the current catalog item.
 */
function handleView(): void {
  emit('view', props.item);
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

    <q-card-section class="marketplace-card__media">
      <q-img
        v-if="imageUrl"
        :src="imageUrl"
        :ratio="16 / 9"
        fit="contain"
        class="marketplace-card__image"
      />
      <div v-else class="marketplace-card__icon">
        <q-icon :name="iconName" size="48px" color="primary" />
      </div>
    </q-card-section>

    <q-card-section class="marketplace-card__body">
      <div class="marketplace-card__name">{{ item.name }}</div>

      <div class="marketplace-card__chips">
        <DetailChip
          v-if="item.vendorName || item.vendor"
          :label="item.vendorName || item.vendor"
          icon="factory"
          color="indigo"
          size="sm"
          outline
        />
        <DetailChip
          v-if="item.model"
          :label="item.model"
          icon="memory"
          color="blue"
          size="sm"
          outline
        />
        <DetailChip
          v-if="item.version"
          :label="item.version"
          icon="sell"
          color="purple"
          size="sm"
          outline
        />
      </div>

      <div v-if="item.category" class="marketplace-card__chips">
        <DetailChip
          :label="item.category"
          icon="category"
          color="teal"
          size="sm"
        />
      </div>

      <div class="marketplace-card__meta">
        <span class="marketplace-card__fields">
          <q-icon name="list_alt" size="16px" />
          {{ fieldsLabel }}
        </span>
        <span v-if="item.hasScripts" class="marketplace-card__scripts">
          <q-icon name="code" size="16px" />
          {{ t.card.hasScripts.value }}
        </span>
      </div>
    </q-card-section>
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

.marketplace-card__media {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--mapex-spacing-md);
}

.marketplace-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 96px;
}

.marketplace-card__image {
  max-width: 100%;
}

.marketplace-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--mapex-spacing-sm);
}

.marketplace-card__name {
  font-weight: 600;
  color: var(--mapex-text-primary);
}

.marketplace-card__chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--mapex-spacing-xs);
}

.marketplace-card__meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--mapex-spacing-md);
  color: var(--mapex-text-secondary);
  font-size: 0.85rem;
}

.marketplace-card__fields,
.marketplace-card__scripts {
  display: inline-flex;
  align-items: center;
  gap: var(--mapex-spacing-xs);
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'DlqSidebar'
});

/** TYPE IMPORTS */
import type { DlqServiceTypeGroup } from '../../interfaces';

/** COMPOSABLES */
import { useDlqLogsPageTranslations } from '@composables/i18n/pages/logs/dlqLogsPage';

/** COMPOSABLES & STORES */
const t = useDlqLogsPageTranslations();

/** PROPS & EMITS */
defineProps<{
  activeServiceType: string | null;
  serviceTypeGroups: DlqServiceTypeGroup[];
  totalCount: number;
}>();

const emit = defineEmits<{
  'update:activeServiceType': [serviceType: string | null];
}>();
</script>

<template>
  <div class="dlq-sidebar">
    <!-- All Failures -->
    <q-list dense class="q-px-xs q-pt-md">
      <q-item
        clickable
        :active="activeServiceType === null"
        active-class="dlq-sidebar__item--active"
        class="dlq-sidebar__item rounded-borders"
        @click="emit('update:activeServiceType', null)"
      >
        <q-item-section avatar style="min-width: 36px;">
          <q-icon name="all_inbox" size="sm" class="dlq-sidebar__all-icon" />
        </q-item-section>
        <q-item-section>
          <span class="text-body2 text-weight-medium">{{ t.sidebar.allFailures.value }}</span>
        </q-item-section>
        <q-item-section v-if="totalCount > 0" side>
          <q-badge
            class="dlq-sidebar__badge-count text-caption"
            rounded
            :label="totalCount"
          />
        </q-item-section>
      </q-item>
    </q-list>

    <q-separator class="q-my-sm q-mx-md" />

    <!-- Service Types -->
    <q-list dense class="q-px-xs">
      <q-item class="rounded-borders" style="min-height: 28px; padding-left: 12px;">
        <q-item-section>
          <span class="dlq-sidebar__section-label">{{ t.sidebar.byService.value }}</span>
        </q-item-section>
      </q-item>

      <q-item
        v-for="group in serviceTypeGroups"
        :key="group.serviceType"
        clickable
        :active="activeServiceType === group.serviceType"
        active-class="dlq-sidebar__item--active"
        class="dlq-sidebar__item rounded-borders"
        :class="{ 'dlq-sidebar__item--disabled': group.count === 0 }"
        @click="emit('update:activeServiceType', group.serviceType)"
      >
        <q-item-section avatar style="min-width: 36px;">
          <q-icon :name="group.icon" size="sm" class="dlq-sidebar__service-icon" />
        </q-item-section>
        <q-item-section>
          <span class="text-body2">{{ group.serviceType }}</span>
        </q-item-section>
        <q-item-section v-if="group.count > 0" side>
          <q-badge
            class="dlq-sidebar__badge-service text-caption"
            rounded
            :label="group.count"
          />
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Info Banner -->
    <div class="dlq-sidebar__info-banner q-mx-sm q-mt-lg q-pa-sm rounded-borders">
      <q-icon name="info" size="xs" class="q-mr-xs" />
      {{ t.sidebar.infoBanner.value }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dlq-sidebar {
  width: 240px;
  height: 100%;
  overflow-y: auto;
  border-right: 1px solid var(--mapex-card-border);
  background: var(--mapex-sidebar-bg);
  flex-shrink: 0;

  &__all-icon {
    color: var(--mapex-text-primary);
  }

  &__section-label {
    font-size: var(--mapex-font-xs);
    font-weight: var(--mapex-font-weight-bold);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--mapex-text-muted);
  }

  &__service-icon {
    color: var(--mapex-text-secondary);
  }

  &__badge-count {
    background: var(--mapex-danger);
    color: var(--mapex-surface-bg);
  }

  &__badge-service {
    background: var(--mapex-text-muted);
    color: var(--mapex-surface-bg);
    font-size: var(--mapex-font-2xs);
  }

  &__info-banner {
    font-size: var(--mapex-font-2xs);
    background: var(--mapex-danger-hover);
    color: var(--mapex-danger);
    border: 1px solid var(--mapex-danger-border);
    line-height: var(--mapex-line-height-relaxed);
  }

  &__item {
    height: 36px;
    margin: 1px 4px;
    border-radius: var(--mapex-radius-sm);
    transition: var(--mapex-transition-fast);

    &--active {
      background: var(--mapex-active-bg) !important;
      color: var(--mapex-primary);
      font-weight: var(--mapex-font-weight-medium);
    }

    &--disabled {
      opacity: 0.4;
    }

    &:hover:not(.dlq-sidebar__item--active) {
      background: var(--mapex-submenu-bg);
    }
  }
}
</style>

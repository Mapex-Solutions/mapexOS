<script setup lang="ts">
defineOptions({
  name: 'GenericSelectorDialog'
});

/** TYPE IMPORTS */
import type {
  GenericSelectorDialogProps,
  GenericSelectorDialogEmits,
  ScrollInfo,
} from './interfaces/genericSelectorDialog.interface';

/** VUE IMPORTS */
import { ref, computed, watch, onBeforeUnmount } from 'vue';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import {
  DEFAULT_DIALOG_WIDTH,
  DEFAULT_SEARCH_DEBOUNCE_MS,
  DEFAULT_SCROLL_THRESHOLD,
  DEFAULT_ITEM_KEY,
  DEFAULT_EMPTY_ICON,
  DEFAULT_EMPTY_TEXT,
  DEFAULT_LOADING_TEXT,
  DEFAULT_CONFIRM_LABEL,
  DEFAULT_CANCEL_LABEL,
  DEFAULT_ITEM_NOUN_SINGULAR,
  DEFAULT_ITEM_NOUN_PLURAL,
  DEFAULT_INFO_BANNER,
  DEFAULT_ACTIVE_ITEM_STYLE,
} from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<GenericSelectorDialogProps>(), {
  iconColor: 'primary',
  itemKey: DEFAULT_ITEM_KEY,
  multiSelect: false,
  selectedIds: () => [],
  loading: false,
  loadingMore: false,
  hasMorePages: false,
  searchPlaceholder: 'Search...',
  width: DEFAULT_DIALOG_WIDTH,
  searchDebounce: DEFAULT_SEARCH_DEBOUNCE_MS,
  emptyIcon: DEFAULT_EMPTY_ICON,
  emptyText: DEFAULT_EMPTY_TEXT,
  loadingText: DEFAULT_LOADING_TEXT,
  confirmLabel: DEFAULT_CONFIRM_LABEL,
  cancelLabel: DEFAULT_CANCEL_LABEL,
  itemNounSingular: DEFAULT_ITEM_NOUN_SINGULAR,
  itemNounPlural: DEFAULT_ITEM_NOUN_PLURAL,
  showFiltersHeader: true,
  showSearch: true,
  scrollThreshold: DEFAULT_SCROLL_THRESHOLD,
});

const emit = defineEmits<GenericSelectorDialogEmits>();

/** COMPOSABLES & STORES */
const tsTitle = useTS({ titleCase: true });
const basePath = 'components.dialogs.genericSelector';

/** STATE */
const searchQuery = ref('');
const selectedItems = ref<any[]>([]);
const debounceTimer = ref<ReturnType<typeof setTimeout> | null>(null);

/** COMPUTED */

/**
 * Dialog visibility model (getter/setter for v-model)
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Card inline style with width and CSS custom properties for active item colors
 */
const cardStyle = computed(() => {
  const activeBg = props.activeItemStyle?.backgroundColor ?? DEFAULT_ACTIVE_ITEM_STYLE.backgroundColor;
  const activeBorder = props.activeItemStyle?.borderColor ?? DEFAULT_ACTIVE_ITEM_STYLE.borderColor;

  return {
    width: `${props.width}px`,
    maxWidth: '90vw',
    display: 'flex',
    flexDirection: 'column' as const,
    height: '85vh',
    maxHeight: '85vh',
    '--selector-active-bg': activeBg,
    '--selector-active-border': activeBorder,
  };
});

/**
 * Set of selected IDs for O(1) lookup
 */
const selectedIdsSet = computed(() => {
  return new Set(selectedItems.value.map(item => item[props.itemKey]));
});

/**
 * Whether the confirm button should be enabled
 */
const canConfirm = computed(() => selectedItems.value.length > 0);

/**
 * Resolved info banner with defaults
 */
const resolvedBanner = computed(() => {
  if (!props.infoBanner) return null;
  return {
    icon: props.infoBanner.icon ?? DEFAULT_INFO_BANNER.icon,
    bgClass: props.infoBanner.bgClass ?? DEFAULT_INFO_BANNER.bgClass,
    textClass: props.infoBanner.textClass ?? DEFAULT_INFO_BANNER.textClass,
    iconColor: props.infoBanner.iconColor ?? DEFAULT_INFO_BANNER.iconColor,
    text: props.infoBanner.text,
  };
});

/**
 * Footer count text (e.g. "5 items", "1 item")
 */
const footerCountText = computed(() => {
  const count = props.totalItems ?? props.items.length;
  const noun = count === 1 ? props.itemNounSingular : props.itemNounPlural;
  return `${count} ${noun}`;
});

/**
 * Whether the filters section should be visible
 * Shows when search is enabled or filters slot is used (checked in template via $slots)
 */
const showFiltersSection = computed(() => props.showSearch);

/** WATCHERS */

/**
 * Initialize selection from props when dialog opens
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    searchQuery.value = '';
    initializeSelection();
  }
});

/**
 * Preserve selection when items change (new pages loaded)
 */
watch(() => props.items, () => {
  // Selection is maintained by ID — no action needed since selectedItems
  // stores full item objects and selectedIdsSet uses itemKey
}, { deep: false });

/** FUNCTIONS */

/**
 * Initialize selected items from props.selectedIds
 * Matches items already in the list; preserves previous selections for items not yet loaded
 *
 * @returns {void}
 */
function initializeSelection(): void {
  if (!props.selectedIds || props.selectedIds.length === 0) {
    selectedItems.value = [];
    return;
  }

  const idsSet = new Set(props.selectedIds);
  selectedItems.value = props.items.filter(item => idsSet.has(item[props.itemKey]));
}

/**
 * Check if an item is currently selected
 * @param {any} item - Item to check
 * @returns {boolean} True if the item is selected
 */
function isSelected(item: any): boolean {
  return selectedIdsSet.value.has(item[props.itemKey]);
}

/**
 * Toggle item selection (add/remove from selected)
 * For single-select: emits immediately and closes dialog
 * For multi-select: toggles in internal array
 *
 * @param {any} item - Item to toggle
 * @returns {void}
 */
function toggleItem(item: any): void {
  if (!props.multiSelect) {
    emit('select', [item]);
    showDialog.value = false;
    return;
  }

  const itemId = item[props.itemKey];
  const index = selectedItems.value.findIndex(i => i[props.itemKey] === itemId);

  if (index >= 0) {
    selectedItems.value.splice(index, 1);
  } else {
    selectedItems.value.push(item);
  }
}

/**
 * Confirm multi-select selection and close dialog
 *
 * @returns {void}
 */
function confirmSelection(): void {
  emit('select', [...selectedItems.value]);
  showDialog.value = false;
}

/**
 * Handle cancel action
 *
 * @returns {void}
 */
function handleCancel(): void {
  emit('cancel');
  showDialog.value = false;
}

/**
 * Handle search input with debounce
 * Updates display value immediately, debounces the emit
 *
 * @param {string | number | null} value - Input value from q-input
 * @returns {void}
 */
function handleSearchInput(value: string | number | null): void {
  const query = String(value ?? '');
  searchQuery.value = query;

  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value);
  }

  debounceTimer.value = setTimeout(() => {
    emit('search', query);
  }, props.searchDebounce);
}

/**
 * Handle search clear — bypass debounce and emit immediately
 *
 * @returns {void}
 */
function handleSearchClear(): void {
  searchQuery.value = '';

  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value);
    debounceTimer.value = null;
  }

  emit('search', '');
}

/**
 * Infinite scroll handler
 * Emits load-more when scroll reaches threshold
 *
 * @param {ScrollInfo} info - Scroll position information from q-scroll-area
 * @returns {void}
 */
function onScroll(info: ScrollInfo): void {
  const { verticalPosition, verticalSize, verticalContainerSize } = info;

  if (verticalPosition + verticalContainerSize >= verticalSize * props.scrollThreshold) {
    if (!props.loadingMore && props.hasMorePages) {
      emit('load-more');
    }
  }
}

/** LIFECYCLE HOOKS */
onBeforeUnmount(() => {
  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value);
  }
});
</script>

<template>
  <q-dialog v-model="showDialog" @escape-key="handleCancel">
    <q-card :style="cardStyle" class="generic-selector-dialog">
      <!-- Header -->
      <q-card-section class="selector-header q-pb-sm">
        <div class="row items-center">
          <q-icon
            v-if="icon"
            :name="icon"
            :color="iconColor"
            size="sm"
            class="q-mr-sm"
          />
          <div class="selector-title">{{ title }}</div>
          <q-space />
          <q-btn
            icon="close"
            flat
            round
            dense
            class="selector-close-btn"
            @click="handleCancel"
          />
        </div>
      </q-card-section>

      <!-- Info Banner -->
      <q-card-section
        v-if="resolvedBanner || $slots['info-banner']"
        class="q-pt-none q-pb-md"
      >
        <slot name="info-banner">
          <div v-if="resolvedBanner" class="selector-banner" :class="[resolvedBanner.bgClass, resolvedBanner.textClass]">
            <q-icon :name="resolvedBanner.icon" :color="resolvedBanner.iconColor" size="sm" />
            <span class="text-caption">{{ resolvedBanner.text }}</span>
          </div>
        </slot>
      </q-card-section>

      <!-- Selected Preview (multi-select) -->
      <slot name="selected-preview" />

      <!-- Filters Section -->
      <q-card-section v-if="showFiltersSection || $slots.filters" class="q-py-md">
        <div v-if="showFiltersHeader" class="selector-section-label q-mb-md">
          <q-icon name="filter_list" size="xs" class="q-mr-xs" />
          {{ tsTitle(`${basePath}.filters`) }}
        </div>
        <div class="row q-col-gutter-md">
          <!-- Search Input -->
          <div v-if="showSearch" class="col-12">
            <q-input
              :model-value="searchQuery"
              outlined
              dense
              :placeholder="searchPlaceholder"
              clearable
              :debounce="0"
              @update:model-value="handleSearchInput"
              @clear="handleSearchClear"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Domain-specific Filters -->
          <slot name="filters" />
        </div>
      </q-card-section>

      <q-separator class="selector-separator" />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="selector-section-label">
          <q-icon v-if="resultsIcon" :name="resultsIcon" size="xs" class="q-mr-xs" />
          {{ tsTitle(`${basePath}.results`) }}
        </div>
      </q-card-section>

      <!-- Items List -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="selector-muted-text q-mt-md">{{ loadingText }}</div>
        </div>

        <!-- Empty State -->
        <div v-else-if="items.length === 0" class="q-pa-md text-center">
          <slot name="empty">
            <q-icon :name="emptyIcon" size="4em" class="selector-empty-icon" />
            <div class="selector-muted-text q-mt-md">{{ emptyText }}</div>
          </slot>
        </div>

        <!-- Items List with Infinite Scroll -->
        <q-scroll-area
          v-else
          style="height: 100%;"
          @scroll="onScroll"
        >
          <q-list separator>
            <q-item
              v-for="item in items"
              :key="item[itemKey]"
              clickable
              :active="isSelected(item)"
              @click="toggleItem(item)"
            >
              <!-- Checkbox for multi-select -->
              <q-item-section v-if="multiSelect" avatar>
                <q-checkbox
                  :model-value="isSelected(item)"
                  color="primary"
                  @click.stop="toggleItem(item)"
                />
              </q-item-section>

              <!-- Item content (slot) -->
              <slot
                name="item"
                :item="item"
                :is-selected="isSelected(item)"
                :toggle="() => toggleItem(item)"
              />
            </q-item>
          </q-list>

          <!-- Load More Indicator -->
          <div v-if="loadingMore" class="q-pa-md text-center">
            <q-spinner color="primary" size="2em" />
          </div>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator class="selector-separator" />
      <q-card-actions class="selector-footer">
        <div class="selector-footer-count">
          <q-icon v-if="footerIcon" :name="footerIcon" size="xs" class="q-mr-xs" />
          {{ footerCountText }}
        </div>
        <q-space />
        <q-btn
          flat
          dense
          :label="cancelLabel"
          size="sm"
          class="selector-btn-cancel"
          @click="handleCancel"
        />
        <q-btn
          v-if="multiSelect"
          flat
          dense
          :label="confirmLabel"
          color="primary"
          size="sm"
          class="selector-btn-confirm"
          :disable="!canConfirm"
          @click="confirmSelection"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style lang="scss" scoped>
.generic-selector-dialog {
  background: var(--mapex-popup-bg) !important;
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-lg) !important;
  box-shadow: var(--mapex-shadow-xl);
}

/* Header */
.selector-header {
  border-bottom: 1px solid var(--mapex-divider);
}

.selector-title {
  font-size: 1.15rem;
  font-weight: 600;
  color: var(--mapex-text-primary);
}

.selector-close-btn {
  color: var(--mapex-text-secondary);
  border-radius: var(--mapex-radius-md);

  &:hover {
    background: var(--mapex-surface-bg);
  }
}

/* Info Banner — theme-aware */
.selector-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 18px;
  border-radius: var(--mapex-radius-md);
  background: rgba(var(--mapex-primary-rgb), 0.1);
  color: var(--mapex-text-primary);
  border: 1px solid rgba(var(--mapex-primary-rgb), 0.2);
}

/* Section labels (Filters, Results) */
.selector-section-label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--mapex-text-secondary);
}

/* Separator */
.selector-separator {
  background: var(--mapex-divider) !important;
}

/* Muted text (loading, empty state) */
.selector-muted-text {
  color: var(--mapex-text-muted);
}

/* Empty state icon */
.selector-empty-icon {
  color: var(--mapex-text-muted);
  opacity: 0.6;
}

/* Footer */
.selector-footer {
  padding: 12px 16px !important;
  background: var(--mapex-surface-bg);
  border-bottom-left-radius: var(--mapex-radius-lg);
  border-bottom-right-radius: var(--mapex-radius-lg);
}

.selector-footer-count {
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
}

.selector-btn-cancel {
  color: var(--mapex-text-secondary);
  border-radius: var(--mapex-radius-md);
}

.selector-btn-confirm {
  border-radius: var(--mapex-radius-md);
}

/* Hover effects for list items */
:deep(.q-item) {
  transition: var(--mapex-transition-base);
  color: var(--mapex-text-primary);
}

:deep(.q-item:hover) {
  background-color: var(--mapex-surface-highlight) !important;
}

:deep(.q-item.q-item--active) {
  background-color: var(--selector-active-bg) !important;
  border-left: 3px solid var(--selector-active-border);
}

/* Item labels should use design tokens */
:deep(.q-item__label--caption) {
  color: var(--mapex-text-secondary) !important;
}

/* Filter input border radius */
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

/* Smooth transitions for badges and chips */
:deep(.q-badge),
:deep(.q-chip) {
  transition: var(--mapex-transition-base);
}

/* Checkbox corner rounding */
:deep(.q-checkbox__inner) {
  border-radius: var(--mapex-radius-xs);
}

/* Separator inside dialog */
:deep(.q-separator) {
  background: var(--mapex-divider) !important;
}
</style>

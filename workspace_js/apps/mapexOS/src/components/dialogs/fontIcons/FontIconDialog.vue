<script setup lang="ts">
defineOptions({
  name: 'FontIconDialog',
});

/** TYPE IMPORTS */
import type { FontIconDialogProps, FontIconDialogEvents, FontIconData } from './interfaces/fontIconDialog.interface';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';

/** COMPOSABLES */
import { useCommonPlaceholders } from '@composables/i18n';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { CATEGORY_ICON_MAP, CATEGORY_DESCRIPTION_MAP } from './constants';
import iconsData from './fonts/icons.json';

/** PROPS & EMITS */
const props = defineProps<FontIconDialogProps>();
const emit = defineEmits<FontIconDialogEvents>();

/** COMPOSABLES & STORES */
const { placeholders } = useCommonPlaceholders();

/** STATE */

/**
 * Search filter text
 */
const filter = ref('');

/**
 * Currently selected icon name
 */
const selectedIcon = ref('');

/**
 * Currently selected category key
 */
const selectedCategory = ref<string>('');

/**
 * All icons loaded from JSON
 */
const allIcons = ref<FontIconData[]>([]);

/** COMPUTED */

/**
 * Two-way computed for dialog visibility
 */
const show = computed<boolean>({
  get() {
    return props.show;
  },
  set(val: boolean) {
    emit('update:show', val);
  },
});

/**
 * Sorted unique category keys
 */
const categories = computed(() =>
  Array.from(
    new Set(allIcons.value.flatMap(i => i.categories)),
  ).sort(),
);

/**
 * Icons filtered by search term or selected category
 */
const filteredIcons = computed(() => {
  let result = allIcons.value;
  if (filter.value.trim()) {
    const term = filter.value.toLowerCase();
    result = result.filter(i => i.name.toLowerCase().includes(term));
  } else if (selectedCategory.value) {
    result = result.filter(i => i.categories.includes(selectedCategory.value));
  }
  return result;
});

/** WATCHERS */

/**
 * Sync selected icon from modelValue prop
 */
watch(() => props.modelValue, (newValue) => {
  selectedIcon.value = newValue;
}, { immediate: true });

/**
 * Clear search when switching categories
 */
watch(() => selectedCategory.value, () => {
  filter.value = '';
});

/** FUNCTIONS */

/**
 * Format a category key into a human-readable title
 *
 * @param {string} cat - Category key (e.g., 'flow_control')
 * @returns {string} Formatted name (e.g., 'Flow Control')
 */
function formatCategoryName(cat: string): string {
  return cat.replace(/_/g, ' ').split(' ')
    .map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ');
}

/**
 * Get the Material icon name for a category
 *
 * @param {string} cat - Category key
 * @returns {string} Material icon name
 */
function getCategoryIcon(cat: string): string {
  return CATEGORY_ICON_MAP[cat] || 'category';
}

/**
 * Get the description text for a category
 *
 * @param {string} cat - Category key
 * @returns {string} Category description
 */
function getCategoryDescription(cat: string): string {
  return CATEGORY_DESCRIPTION_MAP[cat] || 'Material Design icons';
}

/**
 * Select an icon in the grid
 *
 * @param {string} name - Icon name to select
 * @returns {void}
 */
function selectIcon(name: string): void {
  selectedIcon.value = name;
}

/**
 * Confirm selection and close dialog
 *
 * @returns {void}
 */
function confirmSelection(): void {
  if (selectedIcon.value) {
    emit('update:modelValue', selectedIcon.value);
    show.value = false;
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  allIcons.value = iconsData.icons.map((i: any) => ({
    name: i.name,
    categories: i.categories.length ? i.categories : ['uncategorized'],
  }));
  if (!selectedCategory.value && categories.value.length > 0) {
    selectedCategory.value = categories.value[0]!;
  }
});
</script>

<template>
  <q-dialog v-model="show">
    <q-card class="icon-picker-dialog">
      <!-- Header -->
      <q-card-section class="dialog-header">
        <div class="row items-center">
          <q-icon name="category" size="md" color="primary" class="q-mr-sm" />
          <div>
            <div class="text-h5 text-weight-medium" style="color: var(--mapex-text-primary)">
              Select Material Design Icon
            </div>
            <div class="text-body2" style="color: var(--mapex-text-secondary)">
              Choose an icon for your application
            </div>
          </div>
          <q-space />
          <q-btn
            v-close-popup
            flat
            round
            dense
            icon="close"
            color="grey-7"
          />
        </div>
      </q-card-section>

      <!-- Main content area (scrollable) -->
      <q-card-section class="q-pa-none dialog-content">
        <q-separator />
        <div class="row no-wrap q-pa-md full-height">
          <!-- Left sidebar: categories -->
          <div class="col-12 col-md-4 q-pr-md">
            <q-card flat bordered class="stepper-card">
              <!-- Fixed header section -->
              <q-card-section class="bg-primary text-white q-pb-md fixed-header">
                <div class="text-h6 text-weight-bold q-mb-sm">
                  <q-icon size="sm" class="q-mr-xs" name="category" />
                  Icon Categories
                </div>
                <div class="text-caption">Browse icons by category</div>
              </q-card-section>

              <!-- Scrollable categories section -->
              <q-card-section class="categories-scroll-area">
                <div class="progress-steps">
                  <div
                    v-for="(cat, idx) in categories"
                    :key="idx"
                    class="step-item"
                    :class="{ active: selectedCategory === cat }"
                    @click="selectedCategory = cat"
                  >
                    <div class="step-icon-wrapper">
                      <div class="step-icon">
                        <q-icon
                          size="sm"
                          :name="getCategoryIcon(cat)"
                          :color="selectedCategory === cat ? 'white' : 'grey-5'"
                        />
                      </div>
                    </div>
                    <div class="step-content">
                      <div class="step-title">{{ formatCategoryName(cat) }}</div>
                      <div class="step-description">{{ getCategoryDescription(cat) }}</div>
                    </div>
                  </div>
                </div>
              </q-card-section>

              <!-- Fixed footer section -->
              <q-card-section class="fixed-footer">
                <q-separator />
                <div class="current-step-info q-mt-md">
                  <div class="text-caption text-grey-6">
                    <q-icon size="xs" class="q-mr-xs" name="info" />
                    Click on an icon to select it
                  </div>
                  <div class="text-caption text-primary q-mt-xs">
                    <q-icon size="xs" class="q-mr-xs" name="arrow_forward" />
                    Current Category: {{ formatCategoryName(selectedCategory) }}
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>

          <!-- Main content: search + icons grid -->
          <div class="col-12 col-md-8 q-pl-md">
            <q-card flat bordered class="form-card">
              <!-- Section Header -->
              <q-card-section class="section-header q-pb-md">
                <div class="text-h6 text-weight-bold text-primary">
                  <q-icon size="sm" class="q-mr-xs" :name="getCategoryIcon(selectedCategory)" />
                  {{ formatCategoryName(selectedCategory) }}
                </div>
                <div class="text-caption" style="color: var(--mapex-text-secondary)">
                  {{ getCategoryDescription(selectedCategory) }}
                </div>
              </q-card-section>

              <q-card-section class="q-px-md q-py-md icons-body">
                <!-- Search input -->
                <q-input
                  v-model="filter"
                  outlined
                  dense
                  class="q-mb-md"
                  :placeholder="placeholders.search.value"
                >
                  <template #prepend>
                    <q-icon name="search" color="primary" />
                  </template>
                  <template v-if="filter" #append>
                    <q-badge rounded color="primary">
                      {{ filteredIcons.length }}
                    </q-badge>
                  </template>
                </q-input>

                <!-- Icons grid -->
                <q-scroll-area class="icons-scroll-area">
                  <div class="row q-col-gutter-sm q-pa-sm">
                    <template v-if="filteredIcons.length > 0">
                      <div
                        v-for="icon in filteredIcons"
                        :key="icon.name"
                        class="col-4 col-sm-3 col-md-3 col-lg-2"
                      >
                        <q-card
                          clickable
                          flat
                          bordered
                          class="icon-card"
                          :class="{ 'icon-selected': selectedIcon === icon.name }"
                          @click="selectIcon(icon.name)"
                        >
                          <q-card-section class="q-pa-sm column items-center justify-center">
                            <q-icon size="24px" :name="icon.name" />
                          </q-card-section>
                        </q-card>
                      </div>
                    </template>
                    <div v-else class="col-12 column items-center justify-center q-pa-xl">
                      <q-icon size="48px" name="search_off" color="grey-5" />
                      <div class="text-subtitle1 text-grey-7 q-mt-md">No icons found</div>
                      <div class="text-caption text-grey-6">Try a different search term</div>
                    </div>
                  </div>
                </q-scroll-area>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-card-section>

      <!-- Fixed footer with action buttons -->
      <q-separator />
      <q-card-actions align="right" class="dialog-footer">
        <q-btn
          v-close-popup
          flat
          no-caps
          class="q-mr-sm"
          label="Cancel"
          color="grey-7"
        />
        <q-btn
          unelevated
          no-caps
          icon-right="check"
          label="Select Icon"
          color="primary"
          :disable="!selectedIcon"
          @click="confirmSelection"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style lang="scss" scoped>
.icon-picker-dialog {
  width: 1200px;
  max-width: 90vw;
  display: flex;
  flex-direction: column;
  max-height: 85vh;
}

.dialog-header {
  flex-shrink: 0;
  padding: 20px 24px !important;
  background: var(--mapex-surface-bg);
  border-bottom: 1px solid var(--mapex-divider);
}

.dialog-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.dialog-footer {
  flex-shrink: 0;
  padding: 12px 16px !important;
}

.full-height {
  height: 100%;
}

.stepper-card,
.form-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: var(--mapex-radius-md);
}

.fixed-header {
  flex-shrink: 0;
}

.categories-scroll-area {
  flex: 1;
  overflow-y: auto;
  padding-top: 0;
  padding-bottom: 0;
}

.fixed-footer {
  flex-shrink: 0;
  padding-bottom: 16px;
}

.section-header {
  flex-shrink: 0;
  background: var(--mapex-surface-bg);
}

.icons-body {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 8px 0;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 8px 0;
  position: relative;
  cursor: pointer;
}

.step-item:not(:last-child)::after {
  content: '';
  position: absolute;
  left: 19px;
  top: 40px;
  width: 2px;
  height: 24px;
  background-color: var(--mapex-card-border);
  z-index: 0;
}

.step-icon-wrapper {
  position: relative;
  z-index: 1;
}

.step-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--mapex-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--mapex-surface-elevated);
  border: 2px solid var(--mapex-card-border);
  transition: var(--mapex-transition-slow);
}

.step-item.active .step-icon {
  background-color: var(--q-primary);
  border-color: var(--q-primary);
}

.step-content {
  flex: 1;
  padding-top: 2px;
}

.step-title {
  font-weight: 600;
  font-size: 14px;
  color: var(--mapex-text-primary);
  margin-bottom: 4px;
}

.step-item.active .step-title {
  color: var(--q-primary);
}

.step-description {
  font-size: 12px;
  color: var(--mapex-text-secondary);
  line-height: 1.4;
}

.current-step-info {
  background-color: var(--mapex-surface-elevated);
  padding: 12px;
  border-radius: var(--mapex-radius-sm);
  border-left: 3px solid var(--q-primary);
}

.icons-scroll-area {
  flex: 1;
  min-height: 200px;
}

.icon-card {
  height: 56px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: var(--mapex-radius-sm);
  transition: var(--mapex-transition-base);
}

.icon-card:hover {
  background-color: var(--mapex-surface-highlight);
  transform: translateY(-2px);
  box-shadow: var(--mapex-shadow-xs);
}

.icon-card.icon-selected {
  background-color: var(--q-primary) !important;
  color: white !important;
  border-color: var(--q-primary) !important;
  box-shadow: var(--mapex-shadow-xs);
}

.icon-card.icon-selected:hover {
  background-color: var(--q-primary) !important;
  color: white !important;
  transform: none;
}

.icon-card.icon-selected :deep(.q-icon) {
  color: white !important;
}

@media (max-width: 1023px) {
  .icon-picker-dialog {
    width: 95vw;
  }
}
</style>

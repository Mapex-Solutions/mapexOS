<template>
  <div class="data-row-wrapper">
    <!-- DESKTOP/TABLET VERSION (>= 600px) -->
    <q-table
      flat
      hide-header
      hide-pagination
      class="data-row-table"
      :rows="[props.data]"
      :columns="tableColumns"
      :row-key="props.primaryKey"
      :rows-per-page-options="[1]"
  >
    <!-- Custom Row -->
    <template v-slot:body="bodyProps">
      <q-tr
          class="data-row-card"
          :props="bodyProps"
          @click="handleRowClick"
          @dblclick="handleRowDblClick"
      >
        <!-- Avatar Column -->
        <q-td
            v-for="col in bodyProps.cols"
            :props="bodyProps"
            :key="col.name"
            :class="getColumnClass(col)"
            :style="getColumnStyleForTable(col)"
        >
          <div v-if="col.label && col.type !== 'avatar'" class="text-caption text-grey-5 text-weight-medium q-mb-xs">
            {{ col.label }}
          </div>
                    
          <component
              :is="getColumnComponent(col)"
              :value="getColumnValue(col)"
              :column="col"
              :row="props.data"
          />
        </q-td>

        <!-- Actions Menu -->
        <q-td v-if="props.showActions" class="data-row-actions-cell">
          <q-btn
              flat
              dense
              round
              icon="more_vert"
              color="grey-7"
              @click.stop
          >
            <q-menu
              anchor="bottom end"
              self="top end"
              :offset="[0, 4]"
              class="data-row-actions-menu"
            >
              <q-list class="data-row-actions-list">
                <!-- Custom Actions -->
                <template v-for="action in visibleCustomActions" :key="action.key">
                  <q-item
                    v-close-popup
                    clickable
                    class="data-row-action-item"
                    @click.stop="emit('action', action.key, props.data)"
                  >
                    <q-item-section side>
                      <q-icon :name="action.icon" size="22px" :color="action.color || 'teal-7'" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label>{{ action.label }}</q-item-label>
                      <q-item-label v-if="action.description" caption>
                        {{ action.description }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </template>

                <q-separator v-if="visibleCustomActions.length > 0 && (showEdit || showView || showDelete)" class="action-separator" />

                <!-- Edit Action -->
                <q-item
                  v-if="showEdit"
                  v-close-popup
                  clickable
                  class="data-row-action-item"
                  @click.stop="emit('edit', props.data)"
                >
                  <q-item-section side>
                    <q-icon name="edit" size="22px" color="teal-7" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>Edit</q-item-label>
                    <q-item-label caption>Modify this item</q-item-label>
                  </q-item-section>
                </q-item>

                <!-- View Details Action -->
                <q-item
                  v-if="showView"
                  v-close-popup
                  clickable
                  class="data-row-action-item"
                  @click.stop="emit('view', props.data)"
                >
                  <q-item-section side>
                    <q-icon name="visibility" size="22px" color="blue-7" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>View Details</q-item-label>
                    <q-item-label caption>See full information</q-item-label>
                  </q-item-section>
                </q-item>

                <q-separator v-if="(showEdit || showView) && showDelete" class="action-separator" />

                <!-- Delete Action -->
                <q-item
                  v-if="showDelete"
                  v-close-popup
                  clickable
                  class="data-row-action-item data-row-action-item--danger"
                  @click.stop="emit('delete', props.data)"
                >
                  <q-item-section side>
                    <q-icon name="delete_outline" size="22px" color="red-7" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-red-7">Delete</q-item-label>
                    <q-item-label caption>Remove permanently</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </q-tr>
    </template>
  </q-table>

  <!-- MOBILE VERSION (< 600px) -->
  <q-card
      flat
      bordered
      class="data-row-mobile"
  >
    <q-card-section class="q-pa-md" @click="toggleMobileExpand">
      <div class="row items-center q-gutter-sm">
        <!-- Avatar -->
        <component
            v-if="avatarColumn"
            :is="getColumnComponent(avatarColumn)"
            :value="getColumnValue(avatarColumn)"
            :column="avatarColumn"
            :row="props.data"
            :mobile="true"
        />

        <!-- Name + Description -->
        <div class="col" style="min-width: 0;">
          <component
              v-if="nameColumn"
              :is="getColumnComponent(nameColumn)"
              :value="getColumnValue(nameColumn)"
              :column="nameColumn"
              :row="props.data"
              :mobile="true"
          />
        </div>

        <!-- Status Badge -->
        <component
            v-if="statusColumn"
            :is="getColumnComponent(statusColumn)"
            :value="getColumnValue(statusColumn)"
            :column="statusColumn"
            :row="props.data"
            :mobile="true"
        />

        <!-- Actions Menu -->
        <q-btn
            v-if="props.showActions"
            flat
            dense
            round
            size="sm"
            color="grey-7"
            icon="more_vert"
            @click.stop
        >
          <q-menu
            anchor="bottom end"
            self="top end"
            :offset="[0, 4]"
            class="data-row-actions-menu"
          >
            <q-list class="data-row-actions-list">
              <!-- Custom Actions -->
              <template v-for="action in visibleCustomActions" :key="action.key">
                <q-item
                  v-close-popup
                  clickable
                  class="data-row-action-item"
                  @click.stop="emit('action', action.key, props.data)"
                >
                  <q-item-section side>
                    <q-icon :name="action.icon" size="22px" :color="action.color || 'teal-7'" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ action.label }}</q-item-label>
                    <q-item-label v-if="action.description" caption>
                      {{ action.description }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </template>

              <q-separator v-if="visibleCustomActions.length > 0 && (showEdit || showView || showDelete)" class="action-separator" />

              <!-- Edit Action -->
              <q-item
                v-if="showEdit"
                v-close-popup
                clickable
                class="data-row-action-item"
                @click.stop="emit('edit', props.data)"
              >
                <q-item-section side>
                  <q-icon name="edit" size="22px" color="teal-7" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>Edit</q-item-label>
                  <q-item-label caption>Modify this item</q-item-label>
                </q-item-section>
              </q-item>

              <!-- View Details Action -->
              <q-item
                v-if="showView"
                v-close-popup
                clickable
                class="data-row-action-item"
                @click.stop="emit('view', props.data)"
              >
                <q-item-section side>
                  <q-icon name="visibility" size="22px" color="blue-7" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>View Details</q-item-label>
                  <q-item-label caption>See full information</q-item-label>
                </q-item-section>
              </q-item>

              <q-separator v-if="(showEdit || showView) && showDelete" class="action-separator" />

              <!-- Delete Action -->
              <q-item
                v-if="showDelete"
                v-close-popup
                clickable
                class="data-row-action-item data-row-action-item--danger"
                @click.stop="emit('delete', props.data)"
              >
                <q-item-section side>
                  <q-icon name="delete_outline" size="22px" color="red-7" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-red-7">Delete</q-item-label>
                  <q-item-label caption>Remove permanently</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </q-card-section>

    <!-- Expanded Details -->
    <q-slide-transition>
      <q-card-section v-show="mobileExpanded" class="q-pt-none">
        <q-separator class="q-mb-md" />
        <div class="row q-col-gutter-md">
          <template v-for="col in mobileExpandableColumns" :key="col.key">
            <div class="col-6">
              <div class="text-caption text-grey-5 text-weight-medium q-mb-xs">{{ col.label }}</div>
              <component
                  :is="getColumnComponent(col)"
                  :value="getColumnValue(col)"
                  :column="col"
                  :row="props.data"
                  :mobile="true"
              />
            </div>
          </template>
        </div>
      </q-card-section>
    </q-slide-transition>
  </q-card>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'DataRow'
});

import { ref, computed } from 'vue';
import type { DataRowProps, DataRowEmits, DataRowColumn } from './interfaces';
import { AvatarColumn, TextColumn, CodeColumn, ChipColumn, ChipsColumn, BadgeColumn } from './columns';

/** CONSTANTS */
const CLICK_DELAY_MS = 250;

const props = withDefaults(defineProps<DataRowProps>(), {
  primaryKey: 'id',
  showActions: true,
  expandOnClick: true,
});

const emit = defineEmits<DataRowEmits>();

/** STATE */
let clickTimer: ReturnType<typeof setTimeout> | null = null;
const mobileExpanded = ref(false);

/** COMPUTED */
const showEdit = computed(() => {
  if (!props.actions) return true;
  return props.actions.showEdit !== false;
});

const showView = computed(() => {
  if (!props.actions) return true;
  return props.actions.showView !== false;
});

const showDelete = computed(() => {
  if (!props.actions) return true;
  return props.actions.showDelete !== false;
});

const visibleCustomActions = computed(() => {
  const customActions = props.actions?.customActions || [];
  return customActions.filter((action) => {
    if (action.condition) {
      return action.condition(props.data);
    }
    return true;
  });
});

// Convert columns to q-table format
// Following Quasar Table pattern: columns define structure, values extracted on render
const tableColumns = computed(() => {
  return props.columns
    .map(col => ({
      ...col,
      name: col.key,
      field: col.key,
      align: 'left' as const,
      style: col.width ? `width: ${col.width}px; max-width: ${col.width}px;` : '',
    }));
});

// Mobile-specific columns
const avatarColumn = computed(() =>
  props.columns.find(col => col.type === 'avatar')
);

const nameColumn = computed(() =>
  props.columns.find(col => col.key === 'name')
);

const statusColumn = computed(() =>
  props.columns.find(col => col.key === 'status')
);

const mobileExpandableColumns = computed(() =>
  props.columns.filter(col =>
    col.type !== 'avatar' &&
    col.key !== 'name' &&
    col.key !== 'status'
  )
);

/**
 * Handle row click with timer-based single/double click detection
 * Waits CLICK_DELAY_MS to distinguish between single and double click
 * Prevents viewer from opening when user intends to double-click for edit
 *
 * @returns {void}
 */
function handleRowClick(): void {
  if (clickTimer) {
    // Second click within delay = double click
    clearTimeout(clickTimer);
    clickTimer = null;
    emit('dblclick', props.data);
    emit('edit', props.data);
  } else {
    // First click - wait to confirm if single click
    clickTimer = setTimeout(() => {
      clickTimer = null;
      emit('click', props.data);
    }, CLICK_DELAY_MS);
  }
}

/**
 * Handle native dblclick event as fallback safety
 * Clears any pending single click timer to prevent race conditions
 *
 * @returns {void}
 */
function handleRowDblClick(): void {
  if (clickTimer) {
    clearTimeout(clickTimer);
    clickTimer = null;
  }
}

/**
 * Extract value from data object using column key
 * Supports nested properties via dot notation (e.g., 'protocol.type')
 * Following Quasar Table pattern: extract raw value, let column components handle formatting
 */
function getColumnValue(column: DataRowColumn) {
  const keys = column.key.split('.');
  let value = props.data;

  for (const key of keys) {
    value = value?.[key];
  }

  return value;
}

function getColumnClass(column: any) {
  const classes = [`data-row-cell`, `data-row-cell--${column.type}`];

  // Add Quasar responsive visibility classes
  // Maps intuitive aliases to Quasar breakpoint classes
  if (column.visible && column.visible !== 'always') {
    const visibilityClass = mapVisibilityToQuasar(column.visible);
    if (visibilityClass) {
      classes.push(visibilityClass);
    }
  }

  // Add text alignment class
  if (column.align) {
    const alignMap: Record<string, string> = {
      'left': 'text-left',
      'center': 'text-center',
      'right': 'text-right',
    };
    const alignClass = alignMap[column.align];
    if (alignClass) {
      classes.push(alignClass);
    }
  }

  return classes.join(' ');
}

/**
 * Maps intuitive device aliases to Quasar visibility classes
 * Allows using simple names like 'laptop' instead of 'gt-sm'
 *
 * @param visibility - The visibility value (alias or Quasar class)
 * @returns Quasar visibility class or null if 'always'
 */
function mapVisibilityToQuasar(visibility: string): string | null {
  // Intuitive device aliases (PREFERRED)
  const aliasMap: Record<string, string> = {
    'mobile': 'lt-md',      // < 1024px (xs, sm)
    'tablet': 'gt-sm',      // >= 1024px (md, lg, xl) - same as laptop
    'laptop': 'gt-sm',      // >= 1024px (md, lg, xl) - MOST COMMON
    'desktop': 'gt-md',     // >= 1440px (lg, xl)
  };

  // Return mapped alias or pass through Quasar class directly
  return aliasMap[visibility] || visibility;
}

function getColumnStyleForTable(column: any) {
  const styles: string[] = [];

  if (column.width) {
    styles.push(`width: ${column.width}px`);
    styles.push(`max-width: ${column.width}px`);
    styles.push(`min-width: ${column.width}px`);
  }

  return styles.join('; ');
}

function getColumnComponent(column: any) {
  switch (column.type) {
    case 'avatar':
      return AvatarColumn;
    case 'text':
      return TextColumn;
    case 'code':
      return CodeColumn;
    case 'chip':
      return ChipColumn;
    case 'chips':
      return ChipsColumn;
    case 'badge':
      return BadgeColumn;
    default:
      return TextColumn;
  }
}

function toggleMobileExpand() {
  if (props.expandOnClick) {
    mobileExpanded.value = !mobileExpanded.value;
  }
}
</script>

<style scoped lang="scss">
// Remove table appearance
:deep(.data-row-table .q-table__card) {
  box-shadow: none;
}

:deep(.data-row-table .q-table tbody td) {
  border: none;
}

:deep(.data-row-table thead) {
  display: none;
}

// Card-like row styling
:deep(.data-row-card) {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-xs);
  margin-bottom: 8px;
  transition: var(--mapex-transition-base);
  cursor: pointer;
}

:deep(.data-row-card:hover) {
  background-color: var(--mapex-surface-elevated);
  box-shadow: 0 2px 6px var(--mapex-elevation-shadow);
}

:deep(.data-row-card td) {
  padding: 16px 8px;
  vertical-align: middle;
}

:deep(.data-row-card td:first-child) {
  padding-left: 16px;
}

:deep(.data-row-card td:last-child) {
  padding-right: 8px;
}

// Column cells
:deep(.data-row-cell) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

// Quasar's visibility classes are used natively (gt-sm, lt-md, etc.)
// No custom CSS needed - Quasar handles all responsive behavior

// Actions cell - always on the right
:deep(.data-row-actions-cell) {
  width: 56px !important;
  max-width: 56px !important;
  min-width: 56px !important;
  text-align: right !important;
  padding-right: 16px !important;
}

// Expanded area
:deep(.data-row-expanded) {
  background-color: var(--mapex-surface-elevated);
  border-top: 1px solid var(--mapex-card-border);
}

:deep(.data-row-expanded td) {
  padding: 16px;
}

// Mobile version (< 600px)
.data-row-table {
  display: block;

  @media (max-width: 600px) {
    display: none !important;
  }
}

.data-row-mobile {
  display: none;

  @media (max-width: 600px) {
    display: block;
  }
}
</style>

<!-- Non-scoped styles for menu (teleported to body) -->
<style lang="scss">
.data-row-actions-menu {
  border-radius: var(--mapex-radius-lg) !important;
  box-shadow:
    0 4px 6px -1px var(--mapex-elevation-shadow),
    0 2px 4px -1px var(--mapex-elevation-shadow),
    0 20px 25px -5px var(--mapex-elevation-shadow) !important;
  border: 1px solid var(--mapex-card-border) !important;
  overflow: hidden !important;
  background: var(--mapex-popup-bg) !important;

  .data-row-actions-list {
    min-width: 240px;
    padding: 6px !important;
    background: var(--mapex-popup-bg) !important;
  }

  .data-row-action-item {
    border-radius: var(--mapex-radius-md) !important;
    padding: 10px 12px !important;
    background: var(--mapex-popup-bg) !important;
    background-color: var(--mapex-popup-bg) !important;
    transition: background-color 0.15s ease !important;

    // NO hover effect for Edit and View Details
    &:hover,
    &:active,
    &.q-item--active,
    &:focus,
    &.q-focusable:focus,
    &.q-manual-focusable--focused {
      background: var(--mapex-popup-bg) !important;
      background-color: var(--mapex-popup-bg) !important;
    }

    .q-item__section--side {
      min-width: 32px;
      padding-right: 12px;
    }
  }

  // DELETE keeps the red hover effect
  .data-row-action-item--danger {
    background: var(--mapex-popup-bg) !important;
    background-color: var(--mapex-popup-bg) !important;

    &:hover {
      background: var(--mapex-danger-hover) !important;
      background-color: var(--mapex-danger-hover) !important;
    }

    &:active {
      background: var(--mapex-danger-active) !important;
      background-color: var(--mapex-danger-active) !important;
    }

    &.q-item--active,
    &:focus {
      background: var(--mapex-popup-bg) !important;
      background-color: var(--mapex-popup-bg) !important;
    }
  }

  // Separator
  .action-separator {
    margin: 6px 12px !important;
    background-color: var(--mapex-divider) !important;
  }
}
</style>

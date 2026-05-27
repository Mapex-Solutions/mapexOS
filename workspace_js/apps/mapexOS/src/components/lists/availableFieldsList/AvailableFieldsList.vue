<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { AvailableFieldsListProps, AvailableFieldsListEmits } from './interfaces';

defineOptions({
  name: 'AvailableFieldsList'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */

/** COMPOSABLES */

/** UTILS */

/** SERVICES */

/** STORES */

/** LOCAL IMPORTS (constants and handlers ONLY - NO types here!) */
import { DEFAULT_MAX_HEIGHT, SCROLL_THRESHOLD } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<AvailableFieldsListProps>(), {
  maxHeight: DEFAULT_MAX_HEIGHT,
  loading: false,
});

const emit = defineEmits<AvailableFieldsListEmits>();

/** COMPOSABLES & STORES */

/** STATE */

/** COMPUTED */

/**
 * Determines if scroll area should be enabled based on number of fields
 * @returns {boolean} True if fields count exceeds threshold
 */
const shouldScroll = computed(() => {
  return props.fields.length > SCROLL_THRESHOLD;
});

/**
 * Computes the scroll area height style
 * @returns {string} CSS height value or 'auto'
 */
const scrollAreaHeight = computed(() => {
  return shouldScroll.value ? `${props.maxHeight}px` : 'auto';
});

/** WATCHERS */

/** FUNCTIONS */

/**
 * Handle field item click event
 * Emits the field path to parent component
 * @param {string} field - Field path that was clicked
 * @returns {void}
 */
function handleFieldClick(field: string): void {
  emit('field-click', field);
}

/** LIFECYCLE HOOKS */
</script>

<template>
  <div class="available_fields-list">
    <!-- Loading State -->
    <div v-if="loading" class="text-center q-py-lg">
      <q-spinner color="primary" size="md" />
      <div class="text-caption text-grey-7 q-mt-sm">Loading fields...</div>
    </div>

    <!-- Empty State -->
    <div v-else-if="fields.length === 0" class="text-center text-grey-6 q-py-lg">
      <q-icon name="info" size="lg" class="q-mb-sm" />
      <div class="text-body2">No fields available</div>
      <div class="text-caption q-mt-xs">
        Execute a successful test to extract fields from the payload
      </div>
    </div>

    <!-- Fields List -->
    <!-- Short lists render flat; only wrap in a scroll-area when the list
         exceeds the threshold. q-scroll-area collapses to zero with height
         "auto", so we cannot use it unconditionally. -->
    <template v-else>
      <q-scroll-area
        v-if="shouldScroll"
        :style="{ height: scrollAreaHeight }"
        class="fields-scroll-area"
      >
        <q-list bordered separator>
          <q-item
            v-for="(field, index) in fields"
            :key="field"
            clickable
            dense
            @click="handleFieldClick(field)"
          >
            <q-item-section avatar>
              <q-avatar size="24px" color="grey-3" text-color="grey-8">
                {{ index + 1 }}
              </q-avatar>
            </q-item-section>

            <q-item-section>
              <q-item-label class="field-path-text">
                {{ field }}
              </q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-icon name="code" size="xs" color="grey-6" />
            </q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>

      <q-list v-else bordered separator class="fields-scroll-area">
        <q-item
          v-for="(field, index) in fields"
          :key="field"
          clickable
          dense
          @click="handleFieldClick(field)"
        >
          <q-item-section avatar>
            <q-avatar size="24px" color="grey-3" text-color="grey-8">
              {{ index + 1 }}
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label class="field-path-text">
              {{ field }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="code" size="xs" color="grey-6" />
          </q-item-section>
        </q-item>
      </q-list>
    </template>
  </div>
</template>

<style scoped lang="scss">
.available_fields-list {
  .field-path-text {
    font-family: 'Roboto Mono', 'Courier New', monospace;
    font-size: 0.875rem;
    color: var(--mapex-text-primary);
  }

  .fields-scroll-area {
    border-radius: var(--mapex-radius-xs);

    :deep(.q-item) {
      transition: background-color 0.2s ease;

      &:hover {
        background-color: var(--mapex-surface-highlight);
      }
    }
  }
}
</style>

<script setup lang="ts">
defineOptions({
  name: 'Step1Category'
});

/** TYPE IMPORTS */
import type { Trigger, CategoryOption, TriggerCategory } from '../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPONENTS */

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** UTILS */

/** SERVICES */

/** STORES */

/** LOCAL IMPORTS */
import { CATEGORY_OPTIONS } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Trigger];
  'category-selected': [category: TriggerCategory];
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditTriggerTranslations();

/** STATE */

/**
 * Currently selected category option
 */
const selectedCategory = ref(null as CategoryOption | null);

/** COMPUTED */

/** WATCHERS */

/**
 * Watch for category selection changes and emit updates
 */
watch(selectedCategory, (newCategory) => {
  if (newCategory) {
    const updatedTrigger: Trigger = {
      ...props.modelValue,
      category: newCategory.value,
    };
    emit('update:modelValue', updatedTrigger);
    emit('category-selected', newCategory.value);
  }
});

/** FUNCTIONS */

/**
 * Handle category card click
 * @param {CategoryOption} category - Selected category option
 * @returns {void}
 */
function selectCategory(category: CategoryOption): void {
  selectedCategory.value = category;
}

/** LIFECYCLE HOOKS */
</script>

<template>
  <div class="step1-category">
    <div class="text-body1 text-grey-8 q-mb-lg">
      {{ t.steps.step1.intro.value }}
    </div>

    <div class="row q-col-gutter-md">
      <div
        v-for="category in CATEGORY_OPTIONS"
        :key="category.value"
        class="col-12 col-md-6"
      >
        <q-card
          flat
          bordered
          :class="[
            'category-card cursor-pointer transition-all',
            selectedCategory?.value === category.value ? 'selected' : ''
          ]"
          @click="selectCategory(category)"
        >
          <q-card-section class="q-pa-lg">
            <div class="row items-start no-wrap">
              <div class="col">
                <div class="text-center">
                  <div class="category-icon q-mb-md">
                    <q-icon :name="category.icon" size="4rem" color="primary" />
                  </div>
                  <div class="text-h6 text-weight-medium text-dark q-mb-sm">
                    {{ category.emoji }} {{ category.label }}
                  </div>
                  <div class="text-body2 text-grey-7">{{ category.description }}</div>
                </div>
              </div>
              <div v-if="selectedCategory?.value === category.value" class="col-auto">
                <q-icon name="check_circle" color="primary" size="28px" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Info box -->
    <q-banner rounded class="bg-blue-1 text-blue-9 q-mt-lg">
      <template v-slot:avatar>
        <q-icon name="info" color="blue-7" />
      </template>
      <div class="text-body2">
        <strong>{{ t.steps.step1.tip.prefix.value }}:</strong> {{ t.steps.step1.tip.text.value }}
      </div>
    </q-banner>
  </div>
</template>

<style lang="scss" scoped>
.step1-category {
  .category-card {
    border: 2px solid transparent;
    transition: var(--mapex-transition-slow);

    &:hover {
      border-color: var(--mapex-active-border);
      transform: translateY(-2px);
      box-shadow: var(--mapex-shadow-md);
    }

    &.selected {
      border-color: var(--mapex-active-border);
      background-color: var(--mapex-active-bg);
    }
  }

  .category-icon {
    line-height: 1;
  }
}
</style>

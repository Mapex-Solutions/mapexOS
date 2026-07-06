<script setup lang="ts">
defineOptions({
  name: 'Step1Type'
});

/** TYPE IMPORTS */
import type { AssetFormData, AssetTypeOption } from '../../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ASSET_TYPE_OPTIONS } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: AssetFormData;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: AssetFormData];
}>();

/** COMPOSABLES & STORES */
const ts = useTS({ capitalize: true });
const tsRaw = useTS({ capitalize: false });
const bp = 'pages.assets.addAsset.steps.stepType';

/** COMPUTED */
const options = computed<AssetTypeOption[]>(() => ASSET_TYPE_OPTIONS);

/** FUNCTIONS */

/**
 * Whether the given type card matches the current form selection. For the two
 * LoRaWAN cards the kind must also match so Sensor and Gateway are distinct.
 * @param {AssetTypeOption} option - The type card to test.
 * @returns {boolean} True when the card is the active selection.
 */
function isSelected(option: AssetTypeOption): boolean {
  if (props.modelValue.protocol !== option.protocol) {
    return false;
  }
  if (option.protocol === 'LORAWAN') {
    return props.modelValue.lorawanConfig?.kind === option.kind;
  }
  return true;
}

/**
 * Applies a type card: sets the protocol and, for the LoRaWAN cards, the kind.
 * This single up-front choice branches the rest of the wizard.
 * @param {AssetTypeOption} option - The selected type card.
 */
function selectType(option: AssetTypeOption): void {
  const updated: AssetFormData = { ...props.modelValue, protocol: option.protocol };
  if (option.kind) {
    updated.lorawanConfig = { ...props.modelValue.lorawanConfig, kind: option.kind };
  }
  emit('update:modelValue', updated);
}
</script>

<template>
  <div class="step1-type">
    <div class="step1-type__intro">{{ tsRaw(`${bp}.subtitle`) }}</div>

    <div class="row q-col-gutter-md">
      <div v-for="option in options" :key="option.value" class="col-12 col-md-6">
        <q-card
          flat
          bordered
          :class="['type-card cursor-pointer transition-all', isSelected(option) ? 'selected' : '']"
          @click="selectType(option)"
        >
          <q-card-section class="q-pa-lg">
            <div class="row items-start no-wrap">
              <div class="col">
                <div class="text-center">
                  <div class="type-icon q-mb-md">
                    <q-icon :name="option.icon" size="4rem" color="primary" />
                  </div>
                  <div class="text-h6 text-weight-medium q-mb-sm">{{ ts(`${bp}.cards.${option.value}.label`) }}</div>
                  <div class="text-body2 text-grey-7">{{ tsRaw(`${bp}.cards.${option.value}.description`) }}</div>
                </div>
              </div>
              <div v-if="isSelected(option)" class="col-auto">
                <q-icon name="check_circle" color="primary" size="28px" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.step1-type {
  &__intro {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-secondary);
    margin-bottom: var(--mapex-spacing-md);
  }

  .type-card {
    border: 2px solid transparent;
    transition: var(--mapex-transition-slow);

    &:hover {
      border-color: var(--mapex-active-border);
      box-shadow: var(--mapex-shadow-md);
    }

    &.selected {
      border-color: var(--mapex-active-border);
      background-color: var(--mapex-active-bg);
    }
  }
}
</style>

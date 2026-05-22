<script setup lang="ts">
defineOptions({
  name: 'DateRangeFilter'
});

import type { DateRangeFilterProps, DateRangeFilterEmits } from './interfaces';

import { computed } from 'vue';

const props = withDefaults(defineProps<DateRangeFilterProps>(), {
  clearable: true,
  disabled: false,
});

const emit = defineEmits<DateRangeFilterEmits>();

const value = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

const fromDate = computed({
  get: () => value.value?.from ?? null,
  set: (val) => {
    const newValue: { from?: string; to?: string } = { ...value.value };
    if (val) {
      newValue.from = val;
    } else {
      delete newValue.from;
    }
    emit('update:modelValue', newValue);
  },
});

const toDate = computed({
  get: () => value.value?.to ?? null,
  set: (val) => {
    const newValue: { from?: string; to?: string } = { ...value.value };
    if (val) {
      newValue.to = val;
    } else {
      delete newValue.to;
    }
    emit('update:modelValue', newValue);
  },
});

function handleEnter() {
  emit('enter');
}
</script>

<template>
  <div class="row q-col-gutter-sm">
    <div class="col-6">
      <q-input
        v-model="fromDate"
        outlined
        dense
        class="rounded-borders"
        :label="`${label} (From)`"
        :clearable="clearable"
        :disable="disabled"
        mask="####-##-## ##:##:##"
        placeholder="YYYY-MM-DD HH:mm:ss"
        @keyup.enter="handleEnter"
      >
        <template v-slot:prepend>
          <q-icon color="primary" :name="icon" />
        </template>
        <template v-slot:append>
          <q-icon name="event" class="cursor-pointer">
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-date v-model="fromDate" mask="YYYY-MM-DD HH:mm:ss">
                <div class="row items-center justify-end">
                  <q-btn v-close-popup label="Close" color="primary" flat />
                </div>
              </q-date>
            </q-popup-proxy>
          </q-icon>
          <q-icon name="access_time" class="cursor-pointer q-ml-xs">
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-time v-model="fromDate" mask="YYYY-MM-DD HH:mm:ss" with-seconds>
                <div class="row items-center justify-end">
                  <q-btn v-close-popup label="Close" color="primary" flat />
                </div>
              </q-time>
            </q-popup-proxy>
          </q-icon>
        </template>
      </q-input>
    </div>
    <div class="col-6">
      <q-input
        v-model="toDate"
        outlined
        dense
        class="rounded-borders"
        :label="`${label} (To)`"
        :clearable="clearable"
        :disable="disabled"
        mask="####-##-## ##:##:##"
        placeholder="YYYY-MM-DD HH:mm:ss"
        @keyup.enter="handleEnter"
      >
        <template v-slot:prepend>
          <q-icon color="primary" :name="icon" />
        </template>
        <template v-slot:append>
          <q-icon name="event" class="cursor-pointer">
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-date v-model="toDate" mask="YYYY-MM-DD HH:mm:ss">
                <div class="row items-center justify-end">
                  <q-btn v-close-popup label="Close" color="primary" flat />
                </div>
              </q-date>
            </q-popup-proxy>
          </q-icon>
          <q-icon name="access_time" class="cursor-pointer q-ml-xs">
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-time v-model="toDate" mask="YYYY-MM-DD HH:mm:ss" with-seconds>
                <div class="row items-center justify-end">
                  <q-btn v-close-popup label="Close" color="primary" flat />
                </div>
              </q-time>
            </q-popup-proxy>
          </q-icon>
        </template>
      </q-input>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>

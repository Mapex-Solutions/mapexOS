<script setup lang="ts">
defineOptions({
  name: 'FormCard'
});

import type { FormCardProps, FormCardEmits } from './interfaces';
import { computed } from 'vue';

const props = withDefaults(defineProps<FormCardProps>(), {
  navigation: () => ({
    currentStep: 1,
    totalSteps: 1,
    showPreviousButton: true,
    showNextButton: true,
    showSaveButton: true,
  }),
  buttonLabels: () => ({
    previous: 'Previous',
    next: 'Next',
    save: 'Save',
  }),
});

const emit = defineEmits<FormCardEmits>();

const handlePreviousStep = () => {
  emit('previous', props.navigation.currentStep - 1);
};

const handleNextStep = () => {
  emit('next', props.navigation.currentStep + 1);
};

const handleSave = () => {
  emit('save');
};

const shouldShowPreviousButton = computed(() => {
  return props.navigation?.showPreviousButton && (props.navigation?.currentStep || 1) > 1;
});

const shouldShowNextButton = computed(() => {
  return props.navigation?.showNextButton && (props.navigation?.currentStep || 1) < (props.navigation?.totalSteps || 1);
});

const shouldShowSaveButton = computed(() => {
  return props.navigation?.showSaveButton && (props.navigation?.currentStep || 1) >= (props.navigation?.totalSteps || 1);
});
</script>

<template>
  <div class="col-12 col-md-8">
    <q-card class="form-card rounded-borders">

      <!-- FORM CARD HEADER -->
      <q-card-section class="bg-grey-1 q-pb-md">
        <div class="text-h6 text-weight-bold text-primary">
          <q-icon
              size="sm"
              class="q-mr-xs"
              :name="header.icon"
              :color="header.iconColor || 'primary'"
          />
          {{ header.title }}
        </div>
        <div class="text-caption text-grey-7">{{ header.description }}</div>
      </q-card-section>

      <q-card-section class="q-pa-lg">

        <!-- GENERIC SLOT TO PUT ANYTHING HERE, FORM AS AN EXAMPLE -->
        <slot name="form"/>

        <div><q-separator class="q-mt-md" /></div>

        <!-- NAVIGATION BUTTONS -->
        <div class="row justify-between q-mt-lg">
          <q-btn
              v-if="shouldShowPreviousButton"
              flat
              class="rounded-borders"
              color="grey-7"
              icon="arrow_back"
              data-testid="wizard-previous-btn"
              :label="buttonLabels?.previous || 'Previous'"
              @click="handlePreviousStep"
          />
          <div>
            <q-btn
                v-if="shouldShowNextButton"
                class="rounded-borders"
                color="primary"
                icon-right="arrow_forward"
                data-testid="wizard-next-btn"
                :label="buttonLabels?.next || 'Next'"
                :disable="navigation?.disableNextButton || false"
                @click="handleNextStep"
            />
            <q-btn
                v-if="shouldShowSaveButton"
                v-bind="saveButtonId ? { id: saveButtonId } : {}"
                class="rounded-borders"
                color="primary"
                icon-right="save"
                data-testid="wizard-save-btn"
                :label="buttonLabels?.save || 'Save'"
                :disable="navigation?.disableSaveButton || false"
                :loading="navigation?.loadingSaveButton || false"
                @click="handleSave"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<style scoped>
.form-card {
  height: 100%;
}
</style>
<script setup lang="ts">
defineOptions({
  name: 'PageHeader'
});

/** TYPE IMPORTS */
import type { PageHeaderProps } from './interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { InfoModal } from '@components/dialogs/infoModal';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** PROPS & EMITS */
const props = defineProps<PageHeaderProps>();
const emit = defineEmits<{
  (e: 'start-tour'): void;
}>();

/** COMPOSABLES & STORES */
const ts = useTS({ capitalize: true });

/** STATE */
const infoModalOpen = ref(false);

/** COMPUTED */

/**
 * Info button tooltip text (i18n)
 */
const infoTooltipText = computed(() =>
  ts('components.headers.pageHeader.infoTooltip')
);

/** FUNCTIONS */

/**
 * Open the info modal
 */
function openInfoModal(): void {
  infoModalOpen.value = true;
}

/**
 * Emit start-tour event for page tour
 */
function handleStartTour(): void {
  emit('start-tour');
}
</script>

<template>
  <!-- Main header row -->
  <div class="row items-center q-mb-lg">

    <!-- Left: Icon, Title, Description -->
    <div class="col-12 col-md-6">
      <div class="row items-center no-wrap">

        <!-- Optional header icon -->
        <div v-if="props.icon" class="col-auto q-mr-md">
          <q-icon size="xl" :name="props.icon" :color="props.iconColor || 'primary'" />
        </div>
        <div class="col">

          <!-- Main header title with info/tour button -->
          <div class="row items-center q-gutter-sm">
            <div class="text-h4 text-weight-bold text-primary">
              {{ props.title }}
            </div>

            <!-- Tour button (if tour config provided) -->
            <q-btn
              v-if="props.tour?.enabled"
              id="tour-start-btn"
              flat
              dense
              size="sm"
              icon="play_circle_outline"
              label="Tour"
              color="primary"
              class="tour-button"
              no-caps
              @click="handleStartTour"
            />

            <!-- Info button (if info config provided and no tour) -->
            <q-btn
              v-else-if="props.info"
              flat
              round
              dense
              size="sm"
              icon="info"
              color="grey-6"
              class="info-button"
              @click="openInfoModal"
            >
              <AppTooltip :content="infoTooltipText" />
            </q-btn>
          </div>

          <!-- Optional description under the title -->
          <div v-if="props.description" class="text-subtitle1 text-grey-7">
            {{ props.description }}
          </div>

          <!-- Slot for extra content in the header (optional) -->
          <slot name="header-extra"/>
        </div>
      </div>
    </div>
    <!-- Right: Button and extra actions -->
    <div class="col-12 col-md-6 q-mt-sm q-mt-md-none">
      <div class="row justify-end items-center no-wrap">
        <!-- Action button: supports routing or custom handler -->
        <q-btn
            v-if="props.button"
            :id="props.button.id"
            :flat="props.button.flat ?? false"
            :rounded="props.button.rounded ?? true"
            :ripple="props.button.ripple ?? false"
            :unelevated="props.button.unelevated ?? true"
            class="q-px-md"
            :label="props.button.label"
            :icon="props.button.icon"
            :color="props.button.color || 'grey-7'"
            :to="props.button.to"
            @click="props.button.onClick"
        />
        <!-- Slot for additional actions on the right (optional) -->
        <slot name="actions"/>
      </div>
    </div>
  </div>

  <!-- Info Modal -->
  <InfoModal
    v-if="props.info"
    v-model="infoModalOpen"
    :icon="props.icon"
    :title="props.info.title"
    :description="props.info.description"
    :items="props.info.items"
    :docs-url="props.info.docsUrl"
    :docs-label="props.info.docsLabel"
  />
</template>

<style lang="scss" scoped>
.info-button {
  opacity: 0.7;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 1;
  }
}

.tour-button {
  animation: gentle-bounce 2s ease-in-out infinite;

  &:hover {
    animation: none;
  }
}

@keyframes gentle-bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-3px);
  }
}
</style>

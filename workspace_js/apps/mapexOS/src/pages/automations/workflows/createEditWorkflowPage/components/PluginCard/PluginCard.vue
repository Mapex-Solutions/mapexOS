<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginCardProps, PluginCardEmits } from './interfaces/PluginCard.interface';

/** PROPS & EMITS */
defineProps<PluginCardProps>();
const emit = defineEmits<PluginCardEmits>();
</script>

<template>
  <q-card flat bordered class="plugin-card">
    <q-card-section>
      <!-- Header: icon + name + category -->
      <div class="row items-center no-wrap q-mb-sm">
        <q-avatar v-if="brandIconUrl" size="40px" class="q-mr-sm" square>
          <img :src="brandIconUrl" :alt="name" />
        </q-avatar>
        <q-icon v-else :name="icon" size="32px" class="q-mr-sm" :color="color === '#666' ? 'grey-6' : undefined" :style="color !== '#666' ? { color } : undefined" />
        <div class="col">
          <div class="text-subtitle1 text-weight-medium">{{ name }}</div>
          <div class="text-caption text-grey-6">
            {{ author }} &middot; v{{ version }}
          </div>
        </div>
        <q-badge
          outline
          class="text-caption"
          :style="{ color: color, borderColor: color }"
        >
          {{ categoryLabel }}
        </q-badge>
      </div>

      <!-- Description -->
      <div class="text-body2 text-grey-7 q-mb-md plugin-card__description">
        {{ description }}
      </div>

      <!-- Tags -->
      <div v-if="tags.length > 0" class="row q-gutter-xs q-mb-md">
        <q-badge
          v-for="tag in tags"
          :key="tag"
          outline
          color="grey-6"
          class="text-caption"
        >
          {{ tag }}
        </q-badge>
      </div>

      <!-- Footer: node count + actions -->
      <div class="row items-center justify-between">
        <span class="text-caption text-grey-6">
          <q-icon name="account_tree" size="xs" class="q-mr-xs" />
          {{ nodeCount }} {{ nodeCount === 1 ? 'node' : 'nodes' }}
        </span>

        <div class="row q-gutter-xs items-center">
          <!-- Details button -->
          <q-btn
            flat
            dense
            no-caps
            color="grey-7"
            icon="info_outline"
            size="sm"
            :label="detailsLabel"
            @click.stop="emit('details')"
          />

          <!-- Installed badge -->
          <q-btn
            v-if="installed"
            flat
            dense
            no-caps
            color="positive"
            icon="check_circle"
            size="sm"
            disable
            :label="installedLabel"
          />

          <!-- Install button -->
          <q-btn
            v-else
            unelevated
            dense
            no-caps
            color="primary"
            icon="download"
            size="sm"
            :label="installing ? installingLabel : installLabel"
            :loading="installing"
            :disable="installDisabled"
            @click="emit('install')"
          />
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<style lang="scss" scoped>
.plugin-card {
  background: var(--mapex-card-bg);
  border-color: var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-fast);

  &:hover {
    box-shadow: var(--mapex-shadow-sm);
    border-color: var(--mapex-card-hover-border);
  }

  &__description {
    min-height: 40px;
  }
}
</style>

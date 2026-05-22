<template>
  <q-card
      flat
      bordered
      class="asset-row q-card--hover"
  >
    <!-- DESKTOP/TABLET VERSION (>= 600px) -->
    <q-card-section class="asset-row__desktop row items-center q-pa-md">

      <!-- ICON + NAME + DESCRIPTION -->
      <div class="asset-row__main">
        <div class="row items-center q-gutter-x-sm">
          <q-avatar
              size="40px"
              :color="getIconColor()"
              text-color="white"
              :icon="props.asset.icon || 'sensors'"
          />
          <div class="column" style="min-width: 0; flex: 1;">
            <div class="text-subtitle2 text-weight-medium text-grey-9 ellipsis">
              {{ props.asset.name }}
            </div>
            <div v-if="props.asset.description" class="text-caption text-grey-6 ellipsis">
              {{ props.asset.description }}
            </div>
          </div>
        </div>
      </div>

      <!-- UUID (desktop/tablet) -->
      <div class="asset-row__uuid">
        <div class="text-caption text-grey-5 text-weight-medium">UUID</div>
        <code class="text-body2 text-weight-medium text-grey-8 ellipsis">
          {{ props.asset.assetUUID }}
        </code>
      </div>

      <!-- ASSET TYPE -->
      <div class="asset-row__type">
        <div class="text-caption text-grey-5 text-weight-medium">Type</div>
        <DetailChip
            :label="props.asset.assetTypeName || 'Unknown'"
            :color="getTypeColorForChip()"
            size="sm"
            outline
        />
      </div>

      <!-- PROTOCOL -->
      <div class="asset-row__protocol">
        <div class="text-caption text-grey-5 text-weight-medium">Protocol</div>
        <DetailChip
            :label="props.asset.protocol?.type?.toUpperCase() || 'N/A'"
            :color="getProtocolColorForChip()"
            size="sm"
            outline
        />
      </div>

      <!-- STATUS -->
      <div class="asset-row__status">
        <div class="text-caption text-grey-5 text-weight-medium">Status</div>
        <DetailChip
            :label="props.asset.enabled ? 'ACTIVE' : 'INACTIVE'"
            :color="props.asset.enabled ? 'green' : 'red'"
            size="sm"
            outline
        />
      </div>

      <!-- ACTIONS -->
      <div class="asset-row__actions">
        <q-btn
            flat
            dense
            round
            icon="more_vert"
            @click.stop
        >
          <q-menu>
            <q-list style="min-width: 150px">
              <q-item v-close-popup clickable @click.stop="emit('edit', props.asset)">
                <q-item-section avatar>
                  <q-icon name="edit" color="primary"/>
                </q-item-section>
                <q-item-section>Edit</q-item-section>
              </q-item>

              <q-item v-close-popup clickable @click.stop="emit('view', props.asset)">
                <q-item-section avatar>
                  <q-icon name="visibility" color="info"/>
                </q-item-section>
                <q-item-section>View Details</q-item-section>
              </q-item>

              <q-separator />

              <q-item v-close-popup clickable @click.stop="emit('delete', props.asset)">
                <q-item-section avatar>
                  <q-icon name="delete" color="negative"/>
                </q-item-section>
                <q-item-section>Delete</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>

    </q-card-section>

    <!-- MOBILE VERSION (< 600px) -->
    <q-card-section class="asset-row__mobile q-pa-sm">

      <!-- Main Row: Icon + Name + UUID + Expand + Actions -->
      <div class="row items-center q-gutter-x-sm" @click="toggleExpand">
        <q-avatar
            size="40px"
            :color="getIconColor()"
            text-color="white"
            :icon="props.asset.icon || 'sensors'"
        />

        <div class="col column" style="min-width: 0;">
          <div class="text-subtitle2 text-weight-medium text-grey-9 ellipsis">
            {{ props.asset.name }}
          </div>
          <code class="text-caption text-grey-6 ellipsis">
            {{ props.asset.assetUUID }}
          </code>
        </div>

        <q-btn
            flat
            dense
            round
            size="sm"
            :icon="expanded ? 'expand_less' : 'expand_more'"
            color="grey-7"
            @click.stop="toggleExpand"
        />

        <q-btn
            flat
            dense
            round
            icon="more_vert"
            size="sm"
            color="grey-7"
            @click.stop
        >
          <q-menu>
            <q-list style="min-width: 150px">
              <q-item v-close-popup clickable @click.stop="emit('edit', props.asset)">
                <q-item-section avatar>
                  <q-icon name="edit" color="primary"/>
                </q-item-section>
                <q-item-section>Edit</q-item-section>
              </q-item>

              <q-item v-close-popup clickable @click.stop="emit('view', props.asset)">
                <q-item-section avatar>
                  <q-icon name="visibility" color="info"/>
                </q-item-section>
                <q-item-section>View Details</q-item-section>
              </q-item>

              <q-separator />

              <q-item v-close-popup clickable @click.stop="emit('delete', props.asset)">
                <q-item-section avatar>
                  <q-icon name="delete" color="negative"/>
                </q-item-section>
                <q-item-section>Delete</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>

      <!-- Expanded Details -->
      <q-slide-transition>
        <div v-show="expanded" class="q-mt-sm q-pt-sm" :style="{ borderTop: '1px solid var(--mapex-card-border)' }">
          <div v-if="props.asset.description" class="text-caption text-grey-7 q-mb-sm">
            {{ props.asset.description }}
          </div>

          <div class="row q-gutter-sm">
            <DetailChip
                :label="`Type: ${props.asset.assetTypeName || 'Unknown'}`"
                :color="getTypeColorForChip()"
                size="sm"
                outline
            />
            <DetailChip
                :label="`Protocol: ${props.asset.protocol?.type?.toUpperCase() || 'N/A'}`"
                :color="getProtocolColorForChip()"
                size="sm"
                outline
            />
            <DetailChip
                :label="props.asset.enabled ? 'ACTIVE' : 'INACTIVE'"
                :color="props.asset.enabled ? 'green' : 'red'"
                size="sm"
                outline
            />
          </div>
        </div>
      </q-slide-transition>

    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
defineOptions({
  name: 'AssetRow'
});

import { ref } from 'vue';
import { DetailChip } from '@components/chips/DetailChip';
import type { DetailChipColor } from '@components/chips/DetailChip';
import type { AssetRowProps } from './interfaces';

const props = defineProps<{
  asset: AssetRowProps
}>();

const emit = defineEmits<{
  click: [asset: AssetRowProps]
  edit: [asset: AssetRowProps]
  view: [asset: AssetRowProps]
  delete: [asset: AssetRowProps]
}>();

const expanded = ref(false);

const toggleExpand = () => {
  expanded.value = !expanded.value;
  if (expanded.value) {
    emit('click', props.asset);
  }
};

const getIconColor = () => {
  return props.asset.enabled ? 'primary' : 'grey-5';
};

const getTypeColorForChip = (): DetailChipColor => {
  return props.asset.enabled ? 'blue' : 'grey';
};

const getProtocolColorForChip = (): DetailChipColor => {
  const protocolColors: Record<string, DetailChipColor> = {
    'mqtt': 'purple',
    'http': 'blue',
    'lorawan': 'orange',
  };
  return protocolColors[props.asset.protocol?.type?.toLowerCase() || ''] || 'grey';
};

</script>

<style scoped lang="scss">
.asset-row {
  transition: var(--mapex-transition-base);

  &:hover {
    background-color: var(--mapex-surface-highlight);
    box-shadow: var(--mapex-shadow-sm);
  }
}

// Desktop/Tablet Version (>= 600px)
.asset-row__desktop {
  display: flex;
  min-height: 70px;
  gap: 12px;
  align-items: center;
  justify-content: space-between;

  // Main section (Icon + Name + Description)
  .asset-row__main {
    flex: 1 1 300px;
    min-width: 0;
    max-width: 400px;
  }

  // UUID section
  .asset-row__uuid {
    flex: 0 0 200px;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;

    code {
      display: block;
      font-size: 0.75rem;
      font-family: 'Courier New', monospace;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  // Type section
  .asset-row__type {
    flex: 0 0 150px;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  // Protocol section
  .asset-row__protocol {
    flex: 0 0 140px;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  // Status section
  .asset-row__status {
    flex: 0 0 110px;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  // Actions section
  .asset-row__actions {
    flex: 0 0 auto;
  }

  @media (max-width: 1024px) {
    // Hide UUID on tablet
    .asset-row__uuid {
      display: none;
    }

    .asset-row__main {
      flex: 1 1 auto;
    }
  }

  @media (max-width: 600px) {
    // Hide desktop version on mobile
    display: none;
  }
}

// Mobile Version (< 600px)
.asset-row__mobile {
  display: none;

  @media (max-width: 600px) {
    display: block;
  }
}

.q-avatar {
  box-shadow: var(--mapex-shadow-sm);
  flex-shrink: 0;
}
</style>

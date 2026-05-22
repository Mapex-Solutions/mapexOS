<script setup lang="ts">
defineOptions({
  name: 'RouterCard'
});

/** TYPE IMPORTS */
import type { RouterCardProps, RouterCardEmits } from './interfaces/RouterCard.interface';
import type { LakeHouseItem } from '@components/drawers';
import type { RouterFormState } from '../../interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { ConditionalRoutingToggle } from '../ConditionalRoutingToggle';
import { MatchConfiguration } from '../MatchConfiguration';
import { WorkflowConfig } from '../WorkflowConfig';
import {
  LakeHouseSelectorDrawer,
  TriggerSelectorDrawer,
} from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** PROPS & EMITS */
const props = defineProps<RouterCardProps>();
const emit = defineEmits<RouterCardEmits>();

/** STATE - Drawer visibility */
const showLakeHouseDrawer = ref(false);
const showNotificationDrawer = ref(false);

/** COMPUTED */

/**
 * Get the selected data lake ID from router state
 */
const selectedLakeHouseId = computed(() => {
  return props.router.lakeHouse?.lakeHouseId || null;
});

/**
 * Get the selected notification ID from router state
 */
const selectedNotificationId = computed(() => {
  return props.router.notification?.notificationId || null;
});

/** FUNCTIONS */

/**
 * Get router kind option details
 * @param {string} kind - Router kind value
 * @returns {object | undefined} Router kind option
 */
function getRouterKindOption(kind: string) {
  return props.routerKindOptions.find(opt => opt.value === kind);
}

/**
 * Update router fields
 * @param {K} field - Field key to update
 * @param {RouterFormState[K]} value - New value
 */
function updateRouterField<K extends keyof RouterFormState>(field: K, value: RouterFormState[K]) {
  emit('update:router', { ...props.router, [field]: value });
}

/**
 * Handle data lake selection from drawer
 * @param {LakeHouseItem} lakeHouse - Selected data lake
 */
function handleLakeHouseSelect(lakeHouse: LakeHouseItem): void {
  if (lakeHouse.id) {
    emit('update:router', {
      ...props.router,
      lakeHouse: {
        ...props.router.lakeHouse,
        lakeHouseId: lakeHouse.id,
      },
      lakeHouseName: lakeHouse.name || 'Unnamed Data Lake',
    });
  }
}

/**
 * Handle notification/trigger selection from drawer
 * @param {any} trigger - Selected trigger
 */
function handleNotificationSelect(trigger: any): void {
  if (trigger.id) {
    emit('update:router', {
      ...props.router,
      notification: {
        ...props.router.notification,
        notificationId: trigger.id,
      },
      notificationName: trigger.name || 'Unnamed Notification',
    });
  }
}

/**
 * Clear data lake selection
 */
function clearLakeHouse(): void {
  const updated: RouterFormState = {
    ...props.router,
    lakeHouse: {
      ...props.router.lakeHouse,
      lakeHouseId: '',
    },
  };
  delete updated.lakeHouseName;
  emit('update:router', updated);
}

/**
 * Clear notification selection
 */
function clearNotification(): void {
  const updated: RouterFormState = {
    ...props.router,
    notification: {
      ...props.router.notification,
      notificationId: '',
    },
  };
  delete updated.notificationName;
  emit('update:router', updated);
}
</script>

<template>
  <q-card
    :id="`route-group-router-card-${index}`"
    bordered
    class="router-card shadow-2"
  >
    <q-card-section class="bg-grey-2 q-py-sm">
      <div class="row items-center">
        <div class="col">
          <div class="text-subtitle2 text-weight-bold text-primary">
            <q-icon
              size="sm"
              class="q-mr-sm"
              :name="getRouterKindOption(router.kind)?.icon || 'route'"
              :color="getRouterKindOption(router.kind)?.color || 'primary'"
            />
            Router #{{ index + 1 }}
          </div>
        </div>
        <div v-if="index > 0" class="col-auto">
          <q-btn
            flat
            dense
            round
            size="sm"
            icon="delete"
            color="negative"
            @click="emit('remove')"
          >
            <AppTooltip :content="t.createEdit.routersStep.removeRouter.value" />
          </q-btn>
        </div>
      </div>
    </q-card-section>

    <q-card-section>
      <div class="row q-col-gutter-md">
        <!-- Router Kind -->
        <div :id="`route-group-router-kind-${index}`" class="col-12">
          <q-select
            outlined
            dense
            emit-value
            map-options
            :label="t.createEdit.routersStep.routerCard.kind.label.value + ' *'"
            :options="routerKindOptions"
            option-label="label"
            option-value="value"
            :model-value="router.kind"
            @update:model-value="(val) => { updateRouterField('kind', val); emit('kind-change', val); }"
          >
            <template #prepend>
              <q-icon name="category" color="primary" />
            </template>
            <template #option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section avatar>
                  <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ scope.opt.label }}</q-item-label>
                  <q-item-label caption>{{ scope.opt.description }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>
        </div>

        <!-- Destination Configuration based on kind -->

        <!-- Data Lake Selector -->
        <div v-if="router.kind === 'lake_house' && router.lakeHouse" class="col-12">
          <div class="text-caption text-grey-7 q-mb-xs">
            {{ t.createEdit.routersStep.routerCard.lakeHouse.label.value }} *
          </div>
          <q-field
            outlined
            dense
            :hint="t.createEdit.routersStep.routerCard.lakeHouse.required.value"
            class="selector-field"
          >
            <template #prepend>
              <q-icon name="storage" color="purple-6" />
            </template>
            <template #control>
              <div
                class="self-center full-width no-outline cursor-pointer"
                @click="showLakeHouseDrawer = true"
              >
                <span v-if="router.lakeHouseName || router.lakeHouse.lakeHouseId" class="text-body2">
                  {{ router.lakeHouseName || router.lakeHouse.lakeHouseId }}
                </span>
                <span v-else class="text-grey-6 text-body2">
                  Click to select a data lake...
                </span>
              </div>
            </template>
            <template #append>
              <q-btn
                v-if="router.lakeHouseName || router.lakeHouse.lakeHouseId"
                flat
                round
                dense
                size="sm"
                icon="close"
                color="grey-7"
                @click.stop="clearLakeHouse"
              >
                <AppTooltip content="Clear selection" />
              </q-btn>
              <q-btn
                flat
                round
                dense
                size="sm"
                icon="search"
                color="primary"
                @click="showLakeHouseDrawer = true"
              >
                <AppTooltip content="Select data lake" />
              </q-btn>
            </template>
          </q-field>
        </div>

        <!-- Notification/Trigger Selector -->
        <div v-if="router.kind === 'notification' && router.notification" class="col-12">
          <div class="text-caption text-grey-7 q-mb-xs">
            {{ t.createEdit.routersStep.routerCard.notification.label.value }} *
          </div>
          <q-field
            outlined
            dense
            :hint="t.createEdit.routersStep.routerCard.notification.required.value"
            class="selector-field"
          >
            <template #prepend>
              <q-icon name="notifications" color="orange-6" />
            </template>
            <template #control>
              <div
                class="self-center full-width no-outline cursor-pointer"
                @click="showNotificationDrawer = true"
              >
                <span v-if="router.notificationName || router.notification.notificationId" class="text-body2">
                  {{ router.notificationName || router.notification.notificationId }}
                </span>
                <span v-else class="text-grey-6 text-body2">
                  Click to select a notification trigger...
                </span>
              </div>
            </template>
            <template #append>
              <q-btn
                v-if="router.notificationName || router.notification.notificationId"
                flat
                round
                dense
                size="sm"
                icon="close"
                color="grey-7"
                @click.stop="clearNotification"
              >
                <AppTooltip content="Clear selection" />
              </q-btn>
              <q-btn
                flat
                round
                dense
                size="sm"
                icon="search"
                color="primary"
                @click="showNotificationDrawer = true"
              >
                <AppTooltip content="Select notification trigger" />
              </q-btn>
            </template>
          </q-field>
        </div>

        <!-- Workflow Config -->
        <WorkflowConfig
          v-if="router.kind === 'workflow' && router.workflow"
          :workflow="router.workflow"
          :t="t"
          @update:workflow="(val) => emit('update:router', { ...router, workflow: val } as RouterFormState)"
        />

        <!-- Conditional Routing Section -->
        <div :id="`route-group-conditional-routing-${index}`" class="col-12">
          <q-separator class="q-my-lg" />

          <!-- Toggle Card -->
          <ConditionalRoutingToggle
            :model-value="router.hasConditionalRouting"
            :t="t"
            @update:model-value="emit('toggle-conditional', $event)"
          />
        </div>

        <!-- Match Configuration -->
        <div v-if="router.hasConditionalRouting && router.match" class="col-12">
          <MatchConfiguration
            :match="router.match"
            :match-policy-options="matchPolicyOptions"
            :match-operator-options="matchOperatorOptions"
            :t="t"
            @update:match="updateRouterField('match', $event)"
            @add-rule="emit('add-match-rule')"
            @remove-rule="(ruleIndex) => emit('remove-match-rule', ruleIndex)"
          />
        </div>
      </div>
    </q-card-section>

    <!-- Selector Drawers -->
    <LakeHouseSelectorDrawer
      v-model="showLakeHouseDrawer"
      :selected-data-lake-id="selectedLakeHouseId"
      @select="handleLakeHouseSelect"
    />

    <TriggerSelectorDrawer
      v-model="showNotificationDrawer"
      :selected-trigger-id="selectedNotificationId"
      @select="handleNotificationSelect"
    />
  </q-card>
</template>

<style scoped>
.router-card {
  border-radius: var(--mapex-radius-md);
}

.selector-field {
  cursor: pointer;
}

.selector-field :deep(.q-field__control) {
  cursor: pointer;
}

.selector-field :deep(.q-field__native) {
  cursor: pointer;
}
</style>

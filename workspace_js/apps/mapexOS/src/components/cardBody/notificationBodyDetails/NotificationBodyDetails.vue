<template>
  <div class="notification-detail">
    <!-- Slack, Teams, Push, Telegram, Webhook - Layout similar to email destinations -->
    <template v-if="['slack', 'teams', 'push', 'telegram', 'webhook'].includes(notification.type)">
      <div class="column q-gutter-sm q-mb-md">
        <q-card flat bordered class="bg-grey-1">
          <q-card-section class="q-pa-sm">
            <div class="row items-center q-mb-xs">
              <q-icon
                size="16px"
                class="q-mr-xs"
                :name="getMainIcon(notification.type)"
                :color="getMainColor(notification.type)"
              />
              <span class="text-caption text-grey-7 text-weight-medium">
                {{ getMainLabel(notification.type) }}
              </span>
            </div>
            <div class="text-body2 text-weight-bold text-grey-9">
              {{ getMainValue(notification) }}
            </div>
          </q-card-section>
        </q-card>

        <q-card flat bordered class="bg-grey-1 standardized-card">
          <q-card-section class="q-pa-sm">
            <div class="row items-center q-mb-xs">
              <q-icon
                size="16px"
                class="q-mr-xs"
                :name="getSecondaryIcon(notification.type)"
                :color="getMainColor(notification.type)"
              />
              <span class="text-caption text-grey-7 text-weight-medium">
                {{ getSecondaryLabel(notification.type) }}
              </span>
            </div>
            <div class="text-body2 text-weight-bold text-grey-9">
              <div class="chips-container">
                <div class="chips-wrapper">
                  <template v-if="notification.type === 'push'">
                    <q-chip
                      size="sm"
                      class="q-ma-xs"
                      :color="getChipColor(notification.type)"
                      :text-color="getChipTextColor(notification.type)"
                      :icon="getChipIcon(notification.type)"
                    >
                      {{ notification.data.deviceCount?.toLocaleString() }}
                      <q-badge color="orange-3" text-color="orange-9" class="q-ml-xs">ativos</q-badge>
                    </q-chip>
                  </template>
                  <template v-else-if="notification.type === 'webhook'">
                    <q-chip
                      size="sm"
                      class="q-ma-xs"
                      :color="getChipColor(notification.type)"
                      :text-color="getChipTextColor(notification.type)"
                      :icon="getChipIcon(notification.type)"
                    >
                      {{ notification.data.method }} {{ notification.data.url }}
                    </q-chip>
                  </template>
                  <template v-else>
                    <q-chip
                      v-for="(item, index) in visibleItems"
                      size="sm"
                      class="q-ma-xs"
                      :key="index"
                      :color="getChipColor(notification.type)"
                      :text-color="getChipTextColor(notification.type)"
                      :icon="getChipIcon(notification.type)"
                    >
                      {{ item }}
                    </q-chip>
                    <q-chip
                      v-if="hasMoreItems"
                      icon="mdi-dots-horizontal"
                      size="sm"
                      color="primary"
                      text-color="white"
                      class="q-ma-xs cursor-pointer"
                      @click="showAllItems = !showAllItems"
                    >
                      <AppTooltip :offset="[10, 10]">
                        <div class="column q-gutter-xs">
                          <div class="text-weight-bold">{{ showAllItems ? 'Ocultar' : 'Ver mais' }}</div>
                          <div v-if="!showAllItems" class="text-caption">
                            <div v-for="item in hiddenItems" :key="item">{{ item }}</div>
                          </div>
                        </div>
                      </AppTooltip>
                      {{ showAllItems ? 'Ocultar' : `+${hiddenItems.length}` }}
                    </q-chip>
                  </template>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <q-card flat bordered :class="getConfigCardClass(notification.type)">
        <q-card-section class="q-pa-sm">
          <div class="text-caption text-weight-bold q-mb-xs" :class="getConfigTitleClass(notification.type)">
            {{ getConfigTitle(notification.type) }}
          </div>
          <div class="text-body2" :class="getConfigTextClass(notification.type)">
            {{ getConfigText(notification) }}
          </div>
        </q-card-section>
      </q-card>
    </template>

    <!-- Email Notification - Layout original melhorado -->
    <template v-else-if="notification.type === 'email'">
      <div class="column q-gutter-sm q-mb-md">
        <q-card flat bordered class="bg-grey-1">
          <q-card-section class="q-pa-sm">
            <div class="row items-center q-mb-xs">
              <q-icon size="16px" class="q-mr-xs" name="mdi-email-send" color="green-6" />
              <span class="text-caption text-grey-7 text-weight-medium">REMETENTE</span>
            </div>
            <div class="text-body2 text-weight-bold text-grey-9">{{ notification.data.from }}</div>
          </q-card-section>
        </q-card>

        <q-card flat bordered class="bg-grey-1 standardized-card">
          <q-card-section class="q-pa-sm">
            <div class="row items-center q-mb-xs">
              <q-icon size="16px" class="q-mr-xs" name="mdi-email-multiple" color="green-6" />
              <span class="text-caption text-grey-7 text-weight-medium">DESTINATÁRIOS</span>
            </div>
            <div class="text-body2 text-weight-bold text-grey-9">
              <div class="chips-container">
                <div class="chips-wrapper">
                  <q-chip
                    v-for="(email, index) in visibleEmails"
                    icon="mdi-email"
                    size="sm"
                    color="green-2"
                    text-color="green-9"
                    class="q-ma-xs"
                    :key="index"
                  >
                    {{ email }}
                  </q-chip>
                  <q-chip
                    v-if="hasMoreEmails"
                    icon="mdi-dots-horizontal"
                    size="sm"
                    color="primary"
                    text-color="white"
                    class="q-ma-xs cursor-pointer"
                    @click="showAllEmails = !showAllEmails"
                  >
                    <AppTooltip :offset="[10, 10]">
                      <div class="column q-gutter-xs">
                        <div class="text-weight-bold">{{ showAllEmails ? 'Ocultar' : 'Ver mais' }}</div>
                        <div v-if="!showAllEmails" class="text-caption">
                          <div v-for="email in hiddenEmails" :key="email">{{ email }}</div>
                        </div>
                      </div>
                    </AppTooltip>
                    {{ showAllEmails ? 'Ocultar' : `+${hiddenEmails.length}` }}
                  </q-chip>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      
      <q-card flat bordered class="bg-green-1 q-mt-sm">
        <q-card-section class="q-pa-sm">
          <div class="text-caption text-green-8 text-weight-bold q-mb-xs">ENTREGA DE EMAIL</div>
          <div class="text-body2 text-green-9">
            Mensagens serão enviadas via email para {{ notification.data.to.length }}
            destinatário{{ notification.data.to.length > 1 ? 's' : '' }}
            configurado{{ notification.data.to.length > 1 ? 's' : '' }}
          </div>
        </q-card-section>
      </q-card>
    </template>
  </div>
</template>

<script lang="ts" setup>
defineOptions({
  name: 'NotificationBodyDetails'
});

import { computed, ref } from 'vue';
import { AppTooltip } from '@components/tooltips';

const props = defineProps<{ notification: any }>();

const showAllItems = ref(false);
const showAllEmails = ref(false);
const maxItemsPerRow = 4; // Maximum items to show before "ver mais"

// Email specific computed properties
const visibleEmails = computed(() => {
  if (!props.notification.data.to) return [];
  return showAllEmails.value ? props.notification.data.to : props.notification.data.to.slice(0, maxItemsPerRow);
});

const hiddenEmails = computed(() => {
  if (!props.notification.data.to) return [];
  return props.notification.data.to.slice(maxItemsPerRow);
});

const hasMoreEmails = computed(() => {
  return props.notification.data.to && props.notification.data.to.length > maxItemsPerRow;
});

// General items (channels) computed properties
const allItems = computed(() => {
  const data = props.notification.data as { channelsName?: string[], chatNames?: string[] };
  return data.channelsName || data.chatNames || [];
});

const visibleItems = computed(() => {
  return showAllItems.value ? allItems.value : allItems.value.slice(0, maxItemsPerRow);
});

const hiddenItems = computed(() => {
  return allItems.value.slice(maxItemsPerRow);
});

const hasMoreItems = computed(() => {
  return allItems.value.length > maxItemsPerRow;
});

// Helper functions for different notification types
const getMainIcon = (type: string) => {
  switch (type) {
    case 'slack': return 'mdi-domain';
    case 'teams': return 'mdi-account-group';
    case 'push': return 'mdi-cellphone';
    case 'telegram': return 'mdi-send';
    case 'webhook': return 'mdi-webhook';
    default: return 'mdi-bell';
  }
};

const getMainColor = (type: string) => {
  switch (type) {
    case 'slack': return 'purple-6';
    case 'teams': return 'blue-6';
    case 'push': return 'orange-6';
    case 'telegram': return 'cyan-6';
    case 'webhook': return 'indigo-6';
    default: return 'grey-6';
  }
};

const getMainLabel = (type: string) => {
  switch (type) {
    case 'slack': return 'WORKSPACE';
    case 'teams': return 'EQUIPE';
    case 'push': return 'APLICATIVO';
    case 'telegram': return 'BOT';
    case 'webhook': return 'ENDPOINT';
    default: return 'TIPO';
  }
};

const getMainValue = (notification: any) => {
  switch (notification.type) {
    case 'slack': return notification.data.workspace;
    case 'teams': return notification.data.teamName;
    case 'push': return notification.data.appName;
    case 'telegram': return notification.data.botName;
    case 'webhook': return notification.data.name;
    default: return '';
  }
};

const getSecondaryIcon = (type: string) => {
  switch (type) {
    case 'slack': return 'mdi-pound';
    case 'teams': return 'mdi-forum';
    case 'push': return 'mdi-devices';
    case 'telegram': return 'mdi-chat';
    case 'webhook': return 'mdi-api';
    default: return 'mdi-bell';
  }
};

const getSecondaryLabel = (type: string) => {
  switch (type) {
    case 'slack': return 'CANAIS';
    case 'teams': return 'CANAIS';
    case 'push': return 'DISPOSITIVOS';
    case 'telegram': return 'CHATS';
    case 'webhook': return 'CONFIGURAÇÃO';
    default: return 'ITENS';
  }
};

const getChipColor = (type: string) => {
  switch (type) {
    case 'slack': return 'purple-2';
    case 'teams': return 'blue-2';
    case 'push': return 'orange-2';
    case 'telegram': return 'cyan-2';
    case 'webhook': return 'indigo-2';
    default: return 'grey-2';
  }
};

const getChipTextColor = (type: string) => {
  switch (type) {
    case 'slack': return 'purple-9';
    case 'teams': return 'blue-9';
    case 'push': return 'orange-9';
    case 'telegram': return 'cyan-9';
    case 'webhook': return 'indigo-9';
    default: return 'grey-9';
  }
};

const getChipIcon = (type: string) => {
  switch (type) {
    case 'slack': return 'mdi-pound';
    case 'teams': return 'mdi-forum';
    case 'push': return 'mdi-devices';
    case 'telegram': return 'mdi-chat';
    case 'webhook': return 'mdi-link';
    default: return 'mdi-tag';
  }
};

const getConfigCardClass = (type: string) => {
  switch (type) {
    case 'slack': return 'bg-purple-1';
    case 'teams': return 'bg-blue-1';
    case 'push': return 'bg-orange-1';
    case 'telegram': return 'bg-cyan-1';
    case 'webhook': return 'bg-indigo-1';
    default: return 'bg-grey-1';
  }
};

const getConfigTitleClass = (type: string) => {
  switch (type) {
    case 'slack': return 'text-purple-8';
    case 'teams': return 'text-blue-8';
    case 'push': return 'text-orange-8';
    case 'telegram': return 'text-cyan-8';
    case 'webhook': return 'text-indigo-8';
    default: return 'text-grey-8';
  }
};

const getConfigTextClass = (type: string) => {
  switch (type) {
    case 'slack': return 'text-purple-9';
    case 'teams': return 'text-blue-9';
    case 'push': return 'text-orange-9';
    case 'telegram': return 'text-cyan-9';
    case 'webhook': return 'text-indigo-9';
    default: return 'text-grey-9';
  }
};

const getConfigTitle = (type: string) => {
  switch (type) {
    case 'slack': return 'CONFIGURAÇÃO';
    case 'teams': return 'INTEGRAÇÃO';
    case 'push': return 'NOTIFICAÇÕES PUSH';
    case 'telegram': return 'TELEGRAM BOT';
    case 'webhook': return 'WEBHOOKS HTTP';
    default: return 'CONFIGURAÇÃO';
  }
};

const getConfigText = (notification: any) => {
  switch (notification.type) {
    case 'slack': 
      return 'Alertas e atualizações em tempo real serão enviados para este canal do Slack';
    case 'teams': 
      return 'Notificações serão publicadas diretamente no(s) canal(is) especificado(s) do Teams';
    case 'push': 
      return `Notificações push móveis serão enviadas para ${notification.data.deviceCount?.toLocaleString()} dispositivos registrados`;
    case 'telegram':
      return `Mensagens serão enviadas via bot do Telegram para ${notification.data.chatNames?.length || 0} chat(s) configurado(s)`;
    case 'webhook':
      return `Requisições HTTP ${notification.data.method} serão enviadas para o endpoint configurado quando eventos ocorrerem`;
    default: 
      return '';
  }
};
</script>

<style scoped>
.notification-detail {
  width: 100%;
}

.chips-container {
  width: 100%;
}

.chips-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  max-height: 60px; /* Near 2 lines */
  overflow: hidden;
}

.chips-wrapper .q-chip {
  margin: 2px;
}

.cursor-pointer {
  cursor: pointer;
}

.cursor-pointer:hover {
  opacity: 0.8;
}

/* Padronização da altura dos cards */
.standardized-card {
  min-height: 100px; /* Altura mínima padronizada */
}

.standardized-card .q-card-section {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.standardized-card .chips-container {
  flex: 1;
  display: flex;
  align-items: flex-start;
}
</style>
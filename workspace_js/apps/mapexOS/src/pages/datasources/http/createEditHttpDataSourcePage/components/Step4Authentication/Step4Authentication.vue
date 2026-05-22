<template>
  <div class="row q-col-gutter-md">
    <!-- Authentication Type Selection -->
    <div class="col-12">
      <div class="row items-center q-mb-sm">
        <q-icon name="lock" color="primary" class="q-mr-xs" />
        <div class="text-subtitle2 text-weight-medium">{{ t.authentication.title.value }}</div>
      </div>
      <q-select
        v-model="localData.authType"
        outlined
        dense
        class="rounded-borders"
        :label="`${t.authentication.authTypeLabel.value} *`"
        :options="AUTH_TYPE_OPTIONS"
        option-label="label"
        option-value="value"
        emit-value
        map-options
        :rules="[(val: any) => !!val || t.authentication.authTypeRequired.value]"
        @update:model-value="handleAuthTypeChange"
      >
        <template #prepend>
          <q-icon name="security" color="primary" />
        </template>
      </q-select>
    </div>

    <!-- API Key Authentication -->
    <template v-if="localData.authType === 'apiKey'">
      <div class="col-12">
        <q-banner dense class="bg-blue-1 text-blue-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="key" color="blue-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.apiKey.banner.value }}
          </div>
        </q-banner>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.apiKey.headerApiKey"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.authentication.apiKey.headerName.value} *`"
          :placeholder="t.authentication.apiKey.headerNamePlaceholder.value"
          :hint="t.authentication.apiKey.headerNameHint.value"
          :rules="[(val: any) => !!val || t.authentication.apiKey.headerNameRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="label" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.apiKey.valueApiKey"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.authentication.apiKey.value.value} *`"
          :placeholder="t.authentication.apiKey.valuePlaceholder.value"
          :hint="t.authentication.apiKey.valueHint.value"
          :rules="[(val: any) => !!val || t.authentication.apiKey.valueRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="vpn_key" color="primary" />
          </template>
          <template #append>
            <q-icon
              name="content_copy"
              class="cursor-pointer"
              @click="copyToClipboard(localData.apiKey.valueApiKey)"
            >
              <AppTooltip :content="t.authentication.apiKey.copyTooltip.value" />
            </q-icon>
            <q-icon
              name="autorenew"
              class="cursor-pointer q-ml-xs"
              @click="generateApiKey"
            >
              <AppTooltip :content="t.authentication.apiKey.generateTooltip.value" />
            </q-icon>
          </template>
        </q-input>
      </div>
    </template>

    <!-- JWT Authentication -->
    <template v-if="localData.authType === 'jwt'">
      <div class="col-12">
        <q-banner dense class="bg-purple-1 text-purple-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="badge" color="purple-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.jwt.banner.value }}
          </div>
        </q-banner>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.jwt.headerName"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.authentication.jwt.headerName.value} *`"
          :placeholder="t.authentication.jwt.headerNamePlaceholder.value"
          :hint="t.authentication.jwt.headerNameHint.value"
          :rules="[(val: any) => !!val || t.authentication.jwt.headerNameRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="label" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.jwt.secretKey"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.authentication.jwt.secretKey.value} *`"
          :placeholder="t.authentication.jwt.secretKeyPlaceholder.value"
          :hint="t.authentication.jwt.secretKeyHint.value"
          readonly
          :rules="[(val: any) => !!val || t.authentication.jwt.secretKeyRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="vpn_key" color="primary" />
          </template>
          <template #append>
            <q-icon
              name="content_copy"
              class="cursor-pointer"
              @click="copyToClipboard(localData.jwt.secretKey)"
            >
              <AppTooltip :content="t.authentication.jwt.copyTooltip.value" />
            </q-icon>
            <q-icon
              name="autorenew"
              class="cursor-pointer q-ml-xs"
              @click="generateJwtSecret"
            >
              <AppTooltip :content="t.authentication.jwt.generateTooltip.value" />
            </q-icon>
          </template>
        </q-input>
      </div>
    </template>

    <!-- IP Whitelist Authentication -->
    <template v-if="localData.authType === 'ip_whitelist'">
      <div class="col-12">
        <q-banner dense class="bg-green-1 text-green-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="public" color="green-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.ipWhitelist.banner.value }}
          </div>
        </q-banner>
      </div>

      <!-- IP Address and CIDR Mask Input -->
      <div class="col-12 col-sm-8">
        <q-input
          v-model="newIpAddress"
          outlined
          dense
          class="rounded-borders"
          :label="t.authentication.ipWhitelist.ipAddress.value"
          :placeholder="t.authentication.ipWhitelist.ipAddressPlaceholder.value"
          :hint="t.authentication.ipWhitelist.ipAddressHint.value"
          :error="!!ipInputError"
          :error-message="ipInputError"
        >
          <template #prepend>
            <q-icon name="dns" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-4">
        <q-input
          v-model.number="newCidrMask"
          outlined
          dense
          type="number"
          class="rounded-borders"
          :label="t.authentication.ipWhitelist.cidrMask.value"
          :placeholder="t.authentication.ipWhitelist.cidrMaskPlaceholder.value"
          :hint="t.authentication.ipWhitelist.cidrMaskHint.value"
          :min="0"
          :max="128"
        >
          <template #prepend>
            <q-icon name="filter_list" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Add Button -->
      <div class="col-12">
        <q-btn
          dense
          outline
          color="primary"
          icon="add"
          :label="t.authentication.ipWhitelist.addButton.value"
          class="rounded-borders"
          @click="addIpToWhitelist"
        />
      </div>

      <!-- IP Whitelist Table -->
      <div v-if="localData.ipWhitelist.addresses.length > 0" class="col-12">
        <q-card flat bordered class="rounded-borders">
          <q-card-section class="q-pa-none">
            <q-list separator>
              <q-item
                v-for="(ipEntry, index) in localData.ipWhitelist.addresses"
                :key="index"
                dense
              >
                <q-item-section>
                  <q-item-label class="text-body2">
                    <q-icon name="check_circle" color="green" size="xs" class="q-mr-xs" />
                    {{ ipEntry }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn
                    flat
                    dense
                    round
                    icon="delete"
                    color="negative"
                    size="sm"
                    @click="removeIpFromWhitelist(index)"
                  >
                    <AppTooltip :content="t.authentication.ipWhitelist.removeTooltip.value" />
                  </q-btn>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Empty State -->
      <div v-else class="col-12">
        <q-banner dense class="bg-grey-2 text-grey-8 rounded-borders">
          <template #avatar>
            <q-icon name="info" color="grey-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.ipWhitelist.emptyMessage.value }}
          </div>
        </q-banner>
      </div>
    </template>

    <!-- OAuth2 Authentication -->
    <template v-if="localData.authType === 'oauth2'">
      <div class="col-12">
        <q-banner dense class="bg-orange-1 text-orange-9 rounded-borders q-mb-md">
          <template #avatar>
            <q-icon name="verified_user" color="orange-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.oauth2.banner.value }}
          </div>
        </q-banner>
      </div>

      <div class="col-12">
        <q-input
          v-model="localData.oauth2.jwksUrl"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.authentication.oauth2.jwksUrl.value} *`"
          :placeholder="t.authentication.oauth2.jwksUrlPlaceholder.value"
          :hint="t.authentication.oauth2.jwksUrlHint.value"
          :rules="[
            (val: any) => !!val || t.authentication.oauth2.jwksUrlRequired.value,
            (val: any) => /^https?:\/\/.+/.test(val) || t.authentication.oauth2.invalidUrl.value
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="link" color="primary" />
          </template>
        </q-input>
      </div>
    </template>

    <!-- None Authentication -->
    <template v-if="localData.authType === 'none'">
      <div class="col-12">
        <q-banner dense class="bg-red-1 text-red-9 rounded-borders">
          <template #avatar>
            <q-icon name="warning" color="red-6" />
          </template>
          <div class="text-caption">
            {{ t.authentication.none.banner.value }}
          </div>
        </q-banner>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step4Authentication'
});

/** TYPE IMPORTS */
import type { StepEmits, StepProps } from '../../interfaces/httpDataSource.interface';

/** VUE IMPORTS */
import { reactive, watch, ref, nextTick } from 'vue';

/** EXTERNAL IMPORTS */
import { uid, copyToClipboard as qCopyToClipboard } from 'quasar';
import { Address4, Address6 } from 'ip-address';

/** COMPOSABLES */
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifySuccess, notifyInfo } from '@utils/alert/notify';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** LOCAL IMPORTS */
import { AUTH_TYPE_OPTIONS } from '../../constants/httpDataSourceConstants';

/** PROPS & EMITS */
const props = defineProps<StepProps>();
const emit = defineEmits<StepEmits>();

/** COMPOSABLES & STORES */
const t = useHttpDataSourceCreateEditTranslations();
const logger = useLogger('Step4Authentication');

// IP Whitelist state
const newIpAddress = ref('');
const newCidrMask = ref<number | null>(24);
const ipInputError = ref('');

const localData = reactive({
  authType: props.modelValue.authType || null,
  apiKey: props.modelValue.apiKey || { headerApiKey: 'X-API-Key', valueApiKey: '' },
  jwt: props.modelValue.jwt || { secretKey: '', headerName: 'Authorization' },
  ipWhitelist: props.modelValue.ipWhitelist || { addresses: [] as string[] },
  oauth2: props.modelValue.oauth2 || { jwksUrl: '' },
});

watch(() => props.modelValue, (newVal) => {
  localData.authType = newVal.authType || null;
  localData.apiKey = newVal.apiKey || { headerApiKey: 'X-API-Key', valueApiKey: '' };
  localData.jwt = newVal.jwt || { secretKey: '', headerName: 'Authorization' };
  localData.ipWhitelist = newVal.ipWhitelist || { addresses: [] };
  localData.oauth2 = newVal.oauth2 || { jwksUrl: '' };
}, { deep: true, immediate: true });

/**
 * Handle authentication type change
 * Clears previous auth configuration when switching types and updates parent
 * @returns {Promise<void>}
 */
async function handleAuthTypeChange(): Promise<void> {
  logger.debug('Auth type changed to:', localData.authType);

  // Clear fields when auth type changes
  if (localData.authType === 'apiKey') {
    localData.apiKey = { headerApiKey: 'X-API-Key', valueApiKey: '' };
  } else if (localData.authType === 'jwt') {
    localData.jwt = { secretKey: '', headerName: 'Authorization' };
  } else if (localData.authType === 'ip_whitelist') {
    localData.ipWhitelist = { addresses: [] };
  } else if (localData.authType === 'oauth2') {
    localData.oauth2 = { jwksUrl: '' };
  }

  await nextTick();
  updateValue();
}

/**
 * Generate a new API key using GUID
 * Updates the local data and emits to parent component
 * @returns {void}
 */
function generateApiKey(): void {
  // Generate a GUID for API Key
  localData.apiKey.valueApiKey = uid();
  updateValue();
}

/**
 * Generate a new JWT secret key using GUID
 * Updates the local data and emits to parent component
 * @returns {void}
 */
function generateJwtSecret(): void {
  // Generate a GUID for JWT Secret
  localData.jwt.secretKey = uid();
  updateValue();
}

/**
 * Copy text to clipboard with user notification
 * @param {string} text - Text to copy to clipboard
 * @returns {void}
 */
function copyToClipboard(text: string): void {
  qCopyToClipboard(text)
    .then(() => {
      notifySuccess({
        message: t.notifications.copiedToClipboard.value,
        timeout: 1500,
      });
    })
    .catch(() => {
      notifyInfo({
        message: t.notifications.copyFailed.value,
        timeout: 2000,
      });
    });
}

/**
 * Add IP address/CIDR to whitelist
 * Validates IP format (IPv4/IPv6), checks for duplicates, and updates parent
 * @returns {void}
 */
function addIpToWhitelist(): void {
  ipInputError.value = '';

  if (!newIpAddress.value) {
    ipInputError.value = t.authentication.ipWhitelist.ipAddressRequired.value;
    return;
  }

  if (newCidrMask.value === null || newCidrMask.value < 0) {
    ipInputError.value = t.authentication.ipWhitelist.cidrMaskRequired.value;
    return;
  }

  try {
    // Try to validate as IPv4
    let isValid = false;
    let cidrNotation = '';

    try {
      const ipv4 = new Address4(newIpAddress.value);
      if (ipv4.isCorrect()) {
        if (newCidrMask.value > 32) {
          ipInputError.value = t.authentication.ipWhitelist.invalidIpv4Cidr.value;
          return;
        }
        cidrNotation = `${newIpAddress.value}/${newCidrMask.value}`;
        // Validate the full CIDR notation
        const cidrValidation = new Address4(cidrNotation);
        if (cidrValidation.isCorrect()) {
          isValid = true;
        }
      }
    } catch {
      // Not IPv4, try IPv6
    }

    // If not IPv4, try IPv6
    if (!isValid) {
      try {
        const ipv6 = new Address6(newIpAddress.value);
        if (ipv6.isCorrect()) {
          if (newCidrMask.value > 128) {
            ipInputError.value = t.authentication.ipWhitelist.invalidIpv6Cidr.value;
            return;
          }
          cidrNotation = `${newIpAddress.value}/${newCidrMask.value}`;
          // Validate the full CIDR notation
          const cidrValidation = new Address6(cidrNotation);
          if (cidrValidation.isCorrect()) {
            isValid = true;
          }
        }
      } catch {
        // Not IPv6 either
      }
    }

    if (!isValid) {
      ipInputError.value = t.authentication.ipWhitelist.invalidIpFormat.value;
      return;
    }

    // Check for duplicates
    if (localData.ipWhitelist.addresses.includes(cidrNotation)) {
      ipInputError.value = t.authentication.ipWhitelist.duplicateError.value;
      return;
    }

    // Add to whitelist
    localData.ipWhitelist.addresses.push(cidrNotation);

    // Reset inputs
    newIpAddress.value = '';
    newCidrMask.value = 24;
    ipInputError.value = '';

    // Update parent
    updateValue();

    notifySuccess({
      message: t.authentication.ipWhitelist.addedMessage(cidrNotation),
      timeout: 1500,
    });
  } catch (error: any) {
    ipInputError.value = error.message || t.authentication.ipWhitelist.invalidIpFormat.value;
  }
}

/**
 * Remove IP address/CIDR from whitelist
 * @param {number} index - Index of the IP entry to remove
 * @returns {void}
 */
function removeIpFromWhitelist(index: number): void {
  const removed = localData.ipWhitelist.addresses[index] ?? '';
  localData.ipWhitelist.addresses.splice(index, 1);
  updateValue();

  notifyInfo({
    message: t.authentication.ipWhitelist.removedMessage(removed),
    timeout: 1500,
  });
}

/**
 * Emit updated values to parent component
 * Merges local form data with existing model value
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', {
    ...props.modelValue,
    ...localData
  });
}
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
}
</style>

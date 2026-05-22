<script setup lang="ts">
defineOptions({
  name: 'LakeHouseProviderSelection'
});

import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';
import type { Provider } from './interfaces';

import { AppTooltip } from '@components/tooltips';

import {
  PROVIDERS,
  DEFAULT_GCP_DATA,
  DEFAULT_AWS_DATA,
  DEFAULT_AZURE_DATA,
  DEFAULT_MINIO_DATA,
} from './constants';

// this ref is automatically tied to `modelValue` + `update:modelValue`
const modelRef = defineModel<LakeHouseConfigProps>({
  default: () => ({ name: '', status: true, description: '' }),
});

/**
 * Selects a storage provider by updating the localData type
 * and resets the credentials to the provider's default values.
 *
 * @param {Provider} provider - The provider object to select.
 */
const selectProvider = (provider: Provider) => {
  modelRef.value.type = provider.id;
  // Reset credentials when changing provider
  modelRef.value.credentials = getDefaultCredentials(provider.id);
};

/**
 * Returns the default credentials object for the given provider identifier.
 *
 * @param {string} providerId - The identifier of the storage provider.
 * @returns {LakeHouseConfigProps['credentials']} The default credentials for the provider.
 */
const getDefaultCredentials = (providerId: string) => {
  switch (providerId) {
    case 'aws-s3':
      return DEFAULT_AWS_DATA;

    case 'minio':
      return DEFAULT_MINIO_DATA;

    case 'azure-blob':
      return DEFAULT_AZURE_DATA;

    case 'gcp-storage':
      return DEFAULT_GCP_DATA;

    default:
      return {} as any;
  }
};
</script>

<template>
  <!-- Provider Selection -->
  <div>
    <div class="row q-col-gutter-lg justify-center">
      <div
          v-for="provider in PROVIDERS"
          :key="provider.id"
          class="col-xs-12 col-sm-6 col-md-6"
          style="max-width: 380px;"
      >
        <q-card
            flat
            bordered
            class="provider-card cursor-pointer transition-all"
            :class="modelRef.type === provider.id
                ? 'provider-card--selected'
                : 'provider-card--default'"
            @click="selectProvider(provider)"
        >
          <!-- Selection Indicator -->
          <div
              v-if="modelRef.type === provider.id"
              class="absolute-top-right q-ma-md"
          >
            <q-icon
                size="20px"
                color="primary"
                name="check_circle"
            />
          </div>

          <q-card-section class="text-center q-py-lg">
            <!-- Provider Icon -->
            <div class="q-mb-md">
              <q-icon
                  size="60px"
                  text-color="white"
                  :color="provider.iconColor"
                  :name="provider.icon"
              />
            </div>

            <!-- Provider Content -->
            <div class="provider-content">
              <h3 class="text-subtitle1 text-weight-medium text-grey-9 q-mb-xs">
                <div class="row items-center justify-center no-wrap">
                  <span>{{ provider.name }}</span>
                  <q-icon
                      v-if="provider.hasInfo"
                      size="14px"
                      class="q-ml-xs cursor-help"
                      color="primary"
                      name="info"
                  >
                    <AppTooltip
                        :content="provider.infoText || ''"
                        anchor="top middle"
                        self="bottom middle"
                        :offset="[0, 8]"
                    />
                  </q-icon>
                </div>
              </h3>
              <p class="text-body2 text-grey-6 q-px-sm line-height-md">
                {{ provider.description }}
              </p>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>

  <!-- Additional Info Section -->
  <div class="q-mt-xl">
    <q-banner
        dense
        inline-actions
        class="bg-blue-1 text-primary rounded-borders"
    >
      <template v-slot:avatar>
        <q-icon color="primary" name="info"/>
      </template>

      <div>
        <div class="text-weight-medium q-mb-xs">API Compatibility</div>
        <div class="text-body2">
          <span class="text-bold">MinIO</span> and <span class="text-bold">DigitalOcean Spaces Object Store</span> is
          compatible with the Amazon S3 API, allowing you to use the same credentials and configurations.
        </div>
      </div>
    </q-banner>
  </div>

</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.mode-card,
.provider-card {
  border-radius: var(--mapex-radius-lg);
  position: relative;

  &--default {
    border: 2px solid var(--mapex-card-border);

    &:hover {
      border-color: var(--mapex-card-hover-border);
      box-shadow: var(--mapex-shadow-md);
      transform: scale(1.02);
    }
  }

  &--selected {
    border: 2px solid var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.05);
    box-shadow: 0 0 0 3px rgba(var(--q-primary-rgb), 0.1);
  }
}

.provider-icon {
  transition: transform 0.2s ease;
}

.provider-card:hover .provider-icon {
  transform: scale(1.1);
}

.provider-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  flex: 1;
}

.transition-all {
  transition: var(--mapex-transition-base);
}

.line-height-md {
  line-height: 1.4;
}
</style>

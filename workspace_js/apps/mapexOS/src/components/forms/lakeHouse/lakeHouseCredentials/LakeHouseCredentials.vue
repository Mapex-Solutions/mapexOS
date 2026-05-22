<script setup lang="ts">
defineOptions({
  name: 'LakeHouseCredentials'
});

import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';
import { ref, computed } from 'vue';

// this ref is automatically tied to `modelValue` + `update:modelValue`
const modelRef = defineModel<LakeHouseConfigProps>({ default: {} });

// Local state
const testing = ref(false);
const testResult = ref<{ success: boolean; message: string } | null>(null);

// Computed properties
const providerType = computed(() => modelRef.value.type);

const isS3Compatible = computed(() =>
  ['aws-s3', 'minio'].includes(providerType.value),
);

const regionHint = computed(() => {
  switch (providerType.value) {
    case 'aws-s3':
      return 'AWS region (e.g., us-east-1, eu-west-1)';
    case 'minio':
      return 'Region configured in MinIO';
    default:
      return 'Provider region';
  }
});

const regionPlaceholder = computed(() => {
  switch (providerType.value) {
    case 'aws-s3':
    case 'minio':
      return 'us-east-1';
    default:
      return '';
  }
});

const canTestConnection = computed(() => {
  const creds = modelRef.value.credentials;
  switch (providerType.value) {
    case 'aws-s3':
      return !!(creds.accessKey && creds.secretKey && creds.region && creds.bucket);
    case 'minio':
      return !!(creds.accessKey && creds.secretKey && creds.endpoint && creds.bucket);
    case 'azure-blob':
      return !!(creds.accountName && creds.accountKey && creds.containerName);
    case 'gcp-storage':
      return !!(creds.projectId && creds.keyFile && creds.bucket && creds.region);
    default:
      return false;
  }
});

/**
 * Checks if the given string is a valid URL.
 * @param url - The URL string to validate.
 * @returns True if the string is a valid URL, false otherwise.
 */
function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * Checks if the given string is valid JSON.
 * @param str - The string to parse.
 * @returns True if the string is valid JSON, false otherwise.
 */
function isValidJson(str: string): boolean {
  try {
    JSON.parse(str);
    return true;
  } catch {
    return false;
  }
}

/**
 * Tests the connection to the selected data lake provider using current credentials.
 * Simulates an API call, then sets testResult with success status and message.
 */
async function testConnection(): Promise<void> {
  testing.value = true;
  testResult.value = null;

  try {
    // Simulate API call delay
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Mock success/failure (70% chance of success)
    const success = Math.random() > 0.3;

    testResult.value = {
      success,
      message: success
        ? 'Connection established successfully!'
        : 'Connection failed. Please check your credentials.',
    };
  } catch {
    testResult.value = {
      success: false,
      message: 'Error testing connection. Please try again.',
    };
  } finally {
    testing.value = false;
  }
}
</script>

<template>
  <div class="credentials-config">
    <div class="row q-col-gutter-md">
      <!-- AWS S3 / MinIO Credentials -->
      <template v-if="isS3Compatible">
        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.accessKey"
            outlined
            label="Access Key *"
            hint="Your provider access key"
            :rules="[val => !!val || 'Access Key is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.secretKey"
            outlined
            type="password"
            label="Secret Key *"
            hint="Your provider secret key"
            :rules="[val => !!val || 'Secret Key is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.region"
            outlined
            label="Region *"
            :hint="regionHint"
            :placeholder="regionPlaceholder"
            :rules="[val => !!val || 'Region is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.bucket"
            outlined
            label="Bucket *"
            hint="Name of the bucket where data will be stored"
            :rules="[val => !!val || 'Bucket is required']"
          />
        </div>

        <!-- MinIO specific fields -->
        <template v-if="providerType === 'minio'">
          <div class="col-12 col-md-8">
            <q-input
              v-model="modelRef.credentials.endpoint"
              outlined
              placeholder="https://minio.example.com"
              label="Endpoint *"
              hint="MinIO server URL (e.g., https://minio.example.com)"
              :rules="[
                val => !!val || 'Endpoint is required',
                val => isValidUrl(val) || 'Invalid URL'
              ]"
            />
          </div>

          <div class="col-12 col-md-4">
            <q-toggle
              v-model="modelRef.credentials.useSSL"
              color="primary"
              label="Use SSL/TLS"
            />
          </div>
        </template>
      </template>

      <!-- Azure Blob Storage Credentials -->
      <template v-else-if="providerType === 'azure-blob'">
        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.accountName"
            outlined
            label="Account Name *"
            hint="Your Azure storage account name"
            :rules="[val => !!val || 'Account Name is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.accountKey"
            outlined
            type="password"
            label="Account Key *"
            hint="Storage account access key"
            :rules="[val => !!val || 'Account Key is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.containerName"
            outlined
            label="Container Name *"
            hint="Name of the container where data will be stored"
            :rules="[val => !!val || 'Container Name is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.endpoint"
            outlined
            label="Endpoint (optional)"
            hint="Custom endpoint (leave blank to use default)"
          />
        </div>
      </template>

      <!-- Google Cloud Storage Credentials -->
      <template v-else-if="providerType === 'gcp-storage'">
        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.projectId"
            outlined
            label="Project ID *"
            hint="Your Google Cloud project ID"
            :rules="[val => !!val || 'Project ID is required']"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-input
            v-model="modelRef.credentials.region"
            outlined
            placeholder="us-central1"
            label="Region *"
            hint="Google Cloud region (e.g., us-central1)"
            :rules="[val => !!val || 'Region is required']"
          />
        </div>

        <div class="col-12">
          <q-input
            v-model="modelRef.credentials.keyFile"
            outlined
            autogrow
            type="textarea"
            label="Service Account Key (JSON) *"
            hint="Paste the service account JSON key file content here"
            :rules="[
              val => !!val || 'Service Account Key is required',
              val => isValidJson(val) || 'Invalid JSON'
            ]"
          />
        </div>

        <div class="col-12">
          <q-input
            v-model="modelRef.credentials.bucket"
            outlined
            label="Bucket Name *"
            hint="Name of the Google Cloud Storage bucket"
            :rules="[val => !!val || 'Bucket Name is required']"
          />
        </div>
      </template>
    </div>

    <!-- Connection Test -->
    <div class="q-mt-lg">
      <q-card flat bordered class="bg-grey-1">
        <q-card-section>
          <div class="row items-center justify-between">
            <div>
              <div class="text-subtitle2">Test Connection</div>
              <div class="text-body2 text-grey-6">
                Verify credentials are correct
              </div>
            </div>
            <q-btn
              outline
              color="primary"
              :loading="testing"
              :disable="!canTestConnection"
              @click="testConnection"
            >
              Test
            </q-btn>
          </div>

          <div v-if="testResult" class="q-mt-md">
            <q-banner
              dense
              :class="testResult.success ? 'bg-positive text-white' : 'bg-negative text-white'"
            >
              <template v-slot:avatar>
                <q-icon :name="testResult.success ? 'check_circle' : 'error'" />
              </template>
              {{ testResult.message }}
            </q-banner>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>
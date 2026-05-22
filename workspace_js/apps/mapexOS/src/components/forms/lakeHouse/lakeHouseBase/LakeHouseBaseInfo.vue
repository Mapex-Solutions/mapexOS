<!-- Child.vue -->
<script setup lang="ts">
defineOptions({
  name: 'LakeHouseBaseInfo'
});

import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';

// this ref is automatically tied to `modelValue` + `update:modelValue`
const modelRef = defineModel<LakeHouseConfigProps>({
  default: () => ({ name: '', status: true, description: '' }),
});
</script>

<template>
  <div class="row q-col-gutter-md">
    <div class="col-12 col-md-8">
      <q-input
          v-model="modelRef.name"
          outlined
          label="Channel Name *"
          :rules="[val => !!val || 'Name is required']"
      />
    </div>

    <div class="col-12 col-md-4">
      <q-toggle
          v-model="modelRef.status"
          color="positive"
          label="Active Status"
      />
    </div>

    <div class="col-12">
      <q-input
          v-model="modelRef.description"
          outlined
          type="textarea"
          label="Description"
      />
    </div>
  </div>
</template>

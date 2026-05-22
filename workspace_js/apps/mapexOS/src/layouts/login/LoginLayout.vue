<script setup lang="ts">
/** COMPOSABLES */
import { useTheme } from '@composables/theme';

/** COMPOSABLES & STORES */
useTheme();
</script>

<template>
  <q-layout view="lHh Lpr lFf">
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

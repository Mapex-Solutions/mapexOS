<template>
  <router-view />
</template>

<script setup lang="ts">
/**
 * App.vue - Main Application Component
 *
 * CRITICAL: This component is intentionally MINIMAL
 *
 * All authentication, authorization, and organization context logic
 * is handled by the Router Guard (beforeEach.guard.ts).
 *
 * Guard-First Architecture:
 * - Router Guard hydrates auth from storage
 * - Router Guard checks authentication
 * - Router Guard loads organization coverage
 * - Router Guard ensures organization is selected
 * - Router Guard handles permissions (future)
 *
 * This component ONLY renders the router-view.
 * NO auth logic, NO org logic, NO initialization logic.
 *
 * Benefits:
 * - Clear separation of concerns
 * - Guard executes BEFORE component lifecycle
 * - No race conditions
 * - Consistent behavior across all navigation scenarios
 * - Easy to test and maintain
 *
 * NOTE: Monaco Editor workers are configured in @utils/monaco-setup.ts
 * and are lazy-loaded when Monaco components are used.
 */
</script>

<script setup lang="ts">
defineOptions({
  name: 'TabGroups'
});

/** TYPE IMPORTS */
import type { TabGroupsProps, UserGroupInfo } from '../../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useUserDetailTranslations } from '@composables/i18n';

/** PROPS */
const props = defineProps<TabGroupsProps>();

/** COMPOSABLES & STORES */
const t = useUserDetailTranslations();

/** COMPUTED */

/**
 * Groups from user data (passed from parent via API response)
 */
const groups = computed<UserGroupInfo[]>(() => props.user?.groups || []);

/**
 * Loading state from props
 */
const loading = computed(() => props.loading);

/**
 * Groups count
 */
const groupsCount = computed(() => props.user?.groupsCount ?? groups.value.length);
</script>

<template>
  <div class="tab-groups">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-lg">
      <q-spinner color="primary" size="40px" />
    </div>

    <!-- Empty State -->
    <div v-else-if="groups.length === 0" class="section text-center q-pa-xl">
      <q-icon name="groups" size="64px" color="grey-4" class="q-mb-md" />
      <div class="text-h6 text-grey-8 q-mb-sm">{{ t.groups.empty.title.value }}</div>
      <div class="text-body2 text-grey-6">{{ t.groups.empty.description.value }}</div>
    </div>

    <!-- Groups List -->
    <div v-else>
      <div class="section">
        <div class="section-header q-mb-md">
          <q-icon name="groups" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.groups.count.value }}</span>
          <q-badge color="primary" class="q-ml-sm">{{ groupsCount }}</q-badge>
        </div>

        <div class="row q-col-gutter-md">
          <div v-for="group in groups" :key="group.id" class="col-12 col-md-6 col-lg-4">
            <div class="group-card">
              <div class="row items-center q-mb-sm">
                <q-icon name="groups" color="primary" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-medium">{{ group.name }}</div>
              </div>

              <div v-if="group.description" class="text-caption text-grey-7 ellipsis-2-lines">
                {{ group.description }}
              </div>

              <div v-else class="text-caption text-grey-5 text-italic">
                {{ t.groups.noDescription?.value || 'No description' }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.tab-groups {
  .section {
    background: var(--mapex-surface-bg);
    border-radius: var(--mapex-radius-md);
    padding: 16px;
  }

  .section-header {
    display: flex;
    align-items: center;
    border-bottom: 1px solid var(--mapex-card-border);
    padding-bottom: 8px;
  }

  .group-card {
    background: white;
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    padding: 16px;
    height: 100%;
    transition: var(--mapex-transition-base);

    &:hover {
      box-shadow: var(--mapex-shadow-md);
    }
  }

  .ellipsis-2-lines {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
}
</style>

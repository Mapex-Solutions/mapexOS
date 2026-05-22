<template>
  <q-form ref="formRef" greedy>
    <!-- Header -->
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="admin_panel_settings" color="primary" class="q-mr-xs" />
        {{ t.sections.roles.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.roles.value }}
      </div>
    </div>

    <!-- Selected Count & Add Button -->
    <div class="row items-center q-mb-md">
      <div class="col">
        <DetailChip
          :value="`${selectedRoles.length} ${selectedRoles.length === 1 ? t.labels.roleSelected.value : t.labels.rolesSelected.value}`"
          icon="check_circle"
          color="primary"
          size="md"
        />
      </div>
      <div class="col-auto">
        <q-btn
          unelevated
          dense
          icon="add"
          :label="t.labels.addRoles.value"
          color="primary"
          size="sm"
          class="rounded-borders"
          @click="showRoleDrawer = true"
        />
      </div>
    </div>

    <!-- Roles List -->
    <div class="roles-list">
      <!-- Empty State -->
      <div v-if="selectedRoles.length === 0" class="text-center q-pa-xl">
        <q-icon name="admin_panel_settings" size="4em" color="grey-4" />
        <div class="text-grey-6 q-mt-md text-body1">
          {{ t.labels.noRolesSelected.value }}
        </div>
        <div class="text-grey-5 text-caption q-mt-xs">
          {{ t.labels.clickAddRoles.value }}
        </div>
      </div>

      <!-- Roles List -->
      <q-list v-else bordered separator class="rounded-borders">
        <q-item
          v-for="role in selectedRoles"
          :key="role.id"
        >
          <q-item-section avatar>
            <q-avatar
              color="primary"
              text-color="white"
              icon="admin_panel_settings"
              size="md"
            />
          </q-item-section>

          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ role.name }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-btn
              flat
              round
              dense
              icon="close"
              color="red-6"
              @click="removeRole(role)"
            >
              <AppTooltip :content="t.labels.removeRole.value" />
            </q-btn>
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- Info message -->
    <div class="text-caption text-grey-6 q-mt-md">
      <q-icon name="info" size="xs" class="q-mr-xs" />
      {{ t.labels.rolesRequired.value }}
    </div>

    <!-- Role Selector Drawer -->
    <RoleMultiSelectorDrawer
      v-model="showRoleDrawer"
      :selected-role-ids="selectedRoleIds"
      @confirm="onRolesSelected"
      @cancel="showRoleDrawer = false"
    />
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2Roles',
});

/** TYPE IMPORTS */
import type { QForm } from 'quasar';
import type { RoleSelectionItem } from '../../interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { RoleMultiSelectorDrawer } from '@components/drawers/roles';

/** COMPOSABLES */
import { useGroupsTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<{
  /** Selected roles array */
  selectedRoles: RoleSelectionItem[];
}>();

const emit = defineEmits<{
  (e: 'update:selected-roles', roles: RoleSelectionItem[]): void;
}>();

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showRoleDrawer = ref(false);

/** COMPUTED */

/**
 * Selected role IDs for the drawer
 */
const selectedRoleIds = computed(() => props.selectedRoles.map(r => r.id));

/** FUNCTIONS */

/**
 * Handle roles selected from drawer
 *
 * @param {any[]} roles - Selected roles from drawer
 */
function onRolesSelected(roles: any[]): void {
  const mappedRoles: RoleSelectionItem[] = roles.map(r => ({
    id: r.id,
    name: r.name,
  }));
  emit('update:selected-roles', mappedRoles);
  showRoleDrawer.value = false;
}

/**
 * Remove a role from selection
 *
 * @param {RoleSelectionItem} role - Role to remove
 */
function removeRole(role: RoleSelectionItem): void {
  const newRoles = props.selectedRoles.filter(r => r.id !== role.id);
  emit('update:selected-roles', newRoles);
}

/** EXPOSE */
defineExpose({
  formRef,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.roles-list {
  .q-item {
    transition: background-color 0.2s ease;
  }
}
</style>

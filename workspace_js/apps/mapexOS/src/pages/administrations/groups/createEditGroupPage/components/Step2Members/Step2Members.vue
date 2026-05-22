<template>
  <q-form ref="formRef" greedy>
    <!-- Header -->
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="person_add" color="primary" class="q-mr-xs" />
        {{ t.sections.members.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.members.value }}
      </div>
    </div>

    <!-- Selected Count & Add Button -->
    <div class="row items-center q-mb-md">
      <div class="col">
        <DetailChip
          :value="`${displayMembersCount} ${displayMembersCount === 1 ? t.labels.memberSelected.value : t.labels.membersSelected.value}`"
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
          :label="isEditMode ? t.labels.addMoreMembers.value : t.labels.addMembers.value"
          color="primary"
          size="sm"
          class="rounded-borders"
          @click="showUserDrawer = true"
        />
      </div>
    </div>

    <!-- Loading State (Edit Mode) -->
    <div v-if="loadingExistingMembers" class="text-center q-pa-lg">
      <q-spinner color="primary" size="2em" />
      <div class="text-grey-7 q-mt-sm">{{ t.labels.loadingMembers.value }}</div>
    </div>

    <!-- Members List -->
    <div v-else class="members-list">
      <!-- Empty State -->
      <div v-if="displayMembers.length === 0" class="text-center q-pa-xl">
        <q-icon name="group" size="4em" color="grey-4" />
        <div class="text-grey-6 q-mt-md text-body1">
          {{ isEditMode ? t.labels.noMembersInGroup.value : t.labels.noMembersSelectedYet.value }}
        </div>
        <div class="text-grey-5 text-caption q-mt-xs">
          {{ t.labels.clickAddMembers.value }}
        </div>
      </div>

      <!-- Members List with Scroll -->
      <q-scroll-area
        v-else
        style="height: 400px;"
        @scroll="onMembersScroll"
      >
        <q-list bordered separator class="rounded-borders">
          <q-item
            v-for="member in displayMembers"
            :key="member.id"
            :class="{ 'bg-red-1': isPendingRemoval(member.id) }"
          >
            <q-item-section avatar>
              <q-avatar
                color="primary"
                text-color="white"
                size="md"
              >
                {{ getInitials(member) }}
              </q-avatar>
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-medium">
                {{ getDisplayName(member) }}
              </q-item-label>
              <q-item-label caption>
                {{ member.email }}
              </q-item-label>
            </q-item-section>

            <!-- Pending Add Badge -->
            <q-item-section v-if="isPendingAddition(member.id)" side>
              <q-badge color="green-6" :label="t.labels.badgeNew.value.toUpperCase()" />
            </q-item-section>

            <!-- Pending Remove Badge -->
            <q-item-section v-if="isPendingRemoval(member.id)" side>
              <q-badge color="red-6" :label="t.labels.badgeRemoving.value.toUpperCase()" />
            </q-item-section>

            <q-item-section side>
              <q-btn
                flat
                round
                dense
                :icon="isPendingRemoval(member.id) ? 'undo' : 'close'"
                :color="isPendingRemoval(member.id) ? 'green-6' : 'red-6'"
                @click="toggleMemberRemoval(member)"
              >
                <AppTooltip :content="isPendingRemoval(member.id) ? t.labels.undoRemoval.value : t.labels.removeMember.value" />
              </q-btn>
            </q-item-section>
          </q-item>
        </q-list>

        <!-- Load More Spinner (Edit Mode) -->
        <div v-if="loadingMoreMembers" class="text-center q-pa-md">
          <q-spinner color="primary" />
        </div>
      </q-scroll-area>
    </div>

    <!-- Info message -->
    <div class="text-caption text-grey-6 q-mt-md">
      <q-icon name="info" size="xs" class="q-mr-xs" />
      {{ t.labels.membersOptional.value }}
    </div>

    <!-- User Selector Drawer -->
    <UserMultiSelectorDrawer
      v-model="showUserDrawer"
      :exclude-user-ids="allMemberIds"
      @confirm="onUsersSelected"
      @cancel="showUserDrawer = false"
    />
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step3Members',
});

/** TYPE IMPORTS */
import type { QForm } from 'quasar';
import type { UserSelectorItem } from '@components/drawers';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { UserMultiSelectorDrawer } from '@components/drawers';

/** COMPOSABLES */
import { useGroupsTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

const logger = useLogger('Step3Members');

/** LOCAL IMPORTS */
import type { MemberDisplayItem } from './interfaces/Step2Members.interface';

/** PROPS & EMITS */
const props = defineProps<{
  /** Whether in edit mode */
  isEditMode: boolean;

  /** Group ID for edit mode */
  groupId: string | undefined;

  /** Selected member IDs (for compatibility) */
  selectedMembers: string[];
}>();

const emit = defineEmits<{
  (e: 'update:selected-members', members: string[]): void;
  (e: 'pending-changes', changes: { additions: UserSelectorItem[]; removals: string[] }): void;
}>();

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showUserDrawer = ref(false);

// Existing members (from API in edit mode)
const existingMembers = ref<MemberDisplayItem[]>([]);
const loadingExistingMembers = ref(false);
const loadingMoreMembers = ref(false);
const existingMembersPage = ref(1);
const hasMoreExistingMembers = ref(true);

// Pending changes
const pendingAdditions = ref<UserSelectorItem[]>([]);
const pendingRemovals = ref<string[]>([]);

/** COMPUTED */

/**
 * All current member IDs (for excluding in drawer)
 */
const allMemberIds = computed(() => {
  const existingIds = existingMembers.value.map(m => m.id);
  const pendingIds = pendingAdditions.value.map(m => m.id);
  // Exclude pending removals from exclusion list (they can be re-added)
  const activeExistingIds = existingIds.filter(id => !pendingRemovals.value.includes(id));
  return [...activeExistingIds, ...pendingIds];
});

/**
 * Display members list
 */
const displayMembers = computed((): MemberDisplayItem[] => {
  // Start with existing members (excluding fully removed)
  const fromExisting = existingMembers.value.map(m => ({
    id: m.id,
    firstName: m.firstName,
    lastName: m.lastName,
    email: m.email,
  }));

  // Add pending additions
  const fromPending = pendingAdditions.value.map(m => ({
    id: m.id,
    firstName: m.firstName,
    lastName: m.lastName,
    email: m.email,
  }));

  return [...fromExisting, ...fromPending];
});

/**
 * Count of displayed members
 */
const displayMembersCount = computed(() => {
  // Count existing - removals + additions
  const existingCount = existingMembers.value.length - pendingRemovals.value.length;
  const additionsCount = pendingAdditions.value.length;
  return existingCount + additionsCount;
});

/** FUNCTIONS */

/**
 * Check if member is pending addition
 *
 * @param {string} memberId - Member ID
 * @returns {boolean} True if pending addition
 */
function isPendingAddition(memberId: string): boolean {
  return pendingAdditions.value.some(m => m.id === memberId);
}

/**
 * Check if member is pending removal
 *
 * @param {string} memberId - Member ID
 * @returns {boolean} True if pending removal
 */
function isPendingRemoval(memberId: string): boolean {
  return pendingRemovals.value.includes(memberId);
}

/**
 * Toggle member removal status
 *
 * @param {MemberDisplayItem} member - Member to toggle
 */
function toggleMemberRemoval(member: MemberDisplayItem): void {
  // If it's a pending addition, just remove from additions
  if (isPendingAddition(member.id)) {
    const index = pendingAdditions.value.findIndex(m => m.id === member.id);
    if (index >= 0) {
      pendingAdditions.value.splice(index, 1);
    }
    emitChanges();
    return;
  }

  // If it's an existing member, toggle pending removal
  const removalIndex = pendingRemovals.value.indexOf(member.id);
  if (removalIndex >= 0) {
    // Undo removal
    pendingRemovals.value.splice(removalIndex, 1);
  } else {
    // Mark for removal
    pendingRemovals.value.push(member.id);
  }
  emitChanges();
}

/**
 * Handle users selected from drawer
 *
 * @param {UserSelectorItem[]} users - Selected users
 */
function onUsersSelected(users: UserSelectorItem[]): void {
  // Add to pending additions (avoiding duplicates)
  const existingAdditionIds = new Set(pendingAdditions.value.map(m => m.id));
  const existingMemberIds = new Set(existingMembers.value.map(m => m.id));

  users.forEach(user => {
    if (!existingAdditionIds.has(user.id) && !existingMemberIds.has(user.id)) {
      pendingAdditions.value.push(user);
    }
  });

  showUserDrawer.value = false;
  emitChanges();
}

/**
 * Emit pending changes to parent
 */
function emitChanges(): void {
  emit('pending-changes', {
    additions: [...pendingAdditions.value],
    removals: [...pendingRemovals.value],
  });

  // Also emit selected members for compatibility
  const allIds = displayMembers.value
    .filter(m => !isPendingRemoval(m.id))
    .map(m => m.id);
  emit('update:selected-members', allIds);
}

/**
 * Get initials from member
 *
 * @param {MemberDisplayItem} member - Member object
 * @returns {string} Member initials
 */
function getInitials(member: MemberDisplayItem): string {
  const first = member.firstName?.charAt(0) || '';
  const last = member.lastName?.charAt(0) || '';
  if (first || last) {
    return (first + last).toUpperCase();
  }
  return member.email?.charAt(0).toUpperCase() || '?';
}

/**
 * Get display name for member
 *
 * @param {MemberDisplayItem} member - Member object
 * @returns {string} Display name
 */
function getDisplayName(member: MemberDisplayItem): string {
  const name = `${member.firstName || ''} ${member.lastName || ''}`.trim();
  return name || member.email;
}

/**
 * Load existing group members (edit mode only)
 *
 * @param {boolean} append - Whether to append to existing list
 * @returns {Promise<void>}
 */
async function loadExistingMembers(append: boolean = false): Promise<void> {
  if (!props.isEditMode || !props.groupId || !apis.mapexOS?.groups) {
    return;
  }

  if (append) {
    loadingMoreMembers.value = true;
  } else {
    loadingExistingMembers.value = true;
    existingMembersPage.value = 1;
    existingMembers.value = [];
  }

  try {
    const response = await apis.mapexOS.groups.getMembers(
      { groupId: props.groupId },
      { page: existingMembersPage.value, perPage: 20 },
    );

    const mappedMembers: MemberDisplayItem[] = (response.items || []).map((m: any) => ({
      id: m.userId,
      firstName: m.userFirstName || '',
      lastName: m.userLastName || '',
      email: m.userEmail || '',
    }));

    if (append) {
      existingMembers.value = [...existingMembers.value, ...mappedMembers];
    } else {
      existingMembers.value = mappedMembers;
    }

    const totalPages = response.pagination?.totalPages || 1;
    hasMoreExistingMembers.value = existingMembersPage.value < totalPages;

    logger.debug('Loaded group members:', {
      page: existingMembersPage.value,
      loaded: mappedMembers.length,
      total: existingMembers.value.length,
      hasMore: hasMoreExistingMembers.value,
    });

    // Emit changes to parent after loading members
    emitChanges();
  } catch (error: any) {
    logger.error('Failed to load group members:', error);
  } finally {
    loadingExistingMembers.value = false;
    loadingMoreMembers.value = false;
  }
}

/**
 * Handle scroll for infinite loading
 *
 * @param {any} info - Scroll info
 */
function onMembersScroll(info: any): void {
  if (!props.isEditMode || !hasMoreExistingMembers.value || loadingMoreMembers.value) {
    return;
  }

  const scrollPosition = info.verticalPosition;
  const scrollSize = info.verticalSize;
  const containerSize = info.verticalContainerSize;

  if (scrollPosition + containerSize >= scrollSize * 0.8) {
    existingMembersPage.value++;
    void loadExistingMembers(true);
  }
}

/** WATCHERS */

/**
 * Watch for group ID changes to reload members
 */
watch(() => props.groupId, () => {
  if (props.isEditMode && props.groupId) {
    void loadExistingMembers();
  }
}, { immediate: true });

/** LIFECYCLE HOOKS */
onMounted(() => {
  if (props.isEditMode && props.groupId) {
    void loadExistingMembers();
  }
});

/** EXPOSE */
defineExpose({
  formRef,
  pendingAdditions,
  pendingRemovals,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.members-list {
  .q-item {
    transition: background-color 0.2s ease;
  }
}

.bg-red-1 {
  background-color: rgba(var(--q-negative-rgb), 0.08);
}
</style>

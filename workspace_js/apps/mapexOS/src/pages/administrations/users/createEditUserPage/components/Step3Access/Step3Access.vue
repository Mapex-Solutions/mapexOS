<template>
  <q-form ref="formRef" greedy>
    <!-- Section Header -->
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="admin_panel_settings" color="primary" class="q-mr-xs" />
        {{ t.sections.access.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.access.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Access Type Selection -->
      <div class="col-12">
        <div class="text-body2 text-weight-medium q-mb-sm">
          {{ t.fields.accessType.value }} *
        </div>
        <div class="row q-col-gutter-sm items-stretch">
          <div
            v-for="option in ACCESS_TYPE_OPTIONS"
            :key="option.value"
            class="col-12 col-sm-6"
          >
            <q-card
              flat
              bordered
              class="selectable-card cursor-pointer full-height"
              :class="{
                'selectable-card--selected': localData.accessType === option.value,
                'selectable-card--warning': option.warning && localData.accessType === option.value
              }"
              :data-testid="`user-access-type-${option.value}`"
              @click="selectAccessType(option.value)"
            >
              <q-card-section class="q-pa-md full-height">
                <div class="row items-start no-wrap full-height">
                  <q-icon
                    :name="option.icon"
                    size="md"
                    :color="localData.accessType === option.value ? 'primary' : 'grey-6'"
                    class="q-mr-md q-mt-xs"
                  />
                  <div class="col">
                    <div class="row items-center q-mb-xs">
                      <div class="text-subtitle2 text-weight-medium">
                        {{ option.label }}
                      </div>
                      <DetailChip
                        v-if="option.recommended"
                        :label="t.labels.recommended.value"
                        color="positive"
                        size="xs"
                        class="q-ml-sm"
                      />
                    </div>
                    <div class="text-caption text-grey-7">
                      {{ option.description }}
                    </div>
                    <div v-if="option.warning" class="text-caption text-warning q-mt-xs">
                      <q-icon name="warning" size="xs" class="q-mr-xs" />
                      {{ option.warning }}
                    </div>
                  </div>
                  <q-radio
                    v-model="localData.accessType"
                    :val="option.value"
                    color="primary"
                    class="q-ml-sm"
                  />
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>

      <!-- GROUP Selection (when accessType = 'group') -->
      <template v-if="localData.accessType === 'group'">
        <!-- No Organization Context Warning -->
        <div v-if="!hasOrgContext" class="col-12">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            <div class="text-body2">
              {{ t.messages.noOrgContext.value }}
            </div>
          </q-banner>
        </div>

        <!-- Organization (from context - fixed) -->
        <div v-else class="col-12">
          <div class="text-body2 text-weight-medium q-mb-xs">
            {{ t.fields.organization.value }}
          </div>
          <div class="selector-field selector-field--readonly rounded-borders">
            <div class="row items-center no-wrap q-pa-sm">
              <q-icon name="business" color="primary" size="sm" class="q-mr-sm" />
              <q-avatar
                color="primary"
                icon="business"
                text-color="white"
                size="sm"
                class="q-mr-sm"
              />
              <div class="col">
                <div class="text-body2">{{ currentOrg.name }}</div>
                <div class="text-caption text-grey-6">
                  {{ t.hints.orgFromContext.value }}
                </div>
              </div>
              <q-icon name="lock" size="xs" color="grey-5" />
            </div>
          </div>
        </div>

        <!-- Group Access Mode Selection -->
        <div v-if="hasOrgContext" class="col-12">
          <div class="text-body2 text-weight-medium q-mb-sm">
            Group Access Mode *
          </div>
          <div class="row q-col-gutter-sm items-stretch">
            <div
              v-for="option in GROUP_ACCESS_MODE_OPTIONS"
              :key="option.value"
              class="col-12 col-sm-6"
            >
              <q-card
                flat
                bordered
                class="selectable-card cursor-pointer full-height"
                :class="{ 'selectable-card--selected': localData.groupAccessMode === option.value }"
                @click="selectGroupAccessMode(option.value)"
              >
                <q-card-section class="q-pa-md full-height">
                  <div class="row items-start no-wrap full-height">
                    <q-icon
                      :name="option.icon"
                      size="md"
                      :color="localData.groupAccessMode === option.value ? 'primary' : 'grey-6'"
                      class="q-mr-md q-mt-xs"
                    />
                    <div class="col">
                      <div class="row items-center q-mb-xs">
                        <div class="text-subtitle2 text-weight-medium">
                          {{ option.label }}
                        </div>
                        <DetailChip
                          v-if="option.recommended"
                          :label="t.labels.recommended.value"
                          color="positive"
                          size="xs"
                          class="q-ml-sm"
                        />
                      </div>
                      <div class="text-caption text-grey-7">
                        {{ option.description }}
                      </div>
                    </div>
                    <q-radio
                      v-model="localData.groupAccessMode"
                      :val="option.value"
                      color="primary"
                      class="q-ml-sm"
                    />
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>

        <!-- EXISTING GROUP: Group Selector Only -->
        <template v-if="hasOrgContext && localData.groupAccessMode === 'existing'">
          <!-- EDIT MODE: Multiple Groups Display -->
          <template v-if="isEditMode">
            <div class="col-12">
              <div class="text-body2 text-weight-medium q-mb-xs">
                {{ t.fields.groups?.value || 'Groups' }} *
              </div>

              <!-- Selected Groups List -->
              <div v-if="localData.selectedGroups.length > 0" class="selected-groups-list q-mb-sm">
                <q-card
                  v-for="(group, index) in localData.selectedGroups"
                  :key="group.existingGroup?.groupId || index"
                  flat
                  bordered
                  class="selected-group-card q-mb-xs"
                >
                  <q-card-section class="q-pa-sm">
                    <div class="row items-center no-wrap">
                      <q-avatar
                        color="secondary"
                        icon="group"
                        text-color="white"
                        size="sm"
                        class="q-mr-sm"
                      />
                      <div class="col">
                        <div class="text-body2 text-weight-medium">
                          {{ group.existingGroup?.groupName || group.newGroup?.name }}
                        </div>
                        <div v-if="group.mode === 'new'" class="text-caption text-grey-6">
                          {{ t.labels.newGroup?.value || 'New group' }}
                        </div>
                      </div>
                      <q-btn
                        icon="close"
                        flat
                        round
                        dense
                        size="sm"
                        color="grey-6"
                        @click="removeGroupFromList(index)"
                      >
                        <AppTooltip :content="t.actions.remove?.value || 'Remove'" />
                      </q-btn>
                    </div>
                  </q-card-section>
                </q-card>
              </div>

              <!-- Empty State -->
              <div v-else class="text-grey-6 text-body2 q-mb-sm q-pa-md bg-grey-2 rounded-borders text-center">
                {{ t.placeholders.noGroupsSelected?.value || 'No groups selected' }}
              </div>

              <!-- Add Group Button -->
              <q-btn
                outline
                color="primary"
                icon="add"
                :label="t.actions.addGroup?.value || 'Add Group'"
                class="full-width"
                @click="openGroupDrawer"
              />

              <div class="text-caption text-grey-6 q-mt-xs q-ml-sm">
                {{ t.hints.multipleGroups?.value || 'User will inherit roles from all selected groups' }}
              </div>
            </div>
          </template>

          <!-- CREATE MODE: Single Group Selector -->
          <template v-else>
            <div class="col-12">
              <div class="text-body2 text-weight-medium q-mb-xs">
                Select Group *
              </div>
              <div
                class="selector-field rounded-borders cursor-pointer"
                :class="{ 'has-error': !localData.selectedGroup?.existingGroup && showValidation }"
                data-testid="user-group-select-btn"
                @click="openGroupDrawer"
              >
                <div class="row items-center no-wrap q-pa-sm">
                  <q-icon name="group" color="primary" size="sm" class="q-mr-sm" />

                  <!-- Group Selected -->
                  <div v-if="localData.selectedGroup?.existingGroup" class="col">
                    <div class="row items-center no-wrap">
                      <q-avatar
                        color="secondary"
                        icon="group"
                        text-color="white"
                        size="sm"
                        class="q-mr-sm"
                      />
                      <div class="col">
                        <div class="text-body2 text-weight-medium">
                          {{ localData.selectedGroup.existingGroup.groupName }}
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- No Group Selected -->
                  <div v-else class="text-grey-6 col">
                    {{ t.placeholders.selectGroup.value }}
                  </div>

                  <q-btn
                    v-if="localData.selectedGroup?.existingGroup"
                    icon="close"
                    flat
                    round
                    dense
                    size="sm"
                    @click.stop="clearGroup"
                  />
                  <q-icon name="arrow_forward_ios" size="xs" color="grey-6" />
                </div>
              </div>
              <div class="text-caption text-grey-6 q-mt-xs q-ml-sm">
                User will inherit all roles from the selected group
              </div>
            </div>
          </template>

          <!-- Info Banner for Existing Group -->
          <div class="col-12">
            <q-banner rounded class="bg-blue-1 text-blue-9">
              <template #avatar>
                <q-icon name="info" color="blue-6" />
              </template>
              <div class="text-body2">
                {{ t.hints.groupInheritance.value }}
              </div>
            </q-banner>
          </div>
        </template>

        <!-- NEW GROUP: Name + Roles -->
        <template v-if="hasOrgContext && localData.groupAccessMode === 'new'">
          <!-- Group Name -->
          <div class="col-12">
            <q-input
              v-model="newGroupName"
              outlined
              :label="`${t.fields.newGroupName.value} *`"
              :placeholder="t.placeholders.newGroupName.value"
              :rules="[(v: string) => !!v || t.validation.newGroupNameRequired.value, (v: string) => v.length >= 3 || t.validation.newGroupNameMinLength.value]"
              lazy-rules
            >
              <template #prepend>
                <q-icon name="group_add" />
              </template>
            </q-input>
          </div>

          <!-- Group Description -->
          <div class="col-12">
            <q-input
              v-model="newGroupDescription"
              outlined
              :label="t.fields.newGroupDescription.value"
              :placeholder="t.placeholders.newGroupDescription.value"
              type="textarea"
              rows="2"
            >
              <template #prepend>
                <q-icon name="description" />
              </template>
            </q-input>
          </div>

          <!-- Role Selection for New Group -->
          <div class="col-12">
            <div class="text-body2 text-weight-medium q-mb-xs">
              {{ t.fields.newGroupRoles.value }} *
            </div>
            <div
              class="selector-field rounded-borders cursor-pointer"
              :class="{
                'has-error': (!localData.selectedGroup?.newGroup?.roleIds?.length) && showValidation
              }"
              @click="openGroupRolesDrawer()"
            >
              <div class="row items-center no-wrap q-pa-sm">
                <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
                <div v-if="localData.selectedGroup?.newGroup?.roleIds?.length" class="row items-center q-gutter-xs col q-py-xs">
                  <div
                    v-for="(roleName, index) in localData.selectedGroup.newGroup.roleNames"
                    :key="index"
                    class="row items-center no-wrap"
                  >
                    <DetailChip
                      :label="roleName"
                      color="primary"
                      size="sm"
                    />
                    <q-btn
                      icon="close"
                      flat
                      round
                      dense
                      size="xs"
                      color="grey-6"
                      class="q-ml-xs"
                      @click.stop="removeNewGroupRole(index)"
                    >
                      <AppTooltip :content="t.actions.remove?.value || 'Remove'" />
                    </q-btn>
                  </div>
                </div>
                <div v-else class="text-grey-6 col">
                  {{ t.placeholders.selectRoles.value }}
                </div>
                <q-btn
                  v-if="localData.selectedGroup?.newGroup?.roleIds?.length"
                  icon="close"
                  flat
                  round
                  dense
                  size="sm"
                  @click.stop="clearNewGroupRoles"
                >
                  <AppTooltip :content="t.actions.clearAll?.value || 'Clear all'" />
                </q-btn>
                <q-icon name="arrow_forward_ios" size="xs" color="grey-6" />
              </div>
            </div>
            <div class="text-caption text-grey-6 q-mt-xs q-ml-sm">
              These roles will be assigned to the group's membership
            </div>
          </div>

          <!-- Info Banner for New Group -->
          <div class="col-12">
            <q-banner rounded class="bg-amber-1 text-amber-9">
              <template #avatar>
                <q-icon name="info" color="amber-6" />
              </template>
              <div class="text-body2">
                A new group will be created with the specified roles. The user will be automatically added to this group.
              </div>
            </q-banner>
          </div>
        </template>
      </template>

      <!-- DIRECT Assignment (when accessType = 'direct') -->
      <template v-if="localData.accessType === 'direct'">
        <!-- Warning Banner -->
        <div class="col-12">
          <q-banner rounded class="bg-warning text-white">
            <template #avatar>
              <q-icon name="warning" color="white" />
            </template>
            <div class="text-body2">
              {{ t.hints.directWarning.value }}
            </div>
          </q-banner>
        </div>

        <!-- No Organization Context Warning -->
        <div v-if="!hasOrgContext" class="col-12">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            <div class="text-body2">
              {{ t.messages.noOrgContext?.value || 'No organization selected. Please select an organization from the sidebar first.' }}
            </div>
          </q-banner>
        </div>

        <!-- Organization (from context - fixed) -->
        <div v-else class="col-12">
          <div class="text-body2 text-weight-medium q-mb-xs">
            {{ t.fields.organization.value }}
          </div>
          <div class="selector-field selector-field--readonly rounded-borders">
            <div class="row items-center no-wrap q-pa-sm">
              <q-icon name="business" color="primary" size="sm" class="q-mr-sm" />
              <q-avatar
                color="primary"
                icon="business"
                text-color="white"
                size="sm"
                class="q-mr-sm"
              />
              <div class="col">
                <div class="text-body2">{{ currentOrg.name }}</div>
                <div class="text-caption text-grey-6">
                  {{ t.hints.orgFromContext?.value || 'Organization from current context' }}
                </div>
              </div>
              <q-icon name="lock" size="xs" color="grey-5" />
            </div>
          </div>
        </div>

        <!-- EDIT MODE: Multiple Direct Memberships -->
        <template v-if="hasOrgContext && isEditMode">
          <div class="col-12">
            <div class="text-body2 text-weight-medium q-mb-xs">
              {{ t.fields.directMemberships?.value || 'Direct Memberships' }} *
            </div>

            <!-- Memberships List -->
            <div v-if="localData.directMemberships.length > 0" class="memberships-list q-mb-sm">
              <q-card
                v-for="(membership, index) in localData.directMemberships"
                :key="index"
                flat
                bordered
                class="membership-card q-mb-sm"
              >
                <q-card-section class="q-pa-sm">
                  <div class="row items-start no-wrap">
                    <q-avatar
                      color="primary"
                      icon="admin_panel_settings"
                      text-color="white"
                      size="sm"
                      class="q-mr-sm q-mt-xs"
                    />
                    <div class="col">
                      <!-- Roles -->
                      <div class="row items-center q-gutter-xs q-mb-xs">
                        <DetailChip
                          v-for="(roleName, roleIdx) in membership.roleNames"
                          :key="roleIdx"
                          :label="roleName"
                          color="primary"
                          size="sm"
                        />
                      </div>
                      <!-- Scope Selector -->
                      <div class="row items-center q-gutter-sm">
                        <span class="text-caption text-grey-7">{{ t.fields.scope?.value || 'Scope' }}:</span>
                        <q-btn-toggle
                          :model-value="membership.scope"
                          toggle-color="primary"
                          size="xs"
                          dense
                          :options="[
                            { label: t.labels.local?.value || 'Local', value: 'local' },
                            { label: t.labels.recursive?.value || 'Recursive', value: 'recursive' }
                          ]"
                          @update:model-value="updateMembershipScope(index, $event)"
                        />
                      </div>
                    </div>
                    <q-btn
                      icon="close"
                      flat
                      round
                      dense
                      size="sm"
                      color="grey-6"
                      @click="removeMembershipFromList(index)"
                    >
                      <AppTooltip :content="t.actions.remove?.value || 'Remove'" />
                    </q-btn>
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <!-- Empty State (only when not adding) -->
            <div
              v-else-if="!isAddingMembership"
              class="text-grey-6 text-body2 q-mb-sm q-pa-md bg-grey-2 rounded-borders text-center"
            >
              {{ t.placeholders.noMembershipsSelected?.value || 'No direct memberships configured' }}
            </div>

            <!-- INLINE: Add Membership Form -->
            <q-card v-if="isAddingMembership" flat bordered class="add-membership-form q-mb-sm">
              <q-card-section class="q-pa-md">
                <div class="row items-center q-mb-md">
                  <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
                  <div class="text-subtitle2 text-weight-medium">
                    {{ t.dialogs.addMembership?.title?.value || 'Add Direct Membership' }}
                  </div>
                  <q-space />
                  <q-btn
                    icon="close"
                    flat
                    round
                    dense
                    size="sm"
                    color="grey-6"
                    @click="cancelAddMembership"
                  >
                    <AppTooltip :content="t.actions.cancel?.value || 'Cancel'" />
                  </q-btn>
                </div>

                <!-- Organization (readonly) -->
                <div class="q-mb-md">
                  <div class="text-body2 text-weight-medium q-mb-xs">
                    {{ t.fields.organization?.value || 'Organization' }}
                  </div>
                  <div class="selector-field selector-field--readonly rounded-borders">
                    <div class="row items-center no-wrap q-pa-sm">
                      <q-icon name="business" color="primary" size="sm" class="q-mr-sm" />
                      <div class="col">
                        <div class="text-body2">{{ currentOrg.name }}</div>
                      </div>
                      <q-icon name="lock" size="xs" color="grey-5" />
                    </div>
                  </div>
                </div>

                <!-- Roles Selection -->
                <div class="q-mb-md">
                  <div class="text-body2 text-weight-medium q-mb-xs">
                    {{ t.fields.roles?.value || 'Roles' }} *
                  </div>
                  <div
                    class="selector-field rounded-borders cursor-pointer"
                    @click="showNewMembershipRolesDrawer = true"
                  >
                    <div class="row items-center no-wrap q-pa-sm">
                      <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
                      <div v-if="newMembershipRoles.length > 0" class="row items-center q-gutter-xs col q-py-xs">
                        <DetailChip
                          v-for="role in newMembershipRoles"
                          :key="role.id"
                          :label="role.name"
                          color="primary"
                          size="sm"
                        />
                      </div>
                      <div v-else class="text-grey-6 col">
                        {{ t.placeholders.selectRoles?.value || 'Select roles...' }}
                      </div>
                      <q-icon name="arrow_forward_ios" size="xs" color="grey-6" />
                    </div>
                  </div>
                </div>

                <!-- Scope Selection -->
                <div class="q-mb-md">
                  <div class="text-body2 text-weight-medium q-mb-sm">
                    {{ t.fields.scope?.value || 'Scope' }} *
                  </div>
                  <q-btn-toggle
                    v-model="newMembershipScope"
                    toggle-color="primary"
                    :options="[
                      { label: t.labels.local?.value || 'Local', value: 'local' },
                      { label: t.labels.recursive?.value || 'Recursive', value: 'recursive' }
                    ]"
                  />
                  <div class="text-caption text-grey-6 q-mt-xs">
                    {{ newMembershipScope === 'local'
                      ? (t.hints.scopeLocal?.value || 'Access applies to this organization only')
                      : (t.hints.scopeRecursive?.value || 'Access applies to this organization and all sub-organizations')
                    }}
                  </div>
                </div>

                <!-- Action Buttons -->
                <div class="row justify-end q-gutter-sm">
                  <q-btn
                    flat
                    :label="t.actions.cancel?.value || 'Cancel'"
                    @click="cancelAddMembership"
                  />
                  <q-btn
                    color="primary"
                    :label="t.actions.add?.value || 'Add'"
                    :disable="newMembershipRoles.length === 0"
                    @click="confirmAddMembership"
                  />
                </div>
              </q-card-section>
            </q-card>

            <!-- Add Membership Button (only when not adding) -->
            <q-btn
              v-if="!isAddingMembership"
              outline
              color="primary"
              icon="add"
              :label="t.actions.addMembership?.value || 'Add Membership'"
              class="full-width"
              @click="startAddMembership"
            />

            <div class="text-caption text-grey-6 q-mt-xs q-ml-sm">
              {{ t.hints.multipleMemberships?.value || 'Each membership grants roles with its own scope' }}
            </div>
          </div>
        </template>

        <!-- CREATE MODE: Single Role Selection + Scope -->
        <template v-else-if="hasOrgContext">
          <!-- Role Selection -->
          <div class="col-12">
            <div class="text-body2 text-weight-medium q-mb-xs">
              {{ t.fields.roles.value }} *
            </div>
            <div
              class="selector-field rounded-borders cursor-pointer"
              :class="{
                'has-error': (!localData.directMembership?.roleIds?.length) && showValidation
              }"
              @click="openRolesDrawer()"
            >
              <div class="row items-center no-wrap q-pa-sm">
                <q-icon name="admin_panel_settings" color="primary" size="sm" class="q-mr-sm" />
                <div v-if="localData.directMembership?.roleIds?.length" class="row items-center q-gutter-xs col q-py-xs">
                  <div
                    v-for="(roleName, index) in localData.directMembership.roleNames"
                    :key="index"
                    class="row items-center no-wrap"
                  >
                    <DetailChip
                      :label="roleName"
                      color="primary"
                      size="sm"
                    />
                    <q-btn
                      icon="close"
                      flat
                      round
                      dense
                      size="xs"
                      color="grey-6"
                      class="q-ml-xs"
                      @click.stop="removeRole(index)"
                    >
                      <AppTooltip :content="t.actions.remove?.value || 'Remove'" />
                    </q-btn>
                  </div>
                </div>
                <div v-else class="text-grey-6 col">
                  {{ t.placeholders.selectRoles.value }}
                </div>
                <q-btn
                  v-if="localData.directMembership?.roleIds?.length"
                  icon="close"
                  flat
                  round
                  dense
                  size="sm"
                  @click.stop="clearRoles"
                >
                  <AppTooltip :content="t.actions.clearAll?.value || 'Clear all'" />
                </q-btn>
                <q-icon name="arrow_forward_ios" size="xs" color="grey-6" />
              </div>
            </div>
          </div>

          <!-- Scope Selection -->
          <div class="col-12">
            <div class="text-body2 text-weight-medium q-mb-sm">
              {{ t.fields.scope.value }} *
            </div>
            <div class="row q-col-gutter-sm items-stretch">
              <div
                v-for="option in SCOPE_OPTIONS"
                :key="option.value"
                class="col-12 col-sm-6"
              >
                <q-card
                  flat
                  bordered
                  class="selectable-card cursor-pointer full-height"
                  :class="{ 'selectable-card--selected': localData.directMembership?.scope === option.value }"
                  @click="selectScope(option.value)"
                >
                  <q-card-section class="row items-center no-wrap q-pa-md full-height">
                    <q-icon
                      :name="option.icon"
                      size="md"
                      :color="localData.directMembership?.scope === option.value ? 'primary' : 'grey-6'"
                      class="q-mr-md"
                    />
                    <div class="col">
                      <div class="text-subtitle2 text-weight-medium">
                        {{ option.label }}
                      </div>
                      <div class="text-caption text-grey-7">
                        {{ option.description }}
                      </div>
                    </div>
                    <q-radio
                      :model-value="localData.directMembership?.scope"
                      :val="option.value"
                      color="primary"
                      @update:model-value="selectScope"
                    />
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </div>
        </template>
      </template>
    </div>

    <!-- Role Multi Selector Drawer (for Direct) -->
    <RoleMultiSelectorDrawer
      v-model="showRolesDrawer"
      :selected-role-ids="localData.directMembership?.roleIds || []"
      @confirm="handleRolesConfirm"
      @cancel="showRolesDrawer = false"
    />

    <!-- Role Multi Selector Drawer (for New Group) -->
    <RoleMultiSelectorDrawer
      v-model="showGroupRolesDrawer"
      :selected-role-ids="localData.selectedGroup?.newGroup?.roleIds || []"
      @confirm="handleGroupRolesConfirm"
      @cancel="showGroupRolesDrawer = false"
    />

    <!-- Group Selector Drawer -->
    <GroupSelectorDrawer
      v-model="showGroupDrawer"
      :selected-group-id="localData.selectedGroup?.existingGroup?.groupId ?? null"
      @select="handleGroupSelect"
      @cancel="showGroupDrawer = false"
    />

    <!-- Role Multi Selector Drawer (for Add Membership inline form) -->
    <RoleMultiSelectorDrawer
      v-model="showNewMembershipRolesDrawer"
      :selected-role-ids="newMembershipRoles.map(r => r.id)"
      @confirm="handleNewMembershipRolesConfirm"
      @cancel="showNewMembershipRolesDrawer = false"
    />
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step3Access',
});

/** TYPE IMPORTS */
import type { Step3AccessProps } from './interfaces/Step3Access.interface';
import type { QForm } from 'quasar';
import type {
  UserFormData,
  AccessType,
  ScopeType,
  GroupAccessMode,
  SelectedGroupData,
  DirectMembershipData,
} from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch, computed, onMounted } from 'vue';

/** COMPONENTS */
import { RoleMultiSelectorDrawer } from '@components/drawers/roles';
import { GroupSelectorDrawer } from '@components/drawers/groups';
import { DetailChip } from '@components/chips/DetailChip';
import { AppTooltip } from '@components/tooltips';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** COMPOSABLES */
import { useAddUserTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** LOCAL IMPORTS */
import { ACCESS_TYPE_OPTIONS, GROUP_ACCESS_MODE_OPTIONS, SCOPE_OPTIONS } from '../../constants';

const props = withDefaults(defineProps<Step3AccessProps>(), {
  isEditMode: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Pick<UserFormData, 'accessType' | 'selectedGroup' | 'directMembership' | 'selectedGroups' | 'directMemberships'>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddUserTranslations();
const logger = useLogger('Step3Access');
const orgStore = useOrganizationStore();

/** STATE */
const formRef = ref<QForm | null>(null);
const showValidation = ref(false);

// Drawer visibility state
const showRolesDrawer = ref(false);
const showGroupRolesDrawer = ref(false);
const showGroupDrawer = ref(false);

// New group form fields
const newGroupName = ref('');
const newGroupDescription = ref('');

// Add membership dialog state (for edit mode)
const isAddingMembership = ref(false);
const showNewMembershipRolesDrawer = ref(false);
const newMembershipRoles = ref<{ id: string; name: string }[]>([]);
const newMembershipScope = ref<ScopeType>('local');

/** COMPUTED */

/**
 * Current organization from store context
 */
const currentOrg = computed(() => ({
  id: orgStore.selectedOrganizationId,
  name: orgStore.selectedOrganizationName,
}));

/**
 * Check if org context is available (required for both group and direct assignment)
 */
const hasOrgContext = computed(() => !!orgStore.selectedOrganizationId);

/**
 * Whether component is in edit mode (multiple groups/memberships allowed)
 */
const isEditMode = computed(() => props.isEditMode);

// Local reactive state
const localData = reactive<{
  accessType: AccessType;
  groupAccessMode: GroupAccessMode;
  selectedGroup: SelectedGroupData | undefined;
  directMembership: DirectMembershipData;
  // Edit mode: support for multiple groups and memberships
  selectedGroups: SelectedGroupData[];
  directMemberships: DirectMembershipData[];
}>({
  accessType: props.modelValue.accessType || 'group',
  // Determine groupAccessMode: check arrays first (edit mode), then singular (create mode)
  groupAccessMode: props.modelValue.selectedGroups?.[0]?.mode || props.modelValue.selectedGroup?.mode || 'existing',
  selectedGroup: props.modelValue.selectedGroup ?? undefined,
  directMembership: props.modelValue.directMembership ?? {
    orgId: '',
    orgName: '',
    roleIds: [],
    roleNames: [],
    scope: 'local',
  },
  // Initialize from props if available (edit mode)
  selectedGroups: props.modelValue.selectedGroups ?? [],
  directMemberships: props.modelValue.directMemberships ?? [],
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.accessType = newVal.accessType || 'group';

    // Determine groupAccessMode: check arrays first (edit mode), then singular (create mode)
    const groupMode = newVal.selectedGroups?.[0]?.mode || newVal.selectedGroup?.mode || 'existing';
    localData.groupAccessMode = groupMode;

    localData.selectedGroup = newVal.selectedGroup ?? undefined;
    localData.directMembership = newVal.directMembership ?? {
      orgId: '',
      orgName: '',
      roleIds: [],
      roleNames: [],
      scope: 'local',
    };
    // Sync arrays for edit mode
    localData.selectedGroups = newVal.selectedGroups ?? [];
    localData.directMemberships = newVal.directMemberships ?? [];
    // Sync new group form fields
    if (newVal.selectedGroup?.newGroup) {
      newGroupName.value = newVal.selectedGroup.newGroup.name || '';
      newGroupDescription.value = newVal.selectedGroup.newGroup.description || '';
    }
  },
  { deep: true }
);

// Watch new group form fields and update localData
watch(
  [newGroupName, newGroupDescription],
  ([name, description]) => {
    if (localData.groupAccessMode === 'new' && localData.selectedGroup?.newGroup) {
      localData.selectedGroup.newGroup.name = name;
      if (description) {
        localData.selectedGroup.newGroup.description = description;
      }
      updateValue();
    }
  }
);

/** FUNCTIONS */

/**
 * Select access type (group or direct)
 */
function selectAccessType(type: AccessType): void {
  localData.accessType = type;

  // Clear the other type's data when switching
  if (type === 'group') {
    localData.directMembership = {
      orgId: '',
      orgName: '',
      roleIds: [],
      roleNames: [],
      scope: 'local',
    };
  } else {
    localData.selectedGroup = undefined;
    // Set org from current context
    if (currentOrg.value.id && currentOrg.value.name) {
      localData.directMembership = {
        orgId: currentOrg.value.id,
        orgName: currentOrg.value.name,
        roleIds: [],
        roleNames: [],
        scope: 'local',
      };
    }
  }

  updateValue();
  logger.debug('Access type selected:', type);
}

/**
 * Select group access mode (existing or new)
 */
function selectGroupAccessMode(mode: GroupAccessMode): void {
  localData.groupAccessMode = mode;

  // Reset selectedGroup based on mode
  if (mode === 'existing') {
    localData.selectedGroup = {
      mode: 'existing',
      existingGroup: undefined,
      newGroup: undefined,
    };
    // Clear new group fields
    newGroupName.value = '';
    newGroupDescription.value = '';
  } else {
    localData.selectedGroup = {
      mode: 'new',
      existingGroup: undefined,
      newGroup: {
        name: '',
        roleIds: [],
        roleNames: [],
      },
    };
  }

  updateValue();
  logger.debug('Group access mode selected:', mode);
}

/**
 * Open group drawer
 */
function openGroupDrawer(): void {
  showGroupDrawer.value = true;
}

/**
 * Handle group selection from drawer (for existing group mode)
 * In edit mode: adds to array
 * In create mode: replaces single selection
 */
function handleGroupSelect(group: any): void {
  if (props.isEditMode) {
    // Edit mode: add to array
    addGroupToList(group);
  } else {
    // Create mode: single selection
    localData.selectedGroup = {
      mode: 'existing',
      existingGroup: {
        groupId: group.id,
        groupName: group.name,
      },
      newGroup: undefined,
    };
    updateValue();
  }
  logger.debug('Group selected:', { name: group.name, editMode: props.isEditMode });
}

/**
 * Open group roles drawer (for new group mode)
 */
function openGroupRolesDrawer(): void {
  showGroupRolesDrawer.value = true;
}

/**
 * Handle group roles confirmation from drawer (for new group mode)
 */
function handleGroupRolesConfirm(roles: any[]): void {
  if (localData.selectedGroup?.newGroup) {
    localData.selectedGroup.newGroup.roleIds = roles.map(r => r.id);
    localData.selectedGroup.newGroup.roleNames = roles.map(r => r.name);
  }
  updateValue();
  logger.debug('New group roles selected:', roles.length);
}

/**
 * Remove a single role from new group selection
 */
function removeNewGroupRole(index: number): void {
  if (localData.selectedGroup?.newGroup) {
    localData.selectedGroup.newGroup.roleIds.splice(index, 1);
    localData.selectedGroup.newGroup.roleNames.splice(index, 1);
    updateValue();
  }
}

/**
 * Clear all selected roles for new group
 */
function clearNewGroupRoles(): void {
  if (localData.selectedGroup?.newGroup) {
    localData.selectedGroup.newGroup.roleIds = [];
    localData.selectedGroup.newGroup.roleNames = [];
    updateValue();
  }
}

/**
 * Clear selected group (for existing group mode)
 */
function clearGroup(): void {
  if (localData.selectedGroup) {
    localData.selectedGroup.existingGroup = undefined;
  }
  updateValue();
}

/**
 * Open roles drawer
 */
function openRolesDrawer(): void {
  showRolesDrawer.value = true;
}

/**
 * Handle roles confirmation from drawer
 */
function handleRolesConfirm(roles: any[]): void {
  if (localData.directMembership) {
    localData.directMembership.roleIds = roles.map(r => r.id);
    localData.directMembership.roleNames = roles.map(r => r.name);
  }
  updateValue();
  logger.debug('Roles selected:', roles.length);
}

/**
 * Remove a single role from selection
 */
function removeRole(index: number): void {
  if (localData.directMembership) {
    localData.directMembership.roleIds.splice(index, 1);
    localData.directMembership.roleNames.splice(index, 1);
    updateValue();
  }
}

/**
 * Clear all selected roles
 */
function clearRoles(): void {
  if (localData.directMembership) {
    localData.directMembership.roleIds = [];
    localData.directMembership.roleNames = [];
    updateValue();
  }
}

/**
 * Select scope for direct membership
 */
function selectScope(scope: ScopeType): void {
  if (localData.directMembership) {
    localData.directMembership.scope = scope;
    updateValue();
  }
}

// ============================================
// EDIT MODE: Multiple Groups Management
// ============================================

/**
 * Add an existing group to selectedGroups array
 * @param {any} group - Group object from selector
 */
function addGroupToList(group: any): void {
  // Check if already in list
  const exists = localData.selectedGroups.some(
    g => g.existingGroup?.groupId === group.id
  );

  if (!exists) {
    localData.selectedGroups.push({
      mode: 'existing',
      existingGroup: {
        groupId: group.id,
        groupName: group.name,
      },
    });
    updateValue();
    logger.debug('Group added to list:', group.name);
  }
}

/**
 * Remove a group from selectedGroups array
 * @param {number} index - Index of group to remove
 */
function removeGroupFromList(index: number): void {
  localData.selectedGroups.splice(index, 1);
  updateValue();
  logger.debug('Group removed from list, remaining:', localData.selectedGroups.length);
}

/**
 * Add a direct membership to directMemberships array
 * @param {DirectMembershipData} membership - Membership data
 */
function addMembershipToList(membership: DirectMembershipData): void {
  localData.directMemberships.push(membership);
  updateValue();
  logger.debug('Membership added to list');
}

/**
 * Remove a membership from directMemberships array
 * @param {number} index - Index of membership to remove
 */
function removeMembershipFromList(index: number): void {
  localData.directMemberships.splice(index, 1);
  updateValue();
  logger.debug('Membership removed from list, remaining:', localData.directMemberships.length);
}

/**
 * Update scope for a specific membership
 * @param {number} index - Index of membership
 * @param {ScopeType} scope - New scope value
 */
function updateMembershipScope(index: number, scope: ScopeType): void {
  if (localData.directMemberships[index]) {
    localData.directMemberships[index].scope = scope;
    updateValue();
  }
}

// ============================================
// EDIT MODE: Add Membership Dialog
// ============================================

/**
 * Start adding a new membership (show inline form)
 */
function startAddMembership(): void {
  newMembershipRoles.value = [];
  newMembershipScope.value = 'local';
  isAddingMembership.value = true;
}

/**
 * Handle roles selection for new membership
 * @param {any[]} roles - Selected roles
 */
function handleNewMembershipRolesConfirm(roles: any[]): void {
  newMembershipRoles.value = roles.map(r => ({ id: r.id, name: r.name }));
}

/**
 * Confirm and add the new membership
 */
function confirmAddMembership(): void {
  if (newMembershipRoles.value.length === 0) return;

  const newMembership: DirectMembershipData = {
    orgId: currentOrg.value.id || '',
    orgName: currentOrg.value.name || '',
    roleIds: newMembershipRoles.value.map(r => r.id),
    roleNames: newMembershipRoles.value.map(r => r.name),
    scope: newMembershipScope.value,
  };

  addMembershipToList(newMembership);
  isAddingMembership.value = false;
  logger.debug('New membership added with roles:', newMembershipRoles.value.length);
}

/**
 * Cancel adding membership
 */
function cancelAddMembership(): void {
  isAddingMembership.value = false;
  newMembershipRoles.value = [];
  newMembershipScope.value = 'local';
}

/**
 * Emit updated value to parent
 * Supports both legacy single values and new array format
 */
function updateValue(): void {
  const value: Partial<Pick<UserFormData, 'accessType' | 'selectedGroup' | 'directMembership' | 'selectedGroups' | 'directMemberships'>> = {
    accessType: localData.accessType,
  };

  // Legacy single values (for create mode compatibility)
  if (localData.accessType === 'group' || localData.accessType === 'both') {
    value.selectedGroup = localData.selectedGroup;
  }
  if (localData.accessType === 'direct' || localData.accessType === 'both') {
    value.directMembership = localData.directMembership;
  }

  // New array format (for edit mode with multiple)
  if (localData.selectedGroups.length > 0) {
    value.selectedGroups = localData.selectedGroups;
  }
  if (localData.directMemberships.length > 0) {
    value.directMemberships = localData.directMemberships;
  }

  emit('update:modelValue', value as Pick<UserFormData, 'accessType' | 'selectedGroup' | 'directMembership' | 'selectedGroups' | 'directMemberships'>);
}

/**
 * Validate form
 * Supports both create mode (single) and edit mode (multiple)
 */
function validate(): boolean {
  showValidation.value = true;

  // Both access types require org context
  if (!hasOrgContext.value) return false;

  // Helper to validate groups
  const validateGroups = (): boolean => {
    // Check array format first (edit mode)
    if (localData.selectedGroups.length > 0) {
      return localData.selectedGroups.every(g => {
        if (g.mode === 'existing') return !!g.existingGroup?.groupId;
        if (g.mode === 'new') return !!(g.newGroup?.name && g.newGroup.name.length >= 3 && g.newGroup.roleIds?.length);
        return false;
      });
    }
    // Fallback to single (create mode)
    if (localData.groupAccessMode === 'existing') {
      return !!(localData.selectedGroup?.existingGroup?.groupId);
    } else {
      return !!(
        localData.selectedGroup?.newGroup?.name &&
        localData.selectedGroup?.newGroup?.name.length >= 3 &&
        localData.selectedGroup?.newGroup?.roleIds?.length
      );
    }
  };

  // Helper to validate memberships
  const validateMemberships = (): boolean => {
    // Check array format first (edit mode)
    if (localData.directMemberships.length > 0) {
      return localData.directMemberships.every(m => m.roleIds && m.roleIds.length > 0);
    }
    // Fallback to single (create mode)
    return !!(
      localData.directMembership?.orgId &&
      localData.directMembership?.roleIds?.length &&
      localData.directMembership?.scope
    );
  };

  if (localData.accessType === 'group') {
    return validateGroups();
  } else if (localData.accessType === 'direct') {
    return validateMemberships();
  } else if (localData.accessType === 'both') {
    // Both requires at least one valid group OR one valid membership
    return validateGroups() || validateMemberships();
  }

  return false;
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  logger.debug('Step3Access mounted');
});

/** EXPOSE */
defineExpose({
  formRef,
  validate,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

// Utility class for equal height cards in flex rows
.full-height {
  height: 100%;
}

.selector-field {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  min-height: 40px;
  transition: var(--mapex-transition-base);

  &:hover:not(.disabled):not(.selector-field--readonly) {
    border-color: var(--q-primary);
  }

  &.disabled {
    background-color: var(--mapex-submenu-bg);
    cursor: not-allowed;
  }

  &--readonly {
    background-color: var(--mapex-submenu-bg);
    border-style: dashed;
    cursor: default;
  }

  &.has-error {
    border-color: var(--q-negative);
  }
}

// Unified selectable card pattern for all selection cards
// Used for: Access Type, Group Access Mode, Scope Selection
.selectable-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);

  &:hover {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.05);
  }

  &--selected {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.08);
  }

  &--warning {
    border-color: var(--q-warning);
    background-color: rgba(var(--q-warning-rgb), 0.08);
  }
}

// Edit mode: Selected groups list
.selected-groups-list {
  .selected-group-card {
    border-radius: var(--mapex-radius-md);
    transition: var(--mapex-transition-base);

    &:hover {
      border-color: var(--q-secondary);
      background-color: rgba(var(--q-secondary-rgb), 0.05);
    }
  }
}

// Edit mode: Memberships list
.memberships-list {
  .membership-card {
    border-radius: var(--mapex-radius-md);
    transition: var(--mapex-transition-base);

    &:hover {
      border-color: var(--q-primary);
      background-color: rgba(var(--q-primary-rgb), 0.05);
    }
  }
}

// Edit mode: Add membership inline form
.add-membership-form {
  border-radius: var(--mapex-radius-md);
  border-color: var(--q-primary);
  background-color: rgba(var(--q-primary-rgb), 0.02);
}
</style>

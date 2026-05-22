<script setup lang="ts">
/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type { ExternalVariableForm } from './interfaces/ExternalVariables.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { FontIconDialog } from '@components/dialogs/fontIcons';
import { AssetTemplateSelectorDialog } from '@components/dialogs/common/assetTemplateSelectorDialog';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import {
  EXTERNAL_INPUT_TYPE_OPTIONS,
  DEFAULT_VALUE_BY_TYPE,
  ASSET_FIELD_PATH_OPTIONS,
} from '../../constants';

/** COMPOSABLES & STORES */
const {
  externalInputs,
  addExternalInput,
  updateExternalInput,
  removeExternalInput,
  moveExternalInput,
} = useWorkflowEditorState();
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Currently editing variable index (-1 = adding new)
 */
const editingIndex = ref(-1);

/**
 * Form data for add/edit
 */
const form = ref<ExternalVariableForm>({
  field: '',
  label: '',
  icon: 'input',
  type: 'string',
  description: '',
  defaultValue: '',
  required: false,
});

/**
 * Whether the manual icon text input is active (vs FontIconDialog picker)
 */
const manualIconMode = ref(false);

/**
 * Whether the FontIconDialog is open
 */
const showIconPicker = ref(false);

/**
 * Whether the AssetTemplateSelectorDialog is open
 */
const showTemplateSelector = ref(false);

/**
 * Display name of the selected asset template
 */
const selectedTemplateName = ref('');

/** COMPUTED */

/**
 * Whether we are in edit mode (vs add mode)
 */
const isEditing = computed(() => editingIndex.value >= 0);

/**
 * Whether the current type is assetFromTemplate
 */
const isAssetType = computed(() => form.value.type === 'assetFromTemplate');

/**
 * Whether the current type is literal (fixed value)
 */
const isLiteralType = computed(() => form.value.type === 'literal');

/**
 * Whether the current type needs a default value input
 */
const showDefaultValue = computed(() => !isAssetType.value && !isLiteralType.value);

/**
 * Whether the form has the minimum required fields filled
 */
const canSubmit = computed(() => {
  if (!form.value.field.trim() || !form.value.label.trim()) return false;
  if (isAssetType.value) {
    return !!form.value.assetTemplateId && !!form.value.fieldPath;
  }
  if (isLiteralType.value) {
    return form.value.defaultValue !== '' && form.value.defaultValue !== undefined;
  }
  return true;
});

/** FUNCTIONS */

/**
 * Reset form to defaults
 *
 * @returns {void}
 */
function resetForm(): void {
  form.value = {
    field: '',
    label: '',
    icon: 'input',
    type: 'string',
    description: '',
    defaultValue: '',
    required: false,
  };
  selectedTemplateName.value = '';
  editingIndex.value = -1;
}

/**
 * Start editing a variable.
 * For assetFromTemplate type, resolves the template name from API.
 *
 * @param {number} index - Variable index to edit
 * @returns {void}
 */
function startEdit(index: number): void {
  editingIndex.value = index;
  form.value = { ...externalInputs.value[index] } as ExternalVariableForm;

  // Resolve asset template name if editing an assetFromTemplate field
  if (form.value.type === 'assetFromTemplate' && form.value.assetTemplateId) {
    void resolveTemplateName(form.value.assetTemplateId);
  } else {
    selectedTemplateName.value = '';
  }
}

/**
 * Resolve asset template name by ID from API
 *
 * @param {string} templateId - Asset template ID
 * @returns {Promise<void>}
 */
async function resolveTemplateName(templateId: string): Promise<void> {
  try {
    const template = await apis.assets.assetTemplate.getById({
      assetTemplateId: templateId,
    });
    selectedTemplateName.value = template?.name || templateId;
  } catch {
    selectedTemplateName.value = templateId;
  }
}

/**
 * Handle form submission (add or update)
 *
 * @returns {void}
 */
function handleSubmit(): void {
  if (!canSubmit.value) return;

  if (isEditing.value) {
    updateExternalInput(editingIndex.value, { ...form.value });
  } else {
    addExternalInput({ ...form.value });
  }

  resetForm();
}

/**
 * Handle type change — reset default value and clear asset fields
 *
 * @param {string} newType - New variable type
 * @returns {void}
 */
function handleTypeChange(newType: string): void {
  if (newType === 'assetFromTemplate') {
    form.value.defaultValue = '';
    form.value.icon = 'devices';
  } else if (newType === 'literal') {
    form.value.defaultValue = '';
    form.value.icon = 'pin';
    form.value.required = false;
    delete form.value.assetTemplateId;
    delete form.value.fieldPath;
  } else {
    form.value.defaultValue = DEFAULT_VALUE_BY_TYPE[newType] ?? '';
    delete form.value.assetTemplateId;
    delete form.value.fieldPath;
  }
}

/**
 * Handle delete with confirmation
 *
 * @param {number} index - Variable index to delete
 * @returns {void}
 */
function handleDelete(index: number): void {
  removeExternalInput(index);
  if (editingIndex.value === index) {
    resetForm();
  }
}

/**
 * Handle asset template selection from modal
 *
 * @param {import('@mapexos/schemas').AssetTemplateResponse[]} templates - Selected templates
 * @returns {void}
 */
function handleTemplateSelect(templates: AssetTemplateResponse[]): void {
  const first = templates[0];
  if (first?.id) {
    form.value.assetTemplateId = first.id;
    selectedTemplateName.value = first.name ?? '';
  }
}

/**
 * Clear selected asset template
 *
 * @returns {void}
 */
function clearTemplate(): void {
  delete form.value.assetTemplateId;
  selectedTemplateName.value = '';
}

/**
 * Get display label for an asset field path value
 *
 * @param {string} fieldPath - The fieldPath value
 * @returns {string} Human-readable label
 */
function getFieldPathLabel(fieldPath: string): string {
  return ASSET_FIELD_PATH_OPTIONS.find(o => o.value === fieldPath)?.label ?? fieldPath;
}
</script>

<template>
  <div class="row q-col-gutter-md">
    <!-- Sidebar: Add/Edit Form -->
    <div class="col-12 col-md-4">
      <q-card flat bordered class="sticky-sidebar">
        <q-card-section>
          <div class="text-subtitle2 text-weight-medium q-mb-md">
            {{ isEditing ? t.externalInputs.editTitle.value : t.externalInputs.addTitle.value }}
          </div>

          <!-- ── Section: Identity ── -->
          <div class="external-variables__section-label">
            <q-icon name="badge" size="14px" class="q-mr-xs" />
            {{ t.externalInputs.sectionIdentity.value }}
          </div>

          <q-input
            v-model="form.field"
            outlined
            dense
            class="q-mb-md"
            :label="t.externalInputs.field.value"
            :hint="t.externalInputs.fieldHint.value"
            :rules="[(val: string) => !!val || t.validation.fieldNameRequired.value]"
          />

          <q-input
            v-model="form.label"
            outlined
            dense
            class="q-mb-md"
            :label="t.externalInputs.label.value"
            :hint="t.externalInputs.labelHint.value"
            :rules="[(val: string) => !!val || t.validation.nameIsRequired.value]"
          />

          <!-- Icon field -->
          <q-input
            v-if="manualIconMode"
            v-model="form.icon"
            outlined
            dense
            class="q-mb-sm"
            :label="t.externalInputs.icon.value"
            placeholder="e.g. shopping_cart"
          >
            <template #prepend>
              <q-icon :name="form.icon || 'help_outline'" size="20px" />
            </template>
          </q-input>
          <q-input
            v-else
            outlined
            dense
            readonly
            class="q-mb-sm cursor-pointer"
            :model-value="form.icon"
            :label="t.externalInputs.icon.value"
            @click="showIconPicker = true"
          >
            <template #prepend>
              <q-icon :name="form.icon || 'help_outline'" size="20px" color="primary" />
            </template>
            <template #append>
              <q-btn flat dense no-caps size="sm" icon="palette" color="primary" :label="t.externalInputs.pickIcon.value" @click.stop="showIconPicker = true" />
            </template>
          </q-input>
          <q-checkbox
            v-model="manualIconMode"
            dense
            size="xs"
            class="q-mb-sm"
            :label="t.externalInputs.manualIcon.value"
          />

          <!-- FontIconDialog -->
          <FontIconDialog
            v-model="form.icon"
            v-model:show="showIconPicker"
          />

          <!-- ── Section: Type Configuration ── -->
          <q-separator class="q-my-md" />
          <div class="external-variables__section-label">
            <q-icon name="tune" size="14px" class="q-mr-xs" />
            {{ t.externalInputs.sectionTypeConfig.value }}
          </div>

          <q-select
            v-model="form.type"
            outlined
            dense
            emit-value
            map-options
            class="q-mb-md"
            :label="t.externalInputs.type.value"
            :options="[...EXTERNAL_INPUT_TYPE_OPTIONS]"
            @update:model-value="handleTypeChange"
          />

          <!-- Literal type: fixed value input -->
          <template v-if="isLiteralType">
            <q-input
              v-model="form.defaultValue"
              outlined
              dense
              class="q-mb-md"
              :label="t.externalInputs.literalValue.value"
              :hint="t.externalInputs.literalValueHint.value"
            >
              <template #prepend>
                <q-icon name="pin" size="20px" color="deep-purple-4" />
              </template>
            </q-input>

            <div class="external-variables__info q-mb-md">
              <q-icon name="info" color="grey-6" size="xs" class="q-mr-sm" />
              <span>{{ t.externalInputs.literalInfo.value }}</span>
            </div>
          </template>

          <!-- Asset from Template: template selector + field path -->
          <template v-if="isAssetType">
            <!-- Asset Template selector (opens modal) -->
            <q-input
              outlined
              dense
              readonly
              class="q-mb-md cursor-pointer"
              :model-value="selectedTemplateName || form.assetTemplateId || ''"
              :label="t.externalInputs.assetTemplate.value"
              :hint="t.externalInputs.selectAssetTemplate.value"
              @click="showTemplateSelector = true"
            >
              <template #prepend>
                <q-icon name="memory" size="20px" color="teal-6" />
              </template>
              <template #append>
                <q-btn
                  v-if="form.assetTemplateId"
                  flat
                  dense
                  round
                  icon="close"
                  size="xs"
                  color="grey-6"
                  @click.stop="clearTemplate"
                />
                <q-btn
                  flat
                  dense
                  round
                  icon="search"
                  size="sm"
                  color="primary"
                  @click.stop="showTemplateSelector = true"
                />
              </template>
            </q-input>

            <!-- Asset Template Selector Modal -->
            <AssetTemplateSelectorDialog
              v-model="showTemplateSelector"
              :multi-select="false"
              :selected-template-ids="form.assetTemplateId ? [form.assetTemplateId] : []"
              @select="handleTemplateSelect"
            />

            <!-- Field Path selector -->
            <q-select
              v-model="form.fieldPath"
              outlined
              dense
              emit-value
              map-options
              class="q-mb-md"
              :label="t.externalInputs.fieldPath.value"
              :hint="t.externalInputs.fieldPathHint.value"
              :options="[...ASSET_FIELD_PATH_OPTIONS]"
              option-value="value"
              option-label="label"
            >
              <template #prepend>
                <q-icon name="data_object" size="20px" color="teal-6" />
              </template>
              <template #option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section side>
                    <q-icon :name="scope.opt.icon" size="18px" color="grey-7" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ scope.opt.label }}</q-item-label>
                    <q-item-label caption class="external-variables__field-path-caption">{{ scope.opt.value }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-select>

            <!-- Info banner -->
            <div class="external-variables__info q-mb-md">
              <q-icon name="info" color="grey-6" size="xs" class="q-mr-sm" />
              <span>{{ t.externalInputs.assetFromTemplateInfo.value }}</span>
            </div>
          </template>

          <!-- Dynamic default value input (hidden for asset and literal types) -->
          <template v-if="showDefaultValue">
            <q-input
              v-if="form.type === 'string'"
              v-model="form.defaultValue"
              outlined
              dense
              class="q-mb-md"
              :label="t.externalInputs.defaultValue.value"
            />
            <q-input
              v-else-if="form.type === 'number'"
              v-model.number="form.defaultValue"
              outlined
              dense
              type="number"
              class="q-mb-md"
              :label="t.externalInputs.defaultValue.value"
            />
            <q-toggle
              v-else-if="form.type === 'boolean'"
              v-model="form.defaultValue"
              class="q-mb-md"
              :label="t.externalInputs.defaultValue.value"
            />
            <q-input
              v-else
              v-model="form.defaultValue"
              outlined
              dense
              type="textarea"
              autogrow
              class="q-mb-md"
              :label="t.externalInputs.defaultValueJson.value"
            />
          </template>

          <!-- ── Section: Metadata ── -->
          <q-separator class="q-my-md" />
          <div class="external-variables__section-label">
            <q-icon name="info_outline" size="14px" class="q-mr-xs" />
            {{ t.externalInputs.sectionMetadata.value }}
          </div>

          <q-input
            v-model="form.description"
            outlined
            dense
            class="q-mb-md"
            :label="t.externalInputs.description.value"
            :hint="t.externalInputs.descriptionHint.value"
          />

          <!-- Required checkbox (hidden for literal — always "provided") -->
          <q-checkbox
            v-if="!isLiteralType"
            v-model="form.required"
            dense
            class="q-mb-md"
            :label="t.externalInputs.required.value"
          >
            <AppTooltip :content="t.externalInputs.requiredHint.value" />
          </q-checkbox>

          <div class="row q-gutter-sm">
            <q-btn
              unelevated
              no-caps
              :color="isEditing ? 'amber-8' : 'primary'"
              :label="isEditing ? t.externalInputs.update.value : t.externalInputs.add.value"
              :icon="isEditing ? 'save' : 'add'"
              :disable="!canSubmit"
              @click="handleSubmit"
            />
            <q-btn
              v-if="isEditing"
              flat
              no-caps
              color="grey-7"
              :label="t.externalInputs.cancel.value"
              @click="resetForm"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- External Variables List -->
    <div class="col-12 col-md-8">
      <!-- Empty state -->
      <div v-if="externalInputs.length === 0" class="empty-state">
        <q-icon name="input" size="48px" class="q-mb-md" />
        <p>{{ t.externalInputs.emptyTitle.value }}</p>
        <p class="text-caption">{{ t.externalInputs.emptyDescription.value }}</p>

        <div class="external-variables__event-hint q-mt-lg">
          <q-icon name="event" color="blue-6" size="16px" class="q-mr-sm q-mt-xs" />
          <div>
            <div class="text-weight-medium q-mb-xs">{{ t.externalInputs.eventFieldsBannerTitle.value }}</div>
            <div class="text-caption">{{ t.externalInputs.eventFieldsBannerDescription.value }}</div>
          </div>
        </div>
      </div>

      <!-- List -->
      <q-list v-else separator bordered class="rounded-borders">
        <q-item v-for="(variable, index) in externalInputs" :key="index">
          <q-item-section avatar>
            <q-icon :name="variable.icon" color="cyan-6" size="24px" />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ variable.label }}
            </q-item-label>
            <q-item-label caption>
              <span class="variable-field-key">input.{{ variable.field }}</span>
              <q-badge
                :label="variable.type === 'assetFromTemplate' ? 'asset' : variable.type"
                :color="variable.type === 'assetFromTemplate' ? 'teal-7' : variable.type === 'literal' ? 'deep-purple-4' : 'grey-7'"
                class="q-ml-xs q-mr-xs"
              />
              <q-badge
                v-if="variable.type !== 'literal'"
                :label="variable.required ? t.externalInputs.requiredBadge.value : t.externalInputs.optionalBadge.value"
                :color="variable.required ? 'cyan-8' : 'grey-5'"
              />
            </q-item-label>
            <!-- Literal value line -->
            <q-item-label v-if="variable.type === 'literal'" caption class="q-mt-xs">
              <span class="external-variables__asset-detail">
                <q-icon name="pin" size="12px" class="q-mr-xs" />
                {{ variable.defaultValue }}
              </span>
            </q-item-label>
            <!-- Asset details line -->
            <q-item-label v-if="variable.type === 'assetFromTemplate' && variable.fieldPath" caption class="q-mt-xs">
              <span class="external-variables__asset-detail">
                <q-icon name="data_object" size="12px" class="q-mr-xs" />
                {{ getFieldPathLabel(variable.fieldPath) }}
              </span>
            </q-item-label>
            <q-item-label v-if="variable.description" caption class="q-mt-xs text-grey-6">
              {{ variable.description }}
            </q-item-label>
          </q-item-section>
          <q-item-section side>
            <div class="row q-gutter-xs">
              <q-btn flat dense round icon="arrow_upward" size="sm" :disable="index === 0" @click="moveExternalInput(index, 'up')" />
              <q-btn flat dense round icon="arrow_downward" size="sm" :disable="index === externalInputs.length - 1" @click="moveExternalInput(index, 'down')" />
              <q-btn flat dense round icon="edit" size="sm" color="primary" @click="startEdit(index)" />
              <q-btn flat dense round icon="delete" size="sm" color="negative" @click="handleDelete(index)" />
            </div>
          </q-item-section>
        </q-item>
      </q-list>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sticky-sidebar {
  position: sticky;
  top: 80px;
}

.variable-field-key {
  font-family: monospace;
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  height: 100%;
  text-align: center;
  color: var(--mapex-text-secondary);
}

.external-variables {
  &__section-label {
    display: flex;
    align-items: center;
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--mapex-text-muted);
    margin-bottom: 12px;
  }

  &__info {
    display: flex;
    align-items: flex-start;
    padding: 10px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2, var(--mapex-submenu-bg));
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    line-height: 1.4;
  }

  &__event-hint {
    display: flex;
    align-items: flex-start;
    max-width: 480px;
    padding: 12px 16px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2, var(--mapex-submenu-bg));
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    line-height: 1.5;
    text-align: left;
  }

  &__field-path-caption {
    font-family: monospace;
    font-size: 0.7rem;
    color: var(--mapex-text-muted);
  }

  &__asset-detail {
    display: inline-flex;
    align-items: center;
    font-size: 0.7rem;
    color: var(--mapex-text-muted);
  }
}
</style>

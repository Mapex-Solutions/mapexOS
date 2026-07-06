<script setup lang="ts">
defineOptions({
  name: 'AttributesSection'
});

/** TYPE IMPORTS */
import type { AttributesSectionProps, AttributesSectionEmits } from './interfaces';
import type { AssetAttributeForm, AssetAttributeKind, AssetAttributeGeoValue } from '../../interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ATTRIBUTE_KIND_VALUES, defaultValueForKind } from './constants';

/** PROPS & EMITS */
const props = defineProps<AttributesSectionProps>();
const emit = defineEmits<AttributesSectionEmits>();

/** COMPOSABLES & STORES */
const ts = useTS({ capitalize: false });
const bp = 'pages.assets.addAsset.steps.step1.attributes';

/** STATE */
// The committed attributes (a local copy; the prop is never mutated in place).
const rows = ref<AssetAttributeForm[]>(clone(props.modelValue));

// The attribute currently being composed in the static add form.
const newAttr = ref<AssetAttributeForm>(emptyAttribute());

/** COMPUTED */
const kindOptions = computed(() =>
  ATTRIBUTE_KIND_VALUES.map((value) => ({ value, label: ts(`${bp}.kinds.${value}`) }))
);

// An attribute can be committed once it has a label; the value stays optional.
const canAdd = computed(() => newAttr.value.label.trim().length > 0);

/** WATCHERS */
watch(
  () => props.modelValue,
  (value) => {
    if (JSON.stringify(value) !== JSON.stringify(rows.value)) {
      rows.value = clone(value);
    }
  }
);

/** FUNCTIONS */

/**
 * Builds a fresh, empty attribute for the add form (searchable captured as true).
 * @returns {AssetAttributeForm} A blank string attribute.
 */
function emptyAttribute(): AssetAttributeForm {
  return { label: '', kind: 'string', value: defaultValueForKind('string'), searchable: true };
}

/**
 * Deep-clones an attribute list so local edits never touch the prop.
 * @param {AssetAttributeForm[]} list - The list to clone.
 * @returns {AssetAttributeForm[]} A structural copy.
 */
function clone(list: AssetAttributeForm[]): AssetAttributeForm[] {
  return JSON.parse(JSON.stringify(list ?? []));
}

/**
 * Emits the committed rows as a fresh array to the parent v-model.
 */
function emitRows(): void {
  emit('update:modelValue', clone(rows.value));
}

/**
 * Commits the composed attribute to the list below and resets the add form.
 */
function addAttribute(): void {
  if (!canAdd.value) return;
  rows.value.push(JSON.parse(JSON.stringify(newAttr.value)));
  emitRows();
  newAttr.value = emptyAttribute();
}

/**
 * Removes the committed attribute at the given index.
 * @param {number} index - Row index to remove.
 */
function removeRow(index: number): void {
  rows.value.splice(index, 1);
  emitRows();
}

/**
 * Switches the composing attribute's kind and resets its value to that kind's zero.
 * @param {AssetAttributeKind} kind - The newly selected kind.
 */
function onKindChange(kind: AssetAttributeKind): void {
  newAttr.value.kind = kind;
  newAttr.value.value = defaultValueForKind(kind);
}

/**
 * The composing attribute's value as a number for the integer input (null when unset).
 * @returns {number | null} The numeric value.
 */
function numValue(): number | null {
  return typeof newAttr.value.value === 'number' ? newAttr.value.value : null;
}

/**
 * The composing attribute's value as a geo point.
 * @returns {AssetAttributeGeoValue} The geo value.
 */
function geoValue(): AssetAttributeGeoValue {
  return newAttr.value.value as AssetAttributeGeoValue;
}

/**
 * Sets the composing attribute's value verbatim (string/date/boolean inputs).
 * @param {AssetAttributeForm['value']} value - The new value.
 */
function setValue(value: AssetAttributeForm['value']): void {
  newAttr.value.value = value;
}

/**
 * Sets the composing integer value, coercing empty input to null.
 * @param {string | number | null} raw - The raw input value.
 */
function setIntegerValue(raw: string | number | null): void {
  newAttr.value.value = raw === '' || raw === null ? null : Number(raw);
}

/**
 * Sets a geo coordinate on the composing attribute, coercing empty input to null.
 * @param {'lat' | 'lon'} axis - Coordinate to set.
 * @param {string | number | null} raw - The raw input value.
 */
function setGeoValue(axis: 'lat' | 'lon', raw: string | number | null): void {
  geoValue()[axis] = raw === '' || raw === null ? null : Number(raw);
}

/**
 * Localized label for a kind, shown in the committed-attributes list.
 * @param {AssetAttributeKind} kind - The attribute kind.
 * @returns {string} Localized kind label.
 */
function kindLabel(kind: AssetAttributeKind): string {
  return ts(`${bp}.kinds.${kind}`);
}

/**
 * Formats a committed attribute's value for display in the list.
 * @param {AssetAttributeForm} row - The attribute row.
 * @returns {string} A readable value.
 */
function displayValue(row: AssetAttributeForm): string {
  if (row.kind === 'geo') {
    const geo = row.value as AssetAttributeGeoValue;
    return `${geo.lat ?? '—'}, ${geo.lon ?? '—'}`;
  }
  const value = row.value as string | number | boolean | null;
  if (row.kind === 'boolean') return String(value);
  if (value === '' || value === null || value === undefined) return '—';
  return String(value);
}
</script>

<template>
  <div class="attributes-section">
    <div class="attributes-section__header">
      <div class="attributes-section__title">{{ ts(`${bp}.title`) }}</div>
      <div class="attributes-section__subtitle">{{ ts(`${bp}.subtitle`) }}</div>
    </div>

    <!-- Static add form: compose one attribute, then commit it to the list below -->
    <q-card flat bordered class="attributes-section__add q-mb-md">
      <q-card-section class="row q-col-gutter-sm items-start">
        <div class="col-12 col-sm-3">
          <q-input
            v-model="newAttr.label"
            outlined
            dense
            :label="ts(`${bp}.columns.label`)"
            :placeholder="ts(`${bp}.columns.labelPlaceholder`)"
          />
        </div>

        <div class="col-12 col-sm-3">
          <q-select
            :model-value="newAttr.kind"
            outlined
            dense
            emit-value
            map-options
            :options="kindOptions"
            :label="ts(`${bp}.columns.kind`)"
            @update:model-value="(k) => onKindChange(k as AssetAttributeKind)"
          />
        </div>

        <div class="col-12 col-sm-5">
          <!-- Value input adapts to the composing attribute's kind -->
          <q-input
            v-if="newAttr.kind === 'integer'"
            :model-value="numValue()"
            type="number"
            outlined
            dense
            :label="ts(`${bp}.columns.value`)"
            @update:model-value="(v) => setIntegerValue(v)"
          />

          <q-toggle
            v-else-if="newAttr.kind === 'boolean'"
            :model-value="newAttr.value as boolean"
            :label="ts(`${bp}.columns.value`)"
            @update:model-value="(v) => setValue(v)"
          />

          <q-input
            v-else-if="newAttr.kind === 'date'"
            :model-value="(newAttr.value as string) ?? ''"
            outlined
            dense
            :label="ts(`${bp}.columns.value`)"
            :placeholder="ts(`${bp}.columns.datePlaceholder`)"
            @update:model-value="(v) => setValue(v)"
          >
            <template #append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date
                    :model-value="(newAttr.value as string) ?? ''"
                    mask="YYYY-MM-DDTHH:mm:ss[Z]"
                    @update:model-value="(v) => setValue(v)"
                  />
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>

          <div v-else-if="newAttr.kind === 'geo'" class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input
                :model-value="geoValue().lat"
                type="number"
                outlined
                dense
                :label="ts(`${bp}.geo.lat`)"
                @update:model-value="(v) => setGeoValue('lat', v)"
              />
            </div>
            <div class="col-6">
              <q-input
                :model-value="geoValue().lon"
                type="number"
                outlined
                dense
                :label="ts(`${bp}.geo.lon`)"
                @update:model-value="(v) => setGeoValue('lon', v)"
              />
            </div>
          </div>

          <q-input
            v-else
            :model-value="(newAttr.value as string) ?? ''"
            outlined
            dense
            :label="ts(`${bp}.columns.value`)"
            :placeholder="ts(`${bp}.columns.valuePlaceholder`)"
            @update:model-value="(v) => setValue(v)"
          />
        </div>

        <div class="col-12 col-sm-1 flex items-center justify-center">
          <q-btn
            unelevated
            round
            size="sm"
            color="primary"
            icon="add"
            :disable="!canAdd"
            @click="addAttribute"
          >
            <AppTooltip :content="ts(`${bp}.addButton`)" />
          </q-btn>
        </div>
      </q-card-section>
    </q-card>

    <!-- Committed attributes -->
    <div v-if="!rows.length" class="attributes-section__empty">
      {{ ts(`${bp}.empty`) }}
    </div>

    <q-list v-else bordered separator class="attributes-section__list">
      <q-item v-for="(row, index) in rows" :key="index">
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ row.label }}</q-item-label>
          <q-item-label caption>{{ kindLabel(row.kind) }} · {{ displayValue(row) }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-btn
            flat
            dense
            round
            icon="delete"
            color="negative"
            @click="removeRow(index)"
          >
            <AppTooltip :content="ts(`${bp}.remove`)" />
          </q-btn>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<style scoped lang="scss">
.attributes-section {
  &__header {
    margin-bottom: var(--mapex-spacing-sm);
  }

  &__title {
    font-size: var(--mapex-font-md);
    font-weight: var(--mapex-font-weight-semibold);
    color: var(--mapex-text-primary);
  }

  &__subtitle {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-secondary);
  }

  &__add {
    border-radius: var(--mapex-radius-md);
  }

  &__empty {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-secondary);
    font-style: italic;
  }

  &__list {
    border-radius: var(--mapex-radius-md);
    overflow: hidden;
  }
}
</style>

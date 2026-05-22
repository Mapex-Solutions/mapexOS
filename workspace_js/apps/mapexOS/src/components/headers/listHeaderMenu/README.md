# ListHeaderMenu Component

Generic menu button for list pages providing items count display, items per page selection, and column visibility toggles.

## Features

- **Items Count Display**: Shows total items with singular/plural label
- **Items Per Page**: Dropdown to select pagination size
- **Column Visibility**: Checkboxes to toggle column visibility
- **Customizable**: Configure icon, labels, and available options
- **Two-way Binding**: Reactive updates to parent component

## Usage

```vue
<template>
  <ListHeaderMenu
    :items-count="assetsList.length"
    item-label="Asset"
    item-label-plural="Assets"
    icon="devices"
    :items-per-page="itemsPerPage"
    :columns="menuColumns"
    @update:items-per-page="itemsPerPage = $event"
    @update:columns="handleColumnsUpdate"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ListHeaderMenu } from '@components/headers';
import type { ListHeaderMenuColumn } from '@components/headers';

const itemsPerPage = ref(25);

const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'uuid', label: 'UUID', visible: true },
  { key: 'type', label: 'Type', visible: true },
  { key: 'protocol', label: 'Protocol', visible: true },
  { key: 'category', label: 'Category', visible: false },
]);

function handleColumnsUpdate(columns: ListHeaderMenuColumn[]) {
  menuColumns.value = columns;
  // Update your visible columns logic here
}
</script>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `itemsCount` | `number` | **required** | Total number of items to display |
| `itemLabel` | `string` | **required** | Singular label (e.g., "Asset", "Rule", "User") |
| `itemLabelPlural` | `string` | `${itemLabel}s` | Plural label (e.g., "Assets", "Rules", "Users") |
| `icon` | `string` | `'list'` | Icon name for the button |
| `itemsPerPage` | `number` | **required** | Current items per page value |
| `itemsPerPageOptions` | `number[]` | `[10, 25, 50, 100]` | Available items per page options |
| `columns` | `ListHeaderMenuColumn[]` | `[]` | Column visibility configuration |
| `showItemsPerPage` | `boolean` | `true` | Show/hide items per page section |
| `showColumnVisibility` | `boolean` | `true` | Show/hide column visibility section |

## Column Configuration

### ListHeaderMenuColumn Interface

```typescript
interface ListHeaderMenuColumn {
  key: string;      // Unique column identifier
  label: string;    // Display label in menu
  visible: boolean; // Current visibility state
}
```

## Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:itemsPerPage` | `value: number` | Emitted when items per page changes |
| `update:columns` | `columns: ListHeaderMenuColumn[]` | Emitted when column visibility changes |

## Button Label Format

The component automatically formats the button label:
- **Singular**: "1 ASSET"
- **Plural (default)**: "6 ASSETS"
- **Plural (custom)**: "6 ITEMS" (if `itemLabelPlural="Items"`)

## Menu Sections

### Items Per Page Section

Displays when `showItemsPerPage !== false` and `itemsPerPageOptions.length > 0`.

- Shows radio-style selection
- Current selection marked with checkmark
- Closes menu on selection
- Emits `update:itemsPerPage` event

### Column Visibility Section

Displays when `showColumnVisibility !== false` and `columns.length > 0`.

- Shows checkboxes for each column
- Toggle individual columns
- Real-time updates
- Emits `update:columns` event

## Examples

### Basic Usage (Items Count Only)

```vue
<ListHeaderMenu
  :items-count="items.length"
  item-label="Item"
  :items-per-page="25"
  :show-items-per-page="false"
  :show-column-visibility="false"
/>
```

### Items Per Page Only

```vue
<ListHeaderMenu
  :items-count="users.length"
  item-label="User"
  :items-per-page="itemsPerPage"
  :items-per-page-options="[5, 10, 20, 50]"
  :show-column-visibility="false"
  @update:items-per-page="itemsPerPage = $event"
/>
```

### Column Visibility Only

```vue
<ListHeaderMenu
  :items-count="assets.length"
  item-label="Asset"
  :items-per-page="25"
  :columns="columns"
  :show-items-per-page="false"
  @update:columns="handleColumnsUpdate"
/>
```

### Full Configuration

```vue
<ListHeaderMenu
  :items-count="rules.length"
  item-label="Rule"
  item-label-plural="Rules"
  icon="rule"
  :items-per-page="itemsPerPage"
  :items-per-page-options="[10, 25, 50, 100]"
  :columns="menuColumns"
  @update:items-per-page="itemsPerPage = $event"
  @update:columns="handleColumnsUpdate"
/>
```

### Custom Labels

```vue
<ListHeaderMenu
  :items-count="12"
  item-label="Entry"
  item-label-plural="Entries"  <!-- "12 ENTRIES" instead of "12 ENTRYS" -->
  icon="storage"
  :items-per-page="25"
/>
```

## Integration with DataRow

Use together with DataRow component for complete list functionality:

```vue
<template>
  <!-- Header with menu -->
  <div class="row items-center q-mb-md">
    <div class="col">
      <div class="text-subtitle1 text-weight-medium text-primary">
        Assets List
      </div>
    </div>
    <div class="col-auto">
      <ListHeaderMenu
        :items-count="filteredAssets.length"
        item-label="Asset"
        icon="devices"
        :items-per-page="itemsPerPage"
        :columns="menuColumns"
        @update:items-per-page="itemsPerPage = $event"
        @update:columns="handleColumnsUpdate"
      />
    </div>
  </div>

  <!-- Data rows -->
  <div class="row">
    <div
      v-for="asset in paginatedAssets"
      :key="asset.id"
      class="col-12 q-mb-xs"
    >
      <DataRow
        :data="asset"
        :columns="visibleColumns"
        @edit="handleEdit"
        @view="handleView"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { ListHeaderMenu, DataRow } from '@components';
import type { ListHeaderMenuColumn, DataRowColumn } from '@components';

// Menu state
const itemsPerPage = ref(25);
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'uuid', label: 'UUID', visible: true },
  { key: 'type', label: 'Type', visible: true },
  { key: 'protocol', label: 'Protocol', visible: true },
  { key: 'category', label: 'Category', visible: true },
]);

// Column definitions
const allColumns: DataRowColumn[] = [ /* ... */ ];

// Computed visible columns based on menu selection
const visibleColumns = computed(() => {
  return allColumns.filter(col => {
    const menuCol = menuColumns.value.find(mc => mc.key === col.key);
    return !menuCol || menuCol.visible;
  });
});

function handleColumnsUpdate(columns: ListHeaderMenuColumn[]) {
  menuColumns.value = columns;
}
</script>
```

## Styling

The component uses Quasar's default button and menu styles. No custom CSS required.

### Customization

Override button styles if needed:

```vue
<style scoped>
:deep(.q-btn) {
  font-weight: 600; /* Custom font weight */
}
</style>
```

## Accessibility

- Keyboard navigable menu
- Screen reader friendly labels
- ARIA attributes for checkboxes and radio items
- Focus management

## Related Components

- **DataRow**: Row component for displaying list items
- **PageHeader**: Page title and action buttons
- **ListFilter**: Filter component for lists

## Use Cases

This component is ideal for:
- Asset lists
- Rule lists
- User lists
- Device lists
- Any paginated list with optional column visibility

## Browser Support

Works on all modern browsers that support:
- CSS Grid
- CSS Flexbox
- ES2022

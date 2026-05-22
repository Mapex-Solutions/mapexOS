# DataRow Component

Generic row component for displaying data in list views with responsive layouts.

## Features

- **Desktop/Laptop Layout (>= 600px)**: Table-styled cards with all columns visible
- **Mobile Layout (< 600px)**: Compact cards with expandable details
- **Responsive Breakpoints**: Automatic column hiding based on screen size
- **Actions Menu**: Built-in Edit, View, Delete actions
- **Column Types**: Avatar, Text, Code, Chip, Badge
- **Secondary Text**: Support for showing additional text below main content
- **Click Handlers**: Single click, double click, and action events

## Usage

```vue
<template>
  <DataRow
    :data="asset"
    :columns="assetColumns"
    primary-key="id"
    @click="handleClick"
    @dblclick="handleEdit"
    @edit="handleEdit"
    @view="handleView"
    @delete="handleDelete"
  />
</template>

<script setup lang="ts">
import { DataRow } from '@components/cards';
import type { DataRowColumn } from '@components/cards';

const assetColumns: DataRowColumn[] = [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: (value, row) => row.icon || 'sensors',
    color: (value, row) => row.status ? 'primary' : 'grey-5',
  },
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    secondaryKey: 'description', // Shows description below name
  },
  {
    key: 'uuid',
    label: 'UUID',
    type: 'code',
    visible: 'always',
    width: 180,
    ellipsis: true,
  },
  {
    key: 'type',
    label: 'Type',
    type: 'chip',
    visible: 'always',
    width: 150,
    color: 'blue-6',
  },
  {
    key: 'protocol.type',
    label: 'Protocol',
    type: 'chip',
    visible: 'desktop', // Hidden on screens <= 1024px
    width: 120,
    format: (value) => value?.toUpperCase() || 'N/A',
  },
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    visible: 'always',
    width: 100,
    format: (value) => value ? 'ACTIVE' : 'INACTIVE',
    color: (value) => value ? 'green-6' : 'red-6',
  },
];
</script>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `data` | `any` | **required** | The data object to display |
| `columns` | `DataRowColumn[]` | **required** | Column configuration array |
| `primaryKey` | `string` | `'id'` | Unique key for the row |
| `showActions` | `boolean` | `true` | Show/hide actions menu |
| `expandOnClick` | `boolean` | `true` | Enable mobile card expansion on click |

## Column Configuration

### DataRowColumn Interface

```typescript
interface DataRowColumn {
  key: string;                    // Object property path (supports nested: 'protocol.type')
  label: string;                  // Column header label
  type: DataRowColumnType;        // Column type: 'avatar' | 'text' | 'code' | 'chip' | 'badge'
  visible: DataRowColumnVisibility; // Visibility: 'always' | 'desktop' | 'laptop' | 'expandable'
  width?: number;                 // Column width in pixels
  ellipsis?: boolean;             // Enable text truncation with tooltip
  secondaryKey?: string;          // Secondary text key (shown below main text)

  // Formatting & Styling
  format?: (value: any, row: any) => string;           // Value formatter
  color?: string | ((value: any, row: any) => string); // Color value or function
  icon?: string | ((value: any, row: any) => string);  // Icon name or function (for avatar)
}
```

### Column Types

- **`avatar`**: Displays an icon with background color
- **`text`**: Plain text with optional secondary text below
- **`code`**: Monospace text for UUIDs, IDs, etc.
- **`chip`**: Colored chip/tag
- **`badge`**: Colored badge for status indicators

### Visibility Options

- **`always`**: Always visible on all screen sizes
- **`desktop`**: Hidden on screens <= 1024px
- **`laptop`**: Hidden on screens <= 1366px
- **`expandable`**: Only shown in mobile expanded area

## Events

| Event | Payload | Description |
|-------|---------|-------------|
| `click` | `data: any` | Emitted when row is clicked |
| `dblclick` | `data: any` | Emitted when row is double-clicked |
| `edit` | `data: any` | Emitted when Edit action is clicked |
| `view` | `data: any` | Emitted when View action is clicked |
| `delete` | `data: any` | Emitted when Delete action is clicked |
| `expand` | `data: any, expanded: boolean` | Emitted when mobile card is expanded/collapsed |

## Responsive Behavior

### Desktop/Laptop (>= 600px)
- Displays as horizontal row with all visible columns
- Hover effects with shadow
- Actions menu aligned to right edge
- Tooltips on truncated text

### Mobile (< 600px)
- Compact card showing: Avatar, Name+Description, Status, Actions
- Click anywhere on card to expand/collapse
- Expanded area shows remaining columns in 2-column grid
- Smooth slide transition

## Examples

### Basic Usage

```vue
<DataRow
  :data="{ id: 1, name: 'Temperature Sensor', status: true }"
  :columns="columns"
/>
```

### With All Events

```vue
<DataRow
  :data="item"
  :columns="columns"
  @click="console.log('Clicked:', $event)"
  @dblclick="openEditor($event)"
  @edit="openEditor($event)"
  @view="openDetails($event)"
  @delete="confirmDelete($event)"
/>
```

### Nested Property Access

```vue
const columns = [
  {
    key: 'protocol.type',        // Access nested property
    label: 'Protocol',
    type: 'chip',
    format: (value) => value?.toUpperCase(),
  },
];
```

### Dynamic Styling

```vue
const columns = [
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    format: (value) => value ? 'ACTIVE' : 'INACTIVE',
    color: (value) => value ? 'green-6' : 'red-6', // Dynamic color
  },
];
```

## Styling

The component uses scoped SCSS with Quasar utility classes. All styles are automatically applied - no additional CSS needed.

### Customization

You can override styles using deep selectors:

```vue
<style scoped>
:deep(.data-row-card) {
  border-radius: 8px; /* Custom border radius */
}

:deep(.data-row-card:hover) {
  background-color: #f0f0f0; /* Custom hover color */
}
</style>
```

## Related Components

- **ListHeaderMenu**: Menu for items per page and column visibility
- **Column Components**: AvatarColumn, TextColumn, CodeColumn, ChipColumn, BadgeColumn

## Browser Support

Works on all modern browsers that support:
- CSS Grid
- CSS Flexbox
- ES2022
- Media Queries

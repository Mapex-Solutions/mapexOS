# MapexValidationHelpModal Usage

This modal component provides comprehensive documentation for the `$mv` (MapexValidation) library used in validation scripts.

## Installation

The component is already set up and can be imported from the dialogs directory:

```typescript
import { MapexValidationHelpModal } from '@/components/dialogs/mapexValidationHelp';
```

## Basic Usage

```vue
<script setup lang="ts">
import { ref } from 'vue';
import { MapexValidationHelpModal } from '@/components/dialogs/mapexValidationHelp';

const showHelpModal = ref(false);

const openHelp = () => {
  showHelpModal.value = true;
};
</script>

<template>
  <div>
    <!-- Trigger button -->
    <q-btn
      color="primary"
      icon="help"
      label="Validation Help"
      @click="openHelp"
    />

    <!-- Help modal -->
    <MapexValidationHelpModal v-model="showHelpModal" />
  </div>
</template>
```

## Props

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `modelValue` | `boolean` | Yes | Controls the visibility of the modal |

## Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:modelValue` | `boolean` | Emitted when the modal is closed/opened |

## Features

1. **8 Comprehensive Tabs**:
   - Overview - Introduction and basic usage
   - String - String validation methods
   - Number - Number validation methods
   - Boolean - Boolean validation methods
   - Date - Date validation methods
   - Array - Array validation methods
   - Object - Object validation methods
   - Examples - Real-world complete examples

2. **Interactive Features**:
   - Copy-to-clipboard for all code examples
   - Search functionality to filter methods
   - Responsive design (mobile-friendly)
   - Syntax-highlighted code blocks
   - Tooltips on all interactive elements

3. **Search Functionality**:
   - Search across all method names, descriptions, and code
   - Automatic tab switching to show relevant results
   - Keyword-based filtering

## Integration Example

### In a validation script editor:

```vue
<script setup lang="ts">
import { ref } from 'vue';
import { MapexValidationHelpModal } from '@/components/dialogs/mapexValidationHelp';

const showValidationHelp = ref(false);
const validationScript = ref('');

const insertExample = (example: string) => {
  validationScript.value = example;
};
</script>

<template>
  <div class="validation-editor">
    <q-card>
      <q-card-section class="row items-center">
        <div class="text-h6">Validation Script</div>
        <q-space />
        <q-btn
          icon="help_outline"
          label="Help"
          color="primary"
          flat
          @click="showValidationHelp = true"
        >
          <q-tooltip>View MapexValidation documentation</q-tooltip>
        </q-btn>
      </q-card-section>

      <q-card-section>
        <q-input
          v-model="validationScript"
          type="textarea"
          outlined
          rows="10"
          placeholder="Enter your validation script here..."
        />
      </q-card-section>
    </q-card>

    <MapexValidationHelpModal v-model="showValidationHelp" />
  </div>
</template>
```

## Styling

The component uses Quasar's built-in styling system with:
- Primary color scheme
- Grey backgrounds for code blocks
- Rounded corners on cards
- Smooth transitions and hover effects
- Responsive layout that adapts to screen size

## Accessibility

- All buttons have tooltips
- Keyboard navigation supported
- Screen reader friendly
- High contrast code blocks
- Clear visual hierarchy

## Notes

- The modal opens in fullscreen mode for better readability
- Code blocks are scrollable horizontally if content is too wide
- The search is case-insensitive and searches across titles, descriptions, code, and keywords
- Closing the modal (ESC key or close button) will preserve the current tab selection

# Translation Utility

Utility functions for handling i18n translations with text formatting options.

## Overview

The translation utility provides a consistent way to handle translations with automatic text transformations like capitalization, uppercase, lowercase, and title case.

## API

### `useTS(options?)`

Main translation utility with customizable text formatting.

```typescript
interface TranslationOptions {
  capitalize?: boolean;  // Capitalize first letter (default: true)
  uppercase?: boolean;   // Transform to UPPERCASE
  lowercase?: boolean;   // Transform to lowercase
  titleCase?: boolean;   // Transform To Title Case
}
```

### `useRawTS()`

Returns raw translation without any formatting.

### `useTranslationPresets()`

Returns pre-configured translation functions for common use cases.

## Usage Examples

### Basic Usage (Default - Capitalize)

```typescript
import { useTS } from '@utils/translation';

const ts = useTS();

// Simple translation
const saveBtn = ts('common.actions.save'); // "Save"

// With interpolation
const msg = ts('common.messages.deletedSuccessfully', { item: 'User' });
// Result: "User deleted successfully"
```

### Text Transformations

```typescript
import { useTS } from '@utils/translation';

// UPPERCASE
const tsUpper = useTS({ uppercase: true });
const text = tsUpper('common.actions.save'); // "SAVE"

// lowercase
const tsLower = useTS({ lowercase: true });
const text = tsLower('common.actions.save'); // "save"

// Title Case
const tsTitle = useTS({ titleCase: true });
const text = tsTitle('pages.settings.title'); // "System Settings"

// No transformation
const tsRaw = useTS({ capitalize: false });
const text = tsRaw('some.key'); // exact text from JSON
```

### Override Options Per Call

```typescript
const ts = useTS({ capitalize: true }); // default capitalize

// Use default
const text1 = ts('common.actions.save'); // "Save"

// Override to uppercase for this call only
const text2 = ts('common.actions.save', {}, { uppercase: true }); // "SAVE"

// Override to title case
const text3 = ts('pages.settings.title', {}, { titleCase: true }); // "System Settings"
```

### Using Presets

```typescript
import { useTranslationPresets } from '@utils/translation';

const { ts, tsUpper, tsLower, tsTitle, tsRaw } = useTranslationPresets();

const normal = ts('common.actions.save');      // "Save"
const upper = tsUpper('common.actions.save');  // "SAVE"
const lower = tsLower('common.actions.save');  // "save"
const title = tsTitle('pages.settings.title'); // "System Settings"
const raw = tsRaw('some.key');                 // exact from JSON
```

### In Composables

```typescript
// composables/i18n/pages/useSettingsTranslations.ts
import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useSettingsTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    pageTitle: computed(() => ts('pages.administrations.settings.title')),
    pageDescription: computed(() => ts('pages.administrations.settings.description')),

    tabs: {
      general: computed(() => ts('pages.administrations.settings.tabs.general')),
      lists: computed(() => ts('pages.administrations.settings.tabs.lists')),
    },
  };
}
```

### In Vue Components

```vue
<script setup lang="ts">
import { useTS } from '@utils/translation';

const ts = useTS({ capitalize: true });

// Use in template
const pageTitle = computed(() => ts('pages.settings.title'));

// Use in functions
function handleSave() {
  notifySuccess({
    message: ts('pages.settings.messages.savedSuccessfully')
  });
}
</script>

<template>
  <PageHeader :title="pageTitle" />
</template>
```

## Best Practices

### ✅ DO

1. **Use composables** for page-level translations to keep components clean
2. **Use default capitalize** for most UI text
3. **Use uppercase** for emphasis (buttons, alerts)
4. **Use title case** for page titles and headers
5. **Use raw** when you need exact text (codes, technical terms)

```typescript
// ✅ Good - Organized in composable
const { pageTitle, formLabels } = useSettingsTranslations();

// ✅ Good - Use presets
const { ts, tsUpper } = useTranslationPresets();
const saveBtn = ts('common.actions.save');
const alertMsg = tsUpper('alerts.warning');

// ✅ Good - Override when needed
const title = ts('pages.settings.title', {}, { titleCase: true });
```

### ❌ DON'T

1. **Don't use inline $t() with manual formatting** - use the utility instead
2. **Don't repeat transformation logic** in multiple places
3. **Don't mix transformation approaches**

```typescript
// ❌ Bad - Manual formatting
const text = $t('some.key').toUpperCase();

// ❌ Bad - Inline complexity
const text = capitalize($t('some.key'));

// ✅ Good - Use utility
const text = tsUpper('some.key');
```

## Transformation Priority

When multiple transformations are specified, they are applied in this order:

1. `uppercase` (highest priority)
2. `lowercase`
3. `titleCase`
4. `capitalize`

If multiple options are `true`, only the highest priority one is applied.

## TypeScript Support

Full TypeScript support with type inference:

```typescript
import { useTS, type TranslationOptions } from '@utils/translation';

const options: TranslationOptions = {
  capitalize: true,
  titleCase: false,
};

const ts = useTS(options);
// ts is typed as: (key: string, interpolation?: Record<string, unknown>, options?: TranslationOptions) => string
```

## Migration Guide

### From Direct `$t()` or `t()`

```typescript
// Before
const text = $t('common.actions.save');
const upperText = $t('common.actions.save').toUpperCase();

// After
const ts = useTS();
const text = ts('common.actions.save');
const upperText = ts('common.actions.save', {}, { uppercase: true });
```

### From Legacy `useTS` (index.js)

```typescript
// Before (index.js)
const ts = useTS({ useCapitalize: true });
const text = ts('some.key');

// After (index.ts)
const ts = useTS({ capitalize: true });
const text = ts('some.key');
```

## Performance

- **Computed values**: Use `computed()` in composables to cache translations
- **Memoization**: Transformations are lightweight, no additional caching needed
- **Reactive**: All translations update automatically when language changes

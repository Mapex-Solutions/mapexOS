# i18n (Internationalization) Structure

## 📁 Folder Structure

```
src/i18n/
├── en/                          # English locale
│   ├── index.ts                 # Main export
│   ├── common.json              # Shared translations
│   ├── components/              # Component translations
│   │   ├── headers.json
│   │   ├── cards.json
│   │   └── filters.json
│   └── pages/                   # Page translations
│       └── administrations/
│           ├── settings.json
│           ├── groups.json
│           └── roles.json
├── pt-BR/                       # Portuguese (Brazil) locale
│   └── (same structure as en/)
└── index.ts                     # Main i18n config
```

## 🎯 Organization Principles

### 1. **Common Translations** (`common.json`)
Shared across the entire application:
- **actions**: Buttons and action labels (save, edit, delete, etc.)
- **messages**: Success, error, confirmation messages
- **labels**: Generic labels (name, description, status, etc.)
- **status**: Status values (active, inactive, pending, etc.)
- **pagination**: Pagination-related text
- **validation**: Form validation messages

### 2. **Component Translations** (`components/`)
Organized by component type:
- **headers.json**: PageHeader, ListHeaderMenu
- **cards.json**: DataRow, ListCard, EmptyCard
- **filters.json**: ListFilter, DateRange
- **forms.json**: Form components (future)

### 3. **Page Translations** (`pages/`)
Organized by module/feature:
- **pages/administrations/**: Settings, Groups, Roles
- **pages/assets/**: Assets pages (future)
- **pages/notifications/**: Notifications pages (future)

## 🔧 Usage Examples

### 1. Basic Usage (Script Setup)

```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// Using translations
const title = t('pages.administrations.settings.title');
const saveButton = t('common.actions.save');
const deleteMessage = t('common.messages.confirmDelete', { item: 'group' });
</script>
```

### 2. Template Usage

```vue
<template>
  <!-- Simple translation -->
  <q-btn :label="$t('common.actions.save')" />

  <!-- With interpolation -->
  <div>{{ $t('common.messages.deletedSuccessfully', { item: 'User' }) }}</div>

  <!-- Component props -->
  <PageHeader
    :title="$t('pages.administrations.settings.title')"
    :description="$t('pages.administrations.settings.description')"
  />
</template>
```

### 3. Pluralization

```vue
<!-- Automatic pluralization based on count -->
<div>{{ $t('common.pagination.items', { count: 5 }) }}</div>
<!-- Output: "5 items" -->

<div>{{ $t('common.pagination.items', { count: 1 }) }}</div>
<!-- Output: "1 item" -->
```

### 4. Programmatic Usage

```typescript
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

function deleteGroup(name: string) {
  $q.dialog({
    title: t('common.actions.confirmDelete'),
    message: t('pages.administrations.groups.messages.confirmDelete', { name }),
  }).onOk(() => {
    // Delete logic
    $q.notify({
      message: t('pages.administrations.groups.messages.deletedSuccessfully'),
      color: 'positive',
    });
  });
}
```

## 🌍 Adding a New Locale

1. Create folder: `src/i18n/<locale>/`
2. Copy structure from `en/` or `pt-BR/`
3. Translate all JSON files
4. Create `index.ts` with same structure
5. Update `src/i18n/index.ts` to include new locale

```typescript
// src/i18n/index.ts
import es from './es'; // Spanish

export default {
  'en': en,
  'pt-BR': ptBR,
  'es': es, // Add new locale
};
```

## 📝 Adding New Translations

### For Common Translations
Add to `<locale>/common.json`:

```json
{
  "actions": {
    "newAction": "New Action"
  }
}
```

### For Component Translations
Add to `<locale>/components/<component>.json`:

```json
{
  "myComponent": {
    "title": "My Component Title"
  }
}
```

### For Page Translations
Add to `<locale>/pages/<module>/<page>.json`:

```json
{
  "title": "Page Title",
  "sections": {
    "main": "Main Section"
  }
}
```

## 🤖 AI Translation Workflow

### 1. Extract Text from Component
```bash
# Identify all hardcoded text
grep -r "label=\"" src/pages/
```

### 2. Add to English JSON
```json
{
  "myPage": {
    "button": "Click Me",
    "title": "Page Title"
  }
}
```

### 3. Use AI to Translate
```
Prompt: "Translate the following JSON to Portuguese (Brazil), maintaining the same structure:
{
  "myPage": {
    "button": "Click Me",
    "title": "Page Title"
  }
}"
```

### 4. Validate Output
```json
{
  "myPage": {
    "button": "Clique Aqui",
    "title": "Título da Página"
  }
}
```

## 🎨 Best Practices

### ✅ DO
- Use descriptive keys: `pages.administrations.settings.general.form.organizationName`
- Group related translations together
- Use interpolation for dynamic values: `{name}`, `{count}`
- Keep translations in JSON files, not in code
- Use pluralization rules properly

### ❌ DON'T
- Don't hardcode text in components
- Don't create deeply nested structures (max 4 levels)
- Don't duplicate translations across files
- Don't mix languages in the same file

## 🔍 Translation Keys Naming Convention

```
<context>.<subcontext>.<element>
```

Examples:
- `common.actions.save` - Common action
- `components.headers.pageHeader.defaultDescription` - Component text
- `pages.administrations.settings.general.title` - Page section title
- `pages.administrations.groups.messages.deletedSuccessfully` - Page message

## 🚀 Future Enhancements

1. **Lazy Loading**: Load locale files on demand
2. **Date/Number Formatting**: Configure per locale
3. **RTL Support**: For Arabic, Hebrew, etc.
4. **Translation Management**: Tool for non-developers
5. **Auto-detection**: Browser language detection
6. **Fallback**: Automatic fallback to English

## 📚 Resources

- [Vue I18n Docs](https://vue-i18n.intlify.dev/)
- [ICU Message Format](https://unicode-org.github.io/icu/userguide/format_parse/messages/)
- [Pluralization Rules](https://vue-i18n.intlify.dev/guide/essentials/pluralization.html)

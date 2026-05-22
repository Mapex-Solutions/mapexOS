# AvailableFieldsList Component

**Componente de lista formal para exibição de campos disponíveis (available fields) no MapexOS.**

---

## 📋 Descrição

Componente reutilizável que exibe uma lista numerada e scrollável de campos disponíveis extraídos de um StandardizedPayload. Utilizado principalmente em Asset Templates para mostrar quais campos estarão disponíveis para autocomplete em Rules.

---

## 🎯 Uso

### Importação

```typescript
import { AvailableFieldsList } from '@components/lists/availableFieldsList';
```

### Exemplo Básico

```vue
<template>
  <AvailableFieldsList
    :fields="['eventType', 'eventId', 'data.temperature']"
  />
</template>

<script setup lang="ts">
import { AvailableFieldsList } from '@components/lists/availableFieldsList';
</script>
```

### Exemplo Completo (com loading e eventos)

```vue
<template>
  <div>
    <div class="text-subtitle2 q-mb-sm">
      Available Fields ({{ fields.length }})
    </div>

    <AvailableFieldsList
      :fields="fields"
      :loading="isLoading"
      :max-height="400"
      @field-click="handleFieldClick"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { AvailableFieldsList } from '@components/lists/availableFieldsList';
import { copyToClipboard, Notify } from 'quasar';

const fields = ref([
  'eventType',
  'eventId',
  'data.temperature',
  'data.humidity',
  'data.location.lat',
  'data.location.lng',
  'metadata.gateway',
  'created'
]);

const isLoading = ref(false);

/**
 * Handle field click - copy to clipboard
 */
function handleFieldClick(field: string): void {
  copyToClipboard(field)
    .then(() => {
      Notify.create({
        type: 'positive',
        message: `Field "${field}" copied to clipboard`,
        timeout: 2000
      });
    })
    .catch(() => {
      Notify.create({
        type: 'negative',
        message: 'Failed to copy to clipboard',
        timeout: 2000
      });
    });
}
</script>
```

---

## 📦 Props

| Prop | Tipo | Default | Descrição |
|------|------|---------|-----------|
| `fields` | `string[]` | `[]` (required) | Array de campos (paths) a serem exibidos |
| `loading` | `boolean` | `false` | Mostra spinner de loading ao invés da lista |
| `maxHeight` | `number` | `300` | Altura máxima da scroll area (em pixels) quando > 10 campos |

---

## 🎭 Events

| Event | Payload | Descrição |
|-------|---------|-----------|
| `field-click` | `string` | Emitido quando um campo é clicado, retorna o path do campo |

---

## 🎨 Características

### Visual
- ✅ **Numeração sequencial** - Avatar com número de 1 a N
- ✅ **Monospace font** - Campos em `Roboto Mono` (parecem código)
- ✅ **Hover effect** - Background sutil ao passar o mouse
- ✅ **Ícone lateral** - Ícone `code` em cada item
- ✅ **Scrollável** - Altura fixa quando > 10 campos

### Estados
- ✅ **Loading** - Spinner centralizado com mensagem
- ✅ **Empty** - Mensagem quando array vazio
- ✅ **Normal** - Lista numerada e clickável

### Comportamento
- ✅ **Auto-scroll** - Ativa scroll quando `fields.length > 10`
- ✅ **Clickável** - Cada item emite evento `field-click`
- ✅ **Responsivo** - Funciona em mobile, tablet, desktop

---

## 🏗️ Estrutura de Arquivos

```
availableFieldsList/
├── constants/
│   ├── availableFieldsList.constant.ts  # DEFAULT_MAX_HEIGHT, SCROLL_THRESHOLD
│   └── index.ts
├── interfaces/
│   ├── availableFieldsList.interface.ts # Props & Emits
│   └── index.ts
├── AvailableFieldsList.vue              # Componente principal
├── index.ts                              # Exports públicos
└── README.md                             # Esta documentação
```

---

## 📐 Constantes

```typescript
// constants/availableFieldsList.constant.ts

/** Altura máxima padrão da scroll area (em pixels) */
export const DEFAULT_MAX_HEIGHT = 300;

/** Número de campos que ativa o scroll */
export const SCROLL_THRESHOLD = 10;
```

---

## 🔧 Customização

### Altura Customizada

```vue
<AvailableFieldsList
  :fields="fields"
  :max-height="500"  <!-- 500px ao invés de 300px -->
/>
```

### Sem Loading State

```vue
<AvailableFieldsList
  :fields="fields"
  :loading="false"  <!-- Sempre mostra lista ou empty state -->
/>
```

---

## 🎯 Casos de Uso

### 1. Asset Template Testing
Mostrar campos extraídos após executar teste de conversão:

```vue
<AvailableFieldsList
  :fields="normalizedFields"
  :loading="isTestRunning"
  @field-click="copyFieldToClipboard"
/>
```

### 2. Rule Condition Builder
Mostrar campos disponíveis ao criar condições:

```vue
<AvailableFieldsList
  :fields="availableFieldsFromTemplate"
  @field-click="insertFieldIntoCondition"
/>
```

### 3. Documentation/Help
Exibir campos disponíveis em modais de ajuda:

```vue
<q-dialog v-model="showHelp">
  <q-card>
    <q-card-section>
      <div class="text-h6">Available Event Fields</div>
      <AvailableFieldsList :fields="documentedFields" />
    </q-card-section>
  </q-card>
</q-dialog>
```

---

## ✅ Compliance com MapexOS

Este componente segue **RIGOROSAMENTE** os padrões MapexOS:

### Estrutura (Item 1 do guia)
- ✅ Pasta `constants/` com arquivo e index
- ✅ Pasta `interfaces/` com arquivo e index
- ✅ `index.ts` exporta tudo
- ✅ Naming: PascalCase para componente, camelCase para arquivos

### Código (Item 2 do guia)
- ✅ **Seções obrigatórias** - Todos os comentários `/** SECTION */`
- ✅ **Ordem exata** - TYPE IMPORTS → VUE → COMPONENTS → ... → LIFECYCLE
- ✅ **TSDoc completo** - Todas funções e interfaces documentadas
- ✅ **Type imports separados** - Seção 1 só com `type`
- ✅ **Local imports corretos** - Seção 8 sem types

---

## 🧪 Testes Recomendados

```typescript
describe('AvailableFieldsList', () => {
  it('should display all fields', () => {
    const fields = ['field1', 'field2', 'field3'];
    // mount component with fields
    // assert 3 items rendered
  });

  it('should show loading state', () => {
    // mount with loading=true
    // assert spinner is visible
  });

  it('should show empty state', () => {
    // mount with fields=[]
    // assert "No fields available" message
  });

  it('should emit field-click event', () => {
    // mount with fields
    // click on item
    // assert event emitted with correct field
  });

  it('should enable scroll when > 10 fields', () => {
    // mount with 15 fields
    // assert scroll area has fixed height
  });
});
```

---

## 📚 Referências

- **Guia Frontend MapexOS:** `/.claude/agents/frontend/mapexos-ui-generic.md`
- **Componentes Quasar:** [Q-List](https://quasar.dev/vue-components/list-and-list-items), [Q-Scroll-Area](https://quasar.dev/vue-components/scroll-area)
- **Feature Spec:** `/documentations/Implementations/asset-template-available_fields.md`

---

**Criado em:** 2025-01-17
**Autor:** MapexOS Team
**Versão:** 1.0.0

# Error Handling Utility 🛡️

Utilitário centralizado para tratamento de erros de API com suporte completo a i18n.

## 📁 Estrutura

```
src/utils/error/
├── types.ts              # Tipos e interfaces TypeScript
├── handleApiError.ts     # Função principal de tratamento de erros
├── index.ts              # Exports públicos
└── README.md             # Este arquivo
```

## 🚀 Uso Básico

### Importação

```typescript
import { handleApiError } from '@utils/error';
```

### ⚠️ IMPORTANTE: defaultMessage é Obrigatória

A função `handleApiError` **não** usa composables de i18n internamente para evitar problemas de contexto Vue. Você **DEVE** sempre fornecer um `defaultMessage`.

### Exemplo 1: Uso simples com mensagem padrão

```typescript
import { useMyTranslations } from '@composables/i18n/useMyTranslations';

const t = useMyTranslations();

try {
  await apis.users.user.delete({ userId });
} catch (error) {
  handleApiError(error, {
    defaultMessage: t.notifications.deleteFailed
  });
}
```

### Exemplo 2: Com mensagens customizadas

```typescript
import { handleApiError } from '@utils/error';
import { useMyTranslations } from '@composables/i18n/useMyTranslations';

const t = useMyTranslations();

try {
  await apis.assets.assetTemplate.create(data);
} catch (error) {
  handleApiError(error, {
    customMessages: {
      409: t.notifications.alreadyExists,    // Conflito - já existe
      422: t.notifications.validationFailed, // Erro de validação
      network: t.notifications.networkError, // Erro de rede
    },
    defaultMessage: t.notifications.creationFailed
  });
}
```

### Exemplo 3: Configuração completa

```typescript
const t = useMyTranslations();

handleApiError(error, {
  // Mensagens customizadas por código HTTP
  customMessages: {
    400: t.notifications.badRequest,
    409: t.notifications.alreadyExists,
    422: t.notifications.validationFailed,
    network: t.notifications.checkConnection, // Erro de rede
  },

  // Mensagem padrão caso não haja match (OBRIGATÓRIA)
  defaultMessage: t.notifications.operationFailed,

  // Timeout da notificação (em ms)
  timeout: 5000,

  // Desabilitar log no console
  logError: false,

  // Callback customizado para tracking/analytics
  onError: (error) => {
    analytics.track('api_error', {
      endpoint: error.response?.config?.url,
      status: error.response?.status
    });
  }
});
```

### Exemplo 4: Usando traduções globais de erro (opcional)

Se quiser usar as traduções globais de erro para status codes não customizados:

```typescript
import { useErrorTranslations } from '@composables/i18n/common/useErrorTranslations';

const t = useMyTranslations();
const globalErrors = useErrorTranslations();

try {
  await apis.assets.asset.delete({ assetId });
} catch (error) {
  handleApiError(error, {
    customMessages: {
      403: t.notifications.noPermissionToDelete, // Mensagem específica do domínio
      404: globalErrors.http[404],               // Usa mensagem global
      500: globalErrors.http[500],               // Usa mensagem global
      network: globalErrors.http.network,        // Usa mensagem global
    },
    defaultMessage: t.notifications.deleteFailed
  });
}
```

## 🎯 HTTP Status Codes Suportados

As traduções globais (`useErrorTranslations`) fornecem mensagens prontas para:

| Código | Significado | Mensagem (PT-BR) |
|--------|-------------|------------------|
| 400 | Bad Request | Requisição inválida. Verifique os dados enviados. |
| 401 | Unauthorized | Não autorizado. Faça login novamente. |
| 403 | Forbidden | Acesso negado. Você não tem permissão. |
| 404 | Not Found | Recurso não encontrado. |
| 409 | Conflict | Conflito. O recurso já existe. |
| 422 | Unprocessable Entity | Erro de validação. Verifique os dados. |
| 500 | Internal Server Error | Erro interno do servidor. |
| 503 | Service Unavailable | Serviço temporariamente indisponível. |
| network | Network Error | Erro de rede. Verifique sua conexão. |
| unknown | Unknown Error | Ocorreu um erro inesperado. |

## 🌐 Internacionalização (i18n)

### Arquivos de Tradução Global

```
src/i18n/pt-BR/common/errors.json
src/i18n/en-US/common/errors.json
```

### Usando traduções globais no código

```typescript
import { useErrorTranslations } from '@composables/i18n/common/useErrorTranslations';

const errors = useErrorTranslations();

// Acessar mensagens globais
errors.http[422].value // "Erro de validação..."
errors.http.network.value // "Erro de rede..."
errors.http.unknown.value // "Ocorreu um erro inesperado..."
```

## 🔧 Interface TypeScript

### HandleApiErrorOptions

```typescript
interface HandleApiErrorOptions {
  /**
   * Mensagens customizadas mapeadas por código HTTP
   * @example { 409: t.notifications.alreadyExists, 422: t.notifications.validationFailed }
   */
  customMessages?: ErrorMessageMap;

  /**
   * Mensagem padrão se nenhuma customizada for encontrada
   * Se não fornecida, usa mensagem global genérica
   */
  defaultMessage?: ComputedRef<string> | string;

  /**
   * Timeout da notificação em milissegundos
   * @default 5000
   */
  timeout?: number;

  /**
   * Se deve logar o erro no console
   * @default true
   */
  logError?: boolean;

  /**
   * Callback customizado para processar erro
   * Útil para analytics, tracking, etc.
   */
  onError?: (error: ApiError) => void;
}
```

### ErrorMessageMap

```typescript
interface ErrorMessageMap {
  [key: number]: ComputedRef<string> | string; // Códigos HTTP numéricos
  network?: ComputedRef<string> | string;       // Erro de rede
  unknown?: ComputedRef<string> | string;       // Erro desconhecido
}
```

## 📊 Fluxo de Decisão de Mensagens

```
1. Verifica se error.response.status existe
   ↓
2. Procura em customMessages[statusCode]
   ↓ (não encontrou)
3. Procura em globalErrors.http[statusCode]
   ↓ (não encontrou)
4. Usa defaultMessage se fornecida
   ↓ (não fornecida)
5. Usa globalErrors.http.unknown

Se não houver error.response (erro de rede):
   → Usa customMessages.network OU globalErrors.http.network
```

## ✨ Boas Práticas

### 1. ✅ Use mensagens customizadas para erros específicos do domínio

```typescript
// BOM ✅
handleApiError(error, {
  customMessages: {
    409: t.notifications.templateAlreadyExists, // Específico
    422: t.notifications.invalidTemplateData,   // Específico
  }
});

// EVITE ❌
handleApiError(error, {
  customMessages: {
    409: 'Erro', // Muito genérico
  }
});
```

### 2. ✅ Deixe mensagens genéricas usarem traduções globais

```typescript
// BOM ✅ - Não precisa customizar 404, 500, etc.
handleApiError(error, {
  customMessages: {
    409: t.notifications.alreadyExists,
  }
});
```

### 3. ✅ Use logError: false em produção para erros esperados

```typescript
// Para erros esperados (ex: usuário não autenticado)
handleApiError(error, {
  logError: false, // Não poluir console
});
```

### 4. ✅ Use onError para tracking/analytics

```typescript
handleApiError(error, {
  onError: (err) => {
    if (err.response?.status >= 500) {
      // Reportar apenas erros de servidor
      errorTracker.report(err);
    }
  }
});
```

## 🎨 Exemplos Completos

### Criar Asset Template

```typescript
async function submitForm() {
  isCreating.value = true;

  try {
    const created = await apis.assets.assetTemplate.create(templateData);

    notifySuccess({
      message: t.notifications.created.value,
      timeout: 3000
    });

    await router.push('/assets/asset-templates');

  } catch (error) {
    handleApiError(error, {
      customMessages: {
        409: t.notifications.alreadyExists,
        422: t.notifications.validationFailed,
      },
      defaultMessage: t.notifications.creationFailed,
      timeout: 5000
    });
  } finally {
    isCreating.value = false;
  }
}
```

### Deletar Asset

```typescript
async function deleteAsset(assetId: string) {
  try {
    await apis.assets.asset.delete({ assetId });
    notifySuccess({ message: 'Asset deletado com sucesso!' });
  } catch (error) {
    handleApiError(error, {
      customMessages: {
        404: t.notifications.assetNotFound,
        403: t.notifications.noPermissionToDelete,
      },
      defaultMessage: t.notifications.deleteFailed
    });
  }
}
```

### Atualizar Usuário

```typescript
async function updateUser(userId: string, data: UserUpdate) {
  try {
    await apis.users.user.update({ userId }, data);
    notifySuccess({ message: 'Usuário atualizado!' });
  } catch (error) {
    handleApiError(error); // Usa apenas mensagens globais
  }
}
```

## 🔍 Debugging

Para debug, ative os logs:

```typescript
handleApiError(error, {
  logError: true, // default
  onError: (err) => {
    console.log('Status:', err.response?.status);
    console.log('Data:', err.response?.data);
    console.log('Message:', err.message);
  }
});
```

## 📝 Notas Importantes

1. **Sempre use try-catch** ao fazer chamadas de API
2. **ComputedRef vs String**: Ambos são aceitos no `customMessages`
3. **Network errors**: São capturados quando `error.response` é undefined
4. **Default behavior**: Se nenhuma opção for fornecida, usa traduções globais
5. **TypeScript**: Tipos completos garantem uso correto

## 🆘 Troubleshooting

### Mensagem não aparece

```typescript
// Verifique se a tradução existe
const t = useMyTranslations();
console.log(t.notifications.myError.value); // undefined?

// Ou use mensagem global como fallback
handleApiError(error, {
  defaultMessage: 'Erro ao processar requisição'
});
```

### Status code não mapeado

```typescript
// Adicione ao customMessages
handleApiError(error, {
  customMessages: {
    418: 'Sou um bule de chá!', // Status code customizado
  }
});
```

## 📚 Referências

- [HTTP Status Codes - MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [Axios Error Handling](https://axios-http.com/docs/handling_errors)
- [Vue i18n](https://vue-i18n.intlify.dev/)

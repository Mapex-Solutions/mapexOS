# Plugin Marketplace Mock Server

> Express simples que simula o Netlify CDN servindo o catálogo de plugins.
> Usado para testar o fetch do UI (aba Plugins no workflow editor).

---

## 1. Propósito

Em produção, os manifests de plugins ficam como arquivos estáticos no Netlify CDN (`docs.mapexos.com/plugins/`).
Para desenvolvimento local, este server simula o mesmo comportamento: serve JSONs estáticos com CORS.

```
PRODUÇÃO:                                    DEV LOCAL:
Netlify CDN                                  Express mock (porta 3099)
GET docs.mapexos.com/plugins/registry.json   GET localhost:3099/plugins/registry.json
GET docs.mapexos.com/plugins/telegram/...    GET localhost:3099/plugins/telegram/...
```

---

## 2. Estrutura

```
workspace_js/services/plugin-marketplace-mock/
├── package.json
├── tsconfig.json
├── src/
│   └── main.ts                        # Express server (~30 linhas)
└── public/
    └── plugins/
        ├── registry.json              # Índice de todos plugins disponíveis
        └── telegram/
            ├── manifest.json          # Manifest completo do Telegram
            └── icon.svg               # Ícone do plugin
```

---

## 3. Endpoints

Tudo servido via `express.static('./public')`:

| URL | Arquivo | Descrição |
|-----|---------|-----------|
| `GET /plugins/registry.json` | `public/plugins/registry.json` | Índice com todos os plugins |
| `GET /plugins/telegram/manifest.json` | `public/plugins/telegram/manifest.json` | Manifest completo Telegram |
| `GET /plugins/telegram/icon.svg` | `public/plugins/telegram/icon.svg` | Ícone SVG |

---

## 4. registry.json — Formato

Índice leve para listagem no marketplace (sem baixar manifests):

```json
{
  "$schema": "mapex-plugin-registry/v1",
  "version": "1.0.0",
  "updatedAt": "2026-03-14T00:00:00Z",
  "plugins": [
    {
      "id": "telegram",
      "name": "Telegram",
      "version": "1.0.0",
      "category": "messaging",
      "icon": "send",
      "brandIcon": "telegram/icon.svg",
      "color": "#0088CC",
      "description": "Send messages, photos and manage chats via Telegram Bot API",
      "author": "MapexOS",
      "tags": ["messaging", "chat", "bot"],
      "manifestUrl": "telegram/manifest.json",
      "docsUrl": "/docs/1.0.0/en/plugins/telegram",
      "requiresCredentials": true,
      "nodeCount": 2
    }
  ]
}
```

---

## 5. manifest.json — Plugin Telegram (alinhado com SDK)

O manifest segue a interface `WorkflowPlugin` do SDK com campos adicionais para integração:

### Campos existentes no SDK (não mudam):
- `id`, `name`, `version`, `category`, `icon`
- `nodeTypes[]` com: `type`, `label`, `icon`, `color`, `description`, `inputs`, `outputs`, `properties`, `defaults`, `outputHints`
- `NodePropertyDefinition`: `name`, `displayName`, `type`, `default`, `hint`, `required`, `options`, `displayOptions`

### Campos NOVOS (aditivos, não quebram core):

**No `WorkflowPlugin`:**
- `metadata?: PluginMetadata` — brandIcon, color, docs
- `credentials?: PluginCredentialDefinition` — campos de autenticação
- `baseUrl?: string` — URL base da API externa

**No `NodePropertyDefinition`:**
- `type: 'fieldSource'` — novo valor no union (renderiza FieldSourceSelector)
- `allowedSources?: SourceType[]` — quais sources permitir
- `typeOptions?: { multiline?, rows? }` — config do input
- `requestName?: string` — nome do campo no request HTTP
- `isSecret?: boolean` — backend encrypta via Envelope Encryption

**No `PluginNodeType`:**
- `operations?: Record<string, OperationDefinition>` — mapa operation → HTTP request

### Interfaces novas a criar no SDK:

```typescript
interface PluginMetadata {
  brandIcon?: string;
  color?: string;
  docs?: string;
}

interface PluginCredentialDefinition {
  id: string;
  name: string;
  fields: CredentialFieldDefinition[];
  test?: CredentialTestDefinition;
}

interface CredentialFieldDefinition {
  name: string;
  displayName: string;
  type: 'string';
  required: boolean;
  isSecret: boolean;
  hint?: string;
}

interface CredentialTestDefinition {
  method: 'GET' | 'POST';
  path: string;
}

interface OperationDefinition {
  request: {
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
    path: string;
  };
  output?: {
    dataPath?: string;
  };
}
```

---

## 6. Categorias de Plugin

Categorias descrevem o que o plugin **faz**, não que ele é uma "integração":

| Categoria | Exemplos |
|-----------|----------|
| `messaging` | Telegram, Slack, Discord, WhatsApp |
| `email` | Gmail, SendGrid, Mailgun |
| `ai` | OpenAI, Claude, Gemini |
| `payments` | Stripe, PayPal |
| `databases` | PostgreSQL, MySQL, MongoDB |
| `storage` | S3, Google Drive, Dropbox |
| `crm` | HubSpot, Salesforce |
| `analytics` | Google Analytics, Mixpanel |
| `devops` | GitHub, GitLab, Jira |
| `custom` | Plugins do usuário |

Categorias core permanecem: `triggers`, `logic`, `state`, `flow_control`, `timers`, `observability`, `annotations`.

---

## 7. Credentials no JSON do Workflow

Quando o user configura um node com `isSecret: true`, o JSON salvo fica:

```json
{
  "config": {
    "operation": "sendText",
    "chatId": { "type": "event", "value": "payload.chat_id" },
    "text": { "type": "literal", "value": "Hello!" },

    "botToken": "123*****765",
    "botTokenSecret": "U2FsdGVkX1+abc123kJhG7mNpQ9rStUvWxYz..."
  }
}
```

- `botToken` — preview mascarado (frontend exibe)
- `botTokenSecret` — ciphertext (backend decrypta via DEK da org)
- DEK fica em `org_keys` collection (por serviço, isolado)

---

## 8. Como rodar

```bash
cd workspace_js/services/plugin-marketplace-mock
npm install
npm run dev

# Testar:
curl http://localhost:3099/plugins/registry.json
curl http://localhost:3099/plugins/telegram/manifest.json
```

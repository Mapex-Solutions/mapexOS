# Plugin Marketplace — Arquitetura Completa

> Como plugins de integração são publicados, descobertos, instalados e executados no MapexOS.

---

## 1. Visão Geral do Sistema

```
┌─────────────────────────────────────────────────────────────────────┐
│  NETLIFY CDN (público, estático)                                     │
│  web_documentation/.vitepress/public/plugins/                        │
│                                                                      │
│  registry.json          ← índice de TODOS os plugins disponíveis     │
│  telegram/manifest.json ← definição completa do plugin               │
│  telegram/icon.svg      ← ícone do plugin                            │
│  slack/manifest.json                                                 │
│  slack/icon.svg                                                      │
│  openai/manifest.json                                                │
│  ...                                                                 │
│                                                                      │
│  URL: https://docs.mapexos.com/plugins/registry.json                 │
│  Acesso: público, sem auth, cacheável                                │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           │ fetch (GET público)
                           │
┌──────────────────────────▼──────────────────────────────────────────┐
│  MAPEXOS UI (frontend)                                               │
│                                                                      │
│  Workflow Editor → Aba "Plugins"                                     │
│  ├── Installed    → lista do MongoDB (plugins da org)                │
│  └── Marketplace  → fetch registry.json do Netlify CDN               │
│                                                                      │
│  Ações:                                                              │
│  • Browse   → lê do CDN (público)                                    │
│  • Install  → POST /api/v1/plugins/install (JWT da org)              │
│  • Remove   → DELETE /api/v1/plugins/:pluginId (JWT da org)          │
│  • Config   → PATCH /api/v1/plugins/:pluginId/credentials            │
│                                                                      │
│  Após instalar:                                                      │
│  • Plugin aparece no catalog lateral (aba Workflow)                   │
│  • Disponível para drag-and-drop no canvas                           │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           │ POST (JWT auth)
                           │
┌──────────────────────────▼──────────────────────────────────────────┐
│  MAPEXOS BACKEND                                                     │
│                                                                      │
│  1. Recebe install request (pluginId + catalogUrl)                   │
│  2. Baixa manifest.json do Netlify CDN                               │
│  3. Valida schema do manifest                                        │
│  4. Salva manifest no MinIO L2 cache (orgs/{orgId}/plugins/)         │
│  5. Cria doc em MongoDB (installed_plugins collection)               │
│  6. Retorna 201 Created                                              │
│                                                                      │
│  Credentials:                                                        │
│  • Nunca tocam o Netlify                                             │
│  • Inline no JSON (fieldNameSecret, encrypted by DEK)               │
│  • DEK per org em org_keys collection                                │
│  • Decrypted em runtime pelo IntegrationExecutor                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Netlify — Catálogo Estático

### 2.1 Estrutura de Pastas

```
web_documentation/
└── .vitepress/
    └── public/
        └── plugins/
            ├── registry.json              ← índice global
            ├── telegram/
            │   ├── manifest.json          ← definição completa
            │   └── icon.svg               ← ícone 64x64
            ├── slack/
            │   ├── manifest.json
            │   └── icon.svg
            ├── openai/
            │   ├── manifest.json
            │   └── icon.svg
            ├── stripe/
            │   ├── manifest.json
            │   └── icon.svg
            └── ...
```

Após deploy no Netlify, acessíveis via:
```
GET https://docs.mapexos.com/plugins/registry.json
GET https://docs.mapexos.com/plugins/telegram/manifest.json
GET https://docs.mapexos.com/plugins/telegram/icon.svg
```

### 2.2 registry.json — Índice Global

```json
{
  "$schema": "mapex-plugin-registry/v1",
  "version": "1.0.0",
  "updatedAt": "2026-03-13T00:00:00Z",
  "plugins": [
    {
      "id": "integration-telegram",
      "name": "Telegram",
      "version": "1.0.0",
      "category": "integrations",
      "icon": "send",
      "brandIcon": "telegram/icon.svg",
      "color": "#0088CC",
      "description": "Send messages, photos and manage chats via Telegram Bot API",
      "author": "MapexOS",
      "tags": ["messaging", "chat", "bot"],
      "manifestUrl": "telegram/manifest.json",
      "docsUrl": "/docs/1.0.0/en/plugins/telegram",
      "requiresCredentials": true,
      "nodeCount": 1
    },
    {
      "id": "integration-slack",
      "name": "Slack",
      "version": "1.0.0",
      "category": "integrations",
      "icon": "tag",
      "brandIcon": "slack/icon.svg",
      "color": "#4A154B",
      "description": "Send messages, manage channels and users in Slack workspaces",
      "author": "MapexOS",
      "tags": ["messaging", "team", "collaboration"],
      "manifestUrl": "slack/manifest.json",
      "docsUrl": "/docs/1.0.0/en/plugins/slack",
      "requiresCredentials": true,
      "nodeCount": 1
    },
    {
      "id": "integration-openai",
      "name": "OpenAI",
      "version": "1.0.0",
      "category": "integrations",
      "icon": "psychology",
      "brandIcon": "openai/icon.svg",
      "color": "#10A37F",
      "description": "Generate text, images and embeddings using OpenAI models",
      "author": "MapexOS",
      "tags": ["ai", "llm", "gpt", "chatgpt"],
      "manifestUrl": "openai/manifest.json",
      "docsUrl": "/docs/1.0.0/en/plugins/openai",
      "requiresCredentials": true,
      "nodeCount": 1
    }
  ]
}
```

**Campos do registry (leve — só para listagem):**

| Campo | Tipo | Propósito |
|-------|------|-----------|
| `id` | string | Identificador único do plugin |
| `name` | string | Nome para exibição |
| `version` | string | Semver — compara com versão instalada |
| `category` | PluginCategory | Agrupamento no catalog |
| `icon` | string | Material Icons (fallback) |
| `brandIcon` | string | Path relativo ao SVG (no CDN) |
| `color` | string | Cor da marca (hex) |
| `description` | string | Descrição curta para o card |
| `author` | string | Quem publicou |
| `tags` | string[] | Para busca/filtro |
| `manifestUrl` | string | Path relativo ao manifest.json completo |
| `docsUrl` | string | Link para documentação do plugin |
| `requiresCredentials` | boolean | Se precisa configurar credenciais |
| `nodeCount` | number | Quantos node types o plugin registra |

### 2.3 manifest.json — Definição Completa (por plugin)

O manifest é o **JSON completo** que o ManifestLoader converte em `WorkflowPlugin`.
Formato definido em `manifest-schema-v1.md`. Exemplo Telegram em `telegram_plugin.md`.

**Relação entre registry.json e manifest.json:**
```
registry.json = RESUMO (para listagem rápida, sem baixar tudo)
manifest.json = COMPLETO (baixado só quando instala)
```

---

## 3. Trigger vs Integration — Diferença Fundamental

```
TRIGGER (sistema atual)                    INTEGRATION (plugin node no workflow)
───────────────────────                    ───────────────────────────────────
1 credential → 1 chamada                  1 credential → N chamadas (sessão)
Atômico, standalone                        Compartilhado entre steps do DAG
Genérico (todos usam igual)                Múltiplas operações no mesmo fluxo
Config: baixo nível (URL + headers)        Config: alto nível (Resource + Operation)
Execução: Triggers Service (direto)        Execução: Triggers Service (via NATS)

Exemplo:                                   Exemplo:
"Alerta Telegram"                          Workflow com 3 steps Telegram:
  → POST sendMessage                         Step 1: sendText     ──┐
  → FIM                                      Step 2: sendPhoto    ──┼── MESMA credential
                                             Step 3: getChat      ──┘
                                             (3 HTTP calls, 1 bot token resolvido 1x)
```

**Session sharing no workflow — via instance STATE (NATS KV):**

Credentials vivem no JSON de cada serviço (DDD). No workflow, quando múltiplos
nodes do mesmo plugin executam, a credential é decrypted 1x e cacheada no
instance state. Quando o workflow termina, o state é limpo do KV — zero rastro.

```
DAG Execution:
  Node "Telegram: Send Text"
    → Decrypt botToken do node config (BYOK)
    → Cache no instance state: state["credential:telegram"] = decrypted
    → Monta HTTP → Triggers Service executa

  Node "Telegram: Send Photo"
    → Lê state["credential:telegram"] (já decrypted) ← zero decrypt
    → Monta HTTP → Triggers Service executa

  Node "Telegram: Get Chat"
    → Lê state["credential:telegram"] ← zero decrypt
    → Monta HTTP → Triggers Service executa

  Workflow termina → NATS KV limpo → credentials somem da memória
```

---

## 4. Credentials — Envelope Encryption (BYOK)

### 4.1 Princípio: Separação de Responsabilidades (DDD)

Cada serviço é **dono exclusivo** dos seus segredos. Nenhum serviço acessa dados do outro.

```
TRIGGER SERVICE                         WORKFLOW SERVICE
(dono dos triggers)                     (dono dos workflows)
───────────────────                     ───────────────────
Collection: org_keys                    Collection: org_keys
  → DEK por org, encrypted by master     → DEK por org, encrypted by master
  → ISOLADO do Workflow Service           → ISOLADO do Trigger Service

Trigger config JSON {                   Node config JSON {
  endpoint: "...",                        resource: "message",
  method: "POST",                         operation: "sendText",
  headers: {                              chatId: { type: "event", ... },
    "Auth": "{{botToken}}"                botToken: "123*****765",
  },                                      botTokenSecret: "encrypted..."
  botToken: "123*****765",              }
  botTokenSecret: "encrypted..."
}

Trigger Service encrypta/decrypta      Workflow Service encrypta/decrypta
usando SUA org_keys collection          usando SUA org_keys collection
```

### 4.2 Envelope Encryption — Arquitetura de Chaves

Padrão usado por AWS KMS, Google Cloud KMS, HashiCorp Vault:

```
┌─────────────────────────────────────────────────────────────┐
│  CAMADA 1 — Master Key (ENV, nunca sai do server)           │
│  MAPEX_CREDENTIAL_MASTER_KEY = "base64-encoded-32-bytes"    │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  CAMADA 2 — DEK (Data Encryption Key)               │    │
│  │  1 DEK por organização (orgId)                       │    │
│  │  Armazenado em: org_keys collection (per service)    │    │
│  │  Encrypted by: Master Key (AES-256-GCM)              │    │
│  │                                                      │    │
│  │  ┌──────────────────────────────────────────────┐   │    │
│  │  │  CAMADA 3 — Secret (API key, token, etc)     │   │    │
│  │  │  Encrypted by: DEK (AES-256-GCM)             │   │    │
│  │  │  Stored: inline no JSON como fieldNameSecret  │   │    │
│  │  └──────────────────────────────────────────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 Collection: `org_keys` (per service)

Cada serviço que usa BYOK tem sua **própria** collection `org_keys`:

```go
// Go struct
type OrgKey struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    OrgID        string             `bson:"orgId"`
    EncryptedDEK []byte             `bson:"encryptedDEK"` // DEK encrypted by master key
    Algorithm    string             `bson:"algorithm"`     // "AES-256-GCM"
    Version      int                `bson:"version"`       // for key rotation
    CreatedAt    time.Time          `bson:"createdAt"`
    RotatedAt    time.Time          `bson:"rotatedAt"`
}

// Índices
// unique: { orgId, version }
// index:  { orgId }  (busca rápida por org)
```

```
Workflow Service DB:
  ├── workflow_definitions   ← workflows
  ├── workflow_instances     ← executions
  ├── installed_plugins      ← plugins instalados por org
  └── org_keys               ← DEKs por org (Workflow Service)

Trigger Service DB:
  ├── triggers               ← triggers
  ├── trigger_executions     ← logs
  └── org_keys               ← DEKs por org (Trigger Service)
```

**DEK não fica no JSON** — fica exclusivamente em `org_keys`. O campo `fieldNameSecret` no JSON
só tem o ciphertext (string base64). O backend busca o DEK pela orgId na `org_keys` quando precisa decryptar.

### 4.4 `isSecret` Flag no Manifest

No manifest do plugin, properties com dados sensíveis declaram `isSecret: true`:

```json
{
  "credentials": {
    "properties": [
      {
        "name": "botToken",
        "displayName": "Bot Token",
        "type": "string",
        "isSecret": true,
        "required": true,
        "description": "Telegram Bot API token from @BotFather"
      }
    ]
  }
}
```

Quando o backend recebe um node config com campo marcado `isSecret: true`:
1. Encrypta o valor
2. Gera preview mascarado
3. Salva os dois campos no JSON

### 4.5 Formato no JSON — Dois Campos por Secret

Para cada property com `isSecret: true`, o JSON tem **dois campos**:

```json
{
  "id": "node_tg_001",
  "type": "integration/telegram",
  "label": "Send Welcome",
  "config": {
    "resource": "message",
    "operation": "sendText",
    "chatId": { "type": "event", "value": "payload.chat_id" },
    "text": { "type": "state", "value": "welcomeMsg" },
    "parseMode": "HTML",

    "botToken": "123*****765",
    "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
  }
}
```

| Campo | Propósito | Quem lê |
|-------|-----------|---------|
| `botToken` | Preview mascarado (ex: `"123*****765"`) | Frontend (UI exibe) |
| `botTokenSecret` | Ciphertext base64 (string pura, encrypted pelo DEK) | Backend (runtime decrypta) |

**O campo `Secret` NÃO contém metadados de chave.** Algoritmo, versão, orgId — tudo fica
em `org_keys` no MongoDB. O backend resolve: `orgId` (do JWT) → busca DEK em `org_keys` → decrypta.

### 4.6 Flows Completos

#### CRIAR (user configura credential no node):

```
1. User no UI digita: botToken = "123456:ABC-DEF-GHI765"
2. UI envia PATCH /api/v1/workflows/:id
   body: { nodes: [{ config: { botToken: "123456:ABC-DEF-GHI765", ... } }] }

3. Backend detecta: property "botToken" tem isSecret: true no manifest
4. Backend busca DEK: db.org_keys.findOne({ orgId, version: latest })
   → Se não existe: gera novo DEK, encrypta com master key, salva em org_keys
5. Decrypta DEK: AES-256-GCM-decrypt(masterKey, encryptedDEK) → DEK
6. Encrypta valor: AES-256-GCM(DEK, "123456:ABC-DEF-GHI765") → cipherText
7. Gera preview: "123*****765" (primeiros 3 chars + ***** + últimos 3 chars)
8. Salva no JSON do node config:
   {
     "botToken": "123*****765",
     "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
   }
9. Retorna ao frontend com o preview (nunca retorna botTokenSecret)
```

#### LER (UI carrega o workflow):

```
1. Frontend faz GET /api/v1/workflows/:id
2. Backend retorna JSON com:
   - botToken: "123*****765"     ← frontend exibe isso
   - botTokenSecret: OMITIDO    ← backend NUNCA envia Secret fields para o frontend
3. UI renderiza input com "123*****765" (read-only, com botão [Edit])
```

#### EXECUTAR (runtime do workflow):

```
1. IntegrationExecutor lê node config do MongoDB
2. Encontra campo "botTokenSecret" (ciphertext string)
3. Busca DEK: db.org_keys.findOne({ orgId }) → encryptedDEK
4. Decrypta DEK: AES-256-GCM-decrypt(masterKey, encryptedDEK) → DEK
5. Decrypta valor: AES-256-GCM-decrypt(DEK, ct) → "123456:ABC-DEF-GHI765"
6. Cache no instance state: state["credential:telegram"] = { botToken: "123456:..." }
7. Monta HTTP request → envia via NATS ao Triggers Service
8. Próximos nodes do mesmo tipo leem do state (zero decrypt adicional)
9. Workflow termina → NATS KV limpo → credentials somem da memória
```

#### ROTAÇÃO de DEK (per org):

```
1. Admin trigger: rotate DEK for orgId
2. Gera novo DEK (version N+1)
3. Encrypta novo DEK com master key → salva em org_keys (version: N+1)
4. Re-encrypta todos os Secret fields dos workflows dessa org:
   - Decrypta ct com DEK antigo (version N)
   - Re-encrypta ct com DEK novo (version N+1)
   - Atualiza "v" no JSON: v: N → v: N+1
5. DEK antigo (version N) mantido por período de grace (ex: 30 dias)
6. Após grace period: remove DEK antigo de org_keys
```

#### ROTAÇÃO de Master Key:

```
1. Ops: atualiza ENV MAPEX_CREDENTIAL_MASTER_KEY
2. Job re-encrypta TODOS os DEKs em org_keys:
   - Decrypta encryptedDEK com master key antigo
   - Re-encrypta com master key novo
   - Salva encryptedDEK novo em org_keys
3. Os Secret fields NÃO mudam (ct permanece igual, encrypted pelo mesmo DEK)
4. Apenas a camada org_keys é re-encryptada
```

### 4.7 Como Fica no JSON do Workflow (salvo no MongoDB)

```json
{
  "nodes": [
    {
      "id": "node_tg_001",
      "type": "integration/telegram",
      "label": "Send Welcome",
      "position": { "x": 400, "y": 200 },
      "config": {
        "resource": "message",
        "operation": "sendText",
        "chatId": { "type": "event", "value": "payload.chat_id" },
        "text": { "type": "state", "value": "welcomeMsg" },
        "parseMode": "HTML",

        "botToken": "123*****765",
        "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
      }
    },
    {
      "id": "node_tg_002",
      "type": "integration/telegram",
      "label": "Send Photo",
      "position": { "x": 400, "y": 400 },
      "config": {
        "resource": "message",
        "operation": "sendPhoto",
        "chatId": { "type": "event", "value": "payload.chat_id" },
        "photo": { "type": "literal", "value": "https://..." },

        "botToken": "123*****765",
        "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
      }
    }
  ]
}
```

**Nota:** Os dois nodes Telegram têm o mesmo botToken (mesmo bot).
Em runtime, o primeiro node decrypta e coloca no state. O segundo lê do state.

### 4.8 Como Fica no Trigger (mesmo padrão BYOK, domínio diferente)

```json
{
  "_id": "trigger_001",
  "name": "Alerta Telegram",
  "triggerType": "http",
  "config": {
    "endpoint": "https://api.telegram.org/bot{{botToken}}/sendMessage",
    "method": "POST",
    "body": { "chat_id": "123", "text": "Alerta!" },

    "botToken": "123*****765",
    "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
  }
}
```

**Cada serviço trata seus próprios secrets:**
- Triggers Service decrypta usando SUA collection `org_keys`
- Workflow Service decrypta usando SUA collection `org_keys`
- Cada serviço tem sua instância isolada — zero acoplamento

### 4.9 UI — Credential no Node Config

```
┌─────────────────────────────────────────────────┐
│  TELEGRAM                                        │
│──────────────────────────────────────────────────│
│                                                  │
│  Bot Token    [123*****765      ] [Edit] [Test]  │  ← mostra preview
│                                                  │
│  Resource     [▼ Message                  ]      │
│  Operation    [▼ Send Text                ]      │
│  Chat ID      [event ▼] [payload.chat_id ]       │
│  Text         [literal ▼] [Hello!        ]       │
│                                                  │
└─────────────────────────────────────────────────┘

Ao clicar [Edit]:
┌─────────────────────────────────────────────────┐
│  Enter new Bot Token:                            │
│  [________________________________]              │
│                     [Cancel] [Save]              │
└─────────────────────────────────────────────────┘
→ Backend encrypta → salva no JSON → retorna preview "789*****XYZ"
```

### 4.10 Múltiplas Contas (N Instagrams)

Cada workflow pode ter um bot/conta diferente. A credential vive inline no JSON do node:

```
Workflow A: posts em @empresa_oficial
  Node Instagram config.accessToken = "abc*****123"
  Node Instagram config.accessTokenSecret = "encrypted..."

Workflow B: posts em @empresa_promo
  Node Instagram config.accessToken = "xyz*****789"
  Node Instagram config.accessTokenSecret = "encrypted..."

Workflow C: posts em @empresa_oficial (mesmo token — re-digitado)
  Node Instagram config.accessToken = "abc*****123"
  Node Instagram config.accessTokenSecret = "encrypted..."
```

Se o token A muda → user edita nos workflows que usam (A e C).
Trade-off: duplicação vs simplicidade. Zero collection extra de credentials, zero cross-service.

---

## 5. Dados — MongoDB Collections

### 5.1 Collection: `installed_plugins`

Armazena quais plugins cada org instalou. Única collection nova necessária.

```typescript
interface InstalledPlugin {
  _id: ObjectId
  orgId: string                          // Organização que instalou
  pluginId: string                       // "integration-telegram"
  version: string                        // "1.0.0" — versão instalada
  name: string                           // "Telegram" (snapshot no momento da instalação)
  category: string                       // "integrations"
  icon: string                           // "send"
  brandIcon: string                      // URL completa do CDN
  color: string                          // "#0088CC"
  description: string                    // Descrição curta

  // Controle de acesso
  shared: boolean                        // true = todos da org veem, false = só quem instalou
  installedBy: string                    // userId que instalou

  // Metadados
  catalogUrl: string                     // URL do registry de onde veio
  manifestUrl: string                    // URL do manifest no CDN
  minioPath: string                      // Path no MinIO: "orgs/{orgId}/plugins/{pluginId}/manifest.json"

  // Status
  status: 'active' | 'disabled' | 'update_available'
  latestVersion?: string                 // Preenchido quando CDN tem versão mais nova

  // Timestamps
  installedAt: Date
  updatedAt: Date

  // Índices
  // unique: { orgId, pluginId }
  // index: { orgId, status }
  // index: { orgId, shared }
}
```

### 5.2 Collection: `plugin_usage` (opcional, analytics)

```typescript
interface PluginUsage {
  _id: ObjectId
  orgId: string
  pluginId: string
  workflowId: string                     // Qual workflow usa este plugin
  nodeCount: number                      // Quantos nodes desse plugin no workflow
  lastUsedAt: Date
}
```

---

## 6. Backend — Endpoints

### 6.1 Plugin Management (REST API)

Base path: `/api/v1/plugins`

| Method | Path | Auth | Descrição |
|--------|------|------|-----------|
| `GET` | `/` | JWT | Lista plugins instalados na org |
| `POST` | `/install` | JWT | Instala plugin do marketplace na org |
| `DELETE` | `/:pluginId` | JWT | Remove plugin da org |
| `PATCH` | `/:pluginId` | JWT | Atualiza config (shared, status) |
| `POST` | `/:pluginId/check-update` | JWT | Verifica se há versão nova no CDN |
| `POST` | `/:pluginId/update` | JWT | Atualiza para versão mais recente |

### 6.2 Credential dentro do Workflow (não é endpoint separado)

Credentials vivem inline no JSON do workflow (campos `fieldName` + `fieldNameSecret`).
Não existe endpoint `/api/v1/credentials`. O CRUD é parte do workflow CRUD:

```
PATCH /api/v1/workflows/:workflowId
body: { nodes: [...] }  ← inclui campos secret como plaintext (novo valor do user)

O backend ao receber:
1. Cruza properties do manifest: detecta campos com isSecret: true
2. Se campo tem valor plaintext (não mascarado) → encrypta via Envelope Encryption
3. Gera preview mascarado (ex: "123*****765")
4. Salva no MongoDB: fieldName = preview, fieldNameSecret = ciphertext (string base64)
5. Retorna ao frontend SEM os campos Secret (apenas o preview)
```

### 6.3 Flow de Instalação (detalhe)

```
POST /api/v1/plugins/install
Headers: Authorization: Bearer <jwt>
Body: {
  "pluginId": "integration-telegram",
  "catalogUrl": "https://docs.mapexos.com/plugins",
  "shared": true
}
```

**Backend steps:**
```
1. Valida JWT → extrai orgId, userId
2. Verifica permissão: org.plugins.install
3. Checa se já instalado: db.installed_plugins.findOne({ orgId, pluginId })
   → Se sim: retorna 409 Conflict
4. Baixa manifest: GET {catalogUrl}/{pluginId}/manifest.json
   → Se falha: retorna 502 Bad Gateway
5. Valida manifest JSON schema
   → Se inválido: retorna 422 Unprocessable Entity
6. Salva manifest no MinIO: orgs/{orgId}/plugins/{pluginId}/manifest.json
7. Cria doc em MongoDB: installed_plugins
8. Retorna 201 Created com InstalledPlugin
```

### 6.4 Flow de Remoção

```
DELETE /api/v1/plugins/:pluginId
```

**Backend steps:**
```
1. Valida JWT → extrai orgId
2. Verifica permissão: org.plugins.uninstall
3. Checa se plugin está em uso em algum workflow:
   → Se sim: retorna 409 Conflict com lista de workflows
   → (Opção: force=true no query param para forçar)
4. Remove manifest do MinIO (L2 cache)
5. Remove doc de installed_plugins
6. Retorna 200 OK
   (Credentials vivem nos workflows — não há collection separada para limpar)
```

---

## 7. UI — Fluxo do Usuário

### 7.1 Nova Aba "Plugins" no Workflow Editor

```
Workflow Editor Tabs:
┌─────────┬──────┬──────────┬─────────┬───────┐
│ General │ Data │ Workflow │ Plugins │ Debug │
└─────────┴──────┴──────────┴─────────┴───────┘
                              ▲ NOVA
```

### 7.2 Dentro da Aba "Plugins" — Duas Views

```
┌─────────────────────────────────────────────────────────────────┐
│  PLUGINS                                                         │
│                                                                  │
│  ┌────────────────┐  ┌────────────────┐                          │
│  │  🔌 Installed  │  │  🛒 Marketplace │    ← sub-tabs           │
│  └────────────────┘  └────────────────┘                          │
│                                                                  │
│ ─────────────────────────────────────────────────────────────── │
│                                                                  │
│  VIEW: INSTALLED                                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  ✅ Telegram            v1.0.0    Shared    [⚙] [🗑]     │   │
│  │     Send messages via Telegram Bot API                    │   │
│  │     Credentials: My Bot (verified ✓)                      │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │  ✅ OpenAI              v1.0.0    Private   [⚙] [🗑]     │   │
│  │     Generate text using OpenAI models                     │   │
│  │     Credentials: ⚠ Not configured                         │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  VIEW: MARKETPLACE                                               │
│  ┌─────────────┐                                                 │
│  │ 🔍 Search plugins...                                    │    │
│  └─────────────┘                                                 │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐               │
│  │Telegram │ │  Slack  │ │ OpenAI  │ │ Stripe  │               │
│  │  📨     │ │  💬     │ │  🧠     │ │  💳     │               │
│  │ v1.0.0  │ │ v1.0.0  │ │ v1.0.0  │ │ v1.0.0  │               │
│  │         │ │         │ │         │ │         │               │
│  │[Install]│ │[Install]│ │[Installed]│[Install]│               │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘               │
│                                                                  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐               │
│  │ Gmail   │ │ Discord │ │ Notion  │ │Airtable │               │
│  │  ✉️     │ │  🎮     │ │  📝     │ │  📊     │               │
│  │ v1.0.0  │ │ v1.0.0  │ │ v1.0.0  │ │ v1.0.0  │               │
│  │         │ │         │ │         │ │         │               │
│  │[Install]│ │[Install]│ │[Install]│ │[Install]│               │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

### 7.3 Flow: Browse → Install → Use

```
PASSO 1: Browse
─────────────────
User abre aba "Plugins" → sub-tab "Marketplace"
UI faz: GET https://docs.mapexos.com/plugins/registry.json
Renderiza grid de cards com nome, ícone, descrição, versão
Compara com installed_plugins do DB → marca como "Installed" ou "Install"

PASSO 2: Install
─────────────────
User clica "Install" no card do Telegram
UI faz: POST /api/v1/plugins/install
  body: { pluginId: "integration-telegram", catalogUrl: "...", shared: true }
Backend baixa manifest, salva no MinIO, cria doc no MongoDB
UI recebe 201 → card muda para "Installed ✓"

PASSO 3: Configure Credentials (dentro do workflow)
─────────────────────────────────────────────────────
User vai para aba "Workflow" → arrasta node Telegram → abre config
Campos com isSecret: true aparecem como input de senha:
  - Bot Token: [                    ] (input vazio, tipo password)
  - [Test Connection]
User digita o token → salva o workflow
UI faz: PATCH /api/v1/workflows/:id
  body: { nodes: [{ config: { botToken: "123456:ABC-DEF...", ... } }] }
Backend detecta isSecret → encrypta → salva preview + fieldNameSecret no JSON

PASSO 4: Use in Workflow
─────────────────────────
User vai para aba "Workflow"
O ManifestLoader carregou o manifest do MinIO
  → pluginRegistry.registerPlugin(telegramPlugin)
  → Plugin aparece na categoria "Integrations" do catalog lateral
User arrasta "Telegram" para o canvas
Configura: resource, operation, chatId, text...
Salva workflow
```

### 7.4 Flow: Plugin Aparece no Catalog Lateral

```
ANTES de instalar:                    DEPOIS de instalar:

Catalog (lateral):                    Catalog (lateral):
├── Triggers                          ├── Triggers
│   └── Trigger Event                 │   └── Trigger Event
├── Logic & Conditions                ├── Logic & Conditions
│   └── Conditions                    │   └── Conditions
├── Data                              ├── Data
│   ├── Set State                     │   ├── Set State
│   ├── Log                           │   ├── Log
│   └── Code                          │   └── Code
├── Flow Control                      ├── Flow Control
│   ├── Fanout                        │   ├── Fanout
│   ├── Switch                        │   ├── Switch
│   └── ...                           │   └── ...
├── Timers                            ├── Timers
│   ├── Delay                         │   ├── Delay
│   └── ...                           │   └── ...
├── Integrations                      ├── Integrations
│   └── (vazio)                       │   ├── Telegram  ← APARECEU!
│                                     │   └── OpenAI    ← APARECEU!
└── Annotations                       └── Annotations
    └── Group Frame                       └── Group Frame
```

---

## 8. Frontend — Carregamento dos Plugins

### 8.1 Boot Sequence (app startup)

```typescript
// Em CreateEditWorkflowPage.vue ou store init

async function loadPlugins() {
  // 1. Carrega plugins CORE (sempre disponíveis)
  bootWorkflowPlugins(pluginRegistry.registerPlugin)

  // 2. Busca plugins instalados na org (do backend)
  const installedPlugins = await pluginsApi.list()  // GET /api/v1/plugins

  // 3. Para cada plugin instalado, carrega o manifest do MinIO (via backend proxy)
  for (const plugin of installedPlugins) {
    const manifest = await pluginsApi.getManifest(plugin.pluginId)
    // GET /api/v1/plugins/:pluginId/manifest
    // Backend lê do MinIO e retorna o JSON

    // 4. ManifestLoader converte JSON → WorkflowPlugin
    const workflowPlugin = ManifestLoader.load(manifest)

    // 5. Registra no pluginRegistry
    pluginRegistry.registerPlugin(workflowPlugin)
  }
}
```

### 8.2 ManifestLoader — Conversão JSON → WorkflowPlugin

```typescript
// workspace_js/packages/workflow-sdk/src/loaders/ManifestLoader.ts

import type { WorkflowPlugin, PluginNodeType } from '../interfaces'
import { markRaw } from 'vue'
import { GenericWorkflowNode } from '@mapexos/workflow-plugin-core'

class ManifestLoader {
  /**
   * Converte manifest JSON estático em WorkflowPlugin registrável
   *
   * @param manifest - JSON parseado do manifest.json
   * @returns WorkflowPlugin pronto para pluginRegistry.registerPlugin()
   */
  static load(manifest: PluginManifest): WorkflowPlugin {
    return {
      id: manifest.id,
      name: manifest.name,
      version: manifest.version,
      category: manifest.category,
      icon: manifest.icon,

      nodeTypes: manifest.nodeTypes.map(nt => ({
        ...nt,
        configSchema: {},
        canvasComponent: markRaw(GenericWorkflowNode),
        // Não tem configComponent custom → usa DynamicNodeForm (properties[])
        // _pluginMeta carrega dados extras para o executor backend
        _pluginMeta: {
          credentials: manifest.credentials,
          baseUrl: manifest.baseUrl,
          operations: nt.operations,
        },
      })),

      onActivate(context) {
        // Registra translations default (do manifest)
        context.registerTranslations('en-US', {
          nodes: Object.fromEntries(
            manifest.nodeTypes.map(nt => [
              nt.type.split('/')[1],
              { label: nt.label, description: nt.description }
            ])
          )
        })
      },
    }
  }
}
```

---

## 9. Segurança

### 9.1 O Que é Público vs Privado

| Dado | Onde | Acesso |
|------|------|--------|
| registry.json (lista de plugins) | Netlify CDN | Público |
| manifest.json (definição do plugin) | Netlify CDN | Público |
| icon.svg | Netlify CDN | Público |
| installed_plugins (quais plugins a org tem) | MongoDB | Autenticado (orgId) |
| org_keys (DEK por org, encrypted by master key) | MongoDB (per service) | Backend only, nunca exposto |
| `fieldNameSecret` (encrypted by DEK) | Dentro do workflow JSON (MongoDB) | Backend only, nunca enviado ao frontend |
| `fieldName` preview (ex: `"123*****765"`) | Dentro do workflow JSON → Frontend | Safe to display |
| workflow definition JSON | MongoDB (source of truth), MinIO (L2 cache) | Autenticado (orgId) |
| `MAPEX_CREDENTIAL_MASTER_KEY` | ENV (server only) | Nunca sai do server, nunca no JSON |

### 9.2 Por Que Público é Seguro

O manifest.json contém **apenas a definição** — como um manual de instruções:
- Quais campos o plugin tem
- Quais endpoints a API usa
- Como montar o request

**Não contém:**
- Nenhum token ou API key
- Nenhum dado de nenhuma org
- Nenhuma informação sensível

É como publicar a documentação de uma API — o manual é público, a chave de acesso é privada.

### 9.3 Permissões no Backend

```
Ações de Plugin:
  org.plugins.browse      → ver marketplace (qualquer membro)
  org.plugins.install     → instalar plugin (admin)
  org.plugins.uninstall   → remover plugin (admin)

Ações de Workflow (já existentes — credentials são parte do workflow):
  org.workflows.create    → criar workflow (inclui configurar credentials dos nodes)
  org.workflows.update    → atualizar workflow (inclui editar credentials)
  org.workflows.read      → ler workflow (vê apenas preview mascarado)
```

---

## 10. Shared vs Private — Comportamento

### 10.1 Plugin Shared (shared: true)

```
Admin instala Telegram com shared: true
  → TODOS os membros da org veem Telegram no catalog
  → TODOS podem usar em seus workflows
  → Credentials podem ser shared ou per-user
```

### 10.2 Plugin Private (shared: false)

```
User instala Telegram com shared: false
  → SOMENTE esse user vê Telegram no catalog
  → SOMENTE esse user pode usar em workflows
  → Credentials são always private
```

### 10.3 Filtro no Frontend

```typescript
// Ao carregar plugins instalados
const plugins = await pluginsApi.list()
// Backend filtra: WHERE orgId = jwt.orgId AND (shared = true OR installedBy = jwt.userId)
```

---

## 11. Versionamento e Updates

### 11.1 Check for Updates

```
UI periodicamente (ou no load):
1. Busca registry.json do CDN (tem version de cada plugin)
2. Compara com version em installed_plugins no MongoDB
3. Se CDN.version > installed.version → marca status: 'update_available'
4. UI mostra badge "Update available" no card
```

### 11.2 Update Flow

```
User clica "Update" no plugin
POST /api/v1/plugins/:pluginId/update
Backend:
1. Baixa novo manifest.json do CDN
2. Valida schema
3. Substitui manifest no MinIO
4. Atualiza version no MongoDB
5. Frontend recarrega o plugin (unregister → register novo)
```

---

## 12. Resumo — O Que Implementar

### Fase 1: UI Only (teste sem backend)

| # | O quê | Onde |
|---|-------|-----|
| 1 | Criar JSON estático do Telegram plugin | VitePress public/ |
| 2 | Criar registry.json | VitePress public/ |
| 3 | Criar aba "Plugins" no workflow editor | CreateEditWorkflowPage |
| 4 | Sub-tab "Marketplace" → fetch registry.json do Netlify | PluginsTab.vue |
| 5 | Sub-tab "Installed" → lista mock (hardcoded) | PluginsTab.vue |
| 6 | ManifestLoader básico → registra no pluginRegistry | workflow-sdk |
| 7 | Plugin aparece no catalog lateral | Automático (pluginRegistry) |

### Fase 2: Backend — Plugins + BYOK

| # | O quê | Onde |
|---|-------|-----|
| 8 | Collections `installed_plugins` + `org_keys` | MongoDB (Workflow Service) |
| 9 | BYOK crypto package (encrypt/decrypt/mask) | Go shared package |
| 10 | `isSecret` handling + Envelope Encryption no workflow CRUD | Workflow Service (definition_service) |
| 11 | Endpoints REST plugins (install, remove, list) | Workflow Service |
| 12 | MinIO storage para manifests instalados | Workflow Service |
| 13 | Permissões (org.plugins.*) | Permissions package |
| 14 | ENV `MAPEX_CREDENTIAL_MASTER_KEY` | Docker/config |

### Fase 3: Runtime — Execution + Session

| # | O quê | Onde |
|---|-------|-----|
| 15 | IntegrationExecutor decrypt credential Secret fields no runtime | Workflow Service |
| 16 | Session cache no instance state (NATS KV) | Workflow Runtime |
| 17 | Credential fields no DynamicNodeForm (masked input + edit) | Frontend |
| 18 | BYOK para Triggers (mesmo padrão, domínio do Triggers Service) | Triggers Service |

### Fase 4: Produção

| # | O quê | Onde |
|---|-------|-----|
| 19 | Master key rotation job (re-encrypta DEKs em ambos serviços) | Ops |
| 20 | Version check + auto-update | Frontend + Backend |
| 21 | Plugin usage tracking | MongoDB |
| 22 | Publicar 10+ plugins no CDN | VitePress public/ |

---

## 13. Referência — Arquivos Existentes

| Arquivo | Propósito |
|---------|-----------|
| `workflow-sdk/src/interfaces/workflowPlugin.interface.ts` | WorkflowPlugin, PluginNodeType, NodePropertyDefinition |
| `workflow-sdk/src/interfaces/fieldSource.interface.ts` | FieldSourceValue, SourceType |
| `workflow-plugin-core/src/constants/corePlugins.constant.ts` | 6 core plugins, bootWorkflowPlugins() |
| `stores/pluginRegistry/` | registerPlugin, catalog getter, nodeTypeMap |
| `createEditWorkflowPage/CreateEditWorkflowPage.vue` | Tabs do editor (General, Data, Workflow, Debug) |
| `createEditWorkflowPage/components/PluginCatalog/` | Catalog lateral na aba Workflow |
| `createEditWorkflowPage/components/DynamicNodeForm/` | Form renderer para properties[] |
| `implementations/manifest-schema-v1.md` | Schema V1 dos manifests |
| `implementations/telegram_plugin.md` | Exemplo completo Telegram end-to-end |

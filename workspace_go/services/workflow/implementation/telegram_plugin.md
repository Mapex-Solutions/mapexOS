# Telegram Plugin — Implementação Concreta

> Como um plugin de integração se encaixa no sistema existente do MapexOS, passo a passo.

---

## 1. O Que Já Existe (sem mudar nada)

```
WorkflowPlugin interface           → registro no pluginRegistry
PluginNodeType.properties[]        → DynamicNodeForm renderiza (5 types: string, number, boolean, options, json)
PluginNodeType.configComponent     → Vue component custom (para UIs complexas)
FieldSourceSelector                → componente que permite literal/state/event/nodeOutput
NodeConfigPanel                    → 3 tiers: properties → configComponent → JSON fallback
```

**Padrão CORE atual (exemplo: Log node):**
```typescript
{
  type: 'core/log',
  label: 'Log',
  icon: 'article',
  color: 'teal-7',
  description: 'Emit observability event',
  inputs: [{ id: 'in', label: 'In', position: 'top' }],
  outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
  configSchema: {},
  canvasComponent: markRaw(GenericWorkflowNode),
  properties: [
    { name: 'message', displayName: 'Message', type: 'string', default: '' },
    { name: 'level', displayName: 'Level', type: 'options', default: 'info',
      options: [
        { label: 'Info', value: 'info' },
        { label: 'Warn', value: 'warn' },
      ]
    },
  ],
  defaults: { message: '', level: 'info' },
}
```

---

## 2. O Plugin Telegram — Visão Geral

```
┌─────────────────────────────────────────────────────────────────────┐
│  MinIO Storage                                                      │
│  plugins/telegram/                                                   │
│  ├── manifest.json          ← O JSON que define TUDO                │
│  ├── telegram.svg           ← Ícone                                 │
│  └── scripts/               ← Scripts JS (se necessário)            │
│      └── loadOptions/                                                │
│          (nenhum para Telegram — tudo é ID direto)                   │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│  ManifestLoader (runtime JS)                                        │
│  Lê manifest.json → converte para WorkflowPlugin → registra        │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│  pluginRegistry.registerPlugin(telegramPlugin)                      │
│  → nodeTypeMap registra 'integration/telegram'                      │
│  → catalog exibe na categoria 'integrations'                        │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
              ┌───────────┴───────────┐
              ▼                       ▼
        EDITOR (UI)            EXECUTION (Go)
        NodeConfigPanel        IntegrationExecutor
        DynamicNodeForm        → constrói pipeline
        FieldSourceSelector    → despacha via NATS
```

---

## 3. O Manifest JSON — Telegram

```json
{
  "$schema": "mapex-plugin/v1",
  "id": "integration-telegram",
  "name": "Telegram",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "send",

  "metadata": {
    "brandIcon": "telegram.svg",
    "color": "#0088CC",
    "docs": "https://core.telegram.org/bots/api"
  },

  "credentials": {
    "id": "telegramApi",
    "name": "Telegram Bot API",
    "fields": [
      {
        "name": "botToken",
        "displayName": "Bot Token",
        "type": "string",
        "required": true,
        "secret": true,
        "hint": "Obtenha com @BotFather no Telegram"
      }
    ],
    "inject": {
      "path": { "token": "{{botToken}}" }
    },
    "test": {
      "method": "GET",
      "path": "/bot{{botToken}}/getMe"
    }
  },

  "baseUrl": "https://api.telegram.org",

  "nodeTypes": [
    {
      "type": "integration/telegram",
      "label": "Telegram",
      "icon": "send",
      "color": "light-blue-7",
      "description": "Send messages via Telegram Bot API",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [
        { "id": "success", "label": "Success", "position": "bottom", "color": "#4caf50" },
        { "id": "error", "label": "Error", "position": "bottom", "color": "#ef5350" }
      ],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "message",
          "options": [
            { "label": "Message", "value": "message" },
            { "label": "Chat", "value": "chat" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "sendText",
          "options": [
            { "label": "Send Text", "value": "sendText" },
            { "label": "Send Photo", "value": "sendPhoto" },
            { "label": "Edit Text", "value": "editText" },
            { "label": "Delete", "value": "delete" }
          ],
          "displayOptions": { "show": { "resource": ["message"] } }
        },
        {
          "name": "operationChat",
          "displayName": "Operation",
          "type": "options",
          "default": "getInfo",
          "options": [
            { "label": "Get Info", "value": "getInfo" },
            { "label": "Get Admins", "value": "getAdmins" }
          ],
          "displayOptions": { "show": { "resource": ["chat"] } }
        },

        {
          "name": "chatId",
          "displayName": "Chat ID",
          "type": "fieldSource",
          "required": true,
          "requestName": "chat_id",
          "allowedSources": ["literal", "state", "event", "nodeOutput"],
          "hint": "ID numérico do chat ou @username do canal"
        },

        {
          "name": "text",
          "displayName": "Text",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 4 },
          "displayOptions": { "show": { "operation": ["sendText"] } }
        },
        {
          "name": "photo",
          "displayName": "Photo URL",
          "type": "fieldSource",
          "required": true,
          "hint": "URL pública da imagem",
          "displayOptions": { "show": { "operation": ["sendPhoto"] } }
        },
        {
          "name": "caption",
          "displayName": "Caption",
          "type": "fieldSource",
          "typeOptions": { "multiline": true },
          "displayOptions": { "show": { "operation": ["sendPhoto"] } }
        },
        {
          "name": "messageId",
          "displayName": "Message ID",
          "type": "fieldSource",
          "required": true,
          "requestName": "message_id",
          "displayOptions": { "show": { "operation": ["editText", "delete"] } }
        },
        {
          "name": "newText",
          "displayName": "New Text",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 4 },
          "displayOptions": { "show": { "operation": ["editText"] } }
        },

        {
          "name": "parseMode",
          "displayName": "Parse Mode",
          "type": "options",
          "default": "HTML",
          "requestName": "parse_mode",
          "options": [
            { "label": "None", "value": "" },
            { "label": "HTML", "value": "HTML" },
            { "label": "MarkdownV2", "value": "MarkdownV2" }
          ],
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendText", "sendPhoto", "editText"] } }
        },
        {
          "name": "disableNotification",
          "displayName": "Silent Message",
          "type": "boolean",
          "default": false,
          "requestName": "disable_notification",
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendText", "sendPhoto"] } }
        },
        {
          "name": "protectContent",
          "displayName": "Protect Content",
          "type": "boolean",
          "default": false,
          "requestName": "protect_content",
          "hint": "Impede encaminhar e salvar",
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendText", "sendPhoto"] } }
        }
      ],

      "operations": {
        "message/sendText": {
          "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendMessage" },
          "output": { "dataPath": "result" }
        },
        "message/sendPhoto": {
          "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendPhoto" },
          "output": { "dataPath": "result" }
        },
        "message/editText": {
          "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/editMessageText" },
          "output": { "dataPath": "result" }
        },
        "message/delete": {
          "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/deleteMessage" },
          "output": { "dataPath": "result" }
        },
        "chat/getInfo": {
          "request": { "method": "GET", "path": "/bot{{credentials.botToken}}/getChat" },
          "output": { "dataPath": "result" }
        },
        "chat/getAdmins": {
          "request": { "method": "GET", "path": "/bot{{credentials.botToken}}/getChatAdministrators" },
          "output": { "dataPath": "result" }
        }
      },

      "defaults": {
        "resource": "message",
        "operation": "sendText",
        "chatId": { "type": "literal", "value": "" },
        "text": { "type": "literal", "value": "" },
        "parseMode": "HTML",
        "disableNotification": false,
        "protectContent": false
      },

      "outputHints": [
        { "path": "message_id", "description": "ID da mensagem enviada" },
        { "path": "chat.id", "description": "ID do chat" },
        { "path": "date", "description": "Timestamp Unix" }
      ]
    }
  ]
}
```

---

## 4. O Que o ManifestLoader Produz

O `ManifestLoader` lê o JSON e produz um objeto `WorkflowPlugin` idêntico ao padrão core:

```typescript
// O ManifestLoader converte manifest.json → WorkflowPlugin
const telegramPlugin: WorkflowPlugin = {
  id: 'integration-telegram',
  name: 'Telegram',
  version: '1.0.0',
  category: 'integrations',
  icon: 'send',

  onActivate(context) {
    // Registra i18n do manifest (ou defaults em inglês)
    context.registerTranslations('en-US', {
      nodes: {
        telegram: {
          label: 'Telegram',
          description: 'Send messages via Telegram Bot API',
        },
      },
    });
  },

  nodeTypes: [
    {
      type: 'integration/telegram',
      label: 'Telegram',
      icon: 'send',
      color: 'light-blue-7',
      description: 'Send messages via Telegram Bot API',

      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'success', label: 'Success', position: 'bottom', color: '#4caf50' },
        { id: 'error', label: 'Error', position: 'bottom', color: '#ef5350' },
      ],

      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode), // Reutiliza o componente genérico

      // ═══ AQUI ESTÁ A DECISÃO CHAVE ═══
      // properties[] renderiza via DynamicNodeForm (Tier 1)
      // Mas precisamos do type 'fieldSource' que NÃO EXISTE AINDA
      // Duas opções:

      // OPÇÃO A: Estender NodePropertyDefinition com type 'fieldSource'
      //          DynamicNodeForm renderiza FieldSourceSelector (que JÁ EXISTE)
      properties: [
        // ... todas as properties do manifest
      ],

      // OPÇÃO B: Criar IntegrationNodeConfig.vue genérico (configComponent)
      //          Que lê properties[] do manifest e renderiza incluindo fieldSource
      // configComponent: markRaw(IntegrationNodeConfig),

      defaults: {
        resource: 'message',
        operation: 'sendText',
        chatId: { type: 'literal', value: '' },
        text: { type: 'literal', value: '' },
        parseMode: 'HTML',
        disableNotification: false,
        protectContent: false,
      },

      outputHints: [
        { path: 'message_id', description: 'ID da mensagem enviada' },
        { path: 'chat.id', description: 'ID do chat' },
        { path: 'date', description: 'Timestamp Unix' },
      ],

      // O manifest carrega dados extras que o executor precisa
      // Armazenados no _pluginMeta (extensão injetada pelo ManifestLoader)
      _pluginMeta: {
        credentials: { /* ... */ },
        baseUrl: 'https://api.telegram.org',
        operations: { /* ... */ },
      },
    },
  ],
};
```

---

## 5. Como a UI Renderiza — Passo a Passo

### 5.1 Usuário arrasta "Telegram" do catalog para o canvas

```
1. pluginRegistry.catalog (getter)
   → Encontra plugin 'integration-telegram' na categoria 'integrations'
   → Resolve label via i18n: wf.integration-telegram.nodes.telegram.label

2. WorkflowCanvas cria o node Vue Flow
   → type: 'integration/telegram'
   → canvasComponent: GenericWorkflowNode (ícone + label + handles)

3. Node aparece no canvas com:
   → Ícone: 'send' (Material Icons)
   → Cor: light-blue-7
   → 1 input (top), 2 outputs (success + error)
```

### 5.2 Usuário clica no node → abre NodeConfigPanel

```
NodeConfigPanel.vue (3-tier fallback):

1. Busca nodeType: pluginRegistry.getNodeType('integration/telegram')
2. Verifica: nodeType.properties?.length > 0 ? → SIM
3. Renderiza: <DynamicNodeForm :properties="properties" :config="config" />
```

### 5.3 DynamicNodeForm renderiza os campos

**Estado inicial** (resource = "message", operation = "sendText"):

```
┌─────────────────────────────────────────────────┐
│  TELEGRAM                                        │
│─────────────────────────────────────────────────│
│                                                  │
│  Resource        [▼ Message              ]       │  ← type: options
│                                                  │
│  Operation       [▼ Send Text            ]       │  ← type: options (displayOptions: resource=message)
│                                                  │
│  Chat ID         [literal ▼] [___________]       │  ← type: fieldSource → FieldSourceSelector
│                  O usuário pode trocar para:      │
│                  • literal (digita o ID)          │
│                  • state (pega de uma variável)   │
│                  • event (pega do evento trigger)  │
│                  • nodeOutput (de um node antes)  │
│                                                  │
│  Text            [literal ▼]                     │  ← type: fieldSource (multiline)
│                  ┌──────────────────────────┐    │
│                  │ Hello, welcome to the    │    │
│                  │ channel!                 │    │
│                  └──────────────────────────┘    │
│                                                  │
│  Parse Mode      [▼ HTML                 ]       │  ← type: options
│                                                  │
│  Silent Message  [ ] OFF                         │  ← type: boolean
│                                                  │
│  Protect Content [ ] OFF                         │  ← type: boolean
│                                                  │
└─────────────────────────────────────────────────┘
```

**Quando o usuário troca operation para "Send Photo":**

```
displayOptions avalia:
  - text:       show: { operation: ["sendText"] }     → ESCONDE
  - photo:      show: { operation: ["sendPhoto"] }    → MOSTRA
  - caption:    show: { operation: ["sendPhoto"] }    → MOSTRA

┌─────────────────────────────────────────────────┐
│  Resource        [▼ Message              ]       │
│  Operation       [▼ Send Photo           ]       │
│                                                  │
│  Chat ID         [literal ▼] [___________]       │
│                                                  │
│  Photo URL       [literal ▼] [___________]       │  ← apareceu!
│                                                  │
│  Caption         [literal ▼]                     │  ← apareceu!
│                  ┌──────────────────────────┐    │
│                  │                          │    │
│                  └──────────────────────────┘    │
│                                                  │
│  Parse Mode      [▼ HTML                 ]       │
│  Silent Message  [ ] OFF                         │
│  Protect Content [ ] OFF                         │
└─────────────────────────────────────────────────┘
```

### 5.4 FieldSourceSelector em ação

Quando o usuário configura o Chat ID para vir de um evento:

```
Chat ID  [event ▼] [payload.telegram.chat_id]

Config salvo:
{
  "chatId": {
    "type": "event",
    "value": "payload.telegram.chat_id",
    "mode": "manual"
  }
}
```

Quando o Chat ID vem de um node anterior (ex: um Code node):

```
Chat ID  [nodeOutput ▼] [Code_1 ▼] [chatId]

Config salvo:
{
  "chatId": {
    "type": "nodeOutput",
    "value": "chatId",
    "nodeId": "node_abc123"
  }
}
```

---

## 6. Config Final Salvo no Workflow

Quando o usuário salva o workflow, o node fica assim no JSON:

```json
{
  "id": "node_tg_001",
  "type": "integration/telegram",
  "label": "Send Welcome",
  "position": { "x": 400, "y": 200 },
  "config": {
    "resource": "message",
    "operation": "sendText",
    "chatId": {
      "type": "event",
      "value": "payload.telegram.chat_id"
    },
    "text": {
      "type": "state",
      "value": "welcomeMessage"
    },
    "parseMode": "HTML",
    "disableNotification": false,
    "protectContent": false
  }
}
```

**Nota**: `resource`, `operation`, `parseMode` são valores estáticos (strings simples).
`chatId`, `text` são `FieldSourceValue` objects — o executor resolve em runtime.

---

## 7. Como o Backend Executa — Passo a Passo

### 7.1 DAG chega no node `integration/telegram`

```
RuntimeService.executeStep()
  → NodeType = "integration/telegram"
  → Despacha para IntegrationExecutor
```

### 7.2 IntegrationExecutor lê o manifest

```go
func (e *IntegrationExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    // 1. Identifica o plugin
    pluginID := "integration-telegram"                // extraído do node type
    manifest := e.manifestCache.Get(pluginID)         // lê do TieredCache (MinIO → RAM)

    // 2. Pega resource + operation do config
    config := execCtx.ParsedConfig                    // o JSON do node config
    resource := config["resource"].(string)           // "message"
    operation := config["operation"].(string)         // "sendText"

    // 3. Busca a operation definition no manifest
    opKey := resource + "/" + operation                // "message/sendText"
    opDef := manifest.Operations[opKey]
    // opDef = { request: { method: "POST", path: "/bot{{credentials.botToken}}/sendMessage" }, output: { dataPath: "result" } }

    // 4. Resolve FieldSourceValues → valores concretos
    resolvedFields := e.resolveFields(manifest.Properties, config, execCtx)
    // chatId: { type: "event", value: "payload.telegram.chat_id" }
    //   → resolve para "123456789" (pega do eventPayload)
    // text: { type: "state", value: "welcomeMessage" }
    //   → resolve para "Bem-vindo!" (pega do state)

    // 5. Constrói o body do HTTP request
    body := e.buildRequestBody(manifest.Properties, resolvedFields)
    // Aplica requestName mapping:
    //   chatId (requestName: "chat_id") → body["chat_id"] = "123456789"
    //   text (sem requestName) → body["text"] = "Bem-vindo!"
    //   parseMode (requestName: "parse_mode") → body["parse_mode"] = "HTML"
    //   disableNotification (requestName: "disable_notification") → body["disable_notification"] = false

    // body final:
    // {
    //   "chat_id": "123456789",
    //   "text": "Bem-vindo!",
    //   "parse_mode": "HTML",
    //   "disable_notification": false,
    //   "protect_content": false
    // }

    // 6. Resolve credenciais
    creds := e.credentialStore.Get(ctx, pluginID, execCtx.OrgID)
    // creds = { botToken: "123456:ABC-DEF..." }

    // 7. Resolve URL template
    url := manifest.BaseUrl + opDef.Request.Path
    // url = "https://api.telegram.org/bot123456:ABC-DEF.../sendMessage"

    // 8. Constrói pipeline steps
    steps := []PipelineStep{}

    // Telegram NÃO tem preScript → pula

    // Step HTTP (sempre presente)
    steps = append(steps, PipelineStep{
        Service: "triggers",
        Action:  "http",
        Request: HTTPRequest{
            Method:  opDef.Request.Method,   // "POST"
            URL:     url,                     // URL completa com token
            Headers: map[string]string{"Content-Type": "application/json"},
            Body:    body,
        },
    })

    // Telegram NÃO tem postScript → pula

    // 9. Suspende o DAG com pipeline
    return &NodeExecutionResult{
        OutputHandles: []string{},  // vazio = aguardando async
        NodeState: map[string]interface{}{
            "waitType": "plugin_pipeline",
            "pipeline": map[string]interface{}{
                "steps":       steps,
                "currentStep": 0,
                "results":     map[string]interface{}{},
            },
            "output": opDef.Output,  // { dataPath: "result" }
        },
    }, nil
}
```

### 7.3 Pipeline Dispatch → Triggers Service

```
Workflow Service publica no NATS:
  Subject: PLUGIN-HTTP.request
  Payload: {
    instanceId: "inst_xyz",
    nodeId: "node_tg_001",
    request: {
      method: "POST",
      url: "https://api.telegram.org/bot123456:ABC-DEF.../sendMessage",
      headers: { "Content-Type": "application/json" },
      body: {
        "chat_id": "123456789",
        "text": "Bem-vindo!",
        "parse_mode": "HTML",
        "disable_notification": false,
        "protect_content": false
      }
    }
  }
```

### 7.4 Triggers Service Executa → Callback

```
Triggers Service:
  1. Recebe a mensagem NATS
  2. Faz o HTTP POST para api.telegram.org
  3. Recebe resposta: { ok: true, result: { message_id: 42, chat: { id: 123 }, ... } }
  4. Publica callback no NATS: WORKFLOW-RESUME
     Payload: {
       instanceId: "inst_xyz",
       nodeId: "node_tg_001",
       result: { ok: true, result: { message_id: 42, chat: { id: 123 }, date: 1710345600 } }
     }
```

### 7.5 Workflow Resume → KV Checkpoint → Continue DAG

```go
func (s *RuntimeService) handlePluginPipelineResume(instance, nodeID, callbackResult) {
    nodeState := instance.NodeStates[nodeID]
    pipeline := nodeState["pipeline"]

    // Salva resultado do step HTTP
    pipeline.Results["step_0"] = callbackResult

    // Avança pipeline
    pipeline.CurrentStep = 1   // era 0, agora 1

    // KV checkpoint (crash recovery)
    s.checkpoint(instance)

    // Pipeline completo? (1 step, currentStep = 1)
    if pipeline.CurrentStep >= len(pipeline.Steps) {
        // Extrai output via dataPath
        outputConfig := nodeState["output"]   // { dataPath: "result" }
        httpResponse := callbackResult        // { ok: true, result: { message_id: 42, ... } }

        // Aplica dataPath: "result"
        nodeOutput := httpResponse["result"]  // { message_id: 42, chat: { id: 123 }, date: 1710345600 }

        // Salva como node output (disponível para nodes seguintes via nodeOutput)
        instance.NodeOutputs[nodeID] = nodeOutput

        // Continua DAG pelo handle "success"
        return nodeID, nil
    }
}
```

### 7.6 Resultado no Workflow

O node seguinte pode acessar o output:

```json
{
  "type": "nodeOutput",
  "value": "message_id",
  "nodeId": "node_tg_001"
}
```

Resolve para `42` — o ID da mensagem enviada.

---

## 8. O Que Precisa Ser Construído

### 8.1 NENHUMA mudança necessária

| Componente | Status |
|-----------|--------|
| `WorkflowPlugin` interface | ✅ JÁ EXISTE — `category: 'integrations'` já suportado |
| `PluginNodeType` interface | ✅ JÁ EXISTE — `properties`, `defaults`, `outputHints` |
| `GenericWorkflowNode.vue` | ✅ JÁ EXISTE — renderiza qualquer node no canvas |
| `NodeConfigPanel.vue` | ✅ JÁ EXISTE — 3-tier fallback |
| `FieldSourceSelector.vue` | ✅ JÁ EXISTE — literal/state/event/nodeOutput |
| `pluginRegistry` store | ✅ JÁ EXISTE — registerPlugin, catalog, getNodeType |
| `displayOptions` logic | ✅ JÁ EXISTE — DynamicNodeForm.isVisible() |

### 8.2 Extensão necessária (pequena)

| Componente | Mudança |
|-----------|---------|
| `NodePropertyType` | Adicionar `'fieldSource'` ao union type |
| `NodePropertyDefinition` | Adicionar `requestName?`, `allowedSources?`, `typeOptions?` |
| `DynamicNodeForm.vue` | Adicionar case `'fieldSource'` → renderiza `<FieldSourceSelector />` |

**Mudança no DynamicNodeForm.vue** — adicionar UM bloco:

```vue
<!-- NOVO: fieldSource type → FieldSourceSelector -->
<FieldSourceSelector
  v-else-if="prop.type === 'fieldSource'"
  :model-value="form[prop.name] || { type: 'literal', value: '' }"
  :allowed-types="prop.allowedSources || ['literal', 'state', 'event', 'nodeOutput']"
  :label="prop.displayName"
  :state-fields="stateFields"
  :node-output-options="nodeOutputOptions"
  @update:model-value="form[prop.name] = $event"
/>
```

**Isso é tudo.** O FieldSourceSelector já faz todo o trabalho pesado.

### 8.3 Componentes novos

| Componente | O que faz |
|-----------|-----------|
| `ManifestLoader` (JS) | Lê manifest.json → produz WorkflowPlugin |
| `IntegrationExecutor` (Go) | Lê manifest + config → resolve fields → constrói pipeline → suspende |
| `PluginHTTPConsumer` (Go, Triggers) | Recebe request NATS → executa HTTP → callback |
| `PipelineResumeHandler` (Go) | Recebe callback → avança pipeline → checkpoint → continua DAG |

---

## 9. Exemplo de Fluxo Completo — Workflow Real

```
[Start] → [Trigger: Webhook] → [Set State: welcomeMsg] → [Telegram: Send Welcome] → [End]
```

**Config do node Telegram neste workflow:**

```json
{
  "resource": "message",
  "operation": "sendText",
  "chatId": {
    "type": "event",
    "value": "payload.chat_id"
  },
  "text": {
    "type": "state",
    "value": "welcomeMsg"
  },
  "parseMode": "HTML",
  "disableNotification": false,
  "protectContent": false
}
```

**Runtime:**
1. Webhook trigger recebe `{ chat_id: "999", name: "João" }`
2. Set State define `welcomeMsg = "Olá João, bem-vindo!"`
3. Telegram node:
   - Resolve `chatId` → event payload → `"999"`
   - Resolve `text` → state → `"Olá João, bem-vindo!"`
   - Constrói HTTP: POST `/sendMessage` body: `{ chat_id: "999", text: "Olá João...", parse_mode: "HTML" }`
   - Despacha para Triggers Service via NATS
   - Resume → salva output → continua para End

---

## 10. Resumo das Decisões

| Decisão | Escolha | Motivo |
|---------|---------|--------|
| Rendering | `properties[]` (Tier 1) | Telegram é simples, não precisa de configComponent custom |
| Campos de dados | `type: "fieldSource"` | Permite literal/state/event/nodeOutput — MAIS PODEROSO que n8n |
| Campos de config | `type: "options"` / `type: "boolean"` | Resource, operation, parseMode são estáticos |
| Scripts | Zero | Telegram API aceita JSON direto, auto-mapping resolve tudo |
| Visibilidade | `displayOptions.show` | Já existe no DynamicNodeForm, zero mudança |
| Canvas component | `GenericWorkflowNode` | Reutiliza o componente genérico |
| Pipeline | 1 step (HTTP apenas) | Sem preScript, sem postScript |
| Outputs | `success` + `error` | Padrão para integration nodes |
| outputHints | `message_id`, `chat.id`, `date` | Para FieldSourceSelector no node seguinte |

# Plugin Webhook & Trigger — Architecture

> Como plugins recebem dados externos via webhook usando a pipeline
> DataSource → JS Executor → Router → Workflow.
>
> **DECISÃO**: O setup é GUIADO (wizard dentro do trigger node), não escondido.
> O user vê e entende cada entidade (DataSource, RouteGroup) mas o sistema
> pré-preenche defaults inteligentes. Zero magia oculta.

---

## 1. Princípio: Guiado, Não Escondido

### 1.1 Por que NÃO full auto (Saga escondida)

```
Saga falha no Step 3 → user vê "Erro ao ativar workflow"
User: "o que é DataSource? o que quebrou?"
User não sabe debugar porque nunca viu essas entidades
→ Ticket de suporte
```

Se o user não sabe que DataSource existe, ele não sabe consertar quando quebra.

### 1.2 Por que NÃO full manual

```
User quer: "quando alguém manda mensagem no Telegram, processar"
Precisa: criar DataSource, Asset, RouteGroup, linkar, registrar webhook...
User: "eu só queria receber uma mensagem..."
```

Muito passo = user desiste ou erra na configuração.

### 1.3 A abordagem: Wizard guiado dentro do trigger node

O user configura TUDO dentro do trigger node no workflow editor.
O sistema guia passo a passo, pré-preenche defaults, mostra o que vai ser criado.
O user confirma cada passo. As entidades ficam VISÍVEIS na UI.

**Regra: o sistema é inteligente nos defaults, o user confirma cada passo.**

---

## 2. Pipeline de Recebimento (Já Existe)

```
Webhook externo
      │
      ▼
HTTP Gateway (DataSource)
  ✅ Recebe POST /api/v1/events?ds={dataSourceId}
  ✅ Valida auth (5 estratégias: JWT, OAuth2, API Key, IP, None)
  ✅ Rate limiting (token bucket, por segundo/minuto/hora)
  ✅ Publica NATS: processor.js.execute
      │
      ▼
JS Executor (Processor)
  ✅ Resolve Asset via assetBind
  ✅ Roda script de transformação (payload → formato padrão)
  ✅ Enriquece com metadata do Asset
      │
      ▼
Router (RouteGroup)
  ✅ Avalia match rules (eq, neq, gt, in, starts_with, etc.)
  ✅ Roteia para destinos (trigger, ruleEngine, saveEvent, notification)
  ✅ NOVO: kind "workflow" (ver Seção 5)
      │
      ▼
Workflow Service
  ✅ Cria nova instância OU resume instância ativa
  ✅ event = payload do webhook
  ✅ Executa DAG
```

Cada camada é um serviço separado, escala independente, NATS entre cada uma.
Essa pipeline JÁ EXISTE e funciona. O que adicionamos é o `kind: "workflow"` no Router.

---

## 3. Lifecycle — Dois Ciclos Independentes

### 3.1 Credential controla DataSource (1:1)

```
Configurar credential → Cria DataSource + Asset + Registra webhook
Remover credential    → Remove DataSource + Asset + Unregister webhook

1 credential (1 bot) = 1 DataSource = 1 webhook URL
```

### 3.2 Workflow controla RouteGroup (1:1)

```
Ativar workflow   → Cria RouteGroup (linked ao Asset do DataSource)
Desativar workflow → Disable RouteGroup
Deletar workflow  → Delete RouteGroup
```

### 3.3 Diagrama de lifecycle

```
                    CREDENTIAL LIFECYCLE
                    ────────────────────
                    Cria credential
                         │
                    Cria DataSource + Asset
                    Registra webhook no serviço externo
                         │
                         ▼
              ┌─── DataSource ATIVO ───┐
              │   (recebe webhooks)     │
              │                        │
              │   WORKFLOW LIFECYCLE    │
              │   ──────────────────    │
              │   Ativa workflow A      │
              │     → Cria RouteGroup A │
              │   Ativa workflow B      │
              │     → Cria RouteGroup B │
              │   Desativa workflow A   │
              │     → Disable RG A      │
              │   Deleta workflow B     │
              │     → Delete RG B       │
              │                        │
              └────────────────────────┘
                         │
                    Remove credential
                         │
                    Unregister webhook
                    Delete DataSource + Asset
```

### 3.4 Por que é assim

O DataSource é 1:1 com a credential porque **muitos provedores permitem apenas 1 webhook
por credencial** (Telegram: 1 webhook per bot, Stripe: 1 webhook endpoint per account).

Múltiplos workflows podem escutar o MESMO webhook (mesmo bot). O que diferencia é o
**routing** — cada workflow tem seu RouteGroup com match rules específicas.

```
1 Bot Telegram (1 credential, 1 DataSource, 1 webhook URL)
  ├── RouteGroup A: match text starts_with "/pedido"  → Workflow Pedidos
  ├── RouteGroup B: match text starts_with "/ajuda"   → Workflow Suporte
  └── RouteGroup C: match ALL                         → Workflow Logging
```

---

## 4. Trigger Node — Wizard Guiado (3 Steps)

### 4.1 Step 1: Data Source

```
┌─────────────────────────────────────────────────┐
│  TELEGRAM TRIGGER — Setup                        │
│──────────────────────────────────────────────────│
│                                                   │
│  Step 1 of 3: Data Source                         │
│                                                   │
│  Select or create the data source for your bot:   │
│                                                   │
│  ○ Use existing:                                  │
│    [▼ @acme_bot (ds_001)           ]              │
│                                                   │
│  ● Create new:                                    │
│    Name: [Telegram - @acme_bot     ]              │
│    Credential: [▼ My Bot Token     ]              │
│    Auth: [▼ None                   ]              │
│      ↑ pré-preenchido pelo manifest               │
│                                                   │
│                              [Next →]             │
└─────────────────────────────────────────────────┘
```

**"Use existing"** aparece quando já existe um DataSource para a mesma credential.
Isso resolve naturalmente o caso de múltiplos workflows com o mesmo bot.

O manifest do plugin pré-preenche:
- Auth type (Telegram = none, pois valida internamente)
- Nome sugerido

### 4.2 Step 2: Routing

```
┌─────────────────────────────────────────────────┐
│  Step 2 of 3: Routing                            │
│                                                   │
│  Events to listen:                                │
│    ☑ message                                      │
│    ☐ edited_message                               │
│    ☐ callback_query                               │
│    ☐ channel_post                                 │
│    ☐ inline_query                                 │
│    ↑ opções vindas do manifest (trigger.events)   │
│                                                   │
│  ─── Instance Mode ───                            │
│                                                   │
│  How should incoming events be handled?           │
│                                                   │
│  ○ New Instance                                   │
│    Every event starts a new workflow execution.   │
│    Workflow runs, completes, and terminates.       │
│                                                   │
│  ● Running Instance (conversation/session)        │
│    Events are delivered to an active instance      │
│    identified by a correlation field.              │
│                                                   │
│  ─── Running Instance Config ───                  │
│                                                   │
│  Identify session by:                             │
│    [message.chat.id         ]                     │
│    ↑ pré-preenchido pelo manifest                 │
│    (correlationField)                             │
│                                                   │
│  Signal name:                                     │
│    [telegram.message        ]                     │
│    ↑ auto-gerado: {pluginId}.{eventType}          │
│                                                   │
│  If no active session:                            │
│    [▼ Start new instance    ]                     │
│    Opções: Start new | Drop event | Queue         │
│                                                   │
│  ─── Filters (optional) ───                       │
│                                                   │
│  [+ Add Filter]                                   │
│  message.text starts_with "/pedido"               │
│                                                   │
│                      [← Back]  [Next →]           │
└─────────────────────────────────────────────────┘
```

### 4.3 Step 3: Review & Activate

```
┌─────────────────────────────────────────────────┐
│  Step 3 of 3: Review                             │
│                                                   │
│  Summary of what will be configured:              │
│                                                   │
│  Data Source                                      │
│  ☑ Using existing: "Telegram - @acme_bot"         │
│    Auth: None | Rate limit: default               │
│    URL: https://gw.mapexos.io/events?ds=ds_001   │
│                                                   │
│  Route                                            │
│  ☑ Will create route for this workflow:            │
│    Events: message                                │
│    Mode: Running instance                         │
│    Correlation: message.chat.id                   │
│    Filter: message.text starts_with "/pedido"     │
│    On miss: Start new instance                    │
│                                                   │
│  Webhook Registration                             │
│  ☑ Telegram webhook is already registered         │
│    (DataSource ds_001 already active)             │
│    ─── OR ───                                     │
│  ☑ Will register webhook at Telegram:             │
│    POST /bot{token}/setWebhook                    │
│    URL: https://gw.mapexos.io/events?ds=ds_001   │
│                                                   │
│  ⚠ Note: Telegram allows only 1 webhook per bot. │
│    If another service uses this bot's webhook,    │
│    it will be overwritten.                        │
│                                                   │
│                      [← Back]  [Activate →]       │
└─────────────────────────────────────────────────┘
```

---

## 5. Router — kind: "workflow" (NOVO)

### 5.1 RouteGroup config

O wizard cria um RouteGroup com `kind: "workflow"`:

```json
{
  "id": "rg_001",
  "name": "Pedidos → wf_pedidos",
  "enabled": true,
  "managedBy": {
    "workflowId": "wf_pedidos",
    "nodeId": "node_trigger_001",
    "pluginId": "telegram"
  },
  "routers": [
    {
      "kind": "workflow",
      "match": {
        "policy": "all",
        "rules": [
          { "field": "message.text", "op": "starts_with", "value": "/pedido" }
        ]
      },
      "workflow": {
        "workflowId": "wf_pedidos",
        "mode": "newInstance",
        "correlationField": null,
        "signalName": null,
        "onMiss": null
      }
    }
  ]
}
```

### 5.2 Dois modos

#### Mode: newInstance

```json
{
  "kind": "workflow",
  "workflow": {
    "workflowId": "wf_alertas",
    "mode": "newInstance"
  }
}
```

Sempre cria instância nova. Workflow executa, completa, termina.
Não precisa de correlationField nem signalName.

#### Mode: runningInstance

```json
{
  "kind": "workflow",
  "workflow": {
    "workflowId": "wf_suporte",
    "mode": "runningInstance",
    "correlationField": "message.chat.id",
    "signalName": "telegram.message",
    "onMiss": "create"
  }
}
```

Busca instância ativa pelo correlationField. Se encontra, entrega como signal.
Se não encontra, aplica onMiss policy.

### 5.3 onMiss policies

| Policy | Comportamento | Uso |
|--------|---------------|-----|
| `create` | Cria instância nova (fallback) | Bot conversacional — 1a msg cria, próximas resumem |
| `drop` | Descarta o evento | Só aceita se tem instância ativa |
| `queue` | Guarda no KV, entrega quando instância existir | Eventos que chegam antes do workflow estar pronto |

### 5.4 NATS subjects

O Router publica para o Workflow Service:

| Mode | Subject NATS | Payload |
|------|-------------|---------|
| `newInstance` | `workflow.instance.create` | `{ workflowId, event: payload }` |
| `runningInstance` | `workflow.signal.deliver` | `{ workflowId, correlationField, correlationValue, signalName, data: payload }` |

O Router extrai o `correlationValue` do payload usando o `correlationField` configurado.

---

## 6. Workflow Service — Signal Delivery (Temporal-inspired)

### 6.1 Como o Workflow Service recebe

```go
// Consumer de workflow.signal.deliver
func (s *RuntimeService) handleSignalDeliver(msg SignalDeliverMessage) {
    // 1. Busca instância ativa com matching correlation
    instance := s.findActiveInstance(
        msg.WorkflowId,
        msg.CorrelationField,
        msg.CorrelationValue,
    )

    if instance != nil {
        // 2a. Instância encontrada → entrega signal
        s.deliverSignal(instance, msg.SignalName, msg.Data)
        return
    }

    // 2b. Instância não encontrada → aplica onMiss
    switch msg.OnMiss {
    case "create":
        s.createInstance(msg.WorkflowId, msg.Data)
    case "drop":
        // log + discard
    case "queue":
        s.queueSignal(msg)
    }
}
```

### 6.2 Como a instância espera um signal

No DAG, um node do tipo "Wait for Signal" suspende a instância:

```go
// Wait node executor
func (e *WaitSignalExecutor) Execute(ctx *NodeExecutionContext) (*NodeExecutionResult, error) {
    return &NodeExecutionResult{
        Status:   StatusWaiting,
        WaitType: "signal",
        WaitConfig: map[string]interface{}{
            "signalName":       ctx.Config["signalName"],
            "correlationField": ctx.Config["correlationField"],
            "correlationValue": ctx.ResolveFieldSource(ctx.Config["correlationValue"]),
            "timeout":          ctx.Config["timeout"],  // opcional
        },
    }, nil
}
```

O KV fica:
```
instance:inst_001 = {
    status: "waiting",
    currentNode: "node_wait_001",
    waitType: "signal",
    waitConfig: {
        signalName: "telegram.message",
        correlationField: "message.chat.id",
        correlationValue: "987654321"
    }
}
```

### 6.3 Busca de instância por correlation

```go
func (s *RuntimeService) findActiveInstance(
    workflowId string,
    correlationField string,
    correlationValue string,
) *WorkflowInstance {
    // Busca no KV: instâncias do workflowId com status "waiting"
    // e waitConfig.correlationValue == correlationValue
    // e waitConfig.signalName matches

    // Índice no KV: "correlation:{workflowId}:{correlationValue}" → instanceId
    key := fmt.Sprintf("correlation:%s:%s", workflowId, correlationValue)
    instanceId, err := s.kv.Get(key)
    if err != nil {
        return nil
    }
    return s.loadInstance(instanceId)
}
```

**Índice de correlação** no NATS KV:
```
correlation:wf_suporte:987654321 → inst_001
correlation:wf_suporte:123456789 → inst_002
correlation:wf_pedidos:987654321 → inst_003
```

Quando a instância entra em "waiting" com signal, cria a entrada.
Quando a instância resume ou termina, remove a entrada.

### 6.4 Entrega do signal

```go
func (s *RuntimeService) deliverSignal(
    instance *WorkflowInstance,
    signalName string,
    data map[string]interface{},
) {
    nodeId := instance.CurrentNode

    // Salva signal data no node state
    instance.NodeStates[nodeId]["signalData"] = data
    instance.Status = StatusRunning

    // Remove índice de correlação
    corrKey := buildCorrelationKey(instance)
    s.kv.Delete(corrKey)

    // Checkpoint
    s.checkpoint(instance)

    // Continua DAG a partir do node que estava esperando
    s.advanceDAG(instance, nodeId)
}
```

---

## 7. Cenário End-to-End: Bot Conversacional Telegram

### 7.1 Setup (uma vez)

```
1. User instala plugin Telegram
2. User configura credential (bot token)
   → Sistema cria DataSource + Asset
   → Sistema registra webhook no Telegram
3. User cria Workflow "Suporte Bot"
4. User arrasta Telegram Trigger → wizard:
   - Step 1: Use existing DataSource "@acme_bot"
   - Step 2: Events=message, Mode=runningInstance,
             Correlation=message.chat.id, OnMiss=create
   - Step 3: Review → Activate
   → Sistema cria RouteGroup (linked ao Asset)
```

### 7.2 Execução (runtime)

```
═══ Chat 987654321 manda "Oi" ═══

Telegram POST → HTTP Gateway (ds_001)
  → NATS → JS Executor → Router
  → RouteGroup match: kind=workflow, mode=runningInstance
  → Router extrai: correlationValue = "987654321"
  → NATS: workflow.signal.deliver {
      workflowId: "wf_suporte",
      correlationValue: "987654321",
      signalName: "telegram.message",
      data: { message: { chat: { id: 987654321 }, text: "Oi" } },
      onMiss: "create"
    }

Workflow Service:
  → Busca KV: correlation:wf_suporte:987654321 → não existe
  → onMiss = "create" → cria inst_001
  → event = { message: { chat: { id: 987654321 }, text: "Oi" } }
  → Inicia DAG:

    [Trigger] → [AI: classifica] → [Telegram: sendText "Olá! Como posso ajudar?"]
                                                    │
                                          [Wait Signal: telegram.message]
                                           correlationValue: 987654321
                                           timeout: 30min
                                                    │
                                           inst_001 SUSPENDE
                                           KV: correlation:wf_suporte:987654321 → inst_001


═══ Chat 987654321 manda "Quero saber sobre preços" ═══

Telegram POST → HTTP Gateway → JS → Router
  → NATS: workflow.signal.deliver {
      correlationValue: "987654321",
      data: { message: { text: "Quero saber sobre preços" } }
    }

Workflow Service:
  → Busca KV: correlation:wf_suporte:987654321 → inst_001
  → ENCONTROU → deliverSignal(inst_001)
  → inst_001 RESUME:

    [Wait Signal] → signalData = { message: { text: "Quero saber sobre preços" } }
                  → [AI: responde preços] → [Telegram: sendText "Nossos planos são..."]
                                                      │
                                            [Wait Signal: telegram.message]
                                             correlationValue: 987654321
                                             inst_001 SUSPENDE de novo


═══ Chat 987654321 manda "Obrigado, tchau" ═══

Telegram POST → ... → Router → signal.deliver
  → inst_001 RESUME
  → [AI: detecta encerramento] → [Telegram: sendText "Até logo!"]
  → [End]
  → inst_001 TERMINA
  → KV: DELETE correlation:wf_suporte:987654321


═══ Chat 987654321 manda "Oi" de novo ═══

Router → signal.deliver { correlationValue: "987654321" }
Workflow Service:
  → Busca KV: não existe mais
  → onMiss = "create" → cria inst_002
  → Ciclo recomeça
```

### 7.3 Múltiplos chats simultâneos

```
Chat 987654321 ("Oi")      → inst_001 (active, waiting)
Chat 111222333 ("Olá")     → inst_002 (active, waiting)
Chat 444555666 ("/pedido") → não match no RouteGroup do Suporte
                              match no RouteGroup do Pedidos → inst_003

KV state:
  correlation:wf_suporte:987654321 → inst_001
  correlation:wf_suporte:111222333 → inst_002
  correlation:wf_pedidos:444555666 → inst_003
```

Cada chat tem sua instância. Zero conflito. O correlationField garante isolamento.

---

## 8. Manifest — Seção Trigger

### 8.1 Formato no manifest.json

O manifest do plugin declara as capacidades de trigger na seção `triggerTypes`:

```json
{
  "$schema": "mapex-plugin/v1",
  "id": "telegram",
  "name": "Telegram",
  "version": "1.0.0",
  "category": "messaging",
  "dispatch": "http",
  "baseUrl": "https://api.telegram.org",

  "credentials": {
    "id": "telegramApi",
    "name": "Telegram Bot API",
    "fields": [
      {
        "name": "botToken",
        "displayName": "Bot Token",
        "type": "string",
        "required": true,
        "isSecret": true
      }
    ],
    "test": {
      "method": "GET",
      "path": "/bot{{botToken}}/getMe"
    }
  },

  "triggerTypes": [
    {
      "type": "telegram/webhook",
      "label": "Telegram Trigger",
      "icon": "webhook",
      "color": "#0088CC",
      "description": "Receive messages and events from Telegram Bot",

      "outputs": [
        { "id": "event", "label": "Event", "position": "bottom", "color": "#4caf50" }
      ],

      "webhook": {
        "register": {
          "method": "POST",
          "path": "/bot{{credentials.botToken}}/setWebhook",
          "body": {
            "url": "{{system.webhookUrl}}",
            "allowed_updates": "{{config.events}}"
          }
        },
        "unregister": {
          "method": "POST",
          "path": "/bot{{credentials.botToken}}/deleteWebhook"
        }
      },

      "defaults": {
        "dataSourceAuth": "none",
        "correlationField": "message.chat.id",
        "signalName": "telegram.message"
      },

      "events": [
        { "value": "message", "label": "Message" },
        { "value": "edited_message", "label": "Edited Message" },
        { "value": "callback_query", "label": "Callback Query" },
        { "value": "channel_post", "label": "Channel Post" },
        { "value": "inline_query", "label": "Inline Query" }
      ],

      "eventOutputs": {
        "message": {
          "dataPath": "message",
          "outputHints": [
            { "path": "message_id", "description": "Message ID" },
            { "path": "chat.id", "description": "Chat ID" },
            { "path": "chat.type", "description": "Chat type (private, group, supergroup)" },
            { "path": "from.id", "description": "Sender user ID" },
            { "path": "from.first_name", "description": "Sender first name" },
            { "path": "text", "description": "Message text" },
            { "path": "date", "description": "Unix timestamp" }
          ]
        },
        "callback_query": {
          "dataPath": "callback_query",
          "outputHints": [
            { "path": "id", "description": "Callback ID" },
            { "path": "data", "description": "Button data" },
            { "path": "message.chat.id", "description": "Chat ID" },
            { "path": "from.id", "description": "User who clicked" }
          ]
        }
      }
    }
  ],

  "nodeTypes": [
    { "type": "telegram/message", "...": "... (ações — sendText, sendPhoto, etc.)" },
    { "type": "telegram/chat", "...": "... (getInfo, getAdmins)" }
  ]
}
```

### 8.2 Campos do triggerType

| Campo | Tipo | Propósito |
|-------|------|-----------|
| `type` | string | ID único do trigger (ex: `telegram/webhook`) |
| `label` | string | Nome exibido no catalog |
| `icon` | string | Material Icon |
| `color` | string | Cor no canvas |
| `description` | string | Descrição para o tooltip |
| `outputs` | HandleDefinition[] | Handles de saída (geralmente 1: "event") |
| `webhook.register` | object | Como registrar webhook no serviço externo |
| `webhook.unregister` | object | Como remover webhook |
| `defaults.dataSourceAuth` | string | Auth type sugerido para o DataSource |
| `defaults.correlationField` | string | Campo sugerido para correlação |
| `defaults.signalName` | string | Signal name sugerido |
| `events` | OptionItem[] | Tipos de evento que o webhook pode receber |
| `eventOutputs` | Record<string, EventOutput> | outputHints por tipo de evento |

### 8.3 Diferença: triggerTypes vs nodeTypes

```
triggerTypes = RECEBER dados (inbound — webhook, polling)
  → Cria DataSource, RouteGroup
  → É o PRIMEIRO node do workflow (entrada)
  → Declarado separadamente no manifest

nodeTypes = ENVIAR dados (outbound — HTTP calls, API calls)
  → Usa pipeline de execução (dispatch via NATS)
  → São nodes INTERMEDIÁRIOS do workflow
  → Declarado separadamente no manifest
```

---

## 9. Entidades Visíveis ao User

### 9.1 managedBy (em vez de createdBySystem)

Entidades criadas via wizard NÃO são escondidas. São visíveis com tag:

```go
type ManagedByRef struct {
    WorkflowID string `bson:"workflowId" json:"workflowId"`
    NodeID     string `bson:"nodeId"     json:"nodeId"`
    PluginID   string `bson:"pluginId"   json:"pluginId"`
}
```

### 9.2 Comportamento na UI

```
Data Sources (sidebar):
  ├── Telegram - @acme_bot        [managed by: Suporte Bot, Pedidos]
  ├── Stripe - Production         [managed by: Pagamentos]
  └── Custom Webhook (manual)

Routes (sidebar):
  ├── Suporte → wf_suporte        [managed by: Suporte Bot]
  ├── Pedidos → wf_pedidos         [managed by: Pedidos]
  └── IoT alerts → rule_engine     (manual)
```

O user pode:
- **Ver**: entende o que existe e quem criou
- **Editar**: ajustar rate limit, auth, match rules se precisar
- **Deletar**: warning "This is managed by Workflow X. Deleting will break the trigger."

### 9.3 Regras de proteção

| Ação | managedBy presente | Comportamento |
|------|-------------------|---------------|
| Editar | Sim | Permitido com warning: "Changes may affect Workflow X" |
| Deletar | Sim | Confirmação extra: "This will break the trigger of Workflow X" |
| Desativar | Sim | Permitido com warning |
| Ver | Sim | Sempre visível, tag "managed by" |

---

## 10. Flow Completo — Passo a Passo

### 10.1 Primeiro workflow com Telegram Trigger

```
PASSO 1: Instalar plugin
  User: Marketplace → Install "Telegram"
  Sistema: salva manifest no MongoDB (installed_plugins)

PASSO 2: Configurar credential
  User: Plugin Settings → Telegram → Add Credential
  User: digita Bot Token → [Test Connection]
  Sistema:
    a. Testa token: GET /bot{token}/getMe (via Triggers Service, async NATS)
    b. Encrypta token (BYOK envelope encryption)
    c. Salva credential (MongoDB)
    d. Cria DataSource:
       - name: "Telegram - @{bot_username}"
       - auth: none (do manifest defaults)
       - mode: push
       - managedBy: { pluginId: "telegram" }
    e. Cria Asset (linked ao DataSource)
    f. Registra webhook no Telegram:
       POST /bot{token}/setWebhook { url: "https://gw.mapexos.io/events?ds={dsId}" }
       (via Triggers Service, async NATS)

  Resultado: DataSource ativo, recebendo webhooks do Telegram.
  Sem RouteGroups → eventos chegam mas não são processados.

PASSO 3: Criar workflow
  User: New Workflow → "Suporte Bot"
  User: arrasta "Telegram Trigger" para o canvas
  Wizard abre:
    Step 1: "Use existing: Telegram - @acme_bot" → Next
    Step 2: Events=message, Mode=runningInstance,
            Correlation=message.chat.id, OnMiss=create → Next
    Step 3: Review → Activate
  Sistema:
    a. Cria RouteGroup:
       - kind: "workflow"
       - match: null (aceita tudo)
       - workflow: { wfId, mode, correlationField, signalName, onMiss }
       - managedBy: { workflowId, nodeId, pluginId }
    b. Vincula RouteGroup ao Asset do DataSource
    c. Marca workflow como active

  Resultado: webhooks do Telegram agora são roteados para o workflow.
```

### 10.2 Segundo workflow com o MESMO bot

```
PASSO 1: Criar Workflow "Pedidos"
  User: arrasta "Telegram Trigger"
  Wizard:
    Step 1: "Use existing: Telegram - @acme_bot" ← REUSA!
    Step 2: Events=message, Mode=newInstance,
            Filter: message.text starts_with "/pedido"
    Step 3: Review → Activate
  Sistema:
    a. Cria RouteGroup (match: text starts_with "/pedido")
    b. Vincula ao MESMO Asset
    c. Done

  Resultado: mesmo DataSource, 2 RouteGroups.
  /pedido → Workflow Pedidos
  tudo → Workflow Suporte
```

### 10.3 Remover credential

```
User: Plugin Settings → Telegram → Remove "@acme_bot"

Sistema verifica: algum RouteGroup ativo usa este DataSource?
  SIM → 409 "Credential in use by: Suporte Bot, Pedidos"
         "Deactivate these workflows first."
  NÃO →
    a. POST /bot{token}/deleteWebhook (via Triggers Service)
    b. DELETE DataSource
    c. DELETE Asset
    d. DELETE credential
```

---

## 11. Separação de Responsabilidades

```
HTTP Gateway (DataSource):
  ✅ Recebe webhook HTTP
  ✅ Valida auth
  ✅ Rate limiting
  ✅ Publica no NATS
  ❌ Não sabe sobre workflows, instâncias, routing

JS Executor (Processor):
  ✅ Transforma payload
  ✅ Resolve Asset
  ✅ Enriquece dados
  ❌ Não sabe sobre workflows, instâncias, routing

Router:
  ✅ Avalia match rules
  ✅ Roteia para destinos (workflow, trigger, rule, etc.)
  ✅ Extrai correlationValue do payload
  ✅ Publica no NATS subject correto
  ❌ Não sabe sobre instâncias ativas, estado do KV

Workflow Service:
  ✅ Recebe do NATS (create ou signal.deliver)
  ✅ Resolve instância pelo correlationValue (KV)
  ✅ Decide: nova instância vs resume vs drop
  ✅ Executa DAG
  ❌ Não sabe sobre DataSource, auth, routing rules
```

Cada serviço faz UMA coisa. Zero acoplamento. NATS entre todos.

---

## 12. Escala — Enterprise Architecture

```
DataSource → JS Executor → Router → Workflow
(N pods)     (N pods)      (N pods)  (N pods)

                    NATS JetStream
              (buffer + replay + backpressure)
```

| Cenário | O que acontece |
|---------|----------------|
| 50K webhooks/segundo (pico) | DataSource aceita, NATS absorve buffer, JS processa no ritmo |
| Script de transformação com bug | 1 pod JS trava, Kubernetes mata, outros continuam |
| Routing rules complexas | Escala Router separado, não afeta DataSource nem Workflow |
| 1000 instâncias ativas com signal | Workflow Service resolve via KV index, O(1) lookup |
| Bot com 10 workflows escutando | 1 DataSource, 10 RouteGroups, Router avalia em paralelo |

---

## 13. Referência — O Que Muda vs O Que Já Existe

| Entidade / Feature | Status | Onde |
|--------------------|--------|------|
| DataSource entity | **Já existe** | HTTP Gateway service |
| Asset entity | **Já existe** | Assets service |
| RouteGroup entity | **Já existe** | Router service |
| RouteGroup kind: "workflow" | **NOVO** | Router service |
| workflow.instance.create NATS | **NOVO** | Workflow Service consumer |
| workflow.signal.deliver NATS | **NOVO** | Workflow Service consumer |
| Correlation index no KV | **NOVO** | Workflow Service runtime |
| Wait Signal executor | **NOVO** | Workflow Service executors |
| Trigger node wizard (UI) | **NOVO** | Workflow editor frontend |
| triggerTypes no manifest | **NOVO** | Plugin manifest schema |
| managedBy field | **NOVO** (substitui createdBySystem) | DataSource, Asset, RouteGroup entities |

# Integration Plugin Manifest — V1 Proposal

> Alinhado com o DSL existente: WorkflowPlugin + PluginNodeType + NodePropertyDefinition

---

## 1. Como o Manifest se Encaixa no Sistema Existente

```
┌──────────────────────────────────────────────────────────────┐
│  manifest.json (MinIO)                                       │
│  ↓                                                           │
│  ManifestLoader (runtime)                                    │
│  ↓                                                           │
│  WorkflowPlugin object                                       │
│  ↓                                                           │
│  pluginRegistry.registerPlugin(plugin)                       │
│  ↓                                                           │
│  NodeConfigPanel → DynamicNodeForm (renderiza properties[])  │
│                  → FieldSourceSelector (para type=fieldSource)│
└──────────────────────────────────────────────────────────────┘
```

O manifest JSON é um `WorkflowPlugin` serializado + campos de integração (credentials, baseUrl, operations, loadOptions).

---

## 2. Extensão do NodePropertyDefinition

### Types Existentes (core)
`'string' | 'number' | 'boolean' | 'options' | 'json'`

### Types Novos (integrations)
`'fieldSource' | 'multiOptions' | 'group' | 'object' | 'array'`

| Type | Renderiza como | Quando usar |
|------|---------------|-------------|
| `fieldSource` | FieldSourceSelector | Campos de dados dinâmicos (chatId, text, email) |
| `multiOptions` | q-select multiple | Seleção múltipla (labels, scopes) |
| `group` | Seção expansível "Add Field" | Campos opcionais agrupados |
| `object` | Sub-form nested | Estrutura fixa (address, footer) |
| `array` | Lista repetível "Add Item" | Lista de objetos (metadata, mappings) |

### Campos Novos na Property

```typescript
interface IntegrationProperty extends NodePropertyDefinition {
  // Request mapping (como o campo vai pro HTTP request)
  requestName?: string;       // Nome no request (default: name)
  in?: 'body' | 'query' | 'path';  // Onde vai (default: auto by method)

  // Para composite types (group, object, array)
  fields?: IntegrationProperty[];

  // Para fieldSource type
  allowedSources?: SourceType[];    // Default: ['literal','state','event','node_output']

  // Extended typeOptions
  typeOptions?: {
    multiline?: boolean;
    rows?: number;
    secret?: boolean;
    min?: number;
    max?: number;
    precision?: number;
    loadOptions?: string;           // Ref para loadOptions key
    dependsOn?: string[];           // Cascading: recarrega quando deps mudam
    itemLabel?: string;             // Label do item em array
    maxItems?: number;              // Limite de items em array
  };
}
```

### Regra de Uso

| Campo | Type | Motivo |
|-------|------|--------|
| resource, operation | `options` | Sempre estático, define qual operação |
| parseMode, model | `options` | Seletor de configuração |
| chatId, email, text | `fieldSource` | Dados dinâmicos — pode vir de state/event/literal/node_output |
| amount, limit | `fieldSource` | Dados numéricos dinâmicos |
| blocks, filter | `json` | JSON livre |
| metadata, mappings | `array` | Lista de key-value pairs |
| address, shipping | `object` | Objeto nested com campos fixos |
| additionalFields | `group` | Campos opcionais expandíveis |

---

## 3. Formato do Manifest

```jsonc
{
  // ══════ WorkflowPlugin (compatível com o registro existente) ══════
  "id": "integration-telegram",
  "name": "Telegram",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "send",

  // ══════ Extensões de Integração ══════
  "metadata": {
    "brandIcon": "telegram.svg",      // SVG armazenado em MinIO
    "color": "#0088CC",
    "docs": "https://core.telegram.org/bots/api"
  },

  "credentials": {
    "id": "telegramApi",
    "name": "Telegram Bot API",
    "fields": [/* credential fields */],
    "inject": {/* como injetar no request */},
    "test": {/* endpoint de teste */}
  },

  "baseUrl": "https://api.telegram.org",

  "requestDefaults": {
    "headers": { "Content-Type": "application/json" }
  },

  // ══════ LoadOptions — Dynamic dropdown loaders (proxy via backend) ══════
  // Map of resource key → loader definition.
  // Properties reference by key: typeOptions.loadOptions = "getChats"
  // Backend proxies: POST /api/v1/credentials/:id/load_options/:resourceKey
  "loadOptions": {
    "resourceKey": {
      "request": { "method": "GET", "path": "/api/endpoint/{{credentials.field}}" },
      // Simple mode (80%): extract directly
      "dataPath": "result",       // JSON path to array in response
      "valuePath": "id",          // JSON path for value within each item
      "labelPath": "name",        // JSON path for label within each item
      // Transform mode (20%): inline JS ES5 script (optional, replaces simple paths)
      "transform": "function transform(data) { return data.items.map(function(i) { return { label: i.title, value: i.id }; }); }"
    }
  },

  // ══════ Node Types (compatível com PluginNodeType) ══════
  "nodeTypes": [
    {
      "type": "integration/telegram",       // Node type no workflow
      "label": "Telegram",
      "description": "Send messages via Telegram Bot API",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [/* Extended NodePropertyDefinition[] */],

      // Pipeline por resource/operation
      "operations": {
        "message/sendText": {
          "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendMessage" },
          "output": { "dataPath": "result" }
        }
      },

      "defaults": { "resource": "message", "operation": "sendText" },

      "outputHints": [
        { "path": "*", "description": "API response data" }
      ]
    }
  ]
}
```

---

## 4. Os 9 Plugins

---

### ═══════════════════════════════════════
### SIMPLE #1 — Telegram
### ═══════════════════════════════════════

**Patterns**: API key in URL, fieldSource para dados, group para opcionais, requestName

```json
{
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
        "typeOptions": { "secret": true },
        "hint": "Get from @BotFather"
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
      "description": "Send messages via Telegram Bot API",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

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
            { "label": "Send Photo", "value": "sendPhoto" }
          ],
          "displayOptions": { "show": { "resource": ["message"] } }
        },

        {
          "name": "chatId",
          "displayName": "Chat ID",
          "type": "fieldSource",
          "required": true,
          "requestName": "chat_id",
          "allowedSources": ["literal", "state", "event", "node_output"],
          "displayOptions": { "show": { "resource": ["message"] } }
        },
        {
          "name": "text",
          "displayName": "Text",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 4 },
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendText"] } }
        },
        {
          "name": "photo",
          "displayName": "Photo URL",
          "type": "fieldSource",
          "required": true,
          "hint": "URL of the photo to send",
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendPhoto"] } }
        },
        {
          "name": "caption",
          "displayName": "Caption",
          "type": "fieldSource",
          "displayOptions": { "show": { "resource": ["message"], "operation": ["sendPhoto"] } }
        },

        {
          "name": "options",
          "displayName": "Options",
          "type": "group",
          "displayOptions": { "show": { "resource": ["message"] } },
          "fields": [
            {
              "name": "parseMode",
              "displayName": "Parse Mode",
              "type": "options",
              "requestName": "parse_mode",
              "default": "HTML",
              "options": [
                { "label": "HTML", "value": "HTML" },
                { "label": "MarkdownV2", "value": "MarkdownV2" }
              ]
            },
            {
              "name": "disableNotification",
              "displayName": "Silent",
              "type": "boolean",
              "requestName": "disable_notification",
              "default": false
            },
            {
              "name": "protectContent",
              "displayName": "Protect Content",
              "type": "boolean",
              "requestName": "protect_content",
              "default": false
            }
          ]
        },

        {
          "name": "chatIdGet",
          "displayName": "Chat ID",
          "type": "fieldSource",
          "required": true,
          "requestName": "chat_id",
          "in": "path",
          "displayOptions": { "show": { "resource": ["chat"], "operation": ["get"] } }
        },
        {
          "name": "operationChat",
          "displayName": "Operation",
          "type": "options",
          "default": "get",
          "options": [
            { "label": "Get Info", "value": "get" }
          ],
          "displayOptions": { "show": { "resource": ["chat"] } }
        }
      ],

      "operations": {
        "message/sendText": {
          "request": {
            "method": "POST",
            "path": "/bot{{credentials.botToken}}/sendMessage"
          },
          "output": { "dataPath": "result" }
        },
        "message/sendPhoto": {
          "request": {
            "method": "POST",
            "path": "/bot{{credentials.botToken}}/sendPhoto"
          },
          "output": { "dataPath": "result" }
        },
        "chat/get": {
          "request": {
            "method": "GET",
            "path": "/bot{{credentials.botToken}}/getChat"
          },
          "output": { "dataPath": "result" }
        }
      },

      "defaults": {
        "resource": "message",
        "operation": "sendText",
        "chatId": { "type": "literal", "value": "" },
        "text": { "type": "literal", "value": "" }
      },

      "outputHints": [
        { "path": "*", "description": "Telegram API response" }
      ]
    }
  ]
}
```

**Scripts**: Zero. Auto-mapping resolve tudo.
**Execução**: `config.chatId = { type: "state", value: "targetChat" }` → executor resolve para "123456789" → body: `{ "chat_id": "123456789", "text": "Hello" }`

---

### ═══════════════════════════════════════
### SIMPLE #2 — Discord Webhook
### ═══════════════════════════════════════

**Patterns**: URL dinâmico (webhook), nested object (embeds), array dentro de group

```json
{
  "id": "integration-discord",
  "name": "Discord",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "forum",

  "metadata": {
    "brandIcon": "discord.svg",
    "color": "#5865F2",
    "docs": "https://discord.com/developers/docs"
  },

  "credentials": {
    "id": "discordBotApi",
    "name": "Discord Bot",
    "fields": [
      {
        "name": "botToken",
        "displayName": "Bot Token",
        "type": "string",
        "required": true,
        "typeOptions": { "secret": true }
      }
    ],
    "inject": {
      "headers": { "Authorization": "Bot {{botToken}}" }
    },
    "test": {
      "method": "GET",
      "path": "/api/v10/users/@me"
    }
  },

  "baseUrl": "https://discord.com",

  "nodeTypes": [
    {
      "type": "integration/discord",
      "label": "Discord",
      "description": "Send messages and manage Discord",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "message",
          "options": [
            { "label": "Message", "value": "message" },
            { "label": "Webhook", "value": "webhook" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "send",
          "options": [{ "label": "Send", "value": "send" }],
          "displayOptions": { "show": { "resource": ["message", "webhook"] } }
        },

        {
          "name": "channelId",
          "displayName": "Channel ID",
          "type": "fieldSource",
          "required": true,
          "in": "path",
          "displayOptions": { "show": { "resource": ["message"] } }
        },
        {
          "name": "content",
          "displayName": "Message",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 4 },
          "displayOptions": { "show": { "resource": ["message", "webhook"] } }
        },

        {
          "name": "embeds",
          "displayName": "Embeds",
          "type": "array",
          "typeOptions": { "itemLabel": "Embed", "maxItems": 10 },
          "displayOptions": { "show": { "resource": ["message", "webhook"] } },
          "fields": [
            { "name": "title", "displayName": "Title", "type": "fieldSource" },
            { "name": "description", "displayName": "Description", "type": "fieldSource", "typeOptions": { "multiline": true } },
            { "name": "url", "displayName": "URL", "type": "fieldSource" },
            { "name": "color", "displayName": "Color (decimal)", "type": "number" },
            {
              "name": "footer",
              "displayName": "Footer",
              "type": "object",
              "fields": [
                { "name": "text", "displayName": "Text", "type": "fieldSource" },
                { "name": "icon_url", "displayName": "Icon URL", "type": "fieldSource" }
              ]
            }
          ]
        },

        {
          "name": "webhookUrl",
          "displayName": "Webhook URL",
          "type": "fieldSource",
          "required": true,
          "hint": "Full Discord webhook URL",
          "displayOptions": { "show": { "resource": ["webhook"] } }
        },
        {
          "name": "username",
          "displayName": "Username Override",
          "type": "fieldSource",
          "displayOptions": { "show": { "resource": ["webhook"] } }
        },
        {
          "name": "avatar_url",
          "displayName": "Avatar URL",
          "type": "fieldSource",
          "displayOptions": { "show": { "resource": ["webhook"] } }
        }
      ],

      "operations": {
        "message/send": {
          "request": {
            "method": "POST",
            "path": "/api/v10/channels/{{fields.channelId}}/messages"
          }
        },
        "webhook/send": {
          "request": {
            "method": "POST",
            "path": "{{fields.webhookUrl}}"
          },
          "skipCredentials": true
        }
      },

      "defaults": {
        "resource": "message",
        "operation": "send",
        "channelId": { "type": "literal", "value": "" },
        "content": { "type": "literal", "value": "" }
      },

      "outputHints": [
        { "path": "*", "description": "Discord API response" }
      ]
    }
  ]
}
```

**Scripts**: Zero.
**Nota**: `webhook/send` tem `skipCredentials: true` — não injeta bot token.

---

### ═══════════════════════════════════════
### SIMPLE #3 — OpenAI
### ═══════════════════════════════════════

**Patterns**: loadOptions para models, fieldSource para messages JSON, number params

```json
{
  "id": "integration-openai",
  "name": "OpenAI",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "psychology",

  "metadata": {
    "brandIcon": "openai.svg",
    "color": "#412991",
    "docs": "https://platform.openai.com/docs"
  },

  "credentials": {
    "id": "openAiApi",
    "name": "OpenAI API",
    "fields": [
      { "name": "apiKey", "displayName": "API Key", "type": "string", "required": true, "typeOptions": { "secret": true } },
      { "name": "organizationId", "displayName": "Organization ID", "type": "string" }
    ],
    "inject": {
      "headers": {
        "Authorization": "Bearer {{apiKey}}",
        "OpenAI-Organization": "{{organizationId}}"
      }
    },
    "test": { "method": "GET", "path": "/v1/models" }
  },

  "baseUrl": "https://api.openai.com",

  "loadOptions": {
    "getModels": {
      "request": { "method": "GET", "path": "/v1/models" },
      "dataPath": "data",
      "valuePath": "id",
      "labelPath": "id"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/openai",
      "label": "OpenAI",
      "description": "Use GPT models, DALL-E and more",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "chat",
          "options": [
            { "label": "Chat", "value": "chat" },
            { "label": "Image", "value": "image" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "complete",
          "options": [
            { "label": "Complete", "value": "complete" }
          ],
          "displayOptions": { "show": { "resource": ["chat"] } }
        },
        {
          "name": "operationImage",
          "displayName": "Operation",
          "type": "options",
          "default": "generate",
          "options": [
            { "label": "Generate", "value": "generate" }
          ],
          "displayOptions": { "show": { "resource": ["image"] } }
        },

        {
          "name": "model",
          "displayName": "Model",
          "type": "options",
          "required": true,
          "default": "gpt-4o",
          "typeOptions": { "loadOptions": "getModels" },
          "displayOptions": { "show": { "resource": ["chat"] } }
        },
        {
          "name": "messages",
          "displayName": "Messages",
          "type": "json",
          "required": true,
          "default": "[{\"role\": \"user\", \"content\": \"Hello\"}]",
          "hint": "Array of {role, content} objects",
          "displayOptions": { "show": { "resource": ["chat"] } }
        },
        {
          "name": "systemPrompt",
          "displayName": "System Prompt",
          "type": "fieldSource",
          "typeOptions": { "multiline": true, "rows": 4 },
          "hint": "Optional system message prepended to messages",
          "displayOptions": { "show": { "resource": ["chat"] } }
        },

        {
          "name": "chatOptions",
          "displayName": "Options",
          "type": "group",
          "displayOptions": { "show": { "resource": ["chat"] } },
          "fields": [
            { "name": "temperature", "displayName": "Temperature", "type": "number", "default": 1, "typeOptions": { "min": 0, "max": 2, "precision": 1 } },
            { "name": "max_tokens", "displayName": "Max Tokens", "type": "number", "typeOptions": { "min": 1, "max": 128000 } },
            { "name": "top_p", "displayName": "Top P", "type": "number", "default": 1, "typeOptions": { "min": 0, "max": 1, "precision": 2 } },
            { "name": "frequency_penalty", "displayName": "Frequency Penalty", "type": "number", "default": 0, "typeOptions": { "min": -2, "max": 2 } }
          ]
        },

        {
          "name": "prompt",
          "displayName": "Prompt",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 3 },
          "displayOptions": { "show": { "resource": ["image"] } }
        },
        {
          "name": "imageModel",
          "displayName": "Model",
          "type": "options",
          "default": "dall-e-3",
          "options": [
            { "label": "DALL-E 3", "value": "dall-e-3" },
            { "label": "DALL-E 2", "value": "dall-e-2" }
          ],
          "displayOptions": { "show": { "resource": ["image"] } }
        },
        {
          "name": "size",
          "displayName": "Size",
          "type": "options",
          "default": "1024x1024",
          "options": [
            { "label": "1024x1024", "value": "1024x1024" },
            { "label": "1792x1024", "value": "1792x1024" },
            { "label": "1024x1792", "value": "1024x1792" }
          ],
          "displayOptions": { "show": { "resource": ["image"] } }
        },
        {
          "name": "quality",
          "displayName": "Quality",
          "type": "options",
          "default": "standard",
          "options": [
            { "label": "Standard", "value": "standard" },
            { "label": "HD", "value": "hd" }
          ],
          "displayOptions": { "show": { "resource": ["image"] } }
        },
        {
          "name": "n",
          "displayName": "Count",
          "type": "number",
          "default": 1,
          "typeOptions": { "min": 1, "max": 10 },
          "displayOptions": { "show": { "resource": ["image"] } }
        }
      ],

      "operations": {
        "chat/complete": {
          "request": { "method": "POST", "path": "/v1/chat/completions" },
          "output": { "dataPath": "choices[0].message" }
        },
        "image/generate": {
          "request": { "method": "POST", "path": "/v1/images/generations" },
          "output": { "dataPath": "data" }
        }
      },

      "defaults": {
        "resource": "chat",
        "operation": "complete",
        "model": "gpt-4o",
        "messages": "[{\"role\": \"user\", \"content\": \"\"}]"
      },

      "outputHints": [
        { "path": "*", "description": "OpenAI API response" }
      ]
    }
  ]
}
```

**LoadOptions**: 1 (getModels — simple dataPath extraction). Zero transform scripts.

---

### ═══════════════════════════════════════
### MEDIUM #1 — Slack
### ═══════════════════════════════════════

**Patterns**: OAuth2, loadOptions (channels, users), json para blocks, fieldSource

```json
{
  "id": "integration-slack",
  "name": "Slack",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "chat",

  "metadata": {
    "brandIcon": "slack.svg",
    "color": "#4A154B",
    "docs": "https://api.slack.com/methods"
  },

  "credentials": {
    "id": "slackOAuth2",
    "name": "Slack OAuth2",
    "fields": [
      { "name": "clientId", "displayName": "Client ID", "type": "string", "required": true },
      { "name": "clientSecret", "displayName": "Client Secret", "type": "string", "required": true, "typeOptions": { "secret": true } }
    ],
    "inject": {
      "oauth2": {
        "authUrl": "https://slack.com/oauth/v2/authorize",
        "tokenUrl": "https://slack.com/api/oauth.v2.access",
        "scopes": ["chat:write", "channels:read", "users:read"]
      }
    },
    "test": { "method": "POST", "path": "/api/auth.test" }
  },

  "baseUrl": "https://slack.com",

  "loadOptions": {
    "getChannels": {
      "request": { "method": "GET", "path": "/api/conversations.list?types=public_channel,private_channel&limit=200" },
      "dataPath": "channels",
      "valuePath": "id",
      "labelPath": "name"
    },
    "getUsers": {
      "request": { "method": "GET", "path": "/api/users.list?limit=200" },
      "dataPath": "members",
      "valuePath": "id",
      "labelPath": "real_name"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/slack",
      "label": "Slack",
      "description": "Send messages and manage Slack",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "message",
          "options": [
            { "label": "Message", "value": "message" },
            { "label": "Channel", "value": "channel" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "send",
          "options": [
            { "label": "Send", "value": "send" }
          ],
          "displayOptions": { "show": { "resource": ["message"] } }
        },
        {
          "name": "operationChannel",
          "displayName": "Operation",
          "type": "options",
          "default": "getMany",
          "options": [
            { "label": "Get Many", "value": "getMany" }
          ],
          "displayOptions": { "show": { "resource": ["channel"] } }
        },

        {
          "name": "channel",
          "displayName": "Channel",
          "type": "options",
          "required": true,
          "typeOptions": { "loadOptions": "getChannels" },
          "displayOptions": { "show": { "resource": ["message"] } }
        },
        {
          "name": "text",
          "displayName": "Message Text",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 4 },
          "displayOptions": { "show": { "resource": ["message"], "operation": ["send"] } }
        },

        {
          "name": "messageOptions",
          "displayName": "Options",
          "type": "group",
          "displayOptions": { "show": { "resource": ["message"] } },
          "fields": [
            {
              "name": "blocks",
              "displayName": "Block Kit JSON",
              "type": "json",
              "hint": "Rich message blocks (overrides text)"
            },
            {
              "name": "thread_ts",
              "displayName": "Thread Timestamp",
              "type": "fieldSource",
              "hint": "Reply in a specific thread"
            },
            {
              "name": "unfurl_links",
              "displayName": "Unfurl Links",
              "type": "boolean",
              "default": true
            },
            {
              "name": "unfurl_media",
              "displayName": "Unfurl Media",
              "type": "boolean",
              "default": true
            }
          ]
        },

        {
          "name": "types",
          "displayName": "Channel Types",
          "type": "multiOptions",
          "default": ["public_channel"],
          "options": [
            { "label": "Public", "value": "public_channel" },
            { "label": "Private", "value": "private_channel" },
            { "label": "DM", "value": "im" },
            { "label": "Group DM", "value": "mpim" }
          ],
          "displayOptions": { "show": { "resource": ["channel"] } }
        },
        {
          "name": "limit",
          "displayName": "Limit",
          "type": "number",
          "default": 50,
          "typeOptions": { "min": 1, "max": 1000 },
          "displayOptions": { "show": { "resource": ["channel"] } }
        }
      ],

      "operations": {
        "message/send": {
          "request": { "method": "POST", "path": "/api/chat.postMessage" },
          "output": { "dataPath": "message" }
        },
        "channel/getMany": {
          "request": { "method": "GET", "path": "/api/conversations.list" },
          "output": { "dataPath": "channels" }
        }
      },

      "defaults": {
        "resource": "message",
        "operation": "send",
        "channel": "",
        "text": { "type": "literal", "value": "" }
      },

      "outputHints": [
        { "path": "*", "description": "Slack API response" }
      ]
    }
  ]
}
```

---

### ═══════════════════════════════════════
### MEDIUM #2 — Airtable
### ═══════════════════════════════════════

**Patterns**: 3 levels cascading (base→table→fields), array para dynamic fields, preScript

```json
{
  "id": "integration-airtable",
  "name": "Airtable",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "grid_on",

  "metadata": {
    "brandIcon": "airtable.svg",
    "color": "#18BFFF",
    "docs": "https://airtable.com/developers/web/api"
  },

  "credentials": {
    "id": "airtableApi",
    "name": "Airtable Personal Token",
    "fields": [
      { "name": "accessToken", "displayName": "Personal Access Token", "type": "string", "required": true, "typeOptions": { "secret": true } }
    ],
    "inject": {
      "headers": { "Authorization": "Bearer {{accessToken}}" }
    },
    "test": { "method": "GET", "path": "/v0/meta/whoami" }
  },

  "baseUrl": "https://api.airtable.com",

  "loadOptions": {
    "getBases": {
      "request": { "method": "GET", "path": "/v0/meta/bases" },
      "dataPath": "bases",
      "valuePath": "id",
      "labelPath": "name"
    },
    "getTables": {
      "request": { "method": "GET", "path": "/v0/meta/bases/{{dependsOn.baseId}}/tables" },
      "dataPath": "tables",
      "valuePath": "id",
      "labelPath": "name"
    },
    "getFields": {
      "request": { "method": "GET", "path": "/v0/meta/bases/{{dependsOn.baseId}}/tables" },
      "transform": "function transform(data) { var fields = []; for (var i = 0; i < data.tables.length; i++) { var t = data.tables[i]; for (var j = 0; j < t.fields.length; j++) { var f = t.fields[j]; fields.push({ label: t.name + ' / ' + f.name, value: f.id }); } } return fields; }"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/airtable",
      "label": "Airtable",
      "description": "Manage Airtable records",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "record",
          "options": [{ "label": "Record", "value": "record" }]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "create",
          "options": [
            { "label": "Create", "value": "create" },
            { "label": "Get Many", "value": "getMany" },
            { "label": "Update", "value": "update" },
            { "label": "Delete", "value": "delete" }
          ]
        },

        {
          "name": "baseId",
          "displayName": "Base",
          "type": "options",
          "required": true,
          "in": "path",
          "typeOptions": { "loadOptions": "getBases" }
        },
        {
          "name": "tableId",
          "displayName": "Table",
          "type": "options",
          "required": true,
          "in": "path",
          "typeOptions": {
            "loadOptions": "getTables",
            "dependsOn": ["baseId"]
          }
        },

        {
          "name": "recordId",
          "displayName": "Record ID",
          "type": "fieldSource",
          "required": true,
          "in": "path",
          "displayOptions": { "show": { "operation": ["update", "delete"] } }
        },

        {
          "name": "columns",
          "displayName": "Fields",
          "type": "array",
          "required": true,
          "typeOptions": {
            "loadOptions": "getFields",
            "dependsOn": ["baseId", "tableId"],
            "itemLabel": "Field"
          },
          "displayOptions": { "show": { "operation": ["create", "update"] } },
          "fields": [
            { "name": "field", "displayName": "Field Name", "type": "string" },
            { "name": "value", "displayName": "Value", "type": "fieldSource" }
          ]
        },

        {
          "name": "listOptions",
          "displayName": "Options",
          "type": "group",
          "displayOptions": { "show": { "operation": ["getMany"] } },
          "fields": [
            { "name": "maxRecords", "displayName": "Max Records", "type": "number", "default": 100, "typeOptions": { "min": 1, "max": 1000 } },
            { "name": "filterByFormula", "displayName": "Filter Formula", "type": "fieldSource", "hint": "Airtable formula syntax" },
            { "name": "sort", "displayName": "Sort (JSON)", "type": "json", "hint": "[{\"field\":\"Name\",\"direction\":\"asc\"}]" }
          ]
        }
      ],

      "operations": {
        "record/create": {
          "preScript": "scripts/record/create.pre.js",
          "request": { "method": "POST", "path": "/v0/{{fields.baseId}}/{{fields.tableId}}" }
        },
        "record/getMany": {
          "request": { "method": "GET", "path": "/v0/{{fields.baseId}}/{{fields.tableId}}" },
          "postScript": "scripts/record/list.post.js",
          "output": { "dataPath": "records" }
        },
        "record/update": {
          "preScript": "scripts/record/update.pre.js",
          "request": { "method": "PATCH", "path": "/v0/{{fields.baseId}}/{{fields.tableId}}/{{fields.recordId}}" }
        },
        "record/delete": {
          "request": { "method": "DELETE", "path": "/v0/{{fields.baseId}}/{{fields.tableId}}/{{fields.recordId}}" }
        }
      },

      "defaults": {
        "resource": "record",
        "operation": "create",
        "baseId": "",
        "tableId": "",
        "columns": []
      },

      "outputHints": [
        { "path": "*", "description": "Airtable record data" }
      ]
    }
  ]
}
```

**scripts/record/create.pre.js:**
```javascript
module.exports = function(input) {
  const fields = {};
  (input.fields.columns || []).forEach(col => {
    if (col.field && col.value !== undefined) fields[col.field] = col.value;
  });
  return { request: { body: { fields } } };
};
```

---

### ═══════════════════════════════════════
### MEDIUM #3 — Gmail
### ═══════════════════════════════════════

**Patterns**: OAuth2 herdado, preScript para MIME, loadOptions labels, multiOptions

```json
{
  "id": "integration-gmail",
  "name": "Gmail",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "email",

  "metadata": {
    "brandIcon": "gmail.svg",
    "color": "#EA4335",
    "docs": "https://developers.google.com/gmail/api"
  },

  "credentials": {
    "id": "gmailOAuth2",
    "name": "Gmail OAuth2",
    "extends": "googleOAuth2",
    "fields": [
      { "name": "scopes", "type": "string", "default": "https://mail.google.com/", "typeOptions": { "secret": true } }
    ],
    "inject": {
      "oauth2": {
        "authUrl": "https://accounts.google.com/o/oauth2/v2/auth",
        "tokenUrl": "https://oauth2.googleapis.com/token",
        "scopes": ["https://mail.google.com/"]
      }
    },
    "test": { "method": "GET", "path": "/gmail/v1/users/me/profile" }
  },

  "baseUrl": "https://gmail.googleapis.com",

  "loadOptions": {
    "getLabels": {
      "request": { "method": "GET", "path": "/gmail/v1/users/me/labels" },
      "dataPath": "labels",
      "valuePath": "id",
      "labelPath": "name"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/gmail",
      "label": "Gmail",
      "description": "Send emails and manage Gmail",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "email",
          "options": [{ "label": "Email", "value": "email" }]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "send",
          "options": [
            { "label": "Send", "value": "send" },
            { "label": "Get Many", "value": "getMany" }
          ]
        },

        {
          "name": "to",
          "displayName": "To",
          "type": "fieldSource",
          "required": true,
          "hint": "recipient@email.com",
          "displayOptions": { "show": { "operation": ["send"] } }
        },
        {
          "name": "subject",
          "displayName": "Subject",
          "type": "fieldSource",
          "required": true,
          "displayOptions": { "show": { "operation": ["send"] } }
        },
        {
          "name": "body",
          "displayName": "Body",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 8 },
          "displayOptions": { "show": { "operation": ["send"] } }
        },
        {
          "name": "sendOptions",
          "displayName": "Options",
          "type": "group",
          "displayOptions": { "show": { "operation": ["send"] } },
          "fields": [
            { "name": "cc", "displayName": "CC", "type": "fieldSource" },
            { "name": "bcc", "displayName": "BCC", "type": "fieldSource" },
            { "name": "replyTo", "displayName": "Reply To", "type": "fieldSource" },
            { "name": "isHtml", "displayName": "Send as HTML", "type": "boolean", "default": false }
          ]
        },

        {
          "name": "labelIds",
          "displayName": "Labels",
          "type": "multiOptions",
          "default": ["INBOX"],
          "typeOptions": { "loadOptions": "getLabels" },
          "displayOptions": { "show": { "operation": ["getMany"] } }
        },
        {
          "name": "q",
          "displayName": "Search Query",
          "type": "fieldSource",
          "hint": "Gmail search syntax: from:user@example.com",
          "displayOptions": { "show": { "operation": ["getMany"] } }
        },
        {
          "name": "maxResults",
          "displayName": "Max Results",
          "type": "number",
          "default": 10,
          "typeOptions": { "min": 1, "max": 500 },
          "displayOptions": { "show": { "operation": ["getMany"] } }
        }
      ],

      "operations": {
        "email/send": {
          "preScript": "scripts/email/send.pre.js",
          "request": { "method": "POST", "path": "/gmail/v1/users/me/messages/send" }
        },
        "email/getMany": {
          "request": { "method": "GET", "path": "/gmail/v1/users/me/messages" },
          "postScript": "scripts/email/list.post.js",
          "output": { "dataPath": "messages" }
        }
      },

      "defaults": {
        "resource": "email",
        "operation": "send",
        "to": { "type": "literal", "value": "" },
        "subject": { "type": "literal", "value": "" },
        "body": { "type": "literal", "value": "" }
      },

      "outputHints": [
        { "path": "*", "description": "Gmail API response" }
      ]
    }
  ]
}
```

**scripts/email/send.pre.js** — constrói MIME RFC 2822 e base64url encoda.

---

### ═══════════════════════════════════════
### EXTRA HARD #1 — Stripe
### ═══════════════════════════════════════

**Patterns**: form-urlencoded, nested object (address/shipping), array (metadata key-value), preScript obrigatório

```json
{
  "id": "integration-stripe",
  "name": "Stripe",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "payment",

  "metadata": {
    "brandIcon": "stripe.svg",
    "color": "#635BFF",
    "docs": "https://stripe.com/docs/api"
  },

  "credentials": {
    "id": "stripeApi",
    "name": "Stripe API",
    "fields": [
      { "name": "apiKey", "displayName": "Secret Key", "type": "string", "required": true, "typeOptions": { "secret": true } }
    ],
    "inject": {
      "headers": { "Authorization": "Bearer {{apiKey}}" }
    },
    "test": { "method": "GET", "path": "/v1/balance" }
  },

  "baseUrl": "https://api.stripe.com",

  "requestDefaults": {
    "headers": { "Content-Type": "application/x-www-form-urlencoded" }
  },

  "loadOptions": {
    "getCurrencies": {
      "request": { "method": "GET", "path": "/v1/country_specs/US" },
      "transform": "function transform(data) { var currencies = data.supported_payment_currencies || []; return currencies.map(function(c) { return { label: c.toUpperCase(), value: c }; }); }"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/stripe",
      "label": "Stripe",
      "description": "Manage customers, charges, and subscriptions",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "customer",
          "options": [
            { "label": "Customer", "value": "customer" },
            { "label": "Charge", "value": "charge" },
            { "label": "Payment Intent", "value": "paymentIntent" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "create",
          "options": [
            { "label": "Create", "value": "create" },
            { "label": "Get", "value": "get" },
            { "label": "Update", "value": "update" },
            { "label": "Delete", "value": "delete" },
            { "label": "Get Many", "value": "getMany" }
          ]
        },

        {
          "name": "resourceId",
          "displayName": "ID",
          "type": "fieldSource",
          "required": true,
          "in": "path",
          "displayOptions": { "show": { "operation": ["get", "update", "delete"] } }
        },

        {
          "name": "name",
          "displayName": "Name",
          "type": "fieldSource",
          "required": true,
          "displayOptions": { "show": { "resource": ["customer"], "operation": ["create"] } }
        },

        {
          "name": "amount",
          "displayName": "Amount (cents)",
          "type": "fieldSource",
          "required": true,
          "hint": "Amount in cents. 100 = $1.00",
          "displayOptions": { "show": { "resource": ["charge", "paymentIntent"], "operation": ["create"] } }
        },
        {
          "name": "currency",
          "displayName": "Currency",
          "type": "options",
          "required": true,
          "typeOptions": { "loadOptions": "getCurrencies" },
          "displayOptions": { "show": { "resource": ["charge", "paymentIntent"], "operation": ["create"] } }
        },
        {
          "name": "source",
          "displayName": "Source / Payment Method",
          "type": "fieldSource",
          "required": true,
          "displayOptions": { "show": { "resource": ["charge"], "operation": ["create"] } }
        },

        {
          "name": "customerFields",
          "displayName": "Additional Fields",
          "type": "group",
          "displayOptions": { "show": { "resource": ["customer"], "operation": ["create", "update"] } },
          "fields": [
            { "name": "email", "displayName": "Email", "type": "fieldSource" },
            { "name": "phone", "displayName": "Phone", "type": "fieldSource" },
            { "name": "description", "displayName": "Description", "type": "fieldSource" },
            {
              "name": "address",
              "displayName": "Address",
              "type": "object",
              "fields": [
                { "name": "line1", "displayName": "Line 1", "type": "fieldSource" },
                { "name": "line2", "displayName": "Line 2", "type": "fieldSource" },
                { "name": "city", "displayName": "City", "type": "fieldSource" },
                { "name": "state", "displayName": "State", "type": "fieldSource" },
                { "name": "country", "displayName": "Country", "type": "fieldSource", "hint": "ISO 3166-1 alpha-2" },
                { "name": "postal_code", "displayName": "Postal Code", "type": "fieldSource" }
              ]
            },
            {
              "name": "metadata",
              "displayName": "Metadata",
              "type": "array",
              "typeOptions": { "itemLabel": "Item" },
              "fields": [
                { "name": "key", "displayName": "Key", "type": "string" },
                { "name": "value", "displayName": "Value", "type": "fieldSource" }
              ]
            }
          ]
        },

        {
          "name": "chargeFields",
          "displayName": "Additional Fields",
          "type": "group",
          "displayOptions": { "show": { "resource": ["charge"], "operation": ["create"] } },
          "fields": [
            { "name": "description", "displayName": "Description", "type": "fieldSource" },
            { "name": "receipt_email", "displayName": "Receipt Email", "type": "fieldSource" },
            {
              "name": "shipping",
              "displayName": "Shipping",
              "type": "object",
              "fields": [
                { "name": "name", "displayName": "Recipient", "type": "fieldSource" },
                {
                  "name": "address",
                  "displayName": "Address",
                  "type": "object",
                  "fields": [
                    { "name": "line1", "displayName": "Line 1", "type": "fieldSource" },
                    { "name": "city", "displayName": "City", "type": "fieldSource" },
                    { "name": "state", "displayName": "State", "type": "fieldSource" },
                    { "name": "country", "displayName": "Country", "type": "fieldSource" },
                    { "name": "postal_code", "displayName": "Postal Code", "type": "fieldSource" }
                  ]
                }
              ]
            },
            {
              "name": "metadata",
              "displayName": "Metadata",
              "type": "array",
              "typeOptions": { "itemLabel": "Item" },
              "fields": [
                { "name": "key", "displayName": "Key", "type": "string" },
                { "name": "value", "displayName": "Value", "type": "fieldSource" }
              ]
            }
          ]
        },

        {
          "name": "limit",
          "displayName": "Limit",
          "type": "number",
          "default": 10,
          "typeOptions": { "min": 1, "max": 100 },
          "displayOptions": { "show": { "operation": ["getMany"] } }
        }
      ],

      "operations": {
        "customer/create": {
          "preScript": "scripts/customer/create.pre.js",
          "request": { "method": "POST", "path": "/v1/customers" }
        },
        "customer/get": {
          "request": { "method": "GET", "path": "/v1/customers/{{fields.resourceId}}" }
        },
        "customer/update": {
          "preScript": "scripts/customer/update.pre.js",
          "request": { "method": "POST", "path": "/v1/customers/{{fields.resourceId}}" }
        },
        "customer/delete": {
          "request": { "method": "DELETE", "path": "/v1/customers/{{fields.resourceId}}" }
        },
        "customer/getMany": {
          "request": { "method": "GET", "path": "/v1/customers" },
          "output": { "dataPath": "data" }
        },
        "charge/create": {
          "preScript": "scripts/charge/create.pre.js",
          "request": { "method": "POST", "path": "/v1/charges" }
        },
        "charge/get": {
          "request": { "method": "GET", "path": "/v1/charges/{{fields.resourceId}}" }
        },
        "paymentIntent/create": {
          "preScript": "scripts/paymentIntent/create.pre.js",
          "request": { "method": "POST", "path": "/v1/payment_intents" }
        }
      },

      "defaults": {
        "resource": "customer",
        "operation": "create",
        "name": { "type": "literal", "value": "" }
      },

      "outputHints": [
        { "path": "*", "description": "Stripe API response" }
      ]
    }
  ]
}
```

**Por que preScript é obrigatório**: Stripe usa form-urlencoded com bracket notation (`address[line1]`, `metadata[key]`). O preScript transforma:
- `object` fields → bracket notation
- `array` metadata `[{key,value}]` → flat `metadata[key]=value`
- `group` → flatten e merge

---

### ═══════════════════════════════════════
### EXTRA HARD #2 — Google Sheets
### ═══════════════════════════════════════

**Patterns**: OAuth2 herdado, 3-level cascading (spreadsheet→sheet→columns), pre+post scripts, dynamic array

```json
{
  "id": "integration-google-sheets",
  "name": "Google Sheets",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "table_chart",

  "metadata": {
    "brandIcon": "google-sheets.svg",
    "color": "#34A853",
    "docs": "https://developers.google.com/sheets/api"
  },

  "credentials": {
    "id": "googleSheetsOAuth2",
    "name": "Google Sheets OAuth2",
    "extends": "googleOAuth2",
    "fields": [
      { "name": "scopes", "type": "string", "default": "https://www.googleapis.com/auth/spreadsheets", "typeOptions": { "secret": true } }
    ],
    "inject": {
      "oauth2": {
        "authUrl": "https://accounts.google.com/o/oauth2/v2/auth",
        "tokenUrl": "https://oauth2.googleapis.com/token",
        "scopes": ["https://www.googleapis.com/auth/spreadsheets"]
      }
    },
    "test": { "method": "GET", "path": "/v4/spreadsheets?pageSize=1" }
  },

  "baseUrl": "https://sheets.googleapis.com",

  "loadOptions": {
    "getSpreadsheets": {
      "request": { "method": "GET", "path": "https://www.googleapis.com/drive/v3/files?q=mimeType='application/vnd.google-apps.spreadsheet'&pageSize=50" },
      "dataPath": "files",
      "valuePath": "id",
      "labelPath": "name"
    },
    "getSheets": {
      "request": { "method": "GET", "path": "/v4/spreadsheets/{{dependsOn.spreadsheetId}}?fields=sheets.properties" },
      "transform": "function transform(data) { return data.sheets.map(function(s) { return { label: s.properties.title, value: s.properties.title }; }); }"
    },
    "getColumns": {
      "request": { "method": "GET", "path": "/v4/spreadsheets/{{dependsOn.spreadsheetId}}/values/{{dependsOn.sheetName}}!1:1" },
      "transform": "function transform(data) { var row = data.values && data.values[0] || []; return row.map(function(col, i) { return { label: col, value: col }; }); }"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/google-sheets",
      "label": "Google Sheets",
      "description": "Read and write spreadsheet data",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "row",
          "options": [{ "label": "Row", "value": "row" }]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "append",
          "options": [
            { "label": "Append Row", "value": "append" },
            { "label": "Read Rows", "value": "read" },
            { "label": "Update Row", "value": "update" },
            { "label": "Clear", "value": "clear" }
          ]
        },

        {
          "name": "spreadsheetId",
          "displayName": "Spreadsheet",
          "type": "options",
          "required": true,
          "typeOptions": { "loadOptions": "getSpreadsheets" }
        },
        {
          "name": "sheetName",
          "displayName": "Sheet",
          "type": "options",
          "required": true,
          "typeOptions": {
            "loadOptions": "getSheets",
            "dependsOn": ["spreadsheetId"]
          }
        },

        {
          "name": "columns",
          "displayName": "Column Values",
          "type": "array",
          "required": true,
          "typeOptions": {
            "loadOptions": "getColumns",
            "dependsOn": ["spreadsheetId", "sheetName"],
            "itemLabel": "Column"
          },
          "displayOptions": { "show": { "operation": ["append", "update"] } },
          "fields": [
            { "name": "column", "displayName": "Column", "type": "string" },
            { "name": "value", "displayName": "Value", "type": "fieldSource" }
          ]
        },

        {
          "name": "range",
          "displayName": "Range",
          "type": "string",
          "default": "A:Z",
          "hint": "e.g. A1:D10, A:Z for all columns",
          "displayOptions": { "show": { "operation": ["read", "clear"] } }
        },

        {
          "name": "rowIndex",
          "displayName": "Row Number",
          "type": "fieldSource",
          "required": true,
          "hint": "1-based row number to update",
          "displayOptions": { "show": { "operation": ["update"] } }
        }
      ],

      "operations": {
        "row/append": {
          "preScript": "scripts/row/append.pre.js",
          "request": {
            "method": "POST",
            "path": "/v4/spreadsheets/{{fields.spreadsheetId}}/values/{{fields.sheetName}}!A:Z:append?valueInputOption=USER_ENTERED"
          }
        },
        "row/read": {
          "request": {
            "method": "GET",
            "path": "/v4/spreadsheets/{{fields.spreadsheetId}}/values/{{fields.sheetName}}!{{fields.range}}"
          },
          "postScript": "scripts/row/read.post.js"
        },
        "row/update": {
          "preScript": "scripts/row/update.pre.js",
          "request": {
            "method": "PUT",
            "path": "/v4/spreadsheets/{{fields.spreadsheetId}}/values/{{fields.sheetName}}!A{{fields.rowIndex}}?valueInputOption=USER_ENTERED"
          }
        },
        "row/clear": {
          "request": {
            "method": "POST",
            "path": "/v4/spreadsheets/{{fields.spreadsheetId}}/values/{{fields.sheetName}}!{{fields.range}}:clear"
          }
        }
      },

      "defaults": {
        "resource": "row",
        "operation": "append",
        "spreadsheetId": "",
        "sheetName": "",
        "columns": []
      },

      "outputHints": [
        { "path": "*", "description": "Spreadsheet data" }
      ]
    }
  ]
}
```

**scripts/row/read.post.js** — Transforma `[[header1,header2],[val1,val2]]` → `[{header1:val1, header2:val2}]`
**scripts/row/append.pre.js** — Transforma `columns: [{column:"Name",value:"John"}]` → `{values:[["John"]]}`

---

### ═══════════════════════════════════════
### EXTRA HARD #3 — Salesforce
### ═══════════════════════════════════════

**Patterns**: OAuth2 dinâmico (instanceUrl), 3-level cascading, SOQL, preScript+postScript, complexidade máxima

```json
{
  "id": "integration-salesforce",
  "name": "Salesforce",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "cloud",

  "metadata": {
    "brandIcon": "salesforce.svg",
    "color": "#00A1E0",
    "docs": "https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta"
  },

  "credentials": {
    "id": "salesforceOAuth2",
    "name": "Salesforce OAuth2",
    "fields": [
      { "name": "instanceUrl", "displayName": "Instance URL", "type": "string", "required": true, "hint": "https://yourorg.salesforce.com" },
      { "name": "clientId", "displayName": "Client ID", "type": "string", "required": true },
      { "name": "clientSecret", "displayName": "Client Secret", "type": "string", "required": true, "typeOptions": { "secret": true } }
    ],
    "inject": {
      "oauth2": {
        "authUrl": "{{instanceUrl}}/services/oauth2/authorize",
        "tokenUrl": "{{instanceUrl}}/services/oauth2/token",
        "scopes": ["api", "refresh_token"]
      }
    },
    "test": { "method": "GET", "path": "/services/data/v59.0/" }
  },

  "baseUrl": "{{credentials.instanceUrl}}",

  "loadOptions": {
    "getObjects": {
      "request": { "method": "GET", "path": "/services/data/v59.0/sobjects" },
      "transform": "function transform(data) { return data.sobjects.filter(function(s) { return s.queryable && !s.deprecatedAndHidden; }).map(function(s) { return { label: s.label, value: s.name }; }); }"
    },
    "getObjectFields": {
      "request": { "method": "GET", "path": "/services/data/v59.0/sobjects/{{dependsOn.objectType}}/describe" },
      "transform": "function transform(data) { return data.fields.map(function(f) { return { label: f.label, value: f.name }; }); }"
    },
    "getPicklistValues": {
      "request": { "method": "GET", "path": "/services/data/v59.0/sobjects/{{dependsOn.objectType}}/describe" },
      "transform": "function transform(data) { return data.fields.filter(function(f) { return f.type === 'picklist'; }).map(function(f) { return { label: f.label, value: f.name }; }); }"
    }
  },

  "nodeTypes": [
    {
      "type": "integration/salesforce",
      "label": "Salesforce",
      "description": "Manage Salesforce records and queries",

      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [{ "id": "out", "label": "Out", "position": "bottom" }],

      "properties": [
        {
          "name": "resource",
          "displayName": "Resource",
          "type": "options",
          "default": "record",
          "options": [
            { "label": "Record", "value": "record" },
            { "label": "Query", "value": "query" }
          ]
        },
        {
          "name": "operation",
          "displayName": "Operation",
          "type": "options",
          "default": "create",
          "options": [
            { "label": "Create", "value": "create" },
            { "label": "Get", "value": "get" },
            { "label": "Update", "value": "update" },
            { "label": "Delete", "value": "delete" }
          ],
          "displayOptions": { "show": { "resource": ["record"] } }
        },
        {
          "name": "operationQuery",
          "displayName": "Operation",
          "type": "options",
          "default": "soql",
          "options": [
            { "label": "Execute SOQL", "value": "soql" }
          ],
          "displayOptions": { "show": { "resource": ["query"] } }
        },

        {
          "name": "objectType",
          "displayName": "Object Type",
          "type": "options",
          "required": true,
          "typeOptions": { "loadOptions": "getObjects" },
          "displayOptions": { "show": { "resource": ["record"] } }
        },

        {
          "name": "recordId",
          "displayName": "Record ID",
          "type": "fieldSource",
          "required": true,
          "in": "path",
          "displayOptions": { "show": { "resource": ["record"], "operation": ["get", "update", "delete"] } }
        },

        {
          "name": "fieldValues",
          "displayName": "Field Values",
          "type": "array",
          "required": true,
          "typeOptions": {
            "loadOptions": "getObjectFields",
            "dependsOn": ["objectType"],
            "itemLabel": "Field"
          },
          "displayOptions": { "show": { "resource": ["record"], "operation": ["create", "update"] } },
          "fields": [
            { "name": "field", "displayName": "Field API Name", "type": "string" },
            { "name": "value", "displayName": "Value", "type": "fieldSource" }
          ]
        },

        {
          "name": "soqlQuery",
          "displayName": "SOQL Query",
          "type": "fieldSource",
          "required": true,
          "typeOptions": { "multiline": true, "rows": 5 },
          "hint": "SELECT Id, Name FROM Account WHERE...",
          "displayOptions": { "show": { "resource": ["query"] } }
        }
      ],

      "operations": {
        "record/create": {
          "preScript": "scripts/record/create.pre.js",
          "request": {
            "method": "POST",
            "path": "/services/data/v59.0/sobjects/{{fields.objectType}}"
          }
        },
        "record/get": {
          "request": {
            "method": "GET",
            "path": "/services/data/v59.0/sobjects/{{fields.objectType}}/{{fields.recordId}}"
          }
        },
        "record/update": {
          "preScript": "scripts/record/update.pre.js",
          "request": {
            "method": "PATCH",
            "path": "/services/data/v59.0/sobjects/{{fields.objectType}}/{{fields.recordId}}"
          }
        },
        "record/delete": {
          "request": {
            "method": "DELETE",
            "path": "/services/data/v59.0/sobjects/{{fields.objectType}}/{{fields.recordId}}"
          }
        },
        "query/soql": {
          "preScript": "scripts/query/soql.pre.js",
          "request": {
            "method": "GET",
            "path": "/services/data/v59.0/query"
          },
          "postScript": "scripts/query/soql.post.js",
          "output": { "dataPath": "records" }
        }
      },

      "defaults": {
        "resource": "record",
        "operation": "create",
        "objectType": "",
        "fieldValues": []
      },

      "outputHints": [
        { "path": "*", "description": "Salesforce API response" }
      ]
    }
  ]
}
```

**scripts/record/create.pre.js:**
```javascript
module.exports = function(input) {
  const body = {};
  (input.fields.fieldValues || []).forEach(fv => {
    if (fv.field) body[fv.field] = fv.value;
  });
  return { request: { body } };
};
```

**scripts/query/soql.pre.js:**
```javascript
module.exports = function(input) {
  return { request: { query: { q: input.fields.soqlQuery } } };
};
```

**scripts/query/soql.post.js:**
```javascript
module.exports = function(input) {
  const records = (input.response.records || []).map(r => {
    const { attributes, ...data } = r;
    return data;
  });
  return { output: records };
};
```

---

## 5. Matriz de Cobertura — 9 Plugins

| Feature | TEL | DIS | OAI | SLK | AIR | GML | STR | GSH | SF |
|---------|-----|-----|-----|-----|-----|-----|-----|-----|-----|
| **fieldSource** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **options** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **multiOptions** | — | — | — | ✅ | — | ✅ | — | — | — |
| **group** | ✅ | — | ✅ | ✅ | ✅ | ✅ | ✅ | — | — |
| **object** | — | ✅ | — | — | — | — | ✅ | — | — |
| **array** | — | ✅ | — | — | ✅ | — | ✅ | ✅ | ✅ |
| **object nested** | — | ✅ | — | — | — | — | ✅ | — | — |
| **json** | — | — | ✅ | ✅ | ✅ | — | — | — | — |
| **number** | — | ✅ | ✅ | ✅ | ✅ | ✅ | — | — | — |
| **boolean** | ✅ | — | — | ✅ | — | ✅ | — | — | — |
| **displayOptions** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **loadOptions** | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **dependsOn** | — | — | — | — | ✅ | — | — | ✅ | ✅ |
| **preScript** | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| **postScript** | — | — | — | — | ✅ | — | — | ✅ | ✅ |
| **requestName** | ✅ | — | — | — | — | — | — | — | — |
| **API Key** | ✅ | ✅ | ✅ | — | ✅ | — | ✅ | — | — |
| **OAuth2** | — | — | — | ✅ | — | ✅ | — | ✅ | ✅ |
| **OAuth2 extends** | — | — | — | — | — | ✅ | — | ✅ | — |
| **skipCredentials** | — | ✅ | — | — | — | — | — | — | — |
| **baseUrl dinâmico** | — | — | — | — | — | — | — | — | ✅ |

### Complexidade por Plugin

| Plugin | Fields | loadOptions | Transforms | Cascading |
|--------|--------|-------------|------------|-----------|
| Telegram | 8 | 0 | 0 | 0 |
| Discord | 10 | 0 | 0 | 0 |
| OpenAI | 12 | 1 | 0 | 0 |
| Slack | 10 | 2 | 0 | 0 |
| Airtable | 9 | 3 | 1 | 2-level |
| Gmail | 10 | 1 | 0 | 0 |
| Stripe | 18 | 1 | 1 | 0 |
| Google Sheets | 8 | 3 | 2 | 3-level |
| Salesforce | 8 | 3 | 3 | 2-level |

---

## 6. Mudanças Necessárias no Sistema Existente

### 6.1 workflow-sdk (NodePropertyDefinition)

```diff
- type NodePropertyType = 'string' | 'number' | 'boolean' | 'options' | 'json';
+ type NodePropertyType = 'string' | 'number' | 'boolean' | 'options' | 'json'
+   | 'fieldSource' | 'multiOptions' | 'group' | 'object' | 'array';
```

Adicionar campos opcionais:
```typescript
interface NodePropertyDefinition {
  // ... existing fields ...
+ requestName?: string;
+ in?: 'body' | 'query' | 'path';
+ fields?: NodePropertyDefinition[];   // Para group, object, array
+ allowedSources?: SourceType[];       // Para fieldSource
+ typeOptions?: NodePropertyTypeOptions;
}
```

### 6.2 DynamicNodeForm

Estender para renderizar os novos types:
- `fieldSource` → `<FieldSourceSelector />` (JÁ EXISTE)
- `multiOptions` → `<q-select multiple />`
- `group` → seção colapsável com "Add Field" button
- `object` → sub-form nested com campos fixos
- `array` → lista repetível com "Add Item" / "Remove"

### 6.3 IntegrationExecutor (Go)

Novo executor que:
1. Lê manifest do cache
2. Resolve FieldSourceValues (state, event, literal, node_output)
3. Aplica requestName mapping
4. Constrói pipeline (preScript? → HTTP → postScript?)
5. Suspende DAG com waitType=plugin_pipeline

### 6.4 ManifestLoader (JS)

Função que converte manifest.json → WorkflowPlugin:
```typescript
function loadIntegrationManifest(manifest: IntegrationManifest): WorkflowPlugin {
  return {
    id: manifest.id,
    name: manifest.name,
    version: manifest.version,
    category: manifest.category as PluginCategory,
    icon: manifest.icon,
    nodeTypes: manifest.nodeTypes.map(nt => ({
      type: nt.type,
      label: nt.label,
      // ... map all PluginNodeType fields
      properties: nt.properties,
      defaults: nt.defaults,
      outputHints: nt.outputHints,
    })),
  };
}
```

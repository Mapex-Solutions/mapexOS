# Plugin SDK — Interfaces Completas (v1)

> Referência definitiva de todas as interfaces do workflow-sdk para plugins.
> Cada campo documentado com explicação e exemplo real.

---

## 1. NodePropertyType

Tipos de controle de form que o DynamicNodeForm renderiza.

```typescript
/**
 * Tipos de propriedade suportados para forms declarativos de nodes.
 *
 * Cada tipo determina QUAL componente de UI é renderizado no painel de configuração.
 * O DynamicNodeForm lê o type e renderiza o controle correspondente.
 *
 * Tipos básicos (já existem):
 * - 'string'    → QInput text
 * - 'number'    → QInput number
 * - 'boolean'   → QToggle
 * - 'options'   → QSelect single
 * - 'json'      → Monaco Editor (JSON mode)
 *
 * Tipos novos (v1 plugin system):
 * - 'fieldSource'      → FieldSourceSelector (literal/state/event/node_output)
 * - 'collection'       → Grupo expansível de campos opcionais
 * - 'fixedCollection'  → Array de structs (key-value pairs, headers, params)
 * - 'multiOptions'     → QSelect multiple
 * - 'dateTime'         → QInput date + time picker
 * - 'hidden'           → Não renderiza — valor interno armazenado no config
 * - 'notice'           → QBanner informativo — não é campo, é mensagem para o user
 */
type NodePropertyType =
  | 'string'
  | 'number'
  | 'boolean'
  | 'options'
  | 'json'
  | 'fieldSource'
  | 'collection'
  | 'fixedCollection'
  | 'multiOptions'
  | 'dateTime'
  | 'hidden'
  | 'notice';
```

### 1.1 `string`

Renderiza um input de texto simples. Pode virar textarea com `typeOptions.multiline`.

```json
{
  "name": "chatId",
  "displayName": "Chat ID",
  "type": "string",
  "default": "",
  "hint": "ID numérico do chat ou @username do canal",
  "required": true,
  "typeOptions": {
    "placeholder": "Ex: -1001234567890"
  }
}
```

### 1.2 `number`

Renderiza input numérico com validação min/max.

```json
{
  "name": "timeout",
  "displayName": "Timeout (seconds)",
  "type": "number",
  "default": 30,
  "typeOptions": {
    "minValue": 1,
    "maxValue": 300
  }
}
```

### 1.3 `boolean`

Renderiza toggle on/off.

```json
{
  "name": "disableNotification",
  "displayName": "Silent Message",
  "type": "boolean",
  "default": false
}
```

### 1.4 `options`

Renderiza dropdown single-select.

```json
{
  "name": "parseMode",
  "displayName": "Parse Mode",
  "type": "options",
  "default": "HTML",
  "options": [
    { "label": "None", "value": "" },
    { "label": "HTML", "value": "HTML" },
    { "label": "MarkdownV2", "value": "MarkdownV2" }
  ]
}
```

### 1.5 `json`

Renderiza Monaco Editor no modo JSON. Para configurações complexas como body de request, schemas, mappings.

```json
{
  "name": "requestBody",
  "displayName": "Request Body",
  "type": "json",
  "default": "{}",
  "typeOptions": {
    "editor": "json",
    "rows": 10
  }
}
```

### 1.6 `fieldSource` (NOVO)

Renderiza o FieldSourceSelector — permite escolher a ORIGEM do valor: literal, state, event, node_output.
Usado quando o valor de um campo pode vir de diferentes fontes no workflow.

```json
{
  "name": "chatId",
  "displayName": "Chat ID",
  "type": "fieldSource",
  "default": { "type": "literal", "value": "" },
  "required": true,
  "allowedSources": ["literal", "state", "event", "node_output"],
  "hint": "ID numérico do chat ou @username do canal"
}
```

Valor salvo no config:
```json
{ "chatId": { "type": "event", "value": "payload.chat_id" } }
```

### 1.7 `collection` (NOVO)

Grupo de campos opcionais que o user pode expandir. Renderiza como um bloco colapsável
com botão "Add option". Cada sub-campo é independente — user adiciona só o que precisa.

Caso de uso: "Additional Options", "Advanced Settings".

```json
{
  "name": "additionalOptions",
  "displayName": "Additional Options",
  "type": "collection",
  "default": {},
  "placeholder": "Add option",
  "values": [
    {
      "name": "parseMode",
      "displayName": "Parse Mode",
      "type": "options",
      "default": "HTML",
      "options": [
        { "label": "HTML", "value": "HTML" },
        { "label": "MarkdownV2", "value": "MarkdownV2" }
      ]
    },
    {
      "name": "disableNotification",
      "displayName": "Silent Message",
      "type": "boolean",
      "default": false
    },
    {
      "name": "protectContent",
      "displayName": "Protect Content",
      "type": "boolean",
      "default": false
    }
  ]
}
```

Valor salvo no config (só os que o user adicionou):
```json
{ "additionalOptions": { "parseMode": "HTML", "protectContent": true } }
```

### 1.8 `fixedCollection` (NOVO)

Array de structs com schema fixo. Cada item tem as mesmas propriedades definidas em `values`.
Renderiza como lista com botão "Add item" e cada item é um grupo de campos.

Caso de uso: headers HTTP, query params, mapeamentos key-value.

```json
{
  "name": "queryParameters",
  "displayName": "Query Parameters",
  "type": "fixedCollection",
  "default": [],
  "placeholder": "Add parameter",
  "typeOptions": {
    "multipleValues": true
  },
  "values": [
    {
      "name": "key",
      "displayName": "Name",
      "type": "string",
      "default": ""
    },
    {
      "name": "value",
      "displayName": "Value",
      "type": "fieldSource",
      "default": { "type": "literal", "value": "" },
      "allowedSources": ["literal", "state", "event", "node_output"]
    }
  ]
}
```

Valor salvo no config:
```json
{
  "queryParameters": [
    { "key": "limit", "value": { "type": "literal", "value": "10" } },
    { "key": "offset", "value": { "type": "state", "value": "currentPage" } }
  ]
}
```

### 1.9 `multiOptions` (NOVO)

Dropdown multi-select. User pode selecionar múltiplos valores.

```json
{
  "name": "channels",
  "displayName": "Channels",
  "type": "multiOptions",
  "default": [],
  "options": [
    { "label": "#general", "value": "C01234" },
    { "label": "#engineering", "value": "C05678" },
    { "label": "#sales", "value": "C09012" }
  ],
  "hint": "Select one or more channels"
}
```

Valor salvo no config:
```json
{ "channels": ["C01234", "C05678"] }
```

### 1.10 `dateTime` (NOVO)

Renderiza date picker com opção de time. Para agendamento, filtros por data, etc.

```json
{
  "name": "scheduledAt",
  "displayName": "Schedule Send",
  "type": "dateTime",
  "default": "",
  "hint": "Leave empty to send immediately",
  "typeOptions": {
    "dateOnly": false
  }
}
```

Valor salvo no config:
```json
{ "scheduledAt": "2026-03-15T14:30:00Z" }
```

### 1.11 `hidden` (NOVO)

Campo que NÃO renderiza no form. Armazena valores internos no config que o backend precisa
mas o user não deve editar manualmente.

```json
{
  "name": "apiVersion",
  "displayName": "API Version",
  "type": "hidden",
  "default": "v2"
}
```

Valor salvo no config:
```json
{ "apiVersion": "v2" }
```

### 1.12 `notice` (NOVO)

NÃO é um campo real — não salva nada no config. É uma mensagem informativa exibida no form
para orientar o user. Renderiza como QBanner.

```json
{
  "name": "tokenNotice",
  "displayName": "How to get your Bot Token",
  "type": "notice",
  "default": "Open Telegram, search for @BotFather, send /newbot and follow the steps."
}
```

---

## 2. NodePropertyDefinition

Interface completa de uma propriedade declarativa de node.

```typescript
/**
 * Definição declarativa de uma propriedade de node.
 *
 * O DynamicNodeForm lê um array de NodePropertyDefinition e renderiza
 * o formulário automaticamente. Cada propriedade mapeia para config[name].
 *
 * Fluxo:
 * 1. Plugin define properties[] no nodeType
 * 2. DynamicNodeForm recebe properties[] + config
 * 3. Para cada property, renderiza o controle baseado no type
 * 4. User edita → emite update:config → config[name] = novo valor
 */
interface NodePropertyDefinition {
  /**
   * Chave da propriedade no config do node.
   * Mapeia diretamente para config[name].
   *
   * Exemplo: name: "chatId" → config.chatId = "123456"
   *
   * REGRA: camelCase, sem espaços, único dentro do nodeType.
   */
  name: string;

  /**
   * Label exibido no formulário ao lado do campo.
   *
   * Exemplo: "Chat ID", "Bot Token", "Parse Mode"
   */
  displayName: string;

  /**
   * Tipo do controle de UI que será renderizado.
   * Ver seção 1 para detalhes de cada tipo.
   */
  type: NodePropertyType;

  /**
   * Valor padrão quando o node é criado.
   *
   * O tipo do default depende do type:
   * - string/number/boolean/options/hidden → valor primitivo
   * - json → string com JSON (ex: "{}")
   * - fieldSource → FieldSourceValue object (ex: { type: "literal", value: "" })
   * - collection → {} (objeto vazio)
   * - fixedCollection → [] (array vazio)
   * - multiOptions → [] (array vazio)
   * - dateTime → "" (string vazia)
   * - notice → string com a mensagem a exibir
   */
  default: unknown;

  /**
   * Texto de ajuda exibido abaixo do campo.
   * Aparece como caption/subtitle menor.
   *
   * Exemplo: "Obtenha com @BotFather no Telegram"
   */
  hint?: string;

  /**
   * Marca o campo como obrigatório.
   * DynamicNodeForm mostra asterisco e valida antes de salvar.
   *
   * Default: false
   */
  required?: boolean;

  /**
   * Lista de opções para types 'options' e 'multiOptions'.
   *
   * Exemplo:
   * [
   *   { "label": "Send Text", "value": "sendText" },
   *   { "label": "Send Photo", "value": "sendPhoto" }
   * ]
   *
   * Para options dinâmicas carregadas via API, usar typeOptions.loadOptionsMethod.
   */
  options?: Array<{ label: string; value: string | number }>;

  /**
   * Visibilidade condicional — campo só aparece quando outros campos têm valores específicos.
   *
   * Exemplo: mostrar "Text" só quando operation = "sendText":
   * { "show": { "operation": ["sendText"] } }
   *
   * Múltiplas condições = AND lógico:
   * { "show": { "resource": ["message"], "operation": ["sendText", "editText"] } }
   * → Aparece quando resource É "message" E operation É "sendText" OU "editText"
   */
  displayOptions?: {
    show?: Record<string, unknown[]>;
  };

  // ──────────────────────────────────────────────────────────────────────
  // Campos novos — FieldSource
  // ──────────────────────────────────────────────────────────────────────

  /**
   * Fontes permitidas no FieldSourceSelector.
   * Só usado quando type = 'fieldSource'.
   *
   * Se omitido, TODAS as fontes são permitidas.
   *
   * Fontes disponíveis:
   * - 'literal'     → user digita valor fixo
   * - 'state'       → lê variável do state do workflow
   * - 'event'       → lê campo do evento que trigou o workflow
   * - 'node_output' → lê output de um node anterior
   * - 'input'       → lê variável de input externo
   * - 'variable'    → lê variável do contexto
   *
   * Exemplo: campo Chat ID pode vir de qualquer fonte:
   * ["literal", "state", "event", "node_output"]
   */
  allowedSources?: Array<'event' | 'state' | 'input' | 'variable' | 'literal' | 'node_output'>;

  /**
   * Nome do campo no request HTTP enviado à API externa.
   * Usado quando o nome no config (camelCase) difere do nome na API.
   *
   * Exemplo:
   * name: "chatId"           → como fica no config do node
   * requestName: "chat_id"   → como vai no body/query do HTTP request
   *
   * O backend usa requestName ao montar o request. Se omitido, usa name.
   */
  requestName?: string;

  /**
   * Marca o campo como sensível.
   * Quando true, o backend aplica Envelope Encryption:
   * - Encrypta o valor com DEK da org
   * - Salva preview mascarado (ex: "123*****765")
   * - Salva ciphertext no campo fieldNameSecret
   *
   * No form, renderiza como input type="password" com botão [Edit].
   *
   * Exemplo: API keys, tokens, senhas.
   */
  isSecret?: boolean;

  // ──────────────────────────────────────────────────────────────────────
  // Campos novos — Collection / FixedCollection
  // ──────────────────────────────────────────────────────────────────────

  /**
   * Sub-propriedades para types 'collection' e 'fixedCollection'.
   *
   * Para 'collection': cada value é um campo opcional que o user pode adicionar.
   * Para 'fixedCollection': cada value é um campo que aparece em CADA item do array.
   *
   * Exemplo collection (optional group):
   * values: [
   *   { name: "parseMode", type: "options", ... },
   *   { name: "silent", type: "boolean", ... }
   * ]
   * → User clica "Add option" e escolhe qual campo adicionar
   *
   * Exemplo fixedCollection (array of structs):
   * values: [
   *   { name: "key", type: "string", ... },
   *   { name: "value", type: "string", ... }
   * ]
   * → User clica "Add item" e preenche key + value juntos
   */
  values?: NodePropertyDefinition[];

  /**
   * Texto do botão para adicionar item em collection/fixedCollection.
   *
   * Exemplo: "Add option", "Add parameter", "Add header"
   * Default: "Add item"
   */
  placeholder?: string;

  // ──────────────────────────────────────────────────────────────────────
  // TypeOptions — Configuração avançada do controle de UI
  // ──────────────────────────────────────────────────────────────────────

  /**
   * Opções avançadas de renderização e comportamento do campo.
   * Cada propriedade dentro de typeOptions é opcional e depende do type.
   */
  typeOptions?: NodePropertyTypeOptions;
}
```

---

## 3. NodePropertyTypeOptions

Configuração avançada do controle de UI de cada propriedade.

```typescript
/**
 * Opções avançadas para controles de formulário.
 *
 * Estas opções controlam COMO o campo é renderizado e se comporta,
 * independente do tipo básico. Por exemplo, um 'string' pode virar
 * textarea (multiline) ou input de senha (password).
 */
interface NodePropertyTypeOptions {
  /**
   * Transforma input de texto em textarea.
   * Só aplica para type: 'string' e 'fieldSource'.
   *
   * Exemplo: campo de texto de mensagem
   * { "multiline": true, "rows": 4 }
   *
   * Default: false
   */
  multiline?: boolean;

  /**
   * Número de linhas do textarea.
   * Só aplica quando multiline = true.
   *
   * Default: 3
   */
  rows?: number;

  /**
   * Renderiza input como type="password" (pontos no lugar de texto).
   * Usado para campos sensíveis durante digitação.
   *
   * DIFERENTE de isSecret:
   * - password = visual (mascara no form enquanto digita)
   * - isSecret = segurança (backend encrypta o valor no MongoDB)
   * - Geralmente usados juntos em campos de API key/token.
   *
   * Exemplo: API Key
   * { "password": true }
   *
   * Default: false
   */
  password?: boolean;

  /**
   * Tipo de editor especializado.
   * Renderiza Monaco Editor com syntax highlighting para a linguagem.
   *
   * Valores:
   * - 'json' → JSON com validação e autocomplete
   * - 'js'   → JavaScript (usado no node core/code)
   * - 'sql'  → SQL com highlight de keywords
   * - 'html' → HTML com tags
   *
   * Exemplo: campo de script SQL
   * { "editor": "sql" }
   */
  editor?: 'json' | 'js' | 'sql' | 'html';

  /**
   * Permite múltiplos valores (transforma o campo em array).
   * Renderiza botão "Add" abaixo do campo para adicionar mais entradas.
   *
   * Usado em fixedCollection para permitir N items.
   * Também pode ser usado em string/number para arrays simples.
   *
   * Exemplo: múltiplos destinatários
   * { "multipleValues": true }
   *
   * Config salvo: ["user1", "user2", "user3"]
   *
   * Default: false
   */
  multipleValues?: boolean;

  /**
   * Chave que referencia uma entrada no mapa `loadOptions` do manifest root.
   * Substitui o array estático de options[] por dropdown dinâmico.
   *
   * O manifest define os loaders no nível raiz:
   * "loadOptions": { "getChannels": { "request": {...}, "dataPath": "channels", ... } }
   *
   * A property referencia pela chave:
   * { "loadOptions": "getChannels" }
   *
   * O DynamicNodeForm chama:
   * POST /api/v1/credentials/:credentialId/load_options/:resourceKey
   * Body: { "dependsOn": { "field1": "value1" } }
   *
   * O backend decripta a credencial, resolve templates, faz HTTP call,
   * extrai via dataPath/valuePath/labelPath (simples) ou roda transform JS ES5 (complexo),
   * e retorna array de { label, value }.
   *
   * Caso de uso: listar channels do Slack, modelos do OpenAI, tabelas de um DB.
   */
  loadOptions?: string;

  /**
   * Re-carregar options quando estes campos mudam.
   * Usado junto com loadOptions para cascata de selects.
   *
   * Exemplo: tables mudam quando base muda (Airtable)
   * {
   *   "loadOptions": "getTables",
   *   "loadOptionsDependsOn": ["baseId"]
   * }
   *
   * Quando user muda "baseId" → DynamicNodeForm refaz a chamada
   * com dependsOn: { baseId: "appXXX" } e atualiza as options.
   */
  loadOptionsDependsOn?: string[];

  /**
   * Valor mínimo para type: 'number'.
   *
   * Exemplo: timeout mínimo 1 segundo
   * { "minValue": 1 }
   */
  minValue?: number;

  /**
   * Valor máximo para type: 'number'.
   *
   * Exemplo: máximo 100 resultados por página
   * { "maxValue": 100 }
   */
  maxValue?: number;

  /**
   * Texto placeholder no input.
   * Aparece em cinza quando o campo está vazio.
   *
   * Exemplo: "Enter chat ID or @username"
   */
  placeholder?: string;

  /**
   * Para type: 'dateTime' — exibe apenas date sem time picker.
   *
   * true  → só date (2026-03-15)
   * false → date + time (2026-03-15T14:30:00Z)
   *
   * Default: false
   */
  dateOnly?: boolean;
}
```

---

## 4. PluginCategory

Categorias para agrupamento no catalog lateral do workflow editor.

```typescript
/**
 * Categorias de plugins para o catalog.
 *
 * CORE (usado pelos plugins built-in, NÃO mudam):
 * - 'triggers'       → Trigger Event, Start
 * - 'logic'          → Condition
 * - 'state'          → Set State, Log, Code
 * - 'flow_control'   → Fanout, Switch, Merge, Loop, End, Goto
 * - 'timers'         → Delay, Wait Signal, Wait For
 * - 'observability'  → (reservado para métricas/tracing)
 * - 'annotations'    → Text Note, Group Frame
 *
 * PLUGINS (categorias funcionais — descrevem O QUE o plugin faz):
 * - 'messaging'      → Telegram, Slack, Discord, WhatsApp
 * - 'email'          → Gmail, SendGrid, Mailgun
 * - 'ai'             → OpenAI, Claude, Gemini
 * - 'payments'       → Stripe, PayPal
 * - 'databases'      → PostgreSQL, MySQL, MongoDB
 * - 'storage'        → S3, Google Drive, Dropbox
 * - 'crm'            → HubSpot, Salesforce
 * - 'analytics'      → Google Analytics, Mixpanel
 * - 'devops'         → GitHub, GitLab, Jira
 * - 'custom'         → Plugins do usuário
 *
 * O catalog agrupa por categoria e exibe com ícone + label.
 * Novas categorias podem ser adicionadas sem alterar o SDK —
 * o catalog renderiza qualquer string como categoria.
 */
type PluginCategory = string;
```

---

## 5. WorkflowPlugin

Interface principal de registro de plugin.

```typescript
/**
 * Plugin de workflow.
 *
 * Todo plugin implementa esta interface para registrar seus node types no editor.
 * Plugins CORE são definidos em TypeScript. Plugins EXTERNOS vêm de manifest.json
 * e são convertidos pelo ManifestLoader.
 *
 * Ciclo de vida:
 * 1. ManifestLoader.load(json) → WorkflowPlugin
 * 2. pluginRegistry.registerPlugin(plugin) → registra nodeTypes
 * 3. plugin.onActivate(context) → setup (translations, etc)
 * 4. ... user usa os nodes no editor ...
 * 5. plugin.onDeactivate() → cleanup
 */
interface WorkflowPlugin {
  /** Identificador único do plugin. Ex: "telegram", "openai", "stripe" */
  id: string;

  /** Nome para exibição. Ex: "Telegram", "OpenAI", "Stripe" */
  name: string;

  /** Versão semver. Ex: "1.0.0", "2.3.1" */
  version: string;

  /** Categoria para agrupamento no catalog. Ex: "messaging", "ai" */
  category: PluginCategory;

  /** Ícone Material Icons. Ex: "send", "psychology", "payment" */
  icon: string;

  /** Node types registrados pelo plugin. Um plugin pode ter N nodes. */
  nodeTypes: PluginNodeType[];

  /** Chamado quando o plugin é ativado no editor */
  onActivate?: (context: PluginActivationContext) => void | Promise<void>;

  /** Chamado quando o plugin é desativado */
  onDeactivate?: () => void;

  // ──────────────────────────────────────────────────────────────────────
  // Campos novos — Plugin Marketplace (só para plugins de integração)
  // ──────────────────────────────────────────────────────────────────────

  /**
   * Target de dispatch para o DAG.
   * Diz ao Workflow Service para qual subject NATS despachar o node.
   *
   * O DAG NUNCA faz I/O — sempre despacha para um serviço especializado.
   * Cada dispatch target corresponde a um microserviço independente:
   *
   * - "http"   → Triggers Service (MAPEXOS.HTTP.EXECUTE)
   * - "ai"     → AI Service (MAPEXOS.AI.EXECUTE)
   * - "mcp"    → MCP Service (MAPEXOS.MCP.EXECUTE)
   * - "script" → JS Executor (MAPEXOS.SCRIPT.EXECUTE)
   *
   * Se omitido, node executa inline no DAG (só para core nodes:
   * condition, switch, loop, set_state, merge, fanout, delay, end).
   *
   * Exemplos:
   * - Telegram: "http" (HTTP calls via Triggers Service)
   * - OpenAI: "ai" (LLM calls via AI Service)
   * - GitHub MCP: "mcp" (tool calls via MCP Service)
   */
  dispatch?: 'http' | 'ai' | 'mcp' | 'script';

  /**
   * Metadados visuais do plugin para o marketplace.
   * Usado na listagem de plugins disponíveis (registry cards).
   *
   * Exemplo:
   * {
   *   "brandIcon": "telegram.svg",
   *   "color": "#0088CC",
   *   "docs": "https://core.telegram.org/bots/api"
   * }
   */
  metadata?: PluginMetadata;

  /**
   * Definição de credenciais necessárias para autenticação.
   * Descreve quais campos o user precisa preencher (ex: API key, token).
   *
   * Se omitido, plugin não requer autenticação.
   *
   * Exemplo Telegram:
   * {
   *   "id": "telegramApi",
   *   "name": "Telegram Bot API",
   *   "fields": [
   *     { "name": "botToken", "displayName": "Bot Token", "type": "string",
   *       "required": true, "isSecret": true, "hint": "Obtenha com @BotFather" }
   *   ],
   *   "test": { "method": "GET", "path": "/bot{{botToken}}/getMe" }
   * }
   */
  credentials?: PluginCredentialDefinition;

  /**
   * URL base da API externa que o plugin consome.
   * Usado pelo backend para montar requests HTTP.
   *
   * Exemplo:
   * - Telegram: "https://api.telegram.org"
   * - Slack: "https://slack.com/api"
   * - OpenAI: "https://api.openai.com/v1"
   */
  baseUrl?: string;
}
```

---

## 6. PluginMetadata

Metadados visuais para o marketplace.

```typescript
/**
 * Metadados visuais do plugin.
 * Exibidos no card do marketplace e no catalog do editor.
 */
interface PluginMetadata {
  /**
   * Path do ícone SVG da marca (relativo ao plugin no CDN).
   * Usado quando Material Icons não representa a marca.
   *
   * Exemplo: "telegram.svg" → CDN resolve para plugins/telegram/telegram.svg
   */
  brandIcon?: string;

  /**
   * Cor da marca em hex.
   * Usada como accent color no card do marketplace e no node do canvas.
   *
   * Exemplo: "#0088CC" (Telegram blue), "#4A154B" (Slack purple)
   */
  color?: string;

  /**
   * URL da documentação oficial da API que o plugin consome.
   * Link aberto quando user clica "View docs" no painel de configuração.
   *
   * Exemplo: "https://core.telegram.org/bots/api"
   */
  docs?: string;
}
```

---

## 7. PluginCredentialDefinition

Definição dos campos de autenticação do plugin.

```typescript
/**
 * Definição de credenciais necessárias para o plugin funcionar.
 *
 * Descreve:
 * - Quais campos o user precisa preencher (API key, token, senha)
 * - Como testar se as credenciais são válidas
 *
 * O backend usa esta definição para:
 * - Saber quais campos encryptar (isSecret: true)
 * - Montar o test request para validar as credenciais
 * - Injetar as credenciais no request HTTP em runtime
 */
interface PluginCredentialDefinition {
  /**
   * Identificador único das credentials do plugin.
   * Usado internamente para referência.
   *
   * Exemplo: "telegramApi", "slackOAuth", "openAiApi"
   */
  id: string;

  /**
   * Nome descritivo exibido no form.
   *
   * Exemplo: "Telegram Bot API", "Slack OAuth2", "OpenAI API Key"
   */
  name: string;

  /**
   * Campos que o user precisa preencher.
   * Cada field gera um input no formulário de credenciais.
   */
  fields: CredentialFieldDefinition[];

  /**
   * Configuração de test request para validar as credenciais.
   * Quando o user clica [Test Connection], o backend monta e executa
   * este request usando as credenciais fornecidas.
   *
   * O path pode usar template interpolation com os nomes dos fields:
   * "/bot{{botToken}}/getMe" → substitui {{botToken}} pelo valor
   *
   * Se omitido, botão [Test Connection] não aparece.
   */
  test?: CredentialTestDefinition;
}

/**
 * Definição de um campo de credencial.
 */
interface CredentialFieldDefinition {
  /**
   * Nome do campo (camelCase).
   * Usado como chave no config e nos templates de path.
   *
   * Exemplo: "botToken", "apiKey", "clientSecret"
   */
  name: string;

  /**
   * Label exibido no form.
   *
   * Exemplo: "Bot Token", "API Key", "Client Secret"
   */
  displayName: string;

  /**
   * Tipo do input. Atualmente sempre "string".
   * Reservado para expansão futura (ex: "oauth2").
   */
  type: 'string';

  /**
   * Se o campo é obrigatório.
   */
  required: boolean;

  /**
   * Se o valor deve ser encryptado via Envelope Encryption.
   * Quando true:
   * - Backend encrypta com DEK da org
   * - Salva preview mascarado: "123*****765"
   * - Salva ciphertext no campo fieldNameSecret
   * - Frontend renderiza como password input com botão [Edit]
   */
  isSecret: boolean;

  /**
   * Texto de ajuda.
   *
   * Exemplo: "Obtenha com @BotFather no Telegram"
   */
  hint?: string;
}

/**
 * Configuração do test de credenciais.
 */
interface CredentialTestDefinition {
  /**
   * Método HTTP do test request.
   */
  method: 'GET' | 'POST';

  /**
   * Path do request (relativo ao baseUrl do plugin).
   * Suporta template interpolation: {{fieldName}}.
   *
   * Exemplo: "/bot{{botToken}}/getMe"
   * Com baseUrl "https://api.telegram.org":
   * → GET https://api.telegram.org/bot123456:ABC/getMe
   */
  path: string;
}
```

---

## 8. OperationDefinition

Mapeamento operation → HTTP request template.

```typescript
/**
 * Definição de uma operação HTTP do plugin.
 *
 * Cada operation (ex: "sendText", "sendPhoto") mapeia para um request HTTP específico.
 * O backend usa esta definição para montar e executar o request via Triggers Service.
 *
 * As operations ficam no nodeType, indexadas pelo valor do campo operation.
 * Ex: config.operation = "sendText" → operations["sendText"]
 */
interface OperationDefinition {
  /**
   * Template do HTTP request.
   */
  request: {
    /**
     * Método HTTP.
     */
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

    /**
     * Path do request (relativo ao baseUrl do plugin).
     * Suporta template interpolation:
     * - {{credentials.fieldName}} → valor da credential
     * - {{config.fieldName}} → valor do config do node
     *
     * Exemplo: "/bot{{credentials.botToken}}/sendMessage"
     */
    path: string;
  };

  /**
   * Configuração de como extrair o output do response.
   */
  output?: {
    /**
     * Path no JSON de response para extrair como output do node.
     * O Triggers Service retorna o response completo; o Workflow Service
     * extrai o valor em dataPath e salva como output do node.
     *
     * Exemplo: API Telegram retorna { ok: true, result: {...} }
     * dataPath: "result" → output do node = { message_id: 123, ... }
     *
     * Se omitido, response inteiro vira output.
     */
    dataPath?: string;
  };
}
```

---

## 9. PluginNodeType (campos novos)

Campos adicionados ao PluginNodeType existente:

```typescript
/**
 * Extensões do PluginNodeType para plugins de integração.
 * Todos os campos são opcionais — não quebram nodes CORE existentes.
 */
interface PluginNodeType {
  // ... todos os campos existentes permanecem ...

  /**
   * Mapa de operações HTTP do node.
   * Chave = valor do campo operation no config.
   * Valor = definição do request HTTP.
   *
   * Exemplo Telegram:
   * {
   *   "sendText": {
   *     "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendMessage" },
   *     "output": { "dataPath": "result" }
   *   },
   *   "sendPhoto": {
   *     "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendPhoto" },
   *     "output": { "dataPath": "result" }
   *   }
   * }
   *
   * Plugins CORE não usam este campo (executam lógica interna, não HTTP).
   */
  operations?: Record<string, OperationDefinition>;
}
```

---

## 10. HandleDefinition (sem mudanças)

Não precisa de alteração. A interface atual já é genérica o suficiente:

```typescript
interface HandleDefinition {
  id: string;
  label: string;
  position: 'top' | 'bottom' | 'left' | 'right';
  dataType?: string;
  maxConnections?: number | null;
  color?: string;
}
```

---

## 11. Resumo das mudanças no SDK

### Arquivos a alterar:

| Arquivo | Mudança |
|---------|---------|
| `workflowPlugin.interface.ts` | Expandir NodePropertyType union, adicionar campos em NodePropertyDefinition, adicionar NodePropertyTypeOptions, adicionar PluginMetadata, PluginCredentialDefinition, OperationDefinition, expandir WorkflowPlugin, expandir PluginNodeType, expandir PluginCategory |

### O que NÃO muda:

| Interface | Status |
|-----------|--------|
| HandleDefinition | Sem mudanças |
| HandleResolver | Sem mudanças |
| HandleOverrides | Sem mudanças |
| ResolvedHandles | Sem mudanças |
| ValidationResult | Sem mudanças |
| PluginActivationContext | Sem mudanças |
| Disposable | Sem mudanças |
| DynamicNodeFormProps | Sem mudanças |
| DynamicNodeFormEmits | Sem mudanças |
| CatalogGroup | Sem mudanças |
| FieldSourceValue | Sem mudanças |
| SourceType | Sem mudanças |

### Compatibilidade:

Todos os campos novos são **opcionais**. Plugins CORE existentes continuam
funcionando sem alterar uma linha — eles simplesmente não usam os campos novos.

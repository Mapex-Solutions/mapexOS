# Plugin Dispatch Architecture

> Como o DAG despacha trabalho para serviços especializados sem jamais fazer I/O.

---

## 1. Regra Absoluta

**O DAG NUNCA faz I/O. JAMAIS.**

O Workflow Service é um orquestrador puro. Ele sabe:
- Qual é o próximo node no grafo
- Pra qual subject NATS despachar
- Fazer checkpoint no KV entre steps

Ele NÃO sabe:
- Como fazer HTTP call
- Como chamar OpenAI
- Como conectar num MCP server
- Como executar JavaScript
- Como decryptar credentials
- NADA de I/O externo

Se o DAG faz I/O, ele vira bottleneck e perde escalabilidade.

---

## 2. Serviços Especializados

```
┌─────────────────────────────────────────────────────────────┐
│  WORKFLOW SERVICE                                            │
│  Orquestrador puro — zero I/O                                │
│                                                              │
│  node.dispatch →  subject NATS                               │
│  "http"        →  MAPEXOS.HTTP.EXECUTE                       │
│  "ai"          →  MAPEXOS.AI.EXECUTE                         │
│  "mcp"         →  MAPEXOS.MCP.EXECUTE                        │
│  "script"      →  MAPEXOS.SCRIPT.EXECUTE                     │
│  (sem dispatch) → inline (condition, switch, loop, set_state)│
└──────┬─────────────┬─────────────┬─────────────┬────────────┘
       │             │             │             │
  NATS │        NATS │        NATS │        NATS │
       │             │             │             │
  ┌────▼────┐  ┌─────▼────┐  ┌────▼────┐  ┌────▼─────┐
  │Triggers │  │    AI    │  │   MCP   │  │    JS    │
  │Service  │  │ Service  │  │ Service │  │ Executor │
  │         │  │          │  │         │  │          │
  │ HTTP    │  │ OpenAI   │  │ MCP     │  │ V8       │
  │ REST    │  │ Claude   │  │ servers │  │ sandbox  │
  │ SOAP    │  │ Gemini   │  │ RAG     │  │          │
  │ GraphQL │  │ Bedrock  │  │ tools   │  │          │
  │ Webhook │  │ Local    │  │ file    │  │          │
  │         │  │ models   │  │ db      │  │          │
  │         │  │          │  │ search  │  │          │
  │org_keys │  │org_keys  │  │org_keys │  │          │
  └─────────┘  └──────────┘  └─────────┘  └──────────┘
     N pods       N pods       N pods        N pods
  (CPU cheap)  (GPU/$$)     (varies)      (CPU cheap)
```

Cada serviço:
- Escala independente (pods próprios)
- Tem sua collection `org_keys` (credentials isoladas)
- Decrypta credentials localmente
- Rate limiting próprio
- Circuit breaker próprio
- Observability própria

---

## 3. Flow de Execução

### O que o DAG faz (zero I/O):

```
1. Lê próximo node do grafo (já em memória no KV)
2. Lê config do node (já em memória no KV)
3. Resolve dispatch target pelo tipo do node
4. Publica NATS: { orgId, workflowId, instanceId, nodeId, config }
   → config contém o ciphertext das credentials (NÃO decryptado)
5. KV checkpoint: status = "dispatched"
6. Espera response via NATS
7. Recebe response → salva output no KV
8. Avança para próximo node no grafo
```

### O que o serviço especializado faz (todo I/O aqui):

```
1. Consome mensagem NATS
2. Busca DEK em SUA org_keys (pelo orgId)
3. Decrypta credentials do config
4. Executa a operação:
   - Triggers: monta e executa HTTP request
   - AI: chama API do provider (OpenAI, Claude, etc)
   - MCP: conecta no MCP server e executa tool
   - JS: executa script no V8 sandbox
5. Publica response via NATS
```

### Diagrama de sequência:

```
DAG                    NATS                  AI Service            OpenAI
 │                      │                      │                    │
 │ publish(AI.EXECUTE)  │                      │                    │
 │─────────────────────>│                      │                    │
 │                      │  consume             │                    │
 │  checkpoint(KV)      │─────────────────────>│                    │
 │                      │                      │                    │
 │                      │                      │ decrypt(org_keys)  │
 │                      │                      │──────┐             │
 │                      │                      │<─────┘             │
 │                      │                      │                    │
 │                      │                      │ POST /chat/compl.  │
 │                      │                      │───────────────────>│
 │                      │                      │                    │
 │                      │                      │    response        │
 │                      │                      │<───────────────────│
 │                      │                      │                    │
 │                      │  publish(response)   │                    │
 │                      │<─────────────────────│                    │
 │                      │                      │                    │
 │  consume(response)   │                      │                    │
 │<─────────────────────│                      │                    │
 │                      │                      │                    │
 │ save output(KV)      │                      │                    │
 │ next node            │                      │                    │
```

---

## 4. Dispatch Target no Manifest

O manifest do plugin declara `dispatch` para dizer ao DAG pra onde enviar:

```json
{
  "id": "telegram",
  "category": "messaging",
  "dispatch": "http",

  "nodeTypes": [{ "type": "telegram/message", "..." : "..." }]
}
```

```json
{
  "id": "openai",
  "category": "ai",
  "dispatch": "ai",

  "nodeTypes": [{ "type": "openai/chat", "..." : "..." }]
}
```

```json
{
  "id": "github-mcp",
  "category": "devops",
  "dispatch": "mcp",

  "nodeTypes": [{ "type": "github-mcp/tool", "..." : "..." }]
}
```

| dispatch | Subject NATS | Serviço |
|----------|-------------|---------|
| `http` | `MAPEXOS.HTTP.EXECUTE` | Triggers Service |
| `ai` | `MAPEXOS.AI.EXECUTE` | AI Service |
| `mcp` | `MAPEXOS.MCP.EXECUTE` | MCP Service |
| `script` | `MAPEXOS.SCRIPT.EXECUTE` | JS Executor |
| (omitido) | inline | Workflow Service (só core nodes) |

O DAG resolve:
```go
switch manifest.Dispatch {
case "http":
    subject = "MAPEXOS.HTTP.EXECUTE"
case "ai":
    subject = "MAPEXOS.AI.EXECUTE"
case "mcp":
    subject = "MAPEXOS.MCP.EXECUTE"
case "script":
    subject = "MAPEXOS.SCRIPT.EXECUTE"
default:
    // core nodes — inline (condition, switch, loop, set_state)
    executeInline(node)
    return
}
nats.Publish(subject, payload)
```

---

## 5. Payload NATS (padrão para todos)

```json
{
  "orgId": "org_abc123",
  "workflowId": "wf_001",
  "instanceId": "inst_001",
  "nodeId": "node_tg_001",
  "nodeType": "telegram/message",
  "config": {
    "operation": "sendText",
    "chatId": { "type": "event", "value": "payload.chat_id" },
    "text": { "type": "literal", "value": "Hello!" },
    "botToken": "123*****765",
    "botTokenSecret": "U2FsdGVkX1+abc123...encrypted..."
  },
  "resolvedInputs": {
    "chatId": "987654321",
    "text": "Hello!"
  },
  "operation": {
    "request": { "method": "POST", "path": "/bot{{credentials.botToken}}/sendMessage" },
    "output": { "dataPath": "result" }
  },
  "baseUrl": "https://api.telegram.org",
  "timeout": 30000,
  "retryPolicy": { "maxAttempts": 3, "initialInterval": "5s" }
}
```

O Workflow Service resolve os `fieldSource` (event, state, node_output) e envia em `resolvedInputs`.
O serviço de destino só precisa:
1. Decryptar credentials
2. Montar request com resolvedInputs
3. Executar
4. Retornar resultado

---

## 6. Cada Serviço — Responsabilidades

### Triggers Service (já existe)

```
Responsabilidades:
- Executar HTTP requests (REST, SOAP, GraphQL, Webhook)
- Rate limiting por endpoint
- Retry com backoff
- Circuit breaker por host
- Response parsing (JSON, XML, HTML)

org_keys: SIM (credentials HTTP: API keys, OAuth tokens, basic auth)

NATS subjects:
- MAPEXOS.HTTP.EXECUTE (consumer)
- MAPEXOS.HTTP.RESPONSE (producer)
```

### AI Service (NOVO)

```
Responsabilidades:
- Chamar LLM providers (OpenAI, Anthropic, Google, AWS Bedrock)
- Gerenciar modelos locais (Ollama, vLLM)
- Streaming de tokens (se workflow suportar)
- Prompt cache (mesma prompt + model = cache hit)
- Token counting (billing por org)
- Rate limiting por provider (OpenAI tier limits)
- Fallback entre providers (OpenAI down → Claude)

org_keys: SIM (API keys por provider, por org)

NATS subjects:
- MAPEXOS.AI.EXECUTE (consumer)
- MAPEXOS.AI.RESPONSE (producer)

Escala: pods com GPU para modelos locais, CPU para API calls
```

### MCP Service (NOVO)

```
Responsabilidades:
- Conectar em MCP servers (stdio, SSE, HTTP)
- Tool discovery (listar tools disponíveis de um server)
- Tool execution (chamar tool com argumentos)
- Gerenciar lifecycle de MCP servers (connect, disconnect, reconnect)
- RAG pipelines (search, retrieve, rerank)
- File system tools, database tools, search tools

org_keys: SIM (credentials para MCP servers autenticados)

NATS subjects:
- MAPEXOS.MCP.EXECUTE (consumer)
- MAPEXOS.MCP.RESPONSE (producer)

Escala: depende dos MCP servers conectados
```

### JS Executor (já existe)

```
Responsabilidades:
- Executar JavaScript em V8 sandbox
- Timeout por script
- Memory limit por execução
- Acesso controlado a APIs (fetch, crypto)

org_keys: NÃO (scripts não precisam de credentials próprias)

NATS subjects:
- MAPEXOS.SCRIPT.EXECUTE (consumer)
- MAPEXOS.SCRIPT.RESPONSE (producer)
```

---

## 7. Credentials — Cada Serviço é Dono

```
Triggers Service DB:
  ├── triggers
  └── org_keys              ← HTTP credentials (API keys, OAuth, basic auth)

AI Service DB:
  ├── ai_providers          ← providers configurados por org
  └── org_keys              ← AI credentials (OpenAI key, Claude key, etc)

MCP Service DB:
  ├── mcp_servers           ← MCP servers registrados por org
  └── org_keys              ← MCP credentials (tokens para servers autenticados)

Workflow Service DB:
  ├── workflow_definitions
  ├── workflow_instances
  └── installed_plugins
  (SEM org_keys — Workflow Service não decrypta NADA)
```

**Workflow Service NÃO tem org_keys.** Ele envia o ciphertext como está no config.
O serviço de destino é quem decrypta.

---

## 8. Por que MCP é separado de AI

```
AI Service                          MCP Service
─────────────                       ───────────
Chama APIs de LLM                   Conecta em MCP servers
Foco: gerar texto/imagem/embedding  Foco: executar tools genéricos
Provider: OpenAI, Claude, Gemini    Server: qualquer MCP server
I/O: HTTP para cloud APIs           I/O: stdio, SSE, HTTP
Billing: por token                  Billing: por execução
GPU: sim (modelos locais)           GPU: não (geralmente)

Exemplos AI:                        Exemplos MCP:
- Chat completion                   - GitHub (issues, PRs, commits)
- Image generation                  - Database query
- Embeddings                        - File system operations
- Vision (image analysis)           - Web search (RAG)
- Text-to-speech                    - Slack (messages, channels)
- Transcription                     - Google Drive (files)
                                    - Custom tools
```

MCP pode conter RAG, database tools, file tools — nada disso é AI.
AI é especificamente LLM/embeddings/vision. Responsabilidades diferentes,
perfis de escala diferentes, billing diferente.

Um workflow pode usar os dois:
```
[Start] → [AI: classify text] → [MCP: search database] → [AI: summarize] → [End]
             dispatch: ai           dispatch: mcp           dispatch: ai
```

---

## 9. Inline vs Dispatch — Quando cada um

| Node | Dispatch | Motivo |
|------|----------|--------|
| `core/condition` | inline | Pura lógica, zero I/O |
| `core/switch` | inline | Pura lógica, zero I/O |
| `core/loop` | inline | Iterador, zero I/O |
| `core/set_state` | inline | Escreve no KV (infra do DAG) |
| `core/merge` | inline | Join lógico, zero I/O |
| `core/fanout` | inline | Fork lógico, zero I/O |
| `core/delay` | inline | Timer do runtime, zero I/O |
| `core/end` | inline | Termina instância |
| `telegram/message` | http | HTTP call externo |
| `openai/chat` | ai | API call LLM |
| `github-mcp/tool` | mcp | MCP server call |
| `core/code` | script | V8 execution |
| `core/subworkflow` | inline | Dispara outra instância via DAG |

# Matriz de Testes do DAG Walker

Cobertura exaustiva de testes para o runtime do DAG walker — o motor de execucao core que roda toda logica de negocio.

## Visao Geral da Arquitetura de Testes

Testes operam no nivel do `RuntimeService` com executors reais, mas fronteiras externas mockadas (NATS KV, JetStream, MongoDB, plugins). O harness simula callbacks async automaticamente.

```
┌─────────────────────────────────────────────────┐
│                  Test Harness                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │
│  │Definition │  │  Event   │  │  Callbacks   │  │
│  │ Builder   │  │ Payload  │  │ (async sim)  │  │
│  └────┬──────┘  └────┬─────┘  └──────┬───────┘  │
│       │              │               │           │
│  ┌────▼──────────────▼───────────────▼────────┐ │
│  │           RuntimeService.New(deps)          │ │
│  │  ┌─────────────────────────────────────┐   │ │
│  │  │      Executor Registry Real         │   │ │
│  │  │  (todos 18 executors, logica real)  │   │ │
│  │  └─────────────────────────────────────┘   │ │
│  │  ┌─────────────────────────────────────┐   │ │
│  │  │      ValueResolver Real             │   │ │
│  │  │      ConditionEvaluator Real        │   │ │
│  │  └─────────────────────────────────────┘   │ │
│  └────────────────────────────────────────────┘ │
│                      │                           │
│  ┌─────────────── MOCKADO ─────────────────────┐│
│  │ InMemoryStateRepo    (NATS KV)              ││
│  │ CapturingPublisher   (NATS JetStream)       ││
│  │ StaticDefinitionLoader                      ││
│  │ StaticInstanceLoader                        ││
│  │ NoopCredentialService                       ││
│  │ NoopPluginRepo                              ││
│  └─────────────────────────────────────────────┘│
│                      │                           │
│  ┌──────────────── SAIDA ──────────────────────┐│
│  │ WorkflowExecution (estado final)            ││
│  │ ExecutionPath (trace ordenado dos nodes)    ││
│  │ State, NodeOutputs, ErrorInfo               ││
│  │ Dispatches do publisher (capturados)        ││
│  └─────────────────────────────────────────────┘│
└─────────────────────────────────────────────────┘
```

---

## Tipos de Source

Muitos nodes aceitam `FieldValue` que podem vir de diferentes fontes. Cada tipo de source deve ser testado onde um FieldValue for aceito.

| Source | Tipo JSON | Exemplo | Resolve de |
|--------|-----------|---------|------------|
| `literal` | `{ type: "literal", value: "42" }` | Valor estatico | Valor direto |
| `state` | `{ type: "state", value: "counter" }` | Estado do workflow | `execution.State["counter"]` |
| `event` | `{ type: "event", value: "data.temp" }` | Payload do evento | `execution.EventPayload["data"]["temp"]` |
| `input` | `{ type: "input", value: "sensorId" }` | Inputs externos | `execution.ExternalInputs["sensorId"]` |
| `nodeOutput` | `{ type: "nodeOutput", value: "text", nodeId: "n1" }` | Output de node anterior | `execution.NodeOutputs["n1"]["text"]` |

---

## Matriz de Testes por Node

### 1. core/start

Ponto de entrada — sempre sucede, sem config.

| # | Teste | Config | Esperado |
|---|-------|--------|----------|
| 1.1 | Emite "out" | `{}` | OutputHandles: `["out"]` |
| 1.2 | Contexto vazio | `{}` | Sem erro, sem NodeState |
| 1.3 | Sem efeitos colaterais | `{}` | Sem StatePatch, sem NodeOutput |

---

### 2. core/end

Node terminal — opcionalmente falha a execucao.

**Campos de config:**
- `terminateWithError` (bool)
- `errorCode` (string)
- `errorMessage` (FieldValue)

| # | Teste | Config | Esperado |
|---|-------|--------|----------|
| 2.1 | Fim normal | `terminateWithError: false` | OutputHandles: `[]`, status: completed |
| 2.2 | Fim com erro, mensagem literal | `terminateWithError: true, errorCode: "E001", errorMessage: {literal: "falhou"}` | ExecutionError code="E001", message="falhou" |
| 2.3 | Erro, mensagem do state | `errorMessage: {state: "errorMsg"}` | Mensagem resolvida do state |
| 2.4 | Erro, mensagem do event | `errorMessage: {event: "data.error"}` | Mensagem resolvida do event |
| 2.5 | Erro, mensagem do input | `errorMessage: {input: "reason"}` | Mensagem resolvida dos inputs |
| 2.6 | Erro, mensagem do nodeOutput | `errorMessage: {nodeOutput: "result", nodeId: "n1"}` | Mensagem resolvida do output |
| 2.7 | Erro, falha na resolucao | `errorMessage: {state: "naoexiste"}` | Fallback para errorCode como mensagem |
| 2.8 | Erro tem metadados | Qualquer config de erro | ErrorInfo tem nodeId, nodeType, timestamp |

---

### 3. core/set_state

Modifica o estado do workflow com 5 operacoes.

**Campos de config:**
- `operation` (string): `set | increment | decrement | append | remove`
- `targetField` (string)
- `valueSource` (FieldValue) — nao usado para `remove`

#### Operacao: set
| # | Teste | Source | StatePatch esperado |
|---|-------|--------|---------------------|
| 3.1 | Set de literal | `{literal: "hello"}` | `{targetField: "hello"}` |
| 3.2 | Set de state | `{state: "campoExistente"}` | `{targetField: <valor do state>}` |
| 3.3 | Set de event | `{event: "data.value"}` | `{targetField: <valor do event>}` |
| 3.4 | Set de input | `{input: "sensorId"}` | `{targetField: <valor do input>}` |
| 3.5 | Set de nodeOutput | `{nodeOutput: "result", nodeId: "n1"}` | `{targetField: <output do node>}` |

#### Operacao: increment
| # | Teste | Estado inicial | Valor | Esperado |
|---|-------|---------------|-------|----------|
| 3.6 | Incrementar numerico | `counter: 5` | `{literal: "3"}` | `counter: 8` |
| 3.7 | Incrementar nao-numerico | `counter: "abc"` | `{literal: "1"}` | Coerce para 0 + 1 = 1 |

#### Operacao: decrement
| # | Teste | Estado inicial | Valor | Esperado |
|---|-------|---------------|-------|----------|
| 3.8 | Decrementar numerico | `counter: 10` | `{literal: "3"}` | `counter: 7` |
| 3.9 | Decrementar nao-numerico | `counter: "abc"` | `{literal: "1"}` | Coerce para 0 - 1 = -1 |

#### Operacao: append
| # | Teste | Estado inicial | Valor | Esperado |
|---|-------|---------------|-------|----------|
| 3.10 | Append em array | `items: [1,2]` | `{literal: "3"}` | `items: [1,2,3]` |
| 3.11 | Append em campo inexistente | `{}` | `{literal: "1"}` | `items: [1]` (cria array) |
| 3.12 | Append em nao-array | `items: "str"` | `{literal: "1"}` | Erro ou coerce |

#### Operacao: remove
| # | Teste | Estado inicial | Esperado |
|---|-------|---------------|----------|
| 3.13 | Remover campo existente | `counter: 5` | StatePatch `{counter: nil}` (delecao) |
| 3.14 | Remover campo inexistente | `{}` | No-op, sem erro |

#### Erros
| # | Teste | Esperado |
|---|-------|----------|
| 3.15 | Operacao invalida | Erro retornado |
| 3.16 | Falha na resolucao do valor | Erro retornado |

---

### 4. core/condition

Avalia grupo de condicoes → roteia true/false.

**Campos de config:**
- `condition` (ConditionGroup) — arvore de condicoes AND/OR
- Cada item tem `field` (FieldValue) e `value` (FieldValue)

#### Logica
| # | Teste | Condicao | Esperado |
|---|-------|----------|----------|
| 4.1 | Resultado true | Avalia true | OutputHandles: `["true"]` |
| 4.2 | Resultado false | Avalia false | OutputHandles: `["false"]` |
| 4.3 | AND tudo true | `A && B` ambos true | `["true"]` |
| 4.4 | AND um false | `A && B` um false | `["false"]` |
| 4.5 | OR um true | `A \|\| B` um true | `["true"]` |
| 4.6 | OR tudo false | `A \|\| B` ambos false | `["false"]` |
| 4.7 | Grupos aninhados | `(A && B) \|\| C` | Avaliacao correta |

#### Tipos de source nos campos de condicao
| # | Teste | Source do campo | Esperado |
|---|-------|----------------|----------|
| 4.8 | Campo do state | `{state: "temp"}` | Resolve do state |
| 4.9 | Campo do event | `{event: "data.value"}` | Resolve do event |
| 4.10 | Campo do input | `{input: "threshold"}` | Resolve dos inputs |
| 4.11 | Campo do nodeOutput | `{nodeOutput: "result", nodeId: "n1"}` | Resolve do output |
| 4.12 | Campo literal | `{literal: "42"}` | Valor direto |

#### Erros
| # | Teste | Esperado |
|---|-------|----------|
| 4.13 | Campo ausente | Erro retornado |

---

### 5. core/log

Cria uma entrada de log com interpolacao de template.

**Campos de config:**
- `message` (string) — suporta tokens `${state.campo}` e `${event.campo}`
- `level` (string): `debug | info | warn | error`

| # | Teste | Mensagem | Esperado |
|---|-------|----------|----------|
| 5.1 | Mensagem estatica | `"hello world"` | LogEntry message="hello world" |
| 5.2 | Token de state | `"count: ${state.counter}"` | Resolve counter do state |
| 5.3 | Token de event | `"tipo: ${event.eventType}"` | Resolve do event |
| 5.4 | Multiplos tokens | `"${state.a} e ${event.b}"` | Ambos resolvidos |
| 5.5 | Campo ausente | `"${state.naoexiste}"` | Token nao substituido |
| 5.6 | Todos os niveis | Level: debug/info/warn/error | LogEntry com nivel correto |
| 5.7 | Mensagem vazia | `""` | Sem erro, log vazio |
| 5.8 | Sem sintaxe de token | `"texto puro"` | Passado sem alteracao |

---

### 6. core/goto

Portal — sender teleporta para receiver via edges injetadas no grafo.

**Campos de config:**
- `role` (string): `sender | receiver`
- `pairLabel` (string): vincula sender ao receiver
- `pairColor` (string): apenas UI

| # | Teste | Setup | Esperado |
|---|-------|-------|----------|
| 6.1 | Sender com receiver correspondente | Ambos com `pairLabel: "X"` | Sender emite `["out"]`, grafo roteia ao receiver |
| 6.2 | Sender sem receiver | Sender `pairLabel: "X"`, sem receiver | Erro: `GOTO_NO_RECEIVER` |
| 6.3 | Receiver emite "out" | Node receiver | Passthrough, `["out"]` |
| 6.4 | PairLabel errado | Sender "A", receiver "B" | Erro: sem match |
| 6.5 | Multiplos senders um receiver | 2 senders "X", 1 receiver "X" | Ambos roteiam ao mesmo receiver |
| 6.6 | Metadados do erro | Qualquer erro | NodeID e NodeType corretos |

---

### 7. core/switch

Roteamento multi-branch baseado em condicoes de caso.

**Campos de config:**
- `cases` ([]SwitchCase): cada um tem `id`, `condition` (ConditionGroup)
- `matchMode` (string): `first | all`

| # | Teste | Modo | Matches | Esperado |
|---|-------|------|---------|----------|
| 7.1 | Primeiro match, modo "first" | first | Caso A | `["case_A"]` |
| 7.2 | Segundo match, modo "first" | first | Caso B (nao A) | `["case_B"]` |
| 7.3 | Todos matches, modo "all" | all | A, C | `["case_A", "case_C"]` |
| 7.4 | Sem match | first | Nenhum | `["default"]` |
| 7.5 | Multiplos matches, "first" | first | A, B | Para em A: `["case_A"]` |
| 7.6 | Multiplos matches, "all" | all | A, B | Retorna ambos: `["case_A", "case_B"]` |
| 7.7 | Casos vazios | — | — | `["default"]` |
| 7.8 | Erro na condicao do caso | — | Erro | Erro propagado |
| 7.9 | Grupos aninhados nos casos | — | AND/OR aninhados | Avaliacao correta |

---

### 8. core/delay

Suspende execucao por uma duracao — retoma via timer.

**Campos de config:**
- `duration` (int)
- `unit` (string): `seconds | s | minutes | m | hours | h | days | d`

| # | Teste | Duracao | Unidade | Esperado |
|---|-------|---------|---------|----------|
| 8.1 | 5 segundos | 5 | seconds | NodeState: waitType="timer", expiresAt=now+5s |
| 8.2 | 1 minuto | 1 | minutes | expiresAt=now+1m |
| 8.3 | 2 horas | 2 | hours | expiresAt=now+2h |
| 8.4 | 3 dias | 3 | days | expiresAt=now+3d |
| 8.5 | Abreviacao "s" | 10 | s | Mesmo que "seconds" |
| 8.6 | Abreviacao "m" | 5 | m | Mesmo que "minutes" |
| 8.7 | Duracao zero | 0 | s | Erro |
| 8.8 | Duracao negativa | -1 | s | Erro |
| 8.9 | Unidade desconhecida | 5 | weeks | Erro |
| 8.10 | expiresAt correto | 30 | s | Verificar timestamp |

---

### 9. core/code

Suspende para executar JavaScript externamente — retoma via callback.

**Campos de config:**
- `script` (string)
- `timeout` (int, ms)

| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 9.1 | Suspende com script | Script simples | NodeState: waitType="callback", script=<codigo> |
| 9.2 | Callback de sucesso | Resume output={text:"ok"} | NodeOutputs setado, segue edge "success" |
| 9.3 | Callback de erro | Resume error={code,message} | Segue edge "error" ou falha execucao |
| 9.4 | Timeout padrao | Sem config de timeout | expiresAt=now+30s |
| 9.5 | Timeout customizado | timeout=10000 | expiresAt calculado do timeout |
| 9.6 | EnableOutput=true no timeout | Config de timeout | enableOutput=true no NodeState |
| 9.7 | Config ausente | nil config | Erro retornado |
| 9.8 | Campos do NodeState completos | Qualquer | Tem waitType, script, timeout, expiresAt, enableOutput |

---

### 10. core/wait_signal

Suspende esperando entrega de sinal externo.

**Campos de config:**
- `signalName` (string)
- `mappings` ([]SignalMapping): cada um tem `paramName`, `value` (FieldValue)

| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 10.1 | Suspensao basica | signalName="approval" | NodeState: waitType="signal", signalName="approval" |
| 10.2 | Resume com dados | Sinal entrega dados | Execucao continua |
| 10.3 | Mapping do state | `{state: "userId"}` | Param resolvido do state |
| 10.4 | Mapping do event | `{event: "data.id"}` | Param resolvido do event |
| 10.5 | Mapping do input | `{input: "target"}` | Param resolvido dos inputs |
| 10.6 | Mapping do nodeOutput | `{nodeOutput: "id", nodeId: "n1"}` | Param resolvido do output |
| 10.7 | Timeout padrao | Sem config de timeout | expiresAt=now+24h |
| 10.8 | Timeout customizado + EnableOutput | timeout setado | expiresAt e enableOutput corretos |
| 10.9 | signalName ausente | signalName="" | Erro |
| 10.10 | Config ausente | nil | Erro |

---

### 11. core/subworkflow

Dispara um workflow filho — retoma quando filho completa.

**Campos de config:**
- `workflowId` (string)
- `inputMappings` ([]InputMapping): `childParamName`, `value` (FieldValue)
- `outputMappings` ([]OutputMapping): `outputName`, `stateField`

| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 11.1 | Subworkflow basico | workflowId valido | NodeState: waitType="callback", workflowId setado |
| 11.2 | Input do state | `{state: "data"}` | inputData resolvido |
| 11.3 | Input do event | `{event: "payload"}` | inputData resolvido |
| 11.4 | Input do input | `{input: "param"}` | inputData resolvido |
| 11.5 | Input do nodeOutput | `{nodeOutput: "x", nodeId: "n1"}` | inputData resolvido |
| 11.6 | Output mappings no resume | Resume com output → state | StatePatch aplicado |
| 11.7 | Profundidade 0 (root) | depth=0 | Permitido, armazena depth+1 |
| 11.8 | Profundidade = MaxDepth | depth=10 | Erro: recursao maxima |
| 11.9 | Timeout customizado | timeout setado | expiresAt correto |
| 11.10 | workflowId ausente | workflowId="" | Erro |
| 11.11 | Falha na resolucao de input | Source invalido | Erro |
| 11.12 | Profundidade incrementa | depth=3 | NodeState depth=4 |

---

### 12. core/trigger_event

Publica um evento na plataforma — retoma via callback.

**Campos de config:**
- `eventType` (string)
- `payloadMapping` ([]TriggerPayloadField): `key`, `value` (FieldValue)

| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 12.1 | Evento basico | eventType="alert" | NodeState: waitType="callback", eventType setado |
| 12.2 | Payload do state | `{state: "data"}` | Payload resolvido |
| 12.3 | Payload do event | `{event: "source"}` | Payload resolvido |
| 12.4 | Payload do input | `{input: "config"}` | Payload resolvido |
| 12.5 | Payload do nodeOutput | `{nodeOutput: "x", nodeId: "n1"}` | Payload resolvido |
| 12.6 | Multiplos campos de payload | 3 campos mistos | Todos resolvidos |
| 12.7 | Timeout padrao | Sem config | expiresAt=now+30s |
| 12.8 | eventType ausente | eventType="" | Erro |
| 12.9 | Config ausente | nil | Erro |
| 12.10 | Falha na resolucao do payload | Source invalido | Erro |

---

### 13. Plugin (executor generico marketplace)

Resolve manifest, credentials, templates — despacha para servico Triggers.

**Campos de config:**
- `operation` (string)
- `credentialId` (string, opcional)
- `rawConfig` (map) — cada valor pode ser um FieldValue

**Contextos de template:** `{{credentials.*}}`, `{{config.*}}`, `{{wf.state.*}}`, `{{wf.input.*}}`, `{{event.*}}`, `{{manifest.defaults.*}}`

| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 13.1 | Operacao basica | Manifest + operacao validos | NodeState com action resolvida |
| 13.2 | Config do state | `{state: "value"}` | Resolvido no rawConfig |
| 13.3 | Config do event | `{event: "data"}` | Resolvido |
| 13.4 | Config do input | `{input: "param"}` | Resolvido |
| 13.5 | Config do nodeOutput | `{nodeOutput: "x", nodeId: "n1"}` | Resolvido |
| 13.6 | Template de credential | `{{credentials.apiKey}}` | Template resolvido com credential descriptografada |
| 13.7 | Template de config | `{{config.chatId}}` | Template resolvido |
| 13.8 | Template de state | `{{wf.state.counter}}` | Template resolvido |
| 13.9 | Template de event | `{{event.data.id}}` | Template resolvido |
| 13.10 | Defaults do manifest | `{{manifest.defaults.baseUrl}}` | Template resolvido |
| 13.11 | Hooks resolvidos | hooks before/after | Templates nos hooks resolvidos |
| 13.12 | Sem credential | credentialId="" | Credentials vazias, sem erro |
| 13.13 | fetchOptions como literal | `{fetchOptions: "123"}` | Tratado como literal "123" |
| 13.14 | loadOptions como literal | `{loadOptions: "abc"}` | Tratado como literal "abc" |
| 13.15 | Manifest ausente | pluginId desconhecido | Erro |
| 13.16 | Operacao ausente | Nome de operacao invalido | Erro |
| 13.17 | nodeType malformado | Sem "/" no tipo | Erro |
| 13.18 | Falha na descriptografia | credentialId invalido | Erro |
| 13.19 | Timeout padrao | Sem config | expiresAt=now+30s |

---

### 14. core/fanout

Fork de branches paralelas.

**Campos de config:**
- `branches` (int)
- `mode` (string): `waitAll | firstCompleted`

| # | Teste | Branches | Mode | Esperado |
|---|-------|----------|------|----------|
| 14.1 | 2 branches | 2 | — | `["out_1", "out_2"]` |
| 14.2 | 3 branches | 3 | — | `["out_1", "out_2", "out_3"]` |
| 14.3 | 1 branch | 1 | — | `["out_1"]` |
| 14.4 | Max branches | 20 | — | Permitido |
| 14.5 | Max+1 | 21 | — | Erro |
| 14.6 | 0 branches | 0 | — | Erro |
| 14.7 | Negativo | -1 | — | Erro |
| 14.8 | Modo waitAll | 2 | waitAll | Todas branches devem completar |
| 14.9 | Modo firstCompleted | 2 | firstCompleted | Primeira branch completa a execucao |
| 14.10 | Nomeacao dos handles | N | — | Sequencial: out_1, out_2, ..., out_N |

---

### 15. core/merge

Une branches paralelas de volta a caminho unico.

**Campos de config:**
- `branches` (int): contagem esperada
- `strategy` (string): `all | any | first`

| # | Teste | Branches | Strategy | Chegadas | Esperado |
|---|-------|----------|----------|----------|----------|
| 15.1 | Espera todos, 2 branches, 1a chegada | 2 | all | 1 | `[]` (espera) |
| 15.2 | Espera todos, 2 branches, 2a chegada | 2 | all | 2 | `["out"]` |
| 15.3 | Any, 1a chegada | 2 | any | 1 | `["out"]` |
| 15.4 | First, 1a chegada | 2 | first | 1 | `["out"]` |
| 15.5 | 3 branches all | 3 | all | 3 | `["out"]` |
| 15.6 | Strategy padrao | — | (vazio) | — | Usa "all" |
| 15.7 | 0 branches | 0 | — | — | Erro |
| 15.8 | Sem state existente | — | — | — | Inicializa branchCount=0 |
| 15.9 | Verificacao de incremento | — | — | 2 chamadas | branchCount: 0→1→2 |
| 15.10 | State persiste | — | — | — | branchCount no NodeState |

---

### 16. core/sequence

Execucao sequencial de steps.

**Campos de config:**
- `steps` (int)

| # | Teste | Steps | Chamada # | Esperado |
|---|-------|-------|-----------|----------|
| 16.1 | 3 steps, chamada 1 | 3 | 1 | `["step_1"]` |
| 16.2 | 3 steps, chamada 2 | 3 | 2 | `["step_2"]` |
| 16.3 | 3 steps, chamada 3 | 3 | 3 | `["step_3"]` |
| 16.4 | 3 steps, chamada 4 | 3 | 4 | `["done"]` |
| 16.5 | 1 step | 1 | 1,2 | `["step_1"]` depois `["done"]` |
| 16.6 | 0 steps | 0 | 1 | `["done"]` imediatamente |
| 16.7 | currentStep incrementa | 3 | — | 0→1→2→3 |
| 16.8 | Numeracao dos steps | — | — | Comeca em 1, nao 0 |

---

### 17. core/loop

Itera sobre um array, expondo o item atual.

**Campos de config:**
- `source` (FieldValue): deve resolver para array

**Saidas por iteracao:**
- `StatePatch`: `loop_item` (elemento atual), `loop_index` (posicao 0-based)
- `NodeOutput`: `item` (elemento atual), `index` (posicao 0-based)

#### Comportamento de iteracao
| # | Teste | Array | Iteracao | Handles esperados | loop_item | loop_index |
|---|-------|-------|----------|-------------------|-----------|------------|
| 17.1 | 3 itens, iter 0 | [a,b,c] | 0 | `["body"]` | a | 0 |
| 17.2 | 3 itens, iter 1 | [a,b,c] | 1 | `["body"]` | b | 1 |
| 17.3 | 3 itens, iter 2 | [a,b,c] | 2 | `["body"]` | c | 2 |
| 17.4 | 3 itens, iter 3 | [a,b,c] | 3 | `["done"]` | — | — |
| 17.5 | NodeOutput item | [a,b,c] | 0 | — | NodeOutput.item = a | NodeOutput.index = 0 |
| 17.6 | NodeOutput index | [a,b,c] | 1 | — | NodeOutput.item = b | NodeOutput.index = 1 |

#### Tipos de source
| # | Teste | Source | Esperado |
|---|-------|--------|----------|
| 17.7 | Source do state | `{state: "items"}` | Resolve do state |
| 17.8 | Source do event | `{event: "data.list"}` | Resolve do event |
| 17.9 | Source do input | `{input: "targets"}` | Resolve dos inputs |
| 17.10 | Source do nodeOutput | `{nodeOutput: "list", nodeId: "n1"}` | Resolve do output |

#### Casos limites
| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 17.11 | Array vazio | `[]` | `["done"]` imediatamente |
| 17.12 | Item unico | `[x]` | `["body"]` depois `["done"]` |
| 17.13 | Max itens | 10000 itens | Permitido |
| 17.14 | Max+1 itens | 10001 itens | Erro: max iteracoes |
| 17.15 | Source nao-array | `"string"` | Erro: deve ser array |
| 17.16 | Source irresolvivel | Ref invalida | Erro |
| 17.17 | currentIndex persiste | — | Incrementado no NodeState |
| 17.18 | Itens objeto | `[{a:1},{a:2}]` | loop_item eh objeto, nodeOutput.item.a funciona |

---

### 18. core/wait_for

Avalia condicao — procede imediatamente se true, suspende se false.

**Campos de config:**
- `field` (string): campo do state a verificar
- `operator` (string): operador de comparacao
- `compareTo` (FieldValue): valor para comparar

#### Match imediato
| # | Teste | State | Condicao | Esperado |
|---|-------|-------|----------|----------|
| 18.1 | True imediatamente | counter=5 | counter == 5 | `["matched"]` |
| 18.2 | False → suspende | counter=3 | counter == 5 | NodeState: waitType="condition" |

#### Operadores
| # | Teste | Operador | State | CompareTo | Esperado |
|---|-------|----------|-------|-----------|----------|
| 18.3 | Igual | `==` | val=5 | 5 | matched |
| 18.4 | Diferente | `!=` | val=5 | 3 | matched |
| 18.5 | Maior que | `>` | val=10 | 5 | matched |
| 18.6 | Menor que | `<` | val=3 | 5 | matched |
| 18.7 | Maior ou igual | `>=` | val=5 | 5 | matched |
| 18.8 | Menor ou igual | `<=` | val=5 | 5 | matched |

#### Sources do compareTo
| # | Teste | Source | Esperado |
|---|-------|--------|----------|
| 18.9 | CompareTo do state | `{state: "threshold"}` | Resolvido |
| 18.10 | CompareTo do event | `{event: "limit"}` | Resolvido |
| 18.11 | CompareTo do input | `{input: "max"}` | Resolvido |
| 18.12 | CompareTo do nodeOutput | `{nodeOutput: "val", nodeId: "n1"}` | Resolvido |

#### Casos limites
| # | Teste | Cenario | Esperado |
|---|-------|---------|----------|
| 18.13 | Timeout customizado | timeout setado | expiresAt correto |
| 18.14 | EnableOutput=true | timeout + enableOutput | enableOutput no NodeState |
| 18.15 | Falha na avaliacao | Operador invalido | Erro |

---

## Testes de Integracao

Cenarios cross-node que testam o DAG walker de ponta a ponta.

### Fluxos sincronos
| # | Teste | Fluxo | Verificacao |
|---|-------|-------|-------------|
| I.1 | Fanout + Merge (sync) | start → fanout(2) → [setState A, setState B] → merge → end | State tem A e B |
| I.2 | Loop + setState | start → loop([1,2,3]) → body(setState inc "counter") → end | counter=3 |
| I.3 | Roteamento por condicao | start → condition → true: setState "path"="T" → end, false: setState "path"="F" → end | Caminho correto tomado |
| I.4 | Switch 3 casos | start → switch → [case_a, case_b, default] → end | Caso correto roteado |
| I.5 | Sequence 3 steps | start → sequence(3) → [step_1→log, step_2→log, step_3→log] → done → end | 3 logs criados |
| I.6 | Teleporte goto | start → goto_sender("X") →→ goto_receiver("X") → setState → end | State setado apos teleporte |

### Fluxos async
| # | Teste | Fluxo | Verificacao |
|---|-------|-------|-------------|
| I.7 | Code caminho sucesso | start → code → [success: end_ok, error: end_err] | Segue edge success |
| I.8 | Code caminho erro | start → code → [success: end_ok, error: end_err] | Segue edge error |
| I.9 | Fanout + branches async | start → fanout(2) → [code_1, code_2] → merge → end | 2 callbacks, merge em 2 |
| I.10 | Fanout firstCompleted | start → fanout(2,firstCompleted) → [code_1, code_2] → merge → end | 1 callback, outra cancelada |
| I.11 | Loop body async | start → loop([1,2,3]) → body(code) → end | 3 callbacks, loop completa |
| I.12 | Loop async + nodeOutput | start → code(return list) → loop(source=nodeOutput) → body(code) → end | Source resolvido do nodeOutput |

### Combinacoes complexas
| # | Teste | Fluxo | Verificacao |
|---|-------|-------|-------------|
| I.13 | Goto + Code | start → fanout → [goto_sender→receiver→code→end, goto_sender→receiver→end] | Code executa apos teleporte |
| I.14 | Fanout + WaitSignal + Code | start → fanout(2) → [code, wait_signal] → end | 2 nodes async, resume parcial |
| I.15 | Loops aninhados | start → loop([1,2]) → body(loop([a,b]) → body(setState) → end) → end | 4 iteracoes totais (2x2) |
| I.16 | Loop + Fanout no body | start → loop([1,2]) → body(fanout(2) → [setState, setState] → merge) → end | 2 iteracoes x 2 branches |
| I.17 | End com erro | start → code → erro callback → end(terminateWithError) | status=failed, errorInfo setado |
| I.18 | Max loop excedido | start → loop(10001 itens) | Erro: max iteracoes |
| I.19 | Subworkflow no fanout | start → fanout(2) → [subworkflow, code] → merge → end | Ambos async, merge espera |
| I.20 | Workflow complexo completo | start → fanout(3) → [goto→code→end, loop(3)→plugin→end, wait_signal→end] → end | Todas features combinadas |

---

## Total: ~298 casos de teste

| Grupo | Contagem |
|-------|----------|
| core/start | 3 |
| core/end | 8 |
| core/set_state | 16 |
| core/condition | 13 |
| core/log | 8 |
| core/goto | 6 |
| core/switch | 9 |
| core/delay | 10 |
| core/code | 8 |
| core/wait_signal | 10 |
| core/subworkflow | 12 |
| core/trigger_event | 10 |
| plugin (generico) | 19 |
| core/fanout | 10 |
| core/merge | 10 |
| core/sequence | 8 |
| core/loop | 18 |
| core/wait_for | 15 |
| Integracao | 20 |
| **Total** | **~203 unitarios + ~95 integracao = 298** |

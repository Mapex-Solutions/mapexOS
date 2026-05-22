# Plugin Execution — Decisoes Arquiteturais

> Data: 2026-03-23 | Status: Fechado

Documento com todas as decisoes tomadas no debate sobre como o runtime executa nodes de plugins (marketplace). Cada ponto descreve o problema, o que foi debatido, a decisao final e o ganho.

---

## Ponto 1 — Routing de executor

**Problema:** O runtime hoje tem 17 executores core registrados num map. Quando chega um node como `telegram/message`, ele tenta buscar no map, nao encontra e retorna erro `EXECUTOR_NOT_FOUND`. O workflow morre.

**Debate:** Discutimos se deveriamos usar um fallback no registry (tentar buscar, se nao achar usar fallback) ou verificar por prefixo. O fallback gasta CPU tentando buscar algo que nunca vai existir no map.

**Decisao:** Ter um map de nodes core. Se o nodeType esta no map, usa o executor core. Se nao esta, e um plugin node e vai direto para o PluginExecutor. Sem try/catch, sem fallback. Verificacao por prefixo `core/`.

**Ganho:** Zero CPU desperdicada em lookup falho. Roteamento deterministico e rapido.

---

## Ponto 2 — Config parsing para plugin nodes

**Problema:** O `parseNodeConfig` tem um switch/case com todos os 17 node types core. O `default` retorna `nil`, o que significa que plugin nodes perdem toda a informacao de config.

**Debate:** O que extrair do config de um plugin node? Precisamos de informacoes minimas para o executor montar o payload.

**Decisao:** O default case extrai `operation`, `credentialId` e preserva o `rawConfig` inteiro. Nao faz parsing pesado — quem resolve templates e FieldValues e o executor em runtime.

**Ganho:** O graph builder cacheia o config parseado uma vez. O executor recebe tudo que precisa sem reprocessar.

---

## Ponto 3 — Responsabilidade do PluginExecutor

**Problema:** Qual a responsabilidade do executor generico de plugins? Ele so marca o node como "waiting" e suspende? Ou ele tambem resolve credentials, monta o payload e prepara tudo?

**Debate:** Se o executor so suspende, alguem precisa resolver depois (lifecycle? domain service?). Se ele resolve tudo, a responsabilidade fica clara num unico lugar.

**Decisao:** O PluginExecutor faz tudo antes de suspender:
1. Decripta a credential
2. Carrega o manifest do plugin
3. Encontra a operation pelo nome
4. Resolve FieldValues do config (state, event → valores finais)
5. Substitui templates `{{credentials.X}}` e `{{config.X}}`
6. Monta o payload completo
7. Suspende com tudo pronto no NodeState

**Ganho:** Responsabilidade unica. O Triggers Service recebe um payload pronto e so executa. Nenhum servico downstream precisa saber sobre credentials, state ou workflow.

---

## Ponto 4 — Credential resolution

**Problema:** O modulo de credentials vive no Workflow Service. O Triggers Service nao conhece credentialId nem sabe decriptar. Quem decripta e quando?

**Debate:** Inicialmente se pensou em enviar o credentialId para o Triggers resolver. Mas o Triggers nao tem acesso ao modulo de credentials nem a master key de encriptacao.

**Decisao:** O Workflow Service decripta a credential antes do dispatch e envia os valores plain para o Triggers. O Triggers nunca ve credentialId, nunca decripta nada.

**Ganho:** Separacao clara de dominio. O Triggers Service e um executor burro — recebe dados prontos e executa. A seguranca da encriptacao fica isolada no Workflow Service.

---

## Ponto 5 — FieldValue resolution

**Problema:** No config de um plugin node, campos como `chatId` podem ter valores dinamicos: `{ type: "state", value: "myChatVar" }`. Isso significa "pega o valor atual da variavel de estado myChatVar". Alguem precisa transformar isso no valor real (ex: `"5847825355"`) antes de enviar para o Triggers.

**Debate:** O ValueResolver do engine ja faz isso para nodes core como `set_state` e `condition`. A questao era se reutilizamos no PluginExecutor ou se mandamos cru para o Triggers.

**Decisao:** O Workflow Service resolve todos os FieldValues antes do dispatch. O Triggers recebe valores finais, nao referencias.

**Ganho:** Reutilizacao do ValueResolver existente. O Triggers nao precisa saber nada sobre workflow state, event payload ou node outputs.

---

## Ponto 6 — Template resolution

**Problema:** Os templates `{{credentials.botToken}}` e `{{config.chatId}}` no manifest precisam ser substituidos pelos valores reais. Mas templates como `{{before.token}}` (resultado do hook before) so podem ser resolvidos depois do hook executar.

**Debate:** Se tudo fosse resolvido no Workflow, ele precisaria executar os hooks (HTTP calls) — mas a regra e "Workflow Service nunca faz HTTP em runtime". Se tudo fosse no Triggers, ele precisaria saber sobre credentials e state.

**Decisao:** Dois tipos de template, dois resolvedores:
- `{{credentials.X}}` e `{{config.X}}` → Workflow resolve antes do dispatch
- `{{before.X}}` → Triggers resolve em runtime, depois de executar o hook before

**Ganho:** Cada dominio resolve o que e seu. O Workflow nao faz HTTP. O Triggers nao conhece credentials. Templates internos da pipeline (`{{before.X}}`) ficam isolados no Triggers.

---

## Ponto 7 — Montagem do HTTP request

**Problema:** Quem monta a URL final, method, headers, body prontos para execucao?

**Debate:** Inicialmente se pensou em o Workflow montar o request HTTP completo e o Triggers so disparar. Mas com hooks (before → login → token → usa no request), o Triggers precisa resolver `{{before.token}}` no request principal.

**Decisao:** O Workflow monta tudo que pode — resolve credentials, config, baseUrl + path. O Triggers recebe a pipeline (hooks + operation) com templates `{{before.X}}` pendentes e resolve esses internamente.

**Ganho:** O Triggers vira um executor de pipeline generico. Ele nao sabe nada sobre plugins, manifests ou workflows. Recebe uma sequencia de acoes (before → operation → after) e executa.

---

## Ponto 8 — Dispatch NATS

**Problema:** Qual stream/subject usar para enviar a execucao de plugin para o Triggers?

**Debate:** Poderia ser um stream especifico novo ou reutilizar o que ja existe.

**Decisao:** Vamos UM NOVO STREAM pois o stream que temos hoje o router pode usar direto. (DEBATER ANTES DE IMPLEMENTAR)

**Ganho:** Consistencia com o padrao existente de dispatch async.

---

## Ponto 9 — Consumer no Triggers Service

**Problema:** O Triggers Service hoje recebe `TriggerExecuteEvent` com um triggerId e busca config no banco. Para plugin nodes, nao existe trigger cadastrado.

**Debate:** Criar um consumer novo.

**Decisao:** O Triggers recebe um payload completo e pronto (pipeline com hooks + operation). Nao precisa buscar config no banco. Executa a pipeline sequencialmente: before → resolve `{{before.X}}` → operation → after.

**Ganho:** O Triggers nao precisa conhecer o conceito de "plugin". Ele recebe acoes prontas e executa. Qualquer servico pode usar esse mesmo consumer para executar pipelines HTTP/MQTT/NATS.

---

## Ponto 10 — Resume (callback)

**Problema:** Depois que o Triggers executa, como o resultado volta para o workflow? Como mapeia para os handles `success` ou `error` do node?

**Debate:** O mecanismo de resume ja existe para `core/code` e `core/subworkflow`. A questao era como mapear o resultado.

**Decisao:** O Triggers publica no `WORKFLOW-RESUME` com status `success` ou `error`. O runtime recebe, aplica o output no node e segue pelo handle correspondente — `success` ou `error` — que esta ligado aos proximos nodes no DAG.

**Ganho:** Reutilizacao do pipeline de resume existente. Sem logica nova no runtime para plugin nodes.

---

## Ponto 11 — Auth types futuros (OAuth2, login-based)

**Problema:** Como suportar diferentes tipos de autenticacao (API Key, Bearer Token, Basic Auth, OAuth2, login com sessao) sem criar logica especifica para cada tipo?

**Debate:** Inicialmente se pensou em ter um campo `authType` na credential com logica especifica no Workflow para cada tipo. Mas isso acoplaria o Workflow a cada tipo de auth.

**Decisao:** O sistema de hooks resolve isso genericamente. O hook `before` pode fazer qualquer coisa — login, refresh de token, pre-autenticacao. O output do before (ex: token de sessao) fica disponivel como `{{before.token}}` no request principal. Nao precisa de logica especifica por auth type.

**Ganho:** Zero codigo especifico por tipo de auth. Um plugin que precisa de login define um hook before com a chamada de auth. O Telegram que usa token no path nao precisa de hook. Tudo generico via manifest.

---

## Ponto 12 — Dependencias cross-module

**Problema:** O PluginExecutor no modulo runtime precisa acessar funcionalidades de outros modulos: decriptar credentials, buscar manifest, resolver FieldValues.

**Debate:** Como fazer isso sem acoplar o runtime diretamente aos outros modulos?

**Decisao:** Usar ports (interfaces) seguindo a arquitetura hexagonal existente:
- `CredentialServicePort` → metodo `DecryptCredential(ctx, id)` do modulo credentials
- `PluginManifestRepository` → metodo `FindByPluginId(ctx, pluginId)` do modulo plugins
- `ValueResolverPort` → metodo `Resolve(fieldValue, state, event, inputs, nodeOutputs)` do modulo engine

O PluginExecutor recebe essas interfaces via DI. Nunca importa implementacoes concretas.

**Ganho:** Desacoplamento total. O runtime depende de interfaces, nao de modulos concretos. Qualquer modulo pode ser substituido sem alterar o executor.

---

## Decisoes Complementares

### DSL — Manifest Design Decisions

| Decision | Approach | Reason |
|----------|----------|--------|
| Dispatch | Per operation (`type: "http"`) | A plugin can have HTTP + MQTT operations |
| BaseUrl | Available via `{{manifest.defaults.baseUrl}}` template | Each operation can use different URLs |
| fetchOptions | Unified Action contract | Same structure for operations, fetchOptions, hooks |
| Properties | `rendering` + `fetchOptions.rules` separated | Separation of concerns (visual vs data) |
| availableOutputs | Descriptive name | Clear intent for node output hints |
| placeholder | Only in `rendering.placeholder` | Single location, no ambiguity |

### Node Timeout — Async execution safety

Every async node (plugin or core) MUST have a timeout. The timeout lives at the **node level** in the definition, not inside the node config.

**Manifest DSL:**
```json
{
  "nodeTypes": [{
    "type": "telegram/message",
    "timeout": { "duration": 30, "unit": "seconds", "enableOutput": false },
    "properties": [...]
  }]
}
```

**Fields:**
| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `duration` | int | 30 | Timeout duration value |
| `unit` | string | "seconds" | Unit: seconds, minutes, hours, days, months, years |
| `enableOutput` | bool | false | When true, timeout routes to "timeout" output handle instead of failing |

**Resolution priority:**
1. Node instance (`node.timeout`) — user configured in editor
2. Plugin manifest default (`nodeType.timeout`) — plugin author default
3. Platform default — 30s for callback, 24h for signal/condition

**Behavior when timeout expires:**
- `enableOutput: false` → execution fails with `TIMEOUT_EXCEEDED`
- `enableOutput: true` → resumes via "timeout" output handle (dynamic, injected by UI)

**Implementation:** Each executor sets `expiresAt` in NodeState. The Reconciler sweeps `timerExpiresAt` from MongoDB and publishes resume with `isTimeout: true`.

### Script execution — Regra de seguranca

| Contexto | Runtime | Fonte do codigo |
|----------|---------|-----------------|
| `fetchOptions.output.transform` | goja/ES5 (Workflow Service) | Manifest (auditado) |
| `type: "script"` em operations/hooks | V8 isolated-vm (Triggers → JS Executor) | Manifest (auditado) |
| `core/code` node | V8 isolated-vm (JS Workflow Executor) | Usuario (nao auditado, isolado) |

O backend rejeita updates via API onde `type: "script"` — scripts so vem de manifests auditados.

### Body nao-JSON — Content-Type via headers

Sem campo novo no DSL. O Triggers le o `Content-Type` dos headers e serializa o body de acordo:

- Sem header ou `application/json` → JSON (default)
- `application/x-www-form-urlencoded` → form urlencoded
- `multipart/form-data` → multipart
- Qualquer outro → raw

O autor do plugin define o header que precisar. Zero logica especial no DSL.

---

### Credentials — Array with multiple auth methods

`credentials` is an array. Each item has an `id` that identifies the auth method. A plugin can support multiple auth methods (e.g., API Key + OAuth2).

```json
"credentials": [
  { "id": "apiKey", "name": "API Key", "fields": [...] },
  { "id": "oauth2", "name": "OAuth2", "fields": [...] }
]
```

The user chooses which method to use when creating a credential instance. The `credentialId` in the node config points to the created instance.

**Impacto:** Alterar schema Zod, contract Go, entity Go, UI de criacao de credential.

---

### Template Contexts — Reserved Keywords

Templates use `{{context.field}}` syntax. Each context is a reserved keyword that maps to a data source:

| Context | Source | Resolved by | Example |
|---------|--------|-------------|---------|
| `manifest` | Plugin manifest data (defaults, metadata) | Workflow Service | `{{manifest.defaults.baseUrl}}` |
| `credentials` | Decrypted credential data | Workflow Service | `{{credentials.botToken}}` |
| `wf` | Workflow execution data (states, inputs) | Workflow Service | `{{wf.state.counter}}`, `{{wf.input.textToSend}}` |
| `event` | Event payload that triggered the execution | Workflow Service | `{{event.data.temperature}}` |
| `config` | Current node config (form field values) | Workflow Service | `{{config.chatId}}` |
| `before` | Output from the before hook | Triggers Service | `{{before.token}}` |

Workflow Service resolves all contexts except `before` before dispatching to Triggers. The Triggers Service only resolves `{{before.X}}` after executing the before hook.

---

### Hooks — Lifecycle

| Hook | When it runs | Scope |
|------|-------------|-------|
| `before` | Before the operation executes | Global per nodeType (same for all operations) |
| `after` | After the operation completes successfully | Global per nodeType |
| `destroy` | When the workflow completes OR is cancelled | Global per nodeType |

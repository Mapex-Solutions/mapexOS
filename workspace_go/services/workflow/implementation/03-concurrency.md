# Concurrency — Idempotência + CAS Conflict Handling

Decisões: A2, E6

---

## Problema

NATS pode reentreagar uma mensagem se o worker crashar antes do ACK. Dois workers podem receber a mesma mensagem ao mesmo tempo. Branches de fanout podem completar simultaneamente e tentar atualizar o mesmo estado.

Sem proteção, um worker sobrescreve o trabalho do outro silenciosamente.

---

## Como Estamos Resolvendo

### TransitionId NÃO necessário (A2)

Duas proteções existentes cobrem o processamento duplicado sem TransitionId (side effects dependem de idempotência dos serviços downstream — trigger service, js-code service já são idempotentes por instanceId+nodeId):

**Proteção 1 — Estado da instância (detecta duplicada):**
```
Worker recebe callback duplicado:
  1. NATS KV Get("inst:123") → estado atual (~1ms)
  2. Checa: instance.status == "waiting"?
     → NÃO (status já é "completed", worker anterior já processou)
     → ACK e skip

  Esse Get JÁ seria feito de qualquer forma para carregar o state para executar.
  Ou seja: zero custo adicional para detectar duplicada.

  Funciona para todos os tipos:
  - Callback: status != waiting → skip
  - Signal: NodeStates não contém waitType="signal" ou signalName não bate → skip
  - Timer: status != waiting ou NodeStates[nodeId] sem expiresAt → skip
```

**Proteção 2 — NATS KV CAS (previne processamento simultâneo):**
```
2 workers recebem a MESMA mensagem ao mesmo tempo:
  Worker A: NATS KV Get → revision 5
  Worker B: NATS KV Get → revision 5

  Worker A: executa, Update(expectedRev=5) → SUCESSO → revision 6
  Worker B: executa, Update(expectedRev=5) → ERRO (revision agora é 6)

  Worker B detecta conflito → ACK e skip.
  Somente 1 dos 2 consegue. Zero duplicação.
```

**Nota sobre side effects:** Se Worker A publicou request pro trigger service e depois falhou no CAS, o trigger service recebeu um request "fantasma". Isso é seguro porque serviços downstream (trigger, js-code) são idempotentes por `instanceId+nodeId` — request duplicado é detectado e descartado.

### Contrato de idempotência para serviços downstream

O runtime depende de idempotência dos serviços que recebem requests ASYNC. Cada serviço downstream DEVE seguir este contrato:

```
Chave de idempotência: {instanceId}:{nodeId}
  → Composta, identifica univocamente uma execução de um node numa instância.
  → O runtime SEMPRE inclui instanceId e nodeId no payload do request.

TTL da chave: 24 horas
  → Requests duplicados dentro de 24h são detectados e descartados.
  → Após 24h, a chave expira (instância já completou/falhou há muito tempo).
  → Implementar via: TTL index no MongoDB ou KV store do serviço.

Comportamento em retry/duplicata:
  1. Serviço recebe request
  2. Busca chave {instanceId}:{nodeId} no store de idempotência
  3. Se chave NÃO existe:
     → Insere chave com status="processing"
     → Processa request normalmente
     → Atualiza status="done" + resultado
     → Publica resultado no callbackSubject
  4. Se chave existe E status=="done":
     → Retorna resultado cacheado no callbackSubject (re-publish)
     → NÃO re-executa
  5. Se chave existe E status=="processing":
     → Descarta (outro worker está processando)
     → NÃO publica callback (o worker original vai publicar)

Serviços que DEVEM implementar:
  □ js-code service (WORKFLOW-JS-CODE) — execução de scripts
  □ Trigger service (fila workflow) — dispatch de eventos
  □ Workflow service (subworkflow) — criação de child instance

Implementação recomendada:
  MongoDB collection: idempotency_keys
  {
    _id:          "{instanceId}:{nodeId}",
    status:       "processing" | "done",
    result:       <resultado cacheado>,
    callbackSubj: "workflow.resume.callback.{instanceId}",
    createdAt:    ISODate(),         // TTL index: expireAfterSeconds: 86400
  }
```

**Futuro:** Se métricas mostrarem alta taxa de duplicação (> 1%), implementar resumeToken/TransitionId como otimização. Composto: `{instanceId}:{nodeId}:{revision}:{attempt}` → determinístico, persistido na instância, validado no callback.

### CAS (Compare-And-Swap) — O que é (E6)

Operação atômica: "atualizar este valor, MAS só se ninguém mexeu desde a última vez que eu li."

```go
// Lê o estado atual
entry, _ := kv.Get("inst:123")             // retorna dados + revision 5

// Atualiza SÓ SE ninguém mexeu
rev, err := kv.Update("inst:123", newData, entry.Revision())
// Se revision ainda é 5 → SUCESSO (retorna revision 6)
// Se revision mudou    → ERRO (outro processo atualizou primeiro)
```

Funciona como um **lock otimista** — não trava nada, deixa todos trabalharem, e na hora de gravar checa se houve conflito.

---

## Como Implementar

### 2 tipos de conflito CAS

**Tipo 1 — Mensagem duplicada (DISCARD)**

```
Worker A processa callback de instance 123, node "code_1":
  Get → rev 15, status=waiting, currentNode="code_1"
  Executa inline: code_result → set_state → ...
  Cada step faz KV Put (rev 16, 17, 18...)
  Demora 35s total, AckWait=30s

NATS reentrega para Worker B (timeout):
  Get → rev 18, currentNode="set_state_2" (avançou por KV per-step!)
  status != waiting no node "code_1" → DISCARD (ACK sem fazer nada)

Worker A termina (atrasado):
  Próximo KV Put → Update(rev=18) mas rev já é 19+ (Worker B avançou)
  → CAS FALHA → Re-Get → detecta que outro worker avançou → DISCARD
```

**Tipo 2 — Atualização legítima concorrente (RETRY)**

```
Fanout: branch_2 e branch_3 completam simultaneamente via callback:

Worker A (branch_2):
  Get → rev 7
  branch_2=completed → Update(rev=7) → SUCESSO → rev 8

Worker B (branch_3):
  Get → rev 7
  branch_3=completed → Update(rev=7) → FALHA

  branch_3 NÃO é duplicada! É trabalho LEGÍTIMO.
  Re-Get → rev 8, vê branch_2=completed
  branch_3=completed → Update(rev=8) → SUCESSO → rev 9
  → RETRY funcionou ✓
```

### Algoritmo de decisão

```go
func (r *RuntimeService) handleCASConflict(
    instanceID string,
    myNodeID string,      // Node que estou processando
    myBranchID string,    // Branch que estou completando (se fanout)
) CASAction {
    // Re-Get estado fresco do NATS KV
    fresh, rev := r.natsKV.Get("inst:" + instanceID)

    // CASO 1: Estou completando uma branch de fanout
    if myBranchID != "" {
        if fresh.NodeStates[myBranchID] == nil || fresh.NodeStates[myBranchID]["waitType"] == nil {
            return DISCARD  // Outra entrega já completou esta branch
        }
        return RETRY  // Branch ainda waiting → meu trabalho é legítimo
    }

    // CASO 2: Estou processando callback/signal de um node
    if fresh.NodeStates[myNodeID] == nil || fresh.Status != StatusWaiting {
        return DISCARD  // Instance avançou além do meu node → duplicada
    }

    // Instance ainda waiting no meu node, mas revision mudou
    // (ex: outro signal atualizou state mas não resolveu wait_for)
    return RETRY
}
```

### Limites e fallback

```
CAS retry: máximo 3 tentativas
  Tentativa 1: Re-Get → re-apply → Update(freshRev)
  Tentativa 2: Re-Get → re-apply → Update(freshRev)
  Tentativa 3: Re-Get → re-apply → Update(freshRev)

  Falha após 3 → NACK (mensagem volta pro NATS para redelivery)
  NACK esgota RetryPolicy (5 tentativas) → DLQ
```

### Tabela de cenários

```
╔═════════════════════════════════════════╦═════════╦════════════════════════════════════════╗
║ Cenário                                 ║ Ação    ║ Razão                                  ║
╠═════════════════════════════════════════╬═════════╬════════════════════════════════════════╣
║ Callback duplicado (outro worker fez)   ║ DISCARD ║ currentNode avançou                    ║
║ Signal duplicado (status != waiting)    ║ DISCARD ║ Instance não espera mais               ║
║ Timer duplicado (já resumiu)            ║ DISCARD ║ status != waiting                      ║
║ Branch callback + outra branch atualizou║ RETRY   ║ Minha branch ainda waiting             ║
║ Signal + outro signal atualizou state   ║ RETRY   ║ Instance ainda em wait_for             ║
║ 3 retries CAS falharam                 ║ NACK    ║ Concorrência anormal, re-enfileirar    ║
║ NACK esgotou RetryPolicy               ║ DLQ     ║ Mensagem envenenada ou bug             ║
╚═════════════════════════════════════════╩═════════╩════════════════════════════════════════╝
```

### Impacto do KV per-step no CAS

Com KV per-step, a revision avança a cada step inline (~15 vezes por workflow médio). Isso **melhora** a detecção de duplicadas: se Worker A executou 3 steps antes do timeout, a revision avançou 3x. Worker B ao fazer Get já vê um estado mais recente e pode detectar duplicada mais cedo (via currentNodeID ou status check), sem precisar executar nada.

O CAS conflict rate **diminui** com KV per-step porque a janela de ambiguidade é menor — o estado reflete o progresso real, não apenas checkpoints esparsos.

### Frequência em produção

```
CAS conflict rate: < 0.01% das mensagens
  → NATS queue group garante 1 worker por mensagem
  → Conflito apenas em: redelivery por timeout + callbacks simultâneos de branches
  → KV per-step reduz ambiguidade: revision reflete progresso real
  → Custo do retry: 1 NATS KV Get extra (~1ms) — negligível
```

### Checklist de implementação

```
No módulo runtime:
  □ handleCASConflict() — algoritmo DISCARD/RETRY
  □ CAS retry loop (max 3 tentativas)
  □ Métricas: workflow_cas_conflicts_total{action=discard|retry}
  □ Status check em todo callback/signal/timer handler (proteção 1)

Idempotência downstream (em cada serviço):
  □ Collection idempotency_keys com TTL index (24h)
  □ Lookup por {instanceId}:{nodeId} antes de processar
  □ Re-publish resultado cacheado em duplicata com status=="done"
  □ Discard silencioso em duplicata com status=="processing"
```

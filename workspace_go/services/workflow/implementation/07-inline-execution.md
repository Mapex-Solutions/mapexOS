# Inline Execution — Até onde executa sync

Decisão: E2

---

## Problema

Quando um workflow inicia, o runtime precisa decidir: executar cada node individualmente (1 mensagem NATS por node) ou executar vários nodes sequenciais sem parar (inline). Execução per-node gera overhead de NATS (~1ms por hop). Execução inline é ~1000x mais rápida para nodes síncronos, mas precisa de limites para não travar o worker.

---

## Como Estamos Resolvendo

**Executar inline até encontrar um node async, end, ou error.** Fanout é inline (goroutines no mesmo worker). Hard limits protegem contra loops infinitos e workflows gigantes.

### Para inline quando encontra:

- Node **async** (trigger_event, code, subworkflow, delay, wait_signal, wait_for) → checkpoint
- **End** → completion
- **Error** → failure

### NÃO para (continua inline):

- **Fanout** → spawna goroutines no mesmo worker (estado já em memória)
- set_state, condition, switch, log, goto, merge, sequence, loop → tudo inline

### Hard limits

```
Max inline steps:  500        → Se atingir: checkpoint + re-enqueue no WORKFLOW-RESUME
                                (próximo worker continua de onde parou)
Context timeout:   30s        → Se atingir: error (workflow falha)
```

---

## Como Implementar

### Runtime executeInline

```go
func (r *RuntimeService) executeInline(ctx context.Context, instance *WorkflowInstance, graph *ExecutionGraph) {
    steps := 0
    for {
        steps++
        if steps > 500 {
            // Checkpoint + re-enqueue para evitar starvation
            r.checkpoint(instance)
            r.publishResume(instance)
            return
        }

        result, err := executor.Execute(ctx, execCtx)
        if err != nil { ... }

        // KV per-step: persiste estado após CADA step executado.
        // Garante que crash recovery retoma do último step completado,
        // sem re-executar operações destrutivas (increment, append, remove).
        instance.applyResult(result)
        if err := r.kvStore.Put(instance); err != nil {
            // KV Put falhou (network partition, NATS down)
            // → NÃO faz ACK → NATS reentrega → recovery automático
            return
        }

        if result.NodeState != nil && result.NodeState["waitType"] != nil {
            r.checkpoint(instance) // async node → publica checkpoint para Archiver
            return
        }

        nextNodeID := graph.Resolve(currentNodeID, result.OutputHandles)
        if nextNodeID == "" {
            r.complete(instance) // end → completion
            return
        }

        currentNodeID = nextNodeID
    }
}
```

### Na prática

```
99% dos workflows: < 50 nodes inline → nunca atinge limite
Goto loop infinito: detectado pelo max 500 steps → checkpoint/re-enqueue ou error
Workflow gigante (300 nodes set_state): executa 500, checkpoint, continua → funciona
```

### Recovery após crash (KV per-step)

```
Worker A executa inline: start → set_state(+1) → condition → set_state(+1) → [CRASH]
  KV tem: state={counter:2}, executionPath=[start, set_state, condition, set_state]
  Cada step fez KV Put → estado persistido até o último step completado

Worker B recebe redelivery:
  KV Get → state={counter:2}, currentNode="set_state" (último completado)
  Resolve próximo node do graph → continua do step seguinte
  counter = 2 (correto, não 4 como seria se re-executasse do início)
```

Sem KV per-step, o Worker B re-executaria todos os steps inline e o counter seria 4 (incrementado 2x a mais). Com KV per-step, cada operação destrutiva é persistida imediatamente.

### Checklist de implementação

```
No módulo runtime:
  □ executeInline() — loop com max 500 steps + KV Put per step
  □ Context com timeout de 30s
  □ KV Put após cada step (applyResult + kvStore.Put)
  □ Error handling: KV Put falha → return sem ACK → redelivery
  □ Re-enqueue no WORKFLOW-RESUME quando atinge 500 steps
  □ Métricas: workflow_inline_steps_histogram, workflow_kv_put_duration_histogram
```

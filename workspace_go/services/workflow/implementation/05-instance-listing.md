# Instance Listing — Covering Index

Decisão: A10

---

## Problema

`workflow_instances` tem documentos de 5-25KB (state + executionPath). A listagem no UI precisa de ~500 bytes (id, nome, status, node, data). Sem covering index, MongoDB lê 25KB do disco para devolver 500 bytes — desperdício enorme em escala.

Alternativa seria collection separada (denormalizada), mas isso adiciona complexidade de sincronização.

---

## Como Estamos Resolvendo

**Covering index único.** Collection separada NÃO necessária. Data SEMPRE obrigatória em toda query do UI.

### Regras do UI

```
Filtros OBRIGATÓRIOS (sempre presentes):
  orgId       → tenant isolation
  created     → date range (início/fim)

Filtros OPCIONAIS:
  status      → dropdown (running, completed, failed, waiting, cancelled)
  workflowId  → dropdown (qual workflow)

Sort:
  created DESC → default (mais recente primeiro)
  status       → alternativo (agrupar por status)

Campos retornados na listagem:
  _id, orgId, workflowId, workflowName, status, currentNodeId, created, completedAt
```

### Denormalização: workflowName

`workflowName` pertence à `workflow_definitions`, não à `workflow_instances`. Para o covering index funcionar sem join:

```
Runtime cria instância:
  1. Carrega definition do TieredCache (já tem para executar)
  2. Inclui workflowName no state message → Archiver grava no MongoDB
  3. Se definition renomear depois → instâncias antigas mantêm nome da versão que rodaram (correto)
```

---

## Como Implementar

### Covering index

```js
db.workflow_instances.createIndex({
  orgId: 1,            // Equality — sempre presente (tenant)
  created: -1,         // Sort + Range — sempre presente (date range)
  status: 1,           // Post-filter in-index (opcional)
  workflowId: 1,       // Post-filter in-index (opcional)
  workflowName: 1,     // Covering — display na tabela
  currentNodeId: 1,    // Covering — display na tabela
  completedAt: 1       // Covering — display (calcular duração)
}, { name: "idx_listing_covering" })

// Projection obrigatória no query:
{ _id:1, orgId:1, created:1, status:1, workflowId:1, workflowName:1, currentNodeId:1, completedAt:1 }
// Resultado: IXSCAN + PROJECTION_COVERED → ZERO document fetch
```

**IMPORTANTE:** A projection no query DEVE incluir SOMENTE campos que existem no index. Se algum campo fora do index for incluído, MongoDB cai para document fetch silenciosamente. Testar com `explain()` em staging.

### Como funciona cada query

```
╔══════════════════════════════════════════════════════════╦══════════════════════════════╗
║ Query                                                    ║ Behavior                     ║
╠══════════════════════════════════════════════════════════╬══════════════════════════════╣
║ orgId=X, created BETWEEN A AND B                         ║ Index prefix [orgId,created] ║
║ ORDER BY created DESC                                    ║ Range scan + sort grátis     ║
║                                                          ║ ✅ COVERING — zero doc fetch ║
╠══════════════════════════════════════════════════════════╬══════════════════════════════╣
║ orgId=X, created BETWEEN A AND B, status=Y               ║ Index prefix [orgId,created] ║
║ ORDER BY created DESC                                    ║ + post-filter status IN-INDEX ║
║                                                          ║ ✅ COVERING — zero doc fetch ║
╠══════════════════════════════════════════════════════════╬══════════════════════════════╣
║ orgId=X, created BETWEEN A AND B, workflowId=Z           ║ Index prefix [orgId,created] ║
║ ORDER BY created DESC                                    ║ + post-filter workflow IN-IDX ║
║                                                          ║ ✅ COVERING — zero doc fetch ║
╠══════════════════════════════════════════════════════════╬══════════════════════════════╣
║ orgId=X, created BETWEEN A AND B                         ║ Index prefix [orgId,created] ║
║ ORDER BY status ASC                                      ║ Range scan → in-memory sort  ║
║                                                          ║ ✅ COVERING (sort in-memory,  ║
║                                                          ║ OK pois date reduz result)   ║
╠══════════════════════════════════════════════════════════╬══════════════════════════════╣
║ _id = X (detalhe)                                        ║ Default _id index            ║
║                                                          ║ Document fetch completo (OK) ║
╚══════════════════════════════════════════════════════════╩══════════════════════════════╝
```

### Memória do index — 1M+ instâncias

```
Index entry: ~7 campos × ~50 bytes média = ~350 bytes
1M instâncias:  ~350 MB   ← cabe em RAM de produção
10M instâncias: ~3.5 GB   ← cabe em servidores 16GB+
```

### Futuro (se necessário)

Se query `orgId + workflowId sem date range` virar comum:
```js
// Index auxiliar (NÃO covering — apenas para filter+sort)
db.workflow_instances.createIndex({ orgId: 1, workflowId: 1, created: -1 })
```

Com data obrigatória, o covering index único resolve tudo.

### Checklist de implementação

```
No MongoDB seed/migration:
  □ Criar idx_listing_covering no startup

No definitions module:
  □ WorkflowInstance inclui workflowName (denormalizado)

No archiver module:
  □ BulkWrite inclui workflowName no InsertOne

Na API de listing:
  □ Validar: created range é obrigatório (400 se ausente)
  □ Projection SOMENTE campos do index
  □ Testar com explain() em staging
```

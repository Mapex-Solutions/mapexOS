# Definition Cache — TieredCache L1/L2 MinIO + NATS Fanout

Decisão: A4

---

## Problema

O Runtime precisa carregar a WorkflowDefinition (~10-50KB) para executar cada workflow. Com 1M instâncias simultâneas, fazer query no MongoDB para cada execução é inviável.

Problema clássico de **anti-stampede**: 20 pods reiniciam ao mesmo tempo, todos precisam da mesma definition, 20 queries MongoDB simultâneas.

---

## Como Estamos Resolvendo

Usar o mesmo padrão do asset service (módulo `assettemplates`). MinIO como L2 compartilhado elimina o anti-stampede. NATS FANOUT para invalidação cross-pod.

```
WRITE (Create/Update/Delete WorkflowDefinition):
  1. Grava/atualiza no MongoDB (source of truth)
  2. Grava/atualiza no MinIO L2 (cache compartilhado)
     KEY: definitions/{orgId}/{definitionId}.json
  3. Publica NATS FANOUT: "fanout.workflow.definition.invalidate"
     PAYLOAD: { orgId, definitionId }
  4. Todos os pods recebem → deletam L1 local

READ (Runtime precisa da definition para executar):
  L1 (NVMe local) → HIT: retorna (~500µs)
  L1 MISS → L2 (MinIO) → HIT: retorna + popula L1 (~5-50ms)
  L2 MISS → Fallback HTTP → MongoDB → repopula MinIO L2 + L1 (~10-100ms)
```

### Por que isso mata anti-stampede

```
Cenário: 20 pods cold start, todos precisam da mesma definition

SEM MinIO L2 (anti-stampede clássico):
  20 pods × L1 miss → 20 queries MongoDB
  (ou complexa coordenação via NATS pub/sub entre pods)

COM MinIO L2 (nosso padrão):
  20 pods × L1 miss → 20 reads MinIO L2 → 0 queries MongoDB
  MinIO/S3 foi FEITO para servir reads concorrentes (CDN-like)
  Zero coordenação entre pods necessária
```

---

## Como Implementar

### Config do TieredCache

```go
// bootstrap/cache.go
tieredCacheModel.Config{
    EnableL0:  false,                    // RAM desligado (definitions são ~10-50KB)
    EnableL1:  true,                     // NVMe local
    L1Dir:     l1BaseDir + "/" + serviceName + "/definitions",
    L1DefaultTTL: 1 * time.Hour,         // TTL longo (definitions mudam raramente)
    KeyPrefix: "definitions/",           // Mesma chave usada no MinIO L2

    EnableL2: true,
    L2Loader: func(ctx context.Context, key string) ([]byte, error) {
        // key já vem com prefix: "definitions/{orgId}/{definitionId}"
        result, err := minioClient.Get(ctx, key+".json")
        if err != nil {
            return nil, err
        }
        return result.Data, nil
    },
}
```

**IMPORTANTE:** O `KeyPrefix` do TieredCache, as chaves no MinIO, e a invalidação FANOUT DEVEM usar o mesmo padrão de chave: `definitions/{orgId}/{definitionId}`. Se divergirem, a invalidação não acerta a chave correta.

### MinIO storage

```
Bucket: mapex-templates (mesmo do asset service)
Key: definitions/{orgId}/{definitionId}.json

Payload: WorkflowDefinition completa serializada como JSON
  → Nodes, Edges, Variables, RetryPolicy, Metadata
  → ~10-50KB por definition
```

### NATS FANOUT invalidation

```go
// Publicado no Create/Update/Delete de WorkflowDefinition
nats.PublishFanout(ctx, "fanout.workflow.definition.invalidate", payload)

// Payload
{ "orgId": "507f1f77bcf86cd799439011", "definitionId": "607f1f77bcf86cd799439022" }

// Consumer (todos os pods do workflow service)
// Recebe → tieredCache.Invalidate(ctx, "definitions/{orgId}/{definitionId}")
//        → Deleta L1 NVMe local
//        → Próximo read vai pro L2 MinIO (já atualizado)
```

### Fallback endpoint (internal, multi-tenant)

```
GET /internal/definitions/{orgId}/{definitionId}
  → orgId obrigatório no path (multi-tenant isolation)
  → Busca do MongoDB com filtro orgId + definitionId
  → Grava no MinIO L2: definitions/{orgId}/{definitionId}.json (repopula)
  → Retorna JSON
  → Usado quando L2 MinIO também não tem (edge case: MinIO limpo, migration, etc.)
```

### Referência de implementação

```
Copiar padrão de:
  services/assets/src/modules/assettemplates/application/services/assettemplate_service.go
    → writeScripts() (linhas 824-876)
    → publishTemplateInvalidate() (linhas 884-907)
    → GetTemplateByIdForCacheFallback() (linhas 679-710)

  services/assets/src/modules/assettemplates/infrastructure/storage/minio/
    → template_storage_adapter.go (WriteScripts, DeleteScripts)

  services/assets/src/bootstrap/cache.go
    → TieredCache config com L2Loader (linhas 69-104)
```

### Resiliência de invalidação

```
Se FANOUT de invalidação se perder (network issue, pod restart durante delivery):
  → Pod mantém definition stale no L1
  → TTL de 1h é o safety net — após TTL, L1 expira, próximo read vai pro L2 (atualizado)
  → Definitions mudam raramente (< 100 vezes/dia) → window de stale é aceitável
  → Em caso extremo: restart do pod limpa L1 inteiro
  → Aceitável. Futuro: DefinitionVersion no header da instância permite detectar mismatch
```

### Números para 1M+ instâncias

```
WorkflowDefinitions: tipicamente < 10K definitions (mesmo em enterprise)
MinIO L2: 10K × 50KB = 500MB (trivial para MinIO)
L1 NVMe por pod: 10K × 50KB = 500MB (trivial para NVMe)
Invalidation: << 100 msgs/dia (definitions mudam raramente)
Cold start de 20 pods: 20 × MinIO reads = < 1s total
```

### Checklist de implementação

```
No definitions module:
  □ Service: Create/Update → MongoDB → MinIO L2 → NATS FANOUT invalidation
  □ Service: Delete → MongoDB → MinIO L2 delete → NATS FANOUT invalidation
  □ Service: GetByIdForCacheFallback() → MongoDB → repopula MinIO L2

No bootstrap:
  □ cache.go → TieredCache config com L2Loader (MinIO)
  □ FANOUT consumer → invalida L1 local no recebimento

No runtime module:
  □ Carregar definition via TieredCache (L1 → L2 → fallback)
  □ Cachear ExecutionGraph construído (mesma TTL)
```

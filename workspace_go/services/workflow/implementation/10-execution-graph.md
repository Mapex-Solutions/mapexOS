# Execution Graph — BuildGraph + Edge Cases

Decisões: E1, E3-E5

---

## Problema

O frontend gera um JSON com nodes e edges. O runtime precisa de uma estrutura de dados eficiente para:

1. Resolver "dado este node e este handle, qual o próximo node?"
2. Resolver goto pairs (sender → receiver por pairLabel)
3. Filtrar nodes visuais (text_note, group_frame) que não participam da execução

---

## Como Estamos Resolvendo

`ExecutionGraph` com adjacency map, construído uma vez da definition, cacheado no TieredCache.

---

## Como Implementar

### BuildGraph

```go
func BuildGraph(def *WorkflowDefinition) *ExecutionGraph {
    graph := &ExecutionGraph{
        Adjacency: make(map[string]map[string]string),
        Nodes:     make(map[string]*WorkflowNode),
        GotoPairs: make(map[string]string),
    }

    // 1. Indexa nodes (filtra visuais)
    // Nota: Go 1.22+ cria variável nova por iteração, &node é seguro.
    // Projeto usa Go >= 1.25.
    for _, node := range def.Nodes {
        if node.Type == "core/text_note" || node.Type == "core/group_frame" {
            continue  // Visual-only, não participa da execução
        }
        graph.Nodes[node.ID] = &node
    }

    // 2. Resolve goto pairs (receiver por pairLabel)
    for _, node := range graph.Nodes {
        if node.Type == "core/goto" {
            cfg := parseGotoConfig(node.Config)
            if cfg.Role == "receiver" {
                graph.GotoPairs[cfg.PairLabel] = node.ID
            }
        }
    }

    // 3. Constrói adjacency a partir de edges (filtra edges visuais)
    for _, edge := range def.Edges {
        if _, ok := graph.Nodes[edge.Source]; !ok { continue }  // Source é visual (text_note/group_frame) → skip
        if _, ok := graph.Nodes[edge.Target]; !ok { continue }  // Target é visual → skip
        handle := edge.SourceHandle
        if handle == "" { handle = "out" }
        if strings.HasPrefix(handle, "__") { continue }  // Handles visuais (__note_out, __note) → skip
        if graph.Adjacency[edge.Source] == nil {
            graph.Adjacency[edge.Source] = make(map[string]string)
        }
        graph.Adjacency[edge.Source][handle] = edge.Target
    }

    // 4. Injeta goto sender → receiver como edge lógica
    for _, node := range graph.Nodes {
        if node.Type == "core/goto" {
            cfg := parseGotoConfig(node.Config)
            if cfg.Role == "sender" {
                if receiverID, ok := graph.GotoPairs[cfg.PairLabel]; ok {
                    if graph.Adjacency[node.ID] == nil {
                        graph.Adjacency[node.ID] = make(map[string]string)
                    }
                    graph.Adjacency[node.ID]["out"] = receiverID
                }
            }
        }
    }

    return graph
}
```

### Detalhes

- **Filtragem:** `text_note` e `group_frame` excluídos (visual-only)
- **Goto:** sender→receiver resolvido por `pairLabel` matching
- **Ciclos infinitos:** prevenidos pelo max 500 inline steps (E2)

---

## Edge Cases Resolvidos

### E3. State mutation em fanout/merge

Cada branch recebe **cópia** do state (isolation). Merge faz **last-write-wins** por key (branch com maior index ganha). Detalhes nos executors C10 (fanout) e C11 (merge). Ver `09-executors.md`.

### E4. Loop iteration

Max 10000 iterações. `LoopState` tracking no instance state. Output handles: body/done/error. Detalhes no executor C13 (loop). Ver `09-executors.md`.

### E5. Subworkflow recursion

Depth limit 10. Child publica callback no `callbackSubject` do parent. Child failure → parent marca node como error. Detalhes no executor C9 (subworkflow). Ver `09-executors.md`.

---

## Checklist de implementação

```
No módulo runtime:
  ✅ ExecutionGraph struct (domain/entities/execution_graph.go)
  ✅ BuildGraph() function (domain/services/graph_builder.go)
  ✅ parseGotoConfig() helper (domain/services/graph_builder.go)
  ✅ ResolveNextNodes(nodeID, handles), GetNode(nodeID), HasEdge(nodeID, handle)
  ✅ ParsedConfigs — typed config parsing at graph build time (domain/services/config_parsing.go)
  ✅ Timezone resolution — literal FieldValue from WorkflowDefinition
  ✅ GoTo sender validation — executor checks HasEdge("out") for orphaned senders
  □ Cachear graph construído (mesma TTL da definition no TieredCache)
```

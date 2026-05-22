# 🚀 Quick Start - E2E Tests

Guia rápido para rodar testes E2E no MapexOS.

## 📖 TL;DR

```bash
cd workspace_go/packages/e2eTests

# Rodar teste específico
./run-tests.sh mapexos organizations

# Ver ajuda
./run-tests.sh help

# Ver módulos disponíveis
./run-tests.sh list
```

## 🎯 Formato do Comando

```bash
./run-tests.sh [SERVICE] [MODULE] [OPTIONS]
```

### Serviços Disponíveis:
- `mapexos` - organizations, roles, groups, users, memberships
- `assets` - assets, assettemplates
- `router` - routegroups
- `http_gateway` - datasources

## 📝 Exemplos Práticos

### Rodar um módulo específico
```bash
./run-tests.sh mapexos organizations   # ✅ Organizations do mapexos
./run-tests.sh mapexos roles           # ✅ Roles do mapexos
./run-tests.sh mapexos users           # Users do mapexos
./run-tests.sh assets assets           # Assets do serviço assets
./run-tests.sh router routegroups      # Routegroups do router
```

### Rodar todos os testes de um serviço
```bash
./run-tests.sh mapexos                 # Todos os módulos do mapexos
./run-tests.sh assets                  # Todos os módulos do assets
```

### Rodar TODOS os testes
```bash
./run-tests.sh all                     # Todos os testes E2E
```

### Com opções
```bash
./run-tests.sh mapexos users -q        # Sem verbose (quieto)
./run-tests.sh mapexos users -p 4      # 4 workers paralelos
./run-tests.sh mapexos users -t 10m    # Timeout de 10 minutos
```

## 🔧 Comandos Úteis

```bash
./run-tests.sh check    # Verificar se serviços estão rodando
./run-tests.sh list     # Listar serviços/módulos disponíveis
./run-tests.sh help     # Ajuda completa
```

## 🎯 Autocompletar (TAB)

### Habilitar temporariamente
```bash
cd /path/to/e2eTests
source .run-tests-completion.bash

# Agora use TAB:
./run-tests.sh [TAB]           # Lista: all login check list mapexos assets...
./run-tests.sh mapexos [TAB]   # Lista: organizations roles groups users...
```

### Habilitar permanentemente
Adicione ao seu `~/.bashrc`:
```bash
source /path/to/e2eTests/.run-tests-completion.bash
```

## ⚙️ Opções

| Opção | Descrição | Exemplo |
|-------|-----------|---------|
| `-q, --quiet` | Sem verbose | `./run-tests.sh mapexos users -q` |
| `-p N` | N workers paralelos | `./run-tests.sh mapexos users -p 4` |
| `-t TIME` | Timeout customizado | `./run-tests.sh mapexos users -t 10m` |

## 📊 Status dos Testes

| Módulo | Status | Testes |
|--------|--------|--------|
| organizations | ✅ | 16 passing, 1 skipped |
| roles | ✅ | 15 passing, 1 skipped |
| groups | ✅ | ~15 tests |
| users | ✅ | ~12 tests |
| memberships | 🚧 | TODO |
| assets | 🚧 | TODO |
| assettemplates | 🚧 | TODO |
| routegroups | 🚧 | TODO |
| datasources | 🚧 | TODO |

## 🐛 Troubleshooting

### Erro: "mapexos (5000) is NOT running"
```bash
cd workspace_go
make run
```

### Erro: "Failed to generate admin token"
```bash
# Verificar se o usuário admin existe
# Email: admin@mapex.global
# Password: mapex123
```

### Ver logs detalhados
```bash
./run-tests.sh mapexos organizations -v
```

## 📚 Links Úteis

- [README Completo](./README.md) - Documentação completa
- [Conventions](./README.md#convenções) - Convenções de código
- [Fixtures](./README.md#fixtures) - Como usar fixtures

## 💡 Dicas

1. **Use `list`** para ver módulos disponíveis:
   ```bash
   ./run-tests.sh list
   ```

2. **Use TAB** para autocompletar (depois de `source .run-tests-completion.bash`)

3. **Use `-q`** para menos output quando rodar múltiplos testes

4. **Use `-p 4`** para rodar mais rápido em paralelo

5. **Sempre rode `check`** antes de rodar testes:
   ```bash
   ./run-tests.sh check
   ```

## 🚀 Workflow Recomendado

```bash
# 1. Verificar serviços
./run-tests.sh check

# 2. Ver módulos disponíveis (se necessário)
./run-tests.sh list

# 3. Rodar teste específico
./run-tests.sh mapexos organizations

# 4. Rodar todos os testes do serviço
./run-tests.sh mapexos

# 5. (Opcional) Rodar todos os testes E2E
./run-tests.sh all
```

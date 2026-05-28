# MapexOS

> **IoT-first, mas não se limita a IoT.**
> O MapexOS não vê dispositivos ou sensores — ele vê **Assets**.
> Qualquer fonte. Qualquer protocolo. Uma única abstração.
>
> **Connect. Automate. Scale.** — A plataforma aberta para integração de
> dados e automação inteligente.

```
   Fontes                        MapexOS                          Destinos
   ──────                        ───────                          ────────
   Devices ──┐                                              ┌── Webhooks / APIs
   Gateways ─┤   Ingest → Validate → Transform → Route →    ├── Slack / Teams / Email
   APIs ─────┼──        Store / Notify / Automate           ├── NATS / MQTT
   Apps ─────┤                                              └── Plugins customizados
   Terceiros ┘
```

Este repositório contém o código-fonte de todos os serviços do MapexOS,
o frontend Vue 3 e os pacotes compartilhados de Go e TypeScript. A
distribuição executável — Docker Compose com imagens multi-arch
pré-construídas — vive em
[`mapexOSDeploy`](https://github.com/Mapex-Solutions/mapexOSDeploy).

[English version](./README.md) · [Site de documentação](https://mapexos.io)

---

## O que o MapexOS entrega

| | |
|---|---|
| **Ingestão de telemetria** | HTTP e MQTT autenticados. O broker MQTT (um plugin do Mosquitto) toma toda decisão de CONNECT/PUBLISH localmente, com base em um cache de três camadas (Pebble → MinIO → fallback HTTP) — zero round-trip no caminho quente. |
| **Conversão genérica de dados** | Pipelines por asset — *preprocess → validate → convert* — escritos em JavaScript e executados em isolates V8. Aceita qualquer payload de dispositivo e normaliza para o schema de evento da plataforma sem mexer no código fonte. |
| **Roteamento dinâmico de eventos** | Regras de match sobre qualquer campo do evento — `payload.temperature > 30`, `device.tag in [...]`, JSONPath, regex. As regras vivem no MongoDB e recarregam em runtime; sem migração de schema, sem redeploy. |
| **Workflow engine** | Um runtime de DAG inspirado no [Temporal.io](https://temporal.io) — execução determinística, retries, timers, sub-workflows, triggers idempotentes. O estado sobrevive a reinício de processo. |
| **UI de plugins (estilo n8n)** | Nós customizados de workflow são entregues como plugins servidos por um registry de manifestos no estilo CDN. Novos conectores HTTP / MQTT / NATS / Slack / banco de dados aparecem no editor sem rebuild do frontend. |
| **Multi-tenant por design** | Hierarquia de organizações (pai → filhas), isolamento por organização, templates compartilhados que se propagam para organizações descendentes. |
| **RBAC + grupos** | Permissões granulares e participação em grupos, avaliadas centralmente pelo `mapex-iam`. Toda chamada entre serviços carrega um contexto de identidade. |
| **Self-hostable, multi-arch** | Um único `docker compose up -d` sobe a stack inteira. As imagens são publicadas para `linux/amd64` e `linux/arm64` — rodam em servidores Linux, Apple Silicon, Windows com Docker Desktop e Raspberry Pi 4/5. |

---

## O ecossistema MapexOS

O MapexOS é dividido em quatro repositórios abertos. A maioria dos
usuários só precisa do repositório de deploy; os outros são para
contribuidores, operadores do broker e integradores Go.

| Repositório | Papel |
|---|---|
| **[mapexOS](https://github.com/Mapex-Solutions/mapexOS)** *(este repo)* | Código-fonte dos onze serviços de backend e do frontend Vue 3. |
| **[mapexOSDeploy](https://github.com/Mapex-Solutions/mapexOSDeploy)** | Distribuição Docker Compose que puxa as imagens pré-construídas do Docker Hub. **Comece por aqui para rodar a plataforma.** |
| **[mapexMQTTBroker](https://github.com/Mapex-Solutions/mapexMQTTBroker)** | O broker MQTT de produção — Eclipse Mosquitto v2 mais o plugin Go interno que cuida de auth, ACL, presença e ingressão em um único `.so`. |
| **[mapexGoKit](https://github.com/Mapex-Solutions/mapexGoKit)** | Bibliotecas Go compartilhadas usadas por todo serviço Go — middleware HTTP, helpers NATS, observabilidade, validação, contracts. |

---

## O que tem neste repo

```
mapexOS/
├── workspace_go/                # serviços Go (DDD + hexagonal)
│   ├── services/
│   │   ├── mapexIam/            # usuários, organizações, papéis, RBAC, auth
│   │   ├── http_gateway/        # ingestão por webhook, registry de datasources
│   │   ├── assets/              # assets IoT, templates, campos EVA
│   │   ├── router/              # roteamento de eventos, regras de match
│   │   ├── events/              # armazenamento ClickHouse, 7 consumers NATS
│   │   ├── triggers/            # 8 executores (HTTP, MQTT, NATS, NATS-JS, NATS-KV, NATS-OBJ, NATS-RPC, Webhook)
│   │   ├── workflow/            # workflow engine inspirado em Temporal + plugins + credenciais
│   │   └── mapexVault/          # cofre de credenciais, autoridade PKI
│   └── packages/
│       └── contracts/           # DTOs entre serviços (fonte única da verdade)
│
├── workspace_js/                # serviços Node + frontend
│   ├── services/
│   │   ├── js-executor/             # isolates V8 para scripts de evento IoT
│   │   ├── js-workflow-executor/    # isolates V8 para nós de código no workflow
│   │   └── plugin-marketplace-mock/ # CDN estático de manifestos de plugin (dev)
│   ├── apps/
│   │   └── mapexOS/             # SPA Vue 3 + Quasar
│   └── packages/                # schemas Zod, wrappers de API tipados, clientes de infra
│
└── docs/                        # documentação pública finalizada
```

---

## Início rápido

Para **rodar** a plataforma, vá para o `mapexOSDeploy`:

```bash
git clone https://github.com/Mapex-Solutions/mapexOSDeploy.git
cd mapexOSDeploy
docker compose up -d
```

Abra <http://localhost> e faça login com `admin@mapex.local` /
`mapex@123`. O passo a passo completo — incluindo um quickstart de
sensor de temperatura via HTTP ou MQTT — vive em
[`mapexOSDeploy/quickstart/`](https://github.com/Mapex-Solutions/mapexOSDeploy/tree/main/quickstart).

---

## Por que o projeto existe

A maioria das "plataformas IoT" para na ingestão e nos dashboards. A
parte difícil começa depois que o dado chega: rotear para o sistema
downstream certo, rodar workflows com estado que sobrevivem a
reinício de processo, deixar operadores construírem novas integrações
sem redeploy do UI e fazer tudo isso para vários tenants em uma única
stack.

O MapexOS trata esses quatro problemas como cidadãos de primeira
classe — não como puxadinhos em cima de um banco de séries temporais.
O resultado é uma plataforma que dá pra rodar uma frota de verdade,
não uma demo com um widget de temperatura.

---

## Licença

O MapexOS é distribuído sob a **Business Source License 1.1** — veja
[`LICENSE`](./LICENSE) para os termos completos.

Em resumo:

- **Você pode** self-hostar, modificar e rodar o MapexOS em produção
  para a sua organização.
- **Você não pode** oferecer o MapexOS como serviço comercial
  hospedado para terceiros sem um acordo comercial separado.
- Na **Change Date**, a licença converte para **Apache 2.0**.

Para licenciamento comercial, fale com a **Mapex Solutions**.

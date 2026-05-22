# Fase 3 — Cascata do TieredStore L1 -> L2 -> L3

## O que este teste prova

A cascata do TieredStore do plugin do broker (Pebble L1 -> Redis L2 ->
Mongo L3 fallback) mais o caminho de lazy-pull por fanout-invalidate.
A fase encadeia o ciclo de vida da senha da fase 1 (asset criado, L1
aquecido), depois força misses tier a tier e verifica que o log do
broker mostra o hit da camada certa, terminando com um
fanout-invalidate manual que precisa expulsar o L1 e disparar um
re-fetch no próximo CONNECT.

A fase cobre:

1. Prefixo da fase 1 completa — CONNECT por senha funciona, L1 aquecido.
2. Força miss em L1; próximo CONNECT cai em L2 e o broker loga "L2 hit".
3. Força miss em L2; próximo CONNECT cai em L3 fallback e o broker loga "L3 fallback".
4. Publica um fanout-invalidate manual; broker loga "invalidated L1" e o próximo CONNECT re-busca ponta a ponta.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase3_cascade/
```

## Requisitos

- Todos os requisitos da fase 1 (ciclo por senha) mais acesso às
  camadas L1 (Pebble) e L2 (Redis) do plugin do broker para que o saga
  consiga forçar um miss em cada tier.
- Broker MQTT acessível em `tcp://localhost:1883`.
- Log do broker acessível ao runner para que os asserts de hit por
  camada possam fazer grep por `L2 hit`, `L3 fallback` e `invalidated L1`.

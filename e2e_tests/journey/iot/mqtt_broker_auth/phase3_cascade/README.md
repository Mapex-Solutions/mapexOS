# Phase 3 — TieredStore cascade L1 -> L2 -> L3

## What this test proves

The broker plugin's TieredStore cascade (Pebble L1 -> Redis L2 -> Mongo
L3 fallback) plus the fanout-invalidate lazy-pull path. The phase
chains onto the password lifecycle from phase 1 (asset created, L1
warmed), then forces tier misses one at a time and asserts the broker
log surfaces the right cascade hit, followed by a manual fanout
invalidate that must evict L1 and trigger a refetch on the next
CONNECT.

The phase covers:

1. Phase 1 prefix completes — asset password CONNECT succeeds, L1 is warmed.
2. Force L1 miss; next CONNECT hits L2 and the broker logs an "L2 hit".
3. Force L2 miss; next CONNECT hits L3 fallback and the broker logs an "L3 fallback".
4. Publish a manual fanout invalidate; broker logs "invalidated L1" and the next CONNECT re-fetches end-to-end.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase3_cascade/
```

## Requirements

- All requirements of phase 1 (password lifecycle) plus access to the
  broker plugin's L1 (Pebble) and L2 (Redis) tiers so the saga can
  force a miss on each tier.
- MQTT broker reachable on `tcp://localhost:1883`.
- Broker log accessible to the runner so the cascade-hit asserts can
  grep for `L2 hit`, `L3 fallback`, and `invalidated L1`.

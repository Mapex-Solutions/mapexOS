# Phase 1 — Connectivity-driven Email trigger smoke

## What this test proves

Email trigger fires through the healthmonitor → router → triggers
chain and delivers a real message to an in-process SMTP server.

The phase:

1. Starts an in-process SMTP server on `SmtpSinkBindAddr` (default `0.0.0.0:11025`).
2. Creates an Email trigger pointing at that sink.
3. Creates two `kind=trigger` route groups (one for `online`, one for `offline`).
4. Creates an HTTP connectivity asset wired to both route groups.
5. Sends a warm-up heartbeat → asset settles to `online` (silent — no trigger).
6. Force-offline by admin → offline RG fires → SMTP sink captures **1 message**.
7. Asserts the captured message has the configured `from`, `to`, and a subject fragment.
8. Sends a new heartbeat → asset goes `online` → online RG fires → SMTP sink captures **2 messages**.
9. Deletes the asset; Compensate chain rolls everything else back.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_email/phase1_connectivity/...
```

Or via the helper:

```bash
./run-tests.sh saga trigger-email
```

## Requirements

- Live stack with these services running (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Check with `./run-tests.sh check`.
- Port `11025` free on the host (override with `SAGA_TRIGGER_SMTP_BIND_ADDR`).
- When the triggers service runs in Docker, set `SAGA_TRIGGER_SMTP_HOST=host.docker.internal` so the trigger's `smtpHost` reaches the host SMTP listener.
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 logs in as that.

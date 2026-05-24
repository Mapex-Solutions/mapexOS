# Journey: Email trigger

## What this journey proves

The Email (SMTP) trigger of the triggers service works end-to-end.
Each phase exercises a different firing path; both deliver a real
message to an in-process SMTP server and validate envelope + subject.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger fires. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetry to gateway → js-executor → router → trigger fires. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-email
```

## Requirements

Each phase README lists protocol-specific requirements. Common:

- Live stack: `mapexos`, `assets`, `router`, `http_gateway`, `triggers` on default ports.
- Free port `11025` on the host (override `SAGA_TRIGGER_SMTP_BIND_ADDR`).

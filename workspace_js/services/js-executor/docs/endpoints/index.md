# Endpoints

## HTTP API

### App
| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |

### Scripts (base: `/api/v1/scripts`)
| Method | Path | Description |
|---|---|---|
| POST | `/test` | Test script execution |
| GET | `/templates/:templateId/script_test` | Get raw scriptTest from template |
| GET | `/templates/:templateId/sample_payload` | Get processed sample payload |

### Internal (base: `/internal/templates`)
| Method | Path | Description |
|---|---|---|
| GET | `/:orgId/:templateId/script_test` | Internal scriptTest (API key) |
| GET | `/:orgId/:templateId/sample_payload` | Internal sample payload (API key) |

## NATS
| Subject | Description |
|---|---|
| `processor.js.execute` | Script execution requests |

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`

# Endpoints

## HTTP API

### Assets (base: `/api/v1/assets`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List assets (paginated, filtered) |
| GET | `/counter` | Count assets |
| POST | `/` | Create asset |
| GET | `/:assetId` | Get asset by ID |
| PATCH | `/:assetId` | Update asset |
| DELETE | `/:assetId` | Delete asset (HARD delete; drops cert state and L2 row) |

### Asset Templates (base: `/api/v1/asset_templates`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List asset templates (paginated, filtered) |
| GET | `/counter` | Count asset templates |
| POST | `/` | Create asset template |
| GET | `/:assetTemplateId` | Get asset template by ID |
| PATCH | `/:assetTemplateId` | Update asset template |
| DELETE | `/:assetTemplateId` | Delete asset template |
| GET | `/:assetTemplateId/available_fields` | Get available EVA fields |

### MQTT Certificates (base: `/api/v1/mqtt_certs`)
JWT-gated + coverage-injected. All routes wrapped in `RequireCAReady` — return `503 ca_not_ready` until the intermediate CA is mounted from mapexVault.

| Method | Path | Description |
|---|---|---|
| POST | `/` | Issue a new device cert (signs locally with the RAM-cached intermediate CA; fans out asset-invalidate; returns the leaf cert + chain) |
| DELETE | `/:serial` | Revoke a cert (writes a 30-day TTL row to `mqttRevokedCertificates`; fans out asset-invalidate) |
| GET | `/` | List certificates by asset (current + revoked within the 30-day audit window) |

### Health Monitoring (base: `/api/v1/heartbeat`, `/api/v1/health-monitor`)
| Method | Path | Description |
|---|---|---|
| POST | `/heartbeat?ds={dataSourceId}` | Explicit HTTP heartbeat — body `{ assetUUID }`. Publishes to `${env}.mapexos.asset.heartbeat.{orgId}`. |
| GET | `/health-monitor/state/:assetId` | Current online/offline state + last heartbeat timestamp |

### Internal APIs
All routes gated by the standard `ApiKeyAuthMiddleware` (`X-API-Key` header). The mapex-mqtt-broker plugin uses its own internal API key on the same header for L3 lookups.

| Method | Path | Description |
|---|---|---|
| GET | `/internal/assets/:assetUUID` | L3 read-model fallback for every consumer that caches `AssetReadModel` (Router, JS-Executor, Events, mapex-mqtt-broker plugin). Fetches from Mongo and repopulates L2 (MinIO) on the way out. The response carries `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` so the broker plugin decides MQTT CONNECTs locally — there is NO separate auth callout endpoint. |
| GET | `/internal/templates/:templateId` | Fetch template read model and repopulate L2 |

## NATS

### Published (this service)
| Subject | Stream | Description |
|---|---|---|
| `${env}.mapexos.fanout.asset.invalidate` | `MAPEXOS-FANOUT` | Asset cache invalidation fanout — fires on asset CUD AND on cert issue/revoke (cert state rides asset invalidate). |
| `${env}.mapexos.fanout.template.invalidate` | `MAPEXOS-FANOUT` | Template cache invalidation fanout on template CUD. |
| `${env}.mapexos.asset.heartbeat.{orgId}` | `MAPEXOS-ASSETS-HEARTBEAT` | Republished by `http_gateway` on explicit HTTP heartbeat; published by `js-executor` on implicit heartbeat — consumed by `healthmonitor`. Listed here for reference; this service is the consumer end. |
| `${env}.mapexos.healthmonitor.scan.schedule` | `MAPEXOS-ASSETS-HEALTH-MONITOR` | Next scheduled scan (WorkQueue + AllowMsgSchedules + `MsgId=hm-scan` for single pending across pods). |

### Subscribed (this service)
| Subject | Stream | Direction | Description |
|---|---|---|---|
| `${env}.mapexos.asset.heartbeat.>` | `MAPEXOS-ASSETS-HEARTBEAT` | subscribe (DLQ) | Heartbeat ingestion — origin-agnostic; both implicit (`js-executor`) and explicit (`http_gateway`) flow here. MQTT-protocol assets do NOT use this path. |
| `${env}.mapexos.mqtt.presence.advisory` | `MAPEXOS-ASSETS-MQTT-PRESENCE` | subscribe (DLQ) | Single subject for broker presence; two durables (`assets-mqtt-presence`, `assets-mqtt-presence-connect`) gate by `Event=connect\|disconnect`. Payload `healthmonitor.PresenceAdvisory`. |
| `${env}.mapexos.healthmonitor.scan` | `MAPEXOS-ASSETS-HEALTH-MONITOR` | subscribe (DLQ) | Fires the next offline scan; handler republishes the next schedule. |
| `${env}.mapexos.lists.name_updated` | `MAPEXOS-LISTS` | subscribe | Sync denormalized manufacturer/model/category names in templates. |

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`

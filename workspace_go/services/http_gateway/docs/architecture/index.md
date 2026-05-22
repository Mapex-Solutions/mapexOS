# Architecture

## Design
Modular architecture with clean separation between domain logic, application services, and external integrations.

## Project Structure
```
src/
├── modules/
│   ├── datasources/           # Data Source configuration management (CRUD + cache)
│   └── events/                # Event ingestion and authentication
└── shared/
    └── configuration/          # Service configuration
```

## Module Responsibilities
- `datasources`: manages Data Source entities (5 auth strategies, rate limit, asset binding) with cached lookups for performance
- `events`: validates incoming webhook authentication based on Data Source config, publishes events to NATS, and fires auth failure security events

## Request Flow: Event Ingestion
1. `POST /api/v1/events?ds=<dataSourceId>` arrives
2. Data Source ID is validated
3. Data Source configuration is resolved (from cache or database)
4. Request is authenticated according to the Data Source auth type
   - On failure: publishes auth failure event to `events.raw` and returns 401
   - On success: event is published to `processor.js.execute` via NATS
5. Returns 201 Created

## Request Flow: Data Source CRUD
1. Request arrives at `/api/v1/data_sources`
2. User identity is validated (JWT/OAuth2)
3. Request structure and parameters are validated
4. User permissions are verified (e.g., `datasources.list`)
5. Organization context is injected for tenant filtering
6. Requested operation is executed

# MapexOS API Documentation

## Overview
The MapexOS API is the **core IAM and tenant management layer**. It provides flexible roles (no fixed role set), group/membership management, and a hierarchical organization tree (Vendor → Customer → …). This allows enterprises to model real‑world structures without forcing predefined RBAC roles.

## Responsibilities
- Core IAM (users, roles, groups, memberships)
- Organization management and onboarding
- Authentication and internal auth endpoints

## Non‑Responsibilities
- Event ingestion and processing
- Asset and template management
- Rule evaluation

## Primary Data Flow
1. Client authenticates via auth endpoints
2. API serves CRUD and query endpoints for IAM and org data
3. Internal services call internal endpoints using API key auth

## Key Design Decisions
- **No fixed roles**: organizations create and name roles for their business
- **Org hierarchy**: multi‑level tenants (Vendor → Customer → …)
- **Coverage + permissions**: resolved dynamically per org tree

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)

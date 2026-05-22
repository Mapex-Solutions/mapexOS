# Tests

## Prerequisites
- Go 1.25+
- ClickHouse running locally (for integration tests)
- Redis running locally (for integration tests)
- MongoDB running locally (for integration tests)

Unit tests do not require external dependencies and can run in isolation.

## Run
```bash
go test ./... -count=1
```

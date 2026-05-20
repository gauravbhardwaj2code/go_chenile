# Slot Service

Slot Service manages runner availability and customer allocation requests for cooks, maids, and house-help runners.

## Run inside go_chenile

From this service directory:

```bash
GOWORK=$(cd ../../chenile-framework && pwd)/go.work go test ./...
GOWORK=$(cd ../../chenile-framework && pwd)/go.work go run ./cmd/slot-service
```

The service listens on port `8080` by default. Change `config/application.yaml` to use another port.

## Run after copying this service

If the Chenile modules in `go.mod` are published and tagged in the public repository, run:

```bash
go mod tidy
go test ./...
go run ./cmd/slot-service
```

If those public tags are not available yet, point this copied service at a local Chenile framework checkout first:

```bash
bash scripts/use-local-chenile.sh /absolute/path/to/go_chenile/chenile-framework
go test ./...
go run ./cmd/slot-service
```

For the current local checkout, that command is:

```bash
bash scripts/use-local-chenile.sh /Users/gauravbhardwaj/work/go/go_chenile/chenile-framework
```

## API smoke test

In another terminal:

```bash
curl -X POST http://localhost:8080/runners \
  -H "Content-Type: application/json" \
  -d '{"name":"Asha","skills":["cook"],"attributes":{"diet":"veg"},"slots":[{"date":"2026-06-01","start":"09:00","end":"11:00"}]}'

curl -X POST http://localhost:8080/allocations \
  -H "Content-Type: application/json" \
  -d '{"requestId":"req-1","skill":"cook","slot":{"date":"2026-06-01","start":"09:00","end":"11:00"},"constraints":[{"key":"diet","value":"veg","type":"hard"}]}'
```

## Tests

```bash
go test ./...
```

BDD scenarios are in `test/features/slot.feature`.

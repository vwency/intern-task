### Start up

```
docker-compose up --build
```

Or you can manually (bare-metal)

```
go install github.com/go-task/task/v3/cmd/task@latest (possible if not installed)
go mod download
go mod vendor
task run-subpub
```

### Test

```
task test-subpub-server что бы сервер протестировать который запущен
task test-subpub unit-тесты
```

### Stack

1. Graceful shutdown
2. go-kit (gRPC)

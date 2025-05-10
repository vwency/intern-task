### Start up

```
docker-compose up --build
```

Or you can manually (bare-metal)

```
go install github.com/go-task/task/v3/cmd/task@latest (possible if not installed)
go mod download
go mod tidy
go mod vendor
task run-subpub
```

### Test

```
task test-subpub-server     что бы сервер протестировать который запущен(то есть надо запустить)
task test-subpub unit-тесты
```

### Stack

1. Graceful shutdown
2. go-kit (gRPC)
3. Viper configuration

### Структура (Архитектура)

1. Микросервисная архитектура
2. Использование endpoints transport service слоев go-kit
3. Viper configuration
4. Graceful shutdown
5. Boilerplate код

# mwallet

- Lint:

```shell
$ golangci-lint run -v ./...
```

- Test:

```shell
$ go test -v ./...
```

- Build:

```shell
$ go build -v -ldflags="-s -w" -o mwallet cmd/mwallet/main.go
```

- Run:

```shell
$ ./mwallet
```

- Integration test:

```shell
$ docker-compose up -d
$ go test -v -tags integration -run TestTransferPayment ./...
```
# sqlc

## Install sqlc

```sh
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

Config can be found in `./sqlc.yaml`

## Run sqlc

```sh
sqlc generate
```

# Server migrations

## Install migrations tool
```sh
cd ~/go/bin
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
```

## Apply migrations

```sh
migrate -path database/migrations -database "postgresql://postgres:qwerty@localhost:15432/postgres?sslmode=disable" up
```

## Revert migrations

```sh
migrate -path database/migrations -database "postgresql://postgres:qwerty@localhost:15432/postgres?sslmode=disable" down
```

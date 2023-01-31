test: check
	docker-compose up -d
	sleep 2s
	go test -count=1 -v ./...
	docker-compose down

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	staticcheck ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

check: fmt lint vet
.PHONY:check

build: check
	go build
.PHONY:build

clean: go-sqlc-migrations-tests
	rm go-sqlc-migrations-tests

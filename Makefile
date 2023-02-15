test: check sqlc
	$(MAKE) -- "--run_tests"
.PHONY:test

test_with_docker_compose: check sqlc
	docker compose up -d
	sleep 2s
	$(MAKE) -- "--run_tests" || $(MAKE) -- "--docker_compose_down_on_error"
	docker compose down
.PHONY:test_with_docker_compose

--run_tests:
	go test -v -count=1 ./...
.PHONY:--run_tests

--docker_compose_down_on_error:
	docker compose down
	echo "Tests failed!"
	exit 123
.PHONY:--docker_compose_down_on_error

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	staticcheck ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

check: fmt lint vet sqlc_compile_check
.PHONY:check

build: check sqlc
	go build
.PHONY:build

clean: go-sqlc-migrations-tests
	go clean -x

sqlc: sqlc_compile_check
	sqlc generate
.PHONY:sqlc

sqlc_compile_check:
	sqlc compile
.PHONY:sqlc_compile_check

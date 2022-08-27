APP_BIN := "./bin/banners"
MIGRATIONS_BIN := "./bin/goose"
DOCKER_IMG="banners:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up:
	docker-compose -f deployments/docker-compose.yaml up -d --build

down:
	docker-compose -f deployments/docker-compose.yaml down

integration-tests:
	set -e; \
	docker-compose -f deployments/docker-compose.test.yaml up --build -d; \
	status_code=0; \
	docker-compose -f deployments/docker-compose.test.yaml run tests go test -v || status_code=$$?; \
	docker-compose -f deployments/docker-compose.test.yaml down; \
	exit $$status_code

build:
	go build -v -o $(APP_BIN) -ldflags "$(LDFLAGS)" ./cmd/banners

run: build
	$(APP_BIN) -config ./configs/.env

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(APP_BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

build_goose:
	go build -o $(MIGRATIONS_BIN) ./cmd/goose/*.go

migration: build_goose
	./bin/goose --dir=migrations --config=configs/.env.goose create $(name)

migrate: build_goose
	./bin/goose --dir=migrations --config=configs/.env.goose up

rollback: build_goose
	./bin/goose --dir=migrations --config=configs/.env.goose down

seeder: build_goose
	./bin/goose --dir=seeds --config=configs/.env.goose create $(name)

seed: build_goose
	./bin/goose --dir=seeds --config=configs/.env.goose up

.PHONY: tools build run build-img run-img version test lint

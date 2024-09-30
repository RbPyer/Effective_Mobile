BIN_PATH=bin/app
MIGRATOR_BIN_PATH=bin/migrator
APP_PATH=cmd/app/main.go
CONFIG_PATH=config/.env
MIGRATOR_PATH=cmd/migrator/main.go
MIGRATIONS_PATH=migrations
DSN=postgres://oswyndel:my_password@postgres:5432/rest

all: run

run: build
	export CONFIG_PATH=./$(CONFIG_PATH) && ./$(BIN_PATH)

migrator_build:
	CGO_ENABLED=0 GOOS=linux go build -o $(MIGRATOR_BIN_PATH) ./$(MIGRATOR_PATH)

scratch_build:
	CGO_ENABLED=0 GOOS=linux go build -o ./$(BIN_PATH) ./$(APP_PATH)

build:
	go build -o ./$(BIN_PATH) ./$(APP_PATH)

migrate:
	go run $(MIGRATOR_PATH) --storage-path=$(DSN) --migrations-path=$(MIGRATIONS_PATH)
BIN_PATH=bin/app
APP_PATH=cmd/app/main.go
CONFIG_PATH=config/config.yaml

all: run

run: build
	export CONFIG_PATH=./$(CONFIG_PATH) && ./$(BIN_PATH)

scratch_build:
	CGO_ENABLED=0 GOOS=linux go build -o ./$(BIN_PATH) ./$(APP_PATH)

build:
	go build -o ./$(BIN_PATH) ./$(APP_PATH)
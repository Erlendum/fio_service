include .env

OUT_DIR := ./out
MIGRATION_DIR = deployments/migration/
APP := app
DSN_MIGRATION = "user=$(USER) dbname=$(DBNAME) host=$(HOST) password=$(PASSWORD) sslmode=disable port=$(PORT)"

build:
	go build -o $(OUT_DIR)/$(APP) cmd/main.go

clean:
	rm -rf $(OUT_DIR)

run: $(OUT_DIR)/$(APP)
	$(OUT_DIR)/$(APP)

$(OUT_DIR)/$(APP): build


run-db-postgres:
	docker-compose --project-directory deployments up -d --build postgres

stop-db-postgres:
	docker-compose --project-directory deployments stop postgres

goose-install:
	go get github.com/pressly/goose/cmd/goose
	go install github.com/pressly/goose/cmd/goose

swag-gen:
	swag init -g ./internal/delivery/http/v1/person.go

migrate-up:
	goose --dir=$(MIGRATION_DIR) postgres $(DSN_MIGRATION) up

migrate-down:
	goose --dir=$(MIGRATION_DIR) postgres $(DSN_MIGRATION) down

test-all:
	go test ./...

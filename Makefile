ifneq (,$(wildcard ./.env))
    include .env
    export
endif

POSTGRES_CONNECTION_STRING = "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) -p $(POSTGRES_PORT):5432 -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_USER=$(POSTGRES_USER) -d postgres:16-alpine

createdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

dropdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) dropdb --username=$(POSTGRES_USER) $(POSTGRES_DB)

migrateup:
	migrate --path db/migrations -database ${POSTGRES_CONNECTION_STRING} -verbose up

migratedown:
	migrate --path db/migrations -database ${POSTGRES_CONNECTION_STRING} -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

.PHONY: postgres createdb dropdb migrateup migratedown
ifneq (,$(wildcard ./app.env))
    include app.env
    export
endif

postgres:
	docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):5432 -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_USER=$(DB_USER) -d postgres:16-alpine

createdb:
	docker exec -it $(DB_CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER_NAME) dropdb --username=$(DB_USER) $(DB_NAME)

migrateup:
	migrate -path sql/migrations -database $(DB_URL) -verbose up

migratedown:
	migrate -path sql/migrations -database $(DB_URL) -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

mock:
	mockgen -package mockdb -destination pkg/mockdb/store.go github.com/aseerkt/go-simple-bank/pkg/db Store

server:
	go run cmd/server/main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock
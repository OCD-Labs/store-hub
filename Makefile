DB_URL=postgres://root:fde24e52415e@localhost:5434/store_hub?sslmode=disable

server:
	go run ./main.go

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

migration_file:
	migrate create -ext sql -dir db/migration -seq $(file_name)

sqlc:
	sqlc generate

postgres:
	docker run --name store_hub_db -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=fde24e52415e -d postgres:15.2-alpine

createdb:
	docker exec -it store_hub_db createdb --username=root --owner=root store_hub

dropdb:
	docker exec -it store_hub_db dropdb store_hub

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 2

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

test:
	go test -v -cover -short ./...

dev:
	@echo "Starting dev server"
	@air -c ./.air.toml
	@echo "Dev server started"

.PHONY: server db_schema migration_file sqlc postgres createdb dropdb migrateup migratedown migrateup1 migratedown1
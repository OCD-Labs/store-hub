DB_URL=postgres://storehub_db_user:9iDzrJqajhSzbo3QEk8G9Oq94RYCCIyF@dpg-cibbac59aq03rjmp88og-a.oregon-postgres.render.com/storehub_db

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
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover -short ./...

.PHONY: server db_schema migration_file sqlc postgres createdb dropdb migrateup migratedown
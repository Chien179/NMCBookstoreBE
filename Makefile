postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:14-alpine
	
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root nmc_bookstore

dropdb:
	docker exec -it postgres14 dropdb nmc_bookstore

migrateup:
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/nmc_bookstore?sslmode=disable" -verbose up	

migrateup1:
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/nmc_bookstore?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/nmc_bookstore?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/nmc_bookstore?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	swag init
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/Chien179/NMCBookstoreBE/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
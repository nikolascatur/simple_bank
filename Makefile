postgres:
	docker compose -f postgres/docker-compose.yml up -d

down:
	docker compose -f postgres/docker-compose.yml down

createdb:
	docker exec -it postgres-db createdb --username=myuser --owner=myuser simple_bank

dropdb:
	docker exec -it postgres-db dropdb --username=myuser simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm --platform linux/amd64 -v "$$(pwd):/src" -w /src kjconroy/sqlc:1.4.0 generate

.PHONY: createdb postgres down dropdb migrateup migratedown sqlc
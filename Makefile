postgres:
	docker run --name cycloneDB -p 5432:5432 -e DB_HOST=host.docker.internal -e POSTGRES_USER=root -e POSTGRES_PASSWORD=cyclone -d postgres:16-alpine

createdb:
	docker exec -it cycloneDB createdb --username=root --owner=root cyclone

dropdb:
	docker exec -it cycloneDB dropdb cyclone
migrateup:
	migrate -path db/migration -database "postgresql://root:cyclone@localhost:5432/cyclone?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:cyclone@localhost:5432/cyclone?sslmode=disable" -verbose down
sqlc:
	sqlc generate

.PHONY: createdb postgres dropdb migrateup migratedown sqlc
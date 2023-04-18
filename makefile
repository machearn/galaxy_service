postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root galaxy_club

dropdb:
	docker exec -it postgres12 dropdb galaxy_club

migrateup:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/galaxy_club?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/galaxy_club?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc
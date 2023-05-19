postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root galaxy_club

dropdb:
	docker exec -it postgres14 dropdb galaxy_club

migrateup:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/galaxy_club?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/galaxy_club?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v --cover ./...

protoc: 
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/machearn/galaxy_service/db/sqlc Store

token_mock:
	mockgen -package mock_token -destination token/mock/maker.go github.com/machearn/galaxy_service/token/maker TokenMaker

sender_mock:
	mockgen -package mock_sender -destination mail/mock/sender.go github.com/machearn/galaxy_service/mail EmailSender

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

start:
	go run main.go

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test protoc mockdb token_mock sender_mock start
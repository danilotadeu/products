include .env
export

install:
	go install github.com/golang/mock/mockgen@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/golang/mock/
	go get

run:
	go run main.go

.PHONY: mock
mock:
	go generate ./...

.PHONY: test/cov
test/cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test:
	go test ./...

.PHONY: docs
docs:
	swag init

migrateup:
	migrate -path db/migrations -database ${MIGRATE_URL} -verbose up

migratedown:
	migrate -path db/migrations -database ${MIGRATE_URL} -verbose down 1
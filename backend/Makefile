.PHONY: db-sync

dev:
	air
build:
	go build ./...
fmt:
	go fmt ./...
test:
	go test -v ./...
db-sync:
	goose up
	sqlc generate
db-reset:
	goose down all
	goose up
	sqlc generate
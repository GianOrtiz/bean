build:
	go build cmd/api/main.go

test:
	go test ./...

run:
	go run cmd/api/main.go

generate:
	go generate ./...

migrate:
	migrate -database "sqlite3://bean.db" -path internal/db/migrations up

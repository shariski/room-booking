DATABASE_URL ?= postgres://bobobox:bobobox@localhost:5432/bobobox?sslmode=disable
DATABASE_URL_TEST ?= postgres://bobobox:bobobox@localhost:5432/bobobox_test?sslmode=disable

run:
	go run ./cmd/api

build:
	docker build -t bobobox-room-booking .

start:
	docker compose up -d
	migrate -database "$(DATABASE_URL)" -path migrations up

stop:
	docker compose down

test:
	migrate -database "$(DATABASE_URL_TEST)" -path migrations up
	go test ./internal/usecase -v

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-up-test:
	migrate -database "$(DATABASE_URL_TEST)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down

swagger:
	swag init -g cmd/api/main.go

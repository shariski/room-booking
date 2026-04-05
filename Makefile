DATABASE_URL ?= postgres://bobobox:bobobox@localhost:5432/bobobox?sslmode=disable

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down

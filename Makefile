DATABASE_URL ?= postgres://bobobox:bobobox@localhost:5432/bobobox?sslmode=disable
DATABASE_URL_TEST ?= postgres://bobobox:bobobox@localhost:5432/bobobox_test?sslmode=disable

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-up-test:
	migrate -database "$(DATABASE_URL_TEST)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down

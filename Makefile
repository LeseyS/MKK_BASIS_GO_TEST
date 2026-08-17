DB_MIGRATE_URL = mysql://taskuser:taskpass@tcp(localhost:3306)/taskdb?charset=utf8mb4&parseTime=True&loc=Local
MIGRATE_PATH = ./migration/mysql

up:
	docker compose up --build --force-recreate

down:
	docker compose down

test:
	go test ./... -race

integration-test:
	go test ./test/... -tags=integration -count=1 -v

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html

migrate-install:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

run:
	go run ./cmd/app/main.go
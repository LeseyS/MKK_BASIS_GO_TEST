DB_MIGRATE_URL = mysql://taskuser:taskpass@tcp(localhost:3306)/taskdb?charset=utf8mb4&parseTime=True&loc=Local
MIGRATE_PATH = ./migration/mysql

up:
	docker compose up --build --force-recreate

down:
	docker compose down

migrate-install:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all
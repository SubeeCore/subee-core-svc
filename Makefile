ifndef DB_HOST
override DB_HOST = 127.0.0.1
endif

ifndef DB_PORT
override DB_PORT = 5432
endif

ifndef DB_USER
override DB_USER = root
endif

ifndef DB_PASSWORD
override DB_PASSWORD = root
endif

ifndef DB_NAME
override DB_NAME = user-store-db
endif

run:
	docker-compose up

update:
	cd src && go mod tidy

unit-test:
	cd src && go test -v -cover ./...

create-migration:
	cd internal/database/migrations && goose create ${name} sql

migration-up:
	cd internal/database/migrations && goose postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up
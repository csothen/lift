include .env

build-sql-init:
	@./db/setup.sh ${DB_USER} ${DB_NAME}

gqlgen:
	@go run github.com/99designs/gqlgen generate

generate-all: gqlgen
	@go generate ./...

teardown:
	@docker-compose down
	@-./scripts/teardown.sh lift

start: teardown generate-all build-sql-init
	@docker-compose --env-file ./.env build --no-cache
	@docker-compose --env-file ./.env up

db-cli:
	@docker exec -it lift_db psql -d ${DB_NAME} -U ${DB_USER} -h localhost -p 5432
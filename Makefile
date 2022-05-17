gqlgen:
	@go run github.com/99designs/gqlgen generate

generate-all: gqlgen
	@go generate ./...
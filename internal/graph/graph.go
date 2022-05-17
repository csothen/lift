package graph

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/csothen/tmdei-project/internal/services"
)

func NewHandler(s *services.Service) http.Handler {
	resolver := &Resolver{
		s: s,
	}

	c := Config{
		Resolvers: resolver,
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(c))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL playground", endpoint)
}

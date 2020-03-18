//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"
	"github.com/wirelineio/wns/cmd/wnsd-lite/sync"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "9473"

// Server configures and starts the GQL server.
func Server(ctx *sync.Context) {
	if viper.GetBool("gql-server") {
		port := viper.GetString("gql-port")
		if port == "" {
			port = defaultPort
		}

		router := chi.NewRouter()

		// Add CORS middleware around every request
		// See https://github.com/rs/cors for full option listing
		router.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			Debug:          true,
		}).Handler)

		if viper.GetBool("gql-playground") {
			router.Handle("/console", handler.Playground("WNS Lite", "/graphql"))
		}

		router.Handle("/graphql", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{}})))

		err := http.ListenAndServe(":"+port, router)
		if err != nil {
			panic(err)
		}
	}
}

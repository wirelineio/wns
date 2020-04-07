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

	baseGql "github.com/wirelineio/wns/gql"
)

// Server configures and starts the GQL server.
func Server(ctx *sync.Context) {
	if !viper.GetBool("gql-server") {
		return
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          true,
	}).Handler)

	keeper := sync.NewKeeper(ctx)

	router.Handle("/api", handler.GraphQL(baseGql.NewExecutableSchema(baseGql.Config{Resolvers: &Resolver{
		Keeper: keeper,
	}})))

	// TODO(ashwin): Kept for backward compat.
	router.Handle("/graphql", handler.GraphQL(baseGql.NewExecutableSchema(baseGql.Config{Resolvers: &Resolver{
		Keeper: keeper,
	}})))

	if viper.GetBool("gql-playground") {
		router.Handle("/webui", handler.Playground("WNS Lite", "/api"))

		// TODO(ashwin): Kept for backward compat.
		router.Handle("/console", handler.Playground("WNS Lite", "/graphql"))
	}

	err := http.ListenAndServe(":"+viper.GetString("gql-port"), router)
	if err != nil {
		panic(err)
	}
}

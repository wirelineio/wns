//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/wirelineio/wns/x/nameservice"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

// Server configures and starts the GQL server.
func Server(baseApp *bam.BaseApp, cdc *codec.Codec, keeper nameservice.Keeper, accountKeeper auth.AccountKeeper) {
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

	if viper.GetBool("gql-playground") {
		router.Handle("/webui", handler.Playground("Wireline Naming Service", "/api"))

		// TODO(ashwin): Kept for backward compat.
		router.Handle("/console", handler.Playground("Wireline Naming Service", "/graphql"))
	}

	router.Handle("/api", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
		baseApp:       baseApp,
		codec:         cdc,
		keeper:        keeper,
		accountKeeper: accountKeeper,
	}})))

	// TODO(ashwin): Kept for backward compat.
	router.Handle("/graphql", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
		baseApp:       baseApp,
		codec:         cdc,
		keeper:        keeper,
		accountKeeper: accountKeeper,
	}})))

	err := http.ListenAndServe(":"+viper.GetString("gql-port"), router)
	if err != nil {
		panic(err)
	}
}

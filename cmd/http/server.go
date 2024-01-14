package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kaitsubaka/clubhub_franchises/graph"
)

const defaultPort = "8080"

func main() {

	port := viper.GetString("PORT")
	if port == "" {
		port = defaultPort
	}
	resolvers := graph.NewResolver()
	defer resolvers.ShutDown()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kaitsubaka/clubhub_franchises/graph"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common"
)

func Serve() {
	resolvers := graph.NewResolver()
	defer resolvers.ShutDown()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(common.DefaultBaseAddressWOPort, os.Getenv("PORT")), nil); err != nil {
		log.Fatal(err)
	}
}

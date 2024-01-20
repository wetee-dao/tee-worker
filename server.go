package worker

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"wetee.app/worker/graph"
	"wetee.app/worker/mint"
)

const defaultPort = "8880"

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(graph.AuthMiddleware())
	router.Use(cors.AllowAll().Handler)

	router.Handle("/", playground.Handler("Wetee-Worker", "/gql"))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  &graph.Resolver{},
		Directives: graph.NewDirectiveRoot(),
	}))
	router.Handle("/gql", srv)

	router.Post("/appLoader", mint.LoadingHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	http.ListenAndServe(":"+defaultPort, router)
}

func StartServerInCluster() {
	router := chi.NewRouter()
	router.Post("/appLoader", mint.LoadingHandler)

	log.Printf("connect to http://0.0.0.0:443 for InCluster server")
	http.ListenAndServe(":443", router)
}

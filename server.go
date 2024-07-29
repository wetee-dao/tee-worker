package worker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"wetee.app/worker/graph"
	"wetee.app/worker/util"
)

// 启动GraphQL服务器
// StartServer starts the GraphQL server.
func StartServer() {
	port := util.GetEnvInt("PORT", 8880)

	router := chi.NewRouter()
	router.Use(graph.AuthMiddleware())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Handle("/", playground.Handler("Wetee-Worker", "/gql"))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  &graph.Resolver{},
		Directives: graph.NewDirectiveRoot(),
	}))
	router.Handle("/gql", srv)

	if util.IsFileExists(util.WORK_DIR+"/ser.pem") && util.IsFileExists(util.WORK_DIR+"/ser.key") {
		log.Printf("connect to https://localhost:%s/ for GraphQL playground", fmt.Sprint(port))
		http.ListenAndServeTLS(":"+fmt.Sprint(port), util.WORK_DIR+"/ser.pem", util.WORK_DIR+"/ser.key", router)
	} else {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", fmt.Sprint(port))
		http.ListenAndServe(":"+fmt.Sprint(port), router)
	}
}

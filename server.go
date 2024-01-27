package worker

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

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
	router.Post("/appLoader/{AppID}", mint.LoadingHandler)

	log.Printf("connect to http://0.0.0.0:8883 for InCluster server")

	cert, priv := createCertificate()
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}
	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: &tlsCfg}
	server.ListenAndServeTLS("", "")
}

func createCertificate() ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "wetee-worker"},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"localhost"},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

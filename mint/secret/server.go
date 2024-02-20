package secret

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// 启动InCluster服务器
// start server in cluster for confidential
func StartSecretServerInCluster(addr string) {
	router := chi.NewRouter()
	router.Post("/appLoader/{AppID}", LoadingHandler)

	log.Printf("Start http://0.0.0.0:8883 for InCluster server")

	cert, priv, _, _ := GetRemoteReport(addr)
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

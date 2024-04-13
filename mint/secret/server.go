package secret

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wetee.app/worker/mint/proof"
)

// 启动InCluster服务器
// start server in cluster for confidential
func StartSecretServerInCluster(addr string) {
	router := chi.NewRouter()

	// Get root dcap report
	cert, priv, report, _ := proof.GetRemoteReport(addr)

	// Get root dcap report
	router.Get("/report", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"report": hex.EncodeToString(report),
		}
		bt, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	router.Post("/appInfo/{AppID}", AppInfoHandler)
	router.Post("/appLoader/{AppID}", LoadingHandler)

	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}
	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: &tlsCfg}
	log.Printf("Start http://0.0.0.0:8883 for InCluster server")
	server.ListenAndServeTLS("", "")
}

package secret

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
	wtypes "wetee.app/worker/type"
)

// 启动InCluster服务器
// start server in cluster for confidential
func StartSecretServerInCluster(addr string) {
	router := chi.NewRouter()

	router.Get("/report", func(w http.ResponseWriter, r *http.Request) {
		minter, _, _ := mint.GetMintKey()

		// Get root dcap report
		report, t, _ := proof.GetRemoteReport(minter)
		resp := wtypes.TeeParam{
			Time:    t,
			Report:  report,
			Address: minter.Address,
			Data:    nil,
		}

		// Return report
		bt, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	router.Post("/appInfo/{AppID}", AppInfoHandler)
	router.Post("/appLoader/{AppID}", LoadingHandler)

	// TODO
	cert, priv := proof.CreateCertificate(addr)
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

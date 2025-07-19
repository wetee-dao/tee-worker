package libos

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

// 启动InCluster服务器
// start server in cluster for confidential
func StartTEEServer(pk *model.PrivKey) {
	router := chi.NewRouter()
	addr := pk.GetPublic().SS58()
	signer, _ := pk.ToSigner()

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

	// Get worker tee report
	router.Get("/report", func(w http.ResponseWriter, r *http.Request) {
		report, t, err := proof.GetRemoteReport(signer, []byte{})
		if err != nil {
			util.LogWithYellow("SecretServer", "Remote REPORT", err)
			return
		}

		resp := model.TeeParam{
			Time:    t,
			Report:  report,
			Address: pk.GetPublic().Byte(),
			Data:    []byte{},
		}

		// Return report
		bt, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	// Get app info
	router.Post("/appInfo/{AppID}", AppInfoHandler)

	// launch app
	router.Post("/appLaunch/{AppID}", LoadingHandler)

	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: &tlsCfg}
	log.Printf("Start http://0.0.0.0:8883 for InCluster server")
	server.ListenAndServeTLS("", "")
}

package libos

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

func TeeServerSer(pk *model.PrivKey) (cert []byte, key []byte, der []byte, err error) {
	// int tls cert
	return util.Ed25519Cert(
		"localhost",
		[]net.IP{net.ParseIP("127.0.0.1")},
		[]string{"localhost"},
		pk.PrivateKey,
		pk.GetPublic().PublicKey,
	)
}

// 启动InCluster服务器
// start server in cluster for confidential
func StartTEEServer(pk *model.PrivKey) {
	router := chi.NewRouter()
	signer, _ := pk.ToSigner()

	// init tee server cert
	_, _, serverCertDER, err := TeeServerSer(pk)
	if err != nil {
		panic(err)
	}

	// Golang tls.config
	serverCert := tls.Certificate{Certificate: [][]byte{serverCertDER}, PrivateKey: pk.PrivateKey}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS13,

		// skip client verification
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAnyClientCert,
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

	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: tlsConfig}
	log.Printf("Start http://0.0.0.0:8883 for InCluster server")
	server.ListenAndServeTLS("", "")
}

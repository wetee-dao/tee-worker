package libos

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/go-chi/chi/v5"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
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
		resp := &model.TeeCall{
			Tx: &model.TeeCall_Text{
				Text: []byte{},
			},
		}

		err = model.IssueReport(signer, resp)
		if err != nil {
			util.LogWithYellow("SecretServer", "Remote REPORT", err)
			w.WriteHeader(500)
			w.Write([]byte("TEE report error" + err.Error()))
			return
		}

		// Return report
		buf := new(bytes.Buffer)
		err = types.WriteMessage(resp, buf)
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	})

	// Get app info
	router.Post("/appInfo/{AppID}", AppInfoHandler)

	// launch app
	router.Post("/appLaunch/{AppID}", LoadingHandler)

	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: tlsConfig}
	log.Printf("Start http://0.0.0.0:8883 for InCluster server")
	server.ListenAndServeTLS("", "")
}

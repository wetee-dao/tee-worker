package libos

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/go-chi/chi/v5"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
	"wetee.app/worker/internal/util"
)

const domain = "wetee-worker.worker-system.svc.cluster.local"

var sideChain *sidechain.SideChain

func TeeServerSer(pk *model.PrivKey) (cert []byte, key []byte, der []byte, err error) {
	// int tls cert
	return util.Ed25519Cert(
		domain,
		[]net.IP{net.ParseIP("0.0.0.0")},
		[]string{domain},
		pk.PrivateKey,
		pk.GetPublic().PublicKey,
	)
}

// 启动InCluster服务器
// start server in cluster for confidential
func StartTEEServer(pk *model.PrivKey, side *sidechain.SideChain) {
	sideChain = side
	router := chi.NewRouter()
	signer := pk.ToSigner()

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
		// require client cert
		ClientAuth: tls.RequestClientCert,
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
	// router.Post("/info/{AppID}", AppInfoHandler)

	// launch app
	router.Post("/launch/{AppID}", LoadingHandler)

	server := &http.Server{Addr: ":8883", Handler: router, TLSConfig: tlsConfig}
	util.LogWithYellow("TEEServer", "https://"+domain+":8883")
	server.ListenAndServeTLS("", "")
}

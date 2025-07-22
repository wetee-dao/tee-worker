package libos

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
)

func validateTLS(r *http.Request) error {
	if r.TLS == nil {
		return errors.New("no tls")
	}

	if len(r.TLS.PeerCertificates) == 0 {
		return errors.New("no tls")
	}

	cert := r.TLS.PeerCertificates[0]

	pub := cert.PublicKey.(ed25519.PublicKey)
	fmt.Println("Public key:", hex.EncodeToString(pub))

	// opts := x509.VerifyOptions{
	// 	Roots: certPool,
	// }

	// _, err := cert.Verify(opts)
	return nil
}

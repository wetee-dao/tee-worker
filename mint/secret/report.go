package secret

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/edgelesssys/ego/enclave"
)

var (
	// report: 远程报告
	Report []byte
	// cert: 证书
	SslCert []byte
	// priv: 私钥
	SslPriv crypto.PrivateKey
)

// 获取远程报告
// GetRemoteReport get remote report
// return: cert, priv, report, err
func GetRemoteReport(addr string) ([]byte, crypto.PrivateKey, []byte, error) {
	SslCert, SslPriv = CreateCertificate(addr)
	hash := sha256.Sum256(SslCert)
	var err error
	Report, err = enclave.GetRemoteReport(hash[:])
	return SslCert, SslPriv, Report, err
}

// 创建证书
// CreateCertificate create certificate
// return: cert, priv
func CreateCertificate(addr string) ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "wetee.app"},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{addr},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

package proof

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/edgelesssys/ego/enclave"
	"github.com/wetee-dao/go-sdk/core"
	"wetee.app/worker/util"
)

var (
	// report: 远程报告
	Report []byte
	// lastReport: 上次报告时间戳
	LastReport int64
	// cert: 证书
	SslCert []byte
	// priv: 私钥
	SslPriv crypto.PrivateKey
)

// 获取远程报告
// GetRemoteReport get remote report
// return: report, time, err
func GetRemoteReport(minter *core.Signer, data []byte) ([]byte, int64, error) {
	timestamp := time.Now().Unix()
	if Report != nil && LastReport+30 > timestamp && (data == nil || len(data) == 0) {
		return Report, LastReport, nil
	}

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(minter.PublicKey)
	if len(data) > 0 {
		buf.Write(data)
	}
	sig, err := minter.Sign(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}

	// 获取远程报告
	report, err := enclave.GetRemoteReport(sig)
	if err != nil {
		return nil, 0, err
	}

	if data == nil || len(data) == 0 {
		LastReport = timestamp
		Report = report
	}

	return report, timestamp, nil
}

// 创建证书
// CreateCertificate create certificate
// return: cert, priv
func CreateCertificate(addr string) ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "wetee.app"},
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
		DNSNames:     []string{addr},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

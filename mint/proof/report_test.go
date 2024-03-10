package proof

import "testing"

func TestCreateCertificate(t *testing.T) {
	p, _ := CreateCertificate("127.0.0.1:8080")
	p2, _ := CreateCertificate("127.0.0.1:8080")
	if string(p) == string(p2) {
		t.Error("Certificate should not be the same")
	}
}

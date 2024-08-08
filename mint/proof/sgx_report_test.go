package proof

import (
	"testing"
)

// TODO test in not sgx
func TestGetRemoteReport(t *testing.T) {
	// suite := suites.MustFind("Ed25519")
	// privateKey, _, err := types.GenerateKeyPair(suite, rand.Reader)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// bt, err := privateKey.Raw()
	// if err != nil {
	// 	t.Error(err)
	// }

	// var ed25519Key ed25519.PrivateKey = bt
	// kr, err := core.Ed25519PairFromPk(ed25519Key, 42)
	// if err != nil {
	// 	t.Error(err)
	// }

	// report, _, err := GetRemoteReport(&kr)
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(report)
}

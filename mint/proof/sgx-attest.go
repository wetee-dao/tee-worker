package proof

import (
	"fmt"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
)

func VerifyReportProof(reportBytes, certBytes, signer []byte) (*attestation.Report, error) {
	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return nil, err
	}

	// hash := sha256.Sum256(certBytes)
	// if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
	// 	return nil, errors.New("report data does not match the certificate's hash")
	// }

	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).

	fmt.Println("report.SecurityVersion", report.SecurityVersion)
	// if report.SecurityVersion < 2 {
	// 	return nil, errors.New("invalid security version")
	// }
	// if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
	// 	return nil, errors.New("invalid product")
	// }
	// if !bytes.Equal(report.SignerID, signer) {
	// 	return nil, errors.New("invalid signer")
	// }

	// For production, you must also verify that report.Debug == false

	return &report, nil
}

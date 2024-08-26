package proof

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	"wetee.app/worker/util"

	wtypes "wetee.app/worker/type"
)

func VerifyReportProof(workerReport *wtypes.TeeParam) (*wtypes.TeeReport, error) {
	// TODO SEV/TDX not support
	if workerReport.TeeType != 0 {
		return &wtypes.TeeReport{
			CodeSignature: []byte{},
			CodeSigner:    []byte{},
			CodeProductID: []byte{},
		}, nil
	}

	var reportBytes, msgBytes, timestamp = workerReport.Report, workerReport.Data, workerReport.Time

	// decode address
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	// 检查时间戳，超过 30s 签名过期
	if timestamp+30 < time.Now().Unix() {
		return nil, errors.New("report expired")
	}

	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	} else if err != nil {
		return nil, err
	}

	pubkey, err := ed25519.Scheme{}.FromPublicKey(signer)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(signer)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data

	if !pubkey.Verify(buf.Bytes(), sig) {
		return nil, errors.New("invalid sgx report")
	}

	if report.Debug {
		return nil, errors.New("debug mode is not allowed")
	}

	return &wtypes.TeeReport{
		TeeType:       workerReport.TeeType,
		CodeSigner:    report.SignerID,
		CodeSignature: report.UniqueID,
		CodeProductID: report.ProductID,
	}, nil
}

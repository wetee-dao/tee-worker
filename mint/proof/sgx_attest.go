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

func VerifyReportProof(reportBytes, msgBytes, signer []byte, timestamp int64) (*attestation.Report, error) {
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

	return &report, nil
}

func VerifyReportFromTeeParam(workerReport *wtypes.TeeParam) (*attestation.Report, error) {
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return nil, errors.New("VerifyReportFromTeeParam: SS58 decode")
	}

	return VerifyReportProof(workerReport.Report, workerReport.Data, signer, workerReport.Time)
}

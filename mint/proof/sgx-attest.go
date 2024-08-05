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
	"wetee.app/worker/internal/store"
	"wetee.app/worker/util"
)

func VerifyReportProof(reportBytes, msgBytes, signer []byte) (*attestation.Report, error) {
	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return nil, err
	}

	sig := report.Data
	pubkey, err := ed25519.Scheme{}.FromPublicKey(signer)
	if err != nil {
		return nil, err
	}

	if !pubkey.Verify(msgBytes, sig) {
		return nil, errors.New("invalid sgx report")
	}

	if report.Debug {
		return nil, errors.New("debug mode is not allowed")
	}

	return &report, nil
}

func VerifyReportFromTeeParam(workerReport *store.TeeParam) (*attestation.Report, error) {
	// decode time
	timestamp := workerReport.Time

	// 检查时间戳，超过 30s 签名过期
	if timestamp+30 < time.Now().Unix() {
		return nil, errors.New("Report expired")
	}

	// decode report
	report := workerReport.Report

	// decode address
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return nil, errors.New("VerifyReportFromTeeParam: SS58 decode")
	}

	// 构建验证数据
	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(signer)

	return VerifyReportProof(report, buf.Bytes(), signer)
}

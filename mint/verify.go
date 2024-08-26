package mint

import (
	"encoding/hex"
	"errors"

	"github.com/vedhavyas/go-subkey/v2"
	"github.com/wetee-dao/go-sdk/module"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/mint/proof"
	types "wetee.app/worker/type"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyWorker(reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := proof.VerifyReportProof(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := module.GetWorkerCode(m.ChainClient)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyDsecret(reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := proof.VerifyReportProof(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := module.GetDsecretCode(m.ChainClient)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyWorkLibos(wid gtypes.WorkId, reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := proof.VerifyReportProof(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := module.GetWorkCode(m.ChainClient, wid)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}

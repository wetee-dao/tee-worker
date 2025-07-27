package mint

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyWorker(reportData *model.TeeCall) ([]byte, error) {
	// 解码地址
	signer := reportData.Caller

	// report, err := proof.VerifyReportProof(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// 校验 worker 代码版本
	// codeHash, codeSigner, err := module.GetWorkerCode(m.ChainClient)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyDsecret(reportData *model.TeeCall) ([]byte, error) {
	// 解码地址
	signer := reportData.Caller

	// report, err := proof.VerifyReportProof(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// // 校验 worker 代码版本
	// codeHash, codeSigner, err := module.GetDsecretCode(m.ChainClient)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (m *Minter) VerifyWorkLibos(podid uint64, reportData *model.TeeCall) ([]byte, error) {
	// 解码地址
	signer := reportData.Caller

	// report, err := proof.VerifyReportProof(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// // 校验 worker 代码版本
	// codeHash, codeSigner, err := module.GetWorkCode(m.ChainClient, wid)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}

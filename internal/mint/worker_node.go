package mint

import (
	chain "github.com/wetee-dao/ink.go"
)

// RegisterNode register node
// 注册节点
func (c *Minter) RegisterNode(signer *chain.Signer, pubkey []byte) error {
	// var bt [32]byte
	// copy(bt[:], pubkey)

	// call := dsecret.MakeRegisterNodeCall(bt)
	// return c.ChainClient.SignAndSubmit(signer, call, true)
	return nil
}

// 获取全网当前程序的代码版本
// Get CodeSignature
// func (c *Minter) GetCodeSignature() ([]byte, error) {
// 	return dsecret.GetCodeSignatureLatest(c.ChainClient.Api.RPC.State)
// }

// 获取全网当前程序的签名人
// Get GetCodeSigner
// func (c *Minter) GetGetCodeSigner() ([]byte, error) {
// 	return dsecret.GetCodeSignerLatest(c.ChainClient.Api.RPC.State)
// }

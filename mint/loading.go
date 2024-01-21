package mint

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
	"wetee.app/worker/dao"
)

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	// 验证 AppID
	appID := chi.URLParam(r, "AppID")
	wid, err := dao.UnSealAppID(appID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("AppID error: " + err.Error()))
		return
	}

	// 获取数据
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Read body error: " + err.Error()))
		return
	}
	param := &dao.LoadParam{}
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Request body unmarshal error: " + err.Error()))
		return
	}

	// 验证消息
	// 解析地址
	_, pubkeyBytes, err := subkey.SS58Decode(param.Address)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Address decode error: " + err.Error()))
		return
	}

	// 解析公钥
	pubkey, err := sr25519.Scheme{}.FromPublicKey(pubkeyBytes)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Pubkey error: " + err.Error()))
		return
	}

	// 验证签名
	sig, err := hex.DecodeString(param.Signature)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Signature decode error: " + err.Error()))
		return
	}
	ok := pubkey.Verify([]byte(param.Time), sig)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("Signature error"))
		return
	}

	// 验证地址
	address, err := dao.GetSetAppSignerAddress(wid, param.Address)
	if err != nil || address != param.Address {
		w.WriteHeader(500)
		w.Write([]byte("Address error: " + err.Error()))
		return
	}

	s, err := dao.GetSecrets(wid)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Secret error: " + err.Error()))
		return
	}

	bt, _ := json.Marshal(s)
	w.WriteHeader(200)
	w.Write(bt)
}

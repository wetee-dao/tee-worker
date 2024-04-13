package secret

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
	"github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/store"
)

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	// 验证 AppID
	appID := chi.URLParam(r, "AppID")

	// 获取数据
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Read body error" + err.Error()))
		return
	}

	// 解析请求数据
	param := &store.LoadParam{}
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Request body unmarshal error" + err.Error()))
		return
	}

	// 加载应用的加密环境变量和文件
	s, err := loading(appID, param)
	fmt.Println("loading", err)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	bt, _ := json.Marshal(s)
	w.WriteHeader(200)
	w.Write(bt)
}

func loading(appID string, param *store.LoadParam) (*store.Secrets, error) {
	wid, err := VerifyLibOs(appID, param)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLibOs error")
	}
	// 验证地址
	address, err := store.GetSetAppSignerAddress(*wid, param.Address)
	if err != nil || address != param.Address {
		// return nil, errors.New("VerifyAddress error")
	}

	_, err = proof.VerifyReportProof(param.Report, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLocalReport error")
	}

	// 存入Work DCAP信息
	err = store.SetWorkDcapReport(*wid, param.Report)
	if err != nil {
		return nil, errors.Wrap(err, "DCAP Report set error")
	}

	// 获取加密信息
	s, err := store.GetSecrets(*wid)
	if err != nil {
		return nil, errors.Wrap(err, "Secret error")
	}

	return s, nil
}

func VerifyLibOs(appID string, param *store.LoadParam) (*types.WorkId, error) {
	wid, err := store.UnSealAppID(appID)
	if err != nil {
		return nil, errors.Wrap(err, "AppID error")
	}

	// 验证消息
	// 解析地址
	_, pubkeyBytes, err := subkey.SS58Decode(param.Address)
	if err != nil {
		return nil, errors.Wrap(err, "Address decode error")
	}

	// 解析公钥
	pubkey, err := sr25519.Scheme{}.FromPublicKey(pubkeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Pubkey error")
	}

	// 验证签名
	sig, err := hex.DecodeString(param.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "Signature decode error")
	}
	ok := pubkey.Verify([]byte(param.Time), sig)
	if !ok {
		return nil, errors.New("Signature error")
	}

	return &wid, nil
}

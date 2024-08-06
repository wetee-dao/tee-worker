package secret

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/wetee-dao/go-sdk/module"
	"wetee.app/worker/mint"
	wtypes "wetee.app/worker/type"
)

// load app info
// 获取应用消息
func AppInfoHandler(w http.ResponseWriter, r *http.Request) {
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
	param := &wtypes.TeeParam{}
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Request body unmarshal error" + err.Error()))
		return
	}

	// 获取数据
	s, err := GetAppInfo(appID, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("GetAppInfo error" + err.Error()))
		return
	}

	bt, _ := json.Marshal(s)
	w.WriteHeader(200)
	w.Write(bt)
}

// 获取应用消息
// get app info
func GetAppInfo(appID string, param *wtypes.TeeParam) (map[string]string, error) {
	// 验证 report
	wid, err := VerifyLibOs(appID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLibOs error")
	}

	user, err := module.GetAccount(mint.MinterIns.ChainClient, *wid)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"user": hex.EncodeToString(user),
	}, nil
}

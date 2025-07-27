package libos

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// load app info
// 获取应用消息
func AppInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 验证 AppID
	appId := chi.URLParam(r, "AppID")

	// 获取数据
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Read body error" + err.Error()))
		return
	}

	// 解析请求数据
	param := &model.TeeCall{}
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Request body unmarshal error" + err.Error()))
		return
	}

	// 获取数据
	s, err := GetAppInfo(appId, param)
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
func GetAppInfo(appId string, param *model.TeeCall) (map[string]string, error) {
	// 验证 report
	// wid, err := VerifyLibOs(appId, nil)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "VerifyLibOs error")
	// }

	// user, err := module.GetAccount(mint.MinterIns.ChainClient, *wid)
	// if err != nil {
	// 	return nil, err
	// }
	user := []byte{}

	return map[string]string{
		"user": hex.EncodeToString(user),
	}, nil
}

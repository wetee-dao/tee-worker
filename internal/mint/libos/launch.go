package libos

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
	"wetee.app/worker/internal/mint"
	"wetee.app/worker/internal/store"
)

// 加载应用加密文件，加密环境变量
// load app secret file and env
func LoadingHandler(w http.ResponseWriter, r *http.Request) {
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
	param := &model.TeeParam{}
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Request body unmarshal error" + err.Error()))
		return
	}

	// 加载应用的加密环境变量和文件
	s, err := loading(appId, param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	bt, _ := json.Marshal(s)
	w.WriteHeader(200)
	w.Write(bt)
}

// 加载应用加密文件，加密环境变量
// load app secret file and env
func loading(appId string, param *model.TeeParam) (*store.EnvWrap, error) {
	// 验证 libos 完整性信息
	podId, err := VerifyLibOs(appId, param)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLibOs error")
	}

	// 存入 Work DCAP 信息
	err = store.SetPendingTEEReport(podId, param)
	if err != nil {
		return nil, errors.Wrap(err, "DCAP Report set error")
	}

	util.LogWithBlue("LOADED POD", podId)
	err = mint.MinterIns.AddPendingDeployTx(podId)
	if err != nil {
		return nil, errors.Wrap(err, "AddPendingDeployTx error")
	}

	// 获取配置文件
	// 获取加密配置文件
	s := &store.EnvWrap{
		// Sec: *secret,
	}

	return s, nil
}

// VerifyLibOs 函数验证应用程序标识和报告，并返回工作标识或错误
func VerifyLibOs(appId string, report *model.TeeParam) (uint64, error) {
	// 解包应用程序标识
	id, appTime, err := store.UnSealAppID(appId)
	if err != nil {
		return 0, errors.Wrap(err, "AppID error")
	}

	// 应用的状态只能60秒内使用
	if time.Now().Unix() > appTime+60 {
		// return 0, errors.New("AppID time error")
	}

	// 验证工作标识和报告
	_, err = mint.MinterIns.VerifyWorkLibos(id, report)
	if err != nil {
		return 0, errors.Wrap(err, "VerifyWorkLibos error")
	}

	// 返回解包后的工作标识
	return id, nil
}

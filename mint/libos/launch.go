package libos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/mint"
	wtypes "wetee.app/worker/type"
)

// 加载应用加密文件，加密环境变量
// load app secret file and env
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
	param := &wtypes.TeeParam{}
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

// 加载应用加密文件，加密环境变量
// load app secret file and env
func loading(appID string, param *wtypes.TeeParam) (*wtypes.EnvWrap, error) {
	// 验证 libos 完整性信息
	wid, err := VerifyLibOs(appID, param)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLibOs error")
	}

	// 存入 Work DCAP 信息
	bt, _ := json.Marshal(param)
	err = store.SetWorkDcapReport(*wid, bt)
	if err != nil {
		return nil, errors.Wrap(err, "DCAP Report set error")
	}

	// 上传TEE环境变量，设置当前的部署 Key
	secret, err := mint.MinterIns.LaunchFromDsecret(wid, param)
	if err != nil {
		return nil, errors.Wrap(err, "LaunchFromDsecret error")
	}

	// 获取配置文件
	// 获取加密配置文件
	s := &wtypes.EnvWrap{
		Sec: *secret,
	}

	return s, nil
}

// VerifyLibOs 函数验证应用程序标识和报告，并返回工作标识或错误
func VerifyLibOs(appID string, report *wtypes.TeeParam) (*types.WorkId, error) {
	// 解包应用程序标识
	wid, err := store.UnSealAppID(appID)
	if err != nil {
		// 如果解包过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "AppID error")
	}

	// 验证工作标识和报告
	_, err = mint.MinterIns.VerifyWorkLibos(wid, report)
	if err != nil {
		// 如果验证过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "VerifyWorkLibos error")
	}

	// 返回解包后的工作标识
	return &wid, nil
}

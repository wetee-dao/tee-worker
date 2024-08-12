package secret

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/edgelesssys/ego/attestation"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
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
func loading(appID string, param *wtypes.TeeParam) (*store.Secrets, error) {
	// 验证报告是否合理
	report, err := proof.VerifyReportFromTeeParam(param)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLocalReport error")
	}

	// 验证 libos 完整性信息
	wid, err := VerifyLibOs(appID, report)
	if err != nil {
		return nil, errors.Wrap(err, "VerifyLibOs error")
	}

	_, deployAccount, err := subkey.SS58Decode(param.Address)
	if err != nil {
		return nil, errors.Wrap(err, "Address error")
	}

	// 存入 Work DCAP 信息
	err = store.SetWorkDcapReport(*wid, param.Report)
	if err != nil {
		return nil, errors.Wrap(err, "DCAP Report set error")
	}

	err = store.SetWorkDeploy(*wid, deployAccount)
	if err != nil {
		return nil, errors.Wrap(err, "Set deploy error")
	}

	mint.MinterIns.Addlanch(*wid)

	// deploy := param.Address

	// 获取配置文件
	// 获取加密配置文件
	s := &store.Secrets{}

	return s, nil
}

func VerifyLibOs(appID string, report *attestation.Report) (*types.WorkId, error) {
	wid, err := store.UnSealAppID(appID)
	if err != nil {
		return nil, errors.Wrap(err, "AppID error")
	}

	//TODO 验证程序的版本和链上程序的版本是否一样
	return &wid, nil
}

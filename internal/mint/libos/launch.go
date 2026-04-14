package libos

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
	"github.com/wetee-dao/tee-dsecret/pkg/util"

	"wetee.app/worker/internal/mint"
	"wetee.app/worker/internal/store"
)

// 加载应用加密文件，加密环境变量
// load app secret file and env
func (s *TEEServer) launch(req []byte) ([]byte, error) {
	// 解析请求数据
	param := &model.TeeCall{}
	err := protoio.ReadMessage(bytes.NewBuffer(req), param)
	if err != nil {
		return nil, errors.Wrap(err, "ReadMessage error")
	}

	var podId uint64 = 0
	startReq := &model.PodStart{}
	switch call := param.Tx.(type) {
	case *model.TeeCall_PodStart:
		podId, err = VerifyLibOs(string(call.PodStart.AppId), param)
		if err != nil {
			return nil, errors.Wrap(err, "VerifyLibOs error")
		}

		if call.PodStart.Id != podId {
			return nil, errors.New("podid is not match tee call")
		}
		startReq = call.PodStart
	default:
		return nil, errors.New("Tx is not pod start")
	}

	secrets, err := s.side.BroadcastReencryptReq(startReq)
	if err != nil {
		return nil, errors.Wrap(err, "BroadcastDecryptSecret error")
	}

	// 存入 Work DCAP 信息
	err = store.SetPendingTEEReport(podId, param)
	if err != nil {
		return nil, errors.Wrap(err, "DCAP Report set error")
	}

	// sync to chain
	util.LogWithGray("PRELOADED POD", podId)
	err = mint.MinterIns.AddPendingDeployTx(param)
	if err != nil {
		util.LogWithRed("PRELOADED POD EROR", err)
		return nil, errors.Wrap(err, "AddPendingDeployTx error")
	}

	// 写入解密消息
	buf := new(bytes.Buffer)
	protoio.WriteMessage(secrets, buf)

	return buf.Bytes(), nil
}

// VerifyLibOs 函数验证应用程序标识和报告，并返回工作标识或错误
func VerifyLibOs(appId string, report *model.TeeCall) (uint64, error) {
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

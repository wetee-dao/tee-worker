package libos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/panjf2000/gnet/v2"
	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
	"wetee.app/worker/internal/util"
)

// const domain = "wetee-worker.worker-system.svc.cluster.local"
// var sideChain *sidechain.SideChain
// func TeeServerSer(pk *model.PrivKey) (cert []byte, key []byte, der []byte, err error) {
// 	// int tls cert
// 	return util.Ed25519Cert(
// 		domain,
// 		[]net.IP{net.ParseIP("0.0.0.0")},
// 		[]string{domain},
// 		pk.PrivateKey,
// 		pk.GetPublic().PublicKey,
// 	)
// }

type TEEServer struct {
	*gnet.BuiltinEventEngine
	eng       gnet.Engine
	network   string
	addr      string
	multicore bool

	pk   *model.PrivKey
	side *sidechain.SideChain
}

func (s *TEEServer) OnBoot(srv gnet.Engine) (action gnet.Action) {
	s.eng = srv
	util.LogWithYellow("TEEServer", "listening on 8883")
	return
}

func (s *TEEServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	report, err := s.report()
	if err != nil {
		out = s.wrapData(0, 500, []byte("OnOpen failed:"+err.Error()))
	} else {
		out = s.wrapData(0, 0, report)
	}

	return
}

func (s *TEEServer) OnClose(c gnet.Conn, err error) gnet.Action {
	return gnet.Close
}

func (s *TEEServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	id, packet, err := Decode(c)
	if err != nil {
		if errors.Is(err, io.ErrShortBuffer) {
			return gnet.None
		}
		util.LogWithRed("TEEServer", "Decode failed:", err)
		return s.ReturnError(c, id, err)
	}

	// 解析请求
	req := new(model.ApiReq)
	err = protoio.ReadMessage(bytes.NewBuffer(packet), req)
	if err != nil {
		util.LogWithRed("TEEServer", "ReadMessage failed:", err)
		return s.ReturnError(c, id, err)
	}

	// 处理请求
	switch string(req.Url) {
	case "/launch":
		data, err := s.launch(req.Data)
		if err != nil {
			return s.ReturnError(c, id, errors.Wrap(err, "call launch failed"))
		}

		return s.ReturnData(c, id, data)
	default:
		err = fmt.Errorf("unknown URL: %s", string(req.Url))
		return s.ReturnError(c, id, err)
	}
}

// ReturnError 返回错误
func (s *TEEServer) ReturnError(c gnet.Conn, req uint64, err error) gnet.Action {
	bt := s.wrapData(req, 500, []byte(err.Error()))
	c.Write(bt)
	return gnet.None
}

// ReturnData 返回数据
func (s *TEEServer) ReturnData(c gnet.Conn, req uint64, body []byte) gnet.Action {
	bt := s.wrapData(req, 0, body)
	c.Write(bt)
	return gnet.None
}

// wrapData 包装数据
func (s *TEEServer) wrapData(req uint64, code int32, data []byte) []byte {
	result := &model.ApiResp{
		Code: code,
		Data: data,
	}

	buf := new(bytes.Buffer)
	abci.WriteMessage(result, buf)
	bt, _ := Encode(req, buf.Bytes())
	return bt
}

// Stop 停止服务器
func (s *TEEServer) Stop() error {
	return s.eng.Stop(context.Background())
}

// 启动InCluster服务器
// start server in cluster for confidential
func StartTEEServer(pk *model.PrivKey, side *sidechain.SideChain) {
	port := 8883
	multicore := true

	s := &TEEServer{
		network:   "tcp",
		addr:      fmt.Sprintf("0.0.0.0:%d", port),
		multicore: multicore,
		pk:        pk,
		side:      side,
	}

	err := gnet.Run(s, s.network+"://"+s.addr, gnet.WithMulticore(multicore))
	if err != nil {
		util.LogWithRed("TEEServer", "Run failed:", err)
		os.Exit(1)
	}
}

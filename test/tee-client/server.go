package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/panjf2000/gnet/v2"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
	"wetee.app/worker/internal/util"
)

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
	fmt.Printf("Server is listening on %d \n", 8883)
	return
}

func (s *TEEServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	out = []byte("sweetness\n")
	return
}

func (s *TEEServer) OnClose(c gnet.Conn, err error) gnet.Action {
	return gnet.Shutdown
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

	req := new(model.ApiReq)
	err = protoio.ReadMessage(bytes.NewBuffer(packet), req)
	if err != nil {
		util.LogWithRed("TEEServer", "ReadMessage failed: %v", err)
		return s.ReturnError(c, id, err)
	}

	switch string(req.Url) {
	case "/report":
		return s.ReturnData(c, id, []byte("report"))
	case "/launch":
		return s.ReturnData(c, id, []byte("launch"))
	default:
		err = fmt.Errorf("unknown URL: %s", string(req.Url))
		return s.ReturnError(c, id, err)
	}
}

// ReturnError 返回错误
func (s *TEEServer) ReturnError(c gnet.Conn, req uint64, err error) gnet.Action {
	result := &model.ApiResp{
		Code: 500,
		Data: []byte(err.Error()),
	}

	// 写入解密消息
	buf := new(bytes.Buffer)
	abci.WriteMessage(result, buf)
	bt, _ := Encode(req, buf.Bytes())
	c.Write(bt)
	return gnet.None
}

// ReturnData 返回数据
func (s *TEEServer) ReturnData(c gnet.Conn, req uint64, body []byte) gnet.Action {
	result := &model.ApiResp{
		Code: 0,
		Data: body,
	}

	// 写入解密消息
	buf := new(bytes.Buffer)
	abci.WriteMessage(result, buf)
	bt, _ := Encode(req, buf.Bytes())
	c.Write(bt)
	return gnet.None
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
		addr:      fmt.Sprintf(":%d", port),
		multicore: multicore,
		pk:        pk,
		side:      side,
	}

	gnet.Run(s, s.network+"://"+s.addr, gnet.WithMulticore(multicore))
}

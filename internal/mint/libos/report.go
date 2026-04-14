package libos

import (
	"bytes"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
	"wetee.app/worker/internal/util"
)

func (s *TEEServer) report() ([]byte, error) {
	resp := &model.TeeCall{
		Tx: &model.TeeCall_Text{
			Text: []byte{},
		},
	}

	err := model.IssueReport(s.pk.ToSigner(), resp)
	if err != nil {
		util.LogWithYellow("SecretServer", "Remote REPORT", err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = protoio.WriteMessage(resp, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

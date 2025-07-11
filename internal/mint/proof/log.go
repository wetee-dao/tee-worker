package proof

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"
)

// 日志证明
type WorkLogProof struct {
	BlockNumber uint64
	Time        uint64
	Logs        []string
}

var LogBucket = "log"

// 工作量日志列表
// Work Log List
func ListLogsById(podId uint64, page int, size int, isCache bool) ([]WorkLogProof, error) {
	name := fmt.Sprint(podId)
	if isCache {
		name = name + "_cache"
	}

	res, err := model.GetList(LogBucket, name, page, size)
	if err != nil {
		return nil, err
	}

	var list = make([]WorkLogProof, 0, len(res))
	for _, v := range res {
		logProof := WorkLogProof{}
		err = json.Unmarshal(v, &logProof)
		if err != nil {
			return nil, err
		}
		list = append(list, logProof)
	}

	return list, nil
}

// 工作量日志 hash
func GetWorkLogHash(log []string, blockNumber uint64) ([]byte, []byte, error) {
	pf := WorkLogProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Logs:        log,
	}
	bt, err := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	return hash[:], bt, err
}

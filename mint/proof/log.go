package proof

import (
	"encoding/json"
	"fmt"
	"time"

	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"golang.org/x/crypto/blake2b"

	"wetee.app/worker/store"
	"wetee.app/worker/util"
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
func ListLogsById(id gtypes.WorkId, page int, size int) ([]WorkLogProof, error) {
	name := LogBucket + util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id)
	res, err := store.GetList(LogBucket, []byte(name), page, size)
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

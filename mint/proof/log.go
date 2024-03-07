package proof

import (
	"encoding/json"
	"fmt"
	"time"

	gtypes "github.com/wetee-dao/go-sdk/gen/types"
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

func ListLogsById(id gtypes.WorkId, page int, size int) ([]WorkLogProof, error) {
	name := util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id)
	res, err := store.GetLogList([]byte(name), page, size)
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
func GetWorkLogHash(id string, log []string, blockNumber uint64) ([]byte, error) {
	pf := WorkLogProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Logs:        log,
	}
	bt, _ := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	err := store.Addlog([]byte(id), bt)
	return hash[:], err
}

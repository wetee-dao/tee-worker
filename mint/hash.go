package mint

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/blake2b"
	"wetee.app/worker/dao"
)

// TODO 工作量证明查询
func getWorkLogHash(name string, log []string, blockNumber uint64) ([]byte, error) {
	pf := WorkLogProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Logs:        log,
	}
	bt, _ := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	err := dao.Addlog([]byte(name), bt)
	return hash[:], err
}

// TODO 工作量证明查询
func getWorkCrHash(name string, cr map[string][]int64, blockNumber uint64) ([]byte, []uint32, error) {
	pf := WorkCrProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Cr:          cr,
	}
	bt, _ := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	err := dao.Addlog([]byte(name), bt)

	crA := []uint32{0, 0}
	for _, v := range cr {
		crA[0] += uint32(v[0])
		crA[1] += uint32(v[1])
	}
	return hash[:], crA, err
}

// 日志证明
type WorkLogProof struct {
	BlockNumber uint64
	Time        uint64
	Logs        []string
}

// 硬件资源证明
type WorkCrProof struct {
	BlockNumber uint64
	Time        uint64
	Cr          map[string][]int64
}

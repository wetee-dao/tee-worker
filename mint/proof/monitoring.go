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

// 硬件资源证明
type WorkCrProof struct {
	BlockNumber uint64
	Time        uint64
	Cr          map[string][]int64
}

var CrBucket = "cr"

// 工作量证明资源占用列表
func ListMonitoringsById(id gtypes.WorkId, page int, size int) ([]WorkCrProof, error) {
	name := CrBucket + util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id)
	res, err := store.GetList(CrBucket, []byte(name), page, size)
	if err != nil {
		return nil, err
	}

	var list = make([]WorkCrProof, 0, len(res))
	for _, v := range res {
		proof := WorkCrProof{}
		err = json.Unmarshal(v, &proof)
		if err != nil {
			return nil, err
		}
		list = append(list, proof)
	}
	return list, nil
}

// 工作量证明资源占用 hash
func GetWorkCrHash(cr map[string][]int64, blockNumber uint64) ([]byte, []uint32, []byte, error) {
	pf := WorkCrProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Cr:          cr,
	}
	bt, err := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	crA := []uint32{0, 0}
	for _, v := range cr {
		crA[0] += uint32(v[0])
		crA[1] += uint32(v[1])
	}
	return hash[:], crA, bt, err
}

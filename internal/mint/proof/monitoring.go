package proof

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"
)

// 硬件资源证明
type WorkCrProof struct {
	Time uint64
	Cr   map[string][]int64
}

var CrBucket = "cr"

// 工作量证明资源占用列表
func ListMonitoringsById(id uint64, page int, size int, isCache bool) ([]WorkCrProof, error) {
	name := fmt.Sprint(id)
	if isCache {
		name = name + "_cache"
	}
	res, err := model.GetList(CrBucket, name, page, size)
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

func AddMonitor(pod model.Pod, logs []string, crs map[string][]int64) error {
	name := fmt.Sprint(pod.PodId)

	_, bt, err := GetLogHash(logs)
	if err != nil {
		return err
	}

	err = model.AddToList(LogBucket, name, bt)
	if err != nil {
		return err
	}

	_, _, cbt, err := GetWorkCrHash(crs)
	if err != nil {
		return err
	}

	return model.AddToList(CrBucket, name, cbt)
}

// 工作量证明资源占用 hash
func GetWorkCrHash(cr map[string][]int64) ([]byte, []uint32, []byte, error) {
	pf := WorkCrProof{
		Time: uint64(time.Now().Unix()),
		Cr:   cr,
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

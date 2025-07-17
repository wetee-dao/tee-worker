package proof

import (
	"errors"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/utility"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"

	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

func MakeWorkProof(pod model.Pod, logs []string, crs map[string][]int64, now time.Time) (*types.Call, int64, error) {
	name := fmt.Sprint(pod.PodId)

	// 获取log和硬件资源使用量
	var logHash = []byte{}
	var crHash = []byte{}
	var cr = []uint32{0, 0, 0}
	var err error

	err = store.SetCacheId(name, now.Unix())
	if err != nil {
		util.LogError("SetCacheId", err)
		return nil, 0, err
	}

	err = model.DeleteList(LogBucket, name+"_cache")
	if err != nil {
		util.LogError("DeleteLog", err)
		return nil, 0, err
	}

	if len(logs) > 0 {
		// 获取log hash
		// Get log hash
		var bt []byte
		logHash, bt, err = GetWorkLogHash(logs, uint64(pod.LastMintBlockNumber))
		if err != nil {
			util.LogError("getWorkLogHash", err)
			return nil, 0, err
		}
		err = model.AddToList(LogBucket, name, bt)
		if err != nil {
			util.LogError("Addlog", err)
			return nil, 0, err
		}
	}

	err = model.DeleteList(CrBucket, name+"_cache")
	if err != nil {
		util.LogError("DeleteLog", err)
		return nil, 0, err
	}

	if len(crs) > 0 {
		// 获取计算资源hash
		// Get Computing resource hash
		var bt []byte
		crHash, cr, bt, err = GetWorkCrHash(crs, uint64(pod.LastMintBlockNumber))
		if err != nil {
			util.LogError("getWorkCrHash", err)
			return nil, 0, err
		}
		err := model.AddToList(CrBucket, name, bt)
		if err != nil {
			util.LogError("AddCr", err)
			return nil, 0, err
		}
	}

	crProof := gtypes.ComCr{
		Cpu:  cr[0],
		Mem:  cr[1],
		Disk: 0,
	}

	hasHash := false
	if len(logHash) > 0 || len(crHash) > 0 {
		hasHash = true
	}

	fmt.Println("MakeWorkProof ========> ", crProof, hasHash)

	// 获取工作证明
	// Get report of work
	report := [32]byte{}
	reportData, err := store.GetWorkDcapReport(pod.PodId)
	if err != nil {
		// util.LogError("GetWorkDcapReport", err)
	} else {
		hash := blake2b.Sum256(reportData)
		report = hash
	}

	// 所有需要提交的信息都不存在，不继续提交
	// All required submission information is missing, and the submission will not be continued.
	if report == [32]byte{} && crHash == nil && logHash == nil {
		return nil, 0, errors.New("report, crHash and logHash are all nil")
	}

	key, err := model.GetKey("G", "dkg_pub_key")
	if err != nil {
		return nil, 0, errors.New("get G-dkg_pub_key error")
	}
	account, err := types.NewAccountID(key)
	if err != nil {
		return nil, 0, errors.New("get G-dkg_pub_key error")
	}

	err = chains.MainChain.DryStartPod(pod.PodId, types.H256(report), *account)
	if err != nil {
		util.LogError("DryStartPod", err)
		return nil, 0, err
	}

	t := time.Now().UnixMilli()
	call, err := chains.MainChain.TxCallOfStartPod(pod.PodId, types.H256(report), *account)
	if err != nil {
		util.LogError("TxCallOfStartPod", err)
		return nil, 0, err
	}

	return call, t, nil
}

func SubmitWorkProof(client *chain.ChainClient, signer *chain.Signer, proof []gtypes.RuntimeCall) error {
	runtimeCall := utility.MakeBatchCall(proof)
	call, _ := (runtimeCall).AsCall()
	return client.SignAndSubmit(signer, call, true, 0)
}

func CacheWorkProof(podId uint64, logs []string, crs map[string][]int64, now time.Time, BlockNumber uint64) error {
	name := fmt.Sprint(podId)

	// 获取log和硬件资源使用量
	var err error

	err = store.SetCacheId(name+"-cache", now.Unix())
	if err != nil {
		util.LogError("SetCacheId", err)
		return err
	}

	if len(logs) > 0 {
		// 获取 log hash
		// Get log hash
		var bt []byte
		_, bt, err = GetWorkLogHash(logs, BlockNumber)
		if err != nil {
			util.LogError("getWorkLogHash", err)
			return err
		}
		err = model.AddToList(LogBucket, name+"_cache", bt)
		if err != nil {
			util.LogError("Addlog", err)
			return err
		}
	}

	if len(crs) > 0 {
		// 获取计算资源hash
		// Get Computing resource hash
		var bt []byte
		_, _, bt, err = GetWorkCrHash(crs, BlockNumber)
		if err != nil {
			util.LogError("getWorkCrHash", err)
			return err
		}
		err := model.AddToList(CrBucket, name+"_cache", bt)
		if err != nil {
			util.LogError("AddCr", err)
			return err
		}
	}

	fmt.Println("CACHE TEE PROOF ========> ", len(logs), len(crs))

	return nil
}

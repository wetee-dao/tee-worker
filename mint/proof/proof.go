package proof

import (
	"errors"
	"fmt"
	"time"

	chain "github.com/wetee-dao/go-sdk"
	"golang.org/x/crypto/blake2b"

	"github.com/wetee-dao/go-sdk/core"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/utility"
	"github.com/wetee-dao/go-sdk/pallet/weteeworker"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/util"
)

func MakeWorkProof(wid gtypes.WorkId, logs []string, crs map[string][]int64, now time.Time, BlockNumber uint64) (*gtypes.RuntimeCall, error) {
	name := util.GetWorkTypeStr(wid) + "-" + fmt.Sprint(wid.Id)

	// 获取log和硬件资源使用量
	var logHash = []byte{}
	var crHash = []byte{}
	var cr = []uint32{0, 0, 0}
	var err error

	err = store.SetCacheId(name, now.Unix())
	if err != nil {
		util.LogError("SetCacheId", err)
		return nil, err
	}

	err = store.DeleteList(LogBucket, name+"_cache")
	if err != nil {
		util.LogError("DeleteLog", err)
		return nil, err
	}
	if len(logs) > 0 {
		// 获取log hash
		// Get log hash
		var bt []byte
		logHash, bt, err = GetWorkLogHash(logs, BlockNumber)
		if err != nil {
			util.LogError("getWorkLogHash", err)
			return nil, err
		}
		err = store.AddToList(LogBucket, name, bt)
		if err != nil {
			util.LogError("Addlog", err)
			return nil, err
		}
	}

	err = store.DeleteList(CrBucket, name+"_cache")
	if err != nil {
		util.LogError("DeleteLog", err)
		return nil, err
	}
	if len(crs) > 0 {
		// 获取计算资源hash
		// Get Computing resource hash
		var bt []byte
		crHash, cr, bt, err = GetWorkCrHash(crs, BlockNumber)
		if err != nil {
			util.LogError("getWorkCrHash", err)
			return nil, err
		}
		err := store.AddToList(CrBucket, name, bt)
		if err != nil {
			util.LogError("AddCr", err)
			return nil, err
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

	// 获取工作证明
	// Get report of work
	report := []byte{}
	reportData, err := store.GetWorkDcapReport(wid)
	if err != nil {
		util.LogError("GetWorkDcapReport", err)
		report = []byte{}
	} else {
		hash := blake2b.Sum256(reportData)
		report = hash[:]
	}

	// TODO 暂时全部设置为true
	hasReport := true
	if len(report) > 0 {
		hasReport = true
	}

	// 所有需要提交的信息都不存在，不继续提交
	// All required submission information is missing, and the submission will not be continued.
	if report == nil && crHash == nil && logHash == nil {
		return nil, errors.New("report, crHash and logHash are all nil")
	}

	runtimeCall := weteeworker.MakeWorkProofUploadCall(
		wid,
		gtypes.OptionTProofOfWork{
			IsNone: !hasHash,
			IsSome: hasHash,
			AsSomeField0: gtypes.ProofOfWork{
				LogHash: logHash,
				CrHash:  crHash,
				Cr:      crProof,
			},
		},
		gtypes.OptionTByteSlice{
			IsNone:       !hasReport,
			IsSome:       hasReport,
			AsSomeField0: report,
		},
	)

	return &runtimeCall, nil
}

func SubmitWorkProof(client *chain.ChainClient, signer *core.Signer, proof []gtypes.RuntimeCall) error {
	call := utility.MakeBatchCall(proof)
	return client.SignAndSubmit(signer, call, true)
}

func CacheWorkProof(wid gtypes.WorkId, logs []string, crs map[string][]int64, now time.Time, BlockNumber uint64) error {
	name := util.GetWorkTypeStr(wid) + "-" + fmt.Sprint(wid.Id)

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
		err = store.AddToList(LogBucket, name+"_cache", bt)
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
		err := store.AddToList(CrBucket, name+"_cache", bt)
		if err != nil {
			util.LogError("AddCr", err)
			return err
		}
	}

	fmt.Println("cache work proof ========> ", len(logs), len(crs))

	return nil
}

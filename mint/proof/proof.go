package proof

import (
	"errors"
	"fmt"

	chain "github.com/wetee-dao/go-sdk"

	"github.com/wetee-dao/go-sdk/core"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	"github.com/wetee-dao/go-sdk/gen/utility"
	"github.com/wetee-dao/go-sdk/gen/weteeworker"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

func MakeWorkProof(wid gtypes.WorkId, logs []string, crs map[string][]int64, BlockNumber uint64) (*gtypes.RuntimeCall, error) {
	name := util.GetWorkTypeStr(wid) + "-" + fmt.Sprint(wid.Id)

	// 获取log和硬件资源使用量
	var logHash = []byte{}
	var crHash = []byte{}
	var cr = []uint32{0, 0, 0}
	var err error

	if len(logs) > 0 {
		// 获取log hash
		// Get log hash
		var bt []byte
		logHash, bt, err = GetWorkLogHash(logs, BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkLogHash", err)
			return nil, err
		}
		err = store.AddToList(LogBucket, []byte(LogBucket+name), bt)
		if err != nil {
			util.LogWithRed("Addlog", err)
			return nil, err
		}
	}

	if len(crs) > 0 {
		// 获取计算资源hash
		// Get Computing resource hash
		var bt []byte
		crHash, cr, bt, err = GetWorkCrHash(crs, BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkCrHash", err)
			return nil, err
		}
		err := store.AddToList(CrBucket, []byte(CrBucket+name), bt)
		if err != nil {
			util.LogWithRed("AddCr", err)
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
	report, err := store.GetWorkDcapReport(wid)
	if err != nil {
		util.LogWithRed("GetWorkDcapReport", err)
		report = []byte{}
	}

	hasReport := false
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
	return client.SignAndSubmit(signer, call, false)
}

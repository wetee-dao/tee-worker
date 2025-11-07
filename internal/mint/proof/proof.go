package proof

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"
	"golang.org/x/crypto/blake2b"

	inkutil "github.com/wetee-dao/ink.go/util"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

func MakeWorkProof(pod model.Pod, pod_key inkutil.Option[types.AccountID], logs []string, crs map[string][]int64, now time.Time) (*model.TeeCall, error) {
	err := AddMonitor(pod, logs, crs)
	if err != nil {
		return nil, err
	}

	// 获取 TEE 证明
	// Get TEE report of work
	reportHash := [32]byte{}
	report, err := store.GetPendingTEEReport(pod.PodId)
	if err == nil {
		reportData, _ := json.Marshal(report)
		hash := blake2b.Sum256(reportData)
		reportHash = hash
	} else {
		util.LogError("Proof", "GetWorkDcapReport ERROR:", err)
	}

	// 所有需要提交的信息都不存在，不继续提交
	// All required submission information is missing, and the submission will not be continued.
	if reportHash == [32]byte{} {
		return nil, errors.New("report is nil")
	}

	account, err := sidechain.GetDkgPubkey()
	if err != nil {
		return nil, err
	}

	err = chains.MainChain.DryMintPod(pod.PodId, types.H256(reportHash), *account)
	if err != nil {
		util.LogError("DryStartPod", err)
		return nil, err
	}

	podMint := model.TeeCall{
		Time: now.Unix(),
		Tx: &model.TeeCall_PodMint{
			PodMint: &model.PodMint{
				Id:         pod.PodId,
				ReportHash: reportHash[:],
			},
		},
	}

	return &podMint, nil
}

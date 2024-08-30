package mint

import (
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/go-resty/resty/v2"
	chain "github.com/wetee-dao/go-sdk"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteebridge"

	"wetee.app/worker/mint/proof"
	wtypes "wetee.app/worker/type"
	"wetee.app/worker/util"
)

func (m *Minter) trigger(cs map[gtypes.WorkId]ContractStateWrap, clusterId uint64, blockNumber uint64) {
	calls, keys, err := m.listTeeCalls(clusterId)
	if err != nil {
		fmt.Println("Tee trigger listTeeCalls error", err)
		return
	}

	callKey := make(map[gtypes.WorkId][]types.StorageKey)
	callId := make(map[gtypes.WorkId][]types.U128)
	// 为所有的 TEECall 通过 worker id 分组
	for i, call := range calls {
		callKey[call.WorkId] = append(callKey[call.WorkId], keys[i])
		callId[call.WorkId] = append(callId[call.WorkId], call.Id)
	}

	accounts, err := m.GetUsersFromCall(callKey)
	if err != nil {
		fmt.Println("Tee trigger GetUsersFromCall error", err)
		return
	}

	fmt.Println("Find ", len(callId), " trigger")
	for workId, ids := range callId {
		work := cs[workId]

		if blockNumber-work.Version < 20 || work.GetStatus() != 3 {
			// 跳过部署不超过 20 个区块的 TEECall
			fmt.Println("Skip ", util.GetWorkIdFromWorkType(workId), " trigger")
			continue
		}

		account := accounts[workId]
		saddress := AccountToSpace(account[:])
		name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
		client := resty.New()

		msg := fmt.Sprint(clusterId, ids)

		signer, err := m.PrivateKey.ToSigner()
		if err != nil {
			fmt.Println("Tee trigger GetSigner error", err)
			return
		}

		// 获取 worker report
		report, t, err := proof.GetRemoteReport(signer, []byte(msg))
		if err != nil {
			fmt.Println("Tee trigger GetRemoteReport error", err)
			return
		}

		// 获取 worker report
		paramWrap := wtypes.TeeParam{
			Address: signer.Address,
			Time:    t,
			Data:    nil,
			Report:  report,
		}

		// 解决 u128 编码和解码在tee中的错误问题
		var idstr []string
		for _, id := range ids {
			idstr = append(idstr, id.String())
		}

		ps := wtypes.TeeTrigger{
			Tee:       paramWrap,
			ClusterId: clusterId,
			Callids:   idstr,
		}

		bt, _ := json.Marshal(ps)

		// TODO
		// suite := m.PrivateKey.Suite()
		// ciphertext, err := ecies.Encrypt(suite, public, bt, suite.Hash)
		go func() {
			_, err = client.R().SetBody(bt).Post("http://" + name + "." + saddress + ".svc.cluster.local:65535/tee-call")
			if err != nil {
				fmt.Println("http://" + name + "." + saddress + ".svc.cluster.local:65535/tee-call")
				fmt.Println("Tee trigger http error", err)
				return
			}
		}()
	}
}

// 获取 Work 的帐户
func (m *Minter) GetUsersFromCall(calls map[gtypes.WorkId][]types.StorageKey) (map[gtypes.WorkId][32]byte, error) {
	keys := make([]types.StorageKey, 0, len(calls))
	keyMap := make(map[gtypes.WorkId]types.StorageKey)
	accounts := make(map[gtypes.WorkId][32]byte)
	for w := range calls {
		pallet := ""
		if w.Wtype.IsAPP {
			pallet = "WeTEEApp"
		} else if w.Wtype.IsGPU {
			pallet = "WeTEEApu"
		}

		// create key prefix
		key := chain.CreatePrefixedKey(pallet, "AppIdAccounts")
		hashers, err := m.ChainClient.GetHashers(pallet, "AppIdAccounts")
		if err != nil {
			return nil, err
		}

		// write key
		arg, err := codec.Encode(w.Id)
		if err != nil {
			return nil, err
		}
		_, err = hashers[0].Write(arg)
		if err != nil {
			return nil, fmt.Errorf("unable to hash args[%d]: %s Error: %v", 0, arg, err)
		}

		// append hash to key
		key = append(key, hashers[0].Sum(nil)...)
		keys = append(keys, key)
		keyMap[w] = key
	}

	set, err := m.ChainClient.Api.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return nil, err
	}

	for _, elem := range set {
		for _, change := range elem.Changes {
			var d [32]byte

			if err := codec.Decode(change.StorageData, &d); err != nil {
				continue
			}

			keys = append(keys, change.StorageKey)
			for w, key := range keyMap {
				if key.Hex() == change.StorageKey.Hex() {
					accounts[w] = d
				}
			}
		}
	}

	return accounts, nil
}

// list tee calls
func (m *Minter) listTeeCalls(cid uint64) ([]*gtypes.TEECall, []types.StorageKey, error) {
	var pallet, method = "WeTEEBridge", "TEECalls"
	set, err := m.ChainClient.QueryDoubleMapAll(pallet, method, cid, nil)
	if err != nil {
		return nil, nil, err
	}

	var list []*gtypes.TEECall = make([]*gtypes.TEECall, 0, len(set))
	var keys []types.StorageKey = make([]types.StorageKey, 0, len(set))
	for _, elem := range set {
		for _, change := range elem.Changes {
			var d gtypes.TEECall

			if err := codec.Decode(change.StorageData, &d); err != nil {
				fmt.Println(err)
				continue
			}
			keys = append(keys, change.StorageKey)
			list = append(list, &d)
		}
	}

	return list, keys, nil
}

// get tee call
func (m *Minter) getTeeCall(cid uint64, callid types.U128) (*gtypes.TEECall, error) {
	call, ok, err := weteebridge.GetTEECallsLatest(m.ChainClient.Api.RPC.State, cid, callid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &call, nil
}

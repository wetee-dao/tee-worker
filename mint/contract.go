package mint

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"wetee.app/worker/util"

	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	"github.com/wetee-dao/go-sdk/gen/weteeapp"
	"github.com/wetee-dao/go-sdk/gen/weteetask"
	"github.com/wetee-dao/go-sdk/gen/weteeworker"
)

type ContractStateWrap struct {
	BlockHash     string
	ContractState *gtypes.ClusterContractState
	WorkState     *gtypes.ContractState
	App           *gtypes.TeeApp
	Task          *gtypes.TeeTask
	Version       uint64
}

func (m *Minter) GetClusterContracts(clusterID uint64, at *types.Hash) (map[gtypes.WorkId]ContractStateWrap, error) {
	var pallet, method = "WeteeWorker", "ClusterContracts"
	set, err := m.ChainClient.QueryDoubleMapAll(pallet, method, clusterID, at)
	if err != nil {
		return nil, err
	}

	var list = map[gtypes.WorkId]ContractStateWrap{}
	var workIds = make([]gtypes.WorkId, 0, len(set))
	var workContractkeys = make([]types.StorageKey, 0, len(set))

	var appKeys = make([]types.StorageKey, 0, len(set))
	var appIds = make([]gtypes.WorkId, 0, len(set))
	var appVersions = make([]types.StorageKey, 0, len(set))

	var tasKeys = make([]types.StorageKey, 0, len(set))
	var taskIds = make([]gtypes.WorkId, 0, len(set))
	var taskVersions = make([]types.StorageKey, 0, len(set))

	for _, elem := range set {
		for _, change := range elem.Changes {
			var cs gtypes.ClusterContractState
			if err := codec.Decode(change.StorageData, &cs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			list[cs.WorkId] = ContractStateWrap{
				BlockHash:     elem.Block.Hex(),
				ContractState: &cs,
			}

			workIds = append(workIds, cs.WorkId)

			// 获取 work contract的key
			key, err := weteeworker.MakeWorkContractStateStorageKey(cs.WorkId, clusterID)
			if err != nil {
				continue
			}
			workContractkeys = append(workContractkeys, key)

			// 记录 app 相关参数
			if cs.WorkId.Wtype.IsAPP {
				akey, err := weteeapp.MakeTEEAppsStorageKey(cs.User, cs.WorkId.Id)
				if err != nil {
					continue
				}
				appKeys = append(appKeys, akey)
				appIds = append(appIds, cs.WorkId)
				vkey, err := weteeapp.MakeAppVersionStorageKey(cs.WorkId.Id)
				if err != nil {
					continue
				}
				appVersions = append(appVersions, vkey)
			}

			// 记录 task 相关参数
			if cs.WorkId.Wtype.IsTASK {
				tkey, err := weteetask.MakeTEETasksStorageKey(cs.User, cs.WorkId.Id)
				if err != nil {
					continue
				}
				tasKeys = append(tasKeys, tkey)
				taskIds = append(taskIds, cs.WorkId)
				vkey, err := weteetask.MakeTaskVersionStorageKey(cs.WorkId.Id)
				if err != nil {
					continue
				}
				taskVersions = append(taskVersions, vkey)
			}
		}
	}

	// 获取 work contract 的状态
	err = m.GetWorkContracts(workIds, workContractkeys, list, at)
	if err != nil {
		util.LogWithRed("GetWorkContracts", err)
		return nil, err
	}

	// 获取 app 的状态
	err = m.GetApps(appIds, appKeys, list, at)
	if err != nil {
		util.LogWithRed("GetApps", err)
		return nil, err
	}

	// 获取 task 的状态
	err = m.GetTasks(taskIds, tasKeys, list, at)
	if err != nil {
		util.LogWithRed("GetTasks", err)
		return nil, err
	}

	// 获取 app 的版本
	err = m.GetVerions(appIds, appVersions, list, at)
	if err != nil {
		util.LogWithRed("GetVerions APP", err)
		return nil, err
	}

	// 获取 task 的版本
	err = m.GetVerions(taskIds, taskVersions, list, at)
	if err != nil {
		util.LogWithRed("GetVerions TASK", err)
		return nil, err
	}

	return list, nil
}

func (m *Minter) GetWorkContracts(workID []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}
	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workID[IndexOf(wkeys, key)]
			var wcs gtypes.ContractState
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.WorkState = &wcs

			data[workId] = d
		}
	}

	return nil
}

func (m *Minter) GetApps(workID []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}
	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workID[IndexOf(wkeys, key)]
			var wcs gtypes.TeeApp
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.App = &wcs

			data[workId] = d
		}
	}

	return nil
}

func (m *Minter) GetTasks(workID []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workID[IndexOf(wkeys, key)]
			var wcs gtypes.TeeTask
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.Task = &wcs

			data[workId] = d
		}
	}

	return nil
}

func (m *Minter) GetVerions(workID []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}
	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workID[IndexOf(wkeys, key)]
			var wcs uint64
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.Version = wcs
			data[workId] = d
		}
	}

	return nil
}

func IndexOf(list []types.StorageKey, val types.StorageKey) int {
	for i, x := range list {
		if x.Hex() == val.Hex() {
			return i
		}
	}
	return -1
}

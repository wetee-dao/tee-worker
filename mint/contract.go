package mint

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"wetee.app/worker/util"

	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	"github.com/wetee-dao/go-sdk/gen/weteeapp"
	"github.com/wetee-dao/go-sdk/gen/weteegpu"
	"github.com/wetee-dao/go-sdk/gen/weteetask"
	"github.com/wetee-dao/go-sdk/gen/weteeworker"
)

// 合约状态
// Contract State
type ContractStateWrap struct {
	BlockHash     string
	ContractState *gtypes.ClusterContractState
	WorkState     *gtypes.ContractState
	App           *gtypes.TeeApp
	Task          *gtypes.TeeTask
	GpuApp        *gtypes.GpuApp
	Version       uint64
	Settings      []*gtypes.Env
}

// 获取合约状态
// Get Cluster Contracts
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
	var appSettingIds = make([]gtypes.WorkId, 0, len(set))
	var appSettings = make([]types.StorageKey, 0, len(set))

	var tasKeys = make([]types.StorageKey, 0, len(set))
	var taskIds = make([]gtypes.WorkId, 0, len(set))
	var taskVersions = make([]types.StorageKey, 0, len(set))
	var taskSettingIds = make([]gtypes.WorkId, 0, len(set))
	var taskSettings = make([]types.StorageKey, 0, len(set))

	var gpuAppKeys = make([]types.StorageKey, 0, len(set))
	var gpuAppIds = make([]gtypes.WorkId, 0, len(set))
	var gpuAppVersions = make([]types.StorageKey, 0, len(set))
	var gpuAppSettingIds = make([]gtypes.WorkId, 0, len(set))
	var gpuAppSettings = make([]types.StorageKey, 0, len(set))

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

			// 记录 work id
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
				skey, err := m.ChainClient.GetDoubleMapPrefixKey("WeteeApp", "AppSettings", cs.WorkId.Id)
				if err != nil {
					continue
				}

				var keys []types.StorageKey
				keys, err = m.ChainClient.Api.RPC.State.GetKeysLatest(skey)
				if err != nil {
					continue
				}

				for _, k := range keys {
					appSettingIds = append(appSettingIds, cs.WorkId)
					appSettings = append(appSettings, k)
				}
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
				skey, err := m.ChainClient.GetDoubleMapPrefixKey("WeteeTask", "AppSettings", cs.WorkId.Id)
				if err != nil {
					continue
				}

				var keys []types.StorageKey
				keys, err = m.ChainClient.Api.RPC.State.GetKeysLatest(skey)
				if err != nil {
					continue
				}

				for _, k := range keys {
					taskSettingIds = append(taskSettingIds, cs.WorkId)
					taskSettings = append(taskSettings, k)
				}
			}

			// 记录 gpu 相关参数
			if cs.WorkId.Wtype.IsGPU {
				tkey, err := weteegpu.MakeGPUAppsStorageKey(cs.User, cs.WorkId.Id)
				if err != nil {
					continue
				}

				gpuAppKeys = append(gpuAppKeys, tkey)
				gpuAppIds = append(gpuAppIds, cs.WorkId)
				vkey, err := weteegpu.MakeAppVersionStorageKey(cs.WorkId.Id)
				if err != nil {
					continue
				}

				gpuAppVersions = append(gpuAppVersions, vkey)
				skey, err := m.ChainClient.GetDoubleMapPrefixKey("WeteeGpu", "AppSettings", cs.WorkId.Id)
				if err != nil {
					continue
				}

				var keys []types.StorageKey
				keys, err = m.ChainClient.Api.RPC.State.GetKeysLatest(skey)
				if err != nil {
					continue
				}

				for _, k := range keys {
					gpuAppSettingIds = append(gpuAppSettingIds, cs.WorkId)
					gpuAppSettings = append(gpuAppSettings, k)
				}
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

	err = m.GetGpuApps(gpuAppIds, gpuAppKeys, list, at)
	if err != nil {
		util.LogWithRed("GetGpuApps", err)
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

	// 获取 task 的版本
	err = m.GetVerions(gpuAppIds, gpuAppVersions, list, at)
	if err != nil {
		util.LogWithRed("GetVerions GPU", err)
		return nil, err
	}

	// 获取 app 的设置
	err = m.GetSettings(appSettingIds, appSettings, list, at)
	if err != nil {
		util.LogWithRed("GetSettings APP", err)
		return nil, err
	}

	// 获取 task 的设置
	err = m.GetSettings(taskSettingIds, taskSettings, list, at)
	if err != nil {
		util.LogWithRed("GetSettings TASK", err)
		return nil, err
	}

	// 获取 task 的设置
	err = m.GetSettings(gpuAppSettingIds, gpuAppSettings, list, at)
	if err != nil {
		util.LogWithRed("GetSettings GPU", err)
		return nil, err
	}

	return list, nil
}

// 获取 work contract 的状态
// Get Work Contracts
func (m *Minter) GetWorkContracts(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
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

// 获取应用信息
// Get app info
func (m *Minter) GetApps(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
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

// 获取 task 的状态
// Get Task info
func (m *Minter) GetTasks(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
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

// 获取 task 的状态
// Get Task info
func (m *Minter) GetGpuApps(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
			var wcs gtypes.GpuApp
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.GpuApp = &wcs

			data[workId] = d
		}
	}

	return nil
}

// 获取版本信息
// Get Version info
func (m *Minter) GetVerions(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
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

// 获取版本信息
// Get settings
func (m *Minter) GetSettings(workId []gtypes.WorkId, wkeys []types.StorageKey, data map[gtypes.WorkId]ContractStateWrap, at *types.Hash) error {
	wsets, err := m.ChainClient.Api.RPC.State.QueryStorageLatest(wkeys, *at)
	if err != nil {
		return err
	}

	for _, elem := range wsets {
		for _, change := range elem.Changes {
			var key = change.StorageKey
			var workId = workId[IndexOf(wkeys, key)]
			var wcs gtypes.Env
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			d := data[workId]
			d.Settings = append(d.Settings, &wcs)
			data[workId] = d
		}
	}

	return nil
}

func (m *Minter) GetSettingsFromWork(workId gtypes.WorkId, at *types.Hash) ([]*gtypes.Env, error) {
	var pallet, method string
	if workId.Wtype.IsAPP {
		pallet = "WeteeApp"
		method = "AppSettings"
	} else if workId.Wtype.IsTASK {
		pallet = "WeteeTask"
		method = "AppSettings"
	}

	sets, err := m.ChainClient.QueryDoubleMapAll(pallet, method, workId.Id, at)
	if err != nil {
		util.LogWithRed("QueryDoubleMapAll", err)
		return []*gtypes.Env{}, err
	}

	settings := make([]*gtypes.Env, 0, len(sets))
	for _, elem := range sets {
		for _, change := range elem.Changes {
			var wcs gtypes.Env
			if err := codec.Decode(change.StorageData, &wcs); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			settings = append(settings, &wcs)
		}
	}

	return settings, nil
}

// 获取索引
// Get index of key
func IndexOf(list []types.StorageKey, val types.StorageKey) int {
	for i, x := range list {
		if x.Hex() == val.Hex() {
			return i
		}
	}
	return -1
}

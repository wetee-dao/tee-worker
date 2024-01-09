package weteeapp

import (
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "wetee.app/worker/internal/mint/chain/gen/types"
)

// See [`Pallet::create`].
func MakeCreateCall(name0 []byte, image1 []byte, port2 []uint32, cpu3 uint16, memory4 uint16, disk5 uint16, level6 byte, deposit7 types.UCompact) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeApp: true,
		AsWeteeAppField0: &types1.WeteeAppPalletCall{
			IsCreate:         true,
			AsCreateName0:    name0,
			AsCreateImage1:   image1,
			AsCreatePort2:    port2,
			AsCreateCpu3:     cpu3,
			AsCreateMemory4:  memory4,
			AsCreateDisk5:    disk5,
			AsCreateLevel6:   level6,
			AsCreateDeposit7: deposit7,
		},
	}
}

// See [`Pallet::update`].
func MakeUpdateCall(appId0 uint64, name1 []byte, image2 []byte, port3 []uint32) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeApp: true,
		AsWeteeAppField0: &types1.WeteeAppPalletCall{
			IsUpdate:       true,
			AsUpdateAppId0: appId0,
			AsUpdateName1:  name1,
			AsUpdateImage2: image2,
			AsUpdatePort3:  port3,
		},
	}
}

// See [`Pallet::set_settings`].
func MakeSetSettingsCall(appId0 uint64, value1 []types1.AppSettingInput) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeApp: true,
		AsWeteeAppField0: &types1.WeteeAppPalletCall{
			IsSetSettings:       true,
			AsSetSettingsAppId0: appId0,
			AsSetSettingsValue1: value1,
		},
	}
}

// See [`Pallet::recharge`].
func MakeRechargeCall(id0 uint64, deposit1 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeApp: true,
		AsWeteeAppField0: &types1.WeteeAppPalletCall{
			IsRecharge:         true,
			AsRechargeId0:      id0,
			AsRechargeDeposit1: deposit1,
		},
	}
}

// See [`Pallet::stop`].
func MakeStopCall(appId0 uint64) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeApp: true,
		AsWeteeAppField0: &types1.WeteeAppPalletCall{
			IsStop:       true,
			AsStopAppId0: appId0,
		},
	}
}

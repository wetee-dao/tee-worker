package weteesudo

import types "wetee.app/worker/mint/chain/gen/types"

// See [`Pallet::sudo`].
func MakeSudoCall(daoId0 uint64, call1 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeSudo: true,
		AsWeteeSudoField0: &types.WeteeSudoPalletCall{
			IsSudo:       true,
			AsSudoDaoId0: daoId0,
			AsSudoCall1:  &call1,
		},
	}
}

// See [`Pallet::set_sudo_account`].
func MakeSetSudoAccountCall(daoId0 uint64, sudoAccount1 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeSudo: true,
		AsWeteeSudoField0: &types.WeteeSudoPalletCall{
			IsSetSudoAccount:             true,
			AsSetSudoAccountDaoId0:       daoId0,
			AsSetSudoAccountSudoAccount1: sudoAccount1,
		},
	}
}

// See [`Pallet::close_sudo`].
func MakeCloseSudoCall(daoId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeSudo: true,
		AsWeteeSudoField0: &types.WeteeSudoPalletCall{
			IsCloseSudo:       true,
			AsCloseSudoDaoId0: daoId0,
		},
	}
}

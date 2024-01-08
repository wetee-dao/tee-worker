package sudo

import types "wetee.app/worker/internal/worker/chain/gen/types"

// See [`Pallet::sudo`].
func MakeSudoCall(call0 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudo:      true,
			AsSudoCall0: &call0,
		},
	}
}

// See [`Pallet::sudo_unchecked_weight`].
func MakeSudoUncheckedWeightCall(call0 types.RuntimeCall, weight1 types.Weight) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudoUncheckedWeight:        true,
			AsSudoUncheckedWeightCall0:   &call0,
			AsSudoUncheckedWeightWeight1: weight1,
		},
	}
}

// See [`Pallet::set_key`].
func MakeSetKeyCall(new0 types.MultiAddress) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSetKey:     true,
			AsSetKeyNew0: &new0,
		},
	}
}

// See [`Pallet::sudo_as`].
func MakeSudoAsCall(who0 types.MultiAddress, call1 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsSudo: true,
		AsSudoField0: &types.PalletSudoPalletCall{
			IsSudoAs:      true,
			AsSudoAsWho0:  who0,
			AsSudoAsCall1: &call1,
		},
	}
}

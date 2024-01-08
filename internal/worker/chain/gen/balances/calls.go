package balances

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "wetee.app/worker/internal/worker/chain/gen/types"
)

// See [`Pallet::transfer_allow_death`].
func MakeTransferAllowDeathCall(dest0 types.MultiAddress, value1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferAllowDeath:       true,
			AsTransferAllowDeathDest0:  dest0,
			AsTransferAllowDeathValue1: value1,
		},
	}
}

// See [`Pallet::set_balance_deprecated`].
func MakeSetBalanceDeprecatedCall(who0 types.MultiAddress, newFree1 types1.UCompact, oldReserved2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsSetBalanceDeprecated:             true,
			AsSetBalanceDeprecatedWho0:         who0,
			AsSetBalanceDeprecatedNewFree1:     newFree1,
			AsSetBalanceDeprecatedOldReserved2: oldReserved2,
		},
	}
}

// See [`Pallet::force_transfer`].
func MakeForceTransferCall(source0 types.MultiAddress, dest1 types.MultiAddress, value2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceTransfer:        true,
			AsForceTransferSource0: source0,
			AsForceTransferDest1:   dest1,
			AsForceTransferValue2:  value2,
		},
	}
}

// See [`Pallet::transfer_keep_alive`].
func MakeTransferKeepAliveCall(dest0 types.MultiAddress, value1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferKeepAlive:       true,
			AsTransferKeepAliveDest0:  dest0,
			AsTransferKeepAliveValue1: value1,
		},
	}
}

// See [`Pallet::transfer_all`].
func MakeTransferAllCall(dest0 types.MultiAddress, keepAlive1 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransferAll:           true,
			AsTransferAllDest0:      dest0,
			AsTransferAllKeepAlive1: keepAlive1,
		},
	}
}

// See [`Pallet::force_unreserve`].
func MakeForceUnreserveCall(who0 types.MultiAddress, amount1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceUnreserve:        true,
			AsForceUnreserveWho0:    who0,
			AsForceUnreserveAmount1: amount1,
		},
	}
}

// See [`Pallet::upgrade_accounts`].
func MakeUpgradeAccountsCall(who0 [][32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsUpgradeAccounts:     true,
			AsUpgradeAccountsWho0: who0,
		},
	}
}

// See [`Pallet::transfer`].
func MakeTransferCall(dest0 types.MultiAddress, value1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsTransfer:       true,
			AsTransferDest0:  dest0,
			AsTransferValue1: value1,
		},
	}
}

// See [`Pallet::force_set_balance`].
func MakeForceSetBalanceCall(who0 types.MultiAddress, newFree1 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsBalances: true,
		AsBalancesField0: &types.PalletBalancesPalletCall{
			IsForceSetBalance:         true,
			AsForceSetBalanceWho0:     who0,
			AsForceSetBalanceNewFree1: newFree1,
		},
	}
}

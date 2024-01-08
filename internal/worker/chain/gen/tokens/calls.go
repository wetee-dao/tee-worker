package tokens

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "wetee.app/worker/internal/worker/chain/gen/types"
)

// See [`Pallet::transfer`].
func MakeTransferCall(dest0 types.MultiAddress, currencyId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransfer:            true,
			AsTransferDest0:       dest0,
			AsTransferCurrencyId1: currencyId1,
			AsTransferAmount2:     amount2,
		},
	}
}

// See [`Pallet::transfer_all`].
func MakeTransferAllCall(dest0 types.MultiAddress, currencyId1 uint64, keepAlive2 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransferAll:            true,
			AsTransferAllDest0:       dest0,
			AsTransferAllCurrencyId1: currencyId1,
			AsTransferAllKeepAlive2:  keepAlive2,
		},
	}
}

// See [`Pallet::transfer_keep_alive`].
func MakeTransferKeepAliveCall(dest0 types.MultiAddress, currencyId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsTransferKeepAlive:            true,
			AsTransferKeepAliveDest0:       dest0,
			AsTransferKeepAliveCurrencyId1: currencyId1,
			AsTransferKeepAliveAmount2:     amount2,
		},
	}
}

// See [`Pallet::force_transfer`].
func MakeForceTransferCall(source0 types.MultiAddress, dest1 types.MultiAddress, currencyId2 uint64, amount3 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsForceTransfer:            true,
			AsForceTransferSource0:     source0,
			AsForceTransferDest1:       dest1,
			AsForceTransferCurrencyId2: currencyId2,
			AsForceTransferAmount3:     amount3,
		},
	}
}

// See [`Pallet::set_balance`].
func MakeSetBalanceCall(who0 types.MultiAddress, currencyId1 uint64, newFree2 types1.UCompact, newReserved3 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsTokens: true,
		AsTokensField0: &types.OrmlTokensModuleCall{
			IsSetBalance:             true,
			AsSetBalanceWho0:         who0,
			AsSetBalanceCurrencyId1:  currencyId1,
			AsSetBalanceNewFree2:     newFree2,
			AsSetBalanceNewReserved3: newReserved3,
		},
	}
}

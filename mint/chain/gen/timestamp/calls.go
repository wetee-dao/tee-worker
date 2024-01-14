package timestamp

import (
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// See [`Pallet::set`].
func MakeSetCall(now0 types.UCompact) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsTimestamp: true,
		AsTimestampField0: &types1.PalletTimestampPalletCall{
			IsSet:     true,
			AsSetNow0: now0,
		},
	}
}

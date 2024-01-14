package weteetreasury

import (
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// See [`Pallet::spend`].
func MakeSpendCall(daoId0 uint64, beneficiary1 [32]byte, amount2 types.UCompact) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsWeteeTreasury: true,
		AsWeteeTreasuryField0: &types1.WeteeTreasuryPalletCall{
			IsSpend:             true,
			AsSpendDaoId0:       daoId0,
			AsSpendBeneficiary1: beneficiary1,
			AsSpendAmount2:      amount2,
		},
	}
}

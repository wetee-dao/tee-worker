package weteeasset

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "wetee.app/worker/mint/chain/gen/types"
)

// See [`Pallet::create_asset`].
func MakeCreateAssetCall(daoId0 uint64, metadata1 types.DaoAssetMeta, amount2 types1.U128, initDaoAsset3 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsCreateAsset:              true,
			AsCreateAssetDaoId0:        daoId0,
			AsCreateAssetMetadata1:     metadata1,
			AsCreateAssetAmount2:       amount2,
			AsCreateAssetInitDaoAsset3: initDaoAsset3,
		},
	}
}

// See [`Pallet::set_existenial_deposit`].
func MakeSetExistenialDepositCall(daoId0 uint64, existenialDeposit1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsSetExistenialDeposit:                   true,
			AsSetExistenialDepositDaoId0:             daoId0,
			AsSetExistenialDepositExistenialDeposit1: existenialDeposit1,
		},
	}
}

// See [`Pallet::set_metadata`].
func MakeSetMetadataCall(daoId0 uint64, metadata1 types.DaoAssetMeta) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsSetMetadata:          true,
			AsSetMetadataDaoId0:    daoId0,
			AsSetMetadataMetadata1: metadata1,
		},
	}
}

// See [`Pallet::burn`].
func MakeBurnCall(daoId0 uint64, amount1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsBurn:        true,
			AsBurnDaoId0:  daoId0,
			AsBurnAmount1: amount1,
		},
	}
}

// See [`Pallet::transfer`].
func MakeTransferCall(dest0 types.MultiAddress, daoId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsTransfer:        true,
			AsTransferDest0:   dest0,
			AsTransferDaoId1:  daoId1,
			AsTransferAmount2: amount2,
		},
	}
}

// See [`Pallet::join`].
func MakeJoinCall(daoId0 uint64, shareExpect1 uint32, existenialDeposit2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeAsset: true,
		AsWeteeAssetField0: &types.WeteeAssetsPalletCall{
			IsJoin:                   true,
			AsJoinDaoId0:             daoId0,
			AsJoinShareExpect1:       shareExpect1,
			AsJoinExistenialDeposit2: existenialDeposit2,
		},
	}
}

package grandpa

import types "wetee.app/worker/internal/mint/chain/gen/types"

// See [`Pallet::report_equivocation`].
func MakeReportEquivocationCall(equivocationProof0 types.EquivocationProof, keyOwnerProof1 struct{}) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsReportEquivocation:                   true,
			AsReportEquivocationEquivocationProof0: &equivocationProof0,
			AsReportEquivocationKeyOwnerProof1:     keyOwnerProof1,
		},
	}
}

// See [`Pallet::report_equivocation_unsigned`].
func MakeReportEquivocationUnsignedCall(equivocationProof0 types.EquivocationProof, keyOwnerProof1 struct{}) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsReportEquivocationUnsigned:                   true,
			AsReportEquivocationUnsignedEquivocationProof0: &equivocationProof0,
			AsReportEquivocationUnsignedKeyOwnerProof1:     keyOwnerProof1,
		},
	}
}

// See [`Pallet::note_stalled`].
func MakeNoteStalledCall(delay0 uint64, bestFinalizedBlockNumber1 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsGrandpa: true,
		AsGrandpaField0: &types.PalletGrandpaPalletCall{
			IsNoteStalled:                          true,
			AsNoteStalledDelay0:                    delay0,
			AsNoteStalledBestFinalizedBlockNumber1: bestFinalizedBlockNumber1,
		},
	}
}

package weteegov

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "wetee.app/worker/mint/chain/gen/types"
)

// See [`Pallet::submit_proposal`].
func MakeSubmitProposalCall(daoId0 uint64, memberData1 types.MemberData, proposal2 types.RuntimeCall, periodIndex3 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsSubmitProposal:             true,
			AsSubmitProposalDaoId0:       daoId0,
			AsSubmitProposalMemberData1:  memberData1,
			AsSubmitProposalProposal2:    &proposal2,
			AsSubmitProposalPeriodIndex3: periodIndex3,
		},
	}
}

// See [`Pallet::deposit_proposal`].
func MakeDepositProposalCall(daoId0 uint64, proposeId1 uint32, deposit2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsDepositProposal:           true,
			AsDepositProposalDaoId0:     daoId0,
			AsDepositProposalProposeId1: proposeId1,
			AsDepositProposalDeposit2:   deposit2,
		},
	}
}

// See [`Pallet::vote_for_prop`].
func MakeVoteForPropCall(daoId0 uint64, propIndex1 uint32, pledge2 types.Pledge, opinion3 types.Opinion) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsVoteForProp:           true,
			AsVoteForPropDaoId0:     daoId0,
			AsVoteForPropPropIndex1: propIndex1,
			AsVoteForPropPledge2:    pledge2,
			AsVoteForPropOpinion3:   opinion3,
		},
	}
}

// See [`Pallet::cancel_vote`].
func MakeCancelVoteCall(daoId0 uint64, index1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsCancelVote:       true,
			AsCancelVoteDaoId0: daoId0,
			AsCancelVoteIndex1: index1,
		},
	}
}

// See [`Pallet::run_proposal`].
func MakeRunProposalCall(daoId0 uint64, index1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsRunProposal:       true,
			AsRunProposalDaoId0: daoId0,
			AsRunProposalIndex1: index1,
		},
	}
}

// See [`Pallet::unlock`].
func MakeUnlockCall(daoId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsUnlock:       true,
			AsUnlockDaoId0: daoId0,
		},
	}
}

// See [`Pallet::set_max_pre_props`].
func MakeSetMaxPrePropsCall(daoId0 uint64, max1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsSetMaxPreProps:       true,
			AsSetMaxPrePropsDaoId0: daoId0,
			AsSetMaxPrePropsMax1:   max1,
		},
	}
}

// See [`Pallet::update_vote_model`].
func MakeUpdateVoteModelCall(daoId0 uint64, model1 byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsUpdateVoteModel:       true,
			AsUpdateVoteModelDaoId0: daoId0,
			AsUpdateVoteModelModel1: model1,
		},
	}
}

// See [`Pallet::set_periods`].
func MakeSetPeriodsCall(daoId0 uint64, periods1 []types.Period) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGov: true,
		AsWeteeGovField0: &types.WeteeGovPalletCall{
			IsSetPeriods:         true,
			AsSetPeriodsDaoId0:   daoId0,
			AsSetPeriodsPeriods1: periods1,
		},
	}
}

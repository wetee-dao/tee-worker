package messagequeue

import types "wetee.app/worker/mint/chain/gen/types"

// See [`Pallet::reap_page`].
func MakeReapPageCall(messageOrigin0 types.MessageOrigin, pageIndex1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsMessageQueue: true,
		AsMessageQueueField0: &types.PalletMessageQueuePalletCall{
			IsReapPage:               true,
			AsReapPageMessageOrigin0: messageOrigin0,
			AsReapPagePageIndex1:     pageIndex1,
		},
	}
}

// See [`Pallet::execute_overweight`].
func MakeExecuteOverweightCall(messageOrigin0 types.MessageOrigin, page1 uint32, index2 uint32, weightLimit3 types.Weight) types.RuntimeCall {
	return types.RuntimeCall{
		IsMessageQueue: true,
		AsMessageQueueField0: &types.PalletMessageQueuePalletCall{
			IsExecuteOverweight:               true,
			AsExecuteOverweightMessageOrigin0: messageOrigin0,
			AsExecuteOverweightPage1:          page1,
			AsExecuteOverweightIndex2:         index2,
			AsExecuteOverweightWeightLimit3:   weightLimit3,
		},
	}
}

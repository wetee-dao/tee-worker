package system

import types "wetee.app/worker/mint/chain/gen/types"

// See [`Pallet::remark`].
func MakeRemarkCall(remark0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsRemark:        true,
			AsRemarkRemark0: remark0,
		},
	}
}

// See [`Pallet::set_heap_pages`].
func MakeSetHeapPagesCall(pages0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetHeapPages:       true,
			AsSetHeapPagesPages0: pages0,
		},
	}
}

// See [`Pallet::set_code`].
func MakeSetCodeCall(code0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetCode:      true,
			AsSetCodeCode0: code0,
		},
	}
}

// See [`Pallet::set_code_without_checks`].
func MakeSetCodeWithoutChecksCall(code0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetCodeWithoutChecks:      true,
			AsSetCodeWithoutChecksCode0: code0,
		},
	}
}

// See [`Pallet::set_storage`].
func MakeSetStorageCall(items0 []types.TupleOfByteSliceByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsSetStorage:       true,
			AsSetStorageItems0: items0,
		},
	}
}

// See [`Pallet::kill_storage`].
func MakeKillStorageCall(keys0 [][]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsKillStorage:      true,
			AsKillStorageKeys0: keys0,
		},
	}
}

// See [`Pallet::kill_prefix`].
func MakeKillPrefixCall(prefix0 []byte, subkeys1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsKillPrefix:         true,
			AsKillPrefixPrefix0:  prefix0,
			AsKillPrefixSubkeys1: subkeys1,
		},
	}
}

// See [`Pallet::remark_with_event`].
func MakeRemarkWithEventCall(remark0 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsSystem: true,
		AsSystemField0: &types.FrameSystemPalletCall{
			IsRemarkWithEvent:        true,
			AsRemarkWithEventRemark0: remark0,
		},
	}
}

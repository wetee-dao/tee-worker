package system

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for Account
//
//	The full account information for a particular account ID.
func MakeAccountStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "System", "Account", byteArgs...)
}

var AccountResultDefaultBytes, _ = hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080")

func GetAccount(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret types1.AccountInfo, err error) {
	key, err := MakeAccountStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetAccountLatest(state state.State, byteArray320 [32]byte) (ret types1.AccountInfo, err error) {
	key, err := MakeAccountStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ExtrinsicCount id={{false [8]}}
//
//	Total extrinsics count for the current block.
func MakeExtrinsicCountStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "ExtrinsicCount")
}
func GetExtrinsicCount(state state.State, bhash types.Hash) (ret uint32, isSome bool, err error) {
	key, err := MakeExtrinsicCountStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetExtrinsicCountLatest(state state.State) (ret uint32, isSome bool, err error) {
	key, err := MakeExtrinsicCountStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for BlockWeight id={{false [9]}}
//
//	The current weight for the block.
func MakeBlockWeightStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "BlockWeight")
}

var BlockWeightResultDefaultBytes, _ = hex.DecodeString("000000000000")

func GetBlockWeight(state state.State, bhash types.Hash) (ret types1.PerDispatchClass, err error) {
	key, err := MakeBlockWeightStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BlockWeightResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetBlockWeightLatest(state state.State) (ret types1.PerDispatchClass, err error) {
	key, err := MakeBlockWeightStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BlockWeightResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for AllExtrinsicsLen id={{false [8]}}
//
//	Total length (in bytes) for all extrinsics put together, for the current block.
func MakeAllExtrinsicsLenStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "AllExtrinsicsLen")
}
func GetAllExtrinsicsLen(state state.State, bhash types.Hash) (ret uint32, isSome bool, err error) {
	key, err := MakeAllExtrinsicsLenStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAllExtrinsicsLenLatest(state state.State) (ret uint32, isSome bool, err error) {
	key, err := MakeAllExtrinsicsLenStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for BlockHash
//
//	Map of block numbers to block hashes.
func MakeBlockHashStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "System", "BlockHash", byteArgs...)
}

var BlockHashResultDefaultBytes, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

func GetBlockHash(state state.State, bhash types.Hash, uint640 uint64) (ret [32]byte, err error) {
	key, err := MakeBlockHashStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BlockHashResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetBlockHashLatest(state state.State, uint640 uint64) (ret [32]byte, err error) {
	key, err := MakeBlockHashStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BlockHashResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ExtrinsicData
//
//	Extrinsics data for the current block (maps an extrinsic's index to its data).
func MakeExtrinsicDataStorageKey(uint320 uint32) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "System", "ExtrinsicData", byteArgs...)
}

var ExtrinsicDataResultDefaultBytes, _ = hex.DecodeString("00")

func GetExtrinsicData(state state.State, bhash types.Hash, uint320 uint32) (ret []byte, err error) {
	key, err := MakeExtrinsicDataStorageKey(uint320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ExtrinsicDataResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetExtrinsicDataLatest(state state.State, uint320 uint32) (ret []byte, err error) {
	key, err := MakeExtrinsicDataStorageKey(uint320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ExtrinsicDataResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Number id={{false [4]}}
//
//	The current block number being processed. Set by `execute_block`.
func MakeNumberStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "Number")
}

var NumberResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetNumber(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNumberStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NumberResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNumberLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNumberStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NumberResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ParentHash id={{false [12]}}
//
//	Hash of the previous block.
func MakeParentHashStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "ParentHash")
}

var ParentHashResultDefaultBytes, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

func GetParentHash(state state.State, bhash types.Hash) (ret [32]byte, err error) {
	key, err := MakeParentHashStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ParentHashResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetParentHashLatest(state state.State) (ret [32]byte, err error) {
	key, err := MakeParentHashStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ParentHashResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Digest id={{false [14]}}
//
//	Digest of the current block, also part of the block header.
func MakeDigestStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "Digest")
}

var DigestResultDefaultBytes, _ = hex.DecodeString("00")

func GetDigest(state state.State, bhash types.Hash) (ret []types1.DigestItem, err error) {
	key, err := MakeDigestStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DigestResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetDigestLatest(state state.State) (ret []types1.DigestItem, err error) {
	key, err := MakeDigestStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DigestResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Events id={{false [18]}}
//
//	Events deposited for the current block.
//
//	NOTE: The item is unbound and should therefore never be read on chain.
//	It could otherwise inflate the PoV size of a block.
//
//	Events have a large in-memory size. Box the events to not go out-of-memory
//	just in case someone still reads them from within the runtime.
func MakeEventsStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "Events")
}

var EventsResultDefaultBytes, _ = hex.DecodeString("00")

func GetEvents(state state.State, bhash types.Hash) (ret []types1.EventRecord, err error) {
	key, err := MakeEventsStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetEventsLatest(state state.State) (ret []types1.EventRecord, err error) {
	key, err := MakeEventsStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for EventCount id={{false [8]}}
//
//	The number of events in the `Events<T>` list.
func MakeEventCountStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "EventCount")
}

var EventCountResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetEventCount(state state.State, bhash types.Hash) (ret uint32, err error) {
	key, err := MakeEventCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetEventCountLatest(state state.State) (ret uint32, err error) {
	key, err := MakeEventCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for EventTopics
//
//	Mapping between a topic (represented by T::Hash) and a vector of indexes
//	of events in the `<Events<T>>` list.
//
//	All topic vectors have deterministic storage locations depending on the topic. This
//	allows light-clients to leverage the changes trie storage tracking mechanism and
//	in case of changes fetch the list of events of interest.
//
//	The value has the type `(BlockNumberFor<T>, EventIndex)` because if we used only just
//	the `EventIndex` then in case if the topic has the same contents on the next block
//	no notification will be triggered thus the event might be lost.
func MakeEventTopicsStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "System", "EventTopics", byteArgs...)
}

var EventTopicsResultDefaultBytes, _ = hex.DecodeString("00")

func GetEventTopics(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.TupleOfUint64Uint32, err error) {
	key, err := MakeEventTopicsStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventTopicsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetEventTopicsLatest(state state.State, byteArray320 [32]byte) (ret []types1.TupleOfUint64Uint32, err error) {
	key, err := MakeEventTopicsStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(EventTopicsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for LastRuntimeUpgrade id={{false [67]}}
//
//	Stores the `spec_version` and `spec_name` of when the last runtime upgrade happened.
func MakeLastRuntimeUpgradeStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "LastRuntimeUpgrade")
}
func GetLastRuntimeUpgrade(state state.State, bhash types.Hash) (ret types1.LastRuntimeUpgradeInfo, isSome bool, err error) {
	key, err := MakeLastRuntimeUpgradeStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetLastRuntimeUpgradeLatest(state state.State) (ret types1.LastRuntimeUpgradeInfo, isSome bool, err error) {
	key, err := MakeLastRuntimeUpgradeStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for UpgradedToU32RefCount id={{false [47]}}
//
//	True if we have upgraded so that `type RefCount` is `u32`. False (default) if not.
func MakeUpgradedToU32RefCountStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "UpgradedToU32RefCount")
}

var UpgradedToU32RefCountResultDefaultBytes, _ = hex.DecodeString("00")

func GetUpgradedToU32RefCount(state state.State, bhash types.Hash) (ret bool, err error) {
	key, err := MakeUpgradedToU32RefCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UpgradedToU32RefCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetUpgradedToU32RefCountLatest(state state.State) (ret bool, err error) {
	key, err := MakeUpgradedToU32RefCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UpgradedToU32RefCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for UpgradedToTripleRefCount id={{false [47]}}
//
//	True if we have upgraded so that AccountInfo contains three types of `RefCount`. False
//	(default) if not.
func MakeUpgradedToTripleRefCountStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "UpgradedToTripleRefCount")
}

var UpgradedToTripleRefCountResultDefaultBytes, _ = hex.DecodeString("00")

func GetUpgradedToTripleRefCount(state state.State, bhash types.Hash) (ret bool, err error) {
	key, err := MakeUpgradedToTripleRefCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UpgradedToTripleRefCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetUpgradedToTripleRefCountLatest(state state.State) (ret bool, err error) {
	key, err := MakeUpgradedToTripleRefCountStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UpgradedToTripleRefCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ExecutionPhase id={{false [63]}}
//
//	The execution phase of the block.
func MakeExecutionPhaseStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "System", "ExecutionPhase")
}
func GetExecutionPhase(state state.State, bhash types.Hash) (ret types1.Phase, isSome bool, err error) {
	key, err := MakeExecutionPhaseStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetExecutionPhaseLatest(state state.State) (ret types1.Phase, isSome bool, err error) {
	key, err := MakeExecutionPhaseStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

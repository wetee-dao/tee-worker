package grandpa

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/internal/mint/chain/gen/types"
)

// Make a storage key for State id={{false [91]}}
//
//	State of the current authority set.
func MakeStateStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "State")
}

var StateResultDefaultBytes, _ = hex.DecodeString("00")

func GetState(state state.State, bhash types.Hash) (ret types1.StoredState, err error) {
	key, err := MakeStateStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(StateResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetStateLatest(state state.State) (ret types1.StoredState, err error) {
	key, err := MakeStateStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(StateResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for PendingChange id={{false [92]}}
//
//	Pending change: (signaled at, scheduled change).
func MakePendingChangeStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "PendingChange")
}
func GetPendingChange(state state.State, bhash types.Hash) (ret types1.StoredPendingChange, isSome bool, err error) {
	key, err := MakePendingChangeStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetPendingChangeLatest(state state.State) (ret types1.StoredPendingChange, isSome bool, err error) {
	key, err := MakePendingChangeStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for NextForced id={{false [4]}}
//
//	next block number where we can force a change.
func MakeNextForcedStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "NextForced")
}
func GetNextForced(state state.State, bhash types.Hash) (ret uint64, isSome bool, err error) {
	key, err := MakeNextForcedStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetNextForcedLatest(state state.State) (ret uint64, isSome bool, err error) {
	key, err := MakeNextForcedStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Stalled id={{false [95]}}
//
//	`true` if we are currently stalled.
func MakeStalledStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "Stalled")
}
func GetStalled(state state.State, bhash types.Hash) (ret types1.TupleOfUint64Uint64, isSome bool, err error) {
	key, err := MakeStalledStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetStalledLatest(state state.State) (ret types1.TupleOfUint64Uint64, isSome bool, err error) {
	key, err := MakeStalledStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for CurrentSetId id={{false [4]}}
//
//	The number of changes (both in terms of keys and underlying economic responsibilities)
//	in the "set" of Grandpa validators from genesis.
func MakeCurrentSetIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "CurrentSetId")
}

var CurrentSetIdResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetCurrentSetId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeCurrentSetIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CurrentSetIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetCurrentSetIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeCurrentSetIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CurrentSetIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for SetIdSession
//
//	A mapping from grandpa set ID to the index of the *most recent* session for which its
//	members were responsible.
//
//	This is only used for validating equivocation proofs. An equivocation proof must
//	contains a key-ownership proof for a given session, therefore we need a way to tie
//	together sessions and GRANDPA set ids, i.e. we need to validate that a validator
//	was the owner of a given key on a given session, and what the active set ID was
//	during that session.
//
//	TWOX-NOTE: `SetId` is not under user control.
func MakeSetIdSessionStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Grandpa", "SetIdSession", byteArgs...)
}
func GetSetIdSession(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeSetIdSessionStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetSetIdSessionLatest(state state.State, uint640 uint64) (ret uint32, isSome bool, err error) {
	key, err := MakeSetIdSessionStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

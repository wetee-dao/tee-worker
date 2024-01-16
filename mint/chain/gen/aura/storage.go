package aura

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for Authorities id={{false [88]}}
//
//	The current authority set.
func MakeAuthoritiesStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Aura", "Authorities")
}

var AuthoritiesResultDefaultBytes, _ = hex.DecodeString("00")

func GetAuthorities(state state.State, bhash types.Hash) (ret [][32]byte, err error) {
	key, err := MakeAuthoritiesStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AuthoritiesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetAuthoritiesLatest(state state.State) (ret [][32]byte, err error) {
	key, err := MakeAuthoritiesStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AuthoritiesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for CurrentSlot id={{false [92]}}
//
//	The current slot of this block.
//
//	This will be set in `on_initialize`.
func MakeCurrentSlotStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Aura", "CurrentSlot")
}

var CurrentSlotResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetCurrentSlot(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeCurrentSlotStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CurrentSlotResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetCurrentSlotLatest(state state.State) (ret uint64, err error) {
	key, err := MakeCurrentSlotStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(CurrentSlotResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

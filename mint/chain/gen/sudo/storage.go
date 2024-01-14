package sudo

import (
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for Key id={{false []}}
//
//	The `AccountId` of the sudo key.
func MakeKeyStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Sudo", "Key")
}
func GetKey(state state.State, bhash types.Hash) (ret [32]byte, isSome bool, err error) {
	key, err := MakeKeyStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetKeyLatest(state state.State) (ret [32]byte, isSome bool, err error) {
	key, err := MakeKeyStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

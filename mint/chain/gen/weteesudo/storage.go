package weteesudo

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for Account
//
//	WETEE Root account id.
//	组织最高权限 id
func MakeAccountStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeSudo", "Account", byteArgs...)
}
func GetAccount(state state.State, bhash types.Hash, uint640 uint64) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAccountStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAccountLatest(state state.State, uint640 uint64) (ret [32]byte, isSome bool, err error) {
	key, err := MakeAccountStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for CloseDao
//
//	WETEE Root account id.
//	组织最高权限 id
func MakeCloseDaoStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeSudo", "CloseDao", byteArgs...)
}
func GetCloseDao(state state.State, bhash types.Hash, uint640 uint64) (ret bool, isSome bool, err error) {
	key, err := MakeCloseDaoStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetCloseDaoLatest(state state.State, uint640 uint64) (ret bool, isSome bool, err error) {
	key, err := MakeCloseDaoStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for SudoTasks
//
//	sudo模块调用历史
func MakeSudoTasksStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeSudo", "SudoTasks", byteArgs...)
}

var SudoTasksResultDefaultBytes, _ = hex.DecodeString("00")

func GetSudoTasks(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.SudoTask, err error) {
	key, err := MakeSudoTasksStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(SudoTasksResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetSudoTasksLatest(state state.State, uint640 uint64) (ret []types1.SudoTask, err error) {
	key, err := MakeSudoTasksStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(SudoTasksResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

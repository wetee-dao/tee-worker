package weteeasset

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for DaoAssetsInfo
func MakeDaoAssetsInfoStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeAsset", "DaoAssetsInfo", byteArgs...)
}
func GetDaoAssetsInfo(state state.State, bhash types.Hash, uint640 uint64) (ret types1.DaoAssetInfo, isSome bool, err error) {
	key, err := MakeDaoAssetsInfoStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetDaoAssetsInfoLatest(state state.State, uint640 uint64) (ret types1.DaoAssetInfo, isSome bool, err error) {
	key, err := MakeDaoAssetsInfoStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for UsersNumber
func MakeUsersNumberStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeAsset", "UsersNumber", byteArgs...)
}

var UsersNumberResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetUsersNumber(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, err error) {
	key, err := MakeUsersNumberStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UsersNumberResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetUsersNumberLatest(state state.State, uint640 uint64) (ret uint32, err error) {
	key, err := MakeUsersNumberStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(UsersNumberResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ExistentDeposits
func MakeExistentDepositsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeAsset", "ExistentDeposits", byteArgs...)
}

var ExistentDepositsResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000")

func GetExistentDeposits(state state.State, bhash types.Hash, uint640 uint64) (ret types.U128, err error) {
	key, err := MakeExistentDepositsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ExistentDepositsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetExistentDepositsLatest(state state.State, uint640 uint64) (ret types.U128, err error) {
	key, err := MakeExistentDepositsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ExistentDepositsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

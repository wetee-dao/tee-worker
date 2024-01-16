package randomnesscollectiveflip

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for RandomMaterial id={{false [176]}}
//
//	Series of block headers from the last 81 blocks that acts as random seed material. This
//	is arranged as a ring buffer with `block_number % 81` being the index into the `Vec` of
//	the oldest hash.
func MakeRandomMaterialStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "RandomnessCollectiveFlip", "RandomMaterial")
}

var RandomMaterialResultDefaultBytes, _ = hex.DecodeString("00")

func GetRandomMaterial(state state.State, bhash types.Hash) (ret [][32]byte, err error) {
	key, err := MakeRandomMaterialStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(RandomMaterialResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetRandomMaterialLatest(state state.State) (ret [][32]byte, err error) {
	key, err := MakeRandomMaterialStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(RandomMaterialResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

package messagequeue

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types "wetee.app/worker/internal/worker/chain/gen/types"
)

// Make a storage key for BookStateFor
//
//	The index of the first and last (non-empty) pages.
func MakeBookStateForStorageKey(messageOrigin0 types.MessageOrigin) (types1.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(messageOrigin0)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types1.CreateStorageKey(&types.Meta, "MessageQueue", "BookStateFor", byteArgs...)
}

var BookStateForResultDefaultBytes, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000")

func GetBookStateFor(state state.State, bhash types1.Hash, messageOrigin0 types.MessageOrigin) (ret types.BookState, err error) {
	key, err := MakeBookStateForStorageKey(messageOrigin0)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BookStateForResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetBookStateForLatest(state state.State, messageOrigin0 types.MessageOrigin) (ret types.BookState, err error) {
	key, err := MakeBookStateForStorageKey(messageOrigin0)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(BookStateForResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ServiceHead id={{false [45]}}
//
//	The origin at which we should begin servicing.
func MakeServiceHeadStorageKey() (types1.StorageKey, error) {
	return types1.CreateStorageKey(&types.Meta, "MessageQueue", "ServiceHead")
}
func GetServiceHead(state state.State, bhash types1.Hash) (ret types.MessageOrigin, isSome bool, err error) {
	key, err := MakeServiceHeadStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetServiceHeadLatest(state state.State) (ret types.MessageOrigin, isSome bool, err error) {
	key, err := MakeServiceHeadStorageKey()
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Pages
//
//	The map of page indices to pages.
func MakePagesStorageKey(tupleOfMessageOriginUint320 types.MessageOrigin, tupleOfMessageOriginUint321 uint32) (types1.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfMessageOriginUint320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfMessageOriginUint321)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types1.CreateStorageKey(&types.Meta, "MessageQueue", "Pages", byteArgs...)
}
func GetPages(state state.State, bhash types1.Hash, tupleOfMessageOriginUint320 types.MessageOrigin, tupleOfMessageOriginUint321 uint32) (ret types.Page, isSome bool, err error) {
	key, err := MakePagesStorageKey(tupleOfMessageOriginUint320, tupleOfMessageOriginUint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetPagesLatest(state state.State, tupleOfMessageOriginUint320 types.MessageOrigin, tupleOfMessageOriginUint321 uint32) (ret types.Page, isSome bool, err error) {
	key, err := MakePagesStorageKey(tupleOfMessageOriginUint320, tupleOfMessageOriginUint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

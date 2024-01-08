package tokens

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/internal/worker/chain/gen/types"
)

// Make a storage key for TotalIssuance
//
//	The total issuance of a token type.
func MakeTotalIssuanceStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Tokens", "TotalIssuance", byteArgs...)
}

var TotalIssuanceResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000")

func GetTotalIssuance(state state.State, bhash types.Hash, uint640 uint64) (ret types.U128, err error) {
	key, err := MakeTotalIssuanceStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(TotalIssuanceResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetTotalIssuanceLatest(state state.State, uint640 uint64) (ret types.U128, err error) {
	key, err := MakeTotalIssuanceStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(TotalIssuanceResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Locks
//
//	Any liquidity locks of a token type under an account.
//	NOTE: Should only be accessed when setting, changing and freeing a lock.
func MakeLocksStorageKey(tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6410)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6411)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Tokens", "Locks", byteArgs...)
}

var LocksResultDefaultBytes, _ = hex.DecodeString("00")

func GetLocks(state state.State, bhash types.Hash, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret []types1.BalanceLock1, err error) {
	key, err := MakeLocksStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(LocksResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetLocksLatest(state state.State, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret []types1.BalanceLock1, err error) {
	key, err := MakeLocksStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(LocksResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Accounts
//
//	The balance of a token type under an account.
//
//	NOTE: If the total is ever zero, decrease account ref account.
//
//	NOTE: This is only used in the case that this module is used to store
//	balances.
func MakeAccountsStorageKey(tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6410)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6411)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Tokens", "Accounts", byteArgs...)
}

var AccountsResultDefaultBytes, _ = hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

func GetAccounts(state state.State, bhash types.Hash, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret types1.AccountData1, err error) {
	key, err := MakeAccountsStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetAccountsLatest(state state.State, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret types1.AccountData1, err error) {
	key, err := MakeAccountsStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(AccountsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Reserves
//
//	Named reserves on some account balances.
func MakeReservesStorageKey(tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6410)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfByteArray32Uint6411)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Tokens", "Reserves", byteArgs...)
}

var ReservesResultDefaultBytes, _ = hex.DecodeString("00")

func GetReserves(state state.State, bhash types.Hash, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret []types1.ReserveDataReserveIdentifierByteArray8, err error) {
	key, err := MakeReservesStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ReservesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetReservesLatest(state state.State, tupleOfByteArray32Uint6410 [32]byte, tupleOfByteArray32Uint6411 uint64) (ret []types1.ReserveDataReserveIdentifierByteArray8, err error) {
	key, err := MakeReservesStorageKey(tupleOfByteArray32Uint6410, tupleOfByteArray32Uint6411)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ReservesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

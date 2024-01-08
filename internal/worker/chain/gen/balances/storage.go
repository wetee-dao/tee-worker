package balances

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/internal/worker/chain/gen/types"
)

// Make a storage key for TotalIssuance id={{false [6]}}
//
//	The total units issued in the system.
func MakeTotalIssuanceStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Balances", "TotalIssuance")
}

var TotalIssuanceResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000")

func GetTotalIssuance(state state.State, bhash types.Hash) (ret types.U128, err error) {
	key, err := MakeTotalIssuanceStorageKey()
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
func GetTotalIssuanceLatest(state state.State) (ret types.U128, err error) {
	key, err := MakeTotalIssuanceStorageKey()
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

// Make a storage key for InactiveIssuance id={{false [6]}}
//
//	The total units of outstanding deactivated balance in the system.
func MakeInactiveIssuanceStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "Balances", "InactiveIssuance")
}

var InactiveIssuanceResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000")

func GetInactiveIssuance(state state.State, bhash types.Hash) (ret types.U128, err error) {
	key, err := MakeInactiveIssuanceStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(InactiveIssuanceResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetInactiveIssuanceLatest(state state.State) (ret types.U128, err error) {
	key, err := MakeInactiveIssuanceStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(InactiveIssuanceResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Account
//
//	The Balances pallet example of storing the balance of an account.
//
//	# Example
//
//	```nocompile
//	 impl pallet_balances::Config for Runtime {
//	   type AccountStore = StorageMapShim<Self::Account<Runtime>, frame_system::Provider<Runtime>, AccountId, Self::AccountData<Balance>>
//	 }
//	```
//
//	You can also store the balance of an account in the `System` pallet.
//
//	# Example
//
//	```nocompile
//	 impl pallet_balances::Config for Runtime {
//	  type AccountStore = System
//	 }
//	```
//
//	But this comes with tradeoffs, storing account balances in the system pallet stores
//	`frame_system` data alongside the account data contrary to storing account balances in the
//	`Balances` pallet, which uses a `StorageMap` to store balances data only.
//	NOTE: This is only used in the case that this pallet is used to store balances.
func MakeAccountStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Balances", "Account", byteArgs...)
}

var AccountResultDefaultBytes, _ = hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080")

func GetAccount(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret types1.AccountData, err error) {
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
func GetAccountLatest(state state.State, byteArray320 [32]byte) (ret types1.AccountData, err error) {
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

// Make a storage key for Locks
//
//	Any liquidity locks on some account balances.
//	NOTE: Should only be accessed when setting, changing and freeing a lock.
func MakeLocksStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Balances", "Locks", byteArgs...)
}

var LocksResultDefaultBytes, _ = hex.DecodeString("00")

func GetLocks(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.BalanceLock, err error) {
	key, err := MakeLocksStorageKey(byteArray320)
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
func GetLocksLatest(state state.State, byteArray320 [32]byte) (ret []types1.BalanceLock, err error) {
	key, err := MakeLocksStorageKey(byteArray320)
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

// Make a storage key for Reserves
//
//	Named reserves on some account balances.
func MakeReservesStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Balances", "Reserves", byteArgs...)
}

var ReservesResultDefaultBytes, _ = hex.DecodeString("00")

func GetReserves(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.ReserveData, err error) {
	key, err := MakeReservesStorageKey(byteArray320)
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
func GetReservesLatest(state state.State, byteArray320 [32]byte) (ret []types1.ReserveData, err error) {
	key, err := MakeReservesStorageKey(byteArray320)
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

// Make a storage key for Holds
//
//	Holds on account balances.
func MakeHoldsStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Balances", "Holds", byteArgs...)
}

var HoldsResultDefaultBytes, _ = hex.DecodeString("00")

func GetHolds(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.IdAmount, err error) {
	key, err := MakeHoldsStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(HoldsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetHoldsLatest(state state.State, byteArray320 [32]byte) (ret []types1.IdAmount, err error) {
	key, err := MakeHoldsStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(HoldsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Freezes
//
//	Freeze locks on account balances.
func MakeFreezesStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "Balances", "Freezes", byteArgs...)
}

var FreezesResultDefaultBytes, _ = hex.DecodeString("00")

func GetFreezes(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.IdAmount, err error) {
	key, err := MakeFreezesStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(FreezesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetFreezesLatest(state state.State, byteArray320 [32]byte) (ret []types1.IdAmount, err error) {
	key, err := MakeFreezesStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(FreezesResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

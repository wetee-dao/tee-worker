package weteegov

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/internal/worker/chain/gen/types"
)

// Make a storage key for PrePropCount
//
//	Number of public proposals so for.
func MakePrePropCountStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "PrePropCount", byteArgs...)
}

var PrePropCountResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetPrePropCount(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, err error) {
	key, err := MakePrePropCountStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PrePropCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetPrePropCountLatest(state state.State, uint640 uint64) (ret uint32, err error) {
	key, err := MakePrePropCountStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PrePropCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for MaxPreProps
//
//	Maximum number of public proposals at one time.
func MakeMaxPrePropsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "MaxPreProps", byteArgs...)
}

var MaxPrePropsResultDefaultBytes, _ = hex.DecodeString("64000000")

func GetMaxPreProps(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, err error) {
	key, err := MakeMaxPrePropsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MaxPrePropsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetMaxPrePropsLatest(state state.State, uint640 uint64) (ret uint32, err error) {
	key, err := MakeMaxPrePropsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MaxPrePropsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Periods
//
//	投票轨道
func MakePeriodsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "Periods", byteArgs...)
}

var PeriodsResultDefaultBytes, _ = hex.DecodeString("00")

func GetPeriods(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.Period, err error) {
	key, err := MakePeriodsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PeriodsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetPeriodsLatest(state state.State, uint640 uint64) (ret []types1.Period, err error) {
	key, err := MakePeriodsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PeriodsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for DefaultPeriods id={{false [221]}}
//
//	投票轨道
func MakeDefaultPeriodsStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "DefaultPeriods")
}

var DefaultPeriodsResultDefaultBytes, _ = hex.DecodeString("00")

func GetDefaultPeriods(state state.State, bhash types.Hash) (ret []types1.Period, err error) {
	key, err := MakeDefaultPeriodsStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DefaultPeriodsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetDefaultPeriodsLatest(state state.State) (ret []types1.Period, err error) {
	key, err := MakeDefaultPeriodsStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DefaultPeriodsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for PreProps
//
//	The public proposals.
//	Unsorted.
//	The second item is the proposal's hash.
func MakePrePropsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "PreProps", byteArgs...)
}

var PrePropsResultDefaultBytes, _ = hex.DecodeString("00")

func GetPreProps(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.PreProp, err error) {
	key, err := MakePrePropsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PrePropsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetPrePropsLatest(state state.State, uint640 uint64) (ret []types1.PreProp, err error) {
	key, err := MakePrePropsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PrePropsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for DepositOf
//
//	提案
//	Those who have locked a deposit.
//	TWOX-NOTE: Safe, as increasing integer keys are safe.
func MakeDepositOfStorageKey(tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64Uint320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64Uint321)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "DepositOf", byteArgs...)
}
func GetDepositOf(state state.State, bhash types.Hash, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret types1.TupleOfByteArray32SliceU128, isSome bool, err error) {
	key, err := MakeDepositOfStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetDepositOfLatest(state state.State, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret types1.TupleOfByteArray32SliceU128, isSome bool, err error) {
	key, err := MakeDepositOfStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for Props
//
//	全民投票
//	Prop specific information.
func MakePropsStorageKey(tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64Uint320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64Uint321)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "Props", byteArgs...)
}
func GetProps(state state.State, bhash types.Hash, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret types1.Prop, isSome bool, err error) {
	key, err := MakePropsStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetPropsLatest(state state.State, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret types1.Prop, isSome bool, err error) {
	key, err := MakePropsStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for ReserveOf
//
//	Amount of proposal locked.
func MakeReserveOfStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "ReserveOf", byteArgs...)
}

var ReserveOfResultDefaultBytes, _ = hex.DecodeString("00")

func GetReserveOf(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.TupleOfU128Uint64, err error) {
	key, err := MakeReserveOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ReserveOfResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetReserveOfLatest(state state.State, byteArray320 [32]byte) (ret []types1.TupleOfU128Uint64, err error) {
	key, err := MakeReserveOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ReserveOfResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for PropCount
//
//	Number of props so far.
func MakePropCountStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "PropCount", byteArgs...)
}

var PropCountResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetPropCount(state state.State, bhash types.Hash, uint640 uint64) (ret uint32, err error) {
	key, err := MakePropCountStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PropCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetPropCountLatest(state state.State, uint640 uint64) (ret uint32, err error) {
	key, err := MakePropCountStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(PropCountResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for VoteModel
//
//	WETEE 投票模式默认 0，1 TOKEN 1 票
func MakeVoteModelStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "VoteModel", byteArgs...)
}

var VoteModelResultDefaultBytes, _ = hex.DecodeString("00")

func GetVoteModel(state state.State, bhash types.Hash, uint640 uint64) (ret byte, err error) {
	key, err := MakeVoteModelStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(VoteModelResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetVoteModelLatest(state state.State, uint640 uint64) (ret byte, err error) {
	key, err := MakeVoteModelStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(VoteModelResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for VotesOf
//
//	Everyone's voting information.
func MakeVotesOfStorageKey(byteArray320 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(byteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeGov", "VotesOf", byteArgs...)
}

var VotesOfResultDefaultBytes, _ = hex.DecodeString("00")

func GetVotesOf(state state.State, bhash types.Hash, byteArray320 [32]byte) (ret []types1.VoteInfo, err error) {
	key, err := MakeVotesOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(VotesOfResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetVotesOfLatest(state state.State, byteArray320 [32]byte) (ret []types1.VoteInfo, err error) {
	key, err := MakeVotesOfStorageKey(byteArray320)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(VotesOfResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

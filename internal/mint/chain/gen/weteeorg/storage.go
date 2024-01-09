package weteeorg

import (
	"encoding/hex"
	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/internal/mint/chain/gen/types"
)

// Make a storage key for Daos
//
//	All DAOs that have been created.
//	所有组织
func MakeDaosStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "Daos", byteArgs...)
}
func GetDaos(state state.State, bhash types.Hash, uint640 uint64) (ret types1.OrgInfo, isSome bool, err error) {
	key, err := MakeDaosStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetDaosLatest(state state.State, uint640 uint64) (ret types1.OrgInfo, isSome bool, err error) {
	key, err := MakeDaosStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for NextDaoId id={{false [4]}}
//
//	The id of the next dao to be created.
//	获取下一个组织id
func MakeNextDaoIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "NextDaoId")
}

var NextDaoIdResultDefaultBytes, _ = hex.DecodeString("8813000000000000")

func GetNextDaoId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNextDaoIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextDaoIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextDaoIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNextDaoIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextDaoIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for NextAppId id={{false [4]}}
//
//	The id of the next dao to be created.
//	获取下一个组织id
func MakeNextAppIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "NextAppId")
}

var NextAppIdResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetNextAppId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNextAppIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextAppIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextAppIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNextAppIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextAppIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Guilds
//
//	the info of grutypes
//	组织内公会信息
func MakeGuildsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "Guilds", byteArgs...)
}

var GuildsResultDefaultBytes, _ = hex.DecodeString("00")

func GetGuilds(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.GuildInfo, err error) {
	key, err := MakeGuildsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(GuildsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetGuildsLatest(state state.State, uint640 uint64) (ret []types1.GuildInfo, err error) {
	key, err := MakeGuildsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(GuildsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for RoadMaps
//
//	the roadmap info of projects
//	组织内Roadmap信息
func MakeRoadMapsStorageKey(tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (types.StorageKey, error) {
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
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "RoadMaps", byteArgs...)
}

var RoadMapsResultDefaultBytes, _ = hex.DecodeString("00")

func GetRoadMaps(state state.State, bhash types.Hash, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret []types1.QuarterTask, err error) {
	key, err := MakeRoadMapsStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(RoadMapsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetRoadMapsLatest(state state.State, tupleOfUint64Uint320 uint64, tupleOfUint64Uint321 uint32) (ret []types1.QuarterTask, err error) {
	key, err := MakeRoadMapsStorageKey(tupleOfUint64Uint320, tupleOfUint64Uint321)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(RoadMapsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for NextTaskId id={{false [4]}}
//
//	The id of the next dao to be created.
//	获取下一个组织id
func MakeNextTaskIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "NextTaskId")
}

var NextTaskIdResultDefaultBytes, _ = hex.DecodeString("0000000000000000")

func GetNextTaskId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNextTaskIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextTaskIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextTaskIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNextTaskIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextTaskIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Members
//
//	team members
//	团队的成员
func MakeMembersStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "Members", byteArgs...)
}

var MembersResultDefaultBytes, _ = hex.DecodeString("00")

func GetMembers(state state.State, bhash types.Hash, uint640 uint64) (ret [][32]byte, err error) {
	key, err := MakeMembersStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetMembersLatest(state state.State, uint640 uint64) (ret [][32]byte, err error) {
	key, err := MakeMembersStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for GuildMembers
//
//	guild members
//	公会成员
func MakeGuildMembersStorageKey(tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64Uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64Uint641)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "GuildMembers", byteArgs...)
}

var GuildMembersResultDefaultBytes, _ = hex.DecodeString("00")

func GetGuildMembers(state state.State, bhash types.Hash, tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (ret [][32]byte, err error) {
	key, err := MakeGuildMembersStorageKey(tupleOfUint64Uint640, tupleOfUint64Uint641)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(GuildMembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetGuildMembersLatest(state state.State, tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (ret [][32]byte, err error) {
	key, err := MakeGuildMembersStorageKey(tupleOfUint64Uint640, tupleOfUint64Uint641)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(GuildMembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for ProjectMembers
//
//	project members
//	项目成员
func MakeProjectMembersStorageKey(tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64Uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64Uint641)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "ProjectMembers", byteArgs...)
}

var ProjectMembersResultDefaultBytes, _ = hex.DecodeString("00")

func GetProjectMembers(state state.State, bhash types.Hash, tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (ret [][32]byte, err error) {
	key, err := MakeProjectMembersStorageKey(tupleOfUint64Uint640, tupleOfUint64Uint641)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ProjectMembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetProjectMembersLatest(state state.State, tupleOfUint64Uint640 uint64, tupleOfUint64Uint641 uint64) (ret [][32]byte, err error) {
	key, err := MakeProjectMembersStorageKey(tupleOfUint64Uint640, tupleOfUint64Uint641)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(ProjectMembersResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for AppHubs
//
//	apps hubs
//	应用中心
func MakeAppHubsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "AppHubs", byteArgs...)
}
func GetAppHubs(state state.State, bhash types.Hash, uint640 uint64) (ret types1.App, isSome bool, err error) {
	key, err := MakeAppHubsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetAppHubsLatest(state state.State, uint640 uint64) (ret types1.App, isSome bool, err error) {
	key, err := MakeAppHubsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

// Make a storage key for OrgApps
//
//	org apps
//	应用中心
func MakeOrgAppsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "OrgApps", byteArgs...)
}

var OrgAppsResultDefaultBytes, _ = hex.DecodeString("00")

func GetOrgApps(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.OrgApp, err error) {
	key, err := MakeOrgAppsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(OrgAppsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetOrgAppsLatest(state state.State, uint640 uint64) (ret []types1.OrgApp, err error) {
	key, err := MakeOrgAppsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(OrgAppsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for MemberPoint
//
//	point
//	成员贡献点
func MakeMemberPointStorageKey(tupleOfUint64ByteArray320 uint64, tupleOfUint64ByteArray321 [32]byte) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(tupleOfUint64ByteArray320)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	encBytes, err = codec.Encode(tupleOfUint64ByteArray321)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeOrg", "MemberPoint", byteArgs...)
}

var MemberPointResultDefaultBytes, _ = hex.DecodeString("00000000")

func GetMemberPoint(state state.State, bhash types.Hash, tupleOfUint64ByteArray320 uint64, tupleOfUint64ByteArray321 [32]byte) (ret uint32, err error) {
	key, err := MakeMemberPointStorageKey(tupleOfUint64ByteArray320, tupleOfUint64ByteArray321)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MemberPointResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetMemberPointLatest(state state.State, tupleOfUint64ByteArray320 uint64, tupleOfUint64ByteArray321 [32]byte) (ret uint32, err error) {
	key, err := MakeMemberPointStorageKey(tupleOfUint64ByteArray320, tupleOfUint64ByteArray321)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(MemberPointResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

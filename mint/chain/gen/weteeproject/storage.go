package weteeproject

import (
	"encoding/hex"

	state "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	codec "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	types1 "wetee.app/worker/mint/chain/gen/types"
)

// Make a storage key for NextProjectId id={{false [4]}}
//
//	The id of the next dao to be created.
//	获取下一个组织id
func MakeNextProjectIdStorageKey() (types.StorageKey, error) {
	return types.CreateStorageKey(&types1.Meta, "WeteeProject", "NextProjectId")
}

var NextProjectIdResultDefaultBytes, _ = hex.DecodeString("0100000000000000")

func GetNextProjectId(state state.State, bhash types.Hash) (ret uint64, err error) {
	key, err := MakeNextProjectIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextProjectIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetNextProjectIdLatest(state state.State) (ret uint64, err error) {
	key, err := MakeNextProjectIdStorageKey()
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(NextProjectIdResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for DaoProjects
//
//	project board
//	项目看板
func MakeDaoProjectsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeProject", "DaoProjects", byteArgs...)
}

var DaoProjectsResultDefaultBytes, _ = hex.DecodeString("00")

func GetDaoProjects(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.ProjectInfo, err error) {
	key, err := MakeDaoProjectsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DaoProjectsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetDaoProjectsLatest(state state.State, uint640 uint64) (ret []types1.ProjectInfo, err error) {
	key, err := MakeDaoProjectsStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(DaoProjectsResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}

// Make a storage key for Tasks
//
//	project task
//	任务看板
func MakeTasksStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeProject", "Tasks", byteArgs...)
}

var TasksResultDefaultBytes, _ = hex.DecodeString("00")

func GetTasks(state state.State, bhash types.Hash, uint640 uint64) (ret []types1.TaskInfo, err error) {
	key, err := MakeTasksStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(TasksResultDefaultBytes, &ret)
		if err != nil {
			return
		}
	}
	return
}
func GetTasksLatest(state state.State, uint640 uint64) (ret []types1.TaskInfo, err error) {
	key, err := MakeTasksStorageKey(uint640)
	if err != nil {
		return
	}
	var isSome bool
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	if !isSome {
		err = codec.Decode(TasksResultDefaultBytes, &ret)
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
	return types.CreateStorageKey(&types1.Meta, "WeteeProject", "NextTaskId")
}

var NextTaskIdResultDefaultBytes, _ = hex.DecodeString("0100000000000000")

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

// Make a storage key for TaskReviews
//
//	TODO taskDone
//	已完成项目
//	task reviews
//	项目审核报告
func MakeTaskReviewsStorageKey(uint640 uint64) (types.StorageKey, error) {
	byteArgs := [][]byte{}
	encBytes := []byte{}
	var err error
	encBytes, err = codec.Encode(uint640)
	if err != nil {
		return nil, err
	}
	byteArgs = append(byteArgs, encBytes)
	return types.CreateStorageKey(&types1.Meta, "WeteeProject", "TaskReviews", byteArgs...)
}
func GetTaskReviews(state state.State, bhash types.Hash, uint640 uint64) (ret types1.ReviewStatus, isSome bool, err error) {
	key, err := MakeTaskReviewsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorage(key, &ret, bhash)
	if err != nil {
		return
	}
	return
}
func GetTaskReviewsLatest(state state.State, uint640 uint64) (ret types1.ReviewStatus, isSome bool, err error) {
	key, err := MakeTaskReviewsStorageKey(uint640)
	if err != nil {
		return
	}
	isSome, err = state.GetStorageLatest(key, &ret)
	if err != nil {
		return
	}
	return
}

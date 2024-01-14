package weteeproject

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "wetee.app/worker/mint/chain/gen/types"
)

// See [`Pallet::project_join_request`].
func MakeProjectJoinRequestCall(daoId0 uint64, projectId1 uint64, who2 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsProjectJoinRequest:           true,
			AsProjectJoinRequestDaoId0:     daoId0,
			AsProjectJoinRequestProjectId1: projectId1,
			AsProjectJoinRequestWho2:       who2,
		},
	}
}

// See [`Pallet::create_project`].
func MakeCreateProjectCall(daoId0 uint64, name1 []byte, description2 []byte, creator3 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsCreateProject:             true,
			AsCreateProjectDaoId0:       daoId0,
			AsCreateProjectName1:        name1,
			AsCreateProjectDescription2: description2,
			AsCreateProjectCreator3:     creator3,
		},
	}
}

// See [`Pallet::apply_project_funds`].
func MakeApplyProjectFundsCall(daoId0 uint64, projectId1 uint64, amount2 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsApplyProjectFunds:           true,
			AsApplyProjectFundsDaoId0:     daoId0,
			AsApplyProjectFundsProjectId1: projectId1,
			AsApplyProjectFundsAmount2:    amount2,
		},
	}
}

// See [`Pallet::create_task`].
func MakeCreateTaskCall(daoId0 uint64, projectId1 uint64, name2 []byte, description3 []byte, point4 uint16, priority5 byte, maxAssignee6 types.OptionTByte, skills7 types.OptionTByteSlice, assignees8 types.OptionTByteArray32Slice, reviewers9 types.OptionTByteArray32Slice, amount10 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsCreateTask:             true,
			AsCreateTaskDaoId0:       daoId0,
			AsCreateTaskProjectId1:   projectId1,
			AsCreateTaskName2:        name2,
			AsCreateTaskDescription3: description3,
			AsCreateTaskPoint4:       point4,
			AsCreateTaskPriority5:    priority5,
			AsCreateTaskMaxAssignee6: maxAssignee6,
			AsCreateTaskSkills7:      skills7,
			AsCreateTaskAssignees8:   assignees8,
			AsCreateTaskReviewers9:   reviewers9,
			AsCreateTaskAmount10:     amount10,
		},
	}
}

// See [`Pallet::join_task`].
func MakeJoinTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsJoinTask:           true,
			AsJoinTaskDaoId0:     daoId0,
			AsJoinTaskProjectId1: projectId1,
			AsJoinTaskTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::leave_task`].
func MakeLeaveTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsLeaveTask:           true,
			AsLeaveTaskDaoId0:     daoId0,
			AsLeaveTaskProjectId1: projectId1,
			AsLeaveTaskTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::join_task_review`].
func MakeJoinTaskReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsJoinTaskReview:           true,
			AsJoinTaskReviewDaoId0:     daoId0,
			AsJoinTaskReviewProjectId1: projectId1,
			AsJoinTaskReviewTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::leave_task_review`].
func MakeLeaveTaskReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsLeaveTaskReview:           true,
			AsLeaveTaskReviewDaoId0:     daoId0,
			AsLeaveTaskReviewProjectId1: projectId1,
			AsLeaveTaskReviewTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::start_task`].
func MakeStartTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsStartTask:           true,
			AsStartTaskDaoId0:     daoId0,
			AsStartTaskProjectId1: projectId1,
			AsStartTaskTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::request_review`].
func MakeRequestReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsRequestReview:           true,
			AsRequestReviewDaoId0:     daoId0,
			AsRequestReviewProjectId1: projectId1,
			AsRequestReviewTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::task_done`].
func MakeTaskDoneCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsTaskDone:           true,
			AsTaskDoneDaoId0:     daoId0,
			AsTaskDoneProjectId1: projectId1,
			AsTaskDoneTaskId2:    taskId2,
		},
	}
}

// See [`Pallet::make_review`].
func MakeMakeReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64, opinion3 types.ReviewOpinion, meta4 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeProject: true,
		AsWeteeProjectField0: &types.WeteeProjectPalletCall{
			IsMakeReview:           true,
			AsMakeReviewDaoId0:     daoId0,
			AsMakeReviewProjectId1: projectId1,
			AsMakeReviewTaskId2:    taskId2,
			AsMakeReviewOpinion3:   opinion3,
			AsMakeReviewMeta4:      meta4,
		},
	}
}

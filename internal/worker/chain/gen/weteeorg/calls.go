package weteeorg

import types "wetee.app/worker/internal/worker/chain/gen/types"

// See [`Pallet::create_dao`].
func MakeCreateDaoCall(name0 []byte, desc1 []byte, purpose2 []byte, metaData3 []byte, imApi4 []byte, bg5 []byte, logo6 []byte, img7 []byte, homeUrl8 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsCreateDao:          true,
			AsCreateDaoName0:     name0,
			AsCreateDaoDesc1:     desc1,
			AsCreateDaoPurpose2:  purpose2,
			AsCreateDaoMetaData3: metaData3,
			AsCreateDaoImApi4:    imApi4,
			AsCreateDaoBg5:       bg5,
			AsCreateDaoLogo6:     logo6,
			AsCreateDaoImg7:      img7,
			AsCreateDaoHomeUrl8:  homeUrl8,
		},
	}
}

// See [`Pallet::update_dao`].
func MakeUpdateDaoCall(daoId0 uint64, name1 types.OptionTByteSlice, desc2 types.OptionTByteSlice, purpose3 types.OptionTByteSlice, metaData4 types.OptionTByteSlice, imApi5 types.OptionTByteSlice, bg6 types.OptionTByteSlice, logo7 types.OptionTByteSlice, img8 types.OptionTByteSlice, homeUrl9 types.OptionTByteSlice, status10 types.OptionTStatus) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsUpdateDao:          true,
			AsUpdateDaoDaoId0:    daoId0,
			AsUpdateDaoName1:     name1,
			AsUpdateDaoDesc2:     desc2,
			AsUpdateDaoPurpose3:  purpose3,
			AsUpdateDaoMetaData4: metaData4,
			AsUpdateDaoImApi5:    imApi5,
			AsUpdateDaoBg6:       bg6,
			AsUpdateDaoLogo7:     logo7,
			AsUpdateDaoImg8:      img8,
			AsUpdateDaoHomeUrl9:  homeUrl9,
			AsUpdateDaoStatus10:  status10,
		},
	}
}

// See [`Pallet::create_roadmap_task`].
func MakeCreateRoadmapTaskCall(daoId0 uint64, roadmapId1 uint32, name2 []byte, priority3 byte, tags4 types.OptionTByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsCreateRoadmapTask:           true,
			AsCreateRoadmapTaskDaoId0:     daoId0,
			AsCreateRoadmapTaskRoadmapId1: roadmapId1,
			AsCreateRoadmapTaskName2:      name2,
			AsCreateRoadmapTaskPriority3:  priority3,
			AsCreateRoadmapTaskTags4:      tags4,
		},
	}
}

// See [`Pallet::update_roadmap_task`].
func MakeUpdateRoadmapTaskCall(daoId0 uint64, roadmapId1 uint32, taskId2 uint64, priority3 byte, status4 byte, tags5 types.OptionTByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsUpdateRoadmapTask:           true,
			AsUpdateRoadmapTaskDaoId0:     daoId0,
			AsUpdateRoadmapTaskRoadmapId1: roadmapId1,
			AsUpdateRoadmapTaskTaskId2:    taskId2,
			AsUpdateRoadmapTaskPriority3:  priority3,
			AsUpdateRoadmapTaskStatus4:    status4,
			AsUpdateRoadmapTaskTags5:      tags5,
		},
	}
}

// See [`Pallet::create_app`].
func MakeCreateAppCall(name0 []byte, desc1 []byte, icon2 []byte, url3 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsCreateApp:      true,
			AsCreateAppName0: name0,
			AsCreateAppDesc1: desc1,
			AsCreateAppIcon2: icon2,
			AsCreateAppUrl3:  url3,
		},
	}
}

// See [`Pallet::update_app_status`].
func MakeUpdateAppStatusCall(appId0 uint64, status1 types.Status) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsUpdateAppStatus:        true,
			AsUpdateAppStatusAppId0:  appId0,
			AsUpdateAppStatusStatus1: status1,
		},
	}
}

// See [`Pallet::org_integrate_app`].
func MakeOrgIntegrateAppCall(daoId0 uint64, appId1 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsOrgIntegrateApp:       true,
			AsOrgIntegrateAppDaoId0: daoId0,
			AsOrgIntegrateAppAppId1: appId1,
		},
	}
}

// See [`Pallet::update_org_app_status`].
func MakeUpdateOrgAppStatusCall(daoId0 uint64, appId1 uint64, status2 types.Status) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeOrg: true,
		AsWeteeOrgField0: &types.WeteeOrgPalletCall{
			IsUpdateOrgAppStatus:        true,
			AsUpdateOrgAppStatusDaoId0:  daoId0,
			AsUpdateOrgAppStatusAppId1:  appId1,
			AsUpdateOrgAppStatusStatus2: status2,
		},
	}
}

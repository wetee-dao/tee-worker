package mint

import (
	"context"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithEvent(event types.EventRecord, clusterId uint64) error {
	e := event.Event
	ctx := context.Background()
	var err error

	// 处理任务消息
	// Handling Worker Messages
	if e.IsWeteeWorker {
		startEvent := e.AsWeteeWorkerField0
		if startEvent.IsWorkRuning {
			workId := startEvent.AsWorkRuningWorkId1
			user := startEvent.AsWorkRuningUser0
			cid := startEvent.AsWorkRuningClusterId2
			if cid == clusterId {
				version, _ := chain.GetVersion(m.ChainClient, workId)
				envs, err := m.GetEnvs(workId)
				if err != nil {
					return err
				}

				if workId.Wtype.IsAPP {
					appIns := chain.App{
						Client: m.ChainClient,
					}
					app, _ := appIns.GetApp(user[:], workId.Id)
					err = m.CreateApp(&ctx, user[:], workId, app, envs, version)
					util.LogWithRed("===========================================CreateOrUpdateApp error: ", err)
				} else {
					taskIns := chain.Task{
						Client: m.ChainClient,
					}
					task, _ := taskIns.GetTask(user[:], workId.Id)
					err = m.CreateTask(&ctx, user[:], workId, task, envs, version)
					util.LogWithRed("===========================================CreateOrUpdateTask error: ", err)
				}
			}
		}
	}

	// 处理机密应用消息
	// Handling App Messages
	if e.IsWeteeApp {
		appEvent := e.AsWeteeAppField0
		if appEvent.IsWorkStopped {
			workId := appEvent.AsWorkStoppedWorkId1

			err = m.StopApp(workId, "")
			util.LogWithRed("===========================================StopPod error: ", err)
		}
		if appEvent.IsWorkUpdated {
			workId := appEvent.AsWorkUpdatedWorkId1
			user := appEvent.AsWorkUpdatedUser0

			util.LogWithRed("===========================================WorkUpdated: ", workId)
			version, _ := chain.GetVersion(m.ChainClient, workId)
			appIns := chain.App{
				Client: m.ChainClient,
			}
			app, _ := appIns.GetApp(user[:], workId.Id)
			envs, _ := m.GetEnvs(workId)
			err = m.UpdateApp(&ctx, user[:], workId, app, envs, version)
			util.LogWithRed("===========================================CreateOrUpdatePod error: ", err)
		}
	}

	if e.IsWeteeTask {
		taskEvent := e.AsWeteeTaskField0
		if taskEvent.IsTaskStop {
			taskID := taskEvent.AsTaskStopId1
			workId := types.WorkId{Wtype: types.WorkType{
				IsTASK: true,
			}, Id: taskID}

			err = m.StopApp(workId, "")
			util.LogWithRed("===========================================StopTask error: ", err)
		}
	}

	return err
}

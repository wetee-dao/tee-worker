package mint

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wetee-dao/go-sdk/module"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithEvent(event types.EventRecord, clusterId uint64) error {
	e := event.Event
	ctx := context.Background()
	var err error

	// 处理任务消息
	// Handling Worker Messages
	if e.IsWeTEEWorker {
		startEvent := e.AsWeTEEWorkerField0
		if startEvent.IsWorkRuning {
			workId := startEvent.AsWorkRuningWorkId1
			user := startEvent.AsWorkRuningUser0
			cid := startEvent.AsWorkRuningClusterId2
			if cid == clusterId {
				version, _ := module.GetVersion(m.ChainClient, workId)
				settings, err := m.GetSettingsFromWork(workId, nil)
				if err != nil {
					return errors.Wrap(err, "GetSettingsFromWork error")
				}

				if workId.Wtype.IsAPP {
					appIns := module.App{
						Client: m.ChainClient,
					}
					app, _ := appIns.GetApp(user[:], workId.Id)
					err = m.CreateApp(&ctx, user[:], workId, app, settings, version)
					util.LogError("===========================================CreateOrUpdateApp error: ", err)
				} else if workId.Wtype.IsGPU {
					gpuIns := module.GpuApp{
						Client: m.ChainClient,
					}
					gpu, _ := gpuIns.GetApp(user[:], workId.Id)
					err = m.CreateGpuApp(&ctx, user[:], workId, gpu, settings, version)
					util.LogError("===========================================CreateOrUpdateGpuApp error: ", err)
				} else {
					taskIns := module.Task{
						Client: m.ChainClient,
					}
					task, _ := taskIns.GetTask(user[:], workId.Id)
					err = m.CreateTask(&ctx, user[:], workId, task, settings, version)
					util.LogError("===========================================CreateOrUpdateTask error: ", err)
				}
			}
		}
	}

	// 处理机密应用消息
	// Handling App Messages
	if e.IsWeTEEApp {
		appEvent := e.AsWeTEEAppField0
		if appEvent.IsWorkStopped {
			workId := appEvent.AsWorkStoppedWorkId1

			err = m.StopApp(workId, "")
			util.LogError("===========================================StopPod error: ", err)
		}
		if appEvent.IsWorkUpdated {
			workId := appEvent.AsWorkUpdatedWorkId1
			user := appEvent.AsWorkUpdatedUser0

			util.LogError("===========================================WorkUpdated: ", workId)
			version, _ := module.GetVersion(m.ChainClient, workId)
			appIns := module.App{
				Client: m.ChainClient,
			}
			app, _ := appIns.GetApp(user[:], workId.Id)
			envs, _ := m.BuildEnvs(workId)
			err = m.UpdateApp(&ctx, user[:], workId, app, envs, version)
			util.LogError("===========================================CreateOrUpdatePod error: ", err)
		}
	}

	if e.IsWeTEETask {
		taskEvent := e.AsWeTEETaskField0
		if taskEvent.IsTaskStop {
			taskID := taskEvent.AsTaskStopId1
			workId := types.WorkId{Wtype: types.WorkType{
				IsTASK: true,
			}, Id: taskID}

			err = m.StopApp(workId, "")
			util.LogError("===========================================StopTask error: ", err)
		}
	}

	// 处理CPU应用消息
	// Handling GPU App Messages
	if e.IsWeTEEGpu {
		appEvent := e.AsWeTEEGpuField0
		if appEvent.IsWorkStopped {
			workId := appEvent.AsWorkStoppedWorkId1

			err = m.StopApp(workId, "")
			util.LogError("===========================================StopPod error: ", err)
		}
		if appEvent.IsWorkUpdated {
			workId := appEvent.AsWorkUpdatedWorkId1
			user := appEvent.AsWorkUpdatedUser0

			util.LogError("===========================================WorkUpdated: ", workId)
			version, _ := module.GetVersion(m.ChainClient, workId)
			appIns := module.GpuApp{
				Client: m.ChainClient,
			}
			app, _ := appIns.GetApp(user[:], workId.Id)
			envs, _ := m.BuildEnvs(workId)
			err = m.UpdateGpuApp(&ctx, user[:], workId, app, envs, version)
			util.LogError("===========================================CreateOrUpdatePod error: ", err)
		}
	}

	return err
}

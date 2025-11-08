package store

import (
	"encoding/hex"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/util"
)

const PodBucket = "pod"
const PodKey = "run_"

// SetPod
func SetPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.SetJson(PodBucket, key, &pod)
}

// SetPodLastMint
func SetPodLastMint(podId uint64, last uint32) error {
	pod, _ := GetPod(podId)
	pod.LastMintBlockNumber = last

	return SetPod(*pod)
}

// SetPodSkipUtil skip mint util
func SetPodSkipUtil(podId uint64, t uint32, last uint32) error {
	pod, _ := GetPod(podId)
	pod.SkipUtil = t
	pod.LastMintBlockNumber = last

	return SetPod(*pod)
}

// GetPod
func GetPod(podId uint64) (*model.Pod, error) {
	key := PodKey + fmt.Sprint(podId)
	return model.GetJson[model.Pod](PodBucket, key)
}

// Delete ruing pod
func DelPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.DeleteKey(PodBucket, key)
}

// GetPods
func GetPods() ([]*model.Pod, error) {
	return model.GetJsonList[model.Pod](PodBucket, PodKey)
}

var PendingCallKey = "pending_call_"

// Add pending call
func AddPendingCall(call *model.TeeCall) error {
	key := PendingCallKey + hex.EncodeToString(call.Report)
	return model.SetProtoMessage(PodBucket, key, call)
}

// Get pending call
func GetPendingCall() ([]*model.TeeCall, [][]byte, error) {
	return model.GetProtoMessageList[model.TeeCall](PodBucket, PendingCallKey)
}

// Delete pending call
func DeletePendingCalls(keys [][]byte) error {
	tx := model.DBINS.NewTransaction()
	for _, key := range keys {
		err := tx.Delete(key)
		if err != nil {
			util.LogError("DeletePendingCalls", err)
			return err
		}
	}
	return tx.Commit()
}

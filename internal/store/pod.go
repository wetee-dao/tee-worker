package store

import (
	"encoding/hex"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"
)

const PodBucket = "pod"
const PodKey = "run_"

// SetPod
func SetPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.SetJson(PodBucket, key, &pod)
}

// SetPodLastMint
func SetPodLastMint(podId uint64, t uint32) error {
	pod, _ := GetPod(podId)
	pod.LastMintBlockNumber = t

	return SetPod(*pod)
}

// SetPodSkipUtil skip mint util
func SetPodSkipUtil(podId uint64, t uint32) error {
	pod, _ := GetPod(podId)
	pod.SkipUtil = t

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

// Add pending call
func AddPendingCall(call *model.IndexCall) error {
	hash := blake2b.Sum256(call.Call)
	key := "pending_call_" + hex.EncodeToString(hash[:])
	return model.SetProtoMessage(PodBucket, key, call)
}

// Get pending call
func GetPendingCall() ([]*model.IndexCall, [][]byte, error) {
	key := "pending_call_"
	return model.GetProtoMessageList[model.IndexCall](PodBucket, key)
}

// Delete pending call
func DeletePendingCalls(keys [][]byte) error {
	tx := model.DBINS.NewTransaction()
	for _, key := range keys {
		tx.Delete(key)
	}
	return tx.Commit()
}

package mint

import (
	types "wetee.app/worker/type"
)

// HandleWorker handles incoming messages and branches based on the message type
func (m *Minter) HandleWorker(msg *types.Message) error {
	switch msg.Type {
	/// -------------------- Proof -----------------------
	case "upload_cluster_proof_reply":
		err := m.UploadClusterProofreply(msg.Payload, msg.Error, msg.MsgID, msg.OrgId)
		return err
	/// -------------------- Env -----------------------
	case "reencrypt_secret_remote_reply":
		err := m.ReencryptSecretReply(msg.Payload, msg.Error, msg.MsgID, msg.OrgId)
		return err
	/// -------------------- work -----------------------
	case "work_launch_reply":
		err := m.WorkLaunchReply(msg.Payload, msg.Error, msg.MsgID, msg.OrgId)
		return err
	default:
		return nil
	}
}

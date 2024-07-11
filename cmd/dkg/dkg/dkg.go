package dkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
)

// DKG 代表 Rabin DKG 协议的实例。
type DKG struct {
	// Suite 是加密套件。
	Suite suites.Suite
	// NodeSecret 是长期的私钥。
	NodeSecret kyber.Scalar
	// Participants 是参与者的公钥列表。
	Participants []kyber.Point
	// Threshold 是密钥重建所需的最小份额数量。
	Threshold int
	// Host 是 P2P 网络主机。
	Host host.Host
	// NodeID 是当前节点的 ID。
	NodeID peer.ID
	// Shares 是当前节点持有的密钥份额。
	Shares map[peer.ID]*share.PriShare
}

// NewRabinDKG 创建一个新的 Rabin DKG 实例。
func NewRabinDKG(suite suites.Suite, participants []kyber.Point, threshold int) (*DKG, error) {
	// 检查参数。
	if len(participants) < threshold {
		return nil, errors.New("阈值必须小于参与者数量")
	}

	// 创建 DKG 对象。
	dkg := &DKG{
		Suite:        suite,
		Participants: participants,
		Threshold:    threshold,
		Shares:       make(map[peer.ID]*share.PriShare),
	}

	// 生成长期的私钥。
	dkg.NodeSecret = suite.Scalar().Pick(suite.RandomStream())

	return dkg, nil
}

// Run 启动 Rabin DKG 协议。
func (dkg *DKG) Run() error {
	// 初始化 VSS 协议。
	vss, err := rabin.NewDistKeyGenerator(dkg.Suite, dkg.NodeSecret, dkg.Participants, dkg.Threshold)
	if err != nil {
		return fmt.Errorf("初始化 VSS 协议失败: %w", err)
	}

	// 1. 生成密钥份额。
	deals, err := vss.Deals()
	if err != nil {
		return fmt.Errorf("生成密钥份额失败: %w", err)
	}

	// 2. 广播密钥份额。
	for participantID, deal := range deals {
		// 广播消息给指定参与者。
		if err := dkg.BroadcastMessage(peer.ID(fmt.Sprint(participantID)), deal); err != nil {
			return fmt.Errorf("广播密钥份额失败: %w", err)
		}
	}

	// 3. 接收并处理密钥份额。
	for participantID, deal := range deals {
		// 接收消息，处理密钥份额。
		if err := dkg.HandleDealMessage(peer.ID(fmt.Sprint(participantID)), deal); err != nil {
			return fmt.Errorf("处理密钥份额失败: %w", err)
		}
	}

	// 4. 验证密钥份额。
	if err := dkg.VerifyDeals(); err != nil {
		return fmt.Errorf("验证密钥份额失败: %w", err)
	}

	// 5. 生成秘密承诺。
	secretCommits, err := vss.SecretCommits()
	if err != nil {
		return fmt.Errorf("生成秘密承诺失败: %w", err)
	}

	// 6. 广播秘密承诺。
	for participantID, secretCommit := range secretCommits.Commitments {
		// 广播消息给指定参与者。
		if err := dkg.BroadcastMessage(peer.ID(fmt.Sprint(participantID)), secretCommit); err != nil {
			return fmt.Errorf("广播秘密承诺失败: %w", err)
		}
	}

	// 7. 接收并处理秘密承诺。
	// for participantID, secretCommit := range secretCommits {
	// 	// 接收消息，处理秘密承诺。
	// 	if err := dkg.HandleSecretCommitMessage(participantID, secretCommit); err != nil {
	// 		return fmt.Errorf("处理秘密承诺失败: %w", err)
	// 	}
	// }

	// 8. 验证秘密承诺。
	if err := dkg.VerifySecretCommits(); err != nil {
		return fmt.Errorf("验证秘密承诺失败: %w", err)
	}

	// 9. 密钥重建。
	key, err := dkg.ReconstructKey()
	if err != nil {
		return fmt.Errorf("密钥重建失败: %w", err)
	}

	// 密钥已生成。
	fmt.Println("密钥已生成:", key)

	return nil
}

// BroadcastMessage 广播消息给指定参与者。
func (dkg *DKG) BroadcastMessage(participantID peer.ID, message interface{}) error {
	// 获取参与者的连接。
	conn, err := dkg.Host.Network().NewStream(context.Background(), participantID)
	if err != nil {
		return fmt.Errorf("连接到参与者失败: %w", err)
	}
	defer conn.Close()

	// 将消息编码并发送。
	// ...
	return nil
}

// HandleDealMessage 处理密钥份额消息。
func (dkg *DKG) HandleDealMessage(participantID peer.ID, deal *rabin.Deal) error {
	// 解码消息并验证消息的正确性。
	// ...

	// 处理密钥份额。
	// share, err := deal.VerifyAndDecrypt(dkg.NodeID, dkg.Suite)
	// if err != nil {
	// 	return fmt.Errorf("处理密钥份额失败: %w", err)
	// }
	// dkg.Shares[participantID] = share

	return nil
}

// VerifyDeals 验证所有参与者发来的密钥份额。
func (dkg *DKG) VerifyDeals() error {
	// ...
	return nil
}

// HandleSecretCommitMessage 处理秘密承诺消息。
func (dkg *DKG) HandleSecretCommitMessage(participantID peer.ID, secretCommit *rabin.SecretCommits) error {
	// 解码消息并验证消息的正确性。
	// ...

	// 处理秘密承诺。
	// ...
	return nil
}

// VerifySecretCommits 验证所有参与者发来的秘密承诺。
func (dkg *DKG) VerifySecretCommits() error {
	// ...
	return nil
}

// ReconstructKey 从所有参与者的份额中重建密钥。
func (dkg *DKG) ReconstructKey() (kyber.Scalar, error) {
	// ...
	return nil, nil
}

// 剩余代码省略，包含：
// 1. HandleComplaintCommitMessage 处理投诉消息。
// 2. ReconstructShareFromComplaint 处理投诉，重建恶意节点的密钥份额。
// 3. 各种消息编码和解码函数。

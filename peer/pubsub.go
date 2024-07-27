// / Copyright (c) 2022 Sourcenetwork Developers. All rights reserved.
// / copy from https://github.com/sourcenetwork/orbis-go

package peer

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

var _ pubsub.EventTracer = (*pubsubTracer)(nil)

type pubsubTracer struct{}

func (p *pubsubTracer) Trace(evt *pb.TraceEvent) {
	switch evt.Type.String() {
	case pb.TraceEvent_DELIVER_MESSAGE.String():
		pid := peer.ID(string(evt.DeliverMessage.ReceivedFrom))
		fmt.Println("pubsub.tracer: event type ", evt.Type, " from ", pid, " on topic ", *(evt.DeliverMessage.Topic))
	case pb.TraceEvent_PUBLISH_MESSAGE.String():
		fmt.Println("pubsub.tracer: event type ", evt.Type, " on topic", *(evt.PublishMessage.Topic))
	}
}

func (p *Peer) Pub(ctx context.Context, topic string, data []byte) error {
	t, err := p.join(topic)
	if err != nil {
		return fmt.Errorf("join topic: %w", err)
	}

	err = t.Publish(ctx, data)
	if err != nil {
		return fmt.Errorf("publish topic: %w", err)
	}
	return nil
}

func (p *Peer) Sub(ctx context.Context, topic string) (*pubsub.Subscription, error) {
	t, err := p.join(topic)
	if err != nil {
		return nil, fmt.Errorf("join topic: %w", err)
	}

	sub, err := t.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("subscribe topic: %w", err)
	}
	return sub, nil
}

func (p *Peer) join(topic string) (*pubsub.Topic, error) {
	p.topicsLock.Lock()
	defer p.topicsLock.Unlock()

	t, exists := p.topics[topic]
	if exists {
		return t, nil
	}

	t, err := p.pubsub.Join(topic)
	if err != nil {
		return nil, fmt.Errorf("join topic: %w", err)
	}

	p.topics[topic] = t
	return t, nil
}

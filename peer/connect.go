// / Copyright (c) 2022 Sourcenetwork Developers. All rights reserved.
// / copy from https://github.com/sourcenetwork/orbis-go

package peer

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

const (
	reconnectAttempts                  = 20
	reconnectBackOffAttempts           = 10
	reconnectInterval                  = 5 * time.Second
	reconnectBackOffBaseSeconds        = 3
	dialRandomizerIntervalMilliseconds = 3000
)

func (p *Peer) Start(ctx context.Context) {
	for _, peer := range p.bootPeers {
		if err := p.Connect(ctx, peer); err != nil {
			fmt.Println("Can't connect to peer:", peer, err)
		} else {
			fmt.Println("Connected to bootstrap node:", peer)
		}
	}

	go func() {
		subCh, err := p.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged))
		if err != nil {
			fmt.Printf("Error subscribing to peer connectedness changes: %s \n", err)
		}
		defer subCh.Close()
		for {
			select {
			case ev, ok := <-subCh.Out():
				fmt.Println(ev)
				if !ok {
					return
				}

				evt := ev.(event.EvtPeerConnectednessChanged)
				if evt.Connectedness != network.NotConnected {
					continue
				}

				if _, ok := p.bootPeers[evt.Peer]; !ok {
					continue
				}

				paddr := p.bootPeers[evt.Peer]
				go p.reconnectToPeer(ctx, paddr)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (p *Peer) reconnectToPeer(ctx context.Context, paddr peer.AddrInfo) {
	if _, ok := p.reonnecting.Load(paddr.ID.String()); ok {
		fmt.Println("duplicate peer maintainence goroutine", paddr.ID)
		return
	}

	p.reonnecting.Store(paddr.ID.String(), struct{}{})
	defer p.reonnecting.Delete(paddr.ID.String())

	start := time.Now()
	fmt.Printf("Reconnecting to peer %s \n", paddr)
	for i := 0; i < reconnectAttempts; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("peer maintainence goroutine context finished", paddr.ID)
			return
		default:
			// noop fallthrough
		}

		err := p.Connect(ctx, paddr)
		if err == nil {
			fmt.Printf("reconnected to peer %s during regular backoff \n", paddr.ID)
			return //success
		}

		fmt.Printf("Error reconnecting to peer %s: %s, Retrying %d/%d attemps \n", paddr, err, i, reconnectAttempts)
		randomSleep(reconnectInterval)
	}

	fmt.Printf("Failed to reconnect to peer %s. Beginning exponential backoff. Elapsed %s \n", paddr, time.Since(start))
	for i := 0; i < reconnectBackOffAttempts; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// noop fallthrough
		}

		// sleep an exponentially increasing amount
		sleepIntervalSeconds := math.Pow(reconnectBackOffBaseSeconds, float64(i))
		randomSleep(time.Duration(sleepIntervalSeconds) * time.Second)

		err := p.Connect(ctx, paddr)
		if err == nil {
			fmt.Printf("reconnected to peer %s during exponential backoff \n", paddr.ID)
			return //success
		}

		fmt.Printf("Error reconnecting to peer %s: %s, Retrying %d/%d attemps\n", paddr, err, i, reconnectAttempts)
	}
	fmt.Printf("Failed to reconnect to peer %s. Giving up. Elapsed %s\n", paddr, time.Since(start))
}

func (p *Peer) Discover(ctx context.Context) error {
	rendezvous := "wetee"
	d := drouting.NewRoutingDiscovery(p.idht)
	dutil.Advertise(ctx, d, rendezvous)

	fmt.Println("Peer discovery start...")
	peerChan, err := d.FindPeers(ctx, rendezvous)
	if err != nil {
		return fmt.Errorf("Find peers error: %w", err)
	}

	defer fmt.Println("Peer discovery finished...")
	for peer := range peerChan {
		if peer.ID == p.ID() {
			continue
		}

		if len(peer.Addrs) == 0 {
			continue
		}

		err = p.Connect(ctx, peer)
		if err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}

	return nil
}

func randomSleep(interval time.Duration) {
	r := time.Duration(rand.Int63n(dialRandomizerIntervalMilliseconds)) * time.Millisecond
	time.Sleep(r + interval)
}

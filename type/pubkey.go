package types

import (
	gocrypto "crypto"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2pCryptoPb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/suites"
)

type PubKey struct {
	libp2pCrypto.PubKey
	suite suites.Suite
}

func (p *PubKey) Point() kyber.Point {
	buf, _ := p.PubKey.Raw()
	point := p.suite.Point()
	point.UnmarshalBinary(buf)
	return point
}

func (p *PubKey) Std() (gocrypto.PublicKey, error) {
	return libp2pCrypto.PubKeyToStdKey(p.PubKey)
}

func (p *PubKey) String() string {
	bt, err := libp2pCrypto.MarshalPublicKey(p.PubKey)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bt)
}

func (p *PubKey) PeerID() peer.ID {
	peerID, err := peer.IDFromPublicKey(p)
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		return peer.ID("")
	}
	return peerID
}

func (p *PubKey) Byte() ([]byte, error) {
	bt, err := libp2pCrypto.MarshalPublicKey(p.PubKey)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func PublicKeyFromLibp2pHex(hexStr string) (*PubKey, error) {
	buf, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("decode hex: %w", err)
	}

	return PublicKeyFromLibp2pBytes(buf)
}

func PublicKeyFromHex(hexStr string) (*PubKey, error) {
	return PublicKeyFromLibp2pHex("08011220" + hexStr)
}

func PublicKeyFromLibp2pBytes(buf []byte) (*PubKey, error) {
	pk, err := libp2pCrypto.UnmarshalPublicKey(buf)
	if err != nil {
		return nil, fmt.Errorf("public key from bytes: %w", err)
	}

	return PubKeyFromLibP2P(pk)
}

func PubKeyFromLibP2P(pubkey libp2pCrypto.PubKey) (*PubKey, error) {
	suite, err := SuiteForType(pubkey.Type())
	if err != nil {
		return nil, err
	}

	return &PubKey{
		PubKey: pubkey,
		suite:  suite,
	}, nil
}

func PubKeyFromStdPubKey(pubkey gocrypto.PublicKey) (*PubKey, error) {
	var libpk libp2pCrypto.PubKey
	var err error
	switch gopk := pubkey.(type) {
	case ed25519.PublicKey:
		libpk, err = libp2pCrypto.UnmarshalEd25519PublicKey(gopk)
	default:
		return nil, fmt.Errorf("unknown key type")
	}

	if err != nil {
		return nil, err
	}

	return PubKeyFromLibP2P(libpk)
}

func PubKeyFromPoint(suite suites.Suite, point kyber.Point) (*PubKey, error) {

	buf, err := point.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal point: %w", err)
	}

	var pk libp2pCrypto.PubKey

	switch strings.ToLower(suite.String()) {
	case "ed25519":
		pk, err = libp2pCrypto.UnmarshalEd25519PublicKey(buf)
	default:
		return nil, fmt.Errorf("unknown suite: %v", suite)
	}

	if err != nil {
		return nil, fmt.Errorf("unmarshal public key: %w", err)
	}

	return PubKeyFromLibP2P(pk)
}

func SuiteForType(kt libp2pCryptoPb.KeyType) (suites.Suite, error) {
	switch kt {
	case libp2pCryptoPb.KeyType_Ed25519:
		return edwards25519.NewBlakeSHA256Ed25519(), nil
	default:
		return nil, fmt.Errorf("unsupported key type: %v", kt)
	}
}

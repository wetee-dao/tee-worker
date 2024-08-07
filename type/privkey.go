package types

import (
	gocrypto "crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2pCryptoPb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

type PrivKey struct {
	libp2pCrypto.PrivKey
	suite suites.Suite
}

func (p *PrivKey) Scalar() kyber.Scalar {
	switch p.Type() {
	case libp2pCryptoPb.KeyType_Ed25519:
		return p.ed25519Scalar()
	default:
		panic("only ed25519 and secp256k1 private key scalar conversion supported")
	}
}

func (p *PrivKey) Suite() suites.Suite {
	return p.suite
}

func (p *PrivKey) ed25519Scalar() kyber.Scalar {
	// There is a discrepency between LibP2P private keys
	// and "raw" EC scalars. LibP2P private keys is an
	// (x, y) pair, where x is the given "seed" and y is
	// the cooresponding publickey. Where y is computed as
	//
	// h := sha512.Hash(x)
	// s := scalar().SetWithClamp(h)
	// y := point().ScalarBaseMul(x)
	//
	// So to make sure future conversions of this scalar
	// to a public key, like in the DKG setup, we need to
	// convert this key to a scalar using the Hash and Clamp
	// method.
	//
	// To understand clamping, see here:
	// https://neilmadden.blog/2020/05/28/whats-the-curve25519-clamping-all-about/

	buf, err := p.PrivKey.Raw()
	if err != nil {
		panic(err)
	}

	// hash seed and clamp bytes
	digest := sha512.Sum512(buf[:32])
	digest[0] &= 0xf8
	digest[31] &= 0x7f
	digest[31] |= 0x40
	return p.suite.Scalar().SetBytes(digest[:32])
}

func (p *PrivKey) String() string {
	bt, err := libp2pCrypto.MarshalPrivateKey(p.PrivKey)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bt)
}

func (p *PrivKey) GetPublic() *PubKey {
	return &PubKey{
		PubKey: p.PrivKey.GetPublic(),
		suite:  p.suite,
	}
}

func GenerateKeyPair(ste suites.Suite, src io.Reader) (*PrivKey, *PubKey, error) {
	keyType, err := KeyTypeFromString(ste.String())
	if err != nil {
		return nil, nil, err
	}
	sk, pk, err := libp2pCrypto.GenerateKeyPairWithReader(keyType, 0, src)
	if err != nil {
		return nil, nil, err
	}

	return &PrivKey{
			PrivKey: sk,
			suite:   ste,
		}, &PubKey{
			PubKey: pk,
			suite:  ste,
		}, nil
}

// PrivateKeyFromLibp2pHex
func PrivateKeyFromLibp2pHex(key string) (*PrivKey, error) {
	buf, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("private key from hex: %w", err)
	}

	return PrivateKeyFromLibp2pBytes(buf)
}

func PrivateKeyFromHex(key string) (*PrivKey, error) {
	return PrivateKeyFromLibp2pHex("08011240" + key)
}

func PrivateKeyFromLibp2pBytes(buf []byte) (*PrivKey, error) {
	pk, err := libp2pCrypto.UnmarshalPrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("public key from bytes: %w", err)
	}

	return PrivateKeyFromLibP2P(pk)
}

func PrivateKeyFromLibP2P(privkey libp2pCrypto.PrivKey) (*PrivKey, error) {
	suite, err := SuiteForType(privkey.Type())
	if err != nil {
		return nil, err
	}

	return &PrivKey{
		PrivKey: privkey,
		suite:   suite,
	}, nil
}

func PrivateKeyFromStd(privkey gocrypto.PrivateKey) (*PrivKey, error) {
	var libpk libp2pCrypto.PrivKey
	var err error
	switch gopk := privkey.(type) {
	case ed25519.PrivateKey:
		libpk, err = libp2pCrypto.UnmarshalEd25519PrivateKey(gopk)
	default:
		return nil, fmt.Errorf("unknown key type")
	}

	if err != nil {
		return nil, err
	}

	return PrivateKeyFromLibP2P(libpk)
}

func KeyTypeFromString(keyType string) (int, error) {
	switch strings.ToLower(keyType) {
	case "ed25519":
		return libp2pCrypto.Ed25519, nil
	default:
		return 0, fmt.Errorf("unknown key type: %s", keyType)
	}
}

package p2p

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// ID is a hex-encoded crypto.Address
type ID string

// IDByteLength is the length of a crypto.Address. Currently only 20.
// TODO: support other length addresses ?
const IDByteLength = crypto.AddressSize

//------------------------------------------------------------------------------
// Persistent peer ID
// TODO: encrypt on disk

// NodeKey is the persistent peer key.
// It contains the nodes private key for authentication.
type NodeKey struct {
	PrivKey    crypto.PrivKey               `json:"priv_key"` // our priv key
	RSAPrivKey string                       `json:"rsa_priv_key"`
	RSAPubkKey string                       `json:"rsa_pub_key"`
	GroupKeys  map[string]map[string]string `json:"group_keys"`
}

// ID returns the peer's canonical ID - the hash of its public key.
func (nodeKey *NodeKey) ID() ID {
	return PubKeyToID(nodeKey.PubKey())
}

// GetRSAPrivKey retrive rsa privkey
func (nodeKey *NodeKey) GetRSAPrivKey() (*rsa.PrivateKey, error) {
	rk, err := hex.DecodeString(nodeKey.RSAPrivKey)
	if err != nil {
		return nil, errors.New("Decode RSAPrivKey failed:" + err.Error())
	}
	privkey, err := x509.ParsePKCS1PrivateKey(rk)
	if err != nil {
		return nil, errors.New("Parse RSAPrivKey failed:" + err.Error())
	}
	return privkey, nil
}

// PubKey returns the peer's PubKey
func (nodeKey *NodeKey) PubKey() crypto.PubKey {
	return nodeKey.PrivKey.PubKey()
}

// Save rewrite
func (nodeKey *NodeKey) Save(filePath string) error {
	jsonBytes, err := cdc.MarshalJSONIndent(nodeKey, "", "	")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

// PubKeyToID returns the ID corresponding to the given PubKey.
// It's the hex-encoding of the pubKey.Address().
func PubKeyToID(pubKey crypto.PubKey) ID {
	return ID(hex.EncodeToString(pubKey.Address()))
}

// LoadOrGenNodeKey attempts to load the NodeKey from the given filePath.
// If the file does not exist, it generates and saves a new NodeKey.
func LoadOrGenNodeKey(filePath string) (*NodeKey, error) {
	if cmn.FileExists(filePath) {
		nodeKey, err := LoadNodeKey(filePath)
		if err != nil {
			return nil, err
		}
		return nodeKey, nil
	}
	return genNodeKey(filePath)
}

func LoadNodeKey(filePath string) (*NodeKey, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := new(NodeKey)
	err = cdc.UnmarshalJSON(jsonBytes, nodeKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading NodeKey from %v: %v", filePath, err)
	}
	return nodeKey, nil
}

func genNodeKey(filePath string) (*NodeKey, error) {
	privKey := ed25519.GenPrivKey()

	// gen rsk key pair
	rsaSK, err := rsa.GenerateKey(crypto.CReader(), 1024)
	if err != nil {
		return nil, err
	}
	rsaSKbs := x509.MarshalPKCS1PrivateKey(rsaSK)
	rsaPKbs := x509.MarshalPKCS1PublicKey(&rsaSK.PublicKey)

	nodeKey := &NodeKey{
		PrivKey:    privKey,
		RSAPrivKey: hex.EncodeToString(rsaSKbs),
		RSAPubkKey: hex.EncodeToString(rsaPKbs),
	}

	jsonBytes, err := cdc.MarshalJSONIndent(nodeKey, "", "	")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}

//------------------------------------------------------------------------------

// MakePoWTarget returns the big-endian encoding of 2^(targetBits - difficulty) - 1.
// It can be used as a Proof of Work target.
// NOTE: targetBits must be a multiple of 8 and difficulty must be less than targetBits.
func MakePoWTarget(difficulty, targetBits uint) []byte {
	if targetBits%8 != 0 {
		panic(fmt.Sprintf("targetBits (%d) not a multiple of 8", targetBits))
	}
	if difficulty >= targetBits {
		panic(fmt.Sprintf("difficulty (%d) >= targetBits (%d)", difficulty, targetBits))
	}
	targetBytes := targetBits / 8
	zeroPrefixLen := (int(difficulty) / 8)
	prefix := bytes.Repeat([]byte{0}, zeroPrefixLen)
	mod := (difficulty % 8)
	if mod > 0 {
		nonZeroPrefix := byte(1<<(8-mod) - 1)
		prefix = append(prefix, nonZeroPrefix)
	}
	tailLen := int(targetBytes) - len(prefix)
	return append(prefix, bytes.Repeat([]byte{0xFF}, tailLen)...)
}

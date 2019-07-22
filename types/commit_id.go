package types

import (
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// CommitID contains the tree version number and its merkle root.
type CommitID struct {
	Version      int64       `json:"version"`
	Hash         []byte      `json:"hash"`
	ShardingHash cmn.KVPairs `json:"sharding_hash"`
}

func UnmarshalCommitID(bz []byte) CommitID {
	cdc := amino.NewCodec()

	var ch CommitID
	err := cdc.UnmarshalJSON(bz, &ch)
	if err != nil {
		panic(err.Error())
	}

	return ch
}

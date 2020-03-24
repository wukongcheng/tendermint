package types

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/kv"
)

// CommitID contains the tree version number and its merkle root.
type CommitID struct {
	Version      int64       `json:"version"`
	Hash         []byte      `json:"hash"`
	ShardingHash kv.Pairs `json:"sharding_hash"`
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

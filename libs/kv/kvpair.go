package kv

import (
	"bytes"
	"encoding/hex"
	"sort"
	"strings"
)

//----------------------------------------
// KVPair

/*
Defined in types.proto

type Pair struct {
	Key   []byte
	Value []byte
}
*/

type Pairs []Pair

// Sorting
func (kvs Pairs) Len() int { return len(kvs) }
func (kvs Pairs) Less(i, j int) bool {
	switch bytes.Compare(kvs[i].Key, kvs[j].Key) {
	case -1:
		return true
	case 0:
		return bytes.Compare(kvs[i].Value, kvs[j].Value) < 0
	case 1:
		return false
	default:
		panic("invalid comparison result")
	}
}
func (kvs Pairs) Swap(i, j int) { kvs[i], kvs[j] = kvs[j], kvs[i] }
func (kvs Pairs) Sort()         { sort.Sort(kvs) }

func (kvs Pairs) ToString() (str string) {
	kvs.Sort()
	for _, pair := range kvs {
		str += string(pair.Key)
		str += ":"
		str += hex.EncodeToString(pair.Value)
		str += "|"
	}
	return
}

func KVPairsFromString(str string) (kvs Pairs) {
	if len(str) == 0 {
		return
	}

	strs := strings.Split(str, "|")
	for _, s := range strs {
		if len(s) == 0 {
			continue
		}

		kv := strings.Split(s, ":")
		hash, err := hex.DecodeString(kv[1])
		if err != nil {
			panic("invalid hex bytes")
		}

		kvp := Pair{
			Key: []byte(kv[0]),
			Value: hash,
		}
		kvs = append(kvs, kvp)
	}
	return
}

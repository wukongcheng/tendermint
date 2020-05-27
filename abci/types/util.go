package types

import (
	"bytes"
	"sort"
	"github.com/tendermint/tendermint/libs/kv"
)

const (
	DefaultEvent = "origin-tags"
)

//------------------------------------------------------------------------------

// ValidatorUpdates is a list of validators that implements the Sort interface
type ValidatorUpdates []ValidatorUpdate

var _ sort.Interface = (ValidatorUpdates)(nil)

// All these methods for ValidatorUpdates:
//    Len, Less and Swap
// are for ValidatorUpdates to implement sort.Interface
// which will be used by the sort package.
// See Issue https://github.com/tendermint/abci/issues/212

func (v ValidatorUpdates) Len() int {
	return len(v)
}

// XXX: doesn't distinguish same validator with different power
func (v ValidatorUpdates) Less(i, j int) bool {
	return bytes.Compare(v[i].PubKey.Data, v[j].PubKey.Data) <= 0
}

func (v ValidatorUpdates) Swap(i, j int) {
	v1 := v[i]
	v[i] = v[j]
	v[j] = v1
}

func GetTagByKey(events []Event, key string) (kv.Pair, bool) {
	for _, event := range events {
		if event.GetType() != DefaultEvent {
			continue
		}
		for _, tag := range event.Attributes {
			if bytes.Equal(tag.Key, []byte(key)) {
				return tag, true
			}
		}
	}

	return kv.Pair{}, false
}

func TagsToDefaultEvent(events []Event, tags ...kv.Pair) []Event {
	if len(events) == 0 {
		events = append(events, Event{Type: DefaultEvent})
	}
	for i, v := range events {
		if v.Type == DefaultEvent {
			events[i].Attributes = append(events[i].Attributes, tags...)
		}
	}
	return events
}

func GetDefaultTags(events []Event) []kv.Pair {
	for _, v := range events {
		if v.Type == DefaultEvent {
			pairs := make([]kv.Pair, len(v.Attributes))
			copy(pairs, v.Attributes)
			return pairs
		}
	}
	return nil
}

func GetAllTags(events []Event) (pairs []kv.Pair) {
	for _, v := range events {
		ps := make([]kv.Pair, len(v.Attributes))
		copy(ps, v.Attributes)
		pairs = append(pairs, ps...)
	}
	return pairs
}

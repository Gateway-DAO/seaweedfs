package types

import (
	"encoding/json"

	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"github.com/seaweedfs/seaweedfs/weed/storage/types"
)

/**
 * MARK: Generic Events
 */
type Timestamp uint64

type NeedleFunction uint32

const (
	WRITE NeedleFunction = iota
	DELETED
)

type NeedleEvent struct {
	id   types.NeedleId
	hash needle.CRC
	fn   NeedleFunction

	created_at   Timestamp
	last_updated Timestamp
	last_touched Timestamp
}

func (e *NeedleEvent) Id() types.NeedleId {
	return e.id
}

func (e *NeedleEvent) Key() []byte {
	key := []byte{}
	types.NeedleIdToBytes(key, e.id)
	return key
}

func (e *NeedleEvent) Value() ([]byte, error) {
	obj := struct {
		Hash        needle.CRC     `json:"hash"`
		Function    NeedleFunction `json:"fn"`
		CreatedAt   Timestamp      `json:"created_at"`
		LastUpdated Timestamp      `json:"last_updated"`
		LastTouched Timestamp      `json:"last_touched"`
	}{
		Hash:        e.hash,
		Function:    e.fn,
		CreatedAt:   e.created_at,
		LastUpdated: e.last_updated,
		LastTouched: e.last_touched,
	}

	return json.Marshal(obj)
}

func (e *NeedleEvent) Timestamps() (created_at, last_updated, last_touched Timestamp) {
	return e.created_at, e.last_updated, e.last_touched
}

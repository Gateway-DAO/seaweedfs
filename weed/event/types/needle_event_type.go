package types

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"github.com/seaweedfs/seaweedfs/weed/storage/types"
)

/**
 * MARK: Generic Events
 */
type Timestamp uint64

type NeedleFunctionType uint32

const (
	WRITE NeedleFunctionType = iota
	DELETED
)

type NeedleEvent struct {
	id       types.NeedleId
	checksum uint32
	hash     string
	fn       NeedleFunctionType

	created_at   Timestamp
	last_updated Timestamp
	last_touched Timestamp
}

func NewNeedleEvent(id types.NeedleId, checksum needle.CRC, hash string, fn NeedleFunctionType) (n *NeedleEvent) {
	n = new(NeedleEvent)

	n.id = id
	n.checksum = uint32(checksum)
	n.hash = hash
	n.fn = fn

	created_at := Timestamp(time.Now().Unix())
	n.created_at = created_at
	n.last_updated = created_at
	n.last_touched = created_at

	return
}

func (e *NeedleEvent) Id() types.NeedleId {
	return e.id
}

func (e *NeedleEvent) Key() string {
	return e.id.String()
}

func (e *NeedleEvent) Value() ([]byte, error) {
	obj := struct {
		Checksum    string             `json:"checksum"`
		Hash        string             `json:"hash"`
		Function    NeedleFunctionType `json:"fn"`
		CreatedAt   Timestamp          `json:"created_at"`
		LastUpdated Timestamp          `json:"last_updated"`
		LastTouched Timestamp          `json:"last_touched"`
	}{
		Checksum:    strconv.FormatUint(uint64(e.checksum), 16),
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

package event_types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
)

/**
 * MARK: Generic Events
 */
type NanoTimestamp int64

type NeedleEventType uint32

const (
	WRITE NeedleEventType = iota
	DELETE
)

// type NeedleEvent struct {
// 	id string
// 	fn NeedleEventType

// 	dataCenter string
// 	rack       string

// 	needleId types.NeedleId
// 	checksum needle.CRC
// 	hash     *string

// 	created_at   Timestamp
// 	last_updated Timestamp
// 	last_touched Timestamp
// }

type VolumeNeedleEvent struct {
	Type string `json:"type"`
	Hash string `json:"hash,omitempty"`

	Needle struct {
		Id       string          `json:"id"`
		Checksum uint32          `json:"checksum"`
		VolumeId needle.VolumeId `json:"volume_id"`
	} `json:"needle"`

	VolumeServer struct {
		Url        string `json:"url"`
		Rack       string `json:"rack"`
		DataCenter string `json:"data_center"`
	} `json:"volume_server"`

	CreatedAt   NanoTimestamp `json:"created_at"`
	LastUpdated NanoTimestamp `json:"last_updated"`
	LastTouched NanoTimestamp `json:"last_touched"`
}

func NewNeedleEvent(volumeId needle.VolumeId, vsUrl, vsDataCenter, vsRack string, n *needle.Needle, hash *string, eventType NeedleEventType) (ne *VolumeNeedleEvent, err error) {
	ne = new(VolumeNeedleEvent)

	switch eventType {
	case WRITE:
		ne.Type = "WRITE"
	case DELETE:
		ne.Type = "DELETE"
	}

	if hash != nil {
		ne.Hash = *hash
	}

	if n == nil {
		return nil, fmt.Errorf("needle argument is null")
	}
	ne.Needle = struct {
		Id       string          `json:"id"`
		Checksum uint32          `json:"checksum"`
		VolumeId needle.VolumeId `json:"volume_id"`
	}{
		Id:       n.Id.String(),
		Checksum: n.Checksum.Value(),
		VolumeId: volumeId,
	}

	ne.VolumeServer = struct {
		Url        string `json:"url"`
		Rack       string `json:"rack"`
		DataCenter string `json:"data_center"`
	}{
		Url:        vsUrl,
		Rack:       vsRack,
		DataCenter: vsDataCenter,
	}

	now := NanoTimestamp(time.Now().UnixNano())
	ne.CreatedAt = now
	ne.LastUpdated = now
	ne.LastTouched = now

	return
}
func (e *VolumeNeedleEvent) Value() ([]byte, error) {
	return json.Marshal(e)
}

func (e *VolumeNeedleEvent) Timestamps() (created_at, last_updated, last_touched NanoTimestamp) {
	return e.CreatedAt, e.LastUpdated, e.LastTouched
}

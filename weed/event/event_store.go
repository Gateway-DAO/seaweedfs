package event

import (
	"encoding/binary"
	"sync"
)

type NanoTimestamp int64

type VolumeServerEventType uint32

const (
	ALIVE VolumeServerEventType = iota
	WRITE
	DELETE
	VACUUM
)

var needleEventTypes = map[VolumeServerEventType]string{
	ALIVE:  "ALIVE",
	WRITE:  "WRITE",
	DELETE: "DELETE",
	VACUUM: "VACUUM",
}

type VolumeServerEventStore interface {
	RegisterEvent(*VolumeServerEvent) error
	GetLastEvent() (*VolumeServerEvent, error)
	ListAllEvents() ([]*VolumeServerEvent, error)
}

type EventStoreImpl struct {
	sync.RWMutex
}

func timestampToBytes(ts int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(ts))
	return b
}

func bytesToTimestamp(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

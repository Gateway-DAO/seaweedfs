package event

import (
	"encoding/binary"
	"sync"
)

type NanoTimestamp int64

type NeedleEvent uint32

const (
	GENESIS NeedleEvent = iota
	ALIVE
	WRITE
	DELETE
	VACUUM
)

var needleEventTypes = map[NeedleEvent]string{
	GENESIS: "GENESIS",
	ALIVE:   "ALIVE",
	WRITE:   "WRITE",
	DELETE:  "DELETE",
	VACUUM:  "VACUUM",
}

type EventStore interface {
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

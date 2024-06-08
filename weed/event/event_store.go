package event

import "encoding/binary"

type NanoTimestamp int64

type NeedleEventType uint32

const (
	WRITE NeedleEventType = iota
	DELETE
	VACUUM
)

type EventStore interface {
	RegisterEvent(*VolumeServerEvent) error
	ListAllEvents() ([]*VolumeServerEvent, error)
}

func timestampToBytes(ts int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(ts))
	return b
}

func bytesToTimestamp(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

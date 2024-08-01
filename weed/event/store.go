package event

import (
	"encoding/binary"
)

type EventStore[T Event] interface {
	RegisterEvent(T) error
	GetLastEvent() (T, error)
	ListAllEvents() ([]T, error)
}

func timestampToBytes(ts int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(ts))
	return b
}

func bytesToTimestamp(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

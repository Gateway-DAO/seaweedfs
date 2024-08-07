package event

import (
	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
)

type Event interface {
	SetType(t string)
	GetType() string
	isAliveType() bool

	GetServer() *event_pb.Server
	GetProofOfHistory() *event_pb.ProofOfHistory
	SetProofOfHistory(previousHash *string, hash string)

	GetKey() ([]byte, error)
	GetValue() ([]byte, error)
}

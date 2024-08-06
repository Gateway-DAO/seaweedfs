package event

import (
	"encoding/json"
	"time"

	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
	"github.com/gateway-dao/seaweedfs/weed/pb/master_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MasterServerEventType uint32

const (
	MASTER_ALIVE MasterServerEventType = iota
	ASSIGN
)

var MasterServerEvents = map[MasterServerEventType]string{
	MASTER_ALIVE: "MASTER_ALIVE",
	ASSIGN:       "ASSIGN",
}

type MasterServerEvent struct {
	// TODO: consider replacing with protobuf value
	Type      string                 `json:"type"`
	Timestamp *timestamppb.Timestamp `json:"timestamp"`

	Fid            string                   `json:"fid"`
	Locations      []*master_pb.Location    `json:"locations"`
	Server         *event_pb.Server         `json:"server"`
	ProofOfHistory *event_pb.ProofOfHistory `json:"proofOfHistory"`
}

type MasterServerEventKey struct {
	Server string `json:"server"`
	Type   string `json:"type"`
}

func NewMasterServerEvent(
	eventType MasterServerEventType,
	fid *string,
	locations []*master_pb.Location,
	serverPublicUrl string,
) *MasterServerEvent {
	mse := new(MasterServerEvent)

	mse.Type = MasterServerEvents[eventType]
	mse.Server = &event_pb.Server{
		PublicUrl: serverPublicUrl,
	}

	if fid != nil {
		mse.Fid = *fid
	}
	if locations != nil {
		mse.Locations = locations
	}

	mse.Timestamp = timestamppb.New(time.Now())

	return mse
}

func (mse *MasterServerEvent) isAliveType() bool {
	return mse.GetType() == MasterServerEvents[MASTER_ALIVE]
}

func (mse *MasterServerEvent) SetType(t string) {
	mse.Type = t
}

func (mse *MasterServerEvent) GetType() string {
	return mse.Type
}

func (mse *MasterServerEvent) GetServer() *event_pb.Server {
	return mse.Server
}

func (mse *MasterServerEvent) GetProofOfHistory() *event_pb.ProofOfHistory {
	return mse.ProofOfHistory
}

func (mse *MasterServerEvent) SetProofOfHistory(previousHash *string, hash string) {
	mse.ProofOfHistory = &event_pb.ProofOfHistory{
		PreviousHash: previousHash,
		Hash:         hash,
	}
}

func (mse *MasterServerEvent) GetKafkaKey() ([]byte, error) {
	return json.Marshal(MasterServerEventKey{
		Type:   mse.Type,
		Server: mse.Server.PublicUrl,
	})
}

func (mse *MasterServerEvent) GetValue() ([]byte, error) {
	return json.Marshal(mse)
}

package event

import (
	"encoding/json"
	"time"

	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VolumeServerEventType uint32

const (
	ALIVE VolumeServerEventType = iota
	WRITE
	DELETE
	VACUUM
)

var vsEventTypes = map[VolumeServerEventType]string{
	ALIVE:  "ALIVE",
	WRITE:  "WRITE",
	DELETE: "DELETE",
	VACUUM: "VACUUM",
}

type VolumeServerEvent struct {
	*volume_server_pb.VolumeServerEventResponse
}

type VolumeServerEventKafkaKey struct {
	Volume string `json:"volume"`
	Server string `json:"server"`
}

func NewVolumeServerEvent(
	eventType VolumeServerEventType,
	serverMetadata *event_pb.Server,
	volumeMetadata *volume_server_pb.VolumeServerEventResponse_Volume,
	needleMetadata *volume_server_pb.VolumeServerEventResponse_Needle,
) (*VolumeServerEvent, error) {
	ne := new(VolumeServerEvent)
	ne.VolumeServerEventResponse = &volume_server_pb.VolumeServerEventResponse{
		Type:   vsEventTypes[eventType],
		Volume: volumeMetadata,
		Server: serverMetadata,
	}

	if needleMetadata != nil {
		ne.VolumeServerEventResponse.Needle = needleMetadata
	}

	ne.Timestamp = timestamppb.New(time.Now())

	return ne, nil
}

func (vse *VolumeServerEvent) SetType(t string) {
	vse.Type = t
}

func (vse *VolumeServerEvent) GetType() string {
	return vse.Type
}

func (vse *VolumeServerEvent) isAliveType() bool {
	return vse.GetType() == vsEventTypes[ALIVE]
}

func (vse *VolumeServerEvent) GetServer() *event_pb.Server {
	return vse.Server
}

func (vse *VolumeServerEvent) SetProofOfHistory(previousHash *string, hash string) {
	vse.ProofOfHistory = &event_pb.ProofOfHistory{
		PreviousHash: previousHash,
		Hash:         hash,
	}
}

func (vse *VolumeServerEvent) GetProofOfHistory() *event_pb.ProofOfHistory {
	return vse.ProofOfHistory
}

func (vse *VolumeServerEvent) GetKafkaKey() ([]byte, error) {
	kafkaEventKey := VolumeServerEventKafkaKey{
		Server: vse.Server.PublicUrl,
	}
	if vse.Volume != nil {
		kafkaEventKey.Volume = vse.Volume.Id
	}
	return json.Marshal(kafkaEventKey)
}

func (vse *VolumeServerEvent) GetValue() ([]byte, error) {
	return json.Marshal(volume_server_pb.VolumeServerEventResponse{
		Type:           vse.Type,
		Timestamp:      vse.Timestamp,
		Needle:         vse.Needle,
		Volume:         vse.Volume,
		Server:         vse.Server,
		ProofOfHistory: vse.ProofOfHistory,
	})
}

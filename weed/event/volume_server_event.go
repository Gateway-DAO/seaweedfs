package event

import (
	"encoding/json"
	"time"

	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VolumeServerEvent struct {
	volume_server_pb.VolumeServerEventResponse
}

func NewVolumeServerEvent(
	eventType VolumeServerEventType,
	serverMetadata *volume_server_pb.VolumeServerEventResponse_Server,
	volumeMetadata *volume_server_pb.VolumeServerEventResponse_Volume,
	needleMetadata *volume_server_pb.VolumeServerEventResponse_Needle,
) (ne *VolumeServerEvent, err error) {
	ne = new(VolumeServerEvent)

	ne.Type = needleEventTypes[eventType]
	ne.Volume = volumeMetadata
	ne.Server = serverMetadata
	if needleMetadata != nil {
		ne.Needle = needleMetadata
	}

	ne.Timestamp = timestamppb.New(time.Now())

	return
}

func (e *VolumeServerEvent) Value() ([]byte, error) {
	return json.Marshal(e)
}

package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VolumeServerEvent struct {
	volume_server_pb.VolumeServerEventResponse
}

func NewVolumeServerEvent(
	eventType NeedleEventType,
	serverMetadata *volume_server_pb.VolumeServerEventResponse_Server,
	volumeMetadata *volume_server_pb.VolumeServerEventResponse_Volume,
	needleMetadata *volume_server_pb.VolumeServerEventResponse_Volume_Needle,
) (ne *VolumeServerEvent, err error) {
	ne = new(VolumeServerEvent)

	switch eventType {
	case WRITE:
		ne.Type = "WRITE"
	case DELETE:
		ne.Type = "DELETE"
	case VACUUM:
		ne.Type = "VACUUM"
	default:
		return nil, fmt.Errorf("unable to parse event type %d", eventType)
	}

	ne.Volume = volumeMetadata
	ne.Server = serverMetadata
	if needleMetadata != nil {
		ne.Volume.Needle = needleMetadata
	}

	ne.Timestamp = timestamppb.New(time.Now())

	return
}

func (e *VolumeServerEvent) Value() ([]byte, error) {
	return json.Marshal(e)
}

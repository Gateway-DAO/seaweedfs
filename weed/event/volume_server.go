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

func NewVolumeServerEvent(eventType NeedleEventType, needleMetadata *volume_server_pb.VolumeServerEventResponse_Needle, volumeMetadata *volume_server_pb.VolumeServerEventResponse_Volume) (ne *VolumeServerEvent, err error) {
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

	ne.Needle = needleMetadata

	ne.Volume = volumeMetadata

	ne.CreatedAt = timestamppb.New(time.Now())

	return
}

func (e *VolumeServerEvent) Key() []byte {
	return []byte(fmt.Sprintf("%s:%d:%s", e.GetVolume().GetId(), e.GetNeedle().GetId(), e.GetCreatedAt()))
}

func (e *VolumeServerEvent) Value() ([]byte, error) {
	return json.Marshal(e)
}

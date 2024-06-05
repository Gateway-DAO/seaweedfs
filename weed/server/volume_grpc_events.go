package weed_server

import (
	"time"

	"github.com/seaweedfs/seaweedfs/weed/event"
	"github.com/seaweedfs/seaweedfs/weed/event/event_types"
	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (vs *VolumeServer) VolumeServerEvents(req *volume_server_pb.VolumeServerEventsRequest, stream volume_server_pb.VolumeServer_VolumeServerEventsServer) error {
	eventDir := vs.eventsDir

	needleEvents, err := event.ListEvents(eventDir)
	if err != nil {
		return err
	}

	for _, event := range needleEvents {
		if ctxErr := stream.Context().Err(); ctxErr != nil {
			return ctxErr
		}

		parsedEvent := prepareVolumeServerEventResponse(event)
		if streamErr := stream.Send(parsedEvent); streamErr != nil {
			return streamErr
		}

		// DEBUG: Simulate delay
		time.Sleep(1 * time.Second)
	}

	return nil
}

func prepareVolumeServerEventResponse(event *event_types.NeedleEvent) *volume_server_pb.VolumeServerEventResponse {
	resp := &volume_server_pb.VolumeServerEventResponse{
		Type: event.Type,
		Hash: event.Hash,
		Needle: &volume_server_pb.VolumeServerEventResponse_Needle{
			Id:       event.Needle.Id,
			Checksum: event.Needle.Checksum,
			VolumeId: uint32(event.Needle.VolumeId),
		},
		VolumeServer: &volume_server_pb.VolumeServerEventResponse_VolumeServer{
			Url:        event.VolumeServer.Url,
			Rack:       event.VolumeServer.Rack,
			DataCenter: event.VolumeServer.DataCenter,
		},

		CreatedAt:   timestamppb.New(time.Unix(0, int64(event.CreatedAt))),
		LastUpdated: timestamppb.New(time.Unix(0, int64(event.LastUpdated))),
		LastTouched: timestamppb.New(time.Unix(0, int64(event.LastTouched))),
	}

	return resp
}

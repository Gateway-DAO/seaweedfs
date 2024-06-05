package weed_server

import (
	"time"

	"github.com/seaweedfs/seaweedfs/weed/event"
	event_types "github.com/seaweedfs/seaweedfs/weed/event/types"
	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
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

		if req.VolumeId != nil && event.Needle.VolumeId != needle.VolumeId(*req.VolumeId) {
			continue
		}

		parsedEvent := prepareVolumeServerEventResponse(event)
		if streamErr := stream.Send(parsedEvent); streamErr != nil {
			return streamErr
		}
	}

	return nil
}

func prepareVolumeServerEventResponse(ne *event_types.VolumeNeedleEvent) *volume_server_pb.VolumeServerEventResponse {
	resp := &volume_server_pb.VolumeServerEventResponse{
		Type: ne.Type,
		Hash: ne.Hash,
		Needle: &volume_server_pb.VolumeServerEventResponse_Needle{
			Id:       ne.Needle.Id,
			Checksum: ne.Needle.Checksum,
			VolumeId: uint32(ne.Needle.VolumeId),
		},
		VolumeServer: &volume_server_pb.VolumeServerEventResponse_VolumeServer{
			Url:        ne.VolumeServer.Url,
			Rack:       ne.VolumeServer.Rack,
			DataCenter: ne.VolumeServer.DataCenter,
		},

		CreatedAt:   timestamppb.New(time.Unix(0, int64(ne.CreatedAt))),
		LastUpdated: timestamppb.New(time.Unix(0, int64(ne.LastUpdated))),
		LastTouched: timestamppb.New(time.Unix(0, int64(ne.LastTouched))),
	}

	return resp
}

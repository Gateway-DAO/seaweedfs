package weed_server

import (
	"time"

	"github.com/seaweedfs/seaweedfs/weed/event"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"github.com/seaweedfs/seaweedfs/weed/storage"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (vs *VolumeServer) VolumeServerEvents(req *volume_server_pb.VolumeServerEventsRequest, stream volume_server_pb.VolumeServer_VolumeServerEventsServer) error {
	var vol *storage.Volume
	if req.VolumeId != nil {
		vol = vs.store.GetVolume(needle.VolumeId(*req.VolumeId))

		if vol == nil {
			return status.Errorf(codes.NotFound, "volume server does not have volume %d", *req.VolumeId)
		}
	}

	events, err := vs.eventStore.ListAllEvents()
	if err != nil {
		return status.Errorf(codes.Aborted, "error listing vs %s events: %s", vs.store.PublicUrl, err)
	}

	for _, event := range events {
		if event == nil {
			continue
		}
		glog.V(3).Infof("%s iterating through events", vs.store.PublicUrl)
		if ctxErr := stream.Context().Err(); ctxErr != nil {
			return ctxErr
		}

		if vol != nil && event.Volume.Id != vol.Id.String() {
			continue
		}

		parsedEvent := prepareVolumeServerEventResponse(event)
		if streamErr := stream.SendMsg(parsedEvent); streamErr != nil {
			return streamErr
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func prepareVolumeServerEventResponse(ne *event.VolumeServerEvent) *volume_server_pb.VolumeServerEventResponse {
	resp := &volume_server_pb.VolumeServerEventResponse{
		Type:      ne.GetType(),
		Needle:    ne.GetNeedle(),
		Volume:    ne.GetVolume(),
		CreatedAt: ne.GetCreatedAt(),
	}

	return resp
}

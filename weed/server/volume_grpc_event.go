package weed_server

import (
	"context"
	"fmt"

	"github.com/gateway-dao/seaweedfs/weed/event"
	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
	"github.com/gateway-dao/seaweedfs/weed/storage"
	"github.com/gateway-dao/seaweedfs/weed/storage/needle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func registerEvent(eventType event.VolumeServerEventType, vs *VolumeServer, volumeId *needle.VolumeId, needle *needle.Needle) error {
	switch eventType {
	case event.ALIVE:
		glog.V(3).Infof("Emitting ALIVE event for %s", vs.store.Ip)
	case event.WRITE:
		glog.V(3).Infof("Emitting WRITE event for %s", vs.store.Ip)
	case event.DELETE:
		glog.V(3).Infof("Emitting DELETE event for %s", vs.store.Ip)
	case event.VACUUM:
		glog.V(3).Infof("Emitting VACUUM event for %s", vs.store.Ip)
	default:
		return fmt.Errorf("eventType undefined")
	}

	var vse_vol *volume_server_pb.VolumeServerEventResponse_Volume
	var vse_needle *volume_server_pb.VolumeServerEventResponse_Needle
	if volumeId != nil {
		vol := vs.store.GetVolume(*volumeId)
		datSize, idxSize, lastModTime := vol.FileStat()

		vse_vol = &volume_server_pb.VolumeServerEventResponse_Volume{
			Id:           volumeId.String(),
			FileCount:    vol.FileCount(),
			IdxSize:      idxSize,
			DatSize:      datSize,
			DeletedCount: vol.DeletedCount(),
			DeletedSize:  vol.DeletedSize(),
			LastModified: timestamppb.New(lastModTime),
			Replication:  vol.ReplicaPlacement.String(),
		}

		if needle != nil {
			// needle_hash := hash.ToString()
			vse_needle = &volume_server_pb.VolumeServerEventResponse_Needle{
				Id:       uint64(needle.Id),
				Checksum: needle.Checksum.Value(),
				// Hash:     &needle_hash,
				Fid: needle.Id.FileId(uint32(*volumeId)),
			}
		}
	}

	vsStatus, vsStatus_err := vs.VolumeServerStatus(context.Background(), nil)
	if vsStatus_err != nil {
		return fmt.Errorf("unable to load volume server stats, %s", vsStatus_err)
	}

	vsMerkleTree := &volume_server_pb.VolumeServerMerkleTree{
		Digest: vsStatus.GetChecksum(),
		Tree:   map[string]string{},
	}
	for _, diskStatus := range vsStatus.GetDiskStatuses() {
		for key, value := range diskStatus.Checksum {
			vsMerkleTree.Tree[key] = value
		}
	}

	vse, vse_err := event.NewVolumeServerEvent(
		eventType,
		&volume_server_pb.VolumeServerEventResponse_Server{
			Tree:       vsMerkleTree,
			PublicUrl:  vs.store.PublicUrl,
			Rack:       vsStatus.GetRack(),
			DataCenter: vsStatus.GetDataCenter(),
		},
		vse_vol,
		vse_needle,
	)
	if vse_err != nil {
		return vse_err
	}

	if err := vs.eventStore.RegisterEvent(vse); err != nil {
		// Terminate the server if the event is unable to be logged
		glog.Fatalf("unable to register EDV event: %s", err)
	}

	return nil
}

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
	glog.V(3).Infof("volume server has %d events", len(events))

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

		if streamErr := stream.SendMsg(event); streamErr != nil {
			return status.Error(codes.Aborted, streamErr.Error())
		}
	}

	return nil
}

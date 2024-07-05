package weed_server

import (
	"context"
	"fmt"

	"github.com/seaweedfs/seaweedfs/weed/event"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (vs *VolumeServer) registerEvent(eventType event.NeedleEventType, volumeId needle.VolumeId, needle *needle.Needle, hash *string) error {
	switch eventType {
	case event.WRITE:
		glog.V(3).Infof("Emitting WRITE event for %s", vs.store.Ip)
	case event.DELETE:
		glog.V(3).Infof("Emitting DELETE event for %s", vs.store.Ip)
	case event.VACUUM:
		glog.V(3).Infof("Emitting DELETE event for %s", vs.store.Ip)
	default:
		return fmt.Errorf("eventType undefined")
	}

	vol := vs.store.GetVolume(volumeId)
	datSize, idxSize, lastModTime := vol.FileStat()

	var vse_needle *volume_server_pb.VolumeServerEventResponse_Needle
	if needle != nil {
		vse_needle = &volume_server_pb.VolumeServerEventResponse_Needle{
			Id:       uint64(needle.Id),
			Checksum: needle.Checksum.Value(),
			Hash:     hash,
		}
	}

	vsStatus, vsStatus_err := vs.VolumeServerStatus(context.Background(), nil)
	if vsStatus_err != nil {
		return fmt.Errorf("unable to load volume server stats, %s", vsStatus_err)
	}

	vsEventChecksum := &volume_server_pb.VolumeServerEventChecksum{
		Digest: vsStatus.GetChecksum(),
		Tree:   map[string]string{},
	}
	for _, diskStatus := range vsStatus.GetDiskStatuses() {
		for key, value := range diskStatus.Checksum {
			vsEventChecksum.Tree[key] = value
		}
	}

	vse, vse_err := event.NewVolumeServerEvent(
		eventType,
		vse_needle,
		&volume_server_pb.VolumeServerEventResponse_Volume{
			Id:           volumeId.String(),
			IdxSize:      idxSize,
			FileCount:    vol.FileCount(),
			DatSize:      datSize,
			DeletedCount: vol.DeletedCount(),
			DeletedSize:  vol.DeletedSize(),
			LastModified: timestamppb.New(lastModTime),
			Replication:  vol.ReplicaPlacement.String(),

			Server: &volume_server_pb.VolumeServerEventResponse_Volume_Server{
				Checksum:   vsEventChecksum,
				PublicUrl:  vs.store.PublicUrl,
				Rack:       vsStatus.GetRack(),
				DataCenter: vsStatus.GetDataCenter(),
			},
		},
	)
	if vse_err != nil {
		return vse_err
	}

	return vs.eventStore.RegisterEvent(vse)
}

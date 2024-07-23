package weed_server

import (
	"context"
	"fmt"

	"github.com/gateway-dao/seaweedfs/weed/event"
	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
	"github.com/gateway-dao/seaweedfs/weed/storage/needle"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (vs *VolumeServer) registerEvent(eventType event.VolumeServerEventType, volumeId needle.VolumeId, fid *string, needle *needle.Needle, hash *string) error {
	switch eventType {
	case event.WRITE:
		glog.V(3).Infof("Emitting WRITE event for %s", vs.store.Ip)
	case event.DELETE:
		glog.V(3).Infof("Emitting DELETE event for %s", vs.store.Ip)
	case event.VACUUM:
		glog.V(3).Infof("Emitting VACUUM event for %s", vs.store.Ip)
	default:
		return fmt.Errorf("eventType undefined")
	}

	vol := vs.store.GetVolume(volumeId)
	datSize, idxSize, lastModTime := vol.FileStat()

	var vse_needle *volume_server_pb.VolumeServerEventResponse_Needle
	if needle != nil {
		if fid == nil {
			return fmt.Errorf("fid not specified for needle event")
		}
		vse_needle = &volume_server_pb.VolumeServerEventResponse_Needle{
			Id:       uint64(needle.Id),
			Fid:      *fid,
			Checksum: needle.Checksum.Value(),
			Hash:     hash,
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
		&volume_server_pb.VolumeServerEventResponse_Volume{
			Id:           volumeId.String(),
			IdxSize:      idxSize,
			FileCount:    vol.FileCount(),
			DatSize:      datSize,
			DeletedCount: vol.DeletedCount(),
			DeletedSize:  vol.DeletedSize(),
			LastModified: timestamppb.New(lastModTime),
			Replication:  vol.ReplicaPlacement.String(),
		},
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

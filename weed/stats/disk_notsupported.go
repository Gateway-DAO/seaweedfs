//go:build netbsd || plan9
// +build netbsd plan9

package stats

import "github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"

func fillInDiskStatus(status *volume_server_pb.DiskStatus) {
	return
}

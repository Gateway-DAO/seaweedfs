//go:build !windows && !openbsd && !netbsd && !plan9 && !solaris
// +build !windows,!openbsd,!netbsd,!plan9,!solaris

package stats

import (
	"encoding/hex"
	"syscall"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
)

func fillInDiskStatus(disk *volume_server_pb.DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(disk.Dir, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	disk.PercentFree = float32((float64(disk.Free) / float64(disk.All)) * 100)
	disk.PercentUsed = float32((float64(disk.Used) / float64(disk.All)) * 100)

	checksum, err := computeDiskChecksum(disk)
	if err != nil {
		checksum = []byte(err.Error())
	}
	disk.Checksum = hex.EncodeToString(checksum)

	return
}

func computeDiskChecksum(disk *volume_server_pb.DiskStatus) ([]byte, error) {
	timeStart := time.Now()

	hash, err := hashFilteredDirectory(disk.Dir, `\.dat$`)
	if err != nil {
		return nil, err
	}

	timeDuration := float64(time.Since(timeStart).Milliseconds())
	VolumeServerChecksumDuration.Set(timeDuration)

	return hash, nil
}

//go:build !windows && !openbsd && !netbsd && !plan9 && !solaris
// +build !windows,!openbsd,!netbsd,!plan9,!solaris

package stats

import (
	"syscall"
	"time"

	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
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
	if err == nil {
		disk.Checksum = checksum
	}
	// else err != nil {
	// 	checksum = map[string]string{"error": err.Error()}
	// }
	return
}

type DiskChecksum map[string]string

func computeDiskChecksum(disk *volume_server_pb.DiskStatus) (DiskChecksum, error) {
	timeStart := time.Now()

	hashes, err := hashFilteredDirectory(disk.Dir, `\.dat$`)
	if err != nil {
		return nil, err
	}

	formattedHashes := make(DiskChecksum, len(hashes))

	for k, v := range hashes {
		formattedHashes[k] = v.ToString()
	}

	timeDuration := float64(time.Since(timeStart).Milliseconds())
	VolumeServerChecksumDuration.Set(timeDuration)

	return formattedHashes, nil
}

func (d DiskChecksum) convertToString() (result string) {
	for _, v := range d {
		result += v
	}

	return
}

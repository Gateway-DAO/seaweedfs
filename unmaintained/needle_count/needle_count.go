package main

import (
	"fmt"

	// "github.com/seaweedfs/seaweedfs/weed/server/volume"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
)

func countNeedlesForFileId(fileId string) (int, error) {
	fid, err := needle.ParseFileIdFromString(fileId)
	if err != nil {
		return 0, fmt.Errorf("invalid file id: %v", err)
	}

	nId, cookie, err := needle.ParseNeedleIdCookie(fid.GetNeedleIdCookie())

	nv := needle.Needle{
		Id:     nId,
		Cookie: cookie,
	}

	if err != nil {
		return 0, fmt.Errorf("failed to read needle: %v", err)
	}

	if nv.IsChunkedManifest() {
		fmt.Printf("THIS IS A CHUNK")
		return 2, nil
	}

	// This is a single-needle file
	return 1, nil
}

func main() {
	// This is a simplified example. In a real scenario, you'd need to:
	// 1. Connect to the master server
	// 2. Lookup the volume for the given file ID
	// 3. Connect to the volume server
	// 4. Open the volume

	// For demonstration, let's assume we have a volume object
	// v := &volume.Volume{} // This should be properly initialized

	fileId := "1,01aa581f2d" // Replace with actual file ID
	count, err := countNeedlesForFileId(fileId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Number of needles for file ID %s: %d\n", fileId, count)
}

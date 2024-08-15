package storage

import (
	"fmt"
	"testing"

	"github.com/gateway-dao/seaweedfs/weed/stats"
)

func Test_MerkleNodeFromEventTree(t *testing.T) {
	// mn := MerkleNodeFromEventTree(&test_event)
	decoded, _ := stats.DecodeString(test_event.Digest)

	str := string(decoded)
	fmt.Printf("mn.value: %s\n", str)
}

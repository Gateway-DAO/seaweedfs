package storage

import (
	"reflect"
	"testing"

	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
	"github.com/gateway-dao/seaweedfs/weed/stats"
	"github.com/stretchr/testify/assert"
)

var tree = &MerkleTree{
	// root analogous to VolumeServer store
	root: &MerkleNode{
		// children analogous to VolumeServer volume(s)
		children: map[string]*MerkleNode{
			"0": {
				value: stats.Hash("hash of volume 0"),
			},
			"1": {
				value: stats.Hash("hash of volume 1"),
				children: map[string]*MerkleNode{
					"3": {
						value: stats.Hash("volume 3 hash"),
					},
				},
			},
		},
	},
}

var event = event_pb.MerkleTree{
	Digest: "Tk89xUWvQLfcqum2HzfORNoU97hyaCtdeULc3kwVW0A",
	Tree: map[string]*event_pb.MerkleTree{
		"1": {
			Digest: "SE6hGC8xlLJOnL2Da9I8V+UjWwVnVOkh/BZEq7Ayt9k",
			Tree: map[string]*event_pb.MerkleTree{
				"1": {
					Digest: "3Ln6yXb3oEiaTmB+ypJrPK2CcTPYn8ABhNgapItexFs",
					Tree:   map[string]*event_pb.MerkleTree{},
				},
			},
		},
	},
}

func Test_EventToMerkleTree(t *testing.T) {
	tree := FromProto(&event)

	// validate independent node elements
	assert.Equal(t, tree != nil, true)
	assert.Equal(t, tree.root != nil, true)
	compareNodeToEvent(t, tree.root, &event)

	// validate reconstruction
	assert.Equal(t, reflect.DeepEqual(&event, tree.ToProto()), true)
}

func compareNodeToEvent(t *testing.T, n *MerkleNode, e *event_pb.MerkleTree) {
	assert.Equal(t, n == nil, e == nil)

	if n != nil {
		assert.Equal(t, len(n.children), len(e.Tree))
		assert.Equal(t, n.Value().EncodeToString(), e.Digest)

		// merkle node children
		nkeys := make([]string, 0, len(n.children))
		for k := range n.children {
			nkeys = append(nkeys, k)
		}

		// event tree children
		ekeys := make([]string, 0, len(e.Tree))
		for k := range e.Tree {
			ekeys = append(ekeys, k)
		}

		// iteratively validate children
		for ix := range len(n.children) {
			assert.Equal(t, nkeys[ix], ekeys[ix])

			compareNodeToEvent(t, n.children[nkeys[ix]], e.Tree[ekeys[ix]])
		}
	}
}

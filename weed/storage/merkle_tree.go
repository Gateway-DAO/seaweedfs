package storage

import (
	"fmt"
	"strings"

	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
	"github.com/gateway-dao/seaweedfs/weed/stats"
)

// MerkleTree represents a complete Merkle tree with a single root node.
// The root node can be used to verify the integrity of all data represented
// by the tree.
type MerkleTree struct {
	root *MerkleNode
}

func NewMerkleTree(root *MerkleNode) *MerkleTree {
	return &MerkleTree{
		root: root,
	}
}

func (s *Store) MerkleTree() (mt *MerkleTree) {
	return &MerkleTree{
		root: s.MerkleNode(),
	}
}

func (mt *MerkleTree) RootValue() stats.Hash {
	if mt.root == nil {
		return nil
	}

	return mt.root.Value()
}

func (mt *MerkleTree) ToProto() *event_pb.MerkleTree {
	if mt == nil || mt.root == nil {
		return nil
	}

	return mt.root.toProto()
}

func FromProto(e *event_pb.MerkleTree) *MerkleTree {
	if e == nil {
		return nil
	}

	return &MerkleTree{
		root: fromProto(e),
	}
}

func (mt *MerkleNode) ToString(indent string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sDigest: %s\n", indent, mt.Value().EncodeToString()))
	if mt.children != nil && len(mt.children) > 0 {
		sb.WriteString(fmt.Sprintf("%sTree:\n", indent))
		for key, subtree := range mt.children {
			sb.WriteString(fmt.Sprintf("%s  Key: %s\n", indent, key))
			sb.WriteString(subtree.ToString(indent + "    "))
		}
	}
	return sb.String()
}

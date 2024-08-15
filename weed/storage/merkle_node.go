package storage

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/pb/event_pb"
	"github.com/gateway-dao/seaweedfs/weed/stats"
	"github.com/gateway-dao/seaweedfs/weed/storage/needle"
	"github.com/gateway-dao/seaweedfs/weed/storage/types"
)

// MerkleNode represents a node in a Merkle tree, which is a binary tree used
// to verify data integrity. Each node contains a `value` which is the hash
// of the node's data, and potentially a list of `children` representing the
// child nodes in the tree.
//
// If the `value` field is set (non-nil), it represents the hash for this node
// and will be returned directly by the `Value()` method. If `value` is nil,
// the `Value()` method will compute the hash by aggregating the Blake2b
// checksums of all its children. This allows for efficient Merkle tree
// traversal and verification.

type MerkleNode struct {
	isLeaf   bool
	value    stats.Hash
	children map[string]*MerkleNode
}

func NewMerkleNode(value stats.Hash) *MerkleNode {
	return &MerkleNode{
		value: value,
	}
}

func MerkleNodeFromEventTree(e *event_pb.MerkleTree) *MerkleNode {
	digest, err := stats.DecodeString(e.Digest)
	if err != nil {
		glog.Error("Unable to decode event digest as hash bytes")
		return nil
	}
	mn := &MerkleNode{
		value: digest,
	}

	if len(e.Tree) > 0 {
		mn.children = make(map[string]*MerkleNode)
		for k, v := range e.Tree {
			mn.children[k] = MerkleNodeFromEventTree(v)
		}
	} else {
		mn.isLeaf = true
	}

	return mn
}

func (n *MerkleNode) AddChild(key string, nn *MerkleNode) bool {
	if n.children == nil {
		n.children = make(map[string]*MerkleNode)
	}
	if n.children[key] != nil {
		return false
	}

	n.children[key] = nn

	return true
}

func (s *Store) MerkleNode() *MerkleNode {
	storeVolumeInfos := s.VolumeInfos()

	mn := new(MerkleNode)
	mn.children = make(map[string]*MerkleNode, len(storeVolumeInfos))

	wg := sync.WaitGroup{}
	for i, vi := range storeVolumeInfos {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if vi == nil {
				glog.Errorf("VolumeInfo index %d is not found", i)
			}

			vid := vi.Id
			v := s.GetVolume(vid)

			if v == nil {
				glog.Errorf("Volume %d is not found", vid)
			}

			mn.children[vid.String()] = v.MerkleNode()
		}()
	}

	wg.Wait()

	return mn
}

func (v *Volume) MerkleNode() (mn *MerkleNode) {
	mn = new(MerkleNode)
	mn.children = make(map[string]*MerkleNode, v.nm.MaxFileKey())

	glog.V(3).Infof("Computing MerkleTree for Volume %s", v.Id)

	for nId := types.NeedleId(1); nId <= v.nm.MaxFileKey(); nId++ {
		nv, ok := v.nm.Get(nId)

		if !ok {
			glog.Errorf("Unable to query volume %s needleMapper for needleId %d", v.Id, nId)
			continue
		}

		n := &needle.Needle{
			Id:   nId,
			Size: nv.Size,
		}

		err := n.ReadData(v.DataBackend, nv.Offset.ToActualOffset(), nv.Size, v.Version())
		if err != nil {
			glog.Errorf("Unable to read data for needle config")
		}

		hasher, _ := stats.Blake2b()
		hasher.Write(n.Data)
		checksum := hasher.Sum(nil)

		mn.children[nId.String()] = &MerkleNode{
			isLeaf: true,
			value:  checksum,
		}
	}

	return
}

// Value returns the hash value of the MerkleNode. If the `value` field is
// set (non-nil), this method returns it directly. If `value` is nil, the method
// computes the hash by aggregating the Blake2b checksums of its children.
//
// The aggregation is performed by iterating over all child nodes and computing
// the cumulative hash. If a child node is nil, an error is logged using `glog.Error`.
// The method ensures that the node's hash is always consistent with its
// children's hashes, facilitating integrity verification in the Merkle tree.
//
// Returns:
//
//	stats.Hash: The hash of the node, either precomputed and stored in `value`
//	            or computed based on the children's hashes.
func (mn *MerkleNode) Value() stats.Hash {
	if mn == nil {
		return nil
	}

	if mn.isLeaf {
		return mn.value
	}

	hasher, _ := stats.Blake2b()
	if mn.children != nil {
		for _, v := range mn.children {
			if v == nil {
				glog.Error("MerkleNode has nil child")
				continue
			}
			hasher.Write(v.Value())
		}
	}
	checksum := hasher.Sum(nil)

	return checksum
}

func (mn *MerkleNode) toProto() *event_pb.MerkleTree {
	if mn == nil {
		fmt.Println("MerkleNode is null")
		return nil
	}

	tree := new(event_pb.MerkleTree)
	tree.Digest = mn.Value().EncodeToString()
	tree.Tree = make(map[string]*event_pb.MerkleTree)

	for k, child := range mn.children {
		tree.Tree[k] = child.toProto()
	}

	return tree
}

func (mn *MerkleNode) ValidateNode() bool {
	if mn == nil {
		return false
	}

	value := mn.value

	hasher, _ := stats.Blake2b()
	for _, v := range mn.children {
		if v.ValidateNode() {
			hasher.Write(v.value)
		} else {
			return false
		}
	}
	digest := hasher.Sum(nil)

	return reflect.DeepEqual(value, digest)
}

func fromProto(e *event_pb.MerkleTree) *MerkleNode {
	node := new(MerkleNode)
	node.value, _ = stats.DecodeString(e.GetDigest())
	node.children = make(map[string]*MerkleNode)

	for k, v := range e.Tree {
		node.children[k] = fromProto(v)
	}

	return node
}

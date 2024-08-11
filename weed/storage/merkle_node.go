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
	value    stats.Hash
	children map[string]*MerkleNode
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
	mn.value = mn.Value()

	return mn
}

func (v *Volume) MerkleNode() (mn *MerkleNode) {
	mn = new(MerkleNode)
	mn.children = make(map[string]*MerkleNode, v.nm.MaxFileKey())

	glog.V(3).Infof("Computing MerkleTree for Volume %s", v.Id)

	hasher, _ := stats.Blake2b()
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

		hasher.Reset()
		hasher.Write(n.Data)
		checksum := hasher.Sum(nil)

		mn.children[nId.String()] = &MerkleNode{
			value: checksum,
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

	if mn.value != nil {
		return mn.value
	}

	hasher, _ := stats.Blake2b()
	if mn.children != nil {
		for k, v := range mn.children {
			if v == nil {
				glog.Error("MerkleNode has nil child")
			}
			var val []byte = merkleNodeChildrenAsBytes(k, v)
			hasher.Write(val)
		}
	}
	checksum := hasher.Sum(nil)

	mn.value = checksum
	return checksum
}

func merkleNodeChildrenAsBytes(k string, v *MerkleNode) []byte {
	return []byte(k + ":" + v.Value().EncodeToString())
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
		tree.Tree[k] = new(event_pb.MerkleTree)
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
	for k, v := range mn.children {
		if v.ValidateNode() {
			hasher.Write(merkleNodeChildrenAsBytes(k, v))
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

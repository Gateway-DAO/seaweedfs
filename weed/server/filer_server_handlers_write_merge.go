package weed_server

import (
	"github.com/gateway-dao/seaweedfs/weed/operation"
	"github.com/gateway-dao/seaweedfs/weed/pb/filer_pb"
)

func (fs *FilerServer) maybeMergeChunks(so *operation.StorageOption, inputChunks []*filer_pb.FileChunk) (mergedChunks []*filer_pb.FileChunk, err error) {
	//TODO merge consecutive smaller chunks into a large chunk to reduce number of chunks
	return inputChunks, nil
}

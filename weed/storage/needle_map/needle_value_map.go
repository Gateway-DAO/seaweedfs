package needle_map

import (
	. "github.com/gateway-dao/seaweedfs/weed/storage/types"
)

type NeedleValueMap interface {
	Set(key NeedleId, offset Offset, size Size) (oldOffset Offset, oldSize Size)
	Delete(key NeedleId) Size
	Get(key NeedleId) (*NeedleValue, bool)
	AscendingVisit(visit func(NeedleValue) error) error
}

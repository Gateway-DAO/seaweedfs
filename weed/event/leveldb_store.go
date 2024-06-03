package event

import (
	"github.com/seaweedfs/seaweedfs/weed/event/types"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

func connectToEventStore(dir string) (*leveldb.DB, error) {
	glog.Infof("Connecting to leveldb dir: %v", dir)
	return leveldb.OpenFile(dir, nil)
}

func RegisterEvent(dbDir string, n *types.NeedleEvent) error {
	glog.V(3).Infof("Writing to database %s", dbDir)

	db, err := connectToEventStore(dbDir)
	if err != nil {
		return err
	}
	defer db.Close()

	val, ve := n.Value()
	if ve != nil {
		return ve
	}
	db.Put([]byte(n.Key()), val, nil)

	return nil
}

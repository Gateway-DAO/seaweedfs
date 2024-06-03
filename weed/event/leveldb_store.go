package event

import (
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

type EventStoreLevelDB struct {
	dir string
}

func GetLevelDBEventStore(dirPath string) *EventStoreLevelDB {
	glog.V(3).Infof("Connecting to event store: %s", dirPath)
	return &EventStoreLevelDB{
		dir: dirPath,
	}
}

func (e *EventStoreLevelDB) connectToEventStore() (*leveldb.DB, error) {
	glog.Infof("Connecting to leveldb dir: %v", e.dir)
	return leveldb.OpenFile(e.dir, nil)
}

func RegisterEvent(dbDir string) error {
	glog.V(3).Infof("Writing to database %s", dbDir)
	return nil
}

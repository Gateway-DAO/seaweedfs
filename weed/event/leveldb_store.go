package event

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbEventStore struct {
	sync.RWMutex

	Dir  string
	size uint64
}

func NewLevelDbEventStore(eventDir string) (*LevelDbEventStore, error) {
	return &LevelDbEventStore{
		Dir:  eventDir,
		size: 0,
	}, nil
}

func (es *LevelDbEventStore) RegisterEvent(e *VolumeServerEvent) error {
	es.Lock()
	defer es.Unlock()

	dbDir := es.Dir

	glog.V(4).Infof("Writing to database %s", dbDir)

	db, err := leveldb.OpenFile(es.Dir, nil)
	if err != nil {
		return fmt.Errorf("unable to connect to event store: %s", err)
	}
	defer db.Close()

	val, ve := e.Value()
	if ve != nil {
		return ve
	}

	db.Put(
		timestampToBytes(time.Now().UnixNano()),
		val,
		nil,
	)
	es.size++

	return nil
}

func (es *LevelDbEventStore) ListAllEvents() ([]*VolumeServerEvent, error) {
	es.RLock()
	defer es.RUnlock()

	dbDir := es.Dir
	glog.V(4).Infof("Reading database %s", dbDir)

	db, err := leveldb.OpenFile(es.Dir, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to event store: %s", err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	results := make([]*VolumeServerEvent, es.size)

	for iter.Next() {
		key, val := iter.Key(), iter.Value()

		valPtr := new(VolumeServerEvent)
		if err := json.Unmarshal(val, valPtr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value for key %s: %v", string(key), val)
		}

		results = append(results, valPtr)
	}

	// Check for errors encountered during iteration
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("cannot connect to event dir %s", dbDir)
	}

	return results, nil
}

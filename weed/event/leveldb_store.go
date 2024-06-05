package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/event/event_types"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

func connectToEventStore(dir string) (*leveldb.DB, error) {
	glog.Infof("Connecting to leveldb dir: %v", dir)
	return leveldb.OpenFile(dir, nil)
}

func withEventStoreConnection(dir string, handler func(db *leveldb.DB)) error {
	if db, err := connectToEventStore(dir); err == nil {
		handler(db)
		db.Close()

		return nil
	} else {
		return err
	}
}

func RegisterEvent(dbDir string, ne *event_types.NeedleEvent) error {
	glog.V(3).Infof("Writing to database %s", dbDir)

	db, err := connectToEventStore(dbDir)
	if err != nil {
		return err
	}
	defer db.Close()

	val, ve := ne.Value()
	if ve != nil {
		return ve
	}

	db.Put(
		[]byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		val,
		nil,
	)

	return nil
}

func ListEvents(dbDir string) (map[string]*event_types.NeedleEvent, error) {
	glog.V(3).Infof("Reading database %s", dbDir)

	db, err := connectToEventStore(dbDir)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	results := make(map[string]*event_types.NeedleEvent)

	for iter.Next() {
		key, val := iter.Key(), iter.Value()

		valPtr := new(event_types.NeedleEvent)
		if err := json.Unmarshal(val, valPtr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value for key %s: %v", key, val)
		}

		results[string(key)] = valPtr
	}

	// Check for errors encountered during iteration
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("cannot connect to event dir %s", dbDir)
	}

	return results, nil
}

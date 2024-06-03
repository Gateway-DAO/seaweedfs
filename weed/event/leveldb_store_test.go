package event

import (
	"encoding/json"
	"testing"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	eventStore = GetLevelDBEventStore("/Users/srv/Developer/GTW/DFS/seaweedfs/data/events")
	testEvent  = FileEvent{
		// Id: fmt.Sprintf("%d", rand.Int31()),
		Fid:    "leveldb_store_test:3,063de3afa9",
		Status: UPLOADED,
		Hash:   "asdf",
	}
)

func Test_RegisterEvent(t *testing.T) {
	if eventStore.RegisterEvent(testEvent) != nil {
		t.Fail()
	}
}

func Test_RegisteredEvent(t *testing.T) {
	event, get_err := eventStore.GetValue(testEvent.Fid)
	if get_err != nil {
		glog.Errorf("Error getting value: %s", get_err)
		t.Fail()
	}

	if (*event).Metadata() != testEvent.Metadata() {
		glog.Errorf("Queried FileEvent does not match")
		t.Fail()
	}
}

func Test_Cleanup(_ *testing.T) {
	db, _ := leveldb.OpenFile(eventStore.dir, nil)
	defer db.Close()

	mKey, _ := json.Marshal(testEvent.Fid)
	db.Delete(mKey, nil)
}

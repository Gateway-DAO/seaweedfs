package event

import (
	"fmt"
	"testing"

	"github.com/seaweedfs/seaweedfs/weed/event/event_types"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"github.com/seaweedfs/seaweedfs/weed/storage/types"
)

var (
	eventDir = "./test_data/events"
)

func newMockNeedle(id *uint32, checksum *[]byte) (n *needle.Needle) {
	n = new(needle.Needle)

	if id != nil {
		n.Id = types.NeedleId(*id)
	} else {
		n.Id = types.NeedleId(0)
	}

	if checksum != nil {
		n.Checksum = needle.NewCRC(*checksum)
	} else {
		n.Checksum = needle.NewCRC([]byte("test-checksum"))
	}

	return
}

func Test_RegisterWriteEvent(t *testing.T) {
	hash := "stubbed-hash"

	event, err := event_types.NewNeedleEvent(
		needle.VolumeId(1),
		"stubbed-ip",
		"stubbed-data_center",
		"stubbed-rack",
		newMockNeedle(nil, nil),
		&hash,
		event_types.WRITE,
	)

	if err != nil {
		t.Errorf("unable to construct event: %s", err)
	}

	RegisterEvent(eventDir, event)
}

func Test_RegisterDeleteEvent(t *testing.T) {
	event, err := event_types.NewNeedleEvent(
		needle.VolumeId(1),
		"stubbed-ip",
		"stubbed-data_center",
		"stubbed-rack",
		newMockNeedle(nil, nil),
		nil,
		event_types.WRITE,
	)

	if err != nil {
		t.Errorf("unable to construct event: %s", err)
	}

	RegisterEvent(eventDir, event)
}

func Test_GetEvents(t *testing.T) {
	events, err := ListEvents(eventDir)
	if err != nil {
		t.Errorf("unable to delineate events dir: %s", err)
	}

	for k, v := range events {
		fmt.Printf("Key: %s, Value: %+v\n", k, v)
	}
}

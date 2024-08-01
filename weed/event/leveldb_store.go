package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/pb/volume_server_pb"
	"github.com/gateway-dao/seaweedfs/weed/stats"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbEventStore struct {
	EventStoreImpl

	Dir  string
	size uint64

	kafkaStore *KafkaStore
}

func NewLevelDbEventStore(
	eventDir string,
	kafkaBrokers *[]string,
	kafkaTopics *KafkaStoreTopics,
	config *sarama.Config,
) (*LevelDbEventStore, error) {
	es := &LevelDbEventStore{
		Dir:  eventDir,
		size: 0,
	}

	if kafkaBrokers != nil && kafkaTopics != nil {
		for {
			producer, err := sarama.NewSyncProducer(*kafkaBrokers, config)

			if err != nil {
				glog.Errorf("Unable to connect to brokers: %v", err)
				time.Sleep(1790 * time.Millisecond)
				continue
			}

			es.kafkaStore = NewKafkaStore(*kafkaBrokers, *kafkaTopics, config, producer)

			glog.V(3).Infof("connected to brokers: %s", kafkaBrokers)
			break
		}
	}

	return es, nil
}

func (es *LevelDbEventStore) RegisterEvent(e *VolumeServerEvent) error {
	if e == nil {
		return fmt.Errorf("server event is nil")
	}

	// Collect last event's hash
	var lastHash *string
	lastEvent, lastEventErr := es.GetLastEvent()
	if lastEventErr != nil || lastEvent.ProofOfHistory == nil {
		glog.V(3).Infof("unable to find previous event. emitting GENESIS event")
		e.Type = "GENESIS"
	} else {
		lastHash = &lastEvent.ProofOfHistory.Hash
	}

	es.Lock()
	defer es.Unlock()

	hasher, hash_err := stats.Blake2b()
	if hash_err != nil {
		return hash_err
	}
	if e.Type != "GENESIS" && lastHash != nil {
		hasher.Write([]byte(*lastHash))
	}

	val, ve := e.Value()
	if ve != nil {
		return ve
	}
	checksumBytes, err := stats.HashFromString(e.GetServer().GetTree().GetDigest())
	if err != nil {
		glog.Errorf("error decoding server checksum digest")
	}
	hasher.Write(checksumBytes)

	// update with proof of history metadata
	e.ProofOfHistory = &volume_server_pb.VolumeServerEventResponse_ProofOfHistory{
		PreviousHash: lastHash,
		Hash:         stats.Hash(hasher.Sum(nil)).ToString(),
	}
	val, ve = e.Value()
	if ve != nil {
		return ve
	}

	if es.kafkaStore != nil {
		go func() {
			glog.V(3).Infof("writing to kafka stream")

			kafkaKey := EventKafkaKey{
				Server: e.GetServer().PublicUrl,
			}
			if e.GetVolume() != nil {
				kafkaKey.Volume = e.GetVolume().Id
			}
			kafkaEncodedKey, err := json.Marshal(kafkaKey)
			if err != nil {
				glog.Errorf("unable to encode kafkaKey")
			}

			_, _, err = es.kafkaStore.Publish(
				"volume-server",
				kafkaEncodedKey,
				val,
			)
			if err != nil {
				glog.Errorf("unable to publish to kafka topic: %s", err)
			} else {
				glog.Infof("successfully published to topic")
			}
		}()
	}

	dbDir := es.Dir
	glog.V(4).Infof("Writing to database %s", dbDir)

	db, err := leveldb.OpenFile(es.Dir, nil)
	if err != nil {
		return fmt.Errorf("unable to connect to event store: %s", err)
	}
	defer db.Close()

	db.Put(
		timestampToBytes(time.Now().UnixNano()),
		val,
		nil,
	)
	es.size++

	return nil
}

func (es *LevelDbEventStore) GetLastEvent() (*VolumeServerEvent, error) {
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

	if iter.Last() {
		key, val := iter.Key(), iter.Value()

		valPtr := new(VolumeServerEvent)
		if err := json.Unmarshal(val, valPtr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value for key %s: %v", string(key), val)
		}
		glog.V(3).Infof("%v", valPtr)

		return valPtr, nil
	}

	return nil, fmt.Errorf("no events found")
}

func (es *LevelDbEventStore) ListAllEvents() ([]*VolumeServerEvent, error) {
	es.RLock()
	glog.V(3).Info("acquired read lock")
	defer es.RUnlock()
	defer glog.V(3).Infof("released read lock")

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

package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/stats"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbEventStore[T Event] struct {
	EventStore[T]
	mu sync.RWMutex

	Dir  string
	db   *leveldb.DB
	size uint64

	kafkaStore *KafkaStore
	kafkaTopic *string
}

func NewLevelDbEventStore[T Event](
	eventDir string,
	kafkaBrokers *[]string,
	kafkaTopic *string,
	config *sarama.Config,
) (*LevelDbEventStore[T], error) {
	es := &LevelDbEventStore[T]{
		Dir:  eventDir,
		size: 0,
	}

	dbDir := es.Dir
	glog.V(2).Infof("Reading event store %s", dbDir)

	db, err := leveldb.OpenFile(es.Dir, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to event store: %s", err)
	}
	es.db = db

	if kafkaBrokers != nil && kafkaTopic != nil {
		for {
			producer, err := sarama.NewSyncProducer(*kafkaBrokers, config)

			if err != nil {
				glog.Errorf("Unable to connect to brokers: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			es.kafkaStore = NewKafkaStore(*kafkaBrokers, config, producer)
			es.kafkaTopic = kafkaTopic

			glog.V(2).Infof("Connected to kafka brokers: %s", kafkaBrokers)
			break
		}
	}

	return es, nil
}

func (es *LevelDbEventStore[T]) RegisterEvent(e T) error {
	// Collect last event's hash
	var lastHash *string
	lastEvent, lastEventErr := es.GetLastEvent()

	if lastEventErr != nil {
		if e.isAliveType() && errors.Is(lastEventErr, LastEventNotFoundError) {
			glog.V(3).Info("Unable to find previous healthcheck event")
			glog.V(2).Info("Emitting GENESIS event")
			e.SetType("GENESIS")
		} else {
			glog.Errorf("LevelDB store unable to load last event: %s", lastEventErr)
		}
	} else {
		lastHash = &((*lastEvent).GetProofOfHistory().Hash)
		glog.V(4).Infof("lastHash: %s", *lastHash)
	}

	es.mu.Lock()
	defer es.mu.Unlock()

	hasher, hash_err := stats.Blake2b()
	if hash_err != nil {
		return hash_err
	}
	if e.GetType() != "GENESIS" && lastHash != nil {
		hasher.Write([]byte(*lastHash))
	}

	val, ve := e.GetValue()
	if ve != nil {
		return ve
	}
	checksumBytes, err := stats.DecodeString(e.GetServer().GetTree().GetDigest())
	if err != nil {
		glog.Errorf("error decoding server checksum digest")
	}
	hasher.Write(checksumBytes)

	// update with proof of history metadata
	e.SetProofOfHistory(lastHash, stats.Hash(hasher.Sum(nil)).EncodeToString())

	val, ve = e.GetValue()
	if ve != nil {
		return ve
	}

	key, err := e.GetKey()
	if err != nil {
		return fmt.Errorf("error encoding event key")
	}

	if es.kafkaStore != nil && es.kafkaTopic != nil {
		go func() {
			glog.V(4).Infof("Writing to kafka stream")

			_, _, err = es.kafkaStore.Publish(
				*es.kafkaTopic,
				key,
				val,
			)
			if err != nil {
				glog.Errorf("unable to publish to kafka topic: %s", err)
			} else {
				glog.V(3).Infof("Successfully published to topic")
			}
		}()
	} else {
		glog.V(3).Infof("skip publishing kafka event; either kafkaStore or kafkaTopic is nil.")
	}

	db := es.db
	glog.V(3).Infof("Writing to event store %s", es.Dir)
	if err != nil {
		return fmt.Errorf("unable to connect to event store: %s", err)
	}

	db.Put(
		key,
		val,
		nil,
	)
	es.size++

	return nil
}

func (es *LevelDbEventStore[T]) GetLastEvent() (*T, error) {
	dbDir := es.Dir
	glog.V(4).Infof("Reading database %s", dbDir)

	es.mu.RLock()
	defer es.mu.RUnlock()

	iter := es.db.NewIterator(nil, nil)
	defer iter.Release()

	if iter.Last() {
		key, val := iter.Key(), iter.Value()

		valPtr := new(T)
		if err := json.Unmarshal(val, valPtr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value for key %s: %v", string(key), val)
		}
		glog.V(3).Infof("%+v", valPtr)

		return valPtr, nil
	}

	return nil, LastEventNotFoundError
}

func (es *LevelDbEventStore[T]) ListAllEvents() ([]T, error) {
	dbDir := es.Dir
	glog.V(4).Infof("Reading database %s", dbDir)

	glog.V(4).Info("acquired read lock")
	es.mu.RLock()

	defer es.mu.RUnlock()
	defer glog.V(4).Infof("released read lock")

	iter := es.db.NewIterator(nil, nil)
	defer iter.Release()

	results := make([]T, es.size)

	for iter.Next() {
		key, val := iter.Key(), iter.Value()

		valPtr := new(T)
		if err := json.Unmarshal(val, valPtr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value for key %s: %v", string(key), val)
		}

		results = append(results, *valPtr)
	}

	// Check for errors encountered during iteration
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("cannot connect to event dir %s", dbDir)
	}

	return results, nil
}

func (es *LevelDbEventStore[T]) Close() {
	es.db.Close()
}

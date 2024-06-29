package event

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbEventStore struct {
	sync.RWMutex

	Dir  string
	size uint64

	kafkaProducer    sarama.SyncProducer
	kafkaTopixPrefix *string
}

func NewLevelDbEventStore(eventDir string, kafkaBrokers *[]string, kafkaTopicPrefix *string) (*LevelDbEventStore, error) {
	es := &LevelDbEventStore{
		Dir:  eventDir,
		size: 0,
	}

	if kafkaBrokers != nil && kafkaTopicPrefix != nil {
		for {
			config := sarama.NewConfig()
			config.Producer.Return.Successes = true
			producer, err := sarama.NewSyncProducer(*kafkaBrokers, config)

			if err != nil {
				glog.Errorf("Unable to connect to brokers: %v", err)
				time.Sleep(1790 * time.Millisecond)
				continue
			}

			es.kafkaProducer = producer
			es.kafkaTopixPrefix = kafkaTopicPrefix
			glog.V(3).Infof("connected to brokers: %s", kafkaBrokers)
			break
		}
	}

	return es, nil
}

func (es *LevelDbEventStore) sendKafkaMessage(topic, key string, data []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := es.kafkaProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to send message to kafka producers: %s", err)
	}

	glog.V(3).Infof("kafka message successful: %d, %d", partition, offset)

	return partition, offset, nil
}

func (es *LevelDbEventStore) RegisterEvent(e *VolumeServerEvent) error {
	es.Lock()
	defer es.Unlock()

	if e == nil {
		return fmt.Errorf("server event is nil")
	}

	val, ve := e.Value()
	if ve != nil {
		return ve
	}

	if es.kafkaProducer != nil {
		glog.V(3).Infof("write to kafka stream: volume%s", e.Volume.Id)
		go es.sendKafkaMessage(
			"volume_server",
			fmt.Sprintf("volume%s", e.Volume.Id),
			val,
		)
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

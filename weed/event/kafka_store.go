package event

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/seaweedfs/seaweedfs/weed/glog"
)

type KafkaStoreTopics struct {
	Assignment string `toml:"kafka.topic.assignment"`
	Volume     string `toml:"kafka.topic.volume"`
}

type KafkaStore struct {
	brokers []string
	topics  KafkaStoreTopics

	config *sarama.Config

	producer sarama.SyncProducer
}

type EventKafkaKey struct {
	Volume string `json:"volume"`
	Server string `json:"server"`
}

func NewKafkaStore(brokers []string, topics KafkaStoreTopics, config *sarama.Config, producer sarama.SyncProducer) *KafkaStore {
	glog.V(3).Infof("Initializing new kafka store with config: \n%+v", config)

	return &KafkaStore{
		brokers:  brokers,
		topics:   topics,
		config:   config,
		producer: producer,
	}
}

func (ks *KafkaStore) Publish(topic string, key []byte, data []byte) (int32, int64, error) {
	defer ks.producer.Close()

	glog.V(3).Infof("Publishing event to topic %s", topic)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := ks.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to send message to kafka producers: %s", err)
	}

	glog.V(3).Infof("kafka message successful: %d, %d", partition, offset)

	return partition, offset, nil
}

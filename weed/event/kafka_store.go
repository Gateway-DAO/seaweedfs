package event

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/seaweedfs/seaweedfs/weed/glog"
)

type KafkaStore struct {
	brokers  []string
	producer sarama.SyncProducer
}

type EventKafkaKey struct {
	Volume string `json:"volume"`
	Server string `json:"server"`
}

func NewKafkaStore(brokers []string, producer sarama.SyncProducer) *KafkaStore {
	return &KafkaStore{
		brokers:  brokers,
		producer: producer,
	}
}

func (ks *KafkaStore) sendKafkaMessage(topic string, key []byte, data []byte) (int32, int64, error) {
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

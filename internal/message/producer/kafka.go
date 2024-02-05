/*
	카프카 메시지 프로듀서

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package producer

import (
	"context"
	"example.com/message/config"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type Message struct {
	Key   []byte
	Value []byte
}

type LogProducer interface {
	Produce(ctx context.Context, ch <-chan *Message)
	Close()
}

type pd struct {
	config   *config.KafkaConfig
	producer *kafka.Producer
}

func NewProducer(cfg *config.KafkaConfig) LogProducer {
	// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	// batch.size	P	1 .. 2147483647
	conf := cfg.KafkaConfigMap()
	producer, err := kafka.NewProducer(&conf)
	if err != nil {
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &pd{cfg, producer}
}

func (rp *pd) Produce(ctx context.Context, ch <-chan *Message) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped producer")
			return
		case msg := <-ch:
			err := rp.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &rp.config.TopicName, Partition: kafka.PartitionAny},
				Key:            msg.Key,
				Value:          msg.Value,
			}, nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (rp *pd) Close() {
	rp.producer.Flush(int(time.Second * 3))
	rp.producer.Close()
}

package mb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"test-exercise/api/constant"
	"test-exercise/api/messagebroker/handlers"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

const (
	kafka_address      = "kafka.address"
	kafka_max_handlers = "kafka.max_handlers"

	fmt_error_delivery_failed = "kafka produce to topic %v - Delivery failed: %v\n%v"
	fmt_warning_ignored_event = "Kafka produce to topic %v - Ignored event: %s"
	fmt_error_listen          = "kafka error in listening topic: %v error: %v\n%v"
)

type kafkaService struct {
	address    string
	producer   *kafka.Producer
	maxHandler int
}

func NewKafkaService() (*kafkaService, error) {
	service := &kafkaService{
		address:    viper.Get(kafka_address).(string),
		maxHandler: viper.Get(kafka_max_handlers).(int),
	}

	var err error
	if service.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": service.address,
	}); err == nil {
		err = service.createTopics()
	}
	return service, err
}

func (service *kafkaService) Produce(topic string, v interface{}) (err error) {
	var value []byte
	switch v.(type) {
	case []byte:
		value = v.([]byte)

	case string:
		value = []byte(v.(string))

	default:
		if value, err = json.Marshal(v); err != nil {
			return
		}
	}

	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range service.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					err = fmt.Errorf(fmt_error_delivery_failed, topic, m.TopicPartition.Error, string(debug.Stack()))
				} //if err
				return

			default:
				log.Printf(fmt_warning_ignored_event, topic, ev)
			}
		}
	}()

	service.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(kafka.PartitionAny),
		},
		Value: value,
	} //Message

	// wait for delivery report goroutine to finish
	_ = <-doneChan
	return
}

func (service *kafkaService) ListenAndServe(topic, group string, handler handlers.Handler) {
	defer func() {
		if e := recover(); e != nil {
			err := fmt.Errorf("%v", e)
			log.Printf(fmt_error_listen, topic, err, string(debug.Stack()))
			service.ListenAndServe(topic, group, handler)
		} //recover
	}()

	var consumer *kafka.Consumer
	var err error
	if consumer, err = service.startConsumer(topic, group); err != nil {
		panic(err)
	}

	var msgReceiver *kafka.Message
	msgChan := make(chan *kafka.Message, 0)
	var handlerCounter int
	var ok bool
	completeChan := make(chan bool, service.maxHandler)
	var kafkaErr kafka.Error
	restartConsumerTimer := time.NewTimer(time.Minute * 5)

	for {
		select {
		case _ = <-restartConsumerTimer.C:
			if consumer != nil {
				consumer.Close()
			}
			if consumer, err = service.startConsumer(topic, group); err != nil {
				log.Printf(fmt_error_listen, topic, err, string(debug.Stack()))
				consumer = nil
			} //startConsumer
			restartConsumerTimer.Reset(time.Minute * 5)

		case _ = <-completeChan:
			handlerCounter--

		default:
			if handlerCounter < service.maxHandler && consumer != nil {
				if msgReceiver, err = consumer.ReadMessage(time.Second / 5); err != nil {
					if kafkaErr, ok = err.(kafka.Error); !ok || kafkaErr.Code() != kafka.ErrTimedOut {
						log.Printf(fmt_error_listen, topic, kafkaErr, string(debug.Stack()))
					}
					continue
				}
				handlerCounter++

				go func() {
					defer func() {
						completeChan <- true
					}()
					msg := <-msgChan
					if err = handler(msg.Value); err != nil {
						log.Printf(fmt_error_listen, topic, err, string(debug.Stack()))
						return
					}
				}()

				msgChan <- msgReceiver
			}
		} //select
	} //for
}

func (service *kafkaService) createTopics() error {
	var err error
	var aClient *kafka.AdminClient
	if aClient, err = kafka.NewAdminClientFromProducer(service.producer); err != nil {
		return fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	defer aClient.Close()

	var md *kafka.Metadata
	if md, err = aClient.GetMetadata(nil, true, 10000); err != nil {
		return fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}

	topics := make([]kafka.TopicSpecification, 0, 1)
	topic := viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT)
	if _, ok := md.Topics[topic]; !ok {
		topics = append(topics, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 2,
			Config:            make(map[string]string, 0),
		})
	}

	if len(topics) > 0 {
		if _, err = aClient.CreateTopics(context.Background(), topics); err != nil {
			return fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
		} //createTopics
	} //if has topics
	return nil
}

func (service *kafkaService) startConsumer(reqTopic, group string) (consumer *kafka.Consumer, err error) {
	if consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": service.address,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	}); err != nil {
		return
	}

	consumer.Subscribe(reqTopic, nil)
	return
}

package main

import (
	"test-exercise/api/constant"
	"test-exercise/api/messagebroker"
	mb_kafka "test-exercise/api/messagebroker/kafka"
	mb_handlers "test-exercise/api/messagebroker/kafka/handlers"
	"test-exercise/api/repository"
	"test-exercise/api/rest"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	messagebroker.MBroker, err = mb_kafka.NewKafkaService()
	if err != nil {
		panic(err)
	}
	go messagebroker.MBroker.ListenAndServe(viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT), viper.GetString(constant.KAFKA_GROUP), mb_handlers.AddEvent)

	if repository.Repo, err = repository.NewPostgresRepository(); err != nil {
		panic(err)
	}

	rest.ListenAndServe(viper.GetInt("port"))
}

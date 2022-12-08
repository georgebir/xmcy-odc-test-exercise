package main

import (
	"test-exercise/api/constant"
	"test-exercise/api/mb"
	mb_handlers "test-exercise/api/mb/handlers"
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

	mb.Kafka, err = mb.NewKafkaService()
	if err != nil {
		panic(err)
	}
	go mb.Kafka.ListenAndServe(viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT), viper.GetString(constant.KAFKA_GROUP), mb_handlers.AddEvent)

	if repository.Repo, err = repository.NewPostgresRepository(); err != nil {
		panic(err)
	}

	rest.ListenAndServe(5010)
}

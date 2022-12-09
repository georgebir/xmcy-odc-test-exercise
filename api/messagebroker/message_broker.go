package messagebroker

import (
	"test-exercise/api/messagebroker/handlers"
)

var MBroker MessageBroker

type MessageBroker interface {
	Produce(topic string, value interface{}) error
	ListenAndServe(topic, group string, handler handlers.Handler)
}

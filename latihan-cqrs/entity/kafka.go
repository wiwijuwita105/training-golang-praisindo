package entity

const (
	KafkaBrokerAddress = "localhost:9092"
	KafkaProducerPort  = ":8080"
	KafkaConsumerPort  = ":8081"
	DefaultTopic       = "users"
	DefaultGroupID     = "consumer-user"
	WriteGroupID       = "consumer-user-write"
	ReadGroupID        = "consumer-user-read"
)

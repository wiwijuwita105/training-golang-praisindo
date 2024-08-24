package config

const (
	KafkaBrokerAddress = "localhost:9092"
	KafkaProducerPort  = ":8080"
	KafkaConsumerPort  = ":8081"
	UserTopic          = "users"
	WriteGroupID       = "consumer-user-write"
	ReadGroupID        = "consumer-user-read"
	DBWriteConn        = "postgresql://postgres:postgres@localhost:5434/db_cqrs_write"
	DBReadConn         = "postgresql://postgres:postgres@localhost:5433/db_cqrs_read"
)

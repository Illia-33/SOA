package service

type Config struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	KafkaHost  string
	KafkaPort  int
}

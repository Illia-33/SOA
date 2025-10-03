package kafkajobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

const kafka_test_topic = "test-topic"

type kafkaTestSuite struct {
	suite.Suite

	cont     *testcontainers.DockerContainer
	hostname string
	port     int
}

func (s *kafkaTestSuite) SetupSuite() {
	ctx := context.Background()

	cont, err := testcontainers.Run(
		ctx,
		"apache/kafka:4.1.0",
		testcontainers.WithEnv(map[string]string{
			"KAFKA_NODE_ID":                                  "1",
			"KAFKA_PROCESS_ROLES":                            "broker,controller",
			"KAFKA_LISTENERS":                                "PLAINTEXT://:29092,CONTROLLER://:9093",
			"KAFKA_ADVERTISED_LISTENERS":                     "PLAINTEXT://localhost:29092",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_CONTROLLER_QUORUM_VOTERS":                 "1@localhost:9093",
			"KAFKA_CONTROLLER_LISTENER_NAMES":                "CONTROLLER",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
			"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
			"KAFKA_AUTO_CREATE_TOPICS_ENABLE":                "false",
		}),
		testcontainers.WithExposedPorts("29092/tcp"),
		testcontainers.WithHostConfigModifier(
			func(hostConfig *container.HostConfig) {
				hostConfig.PortBindings = nat.PortMap{
					"29092/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "29092"}},
				}
			},
		),
	)
	s.Require().NoError(err, "cannot run kafka container")

	time.Sleep(5 * time.Second)

	s.cont = cont
	s.hostname = "localhost"
	s.port = 29092
}

func (s *kafkaTestSuite) kafkaWriter() kafka.Writer {
	return kafka.Writer{
		Addr:         kafka.TCP(fmt.Sprintf("%s:%d", s.hostname, s.port)),
		RequiredAcks: kafka.RequireAll,
		Topic:        kafka_test_topic,
		// Logger:       log.Default(),
		// ErrorLogger:  log.Default(),
	}
}

func (s *kafkaTestSuite) BeforeTest(suiteName, testName string) {
	status, _, err := s.cont.Exec(context.Background(),
		[]string{
			"/opt/kafka/bin/kafka-topics.sh",
			"--bootstrap-server", "localhost:29092",
			"--create",
			"--topic", kafka_test_topic,
			"--replication-factor", "1",
			"--partitions", "1",
		},
	)
	s.Require().Equal(0, status, "cannot delete test topic")
	s.Require().NoError(err, "cannot delete test topic")
}

func (s *kafkaTestSuite) AfterTest(suiteName, testName string) {
	status, _, err := s.cont.Exec(context.Background(),
		[]string{
			"/opt/kafka/bin/kafka-topics.sh",
			"--bootstrap-server", "localhost:29092",
			"--delete",
			"--topic", kafka_test_topic,
		},
	)
	s.Require().Equal(0, status, "cannot delete test topic")
	s.Require().NoError(err, "cannot create test topic")
}

func TestKafkaTopicWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(kafkaTestSuite))
}

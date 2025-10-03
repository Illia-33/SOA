package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

func TestKafkaConsumerSuite(t *testing.T) {
	suite.Run(t, new(kafkaConsumerTestSuite))
}

type kafkaConsumerTestSuite struct {
	suite.Suite

	cont     *testcontainers.DockerContainer
	hostname string
	port     int
}

func (s *kafkaConsumerTestSuite) SetupSuite() {
	ctx := context.Background()

	cont, err := testcontainers.Run(
		ctx,
		"apache/kafka:4.1.0",
		testcontainers.WithEnv(map[string]string{
			"KAFKA_NODE_ID":                                  "1",
			"KAFKA_PROCESS_ROLES":                            "broker,controller",
			"KAFKA_LISTENERS":                                "PLAINTEXT://:29093,CONTROLLER://:9093",
			"KAFKA_ADVERTISED_LISTENERS":                     "PLAINTEXT://localhost:29093",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_CONTROLLER_QUORUM_VOTERS":                 "1@localhost:9093",
			"KAFKA_CONTROLLER_LISTENER_NAMES":                "CONTROLLER",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
			"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
			"KAFKA_AUTO_CREATE_TOPICS_ENABLE":                "false",
		}),
		testcontainers.WithExposedPorts("29093/tcp"),
		testcontainers.WithHostConfigModifier(
			func(hostConfig *container.HostConfig) {
				hostConfig.PortBindings = nat.PortMap{
					"29093/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "29093"}},
				}
			},
		),
	)
	s.Require().NoError(err, "cannot run kafka container")

	time.Sleep(5 * time.Second)

	s.cont = cont
	s.hostname = "localhost"
	s.port = 29093
}

func (s *kafkaConsumerTestSuite) TestSimple() {
	const msg_cnt = 10

	type testMsg struct {
		Value int `json:"value"`
	}

	ctx := context.Background()
	{
		statusCode, _, err := s.cont.Exec(ctx, []string{
			"/opt/kafka/bin/kafka-topics.sh",
			"--bootstrap-server", "localhost:29093",
			"--create",
			"--topic", "test-simple",
			"--replication-factor", "1",
			"--partitions", "1",
		})
		s.Require().NoError(err)
		s.Require().Equal(0, statusCode)

		for i := range msg_cnt {
			statusCode, _, err := s.cont.Exec(ctx, []string{
				"sh", "-c",
				fmt.Sprintf(`echo '{"value": %d}' | /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:29093 --topic test-simple`, i),
			})

			s.Require().NoError(err)
			s.Require().Equal(0, statusCode)
		}
	}

	consumer, err := NewConsumer[testMsg](
		ConnectionConfig{
			Host: s.hostname,
			Port: s.port,
		},
		ConsumerConfig{
			Topic:   "test-simple",
			GroupId: "test-simple-group",
		},
	)
	s.Require().NoError(err)

	for i := range msg_cnt {
		fetchCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		msg, err := consumer.FetchMessage(fetchCtx)
		s.Require().NoError(err)
		s.Assert().Equal(i, msg.Value.Value)

		err = consumer.CommitMessages(ctx, msg)
		s.Require().NoError(err)
	}
}

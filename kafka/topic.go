package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// CreateCompactTopic creates a topic that is used as state store
func CreateCompactTopic(topic string, numPartitions int, replicationFactor int) error {
    broker := os.Getenv("BROKER") // kafka:29092

    adminClient, err := kafka.NewAdminClient(
        &kafka.ConfigMap{
            "bootstrap.servers": broker,
        },
    )
    if err != nil {
        return fmt.Errorf("cannot create admin client from producer [%#v]", err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    result, err := adminClient.CreateTopics(
        ctx,
        []kafka.TopicSpecification{
            kafka.TopicSpecification{
                Topic: topic,
                Config: map[string]string{
                    "cleanup.policy": "compact",
                },
                ReplicationFactor: replicationFactor,
                NumPartitions:     numPartitions,
            },
        })
    if err != nil && result != nil {
        if result[0].Error.Code() == kafka.ErrTopicAlreadyExists {
            fmt.Println("Topic " + topic + " is already existed !")
            return nil
        }
    }

    fmt.Println("Topic " + topic + " is created successfully !")

    return err
}

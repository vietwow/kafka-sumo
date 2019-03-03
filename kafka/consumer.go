package kafka

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    // "io/ioutil"
    // "github.com/vietwow/kafka-sumo/logging"
    "github.com/vietwow/kafka-sumo/sumologic"
)

type MessageConsumer struct {
    // Topic          string
    // ClientID       string
    Consumer       *kafka.Consumer
    FailedCount    int64
    IgnoredCount   int64
    DeliveredCount int64
}

func newMessageConsumer(topic string, broker string, group string) (*MessageConsumer, error) {
    // Initialize kafka consumer
    fmt.Printf("Creating consumer to broker %v with group %v\n", broker, group)

    c := MessageConsumer{}
    // var c *kafka.Consumer
    var err error
    c.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "enable.auto.commit": false,
        "auto.offset.reset":  "earliest"})

    if err != nil {
        return nil, fmt.Errorf("cannot create kafka consumer error [%#v]", err)
    }

    fmt.Printf("=> Created Consumer %v\n", c)

    err = c.Consumer.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        return nil, fmt.Errorf("cannot subcribe to topic [%s] error [%#v]", topic, err)
    } else {
        fmt.Println("=> Subscribed to topic :", topic)
    }

    return &c, nil
}

func (c *MessageConsumer) ProcessMessage(sClient *sumologic.SumoLogic) error {
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)

    // consumer part
    run := true

    for run == true {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev := c.Consumer.Poll(100)
            if ev == nil {
                continue
            }

            switch e := ev.(type) {
            case *kafka.Message:
                fmt.Printf("%% Message on %s:\n%s\n",
                    e.TopicPartition, string(e.Value))

                if e.Headers != nil {
                    fmt.Printf("%% Headers: %v\n", e.Headers)
                }

                // Parse the received msg
                ProcessEvents(e.Value)

                // Sent to SumoLogic
                go sClient.SendLogs(e.Value)

                // commit offset (using when setting "enable.auto.commit" is false - https://github.com/agis/confluent-kafka-go-GH64/blob/master/main.go)
                tp := kafka.TopicPartition{
                    Topic:     e.TopicPartition.Topic,
                    Partition: 0,
                    Offset:    e.TopicPartition.Offset + 1,
                }

                _, err := c.Consumer.CommitOffsets([]kafka.TopicPartition{tp})
                if err != nil {
                    fmt.Print(err)
                }
            case kafka.Error:
                // Errors should generally be considered as informational, the client will try to automatically recover
                fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
            default:
                fmt.Printf("Ignored %v\n", e)
            }
        }
    }

    return nil
}

//Close close the consumer
func (c *MessageConsumer) Close() {
    c.Close()
}

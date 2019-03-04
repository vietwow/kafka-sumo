package confluent

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    // "io/ioutil"
    // "github.com/vietwow/kafka-sumo/logging"
    // "github.com/vietwow/kafka-sumo/sumologic"
)

type MessageConsumerOption func(*messageConsumerOptions) error

var MessageConsumerOptions messageConsumerOptions

type messageConsumerOptions struct {
    topic  string
    broker string
    group  string
}

func(messageConsumerOptions) Topic(topic string) MessageConsumerOption {
    return func(o *messageConsumerOptions) error {
        o.topic = topic
        return nil
    }
}

func(messageConsumerOptions) Broker(broker string) MessageConsumerOption {
    return func(o *messageConsumerOptions) error {
        o.broker = broker
        return nil
    }
}

func(messageConsumerOptions) Group(group string) MessageConsumerOption {
    return func(o *messageConsumerOptions) error {
        o.group = group
        return nil
    }
}



type MessageConsumer struct {
    // Topic          string
    // ClientID       string
    Consumer       *kafka.Consumer
    FailedCount    int64
    IgnoredCount   int64
    DeliveredCount int64
}

func NewConsumer(opts ...MessageConsumerOption) (*MessageConsumer, error) {
    return newMessageConsumer(opts...)
}

func newMessageConsumer(opts ...MessageConsumerOption) (*MessageConsumer, error) {
    args := messageConsumerOptions{}
    for _, opt := range opts {
        if err := opt(&args); err != nil {
            return nil, err
        }
    }


    // Initialize kafka consumer
    fmt.Printf("Creating consumer to broker %v with group %v\n", args.broker, args.group)

    c := MessageConsumer{}
    // var c *kafka.Consumer
    var err error
    c.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  args.broker,
        "group.id":           args.group,
        // "enable.auto.commit": false,
        "auto.offset.reset":  "earliest"})

    if err != nil {
        return nil, fmt.Errorf("cannot create kafka consumer error [%#v]", err)
    }

    fmt.Printf("=> Created Consumer %v\n", c)

    err = c.Consumer.SubscribeTopics([]string{args.topic}, nil)
    if err != nil {
        return nil, fmt.Errorf("cannot subcribe to topic [%s] error [%#v]", args.topic, err)
    } else {
        fmt.Println("=> Subscribed to topic :", args.topic)
    }

    return &c, nil
}

func (c *MessageConsumer) ProcessMessage(h func(msg *kafka.Message) error) error {
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

                if err := h(e); err != nil {
                    return err
                }
                // Sent to SumoLogic
                // go sClient.SendLogs(e.Value)

                // commit offset (using when setting "enable.auto.commit" is false - https://github.com/agis/confluent-kafka-go-GH64/blob/master/main.go)
                // tp := kafka.TopicPartition{
                //     Topic:     e.TopicPartition.Topic,
                //     Partition: 0,
                //     Offset:    e.TopicPartition.Offset + 1,
                // }

                // _, err := c.Consumer.CommitOffsets([]kafka.TopicPartition{tp})
                // if err != nil {
                //     fmt.Print(err)
                // }
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

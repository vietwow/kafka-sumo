package kafka

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    // "io/ioutil"
    "github.com/vietwow/kafka-sumo/logging"
    "github.com/vietwow/kafka-sumo/sumologic"
)

func newMessageConsumer(topic string, broker string, group string) (*kafka.Consumer, error) {
    // Initialize kafka consumer
    fmt.Printf("Creating consumer to broker %v with group %v\n", broker, group)

    var c *kafka.Consumer
    var err error
    c, err  = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "enable.auto.commit": false,
        "auto.offset.reset":  "earliest"})

    if err != nil {
        return nil, fmt.Printf("cannot create kafka consumer error [%#v]", err)
    }

    fmt.Printf("=> Created Consumer %v\n", c)

    err = c.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        return nil, fmt.Printf("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
    } else {
        fmt.Println("=> Subscribed to topic :", topic)
    }

    return &c, nil
}

func (c *kafka.Consumer) ProcessMessage(sClient *sumologic) error {
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
            ev := c.Poll(100)
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

                _, err = c.CommitOffsets([]kafka.TopicPartition{tp})
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
}

//Close close the consumer
func (c *kafka.Consumer) Close() {
    c.Consumer.Close()
}

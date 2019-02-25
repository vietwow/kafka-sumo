package main

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "io/ioutil"
    "time"
    "strings"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/vietwow/kafka-sumo/sumologic"
    "github.com/vietwow/kafka-sumo/logging"
)

var c *kafka.Consumer

var (
    sURL             = kingpin.Flag("sumologic.url", "SumoLogic Collector URL as give by SumoLogic").Required().Envar("SUMOLOGIC_URL").String()
    sSourceCategory  = kingpin.Flag("sumologic.source.category", "Override default Source Category").Envar("SUMOLOGIC_CAT").Default("").String()
    sSourceName      = kingpin.Flag("sumologic.source.name", "Override default Source Name").Default("").String()
    sSourceHost      = kingpin.Flag("sumologic.source.host", "Override default Source Host").Default("").String()
    version          = "0.0.0"

    logs     = kingpin.Flag("log", logsHelp).Short('l').Default(logging.LevelsInSlice[3], logging.LevelsInSlice[1]).Enums(logging.LevelsInSlice[0:]...)
    logsHelp = fmt.Sprintf("Log levels: -l %s", strings.Join(logging.LevelsInSlice[0:], " -l "))
)

func main() {
    kingpin.Version(version)
    kingpin.Parse()

    //logging init
    logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, logging.SliceToLevels(*logs))

    sClient := sumologic.NewSumoLogic(
        *sURL,
        *sSourceHost,
        *sSourceName,
        *sSourceCategory,
        version,
        2*time.Second)

    topic := os.Getenv("TOPIC") // heroku_logs
    broker := os.Getenv("BROKER") // kafka:29092
    group := os.Getenv("GROUP") // myGroup

    // Initialize kafka consumer
    fmt.Printf("Creating consumer to broker %v with group %v\n", broker, group)

    var err error
    c, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  broker,
        "group.id":           group,
        "auto.offset.reset":  "earliest"})

    if err != nil {
        fmt.Printf("Failed to create consumer: %s\n", err)
        os.Exit(1)
    }

    fmt.Printf("=> Created Consumer %v\n", c)

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)

    err = c.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        fmt.Println("Unable to subscribe to topic " + topic + " due to error - " + err.Error())
        os.Exit(1)
    } else {
        fmt.Println("subscribed to topic :", topic)
    }

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

                // Sent to SumoLogic
                formated := sClient.ProcessEvents(string(e.Value))
                // go sClient.SendLogs(string(e.Value))

                if e.Headers != nil {
                    fmt.Printf("%% Headers: %v\n", e.Headers)
                }
            case kafka.Error:
                // Errors should generally be considered as informational, the client will try to automatically recover
                fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
            default:
                fmt.Printf("Ignored %v\n", e)
            }
        }
    }

    fmt.Printf("Closing consumer\n")
    c.Close()
}

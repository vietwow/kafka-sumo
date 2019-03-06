package main

import (
    "fmt"
    "os"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "io/ioutil"
    "time"
    "strings"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/vietwow/kafka-sumo/sumologic"
    "github.com/vietwow/kafka-sumo/logging"
    "github.com/vietwow/kafka-sumo/kafka"
)

var (
    sURL             = kingpin.Flag("sumologic.url", "SumoLogic Collector URL as give by SumoLogic").Required().Envar("SUMOLOGIC_URL").String()
    sSourceCategory  = kingpin.Flag("sumologic.source.category", "Override default Source Category").Envar("SUMOLOGIC_CAT").Default("").String()
    sSourceName      = kingpin.Flag("sumologic.source.name", "Override default Source Name").Default("").String()
    sSourceHost      = kingpin.Flag("sumologic.source.host", "Override default Source Host").Default("").String()
    version          = "0.0.0"

    verbose  = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

    logs     = kingpin.Flag("log", logsHelp).Short('l').Default(logging.LevelsInSlice[3], logging.LevelsInSlice[1]).Enums(logging.LevelsInSlice[0:]...)
    logsHelp = fmt.Sprintf("Log levels: -l %s", strings.Join(logging.LevelsInSlice[0:], " -l "))
)

func main() {
    kingpin.Version(version)
    kingpin.Parse()

    //logging init
    logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, logging.SliceToLevels(*logs))

    // Init sumologic client
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

    // create topic
    confluent.CreateCompactTopic(broker,topic,0,1)

    c, err := confluent.NewConsumer(
        confluent.MessageConsumerOptions.Topic(topic),
        confluent.MessageConsumerOptions.Broker(broker),
        confluent.MessageConsumerOptions.Group(group),
    )
    if err != nil {
        logging.Error.Printf(err.Error())
    }
    c.ProcessMessage(func(msg *kafka.Message) error {
        sClient.SendLogs(msg.Value)
        return nil
    })
    c.Close()
}
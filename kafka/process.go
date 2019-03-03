package kafka

import (
    "fmt"
    "time"
    "encoding/json"
)
/*

{"syslog.pri":"190","syslog.timestamp":"2019-02-25T16:09:23.296535+00:00","syslog.hostname":"host","syslog.appname":"app","syslog.procid":"sidekiq_pusher.1","message":"[PROJECT_ROOT]/app/models/trading_account.rb:95 in `pnl'","syslog.facility":"local7","syslog.severity":"info","logplex.drain_id":"d.133fd389-68d7-4d1e-a40d-521da502593a","logplex.frame_id":"EBC908AFB0D5433EFE95C18457C0B5F8"}

*/

type log struct {
    Timestamp       int64     `json:"Timestamp"`
    SyslogPri       string    `json:"Syslog.pri"`
    SyslogTimestamp time.Time `json:"Syslog.timestamp"`
    SyslogHostname  string    `json:"Syslog.hostname"`
    SyslogAppname   string    `json:"Syslog.appname"`
    SyslogProcid    string    `json:"Syslog.procid"`
    Message         string    `json:"Message"`
    SyslogFacility  string    `json:"Syslog.facility"`
    SyslogSeverity  string    `json:"Syslog.severity"`
    LogplexDrainID  string    `json:"Logplex.drain_id"`
    LogplexFrameID  string    `json:"Logplex.frame_id"`
}


//ProcessEvents
//Format SlowLog Interface to flat string
func ProcessEvents(msg []byte) {
    // https://github.com/lightstaff/confluent-kafka-go-example/blob/master/main.go
    var log Log
    if err := json.Unmarshal(msg, &log); err != nil {
        fmt.Println(err)
    }
    fmt.Printf("success parse consumed log : message: %s, timestamp: %d, Logplex.DrainID = %v\n", log.Message, log.Timestamp, log.LogplexDrainID)

    // Get byte slice from string.
    // bytes := []byte(msg)

    // Unmarshal string into structs.
    // var log []Log
    // json.Unmarshal(bytes, &log)

    // Loop over structs and display them.
    // for l := range log {
    //     fmt.Printf("Logplex.DrainID = %v, Message = %v", log[l].LogplexDrainID, log[l].Message)
    //     fmt.Println()
    // }
}


package sumologic

import (
	// "bytes"
	// "compress/gzip"
	"encoding/json"
	"fmt"
	"errors"
	"net/http"
	"time"

	"github.com/vietwow/kafka-sumo/logging"
)

type SumoLogicEvents struct {
}

type SumoLogic struct {
	httpClient           http.Client
	forwarderVersion     string
	sumoCategory         string
	sumoName             string
	sumoHost             string
	sumoURL              string
	sumoPostMinimumDelay time.Duration
	timerBetweenPost     time.Time
	customMetadata       string
}

func NewSumoLogic(url string, host string, name string, category string, expVersion string, connectionTimeOut time.Duration) *SumoLogic {
	return &SumoLogic{
		forwarderVersion: expVersion,
		sumoHost:         host,
		sumoName:         name,
		sumoURL:          url,
		sumoCategory:     category,
		httpClient:       http.Client{Timeout: connectionTimeOut},
	}
}

/*

{"syslog.pri":"190","syslog.timestamp":"2019-02-25T16:09:23.296535+00:00","syslog.hostname":"host","syslog.appname":"app","syslog.procid":"sidekiq_pusher.1","message":"[PROJECT_ROOT]/app/models/trading_account.rb:95 in `pnl'","syslog.facility":"local7","syslog.severity":"info","logplex.drain_id":"d.133fd389-68d7-4d1e-a40d-521da502593a","logplex.frame_id":"EBC908AFB0D5433EFE95C18457C0B5F8"}

*/

type Log struct {
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
func (s *SumoLogic) ProcessEvents(msg []byte) {
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

func (s *SumoLogic) SendLogs(logStringToSend []byte) {
	logging.Trace.Println("Attempting to send to Sumo Endpoint: " + s.sumoURL)
	if len(logStringToSend) != 0 {
		request, err := http.NewRequest("POST", s.sumoURL, string(logStringToSend))
		if err != nil {
			logging.Error.Printf("http.NewRequest() error: %v\n", err)
			return
		}
		// request.Header.Set("Content-Type", "application/json")
		request.Header.Add("Content-Encoding", "gzip")
		request.Header.Add("X-Sumo-Client", "redis-forwarder v"+s.forwarderVersion)

		if s.sumoName != "" {
			request.Header.Add("X-Sumo-Name", s.sumoName)
		}
		if s.sumoHost != "" {
			request.Header.Add("X-Sumo-Host", s.sumoHost)
		}
		if s.sumoCategory != "" {
			request.Header.Add("X-Sumo-Category", s.sumoCategory)
		}
		//checking the timer before first POST intent
		for time.Since(s.timerBetweenPost) < s.sumoPostMinimumDelay {
			logging.Trace.Println("Delaying Post because minimum post timer not expired")
			time.Sleep(100 * time.Millisecond)
		}
		response, err := s.httpClient.Do(request)

		if (err != nil) || (response.StatusCode != 200 && response.StatusCode != 302 && response.StatusCode < 500) {
			logging.Info.Println("Endpoint dropped the post send")
			logging.Info.Printf("response.StatusCode is %v and err is %v\n", response.StatusCode, err)
			logging.Info.Println("Waiting for 300 ms to retry")
			time.Sleep(300 * time.Millisecond)
			statusCode := 0
			err := Retry(func(attempt int) (bool, error) {
				var errRetry error
				request, err := http.NewRequest("POST", s.sumoURL, string(logStringToSend))
				if err != nil {
					logging.Error.Printf("http.NewRequest() error: %v\n", err)
				}
				// request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Content-Encoding", "gzip")
				request.Header.Add("X-Sumo-Client", "redis-forwarder v"+s.forwarderVersion)

				if s.sumoName != "" {
					request.Header.Add("X-Sumo-Name", s.sumoName)
				}
				if s.sumoHost != "" {
					request.Header.Add("X-Sumo-Host", s.sumoHost)
				}
				if s.sumoCategory != "" {
					request.Header.Add("X-Sumo-Category", s.sumoCategory)
				}
				//checking the timer before POST (retry intent)
				for time.Since(s.timerBetweenPost) < s.sumoPostMinimumDelay {
					logging.Trace.Println("Delaying Post because minimum post timer not expired")
					time.Sleep(100 * time.Millisecond)
				}
				response, errRetry = s.httpClient.Do(request)

				if errRetry != nil {
					logging.Error.Printf("http.Do() error: %v\n", errRetry)
					logging.Info.Println("Waiting for 300 ms to retry after error")
					time.Sleep(300 * time.Millisecond)
					return attempt < 5, errRetry
				} else if response.StatusCode != 200 && response.StatusCode != 302 && response.StatusCode < 500 {
					logging.Info.Println("Endpoint dropped the post send again")
					logging.Info.Println("Waiting for 300 ms to retry after a retry ...")
					statusCode = response.StatusCode
					time.Sleep(300 * time.Millisecond)
					return attempt < 5, errRetry
				} else if response.StatusCode == 200 {
					logging.Trace.Println("Post of logs successful after retry...")
					s.timerBetweenPost = time.Now()
					statusCode = response.StatusCode
					return true, err
				}
				return attempt < 5, errRetry
			})
			if err != nil {
				logging.Error.Println("Error, Not able to post after retry")
				logging.Error.Printf("http.Do() error: %v\n", err)
				return
			} else if statusCode != 200 {
				logging.Error.Printf("Not able to post after retry, with status code: %d", statusCode)
			}
		} else if response.StatusCode == 200 {
			logging.Trace.Println("Post of logs successful")
			logging.Info.Println("Post of logs successful")
			s.timerBetweenPost = time.Now()
		}

		if response != nil {
			defer response.Body.Close()
		}
	}
}

//------------------Retry Logic Code-------------------------------

// MaxRetries is the maximum number of retries before bailing.
var MaxRetries = 10
var errMaxRetriesReached = errors.New("exceeded retry limit")

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Retry(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func IsMaxRetries(err error) bool {
	return err == errMaxRetriesReached
}
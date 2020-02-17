package logaggregator

import (
	json2 "encoding/json"
	"fmt"
	"lb/common/entity"
	guards "lb/common/guard"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"lb/common/log"
)

const (
	logpath       = "./log.txt"
	appserverkey  = "/etc/lb/certs/key.pem"
	appservercert = "/etc/lb/certs/cert.pem"
)

var entries logEntries

var mutex sync.Mutex

var logger = log.Log

//in production minute, depending on use-case
var tickCh = time.Tick(5 * time.Second)
var writeDelay = time.Second * 2

func Run(port string) {

	http.HandleFunc("/", storeentry)

	f, err := os.OpenFile(logpath, os.O_RDONLY|os.O_CREATE, 0666)
	guards.FailOnError(err, "opening file %s", logpath)
	err = f.Close()
	guards.FailOnError(err, "opening file %s", logpath)

	go http.ListenAndServeTLS(port, appservercert, appserverkey, nil)

	go writeLog()

	logger.Info().Msgf("Log aggregator listening on port %s. Press <ENTER> to exit", port)

	fmt.Scanln()
}

func storeentry(writer http.ResponseWriter, request *http.Request) {
	var entry entity.LogEntry
	err := json2.NewDecoder(request.Body).Decode(&entry)
	if err != nil {
		logger.Err(err).Msg("parsing json from body")
		writer.WriteHeader(http.StatusInternalServerError)
	}

	mutex.Lock()
	entries = append(entries, entry)
	mutex.Unlock()
}

func writeLog() {
	for range tickCh {
		mutex.Lock()

		logFile, err := os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			logger.Err(err).Msgf(logpath)
			continue
		}

		targetTime := time.Now().Add(-writeDelay)

		sort.Sort(entries)
		for i, entry := range entries {
			if entry.Timestamp.Before(targetTime) {
				_, err := logFile.WriteString(formatEntry(entry))
				if err != nil {
					logger.Err(err).Msgf("writing %s to file", entry)
				}
				if i == len(entries)-1 {
					entries = logEntries{}
				}
			} else {
				entries = entries[:i]
				break
			}
		}

		if err := logFile.Close(); err != nil {
			logger.Err(err).Msgf("unable to close file", logpath)
		}

		mutex.Unlock()
	}
}

func formatEntry(entry entity.LogEntry) string {
	return fmt.Sprintf("%v;%v;%v;%v\n", entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Level, entry.Source, entry.Destination)
}

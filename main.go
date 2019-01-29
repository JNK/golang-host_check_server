package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func (ce *Config) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMessages []*json.RawMessage
	err = json.Unmarshal(*objMap["checks"], &rawMessages)
	if err != nil {
		return err
	}

	ce.Checks = make([]ICheck, len(rawMessages))

	var m map[string]interface{}
	for index, rawMessage := range rawMessages {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		if m["type"] == "min-process-count" {
			var p MinProcessCount
			err := json.Unmarshal(*rawMessage, &p)
			checkError(err)
			ce.Checks[index] = &p
		} else if m["type"] == "cmd-exit-code" {
			var p ShellExitCode
			err := json.Unmarshal(*rawMessage, &p)
			checkError(err)
			ce.Checks[index] = &p
		} else if m["type"] == "cmd-result-regex" {
			var p ShellResultRegex
			err := json.Unmarshal(*rawMessage, &p)
			checkError(err)
			ce.Checks[index] = &p
		} else {
			return errors.New("Unsupported type found!")
		}
	}
	return nil
}

var config = Config{}
var healthy = false

func handler(w http.ResponseWriter, r *http.Request) {
	if healthy {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("NOT OK\n"))
	}
}

var defaultCheckInterval int = 5
var defaultConfigFile string = "/etc/check_health.json"
var defaultHandlerPath string = "/health"
var defaultServerExpose string = "0.0.0.0:2424"

func main() {
	var checkInterval int
	var configFilePath string
	var handlerPath string
	var hostDefinition string
	flag.IntVar(&checkInterval, "i", 5, "number of seconds to run checks")
	flag.StringVar(&configFilePath, "c", defaultConfigFile, "path to health check server config file")
	flag.StringVar(&handlerPath, "l", defaultHandlerPath, "path to serve results at")
	flag.StringVar(&hostDefinition, "e", defaultServerExpose, "how to expose server")
	flag.Parse()

	// ARG CHECKS
	if checkInterval <= 0 || checkInterval >= 60 {
		checkError(errors.New("interval must be between 1 and 59 seconds"))
	} else if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		checkError(err)
	} else if !strings.HasPrefix(handlerPath, "/") {
		checkError(errors.New("handler path must start with /"))
	}

	// LOAD CONFIG FILE
	data, err := ioutil.ReadFile(configFilePath)
	checkError(err)

	err = json.Unmarshal(data, &config)
	checkError(err)

	c := cron.New()
	err = c.AddFunc("0/5 * * * * *", func() {
		healthy = CheckHealth(config)
	})
	checkError(err)
	c.Start()

	healthy = CheckHealth(config)

	http.HandleFunc(handlerPath, handler)
	log.Print("host health server starting at: http://" + hostDefinition + handlerPath)
	log.Fatal(http.ListenAndServe(hostDefinition, nil))
}

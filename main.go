package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

		if m["type"] == "process-count-gte" {
			var p ProcessCountGte
			err := json.Unmarshal(*rawMessage, &p)
			checkError(err)
			ce.Checks[index] = &p
		} else if m["type"] == "process-count-lte" {
			var p ProcessCountLte
			err := json.Unmarshal(*rawMessage, &p)
			checkError(err)
			ce.Checks[index] = &p
		} else if m["type"] == "process-count-eq" {
			var p ProcessCountEq
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

func handler(writer http.ResponseWriter, request *http.Request) {
	jsonResult := request.URL.Query()["json"] != nil

	if healthy {
		writer.WriteHeader(200)
		if jsonResult {
			_, _ = writer.Write([]byte("OK\n"))
		} else {
			_, _ = writer.Write([]byte("OK\n"))
		}
	} else {
		writer.WriteHeader(500)
		if jsonResult {
			_, _ = writer.Write([]byte("NOT OK\n"))
		} else {
			_, _ = writer.Write([]byte("NOT OK\n"))
		}
	}
}

var defaultCheckInterval = 5
var defaultConfigFile = "/etc/check_health.json"
var defaultHandlerPath = "/health"
var defaultServerExpose = "0.0.0.0:2424"

func main() {
	var checkInterval int
	var configFilePath string
	var handlerPath string
	var hostDefinition string
	flag.IntVar(&checkInterval, "i", defaultCheckInterval, "number of seconds to run checks")
	flag.StringVar(&configFilePath, "c", defaultConfigFile, "path to health check server config file")
	flag.StringVar(&handlerPath, "l", defaultHandlerPath, "root path to serve server at")
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
	cronString := fmt.Sprintf("0/%v * * * * *", checkInterval)
	err = c.AddFunc(cronString, func() {
		healthy = CheckHealth(config)
	})
	checkError(err)
	c.Start()

	healthy = CheckHealth(config)

	http.HandleFunc(handlerPath, handler)
	log.Print("host health server starting at: http://" + hostDefinition + handlerPath)
	log.Fatal(http.ListenAndServe(hostDefinition, nil))
}

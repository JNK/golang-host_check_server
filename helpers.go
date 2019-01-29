package main

import (
	"time"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckHealth(config Config) bool {
	checkResult := CreateResultSet(config)
	for _, res := range checkResult.Results {
		if !res.Success {
			return false
		}
	}
	return true
}

func CreateResultSet(config Config) ResultSet {
	var results = make([]CheckResult, len(config.Checks))

	for i, check := range config.Checks {
		result, message := check.Validate()
		results[i] = CheckResult{Time: time.Now(), Message: message, Name: check.Name(), Success: result}
	}

	return ResultSet{Results: results}
}

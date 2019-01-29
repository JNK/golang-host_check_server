package main

import "time"

type Config struct {
	Checks []ICheck `json:"checks"`
}

type ICheck interface {
	Validate() (bool, string)
	Name() string
}

type CheckResult struct {
	Name    string
	Success bool
	Time    time.Time
	Message string
}

type ResultSet struct {
	Results []CheckResult
}

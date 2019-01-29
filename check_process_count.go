package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
)

func CountForProcess(name string) int {
	psx, err := ps.Processes()
	checkError(err)

	var count = 0
	for _, process := range psx {
		if process.Executable() == name {
			count += 1
		}
	}
	return count
}

// ProcessCountGte
type ProcessCountGte struct {
	JobName string `json:"name"`
	Process string `json:"process"`
	Minimum int    `json:"min"`
}

func (a *ProcessCountGte) Validate() (bool, string) {
	actualCount := CountForProcess(a.Process)
	result := actualCount >= a.Minimum

	if result {
		return true, "success"
	} else {
		return false, fmt.Sprintf("minimum count: %v; actual count: %v", a.Minimum, actualCount)
	}
}

func (a *ProcessCountGte) Name() string {
	return a.JobName
}

// ProcessCountLte

type ProcessCountLte struct {
	JobName string `json:"name"`
	Process string `json:"process"`
	Maximum int    `json:"max"`
}

func (a *ProcessCountLte) Validate() (bool, string) {
	actualCount := CountForProcess(a.Process)
	result := actualCount <= a.Maximum

	if result {
		return true, "success"
	} else {
		return false, fmt.Sprintf("maximum count: %v; actual count: %v", a.Maximum, actualCount)
	}
}

func (a *ProcessCountLte) Name() string {
	return a.JobName
}

// ProcessCountEq

type ProcessCountEq struct {
	JobName string `json:"name"`
	Process string `json:"process"`
	Count   int    `json:"count"`
}

func (a *ProcessCountEq) Validate() (bool, string) {
	actualCount := CountForProcess(a.Process)
	result := actualCount == a.Count

	if result {
		return true, "success"
	} else {
		return false, fmt.Sprintf("expected count: %v; actual count: %v", a.Count, actualCount)
	}
}

func (a *ProcessCountEq) Name() string {
	return a.JobName
}

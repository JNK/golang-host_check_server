package main

import "github.com/mitchellh/go-ps"

type MinProcessCount struct {
	Process string `json:"process"`
	Minimum int `json:"min"`
}

func (a *MinProcessCount) Validate() bool {
	psx, err := ps.Processes()
	checkError(err)

	var count = 0
	for _, ps := range psx {
		if ps.Executable() == a.Process {
			count += 1
		}
	}

	return count >= a.Minimum
}


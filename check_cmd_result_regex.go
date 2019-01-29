package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type ShellResultRegex struct {
	JobName string `json:"name"`
	Command string `json:"command"`
	Regex   string `json:"regex"`
}

func (a *ShellResultRegex) Validate() (bool, string) {
	commands := strings.Split(a.Command, " ")
	name, arguments := commands[0], commands[1:]

	out, _ := exec.Command(name, arguments...).Output()
	re := regexp.MustCompile(a.Regex)

	result := re.MatchString(string(out))

	if result {
		return true, "success"
	} else {
		return false, fmt.Sprintf("pattern \"%v\" did not match \"%v\"", a.Regex, string(out))
	}
}

func (a *ShellResultRegex) Name() string {
	return a.JobName
}

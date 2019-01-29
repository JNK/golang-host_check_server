package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type ShellResultRegex struct {
	Command string `json:"command"`
	Regex string `json:"regex"`
}

func (a *ShellResultRegex) Validate() bool {
	commands := strings.Split(a.Command, " ")
	name, arguments := commands[0], commands[1:]

	out, err := exec.Command(name, arguments...).Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(a.Regex)

	return re.MatchString(string(out))
}
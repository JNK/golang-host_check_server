package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type ShellExitCode struct {
	JobName string `json:"name"`
	Command string `json:"command"`
	Code    int    `json:"code"`
}

func (a *ShellExitCode) Validate() (bool, string) {
	commands := strings.Split(a.Command, " ")
	name, arguments := commands[0], commands[1:]
	cmd := exec.Command(name, arguments...)

	if err := cmd.Start(); err != nil {
		return false, fmt.Sprintf("error with cmd-exit-code \"%v\": %v", a.Command, err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code := status.ExitStatus()
				if code != a.Code {
					return false, fmt.Sprintf("expected code: %v; got code: %s", a.Code, code)
				}
			}
		} else {
			return false, fmt.Sprintf("failed %s", err)
		}
	}

	return true, "success"
}

func (a *ShellExitCode) Name() string {
	return a.JobName
}

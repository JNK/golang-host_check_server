package main

import (
	"log"
	"os/exec"
	"strings"
	"syscall"
)

type ShellExitCode struct {
	Command string `json:"command"`
	Code int `json:"code"`
}

func (a *ShellExitCode) Validate() bool {
	commands := strings.Split(a.Command, " ")
	name, arguments := commands[0], commands[1:]
	cmd := exec.Command(name, arguments...)

	if err := cmd.Start(); err != nil {
		log.Fatalf("error with cmd-exit-code \"%v\": %v", a.Command, err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return  status.ExitStatus() == a.Code
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	return true
}


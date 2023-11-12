package main

import (
	"os/exec"
)

type step struct {
	name    string
	exe     string
	args    []string
	message string
	proj    string
}

// Constructor function for step
func newStep(name, exec, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exec,
		message: message,
		args:    args,
		proj:    proj,
	}
}

// Define the method execute for step
func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{step: s.name, msg: "failed to execute", cause: err}
	}
	return s.message, nil
}

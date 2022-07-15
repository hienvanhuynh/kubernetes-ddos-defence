package main

import (
	"testing"
)

func TestExecBashCommand(t *testing.T) {

	_, error := execBashCommand("ls -lah")

	if error != nil {
        t.Errorf("Simple command failed")
    }

	_, error = execBashCommand("ls -lah | grep \"[a,\\.].*\"")

	if error != nil {
		t.Errorf("Pipe command failed")
	}
}

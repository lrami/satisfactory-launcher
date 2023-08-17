package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

const (
	TASK_LIST = "tasklist.exe"
)

func isProcessRunning(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}

	cmd := exec.Command(TASK_LIST, "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if bytes.Contains(out, []byte(name)) {
			return true, nil
		}
	}
	return false, nil
}

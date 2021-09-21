package integration

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	_, file, _, _ = runtime.Caller(0)
	ModuleRoot    = filepath.Join(filepath.Dir(file), "../..")
	FlowCmd       = "./cmd/flow"
)

func RunFlowCmd(args ...string) ([]byte, error) {
	// integration testing by way of external process execution

	cmd := MakeFlowCmd(args...)
	fmt.Println("Running command: " + cmd.String())
	out, err := cmd.Output()

	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			err = fmt.Errorf("%w with stderr: %v", err, string(exitErr.Stderr))
		}
	}
	return out, err
}

func MakeFlowCmd(args ...string) *exec.Cmd {
	finalArgs := make([]string, 2+len(args))
	finalArgs[0] = "run"
	finalArgs[1] = FlowCmd
	for i, arg := range args {
		finalArgs[i+2] = arg
	}
	cmd := exec.Command("go", finalArgs...)
	cmd.Dir = ModuleRoot

	return cmd
}

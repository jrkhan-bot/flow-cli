package integration

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, file, _, _ = runtime.Caller(0)
	ModuleRoot    = filepath.Join(filepath.Dir(file), "../..")
	FlowCmd       = "./cmd/flow"
)

func RunFlowCmd(t *testing.T, args ...string) ([]byte, error) {
	// integration testing by way of external process execution
	cmd := MakeFlowCmd(t, args...)
	return cmd.Output()
}

func MakeFlowCmd(t *testing.T, args ...string) *exec.Cmd {
	finalArgs := make([]string, 2+len(args))
	finalArgs[0] = "run"
	finalArgs[1] = FlowCmd
	for i, arg := range args {
		finalArgs[i+2] = arg
	}
	cmd := exec.Command("go", finalArgs...)
	cmd.Dir = ModuleRoot
	t.Log(cmd.String())
	return cmd
}

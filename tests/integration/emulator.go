package integration

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

var EmulatorProcessName = "flow"

// EmulatorProcess is a wrapper around the emulator command with convenience functions for integration testing
type EmulatorProcess struct {
	cmd *exec.Cmd
}

// RunEmulator starts up the emulator in a separate process
func RunEmulator(t *testing.T) (*EmulatorProcess, error) {
	cmd := MakeFlowCmd(t, "emulator")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	emu := &EmulatorProcess{cmd}
	emu.JoinOut(t, out)
	return emu, nil
}

func (e *EmulatorProcess) Stop(t *testing.T) error {
	sig := os.Interrupt
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		// windows needs to send SIGKILL as SIGINT is not supported
		sig = os.Kill
	}

	t.Logf("send %v to process %v", sig, e.cmd.Process.Pid)
	err := e.cmd.Process.Signal(sig)
	if err != nil {
		t.Logf("unable to terminate process with error %v", err)
	}
	_, err = e.cmd.Process.Wait()
	if isWindows {
		// windows needs to clean up children
		e.StopChildrenWindows(t)
	}
	return err
}

func (e *EmulatorProcess) JoinOut(t *testing.T, out io.ReadCloser) {
	outbr := bufio.NewReader(out)
	go func() {
		for {
			line, _ := outbr.ReadString('\n')

			if len(line) > 0 {
				t.Log(line)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (e *EmulatorProcess) StopChildrenWindows(t *testing.T) {
	// stop the child emulator process we missed by just killing the parent
	out, err := exec.Command("powershell.exe", "Get-Process", EmulatorProcessName, "|", "Stop-Process").Output()
	t.Log(out)
	if err != nil {
		t.Log(err)
	}
}

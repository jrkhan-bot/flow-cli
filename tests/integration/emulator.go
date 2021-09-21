package integration

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var EmulatorProcessName = "flow"

// EmulatorProcess is a wrapper around the emulator command with convenience functions for integration testing
type EmulatorProcess struct {
	cmd *exec.Cmd
}

// RunEmulator starts up the emulator in a separate process
func RunEmulator() (*EmulatorProcess, error) {
	cmd := MakeFlowCmd("emulator")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	emu := &EmulatorProcess{cmd}
	emu.JoinOut(out)
	return emu, nil
}

func (e *EmulatorProcess) Stop() error {
	sig := os.Interrupt
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		// windows needs to send SIGKILL as SIGINT is not supported
		sig = os.Kill
	}

	fmt.Printf("send %v to process %v", sig, e.cmd.Process.Pid)
	err := e.cmd.Process.Signal(sig)
	if err != nil {
		fmt.Printf("unable to terminate process with error %v\n", err)
	}
	_, err = e.cmd.Process.Wait()
	if isWindows {
		// windows needs to clean up children
		e.StopChildrenWindows()
	}
	return err
}

func (e *EmulatorProcess) JoinOut(out io.ReadCloser) {
	outbr := bufio.NewReader(out)
	go func() {
		for {
			line, _ := outbr.ReadString('\n')

			if len(line) > 0 {
				fmt.Print(line)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (e *EmulatorProcess) StopChildrenWindows() {
	// stop the child emulator process we missed by just killing the parent
	out, err := exec.Command("powershell.exe", "Get-Process", EmulatorProcessName, "|", "Stop-Process").Output()
	fmt.Print(out)
	if err != nil {
		fmt.Print(err)
	}
}

package integration

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/internal/emulator"
)

// RunEmulator starts up the emulator in a separate process
func RunEmulator() error {
	var cmd = &cobra.Command{
		Use:              "flow",
		TraverseChildren: true,
	}
	cmd.AddCommand(emulator.Cmd)
	command.InitFlags(cmd)
	fmt.Println("🌱  Starting emulator")
	stdOut, _, err := ExecuteAsync(cmd, "emulator")

	if err != nil {
		return err
	}
	// block until the gRPC server has actually started
	scanner := bufio.NewScanner(stdOut)
	for scanner.Scan() {
		// for now, just scanning out for an expected message
		if strings.Contains(scanner.Text(), "Starting gRPC server") {
			break
		}
	}
	return nil
}

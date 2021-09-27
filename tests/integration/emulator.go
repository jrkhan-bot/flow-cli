/*
 * Flow CLI
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/internal/emulator"
)

var (
	testPort = ":3570"
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
	_, _, err := ExecuteAsync(cmd, "emulator")
	if err != nil {
		return err
	}

	// block until the gRPC server has actually started on the expected port
	for {
		// does a channel/event exist to let us know when the server has started?
		// just checking to see if the port is in use for now
		if portInUse(testPort) {
			return nil
		}
		time.Sleep(time.Millisecond * 50)
	}
}

// checkPort checks if a port is free
func portInUse(port string) bool {
	_, err := net.Listen("tcp", testPort)
	// if there is an error, we assume the port is in use
	return err != nil
}

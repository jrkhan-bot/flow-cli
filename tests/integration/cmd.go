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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	_, file, _, _ = runtime.Caller(0)
	ConfigPath    = filepath.Join(filepath.Dir(file), "testdata")
	ModuleRoot    = filepath.Join(filepath.Dir(file), "../..")
	FlowCmd       = "./cmd/flow"
	ConfigFile    = "test_config.json"
)

func ExecuteCommand(root *cobra.Command, args ...string) (result []byte, stdOut string, stdErr string, err error) {
	// create a temp file to capture outout
	cmdName := root.Name()
	file, err := tmpFile(cmdName)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	args = append(args, configArgs()...)
	args = append(args, "-s", file.Name())
	command, outBuffer, errBuffer := ConnectCommand(root, args...)
	fmt.Printf("👟 Running command with args: %v\n", args)
	_, err = command.ExecuteC()
	if err != nil {
		return
	}
	out, err := ioutil.ReadAll(file)

	return out, outBuffer.String(), errBuffer.String(), err
}

func configArgs() []string {
	return []string{"-f", filepath.Join(ConfigPath, ConfigFile)}
}

func ConnectCommand(root *cobra.Command, args ...string) (c *cobra.Command, outBuffer, errBuffer *bytes.Buffer) {
	// since we write directly to os.StdOut in command.outputResult setting this does contain the 'response'
	outBuffer = new(bytes.Buffer)
	errBuffer = new(bytes.Buffer)
	root.SetOut(outBuffer)
	root.SetErr(errBuffer)
	root.SetArgs(args)
	return root, outBuffer, errBuffer
}

// ExecuteAsync executes a command and returns immediately, allowing us to examine output buffers
func ExecuteAsync(root *cobra.Command, args ...string) (outBuffer, errBuffer *bytes.Buffer, err error) {
	args = append(args, configArgs()...)
	var cmd *cobra.Command
	cmd, outBuffer, errBuffer = ConnectCommand(root, args...)
	go func() {
		_, err := cmd.ExecuteC()
		if err != nil {
			panic(err)
		}
	}()
	return
}

func tmpFile(command string) (*os.File, error) {
	return ioutil.TempFile("", command)
}

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

package accounts_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-cli/internal/accounts"
	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/tests/integration"
)

var addressRegex = regexp.MustCompile("Address\\s+(?P<address>0x[^\n]+)")

// TestMain starts and stops the emulator after each test
func TestMain(m *testing.M) {
	err := integration.RunEmulator()
	if err != nil {
		panic("unable to start emulator")
	}
	os.Exit(m.Run())
}

// TestCobraCommand integration test executes the Cobra command in the same process as the test
//
// Pros:
//  `Can run tests within the existing test harness`
//  `Coverage reports with minimal configuration changes`
//  `Faster than spawning additional processes`
//  `No need to deal with os specific signals for stopping emulator`
// Cons:
//  `Minor differences from invoking process directly (e.g. -o format seems to be sticky)``
func TestCobraCommand(t *testing.T) {
	var cmd = &cobra.Command{
		Use:              "flow",
		TraverseChildren: true,
	}

	cmd.AddCommand(accounts.Cmd)
	command.InitFlags(cmd)
	format := "json"
	result, _, stdErr, err := integration.ExecuteCommand(cmd, "accounts", "create", "-o", format)
	assert.NoError(t, err)
	assert.Empty(t, stdErr)
	id, err := getAccountIdFromCreateAccountOut(result, format)
	assert.NoError(t, err)
	t.Log(string(result))
	if id == "" {
		t.Logf("unable to resolve address from output %s", string(result))
		t.FailNow()
	}
	t.Logf("📪 new address: %s", id)
	assert.NoErrorf(t, err, "unable to create account: %v", err)

	// now confirm the address was added
	getResult, _, stdErr, err := integration.ExecuteCommand(cmd, "accounts", "get", id, "-o", format)
	assert.Empty(t, stdErr)
	assert.NoError(t, err)

	// expect the response contains Address, Balance, Keys, and Contracts
	res, err := GetResultJson(getResult)
	assert.NoError(t, err)
	assert.Equal(t, res.Address, id)
}

type AccountJsonResult struct {
	Address string `json:"address"`
}

func GetResultJson(createOutput []byte) (*AccountJsonResult, error) {
	var res AccountJsonResult
	err := json.Unmarshal(createOutput, &res)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall %v with error %w", string(createOutput), err)
	}
	return &res, nil
}

func getAccountIdFromCreateAccountOut(createOutput []byte, format string) (string, error) {
	switch format {
	case "json":
		res, err := GetResultJson(createOutput)
		if err != nil {
			return "", err
		}
		return res.Address, nil
	default:
		addressIndex := addressRegex.SubexpIndex("address")

		matches := addressRegex.FindSubmatch(createOutput)
		if addressIndex >= len(matches) {
			return "", fmt.Errorf("address not found in output: %v", string(createOutput))
		}
		return string(matches[addressIndex]), nil
	}
}

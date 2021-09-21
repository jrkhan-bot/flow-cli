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
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-cli/tests/integration"
)

var addressRegex = regexp.MustCompile("Address\\s+(?P<address>0x[^\n]+)")

func init() {
	// bootstrap configuration if it doesn't exist
	out, err := integration.RunFlowCmd("init")
	// print any command output
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	// print any error outside of config already existing
	if err != nil && !strings.Contains(err.Error(), "configuration already exists") {
		fmt.Println(err)
	}
}

// TestMain starts and stops the emulator after each test
func TestMain(m *testing.M) {
	emu, err := integration.RunEmulator()
	code := m.Run()
	if err != nil {
		panic(fmt.Sprintf("unable to start emulator %v", err))
	}
	err = emu.Stop()
	if err != nil {
		fmt.Printf("unable to stop emulator: %v\n", err)
	}
	os.Exit(code)
}

func TestAccountCreateCommand(t *testing.T) {
	individualTimeout := time.After(20 * time.Second)
	done := make(chan struct{})
	go func() {
		out, err := integration.RunFlowCmd("accounts", "create")
		assert.NoErrorf(t, err, "unable to create account: %v", err)

		id, err := getAccountIdFromCreate(out)
		assert.NoError(t, err)
		if id == "" {
			t.Logf("unable to resolve address from output %s", string(out))
			t.Fail()
		}
		t.Logf("address: %s", id)

		// now confirm the address was added
		out, err = integration.RunFlowCmd("accounts", "get", id)
		assert.NoError(t, err)
		result := string(out)
		t.Log(result)
		// expect the response contains Address, Balance, Keys, and Contracts
		assert.Contains(t, result, "Address")
		assert.Contains(t, result, "Balance")
		assert.Contains(t, result, "Keys")
		assert.Contains(t, result, "Contracts")
		close(done)
	}()

	select {
	case <-individualTimeout:
		t.Error("test timed out")
	case <-done:
	}
}

func getAccountIdFromCreate(createOutput []byte) (string, error) {
	addressIndex := addressRegex.SubexpIndex("address")

	matches := addressRegex.FindSubmatch(createOutput)
	if len(matches) < addressIndex {
		return "", fmt.Errorf("address not found in output: %v", string(createOutput))
	}
	return string(matches[addressIndex]), nil
}

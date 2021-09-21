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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-cli/tests/integration"
)

var addressRegex = regexp.MustCompile("Address\\s+(?P<address>0x[^\n]+)")

func init() {
	// bootstrap configuration if it doesn't exist
	_, _ = integration.RunFlowCmd("init", "--yes")
}

func TestMain(m *testing.M) {
	emu, err := integration.RunEmulator()
	code := m.Run()
	if err != nil {
		panic(fmt.Sprintf("unable to start emulator %v", err))
	}
	err = emu.Stop() // 🛑 the emulator
	if err != nil {
		fmt.Printf("unable to stop emulator: %v\n", err)
	}
	os.Exit(code)
}

func TestAccountCreateCommand(t *testing.T) {
	out, err := integration.RunFlowCmd("accounts", "create")
	if err != nil {
		t.Logf("unable to create account: %v", err)
		t.Fail()
	}
	id := getAccountIdFromCreate(out)
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
}

func getAccountIdFromCreate(createOutput []byte) string {
	addressIndex := addressRegex.SubexpIndex("address")
	matches := addressRegex.FindSubmatch(createOutput)

	return string(matches[addressIndex])
}

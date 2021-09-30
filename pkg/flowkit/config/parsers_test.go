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

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-cli/pkg/flowkit/config"
)

func TestStringToDeployments(t *testing.T) {
	testCases := []struct {
		name, network, account string
		contracts              []string
	}{
		{
			name:      "TestBasic",
			network:   "emulator",
			account:   "emulator-account",
			contracts: []string{"HelloWorld"},
		},
		{
			name:      "TestNoDuplicates",
			network:   "emulator",
			account:   "emulator-account",
			contracts: []string{"HelloWorld", "HelloWorld", "FlowServiceAccount"},
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			deployment := config.StringToDeployment(c.network, c.account, c.contracts)
			assert.Equal(t, c.account, deployment.Account)
			assert.Equal(t, c.network, deployment.Network)
			// check contract names match without duplicates
			assert.ElementsMatch(t, noDupes(c.contracts), namesOf(deployment.Contracts))
		})
	}
}

// noDupes removes duplicate strings
func noDupes(names []string) []string {
	found := map[string]bool{}
	unique := []string{}
	for _, name := range names {
		_, has := found[name]
		if !has {
			found[name] = true
			unique = append(unique, name)
		}
	}
	return unique
}

// namesOf returns the names of supplied contracts
func namesOf(contracts []config.ContractDeployment) []string {
	names := make([]string, len(contracts))
	for i, contract := range contracts {
		names[i] = contract.Name
	}
	return names
}

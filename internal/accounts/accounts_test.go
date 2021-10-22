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

package accounts

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit/util"
)

// One way to handle testing a result serializes as expected - convert it to a map[string]interface{} and make
// many type assertions
func TestAccountsResultOutput(t *testing.T) {

	with10Keys := test.AccountGenerator().New()
	with10Keys.Keys = nKeys(10)

	testCases := []struct {
		Name    string
		Account *flow.Account
	}{
		{Name: "Generated Account", Account: test.AccountGenerator().New()},
		{Name: "With 10 Keys", Account: with10Keys},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			account := tc.Account
			result := AccountResult{Account: tc.Account}

			assert.Equal(t, result.Oneliner(), toOneliner(account, result.include))

			result.include = []string{"contracts"}
			assert.Equal(t, result.JSON(), toJSON(account, result.include))
			assert.Equal(t, result.String(), toString(account, result.include))

			result.include = []string{"contracts", "keys"}
			assert.Equal(t, result.String(), toString(account, result.include))
		})
	}
}

func nKeys(numKeys int) []*flow.AccountKey {
	keyGen := test.AccountKeyGenerator()
	keys := make([]*flow.AccountKey, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = keyGen.New()
	}
	return keys
}

func toJSON(a *flow.Account, includes []string) interface{} {
	result := make(map[string]interface{})
	result["address"] = a.Address
	result["balance"] = fetchBalance(a)

	keys := make([]string, 0)
	for _, key := range a.Keys {
		keys = append(keys, fmt.Sprintf("%x", key.PublicKey.Encode()))
	}

	result["keys"] = keys

	contracts := make([]string, 0, len(a.Contracts))
	for name := range a.Contracts {
		contracts = append(contracts, name)
	}

	result["contracts"] = contracts

	if command.ContainsFlag(includes, "contracts") {
		c := make(map[string]string)
		for name, code := range a.Contracts {
			c[name] = string(code)
		}
		result["code"] = c
	}

	return result
}

func fetchBalance(a *flow.Account) string {
	return cadence.UFix64(a.Balance).String()
}

func toString(a *flow.Account, includes []string) string {
	var b bytes.Buffer
	writer := util.CreateTabWriter(&b)

	_, _ = fmt.Fprintf(writer, "Address\t 0x%s\n", a.Address)
	_, _ = fmt.Fprintf(writer, "Balance\t %s\n", fetchBalance(a))

	_, _ = fmt.Fprintf(writer, "Keys\t %d\n", len(a.Keys))

	for i, key := range a.Keys {
		_, _ = fmt.Fprintf(writer, "\nKey %d\tPublic Key\t %x\n", i, key.PublicKey.Encode())
		_, _ = fmt.Fprintf(writer, "\tWeight\t %d\n", key.Weight)
		_, _ = fmt.Fprintf(writer, "\tSignature Algorithm\t %s\n", key.SigAlgo)
		_, _ = fmt.Fprintf(writer, "\tHash Algorithm\t %s\n", key.HashAlgo)
		_, _ = fmt.Fprintf(writer, "\tRevoked \t %t\n", key.Revoked)
		_, _ = fmt.Fprintf(writer, "\tSequence Number \t %d\n", key.SequenceNumber)
		_, _ = fmt.Fprintf(writer, "\tIndex \t %d\n", key.Index)
		_, _ = fmt.Fprintf(writer, "\n")

		// only show up to 3 keys and then show label to expand more info
		if i == 3 && !command.ContainsFlag(includes, "keys") {
			_, _ = fmt.Fprint(writer, "...keys minimized, use --include keys flag if you want to view all\n\n")
			break
		}
	}

	_, _ = fmt.Fprintf(writer, "Contracts Deployed: %d\n", len(a.Contracts))
	for name := range a.Contracts {
		_, _ = fmt.Fprintf(writer, "Contract: '%s'\n", name)
	}

	if command.ContainsFlag(includes, "contracts") {
		for name, code := range a.Contracts {
			_, _ = fmt.Fprintf(writer, "Contracts '%s':\n", name)
			_, _ = fmt.Fprintln(writer, string(code))
		}
	} else {
		_, _ = fmt.Fprint(writer, "\n\nContracts (hidden, use --include contracts)")
	}

	_ = writer.Flush()

	return b.String()
}

func toOneliner(a *flow.Account, includes []string) string {
	keys := make([]string, 0, len(a.Keys))
	for _, key := range a.Keys {
		keys = append(keys, key.PublicKey.String())
	}

	return fmt.Sprintf("Address: 0x%s, Balance: %s, Public Keys: %s", a.Address, cadence.UFix64(a.Balance), keys)
}

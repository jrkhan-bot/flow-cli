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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/tests"
)

func TestAddContract(t *testing.T) {
	services, state, gw := setup()

	serviceAcc, _ := state.EmulatorServiceAccount()
	serviceAddress := serviceAcc.Address()
	state.Accounts().AddOrUpdate(serviceAcc)

	gw.SendSignedTransaction.Run(func(args mock.Arguments) {
		tx := args.Get(0).(*flowkit.Transaction)
		assert.Equal(t, tx.Signer().Address(), serviceAddress)
		assert.True(t, strings.Contains(string(tx.FlowTransaction().Script), "signer.contracts.add"))

		gw.SendSignedTransaction.Return(tests.NewTransaction(), nil)
	})

	resultAccount := tests.NewAccountWithAddress(serviceAddress.String())
	resultAccount.Contracts = map[string][]byte{"Simple": tests.ContractSimple.Source}
	gw.GetAccount.Run(func(args mock.Arguments) {
		// make sure the account we return contains the newly added contract
		gw.GetAccount.Return(resultAccount, nil)
	})

	// assign package level flags here
	addContractFlags = flagsAddContract{
		Signer:  serviceAcc.Name(),
		Include: []string{"contracts"},
	}
	// reset flags
	defer func() { addContractFlags = flagsAddContract{} }()
	args := []string{"Simple", "contractSimple.cdc"}
	res, err := addContract(args, state.ReaderWriter(), command.GlobalFlags{}, services, state)
	assert.NoError(t, err)

	ar := res.(*AccountResult)
	assert.Contains(t, ar.Contracts, "Simple")

	gw.Mock.AssertNumberOfCalls(t, tests.GetTransactionResultFunc, 1)
	gw.Mock.AssertNumberOfCalls(t, tests.SendSignedTransactionFunc, 1)
}

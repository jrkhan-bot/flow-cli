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
	"encoding/json"
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
	"github.com/onflow/flow-cli/tests"
)

func setup() (*services.Services, *flowkit.State, *tests.TestGateway) {
	readerWriter := tests.ReaderWriter()
	state, err := flowkit.Init(readerWriter, crypto.ECDSA_P256, crypto.SHA3_256)
	if err != nil {
		panic(err)
	}

	gw := tests.DefaultMockGateway()
	s := services.NewServices(gw.Mock, state, output.NewStdoutLogger(output.NoneLog))
	return s, state, gw
}

func TestAccountsCreate(t *testing.T) {
	services, state, gw := setup()

	newAddress := flow.HexToAddress("192440c99cb17282")
	expectedAccount := tests.NewAccountWithAddress(newAddress.String())
	gw.GetAccount.Run(func(args mock.Arguments) {
		gw.GetAccount.Return(expectedAccount, nil)
	})

	gw.GetTransactionResult.Return(
		tests.NewAccountCreateResult(newAddress), nil,
	)

	// assign package level flags here
	createFlags = flagsCreate{
		Signer:    "emulator-account",
		Keys:      []string{},
		SigAlgo:   "ECDSA_P256",
		HashAlgo:  "SHA3_256",
		Contracts: nil,
		Include:   []string{"contracts"},
	}

	// reset flags
	defer func() { createFlags = flagsCreate{} }()

	res, err := create(nil, nil, command.GlobalFlags{}, services, state)

	assert.NoError(t, err)

	// confirm the account returned from create is the one returned from mock gateway
	assert.Equal(t, expectedAccount, res.(*AccountResult).Account)
	out := res.JSON()

	// cycle json
	bytes, err := json.Marshal(out)
	assert.NoError(t, err)
	account := struct {
		Address string            `json:"address"`
		Code    map[string]string `json:"code"`
	}{}
	err = json.Unmarshal(bytes, &account)
	assert.NoError(t, err)

	// confirm the address is the one we created in our mock gateway response
	assert.Equal(t, account.Address, newAddress.Hex())
}

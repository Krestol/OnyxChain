/*
 * Copyright (C) 2019 The onyxchain Authors
 * This file is part of The onyxchain library.
 *
 * The onyxchain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The onyxchain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The onyxchain.  If not, see <http://www.gnu.org/licenses/>.
 */

package test

import (
	"os"
	"testing"

	"github.com/OnyxPay/OnyxChain/common/log"
	"github.com/OnyxPay/OnyxChain/core/types"
	. "github.com/OnyxPay/OnyxChain/smartcontract"
	"github.com/OnyxPay/OnyxChain/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestEqualStackOverflow(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("./Log")
	}()

	code := []byte{
		byte(neovm.PUSH1),    // {1}
		byte(neovm.NEWARRAY), // {[]}
		byte(neovm.DUP),      // {[],[]}
		byte(neovm.DUP),      // {[],[],[]}
		byte(neovm.PUSH0),    // {[],[],[],0}
		byte(neovm.ROT),      // {[],[],0,[]}
		byte(neovm.SETITEM),  // {[[]]}
		byte(neovm.DUP),      // {[[]],[[]]}
		byte(neovm.EQUAL),
	}

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	sc := SmartContract{
		Config:  config,
		Gas:     10000,
		CacheDB: nil,
	}
	engine, _ := sc.NewExecuteEngine(code)
	_, err := engine.Invoke()

	assert.Nil(t, err)
}

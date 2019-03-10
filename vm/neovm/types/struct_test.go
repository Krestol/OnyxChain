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

package types

import (
	"math/big"
	"testing"
)

func TestStruct_Clone(t *testing.T) {
	s := NewStruct(nil)
	//d := NewStruct([]StackItems{s})
	k := NewStruct([]StackItems{NewInteger(big.NewInt(1))})
	for i := 0; i < MAX_STRUCT_DEPTH-2; i++ {
		k = NewStruct([]StackItems{k})
	}
	//k.Add(d)
	s.Add(k)

	if _, err := s.Clone(); err != nil {
		t.Fatal(err)
	}

}

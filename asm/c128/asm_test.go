// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128

import (
	"math"
	"math/cmplx"
	"testing"
)

var (
	inf = math.Inf(1)
)

func same(x, y complex128) bool {
	return x == y || (cmplx.IsNaN(x) && cmplx.IsNaN(y))
}

func guardVector(vec []complex128, gdVal complex128, gdLen int) (guarded []complex128) {
	guarded = make([]complex128, len(vec)+gdLen*2)
	copy(guarded[gdLen:], vec)
	for i := 0; i < gdLen; i++ {
		guarded[i] = gdVal
		guarded[len(guarded)-1-i] = gdVal
	}
	return guarded
}

func isValidGuard(vec []complex128, gdVal complex128, gdLen int) bool {
	for i := 0; i < gdLen; i++ {
		if !same(vec[i], gdVal) || !same(vec[len(vec)-1-i], gdVal) {
			return false
		}
	}
	return true
}

func guardIncVector(vec []complex128, gdVal complex128, inc, gdLen int) (guarded []complex128) {
	if inc < 0 {
		inc = -inc
	}
	inrLen := len(vec) * inc
	guarded = make([]complex128, inrLen+gdLen*2)
	for i := range guarded {
		guarded[i] = gdVal
	}
	for i, v := range vec {
		guarded[gdLen+i*inc] = v
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, vec []complex128, gdVal complex128, inc, gdLen int) {
	srcLen := len(vec) - 2*gdLen
	if inc < 0 {
		srcLen = len(vec) * -inc
	}

	for i := range vec {
		switch {
		case same(vec[i], gdVal):
			// Correct value
		case (i-gdLen)%inc == 0 && (i-gdLen)/inc < len(vec):
			// Ignore input values
		case i < gdLen:
			t.Errorf("Front guard violated at %d %v", i, vec[:gdLen])
		case i > gdLen+srcLen:
			t.Errorf("Back guard violated at %d %v", i-gdLen-srcLen, vec[gdLen+srcLen:])
		default:
			t.Errorf("Internal guard violated at %d %v", i-gdLen, vec[gdLen:gdLen+srcLen])
		}
	}
}

var ( // Offset sets for testing alignment handling in Unitary assembly functions.
	align1 = []int{0, 1}
	align2 = newIncSet(0, 1)
	align3 = newIncToSet(0, 1)
)

type incSet struct {
	x, y int
}

// genInc will generate all (x,y) combinations of the input increment set.
func newIncSet(inc ...int) []incSet {
	n := len(inc)
	is := make([]incSet, n*n)
	for x := range inc {
		for y := range inc {
			is[x*n+y] = incSet{inc[x], inc[y]}
		}
	}
	return is
}

type incToSet struct {
	dst, x, y int
}

// genIncTo will generate all (dst,x,y) combinations of the input increment set.
func newIncToSet(inc ...int) []incToSet {
	n := len(inc)
	is := make([]incToSet, n*n*n)
	for i, dst := range inc {
		for x := range inc {
			for y := range inc {
				is[i*n*n+x*n+y] = incToSet{dst, inc[x], inc[y]}
			}
		}
	}
	return is
}

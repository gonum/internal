// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c64

import (
	"fmt"
	"math"
	"testing"
)

const (
	msgVal   = "%v: unexpected value at %v Got: %v Expected: %v"
	msgGuard = "%v: Guard violated in %s vector %v %v"
)

var (
	inf = float32(math.Inf(1))
)

var dotTests = []struct {
	x, y           []complex64
	exu, exc       complex64
	exuRev, excRev complex64
	n              int
}{
	{
		x:   []complex64{},
		y:   []complex64{},
		n:   0,
		exu: 0, exc: 0,
		exuRev: 0, excRev: 0,
	},
	{
		x:   []complex64{1 + 1i},
		y:   []complex64{1 + 1i},
		n:   1,
		exu: 0 + 2i, exc: 2,
		exuRev: 0 + 2i, excRev: 2,
	},
	{
		x:   []complex64{1 + 2i},
		y:   []complex64{1 + 1i},
		n:   1,
		exu: -1 + 3i, exc: 3 - 1i,
		exuRev: -1 + 3i, excRev: 3 - 1i,
	},
	{
		x:   []complex64{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i, 13 + 14i, 15 + 16i, 17 + 18i, 19 + 20i},
		y:   []complex64{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i, 13 + 14i, 15 + 16i, 17 + 18i, 19 + 20i},
		n:   10,
		exu: -210 + 2860i, exc: 2870 + 0i,
		exuRev: -210 + 1540i, excRev: 1550 + 0i,
	},
	{
		x:   []complex64{1 + 1i, 1 + 1i, 1 + 2i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 3i, 1 + 1i, 1 + 1i, 1 + 4i},
		y:   []complex64{1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i},
		n:   10,
		exu: -22 + 36i, exc: 42 + 4i,
		exuRev: -22 + 36i, excRev: 42 + 4i,
	},
	{
		x:   []complex64{1 + 1i, 1 + 1i, 2 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 2 + 1i},
		y:   []complex64{1 + 2i, 1 + 2i, 1 + 3i, 1 + 2i, 1 + 3i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i},
		n:   10,
		exu: -10 + 37i, exc: 34 + 17i,
		exuRev: -10 + 36i, excRev: 34 + 16i,
	},
	{
		x:   []complex64{1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, complex(inf, 1), 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i},
		y:   []complex64{1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i, 1 + 2i},
		n:   10,
		exu: complex(inf, inf), exc: complex(inf, inf),
		exuRev: complex(inf, inf), excRev: complex(inf, inf),
	},
}

func TestDotcUnitary(t *testing.T) {
	const gd = 1 + 5i
	for i, test := range dotTests {
		for _, align := range align2 {
			prefix := fmt.Sprintf("Test %v (x:%v y:%v)", i, align.x, align.y)
			xgLn, ygLn := 4+align.x, 4+align.y
			xg, yg := guardVector(test.x, gd, xgLn), guardVector(test.y, gd, ygLn)
			x, y := xg[xgLn:len(xg)-xgLn], yg[ygLn:len(yg)-ygLn]
			res := DotcUnitary(x, y)
			if !same(res, test.exc) {
				t.Errorf(msgVal, prefix, i, res, test.exc)
			}
			if !isValidGuard(xg, gd, xgLn) {
				t.Errorf(msgGuard, prefix, "x", xg[:xgLn], xg[len(xg)-xgLn:])
			}
			if !isValidGuard(yg, gd, ygLn) {
				t.Errorf(msgGuard, prefix, "y", yg[:ygLn], yg[len(yg)-ygLn:])
			}
		}
	}
}

func TestDotcInc(t *testing.T) {
	const gd, gdLn = 2 + 5i, 4
	for i, test := range dotTests {
		for _, inc := range newIncSet(1, 2, 5, 10, -1, -2, -5, -10) {
			xg, yg := guardIncVector(test.x, gd, inc.x, gdLn), guardIncVector(test.y, gd, inc.y, gdLn)
			x, y := xg[gdLn:len(xg)-gdLn], yg[gdLn:len(yg)-gdLn]
			want := test.exc
			var ix, iy int
			if inc.x < 0 {
				ix, want = -inc.x*(test.n-1), test.excRev
			}
			if inc.y < 0 {
				iy, want = -inc.y*(test.n-1), test.excRev
			}
			prefix := fmt.Sprintf("Test %v (x:%v y:%v) (ix:%v iy:%v)", i, inc.x, inc.y, ix, iy)
			res := DotcInc(x, y, test.n, inc.x, inc.y, ix, iy)
			if inc.x*inc.y > 0 {
				want = test.exc
			}
			if !same(res, want) {
				t.Errorf(msgVal, prefix, i, res, want)
				t.Error(x, y)
			}
			checkValidIncGuard(t, xg, gd, inc.x, gdLn)
			checkValidIncGuard(t, yg, gd, inc.y, gdLn)
		}
	}
}

func TestDotuUnitary(t *testing.T) {
	const gd = 1 + 5i
	for i, test := range dotTests {
		for _, align := range align2 {
			prefix := fmt.Sprintf("Test %v (x:%v y:%v)", i, align.x, align.y)
			xgLn, ygLn := 4+align.x, 4+align.y
			xg, yg := guardVector(test.x, gd, xgLn), guardVector(test.y, gd, ygLn)
			x, y := xg[xgLn:len(xg)-xgLn], yg[ygLn:len(yg)-ygLn]
			res := DotuUnitary(x, y)
			if !same(res, test.exu) {
				t.Errorf(msgVal, prefix, i, res, test.exu)
			}
			if !isValidGuard(xg, gd, xgLn) {
				t.Errorf(msgGuard, prefix, "x", xg[:xgLn], xg[len(xg)-xgLn:])
			}
			if !isValidGuard(yg, gd, ygLn) {
				t.Errorf(msgGuard, prefix, "y", yg[:ygLn], yg[len(yg)-ygLn:])
			}
		}
	}
}

func TestDotuInc(t *testing.T) {
	const gd, gdLn = 1 + 5i, 4
	for i, test := range dotTests {
		for _, inc := range newIncSet(1, 2, 5, 10, -1, -2, -5, -10) {
			prefix := fmt.Sprintf("Test %v (x:%v y:%v)", i, inc.x, inc.y)
			xg, yg := guardIncVector(test.x, gd, inc.x, gdLn), guardIncVector(test.y, gd, inc.y, gdLn)
			x, y := xg[gdLn:len(xg)-gdLn], yg[gdLn:len(yg)-gdLn]
			want := test.exc
			var ix, iy int
			if inc.x < 0 {
				ix, want = -inc.x*(test.n-1), test.exuRev
			}
			if inc.y < 0 {
				iy, want = -inc.y*(test.n-1), test.exuRev
			}
			res := DotuInc(x, y, test.n, inc.x, inc.y, ix, iy)
			if inc.x*inc.y > 0 {
				want = test.exu
			}
			if !same(res, want) {
				t.Errorf(msgVal, prefix, i, res, want)
			}
			checkValidIncGuard(t, xg, gd, inc.x, gdLn)
			checkValidIncGuard(t, yg, gd, inc.y, gdLn)
		}
	}
}

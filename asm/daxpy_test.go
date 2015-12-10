// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"testing"
)

func TestDaxpyUnitary(t *testing.T) {
	// Test z = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		z, zFront, zBack := newGuardedVector(test.xData, 1)

		DaxpyUnitary(test.alpha, x, y, z)

		prefix := fmt.Sprintf("test %v (z=a*x+y)", i)

		if err := checkGuardsXYZ(xFront, xBack, yFront, yBack, zFront, zBack); err != nil {
			t.Errorf("%v: %v", prefix, err)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, z, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, z)
		}
	}

	// Test y = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)

		DaxpyUnitary(test.alpha, x, y, y)

		prefix := fmt.Sprintf("test %v (y=a*x+y)", i)

		if err := checkGuardsXY(xFront, xBack, yFront, yBack); err != nil {
			t.Errorf("%v: %v", prefix, err)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, y, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, y)
		}
	}

	// Test x = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)

		DaxpyUnitary(test.alpha, x, y, x)

		prefix := fmt.Sprintf("test %v (x=a*x+y)", i)

		if err := checkGuardsXY(xFront, xBack, yFront, yBack); err != nil {
			t.Errorf("%v: %v", prefix, err)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}

	// Test x = alpha * x + x.
	for i, test := range []struct {
		alpha float64
		xData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			want:  []float64{2},
		},
		{
			alpha: 3,
			xData: []float64{2},
			want:  []float64{8},
		},
		{
			alpha: -3,
			xData: []float64{2},
			want:  []float64{-4},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, 1, 2, -3, -4},
		},
		{
			alpha: 3,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, 4, 8, -12, -16},
		},
		{
			alpha: -3,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, -2, -4, 6, 8},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 1, 2, -3, -4, 5},
			want:  []float64{0, -4, -8, 12, 16, -20},
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)

		DaxpyUnitary(test.alpha, x, x, x)

		prefix := fmt.Sprintf("test %v (x=a*x+x)", i)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf("%v: out-of-bounds write to x argument\nfront guard: %v\nback guard: %v", prefix, xFront, xBack)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDaxpyInc(t *testing.T) {
	// Test z = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				for _, incZ := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
					switch {
					case incX > 0 && incY > 0 && incZ > 0:
					case incX < 0 && incY < 0 && incZ < 0:
					default:
						continue
					}
					// All increments have the same sign.

					x, xFront, xBack := newGuardedVector(test.xData, incX)
					y, yFront, yBack := newGuardedVector(test.yData, incY)
					z, zFront, zBack := newGuardedVector(test.xData, incZ)
					n := len(test.xData)

					var ix, iy, iz int
					if incX < 0 {
						ix = (-n + 1) * incX
					}
					if incY < 0 {
						iy = (-n + 1) * incY
					}
					if incZ < 0 {
						iz = (-n + 1) * incZ
					}
					DaxpyInc(test.alpha, x, y, z, uintptr(n),
						uintptr(incX), uintptr(incY), uintptr(incZ),
						uintptr(ix), uintptr(iy), uintptr(iz))

					prefix := fmt.Sprintf("test %v (z=a*x+y), incX = %v, incY = %v, incZ = %v", i, incX, incY, incZ)

					if err := checkGuardsXYZ(xFront, xBack, yFront, yBack, zFront, zBack); err != nil {
						t.Errorf("%v: %v", prefix, err)
					}
					if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
						t.Errorf("%v: modified read-only x argument", prefix)
					}
					if nonStridedWrite(y, incY) || !equalStrided(test.yData, y, incY) {
						t.Errorf("%v: modified read-only y argument", prefix)
					}
					if nonStridedWrite(z, incZ) {
						t.Errorf("%v: modified z argument at non-stride position", prefix)
					}

					if !equalStrided(test.want, z, incZ) {
						t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, z)
					}
				}
			}
		}
	}

	// Test y = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				switch {
				case incX > 0 && incY > 0:
				case incX < 0 && incY < 0:
				default:
					continue
				}
				// All increments have the same sign.

				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)
				n := len(test.xData)

				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}
				DaxpyInc(test.alpha, x, y, y, uintptr(n),
					uintptr(incX), uintptr(incY), uintptr(incY),
					uintptr(ix), uintptr(iy), uintptr(iy))

				prefix := fmt.Sprintf("test %v (y=a*x+y), incX = %v, incY = %v", i, incX, incY)

				if err := checkGuardsXY(xFront, xBack, yFront, yBack); err != nil {
					t.Errorf("%v: %v", prefix, err)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(y, incY) {
					t.Errorf("%v: modified y argument at non-stride position", prefix)
				}

				if !equalStrided(test.want, y, incY) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, y)
				}
			}
		}
	}

	// Test x = alpha * x + y.
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				switch {
				case incX > 0 && incY > 0:
				case incX < 0 && incY < 0:
				default:
					continue
				}
				// All increments have the same sign.

				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)
				n := len(test.xData)

				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}
				DaxpyInc(test.alpha, x, y, x, uintptr(n),
					uintptr(incX), uintptr(incY), uintptr(incX),
					uintptr(ix), uintptr(iy), uintptr(ix))

				prefix := fmt.Sprintf("test %v (x=a*x+y), incX = %v, incY = %v", i, incX, incY)

				if err := checkGuardsXY(xFront, xBack, yFront, yBack); err != nil {
					t.Errorf("%v: %v", prefix, err)
				}
				if nonStridedWrite(y, incY) || !equalStrided(test.yData, y, incY) {
					t.Errorf("%v: modified read-only y argument", prefix)
				}
				if nonStridedWrite(x, incX) {
					t.Errorf("%v: modified x argument at non-stride position", prefix)
				}

				if !equalStrided(test.want, x, incX) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
				}
			}
		}
	}

	// Test x = alpha * x + x.
	for i, test := range []struct {
		alpha float64
		xData []float64

		want []float64
	}{
		// One element
		{
			alpha: 0,
			xData: []float64{2},
			want:  []float64{2},
		},
		{
			alpha: 3,
			xData: []float64{2},
			want:  []float64{8},
		},
		{
			alpha: -3,
			xData: []float64{2},
			want:  []float64{-4},
		},
		// Odd number of elements
		{
			alpha: 0,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, 1, 2, -3, -4},
		},
		{
			alpha: 3,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, 4, 8, -12, -16},
		},
		{
			alpha: -3,
			xData: []float64{0, 1, 2, -3, -4},
			want:  []float64{0, -2, -4, 6, 8},
		},
		// Even number of elements
		{
			alpha: -5,
			xData: []float64{0, 1, 2, -3, -4, 5},
			want:  []float64{0, -4, -8, 12, 16, -20},
		},
	} {
		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			x, xFront, xBack := newGuardedVector(test.xData, incX)
			n := len(test.xData)

			var ix int
			if incX < 0 {
				ix = (-n + 1) * incX
			}
			DaxpyInc(test.alpha, x, x, x, uintptr(n),
				uintptr(incX), uintptr(incX), uintptr(incX),
				uintptr(ix), uintptr(ix), uintptr(ix))

			prefix := fmt.Sprintf("test %v (x=a*x+x), incX = %v", i, incX)

			if !allNaN(xFront) || !allNaN(xBack) {
				t.Errorf("%v: out-of-bounds write to x argument\nfront guard: %v\nback guard: %v", i, xFront, xBack)
			}
			if nonStridedWrite(x, incX) {
				t.Errorf("%v: modified x argument at non-stride position", prefix)
			}

			if !equalStrided(test.want, x, incX) {
				t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
			}
		}
	}
}

// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !amd64 noasm appengine

package asm

// The extra z parameter is needed because of floats.AddScaledTo
func DaxpyUnitary(alpha float64, x, y, z []float64) {
	for i, v := range x {
		z[i] = alpha*v + y[i]
	}
}

func DaxpyInc(alpha float64, x, y, z []float64, n, incX, incY, incZ, ix, iy, iz uintptr) {
	for i := 0; i < int(n); i++ {
		z[iz] = alpha*x[ix] + y[iy]
		ix += incX
		iy += incY
		iz += incZ
	}
}

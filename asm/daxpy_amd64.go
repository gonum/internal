// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package asm

func DaxpyUnitary(alpha float64, x, y, z []float64)

func DaxpyInc(alpha float64, x, y, z []float64, n, incX, incY, incZ, ix, iy, iz uintptr)

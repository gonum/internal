// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define MOVSHDUP_X3_X2    LONG $0xD3160FF3 // @ MOVSHDUP X3, X2
#define MOVSHDUP_X5_X4    LONG $0xE5160FF3 // @ MOVSHDUP X5, X4
#define MOVSHDUP_X7_X6    LONG $0xF7160FF3 // @ MOVSHDUP X7, X6
#define MOVSHDUP_X9_X8    LONG $0x160F45F3; BYTE $0xC1 // @ MOVSHDUP X9, X8

#define MOVSLDUP_X3_X3    LONG $0xDB120FF3 // @ MOVSLDUP X3, X3
#define MOVSLDUP_X5_X5    LONG $0xED120FF3 // @ MOVSLDUP X5, X5
#define MOVSLDUP_X7_X7    LONG $0xFF120FF3 // @ MOVSLDUP X7, X7
#define MOVSLDUP_X9_X9    LONG $0x120F45F3; BYTE $0xC9 // @ MOVSLDUP X9, X9

#define ADDSUBPS_X2_X3    LONG $0xDAD00FF2 // @ ADDSUBPS X2, X3
#define ADDSUBPS_X4_X5    LONG $0xECD00FF2 // @ ADDSUBPS X4, X5
#define ADDSUBPS_X6_X7    LONG $0xFED00FF2 // @ ADDSUBPS X6, X7
#define ADDSUBPS_X8_X9    LONG $0xD00F45F2; BYTE $0xC8 // @ ADDSUBPS X8, X9

// func DotcInc(x, y []complex64, n, incX, incY, ix, iy uintptr) (sum complex64)
TEXT ·DotcInc(SB), NOSPLIT, $0
	MOVQ   x_base+0(FP), SI  // SI := &x
	MOVQ   y_base+24(FP), DI // DI := &y
	PXOR   X0, X0            // psum := 0
	PXOR   X1, X1
	MOVQ   n+48(FP), CX      // CX := n
	CMPQ   CX, $0            // if CX == 0 { return }
	JE     dotc_end
	MOVQ   ix+72(FP), R8
	MOVQ   iy+80(FP), R9
	LEAQ   (SI)(R8*8), SI    // SI = &(SI[ix])
	LEAQ   (DI)(R9*8), DI    // DI = &(DI[iy])
	MOVQ   incX+56(FP), R8   // R8 := incX * sizeof(complex64)
	SHLQ   $3, R8
	MOVQ   incY+64(FP), R9   // R9 := incY * sizeof(complex64)
	SHLQ   $3, R9
	MOVSS  $(-1.0), X15
	SHUFPS $0, X15, X15      // { -1, -1, -1, -1 }

	MOVQ CX, BX
	ANDQ $3, BX    // BX = n % 4
	SHRQ $2, CX    // CX = floor( n / 4 )
	JZ   dotc_tail // if CX == 0 { goto dotc_tail }

	MOVUPS X15, X14        // Copy X15 for pipelining
	LEAQ   (R8)(R8*2), R10 // R10 = R8 * 3
	LEAQ   (R9)(R9*2), R11 // R11 = R9 * 3

dotc_loop: // do {
	MOVSD (SI), X3        // X_i = { imag(x[i]), real(x[i]) }
	MOVSD (SI)(R8*1), X5
	MOVSD (SI)(R8*2), X7
	MOVSD (SI)(R10*1), X9

	// X_(i-1) = { imag(x[i]), imag(x[i]) }
	MOVSHDUP_X3_X2
	MOVSHDUP_X5_X4
	MOVSHDUP_X7_X6
	MOVSHDUP_X9_X8

	// X_i = { real(x[i]), real(x[i]) }
	MOVSLDUP_X3_X3
	MOVSLDUP_X5_X5
	MOVSLDUP_X7_X7
	MOVSLDUP_X9_X9

	// X_(i-1) = { -imag(x[i]), -imag(x[i]) }
	MULPS X15, X2
	MULPS X14, X4
	MULPS X15, X6
	MULPS X14, X8

	// X_j = { imag(y[i]), real(y[i]) }
	MOVSD (DI), X10
	MOVSD (DI)(R9*1), X11
	MOVSD (DI)(R9*2), X12
	MOVSD (DI)(R11*1), X13

	// X_i     = { imag(y[i]) * real(x[i]), real(y[i]) * real(x[i]) }
	MULPS X10, X3
	MULPS X11, X5
	MULPS X12, X7
	MULPS X13, X9

	// X_j = { real(y[i]), imag(y[i]) }
	SHUFPS $0xB1, X10, X10
	SHUFPS $0xB1, X11, X11
	SHUFPS $0xB1, X12, X12
	SHUFPS $0xB1, X13, X13

	// X_(i-1) = { real(y[i]) * imag(x[i]), imag(y[i]) * imag(x[i]) }
	MULPS X10, X2
	MULPS X11, X4
	MULPS X12, X6
	MULPS X13, X8

	// X_i = {
	//	imag(result[i]):  imag(y[i]) * real(x[i]) + real(y[i]) * imag(x[i]),
	//	real(result[i]):  real(y[i]) * real(x[i]) - imag(y[i]) * imag(x[i])  }
	ADDSUBPS_X2_X3
	ADDSUBPS_X4_X5
	ADDSUBPS_X6_X7
	ADDSUBPS_X8_X9

	// psum += X_i
	ADDPS X3, X0
	ADDPS X5, X1
	ADDPS X7, X0
	ADDPS X9, X1

	LEAQ (SI)(R8*4), SI // SI = &(SI[incX*4])
	LEAQ (DI)(R9*4), DI // DI = &(DI[incY*4])

	DECQ CX
	JNZ  dotc_loop // } while --CX > 0

	ADDPS X1, X0   // X0 = { psum_1 + psum_0 }
	CMPQ  BX, $0   // if BX == 0 { return }
	JE    dotc_end

dotc_tail: // do {
	MOVSD  (SI), X3       // X_i = { imag(x[i]), real(x[i]) }
	MOVSHDUP_X3_X2        // X_(i-1) = { imag(x[i]), imag(x[i]) }
	MOVSLDUP_X3_X3        // X_i = { real(x[i]), real(x[i]) }
	MULPS  X15, X2        // X_(i-1) = { -imag(x[i]), imag(x[i]) }
	MOVUPS (DI), X10      // X_j = { imag(y[i]), real(y[i]) }
	MULPS  X10, X3        // X_i = { imag(y[i]) * real(x[i]), real(y[i]) * real(x[i]) }
	SHUFPS $0x1, X10, X10 // X_j = { real(y[i]), imag(y[i]) }
	MULPS  X10, X2        // X_(i-1) = { real(y[i]) * imag(x[i]), imag(y[i]) * imag(x[i]) }

	// X_i = {
	//	imag(result[i]):  imag(y[i])*real(x[i]) + real(y[i])*imag(x[i]),
	//	real(result[i]):  real(y[i])*real(x[i]) - imag(y[i])*imag(x[i]) }
	ADDSUBPS_X2_X3
	ADDPS X3, X0    // sum += X_i
	ADDQ  R8, SI    // SI += incX
	ADDQ  R9, DI    // DI += incY
	DECQ  BX
	JNZ   dotc_tail // } while --BX > 0

dotc_end:
	MOVSD X0, sum+88(FP) // return sum
	RET

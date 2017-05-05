// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define MOVDDUP_SI__X3    LONG $0x1E120FF2 // +	MOVDDUP (SI), X3
#define MOVDDUP_SI_R8__X5    LONG $0x120F42F2; WORD $0x062C // + MOVDDUP (SI)(R8*1), X5
#define MOVDDUP_SI_R8_2__X7    LONG $0x120F42F2; WORD $0x463C // + MOVDDUP (SI)(R8*2), X7
#define MOVDDUP_SI_R9__X9    LONG $0x120F46F2; WORD $0x0E0C // + MOVDDUP (SI)(R9*1), X9

#define MOVDDUP_8_SI__X2    LONG $0x56120FF2; BYTE $0x08 // + MOVDDUP 8(SI), X2
#define MOVDDUP_8_SI_R8__X4    LONG $0x120F42F2; WORD $0x0664; BYTE $0x08 // + MOVDDUP 8(SI)(R8*1), X4
#define MOVDDUP_8_SI_R8_2__X6    LONG $0x120F42F2; WORD $0x4674; BYTE $0x08 // + MOVDDUP 8(SI)(R8*2), X6
#define MOVDDUP_8_SI_R9__X8    LONG $0x120F46F2; WORD $0x0E44; BYTE $0x08 // + MOVDDUP 8(SI)(R9*1), X8

#define ADDSUBPD_X2_X3    LONG $0xDAD00F66 // + ADDSUBPD X2, X3
#define ADDSUBPD_X4_X5    LONG $0xECD00F66 // + ADDSUBPD X4, X5
#define ADDSUBPD_X6_X7    LONG $0xFED00F66 // + ADDSUBPD X6, X7
#define ADDSUBPD_X8_X9    LONG $0xD00F4566; BYTE $0xC8 // + ADDSUBPD X8, X9

// func DotuInc(x, y []complex128, n, incX, incY, ix, iy uintptr) (sum complex128)
TEXT ·DotuInc(SB), NOSPLIT, $0
	MOVQ x_base+0(FP), SI  // SI = &x
	MOVQ y_base+24(FP), DI // DI = &y
	MOVQ n+48(FP), CX      // CX = n
	PXOR X0, X0            // sum = 0
	CMPQ CX, $0            // if CX == 0 { return }
	JE   dot_end
	PXOR X1, X1            // psum = 0
	MOVQ ix+72(FP), R8     // R8 = ix * sizeof(complex128)
	SHLQ $4, R8
	MOVQ iy+80(FP), R10    // R10 = iy * sizeof(complex128)
	SHLQ $4, R10
	LEAQ (SI)(R8*1), SI    // SI = &(SI[ix])
	LEAQ (DI)(R10*1), DI   // DI = &(DI[iy])
	MOVQ incX+56(FP), R8   // R8 = incX
	SHLQ $4, R8            // R8 *=  sizeof(complex128)
	MOVQ incY+64(FP), R10  // R10 = incY
	SHLQ $4, R10           // R10 *=  sizeof(complex128)
	MOVQ CX, BX
	ANDQ $3, CX            // CX = n % 4
	SHRQ $2, BX            // BX = floor( n / 4 )
	JZ   dot_tail          // if n <= 4 { goto dot_tail }
	LEAQ (R8)(R8*2), R9    // R9 = 3 * incX * sizeof(complex128)
	LEAQ (R10)(R10*2), R11 // R11 = 3 * incY * sizeof(complex128)

dot_loop: // do {
	MOVDDUP_SI__X3      // X_(i+1) = { real(x[i], real(x[i]) }
	MOVDDUP_SI_R8__X5
	MOVDDUP_SI_R8_2__X7
	MOVDDUP_SI_R9__X9

	MOVDDUP_8_SI__X2      // X_i = { imag(x[i]), imag(x[i]) }
	MOVDDUP_8_SI_R8__X4
	MOVDDUP_8_SI_R8_2__X6
	MOVDDUP_8_SI_R9__X8

	// X_j = { imag(y[i]), real(y[i]) }
	MOVUPS (DI), X10
	MOVUPS (DI)(R10*1), X11
	MOVUPS (DI)(R10*2), X12
	MOVUPS (DI)(R11*1), X13

	// X_(i+1) = { imag(a) * real(x[i]), real(a) * real(x[i])  }
	MULPD X10, X3
	MULPD X11, X5
	MULPD X12, X7
	MULPD X13, X9

	// X_j     = { real(y[i]), imag(y[i]) }
	SHUFPD $0x1, X10, X10
	SHUFPD $0x1, X11, X11
	SHUFPD $0x1, X12, X12
	SHUFPD $0x1, X13, X13

	// X_i     = { real(a) * imag(x[i]), imag(a) * imag(x[i])  }
	MULPD X10, X2
	MULPD X11, X4
	MULPD X12, X6
	MULPD X13, X8

	// X_(i+1) = {
	//	imag(result[i]):  imag(a)*real(x[i]) + real(a)*imag(x[i]),
	//	real(result[i]):  real(a)*real(x[i]) - imag(a)*imag(x[i])
	//  }
	ADDSUBPD_X2_X3
	ADDSUBPD_X4_X5
	ADDSUBPD_X6_X7
	ADDSUBPD_X8_X9

	// psum += result[i]
	ADDPD X3, X0
	ADDPD X5, X1
	ADDPD X7, X0
	ADDPD X9, X1

	LEAQ (SI)(R8*4), SI  // SI = &(SI[incX*2])
	LEAQ (DI)(R10*4), DI // DI = &(DI[incY*2])

	DECQ  BX
	JNZ   dot_loop // } while --BX > 0
	ADDPD X1, X0   // sum += psum
	CMPQ  CX, $0   // if CX == 0 { return }
	JE    dot_end

dot_tail: // do {
	MOVDDUP_SI__X3        // X_(i+1) = { real(x[i], real(x[i]) }
	MOVDDUP_8_SI__X2      // X_i = { imag(x[i]), imag(x[i]) }
	MOVUPS (DI), X10      // X_j     = {  imag(y[i])          ,  real(y[i])           }
	MULPD  X10, X3        // X_(i+1) = {  imag(a) * real(x[i]),  real(a) * real(x[i]) }
	SHUFPD $0x1, X10, X10 // X_j     = {  real(y[i])          ,  imag(y[i])           }
	MULPD  X10, X2        // X_i     = {  real(a) * imag(x[i]),  imag(a) * imag(x[i]) }

	// X_(i+1) = {
	//	imag(result[i]):  imag(a)*real(x[i]) + real(a)*imag(x[i]),
	//	real(result[i]):  real(a)*real(x[i]) - imag(a)*imag(x[i])
	//  }
	ADDSUBPD_X2_X3
	ADDPD X3, X0   // sum += result[i]
	ADDQ  R8, SI   // SI += incX
	ADDQ  R10, DI  // DI += incY
	DECQ  CX       // --CX
	JNZ   dot_tail // }  while CX > 0

dot_end:
	MOVUPS X0, sum+88(FP)
	RET

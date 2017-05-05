// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define MOVDDUP_SI_AX_8__X3    LONG $0x1C120FF2; BYTE $0xC6 // + MOVDDUP (SI)(AX*8), X3
#define MOVDDUP_16_SI_AX_8__X5    LONG $0x6C120FF2; WORD $0x10C6 // + MOVDDUP 16(SI)(AX*8), X5
#define MOVDDUP_32_SI_AX_8__X7    LONG $0x7C120FF2; WORD $0x20C6 // + MOVDDUP 32(SI)(AX*8), X7
#define MOVDDUP_48_SI_AX_8__X9    LONG $0x120F44F2; WORD $0xC64C; BYTE $0x30 // + MOVDDUP 48(SI)(AX*8), X9

#define MOVDDUP_SI_DX_8__X2    LONG $0x14120FF2; BYTE $0xD6 // + MOVDDUP (SI)(DX*8), X2
#define MOVDDUP_16_SI_DX_8__X4    LONG $0x64120FF2; WORD $0x10D6 // + MOVDDUP 16(SI)(DX*8), X4
#define MOVDDUP_32_SI_DX_8__X6    LONG $0x74120FF2; WORD $0x20D6 // + MOVDDUP 32(SI)(DX*8), X6
#define MOVDDUP_48_SI_DX_8__X8    LONG $0x120F44F2; WORD $0xD644; BYTE $0x30 // + MOVDDUP 48(SI)(DX*8), X8

#define ADDSUBPD_X2_X3    LONG $0xDAD00F66 // + ADDSUBPD X2, X3
#define ADDSUBPD_X4_X5    LONG $0xECD00F66 // + ADDSUBPD X4, X5
#define ADDSUBPD_X6_X7    LONG $0xFED00F66 // + ADDSUBPD X6, X7
#define ADDSUBPD_X8_X9    LONG $0xD00F4566; BYTE $0xC8 // + ADDSUBPD X8, X9

// func DotcUnitary(x, y []complex128) (sum complex128)
TEXT ·DotcUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+0(FP), SI  // SI = &x
	MOVQ    y_base+24(FP), DI // DI = &y
	MOVQ    x_len+8(FP), CX   // CX = min( len(x), len(y) )
	CMPQ    y_len+32(FP), CX
	CMOVQLE y_len+32(FP), CX
	PXOR    X0, X0            // sum = 0
	CMPQ    CX, $0            // if CX == 0 { return }
	JE      dot_end
	XORPS   X1, X1            // psum = 0
	MOVSD   $(-1.0), X15
	SHUFPD  $0, X15, X15      // { -1, -1 }
	MOVAPS  X15, X14          // Copy X15 to X14 for pipelining
	XORQ    AX, AX            // i := 0
	MOVQ    $1, DX            // j := 1
	MOVQ    CX, BX
	ANDQ    $3, CX            // CX = floor( CX / 4 )
	SHRQ    $2, BX            // BX = CX % 4
	JZ      dot_tail          // if BX == 0 { goto dot_tail }

dot_loop: // do {
	MOVDDUP_SI_AX_8__X3    // X_(i+1) = { real(x[i], real(x[i]) }
	MOVDDUP_16_SI_AX_8__X5
	MOVDDUP_32_SI_AX_8__X7
	MOVDDUP_48_SI_AX_8__X9

	MOVDDUP_SI_DX_8__X2    // X_i = { imag(x[i]), imag(x[i]) }
	MOVDDUP_16_SI_DX_8__X4
	MOVDDUP_32_SI_DX_8__X6
	MOVDDUP_48_SI_DX_8__X8

	// X_i = { -imag(x[i]), -imag(x[i]) }
	MULPD X15, X2
	MULPD X14, X4
	MULPD X15, X6
	MULPD X14, X8

	// X_j = { imag(y[i]), real(y[i]) }
	MOVUPS (DI)(AX*8), X10
	MOVUPS 16(DI)(AX*8), X11
	MOVUPS 32(DI)(AX*8), X12
	MOVUPS 48(DI)(AX*8), X13

	// X_(i+1) = { imag(a) * real(x[i]), real(a) * real(x[i])  }
	MULPD X10, X3
	MULPD X11, X5
	MULPD X12, X7
	MULPD X13, X9

	// X_j = { real(y[i]), imag(y[i]) }
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

	ADDQ  $8, AX   // i += 8
	ADDQ  $8, DX   // j += 8
	DECQ  BX
	JNZ   dot_loop // } while --BX > 0
	ADDPD X1, X0   // sum += psum
	CMPQ  CX, $0   // if CX == 0 { return }
	JE    dot_end

dot_tail: // do {
	MOVDDUP_SI_AX_8__X3    // X_(i+1) = {  real(x[i])          ,  real(x[i])           }
	MOVDDUP_SI_DX_8__X2    // X_i     = {  imag(x[i])          ,  imag(x[i])           }
	MULPD  X15, X2         // X_i     = { -imag(x[i])          , -imag(x[i])           }
	MOVUPS (DI)(AX*8), X10 // X_j     = {  imag(y[i])          ,  real(y[i])           }
	MULPD  X10, X3         // X_(i+1) = {  imag(a) * real(x[i]),  real(a) * real(x[i]) }
	SHUFPD $0x1, X10, X10  // X_j     = {  real(y[i])          ,  imag(y[i])           }
	MULPD  X10, X2         // X_i     = {  real(a) * imag(x[i]),  imag(a) * imag(x[i]) }

	// X_(i+1) = {
	//	imag(result[i]):  imag(a)*real(x[i]) + real(a)*imag(x[i]),
	//	real(result[i]):  real(a)*real(x[i]) - imag(a)*imag(x[i])
	//  }
	ADDSUBPD_X2_X3
	ADDPD X3, X0   // psum += result[i]
	ADDQ  $2, AX   // i += 2
	ADDQ  $2, DX   // j += 2
	DECQ  CX
	JNZ   dot_tail // }  while --CX > 0

dot_end:
	MOVUPS X0, sum+48(FP)
	RET

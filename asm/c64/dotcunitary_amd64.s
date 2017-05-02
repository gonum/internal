// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define MOVSLDUP_SI_AX_8__X3    LONG $0x1C120FF3; BYTE $0xC6 // @ MOVSLDUP (SI)(AX*8), X3
#define MOVSLDUP_16_SI_AX_8__X5    LONG $0x6C120FF3; WORD $0x10C6 // @ MOVSLDUP 16(SI)(AX*8), X5
#define MOVSLDUP_32_SI_AX_8__X7    LONG $0x7C120FF3; WORD $0x20C6 // @ MOVSLDUP 32(SI)(AX*8), X7
#define MOVSLDUP_48_SI_AX_8__X9    LONG $0x120F44F3; WORD $0xC64C; BYTE $0x30 // @ MOVSLDUP 48(SI)(AX*8), X9

#define MOVSHDUP_SI_AX_8__X2    LONG $0x14160FF3; BYTE $0xC6 // @ MOVSHDUP (SI)(AX*8), X2
#define MOVSHDUP_16_SI_AX_8__X4    LONG $0x64160FF3; WORD $0x10C6 // @ MOVSHDUP 16(SI)(AX*8), X4
#define MOVSHDUP_32_SI_AX_8__X6    LONG $0x74160FF3; WORD $0x20C6 // @ MOVSHDUP 32(SI)(AX*8), X6
#define MOVSHDUP_48_SI_AX_8__X8    LONG $0x160F44F3; WORD $0xC644; BYTE $0x30 // @ MOVSHDUP 48(SI)(AX*8), X8

#define MOVSHDUP_X3_X2    LONG $0xD3160FF3 // @ MOVSHDUP X3, X2
#define MOVSLDUP_X3_X3    LONG $0xDB120FF3 // @ MOVSLDUP X3, X3

#define ADDSUBPS_X2_X3    LONG $0xDAD00FF2 // @ ADDSUBPS X2, X3
#define ADDSUBPS_X4_X5    LONG $0xECD00FF2 // @ ADDSUBPS X4, X5
#define ADDSUBPS_X6_X7    LONG $0xFED00FF2 // @ ADDSUBPS X6, X7
#define ADDSUBPS_X8_X9    LONG $0xD00F45F2; BYTE $0xC8 // @ ADDSUBPS X8, X9

// func DotcUnitary(x, y []complex64) (sum complex64)
TEXT ·DotcUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+0(FP), SI  // SI = &x
	MOVQ    y_base+24(FP), DI // DI = &y
	PXOR    X0, X0            // psum = 0
	PXOR    X1, X1            // psum_1 = 0
	MOVQ    x_len+8(FP), CX   // CX = min( len(x), len(y) )
	CMPQ    y_len+32(FP), CX
	CMOVQLE y_len+32(FP), CX
	CMPQ    CX, $0            // if CX == 0 { return }
	JE      dotc_end
	XORQ    AX, AX            // i = 0
	MOVSS   $(-1.0), X15
	SHUFPS  $0, X15, X15      // { -1, -1, -1, -1 }

	MOVQ SI, DX
	ANDQ $15, DX      // DX = &x & 15
	JZ   dotc_aligned // if DX == 0 { goto dotc_aligned }

	MOVSD  (SI)(AX*8), X3  // X_i     = { imag(x[i]), real(x[i]) }
	MOVSHDUP_X3_X2         // X_(i-1) = { imag(x[i]), imag(x[i]) }
	MOVSLDUP_X3_X3         // X_i     = { real(x[i]), real(x[i]) }
	MOVSD  (DI)(AX*8), X10 // X_j     = { imag(y[i]), real(y[i]) }
	MULPS  X15, X2         // X_(i-1) = { -imag(x[i]), imag(x[i]) }
	MULPS  X10, X3         // X_i     = { imag(y[i]) * real(x[i]), real(y[i]) * real(x[i]) }
	SHUFPS $0x1, X10, X10  // X_j     = { real(y[i]), imag(y[i]) }
	MULPS  X10, X2         // X_(i-1) = { real(y[i]) * imag(x[i]), imag(y[i]) * imag(x[i]) }

	// X_i = {
	//	imag(result[i]):  imag(y[i])*real(x[i]) + real(y[i])*imag(x[i]),
	//	real(result[i]):  real(y[i])*real(x[i]) - imag(y[i])*imag(x[i]) }
	ADDSUBPS_X2_X3

	MOVAPS X3, X0   // sum = X_i
	INCQ   AX       // i++
	DECQ   CX       // n--
	JZ     dotc_ret // if CX == 0 { goto dotc_ret }

dotc_aligned:
	MOVQ   CX, BX
	ANDQ   $7, BX    // BX = n % 8
	SHRQ   $3, CX    // CX = floor( n / 8 )
	JZ     dotc_tail // if CX == 0 { return }
	MOVUPS X15, X14  // Copy X15 for pipelining

dotc_loop: // do {
	MOVSLDUP_SI_AX_8__X3    // X_i = { real(x[i]), real(x[i]), real(x[i+1]), real(x[i+1]) }
	MOVSLDUP_16_SI_AX_8__X5
	MOVSLDUP_32_SI_AX_8__X7
	MOVSLDUP_48_SI_AX_8__X9

	MOVSHDUP_SI_AX_8__X2    // X_(i-1) = { imag(x[i]), imag(x[i]), imag(x[i]+1), imag(x[i]+1) }
	MOVSHDUP_16_SI_AX_8__X4
	MOVSHDUP_32_SI_AX_8__X6
	MOVSHDUP_48_SI_AX_8__X8

	// X_j = { imag(y[i]), real(y[i]), imag(y[i+1]), real(y[i+1]) }
	MOVUPS (DI)(AX*8), X10
	MOVUPS 16(DI)(AX*8), X11
	MOVUPS 32(DI)(AX*8), X12
	MOVUPS 48(DI)(AX*8), X13

	// X_(i-1) = { -imag(x[i]), -imag(x[i]), -imag(x[i]+1), -imag(x[i]+1) }
	MULPS X15, X2
	MULPS X14, X4
	MULPS X15, X6
	MULPS X14, X8

	// X_i     = {  imag(y[i])   * real(x[i]),   real(y[i])   * real(x[i]),
	// 		imag(y[i+1]) * real(x[i+1]), real(y[i+1]) * real(x[i+1])  }
	MULPS X10, X3
	MULPS X11, X5
	MULPS X12, X7
	MULPS X13, X9

	// X_j = { real(y[i]), imag(y[i]), real(y[i+1]), imag(y[i+1]) }
	SHUFPS $0xB1, X10, X10
	SHUFPS $0xB1, X11, X11
	SHUFPS $0xB1, X12, X12
	SHUFPS $0xB1, X13, X13

	// X_(i-1) = {  real(y[i])   * imag(x[i]),   imag(y[i])   * imag(x[i]),
	//		real(y[i+1]) * imag(x[i+1]), imag(y[i+1]) * imag(x[i+1])  }
	MULPS X10, X2
	MULPS X11, X4
	MULPS X12, X6
	MULPS X13, X8

	// X_i = {
	//	imag(result[i]):   imag(y[i])   * real(x[i])   + real(y[i])   * imag(x[i]),
	//	real(result[i]):   real(y[i])   * real(x[i])   - imag(y[i])   * imag(x[i]),
	//	imag(result[i+1]): imag(y[i+1]) * real(x[i+1]) + real(y[i+1]) * imag(x[i+1]),
	//	real(result[i+1]): real(y[i+1]) * real(x[i+1]) - imag(y[i+1]) * imag(x[i+1]),
	//  }
	ADDSUBPS_X2_X3
	ADDSUBPS_X4_X5
	ADDSUBPS_X6_X7
	ADDSUBPS_X8_X9

	// psum += X_i
	ADDPS X3, X0
	ADDPS X5, X1
	ADDPS X7, X0
	ADDPS X9, X1

	ADDQ $8, AX    // i += 8
	DECQ CX
	JNZ  dotc_loop // } while --CX > 0

	ADDPS X0, X1 // psum_1 = { psum_0[1] + psum_0[0], psum_0[1] + psum_0[0] }
	XORPS X0, X0 // psum_0 = 0

	CMPQ BX, $0   // if BX == 0 { return }
	JE   dotc_end

dotc_tail:
	MOVQ BX, CX
	SHRQ $1, CX        // CX = floor( CX / 2 )
	JZ   dotc_tail_one // if CX == 0 { goto dotc_tail_one }

dotc_tail_two: // do {
	MOVSLDUP_SI_AX_8__X3   // X_i = { real(x[i]), real(x[i]), real(x[i+1]), real(x[i+1]) }
	MOVSHDUP_SI_AX_8__X2   // X_(i-1) = { imag(x[i]), imag(x[i]), imag(x[i]+1), imag(x[i]+1) }
	MOVUPS (DI)(AX*8), X10 // X_j = { imag(y[i]), real(y[i]) }
	MULPS  X15, X2         // X_(i-1) = { -imag(x[i]), imag(x[i]) }
	MULPS  X10, X3         // X_i = { imag(y[i]) * real(x[i]), real(y[i]) * real(x[i]) }
	SHUFPS $0xB1, X10, X10 // X_j = { real(y[i]), imag(y[i]) }
	MULPS  X10, X2         // X_(i-1) = { real(y[i]) * imag(x[i]), imag(y[i]) * imag(x[i]) }

	// X_i = {
	//	imag(result[i]):  imag(y[i])*real(x[i]) + real(y[i])*imag(x[i]),
	//	real(result[i]):  real(y[i])*real(x[i]) - imag(y[i])*imag(x[i]) }
	ADDSUBPS_X2_X3

	ADDPS X3, X0 // sum += X_i

	ADDQ $2, AX        // i += 2
	DECQ CX
	JNZ  dotc_tail_two // } while --CX > 0

	ADDPS X0, X1 // sum = { psum_0[1] + psum_0[0], psum_0[1] + psum_0[0] }
	XORPS X0, X0 // psum_0 = 0

	ANDQ $1, BX
	JZ   dotc_end

dotc_tail_one:
	MOVSD  (SI)(AX*8), X3  // X_i = { imag(x[i]), real(x[i]) }
	MOVSHDUP_X3_X2         // X_(i-1) = { imag(x[i]), imag(x[i]) }
	MOVSLDUP_X3_X3         // X_i = { real(x[i]), real(x[i]) }
	MOVSD  (DI)(AX*8), X10 // X_j = { imag(y[i]), real(y[i]) }
	MULPS  X15, X2         // X_(i-1) = { -imag(x[i]), imag(x[i]) }
	MULPS  X10, X3         // X_i = { imag(y[i]) * real(x[i]), real(y[i]) * real(x[i]) }
	SHUFPS $0x1, X10, X10  // X_j = { real(y[i]), imag(y[i]) }
	MULPS  X10, X2         // X_(i-1) = { real(y[i]) * imag(x[i]), imag(y[i]) * imag(x[i]) }

	// X_i = {
	//	imag(result[i]):  imag(y[i])*real(x[i]) + real(y[i])*imag(x[i]),
	//	real(result[i]):  real(y[i])*real(x[i]) - imag(y[i])*imag(x[i]) }
	ADDSUBPS_X2_X3

	ADDPS X3, X0 // sum += X_i

dotc_end:
	ADDPS   X1, X0 // X0 = { psum_0[1] + psum_1[1], psum_0[0] + psum_1[0] }
	MOVHLPS X1, X1 // X1 = { psum_0[0], psum_0[1] }
	ADDPS   X1, X0 // X0 = { psum_0[1] + psum_1[1], psum_0[0] + psum_1[0] }

dotc_ret:
	MOVSD X0, sum+48(FP) // return sum
	RET

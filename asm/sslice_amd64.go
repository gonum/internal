// +build !noasm

package asm

// These are function definitions for AMD64 optimized routines,
// and fallback that can be used for performance testing.
// See function documentation in sslice.go

// approx 8x faster than Go
func Div32(out, a, b []float32)

// approx 8x faster than Go
func Add32(out, a, b []float32)

// approx 8x faster than Go
func Mul32(out, a, b []float32)

// approx 8x faster than Go
func Sub32(out, a, b []float32)

// approx 8x faster than Go
func Min32(out, a, b []float32)

// approx 8x faster than Go
func Max32(out, a, b []float32)

// approx Xx faster than Go
func ChangeSign32(out, a, b []float32)

// approx 4x faster than Go
func ConstDiv32(out, a []float32, c float32)

// approx 5x faster than Go
func ConstMul32(out, a []float32, c float32)

// approx 5x faster than Go
func ConstAdd32(out, a []float32, c float32)

// approx 13x faster than Go
func AddScaled32(y, x []float32, a float32)

// approx 11x faster than Go
func Sqrt32(out, a []float32)

// approx 18x faster than Go
func Abs32(out, a []float32)

// approx 15x faster than Go
func MinElement32(a []float32) float32

// approx 15x faster than Go
func MaxElement32(a []float32) float32

// approx 8x faster than Go
func Sum32(a []float32) float32

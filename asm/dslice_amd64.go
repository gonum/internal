// +build amd64 !noasm

package asm

// These are function definitions for AMD64 optimized routines,
// and fallback that can be used for performance testing.
// See function documentation in dslice.go

// approx 2x faster than Go
func Div64(out, a, b []float64)

// approx 3x faster than Go
func Add64(out, a, b []float64)

// approx 3x faster than Go
func Mul64(out, a, b []float64)

// approx 3x faster than Go
func Sub64(out, a, b []float64)

// approx 4x faster than Go
func Min64(out, a, b []float64)

// approx 4x faster than Go
func Max64(out, a, b []float64)

// approx Xx faster than Go
func ChangeSign64(out, a, b []float64)

// approx 2x faster than Go
func ConstDiv64(out, a []float64, c float64)

// approx 3x faster than Go
func ConstMul64(out, a []float64, c float64)

// approx 3x faster than Go
func ConstAdd64(out, a []float64, c float64)

// approx 3x faster than Go
func AddScaled64(y, x []float64, a float64)

// approx 2x faster than Go
func Sqrt64(out, a []float64)

// approx 12x faster than Go
func Abs64(out, a []float64)

// approx 6x faster than Go
func MinElement64(a []float64) float64

// approx 6x faster than Go
func MaxElement64(a []float64) float64

// approx 4x faster than Go
func SliceSum64(a []float64) float64

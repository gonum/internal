// +build !amd64 noasm

package asm

import (
	"math"
)

// These are the fallbacks that are used when not on AMD64 platform.

// Div64 divides two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Div64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = a[i] / b[i]
	}
}

// Add64 adds two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Add64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = a[i] + b[i]
	}
}

// Sub64 subtracts two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Sub64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = a[i] - b[i]
	}
}

// Mul64 multiply two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Mul64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = a[i] * b[i]
	}
}

// Min64 returns lowest valus of two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Min64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		if a[i] < b[i] {
			out[i] = a[i]
		} else {
			out[i] = b[i]
		}
	}
}

// Max64 return maximum of two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func Max64(out, a, b []float64) {
	for i := 0; i < len64(out); i++ {
		if a[i] > b[i] {
			out[i] = a[i]
		} else {
			out[i] = b[i]
		}
	}
}

// ChangeSign64 returns a value with the magnitude of a and the sign of b
// for each element in the slice.
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len64(out)  == len(a) == len(b)
func ChangeSign64(out, a, b []float64) {
	const sign = 1 << 63
	for i := 0; i < len64(out); i++ {
		out[i] = math.Float64frombits(math.Float64bits(a[i])&^sign | math.Float64bits(b[i])&sign)

	}
}

// ConstDiv64 will return c / values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len64(out)  == len(a)
func ConstDiv64(out, a []float64, c float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = c / a[i]
	}
}

// ConstMul64 will return c * values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len64(out)  == len(a)
func ConstMul64(out, a []float64, c float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = c * a[i]
	}
}

// ConstAdd64 will return c * values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len64(out)  == len(a)
func ConstAdd64(out, a []float64, c float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = c + a[i]
	}
}

// AddScaled64 adds a scaled narray elementwise.
// y = y + a * x
// Assumptions the assembly can make:
// y != nil, a != nil
// len(x)  == len(y)
func AddScaled64(y, x []float64, a float64) {
	for i, v := range x {
		y[i] += v * a
	}
}

// Sqrt64 will return math.Sqrt(values) of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len64(out)  == len(a)
func Sqrt64(out, a []float64) {
	for i := 0; i < len64(out); i++ {
		out[i] = float64(math.Sqrt(float64(a[i])))
	}
}

// MinElement64 will the smallest value of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) > 0
func MinElement64(a []float64) float64 {
	min := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

// MaxElement64 will the biggest value of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) > 0
func MaxElement64(a []float64) float64 {
	max := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

// Sum64 will return the sum of all elements of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) >= 0
func Sum64(a []float64) float64 {
	sum := float64(0.0)
	for _, v := range a {
		sum += v
	}
	return sum
}

// Abs64 will return math.Abs(values) of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len64(out)  == len(a)
func Abs64(out, a []float64) {
	for i, v := range a {
		out[i] = float64(math.Abs(float64(v)))
	}
}

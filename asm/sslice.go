// +build !amd64 noasm

package asm

import (
	"math"
)

// These are the fallbacks that are used when not on AMD64 platform.

// Div32 divides two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Div32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		out[i] = a[i] / b[i]
	}
}

// Add32 adds two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Add32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		out[i] = a[i] + b[i]
	}
}

// Sub32 subtracts two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Sub32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		out[i] = a[i] - b[i]
	}
}

// Mul32 multiply two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Mul32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		out[i] = a[i] * b[i]
	}
}

// Min32 returns lowest valus of two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Min32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		if a[i] < b[i] {
			out[i] = a[i]
		} else {
			out[i] = b[i]
		}
	}
}

// Max32 return maximum of two slices
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func Max32(out, a, b []float32) {
	for i := 0; i < len(out); i++ {
		if a[i] > b[i] {
			out[i] = a[i]
		} else {
			out[i] = b[i]
		}
	}
}

// ChangeSign32 returns a value with the magnitude of a and the sign of b
// for each element in the slice.
// Assumptions the assembly can make:
// out != nil, a != nil, b != nil
// len(out)  == len(a) == len(b)
func ChangeSign32(out, a, b []float32) {
	const sign = 1 << 31
	for i := 0; i < len(out); i++ {

		out[i] = math.Float32frombits(math.Float32bits(a[i])&^sign | math.Float32bits(b[i])&sign)
	}
}

// ConstDiv32 will return c / values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len(out)  == len(a)
func ConstDiv32(out, a []float32, c float32) {
	for i := 0; i < len(out); i++ {
		out[i] = c / a[i]
	}
}

// ConstMul32 will return c * values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len(out)  == len(a)
func ConstMul32(out, a []float32, c float32) {
	for i := 0; i < len(out); i++ {
		out[i] = c * a[i]
	}
}

// ConstAdd32 will return c * values of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len(out)  == len(a)
func ConstAdd32(out, a []float32, c float32) {
	for i := 0; i < len(out); i++ {
		out[i] = c + a[i]
	}
}

// AddScaled32 adds a scaled narray elementwise.
// y = y + a * x
// Assumptions the assembly can make:
// y != nil, a != nil
// len(x)  == len(y)
func AddScaled32(y, x []float32, a float32) {
	for i, v := range x {
		y[i] += v * a
	}
}

// Sqrt32 will return math.Sqrt(values) of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len(out)  == len(a)
func Sqrt32(out, a []float32) {
	for i := 0; i < len(out); i++ {
		out[i] = float32(math.Sqrt(float64(a[i])))
	}
}

// MinElement32 will the smallest value of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) > 0
func MinElement32(a []float32) float32 {
	min := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

// MaxElement32 will the biggest value of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) > 0
func MaxElement32(a []float32) float32 {
	max := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

// Sum32 will return the sum of all elements of the slice
// Assumptions the assembly can make:
// a != nil
// len(a) >= 0
func Sum32(a []float32) float32 {
	sum := float32(0.0)
	for _, v := range a {
		sum += v
	}
	return sum
}

// Abs32 will return math.Abs(values) of the array
// Assumptions the assembly can make:
// out != nil, a != nil
// len(out)  == len(a)
func Abs32(out, a []float32) {
	for i, v := range a {
		out[i] = float32(math.Abs(float64(v)))
	}
}

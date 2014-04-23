package fft

import (
	"math"
	"math/cmplx"
)

// Return the smallest power of two >= x
// Source: Section 3-2, Hacker's Delight (2nd Edition), Henry S. Warren, Jr.
func clp2(x uint) uint {
	x = x - 1
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

func FFT(input []float64) (output []complex128) {
	inputCmplx := make([]complex128, clp2(uint(len(input))))
	output = make([]complex128, clp2(uint(len(input))))
	
	for i, x := range input {
		inputCmplx[i] = complex(x, 0)
	}
	
	fft(output, inputCmplx, 1)
	return output
}

func fft(output []complex128, input []complex128, stride int) {
	n := len(output)
	
	if n == 1 {
		output[0] = input[0]
		return
	}
	
	p := n / 2
	fft(output[:p], input, stride * 2)
	fft(output[p:], input[stride:], stride * 2)
	
	i := complex(0, 1)
	for k, t := range output[:p] {
		a := complex(-2 * math.Pi * float64(k) / float64(n), 0)
		e := cmplx.Exp(a * i) * output[k + p]
		output[k] = t + e
		output[k + p] = t - e
	}
}

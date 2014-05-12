package fft

import (
	"math"
	"math/cmplx"
)

// Return the smallest power of two >= x
// Source: Section 3-2, Hacker's Delight (2nd Edition), Henry S. Warren, Jr.
func clp2(x uint64) uint64 {
	x -= 1
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	return x + 1
}

// Fast Fourier transform using Cooley-Tukey algorithm.
func FFT(input []float64) (output []complex128) {
	output = make([]complex128, clp2(uint64(len(input))))
	fft(output, input, 1)
	return output
}

func fft(output []complex128, input []float64, stride int) {
	n := len(output)
	p := n / 2

	if n == 1 {
		if len(input) > 0 {
			output[0] = complex(input[0], 0)
		} else {
			output[0] = complex(0, 0)
		}
		return
	}

	var oddInput []float64
	if len(input) > stride {
		oddInput = input[stride:]
	}

	fft(output[:p], input, stride*2)
	fft(output[p:], oddInput, stride*2)

	for k, t := range output[:p] {
		a := complex(0, -2*math.Pi*float64(k)/float64(n))
		e := cmplx.Exp(a) * output[k+p]
		output[k] = t + e
		output[k+p] = t - e
	}
}

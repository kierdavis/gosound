package fft

import (
	"math"
)

// Rolling short-time Fourier transform over an input channel.
func STFT(input chan float64, window []float64, overlapSize int) (output chan []float64) {
	output = make(chan []float64)

	go func() {
		defer close(output)

		windowSize := len(window)
		buffer := make([]float64, windowSize)
		overlap := make([]float64, overlapSize)
		pos := 0

		ok := true
		for ok {
			// Top up the buffer.
			for pos < windowSize {
				buffer[pos], ok = <-input
				pos++
			}
			// Only check if the channel is closed after the end of the loop,
			// since we would want to fill the rest of the buffer with zeroes
			// and do one final FFT anyway.
			// When ok is set to false, the main loop will end.

			// Save the overlap.
			copy(overlap, buffer[windowSize-overlapSize:])

			// Apply the windowing function to the input.
			for i := 0; i < windowSize; i++ {
				buffer[i] *= window[i]
			}

			// FFT it
			cresult := FFT(buffer)

			// Find the squared magnitudes and put them into a result array.
			result := make([]float64, len(cresult))
			for i, c := range cresult {
				x, y := real(c), imag(c)
				result[i] = x*x + y*y
			}

			// Send the result
			output <- result

			// Copy the overlap to the start of the buffer.
			copy(buffer, overlap)
			pos = overlapSize
		}
	}()

	return output
}

func RectangularWindow(size int) (window []float64) {
	window = make([]float64, size)
	for i := 0; i < size; i++ {
		window[i] = 1.0
	}
	return window
}

func TriangularWindow(size int) (window []float64) {
	window = make([]float64, size)
	p := (size - 1) / 2
	for i := 0; i <= p; i++ {
		window[i] = float64(i) / float64(p)
	}
	for i := p + 1; i < size; i++ {
		window[i] = 1.0 - float64(i-p)/float64(p)
	}
	return window
}

func WelchWindow(size int) (window []float64) {
	window = make([]float64, size)
	p := (size - 1) / 2
	for i := 0; i < size; i++ {
		a := float64(i-p) / float64(p+1)
		window[i] = 1.0 - a*a
	}
	return window
}

func CosineWindow(size int, a, b float64) (window []float64) {
	window = make([]float64, size)
	m := (math.Pi * 2.0) / (float64(size) - 1)
	for i := 0; i < size; i++ {
		window[i] = a - b*math.Cos(m*float64(i))
	}
	return window
}

func HanningWindow(size int) (window []float64) {
	return CosineWindow(size, 0.5, 0.5)
}

func HammingWindow(size int) (window []float64) {
	return CosineWindow(size, 0.54, 0.46)
}

// s is the standard deviation
func GaussianWindow(size int, s float64) (window []float64) {
	window = make([]float64, size)
	p := float64(size-1) / 2
	for i := 0; i < size; i++ {
		a := (float64(i) - p) / (s * p)
		window[i] = math.Exp(-0.5 * a * a)
	}
	return window
}

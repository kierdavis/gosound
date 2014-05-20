package filter

import (
	"github.com/kierdavis/gosound/sound"
	"math"
)

// A simple R/C highpass filter.
// Based on https://en.wikipedia.org/wiki/High-pass_filter#Algorithmic_implementation
func HighPassRCFilter(ctx sound.Context, input chan float64, cutoffFreq float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(output)
		
		// Time interval between samples
		dt := 1 / ctx.SampleRate

		// Time constant (of analogue circuit)
		rc := 1 / (2 * math.Pi * cutoffFreq)
		
		alpha := rc / (rc + dt)
		
		lastX := <-input
		lastY := lastX
		output <- lastY

		for x := range input {
			y := alpha * (lastY + x - lastX)
			output <- y

			lastX = x
			lastY = y
		}
	}()

	return output
}

// A simple R/C lowpass filter.
// Based on https://en.wikipedia.org/wiki/Low-pass_filter#Simple_infinite_impulse_response_filter
func LowPassRCFilter(ctx sound.Context, input chan float64, cutoffFreq float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(output)
		
		// Time interval between samples
		dt := 1 / ctx.SampleRate

		// Time constant (of analogue circuit)
		rc := 1 / (2 * math.Pi * cutoffFreq)
		
		alpha := dt / (rc + dt)

		lastY := <-input
		output <- lastY

		for x := range input {
			y := lastY + alpha*(x-lastY)
			output <- y

			lastY = y
		}
	}()

	return output
}


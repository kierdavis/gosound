package filter

import (
	"github.com/kierdavis/gosound/sound"
	"math"
)

// A simple analogue RC filter.
// Based on https://en.wikipedia.org/wiki/High-pass_filter#Algorithmic_implementation
// and https://en.wikipedia.org/wiki/Low-pass_filter#Simple_infinite_impulse_response_filter
func RC(ctx sound.Context, input chan float64, filterType FilterType, cutoffFreqInput chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		// Time interval between samples
		dt := 1.0 / ctx.SampleRate
		
		switch filterType {
		case LowPass:
			lastY := <-input
			output <- lastY

			for x := range input {
				cutoffFreq := <-cutoffFreqInput
				rc := 1.0 / (2.0 * math.Pi * cutoffFreq) // Time constant (of analogue circuit)
				alpha := dt / (rc + dt)
				
				y := lastY + alpha*(x-lastY)
				output <- y

				lastY = y
			}
		
		case HighPass:
			lastX := <-input
			lastY := lastX
			output <- lastY

			for x := range input {
				cutoffFreq := <-cutoffFreqInput
				rc := 1.0 / (2.0 * math.Pi * cutoffFreq) // Time constant (of analogue circuit)
				alpha := rc / (rc + dt)
				
				y := alpha * (lastY + x - lastX)
				output <- y

				lastX = x
				lastY = y
			}
		}
	}()

	return output
}

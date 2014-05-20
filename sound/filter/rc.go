package filter

import (
	"github.com/kierdavis/gosound/sound"
	"math"
)

// A simple analogue RC filter.
// Based on https://en.wikipedia.org/wiki/High-pass_filter#Algorithmic_implementation
// and https://en.wikipedia.org/wiki/Low-pass_filter#Simple_infinite_impulse_response_filter
func RC(ctx sound.Context, input chan float64, filterType FilterType, cutoffFreq float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		// Time interval between samples
		dt := 1 / ctx.SampleRate

		// Time constant (of analogue circuit)
		rc := 1 / (2 * math.Pi * cutoffFreq)
		
		switch filterType {
		case LowPass:
			alpha := dt / (rc + dt)

			lastY := <-input
			output <- lastY

			for x := range input {
				y := lastY + alpha*(x-lastY)
				output <- y

				lastY = y
			}
		
		case HighPass:
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
		}
	}()

	return output
}

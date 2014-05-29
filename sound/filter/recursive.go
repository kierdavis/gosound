package filter

import (
	"github.com/kierdavis/gosound/sound"
)

// Run a recursive filter using input coefficients 'as' and past output
// coefficients 'bs'. 'bs[0]' is ignored.
func Recursive(ctx sound.Context, input chan float64, as, bs []chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		prevInputs := make([]float64, len(as))
		prevOutputs := make([]float64, len(bs))
		prevInputPtr := 0 // points to most recently added prevInput
		prevOutputPtr := 0 // points to most recently added prevOutput
		
		for x := range input {
			y := (<-as[0]) * x
			
			for i, a := range as[1:] {
				y += (<-a) * prevInputs[(prevInputPtr - i + len(prevInputs)) % len(prevInputs)]
			}
			
			for i, b := range bs[1:] {
				y += (<-b) * prevOutputs[(prevOutputPtr - i + len(prevOutputs)) % len(prevOutputs)]
			}
			
			output <- y
			
			prevInputPtr = (prevInputPtr + 1) % len(prevInputs)
			prevInputs[prevInputPtr] = x
			prevOutputPtr = (prevOutputPtr + 1) % len(prevOutputs)
			prevOutputs[prevOutputPtr] = y
		}
	}()
	
	return output
}

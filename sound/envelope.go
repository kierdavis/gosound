package sound

import (
	"fmt"
	"time"
)

// Example:
//     adsrEnvelope = LinearEnvelope(
//       0.0,               // Start level
//       time.Second / 10   // Attack time
//       1.0,               // Peak level
//       time.Second / 10   // Decay time
//       0.2,               // Sustain level
//       time.Second * 2    // Sustain time
//       0.2,               // Sustain level
//       time.Second / 2    // Release time
//       0.0,               // End level
//     )
// Will panic of args are not of the correct types (alternative float64 and time.Duration).
func (ctx Context) LinearEnvelope(args ...interface{}) (output chan float64) {
	if len(args) % 2 != 1 {
		panic("Bad number of arguments")
	}
	
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		x, ok := args[0].(float64)
		if !ok {
			panic("Expected argument 0 to be of type float64")
		}
		
		for i := 1; i < len(args); i += 2 {
			duration, ok := args[i].(time.Duration)
			if !ok {
				panic(fmt.Sprintf("Expected argument %d to be of type time.Duration", i))
			}
			
			y, ok := args[i + 1].(float64)
			if !ok {
				panic(fmt.Sprintf("Expected argument %d to be of type float64", i + 1))
			}
			
			numSamples := (float64(duration) / float64(time.Second)) * ctx.SampleRate
			incr := 1 / numSamples
			
			// Interpolate from x to y across numSamples samples
			for f := 0.0; f < 1.0; f += incr {
				output <- x * (1 - f) + y * f
			}
			
			x = y
		}
		
		
	}()
	
	return output
}

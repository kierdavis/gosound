package filter

import (
	"github.com/kierdavis/gosound/sound"
)

// Run a recursive filter using input coefficients 'as' and past output
// coefficients 'bs'.
// 
// The coefficient of the current input sample should be placed in a[0], that of
// the previous input sample in a[1] and so on. The coefficient of the previous
// output sample should be placed in b[0], that of the output sample before that
// in b[1] and so on.
func Recursive(ctx sound.Context, input chan float64, as, bs []float64) (output chan float64) {
    output = make(chan float64, ctx.StreamBufferSize)
    
    go func() {
        defer close(output)
        
        a0 := 0.0
        if len(as) >= 1 {
            a0 = as[0]
            as = as[1:]
        }
        
        prevInputs := make([]float64, len(as))
        prevOutputs := make([]float64, len(bs))
        prevInputPtr := 0 // points to most recently added prevInput
        prevOutputPtr := 0 // points to most recently added prevOutput
        
        for x := range input {
            y := a0 * x
            
            for i, a := range as {
                y += a * prevInputs[(prevInputPtr - i + len(prevInputs)) % len(prevInputs)]
            }
            
            for i, b := range bs {
                y += b * prevOutputs[(prevOutputPtr - i + len(prevOutputs)) % len(prevOutputs)]
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

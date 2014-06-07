package karplusstrong

import (
    "github.com/kierdavis/gosound/sound"
    "github.com/kierdavis/gosound/sound/filter"
)

// Copy input to output
func pipe(input, output chan float64) {
    for x := range input {
        output <- x
    }
    close(output)
}

// Delay line suitable for use in feedback systems without causing deadlock.
func delay(ctx sound.Context, input chan float64, length uint) (output chan float64) {
    output = make(chan float64, ctx.StreamBufferSize)
    
    go func() {
        buffer := make([]float64, length)
        pos := uint(0)
        
        for {
            output <- buffer[pos]
            buffer[pos] = <-input
            pos = (pos + 1) % length
        }
    }()
    
    return output
}

// input should be finite and preferably short (e.g. one cycle of a triangle wave)
func KarplusStrong(ctx sound.Context, input chan float64, delaySamples uint, cutoff float64, decay float64) (output chan float64) {
    feedback := make(chan float64, ctx.StreamBufferSize)
    
    // Mix the input with the feedback.
    output = ctx.Add(input, feedback)
    
    // Fork off a copy of the output.
    output, outputCopy := ctx.Fork2(output)
    
    // The copy is first passed through a delay line...
    outputCopy = delay(ctx, outputCopy, delaySamples)
    
    // ...then filtered...
    //outputCopy = filter.Chebyshev(ctx, outputCopy, filter.LowPass, cutoff, 0.5, 2)
    outputCopy = filter.RC(ctx, outputCopy, filter.LowPass, ctx.Const(cutoff))
    
    // ...and finally attenuated slightly.
    outputCopy = ctx.Mul(outputCopy, ctx.Const(decay))
    
    // The filtered output copy is fed back into the system.
    go pipe(outputCopy, feedback)
    
    return output
}

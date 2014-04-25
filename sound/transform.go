package sound

import (
	"math"
	"time"
)

// Streams that are not read from will cause the writing end to block! Use this
// to drop samples that aren't needed.
func (ctx Context) Drain(input chan float64) {
	go func() {
		for _ = range input {}
	}()
}

func (ctx Context) Mix(inputs ...chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		for {
			var total float64
			for _, input := range inputs {
				total += <-input
			}
			output <- total
		}
	}()
	
	return output
}

func (ctx Context) MixNormalised(inputs ...chan float64) (output chan float64) {
	return ctx.Scale(ctx.Mix(inputs...), 1 / float64(len(inputs)))
}

func (ctx Context) Append(inputs ...chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for _, input := range inputs {
			for x := range input {
				output <- x
			}
		}
	}()
	
	return output
}

func (ctx Context) SplitAt(input chan float64, length time.Duration, waitForZeroCrossing bool) (before, after chan float64) {
	before = make(chan float64, ctx.StreamBufferSize)
	after = make(chan float64, ctx.StreamBufferSize)
	
	numSamples := uint((float64(length) / float64(time.Second)) * ctx.SampleRate)
	
	go func() {
		var x float64
		
		for x = range input {
			numSamples--
			if numSamples == 0 {
				break
			}
			
			before <- x
		}
		
		if waitForZeroCrossing {
			for x > 1e6 || x < -1e6 {
				x = <-input
				before <- x
			}
		}
		
		close(before)
		
		for x = range input {
			after <- x
		}
		
		close(after)
	}()
	
	return before, after
}

func (ctx Context) Take(input chan float64, length time.Duration, waitForZeroCrossing bool) (output chan float64) {
	before, _ := ctx.SplitAt(input, length, waitForZeroCrossing)
	// If the input is an infinite source then draining it will consume a lot of resources
	//ctx.Drain(after)
	return before
}

func (ctx Context) Drop(input chan float64, length time.Duration, waitForZeroCrossing bool) (output chan float64) {
	before, after := ctx.SplitAt(input, length, waitForZeroCrossing)
	ctx.Drain(before)
	return after
}

func (ctx Context) ToBuffer(input chan float64) (buffer []float64) {
	buffer = make([]float64, 0, ctx.StreamBufferSize)
	
	for x := range input {
		buffer = append(buffer, x)
	}
	
	return buffer
}

func (ctx Context) Offset(input chan float64, offset float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			output <- x + offset
		}
	}()
	
	return output
}

func (ctx Context) OffsetM(input chan float64, offsetModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			offset := <-offsetModulation
			output <- x + offset
		}
	}()
	
	return output
}

func (ctx Context) Scale(input chan float64, scalar float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			output <- x * scalar
		}
	}()
	
	return output
}

func (ctx Context) ScaleM(input chan float64, scalarModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			scalar := <-scalarModulation
			output <- x * scalar
		}
	}()
	
	return output
}

func (ctx Context) Clip(input chan float64, dist float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			output <- math.Max(math.Min(x, dist), -dist)
		}
	}()
	
	return output
}

func (ctx Context) ClipM(input chan float64, distModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			dist := <-distModulation
			output <- math.Max(math.Min(x, dist), -dist)
		}
	}()
	
	return output
}

func (ctx Context) Fork(input chan float64, numOutputs int) (outputs []chan float64) {
	outputs = make([]chan float64, numOutputs)
	for i, _ := range outputs {
		outputs[i] = make(chan float64, ctx.StreamBufferSize)
	}
	
	go func() {
		for x := range input {
			for _, output := range outputs {
				output <- x
			}
		}
		
		for _, output := range outputs {
			close(output)
		}
	}()
	
	return outputs
}

// floating-point GCD. a and b are assumed to be positive.
func fgcd(a, b float64) (gcd float64) {
	if b < 1e-9 {
		return a
	}
	
	return fgcd(b, math.Mod(a, b))
}

func flcm(a, b float64) (lcm float64) {
	return (a * b) / fgcd(a, b)
}

func (ctx Context) Resample(input chan float64, newRate float64) (output chan float64, newCtx Context) {
	output = make(chan float64, ctx.StreamBufferSize)
	intermediate := make(chan float64, ctx.StreamBufferSize)
	
	oldRate := ctx.SampleRate
	intermediateRate := flcm(oldRate, newRate)
	
	// Expand phase: interpolate between input samples to produce an intermediate form at the intermediate sample rate
	go func() {
		defer close(intermediate)
		
		x := <-input
		f := 0.0 // Interpolation factor
		fIncr := oldRate / intermediateRate
		
		for y := range input {
			for f < 1.0 {
				intermediate <- x*(1-f) + y*f
				f += fIncr
			}
			f -= 1.0
			x = y
		}
	}()
	
	// Compress phase: select samples of the intermediate form to go into the output channel
	go func() {
		defer close(output)
		
		f := 0.0
		fIncr := newRate / intermediateRate
		
		for x := range intermediate {
			f += fIncr
			if f >= 1.0 {
				output <- x
				f -= 1.0
			}
		}
	}()
	
	newCtx = ctx
	newCtx.SampleRate = newRate
	
	return output, newCtx
}

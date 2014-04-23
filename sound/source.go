package sound

import (
	"math"
	"math/rand"
)

// Return a stream of values that rise from 0 to 1, 'frequency' times per second.
func (ctx Context) WaveInput(frequency float64) (stream chan float64) {
	stream = make(chan float64, ctx.StreamBufferSize)
	
	incr := frequency / ctx.SampleRate
	
	go func() {
		var x float64
		
		for {
			stream <- x
			x += incr
			if x >= 1 {
				x -= 1
			}
		}
	}()
	
	return stream
}

// Returns a sine wave at a given frequency.
func (ctx Context) Sine(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.WaveInput(frequency)
	
	go func() {
		for x := range input {
			output <- float64(math.Sin(float64(x * 2 * math.Pi)))
		}
	}()
	
	return output
}

// Returns a saw wave at a given frequency.
func (ctx Context) Saw(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.WaveInput(frequency)
	
	go func() {
		for x := range input {
			output <- x*2 - 1
		}
	}()
	
	return output
}

// Returns a square wave at a given frequency.
func (ctx Context) Square(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.WaveInput(frequency)
	
	go func() {
		for x := range input {
			// dear God I hope this conditional is optimised away...
			if x >= 0.5 {
				output <- 1.0
			} else {
				output <- -1.0
			}
		}
	}()
	
	return output
}

func (ctx Context) Silence() (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		for {
			output <- 0.0
		}
	}()
	
	return output
}

func (ctx Context) RandomNoise(seed int64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	r := rand.New(rand.NewSource(seed))
	
	go func() {
		for {
			output <- r.Float64()
		}
	}()
	
	return output
}

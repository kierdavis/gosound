package sound

import (
	"math"
	"math/rand"
)

func (ctx Context) Saw(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	incr := frequency / ctx.SampleRate
	
	go func() {
		x := 0.0
		
		for {
			output <- (x * 2.0) - 1.0
			_, x = math.Modf(x + incr)
		}
	}()
	
	return output
}

func (ctx Context) Triangle(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.Saw(frequency)
	
	go func() {
		for x := range input {
			output <- math.Abs(x)
		}
	}()
	
	return output
}

func (ctx Context) Square(frequency float64, duty float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	duty = (duty - 0.5) * 2 // Range [0..1] -> [-1..1]
	
	input := ctx.Saw(frequency)
	
	go func() {
		for x := range input {
			if x < duty {
				output <- 1.0
			} else {
				output <- -1.0
			}
		}
	}()
	
	return output
}

func (ctx Context) Sine(frequency float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.Saw(frequency)
	
	go func() {
		for x := range input {
			output <- math.Sin(x * math.Pi)
		}
	}()
	
	return output
}

/*
// Return a stream of values that rise from 0 to 1, 'frequency' times per second.
// Phase is in the range [0, 2pi)
func (ctx Context) WaveInput(frequency float64, phase float64) (stream chan float64) {
	stream = make(chan float64, ctx.StreamBufferSize)
	
	incr := frequency / ctx.SampleRate
	
	go func() {
		x := phase / (math.Pi * 2)
		
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
func (ctx Context) SineWithPhase(frequency float64, phase float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.WaveInput(frequency, phase)
	
	go func() {
		for x := range input {
			output <- float64(math.Sin(float64(x * 2 * math.Pi)))
		}
	}()
	
	return output
}

func (ctx Context) Sine(frequency float64) (output chan float64) {
	return ctx.SineWithPhase(frequency, 0.0)
}

// Returns a saw wave at a given frequency.
func (ctx Context) SawWithPhase(frequency float64, phase float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	// Add pi to phase so that the output it starts at 0
	input := ctx.WaveInput(frequency, phase + math.Pi)
	
	go func() {
		for x := range input {
			output <- x*2 - 1
		}
	}()
	
	return output
}

func (ctx Context) Saw(frequency float64) (output chan float64) {
	return ctx.SawWithPhase(frequency, 0.0)
}

// Returns a triangle wave at a given frequency.
func (ctx Context) TriangleWithPhase(frequency float64, phase float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	// Add pi/2 to phase so that the output it starts at 0
	input := ctx.WaveInput(frequency, phase + math.Pi/2)
	
	go func() {
		for x := range input {
			output <- 1 - math.Abs(x*2 - 1)*2
		}
	}()
	
	return output
}

func (ctx Context) Triangle(frequency float64) (output chan float64) {
	return ctx.TriangleWithPhase(frequency, 0.0)
}

func signum(x float64) float64 {
	return x / math.Abs(x)
}

// Returns a square wave at a given frequency.
func (ctx Context) SquareWithPhase(frequency float64, phase float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	input := ctx.WaveInput(frequency, phase)
	
	go func() {
		for x := range input {
			output <- signum(x - 0.5)
		}
	}()
	
	return output
}

func (ctx Context) Square(frequency float64) (output chan float64) {
	return ctx.SquareWithPhase(frequency, 0.0)
}
*/

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

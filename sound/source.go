package sound

import (
	"math"
	"math/rand"
)

func (ctx Context) Saw(frequencyInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(signalOutput)
		
		x := 0.5
		
		for frequency := range frequencyInput {
			signalOutput <- (x * 2.0) - 1.0
			x = math.Mod(x + (frequency / ctx.SampleRate), 1.0)
		}
	}()
	
	return signalOutput
}

func (ctx Context) Triangle(frequencyInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)
	
	saw := ctx.Saw(frequencyInput)
	
	go func() {
		defer close(signalOutput)
		
		for x := range saw {
			signalOutput <- math.Abs(x)
		}
	}()
	
	return signalOutput
}

func (ctx Context) Square(frequencyInput chan float64, dutyInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)
	
	saw := ctx.Saw(frequencyInput)
	
	go func() {
		defer close(signalOutput)
		
		for x := range saw {
			duty := ((<-dutyInput) - 0.5) * 2
			if x < duty {
				signalOutput <- 1.0
			} else {
				signalOutput <- -1.0
			}
		}
	}()
	
	return signalOutput
}

func (ctx Context) Sine(frequencyInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)
	
	saw := ctx.Saw(frequencyInput)
	
	go func() {
		defer close(signalOutput)
		
		for x := range saw {
			signalOutput <- math.Sin(x * math.Pi)
		}
	}()
	
	return signalOutput
}

func (ctx Context) Silence() (output chan float64) {
	return ctx.Const(0)
}

func (ctx Context) RandomNoise(seed int64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	r := rand.New(rand.NewSource(seed))
	
	go func() {
		for {
			output <- r.Float64() * 2.0 - 1.0
		}
	}()
	
	return output
}

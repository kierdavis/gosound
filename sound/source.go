package sound

import (
	"math"
	"math/rand"
)

// SawWithPhase produces a sawtooth wave with frequency modulated by
// 'frequencyInput' and initial phase 'phase'. 'phase' lies in the interval
// [0,1] where a phase of 0 indicates the signal is about to ascend from 0.
func (ctx Context) SawWithPhase(frequencyInput chan float64, phase float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(signalOutput)

		x := math.Mod(phase+0.5, 1.0)

		for frequency := range frequencyInput {
			signalOutput <- x*2.0 - 1.0
			x = math.Mod(x+(frequency/ctx.SampleRate), 1.0)
		}
	}()

	return signalOutput
}

// Saw produces a sawtooth wave with frequency modulated by 'frequencyInput'.
func (ctx Context) Saw(frequencyInput chan float64) (signalOutput chan float64) {
	return ctx.SawWithPhase(frequencyInput, 0.0)
}

// TriangleWithPhase produces a triangle wave with frequency modulated by
// 'frequencyInput' and initial phase 'phase'. 'phase' lies in the interval
// [0,1] where a phase of 0 indicates the signal is about to asend from 0.
func (ctx Context) TriangleWithPhase(frequencyInput chan float64, phase float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	saw := ctx.SawWithPhase(frequencyInput, phase+0.25)

	go func() {
		defer close(signalOutput)

		for x := range saw {
			signalOutput <- math.Abs(x)*2.0 - 1.0
		}
	}()

	return signalOutput
}

// Triangle produces a triangle wave with frequency modulated by
// 'frequencyInput'.
func (ctx Context) Triangle(frequencyInput chan float64) (signalOutput chan float64) {
	return ctx.TriangleWithPhase(frequencyInput, 0.0)
}

// SquareWithPhase produces a square wave with frequency modulated by
// 'frequencyInput' and initial phase 'phase'. 'phase' lies in the interval
// [0,1] where a phase of 0.25 indicates the signal is transitioning from -1 to
// 1. The duty cycle of the wave is modulated by 'dutyInput'.
func (ctx Context) SquareWithPhase(frequencyInput chan float64, dutyInput chan float64, phase float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	saw := ctx.SawWithPhase(frequencyInput, phase)

	go func() {
		defer close(signalOutput)

		for x := range saw {
			threshold := ((<-dutyInput) - 0.5) * 2
			if x < threshold {
				signalOutput <- 1.0
			} else {
				signalOutput <- -1.0
			}
		}
	}()

	return signalOutput
}

// Square produces a square wave with frequency modulated by 'frequencyInput'.
// The duty cycle of the wave is modulated by 'dutyInput'.
func (ctx Context) Square(frequencyInput chan float64, dutyInput chan float64) (signalOutput chan float64) {
	return ctx.SquareWithPhase(frequencyInput, dutyInput, 0.0)
}

func (ctx Context) SineWithPhase(frequencyInput chan float64, phase float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	saw := ctx.SawWithPhase(frequencyInput, phase)

	go func() {
		defer close(signalOutput)

		for x := range saw {
			signalOutput <- math.Sin(x * math.Pi)
		}
	}()

	return signalOutput
}

func (ctx Context) Sine(frequencyInput chan float64) (signalOutput chan float64) {
	return ctx.SineWithPhase(frequencyInput, 0.0)
}

func (ctx Context) Silence() (output chan float64) {
	return ctx.Const(0)
}

func (ctx Context) RandomNoise(seed int64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	r := rand.New(rand.NewSource(seed))

	go func() {
		for {
			output <- r.Float64()*2.0 - 1.0
		}
	}()

	return output
}

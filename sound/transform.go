package sound

import (
	"math"
	"time"
)

func (ctx Context) Clip(signalInput chan float64, distInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(signalOutput)

		for x := range signalInput {
			dist := math.Abs(<-distInput)
			signalOutput <- math.Max(math.Min(x, dist), -dist)
		}
	}()

	return signalOutput
}

func (ctx Context) Resample(input chan float64, ratio float64) (output chan float64, newCtx Context) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(output)

		x := <-input
		f := 0.0
		fIncr := ratio

		for y := range input {
			for f < 1.0 {
				output <- x*(1-f) + y*f
				f += fIncr
			}
			f -= 1.0
			x = y
		}
	}()

	newCtx = ctx
	newCtx.SampleRate *= ratio

	return output, newCtx
}

func (ctx Context) ModulateFrequency(input chan float64, ratio float64) (output chan float64) {
	// An 'f' Hz signal at a sample rate of 'sr' Hz is equivalent to a
	// 'k * f' Hz signal at a sample rate of '1/k * sr' Hz.
	// So we can simply divide the sample rate by 'ratio' and call it the output
	// signal.
	ctx.SampleRate /= ratio

	// However, we want to return a result at the original sample rate. So, let's
	// resample back. Multiplying the sample rate by 'ratio' will return it to
	// its original value.
	output, ctx = ctx.Resample(input, ratio)
	return output
}

func (ctx Context) SplitAtDuration(input chan float64, duration time.Duration, waitForZC bool) (beforeOutput, afterOutput chan float64) {
	count := uint(ctx.SampleRate * float64(duration) / float64(time.Second))
	return ctx.SplitAt(input, count, waitForZC)
}

func (ctx Context) TakeDuration(input chan float64, duration time.Duration, waitForZC bool) (output chan float64) {
	count := uint(ctx.SampleRate * float64(duration) / float64(time.Second))
	return ctx.Take(input, count, waitForZC)
}

func (ctx Context) DropDuration(input chan float64, duration time.Duration, waitForZC bool) (output chan float64) {
	count := uint(ctx.SampleRate * float64(duration) / float64(time.Second))
	return ctx.Drop(input, count, waitForZC)
}

func (ctx Context) Duration(input chan float64) (duration time.Duration) {
	return time.Duration(float64(time.Second) * float64(ctx.Count(input)) / ctx.SampleRate)
}

package sound

import (
	"math"
)

// Convert a high pass filter's cutoff frequency to an alpha value
func (ctx Context) HighPassCutoffToAlpha(frequency float64) (alpha float64) {
	// Time interval between samples
	dt := 1 / ctx.SampleRate

	// Time constant (of analogue circuit)
	rc := 1 / (2 * math.Pi * frequency)

	return rc / (rc + dt)
}

func (ctx Context) HighPassCutoffsToAlphas(frequencyInput chan float64) (alphaOutput chan float64) {
	alphaOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(alphaOutput)

		for frequency := range frequencyInput {
			alphaOutput <- ctx.HighPassCutoffToAlpha(frequency)
		}
	}()

	return alphaOutput
}

func (ctx Context) HighPass(signalInput chan float64, alphaInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(signalOutput)

		lastX := <-signalInput
		lastY := lastX
		signalOutput <- lastY

		for x := range signalInput {
			alpha := <-alphaInput
			y := alpha * (lastY + x - lastX)
			signalOutput <- y

			lastX = x
			lastY = y
		}
	}()

	return signalOutput
}

// Convert a low pass filter's cutoff frequency to an alpha value
func (ctx Context) LowPassCutoffToAlpha(frequency float64) (alpha float64) {
	// Time interval between samples
	dt := 1 / ctx.SampleRate

	// Time constant (of analogue circuit)
	rc := 1 / (2 * math.Pi * frequency)

	return dt / (rc + dt)
}

func (ctx Context) LowPassCutoffsToAlphas(frequencyInput chan float64) (alphaOutput chan float64) {
	alphaOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(alphaOutput)

		for frequency := range frequencyInput {
			alphaOutput <- ctx.LowPassCutoffToAlpha(frequency)
		}
	}()

	return alphaOutput
}

func (ctx Context) LowPass(signalInput chan float64, alphaInput chan float64) (signalOutput chan float64) {
	signalOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(signalOutput)

		lastY := <-signalInput
		signalOutput <- lastY

		for x := range signalInput {
			alpha := <-alphaInput
			y := lastY + alpha*(x-lastY)
			signalOutput <- y

			lastY = y
		}
	}()

	return signalOutput
}

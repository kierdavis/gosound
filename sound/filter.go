package sound

import (
	"math"
)

// Convert a high pass filter's cutoff frequency to an alpha value
func (ctx Context) HighPassCutoffToAlpha(freq float64) (alpha float64) {
	// Time interval between samples
	dt := 1 / ctx.SampleRate
	
	// Time constant (of analogue circuit)
	rc := 1 / (2 * math.Pi * freq)
	
	return rc / (rc + dt)
}

func (ctx Context) HighPassCutoffToAlphaM(freqs chan float64) (alphas chan float64) {
	alphas = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(alphas)
		
		for freq := range freqs {
			alphas <- ctx.HighPassCutoffToAlpha(freq)
		}
	}()
	
	return alphas
}

func (ctx Context) HighPass(input chan float64, alpha float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastX := <-input
		lastY := lastX
		output <- lastY
		
		for x := range input {
			y := alpha * (lastY + x - lastX)
			output <- y
			
			lastX = x
			lastY = y
		}
	}()
	
	return output
}

func (ctx Context) HighPassM(input chan float64, alphaModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastX := <-input
		lastY := lastX
		output <- lastY
		
		for x := range input {
			alpha := <-alphaModulation
			y := alpha * (lastY + x - lastX)
			output <- y
			
			lastX = x
			lastY = y
		}
	}()
	
	return output
}

// Convert a low pass filter's cutoff frequency to an alpha value
func (ctx Context) LowPassCutoffToAlpha(freq float64) (alpha float64) {
	// Time interval between samples
	dt := 1 / ctx.SampleRate
	
	// Time constant (of analogue circuit)
	rc := 1 / (2 * math.Pi * freq)
	
	return dt / (rc + dt)
}

func (ctx Context) LowPassCutoffToAlphaM(freqs chan float64) (alphas chan float64) {
	alphas = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(alphas)
		
		for freq := range freqs {
			alphas <- ctx.LowPassCutoffToAlpha(freq)
		}
	}()
	
	return alphas
}

func (ctx Context) LowPass(input chan float64, alpha float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastY := <-input
		output <- lastY
		
		for x := range input {
			y := lastY + alpha * (x - lastY)
			output <- y
			
			lastY = y
		}
	}()
	
	return output
}

func (ctx Context) LowPassM(input chan float64, alphaModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastY := <-input
		output <- lastY
		
		for x := range input {
			alpha := <-alphaModulation
			y := lastY + alpha * (x - lastY)
			output <- y
			
			lastY = y
		}
	}()
	
	return output
}
